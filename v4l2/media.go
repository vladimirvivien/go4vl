package v4l2

// #include <linux/media.h>
import "C"
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
	var mdi C.struct_media_device_info
	if err := send(fd, C.MEDIA_IOC_DEVICE_INFO, uintptr(unsafe.Pointer(&mdi))); err != nil {
		return MediaDeviceInfo{}, fmt.Errorf("media device info: %w", err)
	}
	return MediaDeviceInfo{
		Driver:           C.GoString((*C.char)(&mdi.driver[0])),
		Model:            C.GoString((*C.char)(&mdi.model[0])),
		Serial:           C.GoString((*C.char)(&mdi.serial[0])),
		BusInfo:          C.GoString((*C.char)(&mdi.bus_info[0])),
		MediaVersion:     VersionInfo{value: uint32(mdi.media_version)},
		HardwareRevision: uint32(mdi.hw_revision),
		DriverVersion:    VersionInfo{value: uint32(mdi.driver_version)},
	}, nil
}
