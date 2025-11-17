package v4l2

// standard.go provides analog video standard enumeration and selection.
//
// Video standards define the analog video signal format (PAL, NTSC, SECAM, etc.)
// used by legacy analog video devices like TV tuners, composite video inputs,
// and analog capture cards.
//
// The V4L2 API provides the following operations:
//   - Enumerate supported video standards
//   - Get/set current video standard
//   - Auto-detect video standard from input signal
//
// Note: Modern digital video devices (HDMI, DisplayPort, etc.) use DV timings
// instead of video standards. See dv_timings.go for digital video interfaces.
//
// Common use cases:
//   - TV tuner cards (PAL/NTSC/SECAM)
//   - Composite/S-Video capture
//   - Analog cameras
//   - Legacy video equipment
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/standard.html
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumstd.html

// #include <linux/videodev2.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// StdId represents a video standard ID (or set of IDs).
// Multiple standards can be OR'd together to form a set.
// Type alias for StandardId already defined in video_info.go
type StdId = StandardId

// Video standard constants - Individual PAL variants
const (
	StdPAL_B  StdId = C.V4L2_STD_PAL_B  // PAL B
	StdPAL_B1 StdId = C.V4L2_STD_PAL_B1 // PAL B1
	StdPAL_G  StdId = C.V4L2_STD_PAL_G  // PAL G
	StdPAL_H  StdId = C.V4L2_STD_PAL_H  // PAL H
	StdPAL_I  StdId = C.V4L2_STD_PAL_I  // PAL I
	StdPAL_D  StdId = C.V4L2_STD_PAL_D  // PAL D
	StdPAL_D1 StdId = C.V4L2_STD_PAL_D1 // PAL D1
	StdPAL_K  StdId = C.V4L2_STD_PAL_K  // PAL K
	StdPAL_M  StdId = C.V4L2_STD_PAL_M  // PAL M (Brazil)
	StdPAL_N  StdId = C.V4L2_STD_PAL_N  // PAL N (Argentina, Paraguay, Uruguay)
	StdPAL_Nc StdId = C.V4L2_STD_PAL_Nc // PAL Nc (Argentina)
	StdPAL_60 StdId = C.V4L2_STD_PAL_60 // PAL 60
)

// Video standard constants - Individual NTSC variants
const (
	StdNTSC_M    StdId = C.V4L2_STD_NTSC_M      // NTSC M (USA, BTSC)
	StdNTSC_M_JP StdId = C.V4L2_STD_NTSC_M_JP   // NTSC M Japan (EIA-J)
	StdNTSC_443  StdId = C.V4L2_STD_NTSC_443    // NTSC 443
	StdNTSC_M_KR StdId = C.V4L2_STD_NTSC_M_KR   // NTSC M Korea (FM A2)
)

// Video standard constants - Individual SECAM variants
const (
	StdSECAM_B  StdId = C.V4L2_STD_SECAM_B  // SECAM B
	StdSECAM_D  StdId = C.V4L2_STD_SECAM_D  // SECAM D
	StdSECAM_G  StdId = C.V4L2_STD_SECAM_G  // SECAM G
	StdSECAM_H  StdId = C.V4L2_STD_SECAM_H  // SECAM H
	StdSECAM_K  StdId = C.V4L2_STD_SECAM_K  // SECAM K
	StdSECAM_K1 StdId = C.V4L2_STD_SECAM_K1 // SECAM K1
	StdSECAM_L  StdId = C.V4L2_STD_SECAM_L  // SECAM L
	StdSECAM_LC StdId = C.V4L2_STD_SECAM_LC // SECAM L'
)

// Video standard constants - ATSC digital standards
const (
	StdATSC_8_VSB  StdId = C.V4L2_STD_ATSC_8_VSB  // ATSC 8-VSB
	StdATSC_16_VSB StdId = C.V4L2_STD_ATSC_16_VSB // ATSC 16-VSB
)

