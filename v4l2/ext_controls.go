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

import (
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// Additional control classes and which-value selectors for VIDIOC_{G,S,TRY}_EXT_CTRLS.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlClassFMTx       CtrlClass = cV4L2_CTRL_CLASS_FM_TX       // FM Modulator controls
	CtrlClassFMRx       CtrlClass = cV4L2_CTRL_CLASS_FM_RX       // FM Receiver controls
	CtrlClassRFTuner    CtrlClass = cV4L2_CTRL_CLASS_RF_TUNER    // RF tuner controls
	CtrlWhichCurVal     CtrlClass = cV4L2_CTRL_WHICH_CUR_VAL     // Get/set current control values
	CtrlWhichDefVal     CtrlClass = cV4L2_CTRL_WHICH_DEF_VAL     // Get default control values
	CtrlWhichRequestVal CtrlClass = cV4L2_CTRL_WHICH_REQUEST_VAL // Bind controls to a Media Request FD
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
	Value any // int32, int64, string, or []byte
}

// NewExtControl creates a new extended control with the specified ID for reading.
// The value is read as an int32; to read a 64-bit control (V4L2_CTRL_TYPE_INTEGER64),
// use NewExtControlWithValue64(id, 0) instead.
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
	RequestFD int32
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

// SetRequestFD binds this control collection to a Media Request file descriptor
// (allocated via AllocMediaRequest) and selects CtrlWhichRequestVal.
func (ec *ExtControls) SetRequestFD(requestFD int32) {
	ec.RequestFD = requestFD
	if requestFD > 0 {
		ec.Class = CtrlWhichRequestVal
	}
}

// GetRequestFD returns the bound Media Request file descriptor (0 if unbound).
func (ec *ExtControls) GetRequestFD() int32 {
	return ec.RequestFD
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

func addCompoundStruct[T any](ec *ExtControls, id CtrlID, name string, val *T) error {
	if ec == nil {
		return fmt.Errorf("%s control collection cannot be nil", name)
	}
	if val == nil {
		return fmt.Errorf("%s cannot be nil", name)
	}
	size := unsafe.Sizeof(*val)
	data := unsafe.Slice((*byte)(unsafe.Pointer(val)), size)
	ec.AddCompound(id, data)
	return nil
}

func getCompoundStruct[T any](ec *ExtControl, name string) (*T, error) {
	if ec == nil {
		return nil, fmt.Errorf("control cannot be nil")
	}
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	var zero T
	expectedSize := int(unsafe.Sizeof(zero))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid %s size: got %d, expected %d", name, len(data), expectedSize)
	}
	out := new(T)
	copy(unsafe.Slice((*byte)(unsafe.Pointer(out)), expectedSize), data)
	return out, nil
}

// AddExtControlStruct adds a typed compound control structure to ec under id.
func AddExtControlStruct[T any](ec *ExtControls, id CtrlID, val *T) error {
	return addCompoundStruct(ec, id, fmt.Sprintf("control 0x%x", id), val)
}

