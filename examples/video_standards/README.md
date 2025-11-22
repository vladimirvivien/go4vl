# Video Standards Example

This example demonstrates how to work with analog video standards (PAL, NTSC, SECAM) using the go4vl library.

## Overview

Video standards define the analog video signal format used by legacy analog video devices like TV tuners, composite video inputs, and analog capture cards. This is different from digital video timings (DV timings) used by modern HDMI/DisplayPort devices.

Common video standards include:
- **PAL** (Phase Alternating Line) - Used in Europe, Asia, Africa, Australia
- **NTSC** (National Television System Committee) - Used in North America, parts of South America, Japan
- **SECAM** (Séquentiel couleur à mémoire) - Used in France, Russia, parts of Africa

## Features

This example shows how to:
- Check if a device supports analog video standards
- Enumerate all supported standards
- Get the current video standard
- Set a new video standard
- Auto-detect the standard from input signal
- Check support for specific standards

## Building

```bash
cd examples/video_standards
go build
```

## Usage

### Display Current Standard and Supported Standards

```bash
./video_standards -d /dev/video0
```

This will show:
- Current video standard
- All supported standards with frame rates and line counts
- Support status for common standards

### Query/Auto-detect Standard from Signal

```bash
./video_standards -d /dev/video0 -query
```

Attempts to auto-detect the video standard from the current input signal. This requires:
- A device that supports standard detection
- An active analog video signal on the input

### Set Video Standard

```bash
# Set to PAL-B/G (Western Europe)
./video_standards -d /dev/video0 -set PAL_BG

# Set to NTSC-M (USA)
./video_standards -d /dev/video0 -set NTSC_M

# Set to SECAM-L (France)
./video_standards -d /dev/video0 -set SECAM_L

# Set to all PAL variants
./video_standards -d /dev/video0 -set PAL
```

Available standard names:
- PAL, PAL_BG, PAL_DK, PAL_I, PAL_M, PAL_N, PAL_B, PAL_G, PAL_H, PAL_D, PAL_K
- NTSC, NTSC_M, NTSC_M_JP, NTSC_M_KR
- SECAM, SECAM_B, SECAM_D, SECAM_G, SECAM_H, SECAM_K, SECAM_L, SECAM_DK
- 525_60 (525 lines, 60 Hz - NTSC family)
- 625_50 (625 lines, 50 Hz - PAL/SECAM family)

## Command Line Flags

- `-d <device>` - Device path (default: `/dev/video0`)
- `-set <standard>` - Set video standard (e.g., PAL, NTSC, PAL_BG, NTSC_M)
- `-query` - Auto-detect standard from input signal

## Video Standards Reference

### PAL Variants

| Standard | Region | Lines | FPS | Description |
|----------|--------|-------|-----|-------------|
| PAL-B | Western Europe | 625 | 25 | Belgium, Netherlands, Switzerland |
| PAL-G | Western Europe | 625 | 25 | Germany, Austria, Portugal |
| PAL-H | Western Europe | 625 | 25 | Belgium |
| PAL-I | UK/Ireland | 625 | 25 | United Kingdom, Ireland |
| PAL-D | Eastern Europe/China | 625 | 25 | China |
| PAL-K | Eastern Europe | 625 | 25 | |
| PAL-M | Brazil | 525 | ~29.97 | Brazil only |
| PAL-N | South America | 625 | 25 | Argentina, Paraguay, Uruguay |

### NTSC Variants

| Standard | Region | Lines | FPS | Description |
|----------|--------|-------|-----|-------------|
| NTSC-M | North America | 525 | ~29.97 | USA, Canada (BTSC audio) |
| NTSC-M-JP | Japan | 525 | ~29.97 | Japan (EIA-J audio) |
| NTSC-M-KR | Korea | 525 | ~29.97 | Korea (FM A2 audio) |
| NTSC-443 | - | 525 | ~29.97 | NTSC with PAL color subcarrier |

### SECAM Variants

| Standard | Region | Lines | FPS | Description |
|----------|--------|-------|-----|-------------|
| SECAM-B | Europe | 625 | 25 | |
| SECAM-D | Eastern Europe | 625 | 25 | |
| SECAM-G | Middle East | 625 | 25 | |
| SECAM-H | - | 625 | 25 | |
| SECAM-K | Eastern Europe | 625 | 25 | |
| SECAM-L | France | 625 | 25 | France, Luxembourg |

