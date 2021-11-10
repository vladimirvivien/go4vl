package v4l2

// #include <linux/videodev2.h>
import "C"

// TimecodeType
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_timecode#timecode-type
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L886
type TimecodeType = uint32

const (
	TimecodeType24FPS TimecodeType = C.V4L2_TC_TYPE_24FPS
	TimecodeType25FPS TimecodeType = C.V4L2_TC_TYPE_25FPS
	TimecodeType30FPS TimecodeType = C.V4L2_TC_TYPE_30FPS
	TimecodeType50FPS TimecodeType = C.V4L2_TC_TYPE_50FPS
	TimecodeType60FPS TimecodeType = C.V4L2_TC_TYPE_60FPS
)

// TimecodeFlag
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_timecode#timecode-flags
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L892
type TimecodeFlag = uint32

const (
	TimecodeFlagDropFrame  TimecodeFlag = C.V4L2_TC_FLAG_DROPFRAME
	TimecodeFlagColorFrame TimecodeFlag = C.V4L2_TC_FLAG_COLORFRAME
)

// Timecode (v4l2_timecode)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_timecode#c.V4L.v4l2_timecode
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L875
type Timecode struct {
	Type    TimecodeType
	Flags   TimecodeFlag
	Frames  uint8
	Seconds uint8
	Minutes uint8
	Hours   uint8
	_       [4]uint8 // userbits
}
