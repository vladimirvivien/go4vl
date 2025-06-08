package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// Streaming with Buffers
// This section defines types and functions related to V4L2 buffer management and streaming I/O.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html

// BufType is a type alias for uint32, representing the type of a V4L2 buffer.
// It specifies the direction of data flow and the kind of data the buffer handles.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html?highlight=v4l2_buf_type#c.V4L.v4l2_buf_type
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L141
type BufType = uint32

// Buffer Type Constants
const (
	// BufTypeVideoCapture is for buffers that capture video data from a device.
	BufTypeVideoCapture BufType = C.V4L2_BUF_TYPE_VIDEO_CAPTURE
	// BufTypeVideoOutput is for buffers that output video data to a device.
	BufTypeVideoOutput BufType = C.V4L2_BUF_TYPE_VIDEO_OUTPUT
	// BufTypeOverlay is for video overlay buffers.
	BufTypeOverlay BufType = C.V4L2_BUF_TYPE_VIDEO_OVERLAY
	// Other buffer types like VBI, Sliced VBI, Meta, etc., are defined in the C headers
	// but not explicitly listed here. They can be added as needed.
)

// IOType is a type alias for uint32, representing the memory I/O method used for V4L2 buffers.
// It specifies how buffer memory is allocated and accessed (e.g., memory mapping, user pointers).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/mmap.html?highlight=v4l2_memory_mmap
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L188
type IOType = uint32

// Memory I/O Type Constants
const (
	// IOTypeMMAP indicates that buffers are allocated by the driver and memory-mapped into user space.
	IOTypeMMAP IOType = C.V4L2_MEMORY_MMAP
	// IOTypeUserPtr indicates that buffers are allocated by the application (user pointers).
	IOTypeUserPtr IOType = C.V4L2_MEMORY_USERPTR
	// IOTypeOverlay indicates that buffers are used for video overlay.
	IOTypeOverlay IOType = C.V4L2_MEMORY_OVERLAY // Typically used with VIDIOC_S_FMT for overlay target.
	// IOTypeDMABuf indicates that buffers are shared via DMA-BUF file descriptors.
	IOTypeDMABuf IOType = C.V4L2_MEMORY_DMABUF
)

// BufFlag is a type alias for uint32, representing flags that describe the state and properties of a V4L2 buffer.
// These flags are used in the `Flags` field of the `Buffer` struct.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html#buffer-flags
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1051
type BufFlag = uint32

// Buffer Flag Constants
const (
	BufFlagMapped              BufFlag = C.V4L2_BUF_FLAG_MAPPED              // Buffer is mapped.
	BufFlagQueued              BufFlag = C.V4L2_BUF_FLAG_QUEUED              // Buffer is currently in an incoming queue.
	BufFlagDone                BufFlag = C.V4L2_BUF_FLAG_DONE                // Buffer is currently in an outgoing queue.
	BufFlagKeyFrame            BufFlag = C.V4L2_BUF_FLAG_KEYFRAME            // Buffer is a keyframe (for coded formats).
	BufFlagPFrame              BufFlag = C.V4L2_BUF_FLAG_PFRAME              // Buffer is a P-frame (for coded formats).
	BufFlagBFrame              BufFlag = C.V4L2_BUF_FLAG_BFRAME              // Buffer is a B-frame (for coded formats).
	BufFlagError               BufFlag = C.V4L2_BUF_FLAG_ERROR               // An error occurred during buffer processing.
	BufFlagInRequest           BufFlag = C.V4L2_BUF_FLAG_IN_REQUEST          // Buffer is part of a submitted request.
	BufFlagTimeCode            BufFlag = C.V4L2_BUF_FLAG_TIMECODE            // Timecode field is valid.
	BufFlagM2MHoldCaptureBuf   BufFlag = C.V4L2_BUF_FLAG_M2M_HOLD_CAPTURE_BUF // For memory-to-memory devices: hold the capture buffer.
	BufFlagPrepared            BufFlag = C.V4L2_BUF_FLAG_PREPARED            // Buffer has been prepared for I/O.
	BufFlagNoCacheInvalidate   BufFlag = C.V4L2_BUF_FLAG_NO_CACHE_INVALIDATE // Do not invalidate cache before I/O.
	BufFlagNoCacheClean        BufFlag = C.V4L2_BUF_FLAG_NO_CACHE_CLEAN        // Do not clean cache after I/O.
	BufFlagTimestampMask       BufFlag = C.V4L2_BUF_FLAG_TIMESTAMP_MASK       // Mask for timestamp type.
	BufFlagTimestampUnknown    BufFlag = C.V4L2_BUF_FLAG_TIMESTAMP_UNKNOWN    // Timestamp type is unknown.
	BufFlagTimestampMonotonic  BufFlag = C.V4L2_BUF_FLAG_TIMESTAMP_MONOTONIC  // Timestamp is monotonic.
	BufFlagTimestampCopy       BufFlag = C.V4L2_BUF_FLAG_TIMESTAMP_COPY       // Timestamp is copied from elsewhere.
	BufFlagTimestampSourceMask BufFlag = C.V4L2_BUF_FLAG_TSTAMP_SRC_MASK       // Mask for timestamp source.
	BufFlagTimestampSourceEOF  BufFlag = C.V4L2_BUF_FLAG_TSTAMP_SRC_EOF       // Timestamp source is end-of-frame.
	BufFlagTimestampSourceSOE  BufFlag = C.V4L2_BUF_FLAG_TSTAMP_SRC_SOE       // Timestamp source is start-of-exposure.
	BufFlagLast                BufFlag = C.V4L2_BUF_FLAG_LAST                // Last buffer produced by the hardware.
	BufFlagRequestFD           BufFlag = C.V4L2_BUF_FLAG_REQUEST_FD           // Buffer contains a request file descriptor (for request API).
)

