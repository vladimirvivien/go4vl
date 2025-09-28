// +build integration

package test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var (
	// Test devices - will be set dynamically to avoid conflicts
	testDevices = []string{}
	testDevice1 = "" // Primary test device (set from testDevices[0])
	testDevice2 = "" // Secondary test device (set from testDevices[1])

	// Global flags
	skipSetup   = flag.Bool("skip-setup", false, "Skip v4l2loopback setup (assume already configured)")
	keepRunning = flag.Bool("keep-running", false, "Keep v4l2loopback loaded after tests")
	useExisting = flag.Bool("use-existing", false, "Use existing v4l2loopback devices if available")
	verbose     = flag.Bool("verbose", false, "Enable verbose logging")

	// Module load parameters
	moduleLoaded = false
	ffmpegProcs  []*exec.Cmd
)

// findAvailableDeviceNumbers finds unused video device numbers
func findAvailableDeviceNumbers(count int) []int {
	available := []int{}
	// Start from 40 to avoid common devices (webcams are usually 0-10)
	for i := 40; i < 100 && len(available) < count; i++ {
		devPath := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(devPath); os.IsNotExist(err) {
			available = append(available, i)
		}
	}
	return available
}

// findExistingV4L2Devices finds existing V4L2 devices that can be used for testing
func findExistingV4L2Devices() []string {
	var devices []string

	// Check common device paths
	for i := 0; i < 20; i++ {
		devPath := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(devPath); err == nil {
			// Try to open it to verify it's a valid V4L2 device
			// Note: We're just checking if it exists and is accessible
			devices = append(devices, devPath)
			if len(devices) >= 2 {
				break // We only need 2 devices max
			}
		}
	}

	return devices
}

// TestMain sets up and tears down the testing environment
func TestMain(m *testing.M) {
	flag.Parse()


	// Find available device numbers to avoid conflicts
	availableDeviceNums := findAvailableDeviceNumbers(2)
	if len(availableDeviceNums) >= 2 {
		testDevices = []string{
			fmt.Sprintf("/dev/video%d", availableDeviceNums[0]),
			fmt.Sprintf("/dev/video%d", availableDeviceNums[1]),
		}
		testDevice1 = testDevices[0]
		testDevice2 = testDevices[1]
		if *verbose {
			log.Printf("Selected test devices: %s, %s", testDevice1, testDevice2)
		}
	} else {
		log.Println("Warning: Could not find 2 available device numbers, using defaults")
		testDevices = []string{"/dev/video42", "/dev/video43"}
		testDevice1 = testDevices[0]
		testDevice2 = testDevices[1]
	}

	// Setup v4l2loopback if needed
	var exitCode int
	if !*skipSetup {
		if err := setupV4L2Loopback(); err != nil {
			log.Printf("Warning: Failed to setup v4l2loopback: %v", err)
			log.Println("Looking for existing devices to use for testing...")

			// Try to find existing devices we can use
			existingDevices := findExistingV4L2Devices()
			if len(existingDevices) > 0 {
				testDevices = existingDevices
				if len(testDevices) > 0 {
					testDevice1 = testDevices[0]
				}
				if len(testDevices) > 1 {
					testDevice2 = testDevices[1]
				}
				log.Printf("Using existing devices for testing: %v", testDevices)
			} else {
				log.Println("No existing V4L2 devices found, tests will skip")
				log.Println("")
				log.Println("To run tests, either:")
				log.Println("  1. Run with sudo: sudo go test -v -tags=integration ./test/...")
				log.Println("  2. Setup v4l2loopback manually and use -skip-setup flag")
				log.Println("  3. Connect a USB webcam or other V4L2 device")
			}
		}
		defer teardownV4L2Loopback()
	} else {
		// If skipping setup, try to find existing devices
		log.Println("Skipping v4l2loopback setup, looking for existing devices...")
		existingDevices := findExistingV4L2Devices()
		if len(existingDevices) > 0 {
			testDevices = existingDevices
			if len(testDevices) > 0 {
				testDevice1 = testDevices[0]
			}
			if len(testDevices) > 1 {
				testDevice2 = testDevices[1]
			}
			log.Printf("Using existing devices: %v", testDevices)
		}
	}

	// Run tests
	exitCode = m.Run()

	os.Exit(exitCode)
}

