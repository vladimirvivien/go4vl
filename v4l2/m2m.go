package v4l2

import (
	"context"
	"errors"
	"fmt"
	"time"

	sys "golang.org/x/sys/unix"
)

// M2MQueueConfig configures either the OUTPUT or CAPTURE queue of an M2MSession.
type M2MQueueConfig struct {
	// BufType selects the queue type (e.g., BufTypeVideoOutput, BufTypeVideoOutputMPlane,
	// BufTypeVideoCapture, BufTypeVideoCaptureMPlane).
	BufType BufType

	// Memory selects the buffer memory type (IOTypeMMAP, IOTypeUserPtr, or IOTypeDMABuf).
	// Defaults to IOTypeMMAP when zero.
	Memory IOType

	// Count is the requested number of buffers for the queue.
	Count uint32

	// NumPlanes is the number of planes per buffer for multi-planar buffer types.
	// Ignored for single-planar buffer types; defaults to 1 if zero on multi-planar queues.
	NumPlanes uint32
}

// M2MQueue tracks the allocated buffers and streaming state of one side (OUTPUT or CAPTURE)
// of a V4L2 memory-to-memory device.
type M2MQueue struct {
	BufType   BufType
	Memory    IOType
	Count     uint32
	NumPlanes uint32
	// Planes holds mapped MMAP slices indexed as Planes[bufIndex][planeIndex].
	// For single-planar queues, each buffer has a single slice at plane index 0.
	Planes    [][][]byte
	Streaming bool
}

// M2MSession manages a dual-queue (OUTPUT + CAPTURE) V4L2 memory-to-memory device
// such as a hardware scaler, color-space converter, ISP, or video codec.
//
// It wraps the raw device file descriptor without spawning background goroutines,
// exposing Fd() for integration with Poller and direct Buffer / DMA-BUF control.
type M2MSession struct {
	fd      uintptr
	ownedFD bool
	cap     Capability
	Output  M2MQueue
	Capture M2MQueue
}

// OpenM2M opens a V4L2 memory-to-memory device node (e.g., "/dev/video10") in
// non-blocking read-write mode and verifies M2M capability.
func OpenM2M(path string) (*M2MSession, error) {
	fd, err := OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK|sys.O_CLOEXEC, 0)
	if err != nil {
		return nil, fmt.Errorf("m2m open %s: %w", path, err)
	}
	sess, err := NewM2M(fd)
	if err != nil {
		_ = CloseDevice(fd)
		return nil, fmt.Errorf("m2m open %s: %w", path, err)
	}
	sess.ownedFD = true
	return sess, nil
}

// NewM2M wraps an already-opened V4L2 M2M file descriptor. Close() will stop
// active queues and release buffers but will not close fd.
func NewM2M(fd uintptr) (*M2MSession, error) {
	cap, err := GetCapability(fd)
	if err != nil {
		return nil, fmt.Errorf("m2m capability: %w", err)
	}
	caps := cap.GetCapabilities()
	if caps&CapVideoMem2Mem == 0 && caps&CapVideoMem2MemMPlane == 0 {
		return nil, fmt.Errorf("m2m: device does not support memory-to-memory: %w", ErrorUnsupportedFeature)
	}
	return &M2MSession{
		fd:  fd,
		cap: cap,
	}, nil
}

// Fd returns the underlying V4L2 device file descriptor for use with Poller or raw ioctls.
func (m *M2MSession) Fd() uintptr {
	return m.fd
}

// Capability returns the queried V4L2 capability of the M2M device.
func (m *M2MSession) Capability() Capability {
	return m.cap
}

// SetupOutput allocates buffers (and maps MMAP planes) for the OUTPUT queue.
func (m *M2MSession) SetupOutput(cfg M2MQueueConfig) error {
	if cfg.BufType == 0 {
		if m.cap.GetCapabilities()&CapVideoMem2MemMPlane != 0 {
			cfg.BufType = BufTypeVideoOutputMPlane
		} else {
			cfg.BufType = BufTypeVideoOutput
		}
	}
	if cfg.BufType != BufTypeVideoOutput && cfg.BufType != BufTypeVideoOutputMPlane {
		return fmt.Errorf("m2m setup output: invalid buffer type %d", cfg.BufType)
	}
	return m.setupQueue(&m.Output, cfg)
}

