package v4l2

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../include/
#include <linux/videodev2.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// GetExtControlValue retrieves the value of a single extended control.
// It takes the file descriptor of the V4L2 device and the CtrlID of the control.
// It returns the control's current value (CtrlValue) and an error if the VIDIOC_G_EXT_CTRLS ioctl call fails.
// Note that VIDIOC_G_EXT_CTRLS is designed for multiple controls, but this function uses it for a single control.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/extended-controls.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1745 (struct v4l2_ext_control)
// and https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1774 (struct v4l2_ext_controls)
func GetExtControlValue(fd uintptr, ctrlID CtrlID) (CtrlValue, error) {
	var v4l2Ctrl C.struct_v4l2_ext_control
	v4l2Ctrl.id = C.uint(ctrlID)
	v4l2Ctrl.size = 0
	if err := send(fd, C.VIDIOC_G_EXT_CTRLS, uintptr(unsafe.Pointer(&v4l2Ctrl))); err != nil {
		return 0, fmt.Errorf("get ext controls: %w", err)
	}
	return *(*CtrlValue)(unsafe.Pointer(&v4l2Ctrl.anon0[0])), nil
}

// SetExtControlValue sets the value of a single extended control.
// It takes the file descriptor, the control's CtrlID, and the desired CtrlValue.
// This function first queries the control's information to validate the new value against its min/max range.
// It uses the VIDIOC_S_CTRL ioctl (not VIDIOC_S_EXT_CTRLS) for setting the single control value,
// which might be unexpected given the "ExtControl" naming.
//
// Returns an error if querying control info fails, the value is out of range, or the set operation fails.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/extended-controls.html
func SetExtControlValue(fd uintptr, id CtrlID, val CtrlValue) error {
	ctrlInfo, err := QueryExtControlInfo(fd, id)
	if err != nil {
		return fmt.Errorf("set ext control value: id %s: %w", id, err)
	}
	if val < ctrlInfo.Minimum || val > ctrlInfo.Maximum {
		return fmt.Errorf("set ext control value: out-of-range failure: val %d: expected ctrl.Min %d, ctrl.Max %d", val, ctrlInfo.Minimum, ctrlInfo.Maximum)
	}

	var v4l2Ctrl C.struct_v4l2_ext_control
	v4l2Ctrl.id = C.uint(id)
	*(*C.int)(unsafe.Pointer(&v4l2Ctrl.anon0[0])) = *(*C.int)(unsafe.Pointer(&val))

	if err := send(fd, C.VIDIOC_S_CTRL, uintptr(unsafe.Pointer(&v4l2Ctrl))); err != nil {
		return fmt.Errorf("set ext control value: id %d: %w", id, err)
	}

	return nil
}

// SetExtControlValues sets the values of one or more extended controls simultaneously.
// It takes the file descriptor, a CtrlClass to specify which class of controls is being set (or 0 for any class),
// and a slice of Control structs containing the IDs and desired values.
// This function uses the VIDIOC_S_EXT_CTRLS ioctl call.
//
// Returns an error if the VIDIOC_S_EXT_CTRLS ioctl call fails.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/extended-controls.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1774 (struct v4l2_ext_controls)
func SetExtControlValues(fd uintptr, whichCtrl CtrlClass, ctrls []Control) error {
	numCtrl := len(ctrls)
	if numCtrl == 0 {
		return nil // Nothing to set
	}

	// Allocate C memory for the array of v4l2_ext_control structs
	cExtCtrlArray := C.make_v4l2_ext_control_array(C.int(numCtrl))
	if cExtCtrlArray == nil {
		return errors.New("failed to allocate memory for extended controls array")
	}
	defer C.free_v4l2_ext_control_array(cExtCtrlArray)

	// Convert Go slice to C array
	// This requires careful handling of pointers and memory if direct slice data access isn't safe/easy.
	// For simplicity and safety, copy element by element.
	// A more optimized approach might use unsafe.Pointer arithmetic if the structs are identical.
	var v4l2Ctrl C.struct_v4l2_ext_control
	for i, ctrl := range ctrls {
		v4l2Ctrl.id = C.uint(ctrl.ID)
		// Assuming CtrlValue is int32 and maps directly to C.int for the union's value field.
		// This part is tricky due to the union C.struct_v4l2_ext_control.anon0
		// The original code writes directly to the anon0 field.
		// For standard integer controls, this is usually fine.
		// For other types (64-bit, string, pointers), this would need more careful handling.
		*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(cExtCtrlArray)) + uintptr(i)*unsafe.Sizeof(v4l2Ctrl))) = C.int(ctrl.Value)
	}


	var v4l2Ctrls C.struct_v4l2_ext_controls
	// The C union for which_ctrl can be set by casting the address of the field.
	// If whichCtrl is 0, it means V4L2_CTRL_WHICH_CUR_VAL, otherwise it's a CtrlClass.
	if whichCtrl == 0 { // Assuming 0 implies V4L2_CTRL_WHICH_CUR_VAL
		v4l2Ctrls.which = C.V4L2_CTRL_WHICH_CUR_VAL
	} else {
		v4l2Ctrls.which = C.uint(whichCtrl)
	}
	v4l2Ctrls.count = C.uint(numCtrl)
	v4l2Ctrls.controls = cExtCtrlArray // Pointer to the C array

	if err := send(fd, C.VIDIOC_S_EXT_CTRLS, uintptr(unsafe.Pointer(&v4l2Ctrls))); err != nil {
		return fmt.Errorf("set ext controls: %w", err)
	}

	return nil
}

