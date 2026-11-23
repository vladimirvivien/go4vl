package v4l2

import (
	"testing"
)

var (
	benchCapUint32 uint32
	benchCapBool   bool
	benchCapDescs  []CapabilityDesc
	benchCapVer    VersionInfo
	benchCapStr    string
)

// BenchmarkCapability_GetCapabilities benchmarks the GetCapabilities method
func BenchmarkCapability_GetCapabilities(b *testing.B) {
	cap := Capability{
		Capabilities:       CapVideoCapture | CapVideoOutput | CapDeviceCapabilities,
		DeviceCapabilities: CapVideoCapture | CapStreaming,
	}

	b.ResetTimer()
	var sink uint32
	for i := 0; i < b.N; i++ {
		sink = cap.GetCapabilities()
	}
	benchCapUint32 = sink
}

// BenchmarkCapability_IsVideoCaptureSupported benchmarks capability checking
func BenchmarkCapability_IsVideoCaptureSupported(b *testing.B) {
	cap := Capability{
		Capabilities: CapVideoCapture | CapStreaming | CapReadWrite,
	}

	b.ResetTimer()
	var sink bool
	for i := 0; i < b.N; i++ {
		sink = cap.IsVideoCaptureSupported()
	}
	benchCapBool = sink
}

// BenchmarkCapability_IsStreamingSupported benchmarks streaming capability check
func BenchmarkCapability_IsStreamingSupported(b *testing.B) {
	cap := Capability{
		Capabilities: CapVideoCapture | CapStreaming | CapReadWrite,
	}

	b.ResetTimer()
	var sink bool
	for i := 0; i < b.N; i++ {
		sink = cap.IsStreamingSupported()
	}
	benchCapBool = sink
}

// BenchmarkCapability_GetDriverCapDescriptions benchmarks description generation
func BenchmarkCapability_GetDriverCapDescriptions(b *testing.B) {
	cap := Capability{
		Capabilities: CapVideoCapture | CapVideoOutput | CapStreaming | CapReadWrite | CapDeviceCapabilities,
	}

	b.ResetTimer()
	var sink []CapabilityDesc
	for i := 0; i < b.N; i++ {
		sink = cap.GetDriverCapDescriptions()
	}
	benchCapDescs = sink
}

// BenchmarkCapability_GetVersionInfo benchmarks version info extraction
func BenchmarkCapability_GetVersionInfo(b *testing.B) {
	cap := Capability{
		Version: (5 << 16) | (15 << 8) | 0,
	}

	b.ResetTimer()
	var sink VersionInfo
	for i := 0; i < b.N; i++ {
		sink = cap.GetVersionInfo()
	}
	benchCapVer = sink
}

// BenchmarkCapability_String benchmarks string formatting
func BenchmarkCapability_String(b *testing.B) {
	cap := Capability{
		Driver:  "uvcvideo",
		Card:    "HD Webcam C920",
		BusInfo: "usb-0000:00:14.0-1",
	}

	b.ResetTimer()
	var sink string
	for i := 0; i < b.N; i++ {
		sink = cap.String()
	}
	benchCapStr = sink
}

// BenchmarkCapability_MultipleChecks benchmarks common usage pattern
func BenchmarkCapability_MultipleChecks(b *testing.B) {
	cap := Capability{
		Capabilities:       CapVideoCapture | CapVideoOutput | CapStreaming | CapDeviceCapabilities,
		DeviceCapabilities: CapVideoCapture | CapStreaming,
	}

	b.ResetTimer()
	var sink bool
	for i := 0; i < b.N; i++ {
		sink = cap.IsVideoCaptureSupported() &&
			cap.IsStreamingSupported() &&
			cap.GetCapabilities() != 0 &&
			cap.IsDeviceCapabilitiesProvided()
	}
	benchCapBool = sink
}
