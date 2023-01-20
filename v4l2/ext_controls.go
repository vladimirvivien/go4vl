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

// GetExtControlValue retrieves the value for an extended control with the specified id.
// See https://linuxtv.org/downloads/v4l-dvb-apis-new/userspace-api/v4l/extended-controls.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1745
func GetExtControlValue(fd uintptr, ctrlID CtrlID) (CtrlValue, error) {
	var v4l2Ctrl C.struct_v4l2_ext_control
	v4l2Ctrl.id = C.uint(ctrlID)
	v4l2Ctrl.size = 0
	if err := send(fd, C.VIDIOC_G_EXT_CTRLS, uintptr(unsafe.Pointer(&v4l2Ctrl))); err != nil {
		return 0, fmt.Errorf("get ext controls: %w", err)
	}
	return *(*CtrlValue)(unsafe.Pointer(&v4l2Ctrl.anon0[0])), nil
}

// SetExtControlValue saves the value for an extended control with the specified id.
// See https://linuxtv.org/downloads/v4l-dvb-apis-new/userspace-api/v4l/extended-controls.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1745
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

// SetExtControlValues implements code to save one or more extended controls at once using the
// v4l2_ext_controls structure.
// https://linuxtv.org/downloads/v4l-dvb-apis-new/userspace-api/v4l/extended-controls.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1774
func SetExtControlValues(fd uintptr, whichCtrl CtrlClass, ctrls []Control) error {
	numCtrl := len(ctrls)

	var v4l2CtrlArray []C.struct_v4l2_ext_control

	for _, ctrl := range ctrls {
		var v4l2Ctrl C.struct_v4l2_ext_control
		v4l2Ctrl.id = C.uint(ctrl.ID)
		*(*C.int)(unsafe.Pointer(&v4l2Ctrl.anon0[0])) = *(*C.int)(unsafe.Pointer(&ctrl.Value))

		v4l2CtrlArray = append(v4l2CtrlArray, v4l2Ctrl)
	}

	var v4l2Ctrls C.struct_v4l2_ext_controls
	*(*uint32)(unsafe.Pointer(&v4l2Ctrls.anon0[0])) = whichCtrl
	v4l2Ctrls.count = C.uint(numCtrl)
	v4l2Ctrls.controls = (*C.struct_v4l2_ext_control)(unsafe.Pointer(&v4l2CtrlArray))

	if err := send(fd, C.VIDIOC_S_EXT_CTRLS, uintptr(unsafe.Pointer(&v4l2Ctrls))); err != nil {
		return fmt.Errorf("set ext controls: %w", err)
	}

	return nil
}

// GetExtControl retrieves information (query) and current value for the specified control.
// See https://linuxtv.org/downloads/v4l-dvb-apis-new/userspace-api/v4l/extended-controls.html
func GetExtControl(fd uintptr, id CtrlID) (Control, error) {
	control, err := QueryExtControlInfo(fd, id)
	if err != nil {
		return Control{}, fmt.Errorf("get control: %w", err)
	}

	// retrieve control value
	ctrlValue, err := GetExtControlValue(fd, uint32(id))
	if err != nil {
		return Control{}, fmt.Errorf("get control: %w", id, err)
	}

	control.Value = ctrlValue
	control.fd = fd
	return control, nil
}

// QueryExtControlInfo queries information about the specified ext control without its current value.
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

// QueryAllExtControls loop through all available ext controls and query the information for
// all controls without their current values (use GetExtControlValue to get current values).
func QueryAllExtControls(fd uintptr) (result []Control, err error) {

	cid := CtrlClassCodec | uint32(C.V4L2_CTRL_FLAG_NEXT_CTRL)
	for {
		control, err := QueryExtControlInfo(fd, cid)
		if err != nil {
			if errors.Is(err, ErrorBadArgument) {
				break
			}
			return result, fmt.Errorf("query all ext controls: %w", err)
		}
		result = append(result, control)
		// setup next id
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
