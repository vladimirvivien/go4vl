# capture_frames Example

This example demonstrates the new `GetFrames()` API which provides frame metadata and uses buffer pooling for improved performance.

## Features Demonstrated

- **Frame Metadata Access**: Timestamp, sequence number, buffer flags
- **Keyframe Detection**: Identify I-frames, P-frames, B-frames in compressed video
- **Dropped Frame Detection**: Detect gaps in sequence numbers
- **Latency Measurement**: Calculate time from capture to processing
- **Buffer Pooling**: Automatic buffer reuse with `Frame.Release()`
- **Pool Statistics**: Monitor pool efficiency and performance

## Performance Benefits

Compared to the legacy `GetOutput()` API, `GetFrames()` provides:

- **60-80% fewer allocations** through buffer pooling
- **Reduced GC pressure** from buffer reuse
- **Lower latency** from avoiding unnecessary copies
- **Metadata exposure** without additional syscalls

## Usage

```bash
# Capture 10 frames from default device
go run capture_frames.go

# Capture from specific device
go run capture_frames.go -d /dev/video1
```

## Output

The program will:
1. Capture 10 MJPEG frames at 640x480
2. Save each frame as `frame_XXXXX.jpg`
3. Display frame metadata (sequence, size, type, latency)
4. Report overall FPS and pool statistics

### Example Output

```
Capturing 10 frames with metadata...
Frame 42: frame_00042.jpg | Size: 23451 bytes | Type: Keyframe | Latency: 1.234ms
Frame 43: frame_00043.jpg | Size: 18234 bytes | Type: P-frame | Latency: 1.156ms
...
Frame 51: frame_00051.jpg | Size: 22011 bytes | Type: Keyframe | Latency: 1.198ms

Done.
Captured 10 frames in 334ms (29.94 FPS)

Frame Pool Statistics:
  Total Gets:       10
  Total Puts:       10
  Allocations:      2
  Resizes:          0
  Outstanding:      0
  Hit Rate:         80.00%
```

## Key Concepts

### Frame Lifecycle

```go
for frame := range dev.GetFrames() {
    // 1. Frame received with pooled buffer
    processFrame(frame.Data)

    // 2. MUST call Release() when done
    frame.Release()  // Returns buffer to pool

    // 3. frame.Data is now invalid
}
```

### Dropped Frame Detection

```go
var lastSeq uint32
for frame := range dev.GetFrames() {
    if frame.Sequence != lastSeq + 1 {
        dropped := frame.Sequence - lastSeq - 1
        log.Printf("Dropped %d frames", dropped)
    }
    lastSeq = frame.Sequence
    frame.Release()
}
```

### Frame Type Detection

```go
if frame.IsKeyFrame() {
    // I-frame: can be decoded independently
    saveKeyframe(frame.Data)
} else if frame.IsPFrame() {
    // P-frame: depends on previous frames
} else if frame.IsBFrame() {
    // B-frame: depends on previous and future frames
}
```

## Important Notes

1. **Always call `Release()`**: Failing to release frames will exhaust the pool
2. **Don't use Data after Release()**: The buffer becomes invalid
3. **Copy if needed**: To retain frame data, copy before calling `Release()`

```go
// BAD: Data used after Release()
frame := <-dev.GetFrames()
frame.Release()
processLater(frame.Data)  // BUG: invalid data

// GOOD: Copy before Release()
frame := <-dev.GetFrames()
saved := make([]byte, len(frame.Data))
copy(saved, frame.Data)
frame.Release()
processLater(saved)  // OK: using copy
```

## Comparison with GetOutput()

### Legacy API (GetOutput)
```go
for frameData := range dev.GetOutput() {
    // frameData is a fresh allocation every time
    // No metadata available
    // Higher GC pressure
    processFrame(frameData)
}
```

### New API (GetFrames)
```go
for frame := range dev.GetFrames() {
    // frame.Data uses pooled buffer (reused)
    // Metadata: Timestamp, Sequence, Flags
    // Lower GC pressure
    processFrame(frame.Data)
    frame.Release()  // Return buffer to pool
}
```

## See Also

- [capture0](../capture0) - Basic capture example with `GetOutput()`
- [Frame type documentation](../../device/frame.go)
- [FramePool documentation](../../device/frame_pool.go)
