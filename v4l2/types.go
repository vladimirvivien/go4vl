package v4l2

import (
	"context"
)

// Device defines the basic interface for a V4L2 (Video4Linux2) device.
// It abstracts common operations applicable to various types of video devices.
type Device interface {
	// Name returns the name or path of the device.
	Name() string
	// Fd returns the file descriptor of the opened device.
	Fd() uintptr
	// Capability returns the capabilities of the device (e.g., video capture, streaming).
	// The Capability struct is defined in capability.go.
	Capability() Capability
	// MemIOType returns the current memory I/O method used by the device (e.g., MMAP, USERPTR).
	// IOType is defined in streaming.go.
	MemIOType() IOType
	// GetOutput returns a read-only channel that delivers output data from the device,
	// typically video frames for capture devices.
	GetOutput() <-chan []byte
	// SetInput provides a channel for sending input data to the device,
	// typically for output devices. The implementation details depend on the specific device.
	SetInput(<-chan []byte)
	// Close releases the resources associated with the device, including closing its file descriptor.
	Close() error
}

// StreamingDevice defines an interface for V4L2 devices that support streaming I/O,
// particularly using memory-mapped buffers. It embeds the base Device interface.
type StreamingDevice interface {
	Device

	// Buffers returns a slice of byte slices, where each inner slice represents a memory-mapped buffer
	// used for streaming. This should be called after buffers have been successfully mapped.
	Buffers() [][]byte
	// BufferType returns the type of buffer used by the device for streaming (e.g., video capture).
	// BufType is defined in streaming.go.
	BufferType() BufType
	// BufferCount returns the number of buffers currently configured for streaming.
	BufferCount() uint32
	// Start initiates the video streaming process.
	// It takes a context for managing cancellation or timeouts of the streaming loop.
	// Implementations should handle buffer allocation, queuing, and starting the capture/output stream.
	Start(ctx context.Context) error
	// Stop terminates the video streaming process.
	// Implementations should handle stopping the stream, unmapping buffers, and other cleanup.
	Stop() error
}
