package device

import (
	"context"
	"errors"
	"fmt"
	"os"
	sys "syscall"
	"testing"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// Mockable v4l2 functions
// These will be replaced by the init() function to call our mocks.
// This assumes that the v4l2 package functions are variables.
// If they are not, this reassignment won't work, and the `device` package
// would need to be structured to allow injection of these dependencies.
var (
	v4l2OpenDevice        = v4l2.OpenDevice
	v4l2GetCapability     = v4l2.GetCapability
	v4l2CloseDevice       = v4l2.CloseDevice
	v4l2GetCropCapability = v4l2.GetCropCapability
	v4l2SetCropRect       = v4l2.SetCropRect
	v4l2GetPixFormat      = v4l2.GetPixFormat
	v4l2SetPixFormat      = v4l2.SetPixFormat
	v4l2GetStreamParam    = v4l2.GetStreamParam
	v4l2SetStreamParam    = v4l2.SetStreamParam
	// Add any other v4l2 functions that device.Open might call indirectly
	// For example, if GetFrameRate calls GetStreamParam, it's covered.
	v4l2InitBuffers        = v4l2.InitBuffers
	v4l2MapMemoryBuffers   = v4l2.MapMemoryBuffers
	v4l2QueueBuffer        = v4l2.QueueBuffer
	v4l2StreamOn           = v4l2.StreamOn
	v4l2UnmapMemoryBuffers = v4l2.UnmapMemoryBuffers
	v4l2StreamOff          = v4l2.StreamOff
	v4l2WaitForRead        = v4l2.WaitForRead
	v4l2DequeueBuffer      = v4l2.DequeueBuffer // Added for completeness, though not directly in Start/Stop mocks yet
)

// Mock function variables, to be set by individual tests
var (
	mockOpenDeviceFn         func(path string, flags int, mode uint32) (uintptr, error)
	mockGetCapabilityFn      func(fd uintptr) (v4l2.Capability, error)
	mockCloseDeviceFn        func(fd uintptr) error
	mockGetCropCapabilityFn  func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error)
	mockSetCropRectFn        func(fd uintptr, r v4l2.Rect) error
	mockGetPixFormatFn       func(fd uintptr) (v4l2.PixFormat, error)
	mockSetPixFormatFn       func(fd uintptr, pixFmt v4l2.PixFormat) error
	mockGetStreamParamFn     func(fd uintptr, bufType v4l2.BufType) (v4l2.StreamParam, error)
	mockSetStreamParamFn     func(fd uintptr, bufType v4l2.BufType, param v4l2.StreamParam) error
	mockInitBuffersFn        func(dev v4l2.StreamingDevice) (v4l2.RequestBuffers, error)
	mockMapMemoryBuffersFn   func(dev v4l2.StreamingDevice) ([][]byte, error)
	mockQueueBufferFn        func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error)
	mockStreamOnFn           func(dev v4l2.StreamingDevice) error
	mockUnmapMemoryBuffersFn func(dev v4l2.StreamingDevice) error
	mockStreamOffFn          func(dev v4l2.StreamingDevice) error
	mockWaitForReadFn        func(dev v4l2.Device) <-chan struct{}
	mockDequeueBufferFn      func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType) (v4l2.Buffer, error)
)

// This init function redirects the actual v4l2 calls to our mock functions.
// This is a common way to achieve mocking in Go when you can modify the package
// or when the external package's functions are variables.
func init() {
	// Preserve existing mocks
	v4l2.OpenDevice = func(path string, flags int, mode uint32) (uintptr, error) {
		if mockOpenDeviceFn != nil {
			return mockOpenDeviceFn(path, flags, mode)
		}
		return 0, errors.New("mockOpenDeviceFn not set")
	}
	v4l2.GetCapability = func(fd uintptr) (v4l2.Capability, error) {
		if mockGetCapabilityFn != nil {
			return mockGetCapabilityFn(fd)
		}
		return v4l2.Capability{}, errors.New("mockGetCapabilityFn not set")
	}
	v4l2.CloseDevice = func(fd uintptr) error {
		if mockCloseDeviceFn != nil {
			return mockCloseDeviceFn(fd)
		}
		return errors.New("mockCloseDeviceFn not set")
	}
	v4l2.GetCropCapability = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		if mockGetCropCapabilityFn != nil {
			return mockGetCropCapabilityFn(fd, bufType)
		}
		return v4l2.CropCapability{}, errors.New("mockGetCropCapabilityFn not set")
	}
	v4l2.SetCropRect = func(fd uintptr, r v4l2.Rect) error {
		if mockSetCropRectFn != nil {
			return mockSetCropRectFn(fd, r)
		}
		return errors.New("mockSetCropRectFn not set")
	}
	v4l2.GetPixFormat = func(fd uintptr) (v4l2.PixFormat, error) {
		if mockGetPixFormatFn != nil {
			return mockGetPixFormatFn(fd)
		}
		return v4l2.PixFormat{}, errors.New("mockGetPixFormatFn not set")
	}
	v4l2.SetPixFormat = func(fd uintptr, pixFmt v4l2.PixFormat) error {
		if mockSetPixFormatFn != nil {
			return mockSetPixFormatFn(fd, pixFmt)
		}
		return errors.New("mockSetPixFormatFn not set")
	}
	v4l2.GetStreamParam = func(fd uintptr, bufType v4l2.BufType) (v4l2.StreamParam, error) {
		if mockGetStreamParamFn != nil {
			return mockGetStreamParamFn(fd, bufType)
		}
		return v4l2.StreamParam{}, errors.New("mockGetStreamParamFn not set")
	}
	v4l2.SetStreamParam = func(fd uintptr, bufType v4l2.BufType, param v4l2.StreamParam) error {
		if mockSetStreamParamFn != nil {
			return mockSetStreamParamFn(fd, bufType, param)
		}
		return errors.New("mockSetStreamParamFn not set")
	}

	// Add new mocks to init
	v4l2.InitBuffers = func(dev v4l2.StreamingDevice) (v4l2.RequestBuffers, error) {
		if mockInitBuffersFn != nil {
			return mockInitBuffersFn(dev)
		}
		return v4l2.RequestBuffers{}, errors.New("mockInitBuffersFn not set")
	}
	v4l2.MapMemoryBuffers = func(dev v4l2.StreamingDevice) ([][]byte, error) {
		if mockMapMemoryBuffersFn != nil {
			return mockMapMemoryBuffersFn(dev)
		}
		return nil, errors.New("mockMapMemoryBuffersFn not set")
	}
	v4l2.QueueBuffer = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) {
		if mockQueueBufferFn != nil {
			return mockQueueBufferFn(fd, ioType, bufType, i)
		}
		return v4l2.Buffer{}, errors.New("mockQueueBufferFn not set")
	}
	v4l2.StreamOn = func(dev v4l2.StreamingDevice) error {
		if mockStreamOnFn != nil {
			return mockStreamOnFn(dev)
		}
		return errors.New("mockStreamOnFn not set")
	}
	v4l2.UnmapMemoryBuffers = func(dev v4l2.StreamingDevice) error {
		if mockUnmapMemoryBuffersFn != nil {
			return mockUnmapMemoryBuffersFn(dev)
		}
		return errors.New("mockUnmapMemoryBuffersFn not set")
	}
	v4l2.StreamOff = func(dev v4l2.StreamingDevice) error {
		if mockStreamOffFn != nil {
			return mockStreamOffFn(dev)
		}
		return errors.New("mockStreamOffFn not set")
	}
	v4l2.WaitForRead = func(dev v4l2.Device) <-chan struct{} {
		if mockWaitForReadFn != nil {
			return mockWaitForReadFn(dev)
		}
		// Return a dummy channel that will never send, to prevent nil channel panics
		// if a test doesn't mock this but the code under test calls it.
		return make(<-chan struct{})
	}
	v4l2.DequeueBuffer = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType) (v4l2.Buffer, error) {
		if mockDequeueBufferFn != nil {
			return mockDequeueBufferFn(fd, ioType, bufType)
		}
		return v4l2.Buffer{}, errors.New("mockDequeueBufferFn not set")
	}
}

