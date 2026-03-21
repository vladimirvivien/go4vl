# Integration Tests for go4vl

This directory contains integration tests that validate the go4vl library against real V4L2 devices and v4l2loopback virtual devices.

For comprehensive testing documentation, see the main [TESTING_GUIDE.md](../TESTING_GUIDE.md) in the repository root.

## What These Tests Cover

The integration test suite validates the complete go4vl functionality:

### Test Files

- **device_test.go** - Device package functionality (opening, capabilities, pixel formats, streaming, context cancellation)
- **v4l2_test.go** - V4L2 types and constants (capability, pixel format, buffer, control operations)
- **integration_test.go** - Full pipeline tests (format negotiation, streaming, start/stop cycles)
- **simple_test.go** - Basic tests that work with any V4L2 device
- **standard_test.go** - Analog video standard tests
- **video_io_test.go** - Video input/output enumeration
- **audio_io_test.go** - Audio input/output enumeration
- **ext_controls_test.go** - Extended controls API
- **dv_timings_test.go** - Digital video timing support
- **tuner_modulator_test.go** - Tuner and modulator support
- **codec_test.go** - Codec capability and command tests

## Running the Integration Tests

### Quick Start

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

### Test Flags

| Flag | Value | Description |
|------|-------|-------------|
| `-use-device` | `auto` | Auto-discover real V4L2 devices |
| `-use-device` | `/dev/video0` | Use specific real device |
| `-use-device-emulation` | `auto` | Auto-discover or load v4l2loopback |
| `-use-device-emulation` | `/dev/video42,/dev/video43` | Use specific loopback devices |
| `-keep-running` | | Keep v4l2loopback loaded after tests |
| `-verbose` | | Enable verbose logging |

**Note:** Flags must be passed after `-args` when using `go test`.

### Running Benchmarks

Benchmarks compare the legacy `GetOutput()` API vs the optimized `GetFrames()` API with buffer pooling.

Due to V4L2 driver limitations, benchmarks must be run **individually**:

```bash
go test -tags=integration -bench=BenchmarkIntegration_GetOutput -benchmem -benchtime=3s -run=^$ ./test -args -use-device-emulation=/dev/video42,/dev/video43
go test -tags=integration -bench=BenchmarkIntegration_GetFrames -benchmem -benchtime=3s -run=^$ ./test -args -use-device-emulation=/dev/video42,/dev/video43
```

## Prerequisites

See [TESTING_GUIDE.md](../TESTING_GUIDE.md#prerequisites) for detailed prerequisites.

**Quick summary:**
- Go 1.25 or later
- For emulation: `v4l2loopback-dkms` and `ffmpeg` packages
- For real devices: video group membership or root access

```bash
# Ubuntu/Debian (for v4l2loopback emulation)
sudo apt-get install -y v4l2loopback-dkms ffmpeg
```

## Environment Variables

```bash
# Force a specific test device (rarely needed)
V4L2_TEST_DEVICE=/dev/video0 go test -v -tags=integration ./test/...
```

## Troubleshooting

### Tests Skip with "No V4L2 device available"

```bash
# Load v4l2loopback
sudo modprobe v4l2loopback devices=2 video_nr=42,43 exclusive_caps=0

# If module not found, install it
sudo apt-get install v4l2loopback-dkms
```

### Other Issues

See the comprehensive [Troubleshooting](../TESTING_GUIDE.md#troubleshooting) section in TESTING_GUIDE.md.

## More Information

For comprehensive testing documentation including unit tests, CI/CD setup, and best practices, see [TESTING_GUIDE.md](../TESTING_GUIDE.md).
