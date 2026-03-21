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
- **Integration Tests** - Test with vivid virtual devices or real hardware

### Key Features

- **Two device modes** - Real hardware (`-use-device=auto`) or v4l2loopback emulation (`-use-device-emulation=auto`)
- **Automatic setup** - TestMain loads v4l2loopback when needed for emulation
- **Graceful degradation** - Tests skip appropriately when devices or features aren't available

## Prerequisites

### 1. Go Installation

Go 1.25 or later is required.

```bash
go version
```

### 2. User Permissions

For integration tests, you need either:
- Root access to load kernel modules (for vivid setup)
- Video group membership for existing devices

```bash
# Add user to video group
sudo usermod -a -G video $USER
# Logout and login for changes to take effect
```

### 3. Testing Tools (for emulated testing)

v4l2loopback emulation is used in CI and on machines without video hardware. If you have a real camera (e.g., on Raspberry Pi), use `--use-device` instead and skip this section.

```bash
# Ubuntu/Debian (x86_64)
sudo apt-get install v4l2loopback-dkms ffmpeg

# Raspberry Pi OS
sudo apt install v4l2loopback-dkms v4l2loopback-utils raspberrypi-kernel-headers

# From source (any Linux, requires kernel headers for your running kernel)
git clone https://github.com/umlaeute/v4l2loopback.git
cd v4l2loopback && make && sudo make install
sudo depmod -a
```

- **v4l2loopback** - Kernel module for virtual video devices
- **ffmpeg** - Feeds test patterns to loopback devices

## Testing Approaches

### 1. v4l2loopback - Virtual Video Device (Recommended for CI)

**v4l2loopback** creates virtual V4L2 devices. Combined with ffmpeg to feed test patterns, it provides realistic testing without hardware.

```bash
# Load the module
sudo modprobe v4l2loopback devices=2 video_nr=42,43 card_label="go4vl_test_1,go4vl_test_2" exclusive_caps=0

# Run tests with specific loopback devices
go test -v -tags=integration ./test/... -args -use-device-emulation=/dev/video42,/dev/video43

# Or auto-discover existing loopback devices
go test -v -tags=integration ./test/... -args -use-device-emulation=auto
```

### 2. Real Hardware

If you have a USB webcam or other V4L2 device:

```bash
# Auto-discover real devices
go test -v -tags=integration ./test/... -args -use-device=auto

# Use a specific device
go test -v -tags=integration ./test/... -args -use-device=/dev/video0
```

## Unit Tests

Unit tests validate the core functionality without requiring V4L2 devices.

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

# Clean test cache
go clean -testcache
```

## Integration Tests

Integration tests validate the complete streaming pipeline with real or virtual V4L2 devices. See `test/README.md` for detailed information.

### Running Integration Tests

```bash
# With v4l2loopback emulation (auto-discover or load)
go test -v -tags=integration ./test/... -args -use-device-emulation=auto

# With specific loopback devices
go test -v -tags=integration ./test/... -args -use-device-emulation=/dev/video42,/dev/video43

# With real hardware (auto-discover)
go test -v -tags=integration ./test/... -args -use-device=auto

# With specific real device
go test -v -tags=integration ./test/... -args -use-device=/dev/video0

# Auto-detect (tries loopback first, then real devices)
go test -v -tags=integration ./test/...

