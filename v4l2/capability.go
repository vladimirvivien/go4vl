package v4l2

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../include/
#include <linux/videodev2.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// V4l2 video capability constants
// see https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L451

const (
	CapVideoCapture       uint32 = C.V4L2_CAP_VIDEO_CAPTURE
	CapVideoOutput        uint32 = C.V4L2_CAP_VIDEO_OUTPUT
	CapVideoOverlay       uint32 = C.V4L2_CAP_VIDEO_OVERLAY
	CapVBICapture         uint32 = C.V4L2_CAP_VBI_CAPTURE
	CapVBIOutput          uint32 = C.V4L2_CAP_VBI_OUTPUT
	CapSlicedVBICapture   uint32 = C.V4L2_CAP_SLICED_VBI_CAPTURE
	CapSlicedVBIOutput    uint32 = C.V4L2_CAP_SLICED_VBI_OUTPUT
	CapRDSCapture         uint32 = C.V4L2_CAP_RDS_CAPTURE
	CapVideoOutputOverlay uint32 = C.V4L2_CAP_VIDEO_OUTPUT_OVERLAY
	CapHWFrequencySeek    uint32 = C.V4L2_CAP_HW_FREQ_SEEK
	CapRDSOutput          uint32 = C.V4L2_CAP_RDS_OUTPUT

	CapVideoCaptureMPlane uint32 = C.V4L2_CAP_VIDEO_CAPTURE_MPLANE
	CapVideoOutputMPlane  uint32 = C.V4L2_CAP_VIDEO_OUTPUT_MPLANE
	CapVideoMem2MemMPlane uint32 = C.V4L2_CAP_VIDEO_M2M_MPLANE
	CapVideoMem2Mem       uint32 = C.V4L2_CAP_VIDEO_M2M

	CapTuner     uint32 = C.V4L2_CAP_TUNER
	CapAudio     uint32 = C.V4L2_CAP_AUDIO
	CapRadio     uint32 = C.V4L2_CAP_RADIO
	CapModulator uint32 = C.V4L2_CAP_MODULATOR

	CapSDRCapture        uint32 = C.V4L2_CAP_SDR_CAPTURE
	CapExtendedPixFormat uint32 = C.V4L2_CAP_EXT_PIX_FORMAT
	CapSDROutput         uint32 = C.V4L2_CAP_SDR_OUTPUT
	CapMetadataCapture   uint32 = C.V4L2_CAP_META_CAPTURE

	CapReadWrite uint32 = C.V4L2_CAP_READWRITE
	CapAsyncIO   uint32 = C.V4L2_CAP_ASYNCIO
	CapStreaming uint32 = C.V4L2_CAP_STREAMING

	CapMetadataOutput     uint32 = C.V4L2_CAP_META_OUTPUT
	CapTouch              uint32 = C.V4L2_CAP_TOUCH
	CapIOMediaController  uint32 = C.V4L2_CAP_IO_MC
	CapDeviceCapabilities uint32 = C.V4L2_CAP_DEVICE_CAPS
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

// Capability represents capabilities retrieved for the device (see v4l2_capability).
// Use attached methods on this type to access capabilities.
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L440
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-querycap.html#c.V4L.v4l2_capability
type Capability struct {
	// Driver name of the driver module
	Driver string

	// Card name of the device card
	Card string

	// BusInfo is the name of the device bus
	BusInfo string

	// Version is the kernel version
	Version uint32

	// Capabilities returns all exported capabilities for the physical device (opened or not)
	Capabilities uint32

	// DeviceCapabilities is the capability for this particular (opened) device or node
	DeviceCapabilities uint32
}

// GetCapability retrieves capability info for device
func GetCapability(fd uintptr) (Capability, error) {
	var v4l2Cap C.struct_v4l2_capability
	if err := send(fd, C.VIDIOC_QUERYCAP, uintptr(unsafe.Pointer(&v4l2Cap))); err != nil {
		return Capability{}, fmt.Errorf("capability: %w", err)
	}
	return Capability{
		Driver:             C.GoString((*C.char)(unsafe.Pointer(&v4l2Cap.driver[0]))),
		Card:               C.GoString((*C.char)(unsafe.Pointer(&v4l2Cap.card[0]))),
		BusInfo:            C.GoString((*C.char)(unsafe.Pointer(&v4l2Cap.bus_info[0]))),
		Version:            uint32(v4l2Cap.version),
		Capabilities:       uint32(v4l2Cap.capabilities),
		DeviceCapabilities: uint32(v4l2Cap.device_caps),
	}, nil
}

// GetCapabilities returns device capabilities if supported
func (c Capability) GetCapabilities() uint32 {
	if c.IsDeviceCapabilitiesProvided() {
		return c.DeviceCapabilities
	}
	return c.Capabilities
}

// IsVideoCaptureSupported returns caps & CapVideoCapture
func (c Capability) IsVideoCaptureSupported() bool {
	return c.Capabilities&CapVideoCapture != 0
}

// IsVideoOutputSupported returns caps & CapVideoOutput
func (c Capability) IsVideoOutputSupported() bool {
	return c.Capabilities&CapVideoOutput != 0
}

// IsVideoOverlaySupported returns caps & CapVideoOverlay
func (c Capability) IsVideoOverlaySupported() bool {
	return c.Capabilities&CapVideoOverlay != 0
}

// IsVideoOutputOverlaySupported returns caps & CapVideoOutputOverlay
func (c Capability) IsVideoOutputOverlaySupported() bool {
	return c.Capabilities&CapVideoOutputOverlay != 0
}

// IsVideoCaptureMultiplanarSupported returns caps & CapVideoCaptureMPlane
func (c Capability) IsVideoCaptureMultiplanarSupported() bool {
	return c.Capabilities&CapVideoCaptureMPlane != 0
}

// IsVideoOutputMultiplanerSupported returns caps & CapVideoOutputMPlane
func (c Capability) IsVideoOutputMultiplanerSupported() bool {
	return c.Capabilities&CapVideoOutputMPlane != 0
}

// IsReadWriteSupported returns caps & CapReadWrite
func (c Capability) IsReadWriteSupported() bool {
	return c.Capabilities&CapReadWrite != 0
}

// IsStreamingSupported returns caps & CapStreaming
func (c Capability) IsStreamingSupported() bool {
	return c.Capabilities&CapStreaming != 0
}

// IsDeviceCapabilitiesProvided returns true if the device returns
// device-specific capabilities (via CapDeviceCapabilities)
// See notes on VL42_CAP_DEVICE_CAPS:
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-querycap.html?highlight=v4l2_cap_device_caps
func (c Capability) IsDeviceCapabilitiesProvided() bool {
	return c.Capabilities&CapDeviceCapabilities != 0
}

// GetDriverCapDescriptions return textual descriptions of driver capabilities
func (c Capability) GetDriverCapDescriptions() []CapabilityDesc {
	var result []CapabilityDesc
	for _, cap := range Capabilities {
		if c.Capabilities&cap.Cap == cap.Cap {
			result = append(result, cap)
		}
	}
	return result
}

// GetDeviceCapDescriptions return textual descriptions of device capabilities
func (c Capability) GetDeviceCapDescriptions() []CapabilityDesc {
	var result []CapabilityDesc
	for _, cap := range Capabilities {
		if c.DeviceCapabilities&cap.Cap == cap.Cap {
			result = append(result, cap)
		}
	}
	return result
}

func (c Capability) GetVersionInfo() VersionInfo {
	return VersionInfo{value: c.Version}
}

// String returns a string value representing driver information
func (c Capability) String() string {
	return fmt.Sprintf("driver: %s; card: %s; bus info: %s", c.Driver, c.Card, c.BusInfo)
}