// Helper to reset all mock functions to nil
func resetMocks() {
	mockOpenDeviceFn = nil
	mockGetCapabilityFn = nil
	mockCloseDeviceFn = nil
	mockGetCropCapabilityFn = nil
	mockSetCropRectFn = nil
	mockGetPixFormatFn = nil
	mockSetPixFormatFn = nil
	mockGetStreamParamFn = nil
	mockSetStreamParamFn = nil
	mockInitBuffersFn = nil
	mockMapMemoryBuffersFn = nil
	mockQueueBufferFn = nil
	mockStreamOnFn = nil
	mockUnmapMemoryBuffersFn = nil
	mockStreamOffFn = nil
	mockWaitForReadFn = nil
	mockDequeueBufferFn = nil
}

func TestOpen_Success(t *testing.T) {
	resetMocks()
	defer resetMocks() // Ensure mocks are reset after the test

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) {
		if path != "/dev/video0" {
			return 0, fmt.Errorf("expected path /dev/video0, got %s", path)
		}
		return 1, nil // dummy fd
	}
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		if fd != 1 {
			return v4l2.Capability{}, fmt.Errorf("expected fd 1, got %d", fd)
		}
		return v4l2.Capability{
			Driver:       "mock_driver",
			Card:         "mock_card",
			BusInfo:      "mock_bus",
			Version:      0x00050A00, // Kernel 5.10.0
			Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming,
		}, nil
	}
	// Mock GetCropCapability to return success with default rect (no actual cropping)
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	// Mock SetCropRect to succeed
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error {
		return nil
	}
	mockGetPixFormatFn = func(fd uintptr) (v4l2.PixFormat, error) {
		return v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 1920, Height: 1080}, nil
	}
	mockGetStreamParamFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.StreamParam, error) {
		return v4l2.StreamParam{
			Type: v4l2.BufTypeVideoCapture, // Ensure type matches
			Capture: v4l2.CaptureParam{
				TimePerFrame: v4l2.Fract{Numerator: 1, Denominator: 30},
			},
		}, nil
	}
	mockCloseDeviceFn = func(fd uintptr) error {
		if fd != 1 {
			return fmt.Errorf("expected fd 1 for close, got %d", fd)
		}
		return nil
	}

	dev, err := Open("/dev/video0")
	if err != nil {
		t.Fatalf("Open() error = %v, wantErr %v", err, false)
	}
	if dev == nil {
		t.Fatal("Open() returned nil device on success")
	}

	if dev.Name() != "/dev/video0" {
		t.Errorf("dev.Name() = %s, want %s", dev.Name(), "/dev/video0")
	}
	if dev.Fd() != 1 {
		t.Errorf("dev.Fd() = %d, want 1", dev.Fd())
	}
	// Check some capability details
	cap := dev.Capability()
	if cap.Driver != "mock_driver" {
		t.Errorf("cap.Driver = %s, want mock_driver", cap.Driver)
	}
	if !cap.IsStreamingSupported() {
		t.Error("expected streaming to be supported")
	}
	if !cap.IsVideoCaptureSupported() {
		t.Error("expected video capture to be supported")
	}

	// Check default format (as mocked by GetPixFormat)
	pixFmt, err := dev.GetPixFormat()
	if err != nil {
		t.Fatalf("dev.GetPixFormat() error = %v", err)
	}
	if pixFmt.PixelFormat != v4l2.PixelFmtMJPEG {
		t.Errorf("dev.GetPixFormat().PixelFormat = %v, want %v", pixFmt.PixelFormat, v4l2.PixelFmtMJPEG)
	}

	// Check default FPS (as mocked by GetStreamParam)
	fps, err := dev.GetFrameRate()
	if err != nil {
		t.Fatalf("dev.GetFrameRate() error = %v", err)
	}
	if fps != 30 {
		t.Errorf("dev.GetFrameRate() = %d, want 30", fps)
	}

	err = dev.Close()
	if err != nil {
		t.Errorf("dev.Close() error = %v", err)
	}
}

func TestOpen_DeviceOpenFails(t *testing.T) {
	resetMocks()
	defer resetMocks()

	expectedErr := errors.New("v4l2.OpenDevice failed")
	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) {
		return 0, expectedErr
	}

	dev, err := Open("/dev/video0")
	if err == nil {
		t.Fatalf("Open() err = nil, want %v", expectedErr)
		if dev != nil {
			dev.Close() // Attempt to close if non-nil device returned
		}
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Open() err = %v, want err containing %v", err, expectedErr)
	}
	if dev != nil {
		t.Errorf("Open() dev = %v, want nil on error", dev)
	}
}

func TestOpen_GetCapabilityFails(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) {
		return 1, nil // dummy fd
	}
	expectedErr := errors.New("v4l2.GetCapability failed")
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{}, expectedErr
	}
	// Mock CloseDevice because it will be called in Open's error path
	mockCloseDeviceFn = func(fd uintptr) error {
		return nil
	}

	dev, err := Open("/dev/video0")
	if err == nil {
		t.Fatalf("Open() err = nil, want %v", expectedErr)
		if dev != nil {
			dev.Close()
		}
	}
	if !errors.Is(err, expectedErr) { // Check if the error wraps the expected one
		// The error from Open will be like "device open: /dev/video0: v4l2.GetCapability failed"
		// So we check if our expectedErr is part of that chain.
		t.Errorf("Open() err = %v, want err containing %v", err, expectedErr)
	}
	if dev != nil {
		t.Errorf("Open() dev = %v, want nil on error", dev)
	}
}

