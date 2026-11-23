//go:build integration

package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func TestV4L2_CoreUnifiedAndSelection(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	fd, err := v4l2.OpenDevice(devPath, sys.O_RDWR|sys.O_NONBLOCK, 0)
	if err != nil {
		t.Skipf("Cannot open %s: %v", devPath, err)
	}
	defer v4l2.CloseDevice(fd)

	t.Run("FormatDescriptionsByType", func(t *testing.T) {
		descs, err := v4l2.GetAllFormatDescriptionsByType(fd, v4l2.BufTypeVideoCapture)
		if err != nil && len(descs) == 0 {
			t.Skipf("GetAllFormatDescriptionsByType unsupported: %v", err)
		}
		t.Logf("Enumerated %d formats via GetAllFormatDescriptionsByType", len(descs))
	})

	t.Run("TryPixFormat", func(t *testing.T) {
		cur, err := v4l2.GetPixFormatByType(fd, v4l2.BufTypeVideoCapture)
		if err != nil {
			t.Skipf("GetPixFormatByType unsupported: %v", err)
		}
		tried, err := v4l2.TryPixFormat(fd, v4l2.BufTypeVideoCapture, cur)
		if err != nil {
			t.Fatalf("TryPixFormat failed: %v", err)
		}
		if tried.Width == 0 || tried.Height == 0 {
			t.Errorf("TryPixFormat returned zero dimensions: %+v", tried)
		}
		t.Logf("TryPixFormat negotiated %dx%d fourcc=0x%x", tried.Width, tried.Height, tried.PixelFormat)
	})

	t.Run("SelectionAPI", func(t *testing.T) {
		sel, err := v4l2.GetSelection(fd, v4l2.BufTypeVideoCapture, v4l2.SelTargetCropBounds)
		if err != nil {
			t.Skipf("VIDIOC_G_SELECTION unsupported on %s: %v", devPath, err)
		}
		t.Logf("Selection CROP_BOUNDS: %+v", sel.Rect)
	})

	t.Run("CreatePrepareRemoveBuffers", func(t *testing.T) {
		cur, err := v4l2.GetPixFormatByType(fd, v4l2.BufTypeVideoCapture)
		if err != nil {
			t.Skipf("GetPixFormatByType unsupported: %v", err)
		}
		res, err := v4l2.CreateBuffers(fd, 2, v4l2.IOTypeMMAP, v4l2.BufTypeVideoCapture, &cur, nil)
		if err != nil {
			t.Skipf("VIDIOC_CREATE_BUFS unsupported on %s: %v", devPath, err)
		}
		defer v4l2.RequestBuffersForType(fd, v4l2.BufTypeVideoCapture, v4l2.IOTypeMMAP, 0)
		t.Logf("CreateBuffers allocated count=%d index=%d caps=0x%x", res.Count, res.Index, res.Capabilities)

		if _, err := v4l2.PrepareBuffer(fd, v4l2.Buffer{
			Type:   v4l2.BufTypeVideoCapture,
			Memory: v4l2.IOTypeMMAP,
			Index:  res.Index,
		}); err != nil {
			t.Logf("VIDIOC_PREPARE_BUF returned (driver-optional): %v", err)
		}
		if err := v4l2.RemoveBuffers(fd, v4l2.BufTypeVideoCapture, res.Index, res.Count); err != nil {
			t.Logf("VIDIOC_REMOVE_BUFS returned (requires Linux 6.10+): %v", err)
		}
	})
}

