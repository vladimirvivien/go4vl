package device

import (
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// Frame represents a captured video frame with metadata and lifecycle management.
// It contains the frame data along with timing information, sequence numbers,
// and buffer flags that provide context about the frame.
//
// Frame uses an internal buffer pool to reduce memory allocations and GC pressure.
// Users should call Release() when done processing to return the buffer to the pool.
//
// # Usage Pattern
//
// Basic usage with automatic cleanup:
//
//	for frame := range dev.GetFrames() {
//	    processFrame(frame.Data)
//	    frame.Release()  // Return buffer to pool
//	}
//
// Saving frame data for later use:
//
//	var saved []byte
//	for frame := range dev.GetFrames() {
//	    saved = make([]byte, len(frame.Data))
//	    copy(saved, frame.Data)  // Must copy before Release()
//	    frame.Release()
//	}
//
// # Memory Management
//
// The Frame.Data slice is backed by a pooled buffer. After calling Release(),
// the Data slice becomes invalid and should not be accessed. If you need to
// retain frame data beyond the Release() call, make a copy first.
//
// # Thread Safety
//
// Frame is NOT safe for concurrent access. Each frame should be processed
// by a single goroutine. The Release() method is safe to call multiple times
// (subsequent calls are no-ops).
type Frame struct {
	// Data contains the raw frame data in the device's pixel format.
	// The format can be determined by calling device.GetPixFormat().
	// This slice is valid only until Release() is called.
	Data []byte

	// Timestamp is the capture time of this frame.
	// The timestamp source depends on the driver configuration.
	Timestamp time.Time

	// Sequence is a monotonically increasing frame counter.
	// Gaps in the sequence number indicate dropped frames.
	Sequence uint32

	// Flags contains buffer status flags (see v4l2.BufFlag* constants).
	// Common flags:
	//   - BufFlagKeyFrame: Frame is a keyframe (I-frame)
	//   - BufFlagPFrame: Frame is a predicted frame
	//   - BufFlagBFrame: Frame is a bidirectional predicted frame
	//   - BufFlagError: Frame contains errors
	Flags uint32

	// Index is the internal buffer index used by the driver.
	// This is primarily for debugging and internal use.
	Index uint32

	// pool is the FramePool this frame's buffer came from
	pool *FramePool

	// released tracks whether Release() has been called
	released bool
}

// Release returns the frame's buffer to the pool for reuse.
// This method should be called as soon as processing is complete to reduce
// memory pressure and improve performance.
//
// After calling Release(), the Data slice becomes invalid and must not be accessed.
// Calling Release() multiple times is safe (subsequent calls are no-ops).
//
// Example:
//
//	frame := <-dev.GetFrames()
//	processFrame(frame.Data)
//	frame.Release()  // Always call when done
func (f *Frame) Release() {
	if f.released || f.pool == nil || f.Data == nil {
		return
	}
	f.released = true

	// Return buffer to pool
	f.pool.Put(f.Data)
	f.Data = nil
}

// IsKeyFrame returns true if this frame is a keyframe (I-frame).
// This is useful for identifying frames that can be decoded independently
// in compressed video streams (H.264, MJPEG, etc.).
func (f *Frame) IsKeyFrame() bool {
	return f.Flags&v4l2.BufFlagKeyFrame != 0
}

// IsPFrame returns true if this frame is a predicted frame (P-frame).
// P-frames depend on previous frames for decoding.
func (f *Frame) IsPFrame() bool {
	return f.Flags&v4l2.BufFlagPFrame != 0
}

// IsBFrame returns true if this frame is a bidirectional predicted frame (B-frame).
// B-frames depend on both previous and future frames for decoding.
func (f *Frame) IsBFrame() bool {
	return f.Flags&v4l2.BufFlagBFrame != 0
}

// HasError returns true if the driver reported an error for this frame.
// Error frames may have incomplete or corrupted data.
func (f *Frame) HasError() bool {
	return f.Flags&v4l2.BufFlagError != 0
}
