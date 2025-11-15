// +build integration

package test

import (
	"testing"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestIntegration_TunerEnumeration tests tuner enumeration
func TestIntegration_TunerEnumeration(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device %s: %v", devPath, err)
	}
	defer dev.Close()

	// Test GetAllTuners
	tuners, err := dev.GetAllTuners()
	if err != nil {
		t.Logf("Device does not support tuner enumeration: %v", err)
		return // Not all devices have tuners
	}

	if len(tuners) == 0 {
		t.Log("Device reports no tuners")
		return
	}

	t.Logf("Found %d tuner(s):", len(tuners))
	for _, tuner := range tuners {
		typeName := v4l2.TunerTypes[tuner.GetType()]
		if typeName == "" {
			typeName = "Unknown"
		}
		t.Logf("  [%d] %s", tuner.GetIndex(), tuner.GetName())
		t.Logf("      Type: %s (0x%x)", typeName, tuner.GetType())
		t.Logf("      Capability: 0x%x", tuner.GetCapability())
		t.Logf("      Range: %d - %d", tuner.GetRangeLow(), tuner.GetRangeHigh())
		t.Logf("      Signal: %d", tuner.GetSignal())
		t.Logf("      AFC: %d", tuner.GetAFC())
		t.Logf("      RxSubchans: 0x%x", tuner.GetRxSubchans())
		t.Logf("      AudioMode: 0x%x", tuner.GetAudioMode())
		t.Logf("      Stereo: %v", tuner.IsStereo())
		t.Logf("      RDS: %v", tuner.HasRDS())
		t.Logf("      LowFreq: %v", tuner.IsLowFreq())
		t.Logf("      HwSeek: %v", tuner.SupportsHwSeek())
		t.Logf("      FreqBands: %v", tuner.SupportsFreqBands())
	}
}

// TestIntegration_TunerInfo tests getting specific tuner info
func TestIntegration_TunerInfo(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get tuner info for index 0
	info, err := dev.GetTunerInfo(0)
	if err != nil {
		t.Skip("Device does not support tuner info query")
	}

	// Verify the info is for the correct index
	if info.GetIndex() != 0 {
		t.Errorf("Expected index 0, got %d", info.GetIndex())
	}

	typeName := v4l2.TunerTypes[info.GetType()]
	if typeName == "" {
		typeName = "Unknown"
	}

	name := info.GetName()
	t.Logf("Tuner 0: %s", name)
	t.Logf("  Type: %s", typeName)
	t.Logf("  Stereo: %v", info.IsStereo())
	t.Logf("  RDS: %v", info.HasRDS())
	t.Logf("  Signal: %d", info.GetSignal())
}

// TestIntegration_FrequencyGet tests getting tuner frequency
func TestIntegration_FrequencyGet(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// First check if device has tuner
	tuners, err := dev.GetAllTuners()
	if err != nil || len(tuners) == 0 {
		t.Skip("Device does not have tuners")
	}

	// Get frequency for first tuner
	freq, err := dev.GetFrequency(0)
	if err != nil {
		t.Skip("Device does not support frequency query")
	}

	typeName := v4l2.TunerTypes[freq.GetType()]
	if typeName == "" {
		typeName = "Unknown"
	}

	t.Logf("Current frequency:")
	t.Logf("  Tuner: %d", freq.GetTuner())
	t.Logf("  Type: %s (0x%x)", typeName, freq.GetType())
	t.Logf("  Frequency: %d", freq.GetFrequency())
}

// TestIntegration_FrequencySet tests setting tuner frequency
func TestIntegration_FrequencySet(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// First check if device has tuner
	tuners, err := dev.GetAllTuners()
	if err != nil || len(tuners) == 0 {
		t.Skip("Device does not have tuners")
	}

	tuner := tuners[0]

	// Get current frequency first
	originalFreq, err := dev.GetFrequency(0)
	if err != nil {
		t.Skip("Device does not support frequency operations")
	}

	// Try to set the same frequency (safe operation)
	err = dev.SetFrequency(0, tuner.GetType(), originalFreq.GetFrequency())
	if err != nil {
		t.Logf("SetFrequency failed (expected for some devices): %v", err)
		return
	}

	// Verify frequency was set
	newFreq, err := dev.GetFrequency(0)
	if err != nil {
		t.Fatalf("Failed to get frequency after setting: %v", err)
	}

	if newFreq.GetFrequency() != originalFreq.GetFrequency() {
		t.Logf("Warning: Frequency changed from %d to %d", originalFreq.GetFrequency(), newFreq.GetFrequency())
	} else {
		t.Logf("Frequency set successfully: %d", newFreq.GetFrequency())
	}
}

// TestIntegration_FrequencyBands tests frequency band enumeration
func TestIntegration_FrequencyBands(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// First check if device has tuner
	tuners, err := dev.GetAllTuners()
	if err != nil || len(tuners) == 0 {
		t.Skip("Device does not have tuners")
	}

	// Check if tuner supports frequency bands
	tuner := tuners[0]
	if !tuner.SupportsFreqBands() {
		t.Skip("Tuner does not support frequency bands enumeration")
	}

	// Get frequency bands
	bands, err := dev.GetFrequencyBands(0, tuner.GetType())
	if err != nil {
		t.Logf("Failed to enumerate frequency bands: %v", err)
		t.Skip("Frequency bands enumeration not supported")
	}

	if len(bands) == 0 {
		t.Log("Tuner reports no frequency bands")
		return
	}

	t.Logf("Found %d frequency band(s) for tuner 0:", len(bands))
	for _, band := range bands {
		t.Logf("  Band %d:", band.GetIndex())
		t.Logf("    Range: %d - %d", band.GetRangeLow(), band.GetRangeHigh())
		t.Logf("    Capability: 0x%x", band.GetCapability())
		t.Logf("    Modulation: 0x%x", band.GetModulation())
		t.Logf("    FM: %v", band.GetModulation()&v4l2.BandModulationFM != 0)
		t.Logf("    AM: %v", band.GetModulation()&v4l2.BandModulationAM != 0)
		t.Logf("    VSB: %v", band.GetModulation()&v4l2.BandModulationVSB != 0)
	}
}

