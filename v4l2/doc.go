// Package v4l2 provides low-level Go bindings for the Video4Linux2 (V4L2) API.
// This package wraps the Linux kernel's V4L2 userspace API, providing direct access
// to video capture and output devices through system calls and ioctls.
//
// # Overview
//
// The v4l2 package is the foundation layer of go4vl, providing direct mappings to
// V4L2 kernel structures and constants. It handles the low-level interactions with
// device drivers through ioctl system calls and memory mapping operations.
//
// Most applications should use the higher-level device package instead of this
// package directly, unless fine-grained control over V4L2 operations is required.
//
// # Architecture
//
// The package is organized into functional areas:
//
//   - Capabilities: Device capability detection and querying
//   - Controls: Device controls for brightness, contrast, etc.
//   - Formats: Pixel formats, frame sizes, and frame intervals
//   - Streaming: Buffer management and streaming I/O operations
//   - Errors: V4L2-specific error handling
//
// # Core Types
//
//   - Capability: Device capabilities and identification
//   - PixFormat: Pixel format configuration
//   - Buffer: Frame buffer management
//   - Control: Device control parameters
//   - StreamParam: Streaming parameters like FPS
//
// # System Requirements
//
//   - Linux kernel 5.10 or later with V4L2 support
//   - CGO enabled for C bindings
//   - Access to /dev/video* devices (typically requires video group membership)
//
// # Basic Usage
//
// Low-level device interaction example:
//
//	// Open device
//	fd, err := v4l2.OpenDevice("/dev/video0", syscall.O_RDWR, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer v4l2.CloseDevice(fd)
//
//	// Query capabilities
//	cap, err := v4l2.GetCapability(fd)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Device: %s\n", cap.Card)
//
//	// Set format
//	format := v4l2.PixFormat{
//	    Width:       640,
//	    Height:      480,
//	    PixelFormat: v4l2.PixelFmtYUYV,
//	    Field:       v4l2.FieldNone,
//	}
//	if err := v4l2.SetPixFormat(fd, format); err != nil {
//	    log.Fatal(err)
//	}
//
// # Memory Management
//
// The package supports multiple I/O methods:
//
//   - Memory-mapped I/O (MMAP): Zero-copy, most efficient
//   - User pointer I/O: Application-managed buffers
//   - Read/Write I/O: Simple but less efficient
//   - DMA buffer I/O: For buffer sharing between devices
//
// Currently, only memory-mapped I/O is fully implemented and tested.
//
// # Error Handling
//
// V4L2 operations return specific error types that indicate the nature of failures:
//
//   - ErrorSystem: System-level errors
//   - ErrorBadArgument: Invalid parameters
//   - ErrorUnsupportedFeature: Feature not supported by device
//   - ErrorNotFound: Resource not found
//   - ErrorBusy: Device or resource busy
//   - ErrorTimeout: Operation timed out
//   - ErrorDeviceRemoved: Device disconnected
//
// # CGO Considerations
//
// This package uses CGO to interface with the Linux kernel headers.
// Key considerations:
//
//   - Build times are longer due to C compilation
//   - Cross-compilation requires appropriate C toolchain
//   - Some Go tools have limited CGO support
//   - Memory passed to C must be carefully managed
//
// # Thread Safety
//
// Most V4L2 operations are NOT thread-safe at the kernel level.
// Applications must synchronize access to the same file descriptor
// across goroutines. The higher-level device package handles this.
//
// # Constants and Flags
//
// The package exports numerous constants that map directly to V4L2 kernel definitions:
//
//   - Capability flags (Cap*)
//   - Pixel formats (PixelFmt*)
//   - Control types (CtrlType*)
//   - Buffer flags (BufFlag*)
//   - Field orders (Field*)
//   - Color spaces (ColorSpace*)
//
// # References
//
//   - V4L2 API Specification: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/v4l2.html
//   - Linux Media Subsystem: https://linuxtv.org/
//   - Kernel Headers: /usr/include/linux/videodev2.h
package v4l2