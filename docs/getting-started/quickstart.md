# Quick Start

Most applications follow a predictable four-step V4L2 lifecycle: open a `/dev/video*` node with format options, start buffer streaming, iterate over decoded or compressed frames, and cleanly stop and close on shutdown. This guide walks through capturing a single compressed frame to disk using `Device.Frames(ctx)`, then generalizes to continuous streaming and safe frame retention.

## Capturing a single snapshot

A minimal snapshot program configures a four-buffer MMAP ring with `MJPEG` output, starts streaming under a context, and exits the range loop after writing the first non-empty frame:

```go
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	devName := flag.String("d", "/dev/video0", "V4L2 video device path")
	flag.Parse()

	dev, err := device.Open(
		*devName,
		device.WithBufferSize(4),
		device.WithPixFormat(v4l2.PixFormat{
			Width:       640,
			Height:      480,
			PixelFormat: v4l2.PixelFmtMJPEG,
		}),
	)
	if err != nil {
		log.Fatalf("open %s: %v", *devName, err)
	}
	defer dev.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		log.Fatalf("start stream: %v", err)
	}
	defer dev.Stop()

	for frame, err := range dev.Frames(ctx) {
		if err != nil {
			log.Fatalf("capture frame: %v", err)
		}
		if err := os.WriteFile("snapshot.jpg", frame.Data, 0o644); err != nil {
			log.Fatalf("write file: %v", err)
		}
		break
	}
}
```

When you run the loop with `range dev.Frames(ctx)`, each step blocks in `v4l2.Poll` until a kernel DMA transfer finishes, dequeues the ready buffer index via `VIDIOC_DQBUF`, and yields a `*device.Frame` whose `Data` slice points straight into the memory-mapped kernel buffer. Breaking out of the loop runs the iterator cleanup hook, automatically re-queueing the buffer with `VIDIOC_QBUF` before `defer dev.Stop()` executes `VIDIOC_STREAMOFF`.

## Continuous capture and keeping frames across iterations

Because `frame.Data` aliases the kernel-mapped ring buffer directly inside `dev.Frames(ctx)`, holding that slice after the loop body finishes is unsafe: once the next iteration begins, the kernel driver may overwrite that memory with a new video frame.

Whenever you want to pass captured bytes to an encoder goroutine, HTTP worker, or image decoder asynchronously, call `frame.Clone()` inside the loop body to obtain an independent heap copy along with all frame metadata (`Sequence`, `Timestamp`, `Flags`, and `Planes`):

```go
for frame, err := range dev.Frames(ctx) {
	if err != nil {
		log.Printf("stream error: %v", err)
		break
	}
	if len(frame.Data) == 0 {
		continue
	}

	// Independent copy safe to pass to worker goroutines.
	owned := frame.Clone()
	workQueue <- owned
}
```

Look closely at the error check inside the loop: transient capture errors or context cancellation are reported cleanly through the second iterator value (`err`), giving you direct control over shutdown without background capture goroutines writing into an unconsumed channel.

## Running with a native or cross-compiled binary

Build directly on a Linux host or cross-compile from any development OS without a C toolchain:

```bash
# Native build on Linux
CGO_ENABLED=0 go build -o snapshot .
./snapshot -d /dev/video0

# Cross-compile for a 64-bit Raspberry Pi from macOS, Windows, or x86_64 Linux
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o snapshot .
```

Next, read [Streaming, I/O & Zero-Copy](../guides/streaming.md) to compare `Frames(ctx)` with the pooled channel API and learn when to use `FramePool`, `WriteFrame`, or `v4l2.Poller`.