// SetupCapture allocates buffers (and maps MMAP planes) for the CAPTURE queue.
func (m *M2MSession) SetupCapture(cfg M2MQueueConfig) error {
	if cfg.BufType == 0 {
		if m.cap.GetCapabilities()&CapVideoMem2MemMPlane != 0 {
			cfg.BufType = BufTypeVideoCaptureMPlane
		} else {
			cfg.BufType = BufTypeVideoCapture
		}
	}
	if cfg.BufType != BufTypeVideoCapture && cfg.BufType != BufTypeVideoCaptureMPlane {
		return fmt.Errorf("m2m setup capture: invalid buffer type %d", cfg.BufType)
	}
	return m.setupQueue(&m.Capture, cfg)
}

func (m *M2MSession) setupQueue(q *M2MQueue, cfg M2MQueueConfig) error {
	if cfg.Memory == 0 {
		cfg.Memory = IOTypeMMAP
	}
	if cfg.Count == 0 {
		cfg.Count = 2
	}
	numPlanes := uint32(1)
	if IsMultiPlanar(cfg.BufType) {
		if cfg.NumPlanes > 0 {
			numPlanes = cfg.NumPlanes
		}
		if numPlanes > MaxPlanes {
			return fmt.Errorf("m2m setup queue: invalid plane count %d", numPlanes)
		}
	}

	req, err := RequestBuffersForType(m.fd, cfg.BufType, cfg.Memory, cfg.Count)
	if err != nil {
		return fmt.Errorf("m2m setup queue: %w", err)
	}

	q.BufType = cfg.BufType
	q.Memory = cfg.Memory
	q.Count = req.Count
	q.NumPlanes = numPlanes
	q.Planes = nil

	if cfg.Memory == IOTypeMMAP {
		planes := make([][][]byte, req.Count)
		for i := uint32(0); i < req.Count; i++ {
			buf, err := QueryBuffer(m.fd, cfg.BufType, cfg.Memory, i, numPlanes)
			if err != nil {
				_ = m.unmapQueuePlanes(planes)
				_, _ = RequestBuffersForType(m.fd, cfg.BufType, cfg.Memory, 0)
				return fmt.Errorf("m2m setup queue: query buffer %d: %w", i, err)
			}
			if IsMultiPlanar(cfg.BufType) {
				bufPlanes := make([][]byte, len(buf.Planes))
				for p, pl := range buf.Planes {
					mapped, err := MapMemoryBuffer(m.fd, int64(pl.Info.MemOffset), int(pl.Length))
					if err != nil {
						planes[i] = bufPlanes
						_ = m.unmapQueuePlanes(planes)
						_, _ = RequestBuffersForType(m.fd, cfg.BufType, cfg.Memory, 0)
						return fmt.Errorf("m2m setup queue: map buffer %d plane %d: %w", i, p, err)
					}
					bufPlanes[p] = mapped
				}
				planes[i] = bufPlanes
			} else {
				mapped, err := MapMemoryBuffer(m.fd, int64(buf.Info.Offset), int(buf.Length))
				if err != nil {
					_ = m.unmapQueuePlanes(planes)
					_, _ = RequestBuffersForType(m.fd, cfg.BufType, cfg.Memory, 0)
					return fmt.Errorf("m2m setup queue: map buffer %d: %w", i, err)
				}
				planes[i] = [][]byte{mapped}
			}
		}
		q.Planes = planes
	}

	return nil
}