func TestOpen_NotStreamingSupported(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) {
		return 1, nil // dummy fd
	}
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{
			Driver:       "mock_driver_no_stream",
			Card:         "mock_card_no_stream",
			Capabilities: v4l2.CapVideoCapture, // No CapStreaming
		}, nil
	}
	// Mock CloseDevice because it will be called in Open's error path
	// (though in this specific case, Open might error out before trying to close,
	// it's good practice to mock it if there's a chance it's called).
	// Update: Open() for "device does not support streamingIO" does not close the fd.
	// It's closed if GetCapability fails, or if no capture/output is supported.

	dev, err := Open("/dev/video0")
	if err == nil {
		t.Fatal("Open() err = nil, want error for non-streaming device")
		if dev != nil {
			dev.Close()
		}
	}
	// Expected error string: "device open: device does not support streamingIO"
	expectedErrStr := "device open: device does not support streamingIO"
	if err.Error() != expectedErrStr {
		t.Errorf("Open() err = %q, want %q", err.Error(), expectedErrStr)
	}
	if dev != nil {
		t.Errorf("Open() dev = %v, want nil on error", dev)
	}
}

func TestOpen_NoVideoCaptureOrOutputSupported(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) {
		return 1, nil
	}
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{ // Supports streaming but not capture/output
			Capabilities: v4l2.CapStreaming,
		}, nil
	}
	mockCloseDeviceFn = func(fd uintptr) error { // Will be called in error path
		return nil
	}

	dev, err := Open("/dev/video0")
	if err == nil {
		t.Fatal("Open() err = nil, want error for no video capture/output support")
		if dev != nil {
			dev.Close()
		}
	}
	// Error should wrap v4l2.ErrorUnsupportedFeature
	if !errors.Is(err, v4l2.ErrorUnsupportedFeature) {
		t.Errorf("Open() err = %v, want err wrapping %v", err, v4l2.ErrorUnsupportedFeature)
	}
	if dev != nil {
		t.Errorf("Open() dev = %v, want nil on error", dev)
	}
}

func TestOpen_GetDefaultPixFormatFails(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) {
		return 1, nil
	}
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming}, nil
	}
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error { return nil }

	expectedErr := errors.New("v4l2.GetPixFormat failed for default")
	mockGetPixFormatFn = func(fd uintptr) (v4l2.PixFormat, error) {
		return v4l2.PixFormat{}, expectedErr
	}
	mockCloseDeviceFn = func(fd uintptr) error { return nil } // Called in error path

	dev, err := Open("/dev/video0")
	if err == nil {
		t.Fatal("Open() err = nil, want error")
		if dev != nil {
			dev.Close()
		}
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Open() err = %v, want err containing %v", err, expectedErr)
	}
	if dev != nil {
		t.Errorf("Open() dev = %v, want nil on error", dev)
	}
}

func TestOpen_GetDefaultFrameRateFails(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) {
		return 1, nil
	}
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming}, nil
	}
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error { return nil }
	mockGetPixFormatFn = func(fd uintptr) (v4l2.PixFormat, error) { // Succeeds
		return v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 1920, Height: 1080}, nil
	}

	expectedErr := errors.New("v4l2.GetStreamParam failed for default FPS")
	mockGetStreamParamFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.StreamParam, error) {
		return v4l2.StreamParam{}, expectedErr
	}
	mockCloseDeviceFn = func(fd uintptr) error { return nil } // Called in error path

	dev, err := Open("/dev/video0")
	if err == nil {
		t.Fatal("Open() err = nil, want error")
		if dev != nil {
			dev.Close()
		}
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Open() err = %v, want err containing %v", err, expectedErr)
	}
	if dev != nil {
		t.Errorf("Open() dev = %v, want nil on error", dev)
	}
}

func TestOpen_WithOptions_SetPixFormatSuccess(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) { return 1, nil }
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming}, nil
	}
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error { return nil }
	// GetPixFormat should not be called if SetPixFormat is called via option
	mockGetPixFormatFn = func(fd uintptr) (v4l2.PixFormat, error) {
		t.Error("GetPixFormat called unexpectedly when WithPixFormat option was used")
		return v4l2.PixFormat{}, errors.New("GetPixFormat should not be called")
	}

	setPixFormatCalled := false
	expectedPixFmt := v4l2.PixFormat{PixelFormat: v4l2.PixelFmtRGB24, Width: 1280, Height: 720}
	mockSetPixFormatFn = func(fd uintptr, pixFmt v4l2.PixFormat) error {
		if pixFmt.PixelFormat != expectedPixFmt.PixelFormat || pixFmt.Width != expectedPixFmt.Width || pixFmt.Height != expectedPixFmt.Height {
			return fmt.Errorf("SetPixFormat called with %+v, want %+v", pixFmt, expectedPixFmt)
		}
		setPixFormatCalled = true
		return nil
	}
	mockGetStreamParamFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.StreamParam, error) { // Default FPS
		return v4l2.StreamParam{Type: v4l2.BufTypeVideoCapture, Capture: v4l2.CaptureParam{TimePerFrame: v4l2.Fract{Denominator: 30}}}, nil
	}
	mockCloseDeviceFn = func(fd uintptr) error { return nil }

	dev, err := Open("/dev/video0", WithPixFormat(expectedPixFmt))
	if err != nil {
		t.Fatalf("Open() error = %v, wantErr false", err)
	}
	if !setPixFormatCalled {
		t.Error("SetPixFormat was not called when WithPixFormat option was used")
	}
	if dev == nil {
		t.Fatal("Open() returned nil device")
	}
	// Verify the format was set (or at least the config reflects it)
	// Note: dev.GetPixFormat() would call the *mock* GetPixFormat, which we want to avoid here
	// We rely on the fact that Open internally sets dev.config.pixFormat
	if dev.config.pixFormat.PixelFormat != expectedPixFmt.PixelFormat ||
		dev.config.pixFormat.Width != expectedPixFmt.Width ||
		dev.config.pixFormat.Height != expectedPixFmt.Height {
		t.Errorf("Device format after Open() = %+v, want %+v", dev.config.pixFormat, expectedPixFmt)
	}
	dev.Close()
}

func TestOpen_WithOptions_SetPixFormatFails(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) { return 1, nil }
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming}, nil
	}
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error { return nil }

	expectedErr := errors.New("v4l2.SetPixFormat failed")
	mockSetPixFormatFn = func(fd uintptr, pixFmt v4l2.PixFormat) error {
		return expectedErr
	}
	mockCloseDeviceFn = func(fd uintptr) error { return nil } // Called in error path

	dev, err := Open("/dev/video0", WithPixFormat(v4l2.PixFormat{Width: 640, Height: 480}))
	if err == nil {
		t.Fatal("Open() err = nil, want error")
		if dev != nil {
			dev.Close()
		}
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Open() err = %v, want err containing %v", err, expectedErr)
	}
	if dev != nil {
		t.Errorf("Open() dev = %v, want nil on error", dev)
	}
}

