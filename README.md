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
package main

import (
    ...
    "github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	// open device
	device, err := v4l2.Open("/dev/video0")
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()

	// configure device with preferred fmt
	if err := device.SetPixFormat(v4l2.PixFormat{
		Width:       640,
		Height:      480,
		PixelFormat: v4l2.PixelFmtMJPEG,
		Field:       v4l2.FieldNone,
	}); err != nil {
		log.Fatalf("failed to set format: %s", err)
	}

	// start a device stream with 3 video buffers
	if err := device.StartStream(3); err != nil {
		log.Fatalf("failed to start stream: %s", err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	// capture video data at 15 fps
	frameChan, err := device.Capture(ctx, 15)
	if err != nil {
		log.Fatal(err)
	}

	// grab 10 frames from frame channel and save them as files
	totalFrames := 10
	count := 0
	for frame := range frameChan {
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
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s: %s", fileName, err)
		}
		count++
		if count >= totalFrames {
			break
		}
	}

	cancel() // stop capture
	if err := device.StopStream(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
}
```

### Other examples
The [./examples](./examples) directory contains additional examples including:

* [device_info](./examples/device_info) - queries and prints devince information
* [webcam](./examples/webcam) - uses the v4l2 package to create a simple webcam that streams images from an attached camera accessible via a web page.

## Roadmap
There is no defined roadmap. The main goal is to port as much functionlities as possible so that 
adopters can use Go to create cool video-based tools on platforms such as the Raspberry Pi.