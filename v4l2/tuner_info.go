package v4l2

// tuner_info.go provides tuner and modulator enumeration and control.
//
// Tuners and modulators are used for radio and TV devices to tune to specific frequencies.
// A tuner receives signals (e.g., FM/AM radio, analog/digital TV, SDR), while a modulator
// transmits signals (e.g., RF modulator for video output).
//
// The V4L2 API provides the following operations:
//   - Enumerate and query tuner/modulator properties
//   - Get/set tuner parameters (audio mode, frequency)
//   - Get/set modulator parameters (frequency, transmission subchannel)
//   - Query signal strength and AFC (Automatic Frequency Control)
//   - Enumerate frequency bands
//
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-tuner.html
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-modulator.html
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-frequency.html

import (
	"fmt"
	"unsafe"
)

// TunerType represents the type of tuner/modulator
type TunerType = uint32

const (
	TunerTypeRadio     TunerType = cV4L2_TUNER_RADIO      // Radio tuner (AM/FM)
	TunerTypeAnalogTV  TunerType = cV4L2_TUNER_ANALOG_TV  // Analog TV tuner
	TunerTypeDigitalTV TunerType = cV4L2_TUNER_DIGITAL_TV // Digital TV tuner
	TunerTypeSDR       TunerType = cV4L2_TUNER_SDR        // Software Defined Radio
	TunerTypeRF        TunerType = cV4L2_TUNER_RF         // RF tuner
)

// TunerTypes maps tuner type constants to human-readable strings
var TunerTypes = map[TunerType]string{
	TunerTypeRadio:     "Radio",
	TunerTypeAnalogTV:  "Analog TV",
	TunerTypeDigitalTV: "Digital TV",
	TunerTypeSDR:       "SDR",
	TunerTypeRF:        "RF",
}

// TunerCapability represents tuner/modulator capability flags
type TunerCapability = uint32

const (
	TunerCapLow            TunerCapability = cV4L2_TUNER_CAP_LOW             // Freq in 1/16 kHz units (vs 1/16 MHz)
	TunerCapNorm           TunerCapability = cV4L2_TUNER_CAP_NORM            // Multi-standard tuner
	TunerCapHwSeekBounded  TunerCapability = cV4L2_TUNER_CAP_HWSEEK_BOUNDED  // Hardware seek bounded
	TunerCapHwSeekWrap     TunerCapability = cV4L2_TUNER_CAP_HWSEEK_WRAP     // Hardware seek wraps around
	TunerCapStereo         TunerCapability = cV4L2_TUNER_CAP_STEREO          // Stereo reception/transmission
	TunerCapLang2          TunerCapability = cV4L2_TUNER_CAP_LANG2           // Bilingual mode (LANG2/SAP)
	TunerCapSAP            TunerCapability = cV4L2_TUNER_CAP_SAP             // SAP (Second Audio Program)
	TunerCapLang1          TunerCapability = cV4L2_TUNER_CAP_LANG1           // Primary language
	TunerCapRDS            TunerCapability = cV4L2_TUNER_CAP_RDS             // Radio Data System
	TunerCapRDSBlockIO     TunerCapability = cV4L2_TUNER_CAP_RDS_BLOCK_IO    // RDS block I/O interface
	TunerCapRDSControls    TunerCapability = cV4L2_TUNER_CAP_RDS_CONTROLS    // RDS controls available
	TunerCapFreqBands      TunerCapability = cV4L2_TUNER_CAP_FREQ_BANDS      // Multiple frequency bands
	TunerCapHwSeekProgLim  TunerCapability = cV4L2_TUNER_CAP_HWSEEK_PROG_LIM // Programmable hardware seek limits
	TunerCap1Hz            TunerCapability = cV4L2_TUNER_CAP_1HZ             // 1 Hz frequency step
)

// TunerRxSubchannel represents received audio subchannels
type TunerRxSubchannel = uint32

