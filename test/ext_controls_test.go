// +build integration

package test

import (
	"testing"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestIntegration_GetExtControls tests retrieving extended controls
func TestIntegration_GetExtControls(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device %s: %v", devPath, err)
	}
	defer dev.Close()

	// Try to get brightness and contrast using extended controls API
	ctrls := v4l2.NewExtControls()

	// Add controls to query
	ctrl1 := v4l2.NewExtControl(v4l2.CtrlBrightness)
	ctrl2 := v4l2.NewExtControl(v4l2.CtrlContrast)
	ctrls.Add(ctrl1)
	ctrls.Add(ctrl2)

	err = dev.GetExtControls(ctrls)
	if err != nil {
		// Some devices may not support extended controls
		t.Logf("Device does not support extended controls: %v", err)
		return
	}

	// If successful, retrieve values
	controls := ctrls.GetControls()
	if len(controls) != 2 {
		t.Errorf("Expected 2 controls, got %d", len(controls))
	}

	brightness := controls[0].GetValue()
	contrast := controls[1].GetValue()

	t.Logf("Brightness: %d", brightness)
	t.Logf("Contrast: %d", contrast)

}

// TestIntegration_SetExtControls tests setting extended controls
func TestIntegration_SetExtControls(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// First, get current values
	getCtrls := v4l2.NewExtControls()

	getCtrl1 := v4l2.NewExtControl(v4l2.CtrlBrightness)
	getCtrls.Add(getCtrl1)

	err = dev.GetExtControls(getCtrls)
	if err != nil {
		t.Logf("Device does not support extended controls: %v", err)
		return
	}

	originalValue := getCtrls.GetControls()[0].GetValue()
	t.Logf("Original brightness: %d", originalValue)

	// Try to set new value
	newValue := originalValue
	if originalValue < 200 {
		newValue = originalValue + 10
	} else {
		newValue = originalValue - 10
	}

	setCtrls := v4l2.NewExtControls()

	setCtrl := v4l2.NewExtControlWithValue(v4l2.CtrlBrightness, newValue)
	setCtrls.Add(setCtrl)

	err = dev.SetExtControls(setCtrls)
	if err != nil {
		t.Logf("Failed to set extended control (may be read-only or out of range): %v", err)
		return
	}

	t.Logf("Successfully set brightness to %d", newValue)

	// Verify the change
	verifyCtrls := v4l2.NewExtControls()

	verifyCtrl := v4l2.NewExtControl(v4l2.CtrlBrightness)
	verifyCtrls.Add(verifyCtrl)

	err = dev.GetExtControls(verifyCtrls)
	if err != nil {
		t.Errorf("Failed to verify control value: %v", err)
	} else {
		actualValue := verifyCtrls.GetControls()[0].GetValue()
		t.Logf("Verified brightness: %d", actualValue)
		if actualValue != newValue {
			t.Logf("Note: Value differs (set=%d, got=%d) - driver may have adjusted it", newValue, actualValue)
		}
	}

	// Restore original value
	restoreCtrls := v4l2.NewExtControls()

	restoreCtrl := v4l2.NewExtControlWithValue(v4l2.CtrlBrightness, originalValue)
	restoreCtrls.Add(restoreCtrl)
	dev.SetExtControls(restoreCtrls)

}

// TestIntegration_TryExtControls tests trying control values without applying
func TestIntegration_TryExtControls(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try a valid value first
	ctrls := v4l2.NewExtControls()

	ctrl := v4l2.NewExtControlWithValue(v4l2.CtrlBrightness, 128)
	ctrls.Add(ctrl)

	err = dev.TryExtControls(ctrls)
	if err != nil {
		t.Logf("Device does not support TRY_EXT_CTRLS: %v", err)
		return
	}

	t.Log("TRY_EXT_CTRLS succeeded for brightness=128")

	// Try an extreme value that might be out of range
	extremeCtrls := v4l2.NewExtControls()

	extremeCtrl := v4l2.NewExtControlWithValue(v4l2.CtrlBrightness, 9999)
	extremeCtrls.Add(extremeCtrl)

	err = dev.TryExtControls(extremeCtrls)
	if err != nil {
		t.Logf("TRY_EXT_CTRLS correctly rejected out-of-range value: %v", err)
	} else {
		t.Log("TRY_EXT_CTRLS accepted extreme value (driver may clamp)")
	}

}

