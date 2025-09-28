// Package device provides a high-level, idiomatic Go interface for Video4Linux2 (V4L2) devices.
// It abstracts the complexity of the low-level V4L2 API and provides a simple, channel-based
// interface for video capture and streaming.
//
// # Overview
//
// The device package is the primary interface for working with V4L2 devices in go4vl.
// It handles device initialization, configuration, buffer management, and streaming,
// allowing developers to focus on processing video frames rather than managing low-level details.
//
// # Key Features
//
//   - Simple device open/close lifecycle management
//   - Automatic capability detection and configuration
//   - Zero-copy frame capture using memory-mapped buffers
//   - Go channels for frame delivery
//   - Context-based cancellation support
//   - Flexible configuration through functional options
//
// # Basic Usage
//
// Opening and using a video device typically follows this pattern:
//
//	// Open the device
//	dev, err := device.Open("/dev/video0",
//	    device.WithBufferSize(4),
//	    device.WithPixFormat(v4l2.PixFormat{
//	        Width: 1920,
//	        Height: 1080,
//	        PixelFormat: v4l2.PixelFmtMJPEG,
//	    }),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer dev.Close()
//
//	// Start streaming
//	ctx := context.Background()
//	if err := dev.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer dev.Stop()
//
//	// Capture frames
//	for frame := range dev.GetOutput() {
//	    // Process frame data
//	    processFrame(frame)
//	}
//
// # Device Discovery
//
// Devices are typically found at /dev/video* paths. You can enumerate available devices:
//
//	devices, err := device.GetDevices()
//	for _, d := range devices {
//	    fmt.Printf("Found device: %s\n", d.Path)
//	}
//
// # Configuration Options
//
// The package provides several configuration options through the functional options pattern:
//
//   - WithBufferSize: Set the number of buffers for streaming (affects latency vs. reliability)
//   - WithPixFormat: Configure pixel format, resolution, and color space
//   - WithFPS: Set the desired frame rate
//   - WithIOType: Select I/O method (currently only memory-mapped I/O is supported)
//
// # Pixel Formats
//
// Common pixel formats include:
//
//   - PixelFmtMJPEG: Motion JPEG (compressed, widely supported)
//   - PixelFmtYUYV: YUV 4:2:2 (uncompressed, good quality)
//   - PixelFmtH264: H.264 video (compressed, efficient)
//   - PixelFmtRGB24: RGB (uncompressed, direct use)
//
// # Error Handling
//
// The package returns detailed errors that can be inspected:
//
//	if errors.Is(err, v4l2.ErrorUnsupportedFeature) {
//	    // Handle unsupported feature
//	}
//
// # Performance Considerations
//
//   - Use memory-mapped I/O for best performance (default)
//   - Configure appropriate buffer count (2-4 buffers typical)
//   - Process frames quickly to avoid drops
//   - Consider frame format impact (MJPEG vs. raw formats)
//
// # Limitations
//
//   - Linux only (V4L2 is Linux-specific)
//   - Requires kernel 5.10 or later
//   - Currently supports only memory-mapped I/O
//   - Single-consumer model (one goroutine reading frames)
//
// # Thread Safety
//
// Device methods are NOT thread-safe and should be called from a single goroutine.
// The output channel returned by GetOutput() is safe for concurrent reads.
package device
