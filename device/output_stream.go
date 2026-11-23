package device

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unsafe"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

// WriteFrame writes a single frame to a streaming video output device
// (BufTypeVideoOutput or single-plane BufTypeVideoOutputMPlane) using the
// device's allocated MMAP, USERPTR, or DMA-BUF buffers. Call Start(ctx) before
// calling WriteFrame.
func (d *Device) WriteFrame(ctx context.Context, data []byte) error {
	if d.bufType == v4l2.BufTypeVideoOutputMPlane {
		return d.WriteFramePlanes(ctx, [][]byte{data})
	}
	if d.config.ioMethod == IOMethodReadWrite {
		return fmt.Errorf("device: WriteFrame() not supported in read/write IO mode; use Write() instead")
	}
	if d.bufType != v4l2.BufTypeVideoOutput {
		return fmt.Errorf("device: WriteFrame() requires a video output device (got bufType %d)", d.bufType)
	}
	if !d.streaming.Load() || !d.streamActive.Load() {
		return fmt.Errorf("device: WriteFrame() called before Start()")
	}
	if len(d.buffers) == 0 {
		return fmt.Errorf("device: WriteFrame() has no allocated buffers")
	}

	var idx uint32
	if d.outputQueued < uint32(len(d.buffers)) {
		idx = d.outputQueued
	} else {
		for {
			if err := ctx.Err(); err != nil {
				return err
			}
			ev, err := v4l2.Poll(ctx, d.fd, v4l2.PollOut, 50*time.Millisecond)
			if err != nil {
				return err
			}
			if ev == 0 {
				continue
			}
			buf, err := v4l2.DequeueBuffer(d.fd, d.config.ioType, d.bufType)
			if err != nil {
				if errors.Is(err, sys.EAGAIN) {
					continue
				}
				return fmt.Errorf("device: output dequeue: %w", err)
			}
			idx = buf.Index
			break
		}
	}

	if int(idx) >= len(d.buffers) {
		return fmt.Errorf("device: output buffer index %d out of range", idx)
	}
	dst := d.buffers[idx]
	if len(data) > len(dst) {
		return fmt.Errorf("device: frame size %d exceeds output buffer length %d", len(data), len(dst))
	}
	copy(dst, data)

	qbuf := v4l2.Buffer{
		Index:     idx,
		Type:      d.bufType,
		Memory:    d.config.ioType,
		BytesUsed: uint32(len(data)),
		Length:    uint32(len(dst)),
	}
	switch d.config.ioType {
	case v4l2.IOTypeUserPtr:
		qbuf.Info.UserPtr = uintptr(unsafe.Pointer(&dst[0]))
	case v4l2.IOTypeDMABuf:
		qbuf.Info.FD = d.dmabufFDs[idx]
	}

	if _, err := v4l2.QueueBufferWith(d.fd, qbuf); err != nil {
		return fmt.Errorf("device: output queue: %w", err)
	}
	if d.outputQueued < uint32(len(d.buffers)) {
		d.outputQueued++
	}

	if !d.outputStreamOn {
		if err := v4l2.StreamOn(d); err != nil {
			return fmt.Errorf("device: output stream on: %w", err)
		}
		d.outputStreamOn = true
	}

	return nil
}

// WriteFramePlanes writes a multi-planar frame to a streaming video output device
// configured with BufTypeVideoOutputMPlane. Call Start(ctx) before calling WriteFramePlanes.
func (d *Device) WriteFramePlanes(ctx context.Context, planes [][]byte) error {
	if d.config.ioMethod == IOMethodReadWrite {
		return fmt.Errorf("device: WriteFramePlanes() not supported in read/write IO mode")
	}
	if d.bufType != v4l2.BufTypeVideoOutputMPlane {
		return fmt.Errorf("device: WriteFramePlanes() requires BufTypeVideoOutputMPlane (got bufType %d)", d.bufType)
	}
	if !d.streaming.Load() || !d.streamActive.Load() {
		return fmt.Errorf("device: WriteFramePlanes() called before Start()")
	}
	if len(d.buffersMPlane) == 0 {
		return fmt.Errorf("device: WriteFramePlanes() has no allocated buffers")
	}
	if uint32(len(planes)) != d.numPlanes {
		return fmt.Errorf("device: WriteFramePlanes() got %d planes, want %d", len(planes), d.numPlanes)
	}

	var idx uint32
	if d.outputQueued < uint32(len(d.buffersMPlane)) {
		idx = d.outputQueued
	} else {
		for {
			if err := ctx.Err(); err != nil {
				return err
			}
			ev, err := v4l2.Poll(ctx, d.fd, v4l2.PollOut, 50*time.Millisecond)
			if err != nil {
				return err
			}
			if ev == 0 {
				continue
			}
			buf, err := v4l2.DequeueBufferMPlane(d.fd, d.config.ioType, d.bufType, d.numPlanes)
			if err != nil {
				if errors.Is(err, sys.EAGAIN) {
					continue
				}
				return fmt.Errorf("device: mplane output dequeue: %w", err)
			}
			idx = buf.Index
			break
		}
	}

	qPlanes := make([]v4l2.Plane, d.numPlanes)
	for p := range d.numPlanes {
		dst := d.buffersMPlane[idx][p]
		src := planes[p]
		if len(src) > len(dst) {
			return fmt.Errorf("device: plane %d size %d exceeds buffer length %d", p, len(src), len(dst))
		}
		copy(dst, src)
		pl := v4l2.Plane{
			BytesUsed: uint32(len(src)),
			Length:    uint32(len(dst)),
		}
		switch d.config.ioType {
		case v4l2.IOTypeUserPtr:
			if len(dst) > 0 {
				pl.Info.UserPtr = uintptr(unsafe.Pointer(&dst[0]))
			}
		case v4l2.IOTypeDMABuf:
			pl.Info.FD = d.dmabufFDs[int(idx)*int(d.numPlanes)+int(p)]
		}
		qPlanes[p] = pl
	}

	if _, err := v4l2.QueueBufferWith(d.fd, v4l2.Buffer{
		Index:  idx,
		Type:   d.bufType,
		Memory: d.config.ioType,
		Planes: qPlanes,
	}); err != nil {
		return fmt.Errorf("device: mplane output queue: %w", err)
	}
	if d.outputQueued < uint32(len(d.buffersMPlane)) {
		d.outputQueued++
	}

	if !d.outputStreamOn {
		if err := v4l2.StreamOn(d); err != nil {
			return fmt.Errorf("device: mplane output stream on: %w", err)
		}
		d.outputStreamOn = true
	}

	return nil
}
