package v4l2

import (
	"strings"
	"testing"
)

// TestPixelFormatConstants tests that all pixel format constants are non-zero
func TestPixelFormatConstants(t *testing.T) {
	formats := []struct {
		name   string
		format FourCCType
	}{
		{"PixelFmtRGB24", PixelFmtRGB24},
		{"PixelFmtGrey", PixelFmtGrey},
		{"PixelFmtYUYV", PixelFmtYUYV},
		{"PixelFmtYYUV", PixelFmtYYUV},
		{"PixelFmtYVYU", PixelFmtYVYU},
		{"PixelFmtUYVY", PixelFmtUYVY},
		{"PixelFmtVYUY", PixelFmtVYUY},
		{"PixelFmtMJPEG", PixelFmtMJPEG},
		{"PixelFmtJPEG", PixelFmtJPEG},
		{"PixelFmtMPEG", PixelFmtMPEG},
		{"PixelFmtH264", PixelFmtH264},
		{"PixelFmtMPEG4", PixelFmtMPEG4},
	}

	for _, tt := range formats {
		t.Run(tt.name, func(t *testing.T) {
			if tt.format == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestPixelFormats_MapComplete tests that PixelFormats map has descriptions
func TestPixelFormats_MapComplete(t *testing.T) {
	// Verify map has entries
	if len(PixelFormats) == 0 {
		t.Error("PixelFormats map should not be empty")
	}

	// Test known formats have descriptions
	tests := []struct {
		format FourCCType
		name   string
	}{
		{PixelFmtRGB24, "RGB24"},
		{PixelFmtGrey, "Grey"},
		{PixelFmtYUYV, "YUYV"},
		{PixelFmtMJPEG, "MJPEG"},
		{PixelFmtH264, "H264"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, exists := PixelFormats[tt.format]
			if !exists {
				t.Errorf("PixelFormats missing entry for %s (0x%08x)", tt.name, tt.format)
			}
			if desc == "" {
				t.Errorf("PixelFormats[%s] description should not be empty", tt.name)
			}
		})
	}
}

// TestIsPixYUVEncoded tests YUV format detection
func TestIsPixYUVEncoded(t *testing.T) {
	tests := []struct {
		name     string
		format   FourCCType
		expected bool
	}{
		{"YUYV is YUV", PixelFmtYUYV, true},
		{"YYUV is YUV", PixelFmtYYUV, true},
		{"YVYU is YUV", PixelFmtYVYU, true},
		{"UYVY is YUV", PixelFmtUYVY, true},
		{"VYUY is YUV", PixelFmtVYUY, true},
		{"RGB24 is not YUV", PixelFmtRGB24, false},
		{"MJPEG is not YUV", PixelFmtMJPEG, false},
		{"H264 is not YUV", PixelFmtH264, false},
		{"Grey is not YUV", PixelFmtGrey, false},
		{"Unknown format is not YUV", FourCCType(0x12345678), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPixYUVEncoded(tt.format)
			if result != tt.expected {
				t.Errorf("IsPixYUVEncoded(0x%08x) = %v, want %v", tt.format, result, tt.expected)
			}
		})
	}
}

// TestColorspaceConstants tests colorspace constant definitions
func TestColorspaceConstants(t *testing.T) {
	colorspaces := []struct {
		name       string
		colorspace ColorspaceType
	}{
		{"ColorspaceDefault", ColorspaceDefault},
		{"ColorspaceSMPTE170M", ColorspaceSMPTE170M},
		{"ColorspaceSMPTE240M", ColorspaceSMPTE240M},
		{"ColorspaceREC709", ColorspaceREC709},
		{"ColorspaceBT878", ColorspaceBT878},
		{"Colorspace470SystemM", Colorspace470SystemM},
		{"Colorspace470SystemBG", Colorspace470SystemBG},
		{"ColorspaceJPEG", ColorspaceJPEG},
		{"ColorspaceSRGB", ColorspaceSRGB},
		{"ColorspaceOPRGB", ColorspaceOPRGB},
		{"ColorspaceBT2020", ColorspaceBT2020},
		{"ColorspaceRaw", ColorspaceRaw},
		{"ColorspaceDCIP3", ColorspaceDCIP3},
	}

	for _, tt := range colorspaces {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify they're defined (they can be zero for default)
			_ = tt.colorspace
		})
	}
}

