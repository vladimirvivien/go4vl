package v4l2

/*
#include <linux/videodev2.h>
*/
import "C"

// ColorspaceType (v4l2_colorspace)
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces.html
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L195
type ColorspaceType = uint32

const (
	ColorspaceDefault     ColorspaceType = C.V4L2_COLORSPACE_DEFAULT      // Default colorspace (driver-dependent)
	ColorspaceSMPTE170M   ColorspaceType = C.V4L2_COLORSPACE_SMPTE170M    // SMPTE 170M (NTSC/PAL/SECAM)
	ColorspaceSMPTE240M   ColorspaceType = C.V4L2_COLORSPACE_SMPTE240M    // SMPTE 240M (obsolete HDTV)
	ColorspaceREC709      ColorspaceType = C.V4L2_COLORSPACE_REC709       // Rec. 709 (HDTV)
	ColorspaceBT878       ColorspaceType = C.V4L2_COLORSPACE_BT878        // BT.878 (obsolete, same as SMPTE170M)
	Colorspace470SystemM  ColorspaceType = C.V4L2_COLORSPACE_470_SYSTEM_M // 470 System M (obsolete NTSC)
	Colorspace470SystemBG ColorspaceType = C.V4L2_COLORSPACE_470_SYSTEM_BG // 470 System BG (obsolete PAL/SECAM)
	ColorspaceJPEG        ColorspaceType = C.V4L2_COLORSPACE_JPEG         // JPEG/sYCC colorspace
	ColorspaceSRGB        ColorspaceType = C.V4L2_COLORSPACE_SRGB         // sRGB colorspace
	ColorspaceOPRGB       ColorspaceType = C.V4L2_COLORSPACE_OPRGB        // opRGB colorspace
	ColorspaceBT2020      ColorspaceType = C.V4L2_COLORSPACE_BT2020       // BT.2020 (UHDTV)
	ColorspaceRaw         ColorspaceType = C.V4L2_COLORSPACE_RAW          // Raw (no colorspace conversion)
	ColorspaceDCIP3       ColorspaceType = C.V4L2_COLORSPACE_DCI_P3       // DCI-P3 (digital cinema)
)

// Colorspaces is a map of colorspace to its respective description
var Colorspaces = map[ColorspaceType]string{
	ColorspaceDefault:     "Default",
	ColorspaceSMPTE170M:   "SMPTE 170M",
	ColorspaceSMPTE240M:   "SMPTE 240M",
	ColorspaceREC709:      "Rec. 709",
	ColorspaceBT878:       "BT.878",
	Colorspace470SystemM:  "470 System M",
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
	YCbCrEncodingDefault        YCbCrEncodingType = C.V4L2_YCBCR_ENC_DEFAULT           // Default YCbCr encoding (driver-dependent)
	YCbCrEncoding601            YCbCrEncodingType = C.V4L2_YCBCR_ENC_601               // ITU-R BT.601 (SDTV)
	YCbCrEncoding709            YCbCrEncodingType = C.V4L2_YCBCR_ENC_709               // Rec. 709 (HDTV)
	YCbCrEncodingXV601          YCbCrEncodingType = C.V4L2_YCBCR_ENC_XV601             // xvYCC 601 (extended gamut BT.601)
	YCbCrEncodingXV709          YCbCrEncodingType = C.V4L2_YCBCR_ENC_XV709             // xvYCC 709 (extended gamut Rec.709)
	YCbCrEncodingSYCC           YCbCrEncodingType = C.V4L2_YCBCR_ENC_SYCC              // sYCC (obsolete, same as XV601)
	YCbCrEncodingBT2020         YCbCrEncodingType = C.V4L2_YCBCR_ENC_BT2020            // BT.2020 (UHDTV)
	YCbCrEncodingBT2020ConstLum YCbCrEncodingType = C.V4L2_YCBCR_ENC_BT2020_CONST_LUM // BT.2020 constant luminance
	YCbCrEncodingSMPTE240M      YCbCrEncodingType = C.V4L2_YCBCR_ENC_SMPTE240M        // SMPTE 240M (obsolete HDTV)
)

// YCbCrEncodings is a map of YCbCr encoding to its description
var YCbCrEncodings = map[YCbCrEncodingType]string{
	YCbCrEncodingDefault:        "Default",
	YCbCrEncoding601:            "ITU-R BT.601",
	YCbCrEncoding709:            "Rec. 709",
	YCbCrEncodingXV601:          "xvYCC 601",
	YCbCrEncodingXV709:          "xvYCC 709",
	YCbCrEncodingSYCC:           "sYCC",
	YCbCrEncodingBT2020:         "BT.2020",
	YCbCrEncodingBT2020ConstLum: "BT.2020 constant luminance",
	YCbCrEncodingSMPTE240M:      "SMPTE 240M",
	HSVEncoding180:              "HSV 0-179",
	HSVEncoding256:              "HSV 0-255",
}

// ColorspaceToYCbCrEnc returns the appropriate YCbCr encoding for a given colorspace
// when only a default YCbCr encoding and the colorspace is known
func ColorspaceToYCbCrEnc(cs ColorspaceType) YCbCrEncodingType {
	switch cs {
	case ColorspaceREC709, ColorspaceDCIP3:
		return YCbCrEncoding709
	case ColorspaceBT2020:
		return YCbCrEncodingBT2020
	case ColorspaceSMPTE240M:
		return YCbCrEncodingSMPTE240M
	default:
		return YCbCrEncoding601
	}
}

