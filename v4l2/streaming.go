package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// Streaming with Buffers
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html

// BufType (v4l2_buf_type)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_buf_type#c.V4L.v4l2_buf_type
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L141
type BufType = uint32

const (
	BufTypeVideoCapture BufType = C.V4L2_BUF_TYPE_VIDEO_CAPTURE
	BufTypeVideoOutput  BufType = C.V4L2_BUF_TYPE_VIDEO_OUTPUT
	BufTypeOverlay      BufType = C.V4L2_BUF_TYPE_VIDEO_OVERLAY
)

// StreamType (v4l2_memory)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/mmap.html?highlight=v4l2_memory_mmap
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L188
type StreamType = uint32

const (
	StreamTypeMMAP    StreamType = C.V4L2_MEMORY_MMAP
	StreamTypeUserPtr StreamType = C.V4L2_MEMORY_USERPTR
	StreamTypeOverlay StreamType = C.V4L2_MEMORY_OVERLAY
	StreamTypeDMABuf  StreamType = C.V4L2_MEMORY_DMABUF
)

// TODO implement vl42_create_buffers

// RequestBuffers (v4l2_requestbuffers) is used to request buffer allocation initializing
// streaming for memory mapped, user pointer, or DMA buffer access.
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L949
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-reqbufs.html?highlight=v4l2_requestbuffers#c.V4L.v4l2_requestbuffers
type RequestBuffers struct {
	Count        uint32
	StreamType   uint32
	Memory       uint32
	Capabilities uint32
	_            [1]uint32
}

// Buffer (v4l2_buffer) is used to send buffers info between application and driver
// after streaming IO has been initialized.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html#c.V4L.v4l2_buffer
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1037
//
type Buffer struct {
	Index      uint32
	StreamType uint32
	BytesUsed  uint32
	Flags      uint32
	Field      uint32
	Timestamp  sys.Timeval
	Timecode   Timecode
	Sequence   uint32
	Memory     uint32
	Info       BufferInfo // union m
	Length     uint32
	Reserved2  uint32
	RequestFD  int32
}

// makeBuffer makes a Buffer value from C.struct_v4l2_buffer
func makeBuffer(v4l2Buf C.struct_v4l2_buffer) Buffer {
	return Buffer{
		Index:      uint32(v4l2Buf.index),
		StreamType: uint32(v4l2Buf._type),
		BytesUsed:  uint32(v4l2Buf.bytesused),
		Flags:      uint32(v4l2Buf.flags),
		Field:      uint32(v4l2Buf.field),
		Timestamp:  *(*sys.Timeval)(unsafe.Pointer(&v4l2Buf.timestamp)),
		Timecode:   *(*Timecode)(unsafe.Pointer(&v4l2Buf.timecode)),
		Sequence:   uint32(v4l2Buf.sequence),
		Memory:     uint32(v4l2Buf.memory),
		Info:       *(*BufferInfo)(unsafe.Pointer(&v4l2Buf.m[0])),
		Length:     uint32(v4l2Buf.length),
		Reserved2:  uint32(v4l2Buf.reserved2),
		RequestFD:  *(*int32)(unsafe.Pointer(&v4l2Buf.anon0[0])),
	}
}

// BufferInfo represents Union of several values in type Buffer
// that are used to service the stream depending on the type of streaming
// selected (MMap, User pointer, planar, file descriptor for DMA)
type BufferInfo struct {
	Offset  uint32
	UserPtr uintptr
	Planes  *Plane
	FD      int32
}

// Plane (see struct v4l2_plane) represents a plane in a multi-planar buffers
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html#c.V4L.v4l2_plane
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L990
//
type Plane struct {
	BytesUsed  uint32
	Length     uint32
	Info       PlaneInfo // union m
	DataOffset uint32
}

// PlaneInfo representes the combination of type
// of type of memory stream that can be serviced for the
// associated plane.
type PlaneInfo struct {
	MemOffset uint32
	UserPtr   uintptr
	FD        int32
}

// StreamOn requests streaming to be turned on for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-streamon.html
func StreamOn(fd uintptr) error {
	bufType := BufTypeVideoCapture
	if err := send(fd, C.VIDIOC_STREAMON, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return fmt.Errorf("stream on: %w", err)
	}
	return nil
}

