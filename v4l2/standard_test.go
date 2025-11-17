package v4l2

import (
	"testing"
)

// TestStdId_Constants verifies video standard ID constants
func TestStdId_Constants(t *testing.T) {
	tests := []struct {
		name  string
		stdId StdId
	}{
		// PAL variants
		{"StdPAL_B", StdPAL_B},
		{"StdPAL_B1", StdPAL_B1},
		{"StdPAL_G", StdPAL_G},
		{"StdPAL_H", StdPAL_H},
		{"StdPAL_I", StdPAL_I},
		{"StdPAL_D", StdPAL_D},
		{"StdPAL_D1", StdPAL_D1},
		{"StdPAL_K", StdPAL_K},
		{"StdPAL_M", StdPAL_M},
		{"StdPAL_N", StdPAL_N},
		{"StdPAL_Nc", StdPAL_Nc},
		{"StdPAL_60", StdPAL_60},
		// NTSC variants
		{"StdNTSC_M", StdNTSC_M},
		{"StdNTSC_M_JP", StdNTSC_M_JP},
		{"StdNTSC_443", StdNTSC_443},
		{"StdNTSC_M_KR", StdNTSC_M_KR},
		// SECAM variants
		{"StdSECAM_B", StdSECAM_B},
		{"StdSECAM_D", StdSECAM_D},
		{"StdSECAM_G", StdSECAM_G},
		{"StdSECAM_H", StdSECAM_H},
		{"StdSECAM_K", StdSECAM_K},
		{"StdSECAM_K1", StdSECAM_K1},
		{"StdSECAM_L", StdSECAM_L},
		{"StdSECAM_LC", StdSECAM_LC},
		// ATSC
		{"StdATSC_8_VSB", StdATSC_8_VSB},
		{"StdATSC_16_VSB", StdATSC_16_VSB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stdId == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestStdId_Groupings verifies standard grouping constants
func TestStdId_Groupings(t *testing.T) {
	tests := []struct {
		name  string
		stdId StdId
	}{
		{"StdPAL_BG", StdPAL_BG},
		{"StdPAL_DK", StdPAL_DK},
		{"StdPAL", StdPAL},
		{"StdNTSC", StdNTSC},
		{"StdSECAM_DK", StdSECAM_DK},
		{"StdSECAM", StdSECAM},
		{"StdB", StdB},
		{"StdG", StdG},
		{"StdH", StdH},
		{"StdL", StdL},
		{"StdGH", StdGH},
		{"StdDK", StdDK},
		{"StdBG", StdBG},
		{"StdMN", StdMN},
		{"Std525_60", Std525_60},
		{"Std625_50", Std625_50},
		{"StdATSC", StdATSC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stdId == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestStdId_SpecialValues verifies special standard values
func TestStdId_SpecialValues(t *testing.T) {
	// StdUnknown is 0 by definition
	if StdUnknown != 0 {
		t.Error("StdUnknown should be zero")
	}
	// StdAll represents all standards (all bits set)
	if StdAll == 0 {
		t.Error("StdAll should not be zero")
	}
}

// TestStdId_BitOperations verifies standard IDs can be OR'd together
func TestStdId_BitOperations(t *testing.T) {
	// Test that groupings contain individual standards
	if (StdPAL_BG & StdPAL_B) == 0 {
		t.Error("StdPAL_BG should contain StdPAL_B")
	}
	if (StdPAL_BG & StdPAL_G) == 0 {
		t.Error("StdPAL_BG should contain StdPAL_G")
	}
	if (StdNTSC & StdNTSC_M) == 0 {
		t.Error("StdNTSC should contain StdNTSC_M")
	}
	if (StdSECAM & StdSECAM_B) == 0 {
		t.Error("StdSECAM should contain StdSECAM_B")
	}
}

// TestStdNames_Coverage verifies StdNames map has entries for common standards
func TestStdNames_Coverage(t *testing.T) {
	tests := []struct {
		name  string
		stdId StdId
	}{
		{"PAL-B", StdPAL_B},
		{"NTSC-M", StdNTSC_M},
		{"SECAM-L", StdSECAM_L},
		{"PAL", StdPAL},
		{"NTSC", StdNTSC},
		{"SECAM", StdSECAM},
		{"525/60", Std525_60},
		{"625/50", Std625_50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, ok := StdNames[tt.stdId]
			if !ok {
				t.Errorf("StdNames missing entry for %s", tt.name)
			}
			if name == "" {
				t.Errorf("StdNames has empty name for %s", tt.name)
			}
		})
	}
}

// TestFract_ToFloat verifies Fract.ToFloat()
func TestFract_ToFloat(t *testing.T) {
	tests := []struct {
		name     string
		fract    Fract
		expected float64
	}{
		{"1/30", Fract{1, 30}, 1.0 / 30.0},
		{"1/25", Fract{1, 25}, 1.0 / 25.0},
		{"1001/30000", Fract{1001, 30000}, 1001.0 / 30000.0},
		{"Zero denominator", Fract{1, 0}, 0},
		{"Zero numerator", Fract{0, 30}, 0},
		{"Both zero", Fract{0, 0}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fract.ToFloat()
			if result != tt.expected {
				t.Errorf("ToFloat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestFract_FrameRate verifies Fract.FrameRate()
func TestFract_FrameRate(t *testing.T) {
	tests := []struct {
		name     string
		fract    Fract
		expected float64
		delta    float64
	}{
		{"NTSC (29.97 fps)", Fract{1001, 30000}, 29.97, 0.01},
		{"PAL (25 fps)", Fract{1, 25}, 25.0, 0.001},
		{"30 fps", Fract{1, 30}, 30.0, 0.001},
		{"Zero period", Fract{0, 30}, 0, 0},
		{"Zero denominator", Fract{1, 0}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fract.FrameRate()
			if abs(result-tt.expected) > tt.delta {
				t.Errorf("FrameRate() = %v, want %v ± %v", result, tt.expected, tt.delta)
			}
		})
	}
}

// Helper function for float comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// TestStandard_NewStandard verifies Standard creation
func TestStandard_NewStandard(t *testing.T) {
	std := NewStandard(5)
	if std.Index() != 5 {
		t.Errorf("NewStandard(5).Index() = %d, want 5", std.Index())
	}
}

// TestStandard_SetIndex verifies Standard.SetIndex()
func TestStandard_SetIndex(t *testing.T) {
	std := NewStandard(0)
	std.SetIndex(10)
	if std.Index() != 10 {
		t.Errorf("After SetIndex(10), Index() = %d, want 10", std.Index())
	}
}

// TestStandard_SetID verifies Standard.SetID()
func TestStandard_SetID(t *testing.T) {
	std := NewStandard(0)
	std.SetID(StdPAL_B)
	if std.ID() != StdPAL_B {
		t.Errorf("After SetID(StdPAL_B), ID() = 0x%x, want 0x%x", std.ID(), StdPAL_B)
	}
}

// TestStandard_FrameRate verifies Standard.FrameRate()
func TestStandard_FrameRate(t *testing.T) {
	std := NewStandard(0)
	// Manually set the frame period (this is normally set by the kernel)
	// For testing, we'll just verify the calculation works
	// Note: We can't directly set frameperiod in the C struct from Go safely,
	// so we'll test the Fract methods separately (already done above)

	// This is more of a smoke test to ensure the method exists and doesn't crash
	_ = std.FrameRate()
}

// TestStandard_String verifies Standard.String()
func TestStandard_String(t *testing.T) {
	std := NewStandard(0)
	std.SetID(StdPAL_BG)

	// Should not crash
	result := std.String()
	if result == "" {
		t.Error("String() should not return empty string")
	}
}

// TestStandard_Name verifies Standard.Name()
func TestStandard_Name(t *testing.T) {
	std := NewStandard(0)

	// Should not crash (may return empty string if not initialized)
	_ = std.Name()
}

// TestStandard_FramePeriod verifies Standard.FramePeriod()
func TestStandard_FramePeriod(t *testing.T) {
	std := NewStandard(0)

	// Should not crash
	fract := std.FramePeriod()
	_ = fract
}

// TestStandard_FrameLines verifies Standard.FrameLines()
func TestStandard_FrameLines(t *testing.T) {
	std := NewStandard(0)

	// Should not crash (will return 0 if not initialized)
	lines := std.FrameLines()
	_ = lines
}

// BenchmarkFract_ToFloat benchmarks Fract.ToFloat()
func BenchmarkFract_ToFloat(b *testing.B) {
	fract := Fract{1001, 30000}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fract.ToFloat()
	}
}

// BenchmarkFract_FrameRate benchmarks Fract.FrameRate()
func BenchmarkFract_FrameRate(b *testing.B) {
	fract := Fract{1001, 30000}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fract.FrameRate()
	}
}

// BenchmarkStandard_String benchmarks Standard.String()
func BenchmarkStandard_String(b *testing.B) {
	std := NewStandard(0)
	std.SetID(StdPAL_BG)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = std.String()
	}
}