// HSVEncodingType (v4l2_hsv_encoding)
// HSV encoding shares the same type space as YCbCr encoding
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L352
type HSVEncodingType = YCbCrEncodingType

const (
	HSVEncoding180 HSVEncodingType = C.V4L2_HSV_ENC_180 // Hue mapped to 0-179
	HSVEncoding256 HSVEncodingType = C.V4L2_HSV_ENC_256 // Hue mapped to 0-255
)

// QuantizationType (v4l2_quantization)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_quantization#c.V4L.v4l2_quantization
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L372
type QuantizationType = uint32

const (
	QuantizationDefault      QuantizationType = C.V4L2_QUANTIZATION_DEFAULT   // Default quantization (driver-dependent)
	QuantizationFullRange    QuantizationType = C.V4L2_QUANTIZATION_FULL_RANGE // Full range (0-255 for 8-bit)
	QuantizationLimitedRange QuantizationType = C.V4L2_QUANTIZATION_LIM_RANGE  // Limited range (16-235 for 8-bit Y, 16-240 for Cb/Cr)
)

// Quantizations is a map of quantization type to its description
var Quantizations = map[QuantizationType]string{
	QuantizationDefault:      "Default",
	QuantizationFullRange:    "Full range",
	QuantizationLimitedRange: "Limited range",
}

// ColorspaceToQuantization returns the appropriate quantization for a given colorspace
func ColorspaceToQuantization(cs ColorspaceType) QuantizationType {
	// RGB and HSV pixel formats use full-range quantization
	switch cs {
	case ColorspaceOPRGB, ColorspaceSRGB, ColorspaceJPEG:
		return QuantizationFullRange
	default:
		return QuantizationLimitedRange
	}
}

// XferFunctionType (v4l2_xfer_func)
// Transfer function (gamma curve / EOTF)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_xfer_func#c.V4L.v4l2_xfer_func
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L259
type XferFunctionType = uint32

const (
	XferFuncDefault   XferFunctionType = C.V4L2_XFER_FUNC_DEFAULT   // Default transfer function (driver-dependent)
	XferFunc709       XferFunctionType = C.V4L2_XFER_FUNC_709       // Rec. 709 transfer function
	XferFuncSRGB      XferFunctionType = C.V4L2_XFER_FUNC_SRGB      // sRGB transfer function
	XferFuncOpRGB     XferFunctionType = C.V4L2_XFER_FUNC_OPRGB     // opRGB transfer function
	XferFuncSMPTE240M XferFunctionType = C.V4L2_XFER_FUNC_SMPTE240M // SMPTE 240M transfer function
	XferFuncNone      XferFunctionType = C.V4L2_XFER_FUNC_NONE      // No transfer function (linear)
	XferFuncDCIP3     XferFunctionType = C.V4L2_XFER_FUNC_DCI_P3    // DCI-P3 transfer function
	XferFuncSMPTE2084 XferFunctionType = C.V4L2_XFER_FUNC_SMPTE2084 // SMPTE 2084 (PQ - Perceptual Quantizer for HDR)
)

// XferFunctions is a map of transfer function type to its description
var XferFunctions = map[XferFunctionType]string{
	XferFuncDefault:   "Default",
	XferFunc709:       "Rec. 709",
	XferFuncSRGB:      "sRGB",
	XferFuncOpRGB:     "opRGB",
	XferFuncSMPTE240M: "SMPTE 240M",
	XferFuncNone:      "None (linear)",
	XferFuncDCIP3:     "DCI-P3",
	XferFuncSMPTE2084: "SMPTE 2084 (PQ)",
}

// ColorspaceToXferFunc returns the appropriate transfer function for a given colorspace
// when only the colorspace and default transfer function are known
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
	case ColorspaceSRGB, ColorspaceJPEG:
		return XferFuncSRGB
	default:
		return XferFunc709
	}
}

// ColorspaceInfo holds complete colorspace information for a format
type ColorspaceInfo struct {
	Colorspace   ColorspaceType
	YCbCrEnc     YCbCrEncodingType
	Quantization QuantizationType
	XferFunc     XferFunctionType
}

// String returns a human-readable description of the colorspace information
func (c ColorspaceInfo) String() string {
	cs := Colorspaces[c.Colorspace]
	ycbcr := YCbCrEncodings[c.YCbCrEnc]
	quant := Quantizations[c.Quantization]
	xfer := XferFunctions[c.XferFunc]
	return cs + " / " + ycbcr + " / " + quant + " / " + xfer
}

// NewColorspaceInfo creates a ColorspaceInfo with defaults filled in based on colorspace
func NewColorspaceInfo(cs ColorspaceType) ColorspaceInfo {
	return ColorspaceInfo{
		Colorspace:   cs,
		YCbCrEnc:     ColorspaceToYCbCrEnc(cs),
		Quantization: ColorspaceToQuantization(cs),
		XferFunc:     ColorspaceToXferFunc(cs),
	}
}

// IsHDR returns true if the colorspace information represents HDR content
func (c ColorspaceInfo) IsHDR() bool {
	return c.XferFunc == XferFuncSMPTE2084 || c.Colorspace == ColorspaceBT2020
}

// IsSDR returns true if the colorspace information represents SDR content
func (c ColorspaceInfo) IsSDR() bool {
	return !c.IsHDR()
}
