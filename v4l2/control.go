package v4l2

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../include/
#include <linux/videodev2.h>
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

// CtrlValue represents the value for a user control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1740
type CtrlValue = int32

// Control (v4l2_control)
//
// This type is used to query/set/get user-specific controls.
// For more information about user controls, see https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html.
//
// Also, see the followings:
//
//   - https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1725
//   - https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ctrl.html
type Control struct {
	fd      uintptr
	Type    CtrlType
	ID      CtrlID
	Value   CtrlValue
	Name    string
	Minimum int32
	Maximum int32
	Step    int32
	Default int32
	flags   uint32
}

type ControlMenuItem struct {
	ID    uint32
	Index uint32
	Value uint32
	Name  string
}

// IsMenu tests whether control Type == CtrlTypeMenu || Type == CtrlIntegerMenu
func (c Control) IsMenu() bool {
	return c.Type == CtrlTypeMenu || c.Type == CtrlTypeIntegerMenu
}

// GetMenuItems returns control menu items if the associated control is a menu.
func (c Control) GetMenuItems() (result []ControlMenuItem, err error) {
	if !c.IsMenu() {
		return result, fmt.Errorf("control is not a menu type")
	}

	for idx := c.Minimum; idx <= c.Maximum; idx++ {
		var qryMenu C.struct_v4l2_querymenu
		qryMenu.id = C.uint(c.ID)
		qryMenu.index = C.uint(idx)
		if err = send(c.fd, C.VIDIOC_QUERYMENU, uintptr(unsafe.Pointer(&qryMenu))); err != nil {
			continue
		}
		result = append(result, makeCtrlMenu(c.Type, qryMenu))
	}

	return result, nil
}

// GetControlValue retrieves the value for a user control with the specified id.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1740
func GetControlValue(fd uintptr, id CtrlID) (CtrlValue, error) {
	var ctrl C.struct_v4l2_control
	ctrl.id = C.uint(id)

	if err := send(fd, C.VIDIOC_G_CTRL, uintptr(unsafe.Pointer(&ctrl))); err != nil {
		return 0, fmt.Errorf("get control value: VIDIOC_G_CTRL: id %d: %w", id, err)
	}

	return CtrlValue(ctrl.value), nil
}

// SetControlValue sets the value for a user control with the specified id.
// This function applies range check based on the values supported by the  control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1740
func SetControlValue(fd uintptr, id CtrlID, val CtrlValue) error {
	ctrlInfo, err := QueryControlInfo(fd, id)
	if err != nil {
		return fmt.Errorf("set control value: id %s: %w", id, err)
	}
	if val < ctrlInfo.Minimum || val > ctrlInfo.Maximum {
		return fmt.Errorf("set control value: out-of-range failure: val %d: expected ctrl.Min %d, ctrl.Max %d", val, ctrlInfo.Minimum, ctrlInfo.Maximum)
	}

	var ctrl C.struct_v4l2_control
	ctrl.id = C.uint(id)
	ctrl.value = C.int(val)

	if err := send(fd, C.VIDIOC_S_CTRL, uintptr(unsafe.Pointer(&ctrl))); err != nil {
		return fmt.Errorf("set control value: id %d: %w", id, err)
	}

	return nil
}

// QueryControlInfo queries information about the specified control without the current value.
func QueryControlInfo(fd uintptr, id CtrlID) (Control, error) {
	// query control information
	var qryCtrl C.struct_v4l2_queryctrl
	qryCtrl.id = C.uint(id)

	if err := send(fd, C.VIDIOC_QUERYCTRL, uintptr(unsafe.Pointer(&qryCtrl))); err != nil {
		return Control{}, fmt.Errorf("query control info: VIDIOC_QUERYCTRL: id %d: %w", id, err)
	}
	control := makeControl(qryCtrl)
	control.fd = fd
	return control, nil
}

// GetControl retrieves the value and information for the user control with the specified id.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html
func GetControl(fd uintptr, id CtrlID) (Control, error) {
	control, err := QueryControlInfo(fd, id)
	if err != nil {
		return Control{}, fmt.Errorf("get control: %w", err)
	}

	// retrieve control value
	ctrlValue, err := GetControlValue(fd, uint32(id))
	if err != nil {
		return Control{}, fmt.Errorf("get control: %w", id, err)
	}

	control.Value = ctrlValue
	control.fd = fd
	return control, nil
}

// QueryAllControls loop through all available user controls and returns information for
// all controls without their current values (use GetControlValue to get current values).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html
func QueryAllControls(fd uintptr) (result []Control, err error) {
	cid := uint32(C.V4L2_CTRL_FLAG_NEXT_CTRL)
	for {
		control, err := QueryControlInfo(fd, cid)
		if err != nil {
			if errors.Is(err, ErrorBadArgument) && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("query all controls: %w", err)
		}
		result = append(result, control)
		// setup next id
		cid = control.ID | uint32(C.V4L2_CTRL_FLAG_NEXT_CTRL)
	}
	return result, nil
}

func makeControl(qryCtrl C.struct_v4l2_queryctrl) Control {
	return Control{
		Type:    CtrlType(qryCtrl._type),
		ID:      uint32(qryCtrl.id),
		Name:    C.GoString((*C.char)(unsafe.Pointer(&qryCtrl.name[0]))),
		Maximum: int32(qryCtrl.maximum),
		Minimum: int32(qryCtrl.minimum),
		Step:    int32(qryCtrl.step),
		Default: int32(qryCtrl.default_value),
		flags:   uint32(qryCtrl.flags),
	}
}

func makeCtrlMenu(cType CtrlType, qryMenu C.struct_v4l2_querymenu) ControlMenuItem {
	item := ControlMenuItem{
		ID:    uint32(qryMenu.id),
		Index: uint32(qryMenu.index),
	}
	if cType == CtrlTypeIntegerMenu {
		val := binary.LittleEndian.Uint64((*[8]byte)(unsafe.Pointer(&qryMenu.anon0[0]))[:])
		item.Name = strconv.FormatInt(int64(val), 10)
	} else {
		item.Name = C.GoString((*C.char)(unsafe.Pointer(&qryMenu.anon0[0])))
	}

	return item
}
