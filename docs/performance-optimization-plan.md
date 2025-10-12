# Frame Processing Pipeline Performance - Implementation Plan

**Roadmap Item:** #1 - Frame Processing Pipeline Performance
**Status:** In Progress
**Started:** 2025-10-11

## Overview

This document tracks the implementation plan for eliminating bottlenecks in the frame processing pipeline to ensure maximum throughput and minimal latency.

---

## Implementation Approach

### Phase 1: Profiling & Bottleneck Identification

#### Step 1: Analyze Current Implementation
- [ ] Read and analyze the current frame capture loop in `device/device.go:startStreamLoop()`
- [ ] Identify the frame copy operation at line ~676-677
- [ ] Map out complete data flow: device → mmap → copy → channel → user
- [ ] Document current memory allocation patterns
- [ ] Review channel buffer sizing and blocking behavior

**Key Code Locations:**
- Frame copy: `device/device.go:676-677`
- Channel send: `device/device.go:680`
- Buffer management: `v4l2/streaming.go:MapMemoryBuffers()`

#### Step 2: Create Profiling Infrastructure ✅ COMPLETED
- [x] Benchmark infrastructure in place (see `benchmark/README.md` for details)

#### Step 3: Identify Specific Bottlenecks
- [ ] Run baseline benchmarks on current implementation
- [ ] Profile CPU usage with pprof
- [ ] Profile memory allocation patterns
- [ ] Analyze execution trace for goroutine blocking
- [ ] Measure channel send/receive latency
- [ ] Document findings in performance report

**Expected Bottlenecks:**
1. Frame copy operation (line 676-677 allocates + copies every frame)
2. Channel send blocking when consumer is slow
3. Memory allocation for each frame (GC pressure)
4. Potential lock contention in multi-device scenarios

---

### Phase 2: Quick Wins Implementation

#### Option A: Zero-Copy Frame Access (Recommended)

**Goal:** Eliminate frame copy by passing buffer references directly to users

**Implementation Plan:**
- [ ] Add `WithZeroCopy()` option to device configuration
- [ ] Create `Frame` type with buffer reference and metadata
  ```go
  type Frame struct {
      Data      []byte  // Reference to mmap buffer (not copy)
      Index     uint32  // Buffer index
      Timestamp time.Time
      Sequence  uint32
      release   func()  // Callback to release buffer
  }
  ```
- [ ] Implement buffer lifecycle tracking
  - [ ] Reference counting per buffer
  - [ ] Release callback mechanism
  - [ ] Automatic re-queue after release
- [ ] Update channel type: `chan []byte` → `chan *Frame` (new API)
- [ ] Add buffer timeout mechanism (force release after N seconds)
- [ ] Update examples to demonstrate zero-copy usage
- [ ] Add safety checks (detect use-after-release)

**API Design:**
```go
// New zero-copy API
dev, _ := device.Open("/dev/video0",
    device.WithZeroCopy(),
    device.WithBufferSize(4),
)

for frame := range dev.GetFrames() { // Returns *Frame instead of []byte
    // Process frame data
    processFrame(frame.Data)

    // Must explicitly release when done
    frame.Release()
}
```

**Pros:**
- Maximum performance (no copy overhead)
- Minimal GC pressure
- Exposes frame metadata naturally

**Cons:**
- Breaking API change (need new method)
- Requires user discipline (must call Release())
- More complex buffer management

**Estimated Effort:** 2-3 hours

---

#### Option B: Frame Pool Implementation

**Goal:** Reuse allocated frame buffers to reduce GC pressure

**Implementation Plan:**
- [ ] Create `framePool` using `sync.Pool`
- [ ] Modify frame copy to use pooled buffers
- [ ] Return buffers to pool after channel receive
- [ ] Add pool metrics (hits, misses, size)
- [ ] Benchmark pool vs non-pool performance

**Implementation:**
```go
var framePool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0) // Capacity set on first use
    },
}

// In capture loop
func (d *Device) captureFrame(buff Buffer) {
    frame := framePool.Get().([]byte)
    if cap(frame) < int(buff.BytesUsed) {
        frame = make([]byte, buff.BytesUsed)
    }
    frame = frame[:buff.BytesUsed]
    copy(frame, d.buffers[buff.Index][:buff.BytesUsed])
    d.output <- frame
}

// User must return to pool (or use finalizer)
func ReturnFrame(frame []byte) {
    framePool.Put(frame[:0])
}
```

**Pros:**
- Compatible with existing API
- Reduces GC pressure
- Easy to implement

