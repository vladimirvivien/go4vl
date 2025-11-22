package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	devicePath = flag.String("d", "/dev/video0", "device path")
	showAll    = flag.Bool("all", false, "show all colorspace information")
	testCS     = flag.String("test", "", "test colorspace (e.g., srgb, rec709, bt2020, dci-p3)")
)

func main() {
	flag.Parse()

	if *showAll {
		showAllColorspaces()
		return
	}

	if *testCS != "" {
		testColorspace(*testCS)
		return
	}

	// Default: show current device colorspace
	showDeviceColorspace(*devicePath)
}

func showDeviceColorspace(devPath string) {
	dev, err := device.Open(devPath)
	if err != nil {
		fmt.Printf("Failed to open device: %v\n", err)
		os.Exit(1)
	}
	defer dev.Close()

	fmt.Printf("Device Colorspace - %s\n", devPath)
	fmt.Println("================================================================================")
	fmt.Println()

	pixFmt, err := v4l2.GetPixFormat(dev.Fd())
	if err != nil {
		fmt.Printf("Failed to get pixel format: %v\n", err)
		os.Exit(1)
	}

	// Create colorspace info
	csInfo := v4l2.ColorspaceInfo{
		Colorspace:   pixFmt.Colorspace,
		YCbCrEnc:     pixFmt.YcbcrEnc,
		Quantization: pixFmt.Quantization,
		XferFunc:     pixFmt.XferFunc,
	}

	fmt.Printf("Current Format:  %s\n", v4l2.PixelFormats[pixFmt.PixelFormat])
	fmt.Printf("Resolution:      %d x %d\n", pixFmt.Width, pixFmt.Height)
	fmt.Println()

	fmt.Println("Colorspace Information:")
	fmt.Println("----------------------")
	fmt.Printf("  Colorspace:       %s\n", v4l2.Colorspaces[csInfo.Colorspace])
	fmt.Printf("  YCbCr Encoding:   %s\n", v4l2.YCbCrEncodings[csInfo.YCbCrEnc])
	fmt.Printf("  Quantization:     %s\n", v4l2.Quantizations[csInfo.Quantization])
	fmt.Printf("  Transfer Func:    %s\n", v4l2.XferFunctions[csInfo.XferFunc])
	fmt.Println()

	fmt.Printf("Complete String:  %s\n", csInfo.String())
	fmt.Printf("HDR Content:      %v\n", csInfo.IsHDR())
	fmt.Printf("SDR Content:      %v\n", csInfo.IsSDR())
}

