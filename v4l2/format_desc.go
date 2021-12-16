package v4l2

/*
#include <linux/videodev2.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// FmtDescFlag image format description flags
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L794
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-fmt.html#fmtdesc-flags
type FmtDescFlag = uint32

const (
	FmtDescFlagCompressed                  FmtDescFlag = C.V4L2_FMT_FLAG_COMPRESSED
	FmtDescFlagEmulated                    FmtDescFlag = C.V4L2_FMT_FLAG_EMULATED
	FmtDescFlagContinuousBytestream        FmtDescFlag = C.V4L2_FMT_FLAG_CONTINUOUS_BYTESTREAM
	FmtDescFlagDynResolution               FmtDescFlag = C.V4L2_FMT_FLAG_DYN_RESOLUTION
	FmtDescFlagEncodedCaptureFrameInterval FmtDescFlag = C.V4L2_FMT_FLAG_ENC_CAP_FRAME_INTERVAL
	FmtDescFlagConfigColorspace            FmtDescFlag = C.V4L2_FMT_FLAG_CSC_COLORSPACE
	FmtDescFlagConfigXferFunc              FmtDescFlag = C.V4L2_FMT_FLAG_CSC_XFER_FUNC
	FmtDescFlagConfigYcbcrEnc              FmtDescFlag = C.V4L2_FMT_FLAG_CSC_YCBCR_ENC
	FmtDescFlagConfigHsvEnc                FmtDescFlag = C.V4L2_FMT_FLAG_CSC_HSV_ENC
	FmtDescFlagConfigQuantization          FmtDescFlag = C.V4L2_FMT_FLAG_CSC_QUANTIZATION
)

var FormatDescriptionFlags = map[FmtDescFlag]string{
	FmtDescFlagCompressed:                  "Compressed",
	FmtDescFlagEmulated:                    "Emulated",
	FmtDescFlagContinuousBytestream:        "Continuous bytestream",
	FmtDescFlagDynResolution:               "Dynamic resolution",
	FmtDescFlagEncodedCaptureFrameInterval: "Encoded capture frame interval",
	FmtDescFlagConfigColorspace:            "Colorspace update supported",
	FmtDescFlagConfigXferFunc:              "Transfer func update supported",
	FmtDescFlagConfigYcbcrEnc:              "YCbCr/HSV update supported",
	FmtDescFlagConfigQuantization:          "Quantization update supported",
}

// FormatDescription  (v4l2_fmtdesc) provides access to the device format description
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L784
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-fmt.html
type FormatDescription struct {
	// Index returns the format number
	Index uint32
	// StreamType type for the buffer (see v4l2_buf_type)
	StreamType BufType
	// Flags is the image description flags (see FmtDescFlag)
	Flags FmtDescFlag
	// Description is a string value for the format description
	Description string
	// PixelFormat stores the four character encoding for the format
	PixelFormat FourCCType
	// MBusCode is the media bus code for drivers that advertise v4l2_cap_io_mc
	MBusCode uint32
}

func (d FormatDescription) String() string {
	return fmt.Sprintf(
		"Format: %s [index: %d, flags: %s, format:%s]",
		d.Description,
		d.Index,
		FormatDescriptionFlags[d.Flags],
		PixelFormats[d.PixelFormat],
	)
}
func makeFormatDescription(fmtDesc C.struct_v4l2_fmtdesc) FormatDescription {
	return FormatDescription{
		Index:       uint32(fmtDesc.index),
		StreamType:  uint32(fmtDesc._type),
		Flags:       uint32(fmtDesc.flags),
		Description: C.GoString((*C.char)(unsafe.Pointer(&fmtDesc.description[0]))),
		PixelFormat: uint32(fmtDesc.pixelformat),
		MBusCode:    uint32(fmtDesc.mbus_code),
	}
}

// GetFormatDescription returns a device format description at index
func GetFormatDescription(fd uintptr, index uint32) (FormatDescription, error) {
	var fmtDesc C.struct_v4l2_fmtdesc
	fmtDesc.index = C.uint(index)
	fmtDesc._type = C.uint(BufTypeVideoCapture)

	if err := send(fd, C.VIDIOC_ENUM_FMT, uintptr(unsafe.Pointer(&fmtDesc))); err != nil {
		return FormatDescription{}, fmt.Errorf("format desc: index %d: %w", index, err)

	}
	return makeFormatDescription(fmtDesc), nil
}

// GetAllFormatDescriptions attempts to retrieve all device format descriptions by
// iterating from 0 up to an index that returns an error. At that point, the function
// will return the collected descriptions and the error.
// So if len(result) > 0, then error could be ignored.
func GetAllFormatDescriptions(fd uintptr) (result []FormatDescription, err error) {
	index := uint32(0)
	for {
		var fmtDesc C.struct_v4l2_fmtdesc
		fmtDesc.index = C.uint(index)
		fmtDesc._type = C.uint(BufTypeVideoCapture)

		if err = send(fd, C.VIDIOC_ENUM_FMT, uintptr(unsafe.Pointer(&fmtDesc))); err != nil {
			if errors.Is(err, ErrorBadArgument) && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("format desc: all: %w", err)
		}
		result = append(result, makeFormatDescription(fmtDesc))
		index++
	}
	return result, nil
}

// GetFormatDescriptionByEncoding returns a FormatDescription that matches the specified encoded pixel format
func GetFormatDescriptionByEncoding(fd uintptr, enc FourCCType) (FormatDescription, error) {
	descs, err := GetAllFormatDescriptions(fd)
	if err != nil {
		return FormatDescription{}, fmt.Errorf("format desc: encoding %s: %s", PixelFormats[enc], err)
	}
	for _, desc := range descs {
		if desc.PixelFormat == enc {
			return desc, nil
		}
	}

	return FormatDescription{}, fmt.Errorf("format desc: driver does not support encoding %d", enc)
}