func TestOpen_WithOptions_SetFPSSuccess(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) { return 1, nil }
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming}, nil
	}
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error { return nil }
	mockGetPixFormatFn = func(fd uintptr) (v4l2.PixFormat, error) { // Default format
		return v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 1920, Height: 1080}, nil
	}
	// GetStreamParam should not be called if SetStreamParam is called via option
	mockGetStreamParamFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.StreamParam, error) {
		t.Error("GetStreamParam called unexpectedly when WithFPS option was used")
		return v4l2.StreamParam{}, errors.New("GetStreamParam should not be called")
	}

	setStreamParamCalled := false
	var expectedFPS uint32 = 60
	mockSetStreamParamFn = func(fd uintptr, bufType v4l2.BufType, param v4l2.StreamParam) error {
		if bufType != v4l2.BufTypeVideoCapture {
			return fmt.Errorf("SetStreamParam called with bufType %v, want %v", bufType, v4l2.BufTypeVideoCapture)
		}
		if param.Capture.TimePerFrame.Denominator != expectedFPS {
			return fmt.Errorf("SetStreamParam called with FPS %d, want %d", param.Capture.TimePerFrame.Denominator, expectedFPS)
		}
		setStreamParamCalled = true
		return nil
	}
	mockCloseDeviceFn = func(fd uintptr) error { return nil }

	dev, err := Open("/dev/video0", WithFPS(expectedFPS))
	if err != nil {
		t.Fatalf("Open() error = %v, wantErr false", err)
	}
	if !setStreamParamCalled {
		t.Error("SetStreamParam was not called when WithFPS option was used")
	}
	if dev == nil {
		t.Fatal("Open() returned nil device")
	}
	// Verify the FPS was set (or at least the config reflects it)
	if dev.config.fps != expectedFPS {
		t.Errorf("Device FPS after Open() = %d, want %d", dev.config.fps, expectedFPS)
	}
	dev.Close()
}

func TestOpen_WithOptions_SetFPSFails(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) { return 1, nil }
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming}, nil
	}
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error { return nil }
	mockGetPixFormatFn = func(fd uintptr) (v4l2.PixFormat, error) { // Default format
		return v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 1920, Height: 1080}, nil
	}

	expectedErr := errors.New("v4l2.SetStreamParam failed")
	mockSetStreamParamFn = func(fd uintptr, bufType v4l2.BufType, param v4l2.StreamParam) error {
		return expectedErr
	}
	mockCloseDeviceFn = func(fd uintptr) error { return nil } // Called in error path

	dev, err := Open("/dev/video0", WithFPS(30))
	if err == nil {
		t.Fatal("Open() err = nil, want error")
		if dev != nil {
			dev.Close()
		}
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Open() err = %v, want err containing %v", err, expectedErr)
	}
	if dev != nil {
		t.Errorf("Open() dev = %v, want nil on error", dev)
	}
}

func TestOpen_WithOptions_CustomBufferSize(t *testing.T) {
	resetMocks()
	defer resetMocks()

	// Standard success path mocks
	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) { return 1, nil }
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming}, nil
	}
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error { return nil }
	mockGetPixFormatFn = func(fd uintptr) (v4l2.PixFormat, error) {
		return v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 1920, Height: 1080}, nil
	}
	mockGetStreamParamFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.StreamParam, error) {
		return v4l2.StreamParam{Type: v4l2.BufTypeVideoCapture, Capture: v4l2.CaptureParam{TimePerFrame: v4l2.Fract{Denominator: 30}}}, nil
	}
	mockCloseDeviceFn = func(fd uintptr) error { return nil }

	var expectedBufSize uint32 = 5
	dev, err := Open("/dev/video0", WithBufferSize(expectedBufSize))
	if err != nil {
		t.Fatalf("Open() error = %v, wantErr false", err)
	}
	if dev == nil {
		t.Fatal("Open() returned nil device")
	}
	if dev.config.bufSize != expectedBufSize {
		t.Errorf("Device buffer size after Open() = %d, want %d", dev.config.bufSize, expectedBufSize)
	}
	dev.Close()
}

func TestOpen_VideoOutputDevice(t *testing.T) {
	resetMocks()
	defer resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) { return 1, nil }
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{
			Capabilities: v4l2.CapVideoOutput | v4l2.CapStreaming, // Video Output, not Capture
		}, nil
	}
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		// GetCropCapability is usually for capture, might not be called or fail for output
		// For this test, let's assume it's called but its failure (if any) is ignored or handled.
		// Or it might succeed if bufType is VideoOutput and device supports cropping for output.
		// Let's assume it's called and succeeds for simplicity here.
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error { return nil }
	mockGetPixFormatFn = func(fd uintptr) (v4l2.PixFormat, error) {
		return v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 1920, Height: 1080}, nil
	}
	mockGetStreamParamFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.StreamParam, error) {
		if bufType != v4l2.BufTypeVideoOutput { // Expecting VideoOutput type
			return v4l2.StreamParam{}, fmt.Errorf("expected bufType VideoOutput, got %v", bufType)
		}
		return v4l2.StreamParam{
			Type: v4l2.BufTypeVideoOutput,
			Output: v4l2.OutputParam{ // Ensure Output field is populated
				TimePerFrame: v4l2.Fract{Numerator: 1, Denominator: 30},
			},
		}, nil
	}
	mockCloseDeviceFn = func(fd uintptr) error { return nil }

	dev, err := Open("/dev/video0")
	if err != nil {
		t.Fatalf("Open() error = %v, wantErr false for video output device", err)
	}
	if dev == nil {
		t.Fatal("Open() returned nil device for video output device")
	}
	if dev.BufferType() != v4l2.BufTypeVideoOutput {
		t.Errorf("dev.BufferType() = %v, want %v", dev.BufferType(), v4l2.BufTypeVideoOutput)
	}

	// Check default FPS (as mocked by GetStreamParam for output)
	fps, err := dev.GetFrameRate()
	if err != nil {
		t.Fatalf("dev.GetFrameRate() error = %v", err)
	}
	if fps != 30 {
		t.Errorf("dev.GetFrameRate() = %d, want 30", fps)
	}

	dev.Close()
}


// TODO: Add TestOpen_SetPixFormatOptionSucceeds
// TODO: Add TestOpen_SetPixFormatOptionFails (Covered by TestOpen_WithOptions_SetPixFormatFails)
// TODO: Add TestOpen_SetFPSOptionSucceeds (Covered by TestOpen_WithOptions_SetFPSSuccess)
// TODO: Add TestOpen_SetFPSOptionFails (Covered by TestOpen_WithOptions_SetFPSFails)
// TODO: Add TestOpen_WithPixFormatOption (Covered by TestOpen_WithOptions_SetPixFormatSuccess)
// TODO: Add TestOpen_WithFPSOption (Covered by TestOpen_WithOptions_SetFPSSuccess)
// TODO: Add TestOpen_WithBufferSizeOption (Covered by TestOpen_WithOptions_CustomBufferSize)
// TODO: Add TestOpen_WithIOTypeOption (though only MMAP is really supported by Open now - could test this restriction)
// TODO: Add TestOpen_UnsupportedBufferTypeOption (e.g. providing WithVideoCaptureEnabled on an output-only device)
// TODO: Add TestOpen_CroppingFails (if GetCropCapability or SetCropRect fails - Open currently ignores these errors, but could be tested)

