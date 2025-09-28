// +build integration

package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestDevice_Open tests device.Open with various options
func TestDevice_Open(t *testing.T) {
	// Check if test devices are available
	if _, err := os.Stat(testDevice1); err != nil {
		t.Skipf("Test device %s not available, skipping test", testDevice1)
	}

	tests := []struct {
		name    string
		device  string
		options []device.Option
		wantErr bool
	}{
		{
			name:   "open with defaults",
			device: testDevice1,
		},
		{
			name:   "open with buffer size",
			device: testDevice1,
			options: []device.Option{
				device.WithBufferSize(4),
			},
		},
		{
			name:   "open with pixel format",
			device: testDevice1,
			options: []device.Option{
				device.WithPixFormat(v4l2.PixFormat{
					Width:       320,
					Height:      240,
					PixelFormat: v4l2.PixelFmtYUYV,
				}),
			},
		},
		{
			name:   "open with FPS",
			device: testDevice1,
			options: []device.Option{
				device.WithFPS(15),
			},
		},
		{
			name:   "open with all options",
			device: testDevice2,
			options: []device.Option{
				device.WithBufferSize(6),
				device.WithPixFormat(v4l2.PixFormat{
					Width:       1280,
					Height:      720,
					PixelFormat: v4l2.PixelFmtYUYV,
				}),
				device.WithFPS(25),
			},
		},
		{
			name:    "open non-existent device",
			device:  "/dev/video99",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip tests that require testDevice2 if it doesn't exist
			if tt.device == testDevice2 {
				if _, err := os.Stat(testDevice2); err != nil {
					t.Skipf("Test device %s not available", testDevice2)
				}
			}

			dev, err := device.Open(tt.device, tt.options...)

			// Handle permission errors gracefully
			if err != nil && strings.Contains(err.Error(), "permission denied") {
				t.Skipf("Permission denied for %s. Add user to video group: sudo usermod -a -G video $USER", tt.device)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				defer dev.Close()

				// Verify device is accessible
				if dev.Name() != tt.device {
					t.Errorf("Device name = %v, want %v", dev.Name(), tt.device)
				}
			}
		})
	}
}

// TestDevice_Capability tests device capability methods
func TestDevice_Capability(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	// Test Capability()
	cap := dev.Capability()

	t.Run("capability fields", func(t *testing.T) {
		// Check basic fields are populated
		if cap.Driver == "" {
			t.Error("Driver field is empty")
		}
		if cap.Card == "" {
			t.Error("Card field is empty")
		}
		if cap.BusInfo == "" {
			t.Error("BusInfo field is empty")
		}

		t.Logf("Device: %s, Driver: %s, Bus: %s", cap.Card, cap.Driver, cap.BusInfo)
	})

	t.Run("capability flags", func(t *testing.T) {
		// v4l2loopback should support these
		if !cap.IsVideoCaptureSupported() {
			t.Error("Expected video capture support")
		}
		if !cap.IsStreamingSupported() {
			t.Error("Expected streaming support")
		}
		if !cap.IsReadWriteSupported() {
			t.Error("Expected read/write support")
		}
	})

	t.Run("version info", func(t *testing.T) {
		version := cap.GetVersionInfo()
		versionStr := version.String()
		if versionStr == "" {
			t.Error("Version string is empty")
		}
		t.Logf("Version: %s", versionStr)
	})

	t.Run("capability descriptions", func(t *testing.T) {
		driverCaps := cap.GetDriverCapDescriptions()
		if len(driverCaps) == 0 {
			t.Error("No driver capability descriptions")
		}

		for _, desc := range driverCaps {
			t.Logf("  %s: 0x%08x", desc.Desc, desc.Cap)
		}
	})
}

// TestDevice_GetSetPixFormat tests pixel format operations
func TestDevice_GetSetPixFormat(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	// Get current format
	origFormat, err := dev.GetPixFormat()
	if err != nil {
		t.Fatalf("GetPixFormat() failed: %v", err)
	}

	t.Logf("Original format: %dx%d, PixelFormat: 0x%08x",
		origFormat.Width, origFormat.Height, origFormat.PixelFormat)

	// Test setting different formats
	testFormats := []v4l2.PixFormat{
		{
			Width:       640,
			Height:      480,
			PixelFormat: v4l2.PixelFmtYUYV,
			Field:       v4l2.FieldNone,
		},
		{
			Width:       320,
			Height:      240,
			PixelFormat: v4l2.PixelFmtYUYV,
			Field:       v4l2.FieldNone,
		},
		{
			Width:       1280,
			Height:      720,
			PixelFormat: v4l2.PixelFmtYUYV,
			Field:       v4l2.FieldNone,
		},
	}

	for _, format := range testFormats {
		t.Run(fmt.Sprintf("%dx%d", format.Width, format.Height), func(t *testing.T) {
			err := dev.SetPixFormat(format)
			if err != nil {
				t.Logf("SetPixFormat(%dx%d) not supported: %v",
					format.Width, format.Height, err)
				return
			}

			// Verify the format was set
			newFormat, err := dev.GetPixFormat()
			if err != nil {
				t.Errorf("GetPixFormat() after set failed: %v", err)
				return
			}

			// Driver may adjust dimensions
			t.Logf("Requested: %dx%d, Got: %dx%d",
				format.Width, format.Height,
				newFormat.Width, newFormat.Height)

			if newFormat.PixelFormat != format.PixelFormat {
				t.Errorf("Pixel format mismatch: got 0x%08x, want 0x%08x",
					newFormat.PixelFormat, format.PixelFormat)
			}
		})
	}
}