func TestV4L2_MultiPlanarHardware(t *testing.T) {
	// Discover any available MPLANE-capable device (e.g., hardware codec/ISP on SoC).
	var mplaneDev string
	for i := 0; i < 32; i++ {
		path := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK, 0)
		if err != nil {
			continue
		}
		cap, err := v4l2.GetCapability(fd)
		_ = v4l2.CloseDevice(fd)
		if err != nil {
			continue
		}
		caps := cap.GetCapabilities()
		if caps&v4l2.CapVideoCaptureMPlane != 0 || caps&v4l2.CapVideoMem2MemMPlane != 0 {
			mplaneDev = path
			break
		}
	}

	if mplaneDev == "" {
		t.Skip("No multi-planar V4L2 device available")
	}

	fd, err := v4l2.OpenDevice(mplaneDev, sys.O_RDWR|sys.O_NONBLOCK, 0)
	if err != nil {
		t.Skipf("Cannot open %s: %v", mplaneDev, err)
	}
	defer v4l2.CloseDevice(fd)

	bufType := v4l2.BufTypeVideoCaptureMPlane
	descs, err := v4l2.GetAllFormatDescriptionsByType(fd, bufType)
	if err != nil && len(descs) == 0 {
		t.Fatalf("GetAllFormatDescriptionsByType(%s) failed: %v", mplaneDev, err)
	}
	t.Logf("%s: enumerated %d MPLANE capture formats", mplaneDev, len(descs))

	cur, err := v4l2.GetPixFormatMPlane(fd, bufType)
	if err != nil {
		t.Fatalf("GetPixFormatMPlane(%s) failed: %v", mplaneDev, err)
	}
	t.Logf("%s: current MPLANE format %dx%d fourcc=0x%08x planes=%d",
		mplaneDev, cur.Width, cur.Height, cur.PixelFormat, cur.NumPlanes)

	tried, err := v4l2.TryPixFormatMPlane(fd, bufType, cur)
	if err != nil {
		t.Fatalf("TryPixFormatMPlane(%s) failed: %v", mplaneDev, err)
	}
	if tried.NumPlanes == 0 {
		t.Fatalf("TryPixFormatMPlane returned 0 planes")
	}

	req, err := v4l2.RequestBuffersForType(fd, bufType, v4l2.IOTypeMMAP, 2)
	if err != nil {
		t.Fatalf("RequestBuffersForType(%s) failed: %v", mplaneDev, err)
	}
	defer v4l2.RequestBuffersForType(fd, bufType, v4l2.IOTypeMMAP, 0)
	if req.Count == 0 {
		t.Fatalf("RequestBuffersForType returned 0 buffers")
	}

	buf, err := v4l2.QueryBuffer(fd, bufType, v4l2.IOTypeMMAP, 0, uint32(tried.NumPlanes))
	if err != nil {
		t.Fatalf("QueryBuffer(%s) failed: %v", mplaneDev, err)
	}
	if len(buf.Planes) != int(tried.NumPlanes) {
		t.Fatalf("QueryBuffer returned %d planes, want %d", len(buf.Planes), tried.NumPlanes)
	}
	t.Logf("%s: queried buffer 0 plane 0 length=%d offset=0x%x",
		mplaneDev, buf.Planes[0].Length, buf.Planes[0].Info.MemOffset)
}

func TestIntegration_DeviceMPlane(t *testing.T) {
	var mplaneDev string
	for i := 0; i < 32; i++ {
		path := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK, 0)
		if err != nil {
			continue
		}
		cap, err := v4l2.GetCapability(fd)
		_ = v4l2.CloseDevice(fd)
		if err != nil {
			continue
		}
		if cap.IsVideoCaptureMultiplanarSupported() || cap.IsM2MMPlaneSupported() {
			mplaneDev = path
			break
		}
	}
	if mplaneDev == "" {
		t.Skip("No multi-planar V4L2 device available")
	}

	dev, err := device.Open(mplaneDev, device.WithBufferType(v4l2.BufTypeVideoCaptureMPlane), device.WithBufferSize(2))
	if err != nil {
		t.Fatalf("device.Open(%s, MPLANE): %v", mplaneDev, err)
	}
	defer dev.Close()

	pixFmt, err := dev.GetPixFormatMPlane()
	if err != nil {
		t.Fatalf("GetPixFormatMPlane(%s): %v", mplaneDev, err)
	}
	if pixFmt.NumPlanes == 0 {
		t.Fatalf("GetPixFormatMPlane(%s): expected NumPlanes > 0", mplaneDev)
	}
	t.Logf("%s: device MPLANE format %dx%d fourcc=0x%08x planes=%d",
		mplaneDev, pixFmt.Width, pixFmt.Height, pixFmt.PixelFormat, pixFmt.NumPlanes)
}
