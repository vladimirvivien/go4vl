package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// FourCCType represents the four character encoding value
type FourCCType = uint32

// Pixel format definitions organized by category
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/pixfmt.html
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L518

// RGB Formats
var (
	PixelFmtRGB332  FourCCType = C.V4L2_PIX_FMT_RGB332
	PixelFmtARGB444 FourCCType = C.V4L2_PIX_FMT_ARGB444
	PixelFmtXRGB444 FourCCType = C.V4L2_PIX_FMT_XRGB444
	PixelFmtRGB555  FourCCType = C.V4L2_PIX_FMT_RGB555
	PixelFmtARGB555 FourCCType = C.V4L2_PIX_FMT_ARGB555
	PixelFmtXRGB555 FourCCType = C.V4L2_PIX_FMT_XRGB555
	PixelFmtRGB565  FourCCType = C.V4L2_PIX_FMT_RGB565
	PixelFmtRGB555X FourCCType = C.V4L2_PIX_FMT_RGB555X
	PixelFmtRGB565X FourCCType = C.V4L2_PIX_FMT_RGB565X
	PixelFmtBGR666  FourCCType = C.V4L2_PIX_FMT_BGR666
	PixelFmtBGR24   FourCCType = C.V4L2_PIX_FMT_BGR24
	PixelFmtRGB24   FourCCType = C.V4L2_PIX_FMT_RGB24
	PixelFmtBGR32   FourCCType = C.V4L2_PIX_FMT_BGR32
	PixelFmtABGR32  FourCCType = C.V4L2_PIX_FMT_ABGR32
	PixelFmtXBGR32  FourCCType = C.V4L2_PIX_FMT_XBGR32
	PixelFmtRGB32   FourCCType = C.V4L2_PIX_FMT_RGB32
	PixelFmtARGB32  FourCCType = C.V4L2_PIX_FMT_ARGB32
	PixelFmtXRGB32  FourCCType = C.V4L2_PIX_FMT_XRGB32
)

// Greyscale Formats
var (
	PixelFmtGrey   FourCCType = C.V4L2_PIX_FMT_GREY
	PixelFmtY4     FourCCType = C.V4L2_PIX_FMT_Y4
	PixelFmtY6     FourCCType = C.V4L2_PIX_FMT_Y6
	PixelFmtY10    FourCCType = C.V4L2_PIX_FMT_Y10
	PixelFmtY12    FourCCType = C.V4L2_PIX_FMT_Y12
	PixelFmtY14    FourCCType = C.V4L2_PIX_FMT_Y14
	PixelFmtY16    FourCCType = C.V4L2_PIX_FMT_Y16
	PixelFmtY16BE  FourCCType = C.V4L2_PIX_FMT_Y16_BE
	PixelFmtY10BPACK FourCCType = C.V4L2_PIX_FMT_Y10BPACK
)

// YUV Packed Formats
var (
	PixelFmtYUYV   FourCCType = C.V4L2_PIX_FMT_YUYV
	PixelFmtYYUV   FourCCType = C.V4L2_PIX_FMT_YYUV
	PixelFmtYVYU   FourCCType = C.V4L2_PIX_FMT_YVYU
	PixelFmtUYVY   FourCCType = C.V4L2_PIX_FMT_UYVY
	PixelFmtVYUY   FourCCType = C.V4L2_PIX_FMT_VYUY
	PixelFmtY41P   FourCCType = C.V4L2_PIX_FMT_Y41P
	PixelFmtYUV444 FourCCType = C.V4L2_PIX_FMT_YUV444
	PixelFmtYUV555 FourCCType = C.V4L2_PIX_FMT_YUV555
	PixelFmtYUV565 FourCCType = C.V4L2_PIX_FMT_YUV565
	PixelFmtYUV32  FourCCType = C.V4L2_PIX_FMT_YUV32
	PixelFmtAYUV32 FourCCType = C.V4L2_PIX_FMT_AYUV32
	PixelFmtXYUV32 FourCCType = C.V4L2_PIX_FMT_XYUV32
	PixelFmtVUYA32 FourCCType = C.V4L2_PIX_FMT_VUYA32
	PixelFmtVUYX32 FourCCType = C.V4L2_PIX_FMT_VUYX32
)

// YUV Planar Formats
var (
	PixelFmtYUV410  FourCCType = C.V4L2_PIX_FMT_YUV410
	PixelFmtYVU410  FourCCType = C.V4L2_PIX_FMT_YVU410
	PixelFmtYUV411P FourCCType = C.V4L2_PIX_FMT_YUV411P
	PixelFmtYUV420  FourCCType = C.V4L2_PIX_FMT_YUV420
	PixelFmtYVU420  FourCCType = C.V4L2_PIX_FMT_YVU420
	PixelFmtYUV422P FourCCType = C.V4L2_PIX_FMT_YUV422P
	PixelFmtYUV444M FourCCType = C.V4L2_PIX_FMT_YUV444M
	PixelFmtYVU444M FourCCType = C.V4L2_PIX_FMT_YVU444M
	PixelFmtYUV420M FourCCType = C.V4L2_PIX_FMT_YUV420M
	PixelFmtYVU420M FourCCType = C.V4L2_PIX_FMT_YVU420M
	PixelFmtYUV422M FourCCType = C.V4L2_PIX_FMT_YUV422M
	PixelFmtYVU422M FourCCType = C.V4L2_PIX_FMT_YVU422M
)

