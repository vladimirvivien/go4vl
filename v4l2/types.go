package v4l2

type Device interface {
	Name() string
	FileDescriptor() uintptr
	Capability() Capability
	Buffers() [][]byte
	BufferType() BufType
	BufferCount() uint32
	MemIOType() IOType
}
