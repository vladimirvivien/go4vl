package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	devPath   = flag.String("d", "/dev/video0", "device path")
	listClass = flag.String("class", "", "list controls for specific class (user, camera, flash, image-source, image-proc, codec)")
	listAll   = flag.Bool("all", false, "list all controls")
	getCtrl   = flag.Uint("get", 0, "get control value by ID")
	setCtrl   = flag.Uint("set", 0, "set control ID")
	setValue  = flag.Int("value", 0, "value to set")
)

func main() {
	flag.Parse()

	dev, err := device.Open(*devPath)
	if err != nil {
		fmt.Printf("Failed to open device %s: %v\n", *devPath, err)
		os.Exit(1)
	}
	defer dev.Close()

	fmt.Printf("V4L2 Control Reference - Device: %s\n", *devPath)
	fmt.Println("=============================================")
	fmt.Println()

	// Handle specific actions
	if *getCtrl != 0 {
		getControlValue(dev, v4l2.CtrlID(*getCtrl))
		return
	}

	if *setCtrl != 0 {
		setControlValue(dev, v4l2.CtrlID(*setCtrl), int32(*setValue))
		return
	}

	// List controls
	if *listClass != "" {
		listControlsByClass(dev, *listClass)
		return
	}

	if *listAll {
		listAllControls(dev)
		return
	}

	// Default: show overview
	showOverview(dev)
}

func showOverview(dev *device.Device) {
	fmt.Println("Control Classes Overview")
	fmt.Println("------------------------")
	fmt.Println()

	classes := []struct {
		name  string
		class v4l2.CtrlClass
		desc  string
	}{
		{"User", v4l2.CtrlClassUser, "Basic picture controls (brightness, contrast, etc.)"},
		{"Camera", v4l2.CtrlClassCamera, "Camera controls (exposure, focus, zoom, pan/tilt)"},
		{"Flash", v4l2.CtrlClassFlash, "Flash and LED controls"},
		{"JPEG", v4l2.CtrlClassJPEG, "JPEG compression settings"},
		{"Image Source", v4l2.CtrlClassImageSource, "Image sensor controls (gain, blanking, test patterns)"},
		{"Image Processing", v4l2.CtrlClassImageProcessing, "Image processing (pixel rate, deinterlacing)"},
		{"Codec", v4l2.CtrlClassCodec, "Codec controls (bitrate, GOP, profiles)"},
		{"Codec Stateless", v4l2.CtrlClassCodecStateless, "Stateless codec controls (H.264, VP8, MPEG2)"},
		{"Digital Video", v4l2.CtrlClassDigitalVideo, "Digital video timing controls"},
		{"Detection", v4l2.CtrlClassDetection, "Detection controls (motion, face)"},
		{"Colorimetry", v4l2.CtrlClassColorimitry, "Color space and HDR controls"},
	}

	for _, c := range classes {
		count := countControlsInClass(dev, c.class)
		status := "Not supported"
		if count > 0 {
			status = fmt.Sprintf("%d controls", count)
		}
		fmt.Printf("  %-18s [%-15s] - %s\n", c.name, status, c.desc)
	}

	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  -all                  List all available controls")
	fmt.Println("  -class <name>         List controls in specific class")
	fmt.Println("  -get <id>             Get control value")
	fmt.Println("  -set <id> -value <n>  Set control value")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  control_reference -class camera")
	fmt.Println("  control_reference -get 9963776")
	fmt.Println("  control_reference -set 9963776 -value 100")
}

func countControlsInClass(dev *device.Device, class v4l2.CtrlClass) int {
	count := 0
	// Try to query controls starting from the class base
	baseID := v4l2.CtrlID(class | 0x900)

	for id := baseID; id < baseID+100; id++ {
		ctrl, err := v4l2.QueryControlInfo(dev.Fd(), id)
		if err != nil {
			continue
		}
		if ctrl.Flags&v4l2.CtrlFlagDisabled == 0 {
			count++
		}
	}
	return count
}