// YUV Semi-Planar (NV) Formats
var (
	PixelFmtNV12    FourCCType = C.V4L2_PIX_FMT_NV12
	PixelFmtNV21    FourCCType = C.V4L2_PIX_FMT_NV21
	PixelFmtNV16    FourCCType = C.V4L2_PIX_FMT_NV16
	PixelFmtNV61    FourCCType = C.V4L2_PIX_FMT_NV61
	PixelFmtNV24    FourCCType = C.V4L2_PIX_FMT_NV24
	PixelFmtNV42    FourCCType = C.V4L2_PIX_FMT_NV42
	PixelFmtNV12M   FourCCType = C.V4L2_PIX_FMT_NV12M
	PixelFmtNV21M   FourCCType = C.V4L2_PIX_FMT_NV21M
	PixelFmtNV16M   FourCCType = C.V4L2_PIX_FMT_NV16M
	PixelFmtNV61M   FourCCType = C.V4L2_PIX_FMT_NV61M
	PixelFmtP010    FourCCType = C.V4L2_PIX_FMT_P010
	PixelFmtP012    FourCCType = C.V4L2_PIX_FMT_P012
)

// Compressed Formats - JPEG
var (
	PixelFmtMJPEG FourCCType = C.V4L2_PIX_FMT_MJPEG
	PixelFmtJPEG  FourCCType = C.V4L2_PIX_FMT_JPEG
)

// Compressed Formats - H.26x
var (
	PixelFmtH263      FourCCType = C.V4L2_PIX_FMT_H263
	PixelFmtH264      FourCCType = C.V4L2_PIX_FMT_H264
	PixelFmtH264NoSC  FourCCType = C.V4L2_PIX_FMT_H264_NO_SC
	PixelFmtH264MVC   FourCCType = C.V4L2_PIX_FMT_H264_MVC
	PixelFmtH264Slice FourCCType = C.V4L2_PIX_FMT_H264_SLICE
	PixelFmtHEVC      FourCCType = C.V4L2_PIX_FMT_HEVC
	PixelFmtHEVCSlice FourCCType = C.V4L2_PIX_FMT_HEVC_SLICE
)

// Compressed Formats - MPEG
var (
	PixelFmtMPEG      FourCCType = C.V4L2_PIX_FMT_MPEG
	PixelFmtMPEG1     FourCCType = C.V4L2_PIX_FMT_MPEG1
	PixelFmtMPEG2     FourCCType = C.V4L2_PIX_FMT_MPEG2
	PixelFmtMPEG2Slice FourCCType = C.V4L2_PIX_FMT_MPEG2_SLICE
	PixelFmtMPEG4     FourCCType = C.V4L2_PIX_FMT_MPEG4
	PixelFmtXVID      FourCCType = C.V4L2_PIX_FMT_XVID
)

// Compressed Formats - VP
var (
	PixelFmtVP8      FourCCType = C.V4L2_PIX_FMT_VP8
	PixelFmtVP8Frame FourCCType = C.V4L2_PIX_FMT_VP8_FRAME
	PixelFmtVP9      FourCCType = C.V4L2_PIX_FMT_VP9
	PixelFmtVP9Frame FourCCType = C.V4L2_PIX_FMT_VP9_FRAME
)

// Compressed Formats - Other
var (
	PixelFmtAV1Frame      FourCCType = C.V4L2_PIX_FMT_AV1_FRAME
	PixelFmtVC1Annex      FourCCType = C.V4L2_PIX_FMT_VC1_ANNEX_G
	PixelFmtVC1AnnexL     FourCCType = C.V4L2_PIX_FMT_VC1_ANNEX_L
	PixelFmtFWHT          FourCCType = C.V4L2_PIX_FMT_FWHT
	PixelFmtFWHTStateless FourCCType = C.V4L2_PIX_FMT_FWHT_STATELESS
)

// Bayer Formats - 8-bit
var (
	PixelFmtSBGGR8 FourCCType = C.V4L2_PIX_FMT_SBGGR8
	PixelFmtSGBRG8 FourCCType = C.V4L2_PIX_FMT_SGBRG8
	PixelFmtSGRBG8 FourCCType = C.V4L2_PIX_FMT_SGRBG8
	PixelFmtSRGGB8 FourCCType = C.V4L2_PIX_FMT_SRGGB8
)

// Bayer Formats - 10-bit
var (
	PixelFmtSBGGR10      FourCCType = C.V4L2_PIX_FMT_SBGGR10
	PixelFmtSGBRG10      FourCCType = C.V4L2_PIX_FMT_SGBRG10
	PixelFmtSGRBG10      FourCCType = C.V4L2_PIX_FMT_SGRBG10
	PixelFmtSRGGB10      FourCCType = C.V4L2_PIX_FMT_SRGGB10
	PixelFmtSBGGR10ALAW8 FourCCType = C.V4L2_PIX_FMT_SBGGR10ALAW8
	PixelFmtSGBRG10ALAW8 FourCCType = C.V4L2_PIX_FMT_SGBRG10ALAW8
	PixelFmtSGRBG10ALAW8 FourCCType = C.V4L2_PIX_FMT_SGRBG10ALAW8
	PixelFmtSRGGB10ALAW8 FourCCType = C.V4L2_PIX_FMT_SRGGB10ALAW8
	PixelFmtSBGGR10DPCM8 FourCCType = C.V4L2_PIX_FMT_SBGGR10DPCM8
	PixelFmtSGBRG10DPCM8 FourCCType = C.V4L2_PIX_FMT_SGBRG10DPCM8
	PixelFmtSGRBG10DPCM8 FourCCType = C.V4L2_PIX_FMT_SGRBG10DPCM8
	PixelFmtSRGGB10DPCM8 FourCCType = C.V4L2_PIX_FMT_SRGGB10DPCM8
	PixelFmtSBGGR10P     FourCCType = C.V4L2_PIX_FMT_SBGGR10P
	PixelFmtSGBRG10P     FourCCType = C.V4L2_PIX_FMT_SGBRG10P
	PixelFmtSGRBG10P     FourCCType = C.V4L2_PIX_FMT_SGRBG10P
	PixelFmtSRGGB10P     FourCCType = C.V4L2_PIX_FMT_SRGGB10P
)

