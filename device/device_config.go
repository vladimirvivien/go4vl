package device

import (
	"github.com/vladimirvivien/go4vl/v4l2"
)

// config holds device configuration parameters.
// This type is unexported and managed by functional options.
type config struct {
	ioType    v4l2.IOType
	pixFormat v4l2.PixFormat
	bufSize   uint32
	fps       uint32
	bufType   uint32
}

// Option is a functional option type for configuring a Device.
// It's a function that takes a pointer to a config struct and modifies it.
type Option func(*config)

// WithIOType creates an Option to set the I/O type for the device.
// Example: WithIOType(v4l2.IOTypeMMAP)
func WithIOType(ioType v4l2.IOType) Option {
	return func(o *config) {
		o.ioType = ioType
	}
}

// WithPixFormat creates an Option to set the pixel format for the device.
// This includes parameters like width, height, and pixel format code.
// Example: WithPixFormat(v4l2.PixFormat{Width: 640, Height: 480, PixelFormat: v4l2.PixelFmtMJPEG})
func WithPixFormat(pixFmt v4l2.PixFormat) Option {
	return func(o *config) {
		o.pixFormat = pixFmt
	}
}

// WithBufferSize creates an Option to set the number of buffers to be used for streaming.
// Example: WithBufferSize(4)
func WithBufferSize(size uint32) Option {
	return func(o *config) {
		o.bufSize = size
	}
}

// WithFPS creates an Option to set the desired frames per second (FPS) for the device.
// Example: WithFPS(30)
func WithFPS(fps uint32) Option {
	return func(o *config) {
		o.fps = fps
	}
}

// WithVideoCaptureEnabled creates an Option to configure the device for video capture.
// This sets the buffer type to v4l2.BufTypeVideoCapture.
func WithVideoCaptureEnabled() Option {
	return func(o *config) {
		o.bufType = v4l2.BufTypeVideoCapture
	}
}

// WithVideoOutputEnabled creates an Option to configure the device for video output.
// This sets the buffer type to v4l2.BufTypeVideoOutput.
func WithVideoOutputEnabled() Option {
	return func(o *config) {
		o.bufType = v4l2.BufTypeVideoOutput
	}
}
