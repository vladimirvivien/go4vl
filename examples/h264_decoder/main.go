// h264_decoder demonstrates hardware H.264 decoding using V4L2 stateful codec API.
//
// This example:
// 1. Reads H.264 data from a file
// 2. Decodes using hardware acceleration
// 3. Saves decoded frames as raw NV12 files
//
// Requirements:
// - A V4L2 hardware decoder device (e.g., /dev/video10 on Raspberry Pi)
// - An H.264 input file
//
// Usage:
//
//	./h264_decoder -d /dev/video10 -i input.h264 -o frames/
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	decoderPath = flag.String("d", "/dev/video10", "decoder device path")
	inputFile   = flag.String("i", "input.h264", "input H.264 file")
	outputDir   = flag.String("o", "frames", "output directory for decoded frames")
	maxFrames   = flag.Int("n", 0, "max frames to decode (0 = all)")
	chunkSize   = flag.Int("s", 65536, "read chunk size in bytes")
)

func main() {
	flag.Parse()

	fmt.Println("H.264 Hardware Decoder Example")
	fmt.Println("==============================")
	fmt.Println()

	// Check if decoder device exists
	if _, err := os.Stat(*decoderPath); os.IsNotExist(err) {
		fmt.Printf("Decoder device %s not found.\n", *decoderPath)
		fmt.Println("\nNote: Hardware decoders are typically available on:")
		fmt.Println("  - Raspberry Pi (VideoCore): /dev/video10")
		fmt.Println("  - Intel Quick Sync: Check with v4l2-ctl --list-devices")
		fmt.Println("  - Rockchip RK3399: /dev/video0 (rkvdec)")
		os.Exit(1)
	}

	// Check if input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		fmt.Printf("Input file %s not found.\n", *inputFile)
		os.Exit(1)
	}

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// Open input file
	inFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Failed to open input file: %v\n", err)
		os.Exit(1)
	}
	defer inFile.Close()

	// Get input file size
	fileInfo, _ := inFile.Stat()
	fmt.Printf("Input: %s (%.2f MB)\n", *inputFile, float64(fileInfo.Size())/(1024*1024))

	// Open decoder
	decoder, err := openDecoder()
	if err != nil {
		fmt.Printf("Failed to open decoder: %v\n", err)
		os.Exit(1)
	}
	defer decoder.Close()

	fmt.Printf("Decoder: %s\n", decoder.Capability().Card)
	fmt.Printf("Output: %s/\n", *outputDir)
	fmt.Println()

	// Set up context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nInterrupted, stopping...")
		cancel()
	}()

	// Start decoding
	if err := runDecoder(ctx, decoder, inFile); err != nil {
		fmt.Printf("Decoding failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nDecoding complete!")
}

func openDecoder() (*device.Decoder, error) {
	return device.OpenDecoder(*decoderPath, device.DecoderConfig{
		InputFormat: v4l2.PixFormat{
			PixelFormat: v4l2.PixelFmtH264,
		},
		OutputFormat: v4l2.PixFormat{
			PixelFormat: v4l2.PixelFmtNV12,
		},
		InputBufferCount:  8,
		OutputBufferCount: 4,
	})
}

func runDecoder(ctx context.Context, decoder *device.Decoder, inFile *os.File) error {
	// Start decoder
	if err := decoder.Start(ctx); err != nil {
		return fmt.Errorf("start decoder: %w", err)
	}
	defer func() {
		decoder.Drain()
		decoder.Stop()
	}()

	frameCount := 0
	bytesRead := int64(0)
	startTime := time.Now()

	fmt.Println("Decoding... (Ctrl+C to stop)")

	// Start file reader goroutine
	readDone := make(chan error, 1)
	go func() {
		chunk := make([]byte, *chunkSize)
		for {
			n, err := inFile.Read(chunk)
			if err == io.EOF {
				readDone <- nil
				return
			}
			if err != nil {
				readDone <- fmt.Errorf("read input: %w", err)
				return
			}

			bytesRead += int64(n)

			// Send chunk to decoder
			select {
			case decoder.GetInput() <- chunk[:n]:
			case <-ctx.Done():
				readDone <- nil
				return
			}
		}
	}()

	// Process decoded frames
	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-readDone:
			if err != nil {
				return err
			}
			// Continue processing remaining frames

		case frame, ok := <-decoder.GetOutput():
			if !ok {
				return nil
			}

			// Save frame to file
			if err := saveFrame(frame, frameCount); err != nil {
				return fmt.Errorf("save frame: %w", err)
			}

			frameCount++

			// Print progress
			elapsed := time.Since(startTime).Seconds()
			if elapsed > 0 {
				fps := float64(frameCount) / elapsed
				fmt.Printf("\rFrames: %d, FPS: %.1f, Read: %.2f MB  ",
					frameCount, fps, float64(bytesRead)/(1024*1024))
			}

			// Check max frames limit
			if *maxFrames > 0 && frameCount >= *maxFrames {
				fmt.Printf("\nReached max frames limit (%d)\n", *maxFrames)
				return nil
			}

		case change := <-decoder.GetResolutionChanges():
			fmt.Printf("\nResolution change: %dx%d (%s)\n",
				change.Width, change.Height,
				v4l2.PixelFormats[change.PixelFormat])

		case err := <-decoder.GetError():
			return fmt.Errorf("decoder error: %w", err)
		}
	}
}

func saveFrame(data []byte, index int) error {
	filename := filepath.Join(*outputDir, fmt.Sprintf("frame_%04d.raw", index))
	return os.WriteFile(filename, data, 0644)
}
