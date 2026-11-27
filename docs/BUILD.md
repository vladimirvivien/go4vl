# Building & Cross-Compiling

Because `go4vl` checks in shared Linux UAPI constants (`v4l2/zconsts_linux.go`) and architecture-specific UAPI definitions (`v4l2/ztypes_linux_{amd64,arm64,arm,riscv64}.go`) while invoking kernel ioctls through `golang.org/x/sys/unix`, standard builds require **no C compiler, no CGO, and no Linux kernel headers**.

## Prerequisites

All application builds and cross-compilations rely solely on the Go toolchain, while UAPI header regeneration uses an optional C toolchain on Linux 6.12+ headers:

* **Go 1.26 or later**
* **Target OS:** `GOOS=linux` (**kernel 5.10 or later** at runtime) with `GOARCH` set to `amd64`, `arm64`, `arm`, or `riscv64`
* **Optional ABI regeneration dependencies:** `gcc` or `clang` plus **Linux 6.12+ `linux-libc-dev`** headers are needed only when regenerating UAPI definitions (`go generate ./v4l2`), or use Docker with `golang:1.26-trixie`

## Native Linux builds

Run a full native build with CGO disabled to verify every package and example compiles cleanly:

```bash
CGO_ENABLED=0 go build ./...
```

To compile a single binary directly from one of the examples:

```bash
CGO_ENABLED=0 go build -o iter_capture ./examples/iter_capture
```

## Cross-compiling for other Linux architectures

Standard Go cross-compilation works from Linux, macOS, or Windows development workstations. Set `CGO_ENABLED=0 GOOS=linux` alongside your target `GOARCH`:

```bash
# 64-bit ARM (Raspberry Pi 4/5 64-bit OS, Jetson, server ARM64)
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o iter_capture ./examples/iter_capture

# 32-bit ARM (Raspberry Pi Zero/3 32-bit OS, ARMv7 boards)
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o iter_capture ./examples/iter_capture

# 64-bit RISC-V Linux
CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -o iter_capture ./examples/iter_capture

# 64-bit x86 Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o iter_capture ./examples/iter_capture
```

On 32-bit `linux/arm`, the checked-in UAPI definitions target 64-bit `timeval`/`timespec` (`_TIME_BITS=64`, Linux 5.6+ `time64` ioctls) alongside 32-bit pointer packing in `v4l2_buffer.m` and `v4l2_plane.m`—no special build flags or host C sysroot are needed.

## Targeting additional Linux architectures

Building for a Linux `GOARCH` outside the four pre-generated targets (such as `386`, `ppc64le`, `s390x`, or `loong64`) requires generating a matching `v4l2/ztypes_linux_<goarch>.go` file from **Linux 6.12+** kernel UAPI headers in a fork (referenced via a `go.mod` `replace` directive) or contributing it upstream:

