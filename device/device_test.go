package device

import (
	"context"
	"testing"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestDevice_StructFields tests that all Device struct fields are accessible
func TestDevice_StructFields(t *testing.T) {
	dev := Device{
		path:    "/dev/video0",
		fd:      3,
		bufType: v4l2.BufTypeVideoCapture,
		cap: v4l2.Capability{
			Driver:       "uvcvideo",
			Card:         "HD Webcam",
			Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming,
		},
		cropCap: v4l2.CropCapability{
			StreamType: v4l2.BufTypeVideoCapture,
		},
		buffers: [][]byte{
			make([]byte, 614400),
			make([]byte, 614400),
		},
		requestedBuf: v4l2.RequestBuffers{
			Count:      2,
			StreamType: v4l2.BufTypeVideoCapture,
			Memory:     v4l2.IOTypeMMAP,
		},
	}
	dev.streaming.Store(false)
	dev.output = make(chan []byte, 2)
	dev.streamErr = make(chan error, 1)

	if dev.path != "/dev/video0" {
		t.Errorf("path: expected /dev/video0, got %s", dev.path)
	}
	if dev.fd != 3 {
		t.Errorf("fd: expected 3, got %d", dev.fd)
	}
	if dev.bufType != v4l2.BufTypeVideoCapture {
		t.Errorf("bufType: expected %d, got %d", v4l2.BufTypeVideoCapture, dev.bufType)
	}
	if dev.streaming.Load() != false {
		t.Error("streaming: expected false")
	}
	if len(dev.buffers) != 2 {
		t.Errorf("buffers: expected 2, got %d", len(dev.buffers))
	}
}

// TestDevice_Name tests the Name() method
func TestDevice_Name(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"video0", "/dev/video0", "/dev/video0"},
		{"video1", "/dev/video1", "/dev/video1"},
		{"custom path", "/dev/v4l/by-id/usb-046d_HD_Webcam", "/dev/v4l/by-id/usb-046d_HD_Webcam"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev := Device{path: tt.path}
			if result := dev.Name(); result != tt.expected {
				t.Errorf("Name() = %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestDevice_Fd tests the Fd() method
func TestDevice_Fd(t *testing.T) {
	tests := []struct {
		name string
		fd   uintptr
	}{
		{"fd 3", 3},
		{"fd 4", 4},
		{"fd 10", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev := Device{fd: tt.fd}
			if result := dev.Fd(); result != tt.fd {
				t.Errorf("Fd() = %d, want %d", result, tt.fd)
			}
		})
	}
}

// TestDevice_Capability tests the Capability() method
func TestDevice_Capability(t *testing.T) {
	cap := v4l2.Capability{
		Driver:       "uvcvideo",
		Card:         "HD Webcam C920",
		BusInfo:      "usb-0000:00:14.0-1",
		Version:      (5 << 16) | (15 << 8) | 0,
		Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming,
	}

	dev := Device{cap: cap}
	result := dev.Capability()

	if result.Driver != "uvcvideo" {
		t.Errorf("Driver = %s, want uvcvideo", result.Driver)
	}
	if result.Card != "HD Webcam C920" {
		t.Errorf("Card = %s, want HD Webcam C920", result.Card)
	}
	if result.Capabilities != (v4l2.CapVideoCapture | v4l2.CapStreaming) {
		t.Errorf("Capabilities = %d, want %d", result.Capabilities, v4l2.CapVideoCapture|v4l2.CapStreaming)
	}
}

// TestDevice_BufferType tests the BufferType() method
func TestDevice_BufferType(t *testing.T) {
	tests := []struct {
		name    string
		bufType v4l2.BufType
	}{
		{"VideoCapture", v4l2.BufTypeVideoCapture},
		{"VideoOutput", v4l2.BufTypeVideoOutput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev := Device{bufType: tt.bufType}
			if result := dev.BufferType(); result != tt.bufType {
				t.Errorf("BufferType() = %d, want %d", result, tt.bufType)
			}
		})
	}
}

// TestDevice_BufferCount tests the BufferCount() method
func TestDevice_BufferCount(t *testing.T) {
	tests := []struct {
		name     string
		bufSize  uint32
		expected v4l2.BufType
	}{
		{"2 buffers", 2, 2},
		{"4 buffers", 4, 4},
		{"8 buffers", 8, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev := Device{config: config{bufSize: tt.bufSize}}
			if result := dev.BufferCount(); result != tt.expected {
				t.Errorf("BufferCount() = %d, want %d", result, tt.expected)
			}
		})
	}
}

// TestDevice_MemIOType tests the MemIOType() method
func TestDevice_MemIOType(t *testing.T) {
	tests := []struct {
		name   string
		ioType v4l2.IOType
	}{
		{"MMAP", v4l2.IOTypeMMAP},
		{"UserPtr", v4l2.IOTypeUserPtr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev := Device{config: config{ioType: tt.ioType}}
			if result := dev.MemIOType(); result != tt.ioType {
				t.Errorf("MemIOType() = %d, want %d", result, tt.ioType)
			}
		})
	}
}

