package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	devPath = flag.String("d", "/dev/video0", "device path")
	pad     = flag.Uint("p", 0, "pad number for multi-pad devices")
)

func main() {
	flag.Parse()

	dev, err := device.Open(*devPath)
	if err != nil {
		fmt.Printf("Failed to open device %s: %v\n", *devPath, err)
		os.Exit(1)
	}
	defer dev.Close()

	fmt.Printf("DV Timings Example - Device: %s\n", *devPath)
	fmt.Println("=====================================")
	fmt.Println()

	// Check if device supports DV timings
	if !checkDVTimingsSupport(dev) {
		fmt.Println("Device does not support DV timings.")
		fmt.Println("\nNote: DV timings are typically supported by:")
		fmt.Println("  - HDMI/DisplayPort capture cards")
		fmt.Println("  - SDI video interfaces")
		fmt.Println("  - DVI capture devices")
		fmt.Println("  - Professional video equipment")
		os.Exit(0)
	}

	// Display DV timing capabilities
	displayCapabilities(dev)

	// Display current DV timings
	displayCurrentTimings(dev)

	// Try to auto-detect timings from signal
	displayDetectedTimings(dev)

	// Enumerate all supported timings
	enumerateTimings(dev)
}

func checkDVTimingsSupport(dev *device.Device) bool {
	_, err := dev.GetDVTimings()
	return err == nil
}

func displayCapabilities(dev *device.Device) {
	fmt.Println("DV Timing Capabilities")
	fmt.Println("----------------------")

	cap, err := dev.GetDVTimingsCap(uint32(*pad))
	if err != nil {
		fmt.Printf("  Could not query capabilities: %v\n", err)
		fmt.Println()
		return
	}

	btCap := cap.GetBTTimingsCap()

	fmt.Printf("  Type: %d\n", cap.GetType())
	fmt.Printf("  Pad: %d\n", cap.GetPad())
	fmt.Printf("  Resolution range: %dx%d to %dx%d\n",
		btCap.GetMinWidth(), btCap.GetMinHeight(),
		btCap.GetMaxWidth(), btCap.GetMaxHeight())
	fmt.Printf("  Pixel clock range: %d - %d Hz (%.1f - %.1f MHz)\n",
		btCap.GetMinPixelClock(), btCap.GetMaxPixelClock(),
		float64(btCap.GetMinPixelClock())/1e6,
		float64(btCap.GetMaxPixelClock())/1e6)

	fmt.Printf("\n  Supported formats:\n")
	if btCap.SupportsInterlaced() {
		fmt.Println("    - Interlaced")
	}
	if btCap.SupportsProgressive() {
		fmt.Println("    - Progressive")
	}
	if btCap.SupportsReducedBlanking() {
		fmt.Println("    - Reduced blanking")
	}
	if btCap.SupportsCustomTimings() {
		fmt.Println("    - Custom timings")
	}

	fmt.Printf("\n  Supported standards:\n")
	if btCap.HasStandard(v4l2.DVStdCEA861) {
		fmt.Println("    - CEA-861 (HDMI/DVI)")
	}
	if btCap.HasStandard(v4l2.DVStdDMT) {
		fmt.Println("    - DMT (VESA Display Monitor Timings)")
	}
	if btCap.HasStandard(v4l2.DVStdCVT) {
		fmt.Println("    - CVT (Coordinated Video Timings)")
	}
	if btCap.HasStandard(v4l2.DVStdGTF) {
		fmt.Println("    - GTF (Generalized Timing Formula)")
	}

	fmt.Println()
}

func displayCurrentTimings(dev *device.Device) {
	fmt.Println("Current DV Timings")
	fmt.Println("------------------")

	timings, err := dev.GetDVTimings()
	if err != nil {
		fmt.Printf("  Could not get current timings: %v\n", err)
		fmt.Println()
		return
	}

	printTimingDetails(timings)
}

func displayDetectedTimings(dev *device.Device) {
	fmt.Println("Auto-Detected Timings")
	fmt.Println("---------------------")

	timings, err := dev.QueryDVTimings()
	if err != nil {
		fmt.Printf("  Could not detect timings: %v\n", err)
		fmt.Println("  (No signal present or auto-detection not supported)\n")
		return
	}

	fmt.Println("  Successfully detected timings from input signal:")
	printTimingDetails(timings)
}

