[![Go Reference](https://pkg.go.dev/badge/github.com/vladimirvivien/go4vl.svg)](https://pkg.go.dev/github.com/vladimirvivien/go4vl) [![Go Report Card](https://goreportcard.com/badge/github.com/vladimirvivien/go4vl)](https://goreportcard.com/report/github.com/vladimirvivien/go4vl) [![Build Status](https://github.com/vladimirvivien/go4vl/actions/workflows/test.yml/badge.svg)](https://github.com/vladimirvivien/go4vl/actions/workflows/test.yml) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# go4vl

![](./docs/go4vl-logo-small.png)

A Go centric abstraction of the library for  `Video for Linux 2`  (v4l2) user API.

----

The `go4vl` project is for working with the Video for Linux 2 API for real-time video.
It hides all the complexities of working with V4L2 and provides idiomatic Go types, like channels, to consume and process captured video frames.

> This project is designed to work with Linux and the Linux Video API only.  It is *NOT* meant to be a portable/cross-platform package.

## Table of Contents

- [Why go4vl?](#why-go4vl)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Overview](#api-overview)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Roadmap](#roadmap)

## Why go4vl?

Working directly with V4L2 in Go requires complex C interop and manual memory management. `go4vl` provides:

* **Idiomatic Go API** - Use channels and standard Go types instead of C structs
* **Zero-copy streaming** - Memory-mapped buffers for efficient video processing
* **Simplified device control** - Easy access to formats, controls, and capture settings
* **Pure Linux focus** - Optimized for Linux video pipelines without cross-platform compromises

## Features

* Capture and control video data from your Go programs
* Idiomatic Go types such as channels to access and stream video data
* Exposes device enumeration and information
* Provides device capture control
* Access to video format information
* Streaming uses zero-copy IO with memory-mapped buffers

## Prerequisites

**Software:**
* Go 1.16 or later
* Linux kernel 5.10.x or later (go4vl works with kernel 5.10.x and newer)
* C compiler (gcc/clang) or cross-compiler for target platform
* V4L2 kernel headers (linux-libc-dev on Debian/Ubuntu, kernel-headers on RHEL/Fedora)
* V4L2 drivers (typically included in Linux kernel)

**Hardware:**
* V4L2-compatible capture device (webcam, capture card, etc.)
* Device accessible via `/dev/videoX`

**Tested platforms:**
* Raspberry Pi 3/4 (32-bit and 64-bit Raspberry Pi OS)
* x86_64 Linux distributions
* ARM Linux systems with V4L2 support

See [docs/BUILD.md](./docs/BUILD.md) for comprehensive build instructions including prerequisite installation, cross-compilation, and troubleshooting.

## Installation

```bash
go get github.com/vladimirvivien/go4vl/v4l2
```

Ensure your user has access to video devices:
```bash
sudo usermod -a -G video $USER
# Log out and back in for changes to take effect
```

**Note**: For detailed build instructions, including installing V4L2 headers, cross-compilation, and using custom kernel headers, see [docs/BUILD.md](./docs/BUILD.md).

## Building

To build the go4vl packages or examples:

```bash
# Build the v4l2 package
go build ./v4l2

# Build an example
go build ./examples/snapshot
```

The build requires V4L2 kernel headers (typically from `linux-libc-dev` package). To use custom kernel headers, override the include path:

```bash
CGO_CFLAGS="-I/path/to/custom/headers" go build ./v4l2
```

For detailed build instructions including prerequisite installation, cross-compilation with Zig or Docker, and troubleshooting, see [docs/BUILD.md](./docs/BUILD.md).

## Quick Start

Capture a single frame and save to file (assumes MJPEG format support):

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/vladimirvivien/go4vl/device"
)

func main() {
	dev, err := device.Open("/dev/video0", device.WithBufferSize(1))
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	if err := dev.Start(context.TODO()); err != nil {
		log.Fatal(err)
	}

	// Capture frame from channel
	frame := <-dev.GetOutput()

	file, err := os.Create("pic.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.Write(frame); err != nil {
		log.Fatal(err)
	}
}
```

See complete example: [examples/snapshot/snap.go](./examples/snapshot/snap.go)

## API Overview

**Core packages:**
* **`device`** - High-level device operations (open, start, capture, close)
* **`v4l2`** - Low-level V4L2 types and ioctls
* **`imgsupport`** - Image format conversion utilities

**Key concepts:**

**Device management:**
```go
// List available devices
devices := device.GetAllDevices()

// Open device with options
dev, err := device.Open("/dev/video0",
    device.WithBufferSize(4),
    device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG}),
)
defer dev.Close()
```

**Format control:**
```go
// Query supported formats
formats := dev.GetFormatDescriptions()

