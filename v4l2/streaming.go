package v4l2

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
	BufTypeVideoCapture BufType = iota + 1 // V4L2_BUF_TYPE_VIDEO_CAPTURE = 1
	BufTypeVideoOutput                     // V4L2_BUF_TYPE_VIDEO_OUTPUT  = 2
	BufTypeOverlay                         // V4L2_BUF_TYPE_VIDEO_OVERLAY = 3
)

// StreamMemoryType (v4l2_memory)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/mmap.html?highlight=v4l2_memory_mmap
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L188
type StreamMemoryType = uint32

const (
	StreamMemoryTypeMMAP    StreamMemoryType = iota + 1 // V4L2_MEMORY_MMAP             = 1,
	StreamMemoryTypeUserPtr                             // V4L2_MEMORY_USERPTR          = 2,
	StreamMemoryTypeOverlay                             // V4L2_MEMORY_OVERLAY          = 3,
	StreamMemoryTypeDMABuf                              // V4L2_MEMORY_DMABUF           = 4,
)

// TODO implement vl42_create_buffers

// RequestBuffers (v4l2_requestbuffers)
// This type is used to allocate buffer/io resources when initializing streaming io for
// memory mapped, user pointer, or DMA buffer access.
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L949
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-reqbufs.html?highlight=v4l2_requestbuffers#c.V4L.v4l2_requestbuffers
type RequestBuffers struct {
	Count        uint32
	StreamType   uint32
	Memory       uint32
	Capabilities uint32
	Reserved     [1]uint32
}

// BufferInfo (v4l2_buffer)
// This type is used to send buffers management info between application and driver
// after streaming IO has been initialized.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html#c.V4L.v4l2_buffer
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1037
//
// BufferInfo represents type v4l2_buffer which contains unions as shown below.
// Remember, the union is represented as an arry of bytes sized as the largest
// member in bytes.
//
//   struct v4l2_buffer {
// 	   __u32			index;
// 	   __u32			type;
// 	   __u32			bytesused;
// 	   __u32			flags;
// 	   __u32			field;
// 	   struct timeval		timestamp;
// 	   struct v4l2_timecode	timecode;
// 	   __u32			sequence;
// 	   __u32			memory;
// 	   union {
// 		  __u32             offset;
// 		  unsigned long     userptr;
// 		  struct v4l2_plane *planes;
// 		  __s32		fd;
// 	   } m;
// 	   __u32	length;
// 	   __u32	reserved2;
// 	   union {
// 		  __s32		request_fd;
// 		  __u32		reserved;
// 	   };
// };
type BufferInfo struct {
	Index      uint32
	StreamType uint32
	BytesUsed  uint32
	Flags      uint32
	Field      uint32
	Timestamp  sys.Timeval
	Timecode   Timecode
	Sequence   uint32
	Memory     uint32
	m          [unsafe.Sizeof(&BufferService{})]byte // union m, cast to BufferService
	Length     uint32
	Reserved2  uint32
	RequestFD  int32
}

func (b BufferInfo) GetService() BufferService {
	m := (*BufferService)(unsafe.Pointer(&b.m[0]))
	return *m
}

// BufferService represents Union of several values in type Buffer
// that are used to service the stream depending on the type of streaming
// selected (MMap, User pointer, planar, file descriptor for DMA)
type BufferService struct {
	Offset  uint32
	UserPtr uintptr
	Planes  *PlaneInfo
	FD      int32
}

// PlaneInfo (v4l2_plane)
// Represents a plane in a multi-planar buffers
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html#c.V4L.v4l2_plane
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L990
//
// PlaneInfo includes a uinion of types as shown below:
// struct v4l2_plane {
// 	__u32			bytesused;
// 	__u32			length;
// 	union {
// 		__u32		mem_offset;
// 		unsigned long	userptr;
// 		__s32		fd;
// 	} m;
// 	__u32			data_offset;
// 	__u32			reserved[11];
// };
type PlaneInfo struct {
	BytesUsed  uint32
	Length     uint32
	m          [unsafe.Sizeof(uintptr(0))]byte // union m, cast to BufferPlaneService
	DataOffset uint32
	Reserved   [11]uint32
}

