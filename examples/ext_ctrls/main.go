package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	var devName string
	var listOnly bool
	var setBrightness int
	var setContrast int
	var setSaturation int
	var setHue int
	var interactive bool

	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.BoolVar(&listOnly, "l", false, "list all extended controls and exit")
	flag.IntVar(&setBrightness, "brightness", -1, "set brightness value")
	flag.IntVar(&setContrast, "contrast", -1, "set contrast value")
	flag.IntVar(&setSaturation, "saturation", -1, "set saturation value")
	flag.IntVar(&setHue, "hue", -1, "set hue value")
	flag.BoolVar(&interactive, "i", false, "interactive mode (adjust controls interactively)")
	flag.Parse()

	dev, err := device.Open(devName)
	if err != nil {
		log.Fatalf("Failed to open device %s: %v", devName, err)
	}
	defer dev.Close()

	fmt.Printf("Device: %s\n", dev.Name())
	fmt.Printf("Driver: %s\n", dev.Capability().Driver)
	fmt.Printf("Card: %s\n\n", dev.Capability().Card)

	// List all extended controls
	if listOnly {
		listAllControls(dev)
		return
	}

	// Show current values
	showCurrentValues(dev)

	// Set individual controls if specified
	if setBrightness >= 0 {
		setControlValue(dev, "Brightness", setBrightness, dev.SetBrightness)
	}
	if setContrast >= 0 {
		setControlValue(dev, "Contrast", setContrast, dev.SetContrast)
	}
	if setSaturation >= 0 {
		setControlValue(dev, "Saturation", setSaturation, dev.SetSaturation)
	}
	if setHue >= 0 {
		setControlValue(dev, "Hue", setHue, dev.SetHue)
	}

	// Interactive mode
	if interactive {
		runInteractiveMode(dev)
	}

	// If any values were set, show the new values
	if setBrightness >= 0 || setContrast >= 0 || setSaturation >= 0 || setHue >= 0 {
		fmt.Println("\nNew values:")
		showCurrentValues(dev)
	}
}

func listAllControls(dev *device.Device) {
	fmt.Println("Extended Controls:")
	fmt.Println("================================================================================")

	ctrls, err := v4l2.QueryAllExtControls(dev.Fd())
	if err != nil {
		fmt.Printf("Device does not support extended controls: %v\n", err)
		fmt.Println("\nNote: Extended controls are typically found on:")
		fmt.Println("  - Webcams")
		fmt.Println("  - Video capture cards")
		fmt.Println("  - USB cameras")
		fmt.Println("  - Hardware encoders/decoders")
		return
	}

	if len(ctrls) == 0 {
		fmt.Println("No extended controls found")
		return
	}

	for _, ctrl := range ctrls {
		printControl(ctrl)

		if ctrl.IsMenu() {
			menus, err := ctrl.GetMenuItems()
			if err == nil {
				for _, m := range menus {
					fmt.Printf("    (%d) %s\n", m.Index, m.Name)
				}
			}
		}
		fmt.Println()
	}

	fmt.Printf("Total controls: %d\n", len(ctrls))
}

func printControl(ctrl v4l2.Control) {
	fmt.Printf("Control: %s (ID: 0x%08x)\n", ctrl.Name, ctrl.ID)
	fmt.Printf("  Type: %s\n", getControlTypeName(ctrl.Type))

	if !ctrl.IsMenu() {
		fmt.Printf("  Range: [%d - %d], Step: %d, Default: %d\n",
			ctrl.Minimum, ctrl.Maximum, ctrl.Step, ctrl.Default)
	}

	flags := getControlFlags(ctrl)
	if len(flags) > 0 {
		fmt.Printf("  Flags: %s\n", flags)
	}

	if ctrl.Value != 0 || ctrl.Type == v4l2.CtrlTypeInt {
		fmt.Printf("  Current: %d\n", ctrl.Value)
	}
}

func getControlTypeName(ctrlType v4l2.CtrlType) string {
	switch ctrlType {
	case v4l2.CtrlTypeInt:
		return "Integer"
	case v4l2.CtrlTypeBool:
		return "Boolean"
	case v4l2.CtrlTypeMenu:
		return "Menu"
	case v4l2.CtrlTypeButton:
		return "Button"
	case v4l2.CtrlTypeInt64:
		return "Integer64"
	case v4l2.CtrlTypeClass:
		return "Control Class"
	case v4l2.CtrlTypeString:
		return "String"
	case v4l2.CtrlTypeBitMask:
		return "Bitmask"
	case v4l2.CtrlTypeIntegerMenu:
		return "Integer Menu"
	default:
		return fmt.Sprintf("Unknown (%d)", ctrlType)
	}
}

