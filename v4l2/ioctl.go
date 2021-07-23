package v4l2

import (
	sys "syscall"
)

// Send sends raw command to driver (via ioctl syscall)
func Send(fd, req, arg uintptr) error {
	return ioctl(fd, req, arg)
}

func ioctl(fd, req, arg uintptr) (err error) {
	if _, _, errno := sys.Syscall(sys.SYS_IOCTL, fd, req, arg); errno != 0 {
		switch errno {
		case sys.EINVAL:
			err = ErrorUnsupported
		default:
			err = errno
		}
		return
	}
	return nil
}
