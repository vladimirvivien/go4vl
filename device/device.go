package device

import (
	"context"
	"errors"
	"fmt"
	"os"
	sys "syscall"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// Device represents a video4linux device.
// It provides methods to interact with the device, such as configuring it,
// starting and stopping video streams, and accessing video frames.
type Device struct {
	// path is the file system path to the device (e.g., /dev/video0).
	path string
	// file is the opened file descriptor for the device.
	file *os.File
	// fd is the file descriptor handle for the opened device.
	fd uintptr
	// config holds the configuration settings for the device.
	config config
	// bufType specifies the type of buffer (e.g., video capture, video output).
	bufType v4l2.BufType
	// cap stores the capabilities of the device.
	cap v4l2.Capability
	// cropCap stores the cropping capabilities of the device.
	cropCap v4l2.CropCapability
	// buffers is a slice of byte slices representing memory-mapped buffers.
	buffers [][]byte
	// requestedBuf stores the buffer request parameters.
	requestedBuf v4l2.RequestBuffers
	// streaming indicates whether the device is currently streaming.
	streaming bool
	// output is a channel that delivers video frames from the device.
	output chan []byte
	// frameDataBuffers is a ring buffer for holding frame data, to reduce allocations.
	frameDataBuffers [][]byte
	// currentFrameDataBufferIndex is the next index to use in frameDataBuffers.
	currentFrameDataBufferIndex int
}

// Open opens the video device at the specified path with the given options.
// It initializes the device, queries its capabilities, and applies any provided configurations.
//
// Parameters:
//   path: The file system path to the video device (e.g., "/dev/video0").
//   options: A variadic list of Option functions to configure the device.
//
// Returns:
//   A pointer to a Device struct if successful, or an error if the device cannot be opened or configured.
func Open(path string, options ...Option) (*Device, error) {
	fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK, 0)
	if err != nil {
		return nil, fmt.Errorf("device open: %w", err)
	}

	dev := &Device{path: path, config: config{}, fd: fd}
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
			return nil, fmt.Errorf("device %s: closing after failure: %w", path, err)
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
		// setup capture parameters and chan for captured data
		dev.bufType = v4l2.BufTypeVideoCapture
		dev.output = make(chan []byte, dev.config.bufSize)
	case cap.IsVideoOutputSupported():
		dev.bufType = v4l2.BufTypeVideoOutput
	default:
		if err := v4l2.CloseDevice(dev.fd); err != nil {
			return nil, fmt.Errorf("device open: %s: closing after failure: %w", path, err)
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

// Close stops the video stream (if active) and closes the underlying device file descriptor.
//
// Returns:
//   An error if stopping the stream or closing the device fails.
func (d *Device) Close() error {
	if d.streaming {
		if err := d.Stop(); err != nil {
			return err
		}
	}
	return v4l2.CloseDevice(d.fd)
}

// Name returns the file system path of the device.
func (d *Device) Name() string {
	return d.path
}

// Fd returns the file descriptor of the opened device.
func (d *Device) Fd() uintptr {
	return d.fd
}

// Buffers returns a slice of byte slices representing the memory-mapped buffers
// used for streaming. This method should be called after streaming has been started;
// otherwise, it may return nil or an empty slice.
func (d *Device) Buffers() [][]byte {
	return d.buffers
}

// Capability returns the capabilities of the video device, such as whether
// it supports video capture, streaming, etc.
func (d *Device) Capability() v4l2.Capability {
	return d.cap
}

// BufferType returns the type of buffer used by the device (e.g., video capture, video output).
// This is a convenience method; for more detailed capability information, use Capability().
func (d *Device) BufferType() v4l2.BufType {
	return d.bufType
}

// BufferCount returns the number of buffers configured for streaming.
// This value might be updated by the driver after streaming starts.
// Note: The current implementation returns d.config.bufSize which is a uint32,
// but the function signature returns v4l2.BufType. This might be an issue.
func (d *Device) BufferCount() v4l2.BufType {
	return d.config.bufSize
}

// MemIOType returns the memory I/O type used by the device (e.g., memory mapping, user pointer).
func (d *Device) MemIOType() v4l2.IOType {
	return d.config.ioType
}

// GetOutput returns a read-only channel that delivers video frames (as byte slices)
// captured from the device. Frames are sent to this channel during active streaming.
//
// **Warning:** The `[]byte` received from this channel is part of an internal ring buffer
// and its content **will be overwritten** by subsequent frame captures.
// If you need to retain the data from a frame beyond the immediate processing scope
// (e.g., after reading the next frame or after the current select case block completes),
// you **must make a copy** of the byte slice. For example:
//
//   frameData := <-dev.GetOutput()
//   myCopy := make([]byte, len(frameData))
//   copy(myCopy, frameData)
//   // Use myCopy for long-term storage or processing
func (d *Device) GetOutput() <-chan []byte {
	return d.output
}

// SetInput sets up an input channel for data to be sent for output to the
// underlying device driver. This is typically used for video output devices.
// The current implementation is a placeholder.
func (d *Device) SetInput(in <-chan []byte) {

}

// GetCropCapability returns the cropping capabilities of the device.
// This includes information like the default cropping rectangle and bounds.
// Returns an error if the device does not support video capture.
func (d *Device) GetCropCapability() (v4l2.CropCapability, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.CropCapability{}, v4l2.ErrorUnsupportedFeature
	}
	return d.cropCap, nil
}

// SetCropRect sets the cropping rectangle for the video device.
// The parameter `r` specifies the desired cropping rectangle.
// Returns an error if the device does not support video capture or if setting the crop rectangle fails.
func (d *Device) SetCropRect(r v4l2.Rect) error {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.ErrorUnsupportedFeature
	}
	if err := v4l2.SetCropRect(d.fd, r); err != nil {
		return fmt.Errorf("device: %w", err)
	}
	return nil
}

