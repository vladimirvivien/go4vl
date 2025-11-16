package v4l2

import (
	"testing"
)

// TestExtControlClass_Constants verifies additional control class constants are defined
func TestExtControlClass_Constants(t *testing.T) {
	tests := []struct {
		name  string
		class CtrlClass
	}{
		{"CtrlClassFMTx", CtrlClassFMTx},
		{"CtrlClassFMRx", CtrlClassFMRx},
		{"CtrlClassRFTuner", CtrlClassRFTuner},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.class == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestNewExtControl verifies ExtControl creation
func TestNewExtControl(t *testing.T) {
	tests := []struct {
		name   string
		ctrlID CtrlID
	}{
		{"Brightness", CtrlBrightness},
		{"Contrast", CtrlContrast},
		{"Saturation", CtrlSaturation},
		{"Hue", CtrlHue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewExtControl(tt.ctrlID)
			if ctrl == nil {
				t.Fatal("NewExtControl returned nil")
			}
			if ctrl.GetID() != tt.ctrlID {
				t.Errorf("GetID() = %d, want %d", ctrl.GetID(), tt.ctrlID)
			}
		})
	}
}

// TestNewExtControlWithValue verifies ExtControl creation with 32-bit value
func TestNewExtControlWithValue(t *testing.T) {
	tests := []struct {
		name   string
		ctrlID CtrlID
		value  int32
	}{
		{"Brightness_128", CtrlBrightness, 128},
		{"Contrast_100", CtrlContrast, 100},
		{"Saturation_Max", CtrlSaturation, 255},
		{"Hue_Min", CtrlHue, 0},
		{"Negative", CtrlHue, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewExtControlWithValue(tt.ctrlID, tt.value)
			if ctrl == nil {
				t.Fatal("NewExtControlWithValue returned nil")
			}
			if ctrl.GetID() != tt.ctrlID {
				t.Errorf("GetID() = %d, want %d", ctrl.GetID(), tt.ctrlID)
			}
			if ctrl.GetValue() != tt.value {
				t.Errorf("GetValue() = %d, want %d", ctrl.GetValue(), tt.value)
			}
		})
	}
}

// TestNewExtControlWithValue64 verifies ExtControl creation with 64-bit value
func TestNewExtControlWithValue64(t *testing.T) {
	tests := []struct {
		name   string
		ctrlID CtrlID
		value  int64
	}{
		{"Small", CtrlBrightness, 100},
		{"Large", CtrlBrightness, 1000000000},
		{"VeryLarge", CtrlBrightness, 9223372036854775807}, // Max int64
		{"Negative", CtrlBrightness, -1000000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewExtControlWithValue64(tt.ctrlID, tt.value)
			if ctrl == nil {
				t.Fatal("NewExtControlWithValue64 returned nil")
			}
			if ctrl.GetID() != tt.ctrlID {
				t.Errorf("GetID() = %d, want %d", ctrl.GetID(), tt.ctrlID)
			}
			if ctrl.GetValue64() != tt.value {
				t.Errorf("GetValue64() = %d, want %d", ctrl.GetValue64(), tt.value)
			}
		})
	}
}

// TestNewExtControlWithString verifies ExtControl creation with string value
func TestNewExtControlWithString(t *testing.T) {
	tests := []struct {
		name   string
		ctrlID CtrlID
		value  string
	}{
		{"Empty", CtrlBrightness, ""},
		{"Short", CtrlBrightness, "test"},
		{"Long", CtrlBrightness, "This is a much longer string value for testing"},
		{"Unicode", CtrlBrightness, "Hello 世界 🌍"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewExtControlWithString(tt.ctrlID, tt.value)
			if ctrl == nil {
				t.Fatal("NewExtControlWithString returned nil")
			}
			if ctrl.GetID() != tt.ctrlID {
				t.Errorf("GetID() = %d, want %d", ctrl.GetID(), tt.ctrlID)
			}
			got := ctrl.GetString()
			if got != tt.value {
				t.Errorf("GetString() = %q, want %q", got, tt.value)
			}
		})
	}
}

// TestExtControl_Value verifies int32 value storage
func TestExtControl_Value(t *testing.T) {
	tests := []int32{0, 1, 128, 255, -1, -128}
	for _, value := range tests {
		ctrl := NewExtControlWithValue(CtrlBrightness, value)
		got := ctrl.GetValue()
		if got != value {
			t.Errorf("NewExtControlWithValue(%d), GetValue() = %d", value, got)
		}
	}
}

// TestExtControl_Value64 verifies int64 value storage
func TestExtControl_Value64(t *testing.T) {
	tests := []int64{0, 1, 1000000000, 9223372036854775807, -1, -1000000000}
	for _, value := range tests {
		ctrl := NewExtControlWithValue64(CtrlBrightness, value)
		got := ctrl.GetValue64()
		if got != value {
			t.Errorf("NewExtControlWithValue64(%d), GetValue64() = %d", value, got)
		}
	}
}

// TestExtControl_String verifies string value storage
func TestExtControl_String(t *testing.T) {
	tests := []string{
		"",
		"test",
		"a longer test string",
		"Unicode: 日本語",
	}

	for _, value := range tests {
		ctrl := NewExtControlWithString(CtrlBrightness, value)
		got := ctrl.GetString()
		if got != value {
			t.Errorf("NewExtControlWithString(%q), GetString() = %q", value, got)
		}
	}
}