1. **Generate `ztypes_linux_<goarch>.go` from the repository root** with `gcc` or `clang` and Linux 6.12+ `linux-libc-dev` headers installed (or run inside the `golang:1.26-trixie` container shown under [Regenerating in Docker](#regenerating-in-docker-optional)):

    ```bash
    # On a Linux host of the target architecture:
    go run ./v4l2/internal/gen -consts=false -arch=<goarch>

    # Or cross-generating from another Linux host with a target cross-compiler:
    CC=<target-triple>-gcc go run ./v4l2/internal/gen -consts=false -arch=<goarch>
    ```

    Passing `-consts=false` skips `v4l2/zconsts_linux.go` and regenerates only `v4l2/ztypes_linux_<goarch>.go`.

2. **Verify ABI layouts and test compilation on the target architecture** so wrapper struct sizes, patched struct layouts, and union offsets match:

    ```bash
    CGO_ENABLED=0 GOOS=linux GOARCH=<goarch> go test -c ./v4l2 -o /dev/null
    CGO_ENABLED=0 go test -v ./v4l2 -run TestABI
    ```

3. **Review word-size and endianness assumptions** before shipping on a new target: inspect `v4l2/streaming.go` for 32-bit union packing (`is32Bit` / pointer-sized `m` fields), `v4l2/timing.go` for `timeval`/`timespec` field widths, and any big-endian byte-order assumptions if targeting a big-endian architecture such as `s390x`.

4. **Compile normally with CGO disabled**:

    ```bash
    CGO_ENABLED=0 GOOS=linux GOARCH=<goarch> go build ./...
    ```

## Multi-stage container builds

Static `CGO_ENABLED=0` binaries cross-compile natively on the build host via Docker BuildKit's `--platform=$BUILDPLATFORM` (avoiding slow QEMU emulation during `go build`) and deploy into minimal `scratch` or distroless images.

```dockerfile
# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26 AS build
ARG TARGETOS=linux
ARG TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/app ./examples/iter_capture

FROM scratch
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]
```

Build single- or multi-architecture images with BuildKit and pass the required `/dev/video*` node (plus `/dev/media*` and `/dev/v4l-subdev*` nodes for CSI-2/ISP pipelines) at runtime:

```bash
# Multi-architecture image build on the host CPU
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t mycam:latest .

# Run with device node access
docker run --rm --device /dev/video0 mycam:latest
```

To extract standalone cross-compiled binaries directly from a container build onto your host (such as with `--output ./build`), see `examples/simplecam/Dockerfile` in the repository.

## Running tests & verification

Unit tests mock the low-level `ioctl` seam and run on any Linux host without video hardware:

```bash
CGO_ENABLED=0 go vet ./...
CGO_ENABLED=0 go test -v ./v4l2 ./device ./v4l2/internal/gen
```

To run the integration test suite against real V4L2 hardware or a loaded `v4l2loopback` device, include the `integration` build tag:

```bash
CGO_ENABLED=0 go test -v -tags=integration ./test/... -args -use-device=/dev/video0
```

See the [Testing & CI Guide](guides/testing.md) for loopback setup and CI flags.

## Regenerating UAPI definitions (contributors)

The generator at `v4l2/internal/gen` invokes `go tool cgo -godefs` directly on `v4l2/internal/gen/defs_consts.go` and `v4l2/internal/gen/defs_types.go`, applies layout rewrites for anonymous C union fields and packed `v4l2_bt_timings.Pixelclock`, and writes standalone `v4l2/zconsts_linux.go` and `v4l2/ztypes_linux_<arch>.go` files.

Maintain these invariants when modifying UAPI definitions:

* **Linux 6.12+ header baseline:** Contributors regenerate against Linux 6.12+ `linux-libc-dev` headers so upstream UAPI enums and newer ioctls resolve directly without C preprocessor fallbacks. Runtime compatibility remains Linux 5.10+ because Linux V4L2 UAPI is strictly additive—kernels without a newer ioctl return `ENOTTY`.
* **Deterministic 64-bit time ABI:** `defs_types.go` defines `_FILE_OFFSET_BITS=64` and `_TIME_BITS=64` so 32-bit architectures (`linux/arm`) emit `time64` struct layouts and ioctl numbers consistently across distros.
* **Non-mainline constants in Go:** Vendor-specific FourCC codes not in mainline Linux headers (such as Raspberry Pi's `MetaFmtBCM2835ISPStats`) are declared directly in Go (for example in `v4l2/format_meta.go`) rather than in `defs_consts.go`.

Configure `v4l2/internal/gen` using the flags and compiler override below:

| Flag / Env | Default | Purpose |
|---|---|---|
| `-arch` | Host `go env GOARCH` | Target Linux `GOARCH` for `ztypes_linux_<arch>.go` |
| `-consts` | `true` | Generate shared `v4l2/zconsts_linux.go` |
| `-types` | `true` | Generate `v4l2/ztypes_linux_<arch>.go` |
| `-out` | `v4l2/` relative to source | Output directory for generated `z*.go` files |
| `CC` | Host `gcc` or `clang` | C compiler binary (required when cross-generating `-arch != host`) |

Add new C symbols to `v4l2/internal/gen/defs_consts.go` or `v4l2/internal/gen/defs_types.go`, then regenerate and verify on a Linux 6.12+ host:

```bash
# Regenerate v4l2/zconsts_linux.go and host v4l2/ztypes_linux_<arch>.go
go generate ./v4l2

# Or cross-generate for another target architecture
CC=aarch64-linux-gnu-gcc go run ./v4l2/internal/gen -consts=false -arch=arm64

# Run generator unit tests and ABI layout assertions
go test -v ./v4l2/internal/gen
CGO_ENABLED=0 go test -v ./v4l2 -run TestABI
```

Run `CGO_ENABLED=0 go build ./...` across `amd64`, `arm64`, `arm`, and `riscv64` after regenerating so that all supported targets compile cleanly.

### Regenerating in Docker (optional)

Docker is optional and not required for application builds or for native regeneration on a Linux 6.12+ workstation. Use the `golang:1.26-trixie` container image (Debian 13, shipping Linux 6.12 LTS `linux-libc-dev` headers) when regenerating from macOS, Windows, an older Linux distro, or when regenerating all four architectures together without installing host cross-compilers:

```bash
# Regenerate host-architecture UAPI files in a pinned Linux 6.12 LTS container
docker run --rm -v "$PWD":/src -w /src golang:1.26-trixie go generate ./v4l2

# Regenerate zconsts_linux.go and all four architecture files in one container pass
docker run --rm -v "$PWD":/src -w /src golang:1.26-trixie sh -c '
  apt-get update && apt-get install -y --no-install-recommends \
    gcc-aarch64-linux-gnu gcc-arm-linux-gnueabihf gcc-riscv64-linux-gnu && \
  go generate ./v4l2 && \
  CC=aarch64-linux-gnu-gcc go run ./v4l2/internal/gen -consts=false -arch=arm64 && \
  CC=arm-linux-gnueabihf-gcc go run ./v4l2/internal/gen -consts=false -arch=arm && \
  CC=riscv64-linux-gnu-gcc go run ./v4l2/internal/gen -consts=false -arch=riscv64
'
```