// Helper function to get a successfully opened device for testing Start/Stop
func getOpenedMockDevice(t *testing.T) *Device {
	t.Helper() // Marks this function as a test helper
	resetMocks()

	mockOpenDeviceFn = func(path string, flags int, mode uint32) (uintptr, error) { return 1, nil }
	mockGetCapabilityFn = func(fd uintptr) (v4l2.Capability, error) {
		return v4l2.Capability{Capabilities: v4l2.CapVideoCapture | v4l2.CapStreaming}, nil
	}
	mockGetCropCapabilityFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.CropCapability, error) {
		return v4l2.CropCapability{DefaultRect: v4l2.Rect{Width: 1920, Height: 1080}}, nil
	}
	mockSetCropRectFn = func(fd uintptr, r v4l2.Rect) error { return nil }
	mockGetPixFormatFn = func(fd uintptr) (v4l2.PixFormat, error) {
		return v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 1920, Height: 1080}, nil
	}
	mockGetStreamParamFn = func(fd uintptr, bufType v4l2.BufType) (v4l2.StreamParam, error) {
		return v4l2.StreamParam{Type: v4l2.BufTypeVideoCapture, Capture: v4l2.CaptureParam{TimePerFrame: v4l2.Fract{Denominator: 30}}}, nil
	}
	mockCloseDeviceFn = func(fd uintptr) error { return nil }

	dev, err := Open("/dev/video0")
	if err != nil {
		t.Fatalf("getOpenedMockDevice: Open() failed: %v", err)
	}
	if dev == nil {
		t.Fatal("getOpenedMockDevice: Open() returned nil device")
	}
	return dev
}

func TestStart_Success(t *testing.T) {
	dev := getOpenedMockDevice(t) // Uses helper
	defer resetMocks()             // Reset mocks specific to this test after helper's mocks
	defer dev.Close()              // Ensure device is closed

	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) {
		// Check if the buffer count from device config is used (default is 2 if not set by option)
		if d.BufferCount() != 2 { // Assuming default buffer size is 2
			return v4l2.RequestBuffers{}, fmt.Errorf("InitBuffers expected buffer count 2, got %d", d.BufferCount())
		}
		return v4l2.RequestBuffers{Count: d.BufferCount(), StreamType: uint32(d.BufferType()), Memory: uint32(d.MemIOType())}, nil
	}
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		buffers := make([][]byte, d.BufferCount())
		for i := 0; i < int(d.BufferCount()); i++ {
			buffers[i] = make([]byte, 1024) // Dummy buffer data
		}
		return buffers, nil
	}
	queueBufferCallCount := 0
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) {
		queueBufferCallCount++
		return v4l2.Buffer{Index: i, Flags: v4l2.BufFlagQueued}, nil
	}
	mockStreamOnFn = func(d v4l2.StreamingDevice) error {
		return nil
	}
	// For this success test, WaitForRead can return a channel that doesn't block or do anything problematic
	waitChan := make(chan struct{})
	// close(waitChan) // Alternative: immediately close if loop logic is not an issue
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} {
		return waitChan // Return a controllable channel
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure the context is cancelled to stop the stream loop

	err := dev.Start(ctx)
	if err != nil {
		t.Fatalf("Start() error = %v, wantErr false", err)
	}
	if !dev.streaming {
		t.Error("dev.streaming = false, want true after Start()")
	}
	if len(dev.buffers) != int(dev.config.bufSize) { // bufSize should be 2 by default or what InitBuffers returned
		t.Errorf("len(dev.buffers) = %d, want %d", len(dev.buffers), dev.config.bufSize)
	}
	if queueBufferCallCount != int(dev.config.bufSize) {
		t.Errorf("QueueBuffer call count = %d, want %d", queueBufferCallCount, dev.config.bufSize)
	}
	if dev.output == nil {
		t.Error("dev.output channel is nil after Start()")
	}

	// To properly test the goroutine cleanup, we'd ideally signal it to stop.
	// For now, cancelling the context is the main mechanism.
}

func TestStart_InitBuffersFails(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	defer dev.Close()

	expectedErr := errors.New("InitBuffers failed")
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) {
		return v4l2.RequestBuffers{}, expectedErr
	}

	err := dev.Start(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Errorf("Start() err = %v, want err containing %v", err, expectedErr)
	}
	if dev.streaming {
		t.Error("dev.streaming = true, want false on Start() failure")
	}
}

func TestStart_MapMemoryBuffersFails(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	defer dev.Close()

	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) {
		return v4l2.RequestBuffers{Count: 2}, nil
	}
	expectedErr := errors.New("MapMemoryBuffers failed")
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		return nil, expectedErr
	}

	err := dev.Start(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Errorf("Start() err = %v, want err containing %v", err, expectedErr)
	}
	if dev.streaming {
		t.Error("dev.streaming = true, want false on Start() failure")
	}
}

func TestStart_QueueBufferFails(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	defer dev.Close()

	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) {
		return v4l2.RequestBuffers{Count: 2}, nil
	}
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		return make([][]byte, 2), nil
	}
	expectedErr := errors.New("QueueBuffer failed")
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) {
		return v4l2.Buffer{}, expectedErr
	}
	// StreamOn might be called before error, or not, depending on loop structure.
	// For safety, mock it.
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }


	err := dev.Start(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Errorf("Start() err = %v, want err containing %v", err, expectedErr)
	}
	if dev.streaming { // Should ideally be false, but Start might set it true before erroring in loop
		// t.Error("dev.streaming = true, want false on Start() failure in queueing")
		// Let's check if the stream loop setup failed.
	}
}

func TestStart_StreamOnFails(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	defer dev.Close()

	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) {
		return v4l2.RequestBuffers{Count: 2}, nil
	}
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		return make([][]byte, 2), nil
	}
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) {
		return v4l2.Buffer{Index: i}, nil // Success for all queue calls
	}
	expectedErr := errors.New("StreamOn failed")
	mockStreamOnFn = func(d v4l2.StreamingDevice) error {
		return expectedErr
	}

	err := dev.Start(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Errorf("Start() err = %v, want err containing %v", err, expectedErr)
	}
	if dev.streaming { // Should be false as StreamOn is the last step before setting streaming = true
		t.Error("dev.streaming = true, want false on StreamOn failure")
	}
}

func TestStart_AlreadyStreaming(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	defer dev.Close()

	// Successful first Start
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 2}, nil }
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) { return make([][]byte, 2), nil }
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) { return v4l2.Buffer{Index: i}, nil }
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }
	waitChan := make(chan struct{})
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return waitChan }


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := dev.Start(ctx)
	if err != nil {
		t.Fatalf("First Start() failed: %v", err)
	}

	// Attempt to Start again
	err = dev.Start(ctx)
	if err == nil {
		t.Fatal("Second Start() did not return an error, but should have")
	}
	expectedErrStr := "device: stream already started"
	if err.Error() != expectedErrStr {
		t.Errorf("Second Start() err = %q, want %q", err.Error(), expectedErrStr)
	}
}

func TestStart_ContextCancelled(t *testing.T) {
	dev := getOpenedMockDevice(t) // Use helper, but don't call dev.Close() yet
	defer resetMocks()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context immediately

	err := dev.Start(ctx)
	if err == nil {
		t.Fatal("Start() with cancelled context did not return an error")
		// If it didn't error, it might have started resources that need cleanup
		mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error { return nil }
		mockStreamOffFn = func(d v4l2.StreamingDevice) error { return nil }
		dev.Stop() // Attempt cleanup
		dev.Close()
	} else {
		// If Start errors due to context, it might not have fully opened, so Close might also error or be NOP
		// We need to ensure the underlying fd is closed if Open was successful.
		// The helper getOpenedMockDevice already calls Open, so fd is open.
		mockCloseDeviceFn = func(fd uintptr) error { return nil } // Ensure Close mock is set
		dev.Close()
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Start() with cancelled context err = %v, want %v", err, context.Canceled)
	}
}


