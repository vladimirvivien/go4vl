package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vladimirvivien/go4vl/v4l2"
)

func deviceCap(device *v4l2.Device) error {
	caps, err := device.GetCapability()
	if err != nil {
		return err
	}

	log.Printf("%#v", caps.String())
	return nil
}

func setDefaultCrop(device *v4l2.Device) error {
	cap, err := device.GetCropCapability()
	if err != nil {
		return err
	}
	log.Printf("device crop capability: %s", cap.String())
	err = device.SetCropRect(cap.DefaultRect)
	if err != nil {
		log.Printf("setcrop unsupported: %s", err)
	}
	return nil
}

func getPixelFormat(device *v4l2.Device) error {
	format, err := device.GetPixFormat()
	if err != nil {
		return fmt.Errorf("default format: %w", err)
	}
	log.Println("got default format")
	log.Printf("pixformat %#v", format)
	return nil
}

func setPixelFormat(device *v4l2.Device) error {
	err := device.SetPixFormat(v4l2.PixFormat{
		Width:       320,
		Height:      240,
		PixelFormat: v4l2.PixelFmtYUYV,
		Field:       v4l2.FieldNone,
	})
	if err != nil {
		return fmt.Errorf("failed to set format: %w", err)
	}
	log.Println("pixel format set")
	return nil
}

func main() {
	var devName string
	flag.StringVar(&devName, "d", "/dev/video0", "device name (path)")
	flag.Parse()
	device, err := v4l2.Open(devName)
	if err != nil {
		log.Fatal(err)
	}
	defer device.Close()

	if err := deviceCap(device); err != nil {
		log.Fatal(err)
	}

	if err := setDefaultCrop(device); err != nil {
		log.Fatal(err)
	}

	if err := getPixelFormat(device); err != nil {
		log.Fatal(err)
	}

	if err := setPixelFormat(device); err != nil {
		log.Fatal(err)
	}
}
