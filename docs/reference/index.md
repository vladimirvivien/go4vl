# Package & V4L2 UAPI Reference Map

Full generated API documentation for all exported types, functions, methods, and constants lives on `pkg.go.dev`:

* **[`github.com/vladimirvivien/go4vl/device`](https://pkg.go.dev/github.com/vladimirvivien/go4vl/device)** — High-level `Device`, functional `Option` helpers, `Frames(ctx)`, `GetFrames()`, `FramePool`, `WriteFrame`, `Encoder`, and `Decoder`.
* **[`github.com/vladimirvivien/go4vl/v4l2`](https://pkg.go.dev/github.com/vladimirvivien/go4vl/v4l2)** — Low-level Linux V4L2, Media Controller, Sub-device, Media Request, and stateless codec UAPI wrappers.

## Kernel subsystem to Go symbol mapping

If you are translating Linux C code or `v4l2-ctl` commands into Go, locate the corresponding Go entry points below:

| Linux Subsystem / C `ioctl` | `v4l2` Package (Low-Level) | `device` Package (High-Level) |
|---|---|---|
| `open` / `close` / sysfs discovery | `OpenDevice`, `CloseDevice` | `GetAllDevicePaths`, `IsDevice`, `Open`, `Device.Close` |
| `VIDIOC_QUERYCAP` | `GetCapability`, `Capability` | `Device.Capability` |
| `VIDIOC_ENUM_FMT`, `FRAMESIZES`, `FRAMEINTERVALS` | `GetAllFormatDescriptions`, `GetFormatFrameSizes`, `GetFormatFrameInterval` | `Device.GetFormatDescriptions`, `GetFormatDescription` |
| `VIDIOC_{G,S}_FMT` (single-planar) | `GetPixFormat`, `SetPixFormat`, `PixFormat` | `WithPixFormat`, `GetPixFormat`, `SetPixFormat` |
| `VIDIOC_{G,S,TRY}_FMT` (`MPLANE`) | `GetPixFormatMPlane`, `SetPixFormatMPlane`, `TryPixFormatMPlane`, `PixFormatMPlane` | `WithPixFormatMPlane`, `GetPixFormatMPlane`, `SetPixFormatMPlane` |
| `VIDIOC_{G,S,TRY}_FMT` (`META`) | `GetMetaFormat`, `SetMetaFormat`, `TryMetaFormat`, `MetaFormat` | — |
| `VIDIOC_{G,S}_PARM` (FPS) | `GetStreamParam`, `SetStreamParam` | `WithFPS`, `GetFrameRate`, `SetFrameRate` |
| `VIDIOC_{G,S}_SELECTION` / `CROPCAP` | `GetSelection`, `SetSelection`, `GetCropRect` | `GetCropCapability`, `SetCropRect` |
| `VIDIOC_{G,S}_EDID` / `SUBDEV_{G,S}_EDID` | `GetEDID`, `SetEDID`, `GetSubdevEDID`, `SetSubdevEDID`, `EDID` | — |
| `VIDIOC_REQBUFS`, `CREATE_BUFS`, `REMOVE_BUFS`, `QUERYBUF`, `PREPARE_BUF` | `InitBuffers`, `CreateBuffers`, `RemoveBuffers`, `QueryBuffer`, `PrepareBuffer`, `MapMemoryBuffers` | Handled inside `Device.Start` / `Stop` |
| `VIDIOC_QBUF`, `DQBUF`, `STREAMON`, `STREAMOFF`, `EXPBUF` | `QueueBufferWith`, `DequeueBuffer`, `StreamOnType`, `StreamOffType`, `ExportBuffer`, `ExportBufferPlane` | `Frames`, `GetFrames`, `WriteFrame`, `WriteFramePlanes`, `ExportBuffer` |
| `poll` / `epoll` multiplexing | `Poll`, `Poller`, `PollEvent` | Used internally by `Frames`, `WriteFrame` |
| `VIDIOC_{QUERYCTRL,G_CTRL,S_CTRL}` | `GetControl`, `GetControlValue`, `SetControlValue`, `QueryAllControls` | `GetControl`, `SetControlValue`, `QueryAllControls` |
| `VIDIOC_{G,S,TRY}_EXT_CTRLS` + stateless codecs | `ExtControls`, `GetExtControls`, `SetExtControls`, `TryExtControls`, `AddH264*`, `AddHEVC*`, `AddVP8*`, `AddVP9*`, `AddAV1*`, `AddMPEG2*`, `AddFWHT*` | `GetExtControls`, `SetExtControls`, `TryExtControls` |
| `VIDIOC_SUBSCRIBE_EVENT` / `DQEVENT` | `SubscribeEvent`, `UnsubscribeEvent`, `DequeueEvent` | `Device.SubscribeEvent`, `UnsubscribeEvent`, `DequeueEvent` |
| Dual-queue M2M & stateful codecs | `OpenM2M`, `NewM2M`, `M2MSession`, `SendEncoderCmd`, `SendDecoderCmd` | `OpenEncoder`, `Encoder`, `OpenDecoder`, `Decoder` |
| `/dev/media*` graph & requests | `GetMediaDeviceInfo`, `GetMediaTopology`, `SetupMediaLink`, `ResolveDevnodePath`, `AllocMediaRequest`, `MediaRequest` | `Device.GetMediaInfo` |
| `/dev/v4l-subdev*` pads & routing | `GetSubdevFormat`, `SetSubdevFormat`, `EnumSubdevMBusCodes`, `EnumSubdevFrameSizes`, `GetSubdevSelection`, `GetSubdevCapability`, `SetSubdevClientCapability`, `GetSubdevRouting`, `SetSubdevRouting` | — |