// GetPixFormat retrieves the current pixel format of the device.
// This includes information like width, height, pixel format code, and field order.
// If the pixel format has not been explicitly set, it queries the device for the default format.
// Returns an error if the device does not support video capture or if querying the format fails.
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

// SetPixFormat sets the pixel format for the video device.
// The parameter `pixFmt` specifies the desired pixel format settings.
// Returns an error if the device does not support video capture or if setting the format fails.
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

// GetFormatDescription returns a description of a specific video format supported by the device.
// The parameter `idx` is the zero-based index of the format description to retrieve.
// Returns an error if the device does not support video capture or if querying the description fails.
func (d *Device) GetFormatDescription(idx uint32) (v4l2.FormatDescription, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.FormatDescription{}, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetFormatDescription(d.fd, idx)
}

// GetFormatDescriptions returns a slice of all video format descriptions supported by the device.
// Returns an error if the device does not support video capture.
func (d *Device) GetFormatDescriptions() ([]v4l2.FormatDescription, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return nil, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetAllFormatDescriptions(d.fd)
}

// GetVideoInputIndex returns the current video input index for the device.
// Returns an error if the device does not support video capture.
func (d *Device) GetVideoInputIndex() (int32, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return 0, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetCurrentVideoInputIndex(d.fd)
}

// GetVideoInputInfo returns information about a specific video input of the device.
// The parameter `index` is the zero-based index of the video input to query.
// Returns an error if the device does not support video capture or if querying the input info fails.
func (d *Device) GetVideoInputInfo(index uint32) (v4l2.InputInfo, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.InputInfo{}, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetVideoInputInfo(d.fd, index)
}

// GetStreamParam returns the streaming parameters for the device.
// This includes parameters like frame rate.
// Returns an error if the device does not support video capture or output, or if querying fails.
func (d *Device) GetStreamParam() (v4l2.StreamParam, error) {
	if !d.cap.IsVideoCaptureSupported() && d.cap.IsVideoOutputSupported() {
		return v4l2.StreamParam{}, v4l2.ErrorUnsupportedFeature
	}
	return v4l2.GetStreamParam(d.fd, d.bufType)
}

// SetStreamParam sets the streaming parameters for the device.
// The parameter `param` specifies the desired streaming parameters.
// Returns an error if the device does not support video capture or output, or if setting the parameters fails.
func (d *Device) SetStreamParam(param v4l2.StreamParam) error {
	if !d.cap.IsVideoCaptureSupported() && d.cap.IsVideoOutputSupported() {
		return v4l2.ErrorUnsupportedFeature
	}
	return v4l2.SetStreamParam(d.fd, d.bufType, param)
}

// SetFrameRate sets the frame rate (frames per second) for the device.
// The parameter `fps` is the desired frame rate.
// Returns an error if the device does not support streaming or if setting the frame rate fails.
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

// GetFrameRate retrieves the current frame rate (frames per second) of the device.
// If the frame rate has not been explicitly set, it queries the device for the current rate.
// Returns an error if querying the stream parameters fails or if the device does not support video capture/output.
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

// GetMediaInfo returns media device information if the device supports the Media API.
func (d *Device) GetMediaInfo() (v4l2.MediaDeviceInfo, error) {
	return v4l2.GetMediaDeviceInfo(d.fd)
}

