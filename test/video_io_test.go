// +build integration

package test

import (
	"testing"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestIntegration_VideoInputEnumeration tests video input enumeration
func TestIntegration_VideoInputEnumeration(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device %s: %v", devPath, err)
	}
	defer dev.Close()

	// Skip if device doesn't support video capture
	if !dev.Capability().IsVideoCaptureSupported() {
		t.Skip("Device does not support video capture")
	}

	// Test GetVideoInputIndex
	currentIdx, err := dev.GetVideoInputIndex()
	if err != nil {
		t.Logf("Device does not support input selection: %v", err)
		return // Not all devices have multiple inputs
	}

	t.Logf("Current input index: %d", currentIdx)

	// Test GetVideoInputDescriptions
	inputs, err := dev.GetVideoInputDescriptions()
	if err != nil {
		t.Fatalf("Failed to enumerate inputs: %v", err)
	}

	if len(inputs) == 0 {
		t.Fatal("Expected at least one input")
	}

	t.Logf("Found %d video input(s):", len(inputs))
	for _, input := range inputs {
		t.Logf("  [%d] %s", input.GetIndex(), input.GetName())
		t.Logf("      Type: %d", input.GetInputType())
		t.Logf("      Status: %s (0x%x)", v4l2.InputStatuses[input.GetStatus()], input.GetStatus())
		t.Logf("      Audioset: 0x%x", input.GetAudioset())
		t.Logf("      Standards: 0x%x", input.GetStandardId())
		t.Logf("      Capabilities: 0x%x", input.GetCapabilities())

		if input.GetIndex() == uint32(currentIdx) {
			t.Logf("      ** ACTIVE **")
		}
	}
}

// TestIntegration_VideoInputInfo tests getting specific input info
func TestIntegration_VideoInputInfo(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	if !dev.Capability().IsVideoCaptureSupported() {
		t.Skip("Device does not support video capture")
	}

	// Get current input
	currentIdx, err := dev.GetVideoInputIndex()
	if err != nil {
		t.Skip("Device does not support input selection")
	}

	// Get info for current input
	info, err := dev.GetVideoInputInfo(uint32(currentIdx))
	if err != nil {
		t.Fatalf("Failed to get input info for index %d: %v", currentIdx, err)
	}

	// Verify the info is for the correct index
	if info.GetIndex() != uint32(currentIdx) {
		t.Errorf("Expected index %d, got %d", currentIdx, info.GetIndex())
	}

	name := info.GetName()
	if name == "" {
		t.Error("Input name should not be empty")
	}

	t.Logf("Input %d: %s", info.GetIndex(), name)
}

// TestIntegration_VideoInputStatus tests input status query
func TestIntegration_VideoInputStatus(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	if !dev.Capability().IsVideoCaptureSupported() {
		t.Skip("Device does not support video capture")
	}

	status, err := dev.GetVideoInputStatus()
	if err != nil {
		t.Skip("Device does not support input status query")
	}

	statusStr, exists := v4l2.InputStatuses[status]
	if !exists {
		t.Errorf("Unknown status value: 0x%x", status)
	}

	t.Logf("Input status: %s (0x%x)", statusStr, status)

	// Status should be one of the known values
	if status != 0 &&
	   status != v4l2.InputStatusNoPower &&
	   status != v4l2.InputStatusNoSignal &&
	   status != v4l2.InputStatusNoColor {
		t.Logf("Warning: Unexpected status value: 0x%x", status)
	}
}

// TestIntegration_VideoInputSelection tests switching inputs
func TestIntegration_VideoInputSelection(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	if !dev.Capability().IsVideoCaptureSupported() {
		t.Skip("Device does not support video capture")
	}

	// Get current input
	originalIdx, err := dev.GetVideoInputIndex()
	if err != nil {
		t.Skip("Device does not support input selection")
	}

	// Get all inputs
	inputs, err := dev.GetVideoInputDescriptions()
	if err != nil {
		t.Fatalf("Failed to enumerate inputs: %v", err)
	}

	if len(inputs) < 2 {
		t.Skip("Device has only one input, cannot test selection")
	}

	// Try to switch to a different input
	var newIdx int32 = -1
	for _, input := range inputs {
		if int32(input.GetIndex()) != originalIdx {
			newIdx = int32(input.GetIndex())
			break
		}
	}

	if newIdx == -1 {
		t.Skip("Could not find alternative input")
	}

	t.Logf("Switching from input %d to input %d", originalIdx, newIdx)

	// Switch input
	err = dev.SetVideoInputIndex(newIdx)
	if err != nil {
		t.Fatalf("Failed to set input index: %v", err)
	}

	// Verify the change
	currentIdx, err := dev.GetVideoInputIndex()
	if err != nil {
		t.Fatalf("Failed to get input index after change: %v", err)
	}

	if currentIdx != newIdx {
		t.Errorf("Expected input %d, got %d", newIdx, currentIdx)
	}

	// Restore original input
	err = dev.SetVideoInputIndex(originalIdx)
	if err != nil {
		t.Fatalf("Failed to restore original input: %v", err)
	}

	t.Logf("Successfully switched input and restored original")
}