func getControlFlags(ctrl v4l2.Control) string {
	// Control flags are stored in the flags field (private)
	// For now, we'll just indicate if it's a special type
	flags := ""
	if ctrl.Type == v4l2.CtrlTypeButton {
		flags += "Button "
	}
	if ctrl.IsMenu() {
		flags += "Menu "
	}
	return flags
}

func showCurrentValues(dev *device.Device) {
	fmt.Println("Current Control Values:")
	fmt.Println("================================================================================")

	// Try to get common controls using high-level API
	controls := []struct {
		name   string
		getter func() (int32, error)
	}{
		{"Brightness", dev.GetBrightness},
		{"Contrast", dev.GetContrast},
		{"Saturation", dev.GetSaturation},
		{"Hue", dev.GetHue},
	}

	anySuccess := false
	for _, ctrl := range controls {
		value, err := ctrl.getter()
		if err == nil {
			fmt.Printf("  %-12s: %d\n", ctrl.name, value)
			anySuccess = true
		}
	}

	if !anySuccess {
		fmt.Println("  No standard controls supported by this device")
		fmt.Println("  Use -l flag to see all available controls")
	}
	fmt.Println()
}

func setControlValue(dev *device.Device, name string, value int, setter func(int32) error) {
	fmt.Printf("Setting %s to %d...\n", name, value)
	err := setter(int32(value))
	if err != nil {
		fmt.Printf("  ✗ Failed to set %s: %v\n", name, err)
	} else {
		fmt.Printf("  ✓ Successfully set %s to %d\n", name, value)
	}
}

func runInteractiveMode(dev *device.Device) {
	fmt.Println("\nInteractive Mode")
	fmt.Println("================================================================================")
	fmt.Println("Demonstrating atomic multi-control operations...")

	// Get current values
	brightness, err1 := dev.GetBrightness()
	contrast, err2 := dev.GetContrast()

	if err1 != nil || err2 != nil {
		fmt.Println("Device doesn't support brightness/contrast controls")
		return
	}

	fmt.Printf("\nCurrent: Brightness=%d, Contrast=%d\n", brightness, contrast)

	// Example 1: Set multiple controls atomically using AddValue
	fmt.Println("\nExample 1: Setting multiple controls atomically...")
	ctrls := v4l2.NewExtControls()
	ctrls.AddValue(v4l2.CtrlBrightness, brightness+10)
	ctrls.AddValue(v4l2.CtrlContrast, contrast+10)

	err := dev.SetExtControls(ctrls)
	if err != nil {
		fmt.Printf("  Failed: %v\n", err)
	} else {
		fmt.Printf("  ✓ Successfully set both controls atomically\n")
		newBrightness, _ := dev.GetBrightness()
		newContrast, _ := dev.GetContrast()
		fmt.Printf("  New values: Brightness=%d, Contrast=%d\n", newBrightness, newContrast)
	}

	// Example 2: Try controls before applying (validation)
	fmt.Println("\nExample 2: Validating control values before applying...")
	testCtrls := v4l2.NewExtControls()
	testCtrls.AddValue(v4l2.CtrlBrightness, brightness+20)
	testCtrls.AddValue(v4l2.CtrlContrast, contrast+20)

	err = dev.TryExtControls(testCtrls)
	if err != nil {
		fmt.Printf("  Values would be rejected: %v\n", err)
	} else {
		fmt.Printf("  ✓ Values are valid and can be applied\n")

		// Apply them
		applyCtrls := v4l2.NewExtControls()
		applyCtrls.AddValue(v4l2.CtrlBrightness, brightness+20)
		applyCtrls.AddValue(v4l2.CtrlContrast, contrast+20)
		dev.SetExtControls(applyCtrls)
		fmt.Printf("  ✓ Applied validated values\n")
	}

	// Restore original values
	fmt.Println("\nRestoring original values...")
	restoreCtrls := v4l2.NewExtControls()
	restoreCtrls.AddValue(v4l2.CtrlBrightness, brightness)
	restoreCtrls.AddValue(v4l2.CtrlContrast, contrast)

	err = dev.SetExtControls(restoreCtrls)
	if err != nil {
		fmt.Printf("  Failed to restore: %v\n", err)
	} else {
		fmt.Printf("  ✓ Restored original values\n")
	}
}
