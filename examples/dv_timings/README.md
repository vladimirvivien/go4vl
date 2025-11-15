# DV Timings Example

This example demonstrates how to work with Digital Video (DV) timings in V4L2 using go4vl. DV timings are used for digital video interfaces like HDMI, DisplayPort, DVI, and SDI.

## What are DV Timings?

DV timings define the video signal characteristics for digital video interfaces:
- **Resolution**: Image dimensions (width x height)
- **Pixel Clock**: Frequency at which pixels are transmitted
- **Frame Rate**: Calculated from pixel clock and total timing parameters
- **Blanking Periods**: Horizontal and vertical sync, front porch, and back porch
- **Sync Polarities**: H-sync and V-sync polarity (positive or negative)
- **Standards**: CEA-861 (HDMI/DVI), DMT (VESA), CVT, GTF
- **Format**: Progressive or interlaced

These timings follow industry standards like BT.656/1120 and are essential for:
- HDMI capture cards
- DisplayPort receivers
- SDI video interfaces
- Professional video equipment

## Features Demonstrated

This example shows how to:

1. **Check DV Timings Support**: Verify if a device supports DV timings
2. **Query Capabilities**: Get supported resolution ranges, pixel clock ranges, and standards
3. **Get Current Timings**: Read the currently configured timings
4. **Auto-Detect Timings**: Automatically detect timings from input signal
5. **Enumerate Timings**: List all supported timing configurations
6. **Display Timing Details**: Show comprehensive timing information including blanking periods and standards

## Prerequisites

- A V4L2 device that supports DV timings (e.g., HDMI capture card, SDI interface)
- The device file (e.g., `/dev/video0`)

## Building

```bash
cd examples/dv_timings
go build
```

## Usage

Basic usage with default device:
```bash
./dv_timings
```

Specify a different device:
```bash
./dv_timings -d /dev/video1
```

Specify a pad number for multi-pad devices:
```bash
./dv_timings -d /dev/video0 -p 1
```

## Command-Line Flags

- `-d <device>`: Device path (default: `/dev/video0`)
- `-p <pad>`: Pad number for multi-pad devices (default: `0`)

## Example Output

```
DV Timings Example - Device: /dev/video0
=====================================

DV Timing Capabilities
----------------------
  Type: 0
  Pad: 0
  Resolution range: 640x480 to 3840x2160
  Pixel clock range: 25000000 - 600000000 Hz (25.0 - 600.0 MHz)

  Supported formats:
    - Interlaced
    - Progressive
    - Reduced blanking

  Supported standards:
    - CEA-861 (HDMI/DVI)
    - DMT (VESA Display Monitor Timings)

Current DV Timings
------------------
  Type: 0 (BT.656/1120)
  Resolution: 1920x1080
  Pixel Clock: 148500000 Hz (148.50 MHz)
  Frame Rate: 60.00 Hz
  Format: Progressive

  Horizontal Blanking:
    Front Porch: 88 pixels
    Sync: 44 pixels
    Back Porch: 148 pixels
    Total: 2200 pixels

  Vertical Blanking:
    Front Porch: 4 lines
    Sync: 5 lines
    Back Porch: 36 lines
    Total: 1125 lines

  Sync Polarities:
    H-Sync: Positive
    V-Sync: Positive

  Standards: 0x1
    - CEA-861 (HDMI/DVI)

  Flags: 0x8
    - Has picture aspect ratio

Auto-Detected Timings
---------------------
  Successfully detected timings from input signal:
  Type: 0 (BT.656/1120)
  Resolution: 1920x1080
  Pixel Clock: 148500000 Hz (148.50 MHz)
  Frame Rate: 60.00 Hz
  ...

Enumerated Supported Timings
----------------------------
  Found 25 supported timing(s):

  [0] 720x480 @ 59.94 Hz - [CEA-861]
  [1] 720x576 @ 50.00 Hz - [CEA-861]
  [2] 1280x720 @ 50.00 Hz - [CEA-861]
  [3] 1280x720 @ 60.00 Hz - [CEA-861]
  [4] 1920x1080 @ 24.00 Hz - [CEA-861]
  [5] 1920x1080 @ 25.00 Hz - [CEA-861]
  [6] 1920x1080 @ 30.00 Hz - [CEA-861]
  [7] 1920x1080 @ 50.00 Hz - [CEA-861]
  [8] 1920x1080 @ 60.00 Hz - [CEA-861]
  [9] 3840x2160 @ 30.00 Hz - [CEA-861]

  ... and 15 more
```