func listAllControls(dev *device.Device) {
	fmt.Println("All Available Controls")
	fmt.Println("======================")
	fmt.Println()

	classes := []struct {
		name  string
		class v4l2.CtrlClass
	}{
		{"User", v4l2.CtrlClassUser},
		{"Camera", v4l2.CtrlClassCamera},
		{"Flash", v4l2.CtrlClassFlash},
		{"JPEG", v4l2.CtrlClassJPEG},
		{"Image Source", v4l2.CtrlClassImageSource},
		{"Image Processing", v4l2.CtrlClassImageProcessing},
		{"Codec", v4l2.CtrlClassCodec},
		{"Codec Stateless", v4l2.CtrlClassCodecStateless},
		{"Digital Video", v4l2.CtrlClassDigitalVideo},
		{"Detection", v4l2.CtrlClassDetection},
		{"Colorimetry", v4l2.CtrlClassColorimitry},
	}

	for _, c := range classes {
		controls := getControlsInClass(dev, c.class)
		if len(controls) > 0 {
			fmt.Printf("%s Controls (%d):\n", c.name, len(controls))
			fmt.Println(strings.Repeat("-", 80))
			for _, ctrl := range controls {
				displayControl(dev, ctrl)
			}
			fmt.Println()
		}
	}
}

func listControlsByClass(dev *device.Device, className string) {
	classMap := map[string]v4l2.CtrlClass{
		"user":          v4l2.CtrlClassUser,
		"camera":        v4l2.CtrlClassCamera,
		"flash":         v4l2.CtrlClassFlash,
		"jpeg":          v4l2.CtrlClassJPEG,
		"image-source":  v4l2.CtrlClassImageSource,
		"image-proc":    v4l2.CtrlClassImageProcessing,
		"codec":         v4l2.CtrlClassCodec,
		"codec-stateless": v4l2.CtrlClassCodecStateless,
		"dv":            v4l2.CtrlClassDigitalVideo,
		"detection":     v4l2.CtrlClassDetection,
		"colorimetry":   v4l2.CtrlClassColorimitry,
	}

	class, ok := classMap[strings.ToLower(className)]
	if !ok {
		fmt.Printf("Unknown class '%s'. Valid classes:\n", className)
		var names []string
		for name := range classMap {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Printf("  - %s\n", name)
		}
		os.Exit(1)
	}

	controls := getControlsInClass(dev, class)
	if len(controls) == 0 {
		fmt.Printf("No controls found in class '%s'\n", className)
		return
	}

	fmt.Printf("%s Controls (%d):\n", className, len(controls))
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	for _, ctrl := range controls {
		displayControl(dev, ctrl)
	}
}

func getControlsInClass(dev *device.Device, class v4l2.CtrlClass) []v4l2.Control {
	var controls []v4l2.Control

	// Try to query controls starting from the class base
	baseID := v4l2.CtrlID(class | 0x900)

	for id := baseID; id < baseID+200; id++ {
		ctrl, err := v4l2.QueryControlInfo(dev.Fd(), id)
		if err != nil {
			continue
		}
		if ctrl.Flags&v4l2.CtrlFlagDisabled == 0 {
			controls = append(controls, ctrl)
		}
	}

	return controls
}

func displayControl(dev *device.Device, ctrl v4l2.Control) {
	fmt.Printf("  ID: 0x%08x (%d)\n", ctrl.ID, ctrl.ID)
	fmt.Printf("  Name: %s\n", ctrl.Name)
	fmt.Printf("  Type: %s\n", getTypeName(ctrl.Type))

	// Get current value for non-write-only controls
	if ctrl.Flags&v4l2.CtrlFlagWriteOnly == 0 {
		if val, err := v4l2.GetControlValue(dev.Fd(), ctrl.ID); err == nil {
			fmt.Printf("  Current Value: %d\n", val)
		}
	}

	// Show range for integer controls
	if ctrl.Type == v4l2.CtrlTypeInt || ctrl.Type == v4l2.CtrlTypeInt64 {
		fmt.Printf("  Range: %d to %d (step: %d)\n", ctrl.Minimum, ctrl.Maximum, ctrl.Step)
		fmt.Printf("  Default: %d\n", ctrl.Default)
	}

	// Show menu items for menu controls
	if ctrl.IsMenu() {
		if items, err := ctrl.GetMenuItems(); err == nil && len(items) > 0 {
			fmt.Printf("  Menu Items:\n")
			for _, item := range items {
				fmt.Printf("    [%d] %s\n", item.Index, item.Name)
			}
		}
	}

	// Show flags
	var flags []string
	if ctrl.Flags&v4l2.CtrlFlagDisabled != 0 {
		flags = append(flags, "disabled")
	}
	if ctrl.Flags&v4l2.CtrlFlagGrabbed != 0 {
		flags = append(flags, "grabbed")
	}
	if ctrl.Flags&v4l2.CtrlFlagReadOnly != 0 {
		flags = append(flags, "read-only")
	}
	if ctrl.Flags&v4l2.CtrlFlagWriteOnly != 0 {
		flags = append(flags, "write-only")
	}
	if ctrl.Flags&v4l2.CtrlFlagVolatile != 0 {
		flags = append(flags, "volatile")
	}
	if ctrl.Flags&v4l2.CtrlFlagInactive != 0 {
		flags = append(flags, "inactive")
	}
	if len(flags) > 0 {
		fmt.Printf("  Flags: %s\n", strings.Join(flags, ", "))
	}

	fmt.Println()
}