**Cons:**
- Still requires frame copy
- Requires user cooperation to return frames
- Pool management overhead

**Estimated Effort:** 1-2 hours

---

#### Option C: Configurable Frame Dropping

**Goal:** Handle slow consumers gracefully without blocking capture

**Implementation Plan:**
- [ ] Add `FrameDropPolicy` enum
  ```go
  type FrameDropPolicy int
  const (
      DropNever   FrameDropPolicy = iota // Block (current behavior)
      DropOldest                          // Discard oldest frame in channel
      DropNewest                          // Skip current frame
  )
  ```
- [ ] Add `WithFrameDropPolicy(policy)` option
- [ ] Implement non-blocking channel send with drop logic
- [ ] Add dropped frame counter metric
- [ ] Expose metrics via `Device.Stats()` method

**Implementation:**
```go
// In capture loop
select {
case d.output <- frame:
    // Frame sent successfully
case <-time.After(0): // Non-blocking
    switch d.config.dropPolicy {
    case DropOldest:
        <-d.output // Discard oldest
        d.output <- frame
        d.droppedFrames++
    case DropNewest:
        // Skip this frame
        d.droppedFrames++
    case DropNever:
        d.output <- frame // Block
    }
}
```

**Pros:**
- Simple to implement
- No API changes required
- Immediate benefit for high-throughput scenarios
- Easy to understand

**Cons:**
- Doesn't eliminate copy overhead
- Frame dropping may not be acceptable for all use cases
- Still has GC pressure

**Estimated Effort:** 1 hour

---

### Phase 3: Advanced Optimizations (Future)

- [ ] Implement lock-free ring buffer (replace channels)
- [ ] Add NUMA-aware buffer allocation
- [ ] Optimize memory layout for cache efficiency
- [ ] Implement batch frame processing
- [ ] Add GPU memory integration (CUDA, OpenCL)
- [ ] Explore io_uring for async I/O

---

## Session Plan

### Session 1: Profiling Infrastructure ✅ COMPLETED (2025-10-12)

**✅ Benchmark Infrastructure Complete**
- Benchmark infrastructure in place and documented
- Tests actual go4vl code paths (device, v4l2 packages)
- See `benchmark/README.md` for usage details

**Next Session Goals:**
- [ ] Run baseline measurements
- [ ] Identify primary bottleneck

**Priority 2: Quick Win Implementation (60-90 min)**
Choose one based on profiling results:
- [ ] Option A: Zero-copy mode (if copy is main bottleneck)
- [ ] Option B: Frame pool (if GC is main bottleneck)
- [ ] Option C: Frame dropping (if channel blocking is main bottleneck)

**Priority 3: Validation (30 min)**
- [ ] Run before/after benchmarks
- [ ] Document performance improvement
- [ ] Create example demonstrating feature
- [ ] Update this document with results

---

## Benchmark Results

### Baseline (Current Implementation)

**Test Configuration:**
- Resolution: TBD
- Frame Rate: TBD
- Duration: TBD
- Device: TBD

**Metrics:**
- FPS Actual: TBD
- CPU Usage: TBD
- Memory/Frame: TBD
- GC Pauses: TBD
- Frame Drops: TBD

### After Optimization

**Test Configuration:**
- Resolution: TBD
- Frame Rate: TBD
- Duration: TBD
- Device: TBD
- Optimization: TBD

**Metrics:**
- FPS Actual: TBD
- CPU Usage: TBD
- Memory/Frame: TBD
- GC Pauses: TBD
- Frame Drops: TBD
- **Improvement: TBD%**

---

## Decision Log

### 2025-10-11 - Initial Planning
- Created implementation plan
- Identified three potential quick wins
- Decided to start with profiling to guide optimization choice

---

## Next Steps

After completing initial optimization:

1. **Measure Impact**
   - Run comprehensive benchmarks
   - Compare against baseline
   - Document improvement percentage

2. **Iterate**
   - Identify remaining bottlenecks
   - Implement next optimization
   - Repeat profiling cycle

3. **Documentation**
   - Update API documentation
   - Create performance tuning guide
   - Add migration guide if API changes

4. **Testing**
   - Add performance regression tests
   - Test on multiple platforms (x86, ARM, ARM64)
   - Test with various devices and resolutions

---

## References

- Current implementation: `device/device.go:636-697`
- Buffer management: `v4l2/streaming.go:290-331`
- Roadmap: `ROADMAP.md` - Enhancement #1
