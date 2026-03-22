package device

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	sys "syscall"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// Decoder represents a V4L2 stateful video decoder device.
// It provides a high-level interface for hardware-accelerated video decoding.
//
// Decoders accept compressed video (e.g., H.264, HEVC) on the input channel
// and produce raw video frames (e.g., NV12, YUV420) on the output channel.
//
// # Dynamic Resolution Changes
//
// V4L2 decoders can signal resolution changes mid-stream. When this happens:
//  1. The decoder signals via GetResolutionChanges() channel
//  2. Application should stop reading, handle the change, and continue
//
// # Usage
//
//	dec, err := device.OpenDecoder("/dev/video10", device.DecoderConfig{
//	    InputFormat: v4l2.PixFormat{
//	        PixelFormat: v4l2.PixelFmtH264,
//	    },
//	    OutputFormat: v4l2.PixFormat{
//	        PixelFormat: v4l2.PixelFmtNV12,
//	    },
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer dec.Close()
//
//	if err := dec.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Send encoded data
//	dec.GetInput() <- encodedData
//
//	// Receive decoded frames
//	frame := <-dec.GetOutput()
//
// # Concurrency Safety
//
// The Decoder is safe for concurrent use:
//   - GetInput() returns a channel for sending data (multiple writers OK)
//   - GetOutput() returns a channel for receiving frames (multiple readers OK)
//   - Start/Stop/Flush/Drain/Close should be called from a single goroutine
type Decoder struct {
	path         string
	fd           uintptr
	cap          v4l2.Capability
	stateMachine *v4l2.CodecStateMachine
	config       DecoderConfig

	// Input/output format (OUTPUT queue = encoded input, CAPTURE queue = raw output)
	inputFormat  v4l2.PixFormat
	outputFormat v4l2.PixFormat

	// Buffer management
	inputBuffers   [][]byte
	outputBuffers  [][]byte
	inputBufCount  uint32
	outputBufCount uint32

	// Channels
	inputChan     chan []byte
	outputChan    chan []byte
	errorChan     chan error
	resChangeChan chan ResolutionChange
	doneChan      chan struct{}

	// State
	streaming atomic.Bool
	mu        sync.Mutex
}

// ResolutionChange describes a detected resolution change in the stream.
type ResolutionChange struct {
	// Width is the new frame width
	Width uint32
	// Height is the new frame height
	Height uint32
	// PixelFormat is the new output pixel format
	PixelFormat v4l2.FourCCType
}

// DecoderConfig holds configuration for opening a decoder device.
type DecoderConfig struct {
	// InputFormat specifies the compressed video format (e.g., H.264, HEVC).
	// PixelFormat is required.
	InputFormat v4l2.PixFormat

	// OutputFormat specifies the raw video output format (e.g., NV12, I420).
	// PixelFormat defaults to NV12 if not set. Width/Height are set by decoder.
	OutputFormat v4l2.PixFormat

	// InputBufferCount is the number of input buffers to allocate.
	// Default is 8 if not specified (decoders typically need more input buffers).
	InputBufferCount uint32

	// OutputBufferCount is the number of output buffers to allocate.
	// Default is 4 if not specified.
	OutputBufferCount uint32
}

// OpenDecoder opens a V4L2 decoder device and configures it for decoding.
//
// Parameters:
//   - path: Device path (e.g., "/dev/video10")
//   - config: Decoder configuration
//
// Returns:
//   - *Decoder: Configured decoder ready for Start()
//   - error: If the device cannot be opened or configured
//
// The decoder must be started with Start() before sending data.
// Call Close() when done to release resources.
func OpenDecoder(path string, config DecoderConfig) (*Decoder, error) {
	// Open device
	fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK, 0)
	if err != nil {
		return nil, fmt.Errorf("decoder open: %w", err)
	}

	// Query capabilities
	cap, err := v4l2.GetCapability(fd)
	if err != nil {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("decoder open: query capability: %w", err)
	}

	// Verify this is a decoder device
	if !cap.IsDecoderSupported() {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("decoder open: device %s is not a decoder", path)
	}

	if !cap.IsStreamingSupported() {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("decoder open: device does not support streaming")
	}

	dec := &Decoder{
		path:           path,
		fd:             fd,
		cap:            cap,
		stateMachine:   v4l2.NewCodecStateMachine(v4l2.CodecTypeDecoder),
		config:         config,
		inputBufCount:  config.InputBufferCount,
		outputBufCount: config.OutputBufferCount,
	}

	// Set defaults
	if dec.inputBufCount == 0 {
		dec.inputBufCount = 8 // Decoders typically need more input buffers
	}
	if dec.outputBufCount == 0 {
		dec.outputBufCount = 4
	}

	// Configure input format
	if err := dec.configureInputFormat(); err != nil {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("decoder open: %w", err)
	}

	// Transition to initialized state
	if err := dec.stateMachine.Initialize(); err != nil {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("decoder open: %w", err)
	}

	return dec, nil
}