## Example Output

### When Device Supports Standards (Analog Capture Card)

```
Video Standards Example - Device: /dev/video0
==========================================

Current Video Standard
----------------------
  Standard ID: 0x0000000000000005
  Standard Name: PAL-B/G
  Family: PAL
  Format: 625 lines / 50 Hz

Supported Video Standards
-------------------------
  [0] PAL-B/G
      ID: 0x0000000000000005
      Frame rate: 25.00 fps
      Frame lines: 625
      Frame period: 1/25 seconds
      Family: PAL

  [1] NTSC-M
      ID: 0x0000000000001000
      Frame rate: 29.97 fps
      Frame lines: 525
      Frame period: 1001/30000 seconds
      Family: NTSC

Common Standard Support
-----------------------
  PAL          [Supported] - All PAL variants
  PAL-B/G      [Supported] - PAL B/G (Western Europe)
  NTSC         [Supported] - All NTSC variants
  NTSC-M       [Supported] - NTSC M (USA, Canada, BTSC)
  525/60       [Supported] - 525 lines, 60 Hz (NTSC)
  625/50       [Supported] - 625 lines, 50 Hz (PAL/SECAM)
```

### When Device Doesn't Support Standards (Digital Webcam)

```
Video Standards Example - Device: /dev/video0
==========================================

Device does not support analog video standards.

Note: Video standards are typically supported by:
  - TV tuner cards
  - Composite/S-Video capture cards
  - Analog cameras
  - Legacy video equipment

Modern digital devices (HDMI, USB webcams, etc.) use DV timings instead.
```

## Devices That Support Video Standards

Video standards are typically supported by:

1. **TV Tuner Cards**
   - PCI/PCIe TV tuners
   - USB TV tuners
   - Capture cards with analog tuners

2. **Analog Video Capture Devices**
   - Composite video (RCA yellow) inputs
   - S-Video inputs
   - Component video inputs

3. **Legacy Video Equipment**
   - Analog security cameras
   - VCRs and DVD players (via analog output)
   - Analog broadcast equipment

4. **Multi-Format Capture Cards**
   - Cards that support both analog and digital inputs
   - May need to select analog input first

## Modern Alternatives

For modern digital video devices, use **DV Timings** instead:
- HDMI capture cards → Use `examples/dv_timings`
- USB webcams → No standards/timings needed (auto-detected)
- DisplayPort/DVI → Use DV timings

See `examples/dv_timings/` for digital video timing examples.

## Notes

- Changing the video standard may also change the current pixel format
- Some devices automatically detect and switch standards
- The driver may adjust the requested standard to the closest supported variant
- Not all devices support standard detection (`-query` mode)
- Some devices may have read-only standards (cannot be changed)

## See Also

- [V4L2 Video Standards Documentation](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/standard.html)
- [V4L2 VIDIOC_ENUMSTD](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enumstd.html)
- [V4L2 VIDIOC_G_STD](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-std.html)
- [PAL on Wikipedia](https://en.wikipedia.org/wiki/PAL)
- [NTSC on Wikipedia](https://en.wikipedia.org/wiki/NTSC)
- [SECAM on Wikipedia](https://en.wikipedia.org/wiki/SECAM)

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

	// Enumerate supported standards
	standards, err := dev.GetAllStandards()
	if err != nil {
		log.Printf("Device doesn't support analog standards: %v", err)
		return
	}

	fmt.Println("Supported standards:")
	for _, std := range standards {
		fmt.Printf("  %s (%.2f fps, %d lines)\n",
			std.Name(), std.FrameRate(), std.FrameLines())
	}

	// Get current standard
	currentStd, err := dev.GetStandard()
	if err == nil {
		fmt.Printf("\nCurrent standard: %s\n", v4l2.StdNames[currentStd])
	}

	// Set to PAL if supported
	if supported, _ := dev.IsStandardSupported(v4l2.StdPAL); supported {
		err = dev.SetStandard(v4l2.StdPAL)
		if err == nil {
			fmt.Println("Successfully set to PAL")
		}
	}
}
```
