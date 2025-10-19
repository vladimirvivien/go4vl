// +build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/vladimirvivien/go4vl/device"
)

// TestDebug_OpenDevice is a minimal test to debug device opening
func TestDebug_OpenDevice(t *testing.T) {
	if testDevice1 == "" {
		t.Skip("No test device available")
	}

	t.Logf("Attempting to open device: %s", testDevice1)

	dev, err := device.Open(testDevice1, device.WithBufferSize(4))
	if err != nil {
		t.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	t.Logf("Device opened successfully")

	// Get capability info
	cap := dev.Capability()
	t.Logf("Device capabilities:")
	t.Logf("  Driver: %s", cap.Driver)
	t.Logf("  Card: %s", cap.Card)
	t.Logf("  Capabilities: 0x%x", cap.Capabilities)
	t.Logf("  Video Capture: %v", cap.IsVideoCaptureSupported())
	t.Logf("  Streaming: %v", cap.IsStreamingSupported())

	// Try to get output channel
	t.Logf("Getting output channel...")
	_ = dev.GetOutput()

	t.Logf("Starting capture...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer dev.Stop()

	t.Logf("Waiting for first frame...")
	select {
	case frame := <-dev.GetOutput():
		t.Logf("SUCCESS! Got frame of %d bytes", len(frame))
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for frame")
	}
}
