package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// FourCCType is a type alias for uint32, representing a Four Character Code (FourCC)
// used to identify pixel formats and other data formats in V4L2.
// Each FourCC is a sequence of four ASCII characters, packed into a 32-bit integer.
type FourCCType = uint32

// Predefined Pixel Format FourCC Constants.
// These constants represent common pixel formats used in video streaming and image capture.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/pixfmt.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L518
var (
	// PixelFmtRGB24 is for 24-bit RGB format (8 bits per R, G, B component).
	PixelFmtRGB24 FourCCType = C.V4L2_PIX_FMT_RGB24
	// PixelFmtGrey is for 8-bit grayscale format.
	PixelFmtGrey FourCCType = C.V4L2_PIX_FMT_GREY
	// PixelFmtYUYV is for YUYV 4:2:2 format (packed YUV).
	PixelFmtYUYV FourCCType = C.V4L2_PIX_FMT_YUYV
	// PixelFmtYYUV is for YYUV 4:2:2 format (packed YUV, alternative to YUYV).
	PixelFmtYYUV FourCCType = C.V4L2_PIX_FMT_YYUV
	// PixelFmtYVYU is for YVYU 4:2:2 format (packed YUV).
	PixelFmtYVYU FourCCType = C.V4L2_PIX_FMT_YVYU
	// PixelFmtUYVY is for UYVY 4:2:2 format (packed YUV).
	PixelFmtUYVY FourCCType = C.V4L2_PIX_FMT_UYVY
	// PixelFmtVYUY is for VYUY 4:2:2 format (packed YUV).
	PixelFmtVYUY FourCCType = C.V4L2_PIX_FMT_VYUY
	// PixelFmtMJPEG is for Motion JPEG format.
	PixelFmtMJPEG FourCCType = C.V4L2_PIX_FMT_MJPEG
	// PixelFmtJPEG is for still JPEG format (JFIF).
	PixelFmtJPEG FourCCType = C.V4L2_PIX_FMT_JPEG
	// PixelFmtMPEG is for MPEG-1/2/4 video elementary streams.
	PixelFmtMPEG FourCCType = C.V4L2_PIX_FMT_MPEG
	// PixelFmtH264 is for H.264 (AVC) video elementary streams.
	PixelFmtH264 FourCCType = C.V4L2_PIX_FMT_H264
	// PixelFmtMPEG4 is for MPEG-4 Part 2 video elementary streams.
	PixelFmtMPEG4 FourCCType = C.V4L2_PIX_FMT_MPEG4
)

// PixelFormats provides a map of common FourCCType constants to their human-readable string descriptions.
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

// IsPixYUVEncoded checks if the given FourCCType pixel format is a YUV (chroma+luminance) format.
// It returns true for common packed YUV formats like YUYV, YYUV, YVYU, UYVY, VYUY.
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

// ColorspaceType is a type alias for uint32, representing the color space of an image or video stream.
// It defines the chromaticity of the red, green, and blue primaries, the white point,
// and the gamma correction function.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L195
// See also https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html
type ColorspaceType = uint32

// Colorspace Type Constants
const (
	ColorspaceDefault ColorspaceType = C.V4L2_COLORSPACE_DEFAULT // Default colorspace, driver picks based on other parameters.
	ColorspaceSMPTE170M ColorspaceType = C.V4L2_COLORSPACE_SMPTE170M // SMPTE 170M colorspace (used for NTSC/PAL SD video).
	ColorspaceSMPTE240M ColorspaceType = C.V4L2_COLORSPACE_SMPTE240M // SMPTE 240M colorspace.
	ColorspaceREC709 ColorspaceType = C.V4L2_COLORSPACE_REC709 // ITU-R BT.709 colorspace (used for HDTV).
	ColorspaceBT878 ColorspaceType = C.V4L2_COLORSPACE_BT878 // Obsolete, do not use.
	Colorspace470SystemM ColorspaceType = C.V4L2_COLORSPACE_470_SYSTEM_M // Obsolete, do not use. (ITU-R BT.470 System M)
	Colorspace470SystemBG ColorspaceType = C.V4L2_COLORSPACE_470_SYSTEM_BG // ITU-R BT.470 System B/G colorspace.
	ColorspaceJPEG ColorspaceType = C.V4L2_COLORSPACE_JPEG // JPEG colorspace (ITU-R BT.601 for YCbCr).
	ColorspaceSRGB ColorspaceType = C.V4L2_COLORSPACE_SRGB // sRGB colorspace.
	ColorspaceOPRGB ColorspaceType = C.V4L2_COLORSPACE_OPRGB // opRGB (Adobe RGB) colorspace.
	ColorspaceBT2020 ColorspaceType = C.V4L2_COLORSPACE_BT2020 // ITU-R BT.2020 colorspace (used for UHDTV).
	ColorspaceRaw ColorspaceType = C.V4L2_COLORSPACE_RAW // Raw sensor data, no specific colorspace.
	ColorspaceDCIP3 ColorspaceType = C.V4L2_COLORSPACE_DCI_P3 // DCI-P3 colorspace (used in digital cinema).
)

