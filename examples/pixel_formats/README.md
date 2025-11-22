# V4L2 Pixel Formats Example

This example demonstrates how to enumerate and query pixel formats supported by V4L2 devices using the go4vl library.

## Overview

Pixel formats define how video data is encoded and stored. V4L2 supports a wide variety of formats including:
- **RGB formats**: Various RGB layouts (RGB332, RGB565, RGB24, RGB32, etc.)
- **YUV formats**: Packed, planar, and semi-planar layouts
- **Greyscale formats**: 8-bit to 16-bit monochrome
- **Bayer formats**: Raw sensor data (8-bit to 16-bit)
- **Compressed formats**: JPEG, H.264, HEVC, VP8, VP9, MPEG, etc.

## Building

```bash
cd examples/pixel_formats
go build
```

## Usage

### Show Current Format

Display the current format configuration of the device:

```bash
./pixel_formats -current
```

Example output:
```
Current Format - Device: /dev/video0
================================================================================

Pixel Format:     Motion-JPEG (FourCC: MJPG)
Resolution:       640 x 480
Field:            none
Bytes Per Line:   0
Size Image:       614400 bytes
Colorspace:       Default
YCbCr Encoding:   Default
Quantization:     Default
Transfer Func:    Default

Category:         JPEG
Flags:            Compressed
```

### List All Supported Formats

Show all formats supported by the device:

```bash
./pixel_formats -all
```

Example output:
```
Supported Formats - Device: /dev/video0
================================================================================

FourCC  Category         Description                    Flags
------  --------         -----------                    -----
MJPG    JPEG             Motion-JPEG                    Compressed
YUYV    YUV Packed       YUYV 4:2:2
NV12    YUV Semi-Planar  Y/CbCr 4:2:0

Total: 3 format(s)
```

### List Formats with Details

Show detailed information about each format:

```bash
./pixel_formats -all -detailed
```

Example output:
```
Format 1:
------------------------------------------------------------
  FourCC:       MJPG (0x47504a4d)
  Description:  Motion-JPEG
  Category:     JPEG
  Flags:        Compressed
  Type:         Compressed

Format 2:
------------------------------------------------------------
  FourCC:       YUYV (0x56595559)
  Description:  YUYV 4:2:2
  Category:     YUV Packed
  Bits/Pixel:   16
  Type:         YUV Packed

Total: 2 format(s)
```

### List Formats by Category

Filter formats by category:

```bash
# YUV formats
./pixel_formats -category yuv

# RGB formats
./pixel_formats -category rgb

# Compressed formats
./pixel_formats -category compressed

# JPEG formats
./pixel_formats -category jpeg

# H.264 formats
./pixel_formats -category h264

# Greyscale formats
./pixel_formats -category greyscale

# Bayer formats
./pixel_formats -category bayer
```

Example output:
```
Formats by Category: YUV - Device: /dev/video0
================================================================================

Found 2 format(s) in category 'yuv'

FourCC  Category         Description                    Flags
------  --------         -----------                    -----
YUYV    YUV Packed       YUYV 4:2:2
NV12    YUV Semi-Planar  Y/CbCr 4:2:0

Total: 2 format(s)
```

### Test Format Support

Check if a specific format is supported:

```bash
# Test if YUYV is supported
./pixel_formats -test YUYV

# Test if H264 is supported
./pixel_formats -test H264

# Test if NV12 is supported
./pixel_formats -test NV12
```

Example output (supported):
```
Testing Format Support: YUYV - Device: /dev/video0
================================================================================

✓ Format YUYV IS SUPPORTED

Description:  YUYV 4:2:2
FourCC:       YUYV (0x56595559)
Category:     YUV Packed
Bits/Pixel:   16
```

Example output (not supported):
```
Testing Format Support: H264 - Device: /dev/video0
================================================================================

✗ Format H264 IS NOT SUPPORTED
```

## Command Line Flags

