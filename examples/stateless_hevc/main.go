package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func main() {
	videoPath := flag.String("video", "/dev/video19", "stateless HEVC decoder video device (e.g., rpi-hevc-dec)")
	mediaPath := flag.String("media", "/dev/media0", "associated media controller device")
	flag.Parse()

	vfd, err := v4l2.OpenDevice(*videoPath, sys.O_RDWR|sys.O_NONBLOCK|sys.O_CLOEXEC, 0)
	if err != nil {
		log.Fatalf("open %s: %v", *videoPath, err)
	}
	defer v4l2.CloseDevice(vfd)

	cap, err := v4l2.GetCapability(vfd)
	if err != nil {
		log.Fatalf("GetCapability(%s): %v", *videoPath, err)
	}
	fmt.Printf("Video device: %s (%s, driver=%s)\n", *videoPath, cap.Card, cap.Driver)

	outFormats, err := v4l2.GetAllFormatDescriptionsByType(vfd, v4l2.BufTypeVideoOutputMPlane)
	if err != nil {
		log.Fatalf("GetAllFormatDescriptionsByType(OUTPUT_MPLANE): %v", err)
	}
	hasHEVCSlice := false
	for _, f := range outFormats {
		fmt.Printf("  OUTPUT format: %s (0x%08x)\n", f.Description, f.PixelFormat)
		if f.PixelFormat == v4l2.PixelFmtHEVCSlice {
			hasHEVCSlice = true
		}
	}
	if !hasHEVCSlice {
		log.Fatalf("%s does not advertise V4L2_PIX_FMT_HEVC_SLICE", *videoPath)
	}

	mfd, err := v4l2.OpenDevice(*mediaPath, sys.O_RDWR|sys.O_CLOEXEC, 0)
	if err != nil {
		log.Fatalf("open %s: %v", *mediaPath, err)
	}
	defer v4l2.CloseDevice(mfd)

	req, err := v4l2.AllocMediaRequest(mfd)
	if err != nil {
		log.Fatalf("AllocMediaRequest(%s): %v", *mediaPath, err)
	}
	defer req.Close()

	// Build a minimal 1920x1080 stateless HEVC SPS/PPS/DecodeParams control batch
	// bound to the allocated MediaRequest and validate it via VIDIOC_TRY_EXT_CTRLS.
	ctrls := v4l2.NewExtControls()
	ctrls.SetRequestFD(req.Fd())

	sps := &v4l2.ControlHEVCSPS{
		Pic_width_in_luma_samples:               1920,
		Pic_height_in_luma_samples:              1080,
		Chroma_format_idc:                       1, // 4:2:0
		Log2_min_luma_coding_block_size_minus3:  0,
		Log2_diff_max_min_luma_coding_block_size: 3,
		Log2_min_luma_transform_block_size_minus2: 0,
		Log2_diff_max_min_luma_transform_block_size: 3,
		Sps_max_dec_pic_buffering_minus1:        1,
	}
	pps := &v4l2.ControlHEVCPPS{
		Init_qp_minus26: 0,
	}
	decodeParams := &v4l2.ControlHEVCDecodeParams{
		Pic_order_cnt_val: 0,
	}

	if err := ctrls.AddHEVCSPS(sps); err != nil {
		log.Fatalf("AddHEVCSPS: %v", err)
	}
	if err := ctrls.AddHEVCPPS(pps); err != nil {
		log.Fatalf("AddHEVCPPS: %v", err)
	}
	if err := ctrls.AddHEVCDecodeParams(decodeParams); err != nil {
		log.Fatalf("AddHEVCDecodeParams: %v", err)
	}

	if err := v4l2.TryExtControls(vfd, ctrls); err != nil {
		log.Fatalf("TryExtControls(HEVC): %v", err)
	}
	fmt.Printf("Validated stateless HEVC SPS/PPS/DecodeParams with request FD %d\n", req.Fd())

	if err := req.Reinit(); err != nil {
		log.Fatalf("req.Reinit: %v", err)
	}
	fmt.Println("Reinitialized MediaRequest successfully")
}
