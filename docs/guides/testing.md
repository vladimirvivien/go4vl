# Testing & CI

Every code change in `go4vl` is backed by fast hardware-free unit tests, C ABI struct/ioctl parity verification against host Linux headers, race detection, and end-to-end streaming integration tests against `v4l2loopback` or physical cameras.

## Running unit tests (no camera required)

Packages `v4l2` and `device` provide an internal syscall mock seam (`ioctlSyscall`) so unit tests exercise error handling, buffer marshaling, frame pooling, MPLANE plane packing, and iterator cancellation without opening a real `/dev/video*` node:

```bash
CGO_ENABLED=0 go test -v -cover ./v4l2 ./device
CGO_ENABLED=0 go vet ./v4l2 ./device
```

To verify concurrent safety under Go's race detector:

```bash
CGO_ENABLED=1 go test -race ./v4l2 ./device
```

## Verifying Linux C ABI parity (`-tags=cgo_verify`)

Whenever UAPI structures or ioctl constants change across architectures, contributors verify binary layout compatibility (struct sizes, field offsets, constant values, and encoded `_IOWR`/`_IOR`/`_IOW` numbers) directly against `/usr/include/linux/videodev2.h`, `media.h`, and `v4l2-subdev.h`:

```bash
# Requires gcc and linux-libc-dev on a Linux host
CGO_ENABLED=1 go test -v -tags=cgo_verify ./v4l2 -run TestABI
```

Also cross-compile all packages across `amd64`, `arm64`, `arm`, and `riscv64` with `CGO_ENABLED=0` to verify 32-bit and 64-bit UAPI parity:

```bash
CGO_ENABLED=0 go build ./...
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ./...
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ./...
CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build ./...
```

## Running integration tests with `v4l2loopback` or hardware

Integration tests under `test/` are gated behind `-tags=integration` and exercise real V4L2 syscalls across `Frames(ctx)`, `GetFrames()`, `WriteFrame`, controls, formats, Media Controller topology, sub-devices, and codecs.

The CI workflow builds `v4l2loopback` cleanly against the running Linux kernel and runs streaming capture tests using an **in-process `Device.WriteFrame` producer** (requiring no external `ffmpeg` process to validate capture streams). To test locally with two loopback nodes:

```bash
sudo modprobe v4l2loopback devices=2 video_nr=42,43 card_label="go4vl_test_1,go4vl_test_2" exclusive_caps=0

CGO_ENABLED=0 go test -v -tags=integration ./test/... \
  -args -use-device-emulation=/dev/video42,/dev/video43 -require-streaming
```

Passing `-require-streaming` ensures that streaming capture/output assertions fail loudly rather than skipping when loopback nodes are unreachable. To test against attached hardware instead (or auto-discover mounted cameras):

```bash
CGO_ENABLED=0 go test -v -tags=integration ./test/... -args -use-device=/dev/video0
# Auto-discover loopback or real hardware:
CGO_ENABLED=0 go test -v -tags=integration ./test/... -args -use-device=auto
```

For full details on loopback troubleshooting, flags, and GitHub Actions CI jobs, see [TESTING_GUIDE.md](https://github.com/vladimirvivien/go4vl/blob/main/TESTING_GUIDE.md).
