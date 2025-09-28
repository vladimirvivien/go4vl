package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// V4L2LoopbackDevice manages a v4l2loopback device for testing
type V4L2LoopbackDevice struct {
	Path         string
	FFmpegCmd    *exec.Cmd
	TestPattern  string
	Width        int
	Height       int
	FPS          int
	PixelFormat  string
}

// NewV4L2LoopbackDevice creates a new test device manager
func NewV4L2LoopbackDevice(devicePath string) *V4L2LoopbackDevice {
	return &V4L2LoopbackDevice{
		Path:        devicePath,
		TestPattern: "testsrc",
		Width:       640,
		Height:      480,
		FPS:         30,
		PixelFormat: "yuyv422",
	}
}

// CheckV4L2LoopbackModule checks if v4l2loopback kernel module is loaded
func CheckV4L2LoopbackModule(t *testing.T) bool {
	cmd := exec.Command("lsmod")
	output, err := cmd.Output()
	if err != nil {
		t.Logf("Failed to run lsmod: %v", err)
		return false
	}
	return strings.Contains(string(output), "v4l2loopback")
}

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

	// Common test device paths
	testDevices := []string{
		"/dev/video20", // v4l2loopback default
		"/dev/video21",
		"/dev/video10", // vivid default
		"/dev/video11",
	}

	for _, device := range testDevices {
		if _, err := os.Stat(device); err == nil {
			// Check if it's a v4l2 device
			cmd := exec.Command("v4l2-ctl", "-d", device, "--info")
			if err := cmd.Run(); err == nil {
				t.Logf("Found test device: %s", device)
				return device
			}
		}
	}

	// Try to find any v4l2loopback device
	devices, _ := filepath.Glob("/dev/video*")
	for _, device := range devices {
		cmd := exec.Command("v4l2-ctl", "-d", device, "--info")
		output, err := cmd.Output()
		if err == nil && strings.Contains(string(output), "v4l2loopback") {
			t.Logf("Found v4l2loopback device: %s", device)
			return device
		}
	}

	return ""
}

// SetupV4L2Loopback ensures v4l2loopback is available and returns device path
func SetupV4L2Loopback(t *testing.T) string {
	// Check if module is loaded
	if !CheckV4L2LoopbackModule(t) {
		// Try to load it if we're in CI
		if os.Getenv("CI") == "true" {
			t.Log("Attempting to load v4l2loopback module...")
			cmd := exec.Command("sudo", "modprobe", "v4l2loopback",
				"devices=1", "video_nr=20", "exclusive_caps=1",
				"card_label=go4vl_test")
			if err := cmd.Run(); err != nil {
				t.Skipf("Failed to load v4l2loopback: %v", err)
			}
			time.Sleep(500 * time.Millisecond) // Give kernel time to create device
		} else {
			t.Skip("v4l2loopback module not loaded. Run: sudo modprobe v4l2loopback")
		}
	}

	// Find the test device
	device := FindTestDevice(t)
	if device == "" {
		t.Skip("No v4l2loopback device found")
	}

	return device
}

// Start starts feeding test pattern to the device
func (d *V4L2LoopbackDevice) Start(t TestLogger) error {
	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Logf("ffmpeg not found, skipping test pattern generation")
		return nil
	}

	// Build ffmpeg command
	args := []string{
		"-re", // Read input at native frame rate
		"-f", "lavfi",
		"-i", fmt.Sprintf("%s=size=%dx%d:rate=%d",
			d.TestPattern, d.Width, d.Height, d.FPS),
		"-pix_fmt", d.PixelFormat,
		"-f", "v4l2",
		d.Path,
	}

	d.FFmpegCmd = exec.Command("ffmpeg", args...)

	// Start in background
	if err := d.FFmpegCmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Give it time to initialize
	time.Sleep(500 * time.Millisecond)

	t.Logf("Started test pattern on %s (%dx%d @ %d FPS)",
		d.Path, d.Width, d.Height, d.FPS)

	return nil
}

// Stop stops the test pattern generator
func (d *V4L2LoopbackDevice) Stop() error {
	if d.FFmpegCmd != nil && d.FFmpegCmd.Process != nil {
		d.FFmpegCmd.Process.Kill()
		d.FFmpegCmd.Wait()
		d.FFmpegCmd = nil
	}
	return nil
}

// StartTestPattern is a convenience function to start test pattern on a device
func StartTestPattern(t TestLogger, devicePath string) func() {
	device := NewV4L2LoopbackDevice(devicePath)

	if err := device.Start(t); err != nil {
		t.Logf("Failed to start test pattern: %v", err)
		return func() {}
	}

	return func() {
		device.Stop()
	}
}

// RequireV4L2Testing checks all requirements and skips test if not met
func RequireV4L2Testing(t *testing.T) string {
	t.Helper()

	// Check for v4l2-ctl tool
	if _, err := exec.LookPath("v4l2-ctl"); err != nil {
		t.Skip("v4l2-ctl not found. Install with: apt-get install v4l-utils")
	}

	// Setup and find device
	return SetupV4L2Loopback(t)
}

// TestPatternType represents different test patterns
type TestPatternType string

const (
	TestPatternBars     TestPatternType = "smptebars"
	TestPatternColor    TestPatternType = "color"
	TestPatternNoise    TestPatternType = "testsrc"
	TestPatternMandel   TestPatternType = "mandelbrot"
	TestPatternLife     TestPatternType = "life"
)

// GenerateTestVideo generates a test video file for testing
func GenerateTestVideo(t *testing.T, filename string, duration int) error {
	cmd := exec.Command("ffmpeg",
		"-f", "lavfi",
		"-i", "testsrc=size=640x480:rate=30",
		"-t", fmt.Sprintf("%d", duration),
		"-pix_fmt", "yuv420p",
		"-y", // Overwrite output
		filename,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("ffmpeg output: %s", output)
		return fmt.Errorf("failed to generate test video: %w", err)
	}

	return nil
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
		// Y should typically be between 16-235 in video range
		if y < 10 || y > 245 {
			invalidSamples++
		}
		samplesChecked++
	}

	if invalidSamples > samplesChecked/2 {
		t.Log("Warning: Many Y values outside typical range")
	}
}