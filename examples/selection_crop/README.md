# Selection & Crop API Example

Demonstrates querying and applying V4L2 selection rectangles (`VIDIOC_G_SELECTION` / `VIDIOC_S_SELECTION`) for crop, compose, bounds, and native sensor dimensions.

## Overview

The V4L2 Selection API supersedes the legacy `VIDIOC_G_CROP`/`VIDIOC_S_CROP` ioctls and is used by capture devices, ISPs, scalers, and hardware codecs to negotiate active crop and compose rectangles.

## Usage

```bash
# Query all supported selection targets on /dev/video0
go run . -d /dev/video0

# Apply a 640x480 crop rectangle at offset (0, 0)
go run . -d /dev/video0 -set -left 0 -top 0 -width 640 -height 480
```

## Flags

- `-d` — V4L2 device path (default: `/dev/video0`)
- `-set` — Apply a crop rectangle via `VIDIOC_S_SELECTION` (default: `false`)
- `-left`, `-top`, `-width`, `-height` — Crop rectangle coordinates when `-set` is enabled
