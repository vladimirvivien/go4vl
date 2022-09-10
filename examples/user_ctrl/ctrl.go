package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	dev "github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	controls = map[string]v4l2.CtrlID{
		"brightness": v4l2.CtrlBrightness,
		"contrast":   v4l2.CtrlContrast,
	}
)

func main() {
	devName := "/dev/video0"
	flag.StringVar(&devName, "d", devName, "device name (path)")
	var list bool
	flag.BoolVar(&list, "list", list, "List current device controls")
	var ctrlName string
	flag.StringVar(&ctrlName, "c", ctrlName, fmt.Sprintf("Contrl name to set or get (supported %v)", controls))
	ctrlVal := math.MinInt32
	flag.IntVar(&ctrlVal, "v", ctrlVal, fmt.Sprintf("Value for selected control (supported %v)", controls))

	flag.Parse()

	// open device
	device, err := dev.Open(devName)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()

	if len(os.Args) < 2 || list {
		listUserControls(device)
		os.Exit(0)
	}

	ctrlID, ok := controls[ctrlName]
	if !ok {
		fmt.Printf("Program does not support ctrl [%s]; supported ctrls: %#v\n", ctrlName, controls)
		os.Exit(1)
	}

	if ctrlName != "" {
		if ctrlVal != math.MinInt32 {
			if err := setUserControlValue(device, ctrlID, ctrlVal); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		ctrl, err := device.GetControl(ctrlID)
		if err != nil {
			fmt.Printf("query controls: %s\n", err)
			os.Exit(1)
		}
		printUserControl(ctrl)
	}
}

func setUserControlValue(device *dev.Device, ctrlID v4l2.CtrlID, val int) error {
	if ctrlID == 0 {
		return fmt.Errorf("invalid control specified")
	}
	return device.SetControlValue(ctrlID, v4l2.CtrlValue(val))
}

func listUserControls(device *dev.Device) {
	ctrls, err := device.QueryAllControls()
	if err != nil {
		log.Fatalf("query controls: %s", err)
	}

	for _, ctrl := range ctrls {
		printUserControl(ctrl)
	}
}

func printUserControl(ctrl v4l2.Control) {
	fmt.Printf("Control id (%d) name: %s\t[min: %d; max: %d; step: %d; default: %d current_val: %d]\n",
		ctrl.ID, ctrl.Name, ctrl.Minimum, ctrl.Maximum, ctrl.Step, ctrl.Default, ctrl.Value)

	if ctrl.IsMenu() {
		menus, err := ctrl.GetMenuItems()
		if err != nil {
			return
		}

		for _, m := range menus {
			fmt.Printf("\tMenu items for %s: %#v\n", ctrl.Name, m)
		}
	}

}
