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

// Capability constants define the various features and functionalities that a V4L2 device can support.
// These flags are used in the Capabilities and DeviceCapabilities fields of the Capability struct
// to indicate what operations the device or driver supports.
//
// Reference: https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L451
const (
	// Video capture/output capabilities

	// CapVideoCapture indicates the device supports video capture via the single-planar API.
	// This is the most common capability for webcams and capture cards.
	CapVideoCapture uint32 = C.V4L2_CAP_VIDEO_CAPTURE

	// CapVideoOutput indicates the device supports video output via the single-planar API.
	// Used for video output devices like displays or video encoders.
	CapVideoOutput uint32 = C.V4L2_CAP_VIDEO_OUTPUT

	// CapVideoOverlay indicates the device supports video overlay onto a display.
	// Allows direct overlay of video onto the screen without going through system memory.
	CapVideoOverlay uint32 = C.V4L2_CAP_VIDEO_OVERLAY

	// VBI (Vertical Blanking Interval) capabilities

	// CapVBICapture indicates support for raw VBI data capture.
	// VBI is used for transmitting data like closed captions in analog TV signals.
	CapVBICapture uint32 = C.V4L2_CAP_VBI_CAPTURE

	// CapVBIOutput indicates support for raw VBI data output.
	CapVBIOutput uint32 = C.V4L2_CAP_VBI_OUTPUT

	// CapSlicedVBICapture indicates support for sliced (parsed) VBI data capture.
	// Provides decoded VBI services like teletext or closed captions.
	CapSlicedVBICapture uint32 = C.V4L2_CAP_SLICED_VBI_CAPTURE

	// CapSlicedVBIOutput indicates support for sliced VBI data output.
	CapSlicedVBIOutput uint32 = C.V4L2_CAP_SLICED_VBI_OUTPUT

	// Radio and RDS capabilities

	// CapRDSCapture indicates support for RDS (Radio Data System) data capture.
	// RDS provides digital information in FM radio broadcasts.
	CapRDSCapture uint32 = C.V4L2_CAP_RDS_CAPTURE

	// CapVideoOutputOverlay indicates support for video output overlay (OSD).
	// Allows overlaying graphics on video output.
	CapVideoOutputOverlay uint32 = C.V4L2_CAP_VIDEO_OUTPUT_OVERLAY

	// CapHWFrequencySeek indicates hardware supports automatic frequency seeking.
	// Used in radio tuners for scanning to the next station.
	CapHWFrequencySeek uint32 = C.V4L2_CAP_HW_FREQ_SEEK

	// CapRDSOutput indicates support for RDS data output.
	CapRDSOutput uint32 = C.V4L2_CAP_RDS_OUTPUT

	// Multi-planar API capabilities

	// CapVideoCaptureMPlane indicates video capture support via the multi-planar API.
	// Used for formats where color components are stored in separate memory planes.
	CapVideoCaptureMPlane uint32 = C.V4L2_CAP_VIDEO_CAPTURE_MPLANE

	// CapVideoOutputMPlane indicates video output support via the multi-planar API.
	CapVideoOutputMPlane uint32 = C.V4L2_CAP_VIDEO_OUTPUT_MPLANE

	// CapVideoMem2MemMPlane indicates memory-to-memory device support via multi-planar API.
	// Used for hardware codecs and format converters.
	CapVideoMem2MemMPlane uint32 = C.V4L2_CAP_VIDEO_M2M_MPLANE

	// CapVideoMem2Mem indicates memory-to-memory device support via single-planar API.
	CapVideoMem2Mem uint32 = C.V4L2_CAP_VIDEO_M2M

	// Tuner and audio capabilities

	// CapTuner indicates the device has a TV or radio tuner.
	CapTuner uint32 = C.V4L2_CAP_TUNER

	// CapAudio indicates the device has audio inputs or outputs.
	CapAudio uint32 = C.V4L2_CAP_AUDIO

	// CapRadio indicates the device is a radio receiver.
	CapRadio uint32 = C.V4L2_CAP_RADIO

	// CapModulator indicates the device has a radio frequency modulator.
	// Used for FM transmitters.
	CapModulator uint32 = C.V4L2_CAP_MODULATOR

	// SDR and metadata capabilities

	// CapSDRCapture indicates support for Software Defined Radio capture.
	CapSDRCapture uint32 = C.V4L2_CAP_SDR_CAPTURE

	// CapExtendedPixFormat indicates the device supports extended pixel format fields.
	// Enables use of additional format flags and modifiers.
	CapExtendedPixFormat uint32 = C.V4L2_CAP_EXT_PIX_FORMAT

	// CapSDROutput indicates support for Software Defined Radio output.
	CapSDROutput uint32 = C.V4L2_CAP_SDR_OUTPUT

	// CapMetadataCapture indicates the device can capture metadata.
	// Used for sensors that provide metadata alongside video frames.
	CapMetadataCapture uint32 = C.V4L2_CAP_META_CAPTURE

	// I/O method capabilities

	// CapReadWrite indicates support for read() and write() I/O methods.
	// Simple but less efficient than streaming I/O.
	CapReadWrite uint32 = C.V4L2_CAP_READWRITE

	// CapAsyncIO indicates support for asynchronous I/O.
	CapAsyncIO uint32 = C.V4L2_CAP_ASYNCIO

	// CapStreaming indicates support for streaming I/O using memory mapping or user pointers.
	// This is the most efficient I/O method for continuous video streaming.
	CapStreaming uint32 = C.V4L2_CAP_STREAMING

	// Additional capabilities

	// CapMetadataOutput indicates the device can output metadata.
	CapMetadataOutput uint32 = C.V4L2_CAP_META_OUTPUT

	// CapTouch indicates the device is a touch device.
	CapTouch uint32 = C.V4L2_CAP_TOUCH

	// CapIOMediaController indicates the device supports the Media Controller API.
	// Used for complex devices with configurable media pipelines.
	CapIOMediaController uint32 = C.V4L2_CAP_IO_MC

	// CapDeviceCapabilities indicates the device provides device-specific capabilities.
	// When set, DeviceCapabilities field should be used instead of Capabilities.
	CapDeviceCapabilities uint32 = C.V4L2_CAP_DEVICE_CAPS
)