// Bayer Formats - 12-bit
var (
	PixelFmtSBGGR12 FourCCType = C.V4L2_PIX_FMT_SBGGR12
	PixelFmtSGBRG12 FourCCType = C.V4L2_PIX_FMT_SGBRG12
	PixelFmtSGRBG12 FourCCType = C.V4L2_PIX_FMT_SGRBG12
	PixelFmtSRGGB12 FourCCType = C.V4L2_PIX_FMT_SRGGB12
	PixelFmtSBGGR12P FourCCType = C.V4L2_PIX_FMT_SBGGR12P
	PixelFmtSGBRG12P FourCCType = C.V4L2_PIX_FMT_SGBRG12P
	PixelFmtSGRBG12P FourCCType = C.V4L2_PIX_FMT_SGRBG12P
	PixelFmtSRGGB12P FourCCType = C.V4L2_PIX_FMT_SRGGB12P
)

// Bayer Formats - 14-bit
var (
	PixelFmtSBGGR14 FourCCType = C.V4L2_PIX_FMT_SBGGR14
	PixelFmtSGBRG14 FourCCType = C.V4L2_PIX_FMT_SGBRG14
	PixelFmtSGRBG14 FourCCType = C.V4L2_PIX_FMT_SGRBG14
	PixelFmtSRGGB14 FourCCType = C.V4L2_PIX_FMT_SRGGB14
	PixelFmtSBGGR14P FourCCType = C.V4L2_PIX_FMT_SBGGR14P
	PixelFmtSGBRG14P FourCCType = C.V4L2_PIX_FMT_SGBRG14P
	PixelFmtSGRBG14P FourCCType = C.V4L2_PIX_FMT_SGRBG14P
	PixelFmtSRGGB14P FourCCType = C.V4L2_PIX_FMT_SRGGB14P
)

// Bayer Formats - 16-bit
var (
	PixelFmtSBGGR16 FourCCType = C.V4L2_PIX_FMT_SBGGR16
	PixelFmtSGBRG16 FourCCType = C.V4L2_PIX_FMT_SGBRG16
	PixelFmtSGRBG16 FourCCType = C.V4L2_PIX_FMT_SGRBG16
	PixelFmtSRGGB16 FourCCType = C.V4L2_PIX_FMT_SRGGB16
)

