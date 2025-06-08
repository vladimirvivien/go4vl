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

// Capability constants represent the various capabilities a V4L2 device can support.
// These values are used in the Capabilities and DeviceCapabilities fields of the Capability struct.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L451
const (
	// CapVideoCapture indicates the device supports video capture (single-planar API).
	CapVideoCapture uint32 = C.V4L2_CAP_VIDEO_CAPTURE
	// CapVideoOutput indicates the device supports video output (single-planar API).
	CapVideoOutput uint32 = C.V4L2_CAP_VIDEO_OUTPUT
	// CapVideoOverlay indicates the device supports video overlay.
	CapVideoOverlay uint32 = C.V4L2_CAP_VIDEO_OVERLAY
	// CapVBICapture indicates the device supports raw VBI (Vertical Blanking Interval) capture.
	CapVBICapture uint32 = C.V4L2_CAP_VBI_CAPTURE
	// CapVBIOutput indicates the device supports raw VBI output.
	CapVBIOutput uint32 = C.V4L2_CAP_VBI_OUTPUT
	// CapSlicedVBICapture indicates the device supports sliced VBI capture.
	CapSlicedVBICapture uint32 = C.V4L2_CAP_SLICED_VBI_CAPTURE
	// CapSlicedVBIOutput indicates the device supports sliced VBI output.
	CapSlicedVBIOutput uint32 = C.V4L2_CAP_SLICED_VBI_OUTPUT
	// CapRDSCapture indicates the device supports RDS (Radio Data System) capture.
	CapRDSCapture uint32 = C.V4L2_CAP_RDS_CAPTURE
	// CapVideoOutputOverlay indicates the device supports video output overlay (OSD).
	CapVideoOutputOverlay uint32 = C.V4L2_CAP_VIDEO_OUTPUT_OVERLAY
	// CapHWFrequencySeek indicates the device supports hardware frequency seeking.
	CapHWFrequencySeek uint32 = C.V4L2_CAP_HW_FREQ_SEEK
	// CapRDSOutput indicates the device supports RDS output.
	CapRDSOutput uint32 = C.V4L2_CAP_RDS_OUTPUT

	// CapVideoCaptureMPlane indicates the device supports video capture (multi-planar API).
	CapVideoCaptureMPlane uint32 = C.V4L2_CAP_VIDEO_CAPTURE_MPLANE
	// CapVideoOutputMPlane indicates the device supports video output (multi-planar API).
	CapVideoOutputMPlane uint32 = C.V4L2_CAP_VIDEO_OUTPUT_MPLANE
	// CapVideoMem2MemMPlane indicates the device supports memory-to-memory video processing (multi-planar API).
	CapVideoMem2MemMPlane uint32 = C.V4L2_CAP_VIDEO_M2M_MPLANE
	// CapVideoMem2Mem indicates the device supports memory-to-memory video processing (single-planar API).
	CapVideoMem2Mem uint32 = C.V4L2_CAP_VIDEO_M2M

	// CapTuner indicates the device has a tuner.
	CapTuner uint32 = C.V4L2_CAP_TUNER
	// CapAudio indicates the device supports audio inputs or outputs.
	CapAudio uint32 = C.V4L2_CAP_AUDIO
	// CapRadio indicates the device is a radio receiver.
	CapRadio uint32 = C.V4L2_CAP_RADIO
	// CapModulator indicates the device has a modulator.
	CapModulator uint32 = C.V4L2_CAP_MODULATOR

	// CapSDRCapture indicates the device supports SDR (Software Defined Radio) capture.
	CapSDRCapture uint32 = C.V4L2_CAP_SDR_CAPTURE
	// CapExtendedPixFormat indicates the device supports extended pixel formats.
	CapExtendedPixFormat uint32 = C.V4L2_CAP_EXT_PIX_FORMAT
	// CapSDROutput indicates the device supports SDR output.
	CapSDROutput uint32 = C.V4L2_CAP_SDR_OUTPUT
	// CapMetadataCapture indicates the device supports metadata capture.
	CapMetadataCapture uint32 = C.V4L2_CAP_META_CAPTURE

	// CapReadWrite indicates the device supports read/write I/O operations.
	CapReadWrite uint32 = C.V4L2_CAP_READWRITE
	// CapAsyncIO indicates the device supports asynchronous I/O operations.
	CapAsyncIO uint32 = C.V4L2_CAP_ASYNCIO
	// CapStreaming indicates the device supports streaming I/O operations.
	CapStreaming uint32 = C.V4L2_CAP_STREAMING

	// CapMetadataOutput indicates the device supports metadata output.
	CapMetadataOutput uint32 = C.V4L2_CAP_META_OUTPUT
	// CapTouch indicates the device is a touch device.
	CapTouch uint32 = C.V4L2_CAP_TOUCH
	// CapIOMediaController indicates the device is part of an I/O Media Controller.
	CapIOMediaController uint32 = C.V4L2_CAP_IO_MC
	// CapDeviceCapabilities indicates that the driver fills the device_caps field with specific capabilities for the opened device node.
	// If not set, device_caps is a copy of the capabilities field.
	CapDeviceCapabilities uint32 = C.V4L2_CAP_DEVICE_CAPS
)

