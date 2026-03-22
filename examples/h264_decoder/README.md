# H.264 Hardware Decoder Example

This example demonstrates hardware-accelerated H.264 video decoding using the V4L2 stateful codec API.

## Overview

The example:
1. Reads H.264 encoded data from a file
2. Decodes using a hardware decoder (V4L2 M2M codec)
3. Saves decoded frames as raw NV12 files

## Requirements

### Hardware

You need a device with a hardware video decoder. Common options include:

| Platform | Decoder Device | Driver |
|----------|---------------|--------|
| Raspberry Pi 4/5 | `/dev/video10` | bcm2835-codec |
| Intel (Quick Sync) | Varies | stateless/stateful |
| Rockchip RK3399 | `/dev/video0` | rkvdec |
| Allwinner | Varies | cedrus |

### Software

- Linux kernel with V4L2 codec support
- An H.264 input file

To check for available decoder devices:

```bash
v4l2-ctl --list-devices
```

To verify decoder capabilities:

```bash
v4l2-ctl -d /dev/video10 --all
```

## Building

```bash
cd examples/h264_decoder
go build
```

## Usage

```bash
# Basic usage
./h264_decoder -d /dev/video10 -i input.h264 -o frames/

# Decode only first 100 frames
./h264_decoder -d /dev/video10 -i input.h264 -o frames/ -n 100

# Help
./h264_decoder -help
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-d` | `/dev/video10` | Decoder device path |
| `-i` | `input.h264` | Input H.264 file |
| `-o` | `frames` | Output directory |
| `-n` | `0` | Max frames (0 = all) |
| `-s` | `65536` | Read chunk size |

## Output

Decoded frames are saved as raw NV12 files:
- `frames/frame_0000.raw`
- `frames/frame_0001.raw`
- etc.

To view a decoded frame:

```bash
# Using FFplay (specify resolution)
ffplay -f rawvideo -pix_fmt nv12 -s 1280x720 frames/frame_0000.raw

# Convert to PNG
ffmpeg -f rawvideo -pix_fmt nv12 -s 1280x720 \
    -i frames/frame_0000.raw frame_0000.png
```

To create a video from decoded frames:

```bash
ffmpeg -f rawvideo -pix_fmt nv12 -s 1280x720 -framerate 30 \
    -i 'frames/frame_%04d.raw' -c:v libx264 output.mp4
```

## Creating Test Input

If you don't have an H.264 file, create one from a video:

```bash
# Extract raw H.264 from MP4
ffmpeg -i video.mp4 -c:v copy -bsf:v h264_mp4toannexb output.h264

# Or encode a video to H.264
ffmpeg -i video.mp4 -c:v libx264 -an output.h264
```

## Dynamic Resolution Changes

The decoder automatically handles resolution changes in the stream. When the resolution changes:

1. The `GetResolutionChanges()` channel receives a notification
2. New frames will have the updated resolution
3. No manual intervention is required

```go
// Monitor resolution changes
go func() {
    for change := range decoder.GetResolutionChanges() {
        fmt.Printf("New resolution: %dx%d\n", change.Width, change.Height)
    }
}()
```

## Troubleshooting

### "Decoder device not found"

Check available devices:

```bash
v4l2-ctl --list-devices
```

Look for devices with "decoder" or "codec" in their name.

### "Device is not a decoder"

The device exists but doesn't support decoding. Check capabilities:

```bash
v4l2-ctl -d /dev/video10 --info
```

### "Failed to decode"

Common issues:
- Input file is not valid H.264 Annex B format
- Decoder doesn't support the H.264 profile/level
- Hardware doesn't support the resolution

Check supported formats:

```bash
v4l2-ctl -d /dev/video10 --list-formats
v4l2-ctl -d /dev/video10 --list-formats-out
```

## Code Overview

### Opening the Decoder

```go
decoder, err := device.OpenDecoder("/dev/video10", device.DecoderConfig{
    InputFormat: v4l2.PixFormat{
        PixelFormat: v4l2.PixelFmtH264,
    },
    OutputFormat: v4l2.PixFormat{
        PixelFormat: v4l2.PixelFmtNV12,
    },
})
```

### Decoding Loop

```go
// Start decoder
decoder.Start(ctx)
defer decoder.Drain()

// Feed encoded data
decoder.GetInput() <- h264Data

// Receive decoded frames
frame := <-decoder.GetOutput()
```

### Handling Resolution Changes

```go
select {
case frame := <-decoder.GetOutput():
    // Process frame
case change := <-decoder.GetResolutionChanges():
    // Handle resolution change
    fmt.Printf("New resolution: %dx%d\n", change.Width, change.Height)
}
```

### Seeking (Flush)

To seek in a stream, flush the decoder and feed new data:

```go
decoder.Flush()  // Clear decoder state
// Now feed data from new position
decoder.GetInput() <- newData
```

## See Also

- [H.264 Encoder Example](../h264_encoder/)
- [V4L2 Codec Documentation](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/dev-decoder.html)
