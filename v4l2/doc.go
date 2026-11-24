//go:generate go run ./internal/gen

// Package v4l2 provides low-level Go bindings for the Linux Video4Linux2 (V4L2),
// Media Controller, and V4L2 Sub-device userspace APIs.
//
// # Overview
//
// The v4l2 package wraps Linux kernel UAPI structures, constants, and ioctl
// calls via golang.org/x/sys/unix and checked-in generated ABI definitions
// (zconsts_linux.go and ztypes_linux_{amd64,arm64,arm,riscv64}.go). It builds
// with CGO_ENABLED=0 and requires no C compiler or kernel headers at build time.
//
// Most applications should use the higher-level github.com/vladimirvivien/go4vl/device
// package unless direct control over ioctl orchestration, Media Controller graphs,
// sub-device routing, epoll multiplexing, or stateless codec requests is required.
//
// # Functional Areas
//
//   - Capabilities & Device Info: GetCapability (VIDIOC_QUERYCAP), OpenDevice, CloseDevice
//   - Formats & Multi-Planar: PixFormat, PixFormatMPlane, MetaFormat, FrameSizes, FrameIntervals
//   - Selection & Crop: GetSelection, SetSelection (VIDIOC_{G,S}_SELECTION), GetCropRect
//   - Streaming I/O & Lifecycle: MMAP, UserPtr, DMA-BUF (ExportBuffer, ExportBufferPlane), Read/Write, PrepareBuffer, CreateBuffers, RemoveBuffers
//   - Readiness & Multiplexing: Poll, Poller (Linux epoll with eventfd cancellation)
//   - Controls & Stateless Codecs: User controls, ExtControls, H.264/HEVC/VP8/VP9/AV1/MPEG2/FWHT compound controls
//   - Memory-to-Memory & Codecs: M2MSession (dual OUTPUT + CAPTURE queues), EncoderCmd, DecoderCmd
//   - Media Controller & Sub-devices: GetMediaDeviceInfo, GetMediaTopology, SetupMediaLink, ResolveDevnodePath, AllocMediaRequest, Subdev formats/selections/capabilities/routing
//   - Display Timings & EDID: GetDVTimings, SetDVTimings, GetEDID, SetEDID, GetSubdevEDID, SetSubdevEDID
//   - Events & Standards: SubscribeEvent, DequeueEvent, video standards, tuners, modulators, audio/video I/O
//
// # System Requirements
//
//   - Linux kernel 5.10 or later with V4L2 / Media Controller drivers
//   - Supported architectures: linux/amd64, linux/arm64, linux/arm, linux/riscv64 (CGO_ENABLED=0)
//   - Read/write access to /dev/video*, /dev/media*, or /dev/v4l-subdev* nodes
//   - Optional (contributors regenerating ABI definitions): gcc and linux-libc-dev for go generate and -tags=cgo_verify
//
// # Basic Usage
//
// Open a device, inspect capabilities, and negotiate format directly via file descriptor:
//
//	fd, err := v4l2.OpenDevice("/dev/video0", unix.O_RDWR|unix.O_NONBLOCK, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer v4l2.CloseDevice(fd)
//
//	cap, err := v4l2.GetCapability(fd)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Device: %s (%s)\n", cap.Card, cap.Driver)
//
//	if err := v4l2.SetPixFormat(fd, v4l2.PixFormat{
//	    Width:       640,
//	    Height:      480,
//	    PixelFormat: v4l2.PixelFmtYUYV,
//	    Field:       v4l2.FieldNone,
//	}); err != nil {
//	    log.Fatal(err)
//	}
//
// # Thread Safety
//
// Raw V4L2 file descriptors and ioctl state transitions are not goroutine-safe
// at the kernel level. Synchronize concurrent calls targeting the same file
// descriptor or isolate queue driving through Poller on a single control loop.
package v4l2
