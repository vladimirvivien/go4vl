package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func main() {
	devPath := flag.String("d", "/dev/video0", "V4L2 device path")
	setCrop := flag.Bool("set", false, "Apply a crop selection rectangle via VIDIOC_S_SELECTION")
	left := flag.Int("left", 0, "Crop left offset")
	top := flag.Int("top", 0, "Crop top offset")
	width := flag.Uint("width", 640, "Crop width")
	height := flag.Uint("height", 480, "Crop height")
	flag.Parse()

	fd, err := v4l2.OpenDevice(*devPath, sys.O_RDWR|sys.O_NONBLOCK, 0)
	if err != nil {
		log.Fatalf("open %s: %v", *devPath, err)
	}
	defer v4l2.CloseDevice(fd)

	cap, err := v4l2.GetCapability(fd)
	if err != nil {
		log.Fatalf("querycap: %v", err)
	}
	fmt.Printf("Device: %s (%s) driver=%s\n", *devPath, cap.Card, cap.Driver)

	caps := cap.GetCapabilities()
	bufType := v4l2.BufTypeVideoCapture
	if caps&v4l2.CapVideoCapture == 0 && (caps&v4l2.CapVideoCaptureMPlane != 0 || caps&v4l2.CapVideoMem2MemMPlane != 0) {
		bufType = v4l2.BufTypeVideoCaptureMPlane
	}

	targets := []struct {
		name   string
		target v4l2.SelectionTarget
	}{
		{"CROP", v4l2.SelTargetCrop},
		{"CROP_DEFAULT", v4l2.SelTargetCropDefault},
		{"CROP_BOUNDS", v4l2.SelTargetCropBounds},
		{"NATIVE_SIZE", v4l2.SelTargetNativeSize},
		{"COMPOSE", v4l2.SelTargetCompose},
		{"COMPOSE_DEFAULT", v4l2.SelTargetComposeDefault},
		{"COMPOSE_BOUNDS", v4l2.SelTargetComposeBounds},
	}

	for _, t := range targets {
		sel, err := v4l2.GetSelection(fd, bufType, t.target)
		if err != nil {
			fmt.Printf("  %-16s: unsupported (%v)\n", t.name, err)
			continue
		}
		fmt.Printf("  %-16s: left=%d top=%d width=%d height=%d\n",
			t.name, sel.Rect.Left, sel.Rect.Top, sel.Rect.Width, sel.Rect.Height)
	}

	if *setCrop {
		req := v4l2.Selection{
			Type:   bufType,
			Target: v4l2.SelTargetCrop,
			Rect: v4l2.Rect{
				Left:   int32(*left),
				Top:    int32(*top),
				Width:  uint32(*width),
				Height: uint32(*height),
			},
		}
		applied, err := v4l2.SetSelection(fd, req)
		if err != nil {
			log.Fatalf("set crop selection: %v", err)
		}
		fmt.Printf("Applied CROP selection: left=%d top=%d width=%d height=%d\n",
			applied.Rect.Left, applied.Rect.Top, applied.Rect.Width, applied.Rect.Height)
	}
}
