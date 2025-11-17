// +build integration

package test

import (
	"testing"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestIntegration_Standards_Enumerate tests enumerating supported video standards
func TestIntegration_Standards_Enumerate(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device %s: %v", devPath, err)
	}
	defer dev.Close()

	// Try to enumerate standards
	standards, err := dev.GetAllStandards()
	if err != nil {
		t.Logf("Failed to enumerate standards: %v", err)
		t.Skip("Standard enumeration not supported (expected for digital-only devices)")
	}

	if len(standards) == 0 {
		t.Log("Device reports no video standards (digital-only device)")
		t.Skip("No video standards available")
	}

	t.Logf("Found %d supported video standard(s):", len(standards))
	for _, std := range standards {
		t.Logf("  Standard %d:", std.Index())
		t.Logf("    ID: 0x%016x", std.ID())
		t.Logf("    Name: %s", std.Name())
		t.Logf("    Frame rate: %.2f fps", std.FrameRate())
		t.Logf("    Frame lines: %d", std.FrameLines())
		framePeriod := std.FramePeriod()
		t.Logf("    Frame period: %d/%d seconds", framePeriod.Numerator, framePeriod.Denominator)
		t.Logf("    String: %s", std.String())
	}
}

// TestIntegration_Standards_Get tests getting current video standard
func TestIntegration_Standards_Get(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to get current standard
	stdId, err := dev.GetStandard()
	if err != nil {
		t.Logf("Device does not support video standards (expected for digital-only devices): %v", err)
		t.Skip("Video standards not supported")
	}

	t.Logf("Current video standard: 0x%016x", stdId)
	if name, ok := v4l2.StdNames[stdId]; ok {
		t.Logf("Standard name: %s", name)
	} else {
		t.Logf("Standard name: (unknown)")
	}

	// Check if it's a known standard
	if stdId == v4l2.StdPAL {
		t.Log("Current standard is PAL")
	} else if stdId == v4l2.StdNTSC {
		t.Log("Current standard is NTSC")
	} else if stdId == v4l2.StdSECAM {
		t.Log("Current standard is SECAM")
	} else if stdId&v4l2.StdPAL != 0 {
		t.Log("Current standard contains PAL variant")
	} else if stdId&v4l2.StdNTSC != 0 {
		t.Log("Current standard contains NTSC variant")
	}
}

// TestIntegration_Standards_SetRestore tests setting video standard and restoring
func TestIntegration_Standards_SetRestore(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current standard to restore later
	originalStd, err := dev.GetStandard()
	if err != nil {
		t.Logf("Device does not support video standards: %v", err)
		t.Skip("Video standards not supported")
	}

	t.Logf("Original standard: 0x%016x", originalStd)

	// Get all supported standards
	standards, err := dev.GetAllStandards()
	if err != nil || len(standards) == 0 {
		t.Skip("Cannot enumerate standards to test setting")
	}

	// Try to set the first enumerated standard
	testStd := standards[0]
	t.Logf("Attempting to set standard %d: %s (0x%016x)",
		testStd.Index(), testStd.Name(), testStd.ID())

	err = dev.SetStandard(testStd.ID())
	if err != nil {
		t.Logf("Failed to set standard (may be read-only): %v", err)
	} else {
		t.Log("Successfully set standard")

		// Verify it was set
		currentStd, err := dev.GetStandard()
		if err != nil {
			t.Errorf("Failed to get standard after setting: %v", err)
		} else {
			t.Logf("Current standard after set: 0x%016x", currentStd)
			if (currentStd & testStd.ID()) != 0 {
				t.Log("Standard was successfully set (matches using bitwise AND)")
			}
		}
	}

	// Restore original standard
	err = dev.SetStandard(originalStd)
	if err != nil {
		t.Logf("Warning: Failed to restore original standard: %v", err)
	} else {
		t.Log("Successfully restored original standard")
	}
}

// TestIntegration_Standards_Query tests auto-detecting video standard
func TestIntegration_Standards_Query(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to query/detect standard
	detectedStd, err := dev.QueryStandard()
	if err != nil {
		t.Logf("Standard detection not supported or no signal: %v", err)
		t.Skip("Cannot detect standard (expected if no analog signal present)")
	}

	t.Logf("Detected standard: 0x%016x", detectedStd)
	if name, ok := v4l2.StdNames[detectedStd]; ok {
		t.Logf("Detected standard name: %s", name)
	}

	// Verify detected standard is in the supported list
	standards, err := dev.GetAllStandards()
	if err == nil {
		found := false
		for _, std := range standards {
			if (std.ID() & detectedStd) != 0 {
				found = true
				t.Logf("Detected standard matches: %s", std.Name())
				break
			}
		}
		if !found {
			t.Log("Warning: Detected standard not in enumerated list")
		}
	}
}

