# V4L2 Control Reference Example

This example demonstrates how to work with V4L2 controls across all control classes using the go4vl library.

## Overview

V4L2 controls are user-configurable device parameters organized into classes. This example shows how to enumerate, query, get, and set controls across all 11 control classes defined in the V4L2 specification.

## Control Classes

V4L2 organizes controls into the following classes:

| Class | Description | Common Controls |
|-------|-------------|-----------------|
| **User** | Basic picture controls | Brightness, Contrast, Saturation, Hue |
| **Camera** | Camera controls | Exposure, Focus, Zoom, Pan, Tilt, Iris |
| **Flash** | Flash and LED controls | Flash Mode, Intensity, Timeout, Strobe |
| **JPEG** | JPEG compression | Compression Quality, Markers, Restart Interval |
| **Image Source** | Image sensor controls | Analog Gain, Blanking, Test Patterns |
| **Image Processing** | Image processing | Pixel Rate, Link Frequency, Deinterlacing |
| **Codec** | Video codec controls | Bitrate, GOP Size, Profiles, Levels |
| **Codec Stateless** | Stateless codec controls | H.264, VP8, VP9, MPEG2 parameters |
| **Digital Video** | Digital video timing | RX/TX mode, Timings, Aspect Ratio |
| **Detection** | Detection controls | Motion Detection, Face Detection |
| **Colorimetry** | Color space and HDR | Color Space, Transfer Function, HDR |

## Features

This example demonstrates:
- **Overview Mode**: Display all control classes and their availability
- **List All Mode**: Enumerate all available controls across all classes
- **Class-Specific Listing**: List controls in a specific class
- **Get Control**: Query current value and metadata for a control
- **Set Control**: Set control value with validation

## Building

```bash
cd examples/control_reference
go build
```

## Usage

### Display Control Classes Overview (Default)

```bash
./control_reference -d /dev/video0
```

Output shows:
- All 11 control classes
- Number of available controls in each class
- Brief description of each class
- Usage examples

### List All Available Controls

```bash
./control_reference -all
```

Shows all controls across all classes with:
- Control ID (hex and decimal)
- Control name
- Control type (Integer, Boolean, Menu, etc.)
- Current value (if readable)
- Range and step (for integer controls)
- Menu items (for menu controls)
- Flags (read-only, write-only, volatile, etc.)

### List Controls in Specific Class

```bash
# Camera controls
./control_reference -class camera

# User controls (brightness, contrast, etc.)
./control_reference -class user

# Codec controls
./control_reference -class codec

# Flash controls
./control_reference -class flash
```

Available class names:
- `user` - User controls
- `camera` - Camera controls
- `flash` - Flash controls
- `jpeg` - JPEG compression
- `image-source` - Image source controls
- `image-proc` - Image processing
- `codec` - Codec controls
- `codec-stateless` - Stateless codec controls
- `dv` - Digital video
- `detection` - Detection controls
- `colorimetry` - Colorimetry controls

### Get Control Value

```bash
# Get control by ID (hex format)
./control_reference -get 0x009a0901

# Get control by ID (decimal format)
./control_reference -get 10094849
```

Shows:
- Control name and type
- Current value
- Range (for integer controls)
- Menu selection (for menu controls)

### Set Control Value

```bash
# Set brightness to 128
./control_reference -set 0x00980900 -value 128

# Set exposure mode to manual (value 1)
./control_reference -set 0x009a0901 -value 1
```

The example will:
- Validate the control exists
- Check if it's read-only
- Validate range (for integer controls)
- Set the value
- Verify the new value

## Command Line Flags

- `-d <device>` - Device path (default: `/dev/video0`)
- `-all` - List all available controls
- `-class <name>` - List controls in specific class
- `-get <id>` - Get control value by ID (hex or decimal)
- `-set <id>` - Set control ID (use with `-value`)
- `-value <n>` - Value to set (use with `-set`)

## Example Output

### Overview Mode