func showAllColorspaces() {
	fmt.Println("V4L2 Colorspace Reference")
	fmt.Println("================================================================================")
	fmt.Println()

	// Colorspaces table
	fmt.Println("1. Colorspaces")
	fmt.Println("--------------")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Constant\tName\tDescription")
	fmt.Fprintln(w, "--------\t----\t-----------")
	fmt.Fprintf(w, "ColorspaceDefault\t%s\tDriver default\n", v4l2.Colorspaces[v4l2.ColorspaceDefault])
	fmt.Fprintf(w, "ColorspaceSMPTE170M\t%s\tNTSC/PAL/SECAM (SDTV)\n", v4l2.Colorspaces[v4l2.ColorspaceSMPTE170M])
	fmt.Fprintf(w, "ColorspaceREC709\t%s\tHDTV (1080p/720p)\n", v4l2.Colorspaces[v4l2.ColorspaceREC709])
	fmt.Fprintf(w, "ColorspaceBT2020\t%s\tUHDTV (4K/8K)\n", v4l2.Colorspaces[v4l2.ColorspaceBT2020])
	fmt.Fprintf(w, "ColorspaceSRGB\t%s\tComputer graphics\n", v4l2.Colorspaces[v4l2.ColorspaceSRGB])
	fmt.Fprintf(w, "ColorspaceDCIP3\t%s\tDigital cinema\n", v4l2.Colorspaces[v4l2.ColorspaceDCIP3])
	fmt.Fprintf(w, "ColorspaceJPEG\t%s\tJPEG/sYCC\n", v4l2.Colorspaces[v4l2.ColorspaceJPEG])
	fmt.Fprintf(w, "ColorspaceOPRGB\t%s\topRGB\n", v4l2.Colorspaces[v4l2.ColorspaceOPRGB])
	fmt.Fprintf(w, "ColorspaceRaw\t%s\tNo conversion\n", v4l2.Colorspaces[v4l2.ColorspaceRaw])
	w.Flush()
	fmt.Println()

	// YCbCr Encodings table
	fmt.Println("2. YCbCr Encodings")
	fmt.Println("------------------")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Constant\tName\tDescription")
	fmt.Fprintln(w, "--------\t----\t-----------")
	fmt.Fprintf(w, "YCbCrEncoding601\t%s\tSDTV (NTSC/PAL)\n", v4l2.YCbCrEncodings[v4l2.YCbCrEncoding601])
	fmt.Fprintf(w, "YCbCrEncoding709\t%s\tHDTV\n", v4l2.YCbCrEncodings[v4l2.YCbCrEncoding709])
	fmt.Fprintf(w, "YCbCrEncodingBT2020\t%s\tUHDTV\n", v4l2.YCbCrEncodings[v4l2.YCbCrEncodingBT2020])
	fmt.Fprintf(w, "YCbCrEncodingXV601\t%s\tExtended gamut SDTV\n", v4l2.YCbCrEncodings[v4l2.YCbCrEncodingXV601])
	fmt.Fprintf(w, "YCbCrEncodingXV709\t%s\tExtended gamut HDTV\n", v4l2.YCbCrEncodings[v4l2.YCbCrEncodingXV709])
	w.Flush()
	fmt.Println()

	// Quantization table
	fmt.Println("3. Quantization Ranges")
	fmt.Println("----------------------")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Constant\tName\tRange (8-bit)")
	fmt.Fprintln(w, "--------\t----\t-------------")
	fmt.Fprintf(w, "QuantizationFullRange\t%s\t0-255\n", v4l2.Quantizations[v4l2.QuantizationFullRange])
	fmt.Fprintf(w, "QuantizationLimitedRange\t%s\tY: 16-235, Cb/Cr: 16-240\n", v4l2.Quantizations[v4l2.QuantizationLimitedRange])
	w.Flush()
	fmt.Println()

	// Transfer Functions table
	fmt.Println("4. Transfer Functions (EOTF/Gamma)")
	fmt.Println("----------------------------------")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Constant\tName\tType")
	fmt.Fprintln(w, "--------\t----\t----")
	fmt.Fprintf(w, "XferFunc709\t%s\tSDR (Rec.709)\n", v4l2.XferFunctions[v4l2.XferFunc709])
	fmt.Fprintf(w, "XferFuncSRGB\t%s\tSDR (sRGB)\n", v4l2.XferFunctions[v4l2.XferFuncSRGB])
	fmt.Fprintf(w, "XferFuncSMPTE2084\t%s\tHDR (PQ)\n", v4l2.XferFunctions[v4l2.XferFuncSMPTE2084])
	fmt.Fprintf(w, "XferFuncDCIP3\t%s\tDigital cinema\n", v4l2.XferFunctions[v4l2.XferFuncDCIP3])
	fmt.Fprintf(w, "XferFuncNone\t%s\tLinear (raw)\n", v4l2.XferFunctions[v4l2.XferFuncNone])
	w.Flush()
	fmt.Println()

	// Common combinations
	fmt.Println("5. Common Colorspace Combinations")
	fmt.Println("----------------------------------")
	showCommonCombination("SDTV (NTSC/PAL)", v4l2.ColorspaceSMPTE170M)
	showCommonCombination("HDTV (720p/1080p)", v4l2.ColorspaceREC709)
	showCommonCombination("UHDTV (4K/8K)", v4l2.ColorspaceBT2020)
	showCommonCombination("sRGB (Computer)", v4l2.ColorspaceSRGB)
	showCommonCombination("DCI-P3 (Cinema)", v4l2.ColorspaceDCIP3)
}

