// Package device provides a high-level, idiomatic Go interface for Video4Linux2
// (V4L2) capture, output, and stateful codec devices.
//
// # Overview
//
// Package device manages device discovery, functional option configuration,
// buffer lifecycle, and streaming over MMAP, UserPtr, DMA-BUF, and direct
// Read/Write I/O. It exposes two capture consumption models:
//
//   - Device.Frames(ctx): Zero-copy iter.Seq2[*Frame, error] iterator backed by v4l2.Poll
//     that yields directly into mapped buffers and re-queues buffers automatically.
//   - Device.GetFrames(): Buffered channel (<-chan *Frame) backed by FramePool for
//     asynchronous pipelines; callers return buffers via Frame.Release().
//
// For video output nodes (such as hardware outputs or v4l2loopback), Device
// provides WriteFrame and WriteFramePlanes. Stateful hardware encoders and
// decoders are exposed via OpenEncoder and OpenDecoder.
//
// # Basic Usage
//
// Capture frames using the zero-copy Frames iterator:
//
//	dev, err := device.Open("/dev/video0",
//	    device.WithBufferSize(4),
//	    device.WithPixFormat(v4l2.PixFormat{
//	        Width:       1280,
//	        Height:      720,
//	        PixelFormat: v4l2.PixelFmtMJPEG,
//	    }),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer dev.Close()
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	if err := dev.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer dev.Stop()
//
//	for frame, err := range dev.Frames(ctx) {
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    processFrame(frame.Data)
//	}
//
// Retain frame bytes across loop iterations by calling frame.Clone().
//
// # Device Discovery
//
// Enumerate available /dev/video*, /dev/media*, and /dev/v4l-subdev* nodes:
//
//	paths, err := device.GetAllDevicePaths()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, p := range paths {
//	    fmt.Println(p)
//	}
//
// # Configuration Options
//
//   - WithBufferSize: Number of streaming ring buffers (typically 2-4)
//   - WithPixFormat: Single-planar PixFormat (width, height, FourCC, colorspace)
//   - WithPixFormatMPlane: Multi-planar PixFormatMPlane configuration
//   - WithFPS: Target stream frame rate
//   - WithBufferType: Explicit v4l2.BufType (capture, output, or MPLANE variants)
//   - WithIOType: Memory mode (IOTypeMMAP, IOTypeUserPtr, IOTypeDMABuf)
//   - WithIOMethod: Streaming ring vs direct Read/Write (IOMethodStreaming, IOMethodReadWrite)
//   - WithTimestampMode: Clock conversion mode (TimestampDefault, TimestampAuto, TimestampMonoToAbs)
//   - WithVideoCaptureEnabled / WithVideoOutputEnabled: Explicit stream direction
//
// # System Constraints & Thread Safety
//
//   - Linux only (kernel 5.10 or later); builds with CGO_ENABLED=0 on amd64, arm64, arm, and riscv64.
//   - Device methods (Open, Start, Stop, Close, control/format accessors) are not
//     concurrent-safe and should be invoked from a single control goroutine.
//   - The channel returned by GetFrames() may be consumed concurrently alongside Frame.Release().
package device
