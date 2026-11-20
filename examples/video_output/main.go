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
	devName := flag.String("d", "/dev/video10", "V4L2 video output device path (e.g., v4l2loopback)")
	width := flag.Uint("w", 640, "frame width")
	height := flag.Uint("h", 480, "frame height")
	count := flag.Int("n", 30, "number of frames to write")
	fps := flag.Int("fps", 30, "target output frame rate")
	flag.Parse()

	dev, err := device.Open(
		*devName,
		device.WithVideoOutputEnabled(),
		device.WithBufferSize(4),
		device.WithPixFormat(v4l2.PixFormat{
			Width:       uint32(*width),
			Height:      uint32(*height),
			PixelFormat: v4l2.PixelFmtYUYV,
			Field:       v4l2.FieldNone,
		}),
	)
	if err != nil {
		log.Fatalf("open output %s: %v", *devName, err)
	}
	defer dev.Close()

	pixFmt, err := dev.GetPixFormat()
	if err != nil {
		log.Fatalf("get output format: %v", err)
	}
	fmt.Printf("Streaming %d YUYV frames (%dx%d, %d bytes/frame) to %s\n",
		*count, pixFmt.Width, pixFmt.Height, pixFmt.SizeImage, *devName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		log.Fatalf("start output stream: %v", err)
	}
	defer dev.Stop()

	frameBuf := make([]byte, pixFmt.SizeImage)
	interval := time.Second / time.Duration(*fps)

	for i := 0; i < *count; i++ {
		// Fill synthetic YUYV moving gradient pattern.
		yVal := byte((i * 8) % 256)
		for p := 0; p+1 < len(frameBuf); p += 2 {
			frameBuf[p] = yVal
			frameBuf[p+1] = 128
		}
		if err := dev.WriteFrame(ctx, frameBuf); err != nil {
			log.Fatalf("write frame %d: %v", i, err)
		}
		time.Sleep(interval)
	}

	fmt.Printf("Wrote %d frames to %s\n", *count, *devName)
}
