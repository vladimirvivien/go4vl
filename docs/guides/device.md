# Device Lifecycle & Options

A V4L2 node must be opened, inspected or configured, initialized for buffer I/O, started, drained on teardown, and closed in strict order. The `device` package captures that boilerplate inside `*device.Device`, using functional options at open time so configuration errors surface immediately before streaming begins.

## Discovering available devices

Before hardcoding `/dev/video0`, you can query mounted character nodes on the system. `device.GetAllDevicePaths()` scans `/sys/class/video4linux` and resolves `/dev/<name>` entries (falling back to scanning `/dev` when sysfs is absent in minimal containers):

```go
paths, err := device.GetAllDevicePaths()
if err != nil {
	log.Fatal(err)
}
for _, p := range paths {
	isDev, _ := device.IsDevice(p)
	fmt.Printf("%s (device=%v)\n", p, isDev)
}
```

Each path can then be opened safely with `device.Open` (or `v4l2.OpenDevice` for Media Controller `/dev/media*` and `/dev/v4l-subdev*` nodes) to query exact driver capabilities before selecting a target node.

## Opening and configuring a `Device`

`device.Open(path, options...)` opens the node, reads `v4l2.Capability` via `VIDIOC_QUERYCAP`, determines the default capture/output buffer type, applies your functional options in order, and validates video I/O support when requested:

```go
dev, err := device.Open(
	"/dev/video0",
	device.WithBufferSize(4),
	device.WithPixFormat(v4l2.PixFormat{
		Width:       1280,
		Height:      720,
		PixelFormat: v4l2.PixelFmtMJPEG,
		Field:       v4l2.FieldNone,
	}),
	device.WithFPS(30),
	device.WithTimestampMode(v4l2.TimestampAuto),
)
if err != nil {
	log.Fatal(err)
}
defer dev.Close()
```

Here is the full set of `device.Option` helpers available at open time:

| Option | Purpose | Notes |
|---|---|---|
| `WithBufferSize(size uint32)` | Sets requested driver buffer count | Typical capture values are `2`–`4` |
| `WithPixFormat(pixFmt v4l2.PixFormat)` | Applies single-planar width, height, and pixel format | Calls `VIDIOC_S_FMT` during `Open` |
| `WithPixFormatMPlane(pixFmt v4l2.PixFormatMPlane)` | Applies multi-planar format configuration | Pair with `WithBufferType(v4l2.BufTypeVideoCaptureMPlane)` or `OutputMPlane` |
| `WithFPS(fps uint32)` | Configures streaming capture/output frame rate | Calls `SetFrameRate` (`VIDIOC_S_PARM`) |
| `WithBufferType(bufType v4l2.BufType)` | Overrides auto-detected buffer queue type | Required to force `MPLANE` or `Output` on dual-cap devices |
| `WithIOType(ioType v4l2.IOType)` | Selects backing memory mode | `v4l2.IOTypeMMAP` (default), `v4l2.IOTypeUserPtr`, `v4l2.IOTypeDMABuf` |
| `WithIOMethod(method IOMethod)` | Chooses streaming ring I/O vs direct syscalls | `IOMethodStreaming` (default) or `IOMethodReadWrite` |
| `WithVideoCaptureEnabled()` | Selects video capture direction explicitly | Useful on loopback or M2M-style nodes |
| `WithVideoOutputEnabled()` | Selects video output direction explicitly | Required for `WriteFrame` on dual-capability nodes |
| `WithTimestampMode(mode v4l2.TimestampMode)` | Controls kernel-to-`time.Time` clock mapping | `TimestampDefault`, `TimestampAuto`, or `TimestampMonoToAbs` |

## Inspecting capabilities and runtime properties

Once open, inspect negotiated device properties without issuing redundant `ioctl` calls:

```go
cap := dev.Capability()
fmt.Printf("Driver=%s Card=%s Bus=%s StreamOK=%v\n",
	cap.Driver, cap.Card, cap.BusInfo, cap.IsStreamingSupported())

pix, err := dev.GetPixFormat()
if err == nil {
	fmt.Printf("Active Format: %dx%d size=%dB\n", pix.Width, pix.Height, pix.SizeImage)
}
```

Use `dev.Fd()` if you need to pass the underlying file descriptor directly into low-level `v4l2` functions such as `v4l2.GetSelection`, `v4l2.Poll`, or `v4l2.Poller`.

## Lifecycle and concurrency guarantees

Follow a strict lifecycle sequence for every streaming session:

1. Call `device.Open(...)` once and keep `defer dev.Close()` paired with it.
2. Call `dev.Start(ctx)` to allocate buffers, map kernel memory (for MMAP), enqueue initial capture buffers, and enable `VIDIOC_STREAMON`.
3. Consume frames via `dev.Frames(ctx)` or `dev.GetFrames()` (or emit frames via `dev.WriteFrame` / `dev.WriteFramePlanes`).
4. Call `dev.Stop()` when finished so buffers are unmapped and released cleanly.

`Device` methods are **not** safe for concurrent calls from multiple goroutines without external synchronization. In the channel-based capture mode (`GetFrames()`), you may read from the returned `<-chan *Frame` safely from consumer goroutines while managing `Start`/`Stop`/`Close` on the owning control goroutine.