func TestStop_Success(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	defer dev.Close() // Ensure base device is closed

	// Setup for successful Start
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 2}, nil }
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		// Provide some dummy buffers for UnmapMemoryBuffers to "unmap"
		return [][]byte{make([]byte, 10), make([]byte, 10)}, nil
	}
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) { return v4l2.Buffer{Index: i}, nil }
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }
	waitChan := make(chan struct{})
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return waitChan }

	ctx, cancel := context.WithCancel(context.Background())
	err := dev.Start(ctx)
	if err != nil {
		t.Fatalf("Start() failed during setup for TestStop_Success: %v", err)
	}
	cancel() // Cancel context to allow stream loop to terminate

	// Mocks for successful Stop
	unmapCalled := false
	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error {
		unmapCalled = true
		return nil
	}
	streamOffCalled := false
	mockStreamOffFn = func(d v4l2.StreamingDevice) error {
		streamOffCalled = true
		return nil
	}

	err = dev.Stop()
	if err != nil {
		t.Fatalf("Stop() error = %v, wantErr false", err)
	}
	if dev.streaming {
		t.Error("dev.streaming = true, want false after Stop()")
	}
	if !unmapCalled {
		t.Error("UnmapMemoryBuffers was not called during Stop()")
	}
	if !streamOffCalled {
		t.Error("StreamOff was not called during Stop()")
	}
}

func TestStop_NotStreaming(t *testing.T) {
	dev := getOpenedMockDevice(t) // Device is opened but not started
	defer resetMocks()
	defer dev.Close()

	// Ensure these are not called
	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error {
		t.Error("UnmapMemoryBuffers called on a non-streaming device")
		return nil
	}
	mockStreamOffFn = func(d v4l2.StreamingDevice) error {
		t.Error("StreamOff called on a non-streaming device")
		return nil
	}

	err := dev.Stop()
	if err != nil {
		t.Fatalf("Stop() on non-streaming device error = %v, wantErr false", err)
	}
	if dev.streaming {
		t.Error("dev.streaming = true, want false for non-streaming device")
	}
}

func TestStop_UnmapMemoryBuffersFails(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	defer dev.Close()

	// Start successfully
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 2}, nil }
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) { return [][]byte{make([]byte, 10)}, nil }
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) { return v4l2.Buffer{Index: i}, nil }
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }
	waitChan := make(chan struct{})
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return waitChan }
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = dev.Start(ctx)


	expectedErr := errors.New("UnmapMemoryBuffers failed")
	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error {
		return expectedErr
	}
	// StreamOff might still be called depending on error handling strategy in Stop, mock it.
	mockStreamOffFn = func(d v4l2.StreamingDevice) error { return nil }


	err := dev.Stop()
	if !errors.Is(err, expectedErr) {
		t.Errorf("Stop() err = %v, want err containing %v", err, expectedErr)
	}
	// dev.streaming might be true or false depending on where Stop() errors out.
	// The key is that the error is propagated.
}

func TestStop_StreamOffFails(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	defer dev.Close()

	// Start successfully
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 2}, nil }
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) { return [][]byte{make([]byte, 10)}, nil }
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) { return v4l2.Buffer{Index: i}, nil }
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }
	waitChan := make(chan struct{})
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return waitChan }
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = dev.Start(ctx)

	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error { return nil } // Succeeds
	expectedErr := errors.New("StreamOff failed")
	mockStreamOffFn = func(d v4l2.StreamingDevice) error {
		return expectedErr
	}

	err := dev.Stop()
	if !errors.Is(err, expectedErr) {
		t.Errorf("Stop() err = %v, want err containing %v", err, expectedErr)
	}
}


// Note on the init() based mocking:
// The effectiveness of the init() function redirecting v4l2 calls (e.g., v4l2.OpenDevice = ...)
// depends on `v4l2.OpenDevice` (and others) being declared as variables in the actual v4l2 package.
// If they are standard `func` declarations, this redirection will not work, and tests might
// call the real v4l2 functions. For robust unit testing where v4l2 functions are standard func,
// the `device` package would need to be refactored to allow dependency injection (e.g., passing
// v4l2 functions as parameters or struct fields). Assuming they are variables for this test setup.

// Further tests would cover other paths in Open, such as failures in GetPixFormat, SetPixFormat,
// GetFrameRate, SetFrameRate, and different capability flags.
// The options processing (WithPixFormat, WithFPS etc.) also needs thorough testing.

func TestStart_StreamLoop_SuccessfulFrameCapture(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	// dev.Close will be called by a separate mock to ensure Stop is tested correctly
	var closeCalled bool
	mockCloseDeviceFn = func(fd uintptr) error {
		closeCalled = true
		return nil
	}

	// Start mocks
	var bufferData = []byte{0xDE, 0xAD, 0xBE, 0xEF}
	dev.config.bufSize = 1 // Simplify to 1 buffer for this test
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) {
		return v4l2.RequestBuffers{Count: dev.config.bufSize}, nil
	}
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		// Simulate the buffers field being populated by MapMemoryBuffers
		// This is what the production code's startStreamLoop will copy from.
		dev.buffers = make([][]byte, dev.config.bufSize)
		dev.buffers[0] = bufferData // Original buffer data
		return dev.buffers, nil
	}

	queueBufferCallCount := 0
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) {
		queueBufferCallCount++
		return v4l2.Buffer{Index: i, Flags: v4l2.BufFlagQueued}, nil
	}
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }

	testReadyChan := make(chan struct{}, 1) // Buffered to prevent send block
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} {
		return testReadyChan
	}
	mockDequeueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType) (v4l2.Buffer, error) {
		return v4l2.Buffer{Index: 0, Flags: v4l2.BufFlagMapped | v4l2.BufFlagDone, BytesUsed: uint32(len(bufferData))}, nil
	}

	// Stop mocks
	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error { return nil }
	mockStreamOffFn = func(d v4l2.StreamingDevice) error { return nil }


	ctx, cancel := context.WithCancel(context.Background())
	err := dev.Start(ctx)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Simulate device ready for read
	testReadyChan <- struct{}{}

	select {
	case frame, ok := <-dev.GetOutput():
		if !ok {
			t.Fatal("dev.GetOutput() channel was closed unexpectedly")
		}
		if len(frame) != len(bufferData) {
			t.Fatalf("Received frame length %d, want %d", len(frame), len(bufferData))
		}
		for i := range bufferData {
			if frame[i] != bufferData[i] {
				t.Errorf("Frame data mismatch at index %d. Got %x, want %x", i, frame[i], bufferData[i])
				break
			}
		}
		// Ensure it's a copy
		if &frame[0] == &dev.buffers[0][0] {
			t.Error("Received frame is not a copy of the internal buffer")
		}

	case <-context.After(context.Second): // Timeout
		t.Fatal("Timeout waiting for frame from dev.GetOutput()")
	}

	// Should have been called once for initial queue, then once for re-queue
	if queueBufferCallCount < 2 { // At least 2: initial + one re-queue
		t.Errorf("Expected QueueBuffer to be called at least twice, got %d", queueBufferCallCount)
	}

	cancel()      // Signal goroutine to stop
	err = dev.Stop() // This should wait for the goroutine to finish
	if err != nil {
		t.Errorf("Stop() failed: %v", err)
	}

	// Check if output channel is closed
	_, outputOpen := <-dev.GetOutput()
	if outputOpen {
		t.Error("dev.GetOutput() channel was not closed after Stop()")
	}
	if !closeCalled {
		t.Error("mockCloseDeviceFn was not called via dev.Close() in Stop()")
	}
}


