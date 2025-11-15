package v4l2

// dv_timings.go provides Digital Video (DV) timing configuration and detection.
//
// DV timings are used for digital video interfaces like HDMI, DisplayPort, DVI, and SDI.
// Unlike analog video standards (PAL/NTSC), DV timings precisely describe the video format
// including resolution, refresh rate, and synchronization parameters.
//
// The V4L2 API provides the following operations:
//   - Enumerate supported DV timings
//   - Get/set current DV timings
//   - Auto-detect DV timings from input signal
//   - Query DV timing capabilities
//
// Common use cases:
//   - HDMI capture cards (1080p, 4K, etc.)
//   - Professional video equipment
//   - DisplayPort/DVI capture
//   - SDI video interfaces
//
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/dv-timings.html
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enum-dv-timings.html

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// DVTimingType represents the type of DV timing
type DVTimingType = uint32

const (
	DVTimingTypeBT6561120 DVTimingType = C.V4L2_DV_BT_656_1120 // BT.656/1120 timing
)

// DVInterlaced represents interlaced vs progressive format
type DVInterlaced = uint32

const (
	DVProgressive       DVInterlaced = C.V4L2_DV_PROGRESSIVE // Progressive scan
	DVInterlacedFormat  DVInterlaced = C.V4L2_DV_INTERLACED  // Interlaced scan
)

// DVPolarity represents sync signal polarities
type DVPolarity = uint32

const (
	DVVSyncPosPolarity DVPolarity = C.V4L2_DV_VSYNC_POS_POL // Positive vertical sync
	DVHSyncPosPolarity DVPolarity = C.V4L2_DV_HSYNC_POS_POL // Positive horizontal sync
)

// DVStandard represents DV timing standards
type DVStandard = uint32

const (
	DVStdCEA861 DVStandard = C.V4L2_DV_BT_STD_CEA861 // CEA-861 Digital TV Profile
	DVStdDMT    DVStandard = C.V4L2_DV_BT_STD_DMT    // VESA Discrete Monitor Timings
	DVStdCVT    DVStandard = C.V4L2_DV_BT_STD_CVT    // VESA Coordinated Video Timings
	DVStdGTF    DVStandard = C.V4L2_DV_BT_STD_GTF    // VESA Generalized Timing Formula
)

// DVFlag represents DV timing flags
type DVFlag = uint32

const (
	DVFlagReducedBlanking      DVFlag = C.V4L2_DV_FL_REDUCED_BLANKING        // Reduced blanking
	DVFlagCanReduceFPS         DVFlag = C.V4L2_DV_FL_CAN_REDUCE_FPS          // Can reduce FPS
	DVFlagReducedFPS           DVFlag = C.V4L2_DV_FL_REDUCED_FPS             // Reduced FPS
	DVFlagHalfLine             DVFlag = C.V4L2_DV_FL_HALF_LINE               // Half line
	DVFlagIsCEVideo            DVFlag = C.V4L2_DV_FL_IS_CE_VIDEO             // Consumer Electronics video
	DVFlagFirstFieldExtraLine  DVFlag = C.V4L2_DV_FL_FIRST_FIELD_EXTRA_LINE  // First field has extra line
	DVFlagHasPictureAspect     DVFlag = C.V4L2_DV_FL_HAS_PICTURE_ASPECT      // Has picture aspect ratio
	DVFlagHasCEA861VIC         DVFlag = C.V4L2_DV_FL_HAS_CEA861_VIC          // Has CEA-861 VIC code
	DVFlagHasHDMIVIC           DVFlag = C.V4L2_DV_FL_HAS_HDMI_VIC            // Has HDMI VIC code
	DVFlagCanDetectReducedFPS  DVFlag = C.V4L2_DV_FL_CAN_DETECT_REDUCED_FPS  // Can detect reduced FPS
)

// DVCapability represents DV timing capabilities
type DVCapability = uint32

const (
	DVCapInterlaced      DVCapability = C.V4L2_DV_BT_CAP_INTERLACED       // Interlaced formats supported
	DVCapProgressive     DVCapability = C.V4L2_DV_BT_CAP_PROGRESSIVE      // Progressive formats supported
	DVCapReducedBlanking DVCapability = C.V4L2_DV_BT_CAP_REDUCED_BLANKING // Reduced blanking supported
	DVCapCustom          DVCapability = C.V4L2_DV_BT_CAP_CUSTOM           // Custom timings supported
)

// BTTimings wraps v4l2_bt_timings structure (BT.656/1120 timings)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/dv-timings.html
type BTTimings struct {
	v4l2BTTimings C.struct_v4l2_bt_timings
}

// Accessor methods for BTTimings

func (bt BTTimings) GetWidth() uint32 {
	return uint32(bt.v4l2BTTimings.width)
}