// configureInputFormat sets up the input format (compressed video).
func (dec *Decoder) configureInputFormat() error {
	// Set input format (OUTPUT queue - compressed video goes here)
	dec.inputFormat = dec.config.InputFormat

	if err := v4l2.SetPixFormatOutput(dec.fd, dec.inputFormat); err != nil {
		return fmt.Errorf("set input format: %w", err)
	}

	// Read back actual format
	actualInput, err := v4l2.GetPixFormatOutput(dec.fd)
	if err != nil {
		return fmt.Errorf("get input format: %w", err)
	}
	dec.inputFormat = actualInput

	return nil
}

// configureOutputFormat sets up the output format after resolution is known.
func (dec *Decoder) configureOutputFormat() error {
	// Set output format (CAPTURE queue - raw video comes from here)
	dec.outputFormat = dec.config.OutputFormat
	if dec.outputFormat.PixelFormat == 0 {
		dec.outputFormat.PixelFormat = v4l2.PixelFmtNV12
	}

	if err := v4l2.SetPixFormatCapture(dec.fd, dec.outputFormat); err != nil {
		return fmt.Errorf("set output format: %w", err)
	}

	// Read back actual format
	actualOutput, err := v4l2.GetPixFormatCapture(dec.fd)
	if err != nil {
		return fmt.Errorf("get output format: %w", err)
	}
	dec.outputFormat = actualOutput

	return nil
}

// Start begins decoding. Call this after OpenDecoder and before sending data.
//
// Parameters:
//   - ctx: Context for cancellation
//
// Returns:
//   - error: If decoding cannot be started
//
// After Start returns successfully:
//   - Send compressed data to GetInput()
//   - Receive decoded frames from GetOutput()
//   - Monitor GetResolutionChanges() for resolution changes
//   - Monitor GetError() for streaming errors
func (dec *Decoder) Start(ctx context.Context) error {
	dec.mu.Lock()
	defer dec.mu.Unlock()

	if dec.streaming.Load() {
		return fmt.Errorf("decoder: already streaming")
	}

	// Subscribe to codec events (important for resolution changes)
	if err := v4l2.SubscribeCodecEvents(dec.fd); err != nil {
		// Non-fatal - some drivers don't support events
	}

	// Allocate input buffers (OUTPUT queue)
	inputReq, err := v4l2.RequestBuffersOutput(dec.fd, dec.inputBufCount, v4l2.IOTypeMMAP)
	if err != nil {
		return fmt.Errorf("decoder start: request input buffers: %w", err)
	}
	dec.inputBufCount = inputReq.Count

	dec.inputBuffers, err = v4l2.MapMemoryBuffersOutput(dec.fd, dec.inputBufCount)
	if err != nil {
		return fmt.Errorf("decoder start: map input buffers: %w", err)
	}

	// Configure output format (may need initial data to determine resolution)
	if err := dec.configureOutputFormat(); err != nil {
		return fmt.Errorf("decoder start: %w", err)
	}

	// Allocate output buffers (CAPTURE queue)
	outputReq, err := v4l2.RequestBuffersCapture(dec.fd, dec.outputBufCount, v4l2.IOTypeMMAP)
	if err != nil {
		return fmt.Errorf("decoder start: request output buffers: %w", err)
	}
	dec.outputBufCount = outputReq.Count

	dec.outputBuffers, err = v4l2.MapMemoryBuffersCapture(dec.fd, dec.outputBufCount)
	if err != nil {
		return fmt.Errorf("decoder start: map output buffers: %w", err)
	}

	// Queue all output buffers
	for i := uint32(0); i < dec.outputBufCount; i++ {
		if _, err := v4l2.QueueBufferCapture(dec.fd, v4l2.IOTypeMMAP, i); err != nil {
			return fmt.Errorf("decoder start: queue output buffer %d: %w", i, err)
		}
	}

	// Create channels
	dec.inputChan = make(chan []byte, dec.inputBufCount)
	dec.outputChan = make(chan []byte, dec.outputBufCount)
	dec.errorChan = make(chan error, 8)
	dec.resChangeChan = make(chan ResolutionChange, 4)
	dec.doneChan = make(chan struct{})

	// Start streaming on both queues
	if err := v4l2.StreamOnOutput(dec.fd); err != nil {
		return fmt.Errorf("decoder start: stream on input: %w", err)
	}
	if err := v4l2.StreamOnCapture(dec.fd); err != nil {
		v4l2.StreamOffOutput(dec.fd)
		return fmt.Errorf("decoder start: stream on output: %w", err)
	}

	// Send decoder start command
	if err := v4l2.StartDecoder(dec.fd); err != nil {
		// Non-fatal - some drivers don't require explicit start
	}

	dec.streaming.Store(true)
	if err := dec.stateMachine.Start(); err != nil {
		return fmt.Errorf("decoder start: %w", err)
	}

	// Start decode loop
	go dec.decodeLoop(ctx)

	return nil
}

