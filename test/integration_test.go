// +build integration

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)


func TestIntegration_DeviceOpen(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device %s: %v", devPath, err)
	}
	defer dev.Close()

	// Verify device properties
	cap := dev.Capability()
	t.Logf("Device: %s", cap.Card)
	t.Logf("Driver: %s", cap.Driver)
	t.Logf("Bus: %s", cap.BusInfo)

	if !cap.IsStreamingSupported() {
		t.Error("Device should support streaming")
	}
}

func TestIntegration_DeviceCapabilities(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	cap := dev.Capability()

	// Log capabilities for debugging
	t.Logf("Capabilities: 0x%08x", cap.Capabilities)

	// Check expected capabilities
	if cap.IsVideoCaptureSupported() {
		t.Log("✓ Video capture supported")
	}

	if cap.IsStreamingSupported() {
		t.Log("✓ Streaming I/O supported")
	}

	// Get format descriptions
	formats, err := dev.GetFormatDescriptions()
	if err != nil {
		t.Logf("Failed to get format descriptions: %v", err)
	} else {
		t.Log("Supported formats:")
		for _, fmt := range formats {
			t.Logf("  - %s (0x%08x)", fmt.Description, fmt.PixelFormat)
		}
	}
}

func TestIntegration_BasicStreaming(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	// Start test pattern
	stopPattern := StartTestPattern(t, devPath)
	defer stopPattern()

	dev, err := device.Open(devPath,
		device.WithBufferSize(4),
		device.WithPixFormat(v4l2.PixFormat{
			Width:       640,
			Height:      480,
			PixelFormat: v4l2.PixelFmtYUYV,
		}),
	)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Start streaming
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		t.Fatalf("Failed to start streaming: %v", err)
	}
	defer dev.Stop()

	// Capture frames
	frameCount := 0
	emptyFrames := 0
	timeout := time.After(3 * time.Second)

	for {
		select {
		case frame := <-dev.GetOutput():
			if len(frame) == 0 {
				emptyFrames++
				continue
			}
			frameCount++

			// Validate frame
			ValidateYUYVFrame(t, frame, 640, 480)
			t.Logf("Frame %d: %d bytes (valid YUYV)", frameCount, len(frame))

			if frameCount >= 10 {
				goto done
			}
		case <-timeout:
			if frameCount == 0 {
				t.Fatal("No frames captured")
			}
			goto done
		}
	}

done:
	t.Logf("Captured %d frames, %d empty frames", frameCount, emptyFrames)

	if frameCount < 5 {
		t.Errorf("Expected at least 5 frames, got %d", frameCount)
	}
}

func TestIntegration_FormatNegotiation(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Try different formats
	testFormats := []v4l2.PixFormat{
		{
			Width:       640,
			Height:      480,
			PixelFormat: v4l2.PixelFmtYUYV,
		},
		{
			Width:       320,
			Height:      240,
			PixelFormat: v4l2.PixelFmtYUYV,
		},
		{
			Width:       1280,
			Height:      720,
			PixelFormat: v4l2.PixelFmtMJPEG,
		},
	}

	for _, format := range testFormats {
		t.Run(fmt.Sprintf("%dx%d", format.Width, format.Height), func(t *testing.T) {
			err := dev.SetPixFormat(format)
			if err != nil {
				t.Logf("Format not supported: %v", err)
				return
			}

			// Get actual format
			actualFormat, err := dev.GetPixFormat()
			if err != nil {
				t.Errorf("Failed to get pixel format: %v", err)
				return
			}

			t.Logf("Requested: %dx%d, Got: %dx%d",
				format.Width, format.Height,
				actualFormat.Width, actualFormat.Height)
		})
	}
}

func TestIntegration_FrameRateControl(t *testing.T) {
	devPath := RequireV4L2Testing(t)

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Test different frame rates
	testRates := []uint32{15, 30, 60}

	for _, fps := range testRates {
		t.Run(fmt.Sprintf("%dFPS", fps), func(t *testing.T) {
			err := dev.SetFrameRate(fps)
			if err != nil {
				t.Logf("Frame rate %d not supported: %v", fps, err)
				return
			}

			actualFPS, err := dev.GetFrameRate()
			if err != nil {
				t.Errorf("Failed to get frame rate: %v", err)
				return
			}

			t.Logf("Requested: %d FPS, Got: %d FPS", fps, actualFPS)
		})
	}
}

func TestIntegration_StopStart(t *testing.T) {
	devPath := RequireV4L2Testing(t)
	stopPattern := StartTestPattern(t, devPath)
	defer stopPattern()

	dev, err := device.Open(devPath, device.WithBufferSize(2))
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	ctx := context.Background()

	// Start-stop cycle multiple times
	for i := 0; i < 3; i++ {
		t.Logf("Cycle %d: Starting", i+1)

		if err := dev.Start(ctx); err != nil {
			t.Fatalf("Cycle %d: Failed to start: %v", i+1, err)
		}

		// Capture a few frames
		frameCount := 0
		timeout := time.After(1 * time.Second)

	captureLoop:
		for {
			select {
			case frame := <-dev.GetOutput():
				if len(frame) > 0 {
					frameCount++
				}
				if frameCount >= 3 {
					break captureLoop
				}
			case <-timeout:
				break captureLoop
			}
		}

		t.Logf("Cycle %d: Captured %d frames", i+1, frameCount)

		if err := dev.Stop(); err != nil {
			t.Errorf("Cycle %d: Failed to stop: %v", i+1, err)
		}

		t.Logf("Cycle %d: Stopped", i+1)
		time.Sleep(100 * time.Millisecond)
	}
}

