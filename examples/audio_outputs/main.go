package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vladimirvivien/go4vl/device"
)

func main() {
	var devName string
	var selectAudioOut int
	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.IntVar(&selectAudioOut, "s", -1, "select audio output by index (optional)")
	flag.Parse()

	dev, err := device.Open(devName)
	if err != nil {
		log.Fatalf("Failed to open device %s: %v", devName, err)
	}
	defer dev.Close()

	fmt.Printf("Device: %s\n", dev.Name())
	fmt.Printf("Driver: %s\n", dev.Capability().Driver)
	fmt.Printf("Card: %s\n\n", dev.Capability().Card)

	// Get current audio output
	currentAudioOut, err := dev.GetCurrentAudioOut()
	if err != nil {
		fmt.Printf("Device does not support audio outputs: %v\n", err)
		fmt.Println("\nNote: Audio outputs are typically found on:")
		fmt.Println("  - Video output devices")
		fmt.Println("  - TV tuner cards")
		fmt.Println("  - Professional video equipment")
		return
	}

	fmt.Printf("Current audio output: [%d] %s\n\n", currentAudioOut.GetIndex(), currentAudioOut.GetName())

	// Enumerate all audio outputs
	audioOuts, err := dev.GetAudioOutDescriptions()
	if err != nil {
		log.Fatalf("Failed to enumerate audio outputs: %v", err)
	}

	if len(audioOuts) == 0 {
		fmt.Println("No audio outputs found")
		return
	}

	fmt.Printf("Available audio outputs (%d):\n", len(audioOuts))
	fmt.Println("================================================================================")
	for _, audioOut := range audioOuts {
		active := ""
		if audioOut.GetIndex() == currentAudioOut.GetIndex() {
			active = " ** ACTIVE **"
		}

		fmt.Printf("[%d] %s%s\n", audioOut.GetIndex(), audioOut.GetName(), active)
		fmt.Printf("    Capability: 0x%08x\n", audioOut.GetCapability())
		fmt.Printf("    Mode:       0x%08x\n", audioOut.GetMode())

		// Display capability flags
		if audioOut.IsStereo() {
			fmt.Println("    ✓ Stereo")
		} else {
			fmt.Println("    • Mono")
		}

		if audioOut.HasAVL() {
			fmt.Println("    ✓ Automatic Volume Level (AVL)")
		}

		fmt.Println()
	}

	// Show current audio output details
	fmt.Println("Current Audio Output Details:")
	fmt.Println("================================================================================")
	fmt.Printf("Index:      %d\n", currentAudioOut.GetIndex())
	fmt.Printf("Name:       %s\n", currentAudioOut.GetName())
	fmt.Printf("Capability: 0x%08x\n", currentAudioOut.GetCapability())
	fmt.Printf("Mode:       0x%08x\n", currentAudioOut.GetMode())
	fmt.Printf("Stereo:     %v\n", currentAudioOut.IsStereo())
	fmt.Printf("AVL:        %v\n", currentAudioOut.HasAVL())
	fmt.Println()

	// Select audio output if requested
	if selectAudioOut >= 0 {
		fmt.Printf("Selecting audio output %d...\n", selectAudioOut)
		err = dev.SetAudioOut(uint32(selectAudioOut))
		if err != nil {
			log.Fatalf("Failed to select audio output %d: %v", selectAudioOut, err)
		}

		// Verify the change
		newAudioOut, err := dev.GetCurrentAudioOut()
		if err != nil {
			log.Fatalf("Failed to verify audio output selection: %v", err)
		}

		if newAudioOut.GetIndex() == uint32(selectAudioOut) {
			fmt.Printf("✓ Successfully switched to audio output %d\n", selectAudioOut)
			fmt.Printf("  Name: %s\n", newAudioOut.GetName())
			fmt.Printf("  Stereo: %v\n", newAudioOut.IsStereo())
			fmt.Printf("  AVL: %v\n", newAudioOut.HasAVL())
		} else {
			fmt.Printf("⚠ Audio output selection may not have worked (got %d, expected %d)\n",
				newAudioOut.GetIndex(), selectAudioOut)
		}
	} else {
		fmt.Println("Tip: Use -s <index> to select a different audio output")
	}
}