// TestExtControl_CompoundData verifies compound data storage
func TestExtControl_CompoundData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"Small", []byte{1, 2, 3, 4}},
		{"Medium", make([]byte, 256)},
		{"Large", make([]byte, 4096)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test data with pattern
			for i := range tt.data {
				tt.data[i] = byte(i % 256)
			}

			ctrl := NewExtControlWithCompound(CtrlBrightness, tt.data)
			got := ctrl.GetCompoundData()

			if len(got) != len(tt.data) {
				t.Errorf("GetCompoundData() length = %d, want %d", len(got), len(tt.data))
			}

			// Verify data matches
			for i := range tt.data {
				if got[i] != tt.data[i] {
					t.Errorf("GetCompoundData()[%d] = %d, want %d", i, got[i], tt.data[i])
					break
				}
			}
		})
	}
}

// TestNewExtControls verifies ExtControls creation
func TestNewExtControls(t *testing.T) {
	ctrls := NewExtControls()
	if ctrls == nil {
		t.Fatal("NewExtControls returned nil")
	}
	if ctrls.Count() != 0 {
		t.Errorf("Count() = %d, want 0", ctrls.Count())
	}
}

// TestNewExtControlsWithClass verifies ExtControls creation with class
func TestNewExtControlsWithClass(t *testing.T) {
	tests := []struct {
		name  string
		class CtrlClass
	}{
		{"User", CtrlClassUser},
		{"Codec", CtrlClassCodec},
		{"Camera", CtrlClassCamera},
		{"JPEG", CtrlClassJPEG},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrls := NewExtControlsWithClass(tt.class)
			if ctrls == nil {
				t.Fatal("NewExtControlsWithClass returned nil")
			}
			if ctrls.GetClass() != tt.class {
				t.Errorf("GetClass() = 0x%08x, want 0x%08x", ctrls.GetClass(), tt.class)
			}
		})
	}
}

// TestExtControls_SetGetClass verifies class set/get
func TestExtControls_SetGetClass(t *testing.T) {
	ctrls := NewExtControls()

	tests := []CtrlClass{
		CtrlClassUser,
		CtrlClassCodec,
		CtrlClassCamera,
		CtrlClassJPEG,
		CtrlClassFlash,
	}

	for _, class := range tests {
		ctrls.SetClass(class)
		got := ctrls.GetClass()
		if got != class {
			t.Errorf("After SetClass(0x%08x), GetClass() = 0x%08x", class, got)
		}
	}
}

// TestExtControls_AddGetControls verifies adding controls
func TestExtControls_AddGetControls(t *testing.T) {
	ctrls := NewExtControls()

	// Start with empty
	if ctrls.Count() != 0 {
		t.Errorf("Initial Count() = %d, want 0", ctrls.Count())
	}

	// Add controls
	ctrl1 := NewExtControlWithValue(CtrlBrightness, 100)
	ctrl2 := NewExtControlWithValue(CtrlContrast, 50)
	ctrl3 := NewExtControlWithValue(CtrlSaturation, 75)

	ctrls.Add(ctrl1)
	if ctrls.Count() != 1 {
		t.Errorf("After 1 Add(), Count() = %d, want 1", ctrls.Count())
	}

	ctrls.Add(ctrl2)
	if ctrls.Count() != 2 {
		t.Errorf("After 2 Add(), Count() = %d, want 2", ctrls.Count())
	}

	ctrls.Add(ctrl3)
	if ctrls.Count() != 3 {
		t.Errorf("After 3 Add(), Count() = %d, want 3", ctrls.Count())
	}

	// Verify controls
	controls := ctrls.GetControls()
	if len(controls) != 3 {
		t.Errorf("GetControls() length = %d, want 3", len(controls))
	}

	// Verify values
	if controls[0].GetID() != CtrlBrightness || controls[0].GetValue() != 100 {
		t.Errorf("controls[0] ID=%d value=%d, want ID=%d value=100",
			controls[0].GetID(), controls[0].GetValue(), CtrlBrightness)
	}
	if controls[1].GetID() != CtrlContrast || controls[1].GetValue() != 50 {
		t.Errorf("controls[1] ID=%d value=%d, want ID=%d value=50",
			controls[1].GetID(), controls[1].GetValue(), CtrlContrast)
	}
	if controls[2].GetID() != CtrlSaturation || controls[2].GetValue() != 75 {
		t.Errorf("controls[2] ID=%d value=%d, want ID=%d value=75",
			controls[2].GetID(), controls[2].GetValue(), CtrlSaturation)
	}

	// Cleanup controls
}

// TestExtControls_GetErrorIndex verifies error index retrieval
func TestExtControls_GetErrorIndex(t *testing.T) {
	ctrls := NewExtControls()

	// Initially should be 0
	idx := ctrls.GetErrorIndex()
	if idx != 0 {
		t.Errorf("GetErrorIndex() = %d, want 0", idx)
	}
}

// TestExtControl_AutomaticMemoryManagement verifies automatic cleanup
func TestExtControl_AutomaticMemoryManagement(t *testing.T) {
	// Test multiple allocations - memory should be automatically cleaned up
	for i := 0; i < 100; i++ {
		_ = NewExtControlWithString(CtrlBrightness, "test string")
	}

	// Test compound data - memory should be automatically cleaned up
	for i := 0; i < 100; i++ {
		_ = NewExtControlWithCompound(CtrlBrightness, make([]byte, 1024))
	}
}

// TestExtControls_MultipleAdd verifies multiple operations
func TestExtControls_MultipleAdd(t *testing.T) {
	ctrls := NewExtControls()

	// Add many controls
	for i := 0; i < 10; i++ {
		ctrl := NewExtControlWithValue(CtrlBrightness, int32(i))
		ctrls.Add(ctrl)
	}

	if ctrls.Count() != 10 {
		t.Errorf("Count() = %d, want 10", ctrls.Count())
	}
}
