package v4l2

import (
	"context"
	"testing"
	"time"
)

// mockDevice implements the Device interface for testing
type mockDevice struct {
	fd uintptr
}

func (m *mockDevice) Name() string {
	return "/dev/video999"
}

func (m *mockDevice) Fd() uintptr {
	return m.fd
}

func (m *mockDevice) Capability() Capability {
	return Capability{}
}

func (m *mockDevice) MemIOType() IOType {
	return IOTypeMMAP
}

func (m *mockDevice) GetOutput() <-chan []byte {
	return nil
}

func (m *mockDevice) SetInput(<-chan []byte) {
}

func (m *mockDevice) Close() error {
	return nil
}

// TestWaitForRead_ContextCancellation tests that WaitForRead respects context cancellation
func TestWaitForRead_ContextCancellation(t *testing.T) {
	// We use a mock fd value since we can't easily create a real V4L2 device in unit tests
	// The key is testing the context cancellation behavior

	// Use a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	dev := &mockDevice{fd: 3} // Mock fd (typical range is 0-15 for FdSet)

	sigChan := WaitForRead(ctx, dev)

	// The channel should close quickly because context is already cancelled
	select {
	case _, ok := <-sigChan:
		if ok {
			t.Error("Expected channel to be closed")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Channel did not close within timeout")
	}
}

// TestWaitForRead_ContextCancellationDuringWait tests cancellation while waiting
func TestWaitForRead_ContextCancellationDuringWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dev := &mockDevice{fd: 4} // Mock fd

	sigChan := WaitForRead(ctx, dev)

	// Let it start waiting
	time.Sleep(50 * time.Millisecond)

	// Cancel the context
	cancel()

	// The channel should close within a reasonable time
	// WaitForRead uses 100ms timeout, so it should respond within ~150ms
	select {
	case _, ok := <-sigChan:
		if ok {
			t.Error("Expected channel to be closed after cancellation")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Channel did not close within timeout after cancellation")
	}
}

// TestWaitForRead_ChannelClosedOnReturn tests that the channel is closed when goroutine exits
func TestWaitForRead_ChannelClosedOnReturn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	dev := &mockDevice{fd: 5}

	sigChan := WaitForRead(ctx, dev)

	// Wait for context to timeout
	time.Sleep(100 * time.Millisecond)

	// Verify channel is closed
	_, ok := <-sigChan
	if ok {
		t.Error("Expected channel to be closed after context timeout")
	}
}

// TestWaitForRead_MultipleCancellations tests behavior with rapid context cancellations
func TestWaitForRead_MultipleCancellations(t *testing.T) {
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithCancel(context.Background())

		dev := &mockDevice{fd: uintptr(6 + i)}
		sigChan := WaitForRead(ctx, dev)

		// Cancel immediately
		cancel()

		// Verify channel closes
		select {
		case _, ok := <-sigChan:
			if ok {
				t.Errorf("Iteration %d: Expected channel to be closed", i)
			}
		case <-time.After(200 * time.Millisecond):
			t.Errorf("Iteration %d: Channel did not close within timeout", i)
		}
	}
}

// TestWaitForRead_ConcurrentContextCancellation tests concurrent cancellations
func TestWaitForRead_ConcurrentContextCancellation(t *testing.T) {
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			dev := &mockDevice{fd: uintptr(3 + (id % 10))}
			sigChan := WaitForRead(ctx, dev)

			// Let it run briefly
			time.Sleep(20 * time.Millisecond)

			// Cancel
			cancel()

			// Verify closure
			select {
			case _, ok := <-sigChan:
				if ok {
					t.Errorf("Goroutine %d: Expected channel to be closed", id)
				}
			case <-time.After(200 * time.Millisecond):
				t.Errorf("Goroutine %d: Channel did not close within timeout", id)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	timeout := time.After(2 * time.Second)
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
			// Success
		case <-timeout:
			t.Fatalf("Timeout waiting for goroutine %d to complete", i)
		}
	}
}

// TestWaitForRead_NoLeaksOnCancel verifies no goroutine leaks occur
func TestWaitForRead_NoLeaksOnCancel(t *testing.T) {
	// This test verifies the goroutine cleanup behavior
	// We create multiple WaitForRead calls and cancel them quickly

	for i := 0; i < 20; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		dev := &mockDevice{fd: uintptr(3 + (i % 10))}

		sigChan := WaitForRead(ctx, dev)

		// Cancel immediately
		cancel()

		// Drain the channel to ensure goroutine can exit
		select {
		case <-sigChan:
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("Iteration %d: Channel did not close, potential goroutine leak", i)
		}
	}

	// If we get here without hanging, no obvious leaks occurred
	// Note: Full leak detection would require runtime.NumGoroutine() checks
	// but those can be flaky in tests
}

// TestWaitForRead_TimeoutBehavior tests the 100ms select timeout
func TestWaitForRead_TimeoutBehavior(t *testing.T) {
	// This test verifies that WaitForRead respects the 100ms timeout
	// and can respond to cancellation even when no data is available

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dev := &mockDevice{fd: 13}
	sigChan := WaitForRead(ctx, dev)

	// Wait for more than one timeout cycle (>100ms)
	time.Sleep(150 * time.Millisecond)

	// Now cancel
	cancel()

	// Should close within one more timeout cycle
	select {
	case _, ok := <-sigChan:
		if ok {
			t.Error("Expected channel to be closed")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Channel should close within 200ms of cancellation")
	}
}
