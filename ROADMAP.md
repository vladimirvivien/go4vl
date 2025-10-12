# go4vl Roadmap

This roadmap outlines the strategic direction and planned enhancements for the go4vl project. The primary goal remains to port as many V4L2 functionalities as possible, providing idiomatic Go types and abstractions for video capture and processing on Linux systems.

## Project Vision

go4vl aims to be the definitive Go library for V4L2 video capture and streaming, enabling developers to build robust video applications on Linux platforms including Raspberry Pi, embedded systems, and standard Linux distributions without dealing with C interop complexities.

---

## Roadmap Summary

- [ ] [1. Frame Processing Pipeline Performance](#1-frame-processing-pipeline-performance)
- [ ] [2. Migration from CGO to purego](#2-migration-from-cgo-to-purego)
- [ ] [3. io.Reader/io.Writer Interface Support](#3-ioreaderiowriter-interface-support)
- [ ] [4. Complete V4L2 Video Capture Feature Parity](#4-complete-v4l2-video-capture-feature-parity)
- [ ] [5. Multi-Planar Format Support](#5-multi-planar-format-support)
- [ ] [6. Enhanced Image Format Conversion Library](#6-enhanced-image-format-conversion-library)
- [ ] [7. User Pointer and DMA-BUF I/O Support](#7-user-pointer-and-dma-buf-io-support)
- [ ] [8. Video Output Device Support](#8-video-output-device-support)
- [ ] [9. Advanced Extended Controls API](#9-advanced-extended-controls-api)
- [ ] [10. Asynchronous Frame Capture with Select/Poll/Epoll](#10-asynchronous-frame-capture-with-selectpollepoll)
- [ ] [11. Frame Metadata and Timestamping](#11-frame-metadata-and-timestamping)
- [ ] [12. Hardware Codec Integration (H264, HEVC, VP8, VP9)](#12-hardware-codec-integration-h264-hevc-vp8-vp9)
- [ ] [13. Media Controller API Integration](#13-media-controller-api-integration)
- [ ] [14. Performance Optimization and Zero-Copy Enhancements](#14-performance-optimization-and-zero-copy-enhancements)
- [ ] [15. WASM Component Model Plugin System](#15-wasm-component-model-plugin-system)

---

## Enhancement Details

### 1. **Frame Processing Pipeline Performance**

**Status:** Not Started
**Goal:** Eliminate bottlenecks in frame processing pipelines to ensure maximum throughput and minimal latency

**Rationale:**
Frame processing performance is critical for real-time video applications. Bottlenecks can occur at multiple stages: frame capture, memory allocation, channel operations, format conversion, and user processing. Identifying and eliminating these bottlenecks ensures go4vl can handle high-resolution, high-frame-rate scenarios without dropping frames or introducing latency.

**Key Performance Concerns:**
- Frame copy overhead in the capture loop (currently copies every frame)
- Channel blocking and buffering strategies
- Memory allocation patterns and GC pressure
- Goroutine scheduling and synchronization
- Lock contention in multi-device scenarios
- Processing pipeline backpressure handling

**Deliverables:**
- Comprehensive performance profiling across different resolutions and frame rates
- Identify and document all bottlenecks in the capture-to-processing pipeline
- Implement lock-free or low-contention data structures where appropriate
- Add configurable frame dropping strategies (drop oldest, drop newest, block)
- Optimize channel buffer sizing and flow control
- Implement frame pool allocator to reduce GC pressure
- Add pipeline backpressure detection and handling
- Create performance monitoring API (dropped frames, latency, throughput)
- Benchmark suite for different pipeline configurations
- Add CPU and memory profiling examples
- Document best practices for high-performance video processing
- Create reference implementations for common high-throughput scenarios

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

**Deliverables:**
- Evaluate purego compatibility with V4L2 ioctl operations
- Create purego-based syscall wrappers for all V4L2 ioctls
- Migrate struct definitions to pure Go (remove C imports)
- Update build system to remove CGO requirement
- Ensure feature parity with current CGO implementation
- Add CI/CD testing for pure Go builds
- Update documentation reflecting pure Go approach
- Create migration guide for existing users

**Challenges:**
- Complex struct alignment and padding requirements
- Union type handling in pure Go
- Performance comparison and optimization
- Maintaining compatibility across kernel versions

---

### 2. **io.Reader/io.Writer Interface Support**

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
- Metadata delivery alongside frame data
- Error handling and EOF semantics
- Buffer management and zero-copy opportunities

---

### 3. **Complete V4L2 Video Capture Feature Parity**

**Status:** In Progress (significant coverage exists, gaps remain)
**Goal:** Achieve 100% coverage of V4L2 video capture APIs and capabilities

**Rationale:**
While go4vl covers core V4L2 functionality, some advanced features remain unimplemented. Complete feature parity ensures users never need to drop to C or use syscall directly.

**Missing Features to Implement:**
- Selection/Compose API (VIDIOC_G_SELECTION, VIDIOC_S_SELECTION)
- Video standards enumeration and selection (NTSC, PAL, SECAM)
- Video timings API for DV and HDMI sources
- Priority and exclusive access control
- Advanced cropping capabilities
- Buffer export/import for cross-device workflows
- Query extended capabilities flags
- All buffer timestamp modes
- Request API support for stateless codecs
- All field order modes
- Complete control flags support
- Menu control integer values

**Deliverables:**
- Audit current V4L2 API coverage against kernel documentation
- Implement remaining VIDIOC_* ioctls relevant to video capture
- Add complete constant definitions for all flags and enums
- Create comprehensive examples demonstrating each capability
- Add integration tests for all implemented features
- Document feature availability by kernel version
- Create feature detection utilities
- Update API reference with complete V4L2 mapping

---

### 4. **Multi-Planar Format Support**

**Status:** Not Started
**Goal:** Add comprehensive support for multi-planar pixel formats (V4L2_BUF_TYPE_VIDEO_CAPTURE_MPLANE)

**Rationale:**
Many modern cameras and video processors use multi-planar formats where Y, U, and V components are stored in separate memory planes. This is essential for advanced codec support and high-performance video processing.

**Deliverables:**
- Implement multi-planar buffer type support in `v4l2/streaming.go`
- Add multi-planar format descriptors and negotiation
- Extend `device.Device` to handle planar buffers
- Create examples demonstrating multi-planar capture
- Update API documentation with multi-planar usage patterns

---

### 2. **Enhanced Image Format Conversion Library**

**Status:** Partially Implemented (YUYV conversion disabled)
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

---

### 3. **User Pointer and DMA-BUF I/O Support**

**Status:** Not Started (only MMAP currently supported)
**Goal:** Implement V4L2_MEMORY_USERPTR and V4L2_MEMORY_DMABUF I/O methods

**Rationale:**
User pointer mode allows application-allocated buffers, enabling zero-copy integration with other libraries. DMA-BUF support enables zero-copy sharing between V4L2 devices and GPU/display subsystems, critical for high-performance pipelines.

**Deliverables:**
- Implement `IOTypeUserPtr` support in streaming layer
- Add user-provided buffer allocation and management
- Implement `IOTypeDMABuf` for inter-device buffer sharing
- Add DMA-BUF file descriptor handling and synchronization
- Create examples for user pointer mode
- Create examples for DMA-BUF sharing with DRM/KMS
- Document memory management best practices for each I/O type
- Add performance comparison benchmarks

---

### 4. **Video Output Device Support**

**Status:** Partially Implemented (output path exists but incomplete)
**Goal:** Complete implementation for video output devices (V4L2_BUF_TYPE_VIDEO_OUTPUT)

**Rationale:**
V4L2 supports video output for display devices, video encoders, and hardware accelerators. Currently go4vl primarily focuses on capture. Output support enables frame injection, video encoding pipelines, and display control.

**Deliverables:**
- Complete `Device.SetInput()` implementation for output devices
- Add output buffer queue management
- Implement frame timing control for output
- Add output format negotiation and validation
- Create video playback/display examples
- Create encoder pipeline examples (with hardware encoders)
- Document output device workflows and use cases

---

### 5. **Advanced Extended Controls API**

**Status:** Basic support implemented
**Goal:** Comprehensive extended controls with validation, compound controls, and control events

**Rationale:**
Extended controls enable access to codec parameters, camera sensors, and advanced device features. Current implementation covers basics but lacks compound control support and event notification.

**Deliverables:**
- Implement compound control support (arrays, structs)
- Add control event subscription and monitoring (VIDIOC_SUBSCRIBE_EVENT)
- Implement control priority and error handling improvements
- Add control validation and bounds checking enhancements
- Create codec control profiles (H264, VP8, VP9, HEVC)
- Add camera control presets (exposure, white balance, focus)
- Create advanced controls example with live adjustment
- Document control classes and codec-specific parameters

---

### 6. **Asynchronous Frame Capture with Select/Poll/Epoll**

**Status:** Currently uses select internally
**Goal:** Expose low-level async I/O options for advanced use cases

**Rationale:**
Current implementation uses `sys.Select` internally but doesn't expose file descriptor for external event loops. Users may want to integrate V4L2 capture into custom event loops or multiplexed I/O systems.

**Deliverables:**
- Add optional async capture mode with exposed file descriptor
- Implement epoll support for better scalability with multiple devices
- Add timeout and deadline control for frame capture
- Create examples showing integration with Go's context patterns
- Create example integrating multiple devices in single event loop
- Document async patterns and best practices
- Add performance benchmarks comparing async methods

---

### 7. **Frame Metadata and Timestamping**

**Status:** Buffer timestamps exist but not fully exposed
**Goal:** Comprehensive frame metadata including timestamps, sequence numbers, and flags

**Rationale:**
Accurate frame timing is critical for video synchronization, frame rate analysis, and time-based processing. Current implementation captures metadata but doesn't expose it through idiomatic interfaces.

**Deliverables:**
- Create `FrameMetadata` struct with timestamp, sequence, and flags
- Extend output channel to deliver metadata alongside frames (or add parallel channel)
- Add frame type detection (keyframe, P-frame, B-frame)
- Implement frame statistics collection (dropped frames, errors)
- Add timestamp source configuration and queries
- Create examples demonstrating frame timing analysis
- Document synchronization strategies for multi-device capture

---

### 8. **Hardware Codec Integration (H264, HEVC, VP8, VP9)**

**Status:** Basic extended controls exist
**Goal:** High-level APIs for hardware encoder/decoder configuration

**Rationale:**
Many platforms (Raspberry Pi, Jetson, modern CPUs) have hardware video codecs accessible via V4L2. Current low-level control support needs higher-level abstractions for practical use.

**Deliverables:**
- Create codec configuration profiles for H264/HEVC/VP8/VP9
- Add encoder preset system (quality vs speed trade-offs)
- Implement bitrate control modes (CBR, VBR, CQP)
- Add GOP structure configuration helpers
- Create decoder capabilities query and setup
- Add examples for hardware encoding (Raspberry Pi H264, etc.)
- Add examples for hardware decoding
- Document platform-specific codec availability and capabilities
- Create codec benchmarking utilities

---

### 9. **Media Controller API Integration**

**Status:** Basic media info query exists
**Goal:** Full Media Controller API for complex video pipelines

**Rationale:**
Complex camera systems (CSI-2, MIPI) require Media Controller API to configure sensor, ISP, and capture device topology. Essential for embedded platforms and advanced camera modules.

**Deliverables:**
- Implement media entity enumeration and topology discovery
- Implement pipeline setup and validation
- Add subdevice format negotiation
- Create CSI-2 camera setup examples (Raspberry Pi Camera Module)
- Create MIPI camera configuration examples
- Document Media Controller concepts and workflows
- Add pipeline visualization and debugging tools

---

### 10. **Performance Optimization and Zero-Copy Enhancements**

**Status:** Basic zero-copy with MMAP implemented
**Goal:** Minimize memory copies and CPU overhead for maximum throughput

**Rationale:**
Video streaming is performance-critical. Even with MMAP, the current implementation copies frames from mapped buffers. For high-resolution or high-frame-rate capture, additional optimizations are needed.

**Deliverables:**
- Implement optional zero-copy frame access with buffer lifecycle management
- Add buffer recycling to avoid allocation overhead
- Implement frame pool management for high-throughput scenarios
- Optimize goroutine scheduling and channel operations
- Add CPU usage profiling and optimization
- Create performance testing suite with various resolutions/formats
- Add benchmark results and optimization guide to documentation
- Create high-throughput example (1080p60, 4K30)

---

### 14. **WASM Component Model Plugin System**

**Status:** Not Started
**Goal:** Enable extensible video processing pipelines using WebAssembly Component Model (WIT)

**Rationale:**
A plugin system allows users to extend go4vl with custom frame processors, filters, encoders, and analyzers without modifying the core library. Using WASM Component Model with WIT (WebAssembly Interface Types) provides:
- Language-agnostic plugins (write in Rust, C, Go, etc.)
- Sandboxed execution for safety and security
- Performance approaching native code
- Cross-platform compatibility
- Hot-reloadable processing pipelines

**Use Cases:**
- Custom image filters and effects
- AI/ML inference on video frames
- Custom codec implementations
- Specialized format converters
- Real-time video analytics
- Edge detection, object tracking, etc.

**Deliverables:**
- Design WIT interface definitions for video processing plugins
- Implement WASM runtime integration (wazero or wasmtime-go)
- Create plugin lifecycle management (load, initialize, execute, unload)
- Add frame buffer passing between host and plugin (zero-copy where possible)
- Implement plugin discovery and registration system
- Create plugin development SDK and templates
- Add example plugins in multiple languages (Rust, TinyGo, C)
- Document plugin API and development workflow
- Create benchmarks comparing plugin vs native performance
- Add plugin marketplace/registry documentation

**WIT Interface Design:**
```wit
// Example WIT interface for frame processor plugin
interface frame-processor {
  // Initialize plugin with configuration
  init: func(config: string) -> result<_, string>

  // Process a single frame
  process-frame: func(input: list<u8>, width: u32, height: u32, format: u32)
    -> result<list<u8>, string>

  // Query plugin capabilities
  get-capabilities: func() -> plugin-info

  // Cleanup resources
  cleanup: func()
}
```

**Integration Points:**
- Pipeline stage insertion (pre-capture, post-capture, pre-encode)
- Format conversion plugins
- Codec plugins for custom formats
- Control plugins for device automation

---

## Additional Improvements

### Documentation
- Create comprehensive API reference with detailed examples
- Add architecture overview and design principles document
- Create platform-specific guides (Raspberry Pi, Jetson, x86)
- Add troubleshooting guide for common issues

### Tooling
- Create device capability inspection CLI tool
- Add format conversion benchmarking tool
- Create video capture GUI example (using fyne or similar)
- Add device stress testing utility

---

## Long-Term Vision

- **Cross-Device Synchronization:** Support synchronized capture from multiple devices with frame alignment
- **Cloud Integration:** Examples for streaming to cloud services (RTMP, WebRTC, HLS)
- **AI/ML Integration:** Examples integrating with Go ML frameworks and ONNX runtime for real-time inference
- **Embedded Optimization:** Specialized builds for resource-constrained devices with reduced memory footprint
- **Time-Code and Genlock Support:** Professional video production features for broadcast applications

---

**Note:** This roadmap represents the project's aspirations and may evolve based on community contributions, user feedback, and emerging requirements.
