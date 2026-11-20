# Dual-Queue M2M Pipeline (`m2m_pipeline`)

This example demonstrates opening a V4L2 memory-to-memory (M2M) device—such as a hardware scaler, ISP, or video codec—with `v4l2.OpenM2M`, enumerating `OUTPUT` and `CAPTURE` formats, allocating mapped buffers for both queues, and registering the device file descriptor with `v4l2.Poller` (`epoll`) for `POLLIN | POLLOUT | POLLPRI` readiness multiplexing.

## Prerequisites

- A Linux host with a V4L2 M2M device (for example, `bcm2835-codec`, `bcm2835-isp`, or `vim2m`)
- Go 1.25+

## Build and Run

```bash
go build -o m2m_pipeline ./examples/m2m_pipeline
./m2m_pipeline -d /dev/video12 -b 2
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-d` | `/dev/video12` | V4L2 M2M device path |
| `-b` | `2` | Number of `MMAP` buffers per queue |
