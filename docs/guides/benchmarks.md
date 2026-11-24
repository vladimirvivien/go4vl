# Benchmarks & Performance Tuning

Real-time video capture leaves tight frame budgets—at 60 frames per second, the kernel and userspace pipeline have just 16.6 ms per frame. `go4vl` includes both micro-benchmarks for allocator and MMAP hot paths and an end-to-end streaming benchmark suite (`benchmark/`) that measures throughput, allocation rates, and CPU profiles on real V4L2 hardware or `v4l2loopback`.

## Micro-benchmarks (`./v4l2` and `./device`)

Measure buffer pool checkout/return overhead, slice copying, and V4L2 streaming ring helpers without opening hardware devices:

```bash
CGO_ENABLED=0 go test -bench=. -benchmem ./device ./v4l2
```

Use `-benchtime=3s` and `-cpuprofile=cpu.prof` / `-memprofile=mem.prof` when tuning zero-allocation code paths and inspect hotspots with `go tool pprof`.

## End-to-end V4L2 benchmark runner

The `benchmark/runner` CLI measures real syscall latency, frame rate stability, dropped frames, and GC impact on physical webcams or loopback virtual cameras using either the zero-copy `Frames(ctx)` iterator (`-api=iter`, default) or the pooled channel pipeline (`-api=frames`).

Run a single benchmark directly against a hardware camera:

```bash
cd benchmark/runner
go run . -single -device /dev/video0 -api=iter -duration 10s
```

Compare against the pooled channel capture path on the same device:

```bash
go run . -single -device /dev/video0 -api=frames -duration 10s
```

When no physical camera is attached, omit `-device` and run with `sudo` (requires `v4l2loopback` and `ffmpeg`) to drive synthetic `MJPEG` or `YUYV` feeds and capture optional `pprof` CPU, memory, and execution traces:

```bash
sudo go run . -single -duration 15s -cpuprofile cpu.prof -memprofile mem.prof
```

## Tuning recommendations

* **Prefer `Device.Frames(ctx)` on single-consumer pipelines:** Because `Frames(ctx)` yields slices backed directly by kernel MMAP memory and re-queues buffers in-place, per-frame heap allocations drop to zero. Only call `frame.Clone()` when a frame must outlive the loop iteration.
* **Right-size your MMAP ring (`WithBufferSize`):** Two buffers (`2`) minimize latency when your processing loop keeps pace with frame intervals; `4`–`6` buffers absorb scheduling jitter or disk/network spikes without dropping frames.
* **Recycle pooled frames reliably (`GetFrames()`):** When reading from `Device.GetFrames()`, always invoke `frame.Release()` promptly so buffers return to `FramePool` rather than triggering fresh multi-hundred-kilobyte heap allocations.
* **Choose hardware-friendly pixel formats:** Requesting compressed `PixelFmtMJPEG` or native planar YUV (`NV12` / `YUYV`) avoids hidden driver or CPU format conversion overhead on USB 2.0 cameras and embedded ISPs.

See [benchmark/README.md](https://github.com/vladimirvivien/go4vl/blob/main/benchmark/README.md) for the complete scenario matrix (`-list`, `-scenario`, and trace analysis workflows).
