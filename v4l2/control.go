package v4l2

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

// CtrlFlag represents flags for V4L2 controls
type CtrlFlag = uint32

// Control flags
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
const (
	CtrlFlagDisabled        CtrlFlag = cV4L2_CTRL_FLAG_DISABLED
	CtrlFlagGrabbed         CtrlFlag = cV4L2_CTRL_FLAG_GRABBED
	CtrlFlagReadOnly        CtrlFlag = cV4L2_CTRL_FLAG_READ_ONLY
	CtrlFlagUpdate          CtrlFlag = cV4L2_CTRL_FLAG_UPDATE
	CtrlFlagInactive        CtrlFlag = cV4L2_CTRL_FLAG_INACTIVE
	CtrlFlagSlider          CtrlFlag = cV4L2_CTRL_FLAG_SLIDER
	CtrlFlagWriteOnly       CtrlFlag = cV4L2_CTRL_FLAG_WRITE_ONLY
	CtrlFlagVolatile        CtrlFlag = cV4L2_CTRL_FLAG_VOLATILE
	CtrlFlagHasPayload      CtrlFlag = cV4L2_CTRL_FLAG_HAS_PAYLOAD
	CtrlFlagExecuteOnWrite  CtrlFlag = cV4L2_CTRL_FLAG_EXECUTE_ON_WRITE
	CtrlFlagModifyLayout    CtrlFlag = cV4L2_CTRL_FLAG_MODIFY_LAYOUT
	CtrlFlagDynamicArray    CtrlFlag = cV4L2_CTRL_FLAG_DYNAMIC_ARRAY
	CtrlFlagNextCtrl        CtrlFlag = cV4L2_CTRL_FLAG_NEXT_CTRL
	CtrlFlagNextCompound    CtrlFlag = cV4L2_CTRL_FLAG_NEXT_COMPOUND
)

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
	Flags   uint32
}

// ControlMenuItem represents a single option in a menu-type V4L2 control.
// Menu controls allow selection from a discrete set of named options,
// such as white balance modes, scene modes, or power line frequency settings.
type ControlMenuItem struct {
	// ID is the control identifier this menu item belongs to
	ID uint32
	// Index is the menu item's position in the menu (starting from 0)
	Index uint32
	// Value is the numeric value associated with this menu option
	Value uint32
	// Name is the human-readable label for this menu option
	Name string
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
		var qryMenu cV4L2QueryMenu
		qryMenu.Id = c.ID
		qryMenu.Index = uint32(idx)
		if err = send(c.fd, cVIDIOC_QUERYMENU, unsafe.Pointer(&qryMenu)); err != nil {
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
	var ctrl cV4L2Control
	ctrl.Id = id

	if err := send(fd, cVIDIOC_G_CTRL, unsafe.Pointer(&ctrl)); err != nil {
		return 0, fmt.Errorf("get control value: VIDIOC_G_CTRL: id %d: %w", id, err)
	}

	return CtrlValue(ctrl.Value), nil
}

// SetControlValue sets the value for a user control with the specified id.
// This function applies range check based on the values supported by the  control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1740
func SetControlValue(fd uintptr, id CtrlID, val CtrlValue) error {
	ctrlInfo, err := QueryControlInfo(fd, id)
	if err != nil {
		return fmt.Errorf("set control value: id %d: %w", id, err)
	}
	if val < ctrlInfo.Minimum || val > ctrlInfo.Maximum {
		return fmt.Errorf("set control value: out-of-range failure: val %d: expected ctrl.Min %d, ctrl.Max %d", val, ctrlInfo.Minimum, ctrlInfo.Maximum)
	}

	var ctrl cV4L2Control
	ctrl.Id = id
	ctrl.Value = val

	if err := send(fd, cVIDIOC_S_CTRL, unsafe.Pointer(&ctrl)); err != nil {
		return fmt.Errorf("set control value: id %d: %w", id, err)
	}

	return nil
}

// QueryControlInfo queries information about the specified control without the current value.
func QueryControlInfo(fd uintptr, id CtrlID) (Control, error) {
	// query control information
	var qryCtrl cV4L2QueryCtrl
	qryCtrl.Id = id

	if err := send(fd, cVIDIOC_QUERYCTRL, unsafe.Pointer(&qryCtrl)); err != nil {
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
		return Control{}, fmt.Errorf("get control: id %d: %w", id, err)
	}

	control.Value = ctrlValue
	control.fd = fd
	return control, nil
}

// QueryAllControls loop through all available user controls and returns information for
// all controls without their current values (use GetControlValue to get current values).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html
func QueryAllControls(fd uintptr) (result []Control, err error) {
	cid := uint32(cV4L2_CTRL_FLAG_NEXT_CTRL)
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
		cid = control.ID | cV4L2_CTRL_FLAG_NEXT_CTRL
	}
	return result, nil
}

func makeControl(qryCtrl cV4L2QueryCtrl) Control {
	return Control{
		Type:    CtrlType(qryCtrl.Type),
		ID:      qryCtrl.Id,
		Name:    cString(qryCtrl.Name[:]),
		Maximum: qryCtrl.Maximum,
		Minimum: qryCtrl.Minimum,
		Step:    qryCtrl.Step,
		Default: qryCtrl.Value,
		Flags:   qryCtrl.Flags,
	}
}

func makeCtrlMenu(cType CtrlType, qryMenu cV4L2QueryMenu) ControlMenuItem {
	item := ControlMenuItem{
		ID:    qryMenu.Id,
		Index: qryMenu.Index,
	}
	if cType == CtrlTypeIntegerMenu {
		val := binary.LittleEndian.Uint64(qryMenu.Name[:8])
		item.Name = strconv.FormatInt(int64(val), 10)
	} else {
		item.Name = cString(qryMenu.Name[:])
	}

	return item
}
