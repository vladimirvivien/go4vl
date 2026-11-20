# Zero-Copy Iterator Capture (`iter_capture`)

This example demonstrates zero-copy video capture using `Device.Frames(ctx)` (`iter.Seq2[*device.Frame, error]`). Each yielded `*device.Frame` slices directly into the driver's memory-mapped buffer without per-frame heap copies or background capture goroutines, and automatically re-queues the buffer (`VIDIOC_QBUF`) after each loop iteration finishes—including on `break`, `return`, or `panic`.

## Prerequisites

- A Linux host with a V4L2 video capture device (such as a webcam, CSI camera, or `v4l2loopback` device)
- Go 1.25+

## Build and Run

```bash
go build -o iter_capture ./examples/iter_capture
./iter_capture -d /dev/video0 -n 30 -b 4
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-d` | `/dev/video0` | V4L2 video capture device path |
| `-n` | `30` | Number of frames to capture |
| `-b` | `4` | Number of `MMAP` buffers to request |
| `-save-last` | `""` | Optional file path to save a `frame.Clone()` copy of the last frame after the loop |

## Buffer Lifetime

`frame.Data` is valid only for the duration of the current loop iteration. If you need to retain frame data across iterations or send it to another goroutine, call `frame.Clone()` inside the loop body.
