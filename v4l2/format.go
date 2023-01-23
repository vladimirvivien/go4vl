package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// FourCCType represents the four character encoding value
type FourCCType = uint32

// Some Predefined pixel format definitions
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/pixfmt.html
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L518
var (
	PixelFmtRGB24 FourCCType = C.V4L2_PIX_FMT_RGB24
	PixelFmtGrey  FourCCType = C.V4L2_PIX_FMT_GREY
	PixelFmtYUYV  FourCCType = C.V4L2_PIX_FMT_YUYV
	PixelFmtYYUV  FourCCType = C.V4L2_PIX_FMT_YYUV
	PixelFmtYVYU  FourCCType = C.V4L2_PIX_FMT_YVYU
	PixelFmtUYVY  FourCCType = C.V4L2_PIX_FMT_UYVY
	PixelFmtVYUY  FourCCType = C.V4L2_PIX_FMT_VYUY
	PixelFmtMJPEG FourCCType = C.V4L2_PIX_FMT_MJPEG
	PixelFmtJPEG  FourCCType = C.V4L2_PIX_FMT_JPEG
	PixelFmtMPEG  FourCCType = C.V4L2_PIX_FMT_MPEG
	PixelFmtH264  FourCCType = C.V4L2_PIX_FMT_H264
	PixelFmtMPEG4 FourCCType = C.V4L2_PIX_FMT_MPEG4
)

