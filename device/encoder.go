package device

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	sys "syscall"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// Encoder represents a V4L2 stateful video encoder device.
// It provides a high-level interface for hardware-accelerated video encoding.
//
// Encoders accept raw video frames (e.g., NV12, YUV420) on the input channel
// and produce compressed video (e.g., H.264, HEVC) on the output channel.
//
// # Usage
//
//	enc, err := device.OpenEncoder("/dev/video11", device.EncoderConfig{
//	    InputFormat: v4l2.PixFormat{
//	        Width:       1920,
//	        Height:      1080,
//	        PixelFormat: v4l2.PixelFmtNV12,
//	    },
//	    OutputFormat: v4l2.PixFormat{
//	        PixelFormat: v4l2.PixelFmtH264,
//	    },
//	    Bitrate: 4000000, // 4 Mbps
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer enc.Close()
//
//	if err := enc.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Send frames to encode
//	enc.GetInput() <- rawFrame
//
//	// Receive encoded data
//	encodedData := <-enc.GetOutput()
//
// # Concurrency Safety
//
// The Encoder is safe for concurrent use:
//   - GetInput() returns a channel for sending frames (multiple writers OK)
//   - GetOutput() returns a channel for receiving encoded data (multiple readers OK)
//   - Start/Stop/Drain/Close should be called from a single goroutine
type Encoder struct {
	path         string
	fd           uintptr
	cap          v4l2.Capability
	stateMachine *v4l2.CodecStateMachine
	config       EncoderConfig

	// Input/output format (OUTPUT queue = raw input, CAPTURE queue = encoded output)
	inputFormat  v4l2.PixFormat
	outputFormat v4l2.PixFormat

	// Buffer management
	inputBuffers  [][]byte
	outputBuffers [][]byte
	inputBufCount uint32
	outputBufCount uint32

	// Channels
	inputChan   chan []byte
	outputChan  chan []byte
	errorChan   chan error
	doneChan    chan struct{}

	// State
	streaming atomic.Bool
	mu        sync.Mutex
}

// EncoderConfig holds configuration for opening an encoder device.
type EncoderConfig struct {
	// InputFormat specifies the raw video format (e.g., NV12, I420).
	// Width and Height are required. PixelFormat defaults to NV12 if not set.
	InputFormat v4l2.PixFormat

	// OutputFormat specifies the compressed video format (e.g., H.264, HEVC).
	// PixelFormat is required. Width/Height are inherited from InputFormat.
	OutputFormat v4l2.PixFormat

	// Bitrate in bits per second (e.g., 4000000 for 4 Mbps).
	// If zero, the driver's default is used.
	Bitrate uint32

	// GOPSize is the number of frames between keyframes.
	// If zero, the driver's default is used.
	GOPSize uint32

	// Profile specifies the encoding profile (codec-specific).
	// For H.264: 66=Baseline, 77=Main, 100=High
	Profile uint32

	// Level specifies the encoding level (codec-specific).
	Level uint32

	// InputBufferCount is the number of input buffers to allocate.
	// Default is 4 if not specified.
	InputBufferCount uint32

	// OutputBufferCount is the number of output buffers to allocate.
	// Default is 4 if not specified.
	OutputBufferCount uint32
}