// PixelFormats provides a map of FourCCType encoding description
var PixelFormats = map[FourCCType]string{
	// RGB formats
	PixelFmtRGB332:  "8-bit RGB 3-3-2",
	PixelFmtARGB444: "16-bit ARGB 4-4-4-4",
	PixelFmtXRGB444: "16-bit XRGB 4-4-4-4",
	PixelFmtRGB555:  "16-bit RGB 5-5-5",
	PixelFmtARGB555: "16-bit ARGB 1-5-5-5",
	PixelFmtXRGB555: "16-bit XRGB 1-5-5-5",
	PixelFmtRGB565:  "16-bit RGB 5-6-5",
	PixelFmtRGB555X: "16-bit RGB 5-5-5 BE",
	PixelFmtRGB565X: "16-bit RGB 5-6-5 BE",
	PixelFmtBGR666:  "18-bit BGR 6-6-6",
	PixelFmtBGR24:   "24-bit BGR 8-8-8",
	PixelFmtRGB24:   "24-bit RGB 8-8-8",
	PixelFmtBGR32:   "32-bit BGR 8-8-8-8",
	PixelFmtABGR32:  "32-bit ABGR 8-8-8-8",
	PixelFmtXBGR32:  "32-bit XBGR 8-8-8-8",
	PixelFmtRGB32:   "32-bit RGB 8-8-8-8",
	PixelFmtARGB32:  "32-bit ARGB 8-8-8-8",
	PixelFmtXRGB32:  "32-bit XRGB 8-8-8-8",

	// Greyscale formats
	PixelFmtGrey:      "8-bit Greyscale",
	PixelFmtY4:        "4-bit Greyscale",
	PixelFmtY6:        "6-bit Greyscale",
	PixelFmtY10:       "10-bit Greyscale",
	PixelFmtY12:       "12-bit Greyscale",
	PixelFmtY14:       "14-bit Greyscale",
	PixelFmtY16:       "16-bit Greyscale",
	PixelFmtY16BE:     "16-bit Greyscale BE",
	PixelFmtY10BPACK: "10-bit Greyscale (packed)",

	// YUV packed formats
	PixelFmtYUYV:   "YUYV 4:2:2",
	PixelFmtYYUV:   "YYUV 4:2:2",
	PixelFmtYVYU:   "YVYU 4:2:2",
	PixelFmtUYVY:   "UYVY 4:2:2",
	PixelFmtVYUY:   "VYUY 4:2:2",
	PixelFmtY41P:   "YUV 4:1:1",
	PixelFmtYUV444: "YUV 4:4:4 (packed)",
	PixelFmtYUV555: "YUV 5:5:5 (packed)",
	PixelFmtYUV565: "YUV 5:6:5 (packed)",
	PixelFmtYUV32:  "32-bit YUV 8-8-8-8",
	PixelFmtAYUV32: "32-bit AYUV 8-8-8-8",
	PixelFmtXYUV32: "32-bit XYUV 8-8-8-8",
	PixelFmtVUYA32: "32-bit VUYA 8-8-8-8",
	PixelFmtVUYX32: "32-bit VUYX 8-8-8-8",

	// YUV planar formats
	PixelFmtYUV410:  "YUV 4:1:0 planar",
	PixelFmtYVU410:  "YVU 4:1:0 planar",
	PixelFmtYUV411P: "YUV 4:1:1 planar",
	PixelFmtYUV420:  "YUV 4:2:0 planar (I420)",
	PixelFmtYVU420:  "YVU 4:2:0 planar (YV12)",
	PixelFmtYUV422P: "YUV 4:2:2 planar",
	PixelFmtYUV444M: "YUV 4:4:4 planar (multiplanar)",
	PixelFmtYVU444M: "YVU 4:4:4 planar (multiplanar)",
	PixelFmtYUV420M: "YUV 4:2:0 planar (multiplanar)",
	PixelFmtYVU420M: "YVU 4:2:0 planar (multiplanar)",
	PixelFmtYUV422M: "YUV 4:2:2 planar (multiplanar)",
	PixelFmtYVU422M: "YVU 4:2:2 planar (multiplanar)",

	// YUV semi-planar (NV) formats
	PixelFmtNV12:  "YUV 4:2:0 semi-planar (NV12)",
	PixelFmtNV21:  "YUV 4:2:0 semi-planar (NV21)",
	PixelFmtNV16:  "YUV 4:2:2 semi-planar (NV16)",
	PixelFmtNV61:  "YUV 4:2:2 semi-planar (NV61)",
	PixelFmtNV24:  "YUV 4:4:4 semi-planar (NV24)",
	PixelFmtNV42:  "YUV 4:4:4 semi-planar (NV42)",
	PixelFmtNV12M: "YUV 4:2:0 semi-planar (NV12M multiplanar)",
	PixelFmtNV21M: "YUV 4:2:0 semi-planar (NV21M multiplanar)",
	PixelFmtNV16M: "YUV 4:2:2 semi-planar (NV16M multiplanar)",
	PixelFmtNV61M: "YUV 4:2:2 semi-planar (NV61M multiplanar)",
	PixelFmtP010:  "YUV 4:2:0 10-bit semi-planar (P010)",
	PixelFmtP012:  "YUV 4:2:0 12-bit semi-planar (P012)",

	// Compressed formats - JPEG
	PixelFmtMJPEG: "Motion-JPEG",
	PixelFmtJPEG:  "JFIF JPEG",

	// Compressed formats - H.26x
	PixelFmtH263:      "H.263",
	PixelFmtH264:      "H.264 / MPEG-4 AVC",
	PixelFmtH264NoSC:  "H.264 without start codes",
	PixelFmtH264MVC:   "H.264 MVC",
	PixelFmtH264Slice: "H.264 parsed slices",
	PixelFmtHEVC:      "H.265 / HEVC",
	PixelFmtHEVCSlice: "H.265 parsed slices",

	// Compressed formats - MPEG
	PixelFmtMPEG:       "MPEG-1/2/4",
	PixelFmtMPEG1:      "MPEG-1 ES",
	PixelFmtMPEG2:      "MPEG-2 ES",
	PixelFmtMPEG2Slice: "MPEG-2 parsed slices",
	PixelFmtMPEG4:      "MPEG-4 Part 2 ES",
	PixelFmtXVID:       "Xvid",

	// Compressed formats - VP
	PixelFmtVP8:      "VP8",
	PixelFmtVP8Frame: "VP8 frame",
	PixelFmtVP9:      "VP9",
	PixelFmtVP9Frame: "VP9 frame",

	// Compressed formats - Other
	PixelFmtAV1Frame:      "AV1 frame",
	PixelFmtVC1Annex:      "VC-1 Annex G",
	PixelFmtVC1AnnexL:     "VC-1 Annex L",
	PixelFmtFWHT:          "Fast Walsh Hadamard Transform",
	PixelFmtFWHTStateless: "FWHT stateless",

	// Bayer formats (selected common ones)
	PixelFmtSBGGR8:  "8-bit Bayer BGGR",
	PixelFmtSGBRG8:  "8-bit Bayer GBRG",
	PixelFmtSGRBG8:  "8-bit Bayer GRBG",
	PixelFmtSRGGB8:  "8-bit Bayer RGGB",
	PixelFmtSBGGR10: "10-bit Bayer BGGR",
	PixelFmtSGBRG10: "10-bit Bayer GBRG",
	PixelFmtSGRBG10: "10-bit Bayer GRBG",
	PixelFmtSRGGB10: "10-bit Bayer RGGB",
	PixelFmtSBGGR12: "12-bit Bayer BGGR",
	PixelFmtSGBRG12: "12-bit Bayer GBRG",
	PixelFmtSGRBG12: "12-bit Bayer GRBG",
	PixelFmtSRGGB12: "12-bit Bayer RGGB",
	PixelFmtSBGGR16: "16-bit Bayer BGGR",
	PixelFmtSGBRG16: "16-bit Bayer GBRG",
	PixelFmtSGRBG16: "16-bit Bayer GRBG",
	PixelFmtSRGGB16: "16-bit Bayer RGGB",
}

