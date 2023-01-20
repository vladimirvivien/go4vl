package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// ========================= V4L2 command encoding =====================
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/asm-generic/ioctl.h

const (
	// ioctl command layout
	iocNone  = 0 // no op
	iocWrite = 1 // userland app is writing, kernel reading
	iocRead  = 2 // userland app is reading, kernel writing

	iocTypeBits   = 8
	iocNumberBits = 8
	iocSizeBits   = 14
	iocOpBits     = 2

	numberPos = 0
	typePos   = numberPos + iocNumberBits
	sizePos   = typePos + iocTypeBits
	opPos     = sizePos + iocSizeBits
)

// ioctl command encoding funcs
func ioEnc(iocMode, iocType, number, size uintptr) uintptr {
	return (iocMode << opPos) |
		(iocType << typePos) |
		(number << numberPos) |
		(size << sizePos)
}

func ioEncR(iocType, number, size uintptr) uintptr {
	return ioEnc(iocRead, iocType, number, size)
}

func ioEncW(iocType, number, size uintptr) uintptr {
	return ioEnc(iocWrite, iocType, number, size)
}

func ioEncRW(iocType, number, size uintptr) uintptr {
	return ioEnc(iocRead|iocWrite, iocType, number, size)
}

// four character pixel format encoding
func fourcc(a, b, c, d uint32) uint32 {
	return (a | b<<8) | c<<16 | d<<24
}

// wrapper for ioctl system call
func ioctl(fd, req, arg uintptr) (err error) {
	if _, _, errno := sys.Syscall(sys.SYS_IOCTL, fd, req, arg); errno != 0 {
		err = errno
		return
	}
	return nil
}

// ========================= Pixel Format =========================
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L682

var (
	PixelFmtMJPEG = fourcc('M', 'J', 'P', 'G') // V4L2_PIX_FMT_MJPEG
)

// Pix format field types
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L89
const (
	FieldAny  uint32 = iota // V4L2_FIELD_ANY
	FieldNone               // V4L2_FIELD_NONE
)

// buff stream types
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L142
const (
	BufTypeVideoCapture uint32 = iota + 1 // V4L2_BUF_TYPE_VIDEO_CAPTURE = 1
	BufTypeVideoOutput                    // V4L2_BUF_TYPE_VIDEO_OUTPUT  = 2
	BufTypeOverlay                        // V4L2_BUF_TYPE_VIDEO_OVERLAY = 3
)

// Format represents C type v4l2_format
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L2324
type Format struct {
	StreamType uint32
	fmt        [200]byte // max union size
}

// PixFormat represents v4l2_pix_format
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L496
type PixFormat struct {
	Width        uint32
	Height       uint32
	PixelFormat  uint32
	Field        uint32
	BytesPerLine uint32
	SizeImage    uint32
	Colorspace   uint32
	Priv         uint32
	Flags        uint32
	YcbcrEnc     uint32
	Quantization uint32
	XferFunc     uint32
}

// setsFormat sets pixel format of device
func setFormat(fd uintptr, pixFmt PixFormat) error {
	format := Format{StreamType: BufTypeVideoCapture}

	// a bit of C union type magic with unsafe.Pointer
	*(*PixFormat)(unsafe.Pointer(&format.fmt[0])) = pixFmt

	// encode command to send
	vidiocSetFormat := ioEncRW('V', 5, uintptr(unsafe.Sizeof(Format{})))

	// send command
	if err := ioctl(fd, vidiocSetFormat, uintptr(unsafe.Pointer(&format))); err != nil {
		return err
	}
	return nil
}

// =========================== Buffers and Streaming ========================== //

// Memory buffer types
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L188
const (
	StreamMemoryTypeMMAP    uint32 = iota + 1 // V4L2_MEMORY_MMAP             = 1,
	StreamMemoryTypeUserPtr                   // V4L2_MEMORY_USERPTR          = 2,
	StreamMemoryTypeOverlay                   // V4L2_MEMORY_OVERLAY          = 3,
	StreamMemoryTypeDMABuf                    // V4L2_MEMORY_DMABUF           = 4,
)

// RequestBuffers represents C type
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L951
type RequestBuffers struct {
	Count        uint32
	StreamType   uint32
	Memory       uint32
	Capabilities uint32
	Reserved     [1]uint32
}

// reqBuffers requests that the device allocates a `count`
// number of internal buffers before they can be mapped into
// the application's address space. The driver will return
// the actual number of buffers allocated in the RequestBuffers
// struct.
func reqBuffers(fd uintptr, count uint32) error {
	reqbuf := RequestBuffers{
		StreamType: BufTypeVideoCapture,
		Count:      count,
		Memory:     StreamMemoryTypeMMAP,
	}
	vidiocReqBufs := ioEncRW('V', 8, uintptr(unsafe.Sizeof(RequestBuffers{})))
	if err := ioctl(fd, vidiocReqBufs, uintptr(unsafe.Pointer(&reqbuf))); err != nil {
		return err
	}
	return nil
}

// ================================ Map device Memory ===============================

// BufferInfo represents C type v4l2_buffer
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L1037
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

// BufferService is embedded union m
// in v4l2_buffer C type.
type BufferService struct {
	Offset  uint32
	UserPtr uintptr
	Planes  uintptr
	FD      int32
}

// Timecode represents C type v4l2_timecode
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L875
type Timecode struct {
	Type     uint32
	Flags    uint32
	Frames   uint8
	Seconds  uint8
	Minutes  uint8
	Hours    uint8
	Userbits [4]uint8
}

