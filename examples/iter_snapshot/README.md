# Single-Frame Snapshot with `Device.Frames` (`iter_snapshot`)

This example captures a single video frame and writes it to disk using `Device.Frames(ctx)` (`iter.Seq2[*device.Frame, error]`). Breaking out of the `for ... range` loop after the first frame automatically returns the underlying `MMAP` buffer to the driver with no manual `Release()` call.

## Prerequisites

- A Linux host with a V4L2 capture device (such as `/dev/video0`)
- Go 1.25+

## Build and Run

```bash
go build -o iter_snapshot ./examples/iter_snapshot
./iter_snapshot -d /dev/video0 -o pic.jpg
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-d` | `/dev/video0` | V4L2 video capture device path |
| `-o` | `snapshot.raw` | Destination file for the captured frame bytes |
