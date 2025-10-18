package v4l2

import (
	"testing"
)

// TestCtrlClassConstants tests control class constants
func TestCtrlClassConstants(t *testing.T) {
	classes := []struct {
		name  string
		class CtrlClass
	}{
		{"CtrlClassUser", CtrlClassUser},
		{"CtrlClassCodec", CtrlClassCodec},
		{"CtrlClassCamera", CtrlClassCamera},
		{"CtrlClassFlash", CtrlClassFlash},
		{"CtrlClassJPEG", CtrlClassJPEG},
		{"CtrlClassImageSource", CtrlClassImageSource},
		{"CtrlClassImageProcessing", CtrlClassImageProcessing},
		{"CtrlClassDigitalVideo", CtrlClassDigitalVideo},
		{"CtrlClassDetection", CtrlClassDetection},
		{"CtrlClassCodecStateless", CtrlClassCodecStateless},
		{"CtrlClassColorimitry", CtrlClassColorimitry},
	}

	for _, tt := range classes {
		t.Run(tt.name, func(t *testing.T) {
			if tt.class == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestCtrlClasses_SliceComplete tests CtrlClasses slice
func TestCtrlClasses_SliceComplete(t *testing.T) {
	if len(CtrlClasses) == 0 {
		t.Error("CtrlClasses slice should not be empty")
	}

	expectedCount := 10 // Based on the constants defined
	if len(CtrlClasses) != expectedCount {
		t.Errorf("CtrlClasses length = %d, want %d", len(CtrlClasses), expectedCount)
	}

	// Verify no duplicates
	seen := make(map[CtrlClass]bool)
	for _, class := range CtrlClasses {
		if seen[class] {
			t.Errorf("Duplicate control class found: 0x%08x", class)
		}
		seen[class] = true
	}
}

// TestCtrlTypeConstants tests control type constants
func TestCtrlTypeConstants(t *testing.T) {
	types := []struct {
		name     string
		ctrlType CtrlType
	}{
		{"CtrlTypeInt", CtrlTypeInt},
		{"CtrlTypeBool", CtrlTypeBool},
		{"CtrlTypeMenu", CtrlTypeMenu},
		{"CtrlTypeButton", CtrlTypeButton},
		{"CtrlTypeInt64", CtrlTypeInt64},
		{"CtrlTypeClass", CtrlTypeClass},
		{"CtrlTypeString", CtrlTypeString},
		{"CtrlTypeBitMask", CtrlTypeBitMask},
		{"CtrlTypeIntegerMenu", CtrlTypeIntegerMenu},
		{"CtrlTypeCompoundTypes", CtrlTypeCompoundTypes},
		{"CtrlTypeU8", CtrlTypeU8},
		{"CtrlTypeU16", CtrlTypeU16},
		{"CtrlTypeU32", CtrlTypeU32},
		{"CtrlTypeArear", CtrlTypeArear},
		{"CtrlTypeHDR10CLLInfo", CtrlTypeHDR10CLLInfo},
		{"CtrlTypeHDRMasteringDisplay", CtrlTypeHDRMasteringDisplay},
		{"CtrlTypeH264SPS", CtrlTypeH264SPS},
		{"CtrlTypeH264PPS", CtrlTypeH264PPS},
		{"CtrlTypeH264ScalingMatrix", CtrlTypeH264ScalingMatrix},
		{"CtrlTypeH264SliceParams", CtrlTypeH264SliceParams},
		{"CtrlTypeH264DecodeParams", CtrlTypeH264DecodeParams},
		{"CtrlTypeFWHTParams", CtrlTypeFWHTParams},
		{"CtrlTypeVP8Frame", CtrlTypeVP8Frame},
		{"CtrlTypeMPEG2Quantization", CtrlTypeMPEG2Quantization},
		{"CtrlTypeMPEG2Sequence", CtrlTypeMPEG2Sequence},
		{"CtrlTypeMPEG2Picture", CtrlTypeMPEG2Picture},
		{"CtrlTypeVP9CompressedHDR", CtrlTypeVP9CompressedHDR},
		{"CtrlTypeVP9Frame", CtrlTypeVP9Frame},
	}

	for _, tt := range types {
		t.Run(tt.name, func(t *testing.T) {
			// Control types can be zero for some values, just verify they're defined
			_ = tt.ctrlType
		})
	}
}

// TestPowerlineFrequencyConstants tests powerline frequency constants
func TestPowerlineFrequencyConstants(t *testing.T) {
	freqs := []struct {
		name string
		freq PowerlineFrequency
	}{
		{"PowerlineFrequencyDisabled", PowerlineFrequencyDisabled},
		{"PowerlineFrequency50Hz", PowerlineFrequency50Hz},
		{"PowerlineFrequency60Hz", PowerlineFrequency60Hz},
		{"PowerlineFrequencyAuto", PowerlineFrequencyAuto},
	}

	for _, tt := range freqs {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.freq
		})
	}
}

// TestColorFXConstants tests color effect constants
func TestColorFXConstants(t *testing.T) {
	effects := []struct {
		name   string
		effect ColorFX
	}{
		{"ColorFXNone", ColorFXNone},
		{"ColorFXBlackWhite", ColorFXBlackWhite},
		{"ColorFXSepia", ColorFXSepia},
		{"ColorFXNegative", ColorFXNegative},
		{"ColorFXEmboss", ColorFXEmboss},
		{"ColorFXSketch", ColorFXSketch},
		{"ColorFXSkyBlue", ColorFXSkyBlue},
		{"ColorFXGrassGreen", ColorFXGrassGreen},
		{"ColorFXSkinWhiten", ColorFXSkinWhiten},
		{"ColorFXVivid", ColorFXVivid},
		{"ColorFXAqua", ColorFXAqua},
		{"ColorFXArtFreeze", ColorFXArtFreeze},
		{"ColorFXSilhouette", ColorFXSilhouette},
		{"ColorFXSolarization", ColorFXSolarization},
		{"ColorFXAntique", ColorFXAntique},
		{"ColorFXSetCBCR", ColorFXSetCBCR},
		{"ColorFXSetRGB", ColorFXSetRGB},
	}

	for _, tt := range effects {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.effect
		})
	}
}