// TestColorspaces_MapComplete tests Colorspaces map completeness
func TestColorspaces_MapComplete(t *testing.T) {
	if len(Colorspaces) == 0 {
		t.Error("Colorspaces map should not be empty")
	}

	// Verify important colorspaces have descriptions
	tests := []ColorspaceType{
		ColorspaceDefault,
		ColorspaceREC709,
		ColorspaceJPEG,
		ColorspaceSRGB,
		ColorspaceBT2020,
	}

	for _, cs := range tests {
		desc, exists := Colorspaces[cs]
		if !exists {
			t.Errorf("Colorspaces missing entry for %d", cs)
		}
		if desc == "" {
			t.Errorf("Colorspaces[%d] description should not be empty", cs)
		}
	}
}

// TestColorspaceToYCbCrEnc tests colorspace to YCbCr encoding conversion
func TestColorspaceToYCbCrEnc(t *testing.T) {
	tests := []struct {
		name       string
		colorspace ColorspaceType
		expected   YCbCrEncodingType
	}{
		{"REC709 -> 709", ColorspaceREC709, YCbCrEncoding709},
		{"DCIP3 -> 709", ColorspaceDCIP3, YCbCrEncoding709},
		{"BT2020 -> BT2020", ColorspaceBT2020, YCbCrEncodingBT2020},
		{"Default -> 601", ColorspaceDefault, YCbCrEncoding601},
		{"SRGB -> 601", ColorspaceSRGB, YCbCrEncoding601},
		{"JPEG -> 601", ColorspaceJPEG, YCbCrEncoding601},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorspaceToYCbCrEnc(tt.colorspace)
			if result != tt.expected {
				t.Errorf("ColorspaceToYCbCrEnc(%d) = %d, want %d", tt.colorspace, result, tt.expected)
			}
		})
	}
}

// TestYCbCrEncodingConstants tests YCbCr encoding constants
func TestYCbCrEncodingConstants(t *testing.T) {
	encodings := []struct {
		name     string
		encoding YCbCrEncodingType
	}{
		{"YCbCrEncodingDefault", YCbCrEncodingDefault},
		{"YCbCrEncoding601", YCbCrEncoding601},
		{"YCbCrEncoding709", YCbCrEncoding709},
		{"YCbCrEncodingXV601", YCbCrEncodingXV601},
		{"YCbCrEncodingXV709", YCbCrEncodingXV709},
		{"YCbCrEncodingBT2020", YCbCrEncodingBT2020},
		{"YCbCrEncodingBT2020ConstLum", YCbCrEncodingBT2020ConstLum},
	}

	for _, tt := range encodings {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.encoding
		})
	}
}

// TestYCbCrEncodings_MapComplete tests YCbCrEncodings map
func TestYCbCrEncodings_MapComplete(t *testing.T) {
	if len(YCbCrEncodings) == 0 {
		t.Error("YCbCrEncodings map should not be empty")
	}

	// Check important encodings
	tests := []YCbCrEncodingType{
		YCbCrEncodingDefault,
		YCbCrEncoding601,
		YCbCrEncoding709,
		YCbCrEncodingBT2020,
	}

	for _, enc := range tests {
		desc, exists := YCbCrEncodings[enc]
		if !exists {
			t.Errorf("YCbCrEncodings missing entry for %d", enc)
		}
		if desc == "" {
			t.Errorf("YCbCrEncodings[%d] description should not be empty", enc)
		}
	}
}

// TestHSVEncodingConstants tests HSV encoding constants
func TestHSVEncodingConstants(t *testing.T) {
	tests := []struct {
		name     string
		encoding HSVEncodingType
	}{
		{"HSVEncoding180", HSVEncoding180},
		{"HSVEncoding256", HSVEncoding256},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.encoding == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestQuantizationConstants tests quantization constants
func TestQuantizationConstants(t *testing.T) {
	quants := []struct {
		name  string
		quant QuantizationType
	}{
		{"QuantizationDefault", QuantizationDefault},
		{"QuantizationFullRange", QuantizationFullRange},
		{"QuantizationLimitedRange", QuantizationLimitedRange},
	}

	for _, tt := range quants {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.quant
		})
	}
}