// OpenEncoder opens a V4L2 encoder device and configures it for encoding.
//
// Parameters:
//   - path: Device path (e.g., "/dev/video11")
//   - config: Encoder configuration
//
// Returns:
//   - *Encoder: Configured encoder ready for Start()
//   - error: If the device cannot be opened or configured
//
// The encoder must be started with Start() before sending frames.
// Call Close() when done to release resources.
func OpenEncoder(path string, config EncoderConfig) (*Encoder, error) {
	// Open device
	fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_NONBLOCK, 0)
	if err != nil {
		return nil, fmt.Errorf("encoder open: %w", err)
	}

	// Query capabilities
	cap, err := v4l2.GetCapability(fd)
	if err != nil {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("encoder open: query capability: %w", err)
	}

	// Verify this is an encoder device
	if !cap.IsEncoderSupported() {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("encoder open: device %s is not an encoder", path)
	}

	if !cap.IsStreamingSupported() {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("encoder open: device does not support streaming")
	}

	enc := &Encoder{
		path:           path,
		fd:             fd,
		cap:            cap,
		stateMachine:   v4l2.NewCodecStateMachine(v4l2.CodecTypeEncoder),
		config:         config,
		inputBufCount:  config.InputBufferCount,
		outputBufCount: config.OutputBufferCount,
	}

	// Set defaults
	if enc.inputBufCount == 0 {
		enc.inputBufCount = 4
	}
	if enc.outputBufCount == 0 {
		enc.outputBufCount = 4
	}

	// Configure formats
	if err := enc.configureFormats(); err != nil {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("encoder open: %w", err)
	}

	// Transition to initialized state
	if err := enc.stateMachine.Initialize(); err != nil {
		v4l2.CloseDevice(fd)
		return nil, fmt.Errorf("encoder open: %w", err)
	}

	return enc, nil
}

// configureFormats sets up input and output formats.
func (enc *Encoder) configureFormats() error {
	// Set input format (OUTPUT queue - raw video goes here)
	enc.inputFormat = enc.config.InputFormat
	if enc.inputFormat.PixelFormat == 0 {
		enc.inputFormat.PixelFormat = v4l2.PixelFmtNV12
	}

	if err := v4l2.SetPixFormatOutput(enc.fd, enc.inputFormat); err != nil {
		return fmt.Errorf("set input format: %w", err)
	}

	// Read back actual format
	actualInput, err := v4l2.GetPixFormatOutput(enc.fd)
	if err != nil {
		return fmt.Errorf("get input format: %w", err)
	}
	enc.inputFormat = actualInput

	// Set output format (CAPTURE queue - encoded data comes from here)
	enc.outputFormat = enc.config.OutputFormat
	if enc.outputFormat.Width == 0 {
		enc.outputFormat.Width = enc.inputFormat.Width
	}
	if enc.outputFormat.Height == 0 {
		enc.outputFormat.Height = enc.inputFormat.Height
	}

	if err := v4l2.SetPixFormat(enc.fd, enc.outputFormat); err != nil {
		return fmt.Errorf("set output format: %w", err)
	}

	// Read back actual format
	actualOutput, err := v4l2.GetPixFormat(enc.fd)
	if err != nil {
		return fmt.Errorf("get output format: %w", err)
	}
	enc.outputFormat = actualOutput

	return nil
}

