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
	dev, err := device.Open(
		devName,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 640, Height: 480}),
	)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer dev.Close()

	// start stream
	ctx, stop := context.WithCancel(context.TODO())
	if err := dev.Start(ctx); err != nil {
		log.Fatalf("failed to start stream: %s", err)
	}

	// process frames from capture channel with metadata
	totalFrames := 10
	count := 0
	log.Printf("Capturing %d frames with metadata...", totalFrames)

	startTime := time.Now()
	var lastSeq uint32

	for frame := range dev.GetFrames() {
		// Detect dropped frames by checking sequence gaps
		if count > 0 && frame.Sequence != lastSeq+1 {
			dropped := frame.Sequence - lastSeq - 1
			log.Printf("WARNING: %d frame(s) dropped (seq jump: %d -> %d)",
				dropped, lastSeq, frame.Sequence)
		}
		lastSeq = frame.Sequence

		// Calculate frame age (latency from capture to processing)
		latency := time.Since(frame.Timestamp)

		// Check frame type for compressed video
		frameType := "Data"
		if frame.IsKeyFrame() {
			frameType = "Keyframe"
		} else if frame.IsPFrame() {
			frameType = "P-frame"
		} else if frame.IsBFrame() {
			frameType = "B-frame"
		}

		fileName := fmt.Sprintf("frame_%05d.jpg", frame.Sequence)
		log.Printf("Frame %d: %s | Size: %d bytes | Type: %s | Latency: %v",
			frame.Sequence, fileName, len(frame.Data), frameType, latency)

		// Save frame to file
		file, err := os.Create(fileName)
		if err != nil {
			log.Printf("failed to create file %s: %s", fileName, err)
			frame.Release() // Always release even on error
			continue
		}

		if _, err := file.Write(frame.Data); err != nil {
			log.Printf("failed to write file %s: %s", fileName, err)
			file.Close()
			frame.Release() // Always release even on error
			continue
		}

		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s: %s", fileName, err)
		}

		// IMPORTANT: Release frame buffer back to pool
		frame.Release()

		count++
		if count >= totalFrames {
			break
		}
	}

	elapsed := time.Since(startTime)
	fps := float64(count) / elapsed.Seconds()

	stop() // stop capture
	fmt.Printf("\nDone.\n")
	fmt.Printf("Captured %d frames in %v (%.2f FPS)\n", count, elapsed, fps)

	// Display pool statistics
	stats := device.DefaultFramePool().Stats()
	fmt.Printf("\nFrame Pool Statistics:\n")
	fmt.Printf("  Total Gets:       %d\n", stats.Gets)
	fmt.Printf("  Total Puts:       %d\n", stats.Puts)
	fmt.Printf("  Allocations:      %d\n", stats.Allocs)
	fmt.Printf("  Resizes:          %d\n", stats.Resizes)
	fmt.Printf("  Outstanding:      %d\n", stats.Outstanding)
	fmt.Printf("  Hit Rate:         %.2f%%\n", stats.HitRate*100)
}
