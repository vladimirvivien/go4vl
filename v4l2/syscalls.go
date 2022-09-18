package v4l2

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	sys "golang.org/x/sys/unix"
)

// OpenDevice offers a simpler file-open operation than the Go API's os.OpenFile  (the Go API's
// operation causes some drivers to return busy). It also applies file validation prior to opening the device.
func OpenDevice(path string, flags int, mode uint32) (uintptr, error) {
	fstat, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("open device: %w", err)
	}

	if (fstat.Mode() | fs.ModeCharDevice) == 0 {
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
	switch parsedErr {
	case ErrorUnsupported, ErrorSystem, ErrorBadArgument:
		return parsedErr
	case ErrorTimeout, ErrorTemporary:
		// TODO add code for automatic retry/recovery
		return errno
	default:
		return errno
	}
}

// WaitForRead returns a channel that can be used to be notified when
// a device's is ready to be read.
func WaitForRead(dev Device) <-chan struct{} {
	sigChan := make(chan struct{})

	go func(fd uintptr) {
		defer close(sigChan)
		var fdsRead sys.FdSet
		fdsRead.Set(int(fd))
		tv := sys.Timeval{Sec: 2, Usec: 0}
		for {
			_, errno := sys.Select(int(fd+1), &fdsRead, nil, nil, &tv)
			if errno == sys.EINTR {
				continue
			}

			sigChan <- struct{}{}
		}
	}(dev.Fd())

	return sigChan
}
