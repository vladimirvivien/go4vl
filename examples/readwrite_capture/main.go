// readwrite_capture demonstrates frame capture using the read/write I/O method.
//
// Unlike streaming mode (Start/GetFrames/Stop), read/write mode uses direct
// read() syscalls for simple, synchronous frame capture.
//
// Usage:
//
//	readwrite_capture -d /dev/video0 -n 5
//	readwrite_capture -d /dev/video0 -n 5 -raw
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	devName := "/dev/video0"
	numFrames := 5
	useRaw := false
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.IntVar(&numFrames, "n", numFrames, "number of frames to capture")
	flag.BoolVar(&useRaw, "raw", useRaw, "use low-level Read() instead of ReadFrame()")
	flag.Parse()

	dev, err := device.Open(
		devName,
		device.WithIOMethod(device.IOMethodReadWrite),
		device.WithPixFormat(v4l2.PixFormat{
			PixelFormat: v4l2.PixelFmtMJPEG,
			Width:       640,
			Height:      480,
		}),
	)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer dev.Close()

	pixFmt, err := dev.GetPixFormat()
	if err != nil {
		log.Fatalf("failed to get pixel format: %s", err)
	}
	log.Printf("Format: %dx%d, SizeImage: %d", pixFmt.Width, pixFmt.Height, pixFmt.SizeImage)

	if useRaw {
		captureRaw(dev, pixFmt, numFrames)
	} else {
		captureFrames(dev, numFrames)
	}
}

func captureRaw(dev *device.Device, pixFmt v4l2.PixFormat, n int) {
	buf := make([]byte, pixFmt.SizeImage)
	for i := 0; i < n; i++ {
		bytesRead, err := dev.Read(buf)
		if err != nil {
			log.Fatalf("Read() failed: %s", err)
		}

		fileName := fmt.Sprintf("frame_%d.jpg", i)
		if err := os.WriteFile(fileName, buf[:bytesRead], 0644); err != nil {
			log.Printf("failed to write %s: %s", fileName, err)
			continue
		}
		log.Printf("Frame %d: %d bytes -> %s", i, bytesRead, fileName)
	}
}

func captureFrames(dev *device.Device, n int) {
	for i := 0; i < n; i++ {
		frame, err := dev.ReadFrame()
		if err != nil {
			log.Fatalf("ReadFrame() failed: %s", err)
		}

		fileName := fmt.Sprintf("frame_%d.jpg", i)
		if err := os.WriteFile(fileName, frame.Data, 0644); err != nil {
			log.Printf("failed to write %s: %s", fileName, err)
			continue
		}
		log.Printf("Frame %d: seq=%d, %d bytes, ts=%v -> %s",
			i, frame.Sequence, len(frame.Data), frame.Timestamp, fileName)
	}
}
