package v4l2

import (
	"fmt"
	"unsafe"
)

// MediaDeviceInfo (media_device_info)
// See https://www.kernel.org/doc/html/latest/userspace-api/media/mediactl/media-ioc-device-info.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/media.h#L29
type MediaDeviceInfo struct {
	Driver           string
	Model            string
	Serial           string
	BusInfo          string
	MediaVersion     VersionInfo
	HardwareRevision uint32
	DriverVersion    VersionInfo
}

// GetMediaDeviceInfo retrieves media information for specified device, if supported.
func GetMediaDeviceInfo(fd uintptr) (MediaDeviceInfo, error) {
	var mdi cMediaDeviceInfo
	if err := send(fd, cMEDIA_IOC_DEVICE_INFO, uintptr(unsafe.Pointer(&mdi))); err != nil {
		return MediaDeviceInfo{}, fmt.Errorf("media device info: %w", err)
	}
	return MediaDeviceInfo{
		Driver:           cString(mdi.Driver[:]),
		Model:            cString(mdi.Model[:]),
		Serial:           cString(mdi.Serial[:]),
		BusInfo:          cString(mdi.Bus_info[:]),
		MediaVersion:     VersionInfo{value: mdi.Media_version},
		HardwareRevision: mdi.Hw_revision,
		DriverVersion:    VersionInfo{value: mdi.Driver_version},
	}, nil
}