// Video standard constants - Common groupings
const (
	StdPAL_BG StdId = C.V4L2_STD_PAL_BG // PAL B/G (Western Europe)
	StdPAL_DK StdId = C.V4L2_STD_PAL_DK // PAL D/K (Eastern Europe, China)
	StdPAL    StdId = C.V4L2_STD_PAL    // All PAL standards

	StdNTSC StdId = C.V4L2_STD_NTSC // All NTSC standards

	StdSECAM_DK StdId = C.V4L2_STD_SECAM_DK // SECAM D/K
	StdSECAM    StdId = C.V4L2_STD_SECAM    // All SECAM standards

	StdB  StdId = C.V4L2_STD_B  // PAL/SECAM B
	StdG  StdId = C.V4L2_STD_G  // PAL/SECAM G
	StdH  StdId = C.V4L2_STD_H  // PAL/SECAM H
	StdL  StdId = C.V4L2_STD_L  // SECAM L/L'
	StdGH StdId = C.V4L2_STD_GH // PAL/SECAM G/H
	StdDK StdId = C.V4L2_STD_DK // PAL/SECAM D/K
	StdBG StdId = C.V4L2_STD_BG // PAL/SECAM B/G
	StdMN StdId = C.V4L2_STD_MN // PAL/NTSC M/N

	Std525_60 StdId = C.V4L2_STD_525_60 // 525 lines, 60 Hz (NTSC)
	Std625_50 StdId = C.V4L2_STD_625_50 // 625 lines, 50 Hz (PAL/SECAM)

	StdATSC StdId = C.V4L2_STD_ATSC // All ATSC standards

	StdUnknown StdId = C.V4L2_STD_UNKNOWN // Unknown standard
	StdAll     StdId = C.V4L2_STD_ALL     // All standards
)

// Standard name mappings for common standards
var StdNames = map[StdId]string{
	// Individual PAL variants
	StdPAL_B:  "PAL-B",
	StdPAL_B1: "PAL-B1",
	StdPAL_G:  "PAL-G",
	StdPAL_H:  "PAL-H",
	StdPAL_I:  "PAL-I",
	StdPAL_D:  "PAL-D",
	StdPAL_D1: "PAL-D1",
	StdPAL_K:  "PAL-K",
	StdPAL_M:  "PAL-M",
	StdPAL_N:  "PAL-N",
	StdPAL_Nc: "PAL-Nc",
	StdPAL_60: "PAL-60",

	// Individual NTSC variants
	StdNTSC_M:    "NTSC-M",
	StdNTSC_M_JP: "NTSC-M-JP",
	StdNTSC_443:  "NTSC-443",
	StdNTSC_M_KR: "NTSC-M-KR",

	// Individual SECAM variants
	StdSECAM_B:  "SECAM-B",
	StdSECAM_D:  "SECAM-D",
	StdSECAM_G:  "SECAM-G",
	StdSECAM_H:  "SECAM-H",
	StdSECAM_K:  "SECAM-K",
	StdSECAM_K1: "SECAM-K1",
	StdSECAM_L:  "SECAM-L",
	StdSECAM_LC: "SECAM-L'",

	// ATSC
	StdATSC_8_VSB:  "ATSC-8-VSB",
	StdATSC_16_VSB: "ATSC-16-VSB",

	// Groupings
	StdPAL_BG:   "PAL-B/G",
	StdPAL_DK:   "PAL-D/K",
	StdPAL:      "PAL",
	StdNTSC:     "NTSC",
	StdSECAM_DK: "SECAM-D/K",
	StdSECAM:    "SECAM",
	StdATSC:     "ATSC",
	Std525_60:   "525/60",
	Std625_50:   "625/50",
	StdUnknown:  "Unknown",
	StdAll:      "All",
}

// Standard wraps v4l2_standard structure
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enumstd.html
type Standard struct {
	v4l2Standard C.struct_v4l2_standard
}

// NewStandard creates a new Standard for the given index
func NewStandard(index uint32) Standard {
	var std C.struct_v4l2_standard
	std.index = C.__u32(index)
	return Standard{v4l2Standard: std}
}

// Index returns the enumeration index (0-based)
func (s *Standard) Index() uint32 {
	return uint32(s.v4l2Standard.index)
}

// SetIndex sets the enumeration index
func (s *Standard) SetIndex(index uint32) {
	s.v4l2Standard.index = C.__u32(index)
}

// ID returns the standard identifier(s)
func (s *Standard) ID() StdId {
	return StdId(s.v4l2Standard.id)
}

// SetID sets the standard identifier
func (s *Standard) SetID(id StdId) {
	s.v4l2Standard.id = C.v4l2_std_id(id)
}

// Name returns the standard name (e.g., "PAL-B/G", "NTSC-M")
func (s *Standard) Name() string {
	return C.GoString((*C.char)(unsafe.Pointer(&s.v4l2Standard.name[0])))
}

// FramePeriod returns the frame period (inverse of frame rate)
func (s *Standard) FramePeriod() Fract {
	return Fract{
		Numerator:   uint32(s.v4l2Standard.frameperiod.numerator),
		Denominator: uint32(s.v4l2Standard.frameperiod.denominator),
	}
}

// FrameRate returns the frame rate in Hz
func (s *Standard) FrameRate() float64 {
	return s.FramePeriod().FrameRate()
}

// FrameLines returns the number of lines per frame
func (s *Standard) FrameLines() uint32 {
	return uint32(s.v4l2Standard.framelines)
}