func (m *M2MSession) unmapQueuePlanes(planes [][][]byte) error {
	var errs []error
	for _, bufPlanes := range planes {
		for _, slice := range bufPlanes {
			if slice != nil {
				if err := UnmapMemoryBuffer(slice); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	return errors.Join(errs...)
}

// StartOutput starts streaming on the OUTPUT queue (VIDIOC_STREAMON).
func (m *M2MSession) StartOutput() error {
	if m.Output.BufType == 0 {
		return fmt.Errorf("m2m start output: queue not initialized")
	}
	if err := StreamOnType(m.fd, m.Output.BufType); err != nil {
		return err
	}
	m.Output.Streaming = true
	return nil
}

// StartCapture starts streaming on the CAPTURE queue (VIDIOC_STREAMON).
func (m *M2MSession) StartCapture() error {
	if m.Capture.BufType == 0 {
		return fmt.Errorf("m2m start capture: queue not initialized")
	}
	if err := StreamOnType(m.fd, m.Capture.BufType); err != nil {
		return err
	}
	m.Capture.Streaming = true
	return nil
}

// StopOutput stops streaming on the OUTPUT queue (VIDIOC_STREAMOFF).
func (m *M2MSession) StopOutput() error {
	if !m.Output.Streaming {
		return nil
	}
	err := StreamOffType(m.fd, m.Output.BufType)
	m.Output.Streaming = false
	return err
}

// StopCapture stops streaming on the CAPTURE queue (VIDIOC_STREAMOFF).
func (m *M2MSession) StopCapture() error {
	if !m.Capture.Streaming {
		return nil
	}
	err := StreamOffType(m.fd, m.Capture.BufType)
	m.Capture.Streaming = false
	return err
}

// QueueOutput enqueues a buffer on the OUTPUT queue. Type and Memory default to
// the configured OUTPUT queue settings when left zero.
func (m *M2MSession) QueueOutput(buf Buffer) (Buffer, error) {
	if buf.Type == 0 {
		buf.Type = m.Output.BufType
	}
	if buf.Memory == 0 {
		buf.Memory = m.Output.Memory
	}
	return QueueBufferWith(m.fd, buf)
}

// QueueCapture enqueues a buffer on the CAPTURE queue. Type and Memory default to
// the configured CAPTURE queue settings when left zero.
func (m *M2MSession) QueueCapture(buf Buffer) (Buffer, error) {
	if buf.Type == 0 {
		buf.Type = m.Capture.BufType
	}
	if buf.Memory == 0 {
		buf.Memory = m.Capture.Memory
	}
	return QueueBufferWith(m.fd, buf)
}

// DequeueOutput dequeues a processed buffer from the OUTPUT queue.
func (m *M2MSession) DequeueOutput() (Buffer, error) {
	return m.dequeueQueue(&m.Output)
}

// DequeueCapture dequeues a filled buffer from the CAPTURE queue.
func (m *M2MSession) DequeueCapture() (Buffer, error) {
	return m.dequeueQueue(&m.Capture)
}

func (m *M2MSession) dequeueQueue(q *M2MQueue) (Buffer, error) {
	if IsMultiPlanar(q.BufType) {
		return DequeueBufferMPlane(m.fd, q.Memory, q.BufType, q.NumPlanes)
	}
	return DequeueBuffer(m.fd, q.Memory, q.BufType)
}

// ExportOutputDMABuf exports an OUTPUT queue buffer plane as a DMA-BUF file descriptor.
func (m *M2MSession) ExportOutputDMABuf(index, plane, flags uint32) (int32, error) {
	return ExportBufferPlane(m.fd, m.Output.BufType, index, plane, flags)
}

// ExportCaptureDMABuf exports a CAPTURE queue buffer plane as a DMA-BUF file descriptor.
func (m *M2MSession) ExportCaptureDMABuf(index, plane, flags uint32) (int32, error) {
	return ExportBufferPlane(m.fd, m.Capture.BufType, index, plane, flags)
}

// SubscribeSourceChange subscribes to V4L2_EVENT_SOURCE_CHANGE on the M2M device.
func (m *M2MSession) SubscribeSourceChange() error {
	return SubscribeEvent(m.fd, NewEventSubscription(EventSourceChange))
}

// SubscribeEOS subscribes to V4L2_EVENT_EOS on the M2M device.
func (m *M2MSession) SubscribeEOS() error {
	return SubscribeEvent(m.fd, NewEventSubscription(EventEOS))
}

// DequeueEvent retrieves a pending V4L2 event from the M2M device.
func (m *M2MSession) DequeueEvent() (*Event, error) {
	return DequeueEvent(m.fd)
}

// Wait polls the M2M file descriptor for any of the requested readiness events
// (typically PollIn | PollOut | PollPri).
func (m *M2MSession) Wait(ctx context.Context, events PollEvent, timeout time.Duration) (PollEvent, error) {
	return Poll(ctx, m.fd, events, timeout)
}

// Close stops both queues, unmaps any MMAP planes, releases kernel buffer
// allocations, and closes the device file descriptor if opened via OpenM2M.
func (m *M2MSession) Close() error {
	var errs []error
	if err := m.StopOutput(); err != nil {
		errs = append(errs, err)
	}
	if err := m.StopCapture(); err != nil {
		errs = append(errs, err)
	}
	if err := m.releaseQueue(&m.Output); err != nil {
		errs = append(errs, err)
	}
	if err := m.releaseQueue(&m.Capture); err != nil {
		errs = append(errs, err)
	}
	if m.ownedFD && m.fd != 0 {
		if err := CloseDevice(m.fd); err != nil {
			errs = append(errs, err)
		}
		m.fd = 0
	}
	return errors.Join(errs...)
}

func (m *M2MSession) releaseQueue(q *M2MQueue) error {
	var errs []error
	if len(q.Planes) > 0 {
		if err := m.unmapQueuePlanes(q.Planes); err != nil {
			errs = append(errs, err)
		}
		q.Planes = nil
	}
	if q.BufType != 0 && q.Count > 0 {
		if _, err := RequestBuffersForType(m.fd, q.BufType, q.Memory, 0); err != nil {
			errs = append(errs, err)
		}
		q.Count = 0
	}
	return errors.Join(errs...)
}