func TestIntegration_ContextCancellation(t *testing.T) {
	devPath := RequireV4L2Testing(t)
	stopPattern := StartTestPattern(t, devPath)
	defer stopPattern()

	dev, err := device.Open(devPath)
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	if err := dev.Start(ctx); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}

	// Capture frames in goroutine
	framesChan := make(chan int)
	go func() {
		count := 0
		for range dev.GetOutput() {
			count++
		}
		framesChan <- count
	}()

	// Let it run briefly
	time.Sleep(500 * time.Millisecond)

	// Cancel context
	cancel()

	// Wait for streaming to stop
	select {
	case frames := <-framesChan:
		t.Logf("Captured %d frames before cancellation", frames)
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for stream to stop after context cancellation")
	}

	// Verify device stopped
	if err := dev.Stop(); err != nil {
		t.Logf("Stop after cancel returned: %v", err)
	}
}

// ============================================================================
// BENCHMARKS
// ============================================================================
// Benchmarks compare GetOutput() vs GetFrames() APIs.
//
// Results: GetFrames() is 2% faster and uses 600x less memory (1 KB vs 600 KB
// per frame) due to buffer pooling, dramatically reducing GC pressure.
//
// Run individually with -run=^$ (v4l2 driver limitation):
//   sudo go test -tags=integration -bench=BenchmarkIntegration_GetOutput -benchmem -run=^$ ./test
//   sudo go test -tags=integration -bench=BenchmarkIntegration_GetFrames -benchmem -run=^$ ./test
// ============================================================================

// BenchmarkIntegration_GetOutput benchmarks the legacy GetOutput() API
// This benchmark measures the complete end-to-end pipeline including:
// - V4L2 buffer dequeue (syscall)
// - Memory copy from mmap buffer to output buffer (direct allocation)
// - Channel send operation
// - V4L2 buffer re-queue (syscall)
func BenchmarkIntegration_GetOutput(b *testing.B) {
	if testDevice1 == "" {
		b.Skip("No test device available")
	}

	dev, err := device.Open(testDevice1, device.WithBufferSize(4))
	if err != nil {
		b.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	_ = dev.GetOutput()

	ctx := context.Background()
	if err := dev.Start(ctx); err != nil {
		b.Fatalf("Failed to start: %v", err)
	}
	defer dev.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		frame := <-dev.GetOutput()
		_ = frame
	}
}

// BenchmarkIntegration_GetFrames benchmarks the optimized GetFrames() API with pooling
// This benchmark measures the complete end-to-end pipeline including:
// - V4L2 buffer dequeue (syscall)
// - Memory copy from mmap buffer to pooled buffer
// - Channel send operation (Frame object)
// - V4L2 buffer re-queue (syscall)
// - Frame.Release() returning buffer to pool
func BenchmarkIntegration_GetFrames(b *testing.B) {
	if testDevice1 == "" {
		b.Skip("No test device available")
	}

	dev, err := device.Open(testDevice1, device.WithBufferSize(4))
	if err != nil {
		b.Fatalf("Failed to open device: %v", err)
	}

	_ = dev.GetFrames()

	ctx := context.Background()
	if err := dev.Start(ctx); err != nil {
		b.Fatalf("Failed to start: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		frame := <-dev.GetFrames()
		_ = frame.Data
		frame.Release()
	}

	b.StopTimer()

	dev.Stop()
	dev.Close()
}

// BenchmarkIntegration_GetFrames_WithMetadata benchmarks GetFrames() with metadata access
// This benchmark tests the realistic usage pattern where users access frame metadata
// in addition to the frame data, demonstrating zero-cost metadata exposure.
func BenchmarkIntegration_GetFrames_WithMetadata(b *testing.B) {
	if testDevice1 == "" {
		b.Skip("No test device available")
	}

	dev, err := device.Open(testDevice1, device.WithBufferSize(4))
	if err != nil {
		b.Fatalf("Failed to open device: %v", err)
	}

	_ = dev.GetFrames()

	ctx := context.Background()
	if err := dev.Start(ctx); err != nil {
		b.Fatalf("Failed to start: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		frame := <-dev.GetFrames()

		_ = frame.Data
		_ = frame.Timestamp
		_ = frame.Sequence
		_ = frame.IsKeyFrame()

		frame.Release()
	}

	b.StopTimer()

	dev.Stop()
	dev.Close()
}

// BenchmarkIntegration_Comparison was removed because running multiple device
// benchmarks in the same process causes v4l2 driver-level conflicts.
// Use the individual benchmarks instead:
//   - BenchmarkIntegration_GetOutput for the legacy API
//   - BenchmarkIntegration_GetFrames for the pooled API
//
// To compare both APIs, run them separately with -run=^$ to skip tests:
//   sudo go test -tags=integration -bench=BenchmarkIntegration_GetOutput -benchmem -run=^$ ./test
//   sudo go test -tags=integration -bench=BenchmarkIntegration_GetFrames -benchmem -run=^$ ./test
//
// Note: Each runs in a separate process, avoiding device conflicts.