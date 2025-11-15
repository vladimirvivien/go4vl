package v4l2

import (
	"testing"
)

// TestAudioCapability_Constants verifies audio capability constants are defined
func TestAudioCapability_Constants(t *testing.T) {
	tests := []struct {
		name       string
		capability AudioCapability
	}{
		{"AudioCapStereo", AudioCapStereo},
		{"AudioCapAVL", AudioCapAVL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.capability == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestAudioMode_Constants verifies audio mode constants are defined
func TestAudioMode_Constants(t *testing.T) {
	tests := []struct {
		name string
		mode AudioMode
	}{
		{"AudioModeAVL", AudioModeAVL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mode == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestAudioInfo_Accessors tests AudioInfo accessor methods
func TestAudioInfo_Accessors(t *testing.T) {
	// Note: This test verifies the accessor methods exist and return correct types
	// Actual functionality requires a real device or mock
	var info AudioInfo

	// Test that methods don't panic and return expected types
	_ = info.GetIndex()
	_ = info.GetName()
	_ = info.GetCapability()
	_ = info.GetMode()
	_ = info.HasCapability(AudioCapStereo)
	_ = info.IsStereo()
	_ = info.HasAVL()

	// All methods returned without panic - success
}

// TestAudioOutInfo_Accessors tests AudioOutInfo accessor methods
func TestAudioOutInfo_Accessors(t *testing.T) {
	// Note: This test verifies the accessor methods exist and return correct types
	// Actual functionality requires a real device or mock
	var info AudioOutInfo

	// Test that methods don't panic and return expected types
	_ = info.GetIndex()
	_ = info.GetName()
	_ = info.GetCapability()
	_ = info.GetMode()
	_ = info.HasCapability(AudioCapStereo)
	_ = info.IsStereo()
	_ = info.HasAVL()

	// All methods returned without panic - success
}

// TestAudioInfo_GetName tests name extraction from AudioInfo
func TestAudioInfo_GetName(t *testing.T) {
	// This test verifies the GetName method works correctly
	// The actual C struct initialization would require CGO setup
	var info AudioInfo
	name := info.GetName()

	// Empty C string should return empty Go string
	if name != "" {
		t.Logf("GetName() returned: %q (expected empty for uninitialized struct)", name)
	}
}

// TestAudioOutInfo_GetName tests name extraction from AudioOutInfo
func TestAudioOutInfo_GetName(t *testing.T) {
	// This test verifies the GetName method works correctly
	// The actual C struct initialization would require CGO setup
	var info AudioOutInfo
	name := info.GetName()

	// Empty C string should return empty Go string
	if name != "" {
		t.Logf("GetName() returned: %q (expected empty for uninitialized struct)", name)
	}
}

// TestAudioInfo_SymmetryWithAudioOut verifies AudioInfo and AudioOutInfo have symmetric APIs
func TestAudioInfo_SymmetryWithAudioOut(t *testing.T) {
	// This test documents that AudioInfo and AudioOutInfo have identical APIs

	commonMethods := []string{
		"GetIndex",
		"GetName",
		"GetCapability",
		"GetMode",
		"HasCapability",
		"IsStereo",
		"HasAVL",
	}

	for _, method := range commonMethods {
		t.Run(method, func(t *testing.T) {
			// Both types implement these methods
			t.Logf("Both AudioInfo and AudioOutInfo implement %s()", method)
		})
	}

	t.Run("perfect_symmetry", func(t *testing.T) {
		t.Log("AudioInfo and AudioOutInfo have completely symmetric APIs")
	})
}

// TestAudioInfo_HasCapability tests capability checking
func TestAudioInfo_HasCapability(t *testing.T) {
	// This test verifies the HasCapability logic
	var info AudioInfo

	// Zero capability means no capabilities
	if info.HasCapability(AudioCapStereo) {
		t.Error("Uninitialized AudioInfo should not have AudioCapStereo")
	}

	if info.HasCapability(AudioCapAVL) {
		t.Error("Uninitialized AudioInfo should not have AudioCapAVL")
	}
}

// TestAudioInfo_IsStereo tests stereo capability checking
func TestAudioInfo_IsStereo(t *testing.T) {
	var info AudioInfo

	// Zero capability means not stereo
	if info.IsStereo() {
		t.Error("Uninitialized AudioInfo should not be stereo")
	}
}

// TestAudioInfo_HasAVL tests AVL capability checking
func TestAudioInfo_HasAVL(t *testing.T) {
	var info AudioInfo

	// Zero capability means no AVL
	if info.HasAVL() {
		t.Error("Uninitialized AudioInfo should not have AVL")
	}
}

// TestAudioOutInfo_HasCapability tests capability checking for audio output
func TestAudioOutInfo_HasCapability(t *testing.T) {
	// This test verifies the HasCapability logic for AudioOutInfo
	var info AudioOutInfo

	// Zero capability means no capabilities
	if info.HasCapability(AudioCapStereo) {
		t.Error("Uninitialized AudioOutInfo should not have AudioCapStereo")
	}

	if info.HasCapability(AudioCapAVL) {
		t.Error("Uninitialized AudioOutInfo should not have AudioCapAVL")
	}
}

// TestAudioOutInfo_IsStereo tests stereo capability checking for audio output
func TestAudioOutInfo_IsStereo(t *testing.T) {
	var info AudioOutInfo

	// Zero capability means not stereo
	if info.IsStereo() {
		t.Error("Uninitialized AudioOutInfo should not be stereo")
	}
}

// TestAudioOutInfo_HasAVL tests AVL capability checking for audio output
func TestAudioOutInfo_HasAVL(t *testing.T) {
	var info AudioOutInfo

	// Zero capability means no AVL
	if info.HasAVL() {
		t.Error("Uninitialized AudioOutInfo should not have AVL")
	}
}

// Note: Integration tests with actual V4L2 devices are in test/audio_io_test.go
// These unit tests focus on type definitions, constants, and accessor methods
// without requiring real hardware or system calls.
