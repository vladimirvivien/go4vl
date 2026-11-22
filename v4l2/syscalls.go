package v4l2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

func cString[T ~int8 | ~uint8](b []T) string {
	if len(b) == 0 {
		return ""
	}
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&b[0])), len(b))
	if i := bytes.IndexByte(raw, 0); i >= 0 {
		return string(raw[:i])
	}
	return string(raw)
}

// OpenDevice offers a simpler file-open operation than the Go API's os.OpenFile  (the Go API's
// operation causes some drivers to return busy). It also applies file validation prior to opening the device.
func OpenDevice(path string, flags int, mode uint32) (uintptr, error) {
	fstat, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("open device: %w", err)
	}

	if fstat.Mode()&fs.ModeCharDevice == 0 {
		return 0, fmt.Errorf("device open: %s: not character device", path)
	}

	return openDev(path, flags, mode)
}

// openDev offers a simpler file open operation than the Go API OpenFile.
// See https://cs.opensource.google/go/go/+/refs/tags/go1.19.1:src/os/file_unix.go;l=205
func openDev(path string, flags int, mode uint32) (uintptr, error) {
	var fd int
	var err error
	for {
		fd, err = sys.Openat(sys.AT_FDCWD, path, flags, mode)
		if err == nil {
			break
		}

		if errors.Is(err, ErrorInterrupted) {
			continue //retry
		}

		return 0, &os.PathError{Op: "open", Path: path, Err: err}
	}
	return uintptr(fd), nil
}

// CloseDevice closes the device.
func CloseDevice(fd uintptr) error {
	return closeDev(fd)
}

func closeDev(fd uintptr) error {
	return sys.Close(int(fd))
}

// ReadDevice reads frame data from a V4L2 device using the read() syscall.
// The device must support V4L2_CAP_READWRITE.
func ReadDevice(fd uintptr, buf []byte) (int, error) {
	n, err := sys.Read(int(fd), buf)
	if err != nil {
		return n, fmt.Errorf("read device: %w", err)
	}
	return n, nil
}

// WriteDevice writes frame data to a V4L2 output device using the write() syscall.
// The device must support V4L2_CAP_READWRITE and V4L2_CAP_VIDEO_OUTPUT.
func WriteDevice(fd uintptr, buf []byte) (int, error) {
	n, err := sys.Write(int(fd), buf)
	if err != nil {
		return n, fmt.Errorf("write device: %w", err)
	}
	return n, nil
}

// ioctl is a wrapper for Syscall(SYS_IOCTL)
func ioctl(fd, req, arg uintptr) (err sys.Errno) {
	for {
		_, _, errno := sys.Syscall(sys.SYS_IOCTL, fd, req, arg)
		switch errno {
		case 0:
			return 0
		case sys.EINTR:
			continue // retry
		default:
			return errno
		}
	}
}

// send sends a request to the kernel (via ioctl syscall)
func send(fd, req, arg uintptr) error {
	errno := ioctl(fd, req, arg)
	if errno == 0 {
		return nil
	}
	parsedErr := parseErrorType(errno)
	switch {
	case errors.Is(parsedErr, ErrorUnsupported),
		errors.Is(parsedErr, ErrorSystem),
		errors.Is(parsedErr, ErrorBadArgument):
		return parsedErr
	case errors.Is(parsedErr, ErrorTimeout),
		errors.Is(parsedErr, ErrorTemporary):
		return errno
	default:
		return errno
	}
}

// WaitForRead returns a channel that is notified when a device is ready to be read.
// The goroutine exits when ctx is cancelled.
func WaitForRead(ctx context.Context, dev Device) <-chan struct{} {
	sigChan := make(chan struct{})

	go func(fd uintptr) {
		defer close(sigChan)
		fds := []sys.PollFd{{
			Fd:     int32(fd),
			Events: sys.POLLIN,
		}}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			fds[0].Revents = 0
			n, err := sys.Poll(fds, 100)
			if err != nil {
				if errors.Is(err, sys.EINTR) {
					continue
				}
				return
			}
			if n == 0 {
				continue
			}

			select {
			case sigChan <- struct{}{}:
			case <-ctx.Done():
				return
			}
		}
	}(dev.Fd())

	return sigChan
}
