package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

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

	if err := printCropInfo(device); err != nil {
		log.Fatal(err)
	}

	if err := printCaptureParam(device); err != nil {
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
	fmt.Printf(template, "Driver name", caps.Driver)
	fmt.Printf(template, "Card name", caps.Card)
	fmt.Printf(template, "Bus info", caps.BusInfo)

	fmt.Printf(template, "Driver version", caps.GetVersionInfo())

	fmt.Printf("\t%-16s : %0x\n", "Driver capabilities", caps.Capabilities)
	for _, desc := range caps.GetDriverCapDescriptions() {
		fmt.Printf("\t\t%s\n", desc.Desc)
	}

	fmt.Printf("\t%-16s : %0x\n", "Device capabilities", caps.Capabilities)
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

func printFormatDesc(dev *v4l2.Device) error {
	descs, err := dev.GetFormatDescriptions()
	if err != nil {
		return fmt.Errorf("format desc: %w", err)
	}
	fmt.Println("Supported formats:")
	for i, desc := range descs{
		frmSizes, err := v4l2.GetFormatFrameSizes(dev.GetFileDescriptor(), desc.PixelFormat)
		if err != nil {
			return fmt.Errorf("format desc: %w", err)
		}
		var sizeStr strings.Builder
		sizeStr.WriteString("Sizes: ")
		for _, size := range frmSizes{
			sizeStr.WriteString(fmt.Sprintf("[%dx%d] ", size.Width, size.Height))
		}
		fmt.Printf(template, fmt.Sprintf("[%0d] %s", i, desc.Description), sizeStr.String())
	}
	return nil
}

func printCropInfo(dev *v4l2.Device) error {
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

func printCaptureParam(dev *v4l2.Device) error {
	params, err := dev.GetCaptureParam()
	if err != nil {
		return fmt.Errorf("streaming capture param: %w", err)
	}
	fmt.Println("Streaming parameters for video capture:")

	tpf := "not specified"
	if params.Capability == v4l2.StreamParamTimePerFrame {
		tpf = "time per frame"
	}
	fmt.Printf(template, "Capability", tpf)

	hiqual := "not specified"
	if params.CaptureMode == v4l2.StreamParamModeHighQuality {
		hiqual = "high quality"
	}
	fmt.Printf(template, "Capture mode", hiqual)

	fmt.Printf(template, "Frames per second", fmt.Sprintf("%d/%d", params.TimePerFrame.Denominator, params.TimePerFrame.Numerator))
	fmt.Printf(template, "Read buffers", fmt.Sprintf("%d",params.ReadBuffers))
	return nil
}