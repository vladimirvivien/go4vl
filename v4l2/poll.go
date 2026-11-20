package v4l2

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	sys "golang.org/x/sys/unix"
)

// PollEvent is a bitmask of file descriptor readiness events used with Poll and Poller.
type PollEvent uint32

const (
	// PollIn indicates data or a capture buffer is ready to read or dequeue.
	PollIn PollEvent = sys.POLLIN

	// PollOut indicates an output buffer is ready to write or dequeue.
	PollOut PollEvent = sys.POLLOUT

	// PollPri indicates a V4L2 event is pending for DequeueEvent (e.g., source change or EOS).
	PollPri PollEvent = sys.POLLPRI

	// PollErr indicates an error condition on the file descriptor.
	PollErr PollEvent = sys.POLLERR

	// PollHup indicates a hangup on the file descriptor.
	PollHup PollEvent = sys.POLLHUP

	// PollNVal indicates the file descriptor is not open.
	PollNVal PollEvent = sys.POLLNVAL
)

// PollResult reports readiness events for a single file descriptor returned by Poller.Wait.
type PollResult struct {
	FD     uintptr
	Events PollEvent
}

// Poll waits synchronously for any of the requested events on fd without spawning a goroutine.
//
// If timeout > 0, Poll returns (0, nil) when the timeout expires before an event occurs.
// If timeout <= 0, Poll waits until an event occurs or ctx is cancelled.
// When ctx is cancelled while waiting, Poll returns (0, ctx.Err()).
func Poll(ctx context.Context, fd uintptr, events PollEvent, timeout time.Duration) (PollEvent, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	var deadline time.Time
	hasTimeout := timeout > 0
	deadlineFromCtx := false
	if hasTimeout {
		deadline = time.Now().Add(timeout)
	}
	if ctxDeadline, ok := ctx.Deadline(); ok {
		if !hasTimeout || ctxDeadline.Before(deadline) {
			deadline = ctxDeadline
			hasTimeout = true
			deadlineFromCtx = true
		}
	}

	fds := []sys.PollFd{{
		Fd:     int32(fd),
		Events: int16(events),
	}}

	const maxCancelSliceMs = 50

	for {
		if err := ctx.Err(); err != nil {
			return 0, err
		}

		waitMs := -1
		if hasTimeout {
			rem := time.Until(deadline)
			if rem <= 0 {
				if err := ctx.Err(); err != nil {
					return 0, err
				}
				if deadlineFromCtx {
					return 0, context.DeadlineExceeded
				}
				return 0, nil
			}
			waitMs = int((rem + time.Millisecond - 1) / time.Millisecond)
		}

		// Bound each poll slice when ctx can be cancelled without a shorter deadline.
		if ctx.Done() != nil && (waitMs < 0 || waitMs > maxCancelSliceMs) {
			waitMs = maxCancelSliceMs
		}

		fds[0].Revents = 0
		n, err := sys.Poll(fds, waitMs)
		if err != nil {
			if errors.Is(err, sys.EINTR) {
				continue
			}
			return 0, fmt.Errorf("poll fd %d: %w", fd, err)
		}
		if n > 0 {
			return PollEvent(uint16(fds[0].Revents)), nil
		}
		if hasTimeout && !time.Now().Before(deadline) {
			if err := ctx.Err(); err != nil {
				return 0, err
			}
			if deadlineFromCtx {
				return 0, context.DeadlineExceeded
			}
			return 0, nil
		}
	}
}

// Poller multiplexes readiness events (PollIn, PollOut, PollPri) across multiple file
// descriptors using Linux epoll, with an internal eventfd for immediate context cancellation.
type Poller struct {
	epfd    int
	wakefd  int
	events  []sys.EpollEvent
	closeMu sync.Mutex
	closed  bool
}

// NewPoller creates a new epoll-backed Poller. Call Close when finished to release
// the underlying epoll and eventfd file descriptors.
func NewPoller() (*Poller, error) {
	epfd, err := sys.EpollCreate1(sys.EPOLL_CLOEXEC)
	if err != nil {
		return nil, fmt.Errorf("epoll create: %w", err)
	}

	wakefd, err := sys.Eventfd(0, sys.EFD_NONBLOCK|sys.EFD_CLOEXEC)
	if err != nil {
		_ = sys.Close(epfd)
		return nil, fmt.Errorf("eventfd create: %w", err)
	}

	wakeEv := sys.EpollEvent{
		Events: sys.EPOLLIN,
		Fd:     int32(wakefd),
	}
	if err := sys.EpollCtl(epfd, sys.EPOLL_CTL_ADD, wakefd, &wakeEv); err != nil {
		_ = sys.Close(wakefd)
		_ = sys.Close(epfd)
		return nil, fmt.Errorf("epoll add wakeup fd: %w", err)
	}

	return &Poller{
		epfd:   epfd,
		wakefd: wakefd,
		events: make([]sys.EpollEvent, 16),
	}, nil
}

