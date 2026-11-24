package v4l2

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// Streaming with Buffers
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html

// BufType specifies the type of V4L2 buffer used for data transfer.
// It determines whether the buffer is for video capture, output, or other purposes.
// This type corresponds to v4l2_buf_type in the kernel API.
//
// References:
//   - https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_buf_type#c.V4L.v4l2_buf_type
//   - https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L141
type BufType = uint32

// Buffer type constants define the direction and purpose of data flow.
const (
	// BufTypeVideoCapture is for video capture buffers (camera/tuner input).
	// Used when reading frames from a video source.
	BufTypeVideoCapture BufType = cV4L2_BUF_TYPE_VIDEO_CAPTURE

	// BufTypeVideoOutput is for video output buffers (display/encoder output).
	// Used when sending frames to a video sink.
	BufTypeVideoOutput BufType = cV4L2_BUF_TYPE_VIDEO_OUTPUT

	// BufTypeOverlay is for video overlay buffers.
	// Used for hardware overlay onto a display.
	BufTypeOverlay BufType = cV4L2_BUF_TYPE_VIDEO_OVERLAY

	// BufTypeVideoCaptureMPlane is for multi-planar video capture buffers.
	BufTypeVideoCaptureMPlane BufType = cV4L2_BUF_TYPE_VIDEO_CAPTURE_MPLANE

	// BufTypeVideoOutputMPlane is for multi-planar video output buffers.
	BufTypeVideoOutputMPlane BufType = cV4L2_BUF_TYPE_VIDEO_OUTPUT_MPLANE

	// BufTypeMetaCapture is for metadata capture buffers.
	BufTypeMetaCapture BufType = cV4L2_BUF_TYPE_META_CAPTURE

	// BufTypeMetaOutput is for metadata output buffers.
	BufTypeMetaOutput BufType = cV4L2_BUF_TYPE_META_OUTPUT
)

// IsMultiPlanar reports whether bufType uses the multi-planar V4L2 API.
func IsMultiPlanar(bufType BufType) bool {
	return bufType == BufTypeVideoCaptureMPlane || bufType == BufTypeVideoOutputMPlane
}

// Buffer capability flags returned in RequestBuffers.Capabilities and CreateBuffersResult.Capabilities.
const (
	BufCapSupportsMMAP         uint32 = cV4L2_BUF_CAP_SUPPORTS_MMAP
	BufCapSupportsUserPtr      uint32 = cV4L2_BUF_CAP_SUPPORTS_USERPTR
	BufCapSupportsDMABuf       uint32 = cV4L2_BUF_CAP_SUPPORTS_DMABUF
	BufCapSupportsRequests     uint32 = cV4L2_BUF_CAP_SUPPORTS_REQUESTS
	BufCapSupportsOrphanedBufs uint32 = cV4L2_BUF_CAP_SUPPORTS_ORPHANED_BUFS
)

// IOType specifies the I/O method for transferring data between userspace and the V4L2 device.
// Different I/O types offer various trade-offs between performance, complexity, and memory usage.
// This type corresponds to v4l2_memory in the kernel API.
//
// References:
//   - https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/mmap.html?highlight=v4l2_memory_mmap
//   - https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L188
type IOType = uint32

// I/O type constants define how buffers are allocated and accessed.
const (
	// IOTypeMMAP uses memory-mapped buffers for zero-copy data transfer.
	// This is the most efficient and commonly supported I/O method.
	// Buffers are allocated by the driver and mapped into userspace.
	IOTypeMMAP IOType = cV4L2_MEMORY_MMAP

	// IOTypeUserPtr uses application-allocated buffers.
	// The application provides memory pointers to the driver.
	// Requires contiguous physical memory in some cases.
	IOTypeUserPtr IOType = cV4L2_MEMORY_USERPTR

	// IOTypeOverlay writes directly to video memory for display.
	// Deprecated in newer kernels.
	IOTypeOverlay IOType = cV4L2_MEMORY_OVERLAY

	// IOTypeDMABuf uses DMA buffer sharing for inter-device communication.
	// Allows zero-copy sharing between V4L2 and other subsystems like DRM.
	// Requires DMABUF support in both driver and hardware.
	IOTypeDMABuf IOType = cV4L2_MEMORY_DMABUF
)