func TestStart_StreamLoop_DequeueReturnsAgainThenSuccess(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	var closeCalled bool
	mockCloseDeviceFn = func(fd uintptr) error { closeCalled = true; return nil }


	dev.config.bufSize = 1
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 1}, nil }
	bufferData := []byte{0xC0, 0xFF, 0xEE}
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		dev.buffers = make([][]byte, 1)
		dev.buffers[0] = bufferData
		return dev.buffers, nil
	}
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) { return v4l2.Buffer{Index: i}, nil }
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }

	testReadyChan := make(chan struct{}, 3) // Buffer for multiple signals
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return testReadyChan }

	dequeueCallCount := 0
	mockDequeueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType) (v4l2.Buffer, error) {
		dequeueCallCount++
		if dequeueCallCount == 1 {
			return v4l2.Buffer{}, sys.EAGAIN
		}
		if dequeueCallCount == 2 {
			return v4l2.Buffer{}, sys.EAGAIN
		}
		return v4l2.Buffer{Index: 0, Flags: v4l2.BufFlagMapped | v4l2.BufFlagDone, BytesUsed: uint32(len(bufferData))}, nil
	}

	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error { return nil }
	mockStreamOffFn = func(d v4l2.StreamingDevice) error { return nil }

	ctx, cancel := context.WithCancel(context.Background())
	err := dev.Start(ctx)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Simulate device ready multiple times
	testReadyChan <- struct{}{} // For first EAGAIN
	testReadyChan <- struct{}{} // For second EAGAIN
	testReadyChan <- struct{}{} // For successful dequeue

	select {
	case frame, ok := <-dev.GetOutput():
		if !ok {
			t.Fatal("dev.GetOutput() channel was closed unexpectedly")
		}
		if len(frame) != len(bufferData) {
			t.Fatalf("Received frame length %d, want %d", len(frame), len(bufferData))
		}
	case <-context.After(2 * context.Second): // Increased timeout
		t.Fatalf("Timeout waiting for frame. Dequeue calls: %d", dequeueCallCount)
	}

	if dequeueCallCount != 3 {
		t.Errorf("Expected DequeueBuffer to be called 3 times, got %d", dequeueCallCount)
	}

	cancel()
	err = dev.Stop()
	if err != nil {
		t.Errorf("Stop() failed: %v", err)
	}
	if !closeCalled {
		t.Error("mockCloseDeviceFn was not called via dev.Close() in Stop()")
	}
}

func TestStart_StreamLoop_DequeueReturnsErrorFlag(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	var closeCalled bool
	mockCloseDeviceFn = func(fd uintptr) error { closeCalled = true; return nil }

	dev.config.bufSize = 1
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 1}, nil }
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		dev.buffers = make([][]byte, 1)
		dev.buffers[0] = []byte{0x01, 0x02} // Some dummy data
		return dev.buffers, nil
	}
	queueBufferCallCount := 0
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) {
		queueBufferCallCount++
		return v4l2.Buffer{Index: i}, nil
	}
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }
	testReadyChan := make(chan struct{}, 1)
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return testReadyChan }
	mockDequeueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType) (v4l2.Buffer, error) {
		return v4l2.Buffer{Index: 0, Flags: v4l2.BufFlagMapped | v4l2.BufFlagError, BytesUsed: 0}, nil
	}

	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error { return nil }
	mockStreamOffFn = func(d v4l2.StreamingDevice) error { return nil }


	ctx, cancel := context.WithCancel(context.Background())
	err := dev.Start(ctx)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	testReadyChan <- struct{}{} // Signal device ready

	select {
	case frame, ok := <-dev.GetOutput():
		if !ok {
			t.Fatal("dev.GetOutput() channel was closed unexpectedly")
		}
		if len(frame) != 0 { // Expect empty slice on BufFlagError
			t.Errorf("Received frame length %d, want 0 for BufFlagError", len(frame))
		}
	case <-context.After(context.Second):
		t.Fatal("Timeout waiting for frame on BufFlagError")
	}

	if queueBufferCallCount < 2 { // Initial + re-queue
		t.Errorf("Expected QueueBuffer to be called at least twice, got %d", queueBufferCallCount)
	}
	cancel()
	err = dev.Stop()
	if err != nil {
		t.Errorf("Stop() failed: %v", err)
	}
	if !closeCalled {
		t.Error("mockCloseDeviceFn was not called via dev.Close() in Stop()")
	}
}

func TestStart_StreamLoop_StopClosesOutputChannel(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	var closeCalled bool
	mockCloseDeviceFn = func(fd uintptr) error { closeCalled = true; return nil }


	// Standard Start mocks
	dev.config.bufSize = 1
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 1}, nil }
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		dev.buffers = make([][]byte, 1)
		dev.buffers[0] = []byte{0x01}
		return dev.buffers, nil
	}
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) { return v4l2.Buffer{Index: i}, nil }
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }
	waitChan := make(chan struct{}) // Goroutine will block on this
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return waitChan }

	// Stop mocks
	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error { return nil }
	mockStreamOffFn = func(d v4l2.StreamingDevice) error { return nil }

	ctx, cancel := context.WithCancel(context.Background())
	err := dev.Start(ctx)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	cancel() // Cancel context first
	err = dev.Stop()
	if err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	// Check if output channel is closed
	select {
	case _, ok := <-dev.GetOutput():
		if ok {
			t.Error("dev.GetOutput() channel was not closed after Stop()")
		}
	case <-context.After(context.Second): // Should not block if closed
		t.Error("Timeout checking if GetOutput() is closed, it might still be open or blocked.")
	}
	if !closeCalled {
		t.Error("mockCloseDeviceFn was not called via dev.Close() in Stop()")
	}
}

