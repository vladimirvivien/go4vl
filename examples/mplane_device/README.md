# mplane_device Example

Captures frames from a multi-planar (`V4L2_BUF_TYPE_VIDEO_CAPTURE_MPLANE`) V4L2 device using the high-level `device` package (`WithBufferType`, `WithPixFormatMPlane`, `GetPixFormatMPlane`, and `Frames(ctx)` with `frame.Planes`).

## Usage

Inspect a multi-planar device's active format and per-plane geometry without starting capture:

```bash
go run ./examples/mplane_device -d /dev/video10 -frames 0
```

Capture 5 frames from a multi-planar capture device:

```bash
go run ./examples/mplane_device -d /dev/video10 -frames 5
```

Request explicit dimensions:

```bash
go run ./examples/mplane_device -d /dev/video10 -w 640 -h 480 -frames 5
```

## Example Output

```
Device: /dev/video10 (bcm2835-codec-decode, driver=bcm2835-codec)
Format: 640x480 pixelformat=0x32315559 planes=1
  plane[0]: bytesperline=640 sizeimage=460800
```

(Host-specific device paths and plane layouts differ by driver and sensor.)