// decodeLoop handles the decoding process.
func (dec *Decoder) decodeLoop(ctx context.Context) {
	defer close(dec.doneChan)
	defer close(dec.outputChan)
	defer close(dec.resChangeChan)

	inputBufIdx := uint32(0)

	for {
		select {
		case <-ctx.Done():
			dec.stopInternal()
			return

		case data, ok := <-dec.inputChan:
			if !ok {
				// Input channel closed, start drain
				dec.drainInternal()
				return
			}

			// Copy data to input buffer
			if inputBufIdx < dec.inputBufCount && len(data) <= len(dec.inputBuffers[inputBufIdx]) {
				copy(dec.inputBuffers[inputBufIdx], data)

				// Queue input buffer with data
				if _, err := v4l2.QueueBufferOutput(dec.fd, v4l2.IOTypeMMAP, inputBufIdx, uint32(len(data))); err != nil {
					dec.sendError(fmt.Errorf("queue input buffer: %w", err))
					continue
				}

				inputBufIdx = (inputBufIdx + 1) % dec.inputBufCount
			}

			// Process output buffers
			dec.processOutputBuffers()

			// Check for events (resolution changes)
			dec.checkEvents()
		}
	}
}

// processOutputBuffers dequeues decoded frames and requeues buffers.
func (dec *Decoder) processOutputBuffers() {
	for {
		buf, err := v4l2.DequeueBufferCapture(dec.fd, v4l2.IOTypeMMAP)
		if err != nil {
			if v4l2.IsTemporaryError(err) {
				return // No more buffers ready
			}
			dec.sendError(fmt.Errorf("dequeue output buffer: %w", err))
			return
		}

		// Send decoded frame
		if buf.BytesUsed > 0 && buf.Index < uint32(len(dec.outputBuffers)) {
			data := make([]byte, buf.BytesUsed)
			copy(data, dec.outputBuffers[buf.Index][:buf.BytesUsed])

			select {
			case dec.outputChan <- data:
			default:
				// Output channel full, drop frame
			}
		}

		// Check for last buffer flag (EOS)
		if buf.Flags&v4l2.BufFlagLast != 0 {
			return
		}

		// Requeue buffer
		if _, err := v4l2.QueueBufferCapture(dec.fd, v4l2.IOTypeMMAP, buf.Index); err != nil {
			dec.sendError(fmt.Errorf("requeue output buffer: %w", err))
		}
	}
}

// checkEvents checks for codec events like resolution changes.
func (dec *Decoder) checkEvents() {
	event, err := v4l2.PollForEvent(dec.fd)
	if err != nil || event == nil {
		return
	}

	if v4l2.IsResolutionChangeEvent(event) {
		// Get new format
		newFormat, err := v4l2.GetPixFormatCapture(dec.fd)
		if err != nil {
			dec.sendError(fmt.Errorf("get new format after resolution change: %w", err))
			return
		}

		dec.outputFormat = newFormat

		// Notify application
		change := ResolutionChange{
			Width:       newFormat.Width,
			Height:      newFormat.Height,
			PixelFormat: newFormat.PixelFormat,
		}

		select {
		case dec.resChangeChan <- change:
		default:
			// Channel full
		}
	}
}

// sendError sends an error to the error channel if there's room.
func (dec *Decoder) sendError(err error) {
	select {
	case dec.errorChan <- err:
	default:
	}
}

