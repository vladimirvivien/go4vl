package v4l2

import (
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// ioctl uses a 32-bit value to encode commands sent to the kernel for device control.
// Requests sent via ioctl uses a 32-bit value with the following layout:
// - lower 16 bits: ioctl command
// - Upper 14 bits: size of the parameter structure
// - MSB 2 bits: are reserved for indicating the ``access mode''.
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/asm-generic/ioctl.h

const (
	// ioctl op direction:
	// Write: userland is writing and kernel is reading.
	// Read:  userland is reading and kernel is writing.
	iocOpNone  = 0
	iocOpWrite = 1
	iocOpRead  = 2

	// ioctl command bit sizes
	iocTypeBits   = 8
	iocNumberBits = 8
	iocSizeBits   = 14
	iocOpBits     = 2

	// ioctl bit layout positions
	numberPos = 0
	typePos   = numberPos + iocNumberBits
	sizePos   = typePos + iocTypeBits
	opPos     = sizePos + iocSizeBits
)

// iocEnc encodes V42L API command as value.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L69
func iocEnc(iocMode, iocType, number, size uintptr) uintptr {
	return (iocMode << opPos) | (iocType << typePos) | (number << numberPos) | (size << sizePos)
}

// iocEncRead encodes ioctl command where program reads result from kernel.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L86
func iocEncRead(iocType, number, size uintptr) uintptr {
	return iocEnc(iocOpRead, iocType, number, size)
}

// iocEncWrite encodes ioctl command where program writes values read by the kernel.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L87
func iocEncWrite(iocType, number, size uintptr) uintptr {
	return iocEnc(iocOpWrite, iocType, number, size)
}

// iocEncReadWrite encodes ioctl command for program reads and program writes.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/asm-generic/ioctl.h#L88
func iocEncReadWrite(iocType, number, size uintptr) uintptr {
	return iocEnc(iocOpRead|iocOpWrite, iocType, number, size)
}

// ioctl is a wrapper for Syscall(SYS_IOCTL)
func ioctl(fd, req, arg uintptr) (err error) {
	if _, _, errno := sys.Syscall(sys.SYS_IOCTL, fd, req, arg); errno != 0 {
		err = errno
		return
	}
	return nil
}

// V4L2 command request values for ioctl
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2510
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/user-func.html

var (
	VidiocQueryCap       = iocEncRead('V', 0, uintptr(unsafe.Sizeof(v4l2Capability{})))          // Represents command VIDIOC_QUERYCAP
	VidiocEnumFmt        = iocEncReadWrite('V', 2, uintptr(unsafe.Sizeof(v4l2FormatDesc{})))     // Represents command VIDIOC_ENUM_FMT
	VidiocGetFormat      = iocEncReadWrite('V', 4, uintptr(unsafe.Sizeof(v4l2Format{})))         // Represents command VIDIOC_G_FMT
	VidiocSetFormat      = iocEncReadWrite('V', 5, uintptr(unsafe.Sizeof(v4l2Format{})))         // Represents command VIDIOC_S_FMT
	VidiocReqBufs        = iocEncReadWrite('V', 8, uintptr(unsafe.Sizeof(RequestBuffers{})))     // Represents command VIDIOC_REQBUFS
	VidiocQueryBuf       = iocEncReadWrite('V', 9, uintptr(unsafe.Sizeof(BufferInfo{})))         // Represents command VIDIOC_QUERYBUF
	VidiocQueueBuf       = iocEncReadWrite('V', 15, uintptr(unsafe.Sizeof(BufferInfo{})))        // Represents command VIDIOC_QBUF
	VidiocDequeueBuf     = iocEncReadWrite('V', 17, uintptr(unsafe.Sizeof(BufferInfo{})))        // Represents command VIDIOC_DQBUF
	VidiocStreamOn       = iocEncWrite('V', 18, uintptr(unsafe.Sizeof(int32(0))))                // Represents command VIDIOC_STREAMON
	VidiocStreamOff      = iocEncWrite('V', 19, uintptr(unsafe.Sizeof(int32(0))))                // Represents command VIDIOC_STREAMOFF
	VidiocCropCap        = iocEncReadWrite('V', 58, uintptr(unsafe.Sizeof(CropCapability{})))    // Represents command VIDIOC_CROPCAP
	VidiocSetCrop        = iocEncWrite('V', 60, uintptr(unsafe.Sizeof(Crop{})))                  // Represents command VIDIOC_S_CROP
	VidiocEnumFrameSizes = iocEncReadWrite('V', 74, uintptr(unsafe.Sizeof(v4l2FrameSizeEnum{}))) // Represents command VIDIOC_ENUM_FRAMESIZES
)

// Send sends a raw ioctl request to the kernel (via ioctl syscall)
func Send(fd, req, arg uintptr) error {
	return ioctl(fd, req, arg)
}