// TestIntegration_VideoOutputEnumeration tests video output enumeration
func TestIntegration_VideoOutputEnumeration(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Skip if device doesn't support video output
	if !dev.Capability().IsVideoOutputSupported() {
		t.Skip("Device does not support video output")
	}

	// Test GetVideoOutputIndex
	currentIdx, err := dev.GetVideoOutputIndex()
	if err != nil {
		t.Logf("Device does not support output selection: %v", err)
		return
	}

	t.Logf("Current output index: %d", currentIdx)

	// Test GetVideoOutputDescriptions
	outputs, err := dev.GetVideoOutputDescriptions()
	if err != nil {
		t.Fatalf("Failed to enumerate outputs: %v", err)
	}

	if len(outputs) == 0 {
		t.Fatal("Expected at least one output")
	}

	t.Logf("Found %d video output(s):", len(outputs))
	for _, output := range outputs {
		t.Logf("  [%d] %s", output.GetIndex(), output.GetName())
		t.Logf("      Type: %d", output.GetOutputType())
		t.Logf("      Audioset: 0x%x", output.GetAudioset())
		t.Logf("      Modulator: %d", output.GetModulator())
		t.Logf("      Standards: 0x%x", output.GetStandardId())
		t.Logf("      Capabilities: 0x%x", output.GetCapabilities())

		if output.GetIndex() == uint32(currentIdx) {
			t.Logf("      ** ACTIVE **")
		}
	}
}

// TestIntegration_VideoOutputInfo tests getting specific output info
func TestIntegration_VideoOutputInfo(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	if !dev.Capability().IsVideoOutputSupported() {
		t.Skip("Device does not support video output")
	}

	// Get current output
	currentIdx, err := dev.GetVideoOutputIndex()
	if err != nil {
		t.Skip("Device does not support output selection")
	}

	// Get info for current output
	info, err := dev.GetVideoOutputInfo(uint32(currentIdx))
	if err != nil {
		t.Fatalf("Failed to get output info for index %d: %v", currentIdx, err)
	}

	// Verify the info is for the correct index
	if info.GetIndex() != uint32(currentIdx) {
		t.Errorf("Expected index %d, got %d", currentIdx, info.GetIndex())
	}

	name := info.GetName()
	if name == "" {
		t.Error("Output name should not be empty")
	}

	t.Logf("Output %d: %s", info.GetIndex(), name)
}

// TestIntegration_VideoOutputStatus tests output status query
func TestIntegration_VideoOutputStatus(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	if !dev.Capability().IsVideoOutputSupported() {
		t.Skip("Device does not support video output")
	}

	status, err := dev.GetVideoOutputStatus()
	if err != nil {
		t.Skip("Device does not support output status query")
	}

	statusStr, exists := v4l2.OutputStatuses[status]
	if !exists {
		t.Errorf("Unknown status value: 0x%x", status)
	}

	t.Logf("Output status: %s (0x%x)", statusStr, status)

	// Output status should always be OK (0) as per V4L2 API
	if status != 0 {
		t.Logf("Warning: Output status is not OK: 0x%x", status)
	}
}

// TestIntegration_UnsupportedFeatures tests error handling for unsupported features
func TestIntegration_UnsupportedFeatures(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	cap := dev.Capability()

	// Test output methods on capture-only device
	if cap.IsVideoCaptureSupported() && !cap.IsVideoOutputSupported() {
		t.Run("output on capture device", func(t *testing.T) {
			_, err := dev.GetVideoOutputIndex()
			if err != v4l2.ErrorUnsupportedFeature {
				t.Errorf("Expected ErrorUnsupportedFeature, got %v", err)
			}

			err = dev.SetVideoOutputIndex(0)
			if err != v4l2.ErrorUnsupportedFeature {
				t.Errorf("Expected ErrorUnsupportedFeature, got %v", err)
			}

			_, err = dev.GetVideoOutputDescriptions()
			if err != v4l2.ErrorUnsupportedFeature {
				t.Errorf("Expected ErrorUnsupportedFeature, got %v", err)
			}
		})
	}

	// Test input methods on output-only device
	if cap.IsVideoOutputSupported() && !cap.IsVideoCaptureSupported() {
		t.Run("input on output device", func(t *testing.T) {
			_, err := dev.GetVideoInputIndex()
			if err != v4l2.ErrorUnsupportedFeature {
				t.Errorf("Expected ErrorUnsupportedFeature, got %v", err)
			}

			err = dev.SetVideoInputIndex(0)
			if err != v4l2.ErrorUnsupportedFeature {
				t.Errorf("Expected ErrorUnsupportedFeature, got %v", err)
			}

			_, err = dev.GetVideoInputDescriptions()
			if err != v4l2.ErrorUnsupportedFeature {
				t.Errorf("Expected ErrorUnsupportedFeature, got %v", err)
			}
		})
	}
}
