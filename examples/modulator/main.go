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
	var setFreq int
	var showBands bool
	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.IntVar(&setFreq, "f", -1, "set modulator frequency (in device units, optional)")
	flag.BoolVar(&showBands, "b", false, "show frequency bands")
	flag.Parse()

	dev, err := device.Open(devName)
	if err != nil {
		log.Fatalf("Failed to open device %s: %v", devName, err)
	}
	defer dev.Close()

	fmt.Printf("Device: %s\n", dev.Name())
	fmt.Printf("Driver: %s\n", dev.Capability().Driver)
	fmt.Printf("Card: %s\n\n", dev.Capability().Card)

	// Enumerate all modulators
	modulators, err := dev.GetAllModulators()
	if err != nil {
		fmt.Printf("Device does not support modulators: %v\n", err)
		fmt.Println("\nNote: Modulators are typically found on:")
		fmt.Println("  - RF modulators for video output")
		fmt.Println("  - TV transmitters")
		fmt.Println("  - FM/AM radio transmitters")
		fmt.Println("  - Software Defined Radio (SDR) devices with TX capability")
		return
	}

	if len(modulators) == 0 {
		fmt.Println("No modulators found")
		return
	}

	fmt.Printf("Available modulators (%d):\n", len(modulators))
	fmt.Println("================================================================================")
	for _, mod := range modulators {
		displayModulator(mod)
	}

	// Get current frequency for first modulator (uses same API as tuner)
	if len(modulators) > 0 {
		freq, err := dev.GetFrequency(0)
		if err == nil {
			fmt.Println("\nCurrent Frequency (Modulator 0):")
			fmt.Println("================================================================================")
			displayFrequency(freq, modulators[0])
		}
	}

	// Show frequency bands if requested
	if showBands && len(modulators) > 0 {
		mod := modulators[0]
		if mod.SupportsFreqBands() {
			bands, err := dev.GetFrequencyBands(0, mod.GetType())
			if err == nil && len(bands) > 0 {
				fmt.Println("\nFrequency Bands (Modulator 0):")
				fmt.Println("================================================================================")
				for _, band := range bands {
					displayBand(band, mod)
				}
			} else if err != nil {
				fmt.Printf("\nFailed to get frequency bands: %v\n", err)
			}
		} else {
			fmt.Println("\nModulator does not support frequency band enumeration")
		}
	}

	// Set frequency if requested
	if setFreq >= 0 && len(modulators) > 0 {
		mod := modulators[0]
		fmt.Printf("\nSetting modulator frequency to %d...\n", setFreq)
		err = dev.SetFrequency(0, mod.GetType(), uint32(setFreq))
		if err != nil {
			log.Fatalf("Failed to set frequency: %v", err)
		}

		// Verify the change
		newFreq, err := dev.GetFrequency(0)
		if err != nil {
			log.Fatalf("Failed to verify frequency: %v", err)
		}

		if newFreq.GetFrequency() == uint32(setFreq) {
			fmt.Printf("✓ Successfully set modulator frequency to %d\n", setFreq)
		} else {
			fmt.Printf("⚠ Frequency may not have changed exactly (got %d, requested %d)\n",
				newFreq.GetFrequency(), setFreq)
		}
	} else if len(modulators) > 0 {
		fmt.Println("\nTip: Use -f <frequency> to set modulator frequency")
		if modulators[0].IsLowFreq() {
			fmt.Println("     Units are 1/16000 kHz (62.5 Hz)")
			fmt.Println("     Example: -f 1608000 for 100.5 MHz (FM)")
		} else {
			fmt.Println("     Units are 1/16 MHz (62.5 kHz)")
		}
		fmt.Println("Tip: Use -b to show available frequency bands")
	}
}