// BufFlag represents buffer status flags that indicate the state of a V4L2 buffer.
// These flags are set by the driver and/or application to track buffer lifecycle
// and provide metadata about the buffer contents.
type BufFlag = uint32

// Buffer flag constants indicate buffer state and frame metadata.
const (
	// BufFlagMapped indicates the buffer is memory-mapped into userspace.
	// Set when the buffer has been successfully mapped with mmap().
	BufFlagMapped BufFlag = cV4L2_BUF_FLAG_MAPPED

	// BufFlagQueued indicates the buffer is queued for I/O.
	// The buffer is owned by the driver and should not be accessed by the application.
	BufFlagQueued BufFlag = cV4L2_BUF_FLAG_QUEUED

	// BufFlagDone indicates the buffer has been processed and contains valid data.
	// For capture: buffer contains a captured frame.
	// For output: buffer has been displayed/transmitted.
	BufFlagDone BufFlag = cV4L2_BUF_FLAG_DONE

	// BufFlagKeyFrame indicates the buffer contains a keyframe (I-frame).
	// Used with compressed video streams.
	BufFlagKeyFrame BufFlag = cV4L2_BUF_FLAG_KEYFRAME

	// BufFlagPFrame indicates the buffer contains a P-frame (predicted frame).
	// Used with compressed video streams.
	BufFlagPFrame BufFlag = cV4L2_BUF_FLAG_PFRAME

	// BufFlagBFrame indicates the buffer contains a B-frame (bi-directional predicted frame).
	// Used with compressed video streams.
	BufFlagBFrame BufFlag = cV4L2_BUF_FLAG_BFRAME

	// BufFlagError indicates an error occurred during capture/output.
	// The buffer contents are likely invalid.
	BufFlagError               BufFlag = cV4L2_BUF_FLAG_ERROR
	BufFlagInRequest           BufFlag = cV4L2_BUF_FLAG_IN_REQUEST
	BufFlagTimeCode            BufFlag = cV4L2_BUF_FLAG_TIMECODE
	BufFlagM2MHoldCaptureBuf   BufFlag = cV4L2_BUF_FLAG_M2M_HOLD_CAPTURE_BUF
	BufFlagPrepared            BufFlag = cV4L2_BUF_FLAG_PREPARED
	BufFlagNoCacheInvalidate   BufFlag = cV4L2_BUF_FLAG_NO_CACHE_INVALIDATE
	BufFlagNoCacheClean        BufFlag = cV4L2_BUF_FLAG_NO_CACHE_CLEAN
	BufFlagTimestampMask       BufFlag = cV4L2_BUF_FLAG_TIMESTAMP_MASK
	BufFlagTimestampUnknown    BufFlag = cV4L2_BUF_FLAG_TIMESTAMP_UNKNOWN
	BufFlagTimestampMonotonic  BufFlag = cV4L2_BUF_FLAG_TIMESTAMP_MONOTONIC
	BufFlagTimestampCopy       BufFlag = cV4L2_BUF_FLAG_TIMESTAMP_COPY
	BufFlagTimestampSourceMask BufFlag = cV4L2_BUF_FLAG_TSTAMP_SRC_MASK
	BufFlagTimestampSourceEOF  BufFlag = cV4L2_BUF_FLAG_TSTAMP_SRC_EOF
	BufFlagTimestampSourceSOE  BufFlag = cV4L2_BUF_FLAG_TSTAMP_SRC_SOE
	BufFlagLast                BufFlag = cV4L2_BUF_FLAG_LAST
	BufFlagRequestFD           BufFlag = cV4L2_BUF_FLAG_REQUEST_FD
)

// CreateBuffersResult holds the driver response from VIDIOC_CREATE_BUFS (v4l2_create_buffers).
type CreateBuffersResult struct {
	Index        uint32
	Count        uint32
	Capabilities uint32
	Flags        uint32
	NumBuffers   uint32
}

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
type Buffer struct {
	Index     uint32
	Type      uint32
	BytesUsed uint32
	Flags     uint32
	Field     uint32
	Timestamp sys.Timeval
	Timecode  Timecode
	Sequence  uint32
	Memory    uint32
	Info      BufferInfo // union m (single-planar)
	Planes    []Plane    // multi-planar plane descriptors (when IsMultiPlanar(Type))
	Length    uint32
	Reserved2 uint32
	RequestFD int32
}