func showCommonCombination(name string, cs v4l2.ColorspaceType) {
	csInfo := v4l2.NewColorspaceInfo(cs)
	fmt.Printf("  %s:\n", name)
	fmt.Printf("    %s\n", csInfo.String())
	fmt.Printf("    HDR: %v\n", csInfo.IsHDR())
	fmt.Println()
}

func testColorspace(name string) {
	var cs v4l2.ColorspaceType
	var csName string

	switch name {
	case "srgb":
		cs = v4l2.ColorspaceSRGB
		csName = "sRGB"
	case "rec709", "709":
		cs = v4l2.ColorspaceREC709
		csName = "Rec. 709"
	case "bt2020", "2020":
		cs = v4l2.ColorspaceBT2020
		csName = "BT.2020"
	case "dci-p3", "dcip3", "p3":
		cs = v4l2.ColorspaceDCIP3
		csName = "DCI-P3"
	case "jpeg":
		cs = v4l2.ColorspaceJPEG
		csName = "JPEG"
	case "raw":
		cs = v4l2.ColorspaceRaw
		csName = "Raw"
	default:
		fmt.Printf("Unknown colorspace: %s\n", name)
		fmt.Println("Supported: srgb, rec709, bt2020, dci-p3, jpeg, raw")
		os.Exit(1)
	}

	fmt.Printf("Testing Colorspace: %s\n", csName)
	fmt.Println("================================================================================")
	fmt.Println()

	csInfo := v4l2.NewColorspaceInfo(cs)

	fmt.Println("Colorspace Details:")
	fmt.Println("------------------")
	fmt.Printf("  Colorspace:       %s\n", v4l2.Colorspaces[csInfo.Colorspace])
	fmt.Printf("  YCbCr Encoding:   %s\n", v4l2.YCbCrEncodings[csInfo.YCbCrEnc])
	fmt.Printf("  Quantization:     %s\n", v4l2.Quantizations[csInfo.Quantization])
	fmt.Printf("  Transfer Func:    %s\n", v4l2.XferFunctions[csInfo.XferFunc])
	fmt.Println()

	fmt.Printf("Complete String:  %s\n", csInfo.String())
	fmt.Printf("HDR Content:      %v\n", csInfo.IsHDR())
	fmt.Printf("SDR Content:      %v\n", csInfo.IsSDR())
	fmt.Println()

	// Use case information
	fmt.Println("Typical Use Cases:")
	fmt.Println("-----------------")
	switch cs {
	case v4l2.ColorspaceSRGB:
		fmt.Println("  • Computer graphics and displays")
		fmt.Println("  • Web content")
		fmt.Println("  • Digital photography (JPEG)")
	case v4l2.ColorspaceREC709:
		fmt.Println("  • HDTV broadcasting (720p, 1080p)")
		fmt.Println("  • Blu-ray discs")
		fmt.Println("  • Streaming video (HD)")
	case v4l2.ColorspaceBT2020:
		fmt.Println("  • UHDTV (4K, 8K)")
		fmt.Println("  • HDR content")
		fmt.Println("  • Wide color gamut displays")
	case v4l2.ColorspaceDCIP3:
		fmt.Println("  • Digital cinema projection")
		fmt.Println("  • Professional video production")
		fmt.Println("  • HDR mastering")
	case v4l2.ColorspaceJPEG:
		fmt.Println("  • JPEG images")
		fmt.Println("  • sYCC color space")
	case v4l2.ColorspaceRaw:
		fmt.Println("  • Raw sensor data")
		fmt.Println("  • No colorspace conversion")
	}
}
