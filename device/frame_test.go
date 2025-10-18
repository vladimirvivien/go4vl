package device

import (
	"testing"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestFrame_Release tests the basic Release() functionality
func TestFrame_Release(t *testing.T) {
	pool := NewFramePool(1024)
	data := pool.Get(512)

	frame := &Frame{
		Data:     data,
		pool:     pool,
		released: false,
	}

	// Verify initial state
	if frame.Data == nil {
		t.Error("Frame.Data should not be nil before Release()")
	}
	if frame.released {
		t.Error("Frame should not be marked as released initially")
	}

	// Release the frame
	frame.Release()

	// Verify post-release state
	if frame.Data != nil {
		t.Error("Frame.Data should be nil after Release()")
	}
	if !frame.released {
		t.Error("Frame should be marked as released after Release()")
	}

	// Verify buffer was returned to pool
	stats := pool.Stats()
	if stats.Puts != 1 {
		t.Errorf("Pool.Puts = %d, want 1", stats.Puts)
	}
}

// TestFrame_ReleaseMultipleTimes tests that multiple Release() calls are safe
func TestFrame_ReleaseMultipleTimes(t *testing.T) {
	pool := NewFramePool(1024)
	data := pool.Get(512)

	frame := &Frame{
		Data:     data,
		pool:     pool,
		released: false,
	}

	// Release multiple times
	frame.Release()
	frame.Release()
	frame.Release()

	// Should only put back to pool once
	stats := pool.Stats()
	if stats.Puts != 1 {
		t.Errorf("Pool.Puts = %d, want 1 (multiple Release() should be no-op)", stats.Puts)
	}
}

// TestFrame_ReleaseNilPool tests Release() with nil pool (should be no-op)
func TestFrame_ReleaseNilPool(t *testing.T) {
	frame := &Frame{
		Data:     make([]byte, 512),
		pool:     nil, // No pool
		released: false,
	}

	// Should not panic
	frame.Release()

	// With nil pool, Release() is a no-op (returns early)
	// The frame is not marked as released since there's nothing to release
	if frame.released {
		t.Error("Frame should not be marked as released with nil pool (no-op)")
	}

	// Data should still be accessible since it wasn't released
	if frame.Data == nil {
		t.Error("Data should not be nil when pool is nil")
	}
}

// TestFrame_ReleaseNilData tests Release() with nil data
func TestFrame_ReleaseNilData(t *testing.T) {
	pool := NewFramePool(1024)

	frame := &Frame{
		Data:     nil, // No data
		pool:     pool,
		released: false,
	}

	// Should not panic
	frame.Release()

	// Should not put anything back to pool
	stats := pool.Stats()
	if stats.Puts != 0 {
		t.Errorf("Pool.Puts = %d, want 0 (nil data shouldn't be returned to pool)", stats.Puts)
	}
}

// TestFrame_IsKeyFrame tests the IsKeyFrame() method
func TestFrame_IsKeyFrame(t *testing.T) {
	tests := []struct {
		name     string
		flags    uint32
		expected bool
	}{
		{"KeyFrame flag set", v4l2.BufFlagKeyFrame, true},
		{"KeyFrame with other flags", v4l2.BufFlagKeyFrame | v4l2.BufFlagMapped, true},
		{"No KeyFrame flag", v4l2.BufFlagMapped, false},
		{"Zero flags", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := &Frame{Flags: tt.flags}
			if result := frame.IsKeyFrame(); result != tt.expected {
				t.Errorf("IsKeyFrame() = %v, want %v (flags=0x%x)", result, tt.expected, tt.flags)
			}
		})
	}
}

// TestFrame_IsPFrame tests the IsPFrame() method
func TestFrame_IsPFrame(t *testing.T) {
	tests := []struct {
		name     string
		flags    uint32
		expected bool
	}{
		{"PFrame flag set", v4l2.BufFlagPFrame, true},
		{"PFrame with other flags", v4l2.BufFlagPFrame | v4l2.BufFlagMapped, true},
		{"No PFrame flag", v4l2.BufFlagMapped, false},
		{"Zero flags", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := &Frame{Flags: tt.flags}
			if result := frame.IsPFrame(); result != tt.expected {
				t.Errorf("IsPFrame() = %v, want %v (flags=0x%x)", result, tt.expected, tt.flags)
			}
		})
	}
}

