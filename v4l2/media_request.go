package v4l2

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// MediaRequest represents a Linux Media Request file descriptor allocated via
// MEDIA_IOC_REQUEST_ALLOC. Requests synchronize per-frame buffers (via Buffer.RequestFD
// and BufFlagRequestFD) with per-frame extended controls (via ExtControls.SetRequestFD)
// for stateless hardware codecs and complex camera pipelines.
type MediaRequest struct {
	fd int32
}

// AllocMediaRequest allocates a new Media Request file descriptor on a Media Controller
// device (e.g., /dev/media0) via MEDIA_IOC_REQUEST_ALLOC.
//
// Call Close when the request is no longer needed.
func AllocMediaRequest(mediaFD uintptr) (*MediaRequest, error) {
	var reqFD int32
	if err := send(mediaFD, cMEDIA_IOC_REQUEST_ALLOC, unsafe.Pointer(&reqFD)); err != nil {
		return nil, fmt.Errorf("media alloc request: %w", err)
	}
	return &MediaRequest{fd: reqFD}, nil
}

// Fd returns the underlying request file descriptor (suitable for Buffer.RequestFD,
// ExtControls.SetRequestFD, or Poller.Add). Returns -1 if the request has been closed.
func (r *MediaRequest) Fd() int32 {
	if r == nil {
		return -1
	}
	return r.fd
}

// Queue atomically submits all buffers and extended controls bound to this request
// for execution via MEDIA_REQUEST_IOC_QUEUE.
func (r *MediaRequest) Queue() error {
	if r == nil || r.fd < 0 {
		return fmt.Errorf("media queue request: request is closed")
	}
	if err := send(uintptr(r.fd), cMEDIA_REQUEST_IOC_QUEUE, nil); err != nil {
		return fmt.Errorf("media queue request (fd %d): %w", r.fd, err)
	}
	return nil
}

// Reinit resets a completed Media Request so it can be reused for another frame
// without allocating a new file descriptor (MEDIA_REQUEST_IOC_REINIT).
func (r *MediaRequest) Reinit() error {
	if r == nil || r.fd < 0 {
		return fmt.Errorf("media reinit request: request is closed")
	}
	if err := send(uintptr(r.fd), cMEDIA_REQUEST_IOC_REINIT, nil); err != nil {
		return fmt.Errorf("media reinit request (fd %d): %w", r.fd, err)
	}
	return nil
}

// Wait blocks until the queued Media Request signals completion (POLLPRI),
// timeout expires, or ctx is cancelled.
func (r *MediaRequest) Wait(ctx context.Context, timeout time.Duration) (PollEvent, error) {
	if r == nil || r.fd < 0 {
		return 0, fmt.Errorf("media wait request: request is closed")
	}
	return Poll(ctx, uintptr(r.fd), PollPri, timeout)
}

// Close closes the underlying Media Request file descriptor. Safe to call multiple times.
func (r *MediaRequest) Close() error {
	if r == nil || r.fd < 0 {
		return nil
	}
	fd := r.fd
	r.fd = -1
	if err := sys.Close(int(fd)); err != nil {
		return fmt.Errorf("media close request (fd %d): %w", fd, err)
	}
	return nil
}
