# Video Outputs Example

This example demonstrates how to enumerate and select video outputs on V4L2 devices that support video output.

## Overview

Video output devices are less common than capture devices. They're typically found in:
- Hardware video encoders
- Memory-to-memory (M2M) codec devices
- Virtual output devices (v4l2loopback)
- Video injection/playback devices

This example shows how to:
- List all available video outputs
- Get information about each output (name, type, capabilities)
- Query the current output status
- Switch between different outputs

## Building

```bash
go build
```

## Usage

### List all outputs on the default device

```bash
./video_outputs
```

### List outputs on a specific device

```bash
./video_outputs -d /dev/video10
```

### Select a specific output

```bash
./video_outputs -s 1
```

This will switch to output index 1 and verify the change.

## Example Output

```
Device: /dev/video10
Driver: v4l2 loopback
Card: Virtual Output

Current output index: 0

Available video outputs (1):
================================================================================
[0] loopback out ** ACTIVE **
    Type:         Analog
    Audioset:     0x00000000
    Modulator:    0
    Standards:    0x0000000000000000
    Capabilities: 0x00000000

Current Output Status:
================================================================================
Status: ok
  ✓ Output OK

Tip: Use -s <index> to select a different output
```

## Hardware Support

### Devices with Video Output

This example requires devices that support V4L2_CAP_VIDEO_OUTPUT capability:

- **v4l2loopback** - Virtual output devices (most common for testing)
- **Hardware encoders** - Some encoding devices support output
- **M2M devices** - Memory-to-memory codec devices
- **Specialized hardware** - Professional video equipment

### Testing with v4l2loopback

The easiest way to test video outputs is with v4l2loopback:

```bash
# Load v4l2loopback in output mode
sudo modprobe v4l2loopback video_nr=10 card_label="Virtual Output"

# Verify the device supports output
v4l2-ctl -d /dev/video10 --all | grep "Video Output"

# Run the example
./video_outputs -d /dev/video10
```

### Why Are Output Devices Rare?

Most V4L2 usage is for video *capture* (webcams, capture cards). Output devices are specialized for:
- Injecting video into the kernel
- Hardware-accelerated encoding
- Video effects and processing pipelines

## Command-Line Flags

- `-d <device>` - Specify the device path (default: `/dev/video0`)
- `-s <index>` - Select output by index (default: no selection, just list)

## Related Examples

- **device_info** - General device information including capabilities
- **video_inputs** - Video input enumeration (the capture equivalent)

## API Reference

This example demonstrates the following go4vl APIs:

### Device-Level Methods
- `dev.GetVideoOutputIndex()` - Get current output index
- `dev.SetVideoOutputIndex(index)` - Select an output
- `dev.GetVideoOutputDescriptions()` - Get all available outputs
- `dev.GetVideoOutputInfo(index)` - Get info for specific output
- `dev.GetVideoOutputStatus()` - Query current output status

### Low-Level v4l2 Methods
- `v4l2.GetCurrentVideoOutputIndex(fd)` - VIDIOC_G_OUTPUT
- `v4l2.SetVideoOutputIndex(fd, index)` - VIDIOC_S_OUTPUT
- `v4l2.GetVideoOutputInfo(fd, index)` - VIDIOC_ENUMOUTPUT
- `v4l2.GetAllVideoOutputInfo(fd)` - Enumerate all outputs
- `v4l2.QueryOutputStatus(fd)` - Query output status

### Types and Constants
- `v4l2.OutputInfo` - Output information structure
- `v4l2.OutputType` - Output type (Modulator, Analog, AnalogVGAOverlay)
- `v4l2.OutputStatus` - Output status (typically always OK)
- `v4l2.OutputStatuses` - Status value to string map

## Differences from Input API

The output API is very similar to the input API, but with some differences:

**Similarities:**
- Enumeration works the same way
- Selection API is identical
- Information structures are parallel

**Differences:**
- Output status is simpler (V4L2 doesn't define detailed status flags)
- Outputs have a `modulator` field instead of `tuner`
- Output types are different (Modulator, Analog, etc.)

## Troubleshooting

### Error: "Device does not support video output"

Most devices are capture-only. You need:
1. A device with V4L2_CAP_VIDEO_OUTPUT capability
2. v4l2loopback for testing
3. Specialized hardware

Check capabilities with:
```bash
v4l2-ctl -d /dev/videoX --all | grep -i "output"
```

### Error: "Failed to get current output"

The device may not support output enumeration or has only one output.

### How to Set Up v4l2loopback for Output

```bash
# Install v4l2loopback
sudo apt install v4l2loopback-dkms  # Debian/Ubuntu

# Load module
sudo modprobe v4l2loopback video_nr=10 card_label="Virtual Output"

# Verify
v4l2-ctl -d /dev/video10 --all
```

## Use Cases

Video output devices are used for:
- **Video injection** - Feeding video into applications that expect a camera
- **Testing** - Creating virtual video sources for development
- **Processing pipelines** - Video effects and transformation chains
- **Hardware encoding** - Offloading encoding to specialized hardware

## See Also

- [V4L2 VIDIOC_ENUMOUTPUT documentation](https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumoutput.html)
- [V4L2 VIDIOC_G_OUTPUT documentation](https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-output.html)
- [v4l2loopback documentation](https://github.com/umlaeute/v4l2loopback)