// Start begins encoding. Call this after OpenEncoder and before sending frames.
//
// Parameters:
//   - ctx: Context for cancellation
//
// Returns:
//   - error: If encoding cannot be started
//
// After Start returns successfully:
//   - Send raw frames to GetInput()
//   - Receive encoded data from GetOutput()
//   - Monitor GetError() for streaming errors
func (enc *Encoder) Start(ctx context.Context) error {
	enc.mu.Lock()
	defer enc.mu.Unlock()

	if enc.streaming.Load() {
		return fmt.Errorf("encoder: already streaming")
	}

	// Subscribe to codec events
	if err := v4l2.SubscribeCodecEvents(enc.fd); err != nil {
		// Non-fatal - some drivers don't support events
	}

	// Allocate and map input buffers (OUTPUT queue)
	inputReq, err := v4l2.RequestBuffersOutput(enc.fd, enc.inputBufCount, v4l2.IOTypeMMAP)
	if err != nil {
		return fmt.Errorf("encoder start: request input buffers: %w", err)
	}
	enc.inputBufCount = inputReq.Count

	enc.inputBuffers, err = v4l2.MapMemoryBuffersOutput(enc.fd, enc.inputBufCount)
	if err != nil {
		return fmt.Errorf("encoder start: map input buffers: %w", err)
	}

	// Allocate and map output buffers (CAPTURE queue)
	outputReq, err := v4l2.RequestBuffersCapture(enc.fd, enc.outputBufCount, v4l2.IOTypeMMAP)
	if err != nil {
		return fmt.Errorf("encoder start: request output buffers: %w", err)
	}
	enc.outputBufCount = outputReq.Count

	enc.outputBuffers, err = v4l2.MapMemoryBuffersCapture(enc.fd, enc.outputBufCount)
	if err != nil {
		return fmt.Errorf("encoder start: map output buffers: %w", err)
	}

	// Queue all output buffers
	for i := uint32(0); i < enc.outputBufCount; i++ {
		if _, err := v4l2.QueueBuffer(enc.fd, v4l2.IOTypeMMAP, v4l2.BufTypeVideoCapture, i); err != nil {
			return fmt.Errorf("encoder start: queue output buffer %d: %w", i, err)
		}
	}

	// Create channels
	enc.inputChan = make(chan []byte, enc.inputBufCount)
	enc.outputChan = make(chan []byte, enc.outputBufCount)
	enc.errorChan = make(chan error, 8)
	enc.doneChan = make(chan struct{})

	// Start streaming on both queues
	if err := v4l2.StreamOnOutput(enc.fd); err != nil {
		return fmt.Errorf("encoder start: stream on input: %w", err)
	}
	if err := v4l2.StreamOnCapture(enc.fd); err != nil {
		v4l2.StreamOffOutput(enc.fd)
		return fmt.Errorf("encoder start: stream on output: %w", err)
	}

	// Send encoder start command
	if err := v4l2.StartEncoder(enc.fd); err != nil {
		// Non-fatal - some drivers don't require explicit start
	}

	enc.streaming.Store(true)
	if err := enc.stateMachine.Start(); err != nil {
		return fmt.Errorf("encoder start: %w", err)
	}

	// Start encode loop
	go enc.encodeLoop(ctx)

	return nil
}

// encodeLoop handles the encoding process.
func (enc *Encoder) encodeLoop(ctx context.Context) {
	defer close(enc.doneChan)
	defer close(enc.outputChan)

	inputBufIdx := uint32(0)

	for {
		select {
		case <-ctx.Done():
			enc.stopInternal()
			return

		case frame, ok := <-enc.inputChan:
			if !ok {
				// Input channel closed, start drain
				enc.drainInternal()
				return
			}

			// Copy frame to input buffer
			if inputBufIdx < enc.inputBufCount {
				copy(enc.inputBuffers[inputBufIdx], frame)

				// Queue input buffer with data
				if _, err := v4l2.QueueBufferOutput(enc.fd, v4l2.IOTypeMMAP, inputBufIdx, uint32(len(frame))); err != nil {
					enc.sendError(fmt.Errorf("queue input buffer: %w", err))
					continue
				}

				inputBufIdx = (inputBufIdx + 1) % enc.inputBufCount
			}

			// Try to dequeue and requeue output buffers
			enc.processOutputBuffers()
		}
	}
}

// processOutputBuffers dequeues encoded data and requeues buffers.
func (enc *Encoder) processOutputBuffers() {
	for {
		buf, err := v4l2.DequeueBuffer(enc.fd, v4l2.IOTypeMMAP, v4l2.BufTypeVideoCapture)
		if err != nil {
			if v4l2.IsTemporaryError(err) {
				return // No more buffers ready
			}
			enc.sendError(fmt.Errorf("dequeue output buffer: %w", err))
			return
		}

		// Send encoded data
		if buf.BytesUsed > 0 && buf.Index < uint32(len(enc.outputBuffers)) {
			data := make([]byte, buf.BytesUsed)
			copy(data, enc.outputBuffers[buf.Index][:buf.BytesUsed])

			select {
			case enc.outputChan <- data:
			default:
				// Output channel full, drop frame
			}
		}

		// Check for last buffer flag (EOS)
		if buf.Flags&v4l2.BufFlagLast != 0 {
			return
		}

		// Requeue buffer
		if _, err := v4l2.QueueBuffer(enc.fd, v4l2.IOTypeMMAP, v4l2.BufTypeVideoCapture, buf.Index); err != nil {
			enc.sendError(fmt.Errorf("requeue output buffer: %w", err))
		}
	}
}