// GetExtControlStruct finds the control with id in ec and decodes its compound payload into T.
func GetExtControlStruct[T any](ec *ExtControls, id CtrlID) (*T, error) {
	if ec == nil {
		return nil, fmt.Errorf("control collection cannot be nil")
	}
	for _, ctrl := range ec.controls {
		if ctrl != nil && ctrl.ID == id {
			return getCompoundStruct[T](ctrl, fmt.Sprintf("control 0x%x", id))
		}
	}
	return nil, fmt.Errorf("control 0x%x not found", id)
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

// prepareForIoctl prepares structures for ioctl calls and pins Go-allocated
// control slices and compound buffers across the syscall.
func prepareForIoctl(ec *ExtControls) (*cV4L2ExtControls, func(), error) {
	count := len(ec.controls)
	if count == 0 {
		return nil, nil, fmt.Errorf("no controls to process")
	}

	var pinner runtime.Pinner
	ctrlSlice := make([]cV4L2ExtControl, count)
	pinner.Pin(&ctrlSlice[0])

	// Retain buffers in a slice so they remain referenced alongside pinning until cleanup runs.
	var pinnedBufs [][]byte

	for i, ctrl := range ec.controls {
		ctrlSlice[i].Id = ctrl.ID
		ctrlSlice[i].Size = 0

		switch v := ctrl.Value.(type) {
		case int32:
			ctrlSlice[i].Value = v
		case int64:
			b := (*[8]byte)(unsafe.Pointer(&ctrlSlice[i].Value))
			binary.NativeEndian.PutUint64(b[:], uint64(v))
		case string:
			buf := make([]byte, len(v)+1)
			copy(buf, v)
			pinner.Pin(&buf[0])
			pinnedBufs = append(pinnedBufs, buf)
			ctrlSlice[i].Size = uint32(len(buf))
			b := (*[8]byte)(unsafe.Pointer(&ctrlSlice[i].Value))
			binary.NativeEndian.PutUint64(b[:], uint64(uintptr(unsafe.Pointer(&buf[0]))))
		case []byte:
			if len(v) == 0 {
				pinner.Unpin()
				return nil, nil, fmt.Errorf("compound data for control %d cannot be empty", ctrl.ID)
			}
			buf := make([]byte, len(v))
			copy(buf, v)
			pinner.Pin(&buf[0])
			pinnedBufs = append(pinnedBufs, buf)
			ctrlSlice[i].Size = uint32(len(buf))
			b := (*[8]byte)(unsafe.Pointer(&ctrlSlice[i].Value))
			binary.NativeEndian.PutUint64(b[:], uint64(uintptr(unsafe.Pointer(&buf[0]))))
		case nil:
			// Reading control - no value to set
		default:
			pinner.Unpin()
			return nil, nil, fmt.Errorf("unsupported control value type for control %d: %T", ctrl.ID, v)
		}
	}

	cleanup := func() {
		pinner.Unpin()
		runtime.KeepAlive(pinnedBufs)
		runtime.KeepAlive(ctrlSlice)
	}

	var v4l2Ctrls cV4L2ExtControls
	v4l2Ctrls.Ctrl_class = ec.Class
	if ec.RequestFD > 0 {
		v4l2Ctrls.Request_fd = ec.RequestFD
		if v4l2Ctrls.Ctrl_class == 0 {
			v4l2Ctrls.Ctrl_class = CtrlWhichRequestVal
		}
	}
	v4l2Ctrls.Count = uint32(count)
	v4l2Ctrls.Controls = &ctrlSlice[0]

	return &v4l2Ctrls, cleanup, nil
}

// updateFromIoctl updates Go structures from C structures after ioctl.
func updateFromIoctl(ec *ExtControls, v4l2Ctrls *cV4L2ExtControls) {
	count := len(ec.controls)
	ctrlSlice := unsafe.Slice(v4l2Ctrls.Controls, count)

	ec.errorIdx = v4l2Ctrls.Error_idx

	for i, ctrl := range ec.controls {
		cCtrl := &ctrlSlice[i]

		switch ctrl.Value.(type) {
		case int64:
			// A 64-bit control: value64 shares the union with value, so its 8 bytes
			// start at Value. The driver returns the current value (get) or the
			// value it applied (set, try).
			b := (*[8]byte)(unsafe.Pointer(&cCtrl.Value))
			ctrl.Value = int64(binary.NativeEndian.Uint64(b[:]))
		case nil:
			// A read of a control added without a value
			if cCtrl.Size == 0 {
				ctrl.Value = cCtrl.Value
			} else {
				var ptr unsafe.Pointer
				*(*[8]byte)(unsafe.Pointer(&ptr)) = *(*[8]byte)(unsafe.Pointer(&cCtrl.Value))
				if ptr != nil {
					src := unsafe.Slice((*byte)(ptr), int(cCtrl.Size))
					dst := make([]byte, len(src))
					copy(dst, src)
					ctrl.Value = dst
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

	if err := send(fd, cVIDIOC_G_EXT_CTRLS, unsafe.Pointer(v4l2Ctrls)); err != nil {
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

	if err := send(fd, cVIDIOC_S_EXT_CTRLS, unsafe.Pointer(v4l2Ctrls)); err != nil {
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

	if err := send(fd, cVIDIOC_TRY_EXT_CTRLS, unsafe.Pointer(v4l2Ctrls)); err != nil {
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

// GetExtControlValue64 retrieves the value of a 64-bit extended control
// (V4L2_CTRL_TYPE_INTEGER64, such as V4L2_CID_PIXEL_RATE) with the specified id.
// GetExtControlValue reads a value as an int32, which truncates such a control.
func GetExtControlValue64(fd uintptr, ctrlID CtrlID) (int64, error) {
	ctrls := NewExtControls()
	ctrl := NewExtControlWithValue64(ctrlID, 0)
	ctrls.Add(ctrl)

	if err := GetExtControls(fd, ctrls); err != nil {
		return 0, err
	}

	return ctrl.GetValue64(), nil
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
	var qryCtrl cV4L2QueryExtCtrl
	qryCtrl.Id = id

	if err := send(fd, cVIDIOC_QUERY_EXT_CTRL, unsafe.Pointer(&qryCtrl)); err != nil {
		return Control{}, fmt.Errorf("query ext control info: VIDIOC_QUERY_EXT_CTRL: id %d: %w", id, err)
	}
	control := makeExtControl(qryCtrl)
	control.fd = fd
	return control, nil
}

// QueryAllExtControlsByClass enumerates extended controls starting from the specified control class
// (pass 0 to enumerate across all control classes, including compound controls when combined with CtrlFlagNextCompound).
func QueryAllExtControlsByClass(fd uintptr, class CtrlClass) (result []Control, err error) {
	cid := class | cV4L2_CTRL_FLAG_NEXT_CTRL
	for {
		control, err := QueryExtControlInfo(fd, cid)
		if err != nil {
			if errors.Is(err, ErrorBadArgument) && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("query all ext controls: %w", err)
		}
		result = append(result, control)
		cid = control.ID | cV4L2_CTRL_FLAG_NEXT_CTRL
	}

	return result, nil
}

// QueryAllExtControls loop through all available ext controls and query the information for
// all controls without their current values (use GetExtControlValue to get current values).
func QueryAllExtControls(fd uintptr) (result []Control, err error) {
	return QueryAllExtControlsByClass(fd, CtrlClassCodec)
}

func makeExtControl(qryCtrl cV4L2QueryExtCtrl) Control {
	return Control{
		Type:    CtrlType(qryCtrl.Type),
		ID:      qryCtrl.Id,
		Name:    cString(qryCtrl.Name[:]),
		Maximum: int32(qryCtrl.Maximum),
		Minimum: int32(qryCtrl.Minimum),
		Step:    int32(qryCtrl.Step),
		Default: int32(qryCtrl.Default_value),
		Flags:   qryCtrl.Flags,

		Minimum64: qryCtrl.Minimum,
		Maximum64: qryCtrl.Maximum,
		Step64:    qryCtrl.Step,
		Default64: qryCtrl.Default_value,
	}
}