func (bt BTTimings) GetHeight() uint32 {
	return uint32(bt.v4l2BTTimings.height)
}

func (bt BTTimings) GetInterlaced() DVInterlaced {
	return DVInterlaced(bt.v4l2BTTimings.interlaced)
}

func (bt BTTimings) GetPolarities() DVPolarity {
	return DVPolarity(bt.v4l2BTTimings.polarities)
}

func (bt BTTimings) GetPixelClock() uint64 {
	return uint64(bt.v4l2BTTimings.pixelclock)
}

func (bt BTTimings) GetHFrontPorch() uint32 {
	return uint32(bt.v4l2BTTimings.hfrontporch)
}

func (bt BTTimings) GetHSync() uint32 {
	return uint32(bt.v4l2BTTimings.hsync)
}

func (bt BTTimings) GetHBackPorch() uint32 {
	return uint32(bt.v4l2BTTimings.hbackporch)
}

func (bt BTTimings) GetVFrontPorch() uint32 {
	return uint32(bt.v4l2BTTimings.vfrontporch)
}

func (bt BTTimings) GetVSync() uint32 {
	return uint32(bt.v4l2BTTimings.vsync)
}

func (bt BTTimings) GetVBackPorch() uint32 {
	return uint32(bt.v4l2BTTimings.vbackporch)
}

func (bt BTTimings) GetILVFrontPorch() uint32 {
	return uint32(bt.v4l2BTTimings.il_vfrontporch)
}

func (bt BTTimings) GetILVSync() uint32 {
	return uint32(bt.v4l2BTTimings.il_vsync)
}

func (bt BTTimings) GetILVBackPorch() uint32 {
	return uint32(bt.v4l2BTTimings.il_vbackporch)
}

func (bt BTTimings) GetStandards() DVStandard {
	return DVStandard(bt.v4l2BTTimings.standards)
}

func (bt BTTimings) GetFlags() DVFlag {
	return DVFlag(bt.v4l2BTTimings.flags)
}

func (bt BTTimings) GetCEA861VIC() uint8 {
	return uint8(bt.v4l2BTTimings.cea861_vic)
}

func (bt BTTimings) GetHDMIVIC() uint8 {
	return uint8(bt.v4l2BTTimings.hdmi_vic)
}

// Helper methods for BTTimings

func (bt BTTimings) IsInterlaced() bool {
	return bt.GetInterlaced() == DVInterlacedFormat
}

func (bt BTTimings) IsProgressive() bool {
	return bt.GetInterlaced() == DVProgressive
}

func (bt BTTimings) HasVSyncPosPolarity() bool {
	return (bt.GetPolarities() & DVVSyncPosPolarity) != 0
}

func (bt BTTimings) HasHSyncPosPolarity() bool {
	return (bt.GetPolarities() & DVHSyncPosPolarity) != 0
}

func (bt BTTimings) HasFlag(flag DVFlag) bool {
	return (bt.GetFlags() & flag) != 0
}

func (bt BTTimings) HasStandard(std DVStandard) bool {
	return (bt.GetStandards() & std) != 0
}

// GetFrameRate calculates the frame rate from pixel clock and total pixels
// Returns frames per second
func (bt BTTimings) GetFrameRate() float64 {
	pixelClock := float64(bt.GetPixelClock())
	if pixelClock == 0 {
		return 0
	}

	// Total horizontal pixels
	hTotal := bt.GetWidth() + bt.GetHFrontPorch() + bt.GetHSync() + bt.GetHBackPorch()

	// Total vertical lines
	vTotal := bt.GetHeight() + bt.GetVFrontPorch() + bt.GetVSync() + bt.GetVBackPorch()

	// For interlaced, add interlaced blanking
	if bt.IsInterlaced() {
		vTotal += bt.GetILVFrontPorch() + bt.GetILVSync() + bt.GetILVBackPorch()
	}

	totalPixels := float64(hTotal) * float64(vTotal)
	if totalPixels == 0 {
		return 0
	}

	return pixelClock / totalPixels
}

// DVTimings wraps v4l2_dv_timings structure
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-dv-timings.html
type DVTimings struct {
	v4l2DVTimings C.struct_v4l2_dv_timings
}

// Accessor methods for DVTimings

func (dv DVTimings) GetType() DVTimingType {
	return DVTimingType(dv.v4l2DVTimings._type)
}

func (dv DVTimings) GetBTTimings() BTTimings {
	// Access union field through type cast
	btTimings := (*C.struct_v4l2_bt_timings)(unsafe.Pointer(&dv.v4l2DVTimings.anon0[0]))
	return BTTimings{v4l2BTTimings: *btTimings}
}

