package v4l2

//#include <linux/videodev2.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// Control (v4l2_control)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1725
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ctrl.html
type Control struct {
	ID    uint32
	Value uint32
}

// GetControl returns control value for specified ID
func GetControl(fd uintptr, id uint32) (Control, error) {
	var ctrl C.struct_v4l2_control
	ctrl.id = C.uint(id)

	if err := send(fd, C.VIDIOC_G_CTRL, uintptr(unsafe.Pointer(&ctrl))); err != nil {
		return Control{}, fmt.Errorf("get control: id %d: %w", id, err)
	}

	return Control{
		ID:    uint32(ctrl.id),
		Value: uint32(ctrl.value),
	}, nil
}

// SetControl applies control value for specified ID
func SetControl(fd uintptr, id, value uint32) error {
	var ctrl C.struct_v4l2_control
	ctrl.id = C.uint(id)
	ctrl.value = C.int(value)

	if err := send(fd, C.VIDIOC_G_CTRL, uintptr(unsafe.Pointer(&ctrl))); err != nil {
		return fmt.Errorf("set control: id %d: value: %d: %w", id, value, err)
	}

	return nil
}
