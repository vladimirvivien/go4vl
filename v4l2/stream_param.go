package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// StreamParamFlag is a type alias for uint32, used for flags within streaming parameters.
// These flags define capabilities or modes related to streaming.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-parm.html#parm-flags
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1214
type StreamParamFlag = uint32

// Stream Parameter Flag Constants
const (
	// StreamParamModeHighQuality indicates that high-quality mode is preferred for streaming.
	StreamParamModeHighQuality StreamParamFlag = C.V4L2_MODE_HIGHQUALITY
	// StreamParamTimePerFrame indicates that the device supports setting the time per frame (frame rate).
	StreamParamTimePerFrame StreamParamFlag = C.V4L2_CAP_TIMEPERFRAME
)

// StreamParam holds streaming parameters for a V4L2 device.
// It corresponds to the `v4l2_streamparm` struct in the Linux kernel, which contains a union
// for capture (`v4l2_captureparm`) and output (`v4l2_outputparm`) parameters.
// The `Type` field indicates which member of the union is valid.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-parm.html#c.V4L.v4l2_streamparm
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2362
type StreamParam struct {
	// Type specifies the buffer type (e.g., video capture, video output) these parameters apply to.
	Type IOType // Corresponds to v4l2_buf_type
	// Capture holds parameters specific to video capture. Valid if Type indicates a capture stream.
	Capture CaptureParam
	// Output holds parameters specific to video output. Valid if Type indicates an output stream.
	Output OutputParam
}

// CaptureParam stores streaming parameters specific to video capture devices.
// It corresponds to the `v4l2_captureparm` struct in the Linux kernel.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-parm.html#c.V4L.v4l2_captureparm
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1205
type CaptureParam struct {
	// Capability indicates capture capabilities, e.g., V4L2_CAP_TIMEPERFRAME. See StreamParamFlag.
	Capability StreamParamFlag
	// CaptureMode is a driver-specific capture mode. V4L2_MODE_HIGHQUALITY is a common flag. See StreamParamFlag.
	CaptureMode StreamParamFlag
	// TimePerFrame is the desired time interval between frames (1/frame rate).
	TimePerFrame Fract
	// ExtendedMode is for driver-specific extended features.
	ExtendedMode uint32
	// ReadBuffers is the recommended minimum number of buffers for read() based I/O.
	ReadBuffers uint32
	// reserved space in C struct
	_ [4]uint32
}

// OutputParam stores streaming parameters specific to video output devices.
// It corresponds to the `v4l2_outputparm` struct in the Linux kernel.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-parm.html#c.V4L.v4l2_outputparm
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1228
type OutputParam struct {
	// Capability indicates output capabilities, e.g., V4L2_CAP_TIMEPERFRAME. See StreamParamFlag.
	Capability StreamParamFlag
	// CaptureMode here refers to `outputmode` in the C struct `v4l2_outputparm`.
	// It is a driver-specific output mode. See StreamParamFlag.
	CaptureMode StreamParamFlag // Maps to C.outputmode
	// TimePerFrame is the desired time interval between frames.
	TimePerFrame Fract
	// ExtendedMode is for driver-specific extended features.
	ExtendedMode uint32
	// WriteBuffers is the recommended minimum number of buffers for write() based I/O.
	WriteBuffers uint32
	// reserved space in C struct
	_ [4]uint32
}

// GetStreamParam retrieves the current streaming parameters for the specified buffer type (e.g., video capture).
// It takes the file descriptor of the V4L2 device and the BufType.
// It returns a StreamParam struct populated with the device's parameters and an error if the VIDIOC_G_PARM ioctl call fails.
// Note: The current implementation populates both Capture and Output fields from the C union
// and sets the returned StreamParam.Type to BufTypeVideoCapture regardless of the input bufType.
// This might not accurately reflect which part of the C union (capture or output) is truly active based on bufType.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-parm.html
func GetStreamParam(fd uintptr, bufType BufType) (StreamParam, error) {
	var v4l2Param C.struct_v4l2_streamparm
	v4l2Param._type = C.uint(bufType)

	if err := send(fd, C.VIDIOC_G_PARM, uintptr(unsafe.Pointer(&v4l2Param))); err != nil {
		return StreamParam{}, fmt.Errorf("stream param: %w", err)
	}

	// The C struct v4l2_streamparm has a union 'parm' for v4l2_captureparm and v4l2_outputparm.
	// Here, we cast to both. The caller should use the one corresponding to v4l2Param._type.
	capture := *(*CaptureParam)(unsafe.Pointer(&v4l2Param.parm[0]))
	// The original code for output parameter extraction was potentially problematic due to direct offset calculation.
	// A safer way depends on how the C union is structured and aligned.
	// Assuming the union overlays capture and output params at the same memory location:
	var output OutputParam
	if bufType == BufTypeVideoOutput || bufType == BufTypeVideoOutputMPlane { // Check if output params are relevant
		output = *(*OutputParam)(unsafe.Pointer(&v4l2Param.parm[0])) // Cast to OutputParam
	}


	// The returned StreamParam.Type is hardcoded to BufTypeVideoCapture in the original code.
	// It should ideally reflect the 'bufType' argument passed to the function.
	return StreamParam{
		Type:    IOType(bufType), // Reflect the queried buffer type
		Capture: capture,
		Output:  output,
	}, nil
}

// SetStreamParam sets the streaming parameters for the specified buffer type.
// It takes the file descriptor, the BufType, and a StreamParam struct containing the desired parameters.
// Depending on the bufType, it will set either the capture or output parameters from the provided StreamParam.
// Returns an error if the VIDIOC_S_PARM ioctl call fails.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-parm.html
func SetStreamParam(fd uintptr, bufType BufType, param StreamParam) error {
	var v4l2Parm C.struct_v4l2_streamparm
	v4l2Parm._type = C.uint(bufType)
	if bufType == BufTypeVideoCapture || bufType == BufTypeVideoCaptureMPlane {
		*(*C.struct_v4l2_captureparm)(unsafe.Pointer(&v4l2Parm.parm[0])) = *(*C.struct_v4l2_captureparm)(unsafe.Pointer(&param.Capture))
	} else if bufType == BufTypeVideoOutput || bufType == BufTypeVideoOutputMPlane {
		// The C struct v4l2_streamparm has a union 'parm'. We cast to the appropriate type.
		*(*C.struct_v4l2_outputparm)(unsafe.Pointer(&v4l2Parm.parm[0])) = *(*C.struct_v4l2_outputparm)(unsafe.Pointer(&param.Output))
	}
	// Not handling other buffer types here, as streamparm is typically for video capture/output.

	if err := send(fd, C.VIDIOC_S_PARM, uintptr(unsafe.Pointer(&v4l2Parm))); err != nil {
		return fmt.Errorf("stream param: %w", err)
	}

	return nil
}