// Colorspaces provides a map of common ColorspaceType constants to their human-readable string descriptions.
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

// YCbCrEncodingType is a type alias for uint32, representing the YCbCr encoding scheme.
// It defines how YCbCr color values are derived from RGB values (e.g., ITU-R BT.601, Rec. 709).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_ycbcr_encoding
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L300
type YCbCrEncodingType = uint32

// YCbCr Encoding Type Constants
const (
	YCbCrEncodingDefault YCbCrEncodingType = C.V4L2_YCBCR_ENC_DEFAULT // Default YCbCr encoding, driver picks based on colorspace.
	YCbCrEncoding601 YCbCrEncodingType = C.V4L2_YCBCR_ENC_601 // ITU-R BT.601 encoding (standard definition video).
	YCbCrEncoding709 YCbCrEncodingType = C.V4L2_YCBCR_ENC_709 // ITU-R BT.709 encoding (high definition video).
	YCbCrEncodingXV601 YCbCrEncodingType = C.V4L2_YCBCR_ENC_XV601 // xvYCC extended gamut for BT.601.
	YCbCrEncodingXV709 YCbCrEncodingType = C.V4L2_YCBCR_ENC_XV709 // xvYCC extended gamut for BT.709.
	_ YCbCrEncodingType = C.V4L2_YCBCR_ENC_SYCC // Obsolete (sYCC).
	YCbCrEncodingBT2020 YCbCrEncodingType = C.V4L2_YCBCR_ENC_BT2020 // ITU-R BT.2020 encoding (ultra-high definition video).
	YCbCrEncodingBT2020ConstLum YCbCrEncodingType = C.V4L2_YCBCR_ENC_BT2020_CONST_LUM // ITU-R BT.2020 constant luminance encoding.
)

// YCbCrEncodings provides a map of YCbCrEncodingType constants to their human-readable string descriptions.
// Note: This map also includes HSVEncodingType descriptions as HSVEncodingType is an alias for YCbCrEncodingType.
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

// ColorspaceToYCbCrEnc determines the appropriate YCbCrEncodingType based on a given ColorspaceType.
// This is useful when a YCbCr encoding is not explicitly specified but can be inferred from the colorspace.
// For example, Rec. 709 colorspace typically uses Rec. 709 YCbCr encoding.
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

// HSVEncodingType is an alias for YCbCrEncodingType, representing the encoding range for HSV colorspaces.
// V4L2 reuses the YCbCr encoding enum for HSV, where the values define the range of the Hue component.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L352
type HSVEncodingType = YCbCrEncodingType

// HSV Encoding Type Constants
const (
	HSVEncoding180 HSVEncodingType = C.V4L2_HSV_ENC_180 // Hue component ranges from 0 to 179.
	HSVEncoding256 HSVEncodingType = C.V4L2_HSV_ENC_256 // Hue component ranges from 0 to 255.
)

// QuantizationType is a type alias for uint32, representing the quantization range of color components.
// It specifies whether color values use the full range (e.g., 0-255 for 8-bit) or a limited range
// (e.g., 16-235 for Y, 16-240 for Cb/Cr in 8-bit video).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_quantization#c.V4L.v4l2_quantization
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L372
type QuantizationType = uint32

// Quantization Type Constants
const (
	QuantizationDefault QuantizationType = C.V4L2_QUANTIZATION_DEFAULT // Default quantization, driver picks based on colorspace.
	QuantizationFullRange QuantizationType = C.V4L2_QUANTIZATION_FULL_RANGE // Full range quantization.
	QuantizationLimitedRange QuantizationType = C.V4L2_QUANTIZATION_LIM_RANGE // Limited range quantization.
)

// Quantizations provides a map of QuantizationType constants to their human-readable string descriptions.
var Quantizations = map[QuantizationType]string{
	QuantizationDefault:      "Default",
	QuantizationFullRange:    "Full range",
	QuantizationLimitedRange: "Limited range",
}

// ColorspaceToQuantization determines the appropriate QuantizationType based on a given ColorspaceType.
// Generally, RGB and JPEG colorspaces use full-range quantization, while others might use limited-range.
// TODO: The original comment mentions RGB/HSV formats should also return full-range. This logic might need review/expansion.
func ColorspaceToQuantization(cs ColorspaceType) QuantizationType {
	switch cs {
	case ColorspaceOPRGB, ColorspaceSRGB, ColorspaceJPEG:
		return QuantizationFullRange
	default:
		return QuantizationLimitedRange
	}
}

