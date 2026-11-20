# Video Output Streaming (`video_output`)

This example demonstrates writing synthetic YUYV video frames to a V4L2 video output device (such as a `v4l2loopback` virtual camera) using `device.WithVideoOutputEnabled()` and `Device.WriteFrame(ctx, data)`.

## Prerequisites

- A Linux host with a V4L2 video output device (for example, `sudo modprobe v4l2loopback video_nr=10 card_label="go4vl-out"`)
- Go 1.25+

## Build and Run

```bash
go build -o video_output ./examples/video_output
./video_output -d /dev/video10 -w 640 -h 480 -n 60 -fps 30
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-d` | `/dev/video10` | V4L2 video output device path |
| `-w` | `640` | Frame width in pixels |
| `-h` | `480` | Frame height in pixels |
| `-n` | `30` | Number of frames to write |
| `-fps` | `30` | Target output frame rate |
