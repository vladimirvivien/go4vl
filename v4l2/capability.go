package v4l2

import (
	"fmt"
	"unsafe"
)

// V4l2 video capability constants
// see https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L451

const (
	CapVideoCapture       = 0x00000001 // V4L2_CAP_VIDEO_CAPTURE
	CapVideoOutput        = 0x00000002 // V4L2_CAP_VIDEO_OUTPUT
	CapVideoOverlay       = 0x00000004 // V4L2_CAP_VIDEO_OVERLAY
	CapVBICapture         = 0x00000010 // V4L2_CAP_VBI_CAPTURE
	CapVBIOutput          = 0x00000020 // V4L2_CAP_VBI_OUTPUT
	CapSlicedVBICapture   = 0x00000040 // V4L2_CAP_SLICED_VBI_CAPTURE
	CapSlicedVBIOutput    = 0x00000080 // V4L2_CAP_SLICED_VBI_OUTPUT
	CapRDSCapture         = 0x00000100 // V4L2_CAP_RDS_CAPTURE
	CapVideoOutputOverlay = 0x00000200 // V4L2_CAP_VIDEO_OUTPUT_OVERLAY
	CapHWFrequencySeek    = 0x00000400 // V4L2_CAP_HW_FREQ_SEEK
	CapRDSOutput          = 0x00000800 // V4L2_CAP_RDS_OUTPUT

	CapVideoCaptureMPlane = 0x00001000 // V4L2_CAP_VIDEO_CAPTURE_MPLANE
	CapVideoOutputMPlane  = 0x00002000 // V4L2_CAP_VIDEO_OUTPUT_MPLANE
	CapVideoMem2MemMPlane = 0x00004000 // V4L2_CAP_VIDEO_M2M_MPLANE
	CapVideoMem2Mem       = 0x00008000 // V4L2_CAP_VIDEO_M2M

	CapTuner     = 0x00010000 // V4L2_CAP_TUNER
	CapAudio     = 0x00020000 // V4L2_CAP_AUDIO
	CapRadio     = 0x00040000 // V4L2_CAP_RADIO
	CapModulator = 0x00080000 // V4L2_CAP_MODULATOR

	CapSDRCapture        = 0x00100000 // V4L2_CAP_SDR_CAPTURE
	CapExtendedPixFormat = 0x00200000 // V4L2_CAP_EXT_PIX_FORMAT
	CapSDROutput         = 0x00400000 // V4L2_CAP_SDR_OUTPUT
	CapMetadataCapture   = 0x00800000 // V4L2_CAP_META_CAPTURE

	CapReadWrite      = 0x01000000 // V4L2_CAP_READWRITE
	CapAsyncIO        = 0x02000000 // V4L2_CAP_ASYNCIO
	CapStreaming      = 0x04000000 // V4L2_CAP_STREAMING
	CapMetadataOutput = 0x08000000 // V4L2_CAP_META_OUTPUT

	CapTouch             = 0x10000000 // V4L2_CAP_TOUCH
	CapIOMediaController = 0x20000000 // V4L2_CAP_IO_MC

	CapDeviceCapabilities = 0x80000000 // V4L2_CAP_DEVICE_CAPS
)

type CapabilityDesc struct {
	Cap  uint32
	Desc string
}

var (
	Capabilities = []CapabilityDesc{
		{Cap: CapVideoCapture, Desc: "video capture (single-planar)"},
		{Cap: CapVideoOutput, Desc: "video output (single-planar)"},
		{Cap: CapVideoOverlay, Desc: "video overlay"},
		{Cap: CapVBICapture, Desc: "raw VBI capture"},
		{Cap: CapVBIOutput, Desc: "raw VBI output"},
		{Cap: CapSlicedVBICapture, Desc: "sliced VBI capture"},
		{Cap: CapSlicedVBIOutput, Desc: "sliced VBI output"},
		{Cap: CapRDSCapture, Desc: "RDS capture"},
		{Cap: CapVideoOutputOverlay, Desc: "video output overlay"},
		{Cap: CapHWFrequencySeek, Desc: "hardware frequency seeking"},
		{Cap: CapRDSOutput, Desc: "RDS output"},

		{Cap: CapVideoCaptureMPlane, Desc: "video capture (multi-planar)"},
		{Cap: CapVideoOutputMPlane, Desc: "video output (multi-planar)"},
		{Cap: CapVideoMem2MemMPlane, Desc: "memory-to-memory video (multi-planar)"},
		{Cap: CapVideoMem2Mem, Desc: "memory-to-memory video (single-planar)"},

		{Cap: CapTuner, Desc: "video tuner"},
		{Cap: CapAudio, Desc: "audio inputs or outputs"},
		{Cap: CapRadio, Desc: "radio receiver"},
		{Cap: CapModulator, Desc: "radio frequency modulator"},

		{Cap: CapSDRCapture, Desc: "SDR capture"},
		{Cap: CapExtendedPixFormat, Desc: "extended pixel format"},
		{Cap: CapSDROutput, Desc: "SDR output"},
		{Cap: CapMetadataCapture, Desc: "metadata capture"},

		{Cap: CapReadWrite, Desc: "read/write IO"},
		{Cap: CapAsyncIO, Desc: "asynchronous IO"},
		{Cap: CapStreaming, Desc: "streaming IO"},
		{Cap: CapMetadataOutput, Desc: "metadata output"},

		{Cap: CapTouch, Desc: "touch capability"},
		{Cap: CapIOMediaController, Desc: "IO media controller"},

		{Cap: CapDeviceCapabilities, Desc: "device capabilities"},
	}
)

