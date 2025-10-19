package device

import (
	"context"
	"errors"
	"fmt"
	sys "syscall"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// captureFrames implements the optimized streaming loop for GetFrames() API.
// It captures frames from the device and sends them as Frame objects with pooled buffers.
// The loop runs in a separate goroutine and uses sys.Select to trigger capture events.
func (d *Device) captureFrames(ctx context.Context) error {
	// Initialize channels
	if d.frames == nil {
		d.frames = make(chan *Frame, d.config.bufSize)
	}
	if d.streamErr == nil {
		d.streamErr = make(chan error, 1)
	}

	// Initial enqueue of buffers for capture
	for i := 0; i < int(d.config.bufSize); i++ {
		_, err := v4l2.QueueBuffer(d.fd, d.config.ioType, d.bufType, uint32(i))
		if err != nil {
			return fmt.Errorf("device: buffer queueing: %w", err)
		}
	}

	if err := v4l2.StreamOn(d); err != nil {
		return fmt.Errorf("device: stream on: %w", err)
	}

	go func() {
		defer close(d.captureDone) // Signal Stop() that goroutine has exited
		defer close(d.frames)
		defer close(d.streamErr)
		defer func() {
			// Mark streaming as stopped when goroutine exits
			d.streaming.Store(false)
			d.streamingMode.Store(0) // Reset mode
		}()

		fd := d.Fd()
		ioMemType := d.MemIOType()
		bufType := d.BufferType()
		waitForRead := v4l2.WaitForRead(ctx, d)
		for {
			select {
			// handle stream capture (read from driver)
			case <-waitForRead:
				buff, err := v4l2.DequeueBuffer(fd, ioMemType, bufType)
				if err != nil {
					if errors.Is(err, sys.EAGAIN) {
						continue
					}
					// Send error and exit gracefully
					select {
					case d.streamErr <- fmt.Errorf("device: stream loop dequeue: %w", err):
					default:
					}
					return
				}

				// Process buffer based on its state
				isMapped := buff.Flags&v4l2.BufFlagMapped != 0
				hasError := buff.Flags&v4l2.BufFlagError != 0
				hasData := buff.BytesUsed > 0

				switch {
				case isMapped && !hasError && hasData:
					// Safety check: ensure we're still streaming before accessing buffers
					// If Stop() was called, buffers may have been unmapped
					if !d.streaming.Load() {
						return
					}

					// Buffer has valid frame data - use pooled buffer
					poolBuf := d.framePool.Get(buff.BytesUsed)

					// Safety check: ensure buffers are still valid before accessing
					if buff.Index >= uint32(len(d.buffers)) || d.buffers[buff.Index] == nil {
						// Buffers have been unmapped, exit gracefully
						d.framePool.Put(poolBuf)
						return
					}

					copy(poolBuf, d.buffers[buff.Index][:buff.BytesUsed])

					// Create Frame object with metadata
					frameObj := &Frame{
						Data:      poolBuf,
						Timestamp: time.Unix(int64(buff.Timestamp.Sec), int64(buff.Timestamp.Usec)*1000),
						Sequence:  buff.Sequence,
						Flags:     buff.Flags,
						Index:     buff.Index,
						pool:      d.framePool,
						released:  false,
					}

					select {
					case d.frames <- frameObj:
						// Frame delivered successfully
					default:
						// Consumer too slow, release buffer back to pool
						frameObj.Release()
						select {
						case d.streamErr <- fmt.Errorf("device: frame dropped (consumer too slow): %d bytes", buff.BytesUsed):
						default:
						}
					}

				case hasError:
					// Buffer has error flag
					select {
					case d.streamErr <- fmt.Errorf("device: buffer error flag: index=%d flags=0x%x", buff.Index, buff.Flags):
					default:
					}

				// Other cases (no data, not mapped) are silently skipped
				}

				if _, err := v4l2.QueueBuffer(fd, ioMemType, bufType, buff.Index); err != nil {
					// Send error and exit gracefully
					select {
					case d.streamErr <- fmt.Errorf("device: stream loop queue: %w: buff: %#v", err, buff):
					default:
					}
					return
				}
			case <-ctx.Done():
				// Context cancelled, exit gracefully
				return
			}
		}
	}()

	return nil
}
