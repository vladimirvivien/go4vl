package v4l2

// TimecodeType
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_timecode#timecode-type
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L886
type TimecodeType = uint32

const (
	TimecodeType24FPS TimecodeType = iota + 1 // V4L2_TC_TYPE_24FPS
	TimecodeType25FPS                         // V4L2_TC_TYPE_25FPS
	TimecodeType30FPS                         // V4L2_TC_TYPE_30FPS
	TimecodeType50FPS                         // V4L2_TC_TYPE_50FPS
	TimecodeType60FPS                         // V4L2_TC_TYPE_60FPS
)

// TimecodeFlag
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_timecode#timecode-flags
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L892
type TimecodeFlag = uint32

const (
	TimecodeFlagDropFrame  TimecodeFlag = 0x0001 // V4L2_TC_FLAG_DROPFRAME	0x0001
	TimecodeFlagColorFrame TimecodeFlag = 0x0002 // V4L2_TC_FLAG_COLORFRAME	0x0002
)

// Timecode (v4l2_timecode)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_timecode#c.V4L.v4l2_timecode
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L875
type Timecode struct {
	Type     TimecodeType
	Flags    TimecodeFlag
	Frames   uint8
	Seconds  uint8
	Minutes  uint8
	Hours    uint8
	Userbits [4]uint8
}
