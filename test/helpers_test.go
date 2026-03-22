// +build integration

package test

import (
	"os"
	"strings"
	"testing"

	"github.com/vladimirvivien/go4vl/device"
)

// TestLogger is an interface that both *testing.T and *testing.B satisfy
type TestLogger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Skip(args ...interface{})
	Skipf(format string, args ...interface{})
	Helper()
}

// FindTestDevice finds an available V4L2 test device
func FindTestDevice(t TestLogger) string {
	// Check environment variable first
	if envDevice := os.Getenv("V4L2_TEST_DEVICE"); envDevice != "" {
		if _, err := os.Stat(envDevice); err == nil {
			t.Logf("Using device from V4L2_TEST_DEVICE: %s", envDevice)
			return envDevice
		}
	}

	// Use global testDevice1 if it was set by TestMain and exists
	if testDevice1 != "" {
		if _, err := os.Stat(testDevice1); err == nil {
			t.Logf("Using test device from TestMain: %s", testDevice1)
			return testDevice1
		}
	}

	// Search common device paths
	commonDevices := []string{
		"/dev/video0",
		"/dev/video1",
		"/dev/video2",
		"/dev/video3",
		"/dev/video10",
		"/dev/video11",
	}

	for _, device := range commonDevices {
		if _, err := os.Stat(device); err == nil {
			t.Logf("Found device: %s", device)
			return device
		}
	}

	return ""
}

// RequireV4L2Testing finds a test device or skips the test
func RequireV4L2Testing(t *testing.T) string {
	t.Helper()

	device := FindTestDevice(t)
	if device == "" {
		t.Skip("No V4L2 device available")
	}
	return device
}

// OpenDeviceOrSkip opens a V4L2 device and skips the test if the device
// is busy or permission is denied. Use this instead of device.Open() directly
// in integration tests to handle transient device contention gracefully.
func OpenDeviceOrSkip(t *testing.T, path string, options ...device.Option) *device.Device {
	t.Helper()
	dev, err := device.Open(path, options...)
	if err != nil {
		if strings.Contains(err.Error(), "device or resource busy") ||
			strings.Contains(err.Error(), "permission denied") {
			t.Skipf("Device %s unavailable: %v", path, err)
		}
		t.Fatalf("Failed to open device %s: %v", path, err)
	}
	t.Cleanup(func() {
		dev.Close()
	})
	return dev
}

// CompareFrames compares two frames for testing
func CompareFrames(frame1, frame2 []byte, tolerance float64) bool {
	if len(frame1) != len(frame2) {
		return false
	}

	differences := 0
	for i := range frame1 {
		if frame1[i] != frame2[i] {
			differences++
		}
	}

	diffPercent := float64(differences) / float64(len(frame1)) * 100
	return diffPercent <= tolerance
}

// ValidateYUYVFrame performs basic validation on a YUYV frame
func ValidateYUYVFrame(t TestLogger, frame []byte, width, height uint32) {
	expectedSize := int(width * height * 2)
	if len(frame) != expectedSize {
		t.Errorf("Invalid YUYV frame size: got %d, expected %d", len(frame), expectedSize)
		return
	}

	// Check for all zeros (black frame might be valid, but often indicates error)
	allZero := true
	for i := 0; i < len(frame) && i < 1000; i++ {
		if frame[i] != 0 {
			allZero = false
			break
		}
	}

	if allZero {
		t.Log("Warning: Frame appears to be all zeros")
	}

	// Basic YUYV validation - Y values should be in reasonable range
	samplesChecked := 0
	invalidSamples := 0
	for i := 0; i < len(frame) && samplesChecked < 100; i += 2 {
		y := frame[i]
		if y < 10 || y > 245 {
			invalidSamples++
		}
		samplesChecked++
	}

	if invalidSamples > samplesChecked/2 {
		t.Log("Warning: Many Y values outside typical range")
	}
}