const (
	TunerSubMono   TunerRxSubchannel = cV4L2_TUNER_SUB_MONO   // Mono audio detected
	TunerSubStereo TunerRxSubchannel = cV4L2_TUNER_SUB_STEREO // Stereo audio detected
	TunerSubLang2  TunerRxSubchannel = cV4L2_TUNER_SUB_LANG2  // LANG2/SAP detected
	TunerSubSAP    TunerRxSubchannel = cV4L2_TUNER_SUB_SAP    // SAP detected
	TunerSubLang1  TunerRxSubchannel = cV4L2_TUNER_SUB_LANG1  // LANG1 detected
	TunerSubRDS    TunerRxSubchannel = cV4L2_TUNER_SUB_RDS    // RDS data available
)

// TunerAudioMode represents the audio mode to select
type TunerAudioMode = uint32

const (
	TunerModeMono       TunerAudioMode = cV4L2_TUNER_MODE_MONO        // Mono audio
	TunerModeStereo     TunerAudioMode = cV4L2_TUNER_MODE_STEREO      // Stereo audio
	TunerModeLang2      TunerAudioMode = cV4L2_TUNER_MODE_LANG2       // LANG2/SAP
	TunerModeSAP        TunerAudioMode = cV4L2_TUNER_MODE_SAP         // SAP (same as LANG2)
	TunerModeLang1      TunerAudioMode = cV4L2_TUNER_MODE_LANG1       // LANG1
	TunerModeLang1Lang2 TunerAudioMode = cV4L2_TUNER_MODE_LANG1_LANG2 // LANG1 + LANG2
)

// TunerAudioModes maps audio mode constants to human-readable strings
var TunerAudioModes = map[TunerAudioMode]string{
	TunerModeMono:       "Mono",
	TunerModeStereo:     "Stereo",
	TunerModeLang2:      "Lang2/SAP",
	TunerModeLang1:      "Lang1",
	TunerModeLang1Lang2: "Lang1+Lang2",
}

// BandModulation represents frequency band modulation types
type BandModulation = uint32

const (
	BandModulationVSB BandModulation = cV4L2_BAND_MODULATION_VSB // Vestigial Sideband
	BandModulationFM  BandModulation = cV4L2_BAND_MODULATION_FM  // Frequency Modulation
	BandModulationAM  BandModulation = cV4L2_BAND_MODULATION_AM  // Amplitude Modulation
)

// TunerInfo wraps v4l2_tuner structure
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-tuner.html
type TunerInfo struct {
	v4l2Tuner cV4L2Tuner
}

// Accessor methods for TunerInfo

func (t TunerInfo) GetIndex() uint32 {
	return t.v4l2Tuner.Index
}

func (t TunerInfo) GetName() string {
	return cString(t.v4l2Tuner.Name[:])
}

func (t TunerInfo) GetType() TunerType {
	return TunerType(t.v4l2Tuner.Type)
}

func (t TunerInfo) GetCapability() TunerCapability {
	return TunerCapability(t.v4l2Tuner.Capability)
}

func (t TunerInfo) GetRangeLow() uint32 {
	return t.v4l2Tuner.Rangelow
}

func (t TunerInfo) GetRangeHigh() uint32 {
	return t.v4l2Tuner.Rangehigh
}

func (t TunerInfo) GetRxSubchans() TunerRxSubchannel {
	return TunerRxSubchannel(t.v4l2Tuner.Rxsubchans)
}

func (t TunerInfo) GetAudioMode() TunerAudioMode {
	return TunerAudioMode(t.v4l2Tuner.Audmode)
}

func (t TunerInfo) GetSignal() int32 {
	return t.v4l2Tuner.Signal
}

func (t TunerInfo) GetAFC() int32 {
	return t.v4l2Tuner.Afc
}

// Capability check helper methods for TunerInfo

func (t TunerInfo) HasCapability(cap TunerCapability) bool {
	return (t.GetCapability() & cap) != 0
}

func (t TunerInfo) IsLowFreq() bool {
	return t.HasCapability(TunerCapLow)
}

