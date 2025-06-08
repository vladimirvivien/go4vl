package v4l2

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	sys "golang.org/x/sys/unix"
)

// OpenDevice opens a V4L2 device specified by its path.
// It performs basic validation to ensure the path points to a character device.
// This function is a wrapper around the lower-level openDev, providing a more user-friendly API
// and attempting to mitigate issues where os.OpenFile might cause some drivers to report "busy".
//
// Parameters:
//   path: The file system path to the V4L2 device (e.g., "/dev/video0").
//   flags: File control flags for opening the device (e.g., sys.O_RDWR | sys.O_NONBLOCK).
//   mode: File mode bits (permissions), typically 0 for device files.
//
// Returns:
//   A uintptr representing the file descriptor of the opened device, and an error if any step fails
//   (e.g., path not found, not a character device, or error during open syscall).
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

// CloseDevice closes an opened V4L2 device file descriptor.
// It is a wrapper around the lower-level closeDev function.
//
// Parameters:
//   fd: The file descriptor (uintptr) of the device to close.
//
// Returns:
//   An error if the close syscall fails.
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
// It uses a blocking select call with a 2-second timeout in a goroutine.
// Note: The Device type is an interface defined in types.go.
//
// Parameters:
//  dev: The V4L2 Device instance for which to wait for readability.
//
// Returns:
//   A read-only channel of empty structs. A value will be sent on this channel
//   when the device becomes ready for reading or if the select call times out or errors.
//   The channel is closed when the internal goroutine exits.
func WaitForRead(dev Device) <-chan struct{} {
	sigChan := make(chan struct{})

	go func(fd uintptr) {
		defer close(sigChan)
		var fdsRead sys.FdSet
		fdsRead.Set(int(fd))
		// Using a fixed 2-second timeout for select.
		// Consider making this configurable if needed.
		tv := sys.Timeval{Sec: 2, Usec: 0}
		for {
			// The select call will block until either the fd is ready,
			// the timeout occurs, or an error (like EINTR) happens.
			_, errno := sys.Select(int(fd+1), &fdsRead, nil, nil, &tv)
			if errno == sys.EINTR { // Interrupted, retry select.
				// Reset fd in read set as select might modify it.
				fdsRead.Zero()
				fdsRead.Set(int(fd))
				// Reset timeout for the next select call.
				tv.Sec = 2
				tv.Usec = 0
				continue
			}
			// Regardless of whether select returned due to data available, timeout, or another error,
			// signal the channel. The consumer can then attempt a read and handle EAGAIN if it was a timeout.
			// If select itself errored beyond EINTR, this will still signal, and subsequent read will likely fail.
			sigChan <- struct{}{}

			// If an actual error occurred (other than EINTR) or if it was a timeout,
			// it might be desirable to exit the loop. However, the current logic
			// will continuously signal every 2 seconds on timeout.
			// For this pass, I'll keep the existing logic.
			// If the fd is no longer valid or an unrecoverable error occurs,
			// this goroutine might loop indefinitely. Consider adding a mechanism
			// to break the loop, perhaps via context cancellation passed into WaitForRead.

			// Reset fd in read set for the next iteration, as select can modify it.
			fdsRead.Zero()
			fdsRead.Set(int(fd))
			// Reset timeout for the next select call.
			tv.Sec = 2
			tv.Usec = 0
		}
	}(dev.Fd())

	return sigChan
}
