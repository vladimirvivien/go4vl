# Benchmark Runner

A Go-based benchmark orchestration tool using [gexe](https://github.com/vladimirvivien/gexe) to automate running performance benchmarks.

## Features

- Runs multiple benchmark scenarios automatically
- Builds benchmark program on-the-fly
- Generates CPU, memory profiles, and execution traces
- Creates summary report
- Pure Go implementation (no shell scripts)

## Usage

### Run All Benchmarks

```bash
cd benchmark/runner
go run main.go
```

### Run Specific Device

```bash
go run main.go -device /dev/video1
```

### Run Single Scenario

```bash
# List available scenarios
go run main.go -list

# Run specific scenario
go run main.go -scenario baseline_720p_mjpeg
```

### Custom Duration and Output

```bash
go run main.go -duration 30s -output my_results
```

## Command-Line Flags

- `-device string`: Video device path (default: /dev/video0)
- `-duration string`: Capture duration per benchmark (default: 10s)
- `-output string`: Output directory (default: results_<timestamp>)
- `-scenario string`: Run specific scenario by name (empty = run all)
- `-list`: List available scenarios and exit

## Available Scenarios

The runner includes predefined benchmark scenarios:

**Baseline Benchmarks:**
- `baseline_480p_mjpeg`: 640x480 @ 30fps MJPEG
- `baseline_720p_mjpeg`: 1280x720 @ 30fps MJPEG
- `baseline_1080p_mjpeg`: 1920x1080 @ 30fps MJPEG

**Format Comparison:**
- `format_480p_yuyv`: 640x480 @ 30fps YUYV (raw)

**Frame Rate Tests:**
- `fps_720p_15fps`: 1280x720 @ 15fps MJPEG
- `fps_720p_30fps`: 1280x720 @ 30fps MJPEG
- `fps_720p_60fps`: 1280x720 @ 60fps MJPEG

## Output Files

For each scenario, the runner generates:

- `<scenario>_results.txt`: Benchmark output with statistics
- `<scenario>_cpu.prof`: CPU profile for pprof analysis
- `<scenario>_mem.prof`: Memory profile for pprof analysis
- `<scenario>_trace.out`: Execution trace for trace viewer
- `summary.txt`: Combined results from all scenarios

## Examples

### Quick Baseline Check

```bash
# Run just 480p baseline for quick validation
go run main.go -scenario baseline_480p_mjpeg -duration 5s
```

### Full Performance Suite

```bash
# Run all scenarios with longer duration
go run main.go -duration 30s -output full_benchmark
```

### Compare Before/After Optimization

```bash
# Baseline
go run main.go -output before_optimization

# ... make code changes ...

# After optimization
go run main.go -output after_optimization

# Compare results
diff before_optimization/summary.txt after_optimization/summary.txt
```

## Adding Custom Scenarios

Edit `main.go` and add to the `scenarios` slice:

```go
var scenarios = []BenchmarkScenario{
    // ... existing scenarios ...
    {Name: "custom_4k_test", Width: 3840, Height: 2160, FPS: 30, Format: "MJPEG"},
}
```

## Integration with CI/CD

The runner can be used in automated testing:

```bash
# Run benchmarks and check for regressions
go run main.go -duration 10s
if [ $? -eq 0 ]; then
    echo "Benchmarks passed"
else
    echo "Benchmarks failed"
    exit 1
fi
```

## Troubleshooting

### Device Not Found

```bash
# Check available devices
ls -l /dev/video*

# Specify correct device
go run main.go -device /dev/video1
```

### Build Errors

```bash
# Ensure dependencies are installed
cd ../frame_capture
go mod tidy
go build
```

### Permission Denied

```bash
# Add user to video group
sudo usermod -a -G video $USER
# Log out and back in

# Or run with sudo (not recommended)
sudo go run main.go
```