// String returns a formatted string representation of the standard
func (s *Standard) String() string {
	name := s.Name()
	if name == "" {
		if stdName, ok := StdNames[s.ID()]; ok {
			name = stdName
		} else {
			name = fmt.Sprintf("0x%016x", s.ID())
		}
	}
	return fmt.Sprintf("%s (%.2f fps, %d lines)",
		name, s.FrameRate(), s.FrameLines())
}

// GetStandard retrieves the currently selected video standard.
// Implements VIDIOC_G_STD ioctl.
//
// Returns the standard ID (which may be a combination of multiple standards).
// Returns error if the device doesn't support video standards.
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-std.html
func GetStandard(fd uintptr) (StdId, error) {
	var stdId C.v4l2_std_id

	if err := send(fd, C.VIDIOC_G_STD, uintptr(unsafe.Pointer(&stdId))); err != nil {
		return 0, fmt.Errorf("v4l2: VIDIOC_G_STD failed: %w", err)
	}

	return StdId(stdId), nil
}

// SetStandard sets the video standard.
// Implements VIDIOC_S_STD ioctl.
//
// The standard ID may be a single standard or a set of standards.
// The driver will choose the best match if multiple standards are specified.
//
// Note: Changing the standard may also change the current format.
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-std.html
func SetStandard(fd uintptr, stdId StdId) error {
	cStdId := C.v4l2_std_id(stdId)

	if err := send(fd, C.VIDIOC_S_STD, uintptr(unsafe.Pointer(&cStdId))); err != nil {
		return fmt.Errorf("v4l2: VIDIOC_S_STD failed: %w", err)
	}

	return nil
}

// QueryStandard auto-detects the video standard from the current input signal.
// Implements VIDIOC_QUERYSTD ioctl.
//
// This ioctl senses which of the supported standards is currently being received.
// It returns a set of all detected standards. If no signal is detected, it returns
// an error (typically ENOLINK).
//
// Note: The device must support standard detection for this to work.
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-querystd.html
func QueryStandard(fd uintptr) (StdId, error) {
	var stdId C.v4l2_std_id

	if err := send(fd, C.VIDIOC_QUERYSTD, uintptr(unsafe.Pointer(&stdId))); err != nil {
		return 0, fmt.Errorf("v4l2: VIDIOC_QUERYSTD failed: %w", err)
	}

	return StdId(stdId), nil
}

// EnumStandard enumerates a video standard by index.
// Implements VIDIOC_ENUMSTD ioctl.
//
// Retrieves information about the video standard at the given index.
// Returns ErrorBadArgument when the index is out of range.
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enumstd.html
func EnumStandard(fd uintptr, index uint32) (Standard, error) {
	std := NewStandard(index)

	if err := send(fd, C.VIDIOC_ENUMSTD, uintptr(unsafe.Pointer(&std.v4l2Standard))); err != nil {
		return Standard{}, fmt.Errorf("v4l2: VIDIOC_ENUMSTD failed for index %d: %w", index, err)
	}

	return std, nil
}

// GetAllStandards enumerates all supported video standards for the device.
//
// Returns a slice of all standards supported by the device.
// Returns an empty slice if no standards are supported (e.g., digital-only devices).
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enumstd.html
func GetAllStandards(fd uintptr) ([]Standard, error) {
	var standards []Standard

	for index := uint32(0); ; index++ {
		std, err := EnumStandard(fd, index)
		if err != nil {
			// EINVAL means we've reached the end of the list
			if errors.Is(err, ErrorBadArgument) {
				break
			}
			return nil, fmt.Errorf("failed to enumerate standard at index %d: %w", index, err)
		}
		standards = append(standards, std)
	}

	return standards, nil
}

// IsStandardSupported checks if a specific standard ID is supported by the device.
//
// This function enumerates all standards and checks if any of them match
// the given standard ID (using bitwise AND since standards can be sets).
func IsStandardSupported(fd uintptr, stdId StdId) (bool, error) {
	standards, err := GetAllStandards(fd)
	if err != nil {
		return false, err
	}

	for _, std := range standards {
		if (std.ID() & stdId) != 0 {
			return true, nil
		}
	}

	return false, nil
}

// GetStandardByID finds the Standard structure matching a specific standard ID.
//
// Returns the first standard that matches (using bitwise AND).
// Returns error if no matching standard is found.
func GetStandardByID(fd uintptr, stdId StdId) (Standard, error) {
	standards, err := GetAllStandards(fd)
	if err != nil {
		return Standard{}, err
	}

	for _, std := range standards {
		if (std.ID() & stdId) != 0 {
			return std, nil
		}
	}

	return Standard{}, fmt.Errorf("standard 0x%016x not supported by device", stdId)
}
