package device

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

// Frames returns a zero-copy iterator over captured video frames.
//
// Each yielded *Frame points directly into the underlying MMAP, USERPTR, or
// DMA-BUF buffer with zero per-frame heap allocations and no background
// goroutines. The buffer is automatically re-queued to the driver (VIDIOC_QBUF)
// immediately after each loop iteration finishes—including on continue, break,
// return, or panic.
//
// Frame.Data is valid only for the duration of the current loop iteration.
// Call frame.Clone() if you need to retain frame data beyond the iteration or
// pass it to another goroutine.
//
// Call Start(ctx) before ranging over Frames(ctx). Once Frames is used for a
// Start/Stop cycle, GetFrames() and GetOutput() cannot be used until Stop() is
// called.
//
// Example:
//
//	if err := dev.Start(ctx); err != nil {
//	    return err
//	}
//	defer dev.Stop()
//
//	for frame, err := range dev.Frames(ctx) {
//	    if err != nil {
//	        return err
//	    }
//	    process(frame.Data)
//	}
func (d *Device) Frames(ctx context.Context) iter.Seq2[*Frame, error] {
	return func(yield func(*Frame, error) bool) {
		if d.config.ioMethod == IOMethodReadWrite {
			yield(nil, fmt.Errorf("device: Frames() requires streaming IO mode"))
			return
		}
		if !d.streaming.Load() || !d.streamActive.Load() {
			yield(nil, fmt.Errorf("device: Frames() called before Start()"))
			return
		}
		if !d.streamingMode.CompareAndSwap(0, 3) && d.streamingMode.Load() != 3 {
			yield(nil, fmt.Errorf("device: Frames() called while another streaming mode is active"))
			return
		}

		fd := d.Fd()
		ioMemType := d.MemIOType()
		bufType := d.BufferType()

		var frame Frame
		checkedOut := -1

		// Ensure any checked-out buffer is returned to the driver even if the
		// caller's loop body panics.
		defer func() {
			if checkedOut >= 0 && d.streaming.Load() {
				_ = d.queueBuffer(uint32(checkedOut))
			}
		}()

		for {
			if !d.streaming.Load() {
				return
			}
			if err := ctx.Err(); err != nil {
				return
			}

			ev, err := v4l2.Poll(ctx, fd, v4l2.PollIn, 50*time.Millisecond)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return
				}
				yield(nil, fmt.Errorf("device: frames poll: %w", err))
				return
			}
			if ev == 0 {
				continue
			}
			if ev&(v4l2.PollErr|v4l2.PollHup|v4l2.PollNVal) != 0 && ev&v4l2.PollIn == 0 {
				yield(nil, fmt.Errorf("device: frames poll error event 0x%x", ev))
				return
			}

			buff, err := v4l2.DequeueBuffer(fd, ioMemType, bufType)
			if err != nil {
				if errors.Is(err, sys.EAGAIN) {
					continue
				}
				yield(nil, fmt.Errorf("device: frames dequeue: %w", err))
				return
			}

			if !d.streaming.Load() {
				return
			}
			if buff.Index >= uint32(len(d.buffers)) || d.buffers[buff.Index] == nil {
				return
			}

			isMapped := (ioMemType == v4l2.IOTypeUserPtr) || (buff.Flags&v4l2.BufFlagMapped != 0)
			hasError := buff.Flags&v4l2.BufFlagError != 0
			hasData := buff.BytesUsed > 0

			keepGoing := true
			if isMapped && !hasError && hasData {
				frame = Frame{
					Data:      d.buffers[buff.Index][:buff.BytesUsed],
					Timestamp: v4l2.BufferTimestamp(buff.Timestamp, buff.Flags, d.config.timestampMode),
					Sequence:  buff.Sequence,
					Flags:     buff.Flags,
					Index:     buff.Index,
				}
				checkedOut = int(buff.Index)
				keepGoing = yield(&frame, nil)
				checkedOut = -1
			} else if hasError {
				checkedOut = int(buff.Index)
				keepGoing = yield(nil, fmt.Errorf("device: buffer error flag: index=%d flags=0x%x", buff.Index, buff.Flags))
				checkedOut = -1
			}

			if !d.streaming.Load() {
				return
			}
			if err := d.queueBuffer(buff.Index); err != nil {
				if keepGoing {
					yield(nil, fmt.Errorf("device: frames requeue: %w", err))
				}
				return
			}
			if !keepGoing {
				return
			}
		}
	}
}
