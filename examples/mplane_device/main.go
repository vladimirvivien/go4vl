package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	var (
		devPath  string
		frames   int
		buffers  uint
		width    uint
		height   uint
		timeout  time.Duration
	)
	flag.StringVar(&devPath, "d", "/dev/video14", "multi-planar V4L2 capture device path")
	flag.IntVar(&frames, "frames", 5, "number of frames to capture (0 to inspect format only)")
	flag.UintVar(&buffers, "buffers", 4, "number of MMAP buffers")
	flag.UintVar(&width, "w", 0, "optional frame width")
	flag.UintVar(&height, "h", 0, "optional frame height")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "capture timeout")
	flag.Parse()

	opts := []device.Option{
		device.WithBufferType(v4l2.BufTypeVideoCaptureMPlane),
		device.WithBufferSize(uint32(buffers)),
	}
	if width > 0 && height > 0 {
		opts = append(opts, device.WithPixFormatMPlane(v4l2.PixFormatMPlane{
			Width:  uint32(width),
			Height: uint32(height),
		}))
	}

	dev, err := device.Open(devPath, opts...)
	if err != nil {
		log.Fatalf("open %s: %v", devPath, err)
	}
	defer dev.Close()

	cap := dev.Capability()
	fmt.Printf("Device: %s (%s, driver=%s)\n", devPath, cap.Card, cap.Driver)

	pixFmt, err := dev.GetPixFormatMPlane()
	if err != nil {
		log.Fatalf("GetPixFormatMPlane: %v", err)
	}
	fmt.Printf("Format: %dx%d pixelformat=0x%08x planes=%d\n",
		pixFmt.Width, pixFmt.Height, pixFmt.PixelFormat, pixFmt.NumPlanes)
	for i := uint8(0); i < pixFmt.NumPlanes; i++ {
		pf := pixFmt.PlaneFmt[i]
		fmt.Printf("  plane[%d]: bytesperline=%d sizeimage=%d\n", i, pf.BytesPerLine, pf.SizeImage)
	}

	if frames <= 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		log.Fatalf("Start: %v", err)
	}
	defer dev.Stop()

	captured := 0
	for frame, err := range dev.Frames(ctx) {
		if err != nil {
			log.Fatalf("Frames: %v", err)
		}
		captured++
		fmt.Printf("frame %d (seq=%d): planes=%d", captured, frame.Sequence, len(frame.Planes))
		for i, p := range frame.Planes {
			fmt.Printf(" plane[%d]=%dB", i, len(p))
		}
		fmt.Println()
		if captured >= frames {
			break
		}
	}
}
