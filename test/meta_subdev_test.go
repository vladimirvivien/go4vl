//go:build integration

package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func TestIntegration_MetadataFormat(t *testing.T) {
	tested := 0
	for i := 0; i < 32; i++ {
		path := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK|sys.O_CLOEXEC, 0)
		if err != nil {
			continue
		}
		cap, err := v4l2.GetCapability(fd)
		if err != nil {
			_ = v4l2.CloseDevice(fd)
			continue
		}
		caps := cap.GetCapabilities()
		var bufType v4l2.BufType
		switch {
		case caps&v4l2.CapMetadataCapture != 0:
			bufType = v4l2.BufTypeMetaCapture
		case caps&v4l2.CapMetadataOutput != 0:
			bufType = v4l2.BufTypeMetaOutput
		default:
			_ = v4l2.CloseDevice(fd)
			continue
		}

		cur, err := v4l2.GetMetaFormat(fd, bufType)
		if err != nil {
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("GetMetaFormat(%s, %d): %v", path, bufType, err)
		}
		tried, err := v4l2.TryMetaFormat(fd, bufType, cur)
		if err != nil {
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("TryMetaFormat(%s, %d): %v", path, bufType, err)
		}
		if cur.BufferSize == 0 {
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("%s: GetMetaFormat returned BufferSize=0 (%+v)", path, cur)
		}
		t.Logf("%s (%s): meta=%s try=%s", path, cap.Driver, cur, tried)
		_ = v4l2.CloseDevice(fd)
		tested++
	}

	if tested == 0 {
		t.Skip("No V4L2 metadata capture/output device available")
	}
}

func TestIntegration_SubdevCapability(t *testing.T) {
	const subdevPath = "/dev/v4l-subdev0"
	if _, err := os.Stat(subdevPath); err != nil {
		t.Skipf("Sub-device %s not available: %v", subdevPath, err)
	}

	fd, err := v4l2.OpenDevice(subdevPath, sys.O_RDWR|sys.O_CLOEXEC, 0)
	if err != nil {
		t.Fatalf("OpenDevice(%s): %v", subdevPath, err)
	}
	defer v4l2.CloseDevice(fd)

	cap, err := v4l2.GetSubdevCapability(fd)
	if err != nil {
		t.Fatalf("GetSubdevCapability(%s): %v", subdevPath, err)
	}
	t.Logf("%s: version=0x%x caps=0x%x ro=%v streams=%v",
		subdevPath, cap.Version, cap.Capabilities, cap.IsReadOnly(), cap.SupportsStreams())

	if clientCaps, err := v4l2.GetSubdevClientCapability(fd); err == nil {
		t.Logf("%s: clientCaps=0x%x", subdevPath, clientCaps)
	}
	if cap.SupportsStreams() {
		if routes, err := v4l2.GetSubdevRouting(fd, v4l2.SubdevFormatActive); err == nil {
			t.Logf("%s: %d active routes", subdevPath, len(routes))
		}
	}
}
