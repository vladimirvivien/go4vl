# Building go4vl

This guide covers how to build go4vl from source, including prerequisite installation and various build scenarios.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Building with Custom Headers](#building-with-custom-headers)
- [Cross-Compilation](#cross-compilation)
- [Platform-Specific Instructions](#platform-specific-instructions)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **Go**: Version 1.16 or later
- **Linux Kernel**: 5.10.x or later
- **C Compiler**: gcc or clang
- **V4L2 Kernel Headers**: Required for CGO compilation

### Installing Prerequisites

#### Ubuntu/Debian

```bash
# Update package list
sudo apt update

# Install build essentials (includes gcc, make, etc.)
sudo apt install build-essential

# Install V4L2 kernel headers
sudo apt install linux-libc-dev

# Optional: Install V4L2 utilities for testing
sudo apt install v4l-utils

# Verify installation
v4l2-ctl --version
```

#### Fedora/RHEL/CentOS

```bash
# Install development tools
sudo dnf groupinstall "Development Tools"

# Install V4L2 kernel headers
sudo dnf install kernel-headers kernel-devel

# Optional: Install V4L2 utilities
sudo dnf install v4l-utils

# Verify installation
v4l2-ctl --version
```

#### Arch Linux

```bash
# Install base development tools
sudo pacman -S base-devel

# Install kernel headers (usually already present)
sudo pacman -S linux-headers

# Optional: Install V4L2 utilities
sudo pacman -S v4l-utils

# Verify installation
v4l2-ctl --version
```

#### Alpine Linux

```bash
# Install build tools
apk add build-base

# Install kernel headers
apk add linux-headers

# Optional: Install V4L2 utilities
apk add v4l-utils
```

### User Permissions

Ensure your user has access to video devices:

```bash
# Add user to video group
sudo usermod -a -G video $USER

# Log out and back in for changes to take effect
# Or verify group membership
groups | grep video
```

## Quick Start

### Basic Build

```bash
# Clone the repository
git clone https://github.com/vladimirvivien/go4vl.git
cd go4vl

# Build the v4l2 package
go build ./v4l2

# Or build a specific example
go build ./examples/snapshot
```

### Running Tests

```bash
# Run unit tests for all packages
go test ./device ./v4l2 ./imgsupport

# Run with verbose output
go test -v ./v4l2

# Run specific test
go test -v ./v4l2 -run TestCapability
```

### Installing as a Dependency

```bash
# Add to your project
go get github.com/vladimirvivien/go4vl/v4l2
```

## Building with Custom Headers

go4vl uses system V4L2 headers by default from `/usr/include`. To use custom or newer kernel headers, override the include path using the `CGO_CFLAGS` environment variable.

### Using Latest Kernel Headers

If you want to use headers from a specific kernel version:

```bash
# Download and extract kernel source
wget https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-6.8.tar.xz
tar xf linux-6.8.tar.xz

# Build with custom headers
CGO_CFLAGS="-I$PWD/linux-6.8/include" go build ./v4l2
```

### Using Headers from Kernel Source Tree

```bash
# If you have kernel sources installed
CGO_CFLAGS="-I/usr/src/linux/include" go build ./v4l2
```

### Making it Permanent

To avoid specifying `CGO_CFLAGS` every time:

```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export CGO_CFLAGS="-I/path/to/custom/headers"

# Reload your shell or source the file
source ~/.bashrc

# Now build normally
go build ./v4l2
```

## Cross-Compilation

### Prerequisites for Cross-Compilation

Cross-compilation requires:
1. Target architecture's C cross-compiler
2. Target architecture's V4L2 headers (from target's sysroot)
3. CGO enabled (`CGO_ENABLED=1`)

### Cross-Compile with GCC

#### ARM 32-bit (armv7)

```bash
# Install cross-compiler
sudo apt install gcc-arm-linux-gnueabihf

# Cross-compile
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=arm \
GOARM=7 \
CC=arm-linux-gnueabihf-gcc \
CGO_CFLAGS="-I/usr/arm-linux-gnueabihf/include" \
go build -o myapp ./examples/snapshot
```

#### ARM 64-bit (aarch64)

```bash
# Install cross-compiler
sudo apt install gcc-aarch64-linux-gnu

# Cross-compile
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=arm64 \
CC=aarch64-linux-gnu-gcc \
CGO_CFLAGS="-I/usr/aarch64-linux-gnu/include" \
go build -o myapp ./examples/snapshot
```

### Cross-Compile with Zig (Recommended)

[Zig](https://ziglang.org/) provides an easier cross-compilation experience with built-in cross-compilers.

#### Install Zig

```bash
# Download from https://ziglang.org/download/
# Or install via package manager
snap install zig --classic --beta
```

#### Build for ARM 32-bit

```bash
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=arm \
GOARM=7 \
CC="zig cc -target arm-linux-musleabihf" \
CXX="zig c++ -target arm-linux-musleabihf" \
go build -o myapp ./examples/snapshot
```

#### Build for ARM 64-bit

```bash
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=arm64 \
CC="zig cc -target aarch64-linux-musl" \
CXX="zig c++ -target aarch64-linux-musl" \
go build -o myapp ./examples/snapshot
```

#### Build for Raspberry Pi with Custom Headers

```bash
# If you have Raspberry Pi sysroot
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=arm \
GOARM=7 \
CC="zig cc -target arm-linux-musleabihf" \
CGO_CFLAGS="-I/path/to/rpi/sysroot/usr/include" \
go build -o myapp ./examples/snapshot
```

### Cross-Compile with Docker

Using Docker provides a consistent build environment with all dependencies.

#### Example Dockerfile

```dockerfile
FROM golang:1.21-bullseye

# Install cross-compilation tools
RUN apt-get update && apt-get install -y \
    gcc-arm-linux-gnueabihf \
    gcc-aarch64-linux-gnu \
    linux-libc-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
COPY . .

# Build for ARM64
RUN CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=arm64 \
    CC=aarch64-linux-gnu-gcc \
    go build -o myapp-arm64 ./examples/snapshot

CMD ["/bin/bash"]
```

#### Build with Docker

```bash
# Build the Docker image
docker build -t go4vl-builder .

# Run container and copy binary
docker run --rm -v $(pwd)/dist:/dist go4vl-builder \
    cp myapp-arm64 /dist/
```

## Platform-Specific Instructions

### Raspberry Pi (On-Device Build)

```bash
# Update system
sudo apt update && sudo apt full-upgrade

# Install prerequisites
sudo apt install golang build-essential linux-libc-dev

# Build
cd go4vl
go build ./v4l2
```

### Raspberry Pi (Cross-Compile from x86_64)

See [Cross-Compile with Zig](#cross-compile-with-zig) or [Cross-Compile with GCC](#cross-compile-with-gcc).

### WSL2 (Windows Subsystem for Linux)

go4vl requires native Linux kernel with V4L2 support. WSL2 with custom kernel supporting V4L2:

```bash
# Install prerequisites in WSL2
sudo apt update
sudo apt install build-essential linux-libc-dev golang

# Build
go build ./v4l2

# Note: You'll need a WSL2 kernel with V4L2 drivers compiled
# See: https://github.com/dorssel/usbipd-win for USB device passthrough
```

## Troubleshooting

### "linux/videodev2.h: No such file or directory"

**Problem**: V4L2 headers not installed.

**Solution**:
```bash
# Ubuntu/Debian
sudo apt install linux-libc-dev

# Fedora/RHEL
sudo dnf install kernel-headers

# Verify installation
ls /usr/include/linux/videodev2.h
```

### "cannot find package" during go get

**Problem**: Network issues or invalid module path.

**Solution**:
```bash
# Clear module cache
go clean -modcache

# Retry
go get github.com/vladimirvivien/go4vl/v4l2

# Or use explicit version
go get github.com/vladimirvivien/go4vl/v4l2@latest
```

### Cross-Compilation Headers Not Found

**Problem**: CGO can't find target architecture headers.

**Solution**:
```bash
# Install cross-architecture headers
sudo apt install linux-libc-dev:armhf  # For ARM 32-bit
sudo apt install linux-libc-dev:arm64  # For ARM 64-bit

# Or specify path explicitly
CGO_CFLAGS="-I/usr/include" go build ...
```

### "permission denied" when accessing /dev/video*

**Problem**: User not in video group.

**Solution**:
```bash
# Add user to video group
sudo usermod -a -G video $USER

# Log out and back in, or use newgrp
newgrp video

# Or temporarily change permissions (not recommended for production)
sudo chmod 666 /dev/video0
```

### CGO Compilation Very Slow

**Problem**: CGO compilation can be slower than pure Go.

**Solution**:
```bash
# Use build cache
go build -x ./v4l2  # Shows what's being rebuilt

# Parallel compilation
GOMAXPROCS=4 go build ./v4l2

# For repeated builds, cache is automatic
```

### Different Results on Different Kernel Versions

**Problem**: V4L2 API varies across kernel versions.

**Solution**:
```bash
# Check kernel version
uname -r

# Check V4L2 driver version
v4l2-ctl --all

# Use specific kernel headers if needed
CGO_CFLAGS="-I/path/to/kernel-X.Y/include" go build ./v4l2
```

## Advanced Build Options

### Static Linking

```bash
# Build static binary (useful for containers)
CGO_ENABLED=1 \
go build -ldflags="-linkmode external -extldflags -static" \
./examples/snapshot
```

### Debug Build

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" ./v4l2

# Use with delve debugger
dlv exec ./myapp
```

### Optimized Build

```bash
# Build with optimizations and trimmed binary
go build -ldflags="-s -w" ./examples/snapshot

# Further compress with upx
upx --best --lzma myapp
```

## Build Scripts

### Example Build Script

Create a `build.sh` script for consistent builds:

```bash
#!/bin/bash
set -e

# Configuration
APP_NAME="myapp"
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

echo "Building ${APP_NAME} ${VERSION}..."

# Native build
echo "Building for native platform..."
go build -ldflags="${LDFLAGS}" -o "${APP_NAME}" ./examples/snapshot

# ARM builds
echo "Building for ARM platforms..."
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 \
    CC="zig cc -target arm-linux-musleabihf" \
    go build -ldflags="${LDFLAGS}" -o "${APP_NAME}-armv7" ./examples/snapshot

CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
    CC="zig cc -target aarch64-linux-musl" \
    go build -ldflags="${LDFLAGS}" -o "${APP_NAME}-arm64" ./examples/snapshot

echo "Build complete!"
ls -lh ${APP_NAME}*
```

## Additional Resources

- [V4L2 API Documentation](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/v4l2.html)
- [Go CGO Documentation](https://pkg.go.dev/cmd/cgo)
- [Cross Compilation Guide](https://go.dev/doc/install/source#environment)
- [go4vl Examples](../examples/README.md)
- [Testing Guide](../TESTING_GUIDE.md)

## Getting Help

- Report issues: [GitHub Issues](https://github.com/vladimirvivien/go4vl/issues)
- API Documentation: [pkg.go.dev](https://pkg.go.dev/github.com/vladimirvivien/go4vl)
- Examples: [examples/](../examples/)
