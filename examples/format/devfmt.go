package main

import (
	"flag"
	"log"
	"strings"

	dev "github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	devName := "/dev/video0"
	width := 640
	height := 480
	format := "yuyv"

	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.IntVar(&width, "w", width, "capture width")
	flag.IntVar(&height, "h", height, "capture height")
	flag.StringVar(&format, "f", format, "pixel format")
	flag.Parse()

	fmtEnc := v4l2.PixelFmtYUYV
	switch strings.ToLower(format) {
	case "mjpeg":
		fmtEnc = v4l2.PixelFmtMJPEG
	case "h264", "h.264":
		fmtEnc = v4l2.PixelFmtH264
	case "yuyv":
		fmtEnc = v4l2.PixelFmtYUYV
	}

	device, err := dev.Open(
		devName,
		dev.WithPixFormat(v4l2.PixFormat{Width: uint32(width), Height: uint32(height), PixelFormat: fmtEnc, Field: v4l2.FieldNone}),
		dev.WithFPS(15),
	)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()

	currFmt, err := device.GetPixFormat()
	if err != nil {
		log.Fatalf("unable to get format: %s", err)
	}
	log.Printf("Current format: %s", currFmt)

	// FPS
	fps, err := device.GetFrameRate()
	if err != nil {
		log.Fatalf("failed to get fps: %s", err)
	}
	log.Printf("current frame rate: %d fps", fps)
	// update fps
	if fps < 30 {
		if err := device.SetFrameRate(30); err != nil {
			log.Fatalf("failed to set frame rate: %s", err)
		}
	}
	fps, err = device.GetFrameRate()
	if err != nil {
		log.Fatalf("failed to get fps: %s", err)
	}
	log.Printf("updated frame rate: %d fps", fps)
}
