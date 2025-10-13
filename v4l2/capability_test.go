package v4l2

import (
	"strings"
	"testing"
)

// TestCapability_GetCapabilities tests the GetCapabilities method
func TestCapability_GetCapabilities(t *testing.T) {
	tests := []struct {
		name               string
		cap                Capability
		expectedResult     uint32
		shouldUseDeviceCap bool
	}{
		{
			name: "modern driver with device capabilities",
			cap: Capability{
				Capabilities:       CapVideoCapture | CapVideoOutput | CapDeviceCapabilities,
				DeviceCapabilities: CapVideoCapture | CapStreaming,
			},
			expectedResult:     CapVideoCapture | CapStreaming,
			shouldUseDeviceCap: true,
		},
		{
			name: "legacy driver without device capabilities",
			cap: Capability{
				Capabilities:       CapVideoCapture | CapStreaming,
				DeviceCapabilities: 0,
			},
			expectedResult:     CapVideoCapture | CapStreaming,
			shouldUseDeviceCap: false,
		},
		{
			name: "multi-function device with device caps",
			cap: Capability{
				Capabilities:       CapVideoCapture | CapVideoOutput | CapVBICapture | CapDeviceCapabilities,
				DeviceCapabilities: CapVideoCapture | CapStreaming,
			},
			expectedResult:     CapVideoCapture | CapStreaming,
			shouldUseDeviceCap: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.GetCapabilities()
			if result != tt.expectedResult {
				t.Errorf("GetCapabilities() = 0x%08x, want 0x%08x", result, tt.expectedResult)
			}

			// Verify it uses the right source
			if tt.shouldUseDeviceCap {
				if result != tt.cap.DeviceCapabilities {
					t.Error("Should have used DeviceCapabilities but didn't")
				}
			} else {
				if result != tt.cap.Capabilities {
					t.Error("Should have used Capabilities but didn't")
				}
			}
		})
	}
}

