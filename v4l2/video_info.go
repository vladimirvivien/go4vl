package v4l2

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
// iterating from input index = 0 until an error (EINVL) is returned.
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
			return result, fmt.Errorf("all video info: %w", err)
		}
		result = append(result, InputInfo{v4l2Input: input})
		index++
	}
	return result, err
}
