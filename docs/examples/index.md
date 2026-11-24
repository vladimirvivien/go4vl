# Examples

All 40+ runnable programs in [`examples/`](https://github.com/vladimirvivien/go4vl/tree/main/examples) build cleanly with `CGO_ENABLED=0` across `amd64`, `arm64`, `arm`, and `riscv64`. Build any example from the repository root with standard `go build` or `go run`:

```bash
CGO_ENABLED=0 go build -o iter_capture ./examples/iter_capture
./iter_capture -d /dev/video0
```

## Recommended capture & streaming

| Example | Description | Link |
|---|---|---|
| `iter_capture` | Zero-copy streaming capture using `Device.Frames(ctx)` (`iter.Seq2`) and safe `frame.Clone()` | [iter_capture](https://github.com/vladimirvivien/go4vl/tree/main/examples/iter_capture) |
| `iter_snapshot` | Capture a single frame to disk using `Device.Frames(ctx)` | [iter_snapshot](https://github.com/vladimirvivien/go4vl/tree/main/examples/iter_snapshot) |
| `mplane_device` | High-level multi-planar (`MPLANE`) capture via `device.Open` and `frame.Planes` | [mplane_device](https://github.com/vladimirvivien/go4vl/tree/main/examples/mplane_device) |
| `video_output` | Stream synthetic `YUYV` frames to a V4L2 output or `v4l2loopback` device | [video_output](https://github.com/vladimirvivien/go4vl/tree/main/examples/video_output) |
| `capture_frames` | Pooled channel capture (`GetFrames()`) with frame metadata and `frame.Release()` | [capture_frames](https://github.com/vladimirvivien/go4vl/tree/main/examples/capture_frames) |
| `capture0` | Capture multiple images to files | [capture0](https://github.com/vladimirvivien/go4vl/tree/main/examples/capture0) |
| `capture1` | Multi-frame capture with explicit `PixFormat` configuration | [capture1](https://github.com/vladimirvivien/go4vl/tree/main/examples/capture1) |
| `snapshot` | Single-frame channel snapshot | [snapshot](https://github.com/vladimirvivien/go4vl/tree/main/examples/snapshot) |
| `userptr_capture` | Application-allocated userspace buffers (`IOTypeUserPtr`) | [userptr_capture](https://github.com/vladimirvivien/go4vl/tree/main/examples/userptr_capture) |
| `dmabuf_export` | Zero-copy `VIDIOC_EXPBUF` DMA-BUF file descriptor export | [dmabuf_export](https://github.com/vladimirvivien/go4vl/tree/main/examples/dmabuf_export) |
| `readwrite_capture` | Direct `Read` / `Write` I/O (`IOMethodReadWrite`) | [readwrite_capture](https://github.com/vladimirvivien/go4vl/tree/main/examples/readwrite_capture) |

## Formats, metadata & device inspection

| Example | Description | Link |
|---|---|---|
| `device_info` | Enumerate `/dev/video*` devices and print driver capabilities | [device_info](https://github.com/vladimirvivien/go4vl/tree/main/examples/device_info) |
| `format` | Query and update active pixel formats | [format](https://github.com/vladimirvivien/go4vl/tree/main/examples/format) |
| `pixel_formats` | Enumerate supported FourCC formats, sizes, and frame intervals | [pixel_formats](https://github.com/vladimirvivien/go4vl/tree/main/examples/pixel_formats) |
| `colorspace` | Query and interpret V4L2 colorspace, transfer function, and quantization | [colorspace](https://github.com/vladimirvivien/go4vl/tree/main/examples/colorspace) |
| `mplane_capture` | Low-level `v4l2` multi-planar negotiation, buffers, and plane DMA-BUF export | [mplane_capture](https://github.com/vladimirvivien/go4vl/tree/main/examples/mplane_capture) |
| `meta_capture` | Discover and configure `BufTypeMetaCapture` / `MetaOutput` nodes (`MetaFormat`) | [meta_capture](https://github.com/vladimirvivien/go4vl/tree/main/examples/meta_capture) |
| `selection_crop` | Query and apply crop/compose selection rectangles (`VIDIOC_G/S_SELECTION`) | [selection_crop](https://github.com/vladimirvivien/go4vl/tree/main/examples/selection_crop) |
| `dv_timings` | Query Digital Video (DV) timings on HDMI/DisplayPort/SDI nodes | [dv_timings](https://github.com/vladimirvivien/go4vl/tree/main/examples/dv_timings) |
| `video_standards` | Query analog/SD video standards (PAL, NTSC, SECAM) | [video_standards](https://github.com/vladimirvivien/go4vl/tree/main/examples/video_standards) |

## Media Controller, sub-devices & hardware codecs

| Example | Description | Link |
|---|---|---|
| `media_topology` | Traverse `/dev/media*` entities, interfaces, pads, and links (`MEDIA_IOC_G_TOPOLOGY`) | [media_topology](https://github.com/vladimirvivien/go4vl/tree/main/examples/media_topology) |
| `subdev_pipeline` | Inspect `/dev/v4l-subdev*` pad formats, media-bus codes, and selections | [subdev_pipeline](https://github.com/vladimirvivien/go4vl/tree/main/examples/subdev_pipeline) |
| `m2m_pipeline` | Dual-queue `v4l2.M2MSession` (`OUTPUT` + `CAPTURE`) multiplexed with `v4l2.Poller` | [m2m_pipeline](https://github.com/vladimirvivien/go4vl/tree/main/examples/m2m_pipeline) |
| `media_request` | Allocate, reinitialize, and bind `MediaRequest` FDs for stateless codecs | [media_request](https://github.com/vladimirvivien/go4vl/tree/main/examples/media_request) |
| `stateless_hevc` | Validate stateless HEVC `SPS`/`PPS`/`DecodeParams` controls on `MediaRequest` | [stateless_hevc](https://github.com/vladimirvivien/go4vl/tree/main/examples/stateless_hevc) |
| `h264_encoder` | Stateful hardware H.264 encoder pipeline (`device.OpenEncoder`) | [h264_encoder](https://github.com/vladimirvivien/go4vl/tree/main/examples/h264_encoder) |
| `h264_decoder` | Stateful hardware H.264 decoder pipeline (`device.OpenDecoder`) | [h264_decoder](https://github.com/vladimirvivien/go4vl/tree/main/examples/h264_decoder) |

## Controls, inputs/outputs, tuners & webcams

| Example | Description | Link |
|---|---|---|
| `user_ctrl` | Read and update user controls (brightness, contrast, saturation) | [user_ctrl](https://github.com/vladimirvivien/go4vl/tree/main/examples/user_ctrl) |
| `ext_ctrls` | Atomic multi-control read/write via `ExtControls` | [ext_ctrls](https://github.com/vladimirvivien/go4vl/tree/main/examples/ext_ctrls) |
| `control_reference` | Browse generated control constants and metadata | [control_reference](https://github.com/vladimirvivien/go4vl/tree/main/examples/control_reference) |
| `video_inputs` / `video_outputs` | Enumerate and select physical video inputs and outputs | [video_inputs](https://github.com/vladimirvivien/go4vl/tree/main/examples/video_inputs) |
| `audio_inputs` / `audio_outputs` | Enumerate and configure V4L2 audio inputs and outputs | [audio_inputs](https://github.com/vladimirvivien/go4vl/tree/main/examples/audio_inputs) |
| `tuner` / `modulator` | Control radio/TV tuners and RF modulators | [tuner](https://github.com/vladimirvivien/go4vl/tree/main/examples/tuner) |
| `simplecam` | Minimal HTTP MJPEG webcam streamer (with multi-stage `Dockerfile`) | [simplecam](https://github.com/vladimirvivien/go4vl/tree/main/examples/simplecam) |
| `webcam` | Full web UI with image controls, format selection, and OpenCV integration | [webcam](https://github.com/vladimirvivien/go4vl/tree/main/examples/webcam) |
| `fileserv` | Serve captured frames over HTTP | [fileserv](https://github.com/vladimirvivien/go4vl/tree/main/examples/fileserv) |