// TestDevice_GetOutput tests the GetOutput() method
func TestDevice_GetOutput(t *testing.T) {
	// Create buffered channel
	output := make(chan []byte, 2)
	dev := Device{output: output}

	result := dev.GetOutput()
	if result != output {
		t.Error("GetOutput() should return the output channel")
	}

	// Test that it's read-only by verifying it's a receive-only channel
	_, ok := (interface{}(result)).(<-chan []byte)
	if !ok {
		t.Error("GetOutput() should return a read-only channel")
	}
}

// TestDevice_GetError tests the GetError() method
func TestDevice_GetError(t *testing.T) {
	// Create buffered channel
	streamErr := make(chan error, 1)
	dev := Device{streamErr: streamErr}

	result := dev.GetError()
	if result != streamErr {
		t.Error("GetError() should return the streamErr channel")
	}

	// Test that it's read-only by verifying it's a receive-only channel
	_, ok := (interface{}(result)).(<-chan error)
	if !ok {
		t.Error("GetError() should return a read-only channel")
	}
}

// TestDevice_Buffers tests the Buffers() method
func TestDevice_Buffers(t *testing.T) {
	buffers := [][]byte{
		make([]byte, 614400),
		make([]byte, 614400),
		make([]byte, 614400),
	}

	dev := Device{buffers: buffers}
	result := dev.Buffers()

	if len(result) != 3 {
		t.Errorf("Buffers() length = %d, want 3", len(result))
	}

	// Verify it's the same slice reference
	if &result[0] != &buffers[0] {
		t.Error("Buffers() should return the same slice reference")
	}
}

// TestDevice_Buffers_Nil tests Buffers() when no buffers are allocated
func TestDevice_Buffers_Nil(t *testing.T) {
	dev := Device{}
	result := dev.Buffers()

	if result != nil {
		t.Errorf("Buffers() = %v, want nil", result)
	}
}

// TestDevice_StreamingFlag tests atomic streaming flag operations
func TestDevice_StreamingFlag(t *testing.T) {
	dev := Device{}

	// Initial state should be false
	if dev.streaming.Load() != false {
		t.Error("Initial streaming state should be false")
	}

	// Set to true
	dev.streaming.Store(true)
	if dev.streaming.Load() != true {
		t.Error("After Store(true), streaming should be true")
	}

	// Set to false
	dev.streaming.Store(false)
	if dev.streaming.Load() != false {
		t.Error("After Store(false), streaming should be false")
	}
}

// TestDevice_GetCropCapability_UnsupportedFeature tests error handling
func TestDevice_GetCropCapability_UnsupportedFeature(t *testing.T) {
	// Device without video capture support
	dev := Device{
		cap: v4l2.Capability{
			Capabilities: v4l2.CapVideoOutput, // Only output, no capture
		},
	}

	_, err := dev.GetCropCapability()
	if err != v4l2.ErrorUnsupportedFeature {
		t.Errorf("GetCropCapability() error = %v, want %v", err, v4l2.ErrorUnsupportedFeature)
	}
}

