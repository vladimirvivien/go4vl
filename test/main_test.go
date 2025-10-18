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

// findExistingLoopbackDevices finds existing v4l2loopback devices specifically
func findExistingLoopbackDevices() []string {
	var devices []string

	// Check for v4l2loopback devices by reading /sys/class/video4linux/videoX/name
	for i := 0; i < 100; i++ {
		devPath := fmt.Sprintf("/dev/video%d", i)
		sysPath := fmt.Sprintf("/sys/class/video4linux/video%d/name", i)

		// Check if device exists
		if _, err := os.Stat(devPath); err != nil {
			continue
		}

		// Check if it's a loopback device
		nameBytes, err := os.ReadFile(sysPath)
		if err != nil {
			continue
		}

		name := strings.TrimSpace(string(nameBytes))
		// v4l2loopback devices have names starting with "Dummy" or containing "loopback"
		if strings.Contains(strings.ToLower(name), "loopback") ||
			strings.Contains(strings.ToLower(name), "dummy") ||
			strings.Contains(strings.ToLower(name), "go4vl_test") {
			devices = append(devices, devPath)
		}
	}

	return devices
}

// checkRequiredBinaries checks if required tools are available
func checkRequiredBinaries() []string {
	var missing []string

	requiredTools := map[string]string{
		"ffmpeg":    "Required for generating test patterns on loopback devices",
		"v4l2-ctl":  "Required for V4L2 device control and inspection",
	}

	for tool, description := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			missing = append(missing, fmt.Sprintf("  - %s: %s", tool, description))
		}
	}

	return missing
}

// TestMain sets up and tears down the testing environment
func TestMain(m *testing.M) {
	flag.Parse()

	// Check for required binaries first
	missingBinaries := checkRequiredBinaries()
	if len(missingBinaries) > 0 {
		log.Println("WARNING: Missing required tools for full test functionality:")
		for _, msg := range missingBinaries {
			log.Println(msg)
		}
		log.Println("")
		log.Println("To install on Debian/Ubuntu:")
		log.Println("  sudo apt-get install ffmpeg v4l-utils")
		log.Println("")
		log.Println("Some tests and benchmarks may be skipped.")
		log.Println("")
	}

	// Determine which device numbers to use
	// Priority: existing loopback devices > available slots > default (42,43)
	existingLoopbackDevices := findExistingLoopbackDevices()

	if len(existingLoopbackDevices) >= 2 {
		// Found existing loopback devices - we'll use their numbers
		testDevices = existingLoopbackDevices[:2]
		testDevice1 = testDevices[0]
		testDevice2 = testDevices[1]
		if *verbose {
			log.Printf("Found existing v4l2loopback devices: %v (will reload with correct params)", testDevices)
		}
	} else {
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
	}

	// Setup v4l2loopback if needed
	var exitCode int
	if !*skipSetup {
		if err := setupV4L2Loopback(); err != nil {
			log.Printf("ERROR: Failed to setup v4l2loopback: %v", err)
			log.Println("")
			log.Println("v4l2loopback is REQUIRED for tests and benchmarks.")
			log.Println("To fix this issue:")
			log.Println("  1. Install v4l2loopback: sudo apt-get install v4l2loopback-dkms")
			log.Println("  2. Run tests with sudo: sudo go test -v -tags=integration ./test/...")
			log.Println("  3. Or manually load module: sudo modprobe v4l2loopback devices=2 video_nr=42,43")
			log.Println("")
			log.Println("Tests will be skipped.")
			// Set devices to empty so all tests/benchmarks skip
			testDevice1 = ""
			testDevice2 = ""
			testDevices = []string{}
		}
		defer teardownV4L2Loopback()
	} else {
		// If skipping setup, verify loopback devices exist
		log.Println("Skipping v4l2loopback setup (-skip-setup flag)")
		if !isModuleLoaded() {
			log.Println("WARNING: v4l2loopback module not loaded, tests will skip")
			testDevice1 = ""
			testDevice2 = ""
			testDevices = []string{}
		} else if len(findExistingLoopbackDevices()) < 2 {
			log.Println("WARNING: Not enough v4l2loopback devices found, tests will skip")
			testDevice1 = ""
			testDevice2 = ""
			testDevices = []string{}
		} else {
			log.Printf("Using existing v4l2loopback devices: %v", testDevices)
		}
	}

	// Run tests
	exitCode = m.Run()

	os.Exit(exitCode)
}