// TODO implement vl42_create_buffers (VIDIOC_CREATE_BUFS) if needed.

// RequestBuffers is used to request buffer allocation for streaming I/O.
// It corresponds to the `v4l2_requestbuffers` struct in the Linux kernel.
// This is a fundamental step in initializing memory-mapped, user pointer, or DMA buffer streaming.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-reqbufs.html#c.V4L.v4l2_requestbuffers
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L949
type RequestBuffers struct {
	// Count is the number of buffers requested from or granted by the driver.
	Count uint32
	// StreamType specifies the type of stream (e.g., video capture). See BufType constants.
	// In Go, this would ideally be of type BufType for better type safety.
	StreamType uint32
	// Memory specifies the I/O method (e.g., MMAP, USERPTR). See IOType constants.
	// In Go, this would ideally be of type IOType for better type safety.
	Memory uint32
	// Capabilities reports buffer capabilities (driver-set, read-only).
	Capabilities uint32
	// _ is a reserved field in C struct for future use.
	_ [1]uint32
}

// Buffer holds information about a single V4L2 buffer, used for queuing and dequeuing.
// It corresponds to the `v4l2_buffer` struct in the Linux kernel.
// This struct is used with ioctls like VIDIOC_QUERYBUF, VIDIOC_QBUF, and VIDIOC_DQBUF.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html#c.V4L.v4l2_buffer
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1037
type Buffer struct {
	// Index is the buffer's zero-based index in the driver's list of buffers.
	Index uint32
	// Type is the buffer type (e.g., video capture). See BufType constants.
	// In Go, this would ideally be of type BufType for better type safety.
	Type uint32
	// BytesUsed is the number of bytes occupied by data in the buffer.
	BytesUsed uint32
	// Flags indicate the buffer's state and properties. See BufFlag constants.
	// In Go, this would ideally be of type BufFlag for better type safety.
	Flags uint32
	// Field specifies the field order for interlaced video. See FieldType constants.
	// In Go, this would ideally be of type FieldType for better type safety.
	Field uint32
	// Timestamp of when the first data byte was captured or output.
	Timestamp sys.Timeval
	// Timecode stores frame timecode information. (Type Timecode is defined elsewhere)
	Timecode Timecode
	// Sequence number of the frame, set by the driver.
	Sequence uint32
	// Memory specifies the I/O method for this buffer. See IOType constants.
	// In Go, this would ideally be of type IOType for better type safety.
	Memory uint32
	// Info contains memory-specific details (offset for MMAP, user pointer, planes, or FD). This is a union in C.
	Info BufferInfo // Represents union `m`
	// Length is the size of the buffer in bytes.
	Length uint32
	// Reserved2 is a reserved field for future use.
	Reserved2 uint32
	// RequestFD is a file descriptor for the request API, used with BufFlagRequestFD.
	// This field is part of an anonymous union in the C struct.
	RequestFD int32
}

