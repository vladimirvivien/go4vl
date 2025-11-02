package v4l2

// video_info.go provides video input and output enumeration and selection.
//
// Video inputs and outputs represent physical or logical video connections on a device.
// For example, a capture card might have multiple inputs (HDMI, composite, S-Video),
// and a video output device might have multiple outputs (HDMI, DisplayPort).
//
// The V4L2 API provides the following operations:
//   - Enumerate available inputs/outputs
//   - Query input/output capabilities and status
//   - Select active input/output
//   - Query audio/video standards association
//
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enuminput.html
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumoutput.html

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// InputStatus
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enuminput.html?highlight=v4l2_input#input-status
type InputStatus = uint32

var (
	InputStatusNoPower  InputStatus = C.V4L2_IN_ST_NO_POWER
	InputStatusNoSignal InputStatus = C.V4L2_IN_ST_NO_SIGNAL
	InputStatusNoColor  InputStatus = C.V4L2_IN_ST_NO_COLOR
)

var InputStatuses = map[InputStatus]string{
	0:                   "ok",
	InputStatusNoPower:  "no power",
	InputStatusNoSignal: "no signal",
	InputStatusNoColor:  "no color",
}

type InputType = uint32

const (
	InputTypeTuner InputType = iota + 1
	InputTypeCamera
	InputTypeTouch
)

type StandardId = uint64

// OutputType represents the type of video output
type OutputType = uint32

const (
	OutputTypeModulator OutputType = iota + 1
	OutputTypeAnalog
	OutputTypeAnalogVGAOverlay
)

// OutputStatus represents the status of a video output
// Note: The V4L2 API does not define output status flags in the same way as input status.
// Most output devices will return 0 (OK) status.
type OutputStatus = uint32

var OutputStatuses = map[OutputStatus]string{
	0: "ok",
}

// InputInfo (v4l2_input)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1649
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enuminput.html
type InputInfo struct {
	v4l2Input C.struct_v4l2_input
}

func (i InputInfo) GetIndex() uint32 {
	return uint32(i.v4l2Input.index)
}

func (i InputInfo) GetName() string {
	return C.GoString((*C.char)(unsafe.Pointer(&i.v4l2Input.name[0])))
}

func (i InputInfo) GetInputType() InputType {
	return InputType(i.v4l2Input._type)
}

func (i InputInfo) GetAudioset() uint32 {
	return uint32(i.v4l2Input.audioset)
}

func (i InputInfo) GetTuner() uint32 {
	return uint32(i.v4l2Input.tuner)
}

func (i InputInfo) GetStandardId() StandardId {
	return StandardId(i.v4l2Input.std)
}

func (i InputInfo) GetStatus() uint32 {
	return uint32(i.v4l2Input.status)
}

func (i InputInfo) GetCapabilities() uint32 {
	return uint32(i.v4l2Input.capabilities)
}

// OutputInfo (v4l2_output)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1746
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumoutput.html
type OutputInfo struct {
	v4l2Output C.struct_v4l2_output
}

func (o OutputInfo) GetIndex() uint32 {
	return uint32(o.v4l2Output.index)
}

func (o OutputInfo) GetName() string {
	return C.GoString((*C.char)(unsafe.Pointer(&o.v4l2Output.name[0])))
}

func (o OutputInfo) GetOutputType() OutputType {
	return OutputType(o.v4l2Output._type)
}

func (o OutputInfo) GetAudioset() uint32 {
	return uint32(o.v4l2Output.audioset)
}

func (o OutputInfo) GetModulator() uint32 {
	return uint32(o.v4l2Output.modulator)
}

func (o OutputInfo) GetStandardId() StandardId {
	return StandardId(o.v4l2Output.std)
}

func (o OutputInfo) GetCapabilities() uint32 {
	return uint32(o.v4l2Output.capabilities)
}

// GetCurrentVideoInputIndex returns the currently selected video input index
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-input.html
func GetCurrentVideoInputIndex(fd uintptr) (int32, error) {
	var index int32
	if err := send(fd, C.VIDIOC_G_INPUT, uintptr(unsafe.Pointer(&index))); err != nil {
		return -1, fmt.Errorf("video input get: %w", err)
	}
	return index, nil
}