// TestIntegration_ExtControlsWithClass tests extended controls with control class
func TestIntegration_ExtControlsWithClass(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try with User control class
	ctrls := v4l2.NewExtControlsWithClass(v4l2.CtrlClassUser)

	// Verify class is set
	if ctrls.GetClass() != v4l2.CtrlClassUser {
		t.Errorf("Expected class User (0x%08x), got 0x%08x", v4l2.CtrlClassUser, ctrls.GetClass())
	}

	ctrl := v4l2.NewExtControl(v4l2.CtrlBrightness)
	ctrls.Add(ctrl)

	err = dev.GetExtControls(ctrls)
	if err != nil {
		t.Logf("Device does not support control class filtering: %v", err)
		return
	}

	brightness := ctrls.GetControls()[0].GetValue()
	t.Logf("Brightness (from User class): %d", brightness)

}

// TestIntegration_ExtControlsMultiple tests getting multiple controls atomically
func TestIntegration_ExtControlsMultiple(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to get multiple common controls
	ctrls := v4l2.NewExtControls()

	controlIDs := []v4l2.CtrlID{
		v4l2.CtrlBrightness,
		v4l2.CtrlContrast,
		v4l2.CtrlSaturation,
		v4l2.CtrlHue,
		v4l2.CtrlGain,
	}

	var controls []*v4l2.ExtControl
	for _, id := range controlIDs {
		ctrl := v4l2.NewExtControl(id)
		ctrls.Add(ctrl)
		controls = append(controls, ctrl)
	}

	err = dev.GetExtControls(ctrls)
	if err != nil {
		t.Logf("Failed to get all controls (some may not be supported): %v", err)
		errorIdx := ctrls.GetErrorIndex()
		t.Logf("Error index: %d (control that failed)", errorIdx)
		return
	}

	// Log all values
	t.Logf("Retrieved %d controls:", ctrls.Count())
	retrievedControls := ctrls.GetControls()
	for _, ctrl := range retrievedControls {
		t.Logf("  Control ID 0x%08x: %d", ctrl.GetID(), ctrl.GetValue())
	}
}

// TestIntegration_ExtControlsErrorHandling tests error handling
func TestIntegration_ExtControlsErrorHandling(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to get a control that definitely doesn't exist
	ctrls := v4l2.NewExtControls()

	// Use an invalid/unsupported control ID
	invalidCtrl := v4l2.NewExtControl(0xFFFFFFFF)
	ctrls.Add(invalidCtrl)

	err = dev.GetExtControls(ctrls)
	if err == nil {
		t.Log("Device accepted invalid control (unexpected)")
	} else {
		t.Logf("Device correctly rejected invalid control: %v", err)
		errorIdx := ctrls.GetErrorIndex()
		t.Logf("Error occurred at index: %d", errorIdx)
	}

}

