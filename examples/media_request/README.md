# Media Request API & Stateless Codec Synchronization (`media_request`)

This example demonstrates allocating and reinitializing a Linux Media Request file descriptor (`v4l2.AllocMediaRequest`, `MediaRequest.Reinit`), discovering the paired stateless codec `/dev/video*` node via `v4l2.GetMediaTopology`, querying `v4l2.BufCapSupportsRequests`, and binding `v4l2.ExtControls` (`SetRequestFD`) and `v4l2.Buffer` (`BufFlagRequestFD`) to the request FD.

## Prerequisites

- A Linux host with a Media Request-capable device (for example, `rpi-hevc-dec`, `cedrus`, `hantro`, or `vicodec` exposed at `/dev/media0`)
- Go 1.25+

## Build and Run

```bash
go build -o media_request ./examples/media_request
./media_request -media /dev/media0
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-media` | `/dev/media0` | Media Controller device path supporting `MEDIA_IOC_REQUEST_ALLOC` |
| `-video` | `""` | Optional paired `/dev/video*` stateless codec node (auto-discovered from topology when omitted) |