// IsPixYUVEncoded returns true if the pixel format is a chrome+luminance YUV format
func IsPixYUVEncoded(pixFmt FourCCType) bool {
	switch pixFmt {
	case
		PixelFmtYUYV,
		PixelFmtYYUV,
		PixelFmtYVYU,
		PixelFmtUYVY,
		PixelFmtVYUY:
		return true
	default:
		return false
	}
}

// FieldType (v4l2_field)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/field-order.html?highlight=v4l2_field#c.v4l2_field
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L88
type FieldType = uint32

const (
	FieldAny                 FieldType = C.V4L2_FIELD_ANY
	FieldNone                FieldType = C.V4L2_FIELD_NONE
	FieldTop                 FieldType = C.V4L2_FIELD_TOP
	FieldBottom              FieldType = C.V4L2_FIELD_BOTTOM
	FieldInterlaced          FieldType = C.V4L2_FIELD_INTERLACED
	FieldSequentialTopBottom FieldType = C.V4L2_FIELD_SEQ_TB
	FieldSequentialBottomTop FieldType = C.V4L2_FIELD_SEQ_BT
	FieldAlternate           FieldType = C.V4L2_FIELD_ALTERNATE
	FieldInterlacedTopBottom FieldType = C.V4L2_FIELD_INTERLACED_TB
	FieldInterlacedBottomTop FieldType = C.V4L2_FIELD_INTERLACED_BT
)

// Fields is a map of FieldType description
var Fields = map[FieldType]string{
	FieldAny:                 "any",
	FieldNone:                "none",
	FieldTop:                 "top",
	FieldBottom:              "bottom",
	FieldInterlaced:          "interlaced",
	FieldSequentialTopBottom: "sequential top-bottom",
	FieldSequentialBottomTop: "Sequential botton-top",
	FieldAlternate:           "alternating",
	FieldInterlacedTopBottom: "interlaced top-bottom",
	FieldInterlacedBottomTop: "interlaced bottom-top",
}

// PixFormatFlag represents format description flags from v4l2_fmtdesc
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-fmt.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L772
type PixFormatFlag = uint32

const (
	FmtFlagCompressed           PixFormatFlag = C.V4L2_FMT_FLAG_COMPRESSED
	FmtFlagEmulated             PixFormatFlag = C.V4L2_FMT_FLAG_EMULATED
	FmtFlagContinuousBytestream PixFormatFlag = C.V4L2_FMT_FLAG_CONTINUOUS_BYTESTREAM
	FmtFlagDynamicResolution    PixFormatFlag = C.V4L2_FMT_FLAG_DYN_RESOLUTION
	FmtFlagEncCapFrameInterval  PixFormatFlag = C.V4L2_FMT_FLAG_ENC_CAP_FRAME_INTERVAL
	FmtFlagCSCColorspace        PixFormatFlag = C.V4L2_FMT_FLAG_CSC_COLORSPACE
	FmtFlagCSCXferFunc          PixFormatFlag = C.V4L2_FMT_FLAG_CSC_XFER_FUNC
	FmtFlagCSCYCbCrEnc          PixFormatFlag = C.V4L2_FMT_FLAG_CSC_YCBCR_ENC
	FmtFlagCSCQuantization      PixFormatFlag = C.V4L2_FMT_FLAG_CSC_QUANTIZATION
	// FmtFlagMetaLineBased requires kernel 6.10+ (not available in Ubuntu 24.04)
	// FmtFlagMetaLineBased        PixFormatFlag = C.V4L2_FMT_FLAG_META_LINE_BASED
)

// FormatFlags is a map of format flag descriptions
var FormatFlags = map[PixFormatFlag]string{
	FmtFlagCompressed:           "Compressed",
	FmtFlagEmulated:             "Emulated (software conversion)",
	FmtFlagContinuousBytestream: "Continuous bytestream",
	FmtFlagDynamicResolution:    "Dynamic resolution change",
	FmtFlagEncCapFrameInterval:  "Encoder frame interval capture",
	FmtFlagCSCColorspace:        "Colorspace conversion supported",
	FmtFlagCSCXferFunc:          "Transfer function conversion supported",
	FmtFlagCSCYCbCrEnc:          "YCbCr encoding conversion supported",
	FmtFlagCSCQuantization:      "Quantization conversion supported",
	// FmtFlagMetaLineBased not included (requires kernel 6.10+)
}

// PixFormat contains video image format from v4l2_pix_format.
// https://www.kernel.org/doc/html/v4.9/media/uapi/v4l/pixfmt-002.html?highlight=v4l2_pix_format
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L496
type PixFormat struct {
	Width        uint32
	Height       uint32
	PixelFormat  FourCCType
	Field        FieldType
	BytesPerLine uint32
	SizeImage    uint32
	Colorspace   ColorspaceType
	Priv         uint32
	Flags        uint32
	YcbcrEnc     YCbCrEncodingType
	HSVEnc       HSVEncodingType
	Quantization QuantizationType
	XferFunc     XferFunctionType
}