// TestDevice_GetCropCapability_Supported tests successful crop capability retrieval
func TestDevice_GetCropCapability_Supported(t *testing.T) {
	cropCap := v4l2.CropCapability{
		StreamType: v4l2.BufTypeVideoCapture,
		Bounds: v4l2.Rect{
			Width:  1920,
			Height: 1080,
		},
	}

	dev := Device{
		cap: v4l2.Capability{
			Capabilities: v4l2.CapVideoCapture,
		},
		cropCap: cropCap,
	}

	result, err := dev.GetCropCapability()
	if err != nil {
		t.Errorf("GetCropCapability() error = %v, want nil", err)
	}
	if result.StreamType != v4l2.BufTypeVideoCapture {
		t.Errorf("StreamType = %d, want %d", result.StreamType, v4l2.BufTypeVideoCapture)
	}
	if result.Bounds.Width != 1920 {
		t.Errorf("Bounds.Width = %d, want 1920", result.Bounds.Width)
	}
}

// TestConfig_StructFields tests config struct field accessibility
func TestConfig_StructFields(t *testing.T) {
	cfg := config{
		ioType: v4l2.IOTypeMMAP,
		pixFormat: v4l2.PixFormat{
			Width:       640,
			Height:      480,
			PixelFormat: v4l2.PixelFmtMJPEG,
		},
		bufSize: 4,
		fps:     30,
		bufType: v4l2.BufTypeVideoCapture,
	}

	if cfg.ioType != v4l2.IOTypeMMAP {
		t.Errorf("ioType: expected %d, got %d", v4l2.IOTypeMMAP, cfg.ioType)
	}
	if cfg.pixFormat.Width != 640 {
		t.Errorf("pixFormat.Width: expected 640, got %d", cfg.pixFormat.Width)
	}
	if cfg.bufSize != 4 {
		t.Errorf("bufSize: expected 4, got %d", cfg.bufSize)
	}
	if cfg.fps != 30 {
		t.Errorf("fps: expected 30, got %d", cfg.fps)
	}
	if cfg.bufType != v4l2.BufTypeVideoCapture {
		t.Errorf("bufType: expected %d, got %d", v4l2.BufTypeVideoCapture, cfg.bufType)
	}
}

// TestWithIOType tests the WithIOType option function
func TestWithIOType(t *testing.T) {
	tests := []struct {
		name   string
		ioType v4l2.IOType
	}{
		{"MMAP", v4l2.IOTypeMMAP},
		{"UserPtr", v4l2.IOTypeUserPtr},
		{"Overlay", v4l2.IOTypeOverlay},
		{"DMABuf", v4l2.IOTypeDMABuf},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config{}
			opt := WithIOType(tt.ioType)
			opt(&cfg)

			if cfg.ioType != tt.ioType {
				t.Errorf("ioType = %d, want %d", cfg.ioType, tt.ioType)
			}
		})
	}
}

// TestWithPixFormat tests the WithPixFormat option function
func TestWithPixFormat(t *testing.T) {
	tests := []struct {
		name      string
		pixFormat v4l2.PixFormat
	}{
		{
			name: "640x480 MJPEG",
			pixFormat: v4l2.PixFormat{
				Width:       640,
				Height:      480,
				PixelFormat: v4l2.PixelFmtMJPEG,
			},
		},
		{
			name: "1920x1080 YUYV",
			pixFormat: v4l2.PixFormat{
				Width:       1920,
				Height:      1080,
				PixelFormat: v4l2.PixelFmtYUYV,
			},
		},
		{
			name: "1280x720 H264",
			pixFormat: v4l2.PixFormat{
				Width:       1280,
				Height:      720,
				PixelFormat: v4l2.PixelFmtH264,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config{}
			opt := WithPixFormat(tt.pixFormat)
			opt(&cfg)

			if cfg.pixFormat.Width != tt.pixFormat.Width {
				t.Errorf("Width = %d, want %d", cfg.pixFormat.Width, tt.pixFormat.Width)
			}
			if cfg.pixFormat.Height != tt.pixFormat.Height {
				t.Errorf("Height = %d, want %d", cfg.pixFormat.Height, tt.pixFormat.Height)
			}
			if cfg.pixFormat.PixelFormat != tt.pixFormat.PixelFormat {
				t.Errorf("PixelFormat = %d, want %d", cfg.pixFormat.PixelFormat, tt.pixFormat.PixelFormat)
			}
		})
	}
}

