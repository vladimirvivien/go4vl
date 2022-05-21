[![Go Report Card](https://goreportcard.com/badge/github.com/vladimirvivien/go4vl)](https://goreportcard.com/report/github.com/vladimirvivien/go4vl)

# go4vl

A Go library for the `Video for Linux 2`  (v4l2) user API.

----

The `go4vl` project is for working with the Video for Linux 2 API for real-time video. 
It hides all the complexities of working with V4L2 and provides idiomatic Go types, like channels, to consume and process captured video frames.

> This project is designed to work with Linux and the Linux Video API.  
> It is *NOT* meant to be a portable/cross-platform capable package for real-time video processing.

## Features

* Capture and control video data from your Go programs
* Idiomatic Go types such as channels to access and stream video data
* Exposes device enumeration and information
* Provides device capture control
* Access to video format information
* Streaming users zero-copy IO using memory mapped buffers

## Compilation Requirements

* Go compiler/tools
* Kernel minimum v5.10.x
* A locally configured C compiler (i.e. gcc)
* Header files for V4L2 (i.e. /usr/include/linux/videodev2.h)

All examples have been tested using a Raspberry PI 3, running 32-bit Raspberry PI OS.
The package should work with no problem on your 64-bit Linux OS.

## Getting started

### System upgrade

To avoid issues with old header files on your machine, upgrade your system to pull down the latest OS packages
with something similar to the following (follow directions for your system for proper upgrade):

```shell
sudo apt update
sudo apt full-upgrade
```

### Using the go4vl package

To include `go4vl` in your own code, `go get` the package:

```bash
go get github.com/vladimirvivien/go4vl/v4l2
```

## Video capture example

The following is a simple example that captures video data from an attached camera device to
and saves the captured frames as JPEG files. 

The example assumes the attached device supports JPEG (MJPEG) output format inherently.

```go
func main() {
	devName := "/dev/video0"
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.Parse()

	// open device
	device, err := device.Open(
		devName,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMPEG, Width: 640, Height: 480}),
	)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()

	// start stream with cancellable context
	ctx, stop := context.WithCancel(context.TODO())
	if err := device.Start(ctx); err != nil {
		log.Fatalf("failed to start stream: %s", err)
	}

	// process frames from capture channel
	totalFrames := 10
	count := 0
	log.Printf("Capturing %d frames...", totalFrames)

	for frame := range device.GetOutput() {
		fileName := fmt.Sprintf("capture_%d.jpg", count)
		file, err := os.Create(fileName)
		if err != nil {
			log.Printf("failed to create file %s: %s", fileName, err)
			continue
		}
		if _, err := file.Write(frame); err != nil {
			log.Printf("failed to write file %s: %s", fileName, err)
			continue
		}
		log.Printf("Saved file: %s", fileName)
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s: %s", fileName, err)
		}
		count++
		if count >= totalFrames {
			break
		}
	}

	stop() // stop capture
	fmt.Println("Done.")
}
```

> Read a detail walk-through about this example [here](./examples/capture0/README.md).

### Other examples
The [./examples](./examples/README.md) directory contains additional examples including:
* [device_info](./examples/device_info/README.md) - queries and prints video device information
* [webcam](./examples/webcam/README.md) - uses the v4l2 package to create a simple webcam that streams images from an attached camera accessible via a web page.

## Roadmap
The main goal is to port as many functionalities as possible so that 
adopters can use Go to create cool video-based tools on platforms such as the Raspberry Pi.