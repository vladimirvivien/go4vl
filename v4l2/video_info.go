package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// InputStatus is a type alias for uint32, representing the status of a V4L2 video input.
// These flags indicate conditions like no power, no signal, etc.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enuminput.html#input-status-flags
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1500
type InputStatus = uint32

// Input Status Flag Constants
var (
	// InputStatusNoPower indicates the input has no power.
	InputStatusNoPower InputStatus = C.V4L2_IN_ST_NO_POWER
	// InputStatusNoSignal indicates no signal is detected on the input.
	InputStatusNoSignal InputStatus = C.V4L2_IN_ST_NO_SIGNAL
	// InputStatusNoColor indicates no color signal is detected (e.g., B&W image on a color input).
	InputStatusNoColor InputStatus = C.V4L2_IN_ST_NO_COLOR
	// Note: A status of 0 means the input is OK. Other status bits might be driver-specific.
)

// InputStatuses provides a map of common InputStatus constants to their human-readable string descriptions.
// A status of 0 is explicitly mapped to "ok".
var InputStatuses = map[InputStatus]string{
	0:                   "ok", // V4L2_IN_ST_OK is typically 0
	InputStatusNoPower:  "no power",
	InputStatusNoSignal: "no signal",
	InputStatusNoColor:  "no color",
}

// InputType is a type alias for uint32, representing the type of a V4L2 video input.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enuminput.html#c.v4l2_inputtype
// See also C.V4L2_INPUT_TYPE_TUNER, C.V4L2_INPUT_TYPE_CAMERA in videodev2.h
type InputType = uint32

// Input Type Constants
const (
	// InputTypeTuner indicates the input is a tuner (e.g., for TV or radio).
	InputTypeTuner InputType = C.V4L2_INPUT_TYPE_TUNER
	// InputTypeCamera indicates the input is a camera.
	InputTypeCamera InputType = C.V4L2_INPUT_TYPE_CAMERA
	// InputTypeTouch indicates the input is a touch device. Note: This constant is defined in this Go package,
	// as C.V4L2_INPUT_TYPE_TOUCH might not be universally available in older kernel headers.
	// Its value here is assigned by iota + 1 relative to other input types from C.
	InputTypeTouch InputType = iota + 1 // Starting from C.V4L2_INPUT_TYPE_CAMERA + 1 (assuming Camera is 2)
)

// StandardId is a type alias for uint64, representing a bitmask of video standards supported by an input.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/standards.html
type StandardId = uint64

// InputInfo provides information about a V4L2 video input.
// It wraps the C struct `v4l2_input` and provides getter methods to access its fields.
// This structure is used with the VIDIOC_ENUMINPUT ioctl to enumerate available inputs
// and to query their properties.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enuminput.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1649
type InputInfo struct {
	v4l2Input C.struct_v4l2_input // Internal C struct
}

// GetIndex returns the zero-based index of the input.
func (i InputInfo) GetIndex() uint32 {
	return uint32(i.v4l2Input.index)
}

// GetName returns the human-readable name of the input (e.g., "Camera 1", "Composite").
func (i InputInfo) GetName() string {
	return C.GoString((*C.char)(unsafe.Pointer(&i.v4l2Input.name[0])))
}

// GetInputType returns the type of the input (e.g., Tuner, Camera). See InputType constants.
func (i InputInfo) GetInputType() InputType {
	return InputType(i.v4l2Input._type)
}

// GetAudioset returns a bitmask indicating which audio inputs are associated with this video input.
// This is relevant if the V4L2 device supports audio.
func (i InputInfo) GetAudioset() uint32 {
	return uint32(i.v4l2Input.audioset)
}

// GetTuner returns the index of the tuner associated with this input, if applicable.
// Only valid if GetInputType() is InputTypeTuner.
func (i InputInfo) GetTuner() uint32 {
	return uint32(i.v4l2Input.tuner)
}

// GetStandardId returns a bitmask of video standards (e.g., PAL, NTSC) supported by this input.
// See StandardId type and V4L2 standard constants (v4l2_std_id).
func (i InputInfo) GetStandardId() StandardId {
	return StandardId(i.v4l2Input.std)
}

// GetStatus returns the current status of the input (e.g., no signal, no power). See InputStatus constants.
func (i InputInfo) GetStatus() uint32 {
	return uint32(i.v4l2Input.status)
}

// GetCapabilities returns capability flags for this input (e.g., custom timings, DV timings).
// See V4L2_IN_CAP_* constants in kernel headers.
func (i InputInfo) GetCapabilities() uint32 {
	return uint32(i.v4l2Input.capabilities)
}

// GetCurrentVideoInputIndex retrieves the index of the currently selected video input.
// It takes the file descriptor of the V4L2 device.
// Returns the 0-based index of the current input, or an error if the VIDIOC_G_INPUT ioctl call fails.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-input.html
func GetCurrentVideoInputIndex(fd uintptr) (int32, error) {
	var index int32
	if err := send(fd, C.VIDIOC_G_INPUT, uintptr(unsafe.Pointer(&index))); err != nil {
		return -1, fmt.Errorf("get current video input index: VIDIOC_G_INPUT failed: %w", err)
	}
	return index, nil
}

// GetVideoInputInfo retrieves information about a specific video input, identified by its index.
// It takes the file descriptor of the V4L2 device and the zero-based index of the input.
// Returns an InputInfo struct populated with the input's details, and an error if the VIDIOC_ENUMINPUT ioctl call fails
// (e.g., if the index is out of bounds).
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-enuminput.html
func GetVideoInputInfo(fd uintptr, index uint32) (InputInfo, error) {
	var input C.struct_v4l2_input
	input.index = C.uint(index)
	if err := send(fd, C.VIDIOC_ENUMINPUT, uintptr(unsafe.Pointer(&input))); err != nil {
		return InputInfo{}, fmt.Errorf("get video input info: VIDIOC_ENUMINPUT failed for index %d: %w", index, err)
	}
	return InputInfo{v4l2Input: input}, nil
}

// GetAllVideoInputInfo retrieves information for all available video inputs on the device.
// It iterates by calling GetVideoInputInfo with increasing indices, starting from 0,
// until an error (typically EINVAL, indicating no more inputs) is encountered.
// It takes the file descriptor of the V4L2 device.
//
// Returns a slice of InputInfo structs and any error encountered during the final failing call.
// If some inputs were successfully retrieved before an error, those will be returned along with the error.
// If the first call fails, it returns an empty slice and the error.
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