// EnumDVTimings wraps v4l2_enum_dv_timings structure
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enum-dv-timings.html
type EnumDVTimings struct {
	v4l2EnumDVTimings C.struct_v4l2_enum_dv_timings
}

// Accessor methods for EnumDVTimings

func (e EnumDVTimings) GetIndex() uint32 {
	return uint32(e.v4l2EnumDVTimings.index)
}

func (e EnumDVTimings) GetPad() uint32 {
	return uint32(e.v4l2EnumDVTimings.pad)
}

func (e EnumDVTimings) GetTimings() DVTimings {
	return DVTimings{v4l2DVTimings: e.v4l2EnumDVTimings.timings}
}

// BTTimingsCap wraps v4l2_bt_timings_cap structure
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-dv-timings-cap.html
type BTTimingsCap struct {
	v4l2BTTimingsCap C.struct_v4l2_bt_timings_cap
}

// Accessor methods for BTTimingsCap

func (btc BTTimingsCap) GetMinWidth() uint32 {
	return uint32(btc.v4l2BTTimingsCap.min_width)
}

func (btc BTTimingsCap) GetMaxWidth() uint32 {
	return uint32(btc.v4l2BTTimingsCap.max_width)
}

func (btc BTTimingsCap) GetMinHeight() uint32 {
	return uint32(btc.v4l2BTTimingsCap.min_height)
}

func (btc BTTimingsCap) GetMaxHeight() uint32 {
	return uint32(btc.v4l2BTTimingsCap.max_height)
}

func (btc BTTimingsCap) GetMinPixelClock() uint64 {
	return uint64(btc.v4l2BTTimingsCap.min_pixelclock)
}

func (btc BTTimingsCap) GetMaxPixelClock() uint64 {
	return uint64(btc.v4l2BTTimingsCap.max_pixelclock)
}

func (btc BTTimingsCap) GetStandards() DVStandard {
	return DVStandard(btc.v4l2BTTimingsCap.standards)
}

func (btc BTTimingsCap) GetCapabilities() DVCapability {
	return DVCapability(btc.v4l2BTTimingsCap.capabilities)
}

// Helper methods for BTTimingsCap

func (btc BTTimingsCap) HasCapability(cap DVCapability) bool {
	return (btc.GetCapabilities() & cap) != 0
}

func (btc BTTimingsCap) SupportsInterlaced() bool {
	return btc.HasCapability(DVCapInterlaced)
}

func (btc BTTimingsCap) SupportsProgressive() bool {
	return btc.HasCapability(DVCapProgressive)
}

func (btc BTTimingsCap) SupportsReducedBlanking() bool {
	return btc.HasCapability(DVCapReducedBlanking)
}

func (btc BTTimingsCap) SupportsCustomTimings() bool {
	return btc.HasCapability(DVCapCustom)
}

func (btc BTTimingsCap) HasStandard(std DVStandard) bool {
	return (btc.GetStandards() & std) != 0
}

// DVTimingsCap wraps v4l2_dv_timings_cap structure
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-dv-timings-cap.html
type DVTimingsCap struct {
	v4l2DVTimingsCap C.struct_v4l2_dv_timings_cap
}

// Accessor methods for DVTimingsCap

func (dvc DVTimingsCap) GetType() DVTimingType {
	return DVTimingType(dvc.v4l2DVTimingsCap._type)
}

func (dvc DVTimingsCap) GetPad() uint32 {
	return uint32(dvc.v4l2DVTimingsCap.pad)
}

func (dvc DVTimingsCap) GetBTTimingsCap() BTTimingsCap {
	// Access union field through type cast
	btCap := (*C.struct_v4l2_bt_timings_cap)(unsafe.Pointer(&dvc.v4l2DVTimingsCap.anon0[0]))
	return BTTimingsCap{v4l2BTTimingsCap: *btCap}
}

// ============================================================================
// DV Timings Operations
// ============================================================================

// GetDVTimings gets the current DV timings
// Implements VIDIOC_G_DV_TIMINGS ioctl
//
// Example:
//
//	timings, err := v4l2.GetDVTimings(fd)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	bt := timings.GetBTTimings()
//	fmt.Printf("Resolution: %dx%d @ %.2f Hz\n", bt.GetWidth(), bt.GetHeight(), bt.GetFrameRate())
func GetDVTimings(fd uintptr) (DVTimings, error) {
	var timings C.struct_v4l2_dv_timings

	if err := send(fd, C.VIDIOC_G_DV_TIMINGS, uintptr(unsafe.Pointer(&timings))); err != nil {
		return DVTimings{}, fmt.Errorf("v4l2: VIDIOC_G_DV_TIMINGS failed: %w", err)
	}

	return DVTimings{v4l2DVTimings: timings}, nil
}

