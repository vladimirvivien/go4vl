# Integration Testing for go4vl

This directory contains comprehensive integration tests that validate the go4vl library against real V4L2 devices and virtual devices using v4l2loopback.

## Test Infrastructure

The test suite includes:
- **TestMain** with automatic v4l2loopback setup/teardown
  - Dynamically selects available device numbers (starting from /dev/video40) to avoid conflicts
  - Falls back to existing devices when v4l2loopback can't be loaded
  - Handles permission errors gracefully with helpful messages
- **Device package tests** (`device_test.go`) - Tests all exported device functionality
- **V4L2 package tests** (`v4l2_test.go`) - Tests all exported v4l2 types and constants
- **Integration tests** (`integration_test.go`) - Complete streaming pipeline tests
- **Simple tests** (`simple_test.go`) - Basic tests that work with any available device
- **Helper utilities** (`helpers.go`) - Shared test utilities and validation functions

## Key Features

- **No environment variables needed** - Tests are controlled solely by the `-tags=integration` build tag
- **Dynamic device selection** - Automatically finds available device numbers to avoid conflicts
- **Graceful degradation** - Tests skip appropriately when devices or permissions aren't available
- **Smart fallback** - Uses existing devices when v4l2loopback setup fails
- **No build tools required** - All tasks use standard Go commands directly
- **No mock devices** - All tests use real or virtual V4L2 devices for authentic testing

## Prerequisites

### 1. User Permissions

Your user must have root access to run the integration tests.

### 2. Install Testing Tools (Recommended)

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y v4l-utils v4l2loopback-dkms ffmpeg

# Fedora
sudo dnf install -y v4l-utils v4l2loopback ffmpeg

# Arch Linux
sudo pacman -S v4l-utils v4l2loopback-dkms ffmpeg
```

## Running Tests

### Common Commands

All testing tasks use standard Go commands directly:

```bash
# Run unit tests
go test -v ./device ./v4l2 ./imgsupport

# Run integration tests (automatic device setup if running with sudo)
go test -v -tags=integration ./test/...

# Run specific test files
go test -v -tags=integration ./test/... -run TestDevice
go test -v -tags=integration ./test/... -run TestV4L2

# Generate coverage report
go test -coverprofile=coverage.out ./device ./v4l2 ./imgsupport
go test -tags=integration -coverprofile=coverage-integration.out ./test/...
go tool cover -html=coverage.out -o coverage.html

# Clean test cache
go clean -testcache
```

### Test Flags

The test suite supports several command-line flags for flexibility:

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

### Manual Setup (Alternative)

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

## Test Categories

### Device Package Tests (`device_test.go`)
Tests all exported functionality from the device package:
- Device opening with various options
- Capability detection and querying
- Pixel format get/set operations
- Frame rate control
- Stream parameter management
- Complete streaming lifecycle
- Context cancellation
- Multiple start/stop cycles
- Video input operations
- Crop capabilities
- File descriptor access

### V4L2 Package Tests (`v4l2_test.go`)
Tests all exported v4l2 types and constants:
- v4l2.Capability struct and methods
- v4l2.PixFormat struct operations
- v4l2.FormatDescription handling
- v4l2.StreamParam configuration
- v4l2.Buffer struct
- v4l2.Control operations
- v4l2.CropCapability
- v4l2.InputInfo
- Pixel format constants
- Buffer flags
- Memory types
- Field types

### Integration Tests (`integration_test.go`)
Full pipeline tests with v4l2loopback:
- Device open and capabilities
- Format negotiation
- Frame rate control
- Streaming with frame validation
- Start/stop cycles
- Context cancellation

### Simple Tests (`simple_test.go`)
Basic tests that work with any V4L2 device:
- Direct device opening
- Basic streaming if possible
- No v4l2loopback requirement

## Environment Variables (Optional)

- `V4L2_TEST_DEVICE=/dev/videoX` - Force a specific test device to use (rarely needed)
- `CI=true` - Enable CI mode for automatic v4l2loopback setup in CI environments

Note: No `RUN_INTEGRATION` variable is needed - tests are controlled by the `-tags=integration` build tag.

## Test Device Configuration

The test suite dynamically selects device numbers to avoid conflicts:
- **Automatic selection**: Starts searching from `/dev/video40` for available numbers
- **Fallback to existing**: Uses `/dev/video0`, `/dev/video1` etc. if they exist
- **Smart conflict avoidance**: Avoids commonly used device numbers:
  - Built-in webcams (typically `/dev/video0`, `/dev/video1`)
  - Common virtual devices (typically `/dev/video10-20`)

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

### Permission Denied

```
Failed to open /dev/video0: permission denied
```

**Solution**: Add user to video group (see Prerequisites)

### v4l2-ctl Not Found

```
v4l2-ctl not found. Install with: apt-get install v4l-utils
```

**Solution**: Install v4l-utils package:
```bash
sudo apt-get install v4l-utils
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

## Testing with Real Hardware

If you have a USB webcam:

```bash
# List all video devices
v4l2-ctl --list-devices

# Test with specific device
V4L2_TEST_DEVICE=/dev/video0 go test -v -tags=integration ./test/...
```

## CI/CD Setup

### GitHub Actions

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2

    - name: Setup Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.21

    - name: Setup V4L2 Testing
      run: |
        sudo apt-get update
        sudo apt-get install -y v4l-utils v4l2loopback-dkms ffmpeg
        sudo modprobe v4l2loopback devices=2 video_nr=42,43

    - name: Run Integration Tests
      run: |
        export CI=true
        go test -v -tags=integration ./test/...
```

### Docker

```dockerfile
FROM golang:1.21

# Install dependencies
RUN apt-get update && apt-get install -y \
    v4l-utils \
    v4l2loopback-dkms \
    ffmpeg \
    kmod

# Copy code
WORKDIR /app
COPY . .

# Run tests
CMD ["sh", "-c", "modprobe v4l2loopback && go test -v -tags=integration ./test/..."]
```

## Benchmarking

Run performance benchmarks:

```bash
# Run all benchmarks
go test -bench=. -tags=integration ./test/...

# Run specific benchmarks
go test -bench=BenchmarkIntegration_FrameCapture -tags=integration ./test/...

# Run benchmarks with memory profiling
go test -bench=. -benchmem -tags=integration ./test/...
```