// CapabilityDesc provides a textual description for a V4L2 capability flag.
type CapabilityDesc struct {
	// Cap is the capability flag (e.g., CapVideoCapture).
	Cap uint32
	// Desc is the human-readable description of the capability.
	Desc string
}

// Capabilities is a predefined list of CapabilityDesc, mapping capability flags to their descriptions.
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

// Capability stores information about the V4L2 device's capabilities.
// It corresponds to the `v4l2_capability` struct in the Linux kernel.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L440
// and https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-querycap.html#c.V4L.v4l2_capability
type Capability struct {
	// Driver is a string identifying the driver module (e.g., "uvcvideo").
	Driver string

	// Card is a string identifying the device card (e.g., "Integrated Camera").
	Card string

	// BusInfo is a string identifying the bus the device is on (e.g., "usb-0000:00:14.0-1").
	BusInfo string

	// Version is the kernel version the driver was compiled for. Use GetVersionInfo() for a parsed representation.
	Version uint32

	// Capabilities is a bitmask of global capabilities of the physical device.
	// These are capabilities of the hardware, regardless of which device node was opened.
	// Use the Cap* constants to check for specific capabilities.
	Capabilities uint32

	// DeviceCapabilities is a bitmask of capabilities specific to the opened device node.
	// If CapDeviceCapabilities is set in the Capabilities field, this field contains device-specific
	// capabilities. Otherwise, it's a copy of the Capabilities field.
	DeviceCapabilities uint32
}

// GetCapability queries the V4L2 device for its capabilities.
// It takes the file descriptor of the opened device.
// It returns a Capability struct populated with the device's information and an error if the query fails.
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

// GetCapabilities returns the effective capabilities for the device.
// If CapDeviceCapabilities is set, it returns DeviceCapabilities; otherwise, it returns Capabilities.
func (c Capability) GetCapabilities() uint32 {
	if c.IsDeviceCapabilitiesProvided() {
		return c.DeviceCapabilities
	}
	return c.Capabilities
}

// IsVideoCaptureSupported checks if the device supports video capture (single-planar API).
func (c Capability) IsVideoCaptureSupported() bool {
	return c.GetCapabilities()&CapVideoCapture != 0
}

// IsVideoOutputSupported checks if the device supports video output (single-planar API).
func (c Capability) IsVideoOutputSupported() bool {
	return c.GetCapabilities()&CapVideoOutput != 0
}

// IsVideoOverlaySupported checks if the device supports video overlay.
func (c Capability) IsVideoOverlaySupported() bool {
	return c.GetCapabilities()&CapVideoOverlay != 0
}

// IsVideoOutputOverlaySupported checks if the device supports video output overlay.
func (c Capability) IsVideoOutputOverlaySupported() bool {
	return c.GetCapabilities()&CapVideoOutputOverlay != 0
}

// IsVideoCaptureMultiplanarSupported checks if the device supports video capture (multi-planar API).
func (c Capability) IsVideoCaptureMultiplanarSupported() bool {
	return c.GetCapabilities()&CapVideoCaptureMPlane != 0
}

// IsVideoOutputMultiplanarSupported checks if the device supports video output (multi-planar API).
func (c Capability) IsVideoOutputMultiplanarSupported() bool {
	return c.GetCapabilities()&CapVideoOutputMPlane != 0
}

// IsReadWriteSupported checks if the device supports read/write I/O.
func (c Capability) IsReadWriteSupported() bool {
	return c.GetCapabilities()&CapReadWrite != 0
}

// IsStreamingSupported checks if the device supports streaming I/O.
func (c Capability) IsStreamingSupported() bool {
	return c.GetCapabilities()&CapStreaming != 0
}

// IsDeviceCapabilitiesProvided checks if the driver fills the DeviceCapabilities field
// with capabilities specific to the opened device node.
// See notes on V4L2_CAP_DEVICE_CAPS:
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-querycap.html?highlight=v4l2_cap_device_caps
func (c Capability) IsDeviceCapabilitiesProvided() bool {
	return c.Capabilities&CapDeviceCapabilities != 0
}

// GetDriverCapDescriptions returns a slice of CapabilityDesc for all global capabilities
// reported in the Capabilities field.
func (c Capability) GetDriverCapDescriptions() []CapabilityDesc {
	var result []CapabilityDesc
	for _, capDesc := range Capabilities {
		if c.Capabilities&capDesc.Cap == capDesc.Cap {
			result = append(result, capDesc)
		}
	}
	return result
}

// GetDeviceCapDescriptions returns a slice of CapabilityDesc for all device-specific capabilities
// reported in the DeviceCapabilities field. This is relevant if IsDeviceCapabilitiesProvided() is true.
func (c Capability) GetDeviceCapDescriptions() []CapabilityDesc {
	var result []CapabilityDesc
	for _, capDesc := range Capabilities {
		if c.DeviceCapabilities&capDesc.Cap == capDesc.Cap {
			result = append(result, capDesc)
		}
	}
	return result
}

// GetVersionInfo parses the raw Version field into a VersionInfo struct.
func (c Capability) GetVersionInfo() VersionInfo {
	return VersionInfo{value: c.Version}
}

// String returns a string representation of the Capability struct,
// including driver, card, and bus information.
func (c Capability) String() string {
	return fmt.Sprintf("driver: %s; card: %s; bus info: %s", c.Driver, c.Card, c.BusInfo)
}