// CapabilityDesc provides a human-readable description for a capability flag.
// Used for displaying device capabilities in a user-friendly format.
type CapabilityDesc struct {
	// Cap is the capability flag constant
	Cap uint32
	// Desc is a human-readable description of the capability
	Desc string
}

// Capabilities provides human-readable descriptions for all V4L2 capability flags.
// This slice can be used to iterate through capabilities and generate reports
// or user interfaces showing device features.
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

// Capability represents the capabilities and identification information of a V4L2 device.
// This struct corresponds to the v4l2_capability structure in the V4L2 API.
//
// The Capability struct provides two sets of capabilities:
//   - Capabilities: The physical device capabilities (all functions the hardware supports)
//   - DeviceCapabilities: The opened device node capabilities (what this specific node can do)
//
// For modern drivers that set CapDeviceCapabilities, use DeviceCapabilities.
// For older drivers, use Capabilities.
//
// References:
//   - https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L440
//   - https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-querycap.html#c.V4L.v4l2_capability
type Capability struct {
	// Driver is the name of the driver module (e.g., "uvcvideo" for USB cameras)
	Driver string

	// Card is a human-readable name of the device (e.g., "HD Webcam C920")
	Card string

	// BusInfo provides information about the device connection (e.g., "usb-0000:00:14.0-1")
	BusInfo string

	// Version encodes the kernel driver version as: (major << 16) | (minor << 8) | patch
	Version uint32

	// Capabilities is a bitmask of all capabilities supported by the physical device.
	// For devices with multiple nodes (e.g., separate capture and metadata nodes),
	// this shows the combined capabilities of all nodes.
	Capabilities uint32

	// DeviceCapabilities is a bitmask of capabilities for this specific opened device node.
	// Only valid when CapDeviceCapabilities is set in Capabilities.
	// This field is preferred over Capabilities for determining what the opened node can do.
	DeviceCapabilities uint32
}

// GetCapability queries the V4L2 device for its capabilities and identification information.
//
// Parameters:
//   - fd: File descriptor of an opened V4L2 device
//
// Returns:
//   - Capability: Device capabilities and identification
//   - error: An error if the ioctl fails
//
// This function issues the VIDIOC_QUERYCAP ioctl to retrieve device information.
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

// GetCapabilities returns the appropriate capability flags for the device.
// If the device provides device-specific capabilities (modern drivers),
// it returns DeviceCapabilities. Otherwise, it returns Capabilities.
//
// Use this method to get the correct capability flags regardless of driver version.
func (c Capability) GetCapabilities() uint32 {
	if c.IsDeviceCapabilitiesProvided() {
		return c.DeviceCapabilities
	}
	return c.Capabilities
}

