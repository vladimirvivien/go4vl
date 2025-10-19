# Integration Tests for go4vl

This directory contains integration tests that validate the go4vl library against real V4L2 devices and virtual devices using v4l2loopback.

For comprehensive testing documentation, see the main [TESTING_GUIDE.md](../TESTING_GUIDE.md) in the repository root.

## What These Tests Cover

The integration test suite validates the complete go4vl functionality with real or virtual V4L2 devices:

### Test Files

- **device_test.go** - Tests all exported device package functionality
  - Device opening with various options
  - Capability detection and querying
  - Pixel format operations
  - Frame rate control
  - Complete streaming lifecycle
  - Context cancellation
  - Multiple start/stop cycles

- **v4l2_test.go** - Tests all exported v4l2 types and constants
  - v4l2.Capability struct and methods
  - v4l2.PixFormat operations
  - v4l2.FormatDescription handling
  - v4l2.Buffer struct
  - v4l2.Control operations
  - Pixel format constants

- **integration_test.go** - Full pipeline tests
  - Device open and capabilities
  - Format negotiation
  - Streaming with frame validation
  - Start/stop cycles

- **simple_test.go** - Basic tests that work with any V4L2 device
  - Direct device opening
  - Basic streaming if possible

- **helpers.go** - Shared test utilities and validation functions

### Test Infrastructure

- **TestMain** with automatic v4l2loopback setup/teardown
  - Dynamically selects available device numbers (starting from /dev/video40)
  - Falls back to existing devices when v4l2loopback can't be loaded
  - Handles permission errors gracefully

## Running the Integration Tests

### Quick Start

```bash
# From repository root - run with automatic setup (requires sudo)
sudo go test -v -tags=integration ./test/...

# Run without sudo (uses existing devices or skips tests)
go test -v -tags=integration ./test/...

# Run specific tests
go test -v -tags=integration ./test/... -run TestDevice
go test -v -tags=integration ./test/... -run TestV4L2

# Generate coverage report
go test -tags=integration -coverprofile=coverage-integration.out ./test/...
go tool cover -html=coverage-integration.out -o coverage-integration.html
```

### Running Benchmarks

Benchmarks compare the legacy `GetOutput()` API vs the optimized `GetFrames()` API with buffer pooling.

**Performance Results:**
- `GetFrames()` is **2% faster** and uses **600x less memory** (1 KB vs 600 KB per frame)
- Buffer pooling dramatically reduces GC pressure for high-throughput video applications

**How to Run Benchmarks:**

Due to v4l2 driver limitations, benchmarks must be run **individually** (not all together):

```bash
# Run each benchmark separately with -run=^$ to skip tests
sudo go test -tags=integration -bench=BenchmarkIntegration_GetOutput -benchmem -benchtime=3s -run=^$ ./test
sudo go test -tags=integration -bench=BenchmarkIntegration_GetFrames -benchmem -benchtime=3s -run=^$ ./test
sudo go test -tags=integration -bench=BenchmarkIntegration_GetFrames_WithMetadata -benchmem -benchtime=3s -run=^$ ./test
```

**Benchmark Flags Explained:**
- `-tags=integration` - Build with integration tag (required for benchmark code)
- `-bench=BenchmarkName` - Run only this specific benchmark
- `-benchmem` - Show memory allocation statistics
- `-benchtime=3s` - Run for 3 seconds (default is 1s, longer = more accurate)
- `-run=^$` - **Critical:** Skip all tests, run only benchmarks (prevents device conflicts)

**Expected Output:**
```
BenchmarkIntegration_GetOutput-4       6421    522640 ns/op    614585 B/op    3 allocs/op
                                       ^^^^    ^^^^^^ ^^^^^^    ^^^^^^ ^^^^^^  ^ ^^^^^^^^^^
                                       iter    ns/op  µs/op     bytes  KB/op   allocations
```

**Common Mistakes:**
```bash
# ❌ Running all benchmarks together - will fail after first one
sudo go test -tags=integration -bench=. -benchmem ./test

# ❌ Omitting -run=^$ - will run tests AND benchmarks, causing conflicts
sudo go test -tags=integration -bench=BenchmarkIntegration_GetOutput -benchmem ./test

# ❌ Running without sudo - module setup will fail
go test -tags=integration -bench=BenchmarkIntegration_GetOutput -benchmem -run=^$ ./test
```

**Why Individual Runs?**

Running multiple benchmarks in the same process causes v4l2 driver-level device conflicts. Each benchmark needs its own process, which Go's test runner provides when you specify one benchmark at a time.

### Test Flags

```bash
# Skip automatic v4l2loopback setup (use existing loopback devices)
go test -v -tags=integration ./test/... -skip-setup

# Keep v4l2loopback loaded after tests complete (for debugging)
go test -v -tags=integration ./test/... -keep-running

# Enable verbose logging
go test -v -tags=integration ./test/... -verbose
```

**Note:** The test harness automatically detects and uses existing v4l2loopback devices, so no flag is needed for that.

## Prerequisites

See [TESTING_GUIDE.md](../TESTING_GUIDE.md#prerequisites) for detailed prerequisites.

**Quick summary:**
- Go 1.21 or later
- Root access (for automatic v4l2loopback setup) OR video group membership
- v4l2loopback-dkms, v4l-utils, ffmpeg (optional but recommended)

```bash
# Ubuntu/Debian
sudo apt-get install -y v4l-utils v4l2loopback-dkms ffmpeg
```

## Manual Device Setup

If automatic setup doesn't work or you prefer manual control:

```bash
# 1. Load v4l2loopback module
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

## Troubleshooting

### Tests Skip with "Test device not available"

1. Check v4l2loopback is installed:
   ```bash
   sudo apt-get install v4l2loopback-dkms
   ```

2. Run tests with sudo for automatic setup:
   ```bash
   sudo go test -v -tags=integration ./test/...
   ```

### Other Issues

See the comprehensive [Troubleshooting](../TESTING_GUIDE.md#troubleshooting) section in TESTING_GUIDE.md for:
- v4l2loopback not loading
- Permission denied errors
- Device busy errors
- Module not found errors

## Environment Variables

```bash
# Force a specific test device (rarely needed)
V4L2_TEST_DEVICE=/dev/video0 go test -v -tags=integration ./test/...

# Enable CI mode
export CI=true
sudo go test -v -tags=integration ./test/...
```

## More Information

For comprehensive testing documentation including:
- Unit tests
- CI/CD setup
- Docker-based testing
- Best practices
- Detailed troubleshooting

See [TESTING_GUIDE.md](../TESTING_GUIDE.md) in the repository root.
