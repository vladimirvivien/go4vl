package v4l2

//#include <linux/videodev2.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type FrameSizeType = uint32

const (
	FrameSizeTypeDiscrete   FrameSizeType = C.V4L2_FRMSIZE_TYPE_DISCRETE
	FrameSizeTypeContinuous FrameSizeType = C.V4L2_FRMSIZE_TYPE_CONTINUOUS
	FrameSizeTypeStepwise   FrameSizeType = C.V4L2_FRMSIZE_TYPE_STEPWISE
)

// FrameSize uses v4l2_frmsizeenum to get supporeted frame size for the driver based for the pixel format.
// Use FrameSizeType to determine which sizes the driver support.
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L829
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-framesizes.html
type FrameSize struct {
	FrameSizeType
	FrameSizeDiscrete
	FrameSizeStepwise
	PixelFormat FourCCType
}

// FrameSizeDiscrete (v4l2_frmsize_discrete)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L815
type FrameSizeDiscrete struct {
	Width  uint32 // width [pixel]
	Height uint32 // height [pixel]
}

// FrameSizeStepwise (v4l2_frmsize_stepwise)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L820
type FrameSizeStepwise struct {
	MinWidth   uint32 // Minimum frame width [pixel]
	MaxWidth   uint32 // Maximum frame width [pixel]
	StepWidth  uint32 // Frame width step size [pixel]
	MinHeight  uint32 // Minimum frame height [pixel]
	MaxHeight  uint32 // Maximum frame height [pixel]
	StepHeight uint32 // Frame height step size [pixel]
}

// getFrameSize retrieves the supported frame size based on the type
func getFrameSize(frmSizeEnum C.struct_v4l2_frmsizeenum) FrameSize {
	frameSize := FrameSize{FrameSizeType: FrameSizeType(frmSizeEnum._type), PixelFormat: FourCCType(frmSizeEnum.pixel_format)}
	switch frameSize.FrameSizeType {
	case FrameSizeTypeDiscrete:
		fsDiscrete := (*FrameSizeDiscrete)(unsafe.Pointer(&frmSizeEnum.anon0[0]))
		frameSize.FrameSizeDiscrete = *fsDiscrete
		frameSize.FrameSizeStepwise.MinWidth = frameSize.FrameSizeDiscrete.Width
		frameSize.FrameSizeStepwise.MinHeight = frameSize.FrameSizeDiscrete.Height
		frameSize.FrameSizeStepwise.MaxWidth = frameSize.FrameSizeDiscrete.Width
		frameSize.FrameSizeStepwise.MaxHeight = frameSize.FrameSizeDiscrete.Height
	case FrameSizeTypeStepwise, FrameSizeTypeContinuous:
		fsStepwise := (*FrameSizeStepwise)(unsafe.Pointer(&frmSizeEnum.anon0[0]))
		frameSize.FrameSizeStepwise = *fsStepwise
		frameSize.FrameSizeDiscrete.Width = frameSize.FrameSizeStepwise.MaxWidth
		frameSize.FrameSizeDiscrete.Height = frameSize.FrameSizeStepwise.MaxHeight
	}
	return frameSize
}

// GetFormatFrameSize returns a supported device frame size for a specified encoding at index
func GetFormatFrameSize(fd uintptr, index uint32, encoding FourCCType) (FrameSize, error) {
	var frmSizeEnum C.struct_v4l2_frmsizeenum
	frmSizeEnum.index = C.uint(index)
	frmSizeEnum.pixel_format = C.uint(encoding)

	if err := send(fd, C.VIDIOC_ENUM_FRAMESIZES, uintptr(unsafe.Pointer(&frmSizeEnum))); err != nil {
		return FrameSize{}, fmt.Errorf("frame size: index %d: %w", index, err)
	}
	return getFrameSize(frmSizeEnum), nil
}

// GetFormatFrameSizes returns all supported device frame sizes for a specified encoding
func GetFormatFrameSizes(fd uintptr, encoding FourCCType) (result []FrameSize, err error) {
	index := uint32(0)
	for {
		var frmSizeEnum C.struct_v4l2_frmsizeenum
		frmSizeEnum.index = C.uint(index)
		frmSizeEnum.pixel_format = C.uint(encoding)

		if err = send(fd, C.VIDIOC_ENUM_FRAMESIZES, uintptr(unsafe.Pointer(&frmSizeEnum))); err != nil {
			if errors.Is(err, ErrorBadArgument) && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("frame sizes: encoding %s: %w", PixelFormats[encoding], err)
		}

		// At index 0, check the frame type, if not discrete exit loop.
		// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-framesizes.html
		result = append(result, getFrameSize(frmSizeEnum))
		if index == 0 && uint32(frmSizeEnum._type) != FrameSizeTypeDiscrete {
			break
		}
		index++
	}
	return result, nil
}

// GetAllFormatFrameSizes returns all supported frame sizes for all supported formats.
// It iterates from format at index 0 until it encounters and error and then stops. For
// each supported format, it retrieves all supported frame sizes.
func GetAllFormatFrameSizes(fd uintptr) (result []FrameSize, err error) {
	formats, err := GetAllFormatDescriptions(fd)
	if len(formats) == 0 && err != nil {
		return nil, fmt.Errorf("frame sizes: %w", err)
	}

	// for each supported format, grab frame size
	for _, format := range formats {
		index := uint32(0)
		for {
			var frmSizeEnum C.struct_v4l2_frmsizeenum
			frmSizeEnum.index = C.uint(index)
			frmSizeEnum.pixel_format = C.uint(format.PixelFormat)

			if err = send(fd, C.VIDIOC_ENUM_FRAMESIZES, uintptr(unsafe.Pointer(&frmSizeEnum))); err != nil {
				if errors.Is(err, ErrorBadArgument) && len(result) > 0 {
					break
				}
				return result, err
			}

			// At index 0, check the frame type, if not discrete exit loop.
			// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-framesizes.html
			result = append(result, getFrameSize(frmSizeEnum))
			if index == 0 && uint32(frmSizeEnum._type) != FrameSizeTypeDiscrete {
				break
			}
			index++
		}
	}
	return result, nil
}
