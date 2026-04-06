package device

import (
	"errors"
	"fmt"
	sys "syscall"
	"time"
	"unsafe"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// startFramesCapture launches the optimized streaming loop for GetFrames() API.
// It captures frames from the device and sends them as Frame objects with pooled buffers.
// The loop runs in a separate goroutine and uses sys.Select to trigger capture events.
// Buffer queueing and StreamOn must have been completed by Start() before calling this.
func (d *Device) startFramesCapture() {
	// Initialize channels
	if d.frames == nil {
		d.frames = make(chan *Frame, d.config.bufSize)
	}
	if d.streamErr == nil {
		d.streamErr = make(chan error, 1)
	}

	ctx := d.startCtx

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
				// For USERPTR, buffers are always valid (app-allocated)
				isMapped := (ioMemType == v4l2.IOTypeUserPtr) || (buff.Flags&v4l2.BufFlagMapped != 0)
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

				var queueErr error
				switch ioMemType {
				case v4l2.IOTypeUserPtr:
					ptr := uintptr(unsafe.Pointer(&d.buffers[buff.Index][0]))
					_, queueErr = v4l2.QueueBufferUserPtr(fd, bufType, buff.Index, ptr, uint32(len(d.buffers[buff.Index])))
				case v4l2.IOTypeDMABuf:
					_, queueErr = v4l2.QueueBufferDMABuf(fd, bufType, buff.Index, d.dmabufFDs[buff.Index], uint32(len(d.buffers[buff.Index])))
				default:
					_, queueErr = v4l2.QueueBuffer(fd, ioMemType, bufType, buff.Index)
				}
				if queueErr != nil {
					select {
					case d.streamErr <- fmt.Errorf("device: stream loop queue: %w: buff: %#v", queueErr, buff):
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
}