func (t TunerInfo) IsStereo() bool {
	return t.HasCapability(TunerCapStereo)
}

func (t TunerInfo) HasRDS() bool {
	return t.HasCapability(TunerCapRDS)
}

func (t TunerInfo) SupportsHwSeek() bool {
	return t.HasCapability(TunerCapHwSeekBounded) || t.HasCapability(TunerCapHwSeekWrap)
}

func (t TunerInfo) SupportsFreqBands() bool {
	return t.HasCapability(TunerCapFreqBands)
}

// ModulatorInfo wraps v4l2_modulator structure
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-modulator.html
type ModulatorInfo struct {
	v4l2Modulator cV4L2Modulator
}

// Accessor methods for ModulatorInfo

func (m ModulatorInfo) GetIndex() uint32 {
	return m.v4l2Modulator.Index
}

func (m ModulatorInfo) GetName() string {
	return cString(m.v4l2Modulator.Name[:])
}

func (m ModulatorInfo) GetCapability() TunerCapability {
	return TunerCapability(m.v4l2Modulator.Capability)
}

func (m ModulatorInfo) GetRangeLow() uint32 {
	return m.v4l2Modulator.Rangelow
}

func (m ModulatorInfo) GetRangeHigh() uint32 {
	return m.v4l2Modulator.Rangehigh
}

func (m ModulatorInfo) GetTxSubchans() uint32 {
	return m.v4l2Modulator.Txsubchans
}

func (m ModulatorInfo) GetType() TunerType {
	return TunerType(m.v4l2Modulator.Type)
}

// Capability check helper methods for ModulatorInfo

func (m ModulatorInfo) HasCapability(cap TunerCapability) bool {
	return (m.GetCapability() & cap) != 0
}

func (m ModulatorInfo) IsLowFreq() bool {
	return m.HasCapability(TunerCapLow)
}

func (m ModulatorInfo) IsStereo() bool {
	return m.HasCapability(TunerCapStereo)
}

func (m ModulatorInfo) HasRDS() bool {
	return m.HasCapability(TunerCapRDS)
}

func (m ModulatorInfo) SupportsFreqBands() bool {
	return m.HasCapability(TunerCapFreqBands)
}

// FrequencyInfo wraps v4l2_frequency structure
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-frequency.html
type FrequencyInfo struct {
	v4l2Frequency cV4L2Frequency
}

// Accessor methods for FrequencyInfo

func (f FrequencyInfo) GetTuner() uint32 {
	return f.v4l2Frequency.Tuner
}

func (f FrequencyInfo) GetType() TunerType {
	return TunerType(f.v4l2Frequency.Type)
}

func (f FrequencyInfo) GetFrequency() uint32 {
	return f.v4l2Frequency.Frequency
}

// FrequencyBandInfo wraps v4l2_frequency_band structure
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enum-freq-bands.html
type FrequencyBandInfo struct {
	v4l2FreqBand cV4L2FrequencyBand
}

// Accessor methods for FrequencyBandInfo

func (fb FrequencyBandInfo) GetTuner() uint32 {
	return fb.v4l2FreqBand.Tuner
}

func (fb FrequencyBandInfo) GetType() TunerType {
	return TunerType(fb.v4l2FreqBand.Type)
}

func (fb FrequencyBandInfo) GetIndex() uint32 {
	return fb.v4l2FreqBand.Index
}

func (fb FrequencyBandInfo) GetCapability() TunerCapability {
	return TunerCapability(fb.v4l2FreqBand.Capability)
}

func (fb FrequencyBandInfo) GetRangeLow() uint32 {
	return fb.v4l2FreqBand.Rangelow
}

func (fb FrequencyBandInfo) GetRangeHigh() uint32 {
	return fb.v4l2FreqBand.Rangehigh
}

func (fb FrequencyBandInfo) GetModulation() BandModulation {
	return BandModulation(fb.v4l2FreqBand.Modulation)
}

func (fb FrequencyBandInfo) HasCapability(cap TunerCapability) bool {
	return (fb.GetCapability() & cap) != 0
}

