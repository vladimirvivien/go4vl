package device

import (
	"context"
	"fmt"
	"sync/atomic"
	sys "syscall"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// Device represents a V4L2 video device and provides high-level methods for video capture and streaming.
// It encapsulates the low-level V4L2 API interactions and manages the device lifecycle including
// configuration, buffer management, and streaming operations.
//
// # Concurrency Safety
//
// The Device struct has limited concurrency support:
//
//   - Configuration methods (Open, Close, Start, Stop, SetPixFormat, SetFrameRate, etc.)
//     MUST be called from a single goroutine. Concurrent calls to these methods will
//     cause race conditions and undefined behavior.
//
//   - GetOutput() and GetError() return channels that are safe to read from concurrently.
//     Multiple goroutines can safely read from these channels.
//
//   - It is safe to call configuration methods while another goroutine reads from the
//     output/error channels, but calling Start()/Stop() concurrently from multiple
//     goroutines is NOT safe.
type Device struct {
	// path is the file system path to the device (e.g., /dev/video0, /dev/video1)
	path string

	// fd is the file descriptor used for low-level V4L2 ioctl operations
	fd uintptr

	// config holds the device configuration including pixel format, FPS, buffer size, and IO type
	config config

	// bufType specifies the buffer type (capture or output) determined by device capabilities
	bufType v4l2.BufType

	// cap stores the device capabilities queried during initialization
	cap v4l2.Capability

	// cropCap stores the cropping capabilities for video capture devices
	cropCap v4l2.CropCapability

	// buffers holds memory-mapped buffers used for zero-copy frame transfer
	buffers [][]byte

	// requestedBuf contains the buffer request parameters negotiated with the driver
	requestedBuf v4l2.RequestBuffers

	// streaming indicates whether the device is currently streaming video data
	// Use atomic operations to access this field for thread-safety
	streaming atomic.Bool

	// streamingMode tracks which API is being used (0=none, 1=GetOutput, 2=GetFrames)
	// This ensures mutual exclusivity between the two streaming approaches
	streamingMode atomic.Int32

	// frames is the channel that delivers Frame objects with metadata to consumers
	frames chan *Frame

	// output is the channel that delivers captured video frames to consumers (legacy API)
	output chan []byte

	// streamErr is the channel that delivers streaming errors to consumers
	streamErr chan error

	// captureDone is closed when the capture goroutine exits
	// Used by Stop() to wait for clean goroutine shutdown before unmapping buffers
	captureDone chan struct{}

	// framePool is the pool used for frame buffer allocation
	framePool *FramePool
}

// Open opens a V4L2 video device at the specified path and prepares it for streaming.
// The device is opened in read-write mode with non-blocking I/O.
//
// Parameters:
//   - path: The file system path to the video device (e.g., "/dev/video0")
//   - options: Optional configuration functions to customize device behavior
//
// Available options:
//   - WithBufferSize(n): Set the number of buffers for streaming (default: 2)
//   - WithPixFormat(fmt): Set the pixel format and frame dimensions
//   - WithFPS(fps): Set the frames per second for capture
//   - WithIOType(io): Set the I/O method (currently only memory-mapped I/O is supported)
//
// The function performs the following initialization steps:
//  1. Opens the device file descriptor
//  2. Queries device capabilities
//  3. Determines buffer type (capture/output) based on capabilities
//  4. Configures pixel format and frame rate
//  5. Resets crop settings to defaults if supported
//
// Returns:
//   - *Device: A configured device ready for streaming
//   - error: An error if the device cannot be opened or configured
//
// Possible errors:
//   - Device not found or permission denied
//   - Device does not support streaming I/O
//   - Unsupported buffer type or pixel format
//   - Failed to set requested frame rate
//
// Example:
//
//	dev, err := device.Open("/dev/video0",
//	    device.WithBufferSize(4),
//	    device.WithPixFormat(v4l2.PixFormat{
//	        Width: 640,
//	        Height: 480,
//	        PixelFormat: v4l2.PixelFmtMJPEG,
//	    }),
//	    device.WithFPS(30),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer dev.Close()
func Open(path string, options ...Option) (*Device, error) {
	fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK, 0)
	if err != nil {
		return nil, fmt.Errorf("device open: %w", err)
	}

	dev := &Device{
		path:      path,
		config:    config{},
		fd:        fd,
		framePool: defaultFramePool, // Use global default pool
	}
	// apply options
	if len(options) > 0 {
		for _, o := range options {
			o(&dev.config)
		}
	}

	// get capability
	cap, err := v4l2.GetCapability(dev.fd)
	if err != nil {
		if err := v4l2.CloseDevice(dev.fd); err != nil {
			return nil, fmt.Errorf("device %s: closing after failure: %s", path, err)
		}
		return nil, fmt.Errorf("device open: %s: %w", path, err)
	}
	dev.cap = cap

	// set preferred device buffer size
	if dev.config.bufSize == 0 {
		dev.config.bufSize = 2
	}

	// only supports streaming IO model right now
	if !dev.cap.IsStreamingSupported() {
		return nil, fmt.Errorf("device open: device does not support streamingIO")
	}

	switch {
	case cap.IsVideoCaptureSupported():
		// setup capture parameters (output channel created when streaming starts)
		dev.bufType = v4l2.BufTypeVideoCapture
	case cap.IsVideoOutputSupported():
		dev.bufType = v4l2.BufTypeVideoOutput
	default:
		if err := v4l2.CloseDevice(dev.fd); err != nil {
			return nil, fmt.Errorf("device open: %s: closing after failure: %s", path, err)
		}
		return nil, fmt.Errorf("device open: %s: %w", path, v4l2.ErrorUnsupportedFeature)
	}

	if dev.config.bufType != 0 && dev.config.bufType != dev.bufType {
		return nil, fmt.Errorf("device open: does not support buffer stream type")
	}

	// ensures IOType is set, only MemMap supported now
	dev.config.ioType = v4l2.IOTypeMMAP

	// reset crop, only if cropping supported
	if cropcap, err := v4l2.GetCropCapability(dev.fd, dev.bufType); err == nil {
		if err := v4l2.SetCropRect(dev.fd, cropcap.DefaultRect); err != nil {
			// ignore errors
		}
	}

	// set pix format
	if dev.config.pixFormat != (v4l2.PixFormat{}) {
		if err := dev.SetPixFormat(dev.config.pixFormat); err != nil {
			return nil, fmt.Errorf("device open: %s: set format: %w", path, err)
		}
	} else {
		dev.config.pixFormat, err = v4l2.GetPixFormat(dev.fd)
		if err != nil {
			return nil, fmt.Errorf("device open: %s: get default format: %w", path, err)
		}
	}

	// set fps
	if dev.config.fps != 0 {
		if err := dev.SetFrameRate(dev.config.fps); err != nil {
			return nil, fmt.Errorf("device open: %s: set fps: %w", path, err)
		}
	} else {
		if dev.config.fps, err = dev.GetFrameRate(); err != nil {
			return nil, fmt.Errorf("device open: %s: get fps: %w", path, err)
		}
	}

	return dev, nil
}