// TestIntegration_Standards_IsSupported tests checking if specific standards are supported
func TestIntegration_Standards_IsSupported(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Test common standards
	tests := []struct {
		name  string
		stdId v4l2.StdId
	}{
		{"PAL", v4l2.StdPAL},
		{"NTSC", v4l2.StdNTSC},
		{"SECAM", v4l2.StdSECAM},
		{"PAL-B/G", v4l2.StdPAL_BG},
		{"NTSC-M", v4l2.StdNTSC_M},
		{"525/60", v4l2.Std525_60},
		{"625/50", v4l2.Std625_50},
	}

	hasAnyStandard := false
	for _, tt := range tests {
		supported, err := dev.IsStandardSupported(tt.stdId)
		if err != nil {
			t.Logf("Error checking if %s is supported: %v", tt.name, err)
			continue
		}
		t.Logf("%s supported: %v", tt.name, supported)
		if supported {
			hasAnyStandard = true
		}
	}

	if !hasAnyStandard {
		t.Log("Device does not support any common analog video standards (digital-only device)")
	}
}

// TestIntegration_Standards_IndividualEnumerate tests enumerating standards one by one
func TestIntegration_Standards_IndividualEnumerate(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to enumerate standards individually
	for index := uint32(0); index < 10; index++ {
		std, err := dev.EnumerateStandard(index)
		if err != nil {
			if index == 0 {
				t.Logf("Standard enumeration not supported: %v", err)
				t.Skip("Video standards not supported")
			}
			// Reached end of list
			t.Logf("Enumerated %d standard(s)", index)
			break
		}

		t.Logf("Standard %d: %s (ID: 0x%016x, %.2f fps, %d lines)",
			std.Index(), std.Name(), std.ID(), std.FrameRate(), std.FrameLines())
	}
}

// TestIntegration_Standards_GetByID tests getting standard by ID
func TestIntegration_Standards_GetByID(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get all standards
	standards, err := dev.GetAllStandards()
	if err != nil || len(standards) == 0 {
		t.Skip("No standards available")
	}

	// Try to get each standard by its ID using v4l2 package function
	for _, originalStd := range standards {
		std, err := v4l2.GetStandardByID(dev.Fd(), originalStd.ID())
		if err != nil {
			t.Errorf("Failed to get standard by ID 0x%016x: %v", originalStd.ID(), err)
			continue
		}

		t.Logf("Found standard by ID 0x%016x: %s", originalStd.ID(), std.Name())

		// Verify IDs match
		if (std.ID() & originalStd.ID()) == 0 {
			t.Errorf("Standard ID mismatch: got 0x%016x, want 0x%016x",
				std.ID(), originalStd.ID())
		}
	}
}

// TestIntegration_Standards_CommonGroupings tests common standard groupings
func TestIntegration_Standards_CommonGroupings(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Test that standard groupings work correctly
	standards, err := dev.GetAllStandards()
	if err != nil || len(standards) == 0 {
		t.Skip("No standards available")
	}

	// Test bitwise operations with groupings
	for _, std := range standards {
		stdId := std.ID()

		// Test if this is a PAL variant
		if (stdId & v4l2.StdPAL) != 0 {
			t.Logf("%s is a PAL variant", std.Name())
		}

		// Test if this is an NTSC variant
		if (stdId & v4l2.StdNTSC) != 0 {
			t.Logf("%s is an NTSC variant", std.Name())
		}

		// Test if this is a SECAM variant
		if (stdId & v4l2.StdSECAM) != 0 {
			t.Logf("%s is a SECAM variant", std.Name())
		}

		// Test 525/60 vs 625/50
		if (stdId & v4l2.Std525_60) != 0 {
			t.Logf("%s is 525 lines / 60Hz", std.Name())
		}
		if (stdId & v4l2.Std625_50) != 0 {
			t.Logf("%s is 625 lines / 50Hz", std.Name())
		}
	}
}