func (p PlaneInfo) GetService() PlaneService {
	m := (*PlaneService)(unsafe.Pointer(&p.m[0]))
	return *m
}

// PlaneService representes the combination of type
// of type of memory stream that can be serviced for the
// associated plane.
type PlaneService struct {
	MemOffset uint32
	UserPtr   uintptr
	FD        int32
}

// StreamOn requests streaming to be turned on for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-streamon.html
func StreamOn(fd uintptr) error {
	bufType := BufTypeVideoCapture
	if err := Send(fd, VidiocStreamOn, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return fmt.Errorf("stream on: %w", err)
	}
	return nil
}

// StreamOff requests streaming to be turned off for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-streamon.html
func StreamOff(fd uintptr) error {
	bufType := BufTypeVideoCapture
	if err := Send(fd, VidiocStreamOff, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return fmt.Errorf("stream off: %w", err)
	}
	return nil
}

// AllocateBuffers sends buffer allocation request to underlying driver
// for video capture when using either mem map, user pointer, or DMA buffers.
func AllocateBuffers(fd uintptr, buffSize uint32) (RequestBuffers, error) {
	req := RequestBuffers{
		Count:      buffSize,
		StreamType: BufTypeVideoCapture,
		Memory:     StreamMemoryTypeMMAP,
	}

	if err := Send(fd, VidiocReqBufs, uintptr(unsafe.Pointer(&req))); err != nil {
		return RequestBuffers{}, fmt.Errorf("request buffers: %w", err)
	}
	if req.Count < 2 {
		return RequestBuffers{}, errors.New("request buffers: insufficient memory on device")
	}

	return req, nil
}

// GetBuffersInfo retrieves information for allocated buffers at provided index.
// This call should take place after buffers are allocated (for mmap for instance).
func GetBufferInfo(fd uintptr, index uint32) (BufferInfo, error) {
	buf := BufferInfo{
		StreamType: BufTypeVideoCapture,
		Memory:     StreamMemoryTypeMMAP,
		Index:      index,
	}

	if err := Send(fd, VidiocQueryBuf, uintptr(unsafe.Pointer(&buf))); err != nil {
		return BufferInfo{}, fmt.Errorf("buffer info: %w", err)
	}

	return buf, nil
}

// MapMemoryBuffer creates a local buffer mapped to the address space of the device specified by fd.
func MapMemoryBuffer(fd uintptr, offset int64, len int) ([]byte, error) {
	return sys.Mmap(int(fd), offset, len, sys.PROT_READ|sys.PROT_WRITE, sys.MAP_SHARED)
}

// UnmapMemoryBuffer removes the buffer that was previously mapped.
func UnmapMemoryBuffer(buf []byte) error {
	return sys.Munmap(buf)
}

// QueueBuffer enqueues a buffer in the device driver (empty for capturing, filled for video output)
// when using either memory map, user pointer, or DMA buffers. BufferInfo is returned with
// additional information about the queued buffer.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html#vidioc-qbuf
func QueueBuffer(fd uintptr, index uint32) (BufferInfo, error) {
	buf := BufferInfo{
		StreamType: BufTypeVideoCapture,
		Memory:     StreamMemoryTypeMMAP,
		Index:      index,
	}

	if err := Send(fd, VidiocQueueBuf, uintptr(unsafe.Pointer(&buf))); err != nil {
		return BufferInfo{}, fmt.Errorf("buffer queue: %w", err)
	}

	return buf, nil
}

// DequeueBuffer dequeues a buffer in the device driver, marking it as consumed by the application,
// when using either memory map, user pointer, or DMA buffers. BufferInfo is returned with
// additional information about the dequeued buffer.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html#vidioc-qbuf
func DequeueBuffer(fd uintptr) (BufferInfo, error) {
	buf := BufferInfo{
		StreamType: BufTypeVideoCapture,
		Memory:     StreamMemoryTypeMMAP,
	}

	if err := Send(fd, VidiocDequeueBuf, uintptr(unsafe.Pointer(&buf))); err != nil {
		return BufferInfo{}, fmt.Errorf("buffer dequeue: %w", err)

	}

	return buf, nil
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
