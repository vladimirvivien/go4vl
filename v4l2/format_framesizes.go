package v4l2

import (
	"errors"
	"fmt"
	"unsafe"
)

//enum v4l2_frmsizetypes {
//V4L2_FRMSIZE_TYPE_DISCRETE	= 1,
//V4L2_FRMSIZE_TYPE_CONTINUOUS	= 2,
//V4L2_FRMSIZE_TYPE_STEPWISE	= 3,
//};
//
//struct v4l2_frmsize_discrete {
//__u32			width;		/* Frame width [pixel] */
//__u32			height;		/* Frame height [pixel] */
//};
//
//struct v4l2_frmsize_stepwise {
//__u32			min_width;	/* Minimum frame width [pixel] */
//__u32			max_width;	/* Maximum frame width [pixel] */
//__u32			step_width;	/* Frame width step size [pixel] */
//__u32			min_height;	/* Minimum frame height [pixel] */
//__u32			max_height;	/* Maximum frame height [pixel] */
//__u32			step_height;	/* Frame height step size [pixel] */
//};
//
//struct v4l2_frmsizeenum {
//__u32			index;		/* Frame size number */
//__u32			pixel_format;	/* Pixel format */
//__u32			type;		/* Frame size type the device supports. */
//
//union {					/* Frame size */
//struct v4l2_frmsize_discrete	discrete;
//struct v4l2_frmsize_stepwise	stepwise;
//};
//
//__u32   reserved[2];			/* Reserved space for future use */
//};

type FrameSizeType = uint32

const (
	FrameSizeTypeDiscrete   FrameSizeType = iota + 1 // V4L2_FRMSIZE_TYPE_DISCRETE	 = 1
	FrameSizeTypeContinuous                          // V4L2_FRMSIZE_TYPE_CONTINUOUS = 2
	FrameSizeTypeStepwise                            // V4L2_FRMSIZE_TYPE_STEPWISE	 = 3
)

// FrameSize is the frame size supported by the driver for the pixel format.
// Use FrameSizeType to determine which sizes the driver support.
type FrameSize struct {
	FrameSizeType
	FrameSizeDiscrete
	FrameSizeStepwise
	PixelFormat FourCCEncoding
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

// v4l2FrameSizes (v4l2_frmsizeenum)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L829
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-framesizes.html
type v4l2FrameSizeEnum struct {
	index         uint32
	pixelFormat   FourCCEncoding
	frameSizeType FrameSizeType
	frameSize     [24]byte // union sized for largest struct: stepwise
	reserved      [2]uint32
}

// getFrameSize retrieves the supported frame size
func (fs v4l2FrameSizeEnum) getFrameSize() FrameSize {
	frameSize := FrameSize{FrameSizeType:fs.frameSizeType, PixelFormat: fs.pixelFormat}
	switch fs.frameSizeType {
	case FrameSizeTypeDiscrete:
		fsDiscrete := (*FrameSizeDiscrete)(unsafe.Pointer(&fs.frameSize[0]))
		frameSize.FrameSizeDiscrete = *fsDiscrete
	case FrameSizeTypeStepwise, FrameSizeTypeContinuous:
		fsStepwise := (*FrameSizeStepwise)(unsafe.Pointer(&fs.frameSize[0]))
		frameSize.FrameSizeStepwise = *fsStepwise
	}
	return frameSize
}

// GetFormatFrameSize returns a supported device frame size for a specified encoding at index
func GetFormatFrameSize(fd uintptr, index uint32, encoding FourCCEncoding) (FrameSize, error) {
	fsEnum := v4l2FrameSizeEnum{index: index, pixelFormat: encoding}
	if err := Send(fd, VidiocEnumFrameSizes, uintptr(unsafe.Pointer(&fsEnum))); err != nil {
		switch {
		case errors.Is(err, ErrorUnsupported):
			return FrameSize{}, fmt.Errorf("frame size: index %d: not found %w", index, err)
		default:
			return FrameSize{}, fmt.Errorf("frame size: %w", err)
		}
	}
	return fsEnum.getFrameSize(), nil
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
	for _, fmt := range formats {
		index := uint32(0)
		for {
			fsEnum := v4l2FrameSizeEnum{index: index, pixelFormat: fmt.GetPixelFormat()}
			if err = Send(fd, VidiocEnumFrameSizes, uintptr(unsafe.Pointer(&fsEnum))); err != nil {
				break
			}

			// At index 0, check the frame type, if not discrete exit loop.
			// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-framesizes.html
			result = append(result, fsEnum.getFrameSize())
			if index == 0 && fsEnum.frameSizeType != FrameSizeTypeDiscrete{
				break
			}

			index++
		}
	}
	return result, err
}