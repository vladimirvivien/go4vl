package v4l2

import (
	"testing"
)

// BenchmarkCapability_GetCapabilities benchmarks the GetCapabilities method
func BenchmarkCapability_GetCapabilities(b *testing.B) {
	cap := Capability{
		Capabilities:       CapVideoCapture | CapVideoOutput | CapDeviceCapabilities,
		DeviceCapabilities: CapVideoCapture | CapStreaming,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cap.GetCapabilities()
	}
}

// BenchmarkCapability_IsVideoCaptureSupported benchmarks capability checking
func BenchmarkCapability_IsVideoCaptureSupported(b *testing.B) {
	cap := Capability{
		Capabilities: CapVideoCapture | CapStreaming | CapReadWrite,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cap.IsVideoCaptureSupported()
	}
}

// BenchmarkCapability_IsStreamingSupported benchmarks streaming capability check
func BenchmarkCapability_IsStreamingSupported(b *testing.B) {
	cap := Capability{
		Capabilities: CapVideoCapture | CapStreaming | CapReadWrite,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cap.IsStreamingSupported()
	}
}

// BenchmarkCapability_GetDriverCapDescriptions benchmarks description generation
func BenchmarkCapability_GetDriverCapDescriptions(b *testing.B) {
	cap := Capability{
		Capabilities: CapVideoCapture | CapVideoOutput | CapStreaming | CapReadWrite | CapDeviceCapabilities,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cap.GetDriverCapDescriptions()
	}
}

// BenchmarkCapability_GetVersionInfo benchmarks version info extraction
func BenchmarkCapability_GetVersionInfo(b *testing.B) {
	cap := Capability{
		Version: (5 << 16) | (15 << 8) | 0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cap.GetVersionInfo()
	}
}

// BenchmarkCapability_String benchmarks string formatting
func BenchmarkCapability_String(b *testing.B) {
	cap := Capability{
		Driver:  "uvcvideo",
		Card:    "HD Webcam C920",
		BusInfo: "usb-0000:00:14.0-1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cap.String()
	}
}

// BenchmarkCapability_MultipleChecks benchmarks common usage pattern
func BenchmarkCapability_MultipleChecks(b *testing.B) {
	cap := Capability{
		Capabilities:       CapVideoCapture | CapVideoOutput | CapStreaming | CapDeviceCapabilities,
		DeviceCapabilities: CapVideoCapture | CapStreaming,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Typical usage: check multiple capabilities
		_ = cap.IsVideoCaptureSupported()
		_ = cap.IsStreamingSupported()
		_ = cap.GetCapabilities()
		_ = cap.IsDeviceCapabilitiesProvided()
	}
}
