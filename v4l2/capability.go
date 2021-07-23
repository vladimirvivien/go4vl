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
	CapVideoOutputOverlay = 0x00000200 // V4L2_CAP_VIDEO_OUTPUT_OVERLAY
	CapReadWrite          = 0x01000000 // V4L2_CAP_READWRITE
	CapAsyncIO            = 0x02000000 // V4L2_CAP_ASYNCIO
	CapStreaming          = 0x04000000 // V4L2_CAP_STREAMING
)

// v4l2Capabiolity type for device (see v4l2_capability
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
// Use methods on this type to access capabilities.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-querycap.html#c.V4L.v4l2_capability
type Capability struct {
	v4l2Cap v4l2Capability
}

// GetCapability retrieves capability info for device
func GetCapability(fd uintptr) (Capability, error) {
	v4l2Cap := v4l2Capability{}
	if err := Send(fd, vidiocQueryCap, uintptr(unsafe.Pointer(&v4l2Cap))); err != nil {
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

func (c Capability) IsVideoCaptureSupported() bool {
	return (c.v4l2Cap.capabilities & CapVideoCapture) != 0
}

func (c Capability) IsVideoOutputSupported() bool {
	return (c.v4l2Cap.capabilities & CapVideoOutput) != 0
}

func (c Capability) IsReadWriteSupported() bool {
	return (c.v4l2Cap.capabilities & CapReadWrite) != 0
}

func (c Capability) IsStreamingSupported() bool {
	return (c.v4l2Cap.capabilities & CapStreaming) != 0
}

func (c Capability) DriverName() string {
	return GoString(c.v4l2Cap.driver[:])
}

func (c Capability) CardName() string {
	return GoString(c.v4l2Cap.card[:])
}

func (c Capability) BusInfo() string {
	return GoString(c.v4l2Cap.busInfo[:])
}

func (c Capability) GetVersion() uint32 {
	return c.v4l2Cap.version
}

func (c Capability) String() string {
	return fmt.Sprintf("driver: %s; card: %s; bus info: %s", c.DriverName(), c.CardName(), c.BusInfo())
}
