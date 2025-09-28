package device_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

// Example demonstrates basic device usage for capturing a single frame.
// This example shows the minimal code needed to capture an image from a camera.
func Example_basicCapture() {
	// Open the device with minimal configuration
	dev, err := device.Open("/dev/video0", device.WithBufferSize(1))
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// Start streaming
	ctx := context.Background()
	if err := dev.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer dev.Stop()

	// Capture a single frame
	frame := <-dev.GetOutput()
	fmt.Printf("Captured frame: %d bytes\n", len(frame))

	// Save to file (assuming JPEG format)
	if err := os.WriteFile("snapshot.jpg", frame, 0644); err != nil {
		log.Fatal(err)
	}
}

// Example_configureDevice demonstrates how to configure device parameters.
// This shows various configuration options available when opening a device.
func Example_configureDevice() {
	// Open device with specific configuration
	dev, err := device.Open("/dev/video0",
		// Set resolution and format
		device.WithPixFormat(v4l2.PixFormat{
			Width:       1920,
			Height:      1080,
			PixelFormat: v4l2.PixelFmtMJPEG,
			Field:       v4l2.FieldNone,
		}),
		// Set frame rate
		device.WithFPS(30),
		// Set buffer count for smooth streaming
		device.WithBufferSize(4),
		// Explicitly set I/O type (though MMAP is default)
		device.WithIOType(v4l2.IOTypeMMAP),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// Query actual configuration (may differ from requested)
	pixFmt, _ := dev.GetPixFormat()
	fps, _ := dev.GetFrameRate()
	fmt.Printf("Actual format: %dx%d @ %d FPS\n", pixFmt.Width, pixFmt.Height, fps)
}

// Example_continuousCapture demonstrates continuous frame capture with cancellation.
// This pattern is useful for video streaming or recording applications.
func Example_continuousCapture() {
	// Open device for streaming
	dev, err := device.Open("/dev/video0",
		device.WithBufferSize(4),
		device.WithPixFormat(v4l2.PixFormat{
			Width:       640,
			Height:      480,
			PixelFormat: v4l2.PixelFmtYUYV,
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// Create context with timeout for controlled capture duration
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start streaming
	if err := dev.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer dev.Stop()

	// Process frames until timeout
	frameCount := 0
	for frame := range dev.GetOutput() {
		frameCount++
		// Process frame here (e.g., encode, analyze, forward)
		fmt.Printf("Frame %d: %d bytes\n", frameCount, len(frame))

		// Example: Stop after 100 frames
		if frameCount >= 100 {
			cancel()
			break
		}
	}
	fmt.Printf("Captured %d frames\n", frameCount)
}

// Example_deviceCapabilities shows how to query and display device capabilities.
// This is useful for understanding what a device supports before using it.
func Example_deviceCapabilities() {
	// Open device just to query capabilities
	dev, err := device.Open("/dev/video0")
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// Get device capabilities
	cap := dev.Capability()

	// Display basic info
	fmt.Printf("Driver: %s\n", cap.Driver)
	fmt.Printf("Card: %s\n", cap.Card)
	fmt.Printf("Bus: %s\n", cap.BusInfo)

	// Check specific capabilities
	if cap.IsVideoCaptureSupported() {
		fmt.Println("✓ Video capture supported")
	}
	if cap.IsStreamingSupported() {
		fmt.Println("✓ Streaming I/O supported")
	}

	// Get all supported formats
	formats, err := dev.GetFormatDescriptions()
	if err == nil {
		fmt.Println("Supported formats:")
		for _, f := range formats {
			fmt.Printf("  - %s\n", f.Description)
		}
	}

	// Get current format
	pixFmt, err := dev.GetPixFormat()
	if err == nil {
		fmt.Printf("Current format: %dx%d\n", pixFmt.Width, pixFmt.Height)
	}
}

// Example_errorHandling demonstrates proper error handling patterns.
// Shows how to handle common errors and recover gracefully.
func Example_errorHandling() {
	// Try to open device
	dev, err := device.Open("/dev/video0", device.WithBufferSize(2))
	if err != nil {
		// Handle specific error types
		if os.IsNotExist(err) {
			fmt.Println("Device not found - is a camera connected?")
			return
		}
		if os.IsPermission(err) {
			fmt.Println("Permission denied - try running with sudo or add user to video group")
			return
		}
		log.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try to set format
	err = dev.SetPixFormat(v4l2.PixFormat{
		Width:       1920,
		Height:      1080,
		PixelFormat: v4l2.PixelFmtMJPEG,
	})
	if err != nil {
		if err == v4l2.ErrorUnsupportedFeature {
			fmt.Println("Format not supported, using device defaults")
		} else {
			log.Printf("Warning: Could not set format: %v", err)
		}
		// Continue with default format
	}

	// Start streaming with context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		log.Fatalf("Failed to start streaming: %v", err)
	}
	defer func() {
		if err := dev.Stop(); err != nil {
			log.Printf("Warning: Error stopping stream: %v", err)
		}
	}()

	// Capture with error recovery
	for {
		select {
		case frame, ok := <-dev.GetOutput():
			if !ok {
				fmt.Println("Stream ended")
				return
			}
			if len(frame) == 0 {
				fmt.Println("Received empty frame, skipping")
				continue
			}
			// Process valid frame
			fmt.Printf("Got frame: %d bytes\n", len(frame))

		case <-ctx.Done():
			fmt.Println("Capture timeout")
			return
		}
	}
}

// ExampleDevice_Start shows how to start video streaming with context control.
func ExampleDevice_Start() {
	dev, err := device.Open("/dev/video0")
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// Use context for cancellation control
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start streaming
	if err := dev.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer dev.Stop()

	// Capture frames
	for frame := range dev.GetOutput() {
		fmt.Printf("Frame received: %d bytes\n", len(frame))
		// Cancel after first frame for this example
		cancel()
		break
	}
}

// ExampleDevice_SetPixFormat demonstrates setting a specific pixel format.
func ExampleDevice_SetPixFormat() {
	dev, err := device.Open("/dev/video0")
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// Configure for 720p MJPEG
	format := v4l2.PixFormat{
		Width:       1280,
		Height:      720,
		PixelFormat: v4l2.PixelFmtMJPEG,
		Field:       v4l2.FieldNone,
	}

	if err := dev.SetPixFormat(format); err != nil {
		log.Fatal(err)
	}

	// Verify the format was set
	actualFormat, err := dev.GetPixFormat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Format set to: %dx%d\n", actualFormat.Width, actualFormat.Height)
}

// ExampleDevice_GetFormatDescriptions shows how to enumerate supported formats.
func ExampleDevice_GetFormatDescriptions() {
	dev, err := device.Open("/dev/video0")
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// Get all supported formats
	formats, err := dev.GetFormatDescriptions()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Supported formats:")
	for i, format := range formats {
		fmt.Printf("%d. %s (0x%08x)\n", i+1, format.Description, format.PixelFormat)
	}
}