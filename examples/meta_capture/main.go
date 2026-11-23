package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func main() {
	devPath := flag.String("d", "auto", "metadata V4L2 device path (or 'auto' to discover)")
	flag.Parse()

	path := *devPath
	var bufType v4l2.BufType
	if path == "auto" {
		var ok bool
		path, bufType, ok = discoverMetaDevice()
		if !ok {
			log.Fatal("no V4L2 metadata capture or output device found")
		}
	}

	fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK|sys.O_CLOEXEC, 0)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	defer v4l2.CloseDevice(fd)

	cap, err := v4l2.GetCapability(fd)
	if err != nil {
		log.Fatalf("GetCapability(%s): %v", path, err)
	}
	if bufType == 0 {
		caps := cap.GetCapabilities()
		switch {
		case caps&v4l2.CapMetadataCapture != 0:
			bufType = v4l2.BufTypeMetaCapture
		case caps&v4l2.CapMetadataOutput != 0:
			bufType = v4l2.BufTypeMetaOutput
		default:
			log.Fatalf("%s (%s) does not support metadata capture or output", path, cap.Card)
		}
	}

	fmt.Printf("Device: %s (%s, driver=%s, bufType=%d)\n", path, cap.Card, cap.Driver, bufType)

	if descs, err := v4l2.GetAllFormatDescriptionsByType(fd, bufType); err == nil {
		fmt.Printf("Supported metadata formats (%d):\n", len(descs))
		for _, d := range descs {
			fmt.Printf("  [%d] %s (0x%08x)\n", d.Index, d.Description, d.PixelFormat)
		}
	}

	cur, err := v4l2.GetMetaFormat(fd, bufType)
	if err != nil {
		log.Fatalf("GetMetaFormat(%s): %v", path, err)
	}
	fmt.Printf("Active format: %s\n", cur)

	tried, err := v4l2.TryMetaFormat(fd, bufType, cur)
	if err != nil {
		log.Fatalf("TryMetaFormat(%s): %v", path, err)
	}
	fmt.Printf("Try format:    %s\n", tried)
}

func discoverMetaDevice() (string, v4l2.BufType, bool) {
	for i := 0; i < 32; i++ {
		path := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK|sys.O_CLOEXEC, 0)
		if err != nil {
			continue
		}
		cap, err := v4l2.GetCapability(fd)
		_ = v4l2.CloseDevice(fd)
		if err != nil {
			continue
		}
		caps := cap.GetCapabilities()
		if caps&v4l2.CapMetadataCapture != 0 {
			return path, v4l2.BufTypeMetaCapture, true
		}
		if caps&v4l2.CapMetadataOutput != 0 {
			return path, v4l2.BufTypeMetaOutput, true
		}
	}
	return "", 0, false
}
