// h264_encoder demonstrates hardware H.264 encoding using V4L2 stateful codec API.
//
// This example:
// 1. Opens a camera device and an encoder device
// 2. Captures raw frames from the camera
// 3. Encodes them to H.264 using hardware acceleration
// 4. Writes the encoded output to a file
//
// Requirements:
// - A V4L2 hardware encoder device (e.g., /dev/video11 on Raspberry Pi)
// - A camera device (e.g., /dev/video0)
//
// Usage:
//
//	./h264_encoder -e /dev/video11 -c /dev/video0 -o output.h264 -t 10
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	encoderPath = flag.String("e", "/dev/video11", "encoder device path")
	cameraPath  = flag.String("c", "/dev/video0", "camera device path")
	outputFile  = flag.String("o", "output.h264", "output file path")
	duration    = flag.Int("t", 10, "recording duration in seconds")
	width       = flag.Uint("w", 1280, "video width")
	height      = flag.Uint("h", 720, "video height")
	fps         = flag.Uint("f", 30, "frames per second")
	bitrate     = flag.Uint("b", 2000000, "bitrate in bits per second")
)

func main() {
	flag.Parse()

	fmt.Println("H.264 Hardware Encoder Example")
	fmt.Println("==============================")
	fmt.Println()

	// Check if encoder device exists
	if _, err := os.Stat(*encoderPath); os.IsNotExist(err) {
		fmt.Printf("Encoder device %s not found.\n", *encoderPath)
		fmt.Println("\nNote: Hardware encoders are typically available on:")
		fmt.Println("  - Raspberry Pi (VideoCore): /dev/video11")
		fmt.Println("  - Intel Quick Sync: Check with v4l2-ctl --list-devices")
		fmt.Println("  - Rockchip RK3399: /dev/video1 (rkvenc)")
		os.Exit(1)
	}

	// Open and configure camera
	camera, err := openCamera()
	if err != nil {
		fmt.Printf("Failed to open camera: %v\n", err)
		os.Exit(1)
	}
	defer camera.Close()

	fmt.Printf("Camera: %s\n", camera.Capability().Card)
	fmt.Printf("  Format: %dx%d\n", *width, *height)

	// Open and configure encoder
	encoder, err := openEncoder()
	if err != nil {
		fmt.Printf("Failed to open encoder: %v\n", err)
		os.Exit(1)
	}
	defer encoder.Close()

	fmt.Printf("Encoder: %s\n", encoder.Capability().Card)
	fmt.Printf("  Input: %s\n", v4l2.PixelFormats[encoder.GetInputFormat().PixelFormat])
	fmt.Printf("  Output: %s\n", v4l2.PixelFormats[encoder.GetOutputFormat().PixelFormat])
	fmt.Printf("  Bitrate: %d bps\n", *bitrate)

	// Open output file
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	fmt.Printf("\nOutput: %s\n", *outputFile)
	fmt.Printf("Duration: %d seconds\n", *duration)
	fmt.Println()

	// Set up context with timeout and signal handling
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*duration)*time.Second)
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nInterrupted, stopping...")
		cancel()
	}()

	// Start encoding
	if err := runEncoder(ctx, camera, encoder, outFile); err != nil {
		fmt.Printf("Encoding failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nEncoding complete!")

	// Print file info
	if info, err := outFile.Stat(); err == nil {
		fmt.Printf("Output file size: %.2f MB\n", float64(info.Size())/(1024*1024))
	}
}

func openCamera() (*device.Device, error) {
	return device.Open(*cameraPath,
		device.WithPixFormat(v4l2.PixFormat{
			Width:       uint32(*width),
			Height:      uint32(*height),
			PixelFormat: v4l2.PixelFmtNV12, // Most encoders expect NV12
		}),
		device.WithFPS(uint32(*fps)),
		device.WithBufferSize(4),
	)
}

func openEncoder() (*device.Encoder, error) {
	return device.OpenEncoder(*encoderPath, device.EncoderConfig{
		InputFormat: v4l2.PixFormat{
			Width:       uint32(*width),
			Height:      uint32(*height),
			PixelFormat: v4l2.PixelFmtNV12,
		},
		OutputFormat: v4l2.PixFormat{
			PixelFormat: v4l2.PixelFmtH264,
		},
		Bitrate:          uint32(*bitrate),
		InputBufferCount: 4,
		OutputBufferCount: 4,
	})
}

func runEncoder(ctx context.Context, camera *device.Device, encoder *device.Encoder, outFile *os.File) error {
	// Start camera
	if err := camera.Start(ctx); err != nil {
		return fmt.Errorf("start camera: %w", err)
	}
	defer camera.Stop()

	// Start encoder
	if err := encoder.Start(ctx); err != nil {
		return fmt.Errorf("start encoder: %w", err)
	}
	defer func() {
		encoder.Drain()
		encoder.Stop()
	}()

	frameCount := 0
	bytesWritten := int64(0)
	startTime := time.Now()

	fmt.Println("Encoding... (Ctrl+C to stop)")

	// Process frames
	for {
		select {
		case <-ctx.Done():
			return nil

		case frame, ok := <-camera.GetOutput():
			if !ok {
				return nil
			}

			// Send frame to encoder
			select {
			case encoder.GetInput() <- frame:
				frameCount++
			case <-ctx.Done():
				return nil
			}

		case encoded, ok := <-encoder.GetOutput():
			if !ok {
				return nil
			}

			// Write encoded data to file
			n, err := outFile.Write(encoded)
			if err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			bytesWritten += int64(n)

			// Print progress
			elapsed := time.Since(startTime).Seconds()
			if frameCount > 0 && int(elapsed)%2 == 0 {
				fps := float64(frameCount) / elapsed
				fmt.Printf("\rFrames: %d, FPS: %.1f, Encoded: %.2f MB  ",
					frameCount, fps, float64(bytesWritten)/(1024*1024))
			}

		case err := <-encoder.GetError():
			return fmt.Errorf("encoder error: %w", err)
		}
	}
}