func (f PixFormat) String() string {
	return fmt.Sprintf(
		"%s [%dx%d]; field=%s; bytes per line=%d; size image=%d; colorspace=%s; YCbCr=%s; Quant=%s; XferFunc=%s",
		PixelFormats[f.PixelFormat],
		f.Width, f.Height,
		Fields[f.Field],
		f.BytesPerLine,
		f.SizeImage,
		Colorspaces[f.Colorspace],
		YCbCrEncodings[f.YcbcrEnc],
		Quantizations[f.Quantization],
		XferFunctions[f.XferFunc],
	)
}

// IsCompressed returns true if the format is compressed
func (f PixFormat) IsCompressed() bool {
	return f.Flags&FmtFlagCompressed != 0
}

// IsEmulated returns true if the format is emulated (software conversion)
func (f PixFormat) IsEmulated() bool {
	return f.Flags&FmtFlagEmulated != 0
}

// SupportsDynamicResolution returns true if the format supports dynamic resolution changes
func (f PixFormat) SupportsDynamicResolution() bool {
	return f.Flags&FmtFlagDynamicResolution != 0
}

// SupportsColorspaceConversion returns true if the format supports colorspace conversion
func (f PixFormat) SupportsColorspaceConversion() bool {
	return f.Flags&FmtFlagCSCColorspace != 0
}

// GetFlags returns a list of flag descriptions for this format
func (f PixFormat) GetFlags() []string {
	var flags []string
	for flag, desc := range FormatFlags {
		if f.Flags&flag != 0 {
			flags = append(flags, desc)
		}
	}
	return flags
}

// IsRGB returns true if the pixel format is an RGB format
func (f PixFormat) IsRGB() bool {
	switch f.PixelFormat {
	case PixelFmtRGB332, PixelFmtARGB444, PixelFmtXRGB444,
		PixelFmtARGB555, PixelFmtXRGB555,
		PixelFmtRGB565, PixelFmtRGB555, PixelFmtRGB555X, PixelFmtRGB565X,
		PixelFmtBGR666, PixelFmtBGR24, PixelFmtRGB24,
		PixelFmtBGR32, PixelFmtABGR32, PixelFmtXBGR32,
		PixelFmtRGB32, PixelFmtARGB32, PixelFmtXRGB32:
		return true
	default:
		return false
	}
}

// IsYUV returns true if the pixel format is a YUV format (packed, planar, or semi-planar)
func (f PixFormat) IsYUV() bool {
	return f.IsYUVPacked() || f.IsYUVPlanar() || f.IsYUVSemiPlanar()
}

// IsYUVPacked returns true if the pixel format is a packed YUV format
func (f PixFormat) IsYUVPacked() bool {
	switch f.PixelFormat {
	case PixelFmtYUYV, PixelFmtYYUV, PixelFmtYVYU, PixelFmtUYVY, PixelFmtVYUY,
		PixelFmtYUV444, PixelFmtYUV555, PixelFmtYUV565, PixelFmtYUV32,
		PixelFmtAYUV32, PixelFmtXYUV32, PixelFmtVUYA32, PixelFmtVUYX32:
		return true
	default:
		return false
	}
}

// IsYUVPlanar returns true if the pixel format is a planar YUV format
func (f PixFormat) IsYUVPlanar() bool {
	switch f.PixelFormat {
	case PixelFmtYUV410, PixelFmtYUV420, PixelFmtYVU410, PixelFmtYVU420,
		PixelFmtYUV422P, PixelFmtYUV411P, PixelFmtY41P,
		PixelFmtYUV444M, PixelFmtYVU444M, PixelFmtYUV422M, PixelFmtYVU422M, PixelFmtYUV420M, PixelFmtYVU420M:
		return true
	default:
		return false
	}
}

// IsYUVSemiPlanar returns true if the pixel format is a semi-planar YUV (NV) format
func (f PixFormat) IsYUVSemiPlanar() bool {
	switch f.PixelFormat {
	case PixelFmtNV12, PixelFmtNV21, PixelFmtNV16, PixelFmtNV61, PixelFmtNV24, PixelFmtNV42,
		PixelFmtNV12M, PixelFmtNV21M, PixelFmtNV16M, PixelFmtNV61M,
		PixelFmtP010, PixelFmtP012:
		return true
	default:
		return false
	}
}

// IsGreyscale returns true if the pixel format is a greyscale format
func (f PixFormat) IsGreyscale() bool {
	switch f.PixelFormat {
	case PixelFmtGrey, PixelFmtY4, PixelFmtY6, PixelFmtY10, PixelFmtY12, PixelFmtY14, PixelFmtY16, PixelFmtY16BE, PixelFmtY10BPACK:
		return true
	default:
		return false
	}
}