// Close closes the device and releases all associated resources.
// If the device is currently streaming, it will be stopped before closing.
// This method should always be called when done with the device, typically using defer.
//
// Returns an error if stopping the stream or closing the file descriptor fails.
func (d *Device) Close() error {
	if d.streaming.Load() {
		if err := d.Stop(); err != nil {
			return err
		}
	}
	return v4l2.CloseDevice(d.fd)
}

// Name returns the file system path of the device (e.g., "/dev/video0").
// This is the same path that was provided to the Open function.
func (d *Device) Name() string {
	return d.path
}

// Fd returns the underlying file descriptor for the device.
// This can be used for low-level V4L2 operations or with select/poll/epoll.
//
// Note: Direct manipulation of the file descriptor may interfere with
// the Device's internal state management.
func (d *Device) Fd() uintptr {
	return d.fd
}

// Buffers returns the memory-mapped buffers used for video streaming.
// These buffers are allocated and mapped during Start() and provide zero-copy
// access to video frames.
//
// Returns nil if called before Start() or after Stop().
// The returned slice should not be modified as it's used internally for streaming.
func (d *Device) Buffers() [][]byte {
	return d.buffers
}

// Capability returns the device capability information including supported
// features, driver info, and device-specific capabilities.
// This information is queried during device initialization.
//
// The capability can be used to check for specific features:
//   - Video capture/output support
//   - Streaming I/O support
//   - Multi-planar format support
//   - Hardware codec support
func (d *Device) Capability() v4l2.Capability {
	return d.cap
}

// BufferType returns the buffer type for this device, which indicates whether
// the device is configured for video capture (BufTypeVideoCapture) or
// video output (BufTypeVideoOutput).
//
// This is determined automatically during device initialization based on
// the device's capabilities.
func (d *Device) BufferType() v4l2.BufType {
	return d.bufType
}

// BufferCount returns the number of buffers allocated for streaming.
// This value may differ from the requested count as the driver may
// adjust it based on hardware constraints.
//
// The actual count is finalized when Start() is called.
// Before streaming, this returns the requested count.
// After streaming starts, it returns the actual allocated count.
func (d *Device) BufferCount() v4l2.BufType {
	return d.config.bufSize
}

// MemIOType returns the memory I/O method used for frame transfer.
// Currently, only memory-mapped I/O (IOTypeMMAP) is supported, which provides
// efficient zero-copy frame access.
//
// Future versions may support additional I/O methods like DMA or user pointers.
func (d *Device) MemIOType() v4l2.IOType {
	return d.config.ioType
}

// GetOutput returns a read-only channel that delivers captured video frames.
// Each frame is delivered as a byte slice containing the raw frame data in the
// configured pixel format.
//
// Deprecated: Use GetFrames() instead. GetFrames() provides the same functionality
// with significantly better performance through buffer pooling, using 540x less memory
// (1 KB vs 614 KB per frame) and reducing GC pressure. This method will be removed
// in a future version.
//
// This method selects the legacy streaming mode. Once called, GetFrames() cannot be used
// on this device instance. Call Stop() to reset the streaming mode.
//
// The channel is created when the device starts streaming and buffered according to
// the configured buffer size. Frames are delivered continuously while streaming.
//
// The channel is closed when Stop() is called or the context is cancelled.
// Consumers should handle channel closure gracefully.
//
// Frame format depends on the configured pixel format (e.g., MJPEG, YUYV, H264).
// Use GetPixFormat() to determine the current format.
//
// Example:
//
//	for frame := range dev.GetOutput() {
//	    // Process frame data
//	    fmt.Printf("Received frame: %d bytes\n", len(frame))
//	}
func (d *Device) GetOutput() <-chan []byte {
	// Set streaming mode to GetOutput (mode 1)
	// CompareAndSwap ensures first-caller-wins semantics
	if d.streamingMode.CompareAndSwap(0, 1) {
		// Successfully claimed GetOutput mode
	} else if d.streamingMode.Load() != 1 {
		// GetFrames() was called first - return nil channel
		// Reading from nil channel will block forever (user will notice the error)
		return nil
	}
	return d.output
}

// GetFrames returns a read-only channel that delivers captured video frames
// with metadata. Each Frame includes the raw frame data, timestamp, sequence number,
// and buffer flags.
//
// This method selects the optimized streaming mode with buffer pooling. Once called,
// GetOutput() cannot be used on this device instance. Call Stop() to reset the streaming mode.
//
// Users MUST call Frame.Release() when done processing to return the buffer to the pool.
//
// The channel is created when the device starts streaming and buffered according to
// the configured buffer size. Frames are delivered continuously while streaming.
//
// The channel is closed when Stop() is called or the context is cancelled.
// Consumers should handle channel closure gracefully.
//
// Example:
//
//	for frame := range dev.GetFrames() {
//	    // Process frame data and metadata
//	    fmt.Printf("Frame %d: %d bytes at %v\n",
//	        frame.Sequence, len(frame.Data), frame.Timestamp)
//
//	    // Check if it's a keyframe
//	    if frame.IsKeyFrame() {
//	        fmt.Println("Keyframe detected")
//	    }
//
//	    // IMPORTANT: Always release when done
//	    frame.Release()
//	}
func (d *Device) GetFrames() <-chan *Frame {
	// Set streaming mode to GetFrames (mode 2)
	// CompareAndSwap ensures first-caller-wins semantics
	if d.streamingMode.CompareAndSwap(0, 2) {
		// Successfully claimed GetFrames mode
	} else if d.streamingMode.Load() != 2 {
		// GetOutput() was called first - return nil channel
		// Reading from nil channel will block forever (user will notice the error)
		return nil
	}
	return d.frames
}

// GetError returns a read-only channel that delivers streaming errors.
// Errors are sent when critical issues occur during streaming, such as:
//   - Failed buffer dequeue operations
//   - Failed buffer requeue operations
//   - Driver errors
//
// The channel is created when Start() is called and closed when Stop() is called
// or the context is cancelled.
//
// Consumers should monitor this channel to detect streaming failures:
//
//	go func() {
//	    for err := range dev.GetError() {
//	        log.Printf("Streaming error: %v", err)
//	    }
//	}()
func (d *Device) GetError() <-chan error {
	return d.streamErr
}

// SetInput sets up an input channel for sending video data to output devices.
// This method is currently not implemented and reserved for future use with
// video output devices.
//
// TODO: Implement for video output device support
func (d *Device) SetInput(in <-chan []byte) {

}

// GetCropCapability returns the cropping capabilities for video capture devices,
// including supported crop boundaries and default crop rectangle.
//
// Returns ErrorUnsupportedFeature if the device doesn't support video capture
// or cropping operations.
func (d *Device) GetCropCapability() (v4l2.CropCapability, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.CropCapability{}, v4l2.ErrorUnsupportedFeature
	}
	return d.cropCap, nil
}

