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
	var selectOutput int
	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.IntVar(&selectOutput, "s", -1, "select output by index (optional)")
	flag.Parse()

	dev, err := device.Open(devName)
	if err != nil {
		log.Fatalf("Failed to open device %s: %v", devName, err)
	}
	defer dev.Close()

	fmt.Printf("Device: %s\n", dev.Name())
	fmt.Printf("Driver: %s\n", dev.Capability().Driver)
	fmt.Printf("Card: %s\n\n", dev.Capability().Card)

	// Check if device supports video output
	if !dev.Capability().IsVideoOutputSupported() {
		log.Fatal("Device does not support video output\n" +
			"Note: Output devices are rare. You may need:\n" +
			"  - v4l2loopback configured for output mode\n" +
			"  - Hardware video encoder/output device\n" +
			"  - Memory-to-memory codec device")
	}

	// Get current output
	currentIdx, err := dev.GetVideoOutputIndex()
	if err != nil {
		log.Fatalf("Failed to get current output: %v\nNote: Device may not support multiple outputs", err)
	}
	fmt.Printf("Current output index: %d\n\n", currentIdx)

	// Enumerate all outputs
	outputs, err := dev.GetVideoOutputDescriptions()
	if err != nil {
		log.Fatalf("Failed to enumerate outputs: %v", err)
	}

	if len(outputs) == 0 {
		log.Fatal("No video outputs found")
	}

	fmt.Printf("Available video outputs (%d):\n", len(outputs))
	fmt.Println("================================================================================")
	for _, output := range outputs {
		active := ""
		if output.GetIndex() == uint32(currentIdx) {
			active = " ** ACTIVE **"
		}

		fmt.Printf("[%d] %s%s\n", output.GetIndex(), output.GetName(), active)
		fmt.Printf("    Type:         %s\n", getOutputTypeName(output.GetOutputType()))
		fmt.Printf("    Audioset:     0x%08x\n", output.GetAudioset())
		fmt.Printf("    Modulator:    %d\n", output.GetModulator())
		fmt.Printf("    Standards:    0x%016x\n", output.GetStandardId())
		fmt.Printf("    Capabilities: 0x%08x\n", output.GetCapabilities())
		fmt.Println()
	}

	// Query current output status
	fmt.Println("Current Output Status:")
	fmt.Println("================================================================================")
	status, err := dev.GetVideoOutputStatus()
	if err != nil {
		fmt.Printf("Status query not supported: %v\n", err)
	} else {
		statusStr, exists := v4l2.OutputStatuses[status]
		if !exists {
			statusStr = fmt.Sprintf("Unknown (0x%x)", status)
		}
		fmt.Printf("Status: %s\n", statusStr)

		// Note: V4L2 doesn't define output status flags like it does for inputs
		if status == 0 {
			fmt.Println("  ✓ Output OK")
		} else {
			fmt.Printf("  ⚠ Unexpected status value: 0x%x\n", status)
		}
	}
	fmt.Println()

	// Select output if requested
	if selectOutput >= 0 {
		fmt.Printf("Selecting output %d...\n", selectOutput)
		err = dev.SetVideoOutputIndex(int32(selectOutput))
		if err != nil {
			log.Fatalf("Failed to select output %d: %v", selectOutput, err)
		}

		// Verify the change
		newIdx, err := dev.GetVideoOutputIndex()
		if err != nil {
			log.Fatalf("Failed to verify output selection: %v", err)
		}

		if newIdx == int32(selectOutput) {
			fmt.Printf("✓ Successfully switched to output %d\n", selectOutput)

			// Show the new output details
			info, err := dev.GetVideoOutputInfo(uint32(selectOutput))
			if err != nil {
				log.Fatalf("Failed to get new output info: %v", err)
			}
			fmt.Printf("  Name: %s\n", info.GetName())
			fmt.Printf("  Type: %s\n", getOutputTypeName(info.GetOutputType()))
		} else {
			fmt.Printf("⚠ Output selection may not have worked (got %d, expected %d)\n", newIdx, selectOutput)
		}
	} else {
		fmt.Println("Tip: Use -s <index> to select a different output")
	}
}

// getOutputTypeName returns a human-readable name for the output type
func getOutputTypeName(outputType v4l2.OutputType) string {
	switch outputType {
	case v4l2.OutputTypeModulator:
		return "Modulator"
	case v4l2.OutputTypeAnalog:
		return "Analog"
	case v4l2.OutputTypeAnalogVGAOverlay:
		return "Analog VGA Overlay"
	default:
		return fmt.Sprintf("Unknown (%d)", outputType)
	}
}
