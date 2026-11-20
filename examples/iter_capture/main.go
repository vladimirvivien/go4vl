package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/vladimirvivien/go4vl/device"
)

func main() {
	devName := flag.String("d", "/dev/video0", "V4L2 video capture device path")
	count := flag.Int("n", 30, "number of frames to capture")
	bufCount := flag.Uint("b", 4, "number of MMAP buffers")
	flag.Parse()

	dev, err := device.Open(*devName, device.WithBufferSize(uint32(*bufCount)))
	if err != nil {
		log.Fatalf("open %s: %v", *devName, err)
	}
	defer dev.Close()

	pixFmt, err := dev.GetPixFormat()
	if err != nil {
		log.Fatalf("get format: %v", err)
	}
	fmt.Printf("Capturing %d frames from %s (%dx%d, size=%d bytes)\n",
		*count, *devName, pixFmt.Width, pixFmt.Height, pixFmt.SizeImage)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		log.Fatalf("start stream: %v", err)
	}
	defer dev.Stop()

	start := time.Now()
	captured := 0
	var totalBytes uint64

	for frame, err := range dev.Frames(ctx) {
		if err != nil {
			log.Fatalf("capture frame: %v", err)
		}
		captured++
		totalBytes += uint64(len(frame.Data))
		if captured == 1 || captured == *count {
			fmt.Printf("  frame #%d: seq=%d bytes=%d index=%d\n",
				captured, frame.Sequence, len(frame.Data), frame.Index)
		}
		if captured >= *count {
			break
		}
	}

	elapsed := time.Since(start)
	fps := float64(captured) / elapsed.Seconds()
	mbps := float64(totalBytes) / (1024 * 1024) / elapsed.Seconds()
	fmt.Printf("Captured %d frames in %v (%.2f fps, %.2f MiB/s)\n", captured, elapsed.Round(time.Millisecond), fps, mbps)
}