```
V4L2 Control Reference - Device: /dev/video0
=============================================

Control Classes Overview
------------------------

  User               [Not supported  ] - Basic picture controls (brightness, contrast, etc.)
  Camera             [3 controls     ] - Camera controls (exposure, focus, zoom, pan/tilt)
  Flash              [Not supported  ] - Flash and LED controls
  JPEG               [Not supported  ] - JPEG compression settings
  Image Source       [Not supported  ] - Image sensor controls (gain, blanking, test patterns)
  Image Processing   [Not supported  ] - Image processing (pixel rate, deinterlacing)
  Codec              [Not supported  ] - Codec controls (bitrate, GOP, profiles)
  ...
```

### Class-Specific Listing

```
camera Controls (3):
================================================================================

  ID: 0x009a0901 (10094849)
  Name: Exposure, Auto
  Type: Menu
  Current Value: 3
  Menu Items:
    [0] Auto Mode
    [1] Manual Mode
    [2] Shutter Priority Mode
    [3] Aperture Priority Mode

  ID: 0x009a0903 (10094851)
  Name: Exposure, Absolute
  Type: Integer
  Current Value: 156
  Range: 3 to 2047 (step: 1)
  Default: 166
```

### Get Control

```
Control: Brightness (0x00980900)
Type: Integer
Current Value: 128
Range: 0 to 255 (step: 1)
```

### Set Control

```
Setting control: Brightness (0x00980900)
Successfully set to: 150
Verification: OK
```

## Control Types

The example handles all V4L2 control types:

- **Integer** - Numeric value with range and step
- **Boolean** - True/false (0/1)
- **Menu** - Selection from named options
- **Button** - Triggers an action
- **Integer64** - 64-bit numeric value
- **String** - Text value
- **Bitmask** - Bitfield value
- **IntegerMenu** - Selection from numeric options
- **U8/U16/U32** - Array controls

## Control Flags

Controls can have various flags:

- **disabled** - Control is disabled
- **grabbed** - Control is grabbed by another application
- **read-only** - Value cannot be changed
- **write-only** - Value cannot be read
- **volatile** - Value may change automatically
- **inactive** - Control is inactive in current configuration

## Common Control IDs

Some common controls you might want to query:

| Control | Class | ID (hex) | ID (decimal) |
|---------|-------|----------|--------------|
| Brightness | User | 0x00980900 | 9963776 |
| Contrast | User | 0x00980901 | 9963777 |
| Saturation | User | 0x00980902 | 9963778 |
| Hue | User | 0x00980903 | 9963779 |
| Exposure Auto | Camera | 0x009a0901 | 10094849 |
| Exposure Absolute | Camera | 0x009a0903 | 10094851 |
| Focus Auto | Camera | 0x009a090c | 10094860 |
| Focus Absolute | Camera | 0x009a090a | 10094858 |
| Zoom Absolute | Camera | 0x009a090d | 10094861 |

## Notes

- Not all devices support all control classes
- Some controls may be read-only or write-only
- The driver may adjust values to valid ranges
- Some controls become active/inactive based on other control values
- Menu controls show available options with their indices

## Related Documentation

- [V4L2 Controls Documentation](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html)
- [V4L2 Extended Controls](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/extended-controls.html)
- [Camera Controls](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-camera.html)
- [Codec Controls](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec.html)

## Code Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/vladimirvivien/go4vl/device"
    "github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
    dev, err := device.Open("/dev/video0")
    if err != nil {
        log.Fatal(err)
    }
    defer dev.Close()

    // Query a control
    ctrl, err := v4l2.QueryControlInfo(dev.Fd(), v4l2.CtrlCameraBrightness)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Control: %s\n", ctrl.Name)
    fmt.Printf("Range: %d to %d (step: %d)\n",
        ctrl.Minimum, ctrl.Maximum, ctrl.Step)

    // Get current value
    val, err := v4l2.GetControlValue(dev.Fd(), v4l2.CtrlCameraBrightness)
    if err == nil {
        fmt.Printf("Current value: %d\n", val)
    }

    // Set new value
    err = v4l2.SetControlValue(dev.Fd(), v4l2.CtrlCameraBrightness, 128)
    if err != nil {
        log.Printf("Failed to set value: %v", err)
    }
}
```

## See Also

- `examples/ext_ctrls/` - Extended controls example
- `v4l2/control_values.go` - All control constant definitions
- `v4l2/ext_ctrls_*.go` - Extended control structures for specific classes
