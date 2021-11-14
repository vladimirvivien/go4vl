package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Flags for capability and capture mode fields
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-parm.html#parm-flags
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1214
type StreamParamFlag = uint32

const (
	StreamParamModeHighQuality StreamParamFlag = C.V4L2_MODE_HIGHQUALITY
	StreamParamTimePerFrame    StreamParamFlag = C.V4L2_CAP_TIMEPERFRAME
)

// CaptureParam (v4l2_captureparam)
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-parm.html#c.V4L.v4l2_captureparm
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1205
type CaptureParam struct {
	Capability   StreamParamFlag
	CaptureMode  StreamParamFlag
	TimePerFrame Fract
	ExtendedMode uint32
	ReadBuffers  uint32
	_            [4]uint32
}

// GetStreamCaptureParam returns streaming capture parameter for the driver (v4l2_streamparm).
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-parm.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2347

func GetStreamCaptureParam(fd uintptr) (CaptureParam, error) {
	var param C.struct_v4l2_streamparm
	param._type = C.uint(BufTypeVideoCapture)

	if err := send(fd, C.VIDIOC_G_PARM, uintptr(unsafe.Pointer(&param))); err != nil {
		return CaptureParam{}, fmt.Errorf("stream param: %w", err)
	}
	return *(*CaptureParam)(unsafe.Pointer(&param.parm[0])), nil
}
