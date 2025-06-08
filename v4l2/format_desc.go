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

// FmtDescFlag is a type alias for uint32, representing flags that provide additional
// information about a V4L2 format description. These flags are part of the FormatDescription struct.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-fmt.html#fmtdesc-flags
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L794
type FmtDescFlag = uint32

// Format Description Flag Constants
const (
	// FmtDescFlagCompressed indicates that the format is compressed.
	FmtDescFlagCompressed FmtDescFlag = C.V4L2_FMT_FLAG_COMPRESSED
	// FmtDescFlagEmulated indicates that the format is emulated by the driver.
	FmtDescFlagEmulated FmtDescFlag = C.V4L2_FMT_FLAG_EMULATED
	// FmtDescFlagContinuousBytestream indicates that the format is a continuous bytestream, not distinct frames.
	FmtDescFlagContinuousBytestream FmtDescFlag = C.V4L2_FMT_FLAG_CONTINUOUS_BYTESTREAM
	// FmtDescFlagDynResolution indicates that the format supports dynamic resolution changes.
	FmtDescFlagDynResolution FmtDescFlag = C.V4L2_FMT_FLAG_DYN_RESOLUTION
	// FmtDescFlagEncodedCaptureFrameInterval indicates that the capture device can vary the frame interval for encoded formats.
	FmtDescFlagEncodedCaptureFrameInterval FmtDescFlag = C.V4L2_FMT_FLAG_ENC_CAP_FRAME_INTERVAL
	// FmtDescFlagConfigColorspace indicates that the colorspace can be configured via VIDIOC_S_EXT_CTRLS.
	FmtDescFlagConfigColorspace FmtDescFlag = C.V4L2_FMT_FLAG_CSC_COLORSPACE
	// FmtDescFlagConfigXferFunc indicates that the transfer function can be configured via VIDIOC_S_EXT_CTRLS.
	FmtDescFlagConfigXferFunc FmtDescFlag = C.V4L2_FMT_FLAG_CSC_XFER_FUNC
	// FmtDescFlagConfigYcbcrEnc indicates that YCbCr encoding can be configured via VIDIOC_S_EXT_CTRLS.
	FmtDescFlagConfigYcbcrEnc FmtDescFlag = C.V4L2_FMT_FLAG_CSC_YCBCR_ENC
	// FmtDescFlagConfigHsvEnc indicates that HSV encoding can be configured via VIDIOC_S_EXT_CTRLS.
	FmtDescFlagConfigHsvEnc FmtDescFlag = C.V4L2_FMT_FLAG_CSC_HSV_ENC // Note: Kernel uses YCBCR for HSV as well.
	// FmtDescFlagConfigQuantization indicates that quantization can be configured via VIDIOC_S_EXT_CTRLS.
	FmtDescFlagConfigQuantization FmtDescFlag = C.V4L2_FMT_FLAG_CSC_QUANTIZATION
)

// FormatDescriptionFlags provides a map of FmtDescFlag constants to their human-readable string descriptions.
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

// FormatDescription describes a V4L2 data format supported by a device.
// It corresponds to the `v4l2_fmtdesc` struct in the Linux kernel.
// This structure is used with the VIDIOC_ENUM_FMT ioctl to enumerate available formats.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-fmt.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L784
type FormatDescription struct {
	// Index is the zero-based index of the format in the enumeration. This is set by the application when calling VIDIOC_ENUM_FMT.
	Index uint32
	// StreamType is the type of data stream (e.g., video capture, video output). This is set by the application. See BufType constants.
	StreamType BufType
	// Flags provide additional information about the format. See FmtDescFlag constants.
	Flags FmtDescFlag
	// Description is a human-readable string describing the format (e.g., "YUYV 4:2:2").
	Description string
	// PixelFormat is the FourCC code identifying the pixel format (e.g., V4L2_PIX_FMT_YUYV).
	PixelFormat FourCCType
	// MBusCode is the media bus code, relevant for devices using the V4L2 subdev API and media controller.
	MBusCode uint32
	// reserved space in C struct
	// _ [4]uint32 // Implicitly handled by CGo struct mapping if present in C.struct_v4l2_fmtdesc
}