// Start begins the video streaming process.
// It takes a context for cancellation.
// This function initializes and maps buffers, queues them, and starts a goroutine
// to continuously dequeue and process video frames.
// Frames are sent to the channel obtained via GetOutput().
// Returns an error if the device does not support streaming, if streaming is already active,
// or if any step in the stream initialization process fails.
func (d *Device) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if !d.cap.IsStreamingSupported() {
		return fmt.Errorf("device: start stream: %w", v4l2.ErrorUnsupportedFeature)
	}

	if d.streaming {
		return fmt.Errorf("device: stream already started")
	}

	// allocate device buffers
	bufReq, err := v4l2.InitBuffers(d)
	if err != nil {
		return fmt.Errorf("device: requested buffer type not be supported: %w", err)
	}

	d.config.bufSize = bufReq.Count
	d.requestedBuf = bufReq

	// for each allocated device buf, map into local space
	if d.buffers, err = v4l2.MapMemoryBuffers(d); err != nil {
		return fmt.Errorf("device: make mapped buffers: %w", err)
	}

	// Initialize frame data buffers (ring buffer)
	d.frameDataBuffers = make([][]byte, d.config.bufSize)
	d.currentFrameDataBufferIndex = 0

	if err := d.startStreamLoop(ctx); err != nil {
		return fmt.Errorf("device: start stream loop: %w", err)
	}

	d.streaming = true

	return nil
}

// Stop terminates the video streaming process.
// It unmaps the memory buffers and turns off the stream.
// Returns an error if unmapping buffers or turning off the stream fails.
func (d *Device) Stop() error {
	if !d.streaming {
		return nil
	}
	if err := v4l2.UnmapMemoryBuffers(d); err != nil {
		return fmt.Errorf("device: stop: %w", err)
	}
	if err := v4l2.StreamOff(d); err != nil {
		return fmt.Errorf("device: stop: %w", err)
	}
	d.streaming = false
	return nil
}

// startStreamLoop sets up the main video capture loop.
// This function is intended to be run as a goroutine.
// It initializes the output channel for frames, queues initial buffers with the driver,
// and starts the video stream. It then enters a loop waiting for frames from the
// device or context cancellation.
//
// Parameters:
//   ctx: A context used to signal cancellation of the stream.
//
// Returns:
//   An error if queuing initial buffers or starting the stream fails.
func (d *Device) startStreamLoop(ctx context.Context) error {
	d.output = make(chan []byte, d.config.bufSize)

	// Initial enqueue of buffers for capture
	for i := 0; i < int(d.config.bufSize); i++ {
		_, err := v4l2.QueueBuffer(d.fd, d.config.ioType, d.bufType, uint32(i))
		if err != nil {
			return fmt.Errorf("device: buffer queueing: %w", err)
		}
	}

	if err := v4l2.StreamOn(d); err != nil {
		return fmt.Errorf("device: stream on: %w", err)
	}

	go func() {
		defer close(d.output)

		fd := d.Fd()
		// var frame []byte // Removed, using d.frameDataBuffers instead
		ioMemType := d.MemIOType()
		bufType := d.BufferType()
		waitForRead := v4l2.WaitForRead(d)
		for {
			select {
			// handle stream capture (read from driver)
			case <-waitForRead:
				buff, err := v4l2.DequeueBuffer(fd, ioMemType, bufType)
				if err != nil {
					if errors.Is(err, sys.EAGAIN) {
						continue
					}
					panic(fmt.Sprintf("device: stream loop dequeue: %s", err))
				}

				// copy mapped buffer (copying avoids polluted data from subsequent dequeue ops)
				if buff.Flags&v4l2.BufFlagMapped != 0 && buff.Flags&v4l2.BufFlagError == 0 {
					// Use the ring buffer
					targetBuf := &d.frameDataBuffers[d.currentFrameDataBufferIndex]
					if *targetBuf == nil || cap(*targetBuf) < int(buff.BytesUsed) {
						*targetBuf = make([]byte, buff.BytesUsed)
					} else {
						*targetBuf = (*targetBuf)[:buff.BytesUsed]
					}

					if n := copy(*targetBuf, d.buffers[buff.Index][:buff.BytesUsed]); n == 0 {
						// This case (n==0 for non-empty source) is unlikely with valid buff.BytesUsed.
						// Sending an empty slice if copy truly yielded nothing.
						d.output <- []byte{}
					} else {
						d.output <- *targetBuf
					}
					d.currentFrameDataBufferIndex = (d.currentFrameDataBufferIndex + 1) % int(d.config.bufSize)
				} else {
					// Handle error or non-mapped buffer by sending an empty slice
					d.output <- []byte{}
				}

				if _, err := v4l2.QueueBuffer(fd, ioMemType, bufType, buff.Index); err != nil {
					panic(fmt.Sprintf("device: stream loop queue: %s: buff: %#v", err, buff))
				}
			case <-ctx.Done():
				d.Stop()
				return
			}
		}
	}()

	return nil
}
