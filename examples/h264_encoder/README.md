# H.264 Hardware Encoder Example

This example demonstrates hardware-accelerated H.264 video encoding using the V4L2 stateful codec API.

## Overview

The example:
1. Opens a camera device to capture raw video frames
2. Opens a hardware encoder device (V4L2 M2M codec)
3. Feeds raw frames (NV12) to the encoder
4. Writes the H.264 encoded output to a file

## Requirements

### Hardware

You need a device with a hardware video encoder. Common options include:

| Platform | Encoder Device | Driver |
|----------|---------------|--------|
| Raspberry Pi 4/5 | `/dev/video11` | bcm2835-codec |
| Intel (Quick Sync) | Varies | stateless/stateful |
| Rockchip RK3399 | `/dev/video1` | rkvenc |
| Allwinner | Varies | cedrus |

### Software

- Linux kernel with V4L2 codec support
- Camera device (USB webcam, CSI camera, etc.)

To check for available encoder devices:

```bash
v4l2-ctl --list-devices
```

To verify encoder capabilities:

```bash
v4l2-ctl -d /dev/video11 --all
```

## Building

```bash
cd examples/h264_encoder
go build
```

## Usage

```bash
# Basic usage (10 second recording)
./h264_encoder -e /dev/video11 -c /dev/video0 -o output.h264 -t 10

# Custom resolution and bitrate
./h264_encoder -e /dev/video11 -c /dev/video0 -o output.h264 \
    -w 1920 -h 1080 -f 30 -b 4000000 -t 30

# Help
./h264_encoder -help
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-e` | `/dev/video11` | Encoder device path |
| `-c` | `/dev/video0` | Camera device path |
| `-o` | `output.h264` | Output file path |
| `-t` | `10` | Recording duration (seconds) |
| `-w` | `1280` | Video width |
| `-h` | `720` | Video height |
| `-f` | `30` | Frames per second |
| `-b` | `2000000` | Bitrate (bits/second) |

## Output

The output is a raw H.264 Annex B bitstream. To play it:

```bash
# Using FFplay
ffplay output.h264

# Using VLC
vlc output.h264

# Convert to MP4
ffmpeg -i output.h264 -c copy output.mp4
```

## Troubleshooting

### "Encoder device not found"

The encoder device path varies by platform. Check available devices:

```bash
v4l2-ctl --list-devices
```

Look for devices with "encoder" or "codec" in their name.

### "Device is not an encoder"

The device exists but doesn't support encoding. It might be:
- A decoder-only device
- A camera device
- A different type of V4L2 device

Check device capabilities:

```bash
v4l2-ctl -d /dev/video11 --info
```

### "Failed to set format"

The encoder may not support the requested format. Try:
- Different resolution (720p instead of 1080p)
- Different pixel format (check what the encoder supports)

List supported formats:

```bash
v4l2-ctl -d /dev/video11 --list-formats
v4l2-ctl -d /dev/video11 --list-formats-out
```

## Code Overview

### Opening the Encoder

```go
encoder, err := device.OpenEncoder("/dev/video11", device.EncoderConfig{
    InputFormat: v4l2.PixFormat{
        Width:       1280,
        Height:      720,
        PixelFormat: v4l2.PixelFmtNV12,
    },
    OutputFormat: v4l2.PixFormat{
        PixelFormat: v4l2.PixelFmtH264,
    },
    Bitrate: 2000000,
})
```

### Encoding Loop

```go
// Start encoder
encoder.Start(ctx)
defer encoder.Drain() // Proper shutdown

// Feed frames
encoder.GetInput() <- rawFrame

// Receive encoded data
encodedData := <-encoder.GetOutput()
```

### Graceful Shutdown

Always call `Drain()` before stopping to ensure all frames are encoded:

```go
encoder.Drain()  // Wait for all frames to be encoded
encoder.Stop()   // Stop the encoder
encoder.Close()  // Release resources
```

## See Also

- [H.264 Decoder Example](../h264_decoder/)
- [V4L2 Codec Documentation](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/dev-encoder.html)