// mmapBuffer first queries the status of the device buffer at idx
// by retrieving BufferInfo which returns the length of the buffer and
// the current offset of the allocated buffers.  That information is
// used to map the device's buffer unto the application's address space.
func mmapBuffer(fd uintptr, idx uint32) ([]byte, error) {
	buf := BufferInfo{
		StreamType: BufTypeVideoCapture,
		Memory:     StreamMemoryTypeMMAP,
		Index:      idx,
	}

	// send ioctl command
	vidiocQueryBuf := ioEncRW('V', 9, uintptr(unsafe.Sizeof(BufferInfo{}))) // VIDIOC_QUERYBUF
	if err := ioctl(fd, vidiocQueryBuf, uintptr(unsafe.Pointer(&buf))); err != nil {
		return nil, err
	}

	// grab m union and place it in type BufferService
	bufSvc := *(*BufferService)(unsafe.Pointer(&buf.m[0]))

	// map the memory and get []byte to access it
	mbuf, err := sys.Mmap(int(fd), int64(bufSvc.Offset), int(buf.Length), sys.PROT_READ|sys.PROT_WRITE, sys.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	return mbuf, nil
}

// =========================== Start device streaming =========================

// startStreaming requests the device to start the capture process and start
// filling device buffers.
func startStreaming(fd uintptr) error {
	bufType := BufTypeVideoCapture
	vidiocStreamOn := ioEncW('V', 18, uintptr(unsafe.Sizeof(int32(0)))) // VIDIOC_STREAMON
	if err := ioctl(fd, vidiocStreamOn, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return err
	}
	return nil
}

// ======================== Queue/Dequeue device buffer =======================

// queueBuffer requests that an empty buffer is enqueued into the device's
// incoming queue at the specified index (so that it can be filled later).
func queueBuffer(fd uintptr, idx uint32) error {
	buf := BufferInfo{
		StreamType: BufTypeVideoCapture,
		Memory:     StreamMemoryTypeMMAP,
		Index:      idx,
	}
	vidiocQueueBuf := ioEncRW('V', 15, uintptr(unsafe.Sizeof(BufferInfo{}))) // VIDIOC_QBUF
	if err := ioctl(fd, vidiocQueueBuf, uintptr(unsafe.Pointer(&buf))); err != nil {
		return err
	}
	return nil
}

// dequeueBuffer is called to dequeue a filled buffer from the devices buffer queue.
// Once a device buffer is dequeued, it is mapped and is ready to be read by the application.
func dequeueBuffer(fd uintptr) (uint32, error) {
	buf := BufferInfo{
		StreamType: BufTypeVideoCapture,
		Memory:     StreamMemoryTypeMMAP,
	}
	vidiocDequeueBuf := ioEncRW('V', 17, uintptr(unsafe.Sizeof(BufferInfo{}))) // VIDIOC_DQBUF
	if err := ioctl(fd, vidiocDequeueBuf, uintptr(unsafe.Pointer(&buf))); err != nil {
		return 0, err
	}
	return buf.BytesUsed, nil
}

// =========================== Start device streaming =========================

// stopStreaming requests the device to stop the streaming process and release
// buffer resources.
func stopStreaming(fd uintptr) error {
	bufType := BufTypeVideoCapture
	vidiocStreamOff := ioEncW('V', 19, uintptr(unsafe.Sizeof(int32(0)))) // VIDIOC_STREAMOFF
	if err := ioctl(fd, vidiocStreamOff, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return err
	}
	return nil
}

// use sys.Select to wait for the device to become read-ready.
func waitForDeviceReady(fd uintptr) error {
	timeval := sys.NsecToTimeval((2 * time.Second).Nanoseconds())
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
			return fmt.Errorf("wait for device ready: timeout")
		default:
			return nil
		}
	}
}

func main() {
	var devName string
	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.Parse()

	// open device
	devFile, err := os.OpenFile(devName, sys.O_RDWR|sys.O_NONBLOCK, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer devFile.Close()
	fd := devFile.Fd()

	// Set the format
	log.Println("setting format to 640x480 MJPEG")
	if err := setFormat(fd, PixFormat{
		Width:       640,
		Height:      480,
		PixelFormat: PixelFmtMJPEG,
		Field:       FieldNone,
	}); err != nil {
		log.Fatal(err)
	}

	// request device to setup 3 buffers
	if err := reqBuffers(fd, 3); err != nil {
		log.Fatal(err)
	}

	// map a device buffer to a local byte slice
	// here we use the latest buffer
	data, err := mmapBuffer(fd, 2)
	if err != nil {
		log.Fatalf("unable to map device buffer: %s", err)
	}

	// now, queue an initial device buffer at the selected index
	//  to be filled with data prior to starting the device stream
	if err := queueBuffer(fd, 2); err != nil {
		log.Fatalf("failed to queue initial buffer: %s", err)
	}

	// now, ask the device to start the stream
	if err := startStreaming(fd); err != nil {
		log.Fatalf("failed to start streaming: %s", err)
	}

	// now wait for the device to be ready for read operation,
	// this means the mapped buffer is ready to be consumed
	if err := waitForDeviceReady(fd); err != nil {
		log.Fatalf("failed during device read-wait: %s", err)
	}

	// dequeue the device buffer so that the local mapped byte slice
	// is filled.
	bufSize, err := dequeueBuffer(fd)
	if err != nil {
		log.Fatalf("failed during device read-wait: %s", err)
	}

	// save mapped buffer bytes to file
	jpgFile, err := os.Create("capture.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer jpgFile.Close()
	if _, err := jpgFile.Write(data[:bufSize]); err != nil {
		log.Fatalf("failed to save file: %s", err)
	}

	// release streaming resources
	if err := stopStreaming(fd); err != nil {
		log.Fatalf("failed to stop stream: %s", err)
	}
}