// IsVideoCaptureSupported checks if the device supports video capture via single-planar API.
// Returns true if the device can capture video from a camera or other video source.
func (c Capability) IsVideoCaptureSupported() bool {
	return c.Capabilities&CapVideoCapture != 0
}

// IsVideoOutputSupported checks if the device supports video output via single-planar API.
// Returns true if the device can output video to a display or encoder.
func (c Capability) IsVideoOutputSupported() bool {
	return c.Capabilities&CapVideoOutput != 0
}

// IsVideoOverlaySupported checks if the device supports video overlay.
// Returns true if the device can overlay video directly onto a display.
func (c Capability) IsVideoOverlaySupported() bool {
	return c.Capabilities&CapVideoOverlay != 0
}

// IsVideoOutputOverlaySupported checks if the device supports video output overlay (OSD).
// Returns true if the device can overlay graphics on video output.
func (c Capability) IsVideoOutputOverlaySupported() bool {
	return c.Capabilities&CapVideoOutputOverlay != 0
}

// IsVideoCaptureMultiplanarSupported checks if the device supports video capture via multi-planar API.
// Returns true for devices that handle formats with separate memory planes for color components.
func (c Capability) IsVideoCaptureMultiplanarSupported() bool {
	return c.Capabilities&CapVideoCaptureMPlane != 0
}

// IsVideoOutputMultiplanerSupported checks if the device supports video output via multi-planar API.
// Returns true for devices that handle formats with separate memory planes for color components.
func (c Capability) IsVideoOutputMultiplanerSupported() bool {
	return c.Capabilities&CapVideoOutputMPlane != 0
}

// IsReadWriteSupported checks if the device supports read() and write() I/O methods.
// Returns true if simple read/write operations are supported (less efficient than streaming).
func (c Capability) IsReadWriteSupported() bool {
	return c.Capabilities&CapReadWrite != 0
}

// IsStreamingSupported checks if the device supports streaming I/O (memory mapping or user pointers).
// Returns true if the efficient streaming I/O method is supported.
// This is the preferred I/O method for continuous video capture or output.
func (c Capability) IsStreamingSupported() bool {
	return c.Capabilities&CapStreaming != 0
}

// IsDeviceCapabilitiesProvided checks if the device provides device-specific capabilities.
// When true, the DeviceCapabilities field contains capabilities for the opened device node,
// which should be used instead of the Capabilities field.
//
// Modern V4L2 drivers set this flag to distinguish between:
//   - Physical device capabilities (all hardware features)
//   - Device node capabilities (what this specific node can do)
//
// Reference: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-querycap.html?highlight=v4l2_cap_device_caps
func (c Capability) IsDeviceCapabilitiesProvided() bool {
	return c.Capabilities&CapDeviceCapabilities != 0
}

// GetDriverCapDescriptions returns human-readable descriptions of all capabilities
// supported by the physical device (from the Capabilities field).
//
// This includes all features the hardware supports, which may span multiple device nodes.
func (c Capability) GetDriverCapDescriptions() []CapabilityDesc {
	var result []CapabilityDesc
	for _, cap := range Capabilities {
		if c.Capabilities&cap.Cap == cap.Cap {
			result = append(result, cap)
		}
	}
	return result
}

// GetDeviceCapDescriptions returns human-readable descriptions of capabilities
// for the specific opened device node (from the DeviceCapabilities field).
//
// Only valid when IsDeviceCapabilitiesProvided() returns true.
// This shows what operations can be performed on the opened device node.
func (c Capability) GetDeviceCapDescriptions() []CapabilityDesc {
	var result []CapabilityDesc
	for _, cap := range Capabilities {
		if c.DeviceCapabilities&cap.Cap == cap.Cap {
			result = append(result, cap)
		}
	}
	return result
}

// GetVersionInfo returns the driver version information decoded from the Version field.
// The version is encoded as: (major << 16) | (minor << 8) | patch
func (c Capability) GetVersionInfo() VersionInfo {
	return VersionInfo{value: c.Version}
}

// String returns a formatted string with device identification information.
// Includes the driver name, device name, and bus information.
//
// Example output: "driver: uvcvideo; card: HD Webcam C920; bus info: usb-0000:00:14.0-1"
func (c Capability) String() string {
	return fmt.Sprintf("driver: %s; card: %s; bus info: %s", c.Driver, c.Card, c.BusInfo)
}
