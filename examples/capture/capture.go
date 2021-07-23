package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	devName := "/dev/video0"
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.Parse()

	// open device
	device, err := v4l2.Open(devName)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()

	// configuration
	if err := device.SetPixFormat(v4l2.PixFormat{
		Width:       640,
		Height:      480,
		PixelFormat: v4l2.PixelFmtMJPEG,
		Field:       v4l2.FieldNone,
	}); err != nil {
		log.Fatalf("failed to set format: %s", err)
	}

	// start stream
	if err := device.StartStream(10); err != nil {
		log.Fatalf("failed to start stream: %s", err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	frameChan, err := device.Capture(ctx, 15)
	if err != nil {
		log.Fatal(err)
	}

	// process frames from capture channel
	totalFrames := 10
	count := 0
	for frame := range frameChan {
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
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s: %s", fileName, err)
		}
		count++
		if count >= totalFrames {
			break
		}
	}

	cancel() // stop capture
	if err := device.StopStream(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done.")

}