// ============================================================================
// Tuner Operations
// ============================================================================

// GetTunerInfo gets information about a specific tuner by index
// Implements VIDIOC_G_TUNER ioctl
//
// Example:
//
//	tuner, err := v4l2.GetTunerInfo(fd, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Tuner: %s, Type: %s\n", tuner.GetName(), v4l2.TunerTypes[tuner.GetType()])
func GetTunerInfo(fd uintptr, index uint32) (TunerInfo, error) {
	var tuner cV4L2Tuner
	tuner.Index = index

	if err := send(fd, cVIDIOC_G_TUNER, unsafe.Pointer(&tuner)); err != nil {
		return TunerInfo{}, fmt.Errorf("v4l2: VIDIOC_G_TUNER failed for index %d: %w", index, err)
	}

	return TunerInfo{v4l2Tuner: tuner}, nil
}

// SetTuner sets tuner parameters (audio mode)
// Implements VIDIOC_S_TUNER ioctl
//
// Example:
//
//	tuner, _ := v4l2.GetTunerInfo(fd, 0)
//	tuner.v4l2Tuner.audmode = C.V4L2_TUNER_MODE_STEREO
//	err := v4l2.SetTuner(fd, tuner)
func SetTuner(fd uintptr, tuner TunerInfo) error {
	if err := send(fd, cVIDIOC_S_TUNER, unsafe.Pointer(&tuner.v4l2Tuner)); err != nil {
		return fmt.Errorf("v4l2: VIDIOC_S_TUNER failed: %w", err)
	}
	return nil
}

// GetAllTuners enumerates all tuners available on the device
// Returns a slice of TunerInfo or an error
//
// Example:
//
//	tuners, err := v4l2.GetAllTuners(fd)
//	for _, tuner := range tuners {
//	    fmt.Printf("Tuner %d: %s\n", tuner.GetIndex(), tuner.GetName())
//	}
func GetAllTuners(fd uintptr) ([]TunerInfo, error) {
	return enumIndexed(false, func(i uint32) (TunerInfo, error) {
		return GetTunerInfo(fd, i)
	})
}

// ============================================================================
// Modulator Operations
// ============================================================================

// GetModulatorInfo gets information about a specific modulator by index
// Implements VIDIOC_G_MODULATOR ioctl
//
// Example:
//
//	mod, err := v4l2.GetModulatorInfo(fd, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Modulator: %s\n", mod.GetName())
func GetModulatorInfo(fd uintptr, index uint32) (ModulatorInfo, error) {
	var modulator cV4L2Modulator
	modulator.Index = index

	if err := send(fd, cVIDIOC_G_MODULATOR, unsafe.Pointer(&modulator)); err != nil {
		return ModulatorInfo{}, fmt.Errorf("v4l2: VIDIOC_G_MODULATOR failed for index %d: %w", index, err)
	}

	return ModulatorInfo{v4l2Modulator: modulator}, nil
}

// SetModulator sets modulator parameters (transmission subchannels)
// Implements VIDIOC_S_MODULATOR ioctl
//
// Example:
//
//	mod, _ := v4l2.GetModulatorInfo(fd, 0)
//	mod.v4l2Modulator.txsubchans = C.V4L2_TUNER_SUB_STEREO
//	err := v4l2.SetModulator(fd, mod)
func SetModulator(fd uintptr, modulator ModulatorInfo) error {
	if err := send(fd, cVIDIOC_S_MODULATOR, unsafe.Pointer(&modulator.v4l2Modulator)); err != nil {
		return fmt.Errorf("v4l2: VIDIOC_S_MODULATOR failed: %w", err)
	}
	return nil
}

// GetAllModulators enumerates all modulators available on the device
// Returns a slice of ModulatorInfo or an error
//
// Example:
//
//	modulators, err := v4l2.GetAllModulators(fd)
//	for _, mod := range modulators {
//	    fmt.Printf("Modulator %d: %s\n", mod.GetIndex(), mod.GetName())
//	}
func GetAllModulators(fd uintptr) ([]ModulatorInfo, error) {
	return enumIndexed(false, func(i uint32) (ModulatorInfo, error) {
		return GetModulatorInfo(fd, i)
	})
}

