# go4vl Extended Roadmap - Strategic Initiatives

This document outlines strategic initiatives and optimization efforts that extend beyond V4L2 API parity. These are enhancements that improve performance, developer experience, and maintainability of the go4vl library.

**Primary Roadmap**: See [ROADMAP.md](./ROADMAP.md) for V4L2 API feature parity tracking (the main project roadmap).

## Purpose

While the primary roadmap focuses on implementing V4L2 features to achieve API parity, this extended roadmap covers:
- Performance optimizations beyond API implementation
- Developer experience improvements
- Build system enhancements (CGO → purego migration)
- Additional abstractions and convenience APIs
- Integration with other subsystems (io.Reader/Writer, etc.)
- Future-looking enhancements (WASM plugins, etc.)

## Project Vision

go4vl aims to be the definitive Go library for V4L2 video capture and streaming, enabling developers to build robust video applications on Linux platforms including Raspberry Pi, embedded systems, and standard Linux distributions without dealing with C interop complexities.

---

## Roadmap Summary

- [x] [1. Frame Processing Pipeline Performance](#1-frame-processing-pipeline-performance) ✅ **Phase 1 Complete**
- [ ] [2. Migration from CGO to purego](#2-migration-from-cgo-to-purego)
- [ ] [3. io.Reader/io.Writer Interface Support](#3-ioreaderiowriter-interface-support)
- [ ] [4. Enhanced Image Format Conversion Library](#4-enhanced-image-format-conversion-library)
- [ ] [5. WASM Component Model Plugin System](#5-wasm-component-model-plugin-system)

---

## Enhancement Details

### 1. **Frame Processing Pipeline Performance**

**Status:** ✅ **Phase 1 Completed** - Memory Allocation Optimization
**Goal:** Eliminate bottlenecks in frame processing pipelines to ensure maximum throughput and minimal latency

**Rationale:**
Frame processing performance is critical for real-time video applications. Bottlenecks can occur at multiple stages: frame capture, memory allocation, channel operations, format conversion, and user processing. Identifying and eliminating these bottlenecks ensures go4vl can handle high-resolution, high-frame-rate scenarios without dropping frames or introducing latency.

**Completed Performance Improvements (Phase 1):**
- ✅ **Memory allocation bottleneck eliminated:** Implemented FramePool with sync.Pool
  - ~1,200x faster buffer allocation (22ns vs 28,000ns per 614KB frame)
  - 99.996% reduction in memory allocated per operation (26 B vs 614 KB)
  - Effective allocation elimination after pool warmup
- ✅ **GC pressure reduction:** Buffer pooling dramatically reduces garbage collection overhead
- ✅ **Separate streaming loops:** Optimized paths for GetOutput() and GetFrames() APIs
- ✅ **Benchmark suite:** Comprehensive benchmarks for pool operations and allocation patterns
- ✅ **Documentation:** Best practices for high-performance video processing
- ✅ **Race condition fixes:** Proper goroutine synchronization in Device.Stop()

**Key Performance Concerns (Remaining - Phase 2):**
- Frame copy overhead in the capture loop (currently copies every frame from mmap)
- Channel blocking and buffering strategies
- Goroutine scheduling and synchronization
- Lock contention in multi-device scenarios
- Processing pipeline backpressure handling

**Remaining Deliverables (Phase 2):**
- Comprehensive performance profiling across different resolutions and frame rates
- Identify remaining bottlenecks in the capture-to-processing pipeline
- Implement lock-free or low-contention data structures where appropriate
- Add configurable frame dropping strategies (drop oldest, drop newest, block)
- Optimize channel buffer sizing and flow control
- Add pipeline backpressure detection and handling
- Create performance monitoring API (dropped frames, latency, throughput)
- Add CPU and memory profiling examples
- Create reference implementations for common high-throughput scenarios
- True zero-copy access (expose mmap buffers directly with lifecycle management)

**Target Performance Metrics:**
- 1080p @ 60fps with < 5% CPU usage
- 4K @ 30fps with minimal frame drops
- Multi-device capture (4+ cameras) without contention
- Sub-millisecond frame delivery latency
- Zero frame drops under sustained load

**Profiling and Optimization Areas:**
- Frame delivery path (device → mmap → copy → channel → user)
- Buffer lifecycle management
- Context switching and goroutine overhead
- Memory allocation patterns
- Cache-line optimization for hot paths

---

### 2. **Migration from CGO to purego**

**Status:** Not Started
**Goal:** Eliminate CGO dependency by migrating to pure Go C bindings using the purego library

**Rationale:**
CGO introduces significant complexity, build overhead, and cross-compilation challenges. It prevents pure Go toolchain benefits like easy cross-compilation, faster builds, and simpler dependency management. The purego library enables calling C functions from pure Go without CGO, making the library more accessible and maintainable.

**Benefits:**
- Faster compilation times (no C compiler required)
- Simplified cross-compilation for ARM/ARM64/x86_64
- Better compatibility with Go's module system
- Reduced binary size in many cases
- Easier to contribute (no C toolchain knowledge required)
- Full support for `go install` without C dependencies
- Better integration with Go's debugging and profiling tools

**Deliverables:**
- Evaluate purego compatibility with V4L2 ioctl operations
- Create purego-based syscall wrappers for all V4L2 ioctls
- Migrate struct definitions to pure Go (remove C imports)
- Update build system to remove CGO requirement
- Ensure feature parity with current CGO implementation
- Add CI/CD testing for pure Go builds
- Update documentation reflecting pure Go approach
- Create migration guide for existing users
- Performance benchmarking (purego vs CGO)

**Challenges:**
- Complex struct alignment and padding requirements
- Union type handling in pure Go
- Performance comparison and optimization
- Maintaining compatibility across kernel versions
- Ensuring exact C struct layout matches

**Priority:** Medium (quality-of-life improvement, doesn't add features)

---

### 3. **io.Reader/io.Writer Interface Support**

**Status:** Not Started
**Goal:** Provide standard io.Reader/io.Writer interfaces alongside existing channel-based API

**Rationale:**
Channels are idiomatic for concurrent frame delivery, but many Go libraries and tools work with io.Reader/io.Writer interfaces. Supporting both patterns increases composability with the broader Go ecosystem (image processing, encoding, streaming libraries).

**Deliverables:**
- Implement `io.Reader` interface for frame capture streams
- Implement `io.Writer` interface for video output devices
- Add frame framing/deframing for Reader/Writer (handle frame boundaries)
- Support both blocking and non-blocking I/O modes
- Create adapter utilities between channel and io interfaces
- Add examples using io.Reader with standard library (io.Copy, bufio, etc.)
- Add examples integrating with popular streaming libraries
- Document performance characteristics of each approach
- Provide guidance on when to use channels vs io interfaces

**API Design Considerations:**
- Frame boundary handling in streaming mode
- Error propagation in io interface vs channel
- Metadata delivery (timestamps, sequence numbers)
- Cancellation and timeout handling
- Buffer size configuration

**Example Usage:**
```go
// Capture as io.Reader
dev, _ := device.Open("/dev/video0")
reader := dev.AsReader()

// Copy directly to encoder
encoder := h264.NewEncoder(writer)
io.Copy(encoder, reader)

// Or use with standard library
scanner := bufio.NewScanner(reader)
scanner.Split(device.SplitFrames) // Custom split function
for scanner.Scan() {
    frame := scanner.Bytes()
    // Process frame
}
```

**Priority:** Low (convenience feature, channels work well)

---

### 4. **Enhanced Image Format Conversion Library**

**Status:** Partial (imgsupport package exists but incomplete)
**Goal:** Complete pixel format conversion utilities for common V4L2 formats

**Rationale:**
Current `imgsupport` package has disabled YUYV conversion. Users frequently need to convert between raw formats (YUYV, NV12, YV12) and standard image formats (JPEG, PNG) for processing or display.

**Deliverables:**
- Complete and test YUYV to JPEG/PNG conversion
- Add NV12, YV12, and I420 format converters
- Implement RGB format family converters (RGB24, RGB32, BGR24)
- Add hardware-accelerated conversion paths where available
- Benchmark and optimize conversion performance
- Create format conversion example application
- Document supported conversion paths and performance characteristics
- Add color space conversion utilities

**Use Cases:**
- Preview generation (raw → JPEG for web display)
- Format normalization (various raw → RGB for processing)
- Thumbnail creation
- Integration with image processing libraries

**Priority:** Medium (common user need for visualization/display)

---

### 5. **WASM Component Model Plugin System**

**Status:** Not Started (Future Vision)
**Goal:** Enable video processing pipelines using WASM Component Model plugins

**Rationale:**
WASM Component Model provides sandboxed, portable, language-agnostic plugin architecture. This would enable users to write frame processing filters in any language (Rust, C++, Go, etc.) and run them safely within go4vl pipelines without CGO or native code vulnerabilities.

**Use Cases:**
- User-provided custom filters (face detection, OCR, object tracking)
- Third-party codec implementations
- Proprietary algorithms without source disclosure
- Hot-reloadable processing pipelines
- Cross-platform filter libraries

**Deliverables:**
- Design WASM Component interface for frame processing
- Implement WASM runtime integration (wazero or wasmer)
- Create frame data passing between Go ↔ WASM
- Add plugin lifecycle management (load, unload, hot-reload)
- Create example WASM filters (blur, edge detection, etc.)
- Document WASM plugin development workflow
- Benchmark WASM vs native performance
- Create plugin marketplace/registry concept

**Example API:**
```go
dev, _ := device.Open("/dev/video0")

// Load WASM plugin
filter, _ := wasm.LoadPlugin("face_detect.wasm")

// Apply in pipeline
for frame := range dev.GetFrames() {
    processed := filter.Process(frame.Data)
    // Use processed frame
    frame.Release()
}
```

**Challenges:**
- Performance overhead of WASM boundary crossings
- Memory management between Go and WASM
- Ensuring real-time performance for video pipelines
- Plugin API versioning and compatibility

**Priority:** Low (future vision, significant complexity)

---

## Relationship to Primary Roadmap

This extended roadmap complements the primary [ROADMAP.md](./ROADMAP.md) by focusing on enhancements that go beyond V4L2 API parity:

**Primary Roadmap** focuses on:
- Implementing V4L2 ioctls and features
- API completeness (controls, formats, events, etc.)
- Device type support (capture, output, M2M, codec)
- Feature parity with kernel documentation

**Extended Roadmap** focuses on:
- Performance optimization beyond basic implementation
- Developer experience improvements
- Build system modernization
- Integration with Go ecosystem (io.Reader/Writer)
- Future-looking capabilities (WASM plugins)

Both roadmaps work together to make go4vl the definitive V4L2 library for Go.

---

## Notes

- Items here should NOT duplicate features in the primary roadmap
- Focus is on enhancements, optimizations, and convenience features
- Priority is generally lower than V4L2 parity work
- Performance work (item #1) is ongoing alongside parity work

---

**Last Updated**: 2025-10-19
**Status**: Extended roadmap cleaned and deduplicated