# Run specific tests
go test -v -tags=integration ./test/... -run TestDevice
go test -v -tags=integration ./test/... -run TestV4L2
```

## Test Flags

| Flag | Value | Description |
|------|-------|-------------|
| `-use-device` | `auto` | Auto-discover real V4L2 devices |
| `-use-device` | `/dev/video0` | Use specific real device |
| `-use-device-emulation` | `auto` | Auto-discover or load v4l2loopback |
| `-use-device-emulation` | `/dev/video42,/dev/video43` | Use specific loopback devices |
| `-keep-running` | | Keep v4l2loopback loaded after tests |
| `-verbose` | | Enable verbose logging |

All flags are passed after `-args`:
```bash
go test -v -tags=integration ./test/... -args -use-device=auto
```

## Environment Variables

### V4L2_TEST_DEVICE

Force a specific test device (rarely needed):

```bash
V4L2_TEST_DEVICE=/dev/video0 go test -v -tags=integration ./test/...
```

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
        go-version: '1.25'

    - name: Run Unit Tests
      run: go test -v ./device ./v4l2 ./imgsupport

  integration-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25'

    - name: Install build dependencies
      run: |
        sudo apt-get update
        sudo apt-get install -y build-essential linux-headers-generic ffmpeg

    - name: Setup v4l2loopback
      run: |
        sudo apt-get install -y linux-headers-$(uname -r) linux-modules-extra-$(uname -r)
        sudo modprobe videodev
        git clone https://github.com/umlaeute/v4l2loopback.git /tmp/v4l2loopback
        cd /tmp/v4l2loopback && make KERNEL_DIR=/lib/modules/$(uname -r)/build
        sudo insmod /tmp/v4l2loopback/v4l2loopback.ko devices=2 video_nr=10,11 card_label="go4vl_test_1,go4vl_test_2" exclusive_caps=0
        sudo chmod 666 /dev/video10 /dev/video11

    - name: Run Integration Tests
      run: |
        go test -v -tags=integration ./test/... -args -use-device-emulation=/dev/video10,/dev/video11
```

## Troubleshooting

### Tests Skip with "No V4L2 device available"

1. **Load v4l2loopback**:
   ```bash
   sudo modprobe v4l2loopback devices=2 video_nr=42,43 exclusive_caps=0
   ```

2. **Install v4l2loopback** (if not found):
   ```bash
   sudo apt-get install v4l2loopback-dkms
   ```

3. **Check devices exist**:
   ```bash
   ls -la /dev/video*
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

### Device Busy

```
Failed to start streaming: device or resource busy
```

**Solutions**:
1. Check if another application is using the device:
   ```bash
   fuser /dev/video*
   lsof /dev/video*
   ```
2. Ensure tests are not running concurrently on the same device
3. Reset loopback devices for a clean state:
   ```bash
   # Kill any ffmpeg processes holding devices
   sudo pkill -9 ffmpeg

   # Unload and reload v4l2loopback
   sudo modprobe -r v4l2loopback
   sudo modprobe v4l2loopback devices=2 video_nr=42,43 card_label="go4vl_test_1,go4vl_test_2" exclusive_caps=0

   # Clear test cache
   go clean -testcache
   ```

### v4l2loopback Not Loading

```
modprobe: FATAL: Module v4l2loopback not found
```

**Solution**: Install v4l2loopback:
```bash
sudo apt-get install v4l2loopback-dkms
# If still failing, ensure kernel headers are installed
sudo apt-get install linux-headers-$(uname -r)
```

### Test Cache Issues

```bash
go clean -testcache
go test -count=1 -v ./device ./v4l2 ./imgsupport
```

## Performance Testing

### Frame Pool Benchmarks

```bash
go test -bench=BenchmarkFramePool -benchmem ./device
go test -bench=Benchmark -benchmem ./device
```

Expected results (640x480 YUYV frame, ~614KB):

```
BenchmarkFramePool_Get-4           45,450,000 ops/sec       22 ns/op      26 B/op    1 allocs/op
BenchmarkDirectAllocation-4            36,260 ops/sec   27,592 ns/op  614,400 B/op    1 allocs/op
```

## Quick Reference

```bash
# Unit Tests
go test -v ./device ./v4l2 ./imgsupport

# Integration Tests (v4l2loopback emulation)
go test -v -tags=integration ./test/... -args -use-device-emulation=auto

# Integration Tests (real device)
go test -v -tags=integration ./test/... -args -use-device=auto

# Integration Tests (auto-detect)
go test -v -tags=integration ./test/...

# Coverage Reports
go test -coverprofile=coverage.out ./device ./v4l2 ./imgsupport
go tool cover -html=coverage.out -o coverage.html

# Benchmarks
go test -bench=. -benchmem ./v4l2
go test -bench=. -benchmem ./device

# Clean Cache
go clean -testcache
```

## References

- [v4l2loopback GitHub](https://github.com/umlaeute/v4l2loopback)
- [V4L2 Testing Tools](https://linuxtv.org/wiki/index.php/V4L2_Test_Suite)
- [V4L2 API Documentation](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/v4l2.html)
- [Go Testing Package](https://pkg.go.dev/testing)
