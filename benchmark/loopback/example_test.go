package loopback_test

import (
	"fmt"
	"log"
	"time"

	"github.com/vladimirvivien/go4vl/benchmark/loopback"
)

// Example demonstrates how to setup and use a v4l2loopback device
func Example() {
	// Check if loopback is available
	if !loopback.IsAvailable() {
		log.Println("Loopback not available: install ffmpeg and v4l2loopback-dkms")
		return
	}

	// Setup loopback device: /dev/video50 at 640x480@30fps with test pattern
	dev, err := loopback.Setup(50, 640, 480, 30, "testsrc")
	if err != nil {
		log.Fatalf("Failed to setup loopback: %v", err)
	}
	defer dev.Close()

	fmt.Printf("Loopback device ready at %s\n", dev.DevicePath)

	// Device is now ready to use with go4vl or other V4L2 applications
	time.Sleep(5 * time.Second)

	// Cleanup happens automatically via defer
	fmt.Println("Cleaning up...")
}

// Example_customPattern shows how to use different test patterns
func Example_customPattern() {
	if !loopback.IsAvailable() {
		return
	}

	// SMPTE color bars
	dev, err := loopback.Setup(50, 1280, 720, 30, "smptebars")
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	fmt.Println("SMPTE color bars streaming...")
	time.Sleep(2 * time.Second)
}

// Example_solidColor shows how to stream a solid color
func Example_solidColor() {
	if !loopback.IsAvailable() {
		return
	}

	// Solid red color
	dev, err := loopback.Setup(50, 640, 480, 30, "color=red")
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	fmt.Println("Red screen streaming...")
	time.Sleep(2 * time.Second)
}
