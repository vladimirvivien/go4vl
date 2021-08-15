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
	FmtDescFlagCompressed           FmtDescFlag = 0x0001                 // V4L2_FMT_FLAG_COMPRESSED
	FmtDescFlagEmulated             FmtDescFlag = 0x0002                 // V4L2_FMT_FLAG_EMULATED
	FmtDescFlagContinuousBytestream FmtDescFlag = 0x0004                 // V4L2_FMT_FLAG_CONTINUOUS_BYTESTREAM
	FmtDescFlagDynResolution        FmtDescFlag = 0x0008                 // V4L2_FMT_FLAG_DYN_RESOLUTION
	FmtDescFlagEncCapFrameInterval  FmtDescFlag = 0x0010                 //  V4L2_FMT_FLAG_ENC_CAP_FRAME_INTERVAL
	FmtDescFlagCscColorspace        FmtDescFlag = 0x0020                 //  V4L2_FMT_FLAG_CSC_COLORSPACE
	FmtDescFlagCscXferFunc          FmtDescFlag = 0x0040                 // V4L2_FMT_FLAG_CSC_XFER_FUNC
	FmtDescFlagCscYcbcrEnc          FmtDescFlag = 0x0080                 //  V4L2_FMT_FLAG_CSC_YCBCR_ENC
	FmtDescFlagCscHsvEnc            FmtDescFlag = FmtDescFlagCscYcbcrEnc // V4L2_FMT_FLAG_CSC_HSV_ENC
	FmtDescFlagCscQuantization      FmtDescFlag = 0x0100                 // V4L2_FMT_FLAG_CSC_QUANTIZATION
)

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

// GetFrameSize return the supported frame size for the format in description.
// NOTE: This method must be used on a FormatDescription value that was created
// with a call to GetFormatDescription or GetAllFormatDescriptions.
func (d FormatDescription) GetFrameSize() (FrameSize, error) {
	if d.fd == 0{
		return FrameSize{}, fmt.Errorf("invalid file descriptor")
	}
	return GetFormatFrameSize(d.fd, d.index, d.pixelFormat)
}

// GetFormatDescription returns a device format description at index
func GetFormatDescription(fd uintptr, index uint32) (FormatDescription, error) {
	desc := v4l2FormatDesc{index: index, bufType: BufTypeVideoCapture}
	if err := Send(fd, VidiocEnumFmt, uintptr(unsafe.Pointer(&desc))); err != nil {
		switch {
		case errors.Is(err, ErrorUnsupported):
			return FormatDescription{}, fmt.Errorf("format desc: index %d: not found %w", index, err)
		default:
			return FormatDescription{}, fmt.Errorf("format desc failed: %w", err)
		}
	}
	return FormatDescription{fd:fd, v4l2FormatDesc:desc}, nil
}

// GetAllFormatDescriptions attempts to retrieve all device format descriptions by
// iterating from 0 upto an index that returns an error. At that point, the function
// will return the collected descriptions and the error.
// So if len(result) > 0, then error could be ignored.
func GetAllFormatDescriptions(fd uintptr) (result []FormatDescription, err error){
	index := uint32(0)
	for {
		desc := v4l2FormatDesc{index: index, bufType: BufTypeVideoCapture}
		if err = Send(fd, VidiocEnumFmt, uintptr(unsafe.Pointer(&desc))); err != nil {
			break
		}
		result = append(result, FormatDescription{fd: fd, v4l2FormatDesc:desc})
		index++
	}
	return result, err
}