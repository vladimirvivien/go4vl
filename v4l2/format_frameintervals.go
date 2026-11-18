package v4l2

import (
	"fmt"
	"unsafe"
)

// FrameIntervalType (v4l2_frmivaltypes)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L845
type FrameIntervalType = uint32

const (
	FrameIntervalTypeDiscrete   FrameIntervalType = cV4L2_FRMIVAL_TYPE_DISCRETE
	FrameIntervalTypeContinuous FrameIntervalType = cV4L2_FRMIVAL_TYPE_CONTINUOUS
	FrameIntervalTypeStepwise   FrameIntervalType = cV4L2_FRMIVAL_TYPE_STEPWISE
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
func getFrameInterval(interval cV4L2FrmIvalEnum) (FrameIntervalEnum, error) {
	frmInterval := FrameIntervalEnum{
		Index:       interval.Index,
		Type:        FrameIntervalType(interval.Type),
		PixelFormat: FourCCType(interval.Format),
		Width:       interval.Width,
		Height:      interval.Height,
	}
	intervalType := interval.Type
	switch intervalType {
	case FrameIntervalTypeDiscrete:
		fiDiscrete := *(*Fract)(unsafe.Pointer(&interval.Discrete))
		frmInterval.Interval.Min = fiDiscrete
		frmInterval.Interval.Max = fiDiscrete
		frmInterval.Interval.Step.Numerator = 1
		frmInterval.Interval.Step.Denominator = 1
	case FrameIntervalTypeStepwise, FrameIntervalTypeContinuous:
		// Stepwise union member starts at offset 0 of the discrete/stepwise union.
		frmInterval.Interval = *(*FrameInterval)(unsafe.Pointer(&interval.Discrete))
	default:
		return FrameIntervalEnum{}, fmt.Errorf("unsupported frame interval type: %d", intervalType)
	}
	return frmInterval, nil
}

// GetFormatFrameInterval returns a supported device frame interval for a specified encoding at index and format
func GetFormatFrameInterval(fd uintptr, index uint32, encoding FourCCType, width, height uint32) (FrameIntervalEnum, error) {
	var interval cV4L2FrmIvalEnum
	interval.Index = index
	interval.Format = encoding
	interval.Width = width
	interval.Height = height

	if err := send(fd, cVIDIOC_ENUM_FRAMEINTERVALS, uintptr(unsafe.Pointer(&interval))); err != nil {
		return FrameIntervalEnum{}, fmt.Errorf("frame interval: index %d: %w", index, err)
	}
	return getFrameInterval(interval)
}
