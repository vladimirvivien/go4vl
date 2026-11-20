package v4l2

import (
	"context"
	"errors"
	"testing"
	"time"

	sys "golang.org/x/sys/unix"
)

func TestPollPipeReadWrite(t *testing.T) {
	var fds [2]int
	if err := sys.Pipe2(fds[:], sys.O_NONBLOCK|sys.O_CLOEXEC); err != nil {
		t.Fatalf("Pipe2: %v", err)
	}
	defer sys.Close(fds[0])
	defer sys.Close(fds[1])

	// Write end should be immediately PollOut ready.
	ev, err := Poll(context.Background(), uintptr(fds[1]), PollOut, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Poll(PollOut): %v", err)
	}
	if ev&PollOut == 0 {
		t.Fatalf("expected PollOut, got 0x%x", ev)
	}

	// Read end should time out with (0, nil) when empty.
	ev, err = Poll(context.Background(), uintptr(fds[0]), PollIn, 15*time.Millisecond)
	if err != nil {
		t.Fatalf("Poll(timeout): %v", err)
	}
	if ev != 0 {
		t.Fatalf("expected 0 events on timeout, got 0x%x", ev)
	}

	// Write bytes via WriteDevice and verify PollIn fires.
	payload := []byte("frame-data")
	n, err := WriteDevice(uintptr(fds[1]), payload)
	if err != nil || n != len(payload) {
		t.Fatalf("WriteDevice: n=%d err=%v", n, err)
	}

	ev, err = Poll(context.Background(), uintptr(fds[0]), PollIn, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Poll(PollIn): %v", err)
	}
	if ev&PollIn == 0 {
		t.Fatalf("expected PollIn, got 0x%x", ev)
	}

	buf := make([]byte, len(payload))
	rn, err := ReadDevice(uintptr(fds[0]), buf)
	if err != nil || rn != len(payload) || string(buf) != string(payload) {
		t.Fatalf("ReadDevice: n=%d data=%q err=%v", rn, string(buf), err)
	}
}

func TestPollContextCancel(t *testing.T) {
	var fds [2]int
	if err := sys.Pipe2(fds[:], sys.O_NONBLOCK|sys.O_CLOEXEC); err != nil {
		t.Fatalf("Pipe2: %v", err)
	}
	defer sys.Close(fds[0])
	defer sys.Close(fds[1])

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := Poll(ctx, uintptr(fds[0]), PollIn, 0)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestPollerMultiFDAndCancel(t *testing.T) {
	var p1, p2 [2]int
	if err := sys.Pipe2(p1[:], sys.O_NONBLOCK|sys.O_CLOEXEC); err != nil {
		t.Fatalf("Pipe2 p1: %v", err)
	}
	defer sys.Close(p1[0])
	defer sys.Close(p1[1])

	if err := sys.Pipe2(p2[:], sys.O_NONBLOCK|sys.O_CLOEXEC); err != nil {
		t.Fatalf("Pipe2 p2: %v", err)
	}
	defer sys.Close(p2[0])
	defer sys.Close(p2[1])

	poller, err := NewPoller()
	if err != nil {
		t.Fatalf("NewPoller: %v", err)
	}
	defer poller.Close()

	if err := poller.Add(uintptr(p1[0]), PollIn); err != nil {
		t.Fatalf("Add p1: %v", err)
	}
	if err := poller.Add(uintptr(p2[0]), PollIn); err != nil {
		t.Fatalf("Add p2: %v", err)
	}

	results := make([]PollResult, 4)

	// Timeout when neither pipe is readable.
	n, err := poller.Wait(context.Background(), results, 15*time.Millisecond)
	if err != nil || n != 0 {
		t.Fatalf("Wait timeout: n=%d err=%v", n, err)
	}

	// Trigger pipe 2 and verify epoll returns it.
	if _, err := WriteDevice(uintptr(p2[1]), []byte("x")); err != nil {
		t.Fatalf("WriteDevice: %v", err)
	}
	n, err = poller.Wait(context.Background(), results, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Wait ready: %v", err)
	}
	if n != 1 || results[0].FD != uintptr(p2[0]) || results[0].Events&PollIn == 0 {
		t.Fatalf("unexpected Wait result: n=%d res=%+v", n, results[:n])
	}

	// Drain pipe 2 and verify immediate wakeup on context cancellation.
	var tmp [1]byte
	_, _ = ReadDevice(uintptr(p2[0]), tmp[:])

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(15 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err = poller.Wait(ctx, results, 5*time.Second)
	elapsed := time.Since(start)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("context cancel wakeup too slow: %v", elapsed)
	}

	// Modify and Remove.
	if err := poller.Modify(uintptr(p1[0]), PollIn|PollPri); err != nil {
		t.Fatalf("Modify: %v", err)
	}
	if err := poller.Remove(uintptr(p1[0])); err != nil {
		t.Fatalf("Remove: %v", err)
	}
}
