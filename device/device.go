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
