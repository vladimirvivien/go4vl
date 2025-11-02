# Video Inputs Example

This example demonstrates how to enumerate and select video inputs on V4L2 devices that support multiple inputs.

## Overview

Many video capture devices support multiple inputs (e.g., composite, S-Video, HDMI). This example shows how to:
- List all available video inputs
- Get information about each input (name, type, status, capabilities)
- Query the current input status
- Switch between different inputs

## Building

```bash
go build
```

## Usage

### List all inputs on the default device

```bash
./video_inputs
```

### List inputs on a specific device

```bash
./video_inputs -d /dev/video1
```

### Select a specific input

```bash
./video_inputs -s 1
```

This will switch to input index 1 and verify the change.

## Example Output

```
Device: /dev/video0
Driver: uvcvideo
Card: HD Pro Webcam C920

Current input index: 0

Available video inputs (1):
================================================================================
[0] Camera 1 ** ACTIVE **
    Type:         Camera
    Status:       OK
    Audioset:     0x00000001
    Tuner:        0
    Standards:    0x0000000000000000
    Capabilities: 0x00000000

Current Input Status:
================================================================================
Status: OK
  ✓ Power OK
  ✓ Signal detected
  ✓ Color information present

Tip: Use -s <index> to select a different input
```

## Hardware Support

### Devices with Multiple Inputs

This example is most useful with:
- **TV tuner cards** - Often have multiple inputs (composite, S-Video, coaxial)
- **Multi-input capture cards** - HDMI capture cards with multiple ports
- **Professional video equipment** - Broadcast/production equipment

### Webcams and Simple Capture Devices

Most USB webcams have only a single input and may not support input enumeration. The example will still work but will show only one input.

## Testing Without Hardware

You can test this example using v4l2loopback:

```bash
# Load v4l2loopback
sudo modprobe v4l2loopback video_nr=10 card_label="Virtual Camera"

# Run the example
./video_inputs -d /dev/video10
```

Note: v4l2loopback devices typically have only one input.

## Command-Line Flags

- `-d <device>` - Specify the device path (default: `/dev/video0`)
- `-s <index>` - Select input by index (default: no selection, just list)

## Related Examples

- **device_info** - General device information including current input
- **format** - Video format capabilities

## API Reference

This example demonstrates the following go4vl APIs:

### Device-Level Methods
- `dev.GetVideoInputIndex()` - Get current input index
- `dev.SetVideoInputIndex(index)` - Select an input
- `dev.GetVideoInputDescriptions()` - Get all available inputs
- `dev.GetVideoInputInfo(index)` - Get info for specific input
- `dev.GetVideoInputStatus()` - Query current input status

### Low-Level v4l2 Methods
- `v4l2.GetCurrentVideoInputIndex(fd)` - VIDIOC_G_INPUT
- `v4l2.SetVideoInputIndex(fd, index)` - VIDIOC_S_INPUT
- `v4l2.GetVideoInputInfo(fd, index)` - VIDIOC_ENUMINPUT
- `v4l2.GetAllVideoInputInfo(fd)` - Enumerate all inputs
- `v4l2.QueryInputStatus(fd)` - Query input status

### Types and Constants
- `v4l2.InputInfo` - Input information structure
- `v4l2.InputType` - Input type (Tuner, Camera, Touch)
- `v4l2.InputStatus` - Input status flags
- `v4l2.InputStatuses` - Status value to string map

## Troubleshooting

### Error: "Device does not support video capture"

The device is not a video capture device. Check with:
```bash
v4l2-ctl -d /dev/videoX --all
```

### Error: "Failed to get current input"

The device may not support input enumeration or has only one input. This is normal for most webcams.

### Error: "Failed to select input"

Some devices don't allow input selection while streaming. Close any applications using the device first.

## See Also

- [V4L2 VIDIOC_ENUMINPUT documentation](https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enuminput.html)
- [V4L2 VIDIOC_G_INPUT documentation](https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-input.html)