// IsBayer returns true if the pixel format is a Bayer pattern format
func (f PixFormat) IsBayer() bool {
	switch f.PixelFormat {
	// 8-bit Bayer
	case PixelFmtSBGGR8, PixelFmtSGBRG8, PixelFmtSGRBG8, PixelFmtSRGGB8,
		// 10-bit Bayer
		PixelFmtSBGGR10, PixelFmtSGBRG10, PixelFmtSGRBG10, PixelFmtSRGGB10,
		PixelFmtSBGGR10P, PixelFmtSGBRG10P, PixelFmtSGRBG10P, PixelFmtSRGGB10P,
		PixelFmtSBGGR10ALAW8, PixelFmtSGBRG10ALAW8, PixelFmtSGRBG10ALAW8, PixelFmtSRGGB10ALAW8,
		PixelFmtSBGGR10DPCM8, PixelFmtSGBRG10DPCM8, PixelFmtSGRBG10DPCM8, PixelFmtSRGGB10DPCM8,
		// 12-bit Bayer
		PixelFmtSBGGR12, PixelFmtSGBRG12, PixelFmtSGRBG12, PixelFmtSRGGB12,
		PixelFmtSBGGR12P, PixelFmtSGBRG12P, PixelFmtSGRBG12P, PixelFmtSRGGB12P,
		// 14-bit Bayer
		PixelFmtSBGGR14, PixelFmtSGBRG14, PixelFmtSGRBG14, PixelFmtSRGGB14,
		PixelFmtSBGGR14P, PixelFmtSGBRG14P, PixelFmtSGRBG14P, PixelFmtSRGGB14P,
		// 16-bit Bayer
		PixelFmtSBGGR16, PixelFmtSGBRG16, PixelFmtSGRBG16, PixelFmtSRGGB16:
		return true
	default:
		return false
	}
}

// IsH264 returns true if the pixel format is an H.264 variant
func (f PixFormat) IsH264() bool {
	switch f.PixelFormat {
	case PixelFmtH264, PixelFmtH264NoSC, PixelFmtH264MVC, PixelFmtH264Slice:
		return true
	default:
		return false
	}
}

// IsHEVC returns true if the pixel format is an HEVC/H.265 variant
func (f PixFormat) IsHEVC() bool {
	switch f.PixelFormat {
	case PixelFmtHEVC, PixelFmtHEVCSlice:
		return true
	default:
		return false
	}
}

// IsMPEG returns true if the pixel format is an MPEG variant
func (f PixFormat) IsMPEG() bool {
	switch f.PixelFormat {
	case PixelFmtMPEG1, PixelFmtMPEG2, PixelFmtMPEG2Slice, PixelFmtMPEG4, PixelFmtXVID:
		return true
	default:
		return false
	}
}

// IsVP returns true if the pixel format is a VP8/VP9 variant
func (f PixFormat) IsVP() bool {
	switch f.PixelFormat {
	case PixelFmtVP8, PixelFmtVP8Frame, PixelFmtVP9, PixelFmtVP9Frame:
		return true
	default:
		return false
	}
}

// IsJPEG returns true if the pixel format is JPEG or Motion-JPEG
func (f PixFormat) IsJPEG() bool {
	switch f.PixelFormat {
	case PixelFmtJPEG, PixelFmtMJPEG:
		return true
	default:
		return false
	}
}

// GetCategory returns a human-readable category for the pixel format
func (f PixFormat) GetCategory() string {
	if f.IsRGB() {
		return "RGB"
	}
	if f.IsYUVPacked() {
		return "YUV Packed"
	}
	if f.IsYUVPlanar() {
		return "YUV Planar"
	}
	if f.IsYUVSemiPlanar() {
		return "YUV Semi-Planar"
	}
	if f.IsGreyscale() {
		return "Greyscale"
	}
	if f.IsBayer() {
		return "Bayer"
	}
	if f.IsJPEG() {
		return "JPEG"
	}
	if f.IsH264() {
		return "H.264"
	}
	if f.IsHEVC() {
		return "HEVC/H.265"
	}
	if f.IsMPEG() {
		return "MPEG"
	}
	if f.IsVP() {
		return "VP8/VP9"
	}
	if f.IsCompressed() {
		return "Compressed"
	}
	return "Other"
}

