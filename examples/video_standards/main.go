package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	devPath   = flag.String("d", "/dev/video0", "device path")
	setStd    = flag.String("set", "", "set video standard (PAL, NTSC, SECAM, PAL_BG, NTSC_M, etc.)")
	queryMode = flag.Bool("query", false, "query/auto-detect current standard from signal")
)

func main() {
	flag.Parse()

	dev, err := device.Open(*devPath)
	if err != nil {
		fmt.Printf("Failed to open device %s: %v\n", *devPath, err)
		os.Exit(1)
	}
	defer dev.Close()

	fmt.Printf("Video Standards Example - Device: %s\n", *devPath)
	fmt.Println("==========================================")
	fmt.Println()

	// Check if device supports video standards
	if !checkStandardsSupport(dev) {
		fmt.Println("Device does not support analog video standards.")
		fmt.Println("\nNote: Video standards are typically supported by:")
		fmt.Println("  - TV tuner cards")
		fmt.Println("  - Composite/S-Video capture cards")
		fmt.Println("  - Analog cameras")
		fmt.Println("  - Legacy video equipment")
		fmt.Println("\nModern digital devices (HDMI, USB webcams, etc.) use DV timings instead.")
		os.Exit(0)
	}

	// Handle query mode
	if *queryMode {
		queryStandard(dev)
		return
	}

	// Handle set mode
	if *setStd != "" {
		setStandard(dev, *setStd)
		return
	}

	// Default: display information
	displayCurrentStandard(dev)
	enumerateStandards(dev)
	checkCommonStandards(dev)
}

func checkStandardsSupport(dev *device.Device) bool {
	standards, err := dev.GetAllStandards()
	if err != nil {
		return false
	}
	return len(standards) > 0
}

func displayCurrentStandard(dev *device.Device) {
	fmt.Println("Current Video Standard")
	fmt.Println("----------------------")

	stdId, err := dev.GetStandard()
	if err != nil {
		fmt.Printf("  Could not get current standard: %v\n", err)
		fmt.Println()
		return
	}

	fmt.Printf("  Standard ID: 0x%016x\n", stdId)

	// Try to find a name for this standard
	if name, ok := v4l2.StdNames[stdId]; ok {
		fmt.Printf("  Standard Name: %s\n", name)
	} else {
		fmt.Printf("  Standard Name: (unknown/custom)\n")
	}

	// Check which family it belongs to
	if (stdId & v4l2.StdPAL) != 0 {
		fmt.Println("  Family: PAL")
	} else if (stdId & v4l2.StdNTSC) != 0 {
		fmt.Println("  Family: NTSC")
	} else if (stdId & v4l2.StdSECAM) != 0 {
		fmt.Println("  Family: SECAM")
	}

	// Check line count and refresh rate
	if (stdId & v4l2.Std525_60) != 0 {
		fmt.Println("  Format: 525 lines / 60 Hz")
	} else if (stdId & v4l2.Std625_50) != 0 {
		fmt.Println("  Format: 625 lines / 50 Hz")
	}

	fmt.Println()
}

func enumerateStandards(dev *device.Device) {
	fmt.Println("Supported Video Standards")
	fmt.Println("-------------------------")

	standards, err := dev.GetAllStandards()
	if err != nil {
		fmt.Printf("  Could not enumerate standards: %v\n", err)
		fmt.Println()
		return
	}

	if len(standards) == 0 {
		fmt.Println("  No standards reported")
		fmt.Println()
		return
	}

	for i, std := range standards {
		fmt.Printf("  [%d] %s\n", i, std.Name())
		fmt.Printf("      ID: 0x%016x\n", std.ID())
		fmt.Printf("      Frame rate: %.2f fps\n", std.FrameRate())
		fmt.Printf("      Frame lines: %d\n", std.FrameLines())

		framePeriod := std.FramePeriod()
		fmt.Printf("      Frame period: %d/%d seconds\n",
			framePeriod.Numerator, framePeriod.Denominator)

		// Show which family this belongs to
		if (std.ID() & v4l2.StdPAL) != 0 {
			fmt.Println("      Family: PAL")
		} else if (std.ID() & v4l2.StdNTSC) != 0 {
			fmt.Println("      Family: NTSC")
		} else if (std.ID() & v4l2.StdSECAM) != 0 {
			fmt.Println("      Family: SECAM")
		}

		fmt.Println()
	}
}

func checkCommonStandards(dev *device.Device) {
	fmt.Println("Common Standard Support")
	fmt.Println("-----------------------")

	commonStandards := []struct {
		name  string
		stdId v4l2.StdId
		desc  string
	}{
		{"PAL", v4l2.StdPAL, "All PAL variants"},
		{"PAL-B/G", v4l2.StdPAL_BG, "PAL B/G (Western Europe)"},
		{"PAL-D/K", v4l2.StdPAL_DK, "PAL D/K (Eastern Europe, China)"},
		{"PAL-I", v4l2.StdPAL_I, "PAL I (UK, Ireland)"},
		{"PAL-M", v4l2.StdPAL_M, "PAL M (Brazil)"},
		{"PAL-N", v4l2.StdPAL_N, "PAL N (Argentina, Paraguay, Uruguay)"},
		{"NTSC", v4l2.StdNTSC, "All NTSC variants"},
		{"NTSC-M", v4l2.StdNTSC_M, "NTSC M (USA, Canada, BTSC)"},
		{"NTSC-M-JP", v4l2.StdNTSC_M_JP, "NTSC M Japan (EIA-J)"},
		{"NTSC-M-KR", v4l2.StdNTSC_M_KR, "NTSC M Korea (FM A2)"},
		{"SECAM", v4l2.StdSECAM, "All SECAM variants"},
		{"SECAM-B", v4l2.StdSECAM_B, "SECAM B"},
		{"SECAM-D/K", v4l2.StdSECAM_DK, "SECAM D/K"},
		{"SECAM-L", v4l2.StdSECAM_L, "SECAM L (France)"},
		{"525/60", v4l2.Std525_60, "525 lines, 60 Hz (NTSC)"},
		{"625/50", v4l2.Std625_50, "625 lines, 50 Hz (PAL/SECAM)"},
	}

	for _, cs := range commonStandards {
		supported, err := dev.IsStandardSupported(cs.stdId)
		if err != nil {
			continue
		}

		status := "Not supported"
		if supported {
			status = "Supported"
		}

		fmt.Printf("  %-12s [%s] - %s\n", cs.name, status, cs.desc)
	}

	fmt.Println()
}

