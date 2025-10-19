package device

import (
	"sync"
	"sync/atomic"
)

// FramePool manages a pool of reusable byte buffers for video frames.
// It reduces memory allocations and GC pressure by reusing buffers across frames.
//
// FramePool is safe for concurrent use by multiple goroutines.
//
// # Usage
//
// The pool is typically used internally by Device, but can be created manually:
//
//	pool := NewFramePool(1024 * 1024)  // 1MB initial capacity
//	buf := pool.Get(640 * 480 * 2)     // Get buffer for YUYV frame
//	// ... use buffer ...
//	pool.Put(buf)                       // Return to pool
//
// # Performance Characteristics
//
// - Get() operations reuse existing buffers when available
// - Buffers grow automatically to accommodate larger frames
// - Put() returns buffers to the pool for future reuse
// - Thread-safe with minimal lock contention
type FramePool struct {
	pool sync.Pool

	// defaultCap is the default capacity for newly allocated buffers
	defaultCap int

	// stats track pool usage for monitoring and debugging
	gets    atomic.Int64
	puts    atomic.Int64
	allocs  atomic.Int64
	resizes atomic.Int64
}

// NewFramePool creates a new FramePool with the specified default capacity.
// The defaultCapacity is used for newly allocated buffers and should be set
// to accommodate typical frame sizes for your use case.
//
// Common default capacities:
//   - 640x480 YUYV: ~614 KB
//   - 1280x720 YUYV: ~1.8 MB
//   - 1920x1080 YUYV: ~4 MB
//   - MJPEG/H264: Variable, typically 100KB - 1MB
//
// Example:
//
//	// For 1080p YUYV video
//	pool := NewFramePool(1920 * 1080 * 2)
func NewFramePool(defaultCapacity int) *FramePool {
	fp := &FramePool{
		defaultCap: defaultCapacity,
	}

	fp.pool.New = func() any {
		// Pre-allocate buffer with default capacity
		buf := make([]byte, 0, fp.defaultCap)
		fp.allocs.Add(1)
		return &buf
	}

	return fp
}

// Get retrieves a buffer from the pool with at least the specified size.
// The returned buffer's length is set to size, and its capacity may be larger.
//
// If the pooled buffer's capacity is insufficient, it will be resized.
// The returned buffer should be passed to Put() when no longer needed.
//
// Example:
//
//	buf := pool.Get(307200)  // Get 300KB buffer
//	copy(buf, frameData)
//	pool.Put(buf)
func (fp *FramePool) Get(size uint32) []byte {
	fp.gets.Add(1)

	bufPtr := fp.pool.Get().(*[]byte)

	// Check if buffer needs resizing
	if cap(*bufPtr) < int(size) {
		fp.resizes.Add(1)

		// Allocate new buffer with extra capacity (2x requested size)
		// This reduces future resize operations
		newCap := int(size) * 2
		if newCap < fp.defaultCap {
			newCap = fp.defaultCap
		}
		*bufPtr = make([]byte, size, newCap)
	} else {
		// Reuse existing buffer, just adjust length
		*bufPtr = (*bufPtr)[:size]
	}

	return *bufPtr
}

// Put returns a buffer to the pool for reuse.
// The buffer should not be used after calling Put().
//
// Calling Put() with a nil or empty slice is safe (it's a no-op).
// The buffer's capacity is preserved, but its length is reset to 0.
//
// Example:
//
//	buf := pool.Get(1024)
//	processData(buf)
//	pool.Put(buf)  // Return to pool
//	// buf is now invalid, don't use it
func (fp *FramePool) Put(buf []byte) {
	if buf == nil || cap(buf) == 0 {
		return
	}

	fp.puts.Add(1)

	// Reset length to 0, preserving capacity
	buf = buf[:0]

	fp.pool.Put(&buf)
}

// Stats returns current pool statistics for monitoring and debugging.
// All statistics are cumulative since pool creation.
type PoolStats struct {
	// Gets is the total number of Get() calls
	Gets int64

	// Puts is the total number of Put() calls
	Puts int64

	// Allocs is the total number of new buffer allocations
	Allocs int64

	// Resizes is the total number of buffer resize operations
	Resizes int64

	// Outstanding is the number of buffers currently in use (Gets - Puts)
	Outstanding int64

	// HitRate is the percentage of Get() calls that reused buffers (0.0 - 1.0)
	HitRate float64
}

// Stats returns current pool usage statistics.
// This is useful for monitoring pool efficiency and tuning default capacity.
//
// Example:
//
//	stats := pool.Stats()
//	fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate * 100)
//	fmt.Printf("Outstanding buffers: %d\n", stats.Outstanding)
func (fp *FramePool) Stats() PoolStats {
	gets := fp.gets.Load()
	puts := fp.puts.Load()
	allocs := fp.allocs.Load()
	resizes := fp.resizes.Load()

	var hitRate float64
	if gets > 0 {
		// Hit rate = (gets - allocs) / gets
		hits := gets - allocs
		if hits < 0 {
			hits = 0
		}
		hitRate = float64(hits) / float64(gets)
	}

	return PoolStats{
		Gets:        gets,
		Puts:        puts,
		Allocs:      allocs,
		Resizes:     resizes,
		Outstanding: gets - puts,
		HitRate:     hitRate,
	}
}

// Reset clears all pool statistics.
// This does not affect pooled buffers, only the statistics counters.
func (fp *FramePool) Reset() {
	fp.gets.Store(0)
	fp.puts.Store(0)
	fp.allocs.Store(0)
	fp.resizes.Store(0)
}

// defaultFramePool is the global default pool used by Device instances.
// It's initialized with a 1MB default capacity suitable for most video formats.
var defaultFramePool = NewFramePool(1024 * 1024)

// DefaultFramePool returns the global default frame pool.
// This pool is shared across all Device instances unless a custom pool is specified.
func DefaultFramePool() *FramePool {
	return defaultFramePool
}
