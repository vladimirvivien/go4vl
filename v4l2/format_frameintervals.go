package v4l2

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../include/
#include <linux/videodev2.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// FrameIntervalType (v4l2_frmivaltypes)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L845
type FrameIntervalType = uint32

const (
	FrameIntervalTypeDiscrete   FrameIntervalType = C.V4L2_FRMIVAL_TYPE_DISCRETE
	FrameIntervalTypeContinuous FrameIntervalType = C.V4L2_FRMIVAL_TYPE_CONTINUOUS
	FrameIntervalTypeStepwise   FrameIntervalType = C.V4L2_FRMIVAL_TYPE_STEPWISE
)

// FrameIntervalEnum is used to store v4l2_frmivalenum values.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L857
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-frameintervals.html
type FrameIntervalEnum struct {
	Index       uint32
	PixelFormat FourCCType
	Width       uint32
	Height      uint32
	Type        FrameIntervalType
	Interval    FrameInterval
}

// FrameInterval stores all frame interval values regardless of its type. This type maps to v4l2_frmival_stepwise.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L851
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-frameintervals.html
type FrameInterval struct {
	Min  Fract
	Max  Fract
	Step Fract
}

// getFrameInterval retrieves the supported frame interval info from following union based on the type:

// 	union {
//	    struct v4l2_fract		discrete;
//	    struct v4l2_frmival_stepwise	stepwise;
//	}

// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-frameintervals.html
func getFrameInterval(interval C.struct_v4l2_frmivalenum) (FrameIntervalEnum, error) {
	frmInterval := FrameIntervalEnum{
		Index:       uint32(interval.index),
		Type:        FrameIntervalType(interval._type),
		PixelFormat: FourCCType(interval.pixel_format),
		Width:       uint32(interval.width),
		Height:      uint32(interval.height),
	}
	intervalType := uint32(interval._type)
	switch intervalType {
	case FrameIntervalTypeDiscrete:
		fiDiscrete := *(*Fract)(unsafe.Pointer(&interval.anon0[0]))
		frmInterval.Interval.Min = fiDiscrete
		frmInterval.Interval.Max = fiDiscrete
		frmInterval.Interval.Step.Numerator = 1
		frmInterval.Interval.Step.Denominator = 1
	case FrameIntervalTypeStepwise, FrameIntervalTypeContinuous:
		// Calculate pointer to stepwise member of union
		frmInterval.Interval = *(*FrameInterval)(unsafe.Pointer(uintptr(unsafe.Pointer(&interval.anon0[0])) + unsafe.Sizeof(Fract{})))
	default:
		return FrameIntervalEnum{}, fmt.Errorf("unsupported frame interval type: %d", intervalType)
	}
	return frmInterval, nil
}

// GetFormatFrameInterval returns a supported device frame interval for a specified encoding at index and format
func GetFormatFrameInterval(fd uintptr, index uint32, encoding FourCCType, width, height uint32) (FrameIntervalEnum, error) {
	var interval C.struct_v4l2_frmivalenum
	interval.index = C.uint(index)
	interval.pixel_format = C.uint(encoding)
	interval.width = C.uint(width)
	interval.height = C.uint(height)

	if err := send(fd, C.VIDIOC_ENUM_FRAMEINTERVALS, uintptr(unsafe.Pointer(&interval))); err != nil {
		return FrameIntervalEnum{}, fmt.Errorf("frame interval: index %d: %w", index, err)
	}
	return getFrameInterval(interval)
}
