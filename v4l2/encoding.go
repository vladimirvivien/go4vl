package v4l2

import "bytes"

// ioctl command API encoding:
// ioctl command encoding uses 32 bits total:
// - command in lower 16 bits
// - size of the parameter structure in the lower 14 bits of the upper 16 bits.
// - The highest 2 bits are reserved for indicating the ``access mode''.
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

// encodes V42L API command
func encode(iocMode, iocType, number, size uintptr) uintptr {
	return (iocMode << opPos) | (iocType << typePos) | (number << numberPos) | (size << sizePos)
}

// encodeRead encodes ioctl read command
func encodeRead(iocType, number, size uintptr) uintptr {
	return encode(iocOpRead, iocType, number, size)
}

// encodeWrite encodes ioctl write command
func encodeWrite(iocType, number, size uintptr) uintptr {
	return encode(iocOpWrite, iocType, number, size)
}

// encodeReadWrite encodes ioctl command for read or write
func encodeReadWrite(iocType, number, size uintptr) uintptr {
	return encode(iocOpRead|iocOpWrite, iocType, number, size)
}

// fourcc implements the four character code encoding found
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L81
// #define v4l2_fourcc(a, b, c, d)\
// 	 ((__u32)(a) | ((__u32)(b) << 8) | ((__u32)(c) << 16) | ((__u32)(d) << 24))
func fourcc(a, b, c, d uint32) uint32 {
	return (a | b<<8) | c<<16 | d<<24
}

// GoString encodes C null-terminated string to Go string
func GoString(s []byte) string {
	null := bytes.Index(s, []byte{0})
	if null < 0 {
		return ""
	}
	return string(s[:null])
}