// TestWithBufferSize tests the WithBufferSize option function
func TestWithBufferSize(t *testing.T) {
	tests := []struct {
		name    string
		bufSize uint32
	}{
		{"1 buffer", 1},
		{"2 buffers", 2},
		{"4 buffers", 4},
		{"8 buffers", 8},
		{"16 buffers", 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config{}
			opt := WithBufferSize(tt.bufSize)
			opt(&cfg)

			if cfg.bufSize != tt.bufSize {
				t.Errorf("bufSize = %d, want %d", cfg.bufSize, tt.bufSize)
			}
		})
	}
}

// TestWithFPS tests the WithFPS option function
func TestWithFPS(t *testing.T) {
	tests := []struct {
		name string
		fps  uint32
	}{
		{"15 FPS", 15},
		{"24 FPS", 24},
		{"30 FPS", 30},
		{"60 FPS", 60},
		{"120 FPS", 120},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config{}
			opt := WithFPS(tt.fps)
			opt(&cfg)

			if cfg.fps != tt.fps {
				t.Errorf("fps = %d, want %d", cfg.fps, tt.fps)
			}
		})
	}
}

// TestWithVideoCaptureEnabled tests the WithVideoCaptureEnabled option function
func TestWithVideoCaptureEnabled(t *testing.T) {
	cfg := config{}
	opt := WithVideoCaptureEnabled()
	opt(&cfg)

	if cfg.bufType != v4l2.BufTypeVideoCapture {
		t.Errorf("bufType = %d, want %d", cfg.bufType, v4l2.BufTypeVideoCapture)
	}
}

// TestWithVideoOutputEnabled tests the WithVideoOutputEnabled option function
func TestWithVideoOutputEnabled(t *testing.T) {
	cfg := config{}
	opt := WithVideoOutputEnabled()
	opt(&cfg)

	if cfg.bufType != v4l2.BufTypeVideoOutput {
		t.Errorf("bufType = %d, want %d", cfg.bufType, v4l2.BufTypeVideoOutput)
	}
}

// TestOptionChaining tests that multiple options can be chained
func TestOptionChaining(t *testing.T) {
	cfg := config{}

	options := []Option{
		WithIOType(v4l2.IOTypeMMAP),
		WithBufferSize(4),
		WithFPS(30),
		WithVideoCaptureEnabled(),
		WithPixFormat(v4l2.PixFormat{
			Width:       1280,
			Height:      720,
			PixelFormat: v4l2.PixelFmtMJPEG,
		}),
	}

	// Apply all options
	for _, opt := range options {
		opt(&cfg)
	}

	// Verify all options were applied
	if cfg.ioType != v4l2.IOTypeMMAP {
		t.Errorf("ioType = %d, want %d", cfg.ioType, v4l2.IOTypeMMAP)
	}
	if cfg.bufSize != 4 {
		t.Errorf("bufSize = %d, want 4", cfg.bufSize)
	}
	if cfg.fps != 30 {
		t.Errorf("fps = %d, want 30", cfg.fps)
	}
	if cfg.bufType != v4l2.BufTypeVideoCapture {
		t.Errorf("bufType = %d, want %d", cfg.bufType, v4l2.BufTypeVideoCapture)
	}
	if cfg.pixFormat.Width != 1280 || cfg.pixFormat.Height != 720 {
		t.Errorf("pixFormat = %dx%d, want 1280x720", cfg.pixFormat.Width, cfg.pixFormat.Height)
	}
}

// TestOptionOverride tests that later options override earlier ones
func TestOptionOverride(t *testing.T) {
	cfg := config{}

	// Apply options with overrides
	WithFPS(30)(&cfg)
	WithFPS(60)(&cfg) // This should override the previous value

	if cfg.fps != 60 {
		t.Errorf("fps = %d, want 60 (should be overridden)", cfg.fps)
	}

	// Test buffer size override
	WithBufferSize(2)(&cfg)
	WithBufferSize(8)(&cfg) // This should override

	if cfg.bufSize != 8 {
		t.Errorf("bufSize = %d, want 8 (should be overridden)", cfg.bufSize)
	}
}