// GetVideoInputInfo returns specified input information for video device
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enuminput.html
func GetVideoInputInfo(fd uintptr, index uint32) (InputInfo, error) {
	var input C.struct_v4l2_input
	input.index = C.uint(index)
	if err := send(fd, C.VIDIOC_ENUMINPUT, uintptr(unsafe.Pointer(&input))); err != nil {
		return InputInfo{}, fmt.Errorf("video input info: index %d: %w", index, err)
	}
	return InputInfo{v4l2Input: input}, nil
}

// GetAllVideoInputInfo returns all input information for device by
// iterating from input index = 0 until an error (EINVAL) is returned.
func GetAllVideoInputInfo(fd uintptr) (result []InputInfo, err error) {
	index := uint32(0)
	for {
		var input C.struct_v4l2_input
		input.index = C.uint(index)
		if err = send(fd, C.VIDIOC_ENUMINPUT, uintptr(unsafe.Pointer(&input))); err != nil {
			errno := err.(sys.Errno)
			if errno.Is(sys.EINVAL) && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("all video input info: %w", err)
		}
		result = append(result, InputInfo{v4l2Input: input})
		index++
	}
	return result, nil
}

// SetVideoInputIndex sets the current video input index
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-s-input.html
func SetVideoInputIndex(fd uintptr, index int32) error {
	if err := send(fd, C.VIDIOC_S_INPUT, uintptr(unsafe.Pointer(&index))); err != nil {
		return fmt.Errorf("video input set: index %d: %w", index, err)
	}
	return nil
}

// GetCurrentVideoOutputIndex returns the currently selected video output index
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-output.html
func GetCurrentVideoOutputIndex(fd uintptr) (int32, error) {
	var index int32
	if err := send(fd, C.VIDIOC_G_OUTPUT, uintptr(unsafe.Pointer(&index))); err != nil {
		return -1, fmt.Errorf("video output get: %w", err)
	}
	return index, nil
}

// SetVideoOutputIndex sets the current video output index
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-s-output.html
func SetVideoOutputIndex(fd uintptr, index int32) error {
	if err := send(fd, C.VIDIOC_S_OUTPUT, uintptr(unsafe.Pointer(&index))); err != nil {
		return fmt.Errorf("video output set: index %d: %w", index, err)
	}
	return nil
}

// GetVideoOutputInfo returns specified output information for video device
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumoutput.html
func GetVideoOutputInfo(fd uintptr, index uint32) (OutputInfo, error) {
	var output C.struct_v4l2_output
	output.index = C.uint(index)
	if err := send(fd, C.VIDIOC_ENUMOUTPUT, uintptr(unsafe.Pointer(&output))); err != nil {
		return OutputInfo{}, fmt.Errorf("video output info: index %d: %w", index, err)
	}
	return OutputInfo{v4l2Output: output}, nil
}

// GetAllVideoOutputInfo returns all output information for device by
// iterating from output index = 0 until an error (EINVAL) is returned.
func GetAllVideoOutputInfo(fd uintptr) (result []OutputInfo, err error) {
	index := uint32(0)
	for {
		var output C.struct_v4l2_output
		output.index = C.uint(index)
		if err = send(fd, C.VIDIOC_ENUMOUTPUT, uintptr(unsafe.Pointer(&output))); err != nil {
			errno := err.(sys.Errno)
			if errno.Is(sys.EINVAL) && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("all video output info: %w", err)
		}
		result = append(result, OutputInfo{v4l2Output: output})
		index++
	}
	return result, nil
}

// QueryInputStatus queries the current status of the selected video input.
// This includes signal detection, power status, and color information.
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enuminput.html#input-status
func QueryInputStatus(fd uintptr) (InputStatus, error) {
	index, err := GetCurrentVideoInputIndex(fd)
	if err != nil {
		return 0, err
	}

	info, err := GetVideoInputInfo(fd, uint32(index))
	if err != nil {
		return 0, err
	}

	return InputStatus(info.GetStatus()), nil
}

// QueryOutputStatus queries the current status of the selected video output.
// Note: The V4L2 API does not provide a status field for outputs like it does for inputs.
// This function returns 0 (OK) if the output can be successfully queried.
func QueryOutputStatus(fd uintptr) (OutputStatus, error) {
	index, err := GetCurrentVideoOutputIndex(fd)
	if err != nil {
		return 0, err
	}

	_, err = GetVideoOutputInfo(fd, uint32(index))
	if err != nil {
		return 0, err
	}

	// Output status is always OK if we can query it successfully
	return 0, nil
}
