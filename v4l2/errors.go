package v4l2

import (
	"errors"
	sys "syscall"
)

// Error variables represent common V4L2 operation failures.
// These errors are returned by V4L2 operations to indicate specific failure conditions.
// Use errors.Is() to check for specific error types in your error handling logic.
var (
	// ErrorSystem indicates a system-level error such as:
	// - EBADF: Invalid file descriptor
	// - ENOMEM: Insufficient memory
	// - ENODEV: Device not found or removed
	// - EIO: I/O error during device communication
	// - ENXIO: Device not configured or unavailable
	// - EFAULT: Bad address in user space
	// These errors are typically unrecoverable and indicate fundamental problems.
	ErrorSystem = errors.New("system error")

	// ErrorBadArgument indicates invalid parameters passed to a V4L2 operation.
	// This corresponds to EINVAL and means the arguments don't meet the requirements
	// of the specific ioctl or function call.
	ErrorBadArgument = errors.New("bad argument error")

	// ErrorTemporary indicates a temporary condition that might resolve if retried.
	// Operations returning this error can potentially succeed on retry after a delay.
	ErrorTemporary = errors.New("temporary error")

	// ErrorTimeout indicates an operation timed out waiting for a condition.
	// This commonly occurs when waiting for frames during capture with a timeout set.
	ErrorTimeout = errors.New("timeout error")

	// ErrorUnsupported indicates the requested ioctl or operation is not supported.
	// This corresponds to ENOTTY and means the device doesn't implement the requested functionality.
	ErrorUnsupported = errors.New("unsupported error")

	// ErrorUnsupportedFeature indicates a specific feature or capability is not supported by the device.
	// This is a higher-level error used when a device lacks required capabilities for an operation.
	ErrorUnsupportedFeature = errors.New("feature unsupported error")

	// ErrorInterrupted indicates the operation was interrupted by a signal.
	// This corresponds to EINTR and the operation can typically be retried.
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
