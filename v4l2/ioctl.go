package v4l2

import (
	sys "golang.org/x/sys/unix"
)

// ioctl is a wrapper for Syscall(SYS_IOCTL)
func ioctl(fd, req, arg uintptr) (err sys.Errno) {
	if _, _, errno := sys.Syscall(sys.SYS_IOCTL, fd, req, arg); errno != 0 {
		if errno != 0 {
			err = errno
			return
		}
	}
	return 0
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
