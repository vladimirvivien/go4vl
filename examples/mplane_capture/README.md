# Multi-Planar Capture Example

Demonstrates querying and negotiating V4L2 multi-planar (`VIDEO_CAPTURE_MPLANE`) formats, requesting multi-planar buffers, inspecting `v4l2_plane` descriptors, mapping plane memory, and optionally exporting a plane as a DMA-BUF file descriptor.

## Overview

Modern Linux SoC ISPs and hardware codecs frequently expose multi-planar capture and output queues (`V4L2_CAP_VIDEO_CAPTURE_MPLANE`, `V4L2_CAP_VIDEO_M2M_MPLANE`) where image planes (such as luma and chroma in `NV12M` or `YUV420M`) reside in separate plane descriptors.

This example uses the low-level `v4l2` APIs directly:
- `v4l2.GetAllFormatDescriptionsByType(fd, v4l2.BufTypeVideoCaptureMPlane)`
- `v4l2.GetPixFormatMPlane` / `v4l2.TryPixFormatMPlane`
- `v4l2.RequestBuffersForType`
- `v4l2.QueryBuffer` + `v4l2.MapMemoryBuffer` / `v4l2.UnmapMemoryBuffer`
- `v4l2.ExportBufferPlane` (with `-export`)

## Usage

```bash
# Inspect and allocate buffers on a multi-planar capture or M2M device
go run . -d /dev/video11

# Request 2 buffers and export plane 0 as a DMA-BUF fd
go run . -d /dev/video11 -b 2 -export
```

## Flags

- `-d` — V4L2 multi-planar device path (default: `/dev/video0`)
- `-b` — Number of buffers to allocate (default: `4`)
- `-export` — Export plane 0 of each buffer as a DMA-BUF fd via `VIDIOC_EXPBUF` (default: `false`)