// GetExtControl retrieves detailed information about a specific extended control, *including* its current value.
// It combines the functionality of QueryExtControlInfo and GetExtControlValue.
// It takes the file descriptor of the V4L2 device and the control's CtrlID.
// It returns a Control struct populated with all attributes and the current value,
// and an error if either querying control info or getting its value fails.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/extended-controls.html
func GetExtControl(fd uintptr, id CtrlID) (Control, error) {
	control, err := QueryExtControlInfo(fd, id)
	if err != nil {
		return Control{}, fmt.Errorf("get ext control: %w", err)
	}

	// retrieve control value
	// GetExtControlValue internally uses VIDIOC_G_EXT_CTRLS which expects an array,
	// but is used here for a single control.
	ctrlValue, err := GetExtControlValue(fd, id)
	if err != nil {
		return Control{}, fmt.Errorf("get ext control: query value for id %d: %w", id, err)
	}

	control.Value = ctrlValue
	control.fd = fd
	return control, nil
}

// QueryExtControlInfo retrieves detailed information about a specific extended control identified by its CtrlID,
// but *without* its current value. To get the current value, use GetExtControl or GetExtControlValue.
// It takes the file descriptor of the V4L2 device and the control ID.
// It returns a Control struct populated with the control's attributes (name, type, min, max, step, default, flags, etc.)
// and an error if the VIDIOC_QUERY_EXT_CTRL ioctl call fails.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-queryctrl.html#extended-control-example
func QueryExtControlInfo(fd uintptr, id CtrlID) (Control, error) {
	// query control information
	var qryCtrl C.struct_v4l2_query_ext_ctrl
	qryCtrl.id = C.uint(id)

	if err := send(fd, C.VIDIOC_QUERY_EXT_CTRL, uintptr(unsafe.Pointer(&qryCtrl))); err != nil {
		return Control{}, fmt.Errorf("query ext control info: VIDIOC_QUERY_EXT_CTRL: id %d: %w", id, err)
	}
	control := makeExtControl(qryCtrl)
	control.fd = fd
	return control, nil
}

// QueryAllExtControls iterates through all available extended controls on the device and retrieves
// their information (name, type, min, max, etc.), but *without* their current values.
// To get current values, use GetExtControlValue or GetExtControl for each specific control.
// It takes the file descriptor of the V4L2 device.
//
// The iteration starts by looking for the "next" control relative to `CtrlClassCodec`.
// This means it will typically list controls from the Codec class onwards.
// To query all controls including User class, the starting `cid` might need to be `C.V4L2_CTRL_FLAG_NEXT_CTRL`
// or `CtrlClassUser | C.V4L2_CTRL_FLAG_NEXT_CTRL`.
//
// It returns a slice of Control structs and an error if querying fails.
// The iteration stops when the driver returns an error indicating no more controls are available (typically ErrorBadArgument).
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-queryctrl.html#iterating-over-controls
func QueryAllExtControls(fd uintptr) (result []Control, err error) {
	// The starting point for querying all extended controls can be V4L2_CTRL_FLAG_NEXT_CTRL.
	// The original code starts with CtrlClassCodec, which might be specific to a use case.
	// For a general "all controls", V4L2_CTRL_FLAG_NEXT_CTRL is more appropriate.
	// However, to match original behavior, using CtrlClassCodec as base for first query.
	cid := CtrlClassCodec | uint32(C.V4L2_CTRL_FLAG_NEXT_CTRL)
	// To query truly all, one might start with:
	// cid := uint32(C.V4L2_CTRL_FLAG_NEXT_CTRL)

	for {
		control, err := QueryExtControlInfo(fd, cid)
		if err != nil {
			// If ErrorBadArgument is returned, it means no more controls with IDs greater than the current one.
			if errors.Is(err, ErrorBadArgument) {
				break // Successfully finished iterating.
			}
			// For other errors, return the error and the controls found so far.
			return result, fmt.Errorf("query all ext controls: iteration error on ID 0x%x: %w", cid, err)
		}
		result = append(result, control)
		// Prepare the ID for the next control.
		cid = control.ID | uint32(C.V4L2_CTRL_FLAG_NEXT_CTRL)
	}

	return result, nil
}

func makeExtControl(qryCtrl C.struct_v4l2_query_ext_ctrl) Control {
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