// TestFrame_IsBFrame tests the IsBFrame() method
func TestFrame_IsBFrame(t *testing.T) {
	tests := []struct {
		name     string
		flags    uint32
		expected bool
	}{
		{"BFrame flag set", v4l2.BufFlagBFrame, true},
		{"BFrame with other flags", v4l2.BufFlagBFrame | v4l2.BufFlagMapped, true},
		{"No BFrame flag", v4l2.BufFlagMapped, false},
		{"Zero flags", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := &Frame{Flags: tt.flags}
			if result := frame.IsBFrame(); result != tt.expected {
				t.Errorf("IsBFrame() = %v, want %v (flags=0x%x)", result, tt.expected, tt.flags)
			}
		})
	}
}

// TestFrame_HasError tests the HasError() method
func TestFrame_HasError(t *testing.T) {
	tests := []struct {
		name     string
		flags    uint32
		expected bool
	}{
		{"Error flag set", v4l2.BufFlagError, true},
		{"Error with other flags", v4l2.BufFlagError | v4l2.BufFlagMapped, true},
		{"No Error flag", v4l2.BufFlagMapped, false},
		{"Zero flags", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := &Frame{Flags: tt.flags}
			if result := frame.HasError(); result != tt.expected {
				t.Errorf("HasError() = %v, want %v (flags=0x%x)", result, tt.expected, tt.flags)
			}
		})
	}
}

// TestFrame_Metadata tests that Frame correctly stores metadata
func TestFrame_Metadata(t *testing.T) {
	now := time.Now()
	frame := &Frame{
		Data:      make([]byte, 1024),
		Timestamp: now,
		Sequence:  42,
		Flags:     v4l2.BufFlagKeyFrame | v4l2.BufFlagMapped,
		Index:     3,
	}

	if !frame.Timestamp.Equal(now) {
		t.Errorf("Timestamp = %v, want %v", frame.Timestamp, now)
	}
	if frame.Sequence != 42 {
		t.Errorf("Sequence = %d, want 42", frame.Sequence)
	}
	if frame.Index != 3 {
		t.Errorf("Index = %d, want 3", frame.Index)
	}
	if !frame.IsKeyFrame() {
		t.Error("Frame should be a keyframe")
	}
}

// TestFrame_DataIntegrity tests that frame data is correctly copied
func TestFrame_DataIntegrity(t *testing.T) {
	pool := NewFramePool(1024)
	data := pool.Get(512)

	// Fill with test data
	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}

	frame := &Frame{
		Data: data,
		pool: pool,
	}

	// Verify data integrity
	for i := 0; i < len(frame.Data); i++ {
		if frame.Data[i] != byte(i%256) {
			t.Errorf("Data[%d] = %d, want %d", i, frame.Data[i], byte(i%256))
			break
		}
	}

	frame.Release()
}

// TestFrame_StructFields tests that all Frame struct fields are accessible
func TestFrame_StructFields(t *testing.T) {
	timestamp := time.Now()
	data := make([]byte, 1024)
	pool := NewFramePool(1024)

	frame := Frame{
		Data:      data,
		Timestamp: timestamp,
		Sequence:  100,
		Flags:     v4l2.BufFlagKeyFrame,
		Index:     5,
		pool:      pool,
		released:  false,
	}

	if len(frame.Data) != 1024 {
		t.Errorf("Data length = %d, want 1024", len(frame.Data))
	}
	if !frame.Timestamp.Equal(timestamp) {
		t.Error("Timestamp mismatch")
	}
	if frame.Sequence != 100 {
		t.Errorf("Sequence = %d, want 100", frame.Sequence)
	}
	if frame.Flags != v4l2.BufFlagKeyFrame {
		t.Errorf("Flags = 0x%x, want 0x%x", frame.Flags, v4l2.BufFlagKeyFrame)
	}
	if frame.Index != 5 {
		t.Errorf("Index = %d, want 5", frame.Index)
	}
	if frame.pool != pool {
		t.Error("pool reference mismatch")
	}
	if frame.released != false {
		t.Error("released should be false")
	}
}

// TestFrame_CombinedFlags tests frames with multiple flags set
func TestFrame_CombinedFlags(t *testing.T) {
	// Frame with multiple flags
	frame := &Frame{
		Flags: v4l2.BufFlagKeyFrame | v4l2.BufFlagMapped | v4l2.BufFlagError,
	}

	if !frame.IsKeyFrame() {
		t.Error("Should be identified as keyframe")
	}
	if !frame.HasError() {
		t.Error("Should be identified as having error")
	}
	if frame.IsPFrame() {
		t.Error("Should not be identified as P-frame")
	}
	if frame.IsBFrame() {
		t.Error("Should not be identified as B-frame")
	}
}
