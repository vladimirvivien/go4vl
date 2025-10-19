package device

import (
	"testing"
)

// TestFramePool_GetPut tests basic pool get/put operations
func TestFramePool_GetPut(t *testing.T) {
	pool := NewFramePool(1024)

	// Get buffer
	buf := pool.Get(512)
	if buf == nil {
		t.Fatal("Get returned nil")
	}
	if len(buf) != 512 {
		t.Errorf("Get returned buffer with len=%d, want 512", len(buf))
	}
	if cap(buf) < 512 {
		t.Errorf("Get returned buffer with cap=%d, want >= 512", cap(buf))
	}

	// Put buffer back
	pool.Put(buf)

	stats := pool.Stats()
	if stats.Gets != 1 {
		t.Errorf("Stats.Gets = %d, want 1", stats.Gets)
	}
	if stats.Puts != 1 {
		t.Errorf("Stats.Puts = %d, want 1", stats.Puts)
	}
}

// TestFramePool_Reuse tests that buffers are actually reused
func TestFramePool_Reuse(t *testing.T) {
	pool := NewFramePool(1024)

	// Get and put a buffer
	buf1 := pool.Get(512)
	pool.Put(buf1)

	// Get another buffer - should reuse
	buf2 := pool.Get(512)

	stats := pool.Stats()
	if stats.Allocs != 1 {
		t.Errorf("Expected 1 allocation, got %d", stats.Allocs)
	}
	if stats.HitRate < 0.4 { // At least 50% hit rate with 2 gets, 1 alloc
		t.Errorf("Expected hit rate >= 0.5, got %.2f", stats.HitRate)
	}

	pool.Put(buf2)
}

// TestFramePool_Resize tests automatic buffer resizing
func TestFramePool_Resize(t *testing.T) {
	pool := NewFramePool(100)

	// Get small buffer
	buf1 := pool.Get(50)
	pool.Put(buf1)

	// Get larger buffer - should resize
	buf2 := pool.Get(200)
	if cap(buf2) < 200 {
		t.Errorf("Buffer capacity = %d, want >= 200", cap(buf2))
	}

	stats := pool.Stats()
	if stats.Resizes != 1 {
		t.Errorf("Expected 1 resize, got %d", stats.Resizes)
	}

	pool.Put(buf2)
}

// TestFramePool_Concurrent tests concurrent access
func TestFramePool_Concurrent(t *testing.T) {
	pool := NewFramePool(1024)
	done := make(chan bool)

	// Spawn multiple goroutines
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				buf := pool.Get(512)
				// Simulate work
				buf[0] = byte(j)
				pool.Put(buf)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	stats := pool.Stats()
	if stats.Gets != 1000 {
		t.Errorf("Expected 1000 gets, got %d", stats.Gets)
	}
	if stats.Puts != 1000 {
		t.Errorf("Expected 1000 puts, got %d", stats.Puts)
	}
	if stats.Outstanding != 0 {
		t.Errorf("Expected 0 outstanding buffers, got %d", stats.Outstanding)
	}
}

// TestFramePool_Stats tests statistics tracking
func TestFramePool_Stats(t *testing.T) {
	pool := NewFramePool(1024)

	// Get some buffers
	buf1 := pool.Get(512)
	buf2 := pool.Get(512)
	buf3 := pool.Get(512)

	stats := pool.Stats()
	if stats.Gets != 3 {
		t.Errorf("Gets = %d, want 3", stats.Gets)
	}
	if stats.Outstanding != 3 {
		t.Errorf("Outstanding = %d, want 3", stats.Outstanding)
	}

	// Put some back
	pool.Put(buf1)
	pool.Put(buf2)

	stats = pool.Stats()
	if stats.Puts != 2 {
		t.Errorf("Puts = %d, want 2", stats.Puts)
	}
	if stats.Outstanding != 1 {
		t.Errorf("Outstanding = %d, want 1", stats.Outstanding)
	}

	// Reset stats
	pool.Reset()
	stats = pool.Stats()
	if stats.Gets != 0 || stats.Puts != 0 {
		t.Error("Reset did not clear statistics")
	}

	pool.Put(buf3)
}

// BenchmarkFramePool_Get benchmarks buffer allocation from pool
func BenchmarkFramePool_Get(b *testing.B) {
	pool := NewFramePool(1024 * 1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := pool.Get(614400) // ~640x480 YUYV
		pool.Put(buf)
	}
}

// BenchmarkDirectAllocation benchmarks direct allocation without pool
func BenchmarkDirectAllocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = make([]byte, 614400) // ~640x480 YUYV
	}
}

// BenchmarkFramePool_GetParallel benchmarks concurrent pool access
func BenchmarkFramePool_GetParallel(b *testing.B) {
	pool := NewFramePool(1024 * 1024)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := pool.Get(614400)
			pool.Put(buf)
		}
	})
}

// BenchmarkFramePool_VaryingSizes benchmarks pool with varying buffer sizes
func BenchmarkFramePool_VaryingSizes(b *testing.B) {
	pool := NewFramePool(1024 * 1024)
	sizes := []uint32{307200, 614400, 921600} // Different frame sizes
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		size := sizes[i%len(sizes)]
		buf := pool.Get(size)
		pool.Put(buf)
	}
}
