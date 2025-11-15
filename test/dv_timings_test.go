// +build integration

package test

import (
	"testing"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestIntegration_DVTimings_Get tests getting current DV timings
func TestIntegration_DVTimings_Get(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device %s: %v", devPath, err)
	}
	defer dev.Close()

	// Try to get DV timings
	timings, err := dev.GetDVTimings()
	if err != nil {
		t.Logf("Device does not support DV timings (expected for most webcams): %v", err)
		t.Skip("DV timings not supported")
	}

	// If we got timings, validate them
	t.Logf("Current DV timings:")
	t.Logf("  Type: %d", timings.GetType())

	bt := timings.GetBTTimings()
	t.Logf("  Resolution: %dx%d", bt.GetWidth(), bt.GetHeight())
	t.Logf("  Pixel Clock: %d Hz", bt.GetPixelClock())
	t.Logf("  Interlaced: %v", bt.IsInterlaced())
	t.Logf("  Progressive: %v", bt.IsProgressive())
	t.Logf("  Frame Rate: %.2f Hz", bt.GetFrameRate())
}

// TestIntegration_DVTimings_Capabilities tests querying DV timing capabilities
func TestIntegration_DVTimings_Capabilities(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to get DV timing capabilities
	cap, err := dev.GetDVTimingsCap(0)
	if err != nil {
		t.Logf("Device does not support DV timing capabilities: %v", err)
		t.Skip("DV timing capabilities not supported")
	}

	// Display capabilities
	btCap := cap.GetBTTimingsCap()
	t.Logf("DV Timing Capabilities:")
	t.Logf("  Type: %d", cap.GetType())
	t.Logf("  Resolution range: %dx%d to %dx%d",
		btCap.GetMinWidth(), btCap.GetMinHeight(),
		btCap.GetMaxWidth(), btCap.GetMaxHeight())
	t.Logf("  Pixel clock range: %d - %d Hz",
		btCap.GetMinPixelClock(), btCap.GetMaxPixelClock())
	t.Logf("  Standards: 0x%x", btCap.GetStandards())
	t.Logf("  Capabilities: 0x%x", btCap.GetCapabilities())
	t.Logf("  Interlaced support: %v", btCap.SupportsInterlaced())
	t.Logf("  Progressive support: %v", btCap.SupportsProgressive())
	t.Logf("  Reduced blanking: %v", btCap.SupportsReducedBlanking())
	t.Logf("  Custom timings: %v", btCap.SupportsCustomTimings())
	t.Logf("  CEA-861 standard: %v", btCap.HasStandard(v4l2.DVStdCEA861))
	t.Logf("  DMT standard: %v", btCap.HasStandard(v4l2.DVStdDMT))
}

// TestIntegration_DVTimings_Enumerate tests enumerating supported DV timings
func TestIntegration_DVTimings_Enumerate(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to enumerate DV timings
	timings, err := dev.GetAllDVTimings(0)
	if err != nil {
		t.Logf("Device does not support DV timing enumeration: %v", err)
		t.Skip("DV timing enumeration not supported")
	}

	if len(timings) == 0 {
		t.Log("Device reports no DV timings")
		return
	}

	t.Logf("Found %d supported DV timing(s):", len(timings))
	for i, timing := range timings {
		dv := timing.GetTimings()
		bt := dv.GetBTTimings()

		t.Logf("  [%d] %dx%d @ %.2f Hz",
			i, bt.GetWidth(), bt.GetHeight(), bt.GetFrameRate())
		t.Logf("      Pixel Clock: %d Hz", bt.GetPixelClock())
		t.Logf("      Interlaced: %v", bt.IsInterlaced())
		t.Logf("      Standards: CEA-861=%v, DMT=%v, CVT=%v, GTF=%v",
			bt.HasStandard(v4l2.DVStdCEA861),
			bt.HasStandard(v4l2.DVStdDMT),
			bt.HasStandard(v4l2.DVStdCVT),
			bt.HasStandard(v4l2.DVStdGTF))

		// Only show first 5 to avoid too much output
		if i >= 4 {
			t.Logf("  ... and %d more", len(timings)-5)
			break
		}
	}
}

// TestIntegration_DVTimings_Query tests auto-detecting DV timings
func TestIntegration_DVTimings_Query(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to query/detect DV timings
	timings, err := dev.QueryDVTimings()
	if err != nil {
		t.Logf("Failed to query DV timings (no signal or not supported): %v", err)
		t.Skip("DV timing query not supported or no signal")
	}

	// If we detected timings, show them
	bt := timings.GetBTTimings()
	t.Logf("Detected DV timings:")
	t.Logf("  Resolution: %dx%d", bt.GetWidth(), bt.GetHeight())
	t.Logf("  Frame Rate: %.2f Hz", bt.GetFrameRate())
	t.Logf("  Pixel Clock: %d Hz", bt.GetPixelClock())
	t.Logf("  Interlaced: %v", bt.IsInterlaced())
	t.Logf("  H-Sync Positive: %v", bt.HasHSyncPosPolarity())
	t.Logf("  V-Sync Positive: %v", bt.HasVSyncPosPolarity())
}