// setupV4L2Loopback loads the v4l2loopback module and starts test patterns
func setupV4L2Loopback() error {
	if *verbose {
		log.Println("Setting up v4l2loopback module...")
	}

	// Check if module is already loaded
	if isModuleLoaded() {
		if *useExisting {
			log.Println("v4l2loopback already loaded, using existing setup")
			return startTestPatterns()
		}
		// Unload existing module to ensure clean state
		exec.Command("sudo", "modprobe", "-r", "v4l2loopback").Run()
		time.Sleep(500 * time.Millisecond)
	}

	// Extract device numbers from testDevices
	var deviceNums []string
	for _, dev := range testDevices {
		// Extract number from /dev/videoXX
		var num int
		fmt.Sscanf(dev, "/dev/video%d", &num)
		deviceNums = append(deviceNums, fmt.Sprintf("%d", num))
	}
	videoNr := strings.Join(deviceNums, ",")

	// Load v4l2loopback with specific device numbers
	cmd := exec.Command("sudo", "modprobe", "v4l2loopback",
		"devices=2",
		fmt.Sprintf("video_nr=%s", videoNr),
		"card_label=go4vl_test_1,go4vl_test_2",
		"exclusive_caps=1",
		"max_buffers=4",
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to load v4l2loopback: %v\nOutput: %s", err, output)
	}

	moduleLoaded = true
	time.Sleep(500 * time.Millisecond) // Give kernel time to create devices

	// Verify devices were created
	for _, dev := range testDevices {
		if _, err := os.Stat(dev); err != nil {
			return fmt.Errorf("device %s not created after module load", dev)
		}
	}

	if *verbose {
		log.Printf("Created test devices: %v", testDevices)
	}

	// Start test patterns
	return startTestPatterns()
}

// teardownV4L2Loopback unloads the module and cleans up
func teardownV4L2Loopback() {
	if *verbose {
		log.Println("Tearing down v4l2loopback...")
	}

	// Stop all test patterns
	stopTestPatterns()

	// Unload module if we loaded it
	if moduleLoaded && !*keepRunning {
		time.Sleep(500 * time.Millisecond) // Let devices settle
		if err := exec.Command("sudo", "modprobe", "-r", "v4l2loopback").Run(); err != nil {
			log.Printf("Warning: Failed to unload v4l2loopback: %v", err)
		} else if *verbose {
			log.Println("v4l2loopback module unloaded")
		}
	}
}

// isModuleLoaded checks if v4l2loopback is loaded
func isModuleLoaded() bool {
	cmd := exec.Command("lsmod")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "v4l2loopback")
}

// startTestPatterns starts ffmpeg test patterns for both devices
func startTestPatterns() error {
	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Println("ffmpeg not found, tests will run without test patterns")
		return nil
	}

	// Start test pattern for device 1 (standard test pattern)
	cmd1 := exec.Command("ffmpeg",
		"-f", "lavfi",
		"-i", "testsrc=size=640x480:rate=30",
		"-pix_fmt", "yuyv422",
		"-f", "v4l2",
		testDevice1,
	)
	if err := cmd1.Start(); err != nil {
		log.Printf("Warning: Failed to start test pattern on %s: %v", testDevice1, err)
	} else {
		ffmpegProcs = append(ffmpegProcs, cmd1)
		if *verbose {
			log.Printf("Started test pattern on %s (PID: %d)", testDevice1, cmd1.Process.Pid)
		}
	}

	// Start test pattern for device 2 (color bars)
	cmd2 := exec.Command("ffmpeg",
		"-f", "lavfi",
		"-i", "smptebars=size=1280x720:rate=25",
		"-pix_fmt", "yuyv422",
		"-f", "v4l2",
		testDevice2,
	)
	if err := cmd2.Start(); err != nil {
		log.Printf("Warning: Failed to start test pattern on %s: %v", testDevice2, err)
	} else {
		ffmpegProcs = append(ffmpegProcs, cmd2)
		if *verbose {
			log.Printf("Started color bars on %s (PID: %d)", testDevice2, cmd2.Process.Pid)
		}
	}

	// Give patterns time to initialize
	time.Sleep(1 * time.Second)
	return nil
}

// stopTestPatterns stops all running ffmpeg processes
func stopTestPatterns() {
	for _, cmd := range ffmpegProcs {
		if cmd != nil && cmd.Process != nil {
			if *verbose {
				log.Printf("Stopping ffmpeg PID %d", cmd.Process.Pid)
			}
			cmd.Process.Kill()
			cmd.Wait()
		}
	}
	ffmpegProcs = nil
}

