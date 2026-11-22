package v4l2

import (
	"errors"
	"testing"

	sys "golang.org/x/sys/unix"
)

func TestParseErrorTypePreservesErrno(t *testing.T) {
	tests := []struct {
		name     string
		errno    sys.Errno
		sentinel error
	}{
		{"ENODEV", sys.ENODEV, ErrorSystem},
		{"ENOMEM", sys.ENOMEM, ErrorSystem},
		{"EIO", sys.EIO, ErrorSystem},
		{"EINVAL", sys.EINVAL, ErrorBadArgument},
		{"ENOTTY", sys.ENOTTY, ErrorUnsupported},
		{"EINTR", sys.EINTR, ErrorInterrupted},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := parseErrorType(tc.errno)
			if !errors.Is(err, tc.sentinel) {
				t.Errorf("errors.Is(%v, %v) = false, want true", err, tc.sentinel)
			}
			if !errors.Is(err, tc.errno) {
				t.Errorf("errors.Is(%v, %v) = false, want true", err, tc.errno)
			}
		})
	}
}
