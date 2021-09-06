package v4l2

import (
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// ioctl uses a 32-bit value to encode commands sent to the kernel for device control.
// Requests sent via ioctl uses the following layout:
// - lower 16 bits: ioctl command
// - Upper 14 bits: size of the parameter structure
// - MSB 2 bits: are reserved for indicating the ``access mode''.
// https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h

const (
	// ioctl command bit sizes
	// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L23
	iocNumberBits = 8
	iocTypeBits   = 8
	iocSizeBits   = 14
	iocDirBits    = 2

	// ioctl bit layout positions
	// see https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L44
	numShift  = 0
	typeShift = numShift + iocNumberBits
	sizeShift = typeShift + iocTypeBits
	dirShift  = sizeShift + iocSizeBits

	// ioctl direction bits
	// These values are from the ioctl.h in linux.
	// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L57
	iocNone  = 0 // no op
	iocWrite = 1 // userland app is writing, kernel reading
	iocRead  = 2 // userland app is reading, kernel writing

)

// iocEnc encodes V42L API command as value.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L69
func iocEnc(dir, iocType, number, size uintptr) uintptr {
	return (dir << dirShift) | (iocType << typeShift) | (number << numShift) | (size << sizeShift)
}

// iocEncR encodes ioctl command where program reads result from kernel.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L86
func iocEncR(iocType, number, size uintptr) uintptr {
	return iocEnc(iocRead, iocType, number, size)
}

// iocEncW encodes ioctl command where program writes values read by the kernel.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L87
func iocEncW(iocType, number, size uintptr) uintptr {
	return iocEnc(iocWrite, iocType, number, size)
}

// iocEncRW encodes ioctl command for program reads and program writes.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L88
func iocEncRW(iocType, number, size uintptr) uintptr {
	return iocEnc(iocRead|iocWrite, iocType, number, size)
}

// ioctl is a wrapper for Syscall(SYS_IOCTL)
func ioctl(fd, req, arg uintptr) (err sys.Errno) {
	if _, _, errno := sys.Syscall(sys.SYS_IOCTL, fd, req, arg); errno != 0 {
		if errno != 0 {
			err = errno
			return
		}
	}
	return 0
}

// V4L2 command request values for ioctl
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2510
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/user-func.html

var (
	VidiocQueryCap       = iocEncR('V', 0, unsafe.Sizeof(v4l2Capability{}))      // Represents command VIDIOC_QUERYCAP
	VidiocEnumFmt        = iocEncRW('V', 2, unsafe.Sizeof(v4l2FormatDesc{}))     // Represents command VIDIOC_ENUM_FMT
	VidiocGetFormat      = iocEncRW('V', 4, unsafe.Sizeof(v4l2Format{}))         // Represents command VIDIOC_G_FMT
	VidiocSetFormat      = iocEncRW('V', 5, unsafe.Sizeof(v4l2Format{}))         // Represents command VIDIOC_S_FMT
	VidiocReqBufs        = iocEncRW('V', 8, unsafe.Sizeof(RequestBuffers{}))     // Represents command VIDIOC_REQBUFS
	VidiocQueryBuf       = iocEncRW('V', 9, unsafe.Sizeof(BufferInfo{}))         // Represents command VIDIOC_QUERYBUF
	VidiocQueueBuf       = iocEncRW('V', 15, unsafe.Sizeof(BufferInfo{}))        // Represents command VIDIOC_QBUF
	VidiocDequeueBuf     = iocEncRW('V', 17, unsafe.Sizeof(BufferInfo{}))        // Represents command VIDIOC_DQBUF
	VidiocStreamOn       = iocEncW('V', 18, unsafe.Sizeof(int32(0)))             // Represents command VIDIOC_STREAMON
	VidiocStreamOff      = iocEncW('V', 19, unsafe.Sizeof(int32(0)))             // Represents command VIDIOC_STREAMOFF
	VidiocEnumInput      = iocEncRW('V', 26, unsafe.Sizeof(v4l2InputInfo{}))     // Represents command VIDIOC_ENUMINPUT
	VidiocGetVideoInput  = iocEncR('V', 38, unsafe.Sizeof(int32(0)))             // Represents command VIDIOC_G_INPUT
	VidiocCropCap        = iocEncRW('V', 58, unsafe.Sizeof(CropCapability{}))    // Represents command VIDIOC_CROPCAP
	VidiocSetCrop        = iocEncW('V', 60, unsafe.Sizeof(Crop{}))               // Represents command VIDIOC_S_CROP
	VidiocEnumFrameSizes = iocEncRW('V', 74, unsafe.Sizeof(v4l2FrameSizeEnum{})) // Represents command VIDIOC_ENUM_FRAMESIZES
)

// Send sends a request to the kernel (via ioctl syscall)
func Send(fd, req, arg uintptr) error {
	errno := ioctl(fd, req, arg)
	if errno == 0 {
		return nil
	}
	parsedErr := parseErrorType(errno)
	switch parsedErr {
	case ErrorUnsupported, ErrorSystem, ErrorBadArgument:
		return parsedErr
	case ErrorTimeout, ErrorTemporary:
		// TODO add code for automatic retry/recovery
		return errno
	default:
		return errno
	}
}
