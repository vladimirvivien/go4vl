package v4l2

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
	FmtDescFlagCompressed                  FmtDescFlag = 0x0001                    // V4L2_FMT_FLAG_COMPRESSED
	FmtDescFlagEmulated                    FmtDescFlag = 0x0002                    // V4L2_FMT_FLAG_EMULATED
	FmtDescFlagContinuousBytestream        FmtDescFlag = 0x0004                    // V4L2_FMT_FLAG_CONTINUOUS_BYTESTREAM
	FmtDescFlagDynResolution               FmtDescFlag = 0x0008                    // V4L2_FMT_FLAG_DYN_RESOLUTION
	FmtDescFlagEncodedCaptureFrameInterval FmtDescFlag = 0x0010                    //  V4L2_FMT_FLAG_ENC_CAP_FRAME_INTERVAL
	FmtDescFlagConfigColorspace            FmtDescFlag = 0x0020                    //  V4L2_FMT_FLAG_CSC_COLORSPACE
	FmtDescFlagConfigXferFunc              FmtDescFlag = 0x0040                    // V4L2_FMT_FLAG_CSC_XFER_FUNC
	FmtDescFlagConfigYcbcrEnc              FmtDescFlag = 0x0080                    //  V4L2_FMT_FLAG_CSC_YCBCR_ENC
	FmtDescFlagConfigHsvEnc                FmtDescFlag = FmtDescFlagConfigYcbcrEnc // V4L2_FMT_FLAG_CSC_HSV_ENC
	FmtDescFlagConfigQuantization          FmtDescFlag = 0x0100                    // V4L2_FMT_FLAG_CSC_QUANTIZATION
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

// v4l2FormatDesc  (v4l2_fmtdesc)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L784
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enum-fmt.html
type v4l2FormatDesc struct {
	index       uint32  // format number
	bufType     BufType // stream type BufType
	flags       FmtDescFlag
	description [32]uint8      // string description
	pixelFormat FourCCEncoding // Format fourcc value
	mbusCode    uint32         // media bus code
	reserved    [3]uint32
}

// FormatDescription provides access to the device format description
// See v4l2FormatDesc
type FormatDescription struct {
	fd uintptr
	v4l2FormatDesc
}

// GetIndex returns the format number
func (d FormatDescription) GetIndex() uint32 {
	return d.index
}

// GetBufType returns the type for the buffer (see v4l2_buf_type)
func (d FormatDescription) GetBufType() BufType {
	return d.bufType
}

// GetFlags returns image description flags (see FmtDescFlag)
func (d FormatDescription) GetFlags() FmtDescFlag {
	return d.flags
}

// GetDescription returns a string value for the format description
func (d FormatDescription) GetDescription() string {
	return toGoString(d.description[:])
}

// GetPixelFormat returns the four character encoding for the format
func (d FormatDescription) GetPixelFormat() FourCCEncoding {
	return d.pixelFormat
}

// GetBusCode returns the media bus code for drivers that advertise v4l2_cap_io_mc
func (d FormatDescription) GetBusCode() uint32 {
	return d.mbusCode
}

// GetFrameSizes return all supported frame sizes for the format description.
func (d FormatDescription) GetFrameSizes() ([]FrameSize, error) {
	if d.fd == 0 {
		return nil, fmt.Errorf("invalid file descriptor")
	}
	return GetFormatFrameSizes(d.fd, d.pixelFormat)
}

// GetFormatDescription returns a device format description at index
func GetFormatDescription(fd uintptr, index uint32) (FormatDescription, error) {
	desc := v4l2FormatDesc{index: index, bufType: BufTypeVideoCapture}
	if err := Send(fd, VidiocEnumFmt, uintptr(unsafe.Pointer(&desc))); err != nil {
		return FormatDescription{}, fmt.Errorf("format desc: index %d: %w", index, err)

	}
	return FormatDescription{fd: fd, v4l2FormatDesc: desc}, nil
}

// GetAllFormatDescriptions attempts to retrieve all device format descriptions by
// iterating from 0 up to an index that returns an error. At that point, the function
// will return the collected descriptions and the error.
// So if len(result) > 0, then error could be ignored.
func GetAllFormatDescriptions(fd uintptr) (result []FormatDescription, err error) {
	index := uint32(0)
	for {
		desc := v4l2FormatDesc{index: index, bufType: BufTypeVideoCapture}
		if err = Send(fd, VidiocEnumFmt, uintptr(unsafe.Pointer(&desc))); err != nil {
			if errors.Is(err, ErrorBadArgument) && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("format desc: all: %w", err)
		}
		result = append(result, FormatDescription{fd: fd, v4l2FormatDesc: desc})
		index++
	}
	return result, nil
}

// GetFormatDescriptionByEncoding returns a FormatDescription that matches the specified encoded pixel format
func GetFormatDescriptionByEncoding(fd uintptr, enc FourCCEncoding)(FormatDescription, error) {
	descs, err := GetAllFormatDescriptions(fd)
	if err != nil {
		return FormatDescription{}, fmt.Errorf("format desc: encoding %s: %s", PixelFormats[enc], err)
	}
	for _, desc := range descs {
		if desc.GetPixelFormat() == enc{
			return desc, nil
		}
	}

	return FormatDescription{}, fmt.Errorf("format desc: driver does not support encoding %d", enc)
}