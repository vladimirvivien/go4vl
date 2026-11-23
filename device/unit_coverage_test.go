package device

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
)

func TestDeviceControls_ErrorWrapping(t *testing.T) {
	fd, err := v4l2.OpenDevice("/dev/null", openDeviceFlags, 0)
	if err != nil {
		t.Fatalf("OpenDevice(/dev/null): %v", err)
	}
	defer v4l2.CloseDevice(fd)

	dev := &Device{path: "/dev/null", fd: fd}

	if _, err := dev.GetControl(v4l2.CtrlBrightness); err == nil || !strings.Contains(err.Error(), "/dev/null") {
		t.Errorf("GetControl() = %v, want error containing device path", err)
	}
	if _, err := dev.QueryAllControls(); err == nil || !strings.Contains(err.Error(), "/dev/null") {
		t.Errorf("QueryAllControls() = %v, want error containing device path", err)
	}
	for name, fn := range map[string]func(v4l2.CtrlValue) error{
		"Brightness": dev.SetControlBrightness,
		"Contrast":   dev.SetControlContrast,
		"Saturation": dev.SetControlSaturation,
		"Hue":        dev.SetControlHue,
	} {
		if err := fn(128); err == nil || !strings.Contains(err.Error(), "/dev/null") {
			t.Errorf("SetControl%s() = %v, want error containing device path", name, err)
		}
	}
}

func TestWriteFrame_ValidationPaths(t *testing.T) {
	ctx := context.Background()

	t.Run("ReadWriteModeRejected", func(t *testing.T) {
		d := &Device{
			bufType: v4l2.BufTypeVideoOutput,
			config:  config{ioMethod: IOMethodReadWrite},
		}
		if err := d.WriteFrame(ctx, []byte{1, 2, 3}); err == nil {
			t.Fatal("expected error in read/write IO mode")
		}
	})

	t.Run("NonOutputBufTypeRejected", func(t *testing.T) {
		d := &Device{bufType: v4l2.BufTypeVideoCapture}
		if err := d.WriteFrame(ctx, []byte{1, 2, 3}); err == nil {
			t.Fatal("expected error for capture bufType")
		}
	})

	t.Run("BeforeStartRejected", func(t *testing.T) {
		d := &Device{bufType: v4l2.BufTypeVideoOutput}
		if err := d.WriteFrame(ctx, []byte{1, 2, 3}); err == nil {
			t.Fatal("expected error before Start()")
		}
	})

	t.Run("NoBuffersRejected", func(t *testing.T) {
		d := &Device{bufType: v4l2.BufTypeVideoOutput}
		d.streaming.Store(true)
		d.streamActive.Store(true)
		if err := d.WriteFrame(ctx, []byte{1, 2, 3}); err == nil {
			t.Fatal("expected error with no allocated buffers")
		}
	})

	t.Run("OversizedFrameRejected", func(t *testing.T) {
		d := &Device{
			bufType: v4l2.BufTypeVideoOutput,
			buffers: [][]byte{make([]byte, 4)},
		}
		d.streaming.Store(true)
		d.streamActive.Store(true)
		if err := d.WriteFrame(ctx, make([]byte, 8)); err == nil {
			t.Fatal("expected error when frame exceeds buffer size")
		}
	})

	t.Run("MPlaneDelegation", func(t *testing.T) {
		d := &Device{bufType: v4l2.BufTypeVideoOutputMPlane}
		if err := d.WriteFrame(ctx, []byte{1, 2, 3}); err == nil {
			t.Fatal("expected error for MPLANE WriteFrame before Start()")
		}
	})
}

func TestStartFramesCapture_ContextCancellation(t *testing.T) {
	fd, err := v4l2.OpenDevice("/dev/null", openDeviceFlags, 0)
	if err != nil {
		t.Fatalf("OpenDevice(/dev/null): %v", err)
	}
	defer v4l2.CloseDevice(fd)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Pre-cancelled so loop exits immediately without blocking.

	dev := &Device{
		fd:          fd,
		bufType:     v4l2.BufTypeVideoCapture,
		config:      config{ioType: v4l2.IOTypeMMAP, bufSize: 2},
		startCtx:    ctx,
		captureDone: make(chan struct{}),
		framePool:   NewFramePool(4),
	}
	dev.streaming.Store(true)
	dev.streamingMode.Store(2)

	dev.startFramesCapture()

	select {
	case <-dev.captureDone:
	case <-time.After(2 * time.Second):
		t.Fatal("startFramesCapture goroutine did not exit after context cancellation")
	}

	if dev.streaming.Load() {
		t.Error("expected streaming=false after capture exit")
	}
	if dev.streamingMode.Load() != 0 {
		t.Errorf("expected streamingMode=0 after capture exit, got %d", dev.streamingMode.Load())
	}
}