// SetCropRect sets the crop rectangle for video capture, allowing selection of
// a sub-region of the full frame.
//
// Parameters:
//   - r: Rectangle specifying the crop region (left, top, width, height)
//
// Returns ErrorUnsupportedFeature if the device doesn't support cropping.
// The actual crop region may be adjusted by the driver to match hardware constraints.
func (d *Device) SetCropRect(r v4l2.Rect) error {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.ErrorUnsupportedFeature
	}
	if err := v4l2.SetCropRect(d.fd, r); err != nil {
		return fmt.Errorf("device: %w", err)
	}
	return nil
}

// GetPixFormat returns the current pixel format configuration including
// frame dimensions, pixel format, field order, and other format parameters.
//
// The format is either explicitly set via SetPixFormat() or the device's
// default format queried during initialization.
//
// Returns ErrorUnsupportedFeature if the device doesn't support video capture.
func (d *Device) GetPixFormat() (v4l2.PixFormat, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.PixFormat{}, v4l2.ErrorUnsupportedFeature
	}

	if d.config.pixFormat == (v4l2.PixFormat{}) {
		pixFmt, err := v4l2.GetPixFormat(d.fd)
		if err != nil {
			return v4l2.PixFormat{}, fmt.Errorf("device: %w", err)
		}
		d.config.pixFormat = pixFmt
	}

	return d.config.pixFormat, nil
}

// SetPixFormat configures the pixel format for video capture including
// frame dimensions, pixel encoding, and other format parameters.
//
// Parameters:
//   - pixFmt: Pixel format specification including width, height, and format
//
// Common pixel formats:
//   - PixelFmtMJPEG: Motion JPEG compressed format
//   - PixelFmtYUYV: YUV 4:2:2 packed format
//   - PixelFmtH264: H.264 compressed video
//
// The driver may adjust the requested format to match hardware capabilities.
// Call GetPixFormat() after setting to retrieve the actual format.
//
// Returns ErrorUnsupportedFeature if the device doesn't support video capture.
func (d *Device) SetPixFormat(pixFmt v4l2.PixFormat) error {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.ErrorUnsupportedFeature
	}

	if err := v4l2.SetPixFormat(d.fd, pixFmt); err != nil {
		return fmt.Errorf("device: %w", err)
	}
	d.config.pixFormat = pixFmt
	return nil
}

// GetFormatDescription returns the format description at the specified index.
// Format descriptions enumerate the pixel formats supported by the device.
//
// Parameters:
//   - idx: Zero-based index of the format to query
//
// Use GetFormatDescriptions() to retrieve all supported formats at once.
// Returns ErrorUnsupportedFeature if the device doesn't support video capture.
func (d *Device) GetFormatDescription(idx uint32) (v4l2.FormatDescription, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.FormatDescription{}, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetFormatDescription(d.fd, idx)
}

// GetFormatDescriptions returns all pixel formats supported by the device.
// Each format description includes the format name, pixel format code,
// and format flags.
//
// This is useful for discovering device capabilities and validating format support
// before calling SetPixFormat().
//
// Returns ErrorUnsupportedFeature if the device doesn't support video capture.
//
// Example:
//
//	formats, err := dev.GetFormatDescriptions()
//	for _, fmt := range formats {
//	    log.Printf("Supported format: %s (0x%08x)\n", fmt.Description, fmt.PixelFormat)
//	}
func (d *Device) GetFormatDescriptions() ([]v4l2.FormatDescription, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return nil, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetAllFormatDescriptions(d.fd)
}

// GetVideoInputIndex returns the currently selected video input index.
// Video devices may have multiple inputs (e.g., composite, S-Video, HDMI).
//
// Returns ErrorUnsupportedFeature if the device doesn't support video capture
// or input selection.
func (d *Device) GetVideoInputIndex() (int32, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return 0, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetCurrentVideoInputIndex(d.fd)
}

// GetVideoInputInfo returns information about a specific video input.
//
// Parameters:
//   - index: Zero-based index of the input to query
//
// The returned InputInfo includes the input name, type, and supported standards.
// Returns ErrorUnsupportedFeature if the device doesn't support video capture.
func (d *Device) GetVideoInputInfo(index uint32) (v4l2.InputInfo, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.InputInfo{}, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetVideoInputInfo(d.fd, index)
}

// SetVideoInputIndex sets the currently selected video input index.
//
// Parameters:
//   - index: Zero-based index of the input to select
//
// Returns ErrorUnsupportedFeature if the device doesn't support video capture
// or input selection.
func (d *Device) SetVideoInputIndex(index int32) error {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.ErrorUnsupportedFeature
	}

	return v4l2.SetVideoInputIndex(d.fd, index)
}

// GetVideoInputDescriptions returns all video inputs supported by the device.
// Each input description includes the input name, type, and capabilities.
//
// This is useful for discovering input options before calling SetVideoInputIndex().
//
// Returns ErrorUnsupportedFeature if the device doesn't support video capture.
//
// Example:
//
//	inputs, err := dev.GetVideoInputDescriptions()
//	for _, in := range inputs {
//	    log.Printf("Input %d: %s (type=%d, status=%s)\n",
//	        in.GetIndex(), in.GetName(), in.GetInputType(),
//	        v4l2.InputStatuses[in.GetStatus()])
//	}
func (d *Device) GetVideoInputDescriptions() ([]v4l2.InputInfo, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return nil, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetAllVideoInputInfo(d.fd)
}

// GetVideoInputStatus returns the current status of the selected video input.
// The status includes signal detection, power status, and color information.
//
// Returns ErrorUnsupportedFeature if the device doesn't support video capture
// or input status queries.
func (d *Device) GetVideoInputStatus() (v4l2.InputStatus, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return 0, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.QueryInputStatus(d.fd)
}

