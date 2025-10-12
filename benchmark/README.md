# go4vl Benchmarking Suite

Performance profiling and benchmarking tools for testing the **actual go4vl code** (device package, v4l2 package, buffer management, frame copying).

## Quick Start

### Loopback Mode (Default)

Automatically sets up v4l2loopback device and benchmarks real go4vl code:

```bash
cd runner
sudo go run main.go -single -duration 5s
```

### Physical Device Mode

Test with your actual webcam:

```bash
cd runner
go run main.go -single -device /dev/video0 -duration 5s
```

## What It Tests

The benchmark exercises **real go4vl code paths**:

1. **device.Open()** - Device initialization and configuration
2. **v4l2.InitBuffers()** - Buffer allocation
3. **v4l2.MapMemoryBuffers()** - Memory mapping (zero-copy)
4. **v4l2.QueueBuffer() / DequeueBuffer()** - Buffer queue/dequeue operations
5. **Frame copying** - Memory allocation and copying from mapped buffers
6. **Channel operations** - Frame delivery via Go channels
7. **Context handling** - Timeout and cancellation

This is **not a synthetic benchmark** - it runs the actual V4L2 syscalls and measures real performance.

## Structure

```
benchmark/
├── loopback/
│   └── main.go      # v4l2loopback device setup utility
└── runner/
    └── main.go      # Benchmark runner
```

## Usage

### Single Benchmark

Run one benchmark with custom parameters:

```bash
sudo go run main.go -single [flags]
```

**Flags:**
- `-single` - Run single benchmark (vs multi-scenario)
- `-duration` - Capture duration (default: 10s)
- `-width` - Frame width (default: 640)
- `-height` - Frame height (default: 480)
- `-fps` - Frames per second (default: 30)
- `-format` - Pixel format: MJPEG, YUYV, H264 (default: MJPEG)
- `-buffers` - Number of buffers (default: 4)
- `-verbose` - Verbose output

**Device Selection:**
- **Default**: Sets up loopback device (requires sudo, ffmpeg, v4l2loopback-dkms)
- `-device /dev/videoN` - Use real device instead
- `-loopback-num N` - Loopback device number (default: 50)
- `-test-pattern` - FFmpeg pattern: testsrc, smptebars, color=<color>

**Profiling:**
- `-cpuprofile <file>` - Write CPU profile
- `-memprofile <file>` - Write memory profile
- `-trace <file>` - Write execution trace

### Multi-Scenario Benchmarks

Run all predefined scenarios:

```bash
sudo go run main.go [flags]
```

**Flags:**
- `-duration` - Duration per scenario (default: 10s)
- `-output` - Output directory (default: results_<timestamp>)
- `-scenario` - Run specific scenario by name
- `-list` - List available scenarios
- `-device /dev/videoN` - Use real device instead of loopback

## Examples

### Quick 5 Second Test

```bash
sudo go run main.go -single -duration 5s
```

### Loopback with Different Test Patterns

```bash
# Default test pattern
sudo go run main.go -single -duration 10s

# SMPTE color bars
sudo go run main.go -single -test-pattern smptebars

# Solid red
sudo go run main.go -single -test-pattern color=red
```

### 1080p Performance Test

```bash
sudo go run main.go -single \
  -width 1920 -height 1080 -fps 60 \
  -duration 30s
```

### Test Physical Webcam

```bash
# Usually /dev/video0
go run main.go -single -device /dev/video0 -duration 10s

# Specific device
go run main.go -single -device /dev/video2 -duration 10s
```

### Full Profiling Session

```bash
sudo go run main.go -single \
  -duration 30s \
  -cpuprofile cpu.prof \
  -memprofile mem.prof \
  -trace trace.out

# Analyze
go tool pprof cpu.prof
go tool pprof mem.prof
go tool trace trace.out
```

### Run All Scenarios

```bash
# With loopback (default)
sudo go run main.go -duration 15s -output baseline_results

# With real device
go run main.go -device /dev/video0 -duration 15s -output webcam_results
```

### Run Specific Scenario

```bash
sudo go run main.go -scenario baseline_720p_mjpeg
```

## Available Scenarios

```bash
$ go run main.go -list
Available benchmark scenarios:
  baseline_480p_mjpeg: 640x480 @ 30 fps (MJPEG)
  baseline_720p_mjpeg: 1280x720 @ 30 fps (MJPEG)
  baseline_1080p_mjpeg: 1920x1080 @ 30 fps (MJPEG)
  format_480p_yuyv: 640x480 @ 30 fps (YUYV)
  fps_720p_15fps: 1280x720 @ 15 fps (MJPEG)
  fps_720p_30fps: 1280x720 @ 30 fps (MJPEG)
  fps_720p_60fps: 1280x720 @ 60 fps (MJPEG)
```

## Key Metrics

### Capture Statistics
- **Frames Captured/Dropped** - Reliability indicator
- **Average FPS** - Should match target FPS
- **Frame Timing** - Min/Max/Avg frame time consistency