// TestDevice_GetFormatDescriptions tests format enumeration
func TestDevice_GetFormatDescriptions(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	formats, err := dev.GetFormatDescriptions()
	if err != nil {
		t.Fatalf("GetFormatDescriptions() failed: %v", err)
	}

	if len(formats) == 0 {
		t.Error("No format descriptions returned")
	}

	for i, fmt := range formats {
		t.Logf("Format %d: %s (0x%08x)", i, fmt.Description, fmt.PixelFormat)

		// Test GetFormatDescription by index
		desc, err := dev.GetFormatDescription(uint32(i))
		if err != nil {
			t.Errorf("GetFormatDescription(%d) failed: %v", i, err)
			continue
		}

		if desc.PixelFormat != fmt.PixelFormat {
			t.Errorf("Format mismatch at index %d: got 0x%08x, want 0x%08x",
				i, desc.PixelFormat, fmt.PixelFormat)
		}
	}
}

// TestDevice_FrameRate tests frame rate operations
func TestDevice_FrameRate(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	// Get current frame rate
	origFPS, err := dev.GetFrameRate()
	if err != nil {
		t.Fatalf("GetFrameRate() failed: %v", err)
	}
	t.Logf("Original FPS: %d", origFPS)

	// Test setting different frame rates
	testRates := []uint32{5, 15, 30, 60}

	for _, fps := range testRates {
		t.Run(fmt.Sprintf("%dFPS", fps), func(t *testing.T) {
			err := dev.SetFrameRate(fps)
			if err != nil {
				t.Logf("SetFrameRate(%d) not supported: %v", fps, err)
				return
			}

			// Verify the rate was set
			newFPS, err := dev.GetFrameRate()
			if err != nil {
				t.Errorf("GetFrameRate() after set failed: %v", err)
				return
			}

			// Driver may adjust FPS
			t.Logf("Requested: %d FPS, Got: %d FPS", fps, newFPS)
		})
	}
}

// TestDevice_StreamParams tests stream parameter operations
func TestDevice_StreamParams(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	// Get stream parameters
	params, err := dev.GetStreamParam()
	if err != nil {
		t.Fatalf("GetStreamParam() failed: %v", err)
	}

	t.Logf("Stream params - TimePerFrame: %d/%d",
		params.Capture.TimePerFrame.Numerator,
		params.Capture.TimePerFrame.Denominator)

	// Modify and set parameters
	newParams := v4l2.StreamParam{
		Capture: v4l2.CaptureParam{
			TimePerFrame: v4l2.Fract{
				Numerator:   1,
				Denominator: 15,
			},
		},
	}

	if err := dev.SetStreamParam(newParams); err != nil {
		t.Logf("SetStreamParam() not fully supported: %v", err)
	}
}

// TestDevice_Streaming tests the streaming lifecycle
func TestDevice_Streaming(t *testing.T) {
	dev, err := device.Open(testDevice1,
		device.WithBufferSize(4),
		device.WithPixFormat(v4l2.PixFormat{
			Width:       640,
			Height:      480,
			PixelFormat: v4l2.PixelFmtYUYV,
		}),
	)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	// Test BufferType
	if bufType := dev.BufferType(); bufType != v4l2.BufTypeVideoCapture {
		t.Errorf("BufferType() = %v, want %v", bufType, v4l2.BufTypeVideoCapture)
	}

	// Test BufferCount before streaming
	bufCount := dev.BufferCount()
	if bufCount == 0 {
		t.Error("BufferCount() returned 0 before streaming")
	}

	// Test MemIOType
	ioType := dev.MemIOType()
	if ioType != v4l2.IOTypeMMAP {
		t.Errorf("MemIOType() = %v, want %v", ioType, v4l2.IOTypeMMAP)
	}

	// Start streaming
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Check buffers after start
	buffers := dev.Buffers()
	if len(buffers) == 0 {
		t.Error("No buffers allocated after Start()")
	}

	// Capture frames
	frameCount := 0
	timeout := time.After(3 * time.Second)

	for frameCount < 10 {
		select {
		case frame := <-dev.GetOutput():
			if len(frame) == 0 {
				t.Log("Received empty frame")
				continue
			}
			frameCount++
			t.Logf("Frame %d: %d bytes", frameCount, len(frame))

			// Validate frame size
			expectedSize := 640 * 480 * 2 // YUYV format
			if len(frame) != expectedSize {
				t.Errorf("Frame size = %d, want %d", len(frame), expectedSize)
			}

		case <-timeout:
			t.Logf("Captured %d frames before timeout", frameCount)
			goto done
		}
	}

done:
	if frameCount == 0 {
		t.Error("No frames captured")
	}

	// Stop streaming
	if err := dev.Stop(); err != nil {
		t.Errorf("Stop() failed: %v", err)
	}

	// Verify can't capture after stop
	select {
	case frame := <-dev.GetOutput():
		if frame != nil {
			t.Error("Received frame after Stop()")
		}
	case <-time.After(100 * time.Millisecond):
		// Expected - no frames after stop
	}
}