// Add registers fd with the poller for the given event mask.
func (p *Poller) Add(fd uintptr, events PollEvent) error {
	ev := sys.EpollEvent{
		Events: uint32(events),
		Fd:     int32(fd),
	}
	if err := sys.EpollCtl(p.epfd, sys.EPOLL_CTL_ADD, int(fd), &ev); err != nil {
		return fmt.Errorf("epoll add fd %d: %w", fd, err)
	}
	return nil
}

// Modify updates the event mask watched for a previously added fd.
func (p *Poller) Modify(fd uintptr, events PollEvent) error {
	ev := sys.EpollEvent{
		Events: uint32(events),
		Fd:     int32(fd),
	}
	if err := sys.EpollCtl(p.epfd, sys.EPOLL_CTL_MOD, int(fd), &ev); err != nil {
		return fmt.Errorf("epoll modify fd %d: %w", fd, err)
	}
	return nil
}

// Remove unregisters fd from the poller.
func (p *Poller) Remove(fd uintptr) error {
	if err := sys.EpollCtl(p.epfd, sys.EPOLL_CTL_DEL, int(fd), nil); err != nil {
		return fmt.Errorf("epoll remove fd %d: %w", fd, err)
	}
	return nil
}

// Wake interrupts an in-progress Wait call.
func (p *Poller) Wake() error {
	p.closeMu.Lock()
	defer p.closeMu.Unlock()
	if p.closed {
		return nil
	}
	var buf [8]byte
	binary.NativeEndian.PutUint64(buf[:], 1)
	_, err := sys.Write(p.wakefd, buf[:])
	if errors.Is(err, sys.EAGAIN) {
		return nil
	}
	return err
}

func (p *Poller) drainWake() {
	var buf [8]byte
	_, _ = sys.Read(p.wakefd, buf[:])
}

// Wait blocks until at least one registered file descriptor is ready, timeout expires,
// or ctx is cancelled. Ready descriptors are written into dst and the count is returned.
//
// Pass timeout <= 0 to wait indefinitely until an event or context cancellation.
func (p *Poller) Wait(ctx context.Context, dst []PollResult, timeout time.Duration) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if len(dst) == 0 {
		return 0, fmt.Errorf("epoll wait: destination slice must not be empty")
	}

	// Register context cancellation hook to wake epoll_wait immediately.
	var stopCancel func() bool
	if ctx.Done() != nil {
		stopCancel = context.AfterFunc(ctx, func() {
			_ = p.Wake()
		})
		defer stopCancel()
	}

	var deadline time.Time
	hasTimeout := timeout > 0
	if hasTimeout {
		deadline = time.Now().Add(timeout)
	}

	// Ensure internal buffer can hold len(dst)+1 (accounting for wakefd).
	need := len(dst) + 1
	if cap(p.events) < need {
		p.events = make([]sys.EpollEvent, need)
	}
	evBuf := p.events[:need]

	for {
		if err := ctx.Err(); err != nil {
			return 0, err
		}

		waitMs := -1
		if hasTimeout {
			rem := time.Until(deadline)
			if rem <= 0 {
				return 0, nil
			}
			waitMs = int((rem + time.Millisecond - 1) / time.Millisecond)
		}

		n, err := sys.EpollWait(p.epfd, evBuf, waitMs)
		if err != nil {
			if errors.Is(err, sys.EINTR) {
				continue
			}
			return 0, fmt.Errorf("epoll wait: %w", err)
		}
		if n == 0 {
			if err := ctx.Err(); err != nil {
				return 0, err
			}
			return 0, nil
		}

		out := 0
		woke := false
		for i := 0; i < n; i++ {
			if int(evBuf[i].Fd) == p.wakefd {
				woke = true
				continue
			}
			if out < len(dst) {
				dst[out] = PollResult{
					FD:     uintptr(evBuf[i].Fd),
					Events: PollEvent(evBuf[i].Events),
				}
				out++
			}
		}
		if woke {
			p.drainWake()
		}
		if out > 0 {
			return out, nil
		}
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		if woke {
			// Woken explicitly via Wake() with no other ready FDs.
			return 0, nil
		}
	}
}

// Close releases the epoll and wakeup eventfd file descriptors.
func (p *Poller) Close() error {
	p.closeMu.Lock()
	defer p.closeMu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	err1 := sys.Close(p.wakefd)
	err2 := sys.Close(p.epfd)
	return errors.Join(err1, err2)
}
