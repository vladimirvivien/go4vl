package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

// Read reads a single frame from the device into buf.
// Returns the number of bytes read.
//
// Only available when the device is opened with IOMethodReadWrite.
// The buffer should be at least SizeImage bytes (from GetPixFormat().SizeImage).
// Because the device is opened with O_NONBLOCK, EAGAIN is returned if no frame
// is ready yet; use ReadContext to wait for readiness with poll().
func (d *Device) Read(buf []byte) (int, error) {
	if d.config.ioMethod != IOMethodReadWrite {
		return 0, fmt.Errorf("device: Read() not supported in streaming IO mode; use Start/Frames/GetFrames instead")
	}

	n, err := v4l2.ReadDevice(d.fd, buf)
	if err != nil {
		if errors.Is(err, sys.EAGAIN) || errors.Is(err, sys.EINTR) {
			return 0, err
		}
		return 0, fmt.Errorf("device: read: %w", err)
	}
	return n, nil
}

// ReadContext waits until a frame is readable (using v4l2.Poll) and reads it into buf,
// transparently retrying on O_NONBLOCK EAGAIN/EINTR until data is read or ctx is cancelled.
// Only available when the device is opened with IOMethodReadWrite.
func (d *Device) ReadContext(ctx context.Context, buf []byte) (int, error) {
	if d.config.ioMethod != IOMethodReadWrite {
		return 0, fmt.Errorf("device: ReadContext() not supported in streaming IO mode; use Start/Frames/GetFrames instead")
	}

	for {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		n, err := v4l2.ReadDevice(d.fd, buf)
		if err == nil {
			return n, nil
		}
		if errors.Is(err, sys.EINTR) {
			continue
		}
		if errors.Is(err, sys.EAGAIN) {
			// Bound each poll wait because drivers such as v4l2loopback do not signal readiness on read()-only fds.
			ev, pollErr := v4l2.Poll(ctx, d.fd, v4l2.PollIn, 50*time.Millisecond)
			if pollErr != nil {
				return 0, pollErr
			}
			if ev&(v4l2.PollErr|v4l2.PollHup|v4l2.PollNVal) != 0 {
				return 0, fmt.Errorf("device: read poll error event 0x%x", ev)
			}
			continue
		}
		return 0, fmt.Errorf("device: read: %w", err)
	}
}

// ReadFrame reads a single frame from the device and returns it as a *Frame
// with metadata, providing parity with the streaming GetFrames() API.
//
// Each call allocates a fresh buffer. The returned Frame.Data is valid until
// the caller discards it (no Release() needed, unlike streaming mode frames).
// Only available when the device is opened with IOMethodReadWrite.
func (d *Device) ReadFrame() (*Frame, error) {
	if d.config.ioMethod != IOMethodReadWrite {
		return nil, fmt.Errorf("device: ReadFrame() not supported in streaming IO mode; use Start/Frames/GetFrames instead")
	}

	buf := make([]byte, d.config.pixFormat.SizeImage)
	n, err := v4l2.ReadDevice(d.fd, buf)
	if err != nil {
		if errors.Is(err, sys.EAGAIN) || errors.Is(err, sys.EINTR) {
			return nil, err
		}
		return nil, fmt.Errorf("device: read frame: %w", err)
	}

	frame := &Frame{
		Data:      buf[:n],
		Timestamp: time.Now(),
		Sequence:  d.readSeq,
	}
	d.readSeq++

	return frame, nil
}

// ReadFrameContext waits until a frame is readable and returns it as a *Frame,
// retrying on EAGAIN/EINTR until data arrives or ctx is cancelled.
// Only available when the device is opened with IOMethodReadWrite.
func (d *Device) ReadFrameContext(ctx context.Context) (*Frame, error) {
	if d.config.ioMethod != IOMethodReadWrite {
		return nil, fmt.Errorf("device: ReadFrameContext() not supported in streaming IO mode; use Start/Frames/GetFrames instead")
	}

	buf := make([]byte, d.config.pixFormat.SizeImage)
	n, err := d.ReadContext(ctx, buf)
	if err != nil {
		return nil, err
	}

	frame := &Frame{
		Data:      buf[:n],
		Timestamp: time.Now(),
		Sequence:  d.readSeq,
	}
	d.readSeq++

	return frame, nil
}

// Write writes a single frame to a video output device using the write() syscall.
// Only available when the device is opened with IOMethodReadWrite and configured
// for video output (BufTypeVideoOutput).
// Because the device is opened with O_NONBLOCK, EAGAIN is returned if the driver
// cannot accept data immediately; use WriteContext to wait with poll().
func (d *Device) Write(buf []byte) (int, error) {
	if d.config.ioMethod != IOMethodReadWrite {
		return 0, fmt.Errorf("device: Write() not supported in streaming IO mode; use Start/WriteFrame instead")
	}
	if d.bufType != v4l2.BufTypeVideoOutput && d.bufType != v4l2.BufTypeVideoOutputMPlane {
		return 0, fmt.Errorf("device: Write() requires a video output device")
	}

	n, err := v4l2.WriteDevice(d.fd, buf)
	if err != nil {
		if errors.Is(err, sys.EAGAIN) || errors.Is(err, sys.EINTR) {
			return 0, err
		}
		return 0, fmt.Errorf("device: write: %w", err)
	}
	return n, nil
}

// WriteContext waits until the output device is writable (using v4l2.Poll) and writes
// buf, retrying on O_NONBLOCK EAGAIN/EINTR until the write succeeds or ctx is cancelled.
// Only available when the device is opened with IOMethodReadWrite and configured
// for video output.
func (d *Device) WriteContext(ctx context.Context, buf []byte) (int, error) {
	if d.config.ioMethod != IOMethodReadWrite {
		return 0, fmt.Errorf("device: WriteContext() not supported in streaming IO mode; use Start/WriteFrame instead")
	}
	if d.bufType != v4l2.BufTypeVideoOutput && d.bufType != v4l2.BufTypeVideoOutputMPlane {
		return 0, fmt.Errorf("device: WriteContext() requires a video output device")
	}

	for {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		n, err := v4l2.WriteDevice(d.fd, buf)
		if err == nil {
			return n, nil
		}
		if errors.Is(err, sys.EINTR) {
			continue
		}
		if errors.Is(err, sys.EAGAIN) {
			// Bound each poll wait to match ReadContext and retry write() on drivers with incomplete poll() support.
			ev, pollErr := v4l2.Poll(ctx, d.fd, v4l2.PollOut, 50*time.Millisecond)
			if pollErr != nil {
				return 0, pollErr
			}
			if ev&(v4l2.PollErr|v4l2.PollHup|v4l2.PollNVal) != 0 {
				return 0, fmt.Errorf("device: write poll error event 0x%x", ev)
			}
			continue
		}
		return 0, fmt.Errorf("device: write: %w", err)
	}
}