// GetVideoOutputIndex returns the currently selected video output index.
// Video devices may have multiple outputs (e.g., HDMI, DisplayPort, composite).
//
// Returns ErrorUnsupportedFeature if the device doesn't support video output
// or output selection.
func (d *Device) GetVideoOutputIndex() (int32, error) {
	if !d.cap.IsVideoOutputSupported() {
		return 0, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetCurrentVideoOutputIndex(d.fd)
}

// SetVideoOutputIndex sets the currently selected video output index.
//
// Parameters:
//   - index: Zero-based index of the output to select
//
// Returns ErrorUnsupportedFeature if the device doesn't support video output
// or output selection.
func (d *Device) SetVideoOutputIndex(index int32) error {
	if !d.cap.IsVideoOutputSupported() {
		return v4l2.ErrorUnsupportedFeature
	}

	return v4l2.SetVideoOutputIndex(d.fd, index)
}

// GetVideoOutputInfo returns information about a specific video output.
//
// Parameters:
//   - index: Zero-based index of the output to query
//
// The returned OutputInfo includes the output name, type, and supported standards.
// Returns ErrorUnsupportedFeature if the device doesn't support video output.
func (d *Device) GetVideoOutputInfo(index uint32) (v4l2.OutputInfo, error) {
	if !d.cap.IsVideoOutputSupported() {
		return v4l2.OutputInfo{}, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetVideoOutputInfo(d.fd, index)
}

// GetVideoOutputDescriptions returns all video outputs supported by the device.
// Each output description includes the output name, type, and capabilities.
//
// This is useful for discovering output options before calling SetVideoOutputIndex().
//
// Returns ErrorUnsupportedFeature if the device doesn't support video output.
//
// Example:
//
//	outputs, err := dev.GetVideoOutputDescriptions()
//	for _, out := range outputs {
//	    log.Printf("Output %d: %s (type=%d)\n",
//	        out.GetIndex(), out.GetName(), out.GetOutputType())
//	}
func (d *Device) GetVideoOutputDescriptions() ([]v4l2.OutputInfo, error) {
	if !d.cap.IsVideoOutputSupported() {
		return nil, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetAllVideoOutputInfo(d.fd)
}

// GetVideoOutputStatus returns the current status of the selected video output.
// The status includes power and signal information.
//
// Returns ErrorUnsupportedFeature if the device doesn't support video output
// or output status queries.
func (d *Device) GetVideoOutputStatus() (v4l2.OutputStatus, error) {
	if !d.cap.IsVideoOutputSupported() {
		return 0, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.QueryOutputStatus(d.fd)
}

// GetAudioInfo returns information about a specific audio input.
//
// Parameters:
//   - index: Zero-based index of the audio input to query
//
// The returned AudioInfo includes the audio input name, capabilities, and mode.
// This can be used to enumerate available audio inputs associated with the device.
//
// Example:
//
//	audio, err := dev.GetAudioInfo(0)
//	if err == nil {
//	    log.Printf("Audio: %s, Stereo: %v\n", audio.GetName(), audio.IsStereo())
//	}
func (d *Device) GetAudioInfo(index uint32) (v4l2.AudioInfo, error) {
	return v4l2.GetAudioInfo(d.fd, index)
}

// GetAudioDescriptions returns all audio inputs supported by the device.
// Each audio description includes the audio name, capabilities, and modes.
//
// This is useful for discovering available audio inputs before selecting one.
//
// Example:
//
//	audios, err := dev.GetAudioDescriptions()
//	for _, audio := range audios {
//	    log.Printf("Audio %d: %s, Stereo: %v, AVL: %v\n",
//	        audio.GetIndex(), audio.GetName(),
//	        audio.IsStereo(), audio.HasAVL())
//	}
func (d *Device) GetAudioDescriptions() ([]v4l2.AudioInfo, error) {
	return v4l2.GetAllAudioInfo(d.fd)
}

// GetCurrentAudio returns the currently selected audio input.
//
// Returns the audio input configuration including index, name, and capabilities.
func (d *Device) GetCurrentAudio() (v4l2.AudioInfo, error) {
	return v4l2.GetCurrentAudio(d.fd)
}

// SetAudio sets the current audio input by index.
//
// Parameters:
//   - index: Zero-based index of the audio input to select
//
// Example:
//
//	// Select the first audio input
//	err := dev.SetAudio(0)
func (d *Device) SetAudio(index uint32) error {
	return v4l2.SetAudio(d.fd, index)
}

// SetAudioMode sets the audio mode for the current audio input.
//
// Parameters:
//   - mode: Audio mode to set (e.g., v4l2.AudioModeAVL for Automatic Volume Level)
//
// Example:
//
//	// Enable Automatic Volume Level
//	err := dev.SetAudioMode(v4l2.AudioModeAVL)
func (d *Device) SetAudioMode(mode v4l2.AudioMode) error {
	return v4l2.SetAudioMode(d.fd, mode)
}

// GetAudioOutInfo returns information about a specific audio output.
//
// Parameters:
//   - index: Zero-based index of the audio output to query
//
// The returned AudioOutInfo includes the audio output name, capabilities, and mode.
// This can be used to enumerate available audio outputs associated with the device.
//
// Example:
//
//	audioOut, err := dev.GetAudioOutInfo(0)
//	if err == nil {
//	    log.Printf("Audio Out: %s, Stereo: %v\n", audioOut.GetName(), audioOut.IsStereo())
//	}
func (d *Device) GetAudioOutInfo(index uint32) (v4l2.AudioOutInfo, error) {
	return v4l2.GetAudioOutInfo(d.fd, index)
}

// GetAudioOutDescriptions returns all audio outputs supported by the device.
// Each audio output description includes the name, capabilities, and modes.
//
// This is useful for discovering available audio outputs before selecting one.
//
// Example:
//
//	audioOuts, err := dev.GetAudioOutDescriptions()
//	for _, audioOut := range audioOuts {
//	    log.Printf("Audio Out %d: %s, Stereo: %v, AVL: %v\n",
//	        audioOut.GetIndex(), audioOut.GetName(),
//	        audioOut.IsStereo(), audioOut.HasAVL())
//	}
func (d *Device) GetAudioOutDescriptions() ([]v4l2.AudioOutInfo, error) {
	return v4l2.GetAllAudioOutInfo(d.fd)
}

// GetCurrentAudioOut returns the currently selected audio output.
//
// Returns the audio output configuration including index, name, and capabilities.
func (d *Device) GetCurrentAudioOut() (v4l2.AudioOutInfo, error) {
	return v4l2.GetCurrentAudioOut(d.fd)
}

// SetAudioOut sets the current audio output by index.
//
// Parameters:
//   - index: Zero-based index of the audio output to select
//
// Example:
//
//	// Select the first audio output
//	err := dev.SetAudioOut(0)
func (d *Device) SetAudioOut(index uint32) error {
	return v4l2.SetAudioOut(d.fd, index)
}

// SetAudioOutMode sets the audio mode for the current audio output.
//
// Parameters:
//   - mode: Audio mode to set (e.g., v4l2.AudioModeAVL for Automatic Volume Level)
//
// Example:
//
//	// Enable Automatic Volume Level on audio output
//	err := dev.SetAudioOutMode(v4l2.AudioModeAVL)
func (d *Device) SetAudioOutMode(mode v4l2.AudioMode) error {
	return v4l2.SetAudioOutMode(d.fd, mode)
}

// GetTunerInfo returns information about a specific tuner.
//
// Parameters:
//   - index: Zero-based index of the tuner to query
//
// The returned TunerInfo includes the tuner name, type, capabilities,
// frequency range, signal strength, and audio mode.
//
// Example:
//
//	tuner, err := dev.GetTunerInfo(0)
//	if err == nil {
//	    log.Printf("Tuner: %s, Type: %s, Signal: %d\n",
//	        tuner.GetName(),
//	        v4l2.TunerTypes[tuner.GetType()],
//	        tuner.GetSignal())
//	}
func (d *Device) GetTunerInfo(index uint32) (v4l2.TunerInfo, error) {
	return v4l2.GetTunerInfo(d.fd, index)
}

// GetAllTuners returns all tuners supported by the device.
// Each tuner description includes the name, type, capabilities, and frequency range.
//
// This is useful for discovering available tuners before tuning.
//
// Example:
//
//	tuners, err := dev.GetAllTuners()
//	for _, tuner := range tuners {
//	    log.Printf("Tuner %d: %s (%s), Range: %d-%d\n",
//	        tuner.GetIndex(), tuner.GetName(),
//	        v4l2.TunerTypes[tuner.GetType()],
//	        tuner.GetRangeLow(), tuner.GetRangeHigh())
//	}
func (d *Device) GetAllTuners() ([]v4l2.TunerInfo, error) {
	return v4l2.GetAllTuners(d.fd)
}

// SetTuner sets tuner parameters such as audio mode.
//
// Parameters:
//   - tuner: TunerInfo struct with desired settings
//
// Example:
//
//	tuner, _ := dev.GetTunerInfo(0)
//	// Modify tuner settings as needed
//	err := dev.SetTuner(tuner)
func (d *Device) SetTuner(tuner v4l2.TunerInfo) error {
	return v4l2.SetTuner(d.fd, tuner)
}

// GetFrequency returns the current frequency for the specified tuner.
//
// Parameters:
//   - tunerIndex: Zero-based index of the tuner
//
// Returns FrequencyInfo containing the frequency in device-specific units.
// Use TunerInfo.IsLowFreq() to determine if units are 1/16 kHz (true) or 1/16 MHz (false).
//
// Example:
//
//	freq, err := dev.GetFrequency(0)
//	if err == nil {
//	    log.Printf("Current frequency: %d\n", freq.GetFrequency())
//	}
func (d *Device) GetFrequency(tunerIndex uint32) (v4l2.FrequencyInfo, error) {
	return v4l2.GetFrequency(d.fd, tunerIndex)
}

// SetFrequency sets the tuner frequency.
//
// Parameters:
//   - tunerIndex: Zero-based index of the tuner
//   - tunerType: Type of tuner (e.g., v4l2.TunerTypeRadio, v4l2.TunerTypeAnalogTV)
//   - frequency: Frequency in device-specific units (check TunerInfo.IsLowFreq())
//
// For radio tuners with TunerCapLow capability:
//   - Units are 1/16000 kHz (62.5 Hz)
//   - Example: 100.5 MHz = 100500 kHz = 100500 * 16 = 1,608,000 units
//
// For tuners without TunerCapLow:
//   - Units are 1/16 MHz (62.5 kHz)
//
// Example:
//
//	// Set FM radio to 100.5 MHz (assuming TunerCapLow)
//	err := dev.SetFrequency(0, v4l2.TunerTypeRadio, 1608000)
func (d *Device) SetFrequency(tunerIndex uint32, tunerType v4l2.TunerType, frequency uint32) error {
	return v4l2.SetFrequency(d.fd, tunerIndex, tunerType, frequency)
}

// GetFrequencyBands returns all frequency bands for the specified tuner.
//
// Parameters:
//   - tunerIndex: Zero-based index of the tuner
//   - tunerType: Type of tuner (e.g., v4l2.TunerTypeRadio)
//
// Example:
//
//	bands, err := dev.GetFrequencyBands(0, v4l2.TunerTypeRadio)
//	for _, band := range bands {
//	    log.Printf("Band %d: %d-%d, Modulation: FM=%v AM=%v\n",
//	        band.GetIndex(),
//	        band.GetRangeLow(), band.GetRangeHigh(),
//	        band.GetModulation()&v4l2.BandModulationFM != 0,
//	        band.GetModulation()&v4l2.BandModulationAM != 0)
//	}
func (d *Device) GetFrequencyBands(tunerIndex uint32, tunerType v4l2.TunerType) ([]v4l2.FrequencyBandInfo, error) {
	return v4l2.GetAllFrequencyBands(d.fd, tunerIndex, tunerType)
}

// GetModulatorInfo returns information about a specific modulator.
//
// Parameters:
//   - index: Zero-based index of the modulator to query
//
// The returned ModulatorInfo includes the modulator name, type, capabilities,
// and frequency range.
//
// Example:
//
//	mod, err := dev.GetModulatorInfo(0)
//	if err == nil {
//	    log.Printf("Modulator: %s, Type: %s\n",
//	        mod.GetName(),
//	        v4l2.TunerTypes[mod.GetType()])
//	}
func (d *Device) GetModulatorInfo(index uint32) (v4l2.ModulatorInfo, error) {
	return v4l2.GetModulatorInfo(d.fd, index)
}

// GetAllModulators returns all modulators supported by the device.
// Each modulator description includes the name, type, capabilities, and frequency range.
//
// This is useful for discovering available modulators before transmission.
//
// Example:
//
//	modulators, err := dev.GetAllModulators()
//	for _, mod := range modulators {
//	    log.Printf("Modulator %d: %s (%s), Range: %d-%d\n",
//	        mod.GetIndex(), mod.GetName(),
//	        v4l2.TunerTypes[mod.GetType()],
//	        mod.GetRangeLow(), mod.GetRangeHigh())
//	}
func (d *Device) GetAllModulators() ([]v4l2.ModulatorInfo, error) {
	return v4l2.GetAllModulators(d.fd)
}

// SetModulator sets modulator parameters such as transmission subchannels.
//
// Parameters:
//   - modulator: ModulatorInfo struct with desired settings
//
// Example:
//
//	mod, _ := dev.GetModulatorInfo(0)
//	// Modify modulator settings as needed
//	err := dev.SetModulator(mod)
func (d *Device) SetModulator(modulator v4l2.ModulatorInfo) error {
	return v4l2.SetModulator(d.fd, modulator)
}

// GetStandard returns the currently selected video standard.
//
// Video standards define the analog video signal format (PAL, NTSC, SECAM, etc.)
// used by legacy analog video devices like TV tuners and composite video inputs.
//
// Returns the current standard ID, which may be a combination of multiple standards.
//
// Note: Modern digital devices (HDMI, etc.) use DV timings instead. This method
// will return an error for devices that don't support analog standards.
//
// Example:
//
//	stdId, err := dev.GetStandard()
//	if err == nil {
//	    log.Printf("Current standard: %s", v4l2.StdNames[stdId])
//	}
func (d *Device) GetStandard() (v4l2.StdId, error) {
	return v4l2.GetStandard(d.fd)
}

// SetStandard sets the video standard for analog video devices.
//
// Parameters:
//   - stdId: Standard identifier (e.g., v4l2.StdPAL, v4l2.StdNTSC, v4l2.StdSECAM)
//
// The standard ID may be a single standard or a set of standards (OR'd together).
// The driver will choose the best match if multiple standards are specified.
//
// Note: Changing the standard may also change the current video format.
//
// Example:
//
//	// Set to PAL-B/G (common in Western Europe)
//	err := dev.SetStandard(v4l2.StdPAL_BG)
//
//	// Set to NTSC-M (USA)
//	err := dev.SetStandard(v4l2.StdNTSC_M)
func (d *Device) SetStandard(stdId v4l2.StdId) error {
	return v4l2.SetStandard(d.fd, stdId)
}

// QueryStandard auto-detects the video standard from the current input signal.
//
// This method senses which of the supported standards is currently being received.
// Returns a set of all detected standards.
//
// Note: The device must support standard detection for this to work.
// Returns an error (typically ENOLINK) if no signal is detected.
//
// Example:
//
//	detected, err := dev.QueryStandard()
//	if err == nil {
//	    log.Printf("Detected standard: %s", v4l2.StdNames[detected])
//	    // Now set it
//	    dev.SetStandard(detected)
//	}
func (d *Device) QueryStandard() (v4l2.StdId, error) {
	return v4l2.QueryStandard(d.fd)
}

// EnumerateStandard retrieves information about a video standard by index.
//
// Parameters:
//   - index: Zero-based index of the standard to query
//
// Returns detailed information about the standard including name, frame rate,
// and line count.
//
// Example:
//
//	std, err := dev.EnumerateStandard(0)
//	if err == nil {
//	    log.Printf("Standard 0: %s", std)
//	}
func (d *Device) EnumerateStandard(index uint32) (v4l2.Standard, error) {
	return v4l2.EnumStandard(d.fd, index)
}

// GetAllStandards enumerates all supported video standards for the device.
//
// Returns a slice of all standards supported by this device.
// Returns an empty slice if the device doesn't support analog standards
// (e.g., digital-only devices like HDMI capture cards).
//
// Example:
//
//	standards, err := dev.GetAllStandards()
//	for _, std := range standards {
//	    log.Printf("Standard %d: %s (%.2f fps, %d lines)\n",
//	        std.Index(), std.Name(), std.FrameRate(), std.FrameLines())
//	}
func (d *Device) GetAllStandards() ([]v4l2.Standard, error) {
	return v4l2.GetAllStandards(d.fd)
}

// IsStandardSupported checks if a specific video standard is supported.
//
// Parameters:
//   - stdId: Standard identifier to check
//
// Returns true if the device supports the specified standard.
//
// Example:
//
//	if supported, _ := dev.IsStandardSupported(v4l2.StdPAL); supported {
//	    log.Println("PAL is supported")
//	}
func (d *Device) IsStandardSupported(stdId v4l2.StdId) (bool, error) {
	return v4l2.IsStandardSupported(d.fd, stdId)
}

// GetDVTimings returns the current Digital Video (DV) timings.
//
// DV timings are used for digital video interfaces like HDMI, DisplayPort, DVI, and SDI.
// They describe the video format including resolution, refresh rate, and synchronization.
//
// Returns the current DV timings configuration.
//
// Example:
//
//	timings, err := dev.GetDVTimings()
//	if err == nil {
//	    bt := timings.GetBTTimings()
//	    log.Printf("Resolution: %dx%d @ %.2f Hz",
//	        bt.GetWidth(), bt.GetHeight(), bt.GetFrameRate())
//	}
func (d *Device) GetDVTimings() (v4l2.DVTimings, error) {
	return v4l2.GetDVTimings(d.fd)
}

// SetDVTimings sets the Digital Video (DV) timings.
//
// Parameters:
//   - timings: DV timings to set
//
// Example:
//
//	// Set timings (typically from enumeration or query)
//	err := dev.SetDVTimings(timings)
func (d *Device) SetDVTimings(timings v4l2.DVTimings) error {
	return v4l2.SetDVTimings(d.fd, timings)
}

// QueryDVTimings attempts to auto-detect DV timings from the input signal.
//
// This is useful for HDMI capture cards that can automatically detect the
// incoming signal format (resolution, refresh rate, etc.).
//
// Returns the detected DV timings, or an error if no valid signal is detected.
//
// Example:
//
//	timings, err := dev.QueryDVTimings()
//	if err != nil {
//	    log.Printf("No signal detected: %v", err)
//	    return
//	}
//	bt := timings.GetBTTimings()
//	log.Printf("Detected: %dx%d @ %.2f Hz",
//	    bt.GetWidth(), bt.GetHeight(), bt.GetFrameRate())
func (d *Device) QueryDVTimings() (v4l2.DVTimings, error) {
	return v4l2.QueryDVTimings(d.fd)
}

// EnumerateDVTimings enumerates a specific DV timing by index.
//
// Parameters:
//   - index: Zero-based index of the DV timing to query
//   - pad: Pad number (use 0 for video nodes, specific pad for subdev nodes)
//
// Returns the enumerated DV timing at the specified index.
//
// Example:
//
//	enumTiming, err := dev.EnumerateDVTimings(0, 0)
//	if err == nil {
//	    timings := enumTiming.GetTimings()
//	    bt := timings.GetBTTimings()
//	    log.Printf("Timing %d: %dx%d @ %.2f Hz",
//	        enumTiming.GetIndex(),
//	        bt.GetWidth(), bt.GetHeight(), bt.GetFrameRate())
//	}
func (d *Device) EnumerateDVTimings(index uint32, pad uint32) (v4l2.EnumDVTimings, error) {
	return v4l2.EnumerateDVTimings(d.fd, index, pad)
}

// GetAllDVTimings enumerates all supported DV timings.
//
// Parameters:
//   - pad: Pad number (use 0 for video nodes)
//
// Returns a slice of all supported DV timings.
//
// Example:
//
//	timings, err := dev.GetAllDVTimings(0)
//	for i, timing := range timings {
//	    bt := timing.GetTimings().GetBTTimings()
//	    log.Printf("Timing %d: %dx%d @ %.2f Hz",
//	        i, bt.GetWidth(), bt.GetHeight(), bt.GetFrameRate())
//	}
func (d *Device) GetAllDVTimings(pad uint32) ([]v4l2.EnumDVTimings, error) {
	return v4l2.GetAllDVTimings(d.fd, pad)
}

// GetDVTimingsCap returns the DV timing capabilities.
//
// This describes the range of supported resolutions, pixel clocks,
// and timing standards supported by the device.
//
// Parameters:
//   - pad: Pad number (use 0 for video nodes)
//
// Example:
//
//	cap, err := dev.GetDVTimingsCap(0)
//	if err == nil {
//	    btCap := cap.GetBTTimingsCap()
//	    log.Printf("Supported resolutions: %dx%d to %dx%d",
//	        btCap.GetMinWidth(), btCap.GetMinHeight(),
//	        btCap.GetMaxWidth(), btCap.GetMaxHeight())
//	    log.Printf("Interlaced: %v, Progressive: %v",
//	        btCap.SupportsInterlaced(), btCap.SupportsProgressive())
//	}
func (d *Device) GetDVTimingsCap(pad uint32) (v4l2.DVTimingsCap, error) {
	return v4l2.GetDVTimingsCap(d.fd, pad)
}

// GetStreamParam returns the current streaming parameters including
// capture/output timing, buffer settings, and capability flags.
//
// For capture devices, this includes the frame interval (FPS) and capture modes.
// Returns ErrorUnsupportedFeature if the device doesn't support streaming parameters.
func (d *Device) GetStreamParam() (v4l2.StreamParam, error) {
	if !d.cap.IsVideoCaptureSupported() && d.cap.IsVideoOutputSupported() {
		return v4l2.StreamParam{}, v4l2.ErrorUnsupportedFeature
	}
	return v4l2.GetStreamParam(d.fd, d.bufType)
}

// SetStreamParam configures streaming parameters for the device.
//
// Parameters:
//   - param: Stream parameters including timing and buffer settings
//
// This is typically used to set frame timing parameters. For simpler FPS control,
// use SetFrameRate() instead.
//
// Returns ErrorUnsupportedFeature if the device doesn't support streaming parameters.
func (d *Device) SetStreamParam(param v4l2.StreamParam) error {
	if !d.cap.IsVideoCaptureSupported() && d.cap.IsVideoOutputSupported() {
		return v4l2.ErrorUnsupportedFeature
	}
	return v4l2.SetStreamParam(d.fd, d.bufType, param)
}

// SetFrameRate sets the frames per second (FPS) for video capture or output.
//
// Parameters:
//   - fps: Desired frames per second (e.g., 30, 60)
//
// The actual frame rate may be adjusted by the driver based on hardware capabilities.
// Use GetFrameRate() to retrieve the actual rate after setting.
//
// Common frame rates:
//   - 15 FPS: Low bandwidth streaming
//   - 30 FPS: Standard video capture
//   - 60 FPS: Smooth motion capture
//
// Returns ErrorUnsupportedFeature if the device doesn't support frame rate control.
func (d *Device) SetFrameRate(fps uint32) error {
	if !d.cap.IsStreamingSupported() {
		return fmt.Errorf("set frame rate: %w", v4l2.ErrorUnsupportedFeature)
	}

	var param v4l2.StreamParam
	switch {
	case d.cap.IsVideoCaptureSupported():
		param.Capture = v4l2.CaptureParam{TimePerFrame: v4l2.Fract{Numerator: 1, Denominator: fps}}
	case d.cap.IsVideoOutputSupported():
		param.Output = v4l2.OutputParam{TimePerFrame: v4l2.Fract{Numerator: 1, Denominator: fps}}
	default:
		return v4l2.ErrorUnsupportedFeature
	}
	if err := d.SetStreamParam(param); err != nil {
		return fmt.Errorf("device: set fps: %w", err)
	}
	d.config.fps = fps
	return nil
}

// GetFrameRate returns the current frames per second (FPS) setting.
// If not explicitly set, returns the device's default frame rate.
//
// The actual frame rate during streaming may vary based on lighting conditions,
// processing load, and hardware limitations.
func (d *Device) GetFrameRate() (uint32, error) {
	if d.config.fps == 0 {
		param, err := d.GetStreamParam()
		if err != nil {
			return 0, fmt.Errorf("device: frame rate: %w", err)
		}
		switch {
		case d.cap.IsVideoCaptureSupported():
			d.config.fps = param.Capture.TimePerFrame.Denominator
		case d.cap.IsVideoOutputSupported():
			d.config.fps = param.Output.TimePerFrame.Denominator
		default:
			return 0, v4l2.ErrorUnsupportedFeature
		}
	}

	return d.config.fps, nil
}

// GetMediaInfo returns media controller information for devices that support
// the Media Controller API, used for complex video pipelines.
//
// Returns an error if the device doesn't support the Media Controller API.
// Most simple webcams don't support this feature.
func (d *Device) GetMediaInfo() (v4l2.MediaDeviceInfo, error) {
	return v4l2.GetMediaDeviceInfo(d.fd)
}

// GetExtControls retrieves multiple extended control values atomically.
//
// This method allows getting multiple control values in a single atomic operation,
// which is particularly useful for codec controls and compound controls.
//
// Example:
//   ctrls := v4l2.NewExtControls()
//   ctrls.Add(v4l2.NewExtControl(v4l2.CtrlBrightness))
//   ctrls.Add(v4l2.NewExtControl(v4l2.CtrlContrast))
//   if err := device.GetExtControls(ctrls); err != nil {
//       return err
//   }
//   brightness := ctrls.GetControls()[0].GetValue()
//   contrast := ctrls.GetControls()[1].GetValue()
func (d *Device) GetExtControls(ctrls *v4l2.ExtControls) error {
	return v4l2.GetExtControls(d.fd, ctrls)
}

// SetExtControls sets multiple extended control values atomically.
//
// This method allows setting multiple control values in a single atomic operation.
// If any control fails, none of the controls are changed.
//
// Example:
//   ctrls := v4l2.NewExtControls()
//   ctrls.Add(v4l2.NewExtControlWithValue(v4l2.CtrlBrightness, 128))
//   ctrls.Add(v4l2.NewExtControlWithValue(v4l2.CtrlContrast, 100))
//   if err := device.SetExtControls(ctrls); err != nil {
//       return err
//   }
func (d *Device) SetExtControls(ctrls *v4l2.ExtControls) error {
	return v4l2.SetExtControls(d.fd, ctrls)
}

// TryExtControls tests whether extended control values would be accepted without actually setting them.
//
// This is useful for validating control values before applying them.
//
// Example:
//   ctrls := v4l2.NewExtControls()
//   ctrls.Add(v4l2.NewExtControlWithValue(v4l2.CtrlBrightness, 128))
//   if err := device.TryExtControls(ctrls); err != nil {
//       // Value would be rejected
//       return err
//   }
//   // Value is valid, safe to apply
func (d *Device) TryExtControls(ctrls *v4l2.ExtControls) error {
	return v4l2.TryExtControls(d.fd, ctrls)
}

// High-level convenience methods for common controls

// SetBrightness sets the brightness control value.
//
// Example:
//   err := device.SetBrightness(128)
func (d *Device) SetBrightness(value int32) error {
	ctrls := v4l2.NewExtControls()
	ctrls.AddValue(v4l2.CtrlBrightness, value)
	return v4l2.SetExtControls(d.fd, ctrls)
}

// GetBrightness gets the current brightness value.
//
// Example:
//   brightness, err := device.GetBrightness()
func (d *Device) GetBrightness() (int32, error) {
	ctrls := v4l2.NewExtControls()
	ctrls.Add(v4l2.NewExtControl(v4l2.CtrlBrightness))
	if err := v4l2.GetExtControls(d.fd, ctrls); err != nil {
		return 0, err
	}
	return ctrls.GetControls()[0].GetValue(), nil
}

// SetContrast sets the contrast control value.
//
// Example:
//   err := device.SetContrast(100)
func (d *Device) SetContrast(value int32) error {
	ctrls := v4l2.NewExtControls()
	ctrls.AddValue(v4l2.CtrlContrast, value)
	return v4l2.SetExtControls(d.fd, ctrls)
}

// GetContrast gets the current contrast value.
//
// Example:
//   contrast, err := device.GetContrast()
func (d *Device) GetContrast() (int32, error) {
	ctrls := v4l2.NewExtControls()
	ctrls.Add(v4l2.NewExtControl(v4l2.CtrlContrast))
	if err := v4l2.GetExtControls(d.fd, ctrls); err != nil {
		return 0, err
	}
	return ctrls.GetControls()[0].GetValue(), nil
}

// SetSaturation sets the saturation control value.
//
// Example:
//   err := device.SetSaturation(64)
func (d *Device) SetSaturation(value int32) error {
	ctrls := v4l2.NewExtControls()
	ctrls.AddValue(v4l2.CtrlSaturation, value)
	return v4l2.SetExtControls(d.fd, ctrls)
}

// GetSaturation gets the current saturation value.
//
// Example:
//   saturation, err := device.GetSaturation()
func (d *Device) GetSaturation() (int32, error) {
	ctrls := v4l2.NewExtControls()
	ctrls.Add(v4l2.NewExtControl(v4l2.CtrlSaturation))
	if err := v4l2.GetExtControls(d.fd, ctrls); err != nil {
		return 0, err
	}
	return ctrls.GetControls()[0].GetValue(), nil
}

// SetHue sets the hue control value.
//
// Example:
//   err := device.SetHue(0)
func (d *Device) SetHue(value int32) error {
	ctrls := v4l2.NewExtControls()
	ctrls.AddValue(v4l2.CtrlHue, value)
	return v4l2.SetExtControls(d.fd, ctrls)
}

// GetHue gets the current hue value.
//
// Example:
//   hue, err := device.GetHue()
func (d *Device) GetHue() (int32, error) {
	ctrls := v4l2.NewExtControls()
	ctrls.Add(v4l2.NewExtControl(v4l2.CtrlHue))
	if err := v4l2.GetExtControls(d.fd, ctrls); err != nil {
		return 0, err
	}
	return ctrls.GetControls()[0].GetValue(), nil
}

// SubscribeEvent subscribes to V4L2 events.
//
// Events allow applications to be notified of device state changes such as
// control value changes, end of stream, source resolution changes, etc.
//
// After subscribing, use DequeueEvent() to retrieve events from the device.
//
// Example - Subscribe to control change events:
//   sub := v4l2.NewControlEventSubscription(v4l2.CtrlBrightness)
//   sub.SetFlags(v4l2.EventSubFlagSendInitial)
//   if err := device.SubscribeEvent(sub); err != nil {
//       return err
//   }
//
// Example - Subscribe to all events:
//   sub := v4l2.NewEventSubscription(v4l2.EventAll)
//   if err := device.SubscribeEvent(sub); err != nil {
//       return err
//   }
func (d *Device) SubscribeEvent(sub *v4l2.EventSubscription) error {
	return v4l2.SubscribeEvent(d.fd, sub)
}

// UnsubscribeEvent unsubscribes from V4L2 events.
//
// The subscription parameter should match the original subscription.
//
// Example:
//   sub := v4l2.NewEventSubscription(v4l2.EventCtrl)
//   if err := device.UnsubscribeEvent(sub); err != nil {
//       return err
//   }
func (d *Device) UnsubscribeEvent(sub *v4l2.EventSubscription) error {
	return v4l2.UnsubscribeEvent(d.fd, sub)
}

// DequeueEvent retrieves a pending event from the device.
//
// This method blocks until an event is available or an error occurs.
// To use non-blocking event retrieval, use select() or poll() on the device file descriptor.
//
// Returns the event and nil on success, or nil and an error if no event is available.
//
// Example:
//   event, err := device.DequeueEvent()
//   if err != nil {
//       return err
//   }
//   switch event.GetType() {
//   case v4l2.EventCtrl:
//       ctrlData := event.GetCtrlData()
//       fmt.Printf("Control changed: ID=%d, Value=%d\n", event.GetID(), ctrlData.Value)
//   case v4l2.EventEOS:
//       fmt.Println("End of stream")
//   }
func (d *Device) DequeueEvent() (*v4l2.Event, error) {
	return v4l2.DequeueEvent(d.fd)
}

// Start begins video streaming from the device. This method allocates buffers,
// memory-maps them, and starts a background goroutine to handle frame capture.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// The streaming process:
//  1. Allocates the requested number of buffers
//  2. Memory-maps buffers for zero-copy access
//  3. Enqueues all buffers to the driver
//  4. Starts the streaming I/O
//  5. Launches a goroutine to handle frame delivery
//
// Frames are delivered through the channel returned by GetOutput().
// The stream continues until Stop() is called or the context is cancelled.
//
// Returns an error if:
//   - The context is already cancelled
//   - Streaming is already active
//   - Buffer allocation fails
//   - The device doesn't support streaming
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	if err := dev.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer dev.Stop()
//
//	for frame := range dev.GetOutput() {
//	    // Process frames until context timeout
//	}
func (d *Device) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if !d.cap.IsStreamingSupported() {
		return fmt.Errorf("device: start stream: %s", v4l2.ErrorUnsupportedFeature)
	}

	if d.streaming.Load() {
		return fmt.Errorf("device: stream already started")
	}

	d.streaming.Store(true)

	// allocate device buffers
	bufReq, err := v4l2.InitBuffers(d)
	if err != nil {
		return fmt.Errorf("device: requested buffer type not be supported: %w", err)
	}

	d.config.bufSize = bufReq.Count
	d.requestedBuf = bufReq

	// for each allocated device buf, map into local space
	if d.buffers, err = v4l2.MapMemoryBuffers(d); err != nil {
		return fmt.Errorf("device: make mapped buffers: %s", err)
	}

	// Create capture done channel for synchronization with Stop()
	d.captureDone = make(chan struct{})

	// Launch appropriate streaming loop based on which API is in use
	mode := d.streamingMode.Load()
	switch mode {
	case 1: // GetOutput() mode
		if err := d.captureRawBytes(ctx); err != nil {
			d.streaming.Store(false)
			return fmt.Errorf("device: start capture raw bytes: %s", err)
		}
	case 2: // GetFrames() mode
		if err := d.captureFrames(ctx); err != nil {
			d.streaming.Store(false)
			return fmt.Errorf("device: start capture frames: %s", err)
		}
	default:
		d.streaming.Store(false)
		return fmt.Errorf("device: no streaming API selected (call GetOutput() or GetFrames() before Start())")
	}

	return nil
}

// Stop halts video streaming and releases streaming resources.
// This method waits for the capture goroutine to exit, then unmaps buffers and stops the stream.
//
// The method ensures proper synchronization:
//  1. Waits for the capture goroutine to exit (with 500ms timeout)
//  2. Unmaps memory buffers (safe after goroutine exits)
//  3. Stops the device stream
//  4. Resets state for next Start()
//
// Safe to call multiple times - returns immediately if not streaming.
// Should always be called when done streaming to free resources.
//
// Returns an error if buffer unmapping or stream stopping fails.
func (d *Device) Stop() error {
	if !d.streaming.Load() {
		return nil
	}

	// Signal goroutines to stop BEFORE waiting
	// This ensures they see streaming=false and can exit cleanly
	d.streaming.Store(false)

	// Wait for capture goroutine to exit before unmapping buffers
	// This prevents segfaults from accessing unmapped memory
	if d.captureDone != nil {
		select {
		case <-d.captureDone:
			// Goroutine exited cleanly
		case <-time.After(500 * time.Millisecond):
			// Timeout - proceed anyway but this indicates a potential issue
			// In practice, context cancellation should cause quick exit
		}
	}

	if err := v4l2.UnmapMemoryBuffers(d); err != nil {
		return fmt.Errorf("device: stop: %w", err)
	}
	if err := v4l2.StreamOff(d); err != nil {
		return fmt.Errorf("device: stop: %w", err)
	}
	d.streamingMode.Store(0) // Reset mode to allow different API on next Start()
	// Set channels to nil so they can be recreated on next Start()
	// The goroutine will close them before exiting
	d.output = nil
	d.frames = nil
	d.streamErr = nil
	d.captureDone = nil
	return nil
}