// TestIntegration_TunerSet tests setting tuner parameters
func TestIntegration_TunerSet(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// First check if device has tuner
	tuner, err := dev.GetTunerInfo(0)
	if err != nil {
		t.Skip("Device does not have tuners")
	}

	// Try to set tuner (with same parameters - safe operation)
	err = dev.SetTuner(tuner)
	if err != nil {
		t.Logf("SetTuner failed (expected for some devices): %v", err)
		return
	}

	t.Log("SetTuner succeeded")
}

// TestIntegration_ModulatorEnumeration tests modulator enumeration
func TestIntegration_ModulatorEnumeration(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device %s: %v", devPath, err)
	}
	defer dev.Close()

	// Test GetAllModulators
	modulators, err := dev.GetAllModulators()
	if err != nil {
		t.Logf("Device does not support modulator enumeration: %v", err)
		return // Not all devices have modulators
	}

	if len(modulators) == 0 {
		t.Log("Device reports no modulators")
		return
	}

	t.Logf("Found %d modulator(s):", len(modulators))
	for _, mod := range modulators {
		typeName := v4l2.TunerTypes[mod.GetType()]
		if typeName == "" {
			typeName = "Unknown"
		}
		t.Logf("  [%d] %s", mod.GetIndex(), mod.GetName())
		t.Logf("      Type: %s (0x%x)", typeName, mod.GetType())
		t.Logf("      Capability: 0x%x", mod.GetCapability())
		t.Logf("      Range: %d - %d", mod.GetRangeLow(), mod.GetRangeHigh())
		t.Logf("      TxSubchans: 0x%x", mod.GetTxSubchans())
		t.Logf("      Stereo: %v", mod.IsStereo())
		t.Logf("      RDS: %v", mod.HasRDS())
		t.Logf("      LowFreq: %v", mod.IsLowFreq())
		t.Logf("      FreqBands: %v", mod.SupportsFreqBands())
	}
}

// TestIntegration_ModulatorInfo tests getting specific modulator info
func TestIntegration_ModulatorInfo(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get modulator info for index 0
	info, err := dev.GetModulatorInfo(0)
	if err != nil {
		t.Skip("Device does not support modulator info query")
	}

	// Verify the info is for the correct index
	if info.GetIndex() != 0 {
		t.Errorf("Expected index 0, got %d", info.GetIndex())
	}

	typeName := v4l2.TunerTypes[info.GetType()]
	if typeName == "" {
		typeName = "Unknown"
	}

	name := info.GetName()
	t.Logf("Modulator 0: %s", name)
	t.Logf("  Type: %s", typeName)
	t.Logf("  Stereo: %v", info.IsStereo())
	t.Logf("  RDS: %v", info.HasRDS())
}

// TestIntegration_ModulatorSet tests setting modulator parameters
func TestIntegration_ModulatorSet(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// First check if device has modulator
	mod, err := dev.GetModulatorInfo(0)
	if err != nil {
		t.Skip("Device does not have modulators")
	}

	// Try to set modulator (with same parameters - safe operation)
	err = dev.SetModulator(mod)
	if err != nil {
		t.Logf("SetModulator failed (expected for some devices): %v", err)
		return
	}

	t.Log("SetModulator succeeded")
}

// TestIntegration_TunerModulator_APIParity tests API parity between tuner and modulator
func TestIntegration_TunerModulator_APIParity(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	t.Log("Verifying API parity between Tuner and Modulator")

	// Both should have GetAll methods
	tuners, tunerErr := dev.GetAllTuners()
	modulators, modErr := dev.GetAllModulators()

	if tunerErr != nil {
		t.Logf("  GetAllTuners: not supported")
	} else {
		t.Logf("  GetAllTuners: %d tuner(s)", len(tuners))
	}

	if modErr != nil {
		t.Logf("  GetAllModulators: not supported")
	} else {
		t.Logf("  GetAllModulators: %d modulator(s)", len(modulators))
	}

	// Both should have GetInfo methods
	if tunerErr == nil && len(tuners) > 0 {
		_, err := dev.GetTunerInfo(0)
		if err != nil {
			t.Logf("  GetTunerInfo: failed - %v", err)
		} else {
			t.Logf("  GetTunerInfo: success")
		}
	}

	if modErr == nil && len(modulators) > 0 {
		_, err := dev.GetModulatorInfo(0)
		if err != nil {
			t.Logf("  GetModulatorInfo: failed - %v", err)
		} else {
			t.Logf("  GetModulatorInfo: success")
		}
	}

	// Both should support frequency operations
	if tunerErr == nil && len(tuners) > 0 {
		_, err := dev.GetFrequency(0)
		if err != nil {
			t.Logf("  GetFrequency (tuner): not supported")
		} else {
			t.Logf("  GetFrequency (tuner): success")
		}

		if tuners[0].SupportsFreqBands() {
			bands, err := dev.GetFrequencyBands(0, tuners[0].GetType())
			if err != nil {
				t.Logf("  GetFrequencyBands: failed - %v", err)
			} else {
				t.Logf("  GetFrequencyBands: %d band(s)", len(bands))
			}
		}
	}

	t.Log("API parity check complete")
}