// String returns a human-readable string representation of the FormatDescription.
// It includes the format description, index, flags, and pixel format FourCC.
func (d FormatDescription) String() string {
	return fmt.Sprintf(
		"Format: %s [index: %d, flags: %s, format:%s]",
		d.Description,
		d.Index,
		FormatDescriptionFlags[d.Flags],
		PixelFormats[d.PixelFormat], // Assumes PixelFormats map is available and contains the FourCC string
	)
}

// makeFormatDescription is an internal helper function to convert a C.struct_v4l2_fmtdesc
// to a Go FormatDescription struct.
func makeFormatDescription(fmtDesc C.struct_v4l2_fmtdesc) FormatDescription {
	return FormatDescription{
		Index:       uint32(fmtDesc.index),
		StreamType:  BufType(fmtDesc._type), // Cast to BufType
		Flags:       FmtDescFlag(fmtDesc.flags), // Cast to FmtDescFlag
		Description: C.GoString((*C.char)(unsafe.Pointer(&fmtDesc.description[0]))),
		PixelFormat: FourCCType(fmtDesc.pixelformat), // Cast to FourCCType
		MBusCode:    uint32(fmtDesc.mbus_code),
	}
}

// GetFormatDescription retrieves a specific format description by its index for a given buffer type.
// It takes the file descriptor of the V4L2 device and the zero-based index of the format.
// It typically defaults to querying for BufTypeVideoCapture.
// Returns a FormatDescription struct and an error if the VIDIOC_ENUM_FMT ioctl call fails (e.g., index out of bounds).
func GetFormatDescription(fd uintptr, index uint32) (FormatDescription, error) {
	var fmtDesc C.struct_v4l2_fmtdesc
	fmtDesc.index = C.uint(index)
	fmtDesc._type = C.uint(BufTypeVideoCapture) // Defaulting to video capture type.

	if err := send(fd, C.VIDIOC_ENUM_FMT, uintptr(unsafe.Pointer(&fmtDesc))); err != nil {
		return FormatDescription{}, fmt.Errorf("format desc: index %d: %w", index, err)
	}
	return makeFormatDescription(fmtDesc), nil
}

// GetAllFormatDescriptions retrieves all available format descriptions for the device's video capture stream.
// It iterates by calling GetFormatDescription with increasing indices until an error (typically ErrorBadArgument,
// indicating no more formats) is encountered.
// It returns a slice of FormatDescription structs and any error encountered during the final failing call.
// If some formats were successfully retrieved before an error, those will be returned along with the error.
// If the first call fails, it returns an empty slice and the error.
func GetAllFormatDescriptions(fd uintptr) (result []FormatDescription, err error) {
	index := uint32(0)
	for {
		var fmtDesc C.struct_v4l2_fmtdesc
		fmtDesc.index = C.uint(index)
		fmtDesc._type = C.uint(BufTypeVideoCapture) // Defaulting to video capture type.

		err = send(fd, C.VIDIOC_ENUM_FMT, uintptr(unsafe.Pointer(&fmtDesc)))
		if err != nil {
			// If ErrorBadArgument is returned, it means we've enumerated all formats.
			// If result has items, this is not a "true" error for the collection process.
			if errors.Is(err, ErrorBadArgument) && len(result) > 0 {
				err = nil // Clear error as we successfully enumerated some formats.
				break
			}
			// For other errors, or if ErrorBadArgument on the first try, return the error.
			return result, fmt.Errorf("format desc: error on index %d: %w", index, err)
		}
		result = append(result, makeFormatDescription(fmtDesc))
		index++
	}
	return result, err
}

// GetFormatDescriptionByEncoding searches through all available format descriptions for one
// that matches the specified FourCCType (pixel format encoding).
// It takes the file descriptor and the desired FourCCType.
// Returns the matching FormatDescription and nil error if found.
// If no matching format is found, or if an error occurs while fetching descriptions,
// it returns an empty FormatDescription and an error.
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
