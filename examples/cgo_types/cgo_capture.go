package main

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../../include/
#include <linux/videodev2.h>
*/
import "C"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

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
	PixelFmtMJPEG uint32 = C.V4L2_PIX_FMT_MJPEG
)

// Pix format field types
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L89
const (
	FieldAny  uint32 = C.V4L2_FIELD_ANY
	FieldNone uint32 = C.V4L2_FIELD_NONE
)

// buff stream types
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L142
const (
	BufTypeVideoCapture uint32 = C.V4L2_BUF_TYPE_VIDEO_CAPTURE
	BufTypeVideoOutput  uint32 = C.V4L2_BUF_TYPE_VIDEO_OUTPUT
	BufTypeOverlay      uint32 = C.V4L2_BUF_TYPE_VIDEO_OVERLAY
)

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
	var v4l2Fmt C.struct_v4l2_format
	v4l2Fmt._type = C.uint(BufTypeVideoCapture)
	*(*C.struct_v4l2_pix_format)(unsafe.Pointer(&v4l2Fmt.fmt[0])) = *(*C.struct_v4l2_pix_format)(unsafe.Pointer(&pixFmt))

	// send command
	if err := ioctl(fd, C.VIDIOC_S_FMT, uintptr(unsafe.Pointer(&v4l2Fmt))); err != nil {
		return err
	}
	log.Printf("setting format to: %dx%d\n", pixFmt.Width, pixFmt.Height)
	return nil
}

func getFormat(fd uintptr) (PixFormat, error) {
	var v4l2Fmt C.struct_v4l2_format
	v4l2Fmt._type = C.uint(BufTypeVideoCapture)

	// send command
	if err := ioctl(fd, C.VIDIOC_G_FMT, uintptr(unsafe.Pointer(&v4l2Fmt))); err != nil {
		return PixFormat{}, err
	}

	var pixFmt PixFormat
	*(*C.struct_v4l2_pix_format)(unsafe.Pointer(&pixFmt)) = *(*C.struct_v4l2_pix_format)(unsafe.Pointer(&v4l2Fmt.fmt[0]))

	return pixFmt, nil

}

// =========================== Buffers and Streaming ========================== //

// Memory buffer types
// https://elixir.bootlin.com/linux/v5.13-rc6/source/include/uapi/linux/videodev2.h#L188
const (
	StreamMemoryTypeMMAP uint32 = C.V4L2_MEMORY_MMAP
)

// reqBuffers requests that the device allocates a `count`
// number of internal buffers before they can be mapped into
// the application's address space. The driver will return
// the actual number of buffers allocated in the RequestBuffers
// struct.
func reqBuffers(fd uintptr, count uint32) error {
	var reqbuf C.struct_v4l2_requestbuffers
	reqbuf.count = C.uint(count)
	reqbuf._type = C.uint(BufTypeVideoCapture)
	reqbuf.memory = C.uint(StreamMemoryTypeMMAP)

	if err := ioctl(fd, C.VIDIOC_REQBUFS, uintptr(unsafe.Pointer(&reqbuf))); err != nil {
		return err
	}
	log.Printf("Request %d buffers OK\n", count)
	return nil
}

// ================================ Map device Memory ===============================

// BufferService is embedded union m
// in v4l2_buffer C type.
type BufferService struct {
	Offset  uint32
	UserPtr uintptr
	Planes  uintptr
	FD      int32
}

// mmapBuffer first queries the status of the device buffer at idx
// by retrieving BufferInfo which returns the length of the buffer and
// the current offset of the allocated buffers.  That information is
// used to map the device's buffer unto the application's address space.
func mmapBuffer(fd uintptr, idx uint32) ([]byte, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(BufTypeVideoCapture)
	v4l2Buf.memory = C.uint(StreamMemoryTypeMMAP)
	v4l2Buf.index = C.uint(idx)

	// send ioctl command
	if err := ioctl(fd, C.VIDIOC_QUERYBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		return nil, err
	}

	// grab m union and place it in type BufferService
	bufSvc := *(*BufferService)(unsafe.Pointer(&v4l2Buf.m[0]))

	// map the memory and get []byte to access it
	mbuf, err := sys.Mmap(int(fd), int64(bufSvc.Offset), int(v4l2Buf.length), sys.PROT_READ|sys.PROT_WRITE, sys.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	return mbuf, nil
}

// =========================== Start device streaming =========================

// startStreaming requests the device to start the capture process and start
// filling device buffers.
func startStreaming(fd uintptr) error {
	bufType := C.uint(BufTypeVideoCapture)
	if err := ioctl(fd, C.VIDIOC_STREAMON, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return err
	}
	return nil
}

// ======================== Queue/Dequeue device buffer =======================

// queueBuffer requests that an empty buffer is enqueued into the device's
// incoming queue at the specified index (so that it can be filled later).
func queueBuffer(fd uintptr, idx uint32) error {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(BufTypeVideoCapture)
	v4l2Buf.memory = C.uint(StreamMemoryTypeMMAP)
	v4l2Buf.index = C.uint(idx)

	if err := ioctl(fd, C.VIDIOC_QBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		return err
	}
	return nil
}

// dequeueBuffer is called to dequeue a filled buffer from the devices buffer queue.
// Once a device buffer is dequeued, it is mapped and is ready to be read by the application.
func dequeueBuffer(fd uintptr) (uint32, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(BufTypeVideoCapture)
	v4l2Buf.memory = C.uint(StreamMemoryTypeMMAP)

	if err := ioctl(fd, C.VIDIOC_DQBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		return 0, err
	}
	return uint32(v4l2Buf.bytesused), nil
}

// =========================== Start device streaming =========================

// stopStreaming requests the device to stop the streaming process and release
// buffer resources.
func stopStreaming(fd uintptr) error {
	bufType := C.uint(BufTypeVideoCapture)

	if err := ioctl(fd, C.VIDIOC_STREAMOFF, uintptr(unsafe.Pointer(&bufType))); err != nil {
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
	if err := setFormat(fd, PixFormat{
		Width:       640,
		Height:      480,
		PixelFormat: PixelFmtMJPEG,
		Field:       FieldNone,
	}); err != nil {
		log.Fatal(err)
	}

	pixFmt, err := getFormat(fd)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Format set: %dx%d pixel format %d [FmtMJPG]", pixFmt.Width, pixFmt.Height, pixFmt.PixelFormat)

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