// TestDevice_TypicalConfiguration tests a typical device setup
func TestDevice_TypicalConfiguration(t *testing.T) {
	// Simulate typical device configuration
	dev := Device{
		path:    "/dev/video0",
		fd:      3,
		bufType: v4l2.BufTypeVideoCapture,
		cap: v4l2.Capability{
			Driver:       "uvcvideo",
			Card:         "HD Webcam C920",
			Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming,
		},
		config: config{
			ioType: v4l2.IOTypeMMAP,
			pixFormat: v4l2.PixFormat{
				Width:       1920,
				Height:      1080,
				PixelFormat: v4l2.PixelFmtMJPEG,
			},
			bufSize: 4,
			fps:     30,
			bufType: v4l2.BufTypeVideoCapture,
		},
	}

	// Verify configuration
	if dev.Name() != "/dev/video0" {
		t.Errorf("Name = %s, want /dev/video0", dev.Name())
	}
	if dev.Fd() != 3 {
		t.Errorf("Fd = %d, want 3", dev.Fd())
	}
	if dev.BufferType() != v4l2.BufTypeVideoCapture {
		t.Errorf("BufferType = %d, want %d", dev.BufferType(), v4l2.BufTypeVideoCapture)
	}
	if dev.BufferCount() != 4 {
		t.Errorf("BufferCount = %d, want 4", dev.BufferCount())
	}
	if dev.MemIOType() != v4l2.IOTypeMMAP {
		t.Errorf("MemIOType = %d, want %d", dev.MemIOType(), v4l2.IOTypeMMAP)
	}

	cap := dev.Capability()
	if cap.Driver != "uvcvideo" {
		t.Errorf("Driver = %s, want uvcvideo", cap.Driver)
	}
	if !cap.IsVideoCaptureSupported() {
		t.Error("Video capture should be supported")
	}
	if !cap.IsStreamingSupported() {
		t.Error("Streaming should be supported")
	}
}

// TestDevice_ZeroValues tests device with zero/nil values
func TestDevice_ZeroValues(t *testing.T) {
	dev := Device{}

	if dev.Name() != "" {
		t.Errorf("Name = %s, want empty string", dev.Name())
	}
	if dev.Fd() != 0 {
		t.Errorf("Fd = %d, want 0", dev.Fd())
	}
	if dev.BufferType() != 0 {
		t.Errorf("BufferType = %d, want 0", dev.BufferType())
	}
	if dev.BufferCount() != 0 {
		t.Errorf("BufferCount = %d, want 0", dev.BufferCount())
	}
	if dev.Buffers() != nil {
		t.Errorf("Buffers = %v, want nil", dev.Buffers())
	}
	if dev.GetOutput() != nil {
		t.Errorf("GetOutput = %v, want nil", dev.GetOutput())
	}
	if dev.GetError() != nil {
		t.Errorf("GetError = %v, want nil", dev.GetError())
	}
	if dev.streaming.Load() != false {
		t.Error("streaming should be false by default")
	}
}

