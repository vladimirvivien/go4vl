package v4l2

import (
	"errors"
	sys "syscall"
)

var (
	ErrorSystem             = errors.New("system error")
	ErrorBadArgument        = errors.New("bad argument error")
	ErrorTemporary          = errors.New("temporary error")
	ErrorTimeout            = errors.New("timeout error")
	ErrorUnsupported        = errors.New("unsupported error")
	ErrorUnsupportedFeature = errors.New("feature unsupported error")
	ErrorInterrupted        = errors.New("interrupted")
)

func parseErrorType(errno sys.Errno) error {
	switch errno {
	case sys.EBADF, sys.ENOMEM, sys.ENODEV, sys.EIO, sys.ENXIO, sys.EFAULT: // structural, terminal
		return ErrorSystem
	case sys.EINTR:
		return ErrorInterrupted
	case sys.EINVAL: // bad argument
		return ErrorBadArgument
	case sys.ENOTTY: // unsupported
		return ErrorUnsupported
	default:
		if errno.Timeout() {
			return ErrorTimeout
		}
		if errno.Temporary() {
			return ErrorTemporary
		}
		return errno
	}
}