func enumerateTimings(dev *device.Device) {
	fmt.Println("Enumerated Supported Timings")
	fmt.Println("----------------------------")

	timings, err := dev.GetAllDVTimings(uint32(*pad))
	if err != nil {
		fmt.Printf("  Could not enumerate timings: %v\n", err)
		fmt.Println()
		return
	}

	if len(timings) == 0 {
		fmt.Println("  No timings enumerated")
		fmt.Println()
		return
	}

	fmt.Printf("  Found %d supported timing(s):\n\n", len(timings))

	for i, enumTiming := range timings {
		dv := enumTiming.GetTimings()
		bt := dv.GetBTTimings()

		fmt.Printf("  [%d] %dx%d @ %.2f Hz", i,
			bt.GetWidth(), bt.GetHeight(), bt.GetFrameRate())

		if bt.IsInterlaced() {
			fmt.Print(" (interlaced)")
		}

		// Show standards
		standards := []string{}
		if bt.HasStandard(v4l2.DVStdCEA861) {
			standards = append(standards, "CEA-861")
		}
		if bt.HasStandard(v4l2.DVStdDMT) {
			standards = append(standards, "DMT")
		}
		if bt.HasStandard(v4l2.DVStdCVT) {
			standards = append(standards, "CVT")
		}
		if bt.HasStandard(v4l2.DVStdGTF) {
			standards = append(standards, "GTF")
		}

		if len(standards) > 0 {
			fmt.Printf(" - %v", standards)
		}

		fmt.Println()

		// Only show first 10 to avoid too much output
		if i >= 9 && len(timings) > 10 {
			fmt.Printf("\n  ... and %d more\n", len(timings)-10)
			break
		}
	}

	fmt.Println()
}

func printTimingDetails(timings v4l2.DVTimings) {
	bt := timings.GetBTTimings()

	fmt.Printf("  Type: %d (BT.656/1120)\n", timings.GetType())
	fmt.Printf("  Resolution: %dx%d\n", bt.GetWidth(), bt.GetHeight())
	fmt.Printf("  Pixel Clock: %d Hz (%.2f MHz)\n",
		bt.GetPixelClock(), float64(bt.GetPixelClock())/1e6)
	fmt.Printf("  Frame Rate: %.2f Hz\n", bt.GetFrameRate())

	if bt.IsInterlaced() {
		fmt.Println("  Format: Interlaced")
	} else {
		fmt.Println("  Format: Progressive")
	}

	fmt.Printf("\n  Horizontal Blanking:\n")
	fmt.Printf("    Front Porch: %d pixels\n", bt.GetHFrontPorch())
	fmt.Printf("    Sync: %d pixels\n", bt.GetHSync())
	fmt.Printf("    Back Porch: %d pixels\n", bt.GetHBackPorch())
	hTotal := bt.GetWidth() + bt.GetHFrontPorch() + bt.GetHSync() + bt.GetHBackPorch()
	fmt.Printf("    Total: %d pixels\n", hTotal)

	fmt.Printf("\n  Vertical Blanking:\n")
	fmt.Printf("    Front Porch: %d lines\n", bt.GetVFrontPorch())
	fmt.Printf("    Sync: %d lines\n", bt.GetVSync())
	fmt.Printf("    Back Porch: %d lines\n", bt.GetVBackPorch())
	vTotal := bt.GetHeight() + bt.GetVFrontPorch() + bt.GetVSync() + bt.GetVBackPorch()

	if bt.IsInterlaced() {
		fmt.Printf("    IL Front Porch: %d lines\n", bt.GetILVFrontPorch())
		fmt.Printf("    IL Sync: %d lines\n", bt.GetILVSync())
		fmt.Printf("    IL Back Porch: %d lines\n", bt.GetILVBackPorch())
		vTotal += bt.GetILVFrontPorch() + bt.GetILVSync() + bt.GetILVBackPorch()
	}
	fmt.Printf("    Total: %d lines\n", vTotal)

	fmt.Printf("\n  Sync Polarities:\n")
	if bt.HasHSyncPosPolarity() {
		fmt.Println("    H-Sync: Positive")
	} else {
		fmt.Println("    H-Sync: Negative")
	}
	if bt.HasVSyncPosPolarity() {
		fmt.Println("    V-Sync: Positive")
	} else {
		fmt.Println("    V-Sync: Negative")
	}

	fmt.Printf("\n  Standards: 0x%x\n", bt.GetStandards())
	if bt.HasStandard(v4l2.DVStdCEA861) {
		fmt.Println("    - CEA-861 (HDMI/DVI)")
	}
	if bt.HasStandard(v4l2.DVStdDMT) {
		fmt.Println("    - DMT (VESA)")
	}
	if bt.HasStandard(v4l2.DVStdCVT) {
		fmt.Println("    - CVT")
	}
	if bt.HasStandard(v4l2.DVStdGTF) {
		fmt.Println("    - GTF")
	}

	fmt.Printf("\n  Flags: 0x%x\n", bt.GetFlags())
	if bt.HasFlag(v4l2.DVFlagReducedBlanking) {
		fmt.Println("    - Reduced blanking")
	}
	if bt.HasFlag(v4l2.DVFlagReducedFPS) {
		fmt.Println("    - Reduced FPS")
	}
	if bt.HasFlag(v4l2.DVFlagIsCEVideo) {
		fmt.Println("    - CE Video")
	}
	if bt.HasFlag(v4l2.DVFlagHasPictureAspect) {
		fmt.Println("    - Has picture aspect ratio")
	}

	if bt.HasFlag(v4l2.DVFlagHasCEA861VIC) {
		fmt.Printf("    - CEA-861 VIC: %d\n", bt.GetCEA861VIC())
	}
	if bt.HasFlag(v4l2.DVFlagHasHDMIVIC) {
		fmt.Printf("    - HDMI VIC: %d\n", bt.GetHDMIVIC())
	}

	fmt.Println()
}