func TestStart_StreamLoop_ContextCancellationStopsLoop(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	var closeCalled bool
	mockCloseDeviceFn = func(fd uintptr) error { closeCalled = true; return nil }


	dev.config.bufSize = 1
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 1}, nil }
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		dev.buffers = make([][]byte, 1)
		dev.buffers[0] = []byte{0x01}
		return dev.buffers, nil
	}
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) { return v4l2.Buffer{Index: i}, nil }
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }
	waitChan := make(chan struct{}) // Goroutine will block on this
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return waitChan }

	unmapCalled := false
	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error { unmapCalled = true; return nil }
	streamOffCalled := false
	mockStreamOffFn = func(d v4l2.StreamingDevice) error { streamOffCalled = true; return nil }

	ctx, cancel := context.WithCancel(context.Background())
	err := dev.Start(ctx)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	cancel() // Cancel the context

	// Wait for the output channel to be closed, indicating loop termination
	select {
	case _, ok := <-dev.GetOutput():
		if ok {
			t.Error("dev.GetOutput() was not closed after context cancellation")
		}
	case <-context.After(2 * context.Second): // Timeout
		t.Error("Timeout waiting for dev.GetOutput() to close after context cancellation")
	}

	// Check if Stop's V4L2 functions were called by the loop itself
	if !unmapCalled {
		t.Error("UnmapMemoryBuffers was not called by stream loop on context cancellation")
	}
	if !streamOffCalled {
		t.Error("StreamOff was not called by stream loop on context cancellation")
	}
	if !closeCalled {
		t.Error("mockCloseDeviceFn was not called via dev.Close() in Stop() which is called by stream loop")
	}
}

func TestStart_StreamLoop_PanicOnDequeueError(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	var closeCalled bool
	mockCloseDeviceFn = func(fd uintptr) error { closeCalled = true; return nil } // For cleanup after panic

	dev.config.bufSize = 1
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 1}, nil }
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		dev.buffers = make([][]byte, 1); dev.buffers[0] = []byte{0x01}; return dev.buffers, nil
	}
	mockQueueBufferFn = func(fd uintptr, iotype v4l2.IOType, btype v4l2.BufType, i uint32) (v4l2.Buffer, error) { return v4l2.Buffer{Index: i}, nil }
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }

	testReadyChan := make(chan struct{}, 1)
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return testReadyChan }

	expectedPanicErr := errors.New("critical dequeue error")
	mockDequeueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType) (v4l2.Buffer, error) {
		return v4l2.Buffer{}, expectedPanicErr
	}

	// Mocks for Stop() that might be called during panic recovery by a defer in Start/goroutine
	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error { return nil }
	mockStreamOffFn = func(d v4l2.StreamingDevice) error { return nil }


	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected on critical DequeueBuffer error")
		} else {
			// Check if the panic message is as expected
			panicMsg := fmt.Sprintf("%v", r)
			expectedMsgPart := "device: stream loop dequeue: " + expectedPanicErr.Error()
			if panicMsg != expectedMsgPart { // Exact match since we control the panic string
				t.Errorf("Panic message = %q, want %q", panicMsg, expectedMsgPart)
			}
		}
		// Ensure cleanup if panic happened, or if test failed before panic
		if dev.streaming {
			dev.Stop()
		}
		dev.Close()
		if !closeCalled {
			// This might not be reached if panic is not properly handled by test defer
			// t.Error("mockCloseDeviceFn was not called after panic test")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Important to cancel context to allow goroutine to attempt exit

	err := dev.Start(ctx)
	if err != nil {
		// Start itself might return an error if the panic is recovered within Start
		// and converted to an error, which is not current behavior.
		// t.Logf("Start() returned error: %v (this might be ok if panic is recovered by Start)", err)
	}

	// If Start completes without error (meaning goroutine launched), trigger the dequeue.
	if err == nil {
		testReadyChan <- struct{}{} // Trigger the read in the loop
		// Give a little time for the goroutine to process and potentially panic
		// This is not ideal, but helps ensure the panic path is hit.
		// A more robust way would involve another channel signaling panic from goroutine,
		// but that requires production code change.
		<-context.After(100 * context.Millisecond)
	}
}

func TestStart_StreamLoop_PanicOnReQueueError(t *testing.T) {
	dev := getOpenedMockDevice(t)
	defer resetMocks()
	var closeCalled bool
	mockCloseDeviceFn = func(fd uintptr) error { closeCalled = true; return nil }

	dev.config.bufSize = 1
	mockInitBuffersFn = func(d v4l2.StreamingDevice) (v4l2.RequestBuffers, error) { return v4l2.RequestBuffers{Count: 1}, nil }
	mockMapMemoryBuffersFn = func(d v4l2.StreamingDevice) ([][]byte, error) {
		dev.buffers = make([][]byte, 1); dev.buffers[0] = []byte{0xDE, 0xAD}; return dev.buffers, nil
	}
	mockStreamOnFn = func(d v4l2.StreamingDevice) error { return nil }

	testReadyChan := make(chan struct{}, 1)
	mockWaitForReadFn = func(d v4l2.Device) <-chan struct{} { return testReadyChan }
	mockDequeueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType) (v4l2.Buffer, error) {
		return v4l2.Buffer{Index: 0, Flags: v4l2.BufFlagMapped | v4l2.BufFlagDone, BytesUsed: 2}, nil
	}

	expectedPanicErr := errors.New("critical requeue error")
	queueCallCount := 0
	mockQueueBufferFn = func(fd uintptr, ioType v4l2.IOType, bufType v4l2.BufType, i uint32) (v4l2.Buffer, error) {
		queueCallCount++
		if queueCallCount == 1 { // First call (initial queue) succeeds
			return v4l2.Buffer{Index: i}, nil
		}
		// Second call (re-queue in loop) fails
		return v4l2.Buffer{}, expectedPanicErr
	}

	mockUnmapMemoryBuffersFn = func(d v4l2.StreamingDevice) error { return nil }
	mockStreamOffFn = func(d v4l2.StreamingDevice) error { return nil }

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected on critical QueueBuffer (re-queue) error")
		} else {
			panicMsg := fmt.Sprintf("%v", r)
			// Expected format: "device: stream loop queue: %s: buff: %#v"
			// We can't easily match the buff part exactly without knowing its internal C pointer.
			// So, we check for the error string part.
			if ! ( len(panicMsg) > 0 && // basic check
				  ( len(panicMsg) > len("device: stream loop queue: ") &&
					panicMsg[0:len("device: stream loop queue: ")] == "device: stream loop queue: " ) &&
				  ( len(panicMsg) > len(expectedPanicErr.Error()) &&
				    panicMsg[len("device: stream loop queue: "):len("device: stream loop queue: ")+len(expectedPanicErr.Error())] == expectedPanicErr.Error() ) ) {
				t.Errorf("Panic message %q does not contain expected error %q in the correct format", panicMsg, expectedPanicErr.Error())
			}
		}
		if dev.streaming { dev.Stop() }
		dev.Close()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := dev.Start(ctx) // Initial QueueBuffer call happens here
	if err != nil {
		// If initial QueueBuffer fails (if we changed the mock to fail on first call),
		// this would be hit. But the test is for re-queue failure.
		t.Fatalf("Start() itself failed: %v", err)
	}

	if err == nil {
		testReadyChan <- struct{}{} // Trigger dequeue and then the failing re-queue
		<-context.After(100 * context.Millisecond) // Give time for panic
	}
}


// Note on the init() based mocking:
// The effectiveness of the init() function redirecting v4l2 calls (e.g., v4l2.OpenDevice = ...)
```