### Memory Statistics
- **Allocs per Frame** - Lower is better (target < 100KB)
- **GC Frequency** - Fewer GC runs = smoother performance
- **Total Allocations** - Memory churn indicator

### Performance Indicators
- **CPU Time per Frame** - Processing efficiency
- **Throughput (MB/s)** - Data handling capacity
- **GC Pause Times** - Latency impact

## Performance Goals

| Resolution | FPS | CPU Usage | Memory/Frame | Frame Drops |
|------------|-----|-----------|--------------|-------------|
| 640x480    | 30  | < 5%      | < 100KB      | 0           |
| 1280x720   | 30  | < 10%     | < 150KB      | 0           |
| 1920x1080  | 30  | < 20%     | < 200KB      | 0           |
| 1920x1080  | 60  | < 35%     | < 200KB      | < 1%        |

## How It Works

### Loopback Mode (Default)
1. **Checks prerequisites** - Verifies ffmpeg and v4l2loopback-dkms installed
2. **Loads v4l2loopback** - Creates /dev/video50 (or specified device number)
3. **Starts FFmpeg** - Generates test pattern at specified resolution/FPS
4. **Runs benchmark** - Captures frames via **actual V4L2 API** and go4vl code
5. **Cleans up** - Stops FFmpeg and unloads v4l2loopback module

### Device Mode
1. **Opens device** - Connects to specified video device via **device.Open()**
2. **Configures format** - Sets resolution, FPS, pixel format via **v4l2** package
3. **Maps buffers** - Memory maps buffers via **v4l2.MapMemoryBuffers()**
4. **Captures frames** - Reads from device via **v4l2.DequeueBuffer()**
5. **Copies frames** - Allocates and copies data via **device package**
6. **Metrics collection** - Tracks performance and memory of actual code

## Prerequisites

### Loopback Mode (Default)

#### Ubuntu/Debian
```bash
sudo apt install ffmpeg v4l2loopback-dkms
sudo dkms autoinstall
```

#### Arch Linux
```bash
sudo pacman -S ffmpeg v4l2loopback-dkms
```

#### Fedora
```bash
sudo dnf install ffmpeg v4l2loopback
```

### Device Mode
No installation needed - just requires a working /dev/videoN device

## Troubleshooting

### Loopback Mode Issues

#### Module Not Found Error
```
ERROR: v4l2loopback kernel module not installed
```

**Solution**: Install v4l2loopback-dkms and rebuild:
```bash
sudo apt install v4l2loopback-dkms
sudo dkms autoinstall
```

#### FFmpeg Not Found
```
ERROR: ffmpeg not found
```

**Solution**:
```bash
sudo apt install ffmpeg
```

#### Permission Denied
Loopback mode requires `sudo` to load/unload kernel modules:
```bash
sudo go run main.go -single
```

#### Device Already Exists
If `/dev/video50` already exists from a previous run:

```bash
# Clean up manually
sudo modprobe -r v4l2loopback
```

Or use a different device number:
```bash
sudo go run main.go -single -loopback-num 51
```

### Device Mode Issues

#### No Device Found
```bash
# List available devices
ls -l /dev/video*

# Check device capabilities
v4l2-ctl --list-devices
```

#### Permission Denied
Add your user to the `video` group:
```bash
sudo usermod -a -G video $USER
# Log out and back in
```

### General

#### Interrupt Handling
Press `Ctrl+C` to cleanly stop benchmarks. The program automatically:
- Stops FFmpeg (loopback mode)
- Unloads v4l2loopback (loopback mode)
- Closes device connections
- Cleans up resources

## Analyzing Results

### Multi-Scenario Output

```bash
cd results_<timestamp>/

# View all results
cat summary.txt

# Individual results
cat baseline_480p_mjpeg_results.txt

# Analyze profiles
go tool pprof baseline_480p_mjpeg_cpu.prof
go tool pprof baseline_480p_mjpeg_mem.prof
go tool trace baseline_480p_mjpeg_trace.out
```

### Using pprof

```bash
# Interactive mode
go tool pprof cpu.prof
> top10         # Top 10 CPU consumers
> list main.*   # Annotated source
> web           # Visual graph (requires graphviz)

# Web UI
go tool pprof -http=:8080 cpu.prof
```

### Using trace

```bash
go tool trace trace.out
```

Opens browser showing:
- Goroutine execution timeline
- GC events and pauses
- System calls and blocking
- Goroutine analysis

## Comparison Workflows

### Before/After Optimization

```bash
# Baseline
sudo go run main.go -output before_opt

# Make code changes...

# After optimization
sudo go run main.go -output after_opt

# Compare
diff before_opt/summary.txt after_opt/summary.txt
```

### Resolution Comparison

```bash
sudo go run main.go -single -width 640 -height 480 -output 480p
sudo go run main.go -single -width 1280 -height 720 -output 720p
sudo go run main.go -single -width 1920 -height 1080 -output 1080p
```

## Next Steps

See `/docs/performance-optimization-plan.md` for the complete optimization roadmap and implementation strategy.