// SetDVTimings sets the DV timings
// Implements VIDIOC_S_DV_TIMINGS ioctl
//
// Example:
//
//	err := v4l2.SetDVTimings(fd, timings)
//	if err != nil {
//	    log.Fatal(err)
//	}
func SetDVTimings(fd uintptr, timings DVTimings) error {
	if err := send(fd, C.VIDIOC_S_DV_TIMINGS, uintptr(unsafe.Pointer(&timings.v4l2DVTimings))); err != nil {
		return fmt.Errorf("v4l2: VIDIOC_S_DV_TIMINGS failed: %w", err)
	}
	return nil
}

// QueryDVTimings attempts to auto-detect DV timings from the input signal
// Implements VIDIOC_QUERY_DV_TIMINGS ioctl
//
// This is useful for HDMI capture cards that can detect the incoming signal format.
//
// Example:
//
//	timings, err := v4l2.QueryDVTimings(fd)
//	if err != nil {
//	    log.Fatal("No signal detected or invalid timings:", err)
//	}
//	fmt.Printf("Detected: %dx%d\n", timings.GetBTTimings().GetWidth(), timings.GetBTTimings().GetHeight())
func QueryDVTimings(fd uintptr) (DVTimings, error) {
	var timings C.struct_v4l2_dv_timings

	if err := send(fd, C.VIDIOC_QUERY_DV_TIMINGS, uintptr(unsafe.Pointer(&timings))); err != nil {
		return DVTimings{}, fmt.Errorf("v4l2: VIDIOC_QUERY_DV_TIMINGS failed: %w", err)
	}

	return DVTimings{v4l2DVTimings: timings}, nil
}

// EnumerateDVTimings enumerates a specific DV timing by index
// Implements VIDIOC_ENUM_DV_TIMINGS ioctl
//
// Example:
//
//	enumTiming, err := v4l2.EnumerateDVTimings(fd, 0, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	timings := enumTiming.GetTimings()
func EnumerateDVTimings(fd uintptr, index uint32, pad uint32) (EnumDVTimings, error) {
	var enumTimings C.struct_v4l2_enum_dv_timings
	enumTimings.index = C.uint(index)
	enumTimings.pad = C.uint(pad)

	if err := send(fd, C.VIDIOC_ENUM_DV_TIMINGS, uintptr(unsafe.Pointer(&enumTimings))); err != nil {
		return EnumDVTimings{}, fmt.Errorf("v4l2: VIDIOC_ENUM_DV_TIMINGS failed for index %d: %w", index, err)
	}

	return EnumDVTimings{v4l2EnumDVTimings: enumTimings}, nil
}

// GetAllDVTimings enumerates all supported DV timings
// Returns a slice of EnumDVTimings or an error
//
// Example:
//
//	timings, err := v4l2.GetAllDVTimings(fd, 0)
//	for _, timing := range timings {
//	    bt := timing.GetTimings().GetBTTimings()
//	    fmt.Printf("Supported: %dx%d @ %.2f Hz\n", bt.GetWidth(), bt.GetHeight(), bt.GetFrameRate())
//	}
func GetAllDVTimings(fd uintptr, pad uint32) ([]EnumDVTimings, error) {
	var result []EnumDVTimings

	for i := uint32(0); i < 256; i++ {
		timing, err := EnumerateDVTimings(fd, i, pad)
		if err != nil {
			// EINVAL indicates no more timings
			if errno, ok := err.(sys.Errno); ok && errno == sys.EINVAL && len(result) > 0 {
				break
			}
			return result, err
		}
		result = append(result, timing)
	}

	return result, nil
}

// GetDVTimingsCap gets the DV timing capabilities
// Implements VIDIOC_DV_TIMINGS_CAP ioctl
//
// Example:
//
//	cap, err := v4l2.GetDVTimingsCap(fd, 0)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	btCap := cap.GetBTTimingsCap()
//	fmt.Printf("Supports: %dx%d to %dx%d\n",
//	    btCap.GetMinWidth(), btCap.GetMinHeight(),
//	    btCap.GetMaxWidth(), btCap.GetMaxHeight())
func GetDVTimingsCap(fd uintptr, pad uint32) (DVTimingsCap, error) {
	var cap C.struct_v4l2_dv_timings_cap
	cap.pad = C.uint(pad)
	cap._type = C.V4L2_DV_BT_656_1120

	if err := send(fd, C.VIDIOC_DV_TIMINGS_CAP, uintptr(unsafe.Pointer(&cap))); err != nil {
		return DVTimingsCap{}, fmt.Errorf("v4l2: VIDIOC_DV_TIMINGS_CAP failed: %w", err)
	}

	return DVTimingsCap{v4l2DVTimingsCap: cap}, nil
}
