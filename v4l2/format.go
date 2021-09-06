package v4l2

import (
	"fmt"
	"unsafe"
)

// FourCCEncoding represents the four character encoding value
type FourCCEncoding = uint32

// Some Predefined pixel format definitions
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/pixfmt.html
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L518
var (
	PixFmtRGB24   = fourcc('R', 'G', 'B', '3') // V4L2_PIX_FMT_RGB24
	PixFmtGrey    = fourcc('G', 'R', 'E', 'Y') // V4L2_PIX_FMT_GREY
	PixelFmtYUYV  = fourcc('Y', 'U', 'Y', 'V') // V4L2_PIX_FMT_YUYV
	PixelFmtYYUV  = fourcc('Y', 'Y', 'U', 'V') // V4L2_PIX_FMT_YYUV
	PixelFmtYVYU  = fourcc('Y', 'V', 'Y', 'U') // V4L2_PIX_FMT_YVYU
	PixelFmtUYVY  = fourcc('U', 'Y', 'V', 'Y') // V4L2_PIX_FMT_UYVY
	PixelFmtVYUY  = fourcc('V', 'Y', 'U', 'Y') // V4L2_PIX_FMT_VYUY
	PixelFmtMJPEG = fourcc('M', 'J', 'P', 'G') // V4L2_PIX_FMT_MJPEG
	PixelFmtJPEG  = fourcc('J', 'P', 'E', 'G') // V4L2_PIX_FMT_JPEG
	PixelFmtMPEG  = fourcc('M', 'P', 'E', 'G') // V4L2_PIX_FMT_MPEG
	PixelFmtH264  = fourcc('H', '2', '6', '4') // V4L2_PIX_FMT_H264
	PixelFmtMPEG4 = fourcc('M', 'P', 'G', '4') // V4L2_PIX_FMT_MPEG4
)

// fourcc implements the four character code encoding found
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L81
// #define v4l2_fourcc(a, b, c, d)\
// 	 ((__u32)(a) | ((__u32)(b) << 8) | ((__u32)(c) << 16) | ((__u32)(d) << 24))
func fourcc(a, b, c, d uint32) FourCCEncoding {
	return (a | b<<8) | c<<16 | d<<24
}

// PixelFormats provides a map of FourCC encoding description
var PixelFormats = map[FourCCEncoding]string{
	PixFmtRGB24:   "24-bit RGB 8-8-8",
	PixFmtGrey:    "8-bit Greyscale",
	PixelFmtYUYV:  "YUYV 4:2:2",
	PixelFmtMJPEG: "Motion-JPEG",
	PixelFmtJPEG:  "JFIF JPEG",
	PixelFmtMPEG:  "MPEG-1/2/4",
	PixelFmtH264:  "H.264",
	PixelFmtMPEG4: "MPEG-4 Part 2 ES",
}

// YcbcrEncoding (v4l2_ycbcr_encoding)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_ycbcr_encoding
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L300
type YcbcrEncoding = uint32

const (
	YcbcrEncDefault YcbcrEncoding = iota // V4L2_YCBCR_ENC_DEFAULT
	YcbcrEnc601                          // V4L2_YCBCR_ENC_601
	YcbcrEnc709                          // V4L2_YCBCR_ENC_709
	YcbcrEncXV601                        // V4L2_YCBCR_ENC_XV601
	YcbcrEncXV709                        // V4L2_YCBCR_ENC_XV709
)

// Quantization (v4l2_quantization)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_quantization#c.V4L.v4l2_quantization
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L372
type Quantization = uint32

const (
	QuantizationDefault   Quantization = iota // V4L2_QUANTIZATION_DEFAULT
	QuantizationFullRange                     // V4L2_QUANTIZATION_FULL_RANGE
	QuantizationLimRange                      // V4L2_QUANTIZATION_LIM_RANGE
)

// XferFunctionType (v4l2_xfer_func)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_xfer_func#c.V4L.v4l2_xfer_func
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L259
type XferFunctionType = uint32

const (
	XferFuncDefault   XferFunctionType = iota // V4L2_XFER_FUNC_DEFAULT     = 0
	XferFunc709                               // V4L2_XFER_FUNC_709         = 1,
	XferFuncSRGB                              // V4L2_XFER_FUNC_SRGB        = 2,
	XferFuncOpRGB                             // V4L2_XFER_FUNC_OPRGB       = 3,
	XferFuncSMPTE240M                         // V4L2_XFER_FUNC_SMPTE240M   = 4,
	XferFuncNone                              // V4L2_XFER_FUNC_NONE        = 5,
	XferFuncDCIP3                             // V4L2_XFER_FUNC_DCI_P3      = 6,
	XferFuncSMPTE2084                         // V4L2_XFER_FUNC_SMPTE2084   = 7,
)

var XferFunctions = map[XferFunctionType]string{
	XferFuncDefault:   "Default",
	XferFunc709:       "Rec. 709",
	XferFuncSRGB:      "sRGB",
	XferFuncOpRGB:     "opRGB",
	XferFuncSMPTE240M: "SMPTE 240M",
	XferFuncNone:      "None",
	XferFuncDCIP3:     "DCI-P3",
	XferFuncSMPTE2084: "SMPTE 2084",
}

// FieldType (v4l2_field)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/field-order.html?highlight=v4l2_field#c.v4l2_field
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L88
type FieldType = uint32

