package main

import (
	"flag"
	"log"
	"strings"

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

	device, err := v4l2.Open(devName)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()

	currFmt, err := device.GetPixFormat()
	if err != nil {
		log.Fatalf("unable to get format: %s", err)
	}
	log.Printf("Current format: %s", currFmt)

	fmtEnc := v4l2.PixelFmtYUYV
	switch strings.ToLower(format) {
	case "mjpeg":
		fmtEnc = v4l2.PixelFmtMJPEG
	case "h264", "h.264":
		fmtEnc = v4l2.PixelFmtH264
	case "yuyv":
		fmtEnc = v4l2.PixelFmtYUYV
	}

	if err := device.SetPixFormat(v4l2.PixFormat{
		Width: uint32(width),
		Height: uint32(height),
		PixelFormat: fmtEnc,
		Field: v4l2.FieldNone,
	}); err != nil {
		log.Fatalf("failed to set format: %s", err)
	}

	currFmt, err = device.GetPixFormat()
	if err != nil {
		log.Fatalf("unable to get format: %s", err)
	}
	log.Printf("Updated format: %s", currFmt)
}