// TestIntegration_DVTimings_Set tests setting DV timings
func TestIntegration_DVTimings_Set(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current timings first
	currentTimings, err := dev.GetDVTimings()
	if err != nil {
		t.Skip("Device does not support DV timings")
	}

	// Try to set the same timings (safe operation)
	err = dev.SetDVTimings(currentTimings)
	if err != nil {
		t.Logf("SetDVTimings failed (may be read-only): %v", err)
		return
	}

	// Verify timings were set
	newTimings, err := dev.GetDVTimings()
	if err != nil {
		t.Fatalf("Failed to get timings after setting: %v", err)
	}

	currentBT := currentTimings.GetBTTimings()
	newBT := newTimings.GetBTTimings()

	if currentBT.GetWidth() != newBT.GetWidth() || currentBT.GetHeight() != newBT.GetHeight() {
		t.Logf("Warning: Timings may have changed: %dx%d -> %dx%d",
			currentBT.GetWidth(), currentBT.GetHeight(),
			newBT.GetWidth(), newBT.GetHeight())
	} else {
		t.Log("SetDVTimings succeeded")
	}
}

// TestIntegration_DVTimings_EnumerateSpecific tests enumerating specific timing by index
func TestIntegration_DVTimings_EnumerateSpecific(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to get timing at index 0
	enumTiming, err := dev.EnumerateDVTimings(0, 0)
	if err != nil {
		t.Logf("Device does not support DV timing enumeration: %v", err)
		t.Skip("DV timing enumeration not supported")
	}

	// Verify index matches
	if enumTiming.GetIndex() != 0 {
		t.Errorf("Expected index 0, got %d", enumTiming.GetIndex())
	}

	timings := enumTiming.GetTimings()
	bt := timings.GetBTTimings()

	t.Logf("Timing at index 0:")
	t.Logf("  Resolution: %dx%d @ %.2f Hz",
		bt.GetWidth(), bt.GetHeight(), bt.GetFrameRate())
	t.Logf("  Pixel Clock: %d Hz", bt.GetPixelClock())
}

// TestIntegration_DVTimings_StandardsAndFlags tests DV timing standards and flags
func TestIntegration_DVTimings_StandardsAndFlags(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	timings, err := dev.GetDVTimings()
	if err != nil {
		t.Skip("Device does not support DV timings")
	}

	bt := timings.GetBTTimings()

	t.Log("DV Timing Standards and Flags:")
	t.Logf("  Standards: 0x%x", bt.GetStandards())
	t.Logf("    CEA-861: %v", bt.HasStandard(v4l2.DVStdCEA861))
	t.Logf("    DMT: %v", bt.HasStandard(v4l2.DVStdDMT))
	t.Logf("    CVT: %v", bt.HasStandard(v4l2.DVStdCVT))
	t.Logf("    GTF: %v", bt.HasStandard(v4l2.DVStdGTF))

	t.Logf("  Flags: 0x%x", bt.GetFlags())
	t.Logf("    Reduced Blanking: %v", bt.HasFlag(v4l2.DVFlagReducedBlanking))
	t.Logf("    Reduced FPS: %v", bt.HasFlag(v4l2.DVFlagReducedFPS))
	t.Logf("    CE Video: %v", bt.HasFlag(v4l2.DVFlagIsCEVideo))
	t.Logf("    Has Picture Aspect: %v", bt.HasFlag(v4l2.DVFlagHasPictureAspect))

	if bt.HasFlag(v4l2.DVFlagHasCEA861VIC) {
		t.Logf("  CEA-861 VIC: %d", bt.GetCEA861VIC())
	}
	if bt.HasFlag(v4l2.DVFlagHasHDMIVIC) {
		t.Logf("  HDMI VIC: %d", bt.GetHDMIVIC())
	}
}

// TestIntegration_DVTimings_Polarities tests sync polarity detection
func TestIntegration_DVTimings_Polarities(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	timings, err := dev.GetDVTimings()
	if err != nil {
		t.Skip("Device does not support DV timings")
	}

	bt := timings.GetBTTimings()

	t.Log("Sync Polarities:")
	t.Logf("  Polarities: 0x%x", bt.GetPolarities())
	t.Logf("  H-Sync Positive: %v", bt.HasHSyncPosPolarity())
	t.Logf("  V-Sync Positive: %v", bt.HasVSyncPosPolarity())
}

// TestIntegration_DVTimings_BlankingInfo tests blanking period information
func TestIntegration_DVTimings_BlankingInfo(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	timings, err := dev.GetDVTimings()
	if err != nil {
		t.Skip("Device does not support DV timings")
	}

	bt := timings.GetBTTimings()

	t.Log("Blanking Information:")
	t.Logf("  Horizontal:")
	t.Logf("    Front Porch: %d", bt.GetHFrontPorch())
	t.Logf("    Sync: %d", bt.GetHSync())
	t.Logf("    Back Porch: %d", bt.GetHBackPorch())
	t.Logf("    Total: %d pixels", bt.GetWidth()+bt.GetHFrontPorch()+bt.GetHSync()+bt.GetHBackPorch())

	t.Logf("  Vertical:")
	t.Logf("    Front Porch: %d", bt.GetVFrontPorch())
	t.Logf("    Sync: %d", bt.GetVSync())
	t.Logf("    Back Porch: %d", bt.GetVBackPorch())
	t.Logf("    Total: %d lines", bt.GetHeight()+bt.GetVFrontPorch()+bt.GetVSync()+bt.GetVBackPorch())

	if bt.IsInterlaced() {
		t.Logf("  Interlaced Vertical:")
		t.Logf("    IL Front Porch: %d", bt.GetILVFrontPorch())
		t.Logf("    IL Sync: %d", bt.GetILVSync())
		t.Logf("    IL Back Porch: %d", bt.GetILVBackPorch())
	}
}