// makeBuffer makes a Buffer value from C.struct_v4l2_buffer
func makeBuffer(v4l2Buf C.struct_v4l2_buffer) Buffer {
	return Buffer{
		Index:     uint32(v4l2Buf.index),
		Type:      uint32(v4l2Buf._type),
		BytesUsed: uint32(v4l2Buf.bytesused),
		Flags:     uint32(v4l2Buf.flags),
		Field:     uint32(v4l2Buf.field),
		Timestamp: *(*sys.Timeval)(unsafe.Pointer(&v4l2Buf.timestamp)),
		Timecode:  *(*Timecode)(unsafe.Pointer(&v4l2Buf.timecode)),
		Sequence:  uint32(v4l2Buf.sequence),
		Memory:    uint32(v4l2Buf.memory),
		Info:      *(*BufferInfo)(unsafe.Pointer(&v4l2Buf.m[0])),
		Length:    uint32(v4l2Buf.length),
		Reserved2: uint32(v4l2Buf.reserved2),
		RequestFD: *(*int32)(unsafe.Pointer(&v4l2Buf.anon0[0])),
	}
}

// BufferInfo represents the union `m` within the `v4l2_buffer` struct.
// It holds different types of information depending on the buffer's memory I/O type (Memory field in Buffer struct).
//   - For IOTypeMMAP: Offset contains the offset of the buffer from the start of the device memory.
//   - For IOTypeUserPtr: UserPtr holds the user-space pointer to the buffer.
//   - For multi-planar buffers: Planes points to an array of Plane structs.
//   - For IOTypeDMABuf: FD holds the file descriptor for a DMA buffer.
type BufferInfo struct {
	// Offset is the offset from the start of device memory, for IOTypeMMAP.
	Offset uint32
	// UserPtr is the user-space pointer to the buffer, for IOTypeUserPtr.
	UserPtr uintptr
	// Planes points to an array of Plane structs, for multi-planar buffers.
	Planes *Plane
	// FD is the file descriptor for a DMA buffer, for IOTypeDMABuf.
	FD int32
}

// Plane describes a single plane in a multi-planar V4L2 buffer.
// It corresponds to the `v4l2_plane` struct in the Linux kernel.
// Multi-planar formats (e.g., NV12, YUV420P) use multiple planes to store different components of a frame.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/buffer.html#c.V4L.v4l2_plane
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L990
type Plane struct {
	// BytesUsed is the number of bytes occupied by data in this plane.
	BytesUsed uint32
	// Length is the size of this plane in bytes.
	Length uint32
	// Info contains memory-specific details for this plane (offset for MMAP, user pointer, or FD for DMABUF). This is a union in C.
	Info PlaneInfo // Represents union `m`
	// DataOffset is an offset to the data from the start of the plane (for formats with padding/header).
	DataOffset uint32
	// reserved fields in C struct
	// _ [11]uint32 // Implicitly handled by CGo struct mapping
}

// PlaneInfo represents the union `m` within the `v4l2_plane` struct.
// It holds different types of information depending on the plane's memory I/O type.
//   - For IOTypeMMAP: MemOffset contains the offset of the plane from the start of the memory mapped area for this buffer.
//   - For IOTypeUserPtr: UserPtr holds the user-space pointer to the plane's data.
//   - For IOTypeDMABuf: FD holds the file descriptor for the plane's DMA buffer.
type PlaneInfo struct {
	// MemOffset is the offset from the start of the buffer's mmap area, for IOTypeMMAP.
	MemOffset uint32
	// UserPtr is the user-space pointer to plane data, for IOTypeUserPtr.
	UserPtr uintptr
	// FD is the file descriptor for the plane's DMA buffer, for IOTypeDMABuf.
	FD int32
}

// StreamOn starts or resumes streaming for the specified device.
// It takes a StreamingDevice (which provides Fd() and BufferType() methods).
// This function calls the VIDIOC_STREAMON ioctl.
//
// Returns an error if the ioctl call fails.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-streamon.html
func StreamOn(dev StreamingDevice) error {
	bufType := dev.BufferType()
	if err := send(dev.Fd(), C.VIDIOC_STREAMON, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return fmt.Errorf("stream on: %w", err)
	}
	return nil
}

// StreamOff stops or pauses streaming for the specified device.
// It takes a StreamingDevice.
// This function calls the VIDIOC_STREAMOFF ioctl.
//
// Returns an error if the ioctl call fails.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-streamon.html
func StreamOff(dev StreamingDevice) error {
	bufType := dev.BufferType()
	if err := send(dev.Fd(), C.VIDIOC_STREAMOFF, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return fmt.Errorf("stream off: %w", err)
	}
	return nil
}

