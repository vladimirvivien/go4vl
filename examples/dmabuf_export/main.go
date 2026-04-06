// dmabuf_export demonstrates exporting V4L2 MMAP buffers as DMA-BUF file descriptors.
//
// This allows other subsystems (GPU, DRM, other V4L2 devices) to access the
// captured video frames via zero-copy buffer sharing.
//
// Usage:
//
//	dmabuf_export -d /dev/video0
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
	devName := "/dev/video0"
	bufCount := 4
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.IntVar(&bufCount, "b", bufCount, "number of buffers")
	flag.Parse()

	dev, err := device.Open(
		devName,
		device.WithPixFormat(v4l2.PixFormat{
			PixelFormat: v4l2.PixelFmtMJPEG,
			Width:       640,
			Height:      480,
		}),
		device.WithBufferSize(uint32(bufCount)),
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
	defer dev.Stop()

	// Export each MMAP buffer as a DMA-BUF fd
	for i := 0; i < bufCount; i++ {
		fd, err := dev.ExportBuffer(uint32(i), 0)
		if err != nil {
			log.Fatalf("failed to export buffer %d: %s", i, err)
		}
		fmt.Printf("Buffer %d exported as DMA-BUF fd=%d\n", i, fd)
		// In a real application, pass fd to GPU/DRM/other device:
		//   gpu.ImportBuffer(fd)
	}

	// Capture a few frames to show streaming still works
	count := 0
	for frame := range dev.GetFrames() {
		fmt.Printf("Frame %d: seq=%d, %d bytes\n", count, frame.Sequence, len(frame.Data))
		frame.Release()
		count++
		if count >= 5 {
			break
		}
	}
}
