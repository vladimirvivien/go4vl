# V4L2 Sub-device Pad Pipeline (`subdev_pipeline`)

This example inspects and negotiates pad formats, media bus codes, frame sizes, and selection rectangles on a V4L2 sub-device (`/dev/v4l-subdev*`) using `v4l2.GetSubdevFormat`, `v4l2.SetSubdevFormat`, `v4l2.EnumSubdevMBusCodes`, `v4l2.EnumSubdevFrameSizes`, and `v4l2.GetSubdevSelection`. When `-video` is provided, it also captures frames from the paired `/dev/video*` node using `Device.Frames(ctx)`.

## Prerequisites

- A Linux host with a V4L2 sub-device (such as a CSI camera sensor at `/dev/v4l-subdev0`)
- Go 1.25+

## Build and Run

```bash
go build -o subdev_pipeline ./examples/subdev_pipeline
./subdev_pipeline -subdev /dev/v4l-subdev0 -pad 0 -video /dev/video0 -n 5
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-subdev` | `/dev/v4l-subdev0` | V4L2 sub-device path |
| `-pad` | `0` | Sub-device pad index |
| `-video` | `""` | Optional paired `/dev/video*` node to stream frames after sub-device inspection |
| `-n` | `5` | Number of frames to capture when `-video` is set |