func TestEncoderAndDecoder_ValidationAndAccessors(t *testing.T) {
	t.Run("OpenEncoder_InvalidConfig", func(t *testing.T) {
		if _, err := OpenEncoder("/dev/nonexistent-video-xyz", EncoderConfig{}); err == nil {
			t.Error("expected error for missing dimensions/format")
		}
		if _, err := OpenEncoder("/dev/null", EncoderConfig{
			InputFormat:  v4l2.PixFormat{Width: 640, Height: 480, PixelFormat: v4l2.PixelFmtNV12},
			OutputFormat: v4l2.PixFormat{PixelFormat: v4l2.PixelFmtH264},
		}); err == nil {
			t.Error("expected error opening /dev/null as encoder")
		}
	})

	t.Run("OpenDecoder_InvalidConfig", func(t *testing.T) {
		if _, err := OpenDecoder("/dev/nonexistent-video-xyz", DecoderConfig{}); err == nil {
			t.Error("expected error for missing input format")
		}
		if _, err := OpenDecoder("/dev/null", DecoderConfig{
			InputFormat:  v4l2.PixFormat{PixelFormat: v4l2.PixelFmtH264},
			OutputFormat: v4l2.PixFormat{PixelFormat: v4l2.PixelFmtNV12},
		}); err == nil {
			t.Error("expected error opening /dev/null as decoder")
		}
	})

	t.Run("Encoder_Accessors", func(t *testing.T) {
		enc := &Encoder{
			inputFormat:  v4l2.PixFormat{Width: 640, Height: 480, PixelFormat: v4l2.PixelFmtNV12},
			outputFormat: v4l2.PixFormat{Width: 640, Height: 480, PixelFormat: v4l2.PixelFmtH264},
			inputChan:    make(chan []byte, 1),
			outputChan:   make(chan []byte, 1),
			errorChan:    make(chan error, 1),
			stateMachine: v4l2.NewCodecStateMachine(v4l2.CodecTypeEncoder),
		}
		if enc.GetInput() == nil || enc.GetOutput() == nil || enc.GetError() == nil {
			t.Error("expected non-nil encoder channels")
		}
		if enc.GetInputFormat().PixelFormat != v4l2.PixelFmtNV12 {
			t.Errorf("unexpected input format: %+v", enc.GetInputFormat())
		}
		if enc.GetOutputFormat().PixelFormat != v4l2.PixelFmtH264 {
			t.Errorf("unexpected output format: %+v", enc.GetOutputFormat())
		}
		if enc.GetState() != v4l2.CodecStateUninitialized {
			t.Errorf("unexpected initial state: %v", enc.GetState())
		}
		if err := enc.Stop(); err != nil {
			t.Errorf("Stop() when not streaming returned error: %v", err)
		}
	})

	t.Run("Decoder_Accessors", func(t *testing.T) {
		dec := &Decoder{
			inputFormat:   v4l2.PixFormat{PixelFormat: v4l2.PixelFmtH264},
			outputFormat:  v4l2.PixFormat{PixelFormat: v4l2.PixelFmtNV12},
			inputChan:     make(chan []byte, 1),
			outputChan:    make(chan []byte, 1),
			errorChan:     make(chan error, 1),
			resChangeChan: make(chan ResolutionChange, 1),
			stateMachine:  v4l2.NewCodecStateMachine(v4l2.CodecTypeDecoder),
		}
		if dec.GetInput() == nil || dec.GetOutput() == nil || dec.GetError() == nil || dec.GetResolutionChanges() == nil {
			t.Error("expected non-nil decoder channels")
		}
		if dec.GetInputFormat().PixelFormat != v4l2.PixelFmtH264 {
			t.Errorf("unexpected input format: %+v", dec.GetInputFormat())
		}
		if dec.GetOutputFormat().PixelFormat != v4l2.PixelFmtNV12 {
			t.Errorf("unexpected output format: %+v", dec.GetOutputFormat())
		}
		if dec.GetState() != v4l2.CodecStateUninitialized {
			t.Errorf("unexpected initial state: %v", dec.GetState())
		}
		if err := dec.Stop(); err != nil {
			t.Errorf("Stop() when not streaming returned error: %v", err)
		}
	})
}
