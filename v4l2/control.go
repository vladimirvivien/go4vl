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

// CtrlValue represents the value of a V4L2 control. It is an alias for int32.
// This type is used when getting or setting control values.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1740
type CtrlValue = int32

// Control represents a V4L2 control (corresponds to `v4l2_control` and `v4l2_queryctrl`).
// It holds information about a specific control, such as its ID, type, name,
// minimum/maximum values, step, and default value. The current value of the control
// is also stored in the Value field when retrieved using GetControl.
//
// For more information about user controls, see:
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html
//
// Kernel struct references:
//   - `v4l2_queryctrl`: https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1617
//   - `v4l2_control`: https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1725
// API call reference:
//   - `VIDIOC_G_CTRL`: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ctrl.html
//   - `VIDIOC_S_CTRL`: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-s-ctrl.html
//   - `VIDIOC_QUERYCTRL`: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-queryctrl.html
type Control struct {
	fd uintptr // unexported file descriptor for internal use (e.g. GetMenuItems)
	// Type is the data type of the control (e.g., integer, boolean, menu). See CtrlType constants.
	Type CtrlType
	// ID is the unique identifier of the control. See CtrlID constants.
	ID CtrlID
	// Value is the current value of the control. This is populated by GetControl.
	Value CtrlValue
	// Name is a human-readable name for the control (e.g., "Brightness").
	Name string
	// Minimum is the minimum value the control can take.
	Minimum int32
	// Maximum is the maximum value the control can take.
	Maximum int32
	// Step is the smallest change that can be made to the control's value.
	Step int32
	// Default is the default value of the control.
	Default int32
	flags   uint32 // unexported, stores V4L2 control flags
}

// ControlMenuItem represents a single item in a menu or integer menu control.
// It corresponds to the `v4l2_querymenu` struct in the Linux kernel.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1699
type ControlMenuItem struct {
	// ID is the control ID to which this menu item belongs.
	ID uint32
	// Index is the numerical index of this menu item.
	Index uint32
	// Value is the value associated with this menu item (for integer menus).
	// For regular menus, this field is not typically used directly by applications;
	// the Name field provides the string representation.
	Value uint32
	// Name is the human-readable string representation of the menu item.
	// For integer menus, this will be the string representation of the integer Value.
	Name string
}

// IsMenu checks if the control type is either CtrlTypeMenu or CtrlTypeIntegerMenu.
// It returns true if the control is a menu type, false otherwise.
func (c Control) IsMenu() bool {
	return c.Type == CtrlTypeMenu || c.Type == CtrlTypeIntegerMenu
}

// GetMenuItems retrieves the list of items for a menu or integer menu control.
// It iterates from the control's Minimum to Maximum value, querying each menu item.
// Returns a slice of ControlMenuItem and an error if the control is not a menu type or if querying fails.
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

// GetControlValue retrieves the current value of a specific control identified by its CtrlID.
// It takes the file descriptor of the V4L2 device and the control ID.
// It returns the control's current value (CtrlValue) and an error if the VIDIOC_G_CTRL ioctl call fails.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ctrl.html
func GetControlValue(fd uintptr, id CtrlID) (CtrlValue, error) {
	var ctrl C.struct_v4l2_control
	ctrl.id = C.uint(id)

	if err := send(fd, C.VIDIOC_G_CTRL, uintptr(unsafe.Pointer(&ctrl))); err != nil {
		return 0, fmt.Errorf("get control value: VIDIOC_G_CTRL: id %d: %w", id, err)
	}

	return CtrlValue(ctrl.value), nil
}

// SetControlValue sets the value of a specific control identified by its CtrlID.
// It takes the file descriptor of the V4L2 device, the control ID, and the desired value (CtrlValue).
// This function first queries the control's information to validate the new value against its min/max range.
// It returns an error if querying control info fails, if the value is out of range, or if the VIDIOC_S_CTRL ioctl call fails.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-s-ctrl.html
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

// QueryControlInfo retrieves detailed information about a specific control identified by its CtrlID,
// but *without* its current value. To get the current value, use GetControl or GetControlValue.
// It takes the file descriptor of the V4L2 device and the control ID.
// It returns a Control struct populated with the control's attributes (name, type, min, max, step, default)
// and an error if the VIDIOC_QUERYCTRL ioctl call fails.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-queryctrl.html
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

// GetControl retrieves detailed information about a specific control, *including* its current value.
// It combines the functionality of QueryControlInfo and GetControlValue.
// It takes the file descriptor of the V4L2 device and the control ID.
// It returns a Control struct populated with all attributes and the current value,
// and an error if either querying control info or getting its value fails.
func GetControl(fd uintptr, id CtrlID) (Control, error) {
	control, err := QueryControlInfo(fd, id)
	if err != nil {
		return Control{}, fmt.Errorf("get control: %w", err)
	}

	// retrieve control value
	ctrlValue, err := GetControlValue(fd, id) // Use CtrlID directly
	if err != nil {
		return Control{}, fmt.Errorf("get control: query value for id %d: %w", id, err)
	}

	control.Value = ctrlValue
	control.fd = fd
	return control, nil
}

// QueryAllControls iterates through all available user controls on the device and retrieves
// their information (name, type, min, max, etc.), but *without* their current values.
// To get current values, use GetControlValue or GetControl for each specific control.
// It takes the file descriptor of the V4L2 device.
// It returns a slice of Control structs and an error if querying fails.
// The iteration stops when the driver returns an error indicating no more controls are available (typically ErrorBadArgument).
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-queryctrl.html#querying-control-details
func QueryAllControls(fd uintptr) (result []Control, err error) {
	cid := uint32(C.V4L2_CTRL_FLAG_NEXT_CTRL) // Start with the first control flag
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