// XferFunctionType is a type alias for uint32, representing the transfer function (gamma correction) of a colorspace.
// It defines how linear light values are mapped to non-linear (e.g., display) values.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_xfer_func#c.V4L.v4l2_xfer_func
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L259 (kernel uses v4l2_xfer_func enum)
type XferFunctionType = uint32

// Transfer Function Type Constants
const (
	XferFuncDefault XferFunctionType = C.V4L2_XFER_FUNC_DEFAULT // Default transfer function, driver picks based on colorspace.
	XferFunc709 XferFunctionType = C.V4L2_XFER_FUNC_709 // ITU-R BT.709 transfer function.
	XferFuncSRGB XferFunctionType = C.V4L2_XFER_FUNC_SRGB // sRGB transfer function.
	XferFuncOpRGB XferFunctionType = C.V4L2_XFER_FUNC_OPRGB // opRGB transfer function.
	XferFuncSMPTE240M XferFunctionType = C.V4L2_XFER_FUNC_SMPTE240M // SMPTE 240M transfer function.
	XferFuncNone XferFunctionType = C.V4L2_XFER_FUNC_NONE // No transfer function (linear light).
	XferFuncDCIP3 XferFunctionType = C.V4L2_XFER_FUNC_DCI_P3 // DCI-P3 transfer function.
	XferFuncSMPTE2084 XferFunctionType = C.V4L2_XFER_FUNC_SMPTE2084 // SMPTE ST 2084 (HDR PQ) transfer function.
)

// XferFunctions provides a map of XferFunctionType constants to their human-readable string descriptions.
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

// ColorspaceToXferFunc determines the appropriate XferFunctionType based on a given ColorspaceType.
// This is useful for inferring the transfer function when it's not explicitly specified.
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

// FieldType is a type alias for uint32, representing the field order of interlaced video frames.
// It specifies how fields (top or bottom) are arranged in a frame or sequence of frames.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/field-order.html?highlight=v4l2_field#c.v4l2_field
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L88
type FieldType = uint32

// Field Order Type Constants
const (
	FieldAny FieldType = C.V4L2_FIELD_ANY // Driver can choose field order.
	FieldNone FieldType = C.V4L2_FIELD_NONE // Progressive (non-interlaced) frame.
	FieldTop FieldType = C.V4L2_FIELD_TOP // Top field only.
	FieldBottom FieldType = C.V4L2_FIELD_BOTTOM // Bottom field only.
	FieldInterlaced FieldType = C.V4L2_FIELD_INTERLACED // Interlaced frame, top field first.
	FieldSequentialTopBottom FieldType = C.V4L2_FIELD_SEQ_TB // Sequential top and bottom fields.
	FieldSequentialBottomTop FieldType = C.V4L2_FIELD_SEQ_BT // Sequential bottom and top fields.
	FieldAlternate FieldType = C.V4L2_FIELD_ALTERNATE // Alternating top and bottom fields.
	FieldInterlacedTopBottom FieldType = C.V4L2_FIELD_INTERLACED_TB // Interlaced frame, top field followed by bottom field.
	FieldInterlacedBottomTop FieldType = C.V4L2_FIELD_INTERLACED_BT // Interlaced frame, bottom field followed by top field.
)

// Fields provides a map of FieldType constants to their human-readable string descriptions.
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

// PixFormat defines the pixel format for a video stream or image.
// It corresponds to the `v4l2_pix_format` struct in the Linux kernel.
// This struct contains detailed information about the image dimensions, pixel encoding,
// field order, colorspace, and other format-specific parameters.
//
// See https://www.kernel.org/doc/html/v4.9/media/uapi/v4l/pixfmt-002.html?highlight=v4l2_pix_format
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L496
type PixFormat struct {
	// Width of the image in pixels.
	Width uint32
	// Height of the image in pixels.
	Height uint32
	// PixelFormat is the FourCC code identifying the pixel encoding (e.g., V4L2_PIX_FMT_RGB24, V4L2_PIX_FMT_YUYV).
	PixelFormat FourCCType
	// Field specifies the field order for interlaced video (e.g., top field first, progressive). See FieldType constants.
	Field FieldType
	// BytesPerLine is the number of bytes per horizontal line of the image. May include padding.
	BytesPerLine uint32
	// SizeImage is the total size in bytes of the image buffer.
	SizeImage uint32
	// Colorspace defines the color space of the image (e.g., sRGB, Rec. 709). See ColorspaceType constants.
	Colorspace ColorspaceType
	// Priv is a private field for driver-specific use. Applications should ignore it.
	Priv uint32
	// Flags can specify additional format properties (currently none are defined for standard pixel formats).
	Flags uint32
	// YcbcrEnc specifies the YCbCr encoding scheme if applicable. See YCbCrEncodingType constants.
	// This field is part of a union in C, used if PixelFormat is YCbCr.
	YcbcrEnc YCbCrEncodingType
	// HSVEnc specifies the HSV encoding scheme if applicable. See HSVEncodingType constants.
	// This field is part of a union in C, used if PixelFormat is HSV.
	HSVEnc HSVEncodingType // Note: In C, this shares memory with YcbcrEnc via a union.
	// Quantization specifies the quantization range (e.g., full range, limited range). See QuantizationType constants.
	Quantization QuantizationType
	// XferFunc specifies the transfer function (gamma correction). See XferFunctionType constants.
	XferFunc XferFunctionType
}

