package v4l2

import (
	"time"

	sys "golang.org/x/sys/unix"
)

// TimestampMode controls how V4L2 buffer timestamps are converted into time.Time values.
type TimestampMode int

const (
	// TimestampDefault returns the raw kernel timestamp as time.Unix(sec, usec*1000)
	// without clock conversion.
	TimestampDefault TimestampMode = iota

	// TimestampAuto inspects BufFlagTimestampMask and converts monotonic kernel
	// timestamps (BufFlagTimestampMonotonic) to wall-clock time while leaving
	// non-monotonic timestamps unchanged.
	TimestampAuto

	// TimestampMonoToAbs treats the kernel timestamp as CLOCK_MONOTONIC and
	// converts it to wall-clock time regardless of buffer flags.
	TimestampMonoToAbs
)

// MonotonicToWallClock converts a CLOCK_MONOTONIC sys.Timeval into an approximate
// wall-clock time.Time using the current offset between CLOCK_REALTIME and CLOCK_MONOTONIC.
func MonotonicToWallClock(ts sys.Timeval) time.Time {
	var monoNow sys.Timespec
	if err := sys.ClockGettime(sys.CLOCK_MONOTONIC, &monoNow); err != nil {
		return time.Unix(int64(ts.Sec), int64(ts.Usec)*1000)
	}
	wallNow := time.Now()
	monoNowNs := monoNow.Sec*int64(time.Second) + monoNow.Nsec
	tsNs := int64(ts.Sec)*int64(time.Second) + int64(ts.Usec)*int64(time.Microsecond)
	return wallNow.Add(time.Duration(tsNs - monoNowNs))
}

// BufferTimestamp converts a V4L2 buffer timestamp and flags into time.Time
// according to mode.
func BufferTimestamp(ts sys.Timeval, flags BufFlag, mode TimestampMode) time.Time {
	switch mode {
	case TimestampMonoToAbs:
		return MonotonicToWallClock(ts)
	case TimestampAuto:
		if flags&BufFlagTimestampMask == BufFlagTimestampMonotonic {
			return MonotonicToWallClock(ts)
		}
		return time.Unix(int64(ts.Sec), int64(ts.Usec)*1000)
	default:
		return time.Unix(int64(ts.Sec), int64(ts.Usec)*1000)
	}
}

// TimestampType returns the masked V4L2_BUF_FLAG_TIMESTAMP_* clock source of the buffer.
func (b Buffer) TimestampType() BufFlag {
	return b.Flags & BufFlagTimestampMask
}

// Time converts the buffer's kernel timestamp to time.Time using mode.
func (b Buffer) Time(mode TimestampMode) time.Time {
	return BufferTimestamp(b.Timestamp, b.Flags, mode)
}

// TimecodeType
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_timecode#timecode-type
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L886
type TimecodeType = uint32

const (
	TimecodeType24FPS TimecodeType = cV4L2_TC_TYPE_24FPS
	TimecodeType25FPS TimecodeType = cV4L2_TC_TYPE_25FPS
	TimecodeType30FPS TimecodeType = cV4L2_TC_TYPE_30FPS
	TimecodeType50FPS TimecodeType = cV4L2_TC_TYPE_50FPS
	TimecodeType60FPS TimecodeType = cV4L2_TC_TYPE_60FPS
)

// TimecodeFlag
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_timecode#timecode-flags
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L892
type TimecodeFlag = uint32

const (
	TimecodeFlagDropFrame  TimecodeFlag = cV4L2_TC_FLAG_DROPFRAME
	TimecodeFlagColorFrame TimecodeFlag = cV4L2_TC_FLAG_COLORFRAME
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