// GetBitsPerPixel returns the average bits per pixel for uncompressed formats
// Returns 0 for compressed formats or unknown formats
func (f PixFormat) GetBitsPerPixel() int {
	switch f.PixelFormat {
	// RGB formats
	case PixelFmtRGB332:
		return 8
	case PixelFmtARGB444, PixelFmtXRGB444:
		return 16
	case PixelFmtRGB555, PixelFmtRGB555X, PixelFmtRGB565, PixelFmtRGB565X,
		PixelFmtARGB555, PixelFmtXRGB555:
		return 16
	case PixelFmtBGR666:
		return 18
	case PixelFmtBGR24, PixelFmtRGB24:
		return 24
	case PixelFmtBGR32, PixelFmtABGR32, PixelFmtXBGR32,
		PixelFmtRGB32, PixelFmtARGB32, PixelFmtXRGB32:
		return 32

	// Greyscale
	case PixelFmtGrey:
		return 8
	case PixelFmtY10, PixelFmtY10BPACK:
		return 10
	case PixelFmtY12:
		return 12
	case PixelFmtY14:
		return 14
	case PixelFmtY16, PixelFmtY16BE:
		return 16

	// Packed YUV
	case PixelFmtYUYV, PixelFmtYYUV, PixelFmtYVYU, PixelFmtUYVY, PixelFmtVYUY, PixelFmtYUV555, PixelFmtYUV565:
		return 16
	case PixelFmtYUV444:
		return 24
	case PixelFmtAYUV32, PixelFmtXYUV32, PixelFmtVUYA32, PixelFmtVUYX32, PixelFmtYUV32:
		return 32

	// Planar YUV (average)
	case PixelFmtYUV410, PixelFmtYVU410:
		return 9 // 4:1:0
	case PixelFmtYUV420, PixelFmtYVU420, PixelFmtYUV420M, PixelFmtYVU420M, PixelFmtNV12, PixelFmtNV21, PixelFmtNV12M, PixelFmtNV21M:
		return 12 // 4:2:0
	case PixelFmtYUV422P, PixelFmtYUV411P, PixelFmtY41P, PixelFmtYUV422M, PixelFmtYVU422M, PixelFmtNV16, PixelFmtNV61, PixelFmtNV16M, PixelFmtNV61M:
		return 16 // 4:2:2
	case PixelFmtYUV444M, PixelFmtYVU444M, PixelFmtNV24, PixelFmtNV42:
		return 24 // 4:4:4

	// Bayer - 8 bit
	case PixelFmtSBGGR8, PixelFmtSGBRG8, PixelFmtSGRBG8, PixelFmtSRGGB8:
		return 8
	// Bayer - 10 bit
	case PixelFmtSBGGR10, PixelFmtSGBRG10, PixelFmtSGRBG10, PixelFmtSRGGB10,
		PixelFmtSBGGR10P, PixelFmtSGBRG10P, PixelFmtSGRBG10P, PixelFmtSRGGB10P,
		PixelFmtSBGGR10ALAW8, PixelFmtSGBRG10ALAW8, PixelFmtSGRBG10ALAW8, PixelFmtSRGGB10ALAW8,
		PixelFmtSBGGR10DPCM8, PixelFmtSGBRG10DPCM8, PixelFmtSGRBG10DPCM8, PixelFmtSRGGB10DPCM8:
		return 10
	// Bayer - 12 bit
	case PixelFmtSBGGR12, PixelFmtSGBRG12, PixelFmtSGRBG12, PixelFmtSRGGB12,
		PixelFmtSBGGR12P, PixelFmtSGBRG12P, PixelFmtSGRBG12P, PixelFmtSRGGB12P:
		return 12
	// Bayer - 14 bit
	case PixelFmtSBGGR14, PixelFmtSGBRG14, PixelFmtSGRBG14, PixelFmtSRGGB14,
		PixelFmtSBGGR14P, PixelFmtSGBRG14P, PixelFmtSGRBG14P, PixelFmtSRGGB14P:
		return 14
	// Bayer - 16 bit
	case PixelFmtSBGGR16, PixelFmtSGBRG16, PixelFmtSGRBG16, PixelFmtSRGGB16:
		return 16

	default:
		// Compressed formats or unknown
		return 0
	}
}

// GetPixFormat retrieves pixel information for the specified driver (via v4l2_format and v4l2_pix_format)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2331
// and https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-fmt.html#ioctl-vidioc-g-fmt-vidioc-s-fmt-vidioc-try-fmt
func GetPixFormat(fd uintptr) (PixFormat, error) {
	var v4l2Format C.struct_v4l2_format
	v4l2Format._type = C.uint(BufTypeVideoCapture)

	if err := send(fd, C.VIDIOC_G_FMT, uintptr(unsafe.Pointer(&v4l2Format))); err != nil {
		return PixFormat{}, fmt.Errorf("pix format failed: %w", err)
	}

	v4l2PixFmt := *(*C.struct_v4l2_pix_format)(unsafe.Pointer(&v4l2Format.fmt[0]))
	return PixFormat{
		Width:        uint32(v4l2PixFmt.width),
		Height:       uint32(v4l2PixFmt.height),
		PixelFormat:  uint32(v4l2PixFmt.pixelformat),
		Field:        uint32(v4l2PixFmt.field),
		BytesPerLine: uint32(v4l2PixFmt.bytesperline),
		SizeImage:    uint32(v4l2PixFmt.sizeimage),
		Colorspace:   uint32(v4l2PixFmt.colorspace),
		Priv:         uint32(v4l2PixFmt.priv),
		Flags:        uint32(v4l2PixFmt.flags),
		YcbcrEnc:     *(*uint32)(unsafe.Pointer(&v4l2PixFmt.anon0[0])),
		HSVEnc:       *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&v4l2PixFmt.anon0[0])) + unsafe.Sizeof(C.uint(0)))),
		Quantization: uint32(v4l2PixFmt.quantization),
		XferFunc:     uint32(v4l2PixFmt.xfer_func),
	}, nil
}

// SetPixFormat sets the pixel format information for the specified driver
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-fmt.html#ioctl-vidioc-g-fmt-vidioc-s-fmt-vidioc-try-fmt
func SetPixFormat(fd uintptr, pixFmt PixFormat) error {
	var v4l2Format C.struct_v4l2_format
	v4l2Format._type = C.uint(BufTypeVideoCapture)
	*(*C.struct_v4l2_pix_format)(unsafe.Pointer(&v4l2Format.fmt[0])) = *(*C.struct_v4l2_pix_format)(unsafe.Pointer(&pixFmt))

	if err := send(fd, C.VIDIOC_S_FMT, uintptr(unsafe.Pointer(&v4l2Format))); err != nil {
		return fmt.Errorf("pix format failed: %w", err)
	}
	return nil
}
