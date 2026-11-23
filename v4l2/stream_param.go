package v4l2

import (
	"fmt"
	"unsafe"
)

// StreamParamFlag is for capability and capture mode fields
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-parm.html#parm-flags
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1214
type StreamParamFlag = uint32

const (
	StreamParamModeHighQuality StreamParamFlag = cV4L2_MODE_HIGHQUALITY
	StreamParamTimePerFrame    StreamParamFlag = cV4L2_CAP_TIMEPERFRAME
)

// StreamParam (v4l2_streamparam)
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-parm.html#c.V4L.v4l2_streamparm
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2362
type StreamParam struct {
	Type    IOType
	Capture CaptureParam
	Output  OutputParam
}

// CaptureParam (v4l2_captureparm)
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

// OutputParam (v4l2_outputparm)
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-parm.html#c.V4L.v4l2_outputparm
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1228
type OutputParam struct {
	Capability   StreamParamFlag
	CaptureMode  StreamParamFlag
	TimePerFrame Fract
	ExtendedMode uint32
	WriteBuffers uint32
	_            [4]uint32
}

// GetStreamParam returns streaming parameters for the driver (v4l2_streamparm).
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-parm.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2362
func GetStreamParam(fd uintptr, bufType BufType) (StreamParam, error) {
	var v4l2Param cV4L2StreamParm
	v4l2Param.Type = bufType

	if err := send(fd, cVIDIOC_G_PARM, unsafe.Pointer(&v4l2Param)); err != nil {
		return StreamParam{}, fmt.Errorf("stream param: %w", err)
	}

	var capture CaptureParam
	var output OutputParam
	switch bufType {
	case BufTypeVideoCapture:
		capture = *(*CaptureParam)(unsafe.Pointer(&v4l2Param.Parm[0]))
	case BufTypeVideoOutput:
		output = *(*OutputParam)(unsafe.Pointer(&v4l2Param.Parm[0]))
	}

	return StreamParam{
		Type:    v4l2Param.Type,
		Capture: capture,
		Output:  output,
	}, nil
}

func SetStreamParam(fd uintptr, bufType BufType, param StreamParam) error {
	var v4l2Parm cV4L2StreamParm
	v4l2Parm.Type = bufType
	if bufType == BufTypeVideoCapture {
		*(*cV4L2CaptureParm)(unsafe.Pointer(&v4l2Parm.Parm[0])) = *(*cV4L2CaptureParm)(unsafe.Pointer(&param.Capture))
	}
	if bufType == BufTypeVideoOutput {
		*(*cV4L2OutputParm)(unsafe.Pointer(&v4l2Parm.Parm[0])) = *(*cV4L2OutputParm)(unsafe.Pointer(&param.Output))
	}

	if err := send(fd, cVIDIOC_S_PARM, unsafe.Pointer(&v4l2Parm)); err != nil {
		return fmt.Errorf("stream param: %w", err)
	}

	return nil
}
