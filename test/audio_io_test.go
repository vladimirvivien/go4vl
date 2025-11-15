// +build integration

package test

import (
	"testing"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestIntegration_AudioInputEnumeration tests audio input enumeration
func TestIntegration_AudioInputEnumeration(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device %s: %v", devPath, err)
	}
	defer dev.Close()

	// Test GetAudioDescriptions
	audios, err := dev.GetAudioDescriptions()
	if err != nil {
		t.Logf("Device does not support audio input enumeration: %v", err)
		return // Not all devices have audio inputs
	}

	if len(audios) == 0 {
		t.Log("Device reports no audio inputs")
		return
	}

	t.Logf("Found %d audio input(s):", len(audios))
	for _, audio := range audios {
		t.Logf("  [%d] %s", audio.GetIndex(), audio.GetName())
		t.Logf("      Capability: 0x%x", audio.GetCapability())
		t.Logf("      Mode: 0x%x", audio.GetMode())
		t.Logf("      Stereo: %v", audio.IsStereo())
		t.Logf("      AVL: %v", audio.HasAVL())
	}
}

// TestIntegration_AudioInputInfo tests getting specific audio input info
func TestIntegration_AudioInputInfo(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get audio info for index 0
	info, err := dev.GetAudioInfo(0)
	if err != nil {
		t.Skip("Device does not support audio input info query")
	}

	// Verify the info is for the correct index
	if info.GetIndex() != 0 {
		t.Errorf("Expected index 0, got %d", info.GetIndex())
	}

	name := info.GetName()
	t.Logf("Audio 0: %s", name)
	t.Logf("  Stereo: %v", info.IsStereo())
	t.Logf("  AVL: %v", info.HasAVL())
}

// TestIntegration_CurrentAudio tests getting current audio input
func TestIntegration_CurrentAudio(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	currentAudio, err := dev.GetCurrentAudio()
	if err != nil {
		t.Skip("Device does not support getting current audio input")
	}

	t.Logf("Current audio input: [%d] %s", currentAudio.GetIndex(), currentAudio.GetName())
	t.Logf("  Capability: 0x%x", currentAudio.GetCapability())
	t.Logf("  Mode: 0x%x", currentAudio.GetMode())
	t.Logf("  Stereo: %v", currentAudio.IsStereo())
	t.Logf("  AVL: %v", currentAudio.HasAVL())
}

// TestIntegration_AudioSelection tests switching audio inputs
func TestIntegration_AudioSelection(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current audio
	originalAudio, err := dev.GetCurrentAudio()
	if err != nil {
		t.Skip("Device does not support audio input selection")
	}
	originalIdx := originalAudio.GetIndex()

	// Get all audio inputs
	audios, err := dev.GetAudioDescriptions()
	if err != nil {
		t.Fatalf("Failed to enumerate audio inputs: %v", err)
	}

	if len(audios) < 2 {
		t.Skip("Device has only one audio input, cannot test selection")
	}

	// Try to switch to a different audio input
	var newIdx uint32 = 0
	for _, audio := range audios {
		if audio.GetIndex() != originalIdx {
			newIdx = audio.GetIndex()
			break
		}
	}

	if newIdx == originalIdx {
		t.Skip("Could not find alternative audio input")
	}

	t.Logf("Switching from audio %d to audio %d", originalIdx, newIdx)

	// Switch audio input
	err = dev.SetAudio(newIdx)
	if err != nil {
		t.Fatalf("Failed to set audio input: %v", err)
	}

	// Verify the change
	currentAudio, err := dev.GetCurrentAudio()
	if err != nil {
		t.Fatalf("Failed to get current audio after change: %v", err)
	}

	if currentAudio.GetIndex() != newIdx {
		t.Errorf("Expected audio %d, got %d", newIdx, currentAudio.GetIndex())
	}

	// Restore original audio input
	err = dev.SetAudio(originalIdx)
	if err != nil {
		t.Fatalf("Failed to restore original audio input: %v", err)
	}

	t.Logf("Successfully restored original audio input %d", originalIdx)
}

// TestIntegration_AudioMode tests setting audio mode
func TestIntegration_AudioMode(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current audio to check if it has AVL capability
	currentAudio, err := dev.GetCurrentAudio()
	if err != nil {
		t.Skip("Device does not support audio input")
	}

	if !currentAudio.HasAVL() {
		t.Skip("Current audio input does not support AVL mode")
	}

	// Try to set AVL mode
	err = dev.SetAudioMode(v4l2.AudioModeAVL)
	if err != nil {
		t.Logf("Failed to set audio mode (this may be expected): %v", err)
		return
	}

	t.Log("Successfully set audio mode to AVL")

	// Verify the mode was set
	newAudio, err := dev.GetCurrentAudio()
	if err != nil {
		t.Fatalf("Failed to get current audio after mode change: %v", err)
	}

	if newAudio.GetMode()&v4l2.AudioModeAVL == 0 {
		t.Error("Audio mode was not set to AVL")
	}
}