// TestIntegration_SubscribeControlEvent tests subscribing to control change events
func TestIntegration_SubscribeControlEvent(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Subscribe to brightness control changes
	sub := v4l2.NewControlEventSubscription(v4l2.CtrlBrightness)
	sub.SetFlags(v4l2.EventSubFlagSendInitial)

	err = dev.SubscribeEvent(sub)
	if err != nil {
		t.Logf("Device does not support event subscription: %v", err)
		return
	}

	t.Logf("Successfully subscribed to brightness control events")

	// Try to dequeue an event (should get initial event due to EventSubFlagSendInitial)
	event, err := dev.DequeueEvent()
	if err != nil {
		t.Logf("No initial event available: %v", err)
	} else {
		t.Logf("Received event type: %d", event.GetType())
		if event.GetType() == v4l2.EventCtrl {
			ctrlData := event.GetCtrlData()
			t.Logf("  Control ID: %d", event.GetID())
			t.Logf("  Changes: 0x%08x", ctrlData.Changes)
			t.Logf("  Value: %d", ctrlData.Value)
		}
	}

	// Unsubscribe
	err = dev.UnsubscribeEvent(sub)
	if err != nil {
		t.Errorf("Failed to unsubscribe: %v", err)
	} else {
		t.Log("Successfully unsubscribed from event")
	}
}

// TestIntegration_SubscribeAllEvents tests subscribing to all events
func TestIntegration_SubscribeAllEvents(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Subscribe to all events
	sub := v4l2.NewEventSubscription(v4l2.EventAll)

	err = dev.SubscribeEvent(sub)
	if err != nil {
		t.Logf("Device does not support subscribing to all events: %v", err)
		return
	}

	t.Log("Successfully subscribed to all events")

	// Unsubscribe
	err = dev.UnsubscribeEvent(sub)
	if err != nil {
		t.Errorf("Failed to unsubscribe from all events: %v", err)
	} else {
		t.Log("Successfully unsubscribed from all events")
	}
}

// TestIntegration_EventTypes tests different event type subscriptions
func TestIntegration_EventTypes(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	eventTypes := []struct {
		name      string
		eventType v4l2.EventType
	}{
		{"VSync", v4l2.EventVSync},
		{"EOS", v4l2.EventEOS},
		{"FrameSync", v4l2.EventFrameSync},
		{"SourceChange", v4l2.EventSourceChange},
		{"MotionDet", v4l2.EventMotionDet},
	}

	for _, et := range eventTypes {
		t.Run(et.name, func(t *testing.T) {
			sub := v4l2.NewEventSubscription(et.eventType)

			err := dev.SubscribeEvent(sub)
			if err != nil {
				t.Logf("Device does not support %s events: %v", et.name, err)
				return
			}

			t.Logf("Successfully subscribed to %s events", et.name)

			// Unsubscribe
			err = dev.UnsubscribeEvent(sub)
			if err != nil {
				t.Errorf("Failed to unsubscribe from %s events: %v", et.name, err)
			}
		})
	}
}

// TestIntegration_HighLevelBrightness tests high-level brightness control methods
func TestIntegration_HighLevelBrightness(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current brightness
	originalBrightness, err := dev.GetBrightness()
	if err != nil {
		t.Logf("Device does not support brightness control: %v", err)
		return
	}

	t.Logf("Original brightness: %d", originalBrightness)

	// Calculate new value
	newBrightness := originalBrightness
	if originalBrightness < 200 {
		newBrightness = originalBrightness + 10
	} else {
		newBrightness = originalBrightness - 10
	}

	// Set new brightness
	err = dev.SetBrightness(newBrightness)
	if err != nil {
		t.Logf("Failed to set brightness (may be read-only): %v", err)
		return
	}

	t.Logf("Successfully set brightness to %d", newBrightness)

	// Verify the change
	verifyBrightness, err := dev.GetBrightness()
	if err != nil {
		t.Errorf("Failed to verify brightness: %v", err)
		return
	}

	if verifyBrightness != newBrightness {
		t.Errorf("Brightness verification failed: got %d, want %d", verifyBrightness, newBrightness)
	}

	// Restore original value
	err = dev.SetBrightness(originalBrightness)
	if err != nil {
		t.Logf("Failed to restore original brightness: %v", err)
	} else {
		t.Logf("Successfully restored brightness to %d", originalBrightness)
	}
}

