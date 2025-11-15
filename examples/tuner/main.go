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
	flag.IntVar(&setFreq, "f", -1, "set frequency (in device units, optional)")
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

	// Enumerate all tuners
	tuners, err := dev.GetAllTuners()
	if err != nil {
		fmt.Printf("Device does not support tuners: %v\n", err)
		fmt.Println("\nNote: Tuners are typically found on:")
		fmt.Println("  - TV tuner cards (analog/digital)")
		fmt.Println("  - FM/AM radio receivers")
		fmt.Println("  - Software Defined Radio (SDR) devices")
		fmt.Println("  - RF receivers")
		return
	}

	if len(tuners) == 0 {
		fmt.Println("No tuners found")
		return
	}

	fmt.Printf("Available tuners (%d):\n", len(tuners))
	fmt.Println("================================================================================")
	for _, tuner := range tuners {
		displayTuner(tuner)
	}

	// Get current frequency for first tuner
	if len(tuners) > 0 {
		freq, err := dev.GetFrequency(0)
		if err == nil {
			fmt.Println("\nCurrent Frequency (Tuner 0):")
			fmt.Println("================================================================================")
			displayFrequency(freq, tuners[0])
		}
	}

	// Show frequency bands if requested
	if showBands && len(tuners) > 0 {
		tuner := tuners[0]
		if tuner.SupportsFreqBands() {
			bands, err := dev.GetFrequencyBands(0, tuner.GetType())
			if err == nil && len(bands) > 0 {
				fmt.Println("\nFrequency Bands (Tuner 0):")
				fmt.Println("================================================================================")
				for _, band := range bands {
					displayBand(band, tuner)
				}
			} else if err != nil {
				fmt.Printf("\nFailed to get frequency bands: %v\n", err)
			}
		} else {
			fmt.Println("\nTuner does not support frequency band enumeration")
		}
	}

	// Set frequency if requested
	if setFreq >= 0 && len(tuners) > 0 {
		tuner := tuners[0]
		fmt.Printf("\nSetting frequency to %d...\n", setFreq)
		err = dev.SetFrequency(0, tuner.GetType(), uint32(setFreq))
		if err != nil {
			log.Fatalf("Failed to set frequency: %v", err)
		}

		// Verify the change
		newFreq, err := dev.GetFrequency(0)
		if err != nil {
			log.Fatalf("Failed to verify frequency: %v", err)
		}

		if newFreq.GetFrequency() == uint32(setFreq) {
			fmt.Printf("✓ Successfully set frequency to %d\n", setFreq)
		} else {
			fmt.Printf("⚠ Frequency may not have changed exactly (got %d, requested %d)\n",
				newFreq.GetFrequency(), setFreq)
		}
	} else if len(tuners) > 0 {
		fmt.Println("\nTip: Use -f <frequency> to tune to a specific frequency")
		if tuners[0].IsLowFreq() {
			fmt.Println("     Units are 1/16000 kHz (62.5 Hz)")
			fmt.Println("     Example: -f 1608000 for 100.5 MHz (FM radio)")
		} else {
			fmt.Println("     Units are 1/16 MHz (62.5 kHz)")
		}
		fmt.Println("Tip: Use -b to show available frequency bands")
	}
}

func displayTuner(tuner v4l2.TunerInfo) {
	typeName := v4l2.TunerTypes[tuner.GetType()]
	if typeName == "" {
		typeName = fmt.Sprintf("Unknown (0x%x)", tuner.GetType())
	}

	fmt.Printf("[%d] %s\n", tuner.GetIndex(), tuner.GetName())
	fmt.Printf("    Type:       %s\n", typeName)
	fmt.Printf("    Capability: 0x%08x\n", tuner.GetCapability())

	// Display frequency range
	rangeLow := tuner.GetRangeLow()
	rangeHigh := tuner.GetRangeHigh()
	if tuner.IsLowFreq() {
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

	// Display signal info
	signal := tuner.GetSignal()
	if signal > 0 {
		fmt.Printf("    Signal:     %d / 65535 (%.1f%%)\n", signal, float64(signal)/655.35)
	} else if signal == 0 {
		fmt.Printf("    Signal:     No signal detected\n")
	}

	afc := tuner.GetAFC()
	if afc != 0 {
		fmt.Printf("    AFC:        %d\n", afc)
	}

	// Display audio info
	audioMode := tuner.GetAudioMode()
	if modeName, ok := v4l2.TunerAudioModes[audioMode]; ok {
		fmt.Printf("    Audio Mode: %s\n", modeName)
	} else {
		fmt.Printf("    Audio Mode: 0x%x\n", audioMode)
	}

	rxSubchans := tuner.GetRxSubchans()
	if rxSubchans != 0 {
		fmt.Printf("    RxSubchans: 0x%x", rxSubchans)
		if rxSubchans&v4l2.TunerSubStereo != 0 {
			fmt.Print(" (Stereo)")
		}
		if rxSubchans&v4l2.TunerSubMono != 0 {
			fmt.Print(" (Mono)")
		}
		if rxSubchans&v4l2.TunerSubRDS != 0 {
			fmt.Print(" (RDS)")
		}
		fmt.Println()
	}

	// Display capability flags
	fmt.Print("    Features:   ")
	features := []string{}
	if tuner.IsStereo() {
		features = append(features, "Stereo")
	}
	if tuner.HasRDS() {
		features = append(features, "RDS")
	}
	if tuner.SupportsHwSeek() {
		features = append(features, "HW Seek")
	}
	if tuner.SupportsFreqBands() {
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

func displayFrequency(freq v4l2.FrequencyInfo, tuner v4l2.TunerInfo) {
	typeName := v4l2.TunerTypes[freq.GetType()]
	if typeName == "" {
		typeName = fmt.Sprintf("Unknown (0x%x)", freq.GetType())
	}

	fmt.Printf("Tuner:      %d\n", freq.GetTuner())
	fmt.Printf("Type:       %s\n", typeName)
	fmt.Printf("Frequency:  %d", freq.GetFrequency())

	if tuner.IsLowFreq() {
		// Units are 1/16000 kHz (62.5 Hz)
		fmt.Printf(" (%.3f MHz)\n", float64(freq.GetFrequency())/16000.0)
	} else {
		// Units are 1/16 MHz (62.5 kHz)
		fmt.Printf(" (%.1f MHz)\n", float64(freq.GetFrequency())/16.0)
	}
}

func displayBand(band v4l2.FrequencyBandInfo, tuner v4l2.TunerInfo) {
	fmt.Printf("Band %d:\n", band.GetIndex())

	// Display frequency range
	rangeLow := band.GetRangeLow()
	rangeHigh := band.GetRangeHigh()
	if tuner.IsLowFreq() {
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