// TestQuantizations_MapComplete tests Quantizations map
func TestQuantizations_MapComplete(t *testing.T) {
	if len(Quantizations) == 0 {
		t.Error("Quantizations map should not be empty")
	}

	for quant, desc := range Quantizations {
		if desc == "" {
			t.Errorf("Quantizations[%d] description should not be empty", quant)
		}
	}
}

// TestColorspaceToQuantization tests colorspace to quantization conversion
func TestColorspaceToQuantization(t *testing.T) {
	tests := []struct {
		name       string
		colorspace ColorspaceType
		expected   QuantizationType
	}{
		{"OPRGB -> Full range", ColorspaceOPRGB, QuantizationFullRange},
		{"SRGB -> Full range", ColorspaceSRGB, QuantizationFullRange},
		{"JPEG -> Full range", ColorspaceJPEG, QuantizationFullRange},
		{"Default -> Limited range", ColorspaceDefault, QuantizationLimitedRange},
		{"REC709 -> Limited range", ColorspaceREC709, QuantizationLimitedRange},
		{"BT2020 -> Limited range", ColorspaceBT2020, QuantizationLimitedRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorspaceToQuantization(tt.colorspace)
			if result != tt.expected {
				t.Errorf("ColorspaceToQuantization(%d) = %d, want %d", tt.colorspace, result, tt.expected)
			}
		})
	}
}

// TestXferFunctionConstants tests transfer function constants
func TestXferFunctionConstants(t *testing.T) {
	xfers := []struct {
		name string
		xfer XferFunctionType
	}{
		{"XferFuncDefault", XferFuncDefault},
		{"XferFunc709", XferFunc709},
		{"XferFuncSRGB", XferFuncSRGB},
		{"XferFuncOpRGB", XferFuncOpRGB},
		{"XferFuncSMPTE240M", XferFuncSMPTE240M},
		{"XferFuncNone", XferFuncNone},
		{"XferFuncDCIP3", XferFuncDCIP3},
		{"XferFuncSMPTE2084", XferFuncSMPTE2084},
	}

	for _, tt := range xfers {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.xfer
		})
	}
}

// TestXferFunctions_MapComplete tests XferFunctions map
func TestXferFunctions_MapComplete(t *testing.T) {
	if len(XferFunctions) == 0 {
		t.Error("XferFunctions map should not be empty")
	}

	for xfer, desc := range XferFunctions {
		if desc == "" {
			t.Errorf("XferFunctions[%d] description should not be empty", xfer)
		}
	}
}

// TestColorspaceToXferFunc tests colorspace to transfer function conversion
func TestColorspaceToXferFunc(t *testing.T) {
	tests := []struct {
		name       string
		colorspace ColorspaceType
		expected   XferFunctionType
	}{
		{"OPRGB -> OpRGB", ColorspaceOPRGB, XferFuncOpRGB},
		{"SMPTE240M -> SMPTE240M", ColorspaceSMPTE240M, XferFuncSMPTE240M},
		{"DCIP3 -> DCIP3", ColorspaceDCIP3, XferFuncDCIP3},
		{"Raw -> None", ColorspaceRaw, XferFuncNone},
		{"SRGB -> SRGB", ColorspaceSRGB, XferFuncSRGB},
		{"JPEG -> SRGB", ColorspaceJPEG, XferFuncSRGB},
		{"Default -> 709", ColorspaceDefault, XferFunc709},
		{"REC709 -> 709", ColorspaceREC709, XferFunc709},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorspaceToXferFunc(tt.colorspace)
			if result != tt.expected {
				t.Errorf("ColorspaceToXferFunc(%d) = %d, want %d", tt.colorspace, result, tt.expected)
			}
		})
	}
}

// TestFieldConstants tests field type constants
func TestFieldConstants(t *testing.T) {
	fields := []struct {
		name  string
		field FieldType
	}{
		{"FieldAny", FieldAny},
		{"FieldNone", FieldNone},
		{"FieldTop", FieldTop},
		{"FieldBottom", FieldBottom},
		{"FieldInterlaced", FieldInterlaced},
		{"FieldSequentialTopBottom", FieldSequentialTopBottom},
		{"FieldSequentialBottomTop", FieldSequentialBottomTop},
		{"FieldAlternate", FieldAlternate},
		{"FieldInterlacedTopBottom", FieldInterlacedTopBottom},
		{"FieldInterlacedBottomTop", FieldInterlacedBottomTop},
	}

	for _, tt := range fields {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.field
		})
	}
}