func getUintptr(m []byte) uintptr {
	if len(m) >= 8 {
		return uintptr(binary.NativeEndian.Uint64(m))
	}
	return uintptr(binary.NativeEndian.Uint32(m))
}

func putUintptr(m []byte, v uintptr) {
	if len(m) >= 8 {
		binary.NativeEndian.PutUint64(m, uint64(v))
		return
	}
	binary.NativeEndian.PutUint32(m, uint32(v))
}

func decodeBufferInfo(memory uint32, m []byte) BufferInfo {
	var info BufferInfo
	switch memory {
	case IOTypeUserPtr:
		info.UserPtr = getUintptr(m)
	case IOTypeDMABuf:
		info.FD = int32(binary.NativeEndian.Uint32(m[:4]))
	default:
		info.Offset = binary.NativeEndian.Uint32(m[:4])
	}
	return info
}

func decodePlaneInfo(memory uint32, m []byte) PlaneInfo {
	var info PlaneInfo
	switch memory {
	case IOTypeUserPtr:
		info.UserPtr = getUintptr(m)
	case IOTypeDMABuf:
		info.FD = int32(binary.NativeEndian.Uint32(m[:4]))
	default:
		info.MemOffset = binary.NativeEndian.Uint32(m[:4])
	}
	return info
}

func encodePlaneInfo(dst []byte, memory uint32, info PlaneInfo) {
	switch memory {
	case IOTypeUserPtr:
		putUintptr(dst, info.UserPtr)
	case IOTypeDMABuf:
		binary.NativeEndian.PutUint32(dst[:4], uint32(info.FD))
	default:
		binary.NativeEndian.PutUint32(dst[:4], info.MemOffset)
	}
}

// makeBuffer makes a Buffer value from cV4L2Buffer
func makeBuffer(v4l2Buf cV4L2Buffer) Buffer {
	return makeBufferWithPlanes(v4l2Buf, nil)
}

