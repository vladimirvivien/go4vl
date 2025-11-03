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
	var selectInput int
	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.IntVar(&selectInput, "s", -1, "select input by index (optional)")
	flag.Parse()

	dev, err := device.Open(devName)
	if err != nil {
		log.Fatalf("Failed to open device %s: %v", devName, err)
	}
	defer dev.Close()

	fmt.Printf("Device: %s\n", dev.Name())
	fmt.Printf("Driver: %s\n", dev.Capability().Driver)
	fmt.Printf("Card: %s\n\n", dev.Capability().Card)

	// Check if device supports video capture
	if !dev.Capability().IsVideoCaptureSupported() {
		log.Fatal("Device does not support video capture")
	}

	// Get current input
	currentIdx, err := dev.GetVideoInputIndex()
	if err != nil {
		log.Fatalf("Failed to get current input: %v\nNote: Device may not support multiple inputs", err)
	}
	fmt.Printf("Current input index: %d\n\n", currentIdx)

	// Enumerate all inputs
	inputs, err := dev.GetVideoInputDescriptions()
	if err != nil {
		log.Fatalf("Failed to enumerate inputs: %v", err)
	}

	if len(inputs) == 0 {
		log.Fatal("No video inputs found")
	}

	fmt.Printf("Available video inputs (%d):\n", len(inputs))
	fmt.Println("================================================================================")
	for _, input := range inputs {
		active := ""
		if input.GetIndex() == uint32(currentIdx) {
			active = " ** ACTIVE **"
		}

		fmt.Printf("[%d] %s%s\n", input.GetIndex(), input.GetName(), active)

		typeName := v4l2.InputTypes[input.GetInputType()]
		if typeName == "" {
			typeName = fmt.Sprintf("Unknown (%d)", input.GetInputType())
		}
		fmt.Printf("    Type:         %s\n", typeName)

		statusName := getInputStatusName(input.GetStatus())
		fmt.Printf("    Status:       %s\n", statusName)
		fmt.Printf("    Audioset:     0x%08x\n", input.GetAudioset())
		fmt.Printf("    Tuner:        %d\n", input.GetTuner())
		fmt.Printf("    Standards:    0x%016x\n", input.GetStandardId())
		fmt.Printf("    Capabilities: 0x%08x\n", input.GetCapabilities())
		fmt.Println()
	}

	// Query current input status
	fmt.Println("Current Input Status:")
	fmt.Println("================================================================================")
	status, err := dev.GetVideoInputStatus()
	if err != nil {
		fmt.Printf("Status query not supported: %v\n", err)
	} else {
		statusStr, exists := v4l2.InputStatuses[status]
		if !exists {
			statusStr = fmt.Sprintf("Unknown (0x%x)", status)
		}
		fmt.Printf("Status: %s\n", statusStr)

		// Decode status flags
		if status == 0 {
			fmt.Println("  ✓ Power OK")
			fmt.Println("  ✓ Signal detected")
			fmt.Println("  ✓ Color information present")
		} else {
			if status&v4l2.InputStatusNoPower != 0 {
				fmt.Println("  ✗ No power")
			}
			if status&v4l2.InputStatusNoSignal != 0 {
				fmt.Println("  ✗ No signal")
			}
			if status&v4l2.InputStatusNoColor != 0 {
				fmt.Println("  ✗ No color")
			}
		}
	}
	fmt.Println()

	// Select input if requested
	if selectInput >= 0 {
		fmt.Printf("Selecting input %d...\n", selectInput)
		err = dev.SetVideoInputIndex(int32(selectInput))
		if err != nil {
			log.Fatalf("Failed to select input %d: %v", selectInput, err)
		}

		// Verify the change
		newIdx, err := dev.GetVideoInputIndex()
		if err != nil {
			log.Fatalf("Failed to verify input selection: %v", err)
		}

		if newIdx == int32(selectInput) {
			fmt.Printf("✓ Successfully switched to input %d\n", selectInput)

			// Show the new input details
			info, err := dev.GetVideoInputInfo(uint32(selectInput))
			if err != nil {
				log.Fatalf("Failed to get new input info: %v", err)
			}
			fmt.Printf("  Name: %s\n", info.GetName())

			typeName := v4l2.InputTypes[info.GetInputType()]
			if typeName == "" {
				typeName = fmt.Sprintf("Unknown (%d)", info.GetInputType())
			}
			fmt.Printf("  Type: %s\n", typeName)
		} else {
			fmt.Printf("⚠ Input selection may not have worked (got %d, expected %d)\n", newIdx, selectInput)
		}
	} else {
		fmt.Println("Tip: Use -s <index> to select a different input")
	}
}

// getInputStatusName returns a human-readable name for the input status
func getInputStatusName(status uint32) string {
	if status == 0 {
		return "OK"
	}

	statusStr, exists := v4l2.InputStatuses[v4l2.InputStatus(status)]
	if exists {
		return statusStr
	}

	// Build combined status string
	var parts []string
	if status&v4l2.InputStatusNoPower != 0 {
		parts = append(parts, "no power")
	}
	if status&v4l2.InputStatusNoSignal != 0 {
		parts = append(parts, "no signal")
	}
	if status&v4l2.InputStatusNoColor != 0 {
		parts = append(parts, "no color")
	}

	if len(parts) > 0 {
		result := parts[0]
		for i := 1; i < len(parts); i++ {
			result += ", " + parts[i]
		}
		return result
	}

	return fmt.Sprintf("Unknown (0x%x)", status)
}
