package v4l2

import (
	"fmt"
	"unsafe"
)

// Flags for capability and capture mode fields
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-parm.html#parm-flags
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1214
type StreamParamFlag = uint32

const (
	StreamParamModeHighQuality StreamParamFlag = 0x0001 // V4L2_MODE_HIGHQUALITY
	StreamParamTimePerFrame    StreamParamFlag = 0x1000 // V4L2_CAP_TIMEPERFRAME
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
	reserved     [4]uint32
}

// v4l2StreamParam (v4l2_streamparam)
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-parm.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2347
//
// Field param points to an embedded union, therefore,  it points to an array sized
// as the largest knwon member of the embedded struct.  See below:
//
// struct v4l2_streamparm {
//	__u32	 type;
//	union {
//		struct v4l2_captureparm	capture;
//		struct v4l2_outputparm	output;
//		__u8	raw_data[200];
//	} parm;
//};
type v4l2StreamParam struct {
	streamType StreamMemoryType
	param      [200]byte // embedded union
}

// getCaptureParam returns CaptureParam value from v4l2StreamParam embedded union
// if p.streamType = BufTypeVideoCapture.
func (p v4l2StreamParam) getCaptureParam() CaptureParam {
	var param CaptureParam
	if p.streamType == BufTypeVideoCapture {
		param = *(*CaptureParam)(unsafe.Pointer(&p.param[0]))
	}
	return param
}

// GetStreamCaptureParam returns streaming capture parameter for the driver.
func GetStreamCaptureParam (fd uintptr)(CaptureParam, error){
	param := v4l2StreamParam{streamType: BufTypeVideoCapture}
	if err := Send(fd, VidiocGetParam, uintptr(unsafe.Pointer(&param))); err != nil {
		return CaptureParam{}, fmt.Errorf("stream param: %w", err)
	}
	return param.getCaptureParam(), nil
}