func makeBufferWithPlanes(v4l2Buf cV4L2Buffer, cPlanes []cV4L2Plane) Buffer {
	buf := Buffer{
		Index:     v4l2Buf.Index,
		Type:      v4l2Buf.Type,
		BytesUsed: v4l2Buf.Bytesused,
		Flags:     v4l2Buf.Flags,
		Field:     v4l2Buf.Field,
		Timestamp: sys.NsecToTimeval(v4l2Buf.Timestamp.Sec*1e9 + v4l2Buf.Timestamp.Usec*1e3),
		Timecode:  *(*Timecode)(unsafe.Pointer(&v4l2Buf.Timecode)),
		Sequence:  v4l2Buf.Sequence,
		Memory:    v4l2Buf.Memory,
		Length:    v4l2Buf.Length,
		Reserved2: v4l2Buf.Reserved2,
		RequestFD: v4l2Buf.Fd,
	}
	if IsMultiPlanar(v4l2Buf.Type) && len(cPlanes) > 0 {
		n := int(v4l2Buf.Length)
		if n > len(cPlanes) {
			n = len(cPlanes)
		}
		buf.Planes = make([]Plane, n)
		for i := 0; i < n; i++ {
			buf.Planes[i] = Plane{
				BytesUsed:  cPlanes[i].Bytesused,
				Length:     cPlanes[i].Length,
				Info:       decodePlaneInfo(v4l2Buf.Memory, cPlanes[i].M[:]),
				DataOffset: cPlanes[i].Offset,
			}
		}
		if len(buf.Planes) > 0 {
			buf.Info.Planes = &buf.Planes[0]
		}
	} else {
		buf.Info = decodeBufferInfo(v4l2Buf.Memory, v4l2Buf.M[:])
	}
	return buf
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
type Plane struct {
	BytesUsed  uint32
	Length     uint32
	Info       PlaneInfo // union m
	DataOffset uint32
}

// PlaneInfo represents the combination of type
// of memory stream that can be serviced for the
// associated plane.
type PlaneInfo struct {
	MemOffset uint32
	UserPtr   uintptr
	FD        int32
}

// StreamOn requests streaming to be turned on for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-streamon.html
func StreamOn(dev StreamingDevice) error {
	bufType := dev.BufferType()
	if err := send(dev.Fd(), cVIDIOC_STREAMON, unsafe.Pointer(&bufType)); err != nil {
		return fmt.Errorf("stream on: %w", err)
	}
	return nil
}

// StreamOff requests streaming to be turned off for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-streamon.html
func StreamOff(dev StreamingDevice) error {
	bufType := dev.BufferType()
	if err := send(dev.Fd(), cVIDIOC_STREAMOFF, unsafe.Pointer(&bufType)); err != nil {
		return fmt.Errorf("stream off: %w", err)
	}
	return nil
}

// InitBuffers sends buffer allocation request (VIDIOC_REQBUFS) to initialize buffer IO
// for video capture or video output when using either mem map, user pointer, or DMA buffers.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-reqbufs.html#vidioc-reqbufs
func InitBuffers(dev StreamingDevice) (RequestBuffers, error) {
	if dev.MemIOType() != IOTypeMMAP && dev.MemIOType() != IOTypeUserPtr && dev.MemIOType() != IOTypeDMABuf {
		return RequestBuffers{}, fmt.Errorf("request buffers: %w", ErrorUnsupported)
	}
	var req cV4L2RequestBuffers
	req.Count = dev.BufferCount()
	req.Type = dev.BufferType()
	req.Memory = dev.MemIOType()

	if err := send(dev.Fd(), cVIDIOC_REQBUFS, unsafe.Pointer(&req)); err != nil {
		return RequestBuffers{}, fmt.Errorf("request buffers: %w: type not supported", err)
	}

	return *(*RequestBuffers)(unsafe.Pointer(&req)), nil
}

// ResetBuffers allocates a buffer of size 0 VIDIOC_REQBUFS(0) to free (or orphan) all
// buffers. Useful when shuttingdown the stream.
// See https://linuxtv.org/downloads/v4l-dvb-apis-new/userspace-api/v4l/vidioc-reqbufs.html
func ResetBuffers(dev StreamingDevice) (RequestBuffers, error) {
	if dev.MemIOType() != IOTypeMMAP && dev.MemIOType() != IOTypeUserPtr && dev.MemIOType() != IOTypeDMABuf {
		return RequestBuffers{}, fmt.Errorf("reset buffers: %w", ErrorUnsupported)
	}
	var req cV4L2RequestBuffers
	req.Count = 0
	req.Type = dev.BufferType()
	req.Memory = dev.MemIOType()

	if err := send(dev.Fd(), cVIDIOC_REQBUFS, unsafe.Pointer(&req)); err != nil {
		return RequestBuffers{}, fmt.Errorf("reset buffers VIDIOC_REQBUFS(0): %w", err)
	}

	return *(*RequestBuffers)(unsafe.Pointer(&req)), nil
}

// GetBuffer retrieves buffer info for allocated buffers at provided index.
// This call should take place after buffers are allocated with RequestBuffers (for mmap for instance).
func GetBuffer(dev StreamingDevice, index uint32) (Buffer, error) {
	var v4l2Buf cV4L2Buffer
	v4l2Buf.Type = dev.BufferType()
	v4l2Buf.Memory = dev.MemIOType()
	v4l2Buf.Index = index

	if err := send(dev.Fd(), cVIDIOC_QUERYBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
		return Buffer{}, fmt.Errorf("query buffer: type not supported: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// mapMemoryBuffer creates a local buffer mapped to the address space of the device specified by fd.
func mapMemoryBuffer(fd uintptr, offset int64, len int) ([]byte, error) {
	data, err := sys.Mmap(int(fd), offset, len, sys.PROT_READ|sys.PROT_WRITE, sys.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("map memory buffer: %w", err)
	}
	return data, nil
}

// MapMemoryBuffers creates mapped memory buffers for specified buffer count of device.
func MapMemoryBuffers(dev StreamingDevice) ([][]byte, error) {
	bufCount := int(dev.BufferCount())
	buffers := make([][]byte, bufCount)
	for i := 0; i < bufCount; i++ {
		buffer, err := GetBuffer(dev, uint32(i))
		if err != nil {
			return nil, fmt.Errorf("mapped buffers: %w", err)
		}

		// TODO check buffer flags for errors etc

		offset := buffer.Info.Offset
		length := buffer.Length
		mappedBuf, err := mapMemoryBuffer(dev.Fd(), int64(offset), int(length))
		if err != nil {
			return nil, fmt.Errorf("mapped buffers: %w", err)
		}
		buffers[i] = mappedBuf
	}
	return buffers, nil
}

// AllocateUserBuffers creates application-managed buffers for USERPTR streaming.
// Each buffer is allocated with the specified size (typically pixFormat.SizeImage).
func AllocateUserBuffers(count int, size uint32) [][]byte {
	buffers := make([][]byte, count)
	for i := range buffers {
		buffers[i] = make([]byte, size)
	}
	return buffers
}

// unmapMemoryBuffer removes the buffer that was previously mapped.
func unmapMemoryBuffer(buf []byte) error {
	if err := sys.Munmap(buf); err != nil {
		return fmt.Errorf("unmap memory buffer: %w", err)
	}
	return nil
}

// UnmapMemoryBuffers unmaps all mapped memory buffer for device
func UnmapMemoryBuffers(dev StreamingDevice) error {
	if dev.Buffers() == nil {
		return fmt.Errorf("unmap buffers: uninitialized buffers")
	}
	for i := 0; i < len(dev.Buffers()); i++ {
		if err := unmapMemoryBuffer(dev.Buffers()[i]); err != nil {
			return fmt.Errorf("unmap buffers: %w", err)
		}
	}
	return nil
}

// QueueBuffer enqueues a buffer in the device driver (as empty for capturing, or filled for video output)
// when using either memory map, user pointer, or DMA buffers. Buffer is returned with
// additional information about the queued buffer.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html#vidioc-qbuf
func QueueBuffer(fd uintptr, ioType IOType, bufType BufType, index uint32) (Buffer, error) {
	var v4l2Buf cV4L2Buffer
	v4l2Buf.Type = bufType
	v4l2Buf.Memory = ioType
	v4l2Buf.Index = index

	if err := send(fd, cVIDIOC_QBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
		return Buffer{}, fmt.Errorf("buffer queue: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// QueueBufferUserPtr enqueues a user-pointer buffer in the device driver.
// The application provides the buffer address and length. Used with IOTypeUserPtr streaming.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html#vidioc-qbuf
func QueueBufferUserPtr(fd uintptr, bufType BufType, index uint32, ptr uintptr, length uint32) (Buffer, error) {
	var v4l2Buf cV4L2Buffer
	v4l2Buf.Type = bufType
	v4l2Buf.Memory = IOTypeUserPtr
	v4l2Buf.Index = index
	v4l2Buf.Length = length
	*(*uintptr)(unsafe.Pointer(&v4l2Buf.M[0])) = ptr

	if err := send(fd, cVIDIOC_QBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
		return Buffer{}, fmt.Errorf("buffer queue userptr: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// QueueBufferDMABuf enqueues a DMA-BUF buffer in the device driver.
// The application provides the DMA-BUF file descriptor and buffer length.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html#vidioc-qbuf
func QueueBufferDMABuf(fd uintptr, bufType BufType, index uint32, dmabufFD int32, length uint32) (Buffer, error) {
	var v4l2Buf cV4L2Buffer
	v4l2Buf.Type = bufType
	v4l2Buf.Memory = IOTypeDMABuf
	v4l2Buf.Index = index
	v4l2Buf.Length = length
	*(*int32)(unsafe.Pointer(&v4l2Buf.M[0])) = dmabufFD

	if err := send(fd, cVIDIOC_QBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
		return Buffer{}, fmt.Errorf("buffer queue dmabuf: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// ExportBufferPlane exports a specific plane of a V4L2 buffer as a DMA-BUF file descriptor (VIDIOC_EXPBUF).
// For single-planar buffers, pass plane = 0.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-expbuf.html
func ExportBufferPlane(devFD uintptr, bufType BufType, index, plane, flags uint32) (int32, error) {
	var exp cV4L2ExportBuffer
	exp.Type = bufType
	exp.Index = index
	exp.Plane = plane
	exp.Flags = flags

	if err := send(devFD, cVIDIOC_EXPBUF, unsafe.Pointer(&exp)); err != nil {
		return 0, fmt.Errorf("export buffer: %w", err)
	}
	return exp.Fd, nil
}

// ExportBuffer exports a V4L2 buffer as a DMA-BUF file descriptor.
// The buffer must have been allocated with MMAP. The returned fd can be
// passed to other subsystems (GPU, DRM, other V4L2 devices).
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-expbuf.html
func ExportBuffer(devFD uintptr, bufType BufType, index uint32, flags uint32) (int32, error) {
	return ExportBufferPlane(devFD, bufType, index, 0, flags)
}

// DequeueBuffer dequeues a buffer in the device driver, marking it as consumed by the application,
// when using either memory map, user pointer, or DMA buffers. Buffer is returned with
// additional information about the dequeued buffer.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html#vidioc-qbuf
func DequeueBuffer(fd uintptr, ioType IOType, bufType BufType) (Buffer, error) {
	var v4l2Buf cV4L2Buffer
	v4l2Buf.Type = bufType
	v4l2Buf.Memory = ioType

	err := send(fd, cVIDIOC_DQBUF, unsafe.Pointer(&v4l2Buf))
	if err != nil {
		return Buffer{}, fmt.Errorf("buffer dequeue: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// DequeueBufferMPlane dequeues a multi-planar buffer (VIDIOC_DQBUF) with space for numPlanes planes.
// numPlanes must be in [1, MaxPlanes].
func DequeueBufferMPlane(fd uintptr, ioType IOType, bufType BufType, numPlanes uint32) (Buffer, error) {
	if !IsMultiPlanar(bufType) {
		return Buffer{}, fmt.Errorf("buffer dequeue mplane: buffer type %d is not multi-planar", bufType)
	}
	if numPlanes == 0 || numPlanes > MaxPlanes {
		return Buffer{}, fmt.Errorf("buffer dequeue mplane: invalid plane count %d", numPlanes)
	}

	var cPlanesBuf [MaxPlanes]cV4L2Plane
	cPlanes := cPlanesBuf[:numPlanes]
	var pinner runtime.Pinner
	pinner.Pin(&cPlanes[0])
	defer func() {
		pinner.Unpin()
		runtime.KeepAlive(&cPlanesBuf)
	}()

	var v4l2Buf cV4L2Buffer
	v4l2Buf.Type = bufType
	v4l2Buf.Memory = ioType
	v4l2Buf.Length = numPlanes
	putUintptr(v4l2Buf.M[:], uintptr(unsafe.Pointer(&cPlanes[0])))

	if err := send(fd, cVIDIOC_DQBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
		return Buffer{}, fmt.Errorf("buffer dequeue mplane: %w", err)
	}

	return makeBufferWithPlanes(v4l2Buf, cPlanes), nil
}

// StreamOnType starts streaming for the specified buffer type directly on fd (VIDIOC_STREAMON).
func StreamOnType(fd uintptr, bufType BufType) error {
	if err := send(fd, cVIDIOC_STREAMON, unsafe.Pointer(&bufType)); err != nil {
		return fmt.Errorf("stream on: %w", err)
	}
	return nil
}

// StreamOffType stops streaming for the specified buffer type directly on fd (VIDIOC_STREAMOFF).
func StreamOffType(fd uintptr, bufType BufType) error {
	if err := send(fd, cVIDIOC_STREAMOFF, unsafe.Pointer(&bufType)); err != nil {
		return fmt.Errorf("stream off: %w", err)
	}
	return nil
}

// RequestBuffersForType requests count buffers for bufType and ioType directly on fd (VIDIOC_REQBUFS).
// Pass count = 0 to free/orphan previously allocated buffers.
func RequestBuffersForType(fd uintptr, bufType BufType, ioType IOType, count uint32) (RequestBuffers, error) {
	var req cV4L2RequestBuffers
	req.Count = count
	req.Type = bufType
	req.Memory = ioType

	if err := send(fd, cVIDIOC_REQBUFS, unsafe.Pointer(&req)); err != nil {
		return RequestBuffers{}, fmt.Errorf("request buffers: %w", err)
	}

	return *(*RequestBuffers)(unsafe.Pointer(&req)), nil
}

// CreateBuffers dynamically allocates additional buffers for bufType and ioType (VIDIOC_CREATE_BUFS).
// For single-planar buffer types, provide pixFmt; for multi-planar buffer types, provide pixFmtMPlane.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-create-bufs.html
func CreateBuffers(fd uintptr, count uint32, ioType IOType, bufType BufType, pixFmt *PixFormat, pixFmtMPlane *PixFormatMPlane) (CreateBuffersResult, error) {
	var cb cV4L2CreateBuffers
	cb.Count = count
	cb.Memory = ioType
	cb.Format.Type = bufType

	if IsMultiPlanar(bufType) {
		if pixFmtMPlane == nil {
			return CreateBuffersResult{}, fmt.Errorf("create buffers: multi-planar buffer type %d requires pixFmtMPlane", bufType)
		}
		*(*cV4L2PixFormatMPlane)(unsafe.Pointer(&cb.Format.Fmt[0])) = pixFormatMPlaneToC(*pixFmtMPlane)
	} else {
		if pixFmt == nil {
			return CreateBuffersResult{}, fmt.Errorf("create buffers: single-planar buffer type %d requires pixFmt", bufType)
		}
		*(*cV4L2PixFormat)(unsafe.Pointer(&cb.Format.Fmt[0])) = pixFormatToC(*pixFmt)
	}

	if err := send(fd, cVIDIOC_CREATE_BUFS, unsafe.Pointer(&cb)); err != nil {
		return CreateBuffersResult{}, fmt.Errorf("create buffers: %w", err)
	}

	return CreateBuffersResult{
		Index:        cb.Index,
		Count:        cb.Count,
		Capabilities: cb.Capabilities,
		Flags:        cb.Flags,
		NumBuffers:   cb.Num_buffers,
	}, nil
}

// QueryBuffer queries buffer metadata at index for bufType and ioType (VIDIOC_QUERYBUF).
// For multi-planar buffer types, numPlanes must be in [1, MaxPlanes]; for single-planar types it is ignored.
func QueryBuffer(fd uintptr, bufType BufType, ioType IOType, index uint32, numPlanes uint32) (Buffer, error) {
	var v4l2Buf cV4L2Buffer
	v4l2Buf.Type = bufType
	v4l2Buf.Memory = ioType
	v4l2Buf.Index = index

	if IsMultiPlanar(bufType) {
		if numPlanes == 0 || numPlanes > MaxPlanes {
			return Buffer{}, fmt.Errorf("query buffer: invalid plane count %d", numPlanes)
		}
		var cPlanesBuf [MaxPlanes]cV4L2Plane
		cPlanes := cPlanesBuf[:numPlanes]
		var pinner runtime.Pinner
		pinner.Pin(&cPlanes[0])
		defer func() {
			pinner.Unpin()
			runtime.KeepAlive(&cPlanesBuf)
		}()

		v4l2Buf.Length = numPlanes
		putUintptr(v4l2Buf.M[:], uintptr(unsafe.Pointer(&cPlanes[0])))

		if err := send(fd, cVIDIOC_QUERYBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
			return Buffer{}, fmt.Errorf("query buffer: %w", err)
		}
		return makeBufferWithPlanes(v4l2Buf, cPlanes), nil
	}

	if err := send(fd, cVIDIOC_QUERYBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
		return Buffer{}, fmt.Errorf("query buffer: %w", err)
	}
	return makeBuffer(v4l2Buf), nil
}

// QueueBufferWith enqueues a buffer with full caller control over v4l2_buffer metadata and planes (VIDIOC_QBUF).
// Supports both single-planar and multi-planar buffer types across MMAP, USERPTR, and DMABUF memory modes.
func QueueBufferWith(fd uintptr, buf Buffer) (Buffer, error) {
	return submitBufferIoctl(fd, cVIDIOC_QBUF, "buffer queue", buf)
}

// PrepareBuffer prepares a buffer for I/O off the hot path (VIDIOC_PREPARE_BUF), allowing
// the driver to validate parameters and perform cache synchronization before VIDIOC_QBUF.
// Supports both single-planar and multi-planar buffer types across MMAP, USERPTR, and DMABUF memory modes.
func PrepareBuffer(fd uintptr, buf Buffer) (Buffer, error) {
	return submitBufferIoctl(fd, cVIDIOC_PREPARE_BUF, "buffer prepare", buf)
}

func submitBufferIoctl(fd, req uintptr, op string, buf Buffer) (Buffer, error) {
	v4l2Buf := cV4L2Buffer{
		Index:     buf.Index,
		Type:      buf.Type,
		Bytesused: buf.BytesUsed,
		Flags:     buf.Flags,
		Field:     buf.Field,
		Timestamp: cTimeval{Sec: int64(buf.Timestamp.Sec), Usec: int64(buf.Timestamp.Usec)},
		Timecode:  *(*cV4L2Timecode)(unsafe.Pointer(&buf.Timecode)),
		Memory:    buf.Memory,
		Length:    buf.Length,
		Reserved2: buf.Reserved2,
		Fd:        buf.RequestFD,
	}

	if IsMultiPlanar(buf.Type) {
		if len(buf.Planes) == 0 || len(buf.Planes) > MaxPlanes {
			return Buffer{}, fmt.Errorf("%s: invalid plane count %d", op, len(buf.Planes))
		}
		var cPlanesBuf [MaxPlanes]cV4L2Plane
		cPlanes := cPlanesBuf[:len(buf.Planes)]
		for i, p := range buf.Planes {
			cPlanes[i] = cV4L2Plane{
				Bytesused: p.BytesUsed,
				Length:    p.Length,
				Offset:    p.DataOffset,
			}
			encodePlaneInfo(cPlanes[i].M[:], buf.Memory, p.Info)
		}

		var pinner runtime.Pinner
		pinner.Pin(&cPlanes[0])
		defer func() {
			pinner.Unpin()
			runtime.KeepAlive(&cPlanesBuf)
		}()

		v4l2Buf.Length = uint32(len(cPlanes))
		putUintptr(v4l2Buf.M[:], uintptr(unsafe.Pointer(&cPlanes[0])))

		if err := send(fd, req, unsafe.Pointer(&v4l2Buf)); err != nil {
			return Buffer{}, fmt.Errorf("%s: %w", op, err)
		}
		return makeBufferWithPlanes(v4l2Buf, cPlanes), nil
	}

	switch buf.Memory {
	case IOTypeUserPtr:
		putUintptr(v4l2Buf.M[:], buf.Info.UserPtr)
	case IOTypeDMABuf:
		binary.NativeEndian.PutUint32(v4l2Buf.M[:4], uint32(buf.Info.FD))
	default:
		binary.NativeEndian.PutUint32(v4l2Buf.M[:4], buf.Info.Offset)
	}

	if err := send(fd, req, unsafe.Pointer(&v4l2Buf)); err != nil {
		return Buffer{}, fmt.Errorf("%s: %w", op, err)
	}
	return makeBuffer(v4l2Buf), nil
}

// RemoveBuffers deletes count buffers starting at index for bufType (VIDIOC_REMOVE_BUFS, Linux 6.10+).
// Complements CreateBuffers for dynamically shrinking a buffer queue without stopping the stream.
func RemoveBuffers(fd uintptr, bufType BufType, index, count uint32) error {
	req := cV4L2RemoveBuffers{
		Index: index,
		Count: count,
		Type:  bufType,
	}
	if err := send(fd, cVIDIOC_REMOVE_BUFS, unsafe.Pointer(&req)); err != nil {
		return fmt.Errorf("remove buffers (index=%d count=%d): %w", index, count, err)
	}
	return nil
}

// MapMemoryBuffer maps a V4L2 MMAP buffer plane at offset with length into process memory.
func MapMemoryBuffer(fd uintptr, offset int64, length int) ([]byte, error) {
	return mapMemoryBuffer(fd, offset, length)
}

// UnmapMemoryBuffer unmaps a buffer previously mapped by MapMemoryBuffer.
func UnmapMemoryBuffer(buf []byte) error {
	return unmapMemoryBuffer(buf)
}