// TestIntegration_AudioOutputEnumeration tests audio output enumeration
func TestIntegration_AudioOutputEnumeration(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device %s: %v", devPath, err)
	}
	defer dev.Close()

	// Test GetAudioOutDescriptions
	audioOuts, err := dev.GetAudioOutDescriptions()
	if err != nil {
		t.Logf("Device does not support audio output enumeration: %v", err)
		return // Not all devices have audio outputs
	}

	if len(audioOuts) == 0 {
		t.Log("Device reports no audio outputs")
		return
	}

	t.Logf("Found %d audio output(s):", len(audioOuts))
	for _, audioOut := range audioOuts {
		t.Logf("  [%d] %s", audioOut.GetIndex(), audioOut.GetName())
		t.Logf("      Capability: 0x%x", audioOut.GetCapability())
		t.Logf("      Mode: 0x%x", audioOut.GetMode())
		t.Logf("      Stereo: %v", audioOut.IsStereo())
		t.Logf("      AVL: %v", audioOut.HasAVL())
	}
}

// TestIntegration_AudioOutputInfo tests getting specific audio output info
func TestIntegration_AudioOutputInfo(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get audio output info for index 0
	info, err := dev.GetAudioOutInfo(0)
	if err != nil {
		t.Skip("Device does not support audio output info query")
	}

	// Verify the info is for the correct index
	if info.GetIndex() != 0 {
		t.Errorf("Expected index 0, got %d", info.GetIndex())
	}

	name := info.GetName()
	t.Logf("Audio Output 0: %s", name)
	t.Logf("  Stereo: %v", info.IsStereo())
	t.Logf("  AVL: %v", info.HasAVL())
}

// TestIntegration_CurrentAudioOut tests getting current audio output
func TestIntegration_CurrentAudioOut(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	currentAudioOut, err := dev.GetCurrentAudioOut()
	if err != nil {
		t.Skip("Device does not support getting current audio output")
	}

	t.Logf("Current audio output: [%d] %s", currentAudioOut.GetIndex(), currentAudioOut.GetName())
	t.Logf("  Capability: 0x%x", currentAudioOut.GetCapability())
	t.Logf("  Mode: 0x%x", currentAudioOut.GetMode())
	t.Logf("  Stereo: %v", currentAudioOut.IsStereo())
	t.Logf("  AVL: %v", currentAudioOut.HasAVL())
}

// TestIntegration_AudioOutputSelection tests switching audio outputs
func TestIntegration_AudioOutputSelection(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current audio output
	originalAudioOut, err := dev.GetCurrentAudioOut()
	if err != nil {
		t.Skip("Device does not support audio output selection")
	}
	originalIdx := originalAudioOut.GetIndex()

	// Get all audio outputs
	audioOuts, err := dev.GetAudioOutDescriptions()
	if err != nil {
		t.Fatalf("Failed to enumerate audio outputs: %v", err)
	}

	if len(audioOuts) < 2 {
		t.Skip("Device has only one audio output, cannot test selection")
	}

	// Try to switch to a different audio output
	var newIdx uint32 = 0
	for _, audioOut := range audioOuts {
		if audioOut.GetIndex() != originalIdx {
			newIdx = audioOut.GetIndex()
			break
		}
	}

	if newIdx == originalIdx {
		t.Skip("Could not find alternative audio output")
	}

	t.Logf("Switching from audio output %d to audio output %d", originalIdx, newIdx)

	// Switch audio output
	err = dev.SetAudioOut(newIdx)
	if err != nil {
		t.Fatalf("Failed to set audio output: %v", err)
	}

	// Verify the change
	currentAudioOut, err := dev.GetCurrentAudioOut()
	if err != nil {
		t.Fatalf("Failed to get current audio output after change: %v", err)
	}

	if currentAudioOut.GetIndex() != newIdx {
		t.Errorf("Expected audio output %d, got %d", newIdx, currentAudioOut.GetIndex())
	}

	// Restore original audio output
	err = dev.SetAudioOut(originalIdx)
	if err != nil {
		t.Fatalf("Failed to restore original audio output: %v", err)
	}

	t.Logf("Successfully restored original audio output %d", originalIdx)
}

// TestIntegration_AudioOutMode tests setting audio output mode
func TestIntegration_AudioOutMode(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current audio output to check if it has AVL capability
	currentAudioOut, err := dev.GetCurrentAudioOut()
	if err != nil {
		t.Skip("Device does not support audio output")
	}

	if !currentAudioOut.HasAVL() {
		t.Skip("Current audio output does not support AVL mode")
	}

	// Try to set AVL mode
	err = dev.SetAudioOutMode(v4l2.AudioModeAVL)
	if err != nil {
		t.Logf("Failed to set audio output mode (this may be expected): %v", err)
		return
	}

	t.Log("Successfully set audio output mode to AVL")

	// Verify the mode was set
	newAudioOut, err := dev.GetCurrentAudioOut()
	if err != nil {
		t.Fatalf("Failed to get current audio output after mode change: %v", err)
	}

	if newAudioOut.GetMode()&v4l2.AudioModeAVL == 0 {
		t.Error("Audio output mode was not set to AVL")
	}
}