// PixelFormats provides a map of FourCCType encoding description
var PixelFormats = map[FourCCType]string{
	PixelFmtRGB24: "24-bit RGB 8-8-8",
	PixelFmtGrey:  "8-bit Greyscale",
	PixelFmtYUYV:  "YUYV 4:2:2",
	PixelFmtMJPEG: "Motion-JPEG",
	PixelFmtJPEG:  "JFIF JPEG",
	PixelFmtMPEG:  "MPEG-1/2/4",
	PixelFmtH264:  "H.264",
	PixelFmtMPEG4: "MPEG-4 Part 2 ES",
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

// ColorspaceType
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L195
type ColorspaceType = uint32

const (
	ColorspaceDefault     ColorspaceType = C.V4L2_COLORSPACE_DEFAULT
	ColorspaceSMPTE170M   ColorspaceType = C.V4L2_COLORSPACE_SMPTE170M
	ColorspaceSMPTE240M   ColorspaceType = C.V4L2_COLORSPACE_SMPTE240M
	ColorspaceREC709      ColorspaceType = C.V4L2_COLORSPACE_REC709
	ColorspaceBT878       ColorspaceType = C.V4L2_COLORSPACE_BT878        //(absolete)
	Colorspace470SystemM  ColorspaceType = C.V4L2_COLORSPACE_470_SYSTEM_M //(absolete)
	Colorspace470SystemBG ColorspaceType = C.V4L2_COLORSPACE_470_SYSTEM_BG
	ColorspaceJPEG        ColorspaceType = C.V4L2_COLORSPACE_JPEG
	ColorspaceSRGB        ColorspaceType = C.V4L2_COLORSPACE_SRGB
	ColorspaceOPRGB       ColorspaceType = C.V4L2_COLORSPACE_OPRGB
	ColorspaceBT2020      ColorspaceType = C.V4L2_COLORSPACE_BT2020
	ColorspaceRaw         ColorspaceType = C.V4L2_COLORSPACE_RAW
	ColorspaceDCIP3       ColorspaceType = C.V4L2_COLORSPACE_DCI_P3
)

// Colorspaces is a map of colorspace to its respective description
var Colorspaces = map[ColorspaceType]string{
	ColorspaceDefault:     "Default",
	ColorspaceREC709:      "Rec. 709",
	Colorspace470SystemBG: "470 System BG",
	ColorspaceJPEG:        "JPEG",
	ColorspaceSRGB:        "sRGB",
	ColorspaceOPRGB:       "opRGB",
	ColorspaceBT2020:      "BT.2020",
	ColorspaceRaw:         "Raw",
	ColorspaceDCIP3:       "DCI-P3",
}

// YCbCrEncodingType (v4l2_ycbcr_encoding)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_ycbcr_encoding
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L300
type YCbCrEncodingType = uint32

const (
	YCbCrEncodingDefault        YCbCrEncodingType = C.V4L2_YCBCR_ENC_DEFAULT
	YCbCrEncoding601            YCbCrEncodingType = C.V4L2_YCBCR_ENC_601
	YCbCrEncoding709            YCbCrEncodingType = C.V4L2_YCBCR_ENC_709
	YCbCrEncodingXV601          YCbCrEncodingType = C.V4L2_YCBCR_ENC_XV601
	YCbCrEncodingXV709          YCbCrEncodingType = C.V4L2_YCBCR_ENC_XV709
	_                           YCbCrEncodingType = C.V4L2_YCBCR_ENC_SYCC //(absolete)
	YCbCrEncodingBT2020         YCbCrEncodingType = C.V4L2_YCBCR_ENC_BT2020
	YCbCrEncodingBT2020ConstLum YCbCrEncodingType = C.V4L2_YCBCR_ENC_BT2020_CONST_LUM
)

var YCbCrEncodings = map[YCbCrEncodingType]string{
	YCbCrEncodingDefault:        "Default",
	YCbCrEncoding601:            "ITU-R 601",
	YCbCrEncoding709:            "Rec. 709",
	YCbCrEncodingXV601:          "xvYCC 601",
	YCbCrEncodingXV709:          "xvYCC 709",
	YCbCrEncodingBT2020:         "BT.2020",
	YCbCrEncodingBT2020ConstLum: "BT.2020 constant luminance",
	HSVEncoding180:              "HSV 0-179",
	HSVEncoding256:              "HSV 0-255",
}

// ColorspaceToYCbCrEnc is used to get the YCbCrEncoding when only a default YCbCr and the colorspace is known
func ColorspaceToYCbCrEnc(cs ColorspaceType) YCbCrEncodingType {
	switch cs {
	case ColorspaceREC709, ColorspaceDCIP3:
		return YCbCrEncoding709
	case ColorspaceBT2020:
		return YCbCrEncodingBT2020
	default:
		return YCbCrEncoding601
	}
}

// HSVEncodingType (v4l2_hsv_encoding)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L352
type HSVEncodingType = YCbCrEncodingType

const (
	HSVEncoding180 HSVEncodingType = C.V4L2_HSV_ENC_180
	HSVEncoding256 HSVEncodingType = C.V4L2_HSV_ENC_256
)

// QuantizationType (v4l2_quantization)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_quantization#c.V4L.v4l2_quantization
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L372
type QuantizationType = uint32

const (
	QuantizationDefault      QuantizationType = C.V4L2_QUANTIZATION_DEFAULT
	QuantizationFullRange    QuantizationType = C.V4L2_QUANTIZATION_FULL_RANGE
	QuantizationLimitedRange QuantizationType = C.V4L2_QUANTIZATION_LIM_RANGE
)

var Quantizations = map[QuantizationType]string{
	QuantizationDefault:      "Default",
	QuantizationFullRange:    "Full range",
	QuantizationLimitedRange: "Limited range",
}

func ColorspaceToQuantization(cs ColorspaceType) QuantizationType {
	// TODO any RGB/HSV pixel formats should also return full-range
	switch cs {
	case ColorspaceOPRGB, ColorspaceSRGB, ColorspaceJPEG:
		return QuantizationFullRange
	default:
		return QuantizationLimitedRange
	}
}

// XferFunctionType (v4l2_xfer_func)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_xfer_func#c.V4L.v4l2_xfer_func
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L259
type XferFunctionType = uint32

const (
	XferFuncDefault   XferFunctionType = C.V4L2_XFER_FUNC_DEFAULT
	XferFunc709       XferFunctionType = C.V4L2_XFER_FUNC_709
	XferFuncSRGB      XferFunctionType = C.V4L2_XFER_FUNC_SRGB
	XferFuncOpRGB     XferFunctionType = C.V4L2_XFER_FUNC_OPRGB
	XferFuncSMPTE240M XferFunctionType = C.V4L2_XFER_FUNC_SMPTE240M
	XferFuncNone      XferFunctionType = C.V4L2_XFER_FUNC_NONE
	XferFuncDCIP3     XferFunctionType = C.V4L2_XFER_FUNC_DCI_P3
	XferFuncSMPTE2084 XferFunctionType = C.V4L2_XFER_FUNC_SMPTE2084
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

// ColorspaceToXferFunc used to get true XferFunc when only colorspace and default XferFuc are known.
func ColorspaceToXferFunc(cs ColorspaceType) XferFunctionType {
	switch cs {
	case ColorspaceOPRGB:
		return XferFuncOpRGB
	case ColorspaceSMPTE240M:
		return XferFuncSMPTE240M
	case ColorspaceDCIP3:
		return XferFuncDCIP3
	case ColorspaceRaw:
		return XferFuncNone
	case ColorspaceSRGB:
		return XferFuncSRGB
	case ColorspaceJPEG:
		return XferFuncSRGB
	default:
		return XferFunc709
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
