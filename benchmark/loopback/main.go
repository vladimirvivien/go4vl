package loopback

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Device represents a v4l2loopback virtual device
type Device struct {
	DevicePath string
	DeviceNum  int
	FFmpegCmd  *exec.Cmd
}

// Setup creates a v4l2loopback device and starts ffmpeg feeding it
func Setup(deviceNum int, width, height, fps int, testPattern string) (*Device, error) {
	dev := &Device{
		DevicePath: fmt.Sprintf("/dev/video%d", deviceNum),
		DeviceNum:  deviceNum,
	}

	// Check if device already exists
	if _, err := os.Stat(dev.DevicePath); err == nil {
		return nil, fmt.Errorf("device %s already exists, unload v4l2loopback first", dev.DevicePath)
	}

	// Check prerequisites
	if err := checkPrerequisites(); err != nil {
		return nil, err
	}

	// Load v4l2loopback module
	if err := loadModule(deviceNum); err != nil {
		return nil, fmt.Errorf("failed to load v4l2loopback: %w", err)
	}

	// Wait for device to appear
	time.Sleep(500 * time.Millisecond)

	// Verify device exists
	if _, err := os.Stat(dev.DevicePath); err != nil {
		unloadModule()
		return nil, fmt.Errorf("device %s not created after loading module", dev.DevicePath)
	}

	// Start ffmpeg
	if err := dev.startFFmpeg(width, height, fps, testPattern); err != nil {
		unloadModule()
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Wait for ffmpeg to initialize
	time.Sleep(1 * time.Second)

	return dev, nil
}

// Close stops ffmpeg and unloads the v4l2loopback module
func (d *Device) Close() error {
	// Stop ffmpeg
	if d.FFmpegCmd != nil && d.FFmpegCmd.Process != nil {
		d.FFmpegCmd.Process.Kill()
		d.FFmpegCmd.Wait()
	}

	// Unload module
	return unloadModule()
}

// checkPrerequisites verifies ffmpeg and v4l2loopback are available
func checkPrerequisites() error {
	// Check ffmpeg
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found: install with 'sudo apt install ffmpeg'")
	}

	// Check v4l2loopback module
	cmd := exec.Command("modinfo", "v4l2loopback")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("v4l2loopback kernel module not installed: install with 'sudo apt install v4l2loopback-dkms'")
	}

	return nil
}

// loadModule loads the v4l2loopback kernel module
func loadModule(deviceNum int) error {
	cmd := exec.Command("modprobe", "v4l2loopback",
		fmt.Sprintf("video_nr=%d", deviceNum),
		"card_label=go4vl_benchmark",
		"exclusive_caps=1")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}

	return nil
}

// unloadModule unloads the v4l2loopback kernel module
func unloadModule() error {
	cmd := exec.Command("modprobe", "-r", "v4l2loopback")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}
	return nil
}

// startFFmpeg starts ffmpeg to feed the loopback device
func (d *Device) startFFmpeg(width, height, fps int, testPattern string) error {
	// Parse test pattern (e.g., "testsrc", "smptebars", "color=red")
	var source string
	if strings.HasPrefix(testPattern, "color=") {
		color := strings.TrimPrefix(testPattern, "color=")
		source = fmt.Sprintf("color=c=%s:s=%dx%d:r=%d", color, width, height, fps)
	} else if testPattern == "" || testPattern == "testsrc" {
		source = fmt.Sprintf("testsrc=size=%dx%d:rate=%d", width, height, fps)
	} else {
		source = fmt.Sprintf("%s=size=%dx%d:rate=%d", testPattern, width, height, fps)
	}

	args := []string{
		"-re",
		"-f", "lavfi",
		"-i", source,
		"-pix_fmt", "yuyv422",
		"-f", "v4l2",
		d.DevicePath,
	}

	d.FFmpegCmd = exec.Command("ffmpeg", args...)

	// Suppress ffmpeg output
	d.FFmpegCmd.Stdout = nil
	d.FFmpegCmd.Stderr = nil

	if err := d.FFmpegCmd.Start(); err != nil {
		return err
	}

	return nil
}

// IsAvailable checks if v4l2loopback and ffmpeg are installed
func IsAvailable() bool {
	return checkPrerequisites() == nil
}