// TestIntegration_HighLevelContrast tests high-level contrast control methods
func TestIntegration_HighLevelContrast(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current contrast
	originalContrast, err := dev.GetContrast()
	if err != nil {
		t.Logf("Device does not support contrast control: %v", err)
		return
	}

	t.Logf("Original contrast: %d", originalContrast)

	// Calculate new value
	newContrast := originalContrast
	if originalContrast < 200 {
		newContrast = originalContrast + 10
	} else {
		newContrast = originalContrast - 10
	}

	// Set new contrast
	err = dev.SetContrast(newContrast)
	if err != nil {
		t.Logf("Failed to set contrast (may be read-only): %v", err)
		return
	}

	t.Logf("Successfully set contrast to %d", newContrast)

	// Restore original value
	err = dev.SetContrast(originalContrast)
	if err != nil {
		t.Logf("Failed to restore original contrast: %v", err)
	}
}

// TestIntegration_HighLevelSaturation tests high-level saturation control methods
func TestIntegration_HighLevelSaturation(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current saturation
	originalSaturation, err := dev.GetSaturation()
	if err != nil {
		t.Logf("Device does not support saturation control: %v", err)
		return
	}

	t.Logf("Original saturation: %d", originalSaturation)

	// Calculate new value
	newSaturation := originalSaturation
	if originalSaturation < 200 {
		newSaturation = originalSaturation + 10
	} else {
		newSaturation = originalSaturation - 10
	}

	// Set new saturation
	err = dev.SetSaturation(newSaturation)
	if err != nil {
		t.Logf("Failed to set saturation (may be read-only): %v", err)
		return
	}

	t.Logf("Successfully set saturation to %d", newSaturation)

	// Restore original value
	err = dev.SetSaturation(originalSaturation)
	if err != nil {
		t.Logf("Failed to restore original saturation: %v", err)
	}
}

// TestIntegration_HighLevelHue tests high-level hue control methods
func TestIntegration_HighLevelHue(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get current hue
	originalHue, err := dev.GetHue()
	if err != nil {
		t.Logf("Device does not support hue control: %v", err)
		return
	}

	t.Logf("Original hue: %d", originalHue)

	// Calculate new value
	newHue := originalHue + 10

	// Set new hue
	err = dev.SetHue(newHue)
	if err != nil {
		t.Logf("Failed to set hue (may be read-only or out of range): %v", err)
		return
	}

	t.Logf("Successfully set hue to %d", newHue)

	// Restore original value
	err = dev.SetHue(originalHue)
	if err != nil {
		t.Logf("Failed to restore original hue: %v", err)
	}
}

// TestIntegration_HighLevelMultipleControls tests setting multiple controls at once
func TestIntegration_HighLevelMultipleControls(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Get original values
	origBrightness, err := dev.GetBrightness()
	if err != nil {
		t.Logf("Device does not support brightness control: %v", err)
		return
	}

	origContrast, err := dev.GetContrast()
	if err != nil {
		t.Logf("Device does not support contrast control: %v", err)
		return
	}

	t.Logf("Original values - Brightness: %d, Contrast: %d", origBrightness, origContrast)

	// Set multiple controls atomically using AddValue
	ctrls := v4l2.NewExtControls()
	ctrls.AddValue(v4l2.CtrlBrightness, origBrightness+5)
	ctrls.AddValue(v4l2.CtrlContrast, origContrast+5)

	err = dev.SetExtControls(ctrls)
	if err != nil {
		t.Logf("Failed to set multiple controls: %v", err)
		return
	}

	t.Log("Successfully set multiple controls atomically")

	// Verify changes
	newBrightness, _ := dev.GetBrightness()
	newContrast, _ := dev.GetContrast()
	t.Logf("New values - Brightness: %d, Contrast: %d", newBrightness, newContrast)

	// Restore original values
	ctrls = v4l2.NewExtControls()
	ctrls.AddValue(v4l2.CtrlBrightness, origBrightness)
	ctrls.AddValue(v4l2.CtrlContrast, origContrast)
	dev.SetExtControls(ctrls)
}