// cleanupDevices attempts to ensure devices are not in use
func cleanupDevices() {
	// Stop any test patterns that might be running
	stopTestPatterns()

	// Also kill any orphaned ffmpeg processes that might be using the devices
	// This is important for cleanup between test runs
	exec.Command("sudo", "pkill", "-9", "ffmpeg").Run() // Ignore errors

	// Give devices time to be released
	time.Sleep(500 * time.Millisecond)
}

// setupV4L2Loopback loads the v4l2loopback module and starts test patterns
func setupV4L2Loopback() error {
	if *verbose {
		log.Println("Setting up v4l2loopback module...")
	}

	// ALWAYS unload and reload the module to ensure correct parameters
	// This is critical because the module might be loaded with wrong exclusive_caps setting
	if isModuleLoaded() {
		if *verbose {
			log.Println("Unloading existing v4l2loopback module to ensure clean state with correct parameters...")
		}

		// Unload module - try multiple times if needed
		for i := 0; i < 3; i++ {
			// Clean up before each attempt
			cleanupDevices()

			// Try to unload
			if err := exec.Command("sudo", "modprobe", "-r", "v4l2loopback").Run(); err != nil {
				if i < 2 {
					if *verbose {
						log.Printf("Unload attempt %d failed: %v, retrying...", i+1, err)
					}
					time.Sleep(1 * time.Second)
				} else {
					return fmt.Errorf("failed to unload v4l2loopback after 3 attempts: %v (make sure to kill all ffmpeg processes)", err)
				}
			} else {
				// Successfully unloaded
				if *verbose {
					log.Println("Successfully unloaded v4l2loopback module")
				}
				break
			}
		}
		time.Sleep(500 * time.Millisecond) // Wait for unload to complete
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

	if *verbose {
		log.Printf("Loading v4l2loopback with video_nr=%s", videoNr)
	}

	// Load v4l2loopback with specific device numbers
	// Note: exclusive_caps=0 allows both producer (ffmpeg) and consumer (tests) to open device
	cmd := exec.Command("sudo", "modprobe", "v4l2loopback",
		"devices=2",
		fmt.Sprintf("video_nr=%s", videoNr),
		"card_label=go4vl_test_1,go4vl_test_2",
		"exclusive_caps=0",
		"max_buffers=4",
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to load v4l2loopback: %v\nOutput: %s", err, output)
	}

	moduleLoaded = true
	time.Sleep(1 * time.Second) // Give kernel time to create devices

	// Verify devices were created
	for _, dev := range testDevices {
		if _, err := os.Stat(dev); err != nil {
			return fmt.Errorf("device %s not created after module load", dev)
		}
	}

	// Verify exclusive_caps parameter is correct
	exclusiveCapsBytes, err := os.ReadFile("/sys/module/v4l2loopback/parameters/exclusive_caps")
	if err == nil {
		exclusiveCaps := strings.TrimSpace(string(exclusiveCapsBytes))
		if *verbose {
			log.Printf("v4l2loopback exclusive_caps: %s", exclusiveCaps)
		}
		// Check that we have at least 2 devices with exclusive_caps=N (0)
		if !strings.Contains(exclusiveCaps, "N") {
			log.Printf("WARNING: exclusive_caps may not be set correctly: %s", exclusiveCaps)
			log.Println("Expected values containing 'N' for non-exclusive mode")
		}
	}

	if *verbose {
		log.Printf("Created test devices: %v", testDevices)
	}

	// Start test patterns for benchmarks to use
	// These will run for the entire test session
	if err := startTestPatterns(); err != nil {
		log.Printf("Warning: Failed to start test patterns: %v", err)
		log.Println("Benchmarks may fail without test patterns")
	}

	return nil
}

// teardownV4L2Loopback unloads the module and cleans up
func teardownV4L2Loopback() {
	if *verbose {
		log.Println("Tearing down v4l2loopback...")
	}

	// Clean up devices first
	cleanupDevices()

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