// TestFields_MapComplete tests Fields map
func TestFields_MapComplete(t *testing.T) {
	if len(Fields) == 0 {
		t.Error("Fields map should not be empty")
	}

	for field, desc := range Fields {
		if desc == "" {
			t.Errorf("Fields[%d] description should not be empty", field)
		}
	}
}

// TestPixFormat_String tests the String method
func TestPixFormat_String(t *testing.T) {
	tests := []struct {
		name     string
		format   PixFormat
		contains []string
	}{
		{
			name: "YUYV format",
			format: PixFormat{
				Width:        640,
				Height:       480,
				PixelFormat:  PixelFmtYUYV,
				Field:        FieldNone,
				BytesPerLine: 1280,
				SizeImage:    614400,
				Colorspace:   ColorspaceDefault,
				YcbcrEnc:     YCbCrEncoding601,
				Quantization: QuantizationDefault,
				XferFunc:     XferFuncDefault,
			},
			contains: []string{"640", "480", "YUYV", "1280", "614400"},
		},
		{
			name: "MJPEG format",
			format: PixFormat{
				Width:        1920,
				Height:       1080,
				PixelFormat:  PixelFmtMJPEG,
				Field:        FieldNone,
				BytesPerLine: 0,
				SizeImage:    1048576,
				Colorspace:   ColorspaceJPEG,
				YcbcrEnc:     YCbCrEncodingDefault,
				Quantization: QuantizationFullRange,
				XferFunc:     XferFuncSRGB,
			},
			contains: []string{"1920", "1080", "Motion-JPEG", "1048576"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format.String()

			if result == "" {
				t.Error("String() should not be empty")
			}

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("String() = %q, should contain %q", result, expected)
				}
			}
		})
	}
}

// TestPixFormat_FieldsPresent tests that PixFormat struct has expected fields
func TestPixFormat_FieldsPresent(t *testing.T) {
	format := PixFormat{
		Width:        1920,
		Height:       1080,
		PixelFormat:  PixelFmtYUYV,
		Field:        FieldNone,
		BytesPerLine: 3840,
		SizeImage:    4147200,
		Colorspace:   ColorspaceDefault,
		Priv:         0,
		Flags:        0,
		YcbcrEnc:     YCbCrEncoding601,
		HSVEnc:       HSVEncoding180,
		Quantization: QuantizationDefault,
		XferFunc:     XferFuncDefault,
	}

	// Verify all fields can be accessed
	if format.Width != 1920 {
		t.Errorf("Width = %d, want 1920", format.Width)
	}
	if format.Height != 1080 {
		t.Errorf("Height = %d, want 1080", format.Height)
	}
	if format.PixelFormat != PixelFmtYUYV {
		t.Errorf("PixelFormat = 0x%08x, want 0x%08x", format.PixelFormat, PixelFmtYUYV)
	}
	if format.Field != FieldNone {
		t.Errorf("Field = %d, want %d", format.Field, FieldNone)
	}
	if format.BytesPerLine != 3840 {
		t.Errorf("BytesPerLine = %d, want 3840", format.BytesPerLine)
	}
	if format.SizeImage != 4147200 {
		t.Errorf("SizeImage = %d, want 4147200", format.SizeImage)
	}
}

// TestPixFormat_CommonResolutions tests common video resolutions
func TestPixFormat_CommonResolutions(t *testing.T) {
	resolutions := []struct {
		name   string
		width  uint32
		height uint32
	}{
		{"QVGA", 320, 240},
		{"VGA", 640, 480},
		{"SVGA", 800, 600},
		{"HD 720p", 1280, 720},
		{"Full HD 1080p", 1920, 1080},
		{"4K UHD", 3840, 2160},
	}

	for _, res := range resolutions {
		t.Run(res.name, func(t *testing.T) {
			format := PixFormat{
				Width:       res.width,
				Height:      res.height,
				PixelFormat: PixelFmtYUYV,
				Field:       FieldNone,
			}

			if format.Width != res.width {
				t.Errorf("Width = %d, want %d", format.Width, res.width)
			}
			if format.Height != res.height {
				t.Errorf("Height = %d, want %d", format.Height, res.height)
			}

			// Verify String() works with different resolutions
			str := format.String()
			if str == "" {
				t.Error("String() should not be empty")
			}
		})
	}
}
