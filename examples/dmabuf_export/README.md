# DMA-BUF Export Example

Demonstrates exporting V4L2 MMAP buffers as DMA-BUF file descriptors for zero-copy sharing with other subsystems (GPU, DRM, other V4L2 devices).

## DMA-BUF Overview

DMA-BUF is a Linux kernel mechanism for sharing buffers between devices via file descriptors:

- **Export**: Share V4L2 capture buffers with GPU/DRM (`VIDIOC_EXPBUF`)
- **Import**: Use external DMA-BUF fds for V4L2 capture (`V4L2_MEMORY_DMABUF`)

This example demonstrates the **export** path.

## Usage

```bash
# Export 4 MMAP buffers as DMA-BUF fds and capture 5 frames
go run . -d /dev/video0

# With custom buffer count
go run . -d /dev/video0 -b 8
```

## Flags

- `-d` — Device path (default: `/dev/video0`)
- `-b` — Number of buffers (default: 4)

## Output

```
Buffer 0 exported as DMA-BUF fd=5
Buffer 1 exported as DMA-BUF fd=6
Buffer 2 exported as DMA-BUF fd=7
Buffer 3 exported as DMA-BUF fd=8
Frame 0: seq=0, 12543 bytes
Frame 1: seq=1, 12601 bytes
```

## Import Mode

For DMA-BUF import (receiving fds from another subsystem):

```go
dev, _ := device.Open("/dev/video0",
    device.WithIOType(v4l2.IOTypeDMABuf),
    device.WithBufferSize(4),
)
dev.AddDMABufferFDs(gpuFD0, gpuFD1, gpuFD2, gpuFD3)
dev.Start(ctx)
for frame := range dev.GetFrames() { ... }
```