func getTypeName(t v4l2.CtrlType) string {
	types := map[v4l2.CtrlType]string{
		v4l2.CtrlTypeInt:         "Integer",
		v4l2.CtrlTypeBool:        "Boolean",
		v4l2.CtrlTypeMenu:        "Menu",
		v4l2.CtrlTypeButton:      "Button",
		v4l2.CtrlTypeInt64:       "Integer64",
		v4l2.CtrlTypeClass:       "Class",
		v4l2.CtrlTypeString:      "String",
		v4l2.CtrlTypeBitMask:     "Bitmask",
		v4l2.CtrlTypeIntegerMenu: "IntegerMenu",
		v4l2.CtrlTypeU8:          "U8",
		v4l2.CtrlTypeU16:         "U16",
		v4l2.CtrlTypeU32:         "U32",
	}
	if name, ok := types[t]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", t)
}

func getControlValue(dev *device.Device, id v4l2.CtrlID) {
	ctrl, err := v4l2.QueryControlInfo(dev.Fd(), id)
	if err != nil {
		fmt.Printf("Error querying control 0x%08x: %v\n", id, err)
		os.Exit(1)
	}

	fmt.Printf("Control: %s (0x%08x)\n", ctrl.Name, ctrl.ID)
	fmt.Printf("Type: %s\n", getTypeName(ctrl.Type))

	if ctrl.Flags&v4l2.CtrlFlagWriteOnly != 0 {
		fmt.Println("This control is write-only")
		return
	}

	val, err := v4l2.GetControlValue(dev.Fd(), id)
	if err != nil {
		fmt.Printf("Error getting value: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current Value: %d\n", val)

	if ctrl.Type == v4l2.CtrlTypeInt || ctrl.Type == v4l2.CtrlTypeInt64 {
		fmt.Printf("Range: %d to %d (step: %d)\n", ctrl.Minimum, ctrl.Maximum, ctrl.Step)
	}

	if ctrl.IsMenu() {
		if items, err := ctrl.GetMenuItems(); err == nil {
			for _, item := range items {
				if item.Index == uint32(val) {
					fmt.Printf("Menu Selection: %s\n", item.Name)
					break
				}
			}
		}
	}
}

func setControlValue(dev *device.Device, id v4l2.CtrlID, value int32) {
	ctrl, err := v4l2.QueryControlInfo(dev.Fd(), id)
	if err != nil {
		fmt.Printf("Error querying control 0x%08x: %v\n", id, err)
		os.Exit(1)
	}

	fmt.Printf("Setting control: %s (0x%08x)\n", ctrl.Name, ctrl.ID)

	if ctrl.Flags&v4l2.CtrlFlagReadOnly != 0 {
		fmt.Println("Error: This control is read-only")
		os.Exit(1)
	}

	// Validate range for integer controls
	if ctrl.Type == v4l2.CtrlTypeInt {
		if value < ctrl.Minimum || value > ctrl.Maximum {
			fmt.Printf("Error: Value %d out of range [%d, %d]\n", value, ctrl.Minimum, ctrl.Maximum)
			os.Exit(1)
		}
	}

	err = v4l2.SetControlValue(dev.Fd(), id, value)
	if err != nil {
		fmt.Printf("Error setting value: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully set to: %d\n", value)

	// Verify
	if ctrl.Flags&v4l2.CtrlFlagWriteOnly == 0 {
		if newVal, err := v4l2.GetControlValue(dev.Fd(), id); err == nil {
			if newVal == value {
				fmt.Println("Verification: OK")
			} else {
				fmt.Printf("Verification: Value is now %d (driver may have adjusted it)\n", newVal)
			}
		}
	}
}
