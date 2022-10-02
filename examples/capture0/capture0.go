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
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.Parse()

	// open device
	device, err := device.Open(
		devName,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 640, Height: 480}),
	)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()

	// start stream
	ctx, stop := context.WithCancel(context.TODO())
	if err := device.Start(ctx); err != nil {
		log.Fatalf("failed to start stream: %s", err)
	}

	// process frames from capture channel
	totalFrames := 10
	count := 0
	log.Printf("Capturing %d frames...", totalFrames)

	for frame := range device.GetOutput() {
		fileName := fmt.Sprintf("capture_%d.jpg", count)
		file, err := os.Create(fileName)
		if err != nil {
			log.Printf("failed to create file %s: %s", fileName, err)
			continue
		}
		if _, err := file.Write(frame); err != nil {
			log.Printf("failed to write file %s: %s", fileName, err)
			continue
		}
		log.Printf("Saved file: %s", fileName)
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s: %s", fileName, err)
		}
		count++
		if count >= totalFrames {
			break
		}
		time.Sleep(1 * time.Second)
	}

	stop() // stop capture
	fmt.Println("Done.")

}
