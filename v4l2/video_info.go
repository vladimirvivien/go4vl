package v4l2

import (
	"fmt"
	sys "golang.org/x/sys/unix"
	"unsafe"
)

// InputStatus
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enuminput.html?highlight=v4l2_input#input-status
type InputStatus = uint32

var (
	InputStatusNoPower  = InputStatus(0x00000001) // V4L2_IN_ST_NO_POWER
	InputStatusNoSignal = InputStatus(0x00000002) // V4L2_IN_ST_NO_SIGNAL
	InputStatusNoColor  = InputStatus(0x00000004) // V4L2_IN_ST_NO_COLOR
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

type v4l2InputInfo struct {
	index        uint32
	name         [32]uint8
	inputType    InputType
	audioset     uint32
	tuner        uint32
	std          StandardId
	status       InputStatus
	capabilities uint32
	reserved     [3]uint32
	_            [4]uint8 // go compiler alignment adjustment for 32-bit platforms (Raspberry pi's, etc)
}

// InputInfo (v4l2_input)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1649
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enuminput.html
type InputInfo struct {
	v4l2InputInfo
}

func (i InputInfo) GetIndex() uint32 {
	return i.index
}

func (i InputInfo) GetName() string {
	return toGoString(i.name[:])
}

func (i InputInfo) GetInputType() InputType {
	return i.inputType
}

func (i InputInfo) GetAudioset() uint32 {
	return i.audioset
}

func (i InputInfo) GetTuner() uint32 {
	return i.tuner
}

func (i InputInfo) GetStandardId() StandardId {
	return i.std
}

func (i InputInfo) GetStatus() uint32 {
	return i.status
}

func (i InputInfo) GetCapabilities() uint32 {
	return i.capabilities
}

// GetCurrentVideoInputIndex returns the currently selected video input index
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-input.html
func GetCurrentVideoInputIndex(fd uintptr) (int32, error) {
	var index int32
	if err := Send(fd, VidiocGetVideoInput, uintptr(unsafe.Pointer(&index))); err != nil {
		return -1, fmt.Errorf("video input get: %w", err)
	}
	return index, nil
}

// GetVideoInputInfo returns specified input information for video device
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enuminput.html
func GetVideoInputInfo(fd uintptr, index uint32) (InputInfo, error) {
	input := v4l2InputInfo{index: index}
	if err := Send(fd, VidiocEnumInput, uintptr(unsafe.Pointer(&input))); err != nil {
		return InputInfo{}, fmt.Errorf("video input info: index %d: %w", index, err)
	}
	return InputInfo{v4l2InputInfo: input}, nil
}

// GetAllVideoInputInfo returns all input information for device by
// iterating from input index = 0 until an error (EINVL) is returned.
func GetAllVideoInputInfo(fd uintptr) (result []InputInfo, err error) {
	index := uint32(0)
	for {
		input := v4l2InputInfo{index: index}
		if err = Send(fd, VidiocEnumInput, uintptr(unsafe.Pointer(&input))); err != nil {
			errno := err.(sys.Errno)
			if errno.Is(sys.EINVAL) && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("all video info: %w", err)
		}
		result = append(result, InputInfo{v4l2InputInfo: input})
		index++
	}
	return result, err
}
