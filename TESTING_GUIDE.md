# Testing Guide for `go4vl`

`go4vl` verifies correctness through fast hardware-free unit tests, ABI layout checks under `CGO_ENABLED=0`, race detector runs, and end-to-end integration tests using `v4l2loopback` or physical V4L2 devices.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Unit Tests & ABI Verification](#unit-tests--abi-verification)
- [Integration Tests (`v4l2loopback` and Hardware)](#integration-tests-v4l2loopback-and-hardware)
- [Integration Test Flags](#integration-test-flags)
- [CI Workflow](#ci-workflow)
- [Troubleshooting](#troubleshooting)
- [Performance Benchmarks](#performance-benchmarks)

## Overview

* **Unit Tests (`./v4l2`, `./device`):** Run without V4L2 hardware or CGO by mocking the `ioctlSyscall` seam and testing buffer encoding, frame pooling, MPLANE layout helpers, and iterator cleanup.
* **ABI Layout Verification (`TestABI`):** Verifies wrapper struct sizes, patched C struct layouts, 64-bit time fields, and union offsets under `CGO_ENABLED=0`.
* **Integration Tests (`-tags=integration ./test/...`):** Exercise live kernel ioctls and streaming against virtual `v4l2loopback` nodes or attached camera/codec hardware. Streaming integration tests use an in-process `Device.WriteFrame` producer on loopback devices rather than requiring `ffmpeg` or the `vivid` kernel driver.

## Prerequisites

* **Go 1.26 or later**
* **Video permissions** for physical hardware or pre-loaded loopback nodes:
  ```bash
  sudo usermod -a -G video "$USER"
  # Log out and back in for group changes to apply
  ```
* **`v4l2loopback`** (for virtual device integration testing on hosts without physical cameras):
  ```bash
  # Ubuntu / Debian
  sudo apt-get update && sudo apt-get install -y v4l2loopback-dkms

  # Raspberry Pi OS
  sudo apt-get install -y v4l2loopback-dkms v4l2loopback-utils raspberrypi-kernel-headers

  # Or build from source against the running kernel headers
  git clone https://github.com/umlaeute/v4l2loopback.git /tmp/v4l2loopback
  cd /tmp/v4l2loopback && make && sudo make install && sudo depmod -a
  ```

## Unit Tests & ABI Verification

Run unit tests, multi-arch `CGO_ENABLED=0` builds, `go vet`, and race checks locally:

```bash
# Standard hardware-free unit tests + coverage
CGO_ENABLED=0 go test -v -cover ./v4l2 ./device
CGO_ENABLED=0 go vet ./v4l2 ./device

# ABI layout and UAPI generator verification
CGO_ENABLED=0 go test -v ./v4l2 -run TestABI
go test -v ./v4l2/internal/gen

# Multi-architecture build check
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ./...
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ./...
CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build ./...

# Race detector (requires CGO=1)
CGO_ENABLED=1 go test -race ./v4l2 ./device
```

## Integration Tests (`v4l2loopback` and Hardware)

### With `v4l2loopback`

Load two virtual loopback devices (using `exclusive_caps=0` so nodes advertise both `VIDEO_CAPTURE` and `VIDEO_OUTPUT` capabilities) and pass `-require-streaming` so streaming assertions execute strictly:

```bash
sudo modprobe v4l2loopback devices=2 video_nr=42,43 card_label="go4vl_test_1,go4vl_test_2" exclusive_caps=0

# Run with explicit loopback paths and strict streaming verification
CGO_ENABLED=0 go test -v -tags=integration ./test/... \
  -args -use-device-emulation=/dev/video42,/dev/video43 -require-streaming

# Or auto-discover existing loopback devices (or auto-load via sudo in TestMain)
CGO_ENABLED=0 go test -v -tags=integration ./test/... \
  -args -use-device-emulation=auto
```

### With Physical V4L2 Hardware

```bash
# Auto-discover mounted /dev/video* nodes
CGO_ENABLED=0 go test -v -tags=integration ./test/... -args -use-device=auto

# Target a specific hardware device
CGO_ENABLED=0 go test -v -tags=integration ./test/... -args -use-device=/dev/video0
```

Set `V4L2_TEST_DEVICE=/dev/video0` if you prefer pinning a device path through the environment.

## Integration Test Flags

Pass test harness flags after `-args`:

| Flag | Value / Default | Description |
|---|---|---|
| `-use-device` | `auto` or `/dev/videoN` | Discover or target physical V4L2 devices |
| `-use-device-emulation` | `auto` or `/dev/videoN,/dev/videoM` | Use or auto-load `v4l2loopback` devices |
| `-require-streaming` | `false` | Fail instead of skipping if streaming capture/output cannot run |
| `-keep-running` | `false` | Keep auto-loaded `v4l2loopback` module loaded after teardown |
| `-verbose` | `false` | Enable verbose harness logging |

## CI Workflow

The authoritative GitHub Actions configuration lives in [`.github/workflows/test.yml`](./.github/workflows/test.yml) and runs two jobs on Go 1.26:

1. **`unit-tests`:** Compiles `CGO_ENABLED=0` across `amd64`, `arm64`, `arm`, and `riscv64` (including integration test compilation and `simplecam`/`webcam` examples), runs `CGO_ENABLED=0 go test -v -cover ./v4l2 ./device` (which includes `TestABI` layout assertions), `go vet`, and `CGO_ENABLED=1 go test -race ./v4l2 ./device`.
2. **`integration-tests`:** Builds and loads `v4l2loopback` from source against the GitHub runner kernel and executes `test/...` with `-use-device-emulation=/dev/video10,/dev/video11 -require-streaming` under both `CGO_ENABLED=0` and `CGO_ENABLED=1 -race`.

## Troubleshooting

* **Loopback device busy or stale state:** Check open file handles with `fuser /dev/video*` or `lsof /dev/video*`, then reload cleanly:
  ```bash
  sudo modprobe -r v4l2loopback
  sudo modprobe v4l2loopback devices=2 video_nr=42,43 card_label="go4vl_test_1,go4vl_test_2" exclusive_caps=0
  go clean -testcache
  ```
* **`v4l2loopback` fails to build or load:** Confirm matching kernel headers exist for `$(uname -r)` (`sudo apt-get install linux-headers-$(uname -r)`).
* **Tests skip with no V4L2 device available:** Verify `ls -la /dev/video*` shows nodes accessible to your user, or pass `-require-streaming` alongside explicit `-use-device-emulation` paths.

## Performance Benchmarks

Run package-level allocator and streaming micro-benchmarks without hardware, or full capture benchmarks with `benchmark/runner`:

```bash
# Hot-path micro-benchmarks
CGO_ENABLED=0 go test -bench=. -benchmem ./v4l2 ./device

# End-to-end streaming benchmark on hardware or loopback
cd benchmark/runner
go run . -single -device /dev/video0 -api=iter -duration 10s
```

See [benchmark/README.md](./benchmark/README.md) and the [Benchmarks & Tuning Guide](https://vladimirvivien.github.io/go4vl/guides/benchmarks/) for pprof profiling options.