func queryStandard(dev *device.Device) {
	fmt.Println("Querying/Auto-detecting Video Standard")
	fmt.Println("---------------------------------------")
	fmt.Println("  Attempting to detect standard from input signal...")
	fmt.Println()

	detectedStd, err := dev.QueryStandard()
	if err != nil {
		fmt.Printf("  Failed to detect standard: %v\n", err)
		fmt.Println("\n  Possible reasons:")
		fmt.Println("    - No signal present on input")
		fmt.Println("    - Device doesn't support standard detection")
		fmt.Println("    - Signal is too weak or unstable")
		os.Exit(1)
	}

	fmt.Printf("  Detected Standard ID: 0x%016x\n", detectedStd)

	if name, ok := v4l2.StdNames[detectedStd]; ok {
		fmt.Printf("  Detected Standard Name: %s\n", name)
	} else {
		fmt.Printf("  Detected Standard Name: (unknown/custom)\n")
	}

	// Check which family
	if (detectedStd & v4l2.StdPAL) != 0 {
		fmt.Println("  Family: PAL")
	} else if (detectedStd & v4l2.StdNTSC) != 0 {
		fmt.Println("  Family: NTSC")
	} else if (detectedStd & v4l2.StdSECAM) != 0 {
		fmt.Println("  Family: SECAM")
	}

	fmt.Println("\n  To apply this standard, run:")
	fmt.Printf("    %s -d %s -set <standard_name>\n", os.Args[0], *devPath)
}

func setStandard(dev *device.Device, stdName string) {
	fmt.Println("Setting Video Standard")
	fmt.Println("----------------------")
	fmt.Printf("  Requested standard: %s\n", stdName)

	// Map standard names to IDs
	standardMap := map[string]v4l2.StdId{
		"PAL":       v4l2.StdPAL,
		"PAL_BG":    v4l2.StdPAL_BG,
		"PAL_B":     v4l2.StdPAL_B,
		"PAL_G":     v4l2.StdPAL_G,
		"PAL_H":     v4l2.StdPAL_H,
		"PAL_I":     v4l2.StdPAL_I,
		"PAL_D":     v4l2.StdPAL_D,
		"PAL_K":     v4l2.StdPAL_K,
		"PAL_M":     v4l2.StdPAL_M,
		"PAL_N":     v4l2.StdPAL_N,
		"PAL_DK":    v4l2.StdPAL_DK,
		"NTSC":      v4l2.StdNTSC,
		"NTSC_M":    v4l2.StdNTSC_M,
		"NTSC_M_JP": v4l2.StdNTSC_M_JP,
		"NTSC_M_KR": v4l2.StdNTSC_M_KR,
		"SECAM":     v4l2.StdSECAM,
		"SECAM_B":   v4l2.StdSECAM_B,
		"SECAM_D":   v4l2.StdSECAM_D,
		"SECAM_G":   v4l2.StdSECAM_G,
		"SECAM_H":   v4l2.StdSECAM_H,
		"SECAM_K":   v4l2.StdSECAM_K,
		"SECAM_L":   v4l2.StdSECAM_L,
		"SECAM_DK":  v4l2.StdSECAM_DK,
		"525_60":    v4l2.Std525_60,
		"625_50":    v4l2.Std625_50,
	}

	stdId, ok := standardMap[stdName]
	if !ok {
		fmt.Printf("  Error: Unknown standard '%s'\n", stdName)
		fmt.Println("\n  Available standards:")
		for name := range standardMap {
			fmt.Printf("    - %s\n", name)
		}
		os.Exit(1)
	}

	// Check if standard is supported
	supported, err := dev.IsStandardSupported(stdId)
	if err != nil {
		fmt.Printf("  Error checking support: %v\n", err)
		os.Exit(1)
	}

	if !supported {
		fmt.Printf("  Warning: Standard %s may not be supported by device\n", stdName)
	}

	// Set the standard
	err = dev.SetStandard(stdId)
	if err != nil {
		fmt.Printf("  Error setting standard: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("  Standard set successfully!")

	// Verify it was set
	currentStd, err := dev.GetStandard()
	if err == nil {
		fmt.Printf("  Current standard: 0x%016x\n", currentStd)
		if (currentStd & stdId) != 0 {
			fmt.Println("  Verification: OK (standard matches)")
		} else {
			fmt.Println("  Verification: WARNING (standard may have been adjusted by driver)")
		}
	}

	fmt.Println("\n  Note: Changing the video standard may also change the current pixel format.")
}