- `-d <device>` - Device path (default: `/dev/video0`)
- `-all` - List all supported formats
- `-category <name>` - List formats by category
- `-current` - Show current format
- `-detailed` - Show detailed format information
- `-test <fourcc>` - Test if format is supported

## Format Categories

The example recognizes the following format categories:

| Category | Description | Examples |
|----------|-------------|----------|
| **rgb** | RGB color formats | RGB332, RGB565, RGB24, RGB32, ARGB32 |
| **yuv** | YUV color formats | YUYV, NV12, YUV420, YV12 |
| **greyscale** | Monochrome formats | GREY, Y10, Y12, Y16 |
| **bayer** | Raw sensor patterns | BGGR, GBRG, GRBG, RGGB (8/10/12/14/16-bit) |
| **jpeg** | JPEG variants | MJPEG, JPEG |
| **h264** | H.264/AVC | H264, H264NoSC, H264MVC, H264Slice |
| **hevc** | HEVC/H.265 | HEVC, HEVCSlice |
| **mpeg** | MPEG variants | MPEG1, MPEG2, MPEG4, XVID |
| **vp** | VP8/VP9 | VP8, VP9, VP8Frame, VP9Frame |
| **compressed** | Any compressed format | All of the above compression formats |

## Format Properties

For each format, the example can display:

- **FourCC**: Four-character code identifying the format
- **Description**: Human-readable format name
- **Category**: Format category (RGB, YUV, compressed, etc.)
- **Bits Per Pixel**: Average bits per pixel for uncompressed formats
- **Flags**: Format flags (compressed, emulated, etc.)
- **Type**: Format characteristics (packed, planar, compressed, etc.)

## Common FourCC Codes

Some commonly encountered pixel formats:

| FourCC | Description | Category |
|--------|-------------|----------|
| YUYV | YUYV 4:2:2 packed | YUV Packed |
| MJPG | Motion-JPEG | JPEG |
| NV12 | Y/CbCr 4:2:0 semi-planar | YUV Semi-Planar |
| RGB3 | 24-bit RGB | RGB |
| BGR3 | 24-bit BGR | RGB |
| H264 | H.264/AVC | H.264 |
| HEVC | HEVC/H.265 | HEVC |
| VP80 | VP8 | VP8/VP9 |
| GREY | 8-bit greyscale | Greyscale |
| BA81 | 8-bit Bayer BGGR | Bayer |

## Use Cases

This example is useful for:

1. **Discovery**: Finding out what formats your camera supports
2. **Debugging**: Checking current format configuration
3. **Validation**: Testing if a specific format is available
4. **Analysis**: Understanding format characteristics and capabilities
5. **Development**: Choosing appropriate formats for your application

## Related Examples

- `examples/format/` - Basic format querying
- `examples/webcam/` - Webcam capture with format selection
- `examples/control_reference/` - V4L2 control enumeration

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

	// Enumerate all formats
	for i := uint32(0); ; i++ {
		fmtDesc, err := v4l2.GetFormatDescription(dev.Fd(), i)
		if err != nil {
			break
		}

		// Create PixFormat to use helper methods
		pixFmt := v4l2.PixFormat{
			PixelFormat: fmtDesc.PixelFormat,
			Flags:       fmtDesc.Flags,
		}

		fmt.Printf("Format: %s\n", fmtDesc.Description)
		fmt.Printf("  Category: %s\n", pixFmt.GetCategory())

		if bpp := pixFmt.GetBitsPerPixel(); bpp > 0 {
			fmt.Printf("  Bits/Pixel: %d\n", bpp)
		}

		if pixFmt.IsCompressed() {
			fmt.Println("  Type: Compressed")
		}
	}
}
```

## See Also

- [V4L2 Format Documentation](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/pixfmt.html)
- [V4L2 Format Enumeration](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-fmt.html)
- [FourCC Codes](http://www.fourcc.org/)