// Set pixel format
currFmt, err := dev.GetPixFormat()
currFmt.PixelFormat = v4l2.PixelFmtYUYV
if err := dev.SetPixFormat(currFmt); err != nil {
    log.Fatal(err)
}
```

**Device controls:**
```go
// Get control value
brightness := dev.GetControl(v4l2.CtrlBrightness)

// Set control value
dev.SetControl(v4l2.CtrlBrightness, 128)
```

**Streaming (Legacy API):**
```go
ctx := context.Background()
dev.Start(ctx)

for frame := range dev.GetOutput() {
    // Process frame bytes
    processFrame(frame)
}
```

**Streaming (Frame API - Recommended):**
```go
ctx := context.Background()
dev.Start(ctx)

for frame := range dev.GetFrames() {
    // Access frame data with metadata
    processFrame(frame.Data)

    // Access metadata
    log.Printf("Frame %d captured at %v", frame.Sequence, frame.Timestamp)

    // Check frame type (for compressed formats like MJPEG, H.264)
    if frame.IsKeyFrame() {
        log.Printf("Keyframe detected")
    }

    // IMPORTANT: Release buffer back to pool
    frame.Release()
}
```


The `GetFrames()` API uses buffer pooling to dramatically reduce allocation overhead and GC pressure, making it ideal for high-throughput video processing. See [examples/capture_frames](./examples/capture_frames/) for a complete example.

Full API documentation: [pkg.go.dev/github.com/vladimirvivien/go4vl](https://pkg.go.dev/github.com/vladimirvivien/go4vl)

## Examples

This repository includes multiple examples demonstrating various capabilities:

* **[snapshot](./examples/snapshot/)** - Capture single frame to file
* **[capture0](./examples/capture0/)** - Capture multiple frames (legacy API)
* **[capture_frames](./examples/capture_frames/)** - Capture with metadata and pooling (recommended)
* **[capture1](./examples/capture1/)** - Capture with specific format
* **[device_info](./examples/device_info/)** - Query device information
* **[format](./examples/format/)** - Query and set formats
* **[user_ctrl](./examples/user_ctrl/)** - Control brightness, contrast, etc.
* **[ext_ctrls](./examples/ext_ctrls/)** - Extended codec controls
* **[simplecam](./examples/simplecam/)** - Web streaming camera
* **[webcam](./examples/webcam/)** - Full-featured webcam with controls

Full examples list: [examples/README.md](./examples/README.md)

## Troubleshooting

**Device not found:**
```bash
# List available video devices
ls -l /dev/video*

# Check device info
v4l2-ctl --list-devices
```

**Permission denied:**
```bash
# Add user to video group
sudo usermod -a -G video $USER
# Log out and back in

# Or temporarily
sudo chmod 666 /dev/video0
```

**No frames received:**
```bash
# Verify device capabilities
v4l2-ctl -d /dev/video0 --all

# Check supported formats
v4l2-ctl -d /dev/video0 --list-formats-ext
```

**Build errors:**
```bash
# Install build essentials
sudo apt install build-essential

# Install kernel headers
sudo apt install linux-headers-$(uname -r)
```

## Testing

go4vl uses real V4L2 devices for testing with v4l2loopback virtual devices.

```bash
# Run unit tests
go test -v ./device ./v4l2 ./imgsupport

# Run integration tests (requires v4l2loopback)
sudo go test -v -tags=integration ./test/...
```

See [TESTING_GUIDE.md](./TESTING_GUIDE.md) for comprehensive testing documentation including:
* v4l2loopback setup
* Virtual device configuration
* CI/CD integration
* Docker-based testing

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure:
* Code follows Go conventions (`go fmt`, `go vet`)
* Tests pass (`go test ./...`)
* New features include tests
* Documentation is updated

Report bugs and request features via [GitHub Issues](https://github.com/vladimirvivien/go4vl/issues).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2021 Vladimir Vivien

## Roadmap

The main goal is to port as many V4L2 functionalities as possible so that adopters can use Go to create video-based tools on platforms such as the Raspberry Pi and other Linux systems.

**Current focus:**
* Extended codec controls (H.264, VP8, MPEG2)
* Advanced streaming modes
* Performance optimizations
* Broader device compatibility testing