// StreamOff requests streaming to be turned off for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-streamon.html
func StreamOff(fd uintptr) error {
	bufType := BufTypeVideoCapture
	if err := send(fd, C.VIDIOC_STREAMOFF, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return fmt.Errorf("stream off: %w", err)
	}
	return nil
}

// InitBuffers sends buffer allocation request to initialize buffer IO
// for video capture when using either mem map, user pointer, or DMA buffers.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-reqbufs.html#vidioc-reqbufs
func InitBuffers(fd uintptr, buffSize uint32) (RequestBuffers, error) {
	var req C.struct_v4l2_requestbuffers
	req.count = C.uint(buffSize)
	req._type = C.uint(BufTypeVideoCapture)
	req.memory = C.uint(StreamTypeMMAP)

	if err := send(fd, C.VIDIOC_REQBUFS, uintptr(unsafe.Pointer(&req))); err != nil {
		return RequestBuffers{}, fmt.Errorf("request buffers: %w", err)
	}
	if req.count < 2 {
		return RequestBuffers{}, errors.New("request buffers: insufficient memory on device")
	}

	return *(*RequestBuffers)(unsafe.Pointer(&req)), nil
}

// GetBuffer retrieves bunffer info for allocated buffers at provided index.
// This call should take place after buffers are allocated (for mmap for instance).
func GetBuffer(fd uintptr, index uint32) (Buffer, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(BufTypeVideoCapture)
	v4l2Buf.memory = C.uint(StreamTypeMMAP)
	v4l2Buf.index = C.uint(index)

	if err := send(fd, C.VIDIOC_QUERYBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		return Buffer{}, fmt.Errorf("query buffer: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// MapMemoryBuffer creates a local buffer mapped to the address space of the device specified by fd.
func MapMemoryBuffer(fd uintptr, offset int64, len int) ([]byte, error) {
	data, err := sys.Mmap(int(fd), offset, len, sys.PROT_READ|sys.PROT_WRITE, sys.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("map memory buffer: %w", err)
	}
	return data, nil
}

// UnmapMemoryBuffer removes the buffer that was previously mapped.
func UnmapMemoryBuffer(buf []byte) error {
	if err := sys.Munmap(buf); err != nil {
		return fmt.Errorf("unmap memory buffer: %w", err)
	}
	return nil
}

// QueueBuffer enqueues a buffer in the device driver (as empty for capturing, or filled for video output)
// when using either memory map, user pointer, or DMA buffers. Buffer is returned with
// additional information about the queued buffer.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html#vidioc-qbuf
func QueueBuffer(fd uintptr, index uint32) (Buffer, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(BufTypeVideoCapture)
	v4l2Buf.memory = C.uint(StreamTypeMMAP)
	v4l2Buf.index = C.uint(index)

	if err := send(fd, C.VIDIOC_QBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		return Buffer{}, fmt.Errorf("buffer queue: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// DequeueBuffer dequeues a buffer in the device driver, marking it as consumed by the application,
// when using either memory map, user pointer, or DMA buffers. Buffer is returned with
// additional information about the dequeued buffer.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html#vidioc-qbuf
func DequeueBuffer(fd uintptr) (Buffer, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(BufTypeVideoCapture)
	v4l2Buf.memory = C.uint(StreamTypeMMAP)

	if err := send(fd, C.VIDIOC_DQBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		return Buffer{}, fmt.Errorf("buffer dequeue: %w", err)

	}

	return makeBuffer(v4l2Buf), nil
}

// WaitForDeviceRead blocks until the specified device is
// ready to be read or has timedout.
func WaitForDeviceRead(fd uintptr, timeout time.Duration) error {
	timeval := sys.NsecToTimeval(timeout.Nanoseconds())
	var fdsRead sys.FdSet
	fdsRead.Set(int(fd))
	for {
		n, err := sys.Select(int(fd+1), &fdsRead, nil, nil, &timeval)
		switch n {
		case -1:
			if err == sys.EINTR {
				continue
			}
			return err
		case 0:
			return errors.New("wait for device ready: timeout")
		default:
			return nil
		}
	}
}
