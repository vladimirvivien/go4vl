# Building & Cross-Compiling

Because `go4vl` checks in shared Linux UAPI constants (`v4l2/zconsts_linux.go`) and architecture-specific UAPI definitions (`v4l2/ztypes_linux_{amd64,arm64,arm,riscv64}.go`) while invoking kernel ioctls through `golang.org/x/sys/unix`, standard builds require **no C compiler, no CGO, and no Linux kernel headers**.

## Prerequisites

* **Go 1.26 or later**
* **Target OS:** `GOOS=linux` with `GOARCH` set to `amd64`, `arm64`, `arm`, or `riscv64`
* **Optional ABI regeneration dependencies:** `gcc` or `clang` plus `linux-libc-dev` are needed only when regenerating UAPI definitions (`go generate ./v4l2`).

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

On 32-bit `linux/arm`, `go4vl`'s checked-in UAPI generator handles 32-bit struct padding, pointer sizes in `v4l2_buffer.m`, and `timeval` widths cleanly—no special build flags or host C sysroot are needed.

## Targeting additional Linux architectures

Building for a Linux `GOARCH` outside the four pre-generated targets (such as `386`, `ppc64le`, `s390x`, or `loong64`) requires generating a matching `v4l2/ztypes_linux_<goarch>.go` file from kernel UAPI headers in a fork (referenced via a `go.mod` `replace` directive) or contributing it upstream:

1. **Generate `ztypes_linux_<goarch>.go` from the repository root** with `gcc` or `clang` and target Linux UAPI headers installed (`linux-libc-dev`):

    ```bash
    # On a Linux host of the target architecture:
    go run ./v4l2/internal/gen -consts=false -arch=<goarch>

    # Or cross-generating from another Linux host with a target cross-compiler:
    CC=<target-triple>-gcc go run ./v4l2/internal/gen -consts=false -arch=<goarch>
    ```

    Passing `-consts=false` skips `v4l2/zconsts_linux.go` and regenerates only `v4l2/ztypes_linux_<goarch>.go`.

2. **Verify ABI layouts on the target architecture** so wrapper struct sizes, patched struct layouts, and union offsets match:

    ```bash
    CGO_ENABLED=0 go test -v ./v4l2 -run TestABI
    ```

3. **Review word-size and endianness assumptions** before shipping on a new target: inspect `v4l2/streaming.go` for 32-bit union packing (`is32Bit` / pointer-sized `m` fields), `v4l2/timing.go` for `timeval`/`timespec` field widths, and any big-endian byte-order assumptions if targeting a big-endian architecture such as `s390x`.

4. **Compile normally with CGO disabled**:

    ```bash
    CGO_ENABLED=0 GOOS=linux GOARCH=<goarch> go build ./...
    ```

## Multi-stage container builds

Static `CGO_ENABLED=0` binaries deploy into minimal `scratch` or distroless containers. Remember to pass the required `/dev/video*` or `/dev/media*` device node into the container at runtime (`docker run --device /dev/video0`):

```dockerfile
FROM golang:1.26 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /out/app ./examples/iter_capture

FROM scratch
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]
```

For a concrete example tailored to web video streaming, see `examples/simplecam/Dockerfile` in the repository.

## Running tests & verification

Unit tests mock the low-level `ioctl` seam and run on any Linux host without video hardware:

```bash
CGO_ENABLED=0 go vet ./...
CGO_ENABLED=0 go test -v ./v4l2 ./device
```

To run the integration test suite against real V4L2 hardware or a loaded `v4l2loopback` device, include the `integration` build tag:

```bash
CGO_ENABLED=0 go test -v -tags=integration ./test/... -args -use-device=/dev/video0
```

See the [Testing & CI Guide](guides/testing.md) for loopback setup and CI flags.

## Regenerating UAPI definitions (contributors)

When adding support for new Linux kernel structs, constants, or ioctls, add the C symbol entries to `v4l2/internal/gen/defs_consts.go` or `v4l2/internal/gen/defs_types.go`, then regenerate on a Linux machine with kernel UAPI headers installed (`sudo apt-get install build-essential linux-libc-dev`):

```bash
# Regenerate v4l2/zconsts_linux.go and host v4l2/ztypes_linux_<arch>.go
go generate ./v4l2

# Or cross-generate for another target architecture
CC=aarch64-linux-gnu-gcc go run ./v4l2/internal/gen -consts=false -arch=arm64

# Verify wrapper struct sizes, patched struct layouts, and union offsets
CGO_ENABLED=0 go test -v ./v4l2 -run TestABI
```

Run `CGO_ENABLED=0 go build ./...` across `amd64`, `arm64`, `arm`, and `riscv64` after regenerating so that 32-bit and 64-bit builds remain aligned.