// InitBuffers requests the driver to allocate buffers for streaming I/O.
// It takes a StreamingDevice, which provides details like buffer count, type, and memory I/O method.
// This function is used for memory-mapped (MMAP) or DMA buffer (DMABUF) I/O.
// It calls the VIDIOC_REQBUFS ioctl.
//
// Returns a RequestBuffers struct (which may contain updated count or capabilities from the driver)
// and an error if the I/O type is not supported or the ioctl call fails.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-reqbufs.html
func InitBuffers(dev StreamingDevice) (RequestBuffers, error) {
	memIOType := dev.MemIOType()
	if memIOType != IOTypeMMAP && memIOType != IOTypeDMABuf {
		return RequestBuffers{}, fmt.Errorf("request buffers: unsupported memory I/O type %d", memIOType)
	}
	var req C.struct_v4l2_requestbuffers
	req.count = C.uint(dev.BufferCount())
	req._type = C.uint(dev.BufferType())
	req.memory = C.uint(memIOType)

	if err := send(dev.Fd(), C.VIDIOC_REQBUFS, uintptr(unsafe.Pointer(&req))); err != nil {
		return RequestBuffers{}, fmt.Errorf("request buffers: VIDIOC_REQBUFS failed: %w", err)
	}

	return *(*RequestBuffers)(unsafe.Pointer(&req)), nil
}

// ResetBuffers requests the driver to free all previously allocated buffers.
// This is done by calling VIDIOC_REQBUFS with a count of 0.
// It's useful for deinitializing streaming or changing buffer parameters.
// It takes a StreamingDevice.
//
// Returns a RequestBuffers struct (which should reflect zero buffers) and an error if the call fails.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-reqbufs.html#releasing-buffers
func ResetBuffers(dev StreamingDevice) (RequestBuffers, error) {
	memIOType := dev.MemIOType()
	if memIOType != IOTypeMMAP && memIOType != IOTypeDMABuf {
		return RequestBuffers{}, fmt.Errorf("reset buffers: unsupported memory I/O type %d", memIOType)
	}
	var req C.struct_v4l2_requestbuffers
	req.count = C.uint(0) // Requesting zero buffers frees allocated ones.
	req._type = C.uint(dev.BufferType())
	req.memory = C.uint(memIOType)

	if err := send(dev.Fd(), C.VIDIOC_REQBUFS, uintptr(unsafe.Pointer(&req))); err != nil {
		return RequestBuffers{}, fmt.Errorf("reset buffers: VIDIOC_REQBUFS(0) failed: %w", err)
	}

	return *(*RequestBuffers)(unsafe.Pointer(&req)), nil
}

// GetBuffer queries the status and information of a specific buffer.
// It takes a StreamingDevice and the zero-based index of the buffer.
// This function is typically called after buffers have been requested via InitBuffers,
// particularly for memory-mapped buffers to get their offset for mmap().
// It calls the VIDIOC_QUERYBUF ioctl.
//
// Returns a Buffer struct populated with the buffer's details and an error if the ioctl call fails.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-querybuf.html
func GetBuffer(dev StreamingDevice, index uint32) (Buffer, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(dev.BufferType())
	v4l2Buf.memory = C.uint(dev.MemIOType())
	v4l2Buf.index = C.uint(index)

	if err := send(dev.Fd(), C.VIDIOC_QUERYBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		// The original error message "query buffer: type not supported: %w" might be misleading.
		// VIDIOC_QUERYBUF can fail for various reasons, e.g. invalid index or if REQBUFS hasn't been called.
		return Buffer{}, fmt.Errorf("query buffer: VIDIOC_QUERYBUF failed for index %d: %w", index, err)
	}

	return makeBuffer(v4l2Buf), nil
}

// mapMemoryBuffer is an internal helper to memory-map a single buffer from the device.
// It uses sys.Mmap with read/write protection and shared mapping.
func mapMemoryBuffer(fd uintptr, offset int64, length int) ([]byte, error) {
	// Note: The PROT_READ|PROT_WRITE and MAP_SHARED flags are common for V4L2 buffer mapping.
	data, err := sys.Mmap(int(fd), offset, length, sys.PROT_READ|sys.PROT_WRITE, sys.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("map memory buffer: mmap failed: %w", err)
	}
	return data, nil
}

// MapMemoryBuffers queries and memory-maps all buffers for a device that uses IOTypeMMAP.
// It takes a StreamingDevice which provides methods to get buffer count, file descriptor, etc.
// For each buffer index, it calls GetBuffer to retrieve buffer information (like offset and length),
// then uses the internal mapMemoryBuffer function to memory-map the buffer.
// The `TODO` for checking buffer flags from the original code is preserved.
//
// Returns a slice of byte slices (each []byte is a memory-mapped buffer) and an error if any
// part of the process (getting buffer info, mmapping) fails.
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

