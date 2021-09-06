package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vladimirvivien/go4vl/v4l2"
)

var template = "\t%-24s : %s\n"

func main() {
	var devName string
	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.Parse()
	device, err := v4l2.Open(devName)
	if err != nil {
		log.Fatal(err)
	}
	defer device.Close()

	if err := printDeviceDriverInfo(device); err != nil {
		log.Fatal(err)
	}

	if err := printVideoInputInfo(device); err != nil {
		log.Fatal(err)
	}

	if err := printFormatInfo(device); err != nil {
		log.Fatal(err)
	}
}

func printDeviceDriverInfo(dev *v4l2.Device) error {
	caps, err := dev.GetCapability()
	if err != nil {
		return fmt.Errorf("driver info: %w", err)
	}

	// print driver info
	fmt.Println("Device Info:")
	fmt.Printf(template, "Driver name", caps.DriverName())
	fmt.Printf(template, "Card name", caps.CardName())
	fmt.Printf(template, "Bus info", caps.BusInfo())

	verVal := caps.GetVersion()
	version := fmt.Sprintf("%d.%d.%d", verVal>>16, (verVal>>8)&0xff, verVal&0xff)
	fmt.Printf(template, "Driver version", version)

	fmt.Printf("\t%-16s : %0x\n", "Driver capabilities", caps.GetCapabilities())
	for _, desc := range caps.GetDriverCapDescriptions() {
		fmt.Printf("\t\t%s\n", desc.Desc)
	}

	fmt.Printf("\t%-16s : %0x\n", "Device capabilities", caps.GetCapabilities())
	for _, desc := range caps.GetDeviceCapDescriptions() {
		fmt.Printf("\t\t%s\n", desc.Desc)
	}

	return nil
}

func printVideoInputInfo(dev *v4l2.Device) error {
	// first get current input
	index, err := dev.GetVideoInputIndex()
	if err != nil {
		return fmt.Errorf("video input info: %w", err)
	}

	fmt.Printf("Video input: %d", index)

	// get the input info
	info, err := dev.GetVideoInputInfo(uint32(index))
	if err != nil {
		return fmt.Errorf("video input info: %w", err)
	}

	// print info
	fmt.Printf(" (%s : %s)\n", info.GetName(), v4l2.InputStatuses[info.GetStatus()])

	return nil
}

func printFormatInfo(dev *v4l2.Device) error {
	pixFmt, err := dev.GetPixFormat()
	if err != nil {
		return fmt.Errorf("video capture format: %w", err)
	}
	fmt.Println("Format video capture:")
	fmt.Printf(template, "WidthxHeight", fmt.Sprintf("%dx%d", pixFmt.Width, pixFmt.Height))
	fmt.Printf(template, "Pixel format", v4l2.PixelFormats[pixFmt.PixelFormat])
	fmt.Printf(template, "Field", v4l2.Fields[pixFmt.Field])
	fmt.Printf(template, "Bytes per line", fmt.Sprintf("%d",pixFmt.BytesPerLine))
	fmt.Printf(template, "Size image", fmt.Sprintf("%d", pixFmt.SizeImage))
	fmt.Printf(template, "Colorspace", v4l2.Colorspaces[pixFmt.Colorspace])

	return nil
}
