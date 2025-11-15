package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vladimirvivien/go4vl/device"
)

func main() {
	var devName string
	var selectAudio int
	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.IntVar(&selectAudio, "s", -1, "select audio input by index (optional)")
	flag.Parse()

	dev, err := device.Open(devName)
	if err != nil {
		log.Fatalf("Failed to open device %s: %v", devName, err)
	}
	defer dev.Close()

	fmt.Printf("Device: %s\n", dev.Name())
	fmt.Printf("Driver: %s\n", dev.Capability().Driver)
	fmt.Printf("Card: %s\n\n", dev.Capability().Card)

	// Get current audio input
	currentAudio, err := dev.GetCurrentAudio()
	if err != nil {
		fmt.Printf("Device does not support audio inputs: %v\n", err)
		fmt.Println("\nNote: Audio inputs are typically found on:")
		fmt.Println("  - TV tuner cards")
		fmt.Println("  - Webcams with microphones")
		fmt.Println("  - Video capture cards with audio")
		return
	}

	fmt.Printf("Current audio input: [%d] %s\n\n", currentAudio.GetIndex(), currentAudio.GetName())

	// Enumerate all audio inputs
	audios, err := dev.GetAudioDescriptions()
	if err != nil {
		log.Fatalf("Failed to enumerate audio inputs: %v", err)
	}

	if len(audios) == 0 {
		fmt.Println("No audio inputs found")
		return
	}

	fmt.Printf("Available audio inputs (%d):\n", len(audios))
	fmt.Println("================================================================================")
	for _, audio := range audios {
		active := ""
		if audio.GetIndex() == currentAudio.GetIndex() {
			active = " ** ACTIVE **"
		}

		fmt.Printf("[%d] %s%s\n", audio.GetIndex(), audio.GetName(), active)
		fmt.Printf("    Capability: 0x%08x\n", audio.GetCapability())
		fmt.Printf("    Mode:       0x%08x\n", audio.GetMode())

		// Display capability flags
		if audio.IsStereo() {
			fmt.Println("    ✓ Stereo")
		} else {
			fmt.Println("    • Mono")
		}

		if audio.HasAVL() {
			fmt.Println("    ✓ Automatic Volume Level (AVL)")
		}

		fmt.Println()
	}

	// Show current audio details
	fmt.Println("Current Audio Input Details:")
	fmt.Println("================================================================================")
	fmt.Printf("Index:      %d\n", currentAudio.GetIndex())
	fmt.Printf("Name:       %s\n", currentAudio.GetName())
	fmt.Printf("Capability: 0x%08x\n", currentAudio.GetCapability())
	fmt.Printf("Mode:       0x%08x\n", currentAudio.GetMode())
	fmt.Printf("Stereo:     %v\n", currentAudio.IsStereo())
	fmt.Printf("AVL:        %v\n", currentAudio.HasAVL())
	fmt.Println()

	// Select audio input if requested
	if selectAudio >= 0 {
		fmt.Printf("Selecting audio input %d...\n", selectAudio)
		err = dev.SetAudio(uint32(selectAudio))
		if err != nil {
			log.Fatalf("Failed to select audio input %d: %v", selectAudio, err)
		}

		// Verify the change
		newAudio, err := dev.GetCurrentAudio()
		if err != nil {
			log.Fatalf("Failed to verify audio input selection: %v", err)
		}

		if newAudio.GetIndex() == uint32(selectAudio) {
			fmt.Printf("✓ Successfully switched to audio input %d\n", selectAudio)
			fmt.Printf("  Name: %s\n", newAudio.GetName())
			fmt.Printf("  Stereo: %v\n", newAudio.IsStereo())
			fmt.Printf("  AVL: %v\n", newAudio.HasAVL())
		} else {
			fmt.Printf("⚠ Audio input selection may not have worked (got %d, expected %d)\n",
				newAudio.GetIndex(), selectAudio)
		}
	} else {
		fmt.Println("Tip: Use -s <index> to select a different audio input")
	}
}
