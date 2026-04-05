// userptr_capture demonstrates frame capture using User Pointer streaming I/O.
//
// Unlike the default MMAP mode where the kernel allocates buffers, USERPTR mode
// uses application-allocated buffers. The consumer API (Start/GetFrames/Stop)
// is identical — the difference is internal.
//
// Usage:
//
//	userptr_capture -d /dev/video0 -n 5
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	devName := "/dev/video0"
	numFrames := 5
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.IntVar(&numFrames, "n", numFrames, "number of frames to capture")
	flag.Parse()

	dev, err := device.Open(
		devName,
		device.WithIOType(v4l2.IOTypeUserPtr),
		device.WithPixFormat(v4l2.PixFormat{
			PixelFormat: v4l2.PixelFmtMJPEG,
			Width:       640,
			Height:      480,
		}),
		device.WithBufferSize(4),
	)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer dev.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		log.Fatalf("failed to start stream: %s", err)
	}

	count := 0
	for frame := range dev.GetFrames() {
		fileName := fmt.Sprintf("frame_%d.jpg", count)
		if err := os.WriteFile(fileName, frame.Data, 0644); err != nil {
			log.Printf("failed to write %s: %s", fileName, err)
			frame.Release()
			continue
		}
		log.Printf("Frame %d: seq=%d, %d bytes, ts=%v -> %s",
			count, frame.Sequence, len(frame.Data), frame.Timestamp, fileName)
		frame.Release()
		count++
		if count >= numFrames {
			break
		}
	}

	if err := dev.Stop(); err != nil {
		log.Fatalf("failed to stop stream: %s", err)
	}
	log.Printf("Captured %d frames", count)
}
