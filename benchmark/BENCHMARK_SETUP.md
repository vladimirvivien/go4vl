# Benchmark Setup Guide

This guide explains how to set up reproducible benchmarking using v4l2loopback and ffmpeg.

## Why Use v4l2loopback?

Using a virtual test device provides:

- **Reproducible results**: Same test pattern every time
- **No hardware dependency**: Works in CI/CD
- **Controlled conditions**: Known resolution, FPS, format
- **No camera quirks**: Eliminates hardware variability

## Prerequisites

### Install Dependencies

```bash
# Ubuntu/Debian
sudo apt install v4l2loopback-dkms ffmpeg v4l-utils

# Arch
sudo pacman -S v4l2loopback-dkms ffmpeg v4l-utils

# Fedora
sudo dnf install v4l2loopback ffmpeg v4l-utils
```

### Verify Installation

```bash
# Check if v4l2loopback module is available
modinfo v4l2loopback

# Check ffmpeg
ffmpeg -version

# Check v4l-utils
v4l2-ctl --version
```

## Quick Start

### Option 1: Automatic Setup (Recommended)

Run the setup script that handles everything:

```bash
cd benchmark
sudo go run setup_loopback.go
```

This will:
1. Load v4l2loopback module
2. Create `/dev/video50`
3. Start ffmpeg test pattern
4. Keep running until you press Ctrl+C

Then in another terminal, run benchmarks:

```bash
cd benchmark/frame_capture
go run main.go -device /dev/video50
```

### Option 2: Manual Setup

```bash
# Load module
sudo modprobe v4l2loopback devices=1 video_nr=50 card_label=benchmark exclusive_caps=1

# Start test pattern
ffmpeg -re -f lavfi -i testsrc=size=640x480:rate=30 -pix_fmt yuyv422 -f v4l2 /dev/video50 &

# Run benchmark
cd benchmark/frame_capture
go run main.go -device /dev/video50

# Cleanup
killall ffmpeg
sudo modprobe -r v4l2loopback
```

### Option 3: Use Existing Camera

If you have a physical webcam:

```bash
# Find your device
ls -l /dev/video*
v4l2-ctl --list-devices

# Run benchmark directly
cd benchmark/frame_capture
go run main.go -device /dev/video0
```

## Setup Script Options

The `setup_loopback.go` script accepts flags:

```bash
# Different device number
sudo go run setup_loopback.go -device 40

# Different resolution and FPS
sudo go run setup_loopback.go -width 1920 -height 1080 -fps 60

# Different test pattern
sudo go run setup_loopback.go -pattern smptebars  # Color bars
sudo go run setup_loopback.go -pattern color      # Solid color
sudo go run setup_loopback.go -pattern testsrc    # Default test pattern
```

## Test Patterns Available

- **testsrc**: Standard test pattern with moving elements
- **smptebars**: SMPTE color bars
- **color**: Solid color (configurable)
- **mandelbrot**: Animated Mandelbrot set
- **life**: Conway's Game of Life

## Troubleshooting

### Module Not Found

```bash
# Install v4l2loopback
sudo apt install v4l2loopback-dkms

# Rebuild module for current kernel
sudo dkms install v4l2loopback/$(modinfo v4l2loopback | grep ^version | awk '{print $2}')
```

### Device Already Exists

```bash
# Unload existing module
sudo modprobe -r v4l2loopback

# Re-run setup
sudo go run setup_loopback.go
```

### Permission Denied

```bash
# Add user to video group
sudo usermod -a -G video $USER
# Log out and back in

# Or temporarily
sudo chmod 666 /dev/video50
```

### FFmpeg Errors

```bash
# Check if device is accessible
v4l2-ctl -d /dev/video50 --all

# Check ffmpeg can write to it
ffmpeg -f lavfi -i testsrc=size=640x480:rate=1 -vframes 1 -f v4l2 /dev/video50
```

## CI/CD Integration

For automated testing in CI/CD:

```yaml
# GitHub Actions example
- name: Setup v4l2loopback
  run: |
    sudo apt-get update
    sudo apt-get install -y v4l2loopback-dkms ffmpeg
    sudo modprobe v4l2loopback devices=1 video_nr=50 exclusive_caps=1

- name: Start test pattern
  run: |
    ffmpeg -re -f lavfi -i testsrc=size=640x480:rate=30 \
           -pix_fmt yuyv422 -f v4l2 /dev/video50 &
    sleep 2

- name: Run benchmarks
  run: |
    cd benchmark/frame_capture
    go run main.go -device /dev/video50 -duration 5s
```

## Multiple Test Devices

To create multiple devices for testing:

```bash
# Load with multiple devices
sudo modprobe v4l2loopback devices=2 video_nr=50,51 card_label=bench1,bench2 exclusive_caps=1

# Start patterns on both
ffmpeg -re -f lavfi -i testsrc=size=640x480:rate=30 -pix_fmt yuyv422 -f v4l2 /dev/video50 &
ffmpeg -re -f lavfi -i smptebars=size=1280x720:rate=30 -pix_fmt yuyv422 -f v4l2 /dev/video51 &
```

## Benchmarking Best Practices

1. **Use v4l2loopback for baseline measurements**
   - Eliminates hardware variability
   - Reproducible results

2. **Test with real hardware for final validation**
   - Verify real-world performance
   - Check device-specific optimizations

3. **Run benchmarks multiple times**
   - Average results across runs
   - Look for consistency

4. **Use consistent test patterns**
   - Same resolution and FPS
   - Same pixel format

5. **Monitor system load**
   - Run benchmarks on idle system
   - Check for background processes

## Next Steps

Once your test device is set up:

1. Run baseline benchmarks (see `frame_capture/README.md`)
2. Analyze profiles with pprof (see main `README.md`)
3. Implement optimizations
4. Re-run benchmarks to measure improvement

## Cleanup

Always cleanup when done:

```bash
# Stop ffmpeg
killall ffmpeg

# Unload module
sudo modprobe -r v4l2loopback

# Verify
lsmod | grep v4l2loopback  # Should be empty
ls /dev/video50  # Should not exist
```