## Common Use Cases

### HDMI Capture

When using an HDMI capture card, you can:
1. Connect an HDMI source (camera, computer, console)
2. Use auto-detection to identify the signal timings
3. Set those timings for capture
4. Enumerate all supported resolutions and frame rates

### Signal Monitoring

Monitor characteristics of incoming digital video signals:
- Resolution and frame rate
- Pixel clock frequency
- Sync polarities
- Video standards (CEA-861, DMT, etc.)
- Interlaced vs progressive format

### Format Validation

Check if a device supports specific timing configurations:
- Query supported resolution ranges
- Enumerate available standard timings
- Verify pixel clock limits

## Understanding the Output

### Resolution and Frame Rate
- **Resolution**: Active video area in pixels (e.g., 1920x1080)
- **Frame Rate**: Calculated from pixel clock and total timing parameters

### Pixel Clock
- Frequency at which pixels are transmitted
- Typically in MHz (e.g., 148.5 MHz for 1080p60)

### Blanking Periods
- **Front Porch**: Time after active video before sync
- **Sync**: Synchronization pulse
- **Back Porch**: Time after sync before active video
- For interlaced formats, additional IL (interlaced) timings are shown

### Standards
- **CEA-861**: Consumer Electronics Association standard (HDMI/DVI)
- **DMT**: VESA Display Monitor Timings
- **CVT**: Coordinated Video Timings
- **GTF**: Generalized Timing Formula

### VIC Codes
- **CEA-861 VIC**: Video Identification Code from CEA-861 standard
- **HDMI VIC**: HDMI-specific Video Identification Code

## API Functions Used

This example demonstrates the following go4vl functions:

```go
// Get current DV timings
timings, err := dev.GetDVTimings()

// Set DV timings
err = dev.SetDVTimings(timings)

// Auto-detect timings from signal
timings, err := dev.QueryDVTimings()

// Enumerate specific timing by index
enumTiming, err := dev.EnumerateDVTimings(index, pad)

// Get all supported timings
timings, err := dev.GetAllDVTimings(pad)

// Get DV timing capabilities
cap, err := dev.GetDVTimingsCap(pad)
```

## Troubleshooting

### "Device does not support DV timings"

This means your device doesn't support the DV timings API. This is normal for:
- Regular webcams
- USB cameras without HDMI input
- Devices using analog video standards

DV timings are specifically for digital video interfaces.

### "Could not detect timings: no signal"

This typically means:
- No input signal is connected
- The connected signal is outside the device's supported range
- The device doesn't support auto-detection

Try connecting a video source or checking cable connections.

### "Could not enumerate timings"

Some devices may not support timing enumeration even if they support DV timings. In this case, you'll need to know the desired timings in advance or use auto-detection.

## Related Examples

- **format** - Shows basic video format configuration
- **capture** - Demonstrates video capture with buffers
- **stream** - Shows continuous video streaming

## References

- [V4L2 DV Timings API](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/dv-timings.html)
- [BT.656/1120 Standards](https://en.wikipedia.org/wiki/ITU-R_BT.656)
- [CEA-861 Standard](https://en.wikipedia.org/wiki/CEA-861)
- [VESA DMT Standard](https://en.wikipedia.org/wiki/VESA_Discrete_Monitor_Timings)
