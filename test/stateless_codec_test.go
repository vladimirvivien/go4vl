//go:build integration

package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func TestIntegration_StatelessHEVCControls(t *testing.T) {
	var hevcDev string
	for i := 0; i < 32; i++ {
		path := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK|sys.O_CLOEXEC, 0)
		if err != nil {
			continue
		}
		descs, err := v4l2.GetAllFormatDescriptionsByType(fd, v4l2.BufTypeVideoOutputMPlane)
		_ = v4l2.CloseDevice(fd)
		if err != nil {
			continue
		}
		for _, d := range descs {
			if d.PixelFormat == v4l2.PixelFmtHEVCSlice {
				hevcDev = path
				break
			}
		}
		if hevcDev != "" {
			break
		}
	}

	if hevcDev == "" {
		t.Skip("No stateless HEVC decoder device (V4L2_PIX_FMT_HEVC_SLICE) available")
	}

	fd, err := v4l2.OpenDevice(hevcDev, sys.O_RDWR|sys.O_NONBLOCK|sys.O_CLOEXEC, 0)
	if err != nil {
		t.Fatalf("OpenDevice(%s): %v", hevcDev, err)
	}
	defer v4l2.CloseDevice(fd)

	ctrls := v4l2.NewExtControls()
	if err := ctrls.AddHEVCSPS(&v4l2.ControlHEVCSPS{
		Pic_width_in_luma_samples:                   1920,
		Pic_height_in_luma_samples:                  1080,
		Chroma_format_idc:                           1,
		Log2_diff_max_min_luma_coding_block_size:    3,
		Log2_diff_max_min_luma_transform_block_size: 3,
		Sps_max_dec_pic_buffering_minus1:            1,
	}); err != nil {
		t.Fatalf("AddHEVCSPS: %v", err)
	}
	if err := ctrls.AddHEVCPPS(&v4l2.ControlHEVCPPS{}); err != nil {
		t.Fatalf("AddHEVCPPS: %v", err)
	}
	if err := ctrls.AddHEVCDecodeParams(&v4l2.ControlHEVCDecodeParams{}); err != nil {
		t.Fatalf("AddHEVCDecodeParams: %v", err)
	}

	if err := v4l2.TryExtControls(fd, ctrls); err != nil {
		t.Fatalf("TryExtControls(%s): %v", hevcDev, err)
	}
	t.Logf("%s: validated stateless HEVC SPS/PPS/DecodeParams via VIDIOC_TRY_EXT_CTRLS", hevcDev)
}
