package device

import (
	"github.com/vladimirvivien/go4vl/v4l2"
)

// IOMethod defines the I/O method used by the device for frame transfer.
type IOMethod int

const (
	// IOMethodStreaming uses continuous streaming with memory-mapped buffers.
	// This is the default. Use Start/GetFrames/GetOutput/Stop to capture frames.
	// Requires V4L2_CAP_STREAMING.
	IOMethodStreaming IOMethod = iota

	// IOMethodReadWrite uses direct read()/write() syscalls for frame transfer.
	// Use Read/ReadFrame methods directly after Open. No Start/Stop needed.
	// Simpler but less efficient than streaming. Requires V4L2_CAP_READWRITE.
	IOMethodReadWrite
)

// config holds the internal configuration for a Device.
// These settings are applied during device initialization and streaming setup.
type config struct {
	// ioMethod specifies the I/O method category (streaming or read/write)
	ioMethod IOMethod

	// ioType specifies the memory type for streaming I/O (MMAP, user pointer, etc.)
	// Only applies when ioMethod is IOMethodStreaming.
	ioType v4l2.IOType

	// pixFormat defines the pixel format, frame dimensions, and color space
	pixFormat v4l2.PixFormat

	// bufSize is the number of buffers to allocate for streaming
	bufSize uint32

	// fps is the desired frames per second
	fps uint32

	// bufType specifies the buffer type (capture or output)
	bufType uint32
}

// Option is a functional option for configuring a Device during initialization.
// Options are applied in the order they are provided to the Open function.
type Option func(*config)

// WithIOType sets the buffer memory type for streaming I/O.
// Only applies when using IOMethodStreaming (the default).
//
// Available I/O types:
//   - IOTypeMMAP: Memory-mapped I/O (default, zero-copy from kernel)
//   - IOTypeUserPtr: User pointer I/O (application-allocated buffers)
//   - IOTypeDMABuf: DMA buffer sharing (not yet supported)
//
// Example:
//
//	device.Open("/dev/video0", device.WithIOType(v4l2.IOTypeMMAP))
//
// Note: This option only applies when using streaming I/O mode (the default).
// It is ignored in read/write I/O mode. Automatically set to IOTypeMMAP if not specified.
func WithIOType(ioType v4l2.IOType) Option {
	return func(o *config) {
		o.ioType = ioType
	}
}

// WithIOMethod sets the I/O method for the device.
// IOMethodStreaming (default) uses Start/GetFrames/Stop for continuous capture.
// IOMethodReadWrite uses Read/ReadFrame for direct frame access.
//
// Example:
//
//	// Read/write mode — simple synchronous reads
//	device.Open("/dev/video0", device.WithIOMethod(device.IOMethodReadWrite))
//
//	// Streaming mode (default) — continuous capture via channels
//	device.Open("/dev/video0") // or device.WithIOMethod(device.IOMethodStreaming)
func WithIOMethod(method IOMethod) Option {
	return func(o *config) {
		o.ioMethod = method
	}
}

// WithPixFormat configures the pixel format for video capture or output.
// This includes frame dimensions, pixel encoding, field order, and color space.
//
// The pixel format determines:
//   - Width and Height: Frame dimensions in pixels
//   - PixelFormat: Encoding format (MJPEG, YUYV, H264, etc.)
//   - Field: Interlacing mode (progressive, interlaced, etc.)
//   - ColorSpace: Color space information
//   - BytesPerLine: Bytes per scan line (may be auto-calculated)
//   - SizeImage: Total image size in bytes (may be auto-calculated)
//
// Common configurations:
//
//	// 1080p MJPEG capture
//	device.WithPixFormat(v4l2.PixFormat{
//	    Width:       1920,
//	    Height:      1080,
//	    PixelFormat: v4l2.PixelFmtMJPEG,
//	    Field:       v4l2.FieldNone,
//	})
//
//	// 720p YUV capture
//	device.WithPixFormat(v4l2.PixFormat{
//	    Width:       1280,
//	    Height:      720,
//	    PixelFormat: v4l2.PixelFmtYUYV,
//	})
//
// Note: The driver may adjust the requested format to match hardware capabilities.
// Always check the actual format with GetPixFormat() after opening the device.
func WithPixFormat(pixFmt v4l2.PixFormat) Option {
	return func(o *config) {
		o.pixFormat = pixFmt
	}
}

// WithBufferSize sets the number of buffers to use for video streaming.
// More buffers can help prevent frame drops but increase memory usage and latency.
//
// Typical values:
//   - 1: Minimal latency, higher risk of drops
//   - 2: Default, balanced performance
//   - 4: Better reliability, slightly higher latency
//   - 8+: High reliability for slow consumers
//
// The actual number allocated may differ as the driver can adjust this value.
// Query the actual count with BufferCount() after starting the stream.
//
// Example:
//
//	// Use 4 buffers for reliable streaming
//	device.Open("/dev/video0", device.WithBufferSize(4))
//
// Considerations:
//   - Each buffer holds one complete frame
//   - Memory usage = buffer_count × frame_size
//   - More buffers help absorb processing delays
//   - Fewer buffers reduce latency
func WithBufferSize(size uint32) Option {
	return func(o *config) {
		o.bufSize = size
	}
}

// WithFPS sets the desired frames per second for video capture or output.
// The actual frame rate may be adjusted by the driver based on hardware capabilities.
//
// Common frame rates:
//   - 15 FPS: Low bandwidth, suitable for monitoring
//   - 24 FPS: Cinematic look
//   - 30 FPS: Standard video capture (default for many devices)
//   - 60 FPS: Smooth motion capture
//   - 120+ FPS: High-speed capture (if supported)
//
// Example:
//
//	// Configure 30 FPS capture
//	device.Open("/dev/video0", device.WithFPS(30))
//
// Note:
//   - The achievable FPS depends on resolution and pixel format
//   - Higher resolutions may limit maximum FPS
//   - Compressed formats (MJPEG, H264) typically allow higher FPS
//   - Query actual FPS with GetFrameRate() after configuration
func WithFPS(fps uint32) Option {
	return func(o *config) {
		o.fps = fps
	}
}

// WithVideoCaptureEnabled explicitly configures the device for video capture mode.
// This is typically auto-detected based on device capabilities, but can be set
// explicitly if needed.
//
// Use this option when:
//   - The device supports both capture and output
//   - You want to ensure capture mode is selected
//   - Auto-detection is not working correctly
//
// Example:
//
//	device.Open("/dev/video0", device.WithVideoCaptureEnabled())
//
// Note: Most devices are either capture OR output devices, making this option
// unnecessary in typical use cases.
func WithVideoCaptureEnabled() Option {
	return func(o *config) {
		o.bufType = v4l2.BufTypeVideoCapture
	}
}

// WithVideoOutputEnabled explicitly configures the device for video output mode.
// This is typically auto-detected based on device capabilities, but can be set
// explicitly if needed.
//
// Use this option when:
//   - The device supports both capture and output
//   - You want to ensure output mode is selected
//   - You're working with a video output device (display, encoder)
//
// Example:
//
//	device.Open("/dev/video1", device.WithVideoOutputEnabled())
//
// Note: Video output support is less common than capture. Most USB cameras
// are capture-only devices.
func WithVideoOutputEnabled() Option {
	return func(o *config) {
		o.bufType = v4l2.BufTypeVideoOutput
	}
}
