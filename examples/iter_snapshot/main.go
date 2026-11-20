package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vladimirvivien/go4vl/device"
)

func main() {
	devName := flag.String("d", "/dev/video0", "V4L2 video capture device path")
	outFile := flag.String("o", "snapshot.raw", "output file path for the captured frame")
	flag.Parse()

	dev, err := device.Open(*devName, device.WithBufferSize(2))
	if err != nil {
		log.Fatalf("open %s: %v", *devName, err)
	}
	defer dev.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		log.Fatalf("start stream: %v", err)
	}
	defer dev.Stop()

	for frame, err := range dev.Frames(ctx) {
		if err != nil {
			log.Fatalf("capture frame: %v", err)
		}
		if err := os.WriteFile(*outFile, frame.Data, 0o644); err != nil {
			log.Fatalf("write %s: %v", *outFile, err)
		}
		fmt.Printf("Saved frame #%d (%d bytes) to %s\n", frame.Sequence, len(frame.Data), *outFile)
		break // Buffer is automatically re-queued on break.
	}
}