// v4l2Capability type for device (see v4l2_capability)
// This type stores the capability information returned by the device.
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-querycap.html#c.V4L.v4l2_capability
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L440
type v4l2Capability struct {
	driver       [16]uint8
	card         [32]uint8
	busInfo      [32]uint8
	version      uint32
	capabilities uint32
	deviceCaps   uint32
	reserved     [3]uint32
}

// Capability represents capabilities retrieved for the device.
// Use attached methods on this type to access capabilities.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-querycap.html#c.V4L.v4l2_capability
type Capability struct {
	v4l2Cap v4l2Capability
}

// GetCapability retrieves capability info for device
func GetCapability(fd uintptr) (Capability, error) {
	v4l2Cap := v4l2Capability{}
	if err := Send(fd, VidiocQueryCap, uintptr(unsafe.Pointer(&v4l2Cap))); err != nil {
		return Capability{}, fmt.Errorf("capability: %w", err)
	}
	return Capability{v4l2Cap: v4l2Cap}, nil
}

// GetCapabilities returns the capability mask as a union of
// all exported capabilities for the physical device (opened or not).
// Use this method to access capabilities.
func (c Capability) GetCapabilities() uint32 {
	return c.v4l2Cap.capabilities
}

// GetDeviceCaps returns the capability mask for the open device.
// This is a subset of capabilities returned by GetCapabilities.
func (c Capability) GetDeviceCaps() uint32 {
	return c.v4l2Cap.deviceCaps
}

// IsVideoCaptureSupported returns caps & CapVideoCapture
func (c Capability) IsVideoCaptureSupported() bool {
	return (c.v4l2Cap.capabilities & CapVideoCapture) != 0
}

// IsVideoOutputSupported returns caps & CapVideoOutput
func (c Capability) IsVideoOutputSupported() bool {
	return (c.v4l2Cap.capabilities & CapVideoOutput) != 0
}

// IsVideoOverlaySupported returns caps & CapVideoOverlay
func (c Capability) IsVideoOverlaySupported() bool {
	return (c.v4l2Cap.capabilities & CapVideoOverlay) != 0
}

// IsVideoOutputOverlaySupported returns caps & CapVideoOutputOverlay
func (c Capability) IsVideoOutputOverlaySupported() bool {
	return (c.v4l2Cap.capabilities & CapVideoOutputOverlay) != 0
}

// IsVideoCaptureMultiplanarSupported returns caps & CapVideoCaptureMPlane
func (c Capability) IsVideoCaptureMultiplanarSupported() bool {
	return (c.v4l2Cap.capabilities & CapVideoCaptureMPlane) != 0
}

// IsVideoOutputMultiplanerSupported returns caps & CapVideoOutputMPlane
func (c Capability) IsVideoOutputMultiplanerSupported() bool {
	return (c.v4l2Cap.capabilities & CapVideoOutputMPlane) != 0
}

// IsReadWriteSupported returns caps & CapReadWrite
func (c Capability) IsReadWriteSupported() bool {
	return (c.v4l2Cap.capabilities & CapReadWrite) != 0
}

// IsStreamingSupported returns caps & CapStreaming
func (c Capability) IsStreamingSupported() bool {
	return (c.v4l2Cap.capabilities & CapStreaming) != 0
}

// IsDeviceCapabilitiesProvided returns true if the device returns
// device-specific capabilities (via CapDeviceCapabilities)
// See notes on VL42_CAP_DEVICE_CAPS:
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-querycap.html?highlight=v4l2_cap_device_caps
func (c Capability) IsDeviceCapabilitiesProvided() bool {
	return (c.v4l2Cap.capabilities & CapDeviceCapabilities) != 0
}

// GetDriverCapDescriptions return textual descriptions of driver capabilities
func (c Capability) GetDriverCapDescriptions() []CapabilityDesc {
	var result []CapabilityDesc
	for _, cap := range Capabilities {
		if c.GetCapabilities() & cap.Cap == cap.Cap {
			result = append(result, cap)
		}
	}
	return result
}

// GetDeviceCapDescriptions return textual descriptions of device capabilities
func (c Capability) GetDeviceCapDescriptions() []CapabilityDesc {
	var result []CapabilityDesc
	for _, cap := range Capabilities {
		if c.GetDeviceCaps() & cap.Cap == cap.Cap {
			result = append(result, cap)
		}
	}
	return result
}

// DriverName returns a string value for the driver name
func (c Capability) DriverName() string {
	return toGoString(c.v4l2Cap.driver[:])
}

// CardName returns a string value for device's card
func (c Capability) CardName() string {
	return toGoString(c.v4l2Cap.card[:])
}

// BusInfo returns the device's bus info
func (c Capability) BusInfo() string {
	return toGoString(c.v4l2Cap.busInfo[:])
}

// GetVersion returns the device's version
func (c Capability) GetVersion() uint32 {
	return c.v4l2Cap.version
}

// String returns a string value representing driver information
func (c Capability) String() string {
	return fmt.Sprintf("driver: %s; card: %s; bus info: %s", c.DriverName(), c.CardName(), c.BusInfo())
}
