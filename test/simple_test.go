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
)

// TestSimple_DirectDeviceOpen tests opening a device without helper requirements
func TestSimple_DirectDeviceOpen(t *testing.T) {

	// Try to find any video device
	var devPath string
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(path); err == nil {
			devPath = path
			break
		}
	}

	if devPath == "" {
		t.Skip("No video devices found")
	}

	t.Logf("Testing with device: %s", devPath)

	// Try to open the device
	dev, err := device.Open(devPath)
	if err != nil {
		// This might fail if device is busy or not a capture device
		t.Logf("Failed to open %s: %v", devPath, err)

		// Check for permission errors
		errStr := err.Error()
		if strings.Contains(errStr, "permission denied") || strings.Contains(errStr, "Permission denied") {
			t.Log("")
			t.Log("🔒 Permission denied. To fix this:")
			t.Log("  1. Add user to video group: sudo usermod -a -G video $USER")
			t.Log("  2. Logout and login again (or run: newgrp video)")
			t.Log("  3. Or run tests with sudo (not recommended)")
			t.Log("")
		}
		t.Skip("Device not available for testing")
	}
	defer dev.Close()

	// Check basic properties
	cap := dev.Capability()
	t.Logf("Device opened successfully!")
	t.Logf("  Driver: %s", cap.Driver)
	t.Logf("  Card: %s", cap.Card)
	t.Logf("  Bus: %s", cap.BusInfo)
	t.Logf("  Version: %s", cap.GetVersionInfo().String())

	// Check capabilities
	if cap.IsVideoCaptureSupported() {
		t.Log("  ✓ Video capture supported")
	} else {
		t.Log("  ✗ Video capture NOT supported")
	}

	if cap.IsStreamingSupported() {
		t.Log("  ✓ Streaming I/O supported")
	} else {
		t.Log("  ✗ Streaming I/O NOT supported")
	}

	// Try to get format
	if pixFmt, err := dev.GetPixFormat(); err == nil {
		t.Logf("  Current format: %dx%d", pixFmt.Width, pixFmt.Height)
	}
}

// TestSimple_StreamIfPossible attempts to stream from a device if available
func TestSimple_StreamIfPossible(t *testing.T) {

	// Find a device
	var devPath string
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(path); err == nil {
			// Try to open it
			if dev, err := device.Open(path); err == nil {
				if dev.Capability().IsVideoCaptureSupported() && dev.Capability().IsStreamingSupported() {
					dev.Close()
					devPath = path
					break
				}
				dev.Close()
			}
		}
	}

	if devPath == "" {
		t.Skip("No suitable capture device found")
	}

	t.Logf("Found capture device: %s", devPath)

	// Open and try to stream
	dev, err := device.Open(devPath, device.WithBufferSize(2))
	if err != nil {
		t.Skipf("Cannot open device: %v", err)
	}
	defer dev.Close()

	// Start streaming with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		t.Logf("Cannot start streaming: %v", err)
		t.Skip("Device cannot stream (might be in use)")
	}
	defer dev.Stop()

	// Try to get one frame
	select {
	case frame := <-dev.GetOutput():
		if len(frame) > 0 {
			t.Logf("✓ Successfully captured frame: %d bytes", len(frame))
		} else {
			t.Log("✗ Received empty frame")
		}
	case <-time.After(1 * time.Second):
		t.Log("✗ Timeout waiting for frame (device might not be producing frames)")
	}
}