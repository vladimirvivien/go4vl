# Media Controller Graph Inspection (`media_topology`)

This example inspects a Linux Media Controller device (`/dev/media*`) using `v4l2.GetMediaDeviceInfo`, `v4l2.GetMediaTopology` (`MEDIA_IOC_G_TOPOLOGY`), and `v4l2.ResolveDevnodePath` to print all entities, attached `/dev/video*` and `/dev/v4l-subdev*` interfaces, pads, and data links.

## Prerequisites

- A Linux host exposing one or more `/dev/media*` devices (such as a CSI camera, ISP, or hardware codec)
- Go 1.25+

## Build and Run

```bash
go build -o media_topology ./examples/media_topology
./media_topology -d /dev/media0
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-d` | `/dev/media0` | Media Controller device path |
