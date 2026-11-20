//go:build integration
// +build integration

package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

func TestIntegration_DeviceFramesIter(t *testing.T) {
	devPath := testDevice1
	if _, err := os.Stat(devPath); err != nil {
		t.Skipf("Test device %s not available: %v", devPath, err)
	}

	dev, err := device.Open(devPath, device.WithBufferSize(4))
	if err != nil {
		t.Fatalf("device.Open(%s): %v", devPath, err)
	}
	defer dev.Close()

	if !dev.Capability().IsVideoCaptureSupported() {
		t.Skipf("Device %s does not support single-planar video capture", devPath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer dev.Stop()

	const targetFrames = 5
	got := 0
	var cloned *device.Frame

	for frame, err := range dev.Frames(ctx) {
		if err != nil {
			t.Fatalf("Frames iteration error: %v", err)
		}
		if len(frame.Data) == 0 {
			t.Fatalf("frame %d has empty Data", got)
		}
		if got == 0 {
			cloned = frame.Clone()
		}
		got++
		if got >= targetFrames {
			break
		}
	}

	if got < targetFrames {
		t.Fatalf("captured %d frames, want %d", got, targetFrames)
	}
	if cloned == nil || len(cloned.Data) == 0 {
		t.Fatalf("expected non-empty cloned frame")
	}
}

func TestIntegration_VideoOutputWriteFrame(t *testing.T) {
	devPath := testDevice1
	if _, err := os.Stat(devPath); err != nil {
		t.Skipf("Test device %s not available: %v", devPath, err)
	}

	// Only run if device supports video output (e.g., v4l2loopback).
	fd, err := v4l2.OpenDevice(devPath, 0x2|0x800, 0)
	if err != nil {
		t.Skipf("cannot open %s: %v", devPath, err)
	}
	cap, err := v4l2.GetCapability(fd)
	_ = v4l2.CloseDevice(fd)
	if err != nil || !cap.IsVideoOutputSupported() {
		t.Skipf("Device %s does not support video output", devPath)
	}

	outDev, err := device.Open(
		devPath,
		device.WithVideoOutputEnabled(),
		device.WithBufferSize(4),
		device.WithPixFormat(v4l2.PixFormat{
			Width:       640,
			Height:      480,
			PixelFormat: v4l2.PixelFmtYUYV,
			Field:       v4l2.FieldNone,
		}),
	)
	if err != nil {
		t.Fatalf("open output device %s: %v", devPath, err)
	}
	defer outDev.Close()

	pixFmt, err := outDev.GetPixFormat()
	if err != nil {
		t.Fatalf("GetPixFormat: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := outDev.Start(ctx); err != nil {
		t.Fatalf("Start output: %v", err)
	}
	defer outDev.Stop()

	frame := make([]byte, pixFmt.SizeImage)
	for i := range frame {
		frame[i] = byte(i % 251)
	}

	for i := 0; i < 6; i++ {
		if err := outDev.WriteFrame(ctx, frame); err != nil {
			t.Fatalf("WriteFrame %d: %v", i, err)
		}
	}
}