// sendError sends an error to the error channel if there's room.
func (enc *Encoder) sendError(err error) {
	select {
	case enc.errorChan <- err:
	default:
	}
}

// Drain signals the encoder to finish encoding all queued frames.
// Call this before Stop() for a clean shutdown.
//
// After Drain():
//   - No new frames should be sent to GetInput()
//   - Continue reading from GetOutput() until the channel closes
//   - The encoder will signal EOS when all frames are encoded
func (enc *Encoder) Drain() error {
	enc.mu.Lock()
	defer enc.mu.Unlock()

	if !enc.streaming.Load() {
		return nil
	}

	return enc.drainInternal()
}

// drainInternal performs the drain operation (must hold mutex).
func (enc *Encoder) drainInternal() error {
	if err := enc.stateMachine.StartDrain(); err != nil {
		return err
	}

	// Send stop command to encoder
	if err := v4l2.StopEncoder(enc.fd); err != nil {
		return fmt.Errorf("encoder drain: %w", err)
	}

	return nil
}

// Stop stops the encoder immediately without draining.
// Use Drain() first for a graceful shutdown.
func (enc *Encoder) Stop() error {
	enc.mu.Lock()
	defer enc.mu.Unlock()

	if !enc.streaming.Load() {
		return nil
	}

	return enc.stopInternal()
}

// stopInternal performs the stop operation (must hold mutex or be called from encodeLoop).
func (enc *Encoder) stopInternal() error {
	enc.streaming.Store(false)

	// Stop streaming
	v4l2.StreamOffOutput(enc.fd)
	v4l2.StreamOffCapture(enc.fd)

	// Unsubscribe from events
	v4l2.UnsubscribeCodecEvents(enc.fd)

	// Unmap buffers
	if enc.inputBuffers != nil {
		for _, buf := range enc.inputBuffers {
			sys.Munmap(buf)
		}
		enc.inputBuffers = nil
	}
	if enc.outputBuffers != nil {
		for _, buf := range enc.outputBuffers {
			sys.Munmap(buf)
		}
		enc.outputBuffers = nil
	}

	// Close channels
	if enc.inputChan != nil {
		close(enc.inputChan)
		enc.inputChan = nil
	}
	if enc.errorChan != nil {
		close(enc.errorChan)
		enc.errorChan = nil
	}

	enc.stateMachine.Stop()

	return nil
}

// Close closes the encoder and releases all resources.
func (enc *Encoder) Close() error {
	if enc.streaming.Load() {
		enc.Stop()
	}

	// Wait for encode loop to finish
	if enc.doneChan != nil {
		<-enc.doneChan
	}

	enc.stateMachine.Uninitialize()

	return v4l2.CloseDevice(enc.fd)
}

// GetInput returns a channel for sending raw frames to encode.
// Frames should be in the configured input format (e.g., NV12).
func (enc *Encoder) GetInput() chan<- []byte {
	return enc.inputChan
}

// GetOutput returns a channel for receiving encoded data.
// Data is in the configured output format (e.g., H.264 NAL units).
func (enc *Encoder) GetOutput() <-chan []byte {
	return enc.outputChan
}

// GetError returns a channel for receiving encoding errors.
func (enc *Encoder) GetError() <-chan error {
	return enc.errorChan
}

// GetState returns the current encoder state.
func (enc *Encoder) GetState() v4l2.CodecState {
	return enc.stateMachine.GetState()
}

// GetInputFormat returns the configured input format.
func (enc *Encoder) GetInputFormat() v4l2.PixFormat {
	return enc.inputFormat
}

// GetOutputFormat returns the configured output format.
func (enc *Encoder) GetOutputFormat() v4l2.PixFormat {
	return enc.outputFormat
}

// Fd returns the underlying file descriptor.
func (enc *Encoder) Fd() uintptr {
	return enc.fd
}

// Capability returns the device capabilities.
func (enc *Encoder) Capability() v4l2.Capability {
	return enc.cap
}
