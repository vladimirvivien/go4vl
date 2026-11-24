# Examples

## Basic Capture Examples

* [snapshot](./snapshot/) - A simple example to capture a single frame and save it to a file.
* [iter_snapshot](./iter_snapshot/) - Capture a single frame to disk using the zero-copy `Device.Frames(ctx)` iterator.
* [capture0](./capture0) - Captures multiple images and saves them to files using `Device.Frames(ctx)`.
* [capture_frames](./capture_frames) - Capture frames with metadata and buffer pooling for high-performance streaming.
* [iter_capture](./iter_capture/) - **Recommended**: Zero-copy frame capture using `Device.Frames(ctx)` (`iter.Seq2`) with automatic buffer re-queueing and `frame.Clone()`.
* [mplane_device](./mplane_device/) - High-level multi-planar (`MPLANE`) capture using `device.Open`, `GetPixFormatMPlane`, and `Device.Frames(ctx)` (`frame.Planes`).
* [capture1](./capture1) - Shows how to capture multiple images using specified image format.
* [video_output](./video_output/) - Stream synthetic YUYV frames to a V4L2 video output or `v4l2loopback` device.
* [m2m_pipeline](./m2m_pipeline/) - Dual-queue M2M (`OUTPUT` + `CAPTURE`) session multiplexed with `v4l2.Poller` (`epoll`).

## Device Information and Control

* [device_info](./device_info) - Uses go4vl to query and print device and format information.
* [format](./format) - Shows how to query and apply device and format information.
* [mplane_capture](./mplane_capture/) - Query and negotiate multi-planar (`MPLANE`) formats, plane buffers, and DMA-BUF plane export.
* [meta_capture](./meta_capture/) - Discover and inspect V4L2 metadata capture/output nodes (`GetMetaFormat`, `TryMetaFormat`).
* [selection_crop](./selection_crop/) - Query and set V4L2 selection rectangles (`VIDIOC_G_SELECTION` / `VIDIOC_S_SELECTION`).
* [media_topology](./media_topology/) - Inspect `/dev/media*` entities, interfaces, pads, and links via `MEDIA_IOC_G_TOPOLOGY`.
* [subdev_pipeline](./subdev_pipeline/) - Query and negotiate `/dev/v4l-subdev*` pad formats, media-bus codes, and selection rectangles.
* [media_request](./media_request/) - Allocate, reinitialize, and bind Media Request FDs (`MEDIA_IOC_REQUEST_ALLOC`) for stateless hardware codecs.
* [stateless_hevc](./stateless_hevc/) - Validate stateless HEVC SPS/PPS/DecodeParams controls bound to a `MediaRequest` on hardware HEVC decoders.
* [video_inputs](./video_inputs/) - Enumerate and select video inputs (for devices with multiple inputs).
* [video_outputs](./video_outputs/) - Enumerate and select video outputs (for output devices and v4l2loopback).
* [audio_inputs](./audio_inputs/) - Enumerate and select audio inputs (for webcams with mics, TV tuner cards).
* [audio_outputs](./audio_outputs/) - Enumerate and select audio outputs (for video output devices).
* [tuner](./tuner/) - Enumerate tuners, get/set frequency, and query signal strength (for radio/TV tuners, SDR devices).
* [modulator](./modulator/) - Enumerate modulators and control transmission frequency (for RF modulators, transmitters).
* [dv_timings](./dv_timings/) - Query and detect Digital Video (DV) timings (for HDMI/DisplayPort/SDI capture cards).
* [ext_ctrls](./ext_ctrls/) Shows how to query and apply extended controls.
* [user_ctrl](./user_ctrl/) Shows how to query and apply user controls.

## Web Streaming Examples

* [simplecam](./simplecam/) A functional webcam program that streams video to web page.
* [webcam](./webcam) - Builds on simplecam and adds image control, format control, and face detection.

## Building the example code

`go4vl` builds with `CGO_ENABLED=0` on Linux (`amd64`, `arm64`, `arm`, and `riscv64`) with no C compiler or kernel headers required. For full details, see [docs/BUILD.md](../docs/BUILD.md) and the online [Examples documentation](https://vladimirvivien.github.io/go4vl/examples/).

### Native build

```bash
CGO_ENABLED=0 go build ./examples/iter_capture
```

### Cross-compile for Linux ARM64 (e.g., Raspberry Pi)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o iter_capture ./examples/iter_capture
```

### Cross-compile with Docker

See [examples/simplecam/Dockerfile](./simplecam/Dockerfile) for a multi-stage `CGO_ENABLED=0` container build.