func displayModulator(mod v4l2.ModulatorInfo) {
	typeName := v4l2.TunerTypes[mod.GetType()]
	if typeName == "" {
		typeName = fmt.Sprintf("Unknown (0x%x)", mod.GetType())
	}

	fmt.Printf("[%d] %s\n", mod.GetIndex(), mod.GetName())
	fmt.Printf("    Type:       %s\n", typeName)
	fmt.Printf("    Capability: 0x%08x\n", mod.GetCapability())

	// Display frequency range
	rangeLow := mod.GetRangeLow()
	rangeHigh := mod.GetRangeHigh()
	if mod.IsLowFreq() {
		// Units are 1/16000 kHz (62.5 Hz)
		fmt.Printf("    Range:      %d - %d (units: 62.5 Hz)\n", rangeLow, rangeHigh)
		fmt.Printf("                %.3f - %.3f MHz\n",
			float64(rangeLow)/16000.0, float64(rangeHigh)/16000.0)
	} else {
		// Units are 1/16 MHz (62.5 kHz)
		fmt.Printf("    Range:      %d - %d (units: 62.5 kHz)\n", rangeLow, rangeHigh)
		fmt.Printf("                %.1f - %.1f MHz\n",
			float64(rangeLow)/16.0, float64(rangeHigh)/16.0)
	}

	// Display transmission subchannels
	txSubchans := mod.GetTxSubchans()
	if txSubchans != 0 {
		fmt.Printf("    TxSubchans: 0x%x", txSubchans)
		if txSubchans&v4l2.TunerSubStereo != 0 {
			fmt.Print(" (Stereo)")
		}
		if txSubchans&v4l2.TunerSubMono != 0 {
			fmt.Print(" (Mono)")
		}
		if txSubchans&v4l2.TunerSubRDS != 0 {
			fmt.Print(" (RDS)")
		}
		fmt.Println()
	}

	// Display capability flags
	fmt.Print("    Features:   ")
	features := []string{}
	if mod.IsStereo() {
		features = append(features, "Stereo")
	}
	if mod.HasRDS() {
		features = append(features, "RDS")
	}
	if mod.SupportsFreqBands() {
		features = append(features, "Freq Bands")
	}
	if len(features) > 0 {
		for i, feat := range features {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(feat)
		}
	} else {
		fmt.Print("None")
	}
	fmt.Println()
	fmt.Println()
}

func displayFrequency(freq v4l2.FrequencyInfo, mod v4l2.ModulatorInfo) {
	typeName := v4l2.TunerTypes[freq.GetType()]
	if typeName == "" {
		typeName = fmt.Sprintf("Unknown (0x%x)", freq.GetType())
	}

	fmt.Printf("Modulator:  %d\n", freq.GetTuner())
	fmt.Printf("Type:       %s\n", typeName)
	fmt.Printf("Frequency:  %d", freq.GetFrequency())

	if mod.IsLowFreq() {
		// Units are 1/16000 kHz (62.5 Hz)
		fmt.Printf(" (%.3f MHz)\n", float64(freq.GetFrequency())/16000.0)
	} else {
		// Units are 1/16 MHz (62.5 kHz)
		fmt.Printf(" (%.1f MHz)\n", float64(freq.GetFrequency())/16.0)
	}
}

func displayBand(band v4l2.FrequencyBandInfo, mod v4l2.ModulatorInfo) {
	fmt.Printf("Band %d:\n", band.GetIndex())

	// Display frequency range
	rangeLow := band.GetRangeLow()
	rangeHigh := band.GetRangeHigh()
	if mod.IsLowFreq() {
		fmt.Printf("  Range:      %d - %d (%.3f - %.3f MHz)\n",
			rangeLow, rangeHigh,
			float64(rangeLow)/16000.0, float64(rangeHigh)/16000.0)
	} else {
		fmt.Printf("  Range:      %d - %d (%.1f - %.1f MHz)\n",
			rangeLow, rangeHigh,
			float64(rangeLow)/16.0, float64(rangeHigh)/16.0)
	}

	fmt.Printf("  Capability: 0x%08x\n", band.GetCapability())

	// Display modulation
	modulation := band.GetModulation()
	fmt.Printf("  Modulation: 0x%x", modulation)
	mods := []string{}
	if modulation&v4l2.BandModulationFM != 0 {
		mods = append(mods, "FM")
	}
	if modulation&v4l2.BandModulationAM != 0 {
		mods = append(mods, "AM")
	}
	if modulation&v4l2.BandModulationVSB != 0 {
		mods = append(mods, "VSB")
	}
	if len(mods) > 0 {
		fmt.Print(" (")
		for i, mod := range mods {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(mod)
		}
		fmt.Print(")")
	}
	fmt.Println()
	fmt.Println()
}