// TestDevice_MultipleStartStop tests multiple start/stop cycles
func TestDevice_MultipleStartStop(t *testing.T) {
	dev, err := device.Open(testDevice2, device.WithBufferSize(2))
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	ctx := context.Background()

	for i := 0; i < 3; i++ {
		t.Logf("Cycle %d: Starting", i+1)

		if err := dev.Start(ctx); err != nil {
			t.Fatalf("Cycle %d: Start() failed: %v", i+1, err)
		}

		// Capture a few frames
		frames := 0
		timeout := time.After(1 * time.Second)

	capture:
		for frames < 3 {
			select {
			case frame := <-dev.GetOutput():
				if len(frame) > 0 {
					frames++
				}
			case <-timeout:
				break capture
			}
		}

		t.Logf("Cycle %d: Captured %d frames", i+1, frames)

		if err := dev.Stop(); err != nil {
			t.Errorf("Cycle %d: Stop() failed: %v", i+1, err)
		}

		t.Logf("Cycle %d: Stopped", i+1)
		time.Sleep(100 * time.Millisecond)
	}
}

// TestDevice_ContextCancellation tests context cancellation during streaming
func TestDevice_ContextCancellation(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	ctx, cancel := context.WithCancel(context.Background())

	if err := dev.Start(ctx); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Capture frames in background
	frameCount := 0
	done := make(chan bool)

	go func() {
		for range dev.GetOutput() {
			frameCount++
		}
		done <- true
	}()

	// Let it run briefly
	time.Sleep(500 * time.Millisecond)

	// Cancel context
	cancel()

	// Wait for streaming to stop
	select {
	case <-done:
		t.Logf("Stream stopped after cancellation, captured %d frames", frameCount)
	case <-time.After(2 * time.Second):
		t.Error("Stream didn't stop after context cancellation")
	}
}

// TestDevice_VideoInput tests video input operations
func TestDevice_VideoInput(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	// Get current input index
	index, err := dev.GetVideoInputIndex()
	if err != nil {
		// v4l2loopback may not support this
		t.Logf("GetVideoInputIndex() not supported: %v", err)
		return
	}

	t.Logf("Current video input index: %d", index)

	// Get input info
	info, err := dev.GetVideoInputInfo(0)
	if err != nil {
		t.Logf("GetVideoInputInfo(0) not supported: %v", err)
		return
	}

	t.Logf("Input 0: Name=%s, Type=%d", info.GetName(), info.GetInputType())
}

// TestDevice_CropCapability tests cropping operations
func TestDevice_CropCapability(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	cropCap, err := dev.GetCropCapability()
	if err != nil {
		t.Logf("GetCropCapability() not supported: %v", err)
		return
	}

	t.Logf("Crop bounds: %dx%d+%d+%d",
		cropCap.Bounds.Width, cropCap.Bounds.Height,
		cropCap.Bounds.Left, cropCap.Bounds.Top)

	t.Logf("Default rect: %dx%d+%d+%d",
		cropCap.DefaultRect.Width, cropCap.DefaultRect.Height,
		cropCap.DefaultRect.Left, cropCap.DefaultRect.Top)

	// Try to set crop
	newRect := v4l2.Rect{
		Left:   10,
		Top:    10,
		Width:  320,
		Height: 240,
	}

	if err := dev.SetCropRect(newRect); err != nil {
		t.Logf("SetCropRect() not supported: %v", err)
	}
}

// TestDevice_GetMediaInfo tests media controller info
func TestDevice_GetMediaInfo(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	info, err := dev.GetMediaInfo()
	if err != nil {
		// Most devices don't support media controller
		t.Logf("GetMediaInfo() not supported: %v", err)
		return
	}

	t.Logf("Media device: %s", info.Driver)
}

// TestDevice_FileDescriptor tests Fd() method
func TestDevice_FileDescriptor(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	fd := dev.Fd()
	if fd == 0 {
		t.Error("Fd() returned 0")
	}
	t.Logf("File descriptor: %d", fd)
}