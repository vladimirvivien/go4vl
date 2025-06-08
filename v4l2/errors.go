package v4l2

import (
	"errors"
	sys "syscall"
)

// Predefined error variables for common V4L2 issues.
var (
	// ErrorSystem indicates a general system error occurred, often related to device access or memory.
	// Corresponds to syscall errors like EBADF, ENOMEM, ENODEV, EIO, ENXIO, EFAULT.
	ErrorSystem = errors.New("system error")
	// ErrorBadArgument indicates that an invalid argument was passed to a V4L2 ioctl.
	// Corresponds to syscall.EINVAL.
	ErrorBadArgument = errors.New("bad argument error")
	// ErrorTemporary indicates that a temporary condition prevented the operation from completing.
	// Retrying the operation may succeed.
	ErrorTemporary = errors.New("temporary error")
	// ErrorTimeout indicates that an operation timed out.
	ErrorTimeout = errors.New("timeout error")
	// ErrorUnsupported indicates that a requested operation or feature is not supported by the driver or device.
	// Corresponds to syscall.ENOTTY when an ioctl is not supported.
	ErrorUnsupported = errors.New("unsupported error")
	// ErrorUnsupportedFeature indicates that a specific feature within a supported operation is not available.
	ErrorUnsupportedFeature = errors.New("feature unsupported error")
	// ErrorInterrupted indicates that a blocking operation was interrupted by a signal.
	// Corresponds to syscall.EINTR.
	ErrorInterrupted = errors.New("interrupted")
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