const (
	FieldAny                 FieldType = iota // V4L2_FIELD_ANY
	FieldNone                                 // V4L2_FIELD_NONE
	FieldTop                                  // V4L2_FIELD_TOP
	FieldBottom                               // V4L2_FIELD_BOTTOM
	FieldInterlaced                           // V4L2_FIELD_INTERLACED
	FieldSequentialTopBottom                  // V4L2_FIELD_SEQ_TB
	FieldSequentialBottomTop                  // V4L2_FIELD_SEQ_BT
	FieldAlternate                            // V4L2_FIELD_ALTERNATE
	FieldInterlacedTopBottom                  // V4L2_FIELD_INTERLACED_TB
	FieldInterlacedBottomTop                  // V4L2_FIELD_INTERLACED_BT
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

// ColorspaceType
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L195
type ColorspaceType = uint32

const (
	ColorspaceTypeDefault ColorspaceType = iota //V4L2_COLORSPACE_DEFAULT
	ColorspaceSMPTE170M                         //V4L2_COLORSPACE_SMPTE170M
	ColorspaceSMPTE240M                         // V4L2_COLORSPACE_SMPTE240M
	ColorspaceREC709                            // V4L2_COLORSPACE_REC709
	ColorspaceBT878                             // V4L2_COLORSPACE_BT878 (absolete)
	Colorspace470SystemM                        // V4L2_COLORSPACE_470_SYSTEM_M (absolete)
	Colorspace470SystemBG                       // V4L2_COLORSPACE_470_SYSTEM_BG
	ColorspaceJPEG                              // V4L2_COLORSPACE_JPEG
	ColorspaceSRGB                              // V4L2_COLORSPACE_SRGB
	ColorspaceOPRGB                             // V4L2_COLORSPACE_OPRGB
	ColorspaceBT2020                            // V4L2_COLORSPACE_BT2020
	ColorspaceRaw                               // V4L2_COLORSPACE_RAW
	ColorspaceDCIP3                             // V4L2_COLORSPACE_DCI_P3
)

// Colorspaces is a map of colorspace to its respective description
var Colorspaces = map[ColorspaceType]string{
	ColorspaceTypeDefault: "Default",
	ColorspaceREC709:      "Rec. 709",
	Colorspace470SystemBG: "470 System BG",
	ColorspaceJPEG:        "JPEG",
	ColorspaceSRGB:        "sRGB",
	ColorspaceOPRGB:       "opRGB",
	ColorspaceBT2020:      "BT.2020",
	ColorspaceRaw:         "Raw",
	ColorspaceDCIP3:       "DCI-P3",
}

// PixFormat (v4l2_pix_format)
// https://www.kernel.org/doc/html/v4.9/media/uapi/v4l/pixfmt-002.html?highlight=v4l2_pix_format
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L496
type PixFormat struct {
	Width        uint32
	Height       uint32
	PixelFormat  FourCCEncoding
	Field        FieldType
	BytesPerLine uint32
	SizeImage    uint32
	Colorspace   ColorspaceType
	Priv         uint32
	Flags        uint32
	YcbcrEnc     YcbcrEncoding
	Quantization Quantization
	XferFunc     XferFunctionType
}

// v4l2Format (v4l2_format)
// https://www.kernel.org/doc/html/v4.9/media/uapi/v4l/vidioc-g-fmt.html?highlight=v4l2_format
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2303
//
// field fmt is a union, thus it's constructed as an appropriately sized array:
//
// struct v4l2_format {
// 	__u32	 type;
// 	union {
// 		struct v4l2_pix_format		    pix;
// 		struct v4l2_pix_format_mplane	pix_mp;
// 		struct v4l2_window		        win;
// 		struct v4l2_vbi_format		    vbi;
// 		struct v4l2_sliced_vbi_format	sliced;
// 		struct v4l2_sdr_format	 	    sdr;
// 		struct v4l2_meta_format		    meta;
// 		__u8	raw_data[200];   /* user-defined */
// 	} fmt;
// };
type v4l2Format struct {
	StreamType uint32
	fmt        [200]byte
}

// getPixFormat returns the PixFormat by casting the pointer to the union type
func (f v4l2Format) getPixFormat() PixFormat {
	pixfmt := (*PixFormat)(unsafe.Pointer(&f.fmt[0]))
	return *pixfmt
}

// setPixFormat sets the PixFormat by casting the pointer to the fmt union and set its value
func (f v4l2Format) setPixFormat(newPix PixFormat) {
	*(*PixFormat)(unsafe.Pointer(&f.fmt[0])) = newPix
}

// GetPixFormat retrieves pixel information for the specified driver
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-fmt.html#ioctl-vidioc-g-fmt-vidioc-s-fmt-vidioc-try-fmt
func GetPixFormat(fd uintptr) (PixFormat, error) {
	format := v4l2Format{StreamType: BufTypeVideoCapture}
	if err := Send(fd, VidiocGetFormat, uintptr(unsafe.Pointer(&format))); err != nil {
		return PixFormat{}, fmt.Errorf("pix format failed: %w", err)
	}

	return format.getPixFormat(), nil
}

// SetPixFormat sets the pixel format information for the specified driver
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-fmt.html#ioctl-vidioc-g-fmt-vidioc-s-fmt-vidioc-try-fmt
func SetPixFormat(fd uintptr, pixFmt PixFormat) error {
	format := v4l2Format{StreamType: BufTypeVideoCapture}
	format.setPixFormat(pixFmt)

	if err := Send(fd, VidiocSetFormat, uintptr(unsafe.Pointer(&format))); err != nil {
		return fmt.Errorf("pix format failed: %w", err)
	}
	return nil
}