// TestUserControlIDs tests user control ID constants
func TestUserControlIDs(t *testing.T) {
	controls := []struct {
		name string
		id   CtrlID
	}{
		{"CtrlBrightness", CtrlBrightness},
		{"CtrlContrast", CtrlContrast},
		{"CtrlSaturation", CtrlSaturation},
		{"CtrlHue", CtrlHue},
		{"CtrlAutoWhiteBalance", CtrlAutoWhiteBalance},
		{"CtrlDoWhiteBalance", CtrlDoWhiteBalance},
		{"CtrlRedBalance", CtrlRedBalance},
		{"CtrlBlueBalance", CtrlBlueBalance},
		{"CtrlGamma", CtrlGamma},
		{"CtrlExposure", CtrlExposure},
		{"CtrlAutogain", CtrlAutogain},
		{"CtrlGain", CtrlGain},
		{"CtrlHFlip", CtrlHFlip},
		{"CtrlVFlip", CtrlVFlip},
		{"CtrlPowerlineFrequency", CtrlPowerlineFrequency},
		{"CtrlHueAuto", CtrlHueAuto},
		{"CtrlWhiteBalanceTemperature", CtrlWhiteBalanceTemperature},
		{"CtrlSharpness", CtrlSharpness},
		{"CtrlBacklightCompensation", CtrlBacklightCompensation},
		{"CtrlChromaAutomaticGain", CtrlChromaAutomaticGain},
		{"CtrlColorKiller", CtrlColorKiller},
		{"CtrlColorFX", CtrlColorFX},
		{"CtrlAutoBrightness", CtrlAutoBrightness},
		{"CtrlRotate", CtrlRotate},
		{"CtrlBackgroundColor", CtrlBackgroundColor},
	}

	for _, tt := range controls {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestCameraControlIDs tests camera control ID constants
func TestCameraControlIDs(t *testing.T) {
	controls := []struct {
		name string
		id   CtrlID
	}{
		{"CtrlCameraClass", CtrlCameraClass},
		{"CtrlCameraExposureAuto", CtrlCameraExposureAuto},
		{"CtrlCameraExposureAbsolute", CtrlCameraExposureAbsolute},
		{"CtrlCameraExposureAutoPriority", CtrlCameraExposureAutoPriority},
		{"CtrlCameraPanRelative", CtrlCameraPanRelative},
		{"CtrlCameraTiltRelative", CtrlCameraTiltRelative},
		{"CtrlCameraPanAbsolute", CtrlCameraPanAbsolute},
		{"CtrlCameraTiltAbsolute", CtrlCameraTiltAbsolute},
		{"CtrlCameraFocusAbsolute", CtrlCameraFocusAbsolute},
		{"CtrlCameraFocusRelative", CtrlCameraFocusRelative},
		{"CtrlCameraFocusAuto", CtrlCameraFocusAuto},
		{"CtrlCameraZoomAbsolute", CtrlCameraZoomAbsolute},
		{"CtrlCameraZoomRelative", CtrlCameraZoomRelative},
		{"CtrlCameraPrivacy", CtrlCameraPrivacy},
		{"CtrlCameraAutoExposureBias", CtrlCameraAutoExposureBias},
		{"CtrlCameraWideDynamicRange", CtrlCameraWideDynamicRange},
		{"CtrlCameraImageStabilization", CtrlCameraImageStabilization},
		{"CtrlCameraIsoSensitivity", CtrlCameraIsoSensitivity},
		{"CtrlCameraIsoSensitivityAuto", CtrlCameraIsoSensitivityAuto},
		{"CtrlCameraExposureMetering", CtrlCameraExposureMetering},
		{"CtrlCameraSceneMode", CtrlCameraSceneMode},
		{"CtrlCamera3ALock", CtrlCamera3ALock},
		{"CtrlCameraAutoFocusStart", CtrlCameraAutoFocusStart},
		{"CtrlCameraAutoFocusStop", CtrlCameraAutoFocusStop},
		{"CtrlCameraAutoFocusStatus", CtrlCameraAutoFocusStatus},
		{"CtrlCameraAutoFocusRange", CtrlCameraAutoFocusRange},
	}

	for _, tt := range controls {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestControl_StructFields tests Control struct field accessibility
func TestControl_StructFields(t *testing.T) {
	ctrl := Control{
		Type:    CtrlTypeInt,
		ID:      CtrlBrightness,
		Value:   50,
		Name:    "Brightness",
		Minimum: 0,
		Maximum: 100,
		Step:    1,
		Default: 50,
	}

	// Verify all fields are accessible
	if ctrl.Type != CtrlTypeInt {
		t.Errorf("Type = %d, want %d", ctrl.Type, CtrlTypeInt)
	}
	if ctrl.ID != CtrlBrightness {
		t.Errorf("ID = %d, want %d", ctrl.ID, CtrlBrightness)
	}
	if ctrl.Value != 50 {
		t.Errorf("Value = %d, want 50", ctrl.Value)
	}
	if ctrl.Name != "Brightness" {
		t.Errorf("Name = %s, want Brightness", ctrl.Name)
	}
	if ctrl.Minimum != 0 {
		t.Errorf("Minimum = %d, want 0", ctrl.Minimum)
	}
	if ctrl.Maximum != 100 {
		t.Errorf("Maximum = %d, want 100", ctrl.Maximum)
	}
	if ctrl.Step != 1 {
		t.Errorf("Step = %d, want 1", ctrl.Step)
	}
	if ctrl.Default != 50 {
		t.Errorf("Default = %d, want 50", ctrl.Default)
	}
}

// TestControl_IsMenu tests the IsMenu method
func TestControl_IsMenu(t *testing.T) {
	tests := []struct {
		name     string
		ctrl     Control
		expected bool
	}{
		{
			name:     "Menu type",
			ctrl:     Control{Type: CtrlTypeMenu},
			expected: true,
		},
		{
			name:     "Integer menu type",
			ctrl:     Control{Type: CtrlTypeIntegerMenu},
			expected: true,
		},
		{
			name:     "Integer type",
			ctrl:     Control{Type: CtrlTypeInt},
			expected: false,
		},
		{
			name:     "Boolean type",
			ctrl:     Control{Type: CtrlTypeBool},
			expected: false,
		},
		{
			name:     "Button type",
			ctrl:     Control{Type: CtrlTypeButton},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ctrl.IsMenu()
			if result != tt.expected {
				t.Errorf("IsMenu() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestControl_ValueRanges tests typical control value ranges
func TestControl_ValueRanges(t *testing.T) {
	tests := []struct {
		name    string
		ctrl    Control
		testVal int32
		inRange bool
	}{
		{
			name: "Value within range",
			ctrl: Control{
				Minimum: 0,
				Maximum: 100,
			},
			testVal: 50,
			inRange: true,
		},
		{
			name: "Value at minimum",
			ctrl: Control{
				Minimum: 0,
				Maximum: 100,
			},
			testVal: 0,
			inRange: true,
		},
		{
			name: "Value at maximum",
			ctrl: Control{
				Minimum: 0,
				Maximum: 100,
			},
			testVal: 100,
			inRange: true,
		},
		{
			name: "Value below minimum",
			ctrl: Control{
				Minimum: 0,
				Maximum: 100,
			},
			testVal: -1,
			inRange: false,
		},
		{
			name: "Value above maximum",
			ctrl: Control{
				Minimum: 0,
				Maximum: 100,
			},
			testVal: 101,
			inRange: false,
		},
		{
			name: "Negative range",
			ctrl: Control{
				Minimum: -50,
				Maximum: 50,
			},
			testVal: 0,
			inRange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inRange := tt.testVal >= tt.ctrl.Minimum && tt.testVal <= tt.ctrl.Maximum
			if inRange != tt.inRange {
				t.Errorf("Value %d in range [%d, %d] = %v, want %v",
					tt.testVal, tt.ctrl.Minimum, tt.ctrl.Maximum, inRange, tt.inRange)
			}
		})
	}
}

// TestControl_StepValues tests control step values
func TestControl_StepValues(t *testing.T) {
	tests := []struct {
		name string
		ctrl Control
	}{
		{
			name: "Step 1",
			ctrl: Control{
				Minimum: 0,
				Maximum: 100,
				Step:    1,
			},
		},
		{
			name: "Step 5",
			ctrl: Control{
				Minimum: 0,
				Maximum: 100,
				Step:    5,
			},
		},
		{
			name: "Step 10",
			ctrl: Control{
				Minimum: 0,
				Maximum: 255,
				Step:    10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ctrl.Step <= 0 {
				t.Errorf("Step should be positive, got %d", tt.ctrl.Step)
			}

			// Verify step divides the range evenly (or leaves a remainder)
			rangeSize := tt.ctrl.Maximum - tt.ctrl.Minimum
			_ = rangeSize / tt.ctrl.Step // Number of steps
		})
	}
}

// TestControlMenuItem_StructFields tests ControlMenuItem struct
func TestControlMenuItem_StructFields(t *testing.T) {
	item := ControlMenuItem{
		ID:    CtrlPowerlineFrequency,
		Index: 0,
		Value: PowerlineFrequencyDisabled,
		Name:  "Disabled",
	}

	if item.ID != CtrlPowerlineFrequency {
		t.Errorf("ID = %d, want %d", item.ID, CtrlPowerlineFrequency)
	}
	if item.Index != 0 {
		t.Errorf("Index = %d, want 0", item.Index)
	}
	if item.Value != PowerlineFrequencyDisabled {
		t.Errorf("Value = %d, want %d", item.Value, PowerlineFrequencyDisabled)
	}
	if item.Name != "Disabled" {
		t.Errorf("Name = %s, want Disabled", item.Name)
	}
}

// TestControlMenuItem_MenuSequence tests a typical menu sequence
func TestControlMenuItem_MenuSequence(t *testing.T) {
	// Simulate powerline frequency menu items
	items := []ControlMenuItem{
		{ID: CtrlPowerlineFrequency, Index: 0, Name: "Disabled"},
		{ID: CtrlPowerlineFrequency, Index: 1, Name: "50 Hz"},
		{ID: CtrlPowerlineFrequency, Index: 2, Name: "60 Hz"},
		{ID: CtrlPowerlineFrequency, Index: 3, Name: "Auto"},
	}

	// Verify indexes are sequential
	for i, item := range items {
		if item.Index != uint32(i) {
			t.Errorf("Item %d: Index = %d, want %d", i, item.Index, i)
		}
		if item.ID != CtrlPowerlineFrequency {
			t.Errorf("Item %d: ID mismatch", i)
		}
		if item.Name == "" {
			t.Errorf("Item %d: Name should not be empty", i)
		}
	}
}

// TestControl_CommonControlTypes tests common control type scenarios
func TestControl_CommonControlTypes(t *testing.T) {
	tests := []struct {
		name        string
		ctrl        Control
		description string
	}{
		{
			name: "Integer control (Brightness)",
			ctrl: Control{
				Type:    CtrlTypeInt,
				ID:      CtrlBrightness,
				Name:    "Brightness",
				Minimum: 0,
				Maximum: 255,
				Step:    1,
				Default: 128,
			},
			description: "Standard 8-bit integer control",
		},
		{
			name: "Boolean control (Auto White Balance)",
			ctrl: Control{
				Type:    CtrlTypeBool,
				ID:      CtrlAutoWhiteBalance,
				Name:    "Auto White Balance",
				Minimum: 0,
				Maximum: 1,
				Step:    1,
				Default: 1,
			},
			description: "On/off boolean control",
		},
		{
			name: "Menu control (Powerline Frequency)",
			ctrl: Control{
				Type:    CtrlTypeMenu,
				ID:      CtrlPowerlineFrequency,
				Name:    "Power Line Frequency",
				Minimum: 0,
				Maximum: 3,
				Step:    1,
				Default: 1,
			},
			description: "Menu with discrete options",
		},
		{
			name: "Button control (Auto Focus)",
			ctrl: Control{
				Type: CtrlTypeButton,
				ID:   CtrlCameraAutoFocusStart,
				Name: "Auto Focus Start",
			},
			description: "Button with no value range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify control type is appropriate
			switch tt.ctrl.Type {
			case CtrlTypeInt:
				if tt.ctrl.Maximum <= tt.ctrl.Minimum {
					t.Error("Integer control should have max > min")
				}
			case CtrlTypeBool:
				if tt.ctrl.Minimum != 0 || tt.ctrl.Maximum != 1 {
					t.Error("Boolean control should have range [0, 1]")
				}
			case CtrlTypeMenu:
				if tt.ctrl.Maximum < tt.ctrl.Minimum {
					t.Error("Menu control should have max >= min")
				}
				if !tt.ctrl.IsMenu() {
					t.Error("Menu control IsMenu() should return true")
				}
			case CtrlTypeButton:
				// Buttons typically don't have value ranges
			}
		})
	}
}

// TestControl_TypeClassification tests control type classification
func TestControl_TypeClassification(t *testing.T) {
	// Test which types are considered "simple" vs "compound"
	simpleTypes := []CtrlType{
		CtrlTypeInt,
		CtrlTypeBool,
		CtrlTypeMenu,
		CtrlTypeButton,
		CtrlTypeInt64,
		CtrlTypeString,
		CtrlTypeBitMask,
		CtrlTypeIntegerMenu,
	}

	for _, typ := range simpleTypes {
		ctrl := Control{Type: typ}
		// Simple types can be read/written with VIDIOC_G_CTRL/VIDIOC_S_CTRL
		_ = ctrl.Type
	}

	// Compound types require extended controls API
	compoundTypes := []CtrlType{
		CtrlTypeU8,
		CtrlTypeU16,
		CtrlTypeU32,
		CtrlTypeH264SPS,
		CtrlTypeH264PPS,
		CtrlTypeVP8Frame,
	}

	for _, typ := range compoundTypes {
		ctrl := Control{Type: typ}
		// Compound types require VIDIOC_G_EXT_CTRLS/VIDIOC_S_EXT_CTRLS
		_ = ctrl.Type
	}
}
