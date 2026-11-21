package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func main() {
	subdevPath := flag.String("subdev", "/dev/v4l-subdev0", "V4L2 sub-device path")
	pad := flag.Uint("pad", 0, "sub-device pad index")
	videoPath := flag.String("video", "", "optional paired /dev/video* capture node to stream frames after inspection")
	frames := flag.Int("n", 5, "number of frames to capture when -video is set")
	flag.Parse()

	fd, err := v4l2.OpenDevice(*subdevPath, sys.O_RDWR|sys.O_CLOEXEC, 0)
	if err != nil {
		log.Fatalf("open %s: %v", *subdevPath, err)
	}
	defer v4l2.CloseDevice(fd)

	padIdx := uint32(*pad)
	fmt.Printf("Inspecting V4L2 sub-device %s (pad %d)\n", *subdevPath, padIdx)

	activeFmt, err := v4l2.GetSubdevFormat(fd, padIdx, v4l2.SubdevFormatActive)
	if err != nil {
		log.Fatalf("get active format: %v", err)
	}
	fmt.Printf("Active format: %dx%d code=%s (0x%04x) field=%d colorspace=%d\n",
		activeFmt.Format.Width,
		activeFmt.Format.Height,
		mbusName(activeFmt.Format.Code),
		activeFmt.Format.Code,
		activeFmt.Format.Field,
		activeFmt.Format.Colorspace,
	)

	// Validate TRY format negotiation without mutating active hardware state.
	tryFmt := activeFmt
	tryFmt.Which = v4l2.SubdevFormatTry
	if negotiated, err := v4l2.SetSubdevFormat(fd, tryFmt); err == nil {
		fmt.Printf("TRY format negotiated: %dx%d code=%s\n",
			negotiated.Format.Width, negotiated.Format.Height, mbusName(negotiated.Format.Code))
	}

	codes, err := v4l2.EnumSubdevMBusCodes(fd, padIdx, v4l2.SubdevFormatActive)
	if err == nil {
		fmt.Printf("Supported media bus codes (%d):\n", len(codes))
		for _, c := range codes {
			fmt.Printf("  - %s (0x%04x)\n", mbusName(c.Code), c.Code)
			if sizes, err := v4l2.EnumSubdevFrameSizes(fd, padIdx, c.Code, v4l2.SubdevFormatActive); err == nil {
				for _, sz := range sizes {
					fmt.Printf("      size: %dx%d .. %dx%d\n",
						sz.MinWidth, sz.MinHeight, sz.MaxWidth, sz.MaxHeight)
				}
			}
		}
	}

	for _, target := range []struct {
		name   string
		target v4l2.SelectionTarget
	}{
		{"NATIVE_SIZE", v4l2.SelTargetNativeSize},
		{"CROP_BOUNDS", v4l2.SelTargetCropBounds},
		{"CROP", v4l2.SelTargetCrop},
	} {
		if sel, err := v4l2.GetSubdevSelection(fd, padIdx, v4l2.SubdevFormatActive, target.target); err == nil {
			fmt.Printf("Selection %s: (%d,%d) %dx%d\n",
				target.name, sel.Rect.Left, sel.Rect.Top, sel.Rect.Width, sel.Rect.Height)
		}
	}

	if *videoPath == "" {
		return
	}

	dev, err := device.Open(*videoPath, device.WithBufferSize(4))
	if err != nil {
		log.Fatalf("open video %s: %v", *videoPath, err)
	}
	defer dev.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		log.Fatalf("start %s: %v", *videoPath, err)
	}
	defer dev.Stop()

	got := 0
	for frame, err := range dev.Frames(ctx) {
		if err != nil {
			log.Fatalf("capture: %v", err)
		}
		got++
		fmt.Printf("Captured frame #%d from %s: seq=%d bytes=%d\n", got, *videoPath, frame.Sequence, len(frame.Data))
		if got >= *frames {
			break
		}
	}
}

func mbusName(code v4l2.MBusFormatCode) string {
	if name, ok := v4l2.MBusFormatNames[code]; ok {
		return name
	}
	return fmt.Sprintf("0x%04x", code)
}
