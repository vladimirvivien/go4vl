//go:build integration

package test

import (
	"path/filepath"
	"testing"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func TestIntegration_MediaRequestLifecycle(t *testing.T) {
	matches, _ := filepath.Glob("/dev/media*")
	if len(matches) == 0 {
		t.Skip("No /dev/media* devices available")
	}

	tested := 0
	for _, path := range matches {
		fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_CLOEXEC, 0)
		if err != nil {
			continue
		}

		req, err := v4l2.AllocMediaRequest(fd)
		if err != nil {
			_ = v4l2.CloseDevice(fd)
			continue
		}
		tested++

		if req.Fd() <= 0 {
			_ = req.Close()
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("%s: expected positive request FD, got %d", path, req.Fd())
		}
		if err := req.Reinit(); err != nil {
			_ = req.Close()
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("%s: Reinit() failed: %v", path, err)
		}
		if err := req.Close(); err != nil {
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("%s: Close() failed: %v", path, err)
		}

		info, _ := v4l2.GetMediaDeviceInfo(fd)
		t.Logf("%s (%s): allocated, reinitialized, and closed Media Request successfully", path, info.Driver)
		_ = v4l2.CloseDevice(fd)
	}

	if tested == 0 {
		t.Skip("No /dev/media* device on this host supports MEDIA_IOC_REQUEST_ALLOC")
	}
}