// String returns a human-readable string representation of the PixFormat struct.
// It includes details like pixel format, dimensions, field order, colorspace, YCbCr encoding,
// quantization, and transfer function.
func (f PixFormat) String() string {
	return fmt.Sprintf(
		"%s [%dx%d]; field=%s; bytes per line=%d; size image=%d; colorspace=%s; YCbCr=%s; Quant=%s; XferFunc=%s",
		PixelFormats[f.PixelFormat], // Assumes PixelFormats map contains the description for f.PixelFormat
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

// GetPixFormat retrieves the current pixel format information for the device's video capture stream.
// It takes the file descriptor of the V4L2 device.
// It returns a PixFormat struct populated with the current format details and an error if the VIDIOC_G_FMT ioctl call fails.
// The `_type` field in the underlying C struct is set to `V4L2_BUF_TYPE_VIDEO_CAPTURE`.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-fmt.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2331 (struct v4l2_format)
func GetPixFormat(fd uintptr) (PixFormat, error) {
	var v4l2Format C.struct_v4l2_format
	v4l2Format._type = C.uint(BufTypeVideoCapture) // Assuming video capture, adjust if other types are needed.

	if err := send(fd, C.VIDIOC_G_FMT, uintptr(unsafe.Pointer(&v4l2Format))); err != nil {
		return PixFormat{}, fmt.Errorf("pix format failed: %w", err)
	}

	// Extract the v4l2_pix_format union member
	v4l2PixFmt := *(*C.struct_v4l2_pix_format)(unsafe.Pointer(&v4l2Format.fmt[0]))
	return PixFormat{
		Width:        uint32(v4l2PixFmt.width),
		Height:       uint32(v4l2PixFmt.height),
		PixelFormat:  FourCCType(v4l2PixFmt.pixelformat),
		Field:        FieldType(v4l2PixFmt.field),
		BytesPerLine: uint32(v4l2PixFmt.bytesperline),
		SizeImage:    uint32(v4l2PixFmt.sizeimage),
		Colorspace:   ColorspaceType(v4l2PixFmt.colorspace),
		Priv:         uint32(v4l2PixFmt.priv),
		Flags:        uint32(v4l2PixFmt.flags),
		// Correctly access union members for YCbCr/HSV encoding.
		// The C struct v4l2_pix_format has a union for ycbcr_enc and hsv_enc.
		// This Go struct has separate fields. Assuming only one is relevant based on colorspace/pixel format.
		// The original code reads both from the same location with an offset for HSV, which might be problematic
		// if the C union isn't structured exactly that way or if only one is valid at a time.
		// For simplicity, this mapping might need adjustment based on how drivers populate this union.
		YcbcrEnc:     YCbCrEncodingType(v4l2PixFmt.ycbcr_enc), // Direct mapping if ycbcr_enc is the active union part
		HSVEnc:       HSVEncodingType(v4l2PixFmt.hsv_enc),     // Direct mapping if hsv_enc is the active union part
		Quantization: QuantizationType(v4l2PixFmt.quantization),
		XferFunc:     XferFunctionType(v4l2PixFmt.xfer_func),
	}, nil
}

// SetPixFormat sets the pixel format information for the device's video capture stream.
// It takes the file descriptor and a PixFormat struct containing the desired format settings.
// The `_type` field in the underlying C struct is set to `V4L2_BUF_TYPE_VIDEO_CAPTURE`.
// Returns an error if the VIDIOC_S_FMT ioctl call fails.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-fmt.html
func SetPixFormat(fd uintptr, pixFmt PixFormat) error {
	var v4l2Format C.struct_v4l2_format
	v4l2Format._type = C.uint(BufTypeVideoCapture) // Assuming video capture
	*(*C.struct_v4l2_pix_format)(unsafe.Pointer(&v4l2Format.fmt[0])) = *(*C.struct_v4l2_pix_format)(unsafe.Pointer(&pixFmt))

	if err := send(fd, C.VIDIOC_S_FMT, uintptr(unsafe.Pointer(&v4l2Format))); err != nil {
		return fmt.Errorf("pix format failed: %w", err)
	}
	return nil
}