// TestDevice_ContextCancellation tests Start with cancelled context
func TestDevice_ContextCancellation(t *testing.T) {
	dev := Device{
		cap: v4l2.Capability{
			Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := dev.Start(ctx)
	if err == nil {
		t.Error("Start() with cancelled context should return error")
	}
	if err != context.Canceled {
		t.Errorf("Start() error = %v, want %v", err, context.Canceled)
	}
}

// TestDevice_PixelFormatZeroValue tests GetPixFormat behavior with zero value
func TestDevice_PixelFormatZeroValue(t *testing.T) {
	dev := Device{
		cap: v4l2.Capability{
			Capabilities: v4l2.CapVideoOutput, // No capture support
		},
	}

	_, err := dev.GetPixFormat()
	if err != v4l2.ErrorUnsupportedFeature {
		t.Errorf("GetPixFormat() without capture support should return ErrorUnsupportedFeature, got %v", err)
	}
}

// TestDevice_SetPixFormat_UnsupportedFeature tests SetPixFormat error handling
func TestDevice_SetPixFormat_UnsupportedFeature(t *testing.T) {
	dev := Device{
		cap: v4l2.Capability{
			Capabilities: v4l2.CapVideoOutput, // No capture support
		},
	}

	err := dev.SetPixFormat(v4l2.PixFormat{})
	if err != v4l2.ErrorUnsupportedFeature {
		t.Errorf("SetPixFormat() without capture support should return ErrorUnsupportedFeature, got %v", err)
	}
}

// TestDevice_SetCropRect_UnsupportedFeature tests SetCropRect error handling
func TestDevice_SetCropRect_UnsupportedFeature(t *testing.T) {
	dev := Device{
		cap: v4l2.Capability{
			Capabilities: v4l2.CapVideoOutput, // No capture support
		},
	}

	err := dev.SetCropRect(v4l2.Rect{})
	if err != v4l2.ErrorUnsupportedFeature {
		t.Errorf("SetCropRect() without capture support should return ErrorUnsupportedFeature, got %v", err)
	}
}

// TestDevice_StreamingMutualExclusivity tests that GetOutput and GetFrames are mutually exclusive
func TestDevice_StreamingMutualExclusivity(t *testing.T) {
	t.Run("GetOutput_then_GetFrames", func(t *testing.T) {
		dev := Device{}
		dev.output = make(chan []byte, 2)

		// First call to GetOutput() should succeed
		ch1 := dev.GetOutput()
		if ch1 == nil {
			t.Error("First GetOutput() should return valid channel")
		}
		if dev.streamingMode.Load() != 1 {
			t.Errorf("streamingMode = %d, want 1 (GetOutput)", dev.streamingMode.Load())
		}

		// Second call to GetOutput() should still work (same mode)
		ch2 := dev.GetOutput()
		if ch2 == nil {
			t.Error("Second GetOutput() should return valid channel")
		}

		// Call to GetFrames() should fail (different mode)
		framesCh := dev.GetFrames()
		if framesCh != nil {
			t.Error("GetFrames() after GetOutput() should return nil channel")
		}
	})

	t.Run("GetFrames_then_GetOutput", func(t *testing.T) {
		dev := Device{}
		dev.frames = make(chan *Frame, 2)

		// First call to GetFrames() should succeed
		ch1 := dev.GetFrames()
		if ch1 == nil {
			t.Error("First GetFrames() should return valid channel")
		}
		if dev.streamingMode.Load() != 2 {
			t.Errorf("streamingMode = %d, want 2 (GetFrames)", dev.streamingMode.Load())
		}

		// Second call to GetFrames() should still work (same mode)
		ch2 := dev.GetFrames()
		if ch2 == nil {
			t.Error("Second GetFrames() should return valid channel")
		}

		// Call to GetOutput() should fail (different mode)
		outputCh := dev.GetOutput()
		if outputCh != nil {
			t.Error("GetOutput() after GetFrames() should return nil channel")
		}
	})

	t.Run("Stop_resets_mode", func(t *testing.T) {
		dev := Device{}
		dev.output = make(chan []byte, 2)
		dev.streaming.Store(true) // Simulate streaming state

		// Set to GetOutput mode
		ch := dev.GetOutput()
		if ch == nil {
			t.Error("GetOutput() should succeed initially")
		}

		// Manually reset streaming mode (simulates what Stop() does when streaming is active)
		dev.streamingMode.Store(0)

		if dev.streamingMode.Load() != 0 {
			t.Errorf("After mode reset, streamingMode = %d, want 0", dev.streamingMode.Load())
		}

		// Now GetFrames() should work
		dev.frames = make(chan *Frame, 2)
		framesCh := dev.GetFrames()
		if framesCh == nil {
			t.Error("GetFrames() should succeed after mode reset")
		}
		if dev.streamingMode.Load() != 2 {
			t.Errorf("streamingMode = %d, want 2 (GetFrames)", dev.streamingMode.Load())
		}
	})
}
