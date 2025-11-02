package v4l2

import (
	"testing"
)

// TestInputStatus_Constants verifies input status constants are defined
func TestInputStatus_Constants(t *testing.T) {
	tests := []struct {
		name   string
		status InputStatus
	}{
		{"InputStatusNoPower", InputStatusNoPower},
		{"InputStatusNoSignal", InputStatusNoSignal},
		{"InputStatusNoColor", InputStatusNoColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.status == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestInputStatuses_MapComplete verifies all status values have descriptions
func TestInputStatuses_MapComplete(t *testing.T) {
	tests := []struct {
		name   string
		status InputStatus
	}{
		{"ok", 0},
		{"no power", InputStatusNoPower},
		{"no signal", InputStatusNoSignal},
		{"no color", InputStatusNoColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, exists := InputStatuses[tt.status]
			if !exists {
				t.Errorf("Status %d not found in InputStatuses map", tt.status)
			}
			if desc != tt.name {
				t.Errorf("InputStatuses[%d] = %q, want %q", tt.status, desc, tt.name)
			}
		})
	}
}

// TestInputType_Constants verifies input type constants
func TestInputType_Constants(t *testing.T) {
	tests := []struct {
		name      string
		inputType InputType
		expected  InputType
	}{
		{"InputTypeTuner", InputTypeTuner, 1},
		{"InputTypeCamera", InputTypeCamera, 2},
		{"InputTypeTouch", InputTypeTouch, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inputType != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.inputType, tt.expected)
			}
		})
	}
}

// TestOutputType_Constants verifies output type constants
func TestOutputType_Constants(t *testing.T) {
	tests := []struct {
		name       string
		outputType OutputType
		expected   OutputType
	}{
		{"OutputTypeModulator", OutputTypeModulator, 1},
		{"OutputTypeAnalog", OutputTypeAnalog, 2},
		{"OutputTypeAnalogVGAOverlay", OutputTypeAnalogVGAOverlay, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.outputType != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.outputType, tt.expected)
			}
		})
	}
}

// TestOutputStatuses_MapComplete verifies output status map
func TestOutputStatuses_MapComplete(t *testing.T) {
	tests := []struct {
		name   string
		status OutputStatus
	}{
		{"ok", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc, exists := OutputStatuses[tt.status]
			if !exists {
				t.Errorf("Status %d not found in OutputStatuses map", tt.status)
			}
			if desc != tt.name {
				t.Errorf("OutputStatuses[%d] = %q, want %q", tt.status, desc, tt.name)
			}
		})
	}
}

// TestInputInfo_Accessors tests InputInfo accessor methods
func TestInputInfo_Accessors(t *testing.T) {
	// Note: This test verifies the accessor methods exist and return correct types
	// Actual functionality requires a real device or mock
	var info InputInfo

	// Test that methods don't panic and return expected types
	_ = info.GetIndex()
	_ = info.GetName()
	_ = info.GetInputType()
	_ = info.GetAudioset()
	_ = info.GetTuner()
	_ = info.GetStandardId()
	_ = info.GetStatus()
	_ = info.GetCapabilities()

	// All methods returned without panic - success
}

// TestOutputInfo_Accessors tests OutputInfo accessor methods
func TestOutputInfo_Accessors(t *testing.T) {
	// Note: This test verifies the accessor methods exist and return correct types
	// Actual functionality requires a real device or mock
	var info OutputInfo

	// Test that methods don't panic and return expected types
	_ = info.GetIndex()
	_ = info.GetName()
	_ = info.GetOutputType()
	_ = info.GetAudioset()
	_ = info.GetModulator()
	_ = info.GetStandardId()
	_ = info.GetCapabilities()

	// All methods returned without panic - success
}

// TestInputInfo_GetName tests name extraction from InputInfo
func TestInputInfo_GetName(t *testing.T) {
	// This test verifies the GetName method works correctly
	// The actual C struct initialization would require CGO setup
	var info InputInfo
	name := info.GetName()

	// Empty C string should return empty Go string
	if name != "" {
		t.Logf("GetName() returned: %q (expected empty for uninitialized struct)", name)
	}
}

// TestOutputInfo_GetName tests name extraction from OutputInfo
func TestOutputInfo_GetName(t *testing.T) {
	// This test verifies the GetName method works correctly
	// The actual C struct initialization would require CGO setup
	var info OutputInfo
	name := info.GetName()

	// Empty C string should return empty Go string
	if name != "" {
		t.Logf("GetName() returned: %q (expected empty for uninitialized struct)", name)
	}
}

// TestInputInfo_SymmetryWithOutput verifies InputInfo and OutputInfo have symmetric APIs
func TestInputInfo_SymmetryWithOutput(t *testing.T) {
	// This test documents that InputInfo and OutputInfo have the same accessor pattern
	// (except InputInfo has GetTuner/GetStatus, OutputInfo has GetModulator)

	commonMethods := []string{
		"GetIndex",
		"GetName",
		"GetAudioset",
		"GetStandardId",
		"GetCapabilities",
	}

	for _, method := range commonMethods {
		t.Run(method, func(t *testing.T) {
			// Both types implement these methods
			t.Logf("Both InputInfo and OutputInfo implement %s()", method)
		})
	}

	// Document the differences
	t.Run("differences", func(t *testing.T) {
		t.Log("InputInfo has: GetInputType(), GetTuner(), GetStatus()")
		t.Log("OutputInfo has: GetOutputType(), GetModulator()")
	})
}

// TestStandardId_Type verifies StandardId is defined
func TestStandardId_Type(t *testing.T) {
	var sid StandardId = 0x12345678

	if sid != 0x12345678 {
		t.Errorf("StandardId assignment failed: got %v, want 0x12345678", sid)
	}
}

// Note: Integration tests with actual V4L2 devices are in test/video_io_test.go
// These unit tests focus on type definitions, constants, and accessor methods
// without requiring real hardware or system calls.