// TestCapability_IsVideoCaptureSupported tests video capture capability detection
func TestCapability_IsVideoCaptureSupported(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		expected bool
	}{
		{
			name:     "video capture supported",
			cap:      Capability{Capabilities: CapVideoCapture},
			expected: true,
		},
		{
			name:     "video capture not supported",
			cap:      Capability{Capabilities: CapVideoOutput},
			expected: false,
		},
		{
			name:     "video capture with other caps",
			cap:      Capability{Capabilities: CapVideoCapture | CapStreaming | CapReadWrite},
			expected: true,
		},
		{
			name:     "no capabilities",
			cap:      Capability{Capabilities: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.IsVideoCaptureSupported()
			if result != tt.expected {
				t.Errorf("IsVideoCaptureSupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCapability_IsVideoOutputSupported tests video output capability detection
func TestCapability_IsVideoOutputSupported(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		expected bool
	}{
		{
			name:     "video output supported",
			cap:      Capability{Capabilities: CapVideoOutput},
			expected: true,
		},
		{
			name:     "video output not supported",
			cap:      Capability{Capabilities: CapVideoCapture},
			expected: false,
		},
		{
			name:     "video output with other caps",
			cap:      Capability{Capabilities: CapVideoOutput | CapStreaming},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.IsVideoOutputSupported()
			if result != tt.expected {
				t.Errorf("IsVideoOutputSupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCapability_IsStreamingSupported tests streaming I/O capability detection
func TestCapability_IsStreamingSupported(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		expected bool
	}{
		{
			name:     "streaming supported",
			cap:      Capability{Capabilities: CapStreaming},
			expected: true,
		},
		{
			name:     "streaming not supported",
			cap:      Capability{Capabilities: CapReadWrite},
			expected: false,
		},
		{
			name:     "streaming with video capture",
			cap:      Capability{Capabilities: CapVideoCapture | CapStreaming},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.IsStreamingSupported()
			if result != tt.expected {
				t.Errorf("IsStreamingSupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCapability_IsReadWriteSupported tests read/write I/O capability detection
func TestCapability_IsReadWriteSupported(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		expected bool
	}{
		{
			name:     "read/write supported",
			cap:      Capability{Capabilities: CapReadWrite},
			expected: true,
		},
		{
			name:     "read/write not supported",
			cap:      Capability{Capabilities: CapStreaming},
			expected: false,
		},
		{
			name:     "read/write with streaming",
			cap:      Capability{Capabilities: CapReadWrite | CapStreaming},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.IsReadWriteSupported()
			if result != tt.expected {
				t.Errorf("IsReadWriteSupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCapability_IsDeviceCapabilitiesProvided tests device capabilities flag detection
func TestCapability_IsDeviceCapabilitiesProvided(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		expected bool
	}{
		{
			name:     "device capabilities provided",
			cap:      Capability{Capabilities: CapDeviceCapabilities},
			expected: true,
		},
		{
			name:     "device capabilities not provided",
			cap:      Capability{Capabilities: CapVideoCapture},
			expected: false,
		},
		{
			name: "device capabilities with other flags",
			cap: Capability{
				Capabilities: CapVideoCapture | CapStreaming | CapDeviceCapabilities,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.IsDeviceCapabilitiesProvided()
			if result != tt.expected {
				t.Errorf("IsDeviceCapabilitiesProvided() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCapability_IsVideoOverlaySupported tests video overlay capability detection
func TestCapability_IsVideoOverlaySupported(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		expected bool
	}{
		{
			name:     "video overlay supported",
			cap:      Capability{Capabilities: CapVideoOverlay},
			expected: true,
		},
		{
			name:     "video overlay not supported",
			cap:      Capability{Capabilities: CapVideoCapture},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.IsVideoOverlaySupported()
			if result != tt.expected {
				t.Errorf("IsVideoOverlaySupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCapability_IsVideoOutputOverlaySupported tests video output overlay capability detection
func TestCapability_IsVideoOutputOverlaySupported(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		expected bool
	}{
		{
			name:     "video output overlay supported",
			cap:      Capability{Capabilities: CapVideoOutputOverlay},
			expected: true,
		},
		{
			name:     "video output overlay not supported",
			cap:      Capability{Capabilities: CapVideoCapture},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.IsVideoOutputOverlaySupported()
			if result != tt.expected {
				t.Errorf("IsVideoOutputOverlaySupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCapability_IsVideoCaptureMultiplanarSupported tests multi-planar capture capability
func TestCapability_IsVideoCaptureMultiplanarSupported(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		expected bool
	}{
		{
			name:     "multi-planar capture supported",
			cap:      Capability{Capabilities: CapVideoCaptureMPlane},
			expected: true,
		},
		{
			name:     "multi-planar capture not supported",
			cap:      Capability{Capabilities: CapVideoCapture},
			expected: false,
		},
		{
			name:     "both single and multi-planar",
			cap:      Capability{Capabilities: CapVideoCapture | CapVideoCaptureMPlane},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.IsVideoCaptureMultiplanarSupported()
			if result != tt.expected {
				t.Errorf("IsVideoCaptureMultiplanarSupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCapability_IsVideoOutputMultiplanarSupported tests multi-planar output capability
func TestCapability_IsVideoOutputMultiplanerSupported(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		expected bool
	}{
		{
			name:     "multi-planar output supported",
			cap:      Capability{Capabilities: CapVideoOutputMPlane},
			expected: true,
		},
		{
			name:     "multi-planar output not supported",
			cap:      Capability{Capabilities: CapVideoOutput},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.IsVideoOutputMultiplanerSupported()
			if result != tt.expected {
				t.Errorf("IsVideoOutputMultiplanerSupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCapability_GetDriverCapDescriptions tests driver capability descriptions
func TestCapability_GetDriverCapDescriptions(t *testing.T) {
	tests := []struct {
		name          string
		cap           Capability
		expectedCount int
		checkCaps     []uint32
	}{
		{
			name:          "single capability",
			cap:           Capability{Capabilities: CapVideoCapture},
			expectedCount: 1,
			checkCaps:     []uint32{CapVideoCapture},
		},
		{
			name:          "multiple capabilities",
			cap:           Capability{Capabilities: CapVideoCapture | CapStreaming | CapReadWrite},
			expectedCount: 3,
			checkCaps:     []uint32{CapVideoCapture, CapStreaming, CapReadWrite},
		},
		{
			name:          "no capabilities",
			cap:           Capability{Capabilities: 0},
			expectedCount: 0,
			checkCaps:     []uint32{},
		},
		{
			name:          "complex device",
			cap:           Capability{Capabilities: CapVideoCapture | CapVideoOutput | CapStreaming | CapDeviceCapabilities},
			expectedCount: 4,
			checkCaps:     []uint32{CapVideoCapture, CapVideoOutput, CapStreaming, CapDeviceCapabilities},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			descs := tt.cap.GetDriverCapDescriptions()

			if len(descs) != tt.expectedCount {
				t.Errorf("GetDriverCapDescriptions() returned %d descriptions, want %d", len(descs), tt.expectedCount)
			}

			// Verify expected caps are present
			foundCaps := make(map[uint32]bool)
			for _, desc := range descs {
				foundCaps[desc.Cap] = true
				if desc.Desc == "" {
					t.Errorf("Description for cap 0x%08x should not be empty", desc.Cap)
				}
			}

			for _, expectedCap := range tt.checkCaps {
				if !foundCaps[expectedCap] {
					t.Errorf("Expected capability 0x%08x not found in descriptions", expectedCap)
				}
			}
		})
	}
}

// TestCapability_GetDeviceCapDescriptions tests device capability descriptions
func TestCapability_GetDeviceCapDescriptions(t *testing.T) {
	tests := []struct {
		name          string
		cap           Capability
		expectedCount int
		checkCaps     []uint32
	}{
		{
			name: "device capabilities only",
			cap: Capability{
				Capabilities:       CapVideoCapture | CapVideoOutput | CapDeviceCapabilities,
				DeviceCapabilities: CapVideoCapture | CapStreaming,
			},
			expectedCount: 2,
			checkCaps:     []uint32{CapVideoCapture, CapStreaming},
		},
		{
			name: "no device capabilities set",
			cap: Capability{
				Capabilities:       CapVideoCapture,
				DeviceCapabilities: 0,
			},
			expectedCount: 0,
			checkCaps:     []uint32{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			descs := tt.cap.GetDeviceCapDescriptions()

			if len(descs) != tt.expectedCount {
				t.Errorf("GetDeviceCapDescriptions() returned %d descriptions, want %d", len(descs), tt.expectedCount)
			}

			// Verify expected caps are present
			foundCaps := make(map[uint32]bool)
			for _, desc := range descs {
				foundCaps[desc.Cap] = true
			}

			for _, expectedCap := range tt.checkCaps {
				if !foundCaps[expectedCap] {
					t.Errorf("Expected capability 0x%08x not found in descriptions", expectedCap)
				}
			}
		})
	}
}

// TestCapability_GetVersionInfo tests version info extraction
func TestCapability_GetVersionInfo(t *testing.T) {
	tests := []struct {
		name          string
		version       uint32
		expectedMajor uint32
		expectedMinor uint32
		expectedPatch uint32
	}{
		{
			name:          "kernel 5.15.0",
			version:       (5 << 16) | (15 << 8) | 0,
			expectedMajor: 5,
			expectedMinor: 15,
			expectedPatch: 0,
		},
		{
			name:          "kernel 6.1.25",
			version:       (6 << 16) | (1 << 8) | 25,
			expectedMajor: 6,
			expectedMinor: 1,
			expectedPatch: 25,
		},
		{
			name:          "kernel 4.19.255",
			version:       (4 << 16) | (19 << 8) | 255,
			expectedMajor: 4,
			expectedMinor: 19,
			expectedPatch: 255,
		},
		{
			name:          "zero version",
			version:       0,
			expectedMajor: 0,
			expectedMinor: 0,
			expectedPatch: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cap := Capability{Version: tt.version}
			versionInfo := cap.GetVersionInfo()

			if versionInfo.Major() != tt.expectedMajor {
				t.Errorf("Major() = %d, want %d", versionInfo.Major(), tt.expectedMajor)
			}
			if versionInfo.Minor() != tt.expectedMinor {
				t.Errorf("Minor() = %d, want %d", versionInfo.Minor(), tt.expectedMinor)
			}
			if versionInfo.Patch() != tt.expectedPatch {
				t.Errorf("Patch() = %d, want %d", versionInfo.Patch(), tt.expectedPatch)
			}

			// Test String() method
			versionStr := versionInfo.String()
			if versionStr == "" && tt.version != 0 {
				t.Error("String() should not be empty for non-zero version")
			}
		})
	}
}

// TestCapability_String tests the String method
func TestCapability_String(t *testing.T) {
	tests := []struct {
		name     string
		cap      Capability
		contains []string
	}{
		{
			name: "basic capability",
			cap: Capability{
				Driver:  "uvcvideo",
				Card:    "HD Webcam",
				BusInfo: "usb-0000:00:14.0-1",
			},
			contains: []string{"uvcvideo", "HD Webcam", "usb-0000:00:14.0-1"},
		},
		{
			name: "v4l2loopback",
			cap: Capability{
				Driver:  "v4l2loopback",
				Card:    "Dummy video device",
				BusInfo: "platform:v4l2loopback-000",
			},
			contains: []string{"v4l2loopback", "Dummy video device", "platform:v4l2loopback-000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.String()

			if result == "" {
				t.Error("String() should not be empty")
			}

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("String() = %q, should contain %q", result, expected)
				}
			}
		})
	}
}

// TestCapabilityDesc_Complete tests all capability descriptions are complete
func TestCapabilityDesc_Complete(t *testing.T) {
	// Verify all entries in Capabilities have valid data
	for i, capDesc := range Capabilities {
		t.Run(capDesc.Desc, func(t *testing.T) {
			if capDesc.Cap == 0 {
				t.Errorf("Capabilities[%d]: Cap should not be zero", i)
			}
			if capDesc.Desc == "" {
				t.Errorf("Capabilities[%d]: Desc should not be empty", i)
			}
		})
	}

	// Verify we have a reasonable number of capability descriptions
	if len(Capabilities) < 10 {
		t.Errorf("Expected at least 10 capability descriptions, got %d", len(Capabilities))
	}
}

// TestCapabilityConstants tests that capability constants are defined correctly
func TestCapabilityConstants(t *testing.T) {
	tests := []struct {
		name  string
		value uint32
	}{
		{"CapVideoCapture", CapVideoCapture},
		{"CapVideoOutput", CapVideoOutput},
		{"CapVideoOverlay", CapVideoOverlay},
		{"CapVBICapture", CapVBICapture},
		{"CapVBIOutput", CapVBIOutput},
		{"CapSlicedVBICapture", CapSlicedVBICapture},
		{"CapSlicedVBIOutput", CapSlicedVBIOutput},
		{"CapRDSCapture", CapRDSCapture},
		{"CapVideoOutputOverlay", CapVideoOutputOverlay},
		{"CapHWFrequencySeek", CapHWFrequencySeek},
		{"CapRDSOutput", CapRDSOutput},
		{"CapVideoCaptureMPlane", CapVideoCaptureMPlane},
		{"CapVideoOutputMPlane", CapVideoOutputMPlane},
		{"CapVideoMem2MemMPlane", CapVideoMem2MemMPlane},
		{"CapVideoMem2Mem", CapVideoMem2Mem},
		{"CapTuner", CapTuner},
		{"CapAudio", CapAudio},
		{"CapRadio", CapRadio},
		{"CapModulator", CapModulator},
		{"CapSDRCapture", CapSDRCapture},
		{"CapExtendedPixFormat", CapExtendedPixFormat},
		{"CapSDROutput", CapSDROutput},
		{"CapMetadataCapture", CapMetadataCapture},
		{"CapReadWrite", CapReadWrite},
		{"CapAsyncIO", CapAsyncIO},
		{"CapStreaming", CapStreaming},
		{"CapMetadataOutput", CapMetadataOutput},
		{"CapTouch", CapTouch},
		{"CapIOMediaController", CapIOMediaController},
		{"CapDeviceCapabilities", CapDeviceCapabilities},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}
