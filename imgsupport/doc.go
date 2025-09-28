// Package imgsupport provides image format conversion utilities for V4L2 captured frames.
// This package contains helper functions to convert between different pixel formats
// commonly used in video capture applications.
//
// # Overview
//
// The imgsupport package bridges the gap between raw V4L2 pixel formats and
// standard Go image formats. It provides converters for transforming raw video
// frames into formats suitable for display, storage, or further processing.
//
// # Supported Conversions
//
// Currently supported format conversions:
//   - YUYV to JPEG: Convert YUV 4:2:2 packed format to JPEG (experimental)
//
// # Pixel Format Background
//
// V4L2 devices often capture video in YUV formats rather than RGB because:
//   - YUV is more efficient for video compression
//   - Many camera sensors natively output YUV
//   - YUV separates luminance from chrominance, useful for video processing
//
// Common YUV formats:
//   - YUYV (YUV 4:2:2): Packed format with 2 pixels sharing color data
//   - NV12 (YUV 4:2:0): Planar format common in hardware encoders
//   - I420 (YUV 4:2:0): Planar format used in many video codecs
//
// # Usage Example
//
//	// Capture YUYV frame from device
//	frame := <-device.GetOutput()
//
//	// Convert to JPEG
//	jpegData, err := imgsupport.Yuyv2Jpeg(640, 480, frame)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Save JPEG file
//	os.WriteFile("output.jpg", jpegData, 0644)
//
// # Performance Considerations
//
// Format conversion can be CPU-intensive. Consider:
//   - Using hardware-accelerated formats when possible (MJPEG, H264)
//   - Implementing conversions in parallel for multiple frames
//   - Caching conversion results if frames are reused
//   - Using lower resolutions when real-time performance is critical
//
// # Future Enhancements
//
// Planned additions to this package:
//   - RGB24/BGR24 conversion support
//   - NV12/I420 planar format support
//   - Hardware-accelerated conversions (when available)
//   - Direct integration with Go's image package types
//
// # Limitations
//
// - YUYV to JPEG conversion is currently experimental and disabled
// - No support for planar YUV formats yet
// - No hardware acceleration support
// - Conversion quality/speed trade-offs not configurable
package imgsupport