// unmapMemoryBuffer is an internal helper to unmap a single memory-mapped buffer.
// It takes a byte slice representing the mapped buffer and calls sys.Munmap.
func unmapMemoryBuffer(buf []byte) error {
	if err := sys.Munmap(buf); err != nil {
		return fmt.Errorf("unmap memory buffer: munmap failed: %w", err)
	}
	return nil
}

// UnmapMemoryBuffers unmaps all previously memory-mapped buffers for a device.
// It takes a StreamingDevice, which is expected to provide access to the slice of
// mapped buffers via its Buffers() method.
// It iterates through the buffers obtained from `dev.Buffers()` and calls the internal
// `unmapMemoryBuffer` for each.
//
// Returns an error if the device's buffer list (from `dev.Buffers()`) is nil,
// or if any individual unmap operation fails. If an unmap fails, it returns immediately
// with the error for that specific buffer.
// Note: It's crucial that dev.Buffers() returns the actual slice of mapped byte slices.
func UnmapMemoryBuffers(dev StreamingDevice) error {
	mappedBuffers := dev.Buffers() // Assumes dev.Buffers() returns the [][]byte of mapped buffers.
	if mappedBuffers == nil {
		return fmt.Errorf("unmap buffers: buffers slice is nil, nothing to unmap")
	}
	for i := 0; i < len(mappedBuffers); i++ {
		if mappedBuffers[i] == nil {
			// This case should ideally not happen if buffers were mapped correctly.
			// A log warning here might be useful in practice.
			continue // Skip nil buffers if any
		}
		if err := unmapMemoryBuffer(mappedBuffers[i]); err != nil {
			// Collect first error; consider if all should be attempted and errors aggregated.
			return fmt.Errorf("unmap buffers: error unmapping buffer at index %d: %w", i, err)
		}
	}
	return nil
}

// QueueBuffer queues an empty (or filled) buffer onto the driver's incoming queue.
// fd: file descriptor of the device.
// ioType: type of memory I/O (e.g., MMAP).
// bufType: type of buffer (e.g., VideoCapture).
// index: index of the buffer to queue.
// It returns the queued buffer information and an error if any.
// For multi-planar API, the `Planes` field in `Buffer.Info` would need to be populated by the caller
// before calling this function if the `v4l2_buffer` struct's `m.planes` field is used by the driver for QBUF.
// This implementation assumes single-planar usage or that `m.planes` is handled externally if needed.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html
func QueueBuffer(fd uintptr, ioType IOType, bufType BufType, index uint32) (Buffer, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(bufType)
	v4l2Buf.memory = C.uint(ioType)
	v4l2Buf.index = C.uint(index)
	// For multi-planar API, if C.V4L2_BUF_TYPE_VIDEO_CAPTURE_MPLANE or C.V4L2_BUF_TYPE_VIDEO_OUTPUT_MPLANE is used,
	// the caller might need to prepare v4l2Buf.m.planes and set v4l2Buf.length to the number of planes.
	// This function currently does not explicitly handle setting m.planes.

	if err := send(fd, C.VIDIOC_QBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		return Buffer{}, fmt.Errorf("buffer queue: VIDIOC_QBUF failed: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// DequeueBuffer dequeues a filled (or empty) buffer from the driver's outgoing queue.
// fd: file descriptor of the device.
// ioType: type of memory I/O (e.g., MMAP).
// bufType: type of buffer (e.g., VideoCapture).
// It returns the dequeued buffer information and an error if any.
// Note: If no buffer is available, this call can block if the device was opened in blocking mode,
// or return EAGAIN if opened in non-blocking mode.
// For multi-planar API, the driver will populate `v4l2Buf.m.planes` if applicable.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-qbuf.html
func DequeueBuffer(fd uintptr, ioType IOType, bufType BufType) (Buffer, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(bufType)
	v4l2Buf.memory = C.uint(ioType)
	// For multi-planar API, if C.V4L2_BUF_TYPE_VIDEO_CAPTURE_MPLANE or C.V4L2_BUF_TYPE_VIDEO_OUTPUT_MPLANE is used,
	// the driver will fill in v4l2Buf.m.planes. The caller should ensure enough space is allocated if providing a pointer.
	// However, makeBuffer will handle the C to Go conversion including planes if they are populated by the driver.

	err := send(fd, C.VIDIOC_DQBUF, uintptr(unsafe.Pointer(&v4l2Buf)))
	if err != nil {
		return Buffer{}, fmt.Errorf("buffer dequeue: VIDIOC_DQBUF failed: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}
