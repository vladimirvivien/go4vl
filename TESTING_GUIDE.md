# Testing Guide for go4vl

This guide provides a summary of testing approaches for go4vl. For detailed instructions, see `test/README.md`.

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

#### Setup

```bash
# Load the module with options
sudo modprobe v4l2loopback devices=2 video_nr=20,21 card_label="Fake Camera 1","Fake Camera 2" exclusive_caps=1

# Verify devices created
ls -la /dev/video20 /dev/video21
v4l2-ctl --list-devices

# Remove module when done
sudo modprobe -r v4l2loopback
```

#### Feeding Test Data

```bash
# Using ffmpeg to feed test pattern
ffmpeg -re -f lavfi -i testsrc=size=640x480:rate=30 -pix_fmt yuyv422 -f v4l2 /dev/video20

# Using gstreamer
gst-launch-1.0 videotestsrc ! v4l2sink device=/dev/video20

# Feed an MP4 file
ffmpeg -re -i test.mp4 -pix_fmt yuyv422 -f v4l2 /dev/video20

# Feed static image
ffmpeg -loop 1 -i test.jpg -pix_fmt yuyv422 -f v4l2 /dev/video20
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

## Testing Strategy

The go4vl test suite uses real V4L2 devices (no mock implementations):

1. **Unit Tests** - Test packages in isolation
   ```bash
   go test -v ./device ./v4l2 ./imgsupport
   ```

2. **Integration Tests** - Test with v4l2loopback virtual devices
   ```bash
   go test -v -tags=integration ./test/...
   ```

### Key Features

- **Automatic v4l2loopback setup** - TestMain handles module loading/unloading
- **Dynamic device selection** - Avoids conflicts with existing devices (starts from /dev/video40)
- **Fallback to existing devices** - Uses available devices when v4l2loopback setup fails
- **No environment variables required** - Uses build tags for test selection
- **No build tools needed** - All commands use standard Go tooling

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

```dockerfile
# Dockerfile.test
FROM golang:1.21-bullseye

# Install V4L2 tools and dependencies
RUN apt-get update && apt-get install -y \
    v4l-utils \
    v4l2loopback-dkms \
    ffmpeg \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Set up working directory
WORKDIR /app

# Copy source code
COPY . .

# Install v4l2loopback module
RUN apt-get update && apt-get install -y kmod

# Test script
COPY <<'EOF' /test.sh
#!/bin/bash
# Load v4l2loopback if possible (requires privileged mode)
if [ -w /dev ]; then
    modprobe v4l2loopback devices=1 video_nr=20 || true
fi

# Run tests
go test -v ./...
EOF

RUN chmod +x /test.sh

CMD ["/test.sh"]
```

Run with:
```bash
# Build test image
docker build -f Dockerfile.test -t go4vl-test .

# Run tests (privileged needed for kernel modules)
docker run --privileged --device=/dev/video0 go4vl-test
```

## Quick Start

```bash
# Run unit tests
go test -v ./device ./v4l2 ./imgsupport

# Run integration tests (with automatic setup)
sudo go test -v -tags=integration ./test/...

# Generate coverage report
go test -coverprofile=coverage.out ./device ./v4l2 ./imgsupport
go tool cover -html=coverage.out -o coverage.html
```

## Best Practices

1. **Real V4L2 testing only** - No mock devices, uses v4l2loopback or real hardware
2. **Automatic setup/teardown** - TestMain handles v4l2loopback lifecycle
3. **Dynamic device selection** - Avoids conflicts with existing devices
4. **Build tags for test separation** - Use `-tags=integration` for integration tests
5. **Graceful degradation** - Tests skip when devices unavailable

## Troubleshooting

### v4l2loopback not loading
```bash
# Check kernel module support
lsmod | grep v4l2loopback

# Check dmesg for errors
sudo dmesg | grep v4l2loopback

# Ensure headers installed
sudo apt-get install linux-headers-$(uname -r)
```

### Permission issues
```bash
# Add user to video group
sudo usermod -a -G video $USER

# Logout and login again
```

### No frames received
```bash
# Check if producer is running
ps aux | grep ffmpeg

# Verify device is readable
v4l2-ctl -d /dev/video20 --all
```

## References

- [v4l2loopback GitHub](https://github.com/umlaeute/v4l2loopback)
- [vivid Documentation](https://www.kernel.org/doc/html/latest/admin-guide/media/vivid.html)
- [V4L2 Testing Tools](https://linuxtv.org/wiki/index.php/V4L2_Test_Suite)