// Flush flushes the decoder for seeking.
// This discards all pending data and resets the decoder state.
func (dec *Decoder) Flush() error {
	dec.mu.Lock()
	defer dec.mu.Unlock()

	if !dec.streaming.Load() {
		return nil
	}

	if err := dec.stateMachine.StartFlush(); err != nil {
		return err
	}

	if err := v4l2.FlushDecoder(dec.fd); err != nil {
		return fmt.Errorf("decoder flush: %w", err)
	}

	if err := dec.stateMachine.CompleteFlush(); err != nil {
		return err
	}

	return nil
}

// Drain signals the decoder to finish decoding all queued data.
// Call this before Stop() for a clean shutdown.
//
// After Drain():
//   - No new data should be sent to GetInput()
//   - Continue reading from GetOutput() until the channel closes
//   - The decoder will signal EOS when all data is decoded
func (dec *Decoder) Drain() error {
	dec.mu.Lock()
	defer dec.mu.Unlock()

	if !dec.streaming.Load() {
		return nil
	}

	return dec.drainInternal()
}

// drainInternal performs the drain operation.
func (dec *Decoder) drainInternal() error {
	if err := dec.stateMachine.StartDrain(); err != nil {
		return err
	}

	// Send stop command to decoder
	if err := v4l2.StopDecoder(dec.fd); err != nil {
		return fmt.Errorf("decoder drain: %w", err)
	}

	return nil
}

// Stop stops the decoder immediately without draining.
// Use Drain() first for a graceful shutdown.
func (dec *Decoder) Stop() error {
	dec.mu.Lock()
	defer dec.mu.Unlock()

	if !dec.streaming.Load() {
		return nil
	}

	return dec.stopInternal()
}

// stopInternal performs the stop operation.
func (dec *Decoder) stopInternal() error {
	dec.streaming.Store(false)

	// Stop streaming
	v4l2.StreamOffOutput(dec.fd)
	v4l2.StreamOffCapture(dec.fd)

	// Unsubscribe from events
	v4l2.UnsubscribeCodecEvents(dec.fd)

	// Unmap buffers
	if dec.inputBuffers != nil {
		for _, buf := range dec.inputBuffers {
			sys.Munmap(buf)
		}
		dec.inputBuffers = nil
	}
	if dec.outputBuffers != nil {
		for _, buf := range dec.outputBuffers {
			sys.Munmap(buf)
		}
		dec.outputBuffers = nil
	}

	// Close channels
	if dec.inputChan != nil {
		close(dec.inputChan)
		dec.inputChan = nil
	}
	if dec.errorChan != nil {
		close(dec.errorChan)
		dec.errorChan = nil
	}

	dec.stateMachine.Stop()

	return nil
}

// Close closes the decoder and releases all resources.
func (dec *Decoder) Close() error {
	if dec.streaming.Load() {
		dec.Stop()
	}

	// Wait for decode loop to finish
	if dec.doneChan != nil {
		<-dec.doneChan
	}

	dec.stateMachine.Uninitialize()

	return v4l2.CloseDevice(dec.fd)
}

// GetInput returns a channel for sending compressed data to decode.
// Data should be in the configured input format (e.g., H.264 NAL units).
func (dec *Decoder) GetInput() chan<- []byte {
	return dec.inputChan
}

// GetOutput returns a channel for receiving decoded frames.
// Frames are in the configured output format (e.g., NV12).
func (dec *Decoder) GetOutput() <-chan []byte {
	return dec.outputChan
}

// GetError returns a channel for receiving decoding errors.
func (dec *Decoder) GetError() <-chan error {
	return dec.errorChan
}

// GetResolutionChanges returns a channel for resolution change notifications.
// The decoder signals on this channel when the stream resolution changes.
func (dec *Decoder) GetResolutionChanges() <-chan ResolutionChange {
	return dec.resChangeChan
}

// GetState returns the current decoder state.
func (dec *Decoder) GetState() v4l2.CodecState {
	return dec.stateMachine.GetState()
}

// GetInputFormat returns the configured input format.
func (dec *Decoder) GetInputFormat() v4l2.PixFormat {
	return dec.inputFormat
}

// GetOutputFormat returns the current output format.
// This may change during decoding due to resolution changes.
func (dec *Decoder) GetOutputFormat() v4l2.PixFormat {
	return dec.outputFormat
}

// Fd returns the underlying file descriptor.
func (dec *Decoder) Fd() uintptr {
	return dec.fd
}

// Capability returns the device capabilities.
func (dec *Decoder) Capability() v4l2.Capability {
	return dec.cap
}
