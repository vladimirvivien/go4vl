package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	device2 "github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var template = "\t%-24s : %s\n"

func main() {
	var devName string
	var devList bool
	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.BoolVar(&devList, "l", false, "list all devices")
	flag.Parse()

	if devList {
		if err := listDevices(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	device, err := device2.Open(devName)
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

	if err := printCropInfo(device); err != nil {
		log.Fatal(err)
	}

	if device.Capability().IsVideoCaptureSupported() {
		if err := printCaptureParam(device); err != nil {
			log.Fatal(err)
		}
	}

	if device.Capability().IsVideoOutputSupported() {
		if err := printOutputParam(device); err != nil {
			log.Fatal(err)
		}
	}

}

func listDevices() error {
	paths, err := device2.GetAllDevicePaths()
	if err != nil {
		return err
	}
	for _, path := range paths {
		dev, err := device2.Open(path)
		if err != nil {
			log.Print(err)
			continue
		}

		var busInfo, card string
		cap := dev.Capability()

		// is a media device?
		if mdi, err := dev.GetMediaInfo(); err == nil {
			if mdi.BusInfo != "" {
				busInfo = mdi.BusInfo
			} else {
				busInfo = "platform: " + mdi.Driver
			}
			if mdi.Model != "" {
				card = mdi.Model
			} else {
				card = mdi.Driver
			}
		} else {
			busInfo = cap.BusInfo
			card = cap.Card
		}

		// close device
		if err := dev.Close(); err != nil {
			log.Print(err)
			continue
		}

		fmt.Printf("v4l2Device [%s]: %s: %s\n", path, card, busInfo)

	}
	return nil
}

func printDeviceDriverInfo(dev *device2.Device) error {
	caps := dev.Capability()

	// print driver info
	fmt.Println("v4l2Device Info:")
	fmt.Printf(template, "Driver name", caps.Driver)
	fmt.Printf(template, "Card name", caps.Card)
	fmt.Printf(template, "Bus info", caps.BusInfo)

	fmt.Printf(template, "Driver version", caps.GetVersionInfo())

	fmt.Printf("\t%-16s : %0x\n", "Driver capabilities", caps.Capabilities)
	for _, desc := range caps.GetDriverCapDescriptions() {
		fmt.Printf("\t\t%s\n", desc.Desc)
	}

	fmt.Printf("\t%-16s : %0x\n", "v4l2Device capabilities", caps.Capabilities)
	for _, desc := range caps.GetDeviceCapDescriptions() {
		fmt.Printf("\t\t%s\n", desc.Desc)
	}

	return nil
}

func printVideoInputInfo(dev *device2.Device) error {
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

func printFormatInfo(dev *device2.Device) error {
	pixFmt, err := dev.GetPixFormat()
	if err != nil {
		return fmt.Errorf("video capture format: %w", err)
	}
	fmt.Println("Video format for capture (default):")
	fmt.Printf(template, "Width x Height", fmt.Sprintf("%d x %d", pixFmt.Width, pixFmt.Height))
	fmt.Printf(template, "Pixel format", v4l2.PixelFormats[pixFmt.PixelFormat])
	fmt.Printf(template, "Field", v4l2.Fields[pixFmt.Field])
	fmt.Printf(template, "Bytes per line", fmt.Sprintf("%d", pixFmt.BytesPerLine))
	fmt.Printf(template, "Size image", fmt.Sprintf("%d", pixFmt.SizeImage))
	fmt.Printf(template, "Colorspace", v4l2.Colorspaces[pixFmt.Colorspace])

	// xferfunc
	xfunc := v4l2.XferFunctions[pixFmt.XferFunc]
	if pixFmt.XferFunc == v4l2.XferFuncDefault {
		xfunc = fmt.Sprintf("%s (map to %s)", xfunc, v4l2.XferFunctions[v4l2.ColorspaceToXferFunc(pixFmt.XferFunc)])
	}
	fmt.Printf(template, "Transfer function", xfunc)

	// ycbcr
	ycbcr := v4l2.YCbCrEncodings[pixFmt.YcbcrEnc]
	if pixFmt.YcbcrEnc == v4l2.YCbCrEncodingDefault {
		ycbcr = fmt.Sprintf("%s (map to %s)", ycbcr, v4l2.YCbCrEncodings[v4l2.ColorspaceToYCbCrEnc(pixFmt.YcbcrEnc)])
	}
	fmt.Printf(template, "YCbCr/HSV encoding", ycbcr)

	// quant
	quant := v4l2.Quantizations[pixFmt.Quantization]
	if pixFmt.Quantization == v4l2.QuantizationDefault {
		if v4l2.IsPixYUVEncoded(pixFmt.PixelFormat) {
			quant = fmt.Sprintf("%s (map to %s)", quant, v4l2.Quantizations[v4l2.QuantizationLimitedRange])
		} else {
			quant = fmt.Sprintf("%s (map to %s)", quant, v4l2.Quantizations[v4l2.QuantizationFullRange])
		}
	}
	fmt.Printf(template, "Quantization", quant)

	// format desc
	return printFormatDesc(dev)
}

func printFormatDesc(dev *device2.Device) error {
	descs, err := dev.GetFormatDescriptions()
	if err != nil {
		return fmt.Errorf("format desc: %w", err)
	}
	fmt.Println("Supported formats:")
	for i, desc := range descs {
		frmSizes, err := v4l2.GetFormatFrameSizes(dev.Fd(), desc.PixelFormat)
		if err != nil {
			return fmt.Errorf("format desc: %w", err)
		}
		var sizeStr strings.Builder
		sizeStr.WriteString("Sizes: ")
		for _, size := range frmSizes {
			sizeStr.WriteString(fmt.Sprintf("[%dx%d] ", size.Size.MinWidth, size.Size.MinHeight))
		}
		fmt.Printf(template, fmt.Sprintf("[%0d] %s", i, desc.Description), sizeStr.String())
	}
	return nil
}

func printCropInfo(dev *device2.Device) error {
	crop, err := dev.GetCropCapability()
	if err != nil {
		return fmt.Errorf("crop capability: %w", err)
	}

	fmt.Println("Crop capability for video capture:")
	fmt.Printf(
		template,
		"Bounds:",
		fmt.Sprintf(
			"left=%d; top=%d; width=%d; heigh=%d",
			crop.Bounds.Left, crop.Bounds.Top, crop.Bounds.Width, crop.Bounds.Height,
		),
	)
	fmt.Printf(
		template,
		"Default:",
		fmt.Sprintf(
			"left=%d; top=%d; width=%d; heigh=%d",
			crop.DefaultRect.Left, crop.DefaultRect.Top, crop.DefaultRect.Width, crop.DefaultRect.Height,
		),
	)
	fmt.Printf(template, "Pixel aspect", fmt.Sprintf("%d/%d", crop.PixelAspect.Numerator, crop.PixelAspect.Denominator))
	return nil
}

func printCaptureParam(dev *device2.Device) error {
	params, err := dev.GetStreamParam()
	if err != nil {
		return fmt.Errorf("stream capture param: %w", err)
	}
	fmt.Println("Stream capture parameters:")

	tpf := "not specified"
	if params.Capture.Capability == v4l2.StreamParamTimePerFrame {
		tpf = "time per frame"
	}
	fmt.Printf(template, "Capability", tpf)

	hiqual := "not specified"
	if params.Capture.CaptureMode == v4l2.StreamParamModeHighQuality {
		hiqual = "high quality"
	}
	fmt.Printf(template, "Capture mode", hiqual)

	fmt.Printf(template, "Frames per second", fmt.Sprintf("%d/%d", params.Capture.TimePerFrame.Denominator, params.Capture.TimePerFrame.Numerator))
	fmt.Printf(template, "Read buffers", fmt.Sprintf("%d", params.Capture.ReadBuffers))
	return nil
}

func printOutputParam(dev *device2.Device) error {
	params, err := dev.GetStreamParam()
	if err != nil {
		return fmt.Errorf("stream output param: %w", err)
	}
	fmt.Println("Stream output parameters:")

	tpf := "not specified"
	if params.Output.Capability == v4l2.StreamParamTimePerFrame {
		tpf = "time per frame"
	}
	fmt.Printf(template, "Capability", tpf)

	hiqual := "not specified"
	if params.Output.CaptureMode == v4l2.StreamParamModeHighQuality {
		hiqual = "high quality"
	}
	fmt.Printf(template, "Output mode", hiqual)

	fmt.Printf(template, "Frames per second", fmt.Sprintf("%d/%d", params.Output.TimePerFrame.Denominator, params.Output.TimePerFrame.Numerator))
	fmt.Printf(template, "Write buffers", fmt.Sprintf("%d", params.Output.WriteBuffers))
	return nil
}
