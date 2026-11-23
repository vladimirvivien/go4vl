# Building `go4vl`

`go4vl` uses generated Linux UAPI definitions (`v4l2/zconsts_linux.go`, `v4l2/ztypes_linux_amd64.go`, `v4l2/ztypes_linux_arm64.go`) and `golang.org/x/sys/unix` syscalls. Normal builds and cross-compiles require **no C compiler, no CGO, and no Linux kernel headers**.

## Prerequisites

- **Go 1.26+**
- **Linux target** (`GOOS=linux`, `GOARCH=amd64` or `GOARCH=arm64`)
- *(Optional)* A C compiler (`gcc` or `clang`) and `linux-libc-dev` are only needed when regenerating UAPI definitions (`go generate ./v4l2`) or running `-tags=cgo_verify` ABI parity tests.

## Native Build

```bash
CGO_ENABLED=0 go build ./...
```

Build a specific example:

```bash
CGO_ENABLED=0 go build -o iter_capture ./examples/iter_capture
```

## Cross-Compilation

Because runtime CGO is not used, standard Go cross-compilation works from Linux, macOS, or Windows hosts:

```bash
# 64-bit ARM (e.g., Raspberry Pi 4/5 64-bit OS)
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o iter_capture ./examples/iter_capture

# 64-bit x86 Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o iter_capture ./examples/iter_capture
```

## Docker Build

Multi-stage static builds work with the official Go image and `scratch`:

```dockerfile
FROM golang:1.26 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /out/app ./examples/iter_capture

FROM scratch
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]
```

## Testing & Verification

```bash
# Unit tests and vet (no CGO required)
CGO_ENABLED=0 go test -v ./v4l2 ./device
CGO_ENABLED=0 go vet ./v4l2 ./device

# Integration tests against real hardware or v4l2loopback
CGO_ENABLED=0 go test -v -tags=integration ./test/... -args -use-device=/dev/video0
```

## Regenerating UAPI Definitions & ABI Parity

Contributors adding new kernel structs or ioctls regenerate the checked-in ABI files on a Linux host with kernel headers installed (`sudo apt install build-essential linux-libc-dev`):

```bash
# Regenerate zconsts_linux.go and ztypes_linux_<arch>.go
go generate ./v4l2

# Verify Go struct sizes, field offsets, ioctl numbers, and constants against host C headers
CGO_ENABLED=1 go test -v -tags=cgo_verify ./v4l2 -run TestABI
```
