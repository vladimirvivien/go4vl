package v4l2

// ext_controls.go provides V4L2 Extended Controls API support.
//
// Extended controls allow advanced control of devices with support for:
// - Compound controls (arrays, strings, structs)
// - Control classes (user, codec, camera, etc.)
// - Atomic multi-control operations
// - 64-bit values
// - Complex codec parameters (H.264 SPS/PPS, etc.)
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ext-ctrls.html
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/extended-controls.html

/*
#include <linux/videodev2.h>
#include <linux/v4l2-controls.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// Additional control classes not defined in control_values.go
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlClassFMTx    CtrlClass = C.V4L2_CTRL_CLASS_FM_TX    // FM Modulator controls
	CtrlClassFMRx    CtrlClass = C.V4L2_CTRL_CLASS_FM_RX    // FM Receiver controls
	CtrlClassRFTuner CtrlClass = C.V4L2_CTRL_CLASS_RF_TUNER // RF tuner controls
)

// ExtControl represents a single extended control.
// Extended controls support simple values, 64-bit values, strings, and compound types.
//
// Memory is managed automatically - no manual Free() calls needed.
// Use ExtControls.Add() methods to add controls to a collection.
//
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type ExtControl struct {
	ID    CtrlID
	Value interface{} // int32, int64, string, or []byte
}

// NewExtControl creates a new extended control with the specified ID for reading.
func NewExtControl(id CtrlID) *ExtControl {
	return &ExtControl{
		ID: id,
	}
}

// NewExtControlWithValue creates a new extended control with an int32 value.
func NewExtControlWithValue(id CtrlID, value int32) *ExtControl {
	return &ExtControl{
		ID:    id,
		Value: value,
	}
}

// NewExtControlWithValue64 creates a new extended control with an int64 value.
func NewExtControlWithValue64(id CtrlID, value int64) *ExtControl {
	return &ExtControl{
		ID:    id,
		Value: value,
	}
}

// NewExtControlWithString creates a new extended control with a string value.
func NewExtControlWithString(id CtrlID, value string) *ExtControl {
	return &ExtControl{
		ID:    id,
		Value: value,
	}
}

// NewExtControlWithCompound creates a new extended control with compound data.
func NewExtControlWithCompound(id CtrlID, data []byte) *ExtControl {
	return &ExtControl{
		ID:    id,
		Value: data,
	}
}

// GetID returns the control ID.
func (ec *ExtControl) GetID() CtrlID {
	return ec.ID
}

// GetValue returns the value as int32. Panics if value is not int32.
func (ec *ExtControl) GetValue() int32 {
	if v, ok := ec.Value.(int32); ok {
		return v
	}
	return 0
}

// GetValue64 returns the value as int64. Panics if value is not int64.
func (ec *ExtControl) GetValue64() int64 {
	if v, ok := ec.Value.(int64); ok {
		return v
	}
	return 0
}

// GetString returns the value as string. Returns empty string if not a string.
func (ec *ExtControl) GetString() string {
	if v, ok := ec.Value.(string); ok {
		return v
	}
	return ""
}

// GetCompoundData returns the value as []byte. Returns nil if not compound data.
func (ec *ExtControl) GetCompoundData() []byte {
	if v, ok := ec.Value.([]byte); ok {
		return v
	}
	return nil
}

// ExtControls represents a collection of extended controls.
// It supports atomic multi-control operations and control class filtering.
//
// Memory is managed automatically - no manual Free() calls needed.
//
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type ExtControls struct {
	Class     CtrlClass
	controls  []*ExtControl
	errorIdx  uint32 // Index of control that caused error (after ioctl)
}

// NewExtControls creates a new extended controls collection.
func NewExtControls() *ExtControls {
	return &ExtControls{
		controls: make([]*ExtControl, 0),
	}
}

// NewExtControlsWithClass creates a new extended controls collection for a specific class.
func NewExtControlsWithClass(class CtrlClass) *ExtControls {
	return &ExtControls{
		Class:    class,
		controls: make([]*ExtControl, 0),
	}
}

// SetClass sets the control class for this collection.
func (ec *ExtControls) SetClass(class CtrlClass) {
	ec.Class = class
}

// GetClass returns the control class.
func (ec *ExtControls) GetClass() CtrlClass {
	return ec.Class
}

// Add adds an extended control to the collection.
func (ec *ExtControls) Add(ctrl *ExtControl) {
	ec.controls = append(ec.controls, ctrl)
}

// AddValue adds a control with int32 value.
func (ec *ExtControls) AddValue(id CtrlID, value int32) {
	ec.controls = append(ec.controls, NewExtControlWithValue(id, value))
}

// AddValue64 adds a control with int64 value.
func (ec *ExtControls) AddValue64(id CtrlID, value int64) {
	ec.controls = append(ec.controls, NewExtControlWithValue64(id, value))
}

// AddString adds a control with string value.
func (ec *ExtControls) AddString(id CtrlID, value string) {
	ec.controls = append(ec.controls, NewExtControlWithString(id, value))
}

// AddCompound adds a control with compound data.
func (ec *ExtControls) AddCompound(id CtrlID, data []byte) {
	ec.controls = append(ec.controls, NewExtControlWithCompound(id, data))
}

// GetControls returns all controls in the collection.
func (ec *ExtControls) GetControls() []*ExtControl {
	return ec.controls
}

// GetErrorIndex returns the index of the control that caused an error (if any).
func (ec *ExtControls) GetErrorIndex() uint32 {
	return ec.errorIdx
}

// Count returns the number of controls in the collection.
func (ec *ExtControls) Count() int {
	return len(ec.controls)
}

// prepareForIoctl prepares C structures for ioctl call.
// Returns cleanup function to free all allocated memory.
func prepareForIoctl(ec *ExtControls) (*C.struct_v4l2_ext_controls, func(), error) {
	count := len(ec.controls)
	if count == 0 {
		return nil, nil, fmt.Errorf("no controls to process")
	}

	// Allocate C array for controls using malloc
	ctrlArray := (*C.struct_v4l2_ext_control)(C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof(C.struct_v4l2_ext_control{}))))
	if ctrlArray == nil {
		return nil, nil, fmt.Errorf("failed to allocate memory for control array")
	}

	// Track allocated memory for cleanup
	allocatedStrings := make([]*C.char, 0)
	allocatedCompound := make([]unsafe.Pointer, 0)

	// Cleanup function - frees all C memory allocated in this function
	cleanup := func() {
		for _, str := range allocatedStrings {
			C.free(unsafe.Pointer(str))
		}
		for _, ptr := range allocatedCompound {
			C.free(ptr)
		}
		C.free(unsafe.Pointer(ctrlArray))
	}

	// Convert C pointer to Go slice using the "large array trick":
	// 1 << 30 = 1073741824 elements - we create a very large fixed-size array type
	// Then immediately slice it to the actual count. This is safe because:
	// - We only access [:count] elements (which we know exist)
	// - Go won't bounds-check the large array size (it trusts our unsafe conversion)
	// - This is a standard Go pattern for C array interop
	ctrlSlice := (*[1 << 30]C.struct_v4l2_ext_control)(unsafe.Pointer(ctrlArray))[:count:count]
	for i, ctrl := range ec.controls {
		ctrlSlice[i].id = C.__u32(ctrl.ID)
		ctrlSlice[i].size = 0

		// Set value based on type
		switch v := ctrl.Value.(type) {
		case int32:
			// Cast union field anon0 to *int32 and set the value
			*(*C.__s32)(unsafe.Pointer(&ctrlSlice[i].anon0[0])) = C.__s32(v)
		case int64:
			// Cast union field anon0 to *int64 and set the value
			*(*C.__s64)(unsafe.Pointer(&ctrlSlice[i].anon0[0])) = C.__s64(v)
		case string:
			// Allocate C string and store pointer in union
			cstr := C.CString(v)
			allocatedStrings = append(allocatedStrings, cstr)
			ctrlSlice[i].size = C.__u32(len(v) + 1)
			// Cast union field to **char and store the string pointer
			*(**C.char)(unsafe.Pointer(&ctrlSlice[i].anon0[0])) = cstr
		case []byte:
			if len(v) == 0 {
				cleanup()
				return nil, nil, fmt.Errorf("compound data for control %d cannot be empty", ctrl.ID)
			}
			// Allocate C memory for compound data
			ptr := C.malloc(C.size_t(len(v)))
			if ptr == nil {
				cleanup()
				return nil, nil, fmt.Errorf("failed to allocate memory for compound data")
			}
			allocatedCompound = append(allocatedCompound, ptr)
			// Copy Go bytes to C memory
			C.memcpy(ptr, unsafe.Pointer(&v[0]), C.size_t(len(v)))
			ctrlSlice[i].size = C.__u32(len(v))
			// Cast union field to unsafe.Pointer* and store the data pointer
			*(*unsafe.Pointer)(unsafe.Pointer(&ctrlSlice[i].anon0[0])) = ptr
		case nil:
			// Reading control - no value to set
		default:
			cleanup()
			return nil, nil, fmt.Errorf("unsupported control value type for control %d: %T", ctrl.ID, v)
		}
	}

	// Set up v4l2_ext_controls structure
	var v4l2Ctrls C.struct_v4l2_ext_controls
	// Access the 'which' field (first field in struct, in union with ctrl_class)
	*(*C.__u32)(unsafe.Pointer(&v4l2Ctrls)) = C.__u32(ec.Class)
	v4l2Ctrls.count = C.__u32(count)
	v4l2Ctrls.controls = ctrlArray

	return &v4l2Ctrls, cleanup, nil
}

// updateFromIoctl updates Go structures from C structures after ioctl.
func updateFromIoctl(ec *ExtControls, v4l2Ctrls *C.struct_v4l2_ext_controls) {
	count := len(ec.controls)
	// Convert C pointer to Go slice using the "large array trick" (see prepareForIoctl for explanation)
	ctrlSlice := (*[1 << 30]C.struct_v4l2_ext_control)(unsafe.Pointer(v4l2Ctrls.controls))[:count:count]

	// Update error index from C structure
	ec.errorIdx = uint32(v4l2Ctrls.error_idx)

	// Copy values back from C array
	for i, ctrl := range ec.controls {
		cCtrl := &ctrlSlice[i]

		// Only update value if it was a read operation (Value was nil)
		if ctrl.Value == nil {
			// Determine value type based on size
			if cCtrl.size == 0 {
				// Simple 32-bit value: cast union field to *int32 and read
				ctrl.Value = int32(*(*C.__s32)(unsafe.Pointer(&cCtrl.anon0[0])))
			} else {
				// Could be string or compound - cast union field to unsafe.Pointer* and read pointer
				ptr := *(*unsafe.Pointer)(unsafe.Pointer(&cCtrl.anon0[0]))
				if ptr != nil {
					// Copy C bytes to Go slice (works for both strings and compound data)
					ctrl.Value = C.GoBytes(ptr, C.int(cCtrl.size))
				}
			}
		}
	}
}

// GetExtControls retrieves values for a collection of extended controls.
// Memory is managed automatically - no Free() calls needed.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ext-ctrls.html
func GetExtControls(fd uintptr, ctrls *ExtControls) error {
	v4l2Ctrls, cleanup, err := prepareForIoctl(ctrls)
	if err != nil {
		return fmt.Errorf("v4l2: VIDIOC_G_EXT_CTRLS prepare failed: %w", err)
	}
	defer cleanup()

	if err := send(fd, C.VIDIOC_G_EXT_CTRLS, uintptr(unsafe.Pointer(v4l2Ctrls))); err != nil {
		updateFromIoctl(ctrls, v4l2Ctrls) // Update error index
		return fmt.Errorf("v4l2: VIDIOC_G_EXT_CTRLS failed: error at index %d: %w", ctrls.GetErrorIndex(), err)
	}

	updateFromIoctl(ctrls, v4l2Ctrls)
	return nil
}

// SetExtControls sets values for a collection of extended controls atomically.
// Memory is managed automatically - no Free() calls needed.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ext-ctrls.html
func SetExtControls(fd uintptr, ctrls *ExtControls) error {
	v4l2Ctrls, cleanup, err := prepareForIoctl(ctrls)
	if err != nil {
		return fmt.Errorf("v4l2: VIDIOC_S_EXT_CTRLS prepare failed: %w", err)
	}
	defer cleanup()

	if err := send(fd, C.VIDIOC_S_EXT_CTRLS, uintptr(unsafe.Pointer(v4l2Ctrls))); err != nil {
		updateFromIoctl(ctrls, v4l2Ctrls) // Update error index
		return fmt.Errorf("v4l2: VIDIOC_S_EXT_CTRLS failed: error at index %d: %w", ctrls.GetErrorIndex(), err)
	}

	updateFromIoctl(ctrls, v4l2Ctrls)
	return nil
}

// TryExtControls validates values for extended controls without applying them.
// Memory is managed automatically - no Free() calls needed.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ext-ctrls.html
func TryExtControls(fd uintptr, ctrls *ExtControls) error {
	v4l2Ctrls, cleanup, err := prepareForIoctl(ctrls)
	if err != nil {
		return fmt.Errorf("v4l2: VIDIOC_TRY_EXT_CTRLS prepare failed: %w", err)
	}
	defer cleanup()

	if err := send(fd, C.VIDIOC_TRY_EXT_CTRLS, uintptr(unsafe.Pointer(v4l2Ctrls))); err != nil {
		updateFromIoctl(ctrls, v4l2Ctrls) // Update error index
		return fmt.Errorf("v4l2: VIDIOC_TRY_EXT_CTRLS failed: error at index %d: %w", ctrls.GetErrorIndex(), err)
	}

	updateFromIoctl(ctrls, v4l2Ctrls)
	return nil
}

// Legacy simple control functions for backward compatibility
// These wrap the simpler single-control operations

// GetExtControlValue retrieves the value for an extended control with the specified id.
// See https://linuxtv.org/downloads/v4l-dvb-apis-new/userspace-api/v4l/extended-controls.html
func GetExtControlValue(fd uintptr, ctrlID CtrlID) (CtrlValue, error) {
	ctrls := NewExtControls()
	ctrl := NewExtControl(ctrlID)
	ctrls.Add(ctrl)

	if err := GetExtControls(fd, ctrls); err != nil {
		return 0, err
	}

	return ctrl.GetValue(), nil
}

// SetExtControlValue saves the value for an extended control with the specified id.
// See https://linuxtv.org/downloads/v4l-dvb-apis-new/userspace-api/v4l/extended-controls.html
func SetExtControlValue(fd uintptr, id CtrlID, val CtrlValue) error {
	ctrlInfo, err := QueryExtControlInfo(fd, id)
	if err != nil {
		return fmt.Errorf("set ext control value: id %d: %w", id, err)
	}
	if val < ctrlInfo.Minimum || val > ctrlInfo.Maximum {
		return fmt.Errorf("set ext control value: out-of-range failure: val %d: expected ctrl.Min %d, ctrl.Max %d", val, ctrlInfo.Minimum, ctrlInfo.Maximum)
	}

	ctrls := NewExtControls()
	ctrl := NewExtControlWithValue(id, val)
	ctrls.Add(ctrl)

	return SetExtControls(fd, ctrls)
}

// SetExtControlValues implements code to save one or more extended controls at once using the
// v4l2_ext_controls structure.
// https://linuxtv.org/downloads/v4l-dvb-apis-new/userspace-api/v4l/extended-controls.html
func SetExtControlValues(fd uintptr, whichCtrl CtrlClass, ctrls []Control) error {
	extCtrls := NewExtControlsWithClass(whichCtrl)

	for _, ctrl := range ctrls {
		extCtrl := NewExtControlWithValue(ctrl.ID, ctrl.Value)
		extCtrls.Add(extCtrl)
	}

	return SetExtControls(fd, extCtrls)
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
		return Control{}, fmt.Errorf("get control: id %d: %w", id, err)
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
			if errors.Is(err, ErrorBadArgument) && len(result) > 0 {
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