// ============================================================================
// Frequency Operations
// ============================================================================

// GetFrequency gets the current tuner/modulator frequency
// Implements VIDIOC_G_FREQUENCY ioctl
//
// Example:
//
//	freq, err := v4l2.GetFrequency(fd, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Current frequency: %d\n", freq.GetFrequency())
func GetFrequency(fd uintptr, tunerIndex uint32) (FrequencyInfo, error) {
	var frequency cV4L2Frequency
	frequency.Tuner = tunerIndex

	if err := send(fd, cVIDIOC_G_FREQUENCY, unsafe.Pointer(&frequency)); err != nil {
		return FrequencyInfo{}, fmt.Errorf("v4l2: VIDIOC_G_FREQUENCY failed for tuner %d: %w", tunerIndex, err)
	}

	return FrequencyInfo{v4l2Frequency: frequency}, nil
}

// SetFrequency sets the tuner/modulator frequency
// Implements VIDIOC_S_FREQUENCY ioctl
//
// Example:
//
//	// Set FM radio to 100.5 MHz (assuming 62.5 Hz units, which is 1/16000 kHz)
//	// 100.5 MHz = 100500 kHz = 100500 * 16 units = 1608000 units
//	err := v4l2.SetFrequency(fd, 0, v4l2.TunerTypeRadio, 1608000)
func SetFrequency(fd uintptr, tunerIndex uint32, tunerType TunerType, frequency uint32) error {
	var freq cV4L2Frequency
	freq.Tuner = tunerIndex
	freq.Type = tunerType
	freq.Frequency = frequency

	if err := send(fd, cVIDIOC_S_FREQUENCY, unsafe.Pointer(&freq)); err != nil {
		return fmt.Errorf("v4l2: VIDIOC_S_FREQUENCY failed: %w", err)
	}
	return nil
}

// ============================================================================
// Frequency Band Operations
// ============================================================================

// GetFrequencyBandInfo gets information about a specific frequency band
// Implements VIDIOC_ENUM_FREQ_BANDS ioctl
//
// Example:
//
//	band, err := v4l2.GetFrequencyBandInfo(fd, 0, v4l2.TunerTypeRadio, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Band: %d - %d\n", band.GetRangeLow(), band.GetRangeHigh())
func GetFrequencyBandInfo(fd uintptr, tunerIndex uint32, tunerType TunerType, bandIndex uint32) (FrequencyBandInfo, error) {
	var band cV4L2FrequencyBand
	band.Tuner = tunerIndex
	band.Type = tunerType
	band.Index = bandIndex

	if err := send(fd, cVIDIOC_ENUM_FREQ_BANDS, unsafe.Pointer(&band)); err != nil {
		return FrequencyBandInfo{}, fmt.Errorf("v4l2: VIDIOC_ENUM_FREQ_BANDS failed for tuner %d band %d: %w", tunerIndex, bandIndex, err)
	}

	return FrequencyBandInfo{v4l2FreqBand: band}, nil
}

// GetAllFrequencyBands enumerates all frequency bands for a specific tuner
// Returns a slice of FrequencyBandInfo or an error
//
// Example:
//
//	bands, err := v4l2.GetAllFrequencyBands(fd, 0, v4l2.TunerTypeRadio)
//	for _, band := range bands {
//	    fmt.Printf("Band %d: %d - %d\n", band.GetIndex(), band.GetRangeLow(), band.GetRangeHigh())
//	}
func GetAllFrequencyBands(fd uintptr, tunerIndex uint32, tunerType TunerType) ([]FrequencyBandInfo, error) {
	return enumIndexed(false, func(i uint32) (FrequencyBandInfo, error) {
		return GetFrequencyBandInfo(fd, tunerIndex, tunerType, i)
	})
}
