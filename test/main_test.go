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
	// Test devices - set by TestMain based on flags and discovery
	testDevices = []string{}
	testDevice1 = "" // Primary test device
	testDevice2 = "" // Secondary test device

	// Flags
	useDevice          = flag.String("use-device", "", "Real V4L2 device: 'auto' to discover or device path (e.g. /dev/video0)")
	useDeviceEmulation = flag.String("use-device-emulation", "", "Loopback device: 'auto' to discover or comma-separated paths (e.g. /dev/video42,/dev/video43)")
	keepRunning        = flag.Bool("keep-running", false, "Keep v4l2loopback loaded after tests")
	verbose            = flag.Bool("verbose", false, "Enable verbose logging")

	// Track whether we loaded v4l2loopback (so we can unload on teardown)
	loopbackLoaded bool

)

// findExistingV4L2Devices finds existing V4L2 devices that can be used for testing
func findExistingV4L2Devices() []string {
	var devices []string
	for i := 0; i < 20; i++ {
		devPath := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(devPath); err == nil {
			devices = append(devices, devPath)
			if len(devices) >= 2 {
				break
			}
		}
	}
	return devices
}

// findLoopbackDevices finds v4l2loopback or vivid devices via sysfs
func findLoopbackDevices() []string {
	var devices []string
	for i := 0; i < 100; i++ {
		devPath := fmt.Sprintf("/dev/video%d", i)
		sysPath := fmt.Sprintf("/sys/class/video4linux/video%d/name", i)

		if _, err := os.Stat(devPath); err != nil {
			continue
		}

		nameBytes, err := os.ReadFile(sysPath)
		if err != nil {
			continue
		}

		name := strings.ToLower(strings.TrimSpace(string(nameBytes)))
		if strings.Contains(name, "loopback") ||
			strings.Contains(name, "dummy") ||
			strings.Contains(name, "go4vl_test") ||
			strings.Contains(name, "vivid") {
			devices = append(devices, devPath)
		}
	}
	return devices
}

// isLoopbackLoaded checks if v4l2loopback kernel module is loaded
func isLoopbackLoaded() bool {
	output, err := exec.Command("lsmod").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "v4l2loopback")
}

// isCI returns true if running in a CI environment
func isCI() bool {
	return os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true"
}

// findAvailableDeviceNumbers finds unused video device numbers
func findAvailableDeviceNumbers(count int) []int {
	var available []int
	for i := 40; i < 100 && len(available) < count; i++ {
		devPath := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(devPath); os.IsNotExist(err) {
			available = append(available, i)
		}
	}
	return available
}

// setupLoopback loads the v4l2loopback module and starts ffmpeg test patterns
func setupLoopback() error {
	if *verbose {
		log.Println("Setting up v4l2loopback...")
	}

	// If already loaded, just find devices
	if isLoopbackLoaded() {
		if *verbose {
			log.Println("v4l2loopback already loaded")
		}
		return nil
	}

	// Determine device numbers
	deviceNums := findAvailableDeviceNumbers(2)
	if len(deviceNums) < 2 {
		deviceNums = []int{42, 43}
	}

	videoNr := fmt.Sprintf("%d,%d", deviceNums[0], deviceNums[1])

	cmd := exec.Command("sudo", "modprobe", "v4l2loopback",
		"devices=2",
		fmt.Sprintf("video_nr=%s", videoNr),
		"card_label=go4vl_test_1,go4vl_test_2",
		"exclusive_caps=0",
		"max_buffers=4",
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to load v4l2loopback: %v\n%s", err, output)
	}

	loopbackLoaded = true
	time.Sleep(1 * time.Second)

	// Set device paths
	testDevice1 = fmt.Sprintf("/dev/video%d", deviceNums[0])
	testDevice2 = fmt.Sprintf("/dev/video%d", deviceNums[1])
	testDevices = []string{testDevice1, testDevice2}

	// Verify devices were created
	for _, dev := range testDevices {
		if _, err := os.Stat(dev); err != nil {
			return fmt.Errorf("device %s not created after module load", dev)
		}
	}

	if *verbose {
		log.Printf("Created loopback devices: %v", testDevices)
	}

	return nil
}

// teardownLoopback unloads the v4l2loopback module
func teardownLoopback() {
	if !loopbackLoaded || *keepRunning {
		return
	}

	if *verbose {
		log.Println("Unloading v4l2loopback...")
	}

	time.Sleep(500 * time.Millisecond)
	if err := exec.Command("sudo", "modprobe", "-r", "v4l2loopback").Run(); err != nil {
		log.Printf("Warning: Failed to unload v4l2loopback: %v", err)
	} else if *verbose {
		log.Println("v4l2loopback unloaded")
	}
}

// setDevices sets testDevice1 and testDevice2 from a device list
func setDevices(devices []string) {
	testDevices = devices
	if len(devices) >= 1 {
		testDevice1 = devices[0]
	}
	if len(devices) >= 2 {
		testDevice2 = devices[1]
	}
}

// TestMain sets up and tears down the testing environment
func TestMain(m *testing.M) {
	flag.Parse()

	switch {
	case *useDevice != "":
		// Mode: real device — no ffmpeg, no module loading
		if *useDevice == "auto" {
			devices := findExistingV4L2Devices()
			if len(devices) == 0 {
				log.Println("No V4L2 devices found")
			} else {
				setDevices(devices)
				log.Printf("Discovered real devices: %v", testDevices)
			}
		} else {
			setDevices([]string{*useDevice})
			log.Printf("Using specified device: %s", testDevice1)
		}

	case *useDeviceEmulation != "":
		// Mode: v4l2loopback emulation
		if *useDeviceEmulation == "auto" {
			devices := findLoopbackDevices()
			if len(devices) >= 2 {
				setDevices(devices)
				log.Printf("Found existing loopback devices: %v", testDevices)
			} else {
				if err := setupLoopback(); err != nil {
					log.Printf("Failed to setup v4l2loopback: %v", err)
					log.Println("Install v4l2loopback-dkms or load the module manually")
				}
			}
		}
		defer teardownLoopback()

	default:
		// Mode: auto-detect — try loopback/vivid first, then real devices
		devices := findLoopbackDevices()
		if len(devices) >= 1 {
			setDevices(devices)
			log.Printf("Auto-detected loopback devices: %v", testDevices)
		} else {
			// Try loading v4l2loopback
			if err := setupLoopback(); err == nil && testDevice1 != "" {
				log.Printf("Loaded v4l2loopback, using devices: %v", testDevices)
			} else {
				// Fall back to real devices
				devices = findExistingV4L2Devices()
				if len(devices) >= 1 {
					setDevices(devices)
					log.Printf("Using real devices: %v", testDevices)
				} else {
					log.Println("No V4L2 devices found, hardware-dependent tests will skip")
				}
			}
		}
		defer teardownLoopback()
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
