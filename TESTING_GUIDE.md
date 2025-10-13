# Testing Guide for go4vl

This comprehensive guide covers all aspects of testing the go4vl library, including unit tests, integration tests, and CI/CD setup.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Testing Approaches](#testing-approaches)
- [Unit Tests](#unit-tests)
- [Integration Tests](#integration-tests)
- [Test Flags](#test-flags)
- [Environment Variables](#environment-variables)
- [CI/CD Setup](#cicd-setup)
- [Troubleshooting](#troubleshooting)
- [References](#references)

## Overview

The go4vl test suite uses real V4L2 devices (no mock implementations) to ensure authentic testing:

- **Unit Tests** - Test packages in isolation without requiring V4L2 devices
- **Integration Tests** - Test with v4l2loopback virtual devices or real hardware

### Key Features

- **Automatic v4l2loopback setup** - TestMain handles module loading/unloading
- **Dynamic device selection** - Avoids conflicts with existing devices (starts from /dev/video40)
- **Fallback to existing devices** - Uses available devices when v4l2loopback setup fails
- **No environment variables required** - Uses build tags for test selection
- **No build tools needed** - All commands use standard Go tooling
- **Graceful degradation** - Tests skip appropriately when devices or permissions aren't available

## Prerequisites

### 1. Go Installation

Go 1.21 or later is required.

```bash
# Check Go version
go version
```

### 2. User Permissions

For integration tests, you need either:
- Root access to load kernel modules (recommended for automatic setup)
- Video group membership for existing devices

```bash
# Add user to video group
sudo usermod -a -G video $USER

# Logout and login for changes to take effect
```

### 3. Install Testing Tools

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y v4l-utils v4l2loopback-dkms ffmpeg

# Fedora
sudo dnf install -y v4l-utils v4l2loopback ffmpeg

# Arch Linux
sudo pacman -S v4l-utils v4l2loopback-dkms ffmpeg
```

#### Package Descriptions

- **v4l-utils** - Command-line tools for V4L2 devices (v4l2-ctl, etc.)
- **v4l2loopback-dkms** - Kernel module for virtual video devices
- **ffmpeg** - Video processing tool for feeding test patterns

## Testing Approaches

### 1. v4l2loopback - Virtual Video Device (Recommended)

**v4l2loopback** is a kernel module that creates virtual V4L2 loopback devices. It's the most realistic way to test V4L2 code without hardware.

#### Installation

```bash
# Ubuntu/Debian
sudo apt-get install v4l2loopback-dkms v4l2loopback-utils

# Fedora
sudo dnf install v4l2loopback

# From source
git clone https://github.com/umlaeute/v4l2loopback
cd v4l2loopback
make && sudo make install
```

#### Manual Setup

```bash
# Load the module with options (for manual testing)
sudo modprobe v4l2loopback devices=2 video_nr=42,43 card_label="Test Camera 1","Test Camera 2" exclusive_caps=1

# Verify devices created
ls -la /dev/video42 /dev/video43
v4l2-ctl --list-devices

# Remove module when done
sudo modprobe -r v4l2loopback
```

#### Feeding Test Data

```bash
# Using ffmpeg to feed test pattern
ffmpeg -re -f lavfi -i testsrc=size=640x480:rate=30 -pix_fmt yuyv422 -f v4l2 /dev/video42

# Using gstreamer
gst-launch-1.0 videotestsrc ! v4l2sink device=/dev/video42

# Feed an MP4 file
ffmpeg -re -i test.mp4 -pix_fmt yuyv422 -f v4l2 /dev/video42

# Feed static image
ffmpeg -loop 1 -i test.jpg -pix_fmt yuyv422 -f v4l2 /dev/video42

# Feed SMPTE color bars
ffmpeg -f lavfi -i smptebars=size=1280x720:rate=25 -pix_fmt yuyv422 -f v4l2 /dev/video43
```

### 2. vivid - Virtual Video Test Driver

The **vivid** (Virtual Video Test Driver) is a kernel module that creates virtual video capture and output devices with extensive testing capabilities.

#### Setup

```bash
# Load vivid module
sudo modprobe vivid

# Configure vivid with specific options
sudo modprobe vivid n_devs=2 \
    vid_cap_nr=10,11 \
    vid_out_nr=12,13 \
    node_types=0x1,0x1

# Check created devices
v4l2-ctl -d /dev/video10 --info
```

#### Features

- Multiple test patterns
- Various resolutions and formats
- Control simulation
- Error injection capabilities

### 3. Real Hardware

If you have a USB webcam or other V4L2 device:

```bash
# List all video devices
v4l2-ctl --list-devices

# Test with specific device
V4L2_TEST_DEVICE=/dev/video0 go test -v -tags=integration ./test/...
```

## Unit Tests

Unit tests validate the core functionality without requiring V4L2 devices. They test data structures, methods, and business logic in isolation.

### Packages with Unit Tests

- **v4l2** - Core V4L2 API types and operations
- **device** - High-level device management
- **imgsupport** - Image format support utilities

### Running Unit Tests

```bash
# Run all unit tests
go test -v ./device ./v4l2 ./imgsupport

# Run specific package tests
go test -v ./v4l2
go test -v ./device

# Run with coverage
go test -coverprofile=coverage.out ./device ./v4l2 ./imgsupport
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test -v ./v4l2 -run TestCapability_GetCapabilities

# Clean test cache
go clean -testcache
```

### Unit Test Coverage

The unit test suite includes comprehensive coverage:

#### v4l2 Package Tests

1. **capability_test.go** (662 lines) - 15 test functions with 53+ subtests
   - Capability detection methods
   - Version parsing
   - String formatting
   - 29 capability constants validation

2. **syscalls_test.go** (230 lines) - 7 test functions
   - WaitForRead context handling
   - Context cancellation
   - Goroutine cleanup validation

3. **capability_bench_test.go** (101 lines) - 7 benchmark functions
   - Performance baselines (~0.26 ns/op with 0 allocations)

4. **format_test.go** (528 lines) - 24 test functions with 95+ subtests
   - Pixel formats and colorspaces
   - YCbCr encodings
   - Quantization and transfer functions

5. **streaming_test.go** (489 lines) - 18 test functions with 75+ subtests
   - Buffer types and I/O types
   - Buffer lifecycle management
   - Keyframe detection

6. **control_test.go** (489 lines) - 16 test functions
   - Control classes and types
   - Control value ranges
   - Menu items

7. **crop_test.go** (585 lines) - 17 test functions
   - Cropping operations
   - Aspect ratios
   - Digital zoom scenarios

#### device Package Tests

1. **device_test.go** (742 lines) - 31 test functions
   - Device struct methods
   - Configuration options
   - Error handling
   - Concurrent streaming flag operations

## Integration Tests

Integration tests validate the complete streaming pipeline with real or virtual V4L2 devices. See `test/README.md` for detailed information about the integration test suite.

### Running Integration Tests

```bash
# Run with automatic setup (requires sudo for module loading)
sudo go test -v -tags=integration ./test/...

# Run without sudo (uses existing devices or skips)
go test -v -tags=integration ./test/...

# Run specific test files
go test -v -tags=integration ./test/... -run TestDevice
go test -v -tags=integration ./test/... -run TestV4L2

# Generate coverage report
go test -tags=integration -coverprofile=coverage-integration.out ./test/...
go tool cover -html=coverage-integration.out -o coverage-integration.html
```

### Integration Test Categories

The integration test suite includes:

- **Device package tests** - All exported device functionality
- **V4L2 package tests** - All exported v4l2 types and constants
- **Full pipeline tests** - Complete streaming workflows
- **Simple tests** - Basic tests that work with any available device

### Manual Device Setup (Alternative)

If you prefer manual setup or automatic setup fails:

```bash
# 1. Load v4l2loopback module with test devices
sudo modprobe v4l2loopback devices=2 video_nr=42,43 exclusive_caps=1

# 2. Start test patterns
ffmpeg -f lavfi -i testsrc=size=640x480:rate=30 -pix_fmt yuyv422 -f v4l2 /dev/video42 &
ffmpeg -f lavfi -i smptebars=size=1280x720:rate=25 -pix_fmt yuyv422 -f v4l2 /dev/video43 &

# 3. Run tests with skip-setup flag
go test -v -tags=integration ./test/... -skip-setup

# 4. Clean up when done
killall ffmpeg
sudo modprobe -r v4l2loopback
```

## Test Flags

The integration test suite supports several command-line flags:

```bash
# Skip automatic v4l2loopback setup (use existing devices)
go test -v -tags=integration ./test/... -skip-setup

# Keep v4l2loopback loaded after tests complete
go test -v -tags=integration ./test/... -keep-running

# Use existing v4l2loopback if already loaded
go test -v -tags=integration ./test/... -use-existing

# Enable verbose logging
go test -v -tags=integration ./test/... -verbose
```

## Environment Variables

The test suite supports optional environment variables:

### V4L2_TEST_DEVICE

Force a specific test device to use (rarely needed):

```bash
V4L2_TEST_DEVICE=/dev/video0 go test -v -tags=integration ./test/...
```

### CI

Enable CI mode for automatic v4l2loopback setup in CI environments:

```bash
export CI=true
sudo go test -v -tags=integration ./test/...
```

### Notes

- No `RUN_INTEGRATION` variable is needed - tests are controlled by the `-tags=integration` build tag
- The test suite automatically selects device numbers starting from `/dev/video40` to avoid conflicts

## CI/CD Setup

### GitHub Actions Workflow

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Run Unit Tests
      run: go test -v ./device ./v4l2 ./imgsupport

    - name: Generate Coverage Report
      run: |
        go test -coverprofile=coverage.out ./device ./v4l2 ./imgsupport
        go tool cover -html=coverage.out -o coverage.html

    - name: Upload Coverage
      uses: actions/upload-artifact@v3
      with:
        name: coverage
        path: coverage.html

  integration-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Dependencies
      run: |
        sudo apt-get update
        sudo apt-get install -y v4l2loopback-dkms v4l-utils ffmpeg

    - name: Run Integration Tests
      run: |
        export CI=true
        # TestMain will handle v4l2loopback setup automatically
        sudo go test -v -tags=integration ./test/...
```

### Docker-based Testing

#### Dockerfile

```dockerfile
# Dockerfile.test
FROM golang:1.21

# Install V4L2 tools and dependencies
RUN apt-get update && apt-get install -y \
    v4l-utils \
    v4l2loopback-dkms \
    ffmpeg \
    build-essential \
    kmod \
    && rm -rf /var/lib/apt/lists/*

# Set up working directory
WORKDIR /app

# Copy source code
COPY . .

# Test script
COPY <<'EOF' /test.sh
#!/bin/bash
set -e

# Load v4l2loopback if possible (requires privileged mode)
if [ -w /dev ]; then
    modprobe v4l2loopback devices=2 video_nr=42,43 exclusive_caps=1 || true
fi

# Run unit tests
echo "Running unit tests..."
go test -v ./device ./v4l2 ./imgsupport

# Run integration tests
echo "Running integration tests..."
go test -v -tags=integration ./test/...
EOF

RUN chmod +x /test.sh

CMD ["/test.sh"]
```

#### Running Docker Tests

```bash
# Build test image
docker build -f Dockerfile.test -t go4vl-test .

# Run tests (privileged needed for kernel modules)
docker run --privileged go4vl-test

# Run with volume mount for development
docker run --privileged -v $(pwd):/app go4vl-test

# Run with specific device
docker run --privileged --device=/dev/video0 go4vl-test
```

## Troubleshooting

### Tests Skip with "Test device not available"

The tests automatically skip when devices aren't available. To fix:

1. **Check v4l2loopback is installed**:
   ```bash
   sudo apt-get install v4l2loopback-dkms
   ```

2. **Check if module loads manually**:
   ```bash
   sudo modprobe v4l2loopback devices=1 video_nr=42
   ls -la /dev/video42
   sudo modprobe -r v4l2loopback
   ```

3. **Run tests with sudo** (for automatic setup):
   ```bash
   sudo go test -v -tags=integration ./test/...
   ```

### v4l2loopback not loading

```bash
# Check kernel module support
lsmod | grep v4l2loopback

# Check dmesg for errors
sudo dmesg | grep v4l2loopback

# Ensure headers installed
sudo apt-get install linux-headers-$(uname -r)
```

### v4l2loopback Module Not Loading

```
modprobe: FATAL: Module v4l2loopback not found
```

**Solutions**:
```bash
# Install module
sudo apt-get install v4l2loopback-dkms

# If still failing, install kernel headers
sudo apt-get install linux-headers-$(uname -r)

# Rebuild module
sudo dpkg-reconfigure v4l2loopback-dkms
```

### Permission Denied

```
Failed to open /dev/video0: permission denied
```

**Solution**: Add user to video group:
```bash
sudo usermod -a -G video $USER
# Logout and login again
```

### v4l2-ctl Not Found

```
v4l2-ctl not found. Install with: apt-get install v4l-utils
```

**Solution**: Install v4l-utils package:
```bash
sudo apt-get install v4l-utils
```

### Device Busy

```
Failed to start streaming: device or resource busy
```

**Solutions**:
1. Check if another application is using the device:
   ```bash
   fuser /dev/video42
   lsof /dev/video42
   ```
2. Use different device numbers:
   ```bash
   sudo modprobe v4l2loopback video_nr=50,51
   ```

### No frames received

```bash
# Check if producer is running
ps aux | grep ffmpeg

# Verify device is readable
v4l2-ctl -d /dev/video42 --all

# Check device capabilities
v4l2-ctl -d /dev/video42 --all

# Try reading directly
ffmpeg -f v4l2 -i /dev/video42 -frames:v 1 test.jpg
```

### Test Cache Issues

```bash
# Clean test cache
go clean -testcache

# Force test re-run
go test -count=1 -v ./device ./v4l2 ./imgsupport
```

## Best Practices

1. **Real V4L2 testing only** - No mock devices, uses v4l2loopback or real hardware
2. **Automatic setup/teardown** - TestMain handles v4l2loopback lifecycle
3. **Dynamic device selection** - Avoids conflicts with existing devices
4. **Build tags for test separation** - Use `-tags=integration` for integration tests
5. **Graceful degradation** - Tests skip when devices unavailable
6. **Run unit tests frequently** - They're fast and don't require special setup
7. **Run integration tests before commits** - Validates complete functionality
8. **Use coverage reports** - Identify untested code paths
9. **Test with real hardware** - When possible, test with actual webcams

## Quick Reference

```bash
# Unit Tests
go test -v ./device ./v4l2 ./imgsupport

# Integration Tests (automatic setup)
sudo go test -v -tags=integration ./test/...

# Integration Tests (manual setup)
sudo modprobe v4l2loopback devices=2 video_nr=42,43 exclusive_caps=1
ffmpeg -f lavfi -i testsrc=size=640x480:rate=30 -pix_fmt yuyv422 -f v4l2 /dev/video42 &
go test -v -tags=integration ./test/... -skip-setup
killall ffmpeg
sudo modprobe -r v4l2loopback

# Coverage Reports
go test -coverprofile=coverage.out ./device ./v4l2 ./imgsupport
go test -tags=integration -coverprofile=coverage-integration.out ./test/...
go tool cover -html=coverage.out -o coverage.html

# Benchmarks
go test -bench=. -benchmem ./v4l2
go test -bench=. -benchmem -tags=integration ./test/...

# Specific Tests
go test -v ./v4l2 -run TestCapability
go test -v -tags=integration ./test/... -run TestDevice

# Clean Cache
go clean -testcache
```

## References

- [v4l2loopback GitHub](https://github.com/umlaeute/v4l2loopback)
- [vivid Documentation](https://www.kernel.org/doc/html/latest/admin-guide/media/vivid.html)
- [V4L2 Testing Tools](https://linuxtv.org/wiki/index.php/V4L2_Test_Suite)
- [V4L2 API Documentation](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/v4l2.html)
- [Go Testing Package](https://pkg.go.dev/testing)
