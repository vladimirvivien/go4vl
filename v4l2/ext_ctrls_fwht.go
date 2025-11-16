package v4l2

/*
#include <linux/videodev2.h>
#include <linux/v4l2-controls.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ControlFWHTParams (v4l2_ctrl_fwht_params)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1659
type ControlFWHTParams struct {
	BackwardRefTimestamp uint64
	Version              uint32
	Width                uint32
	Height               uint32
	Flags                uint32
	// Colorspace           ColorspaceType
	// XFerFunc             XferFunctionType
	// YCbCrEncoding        YCbCrEncodingType
	// Quantization         QuantizationType
}

// FWHT Stateless Codec Control IDs
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlFWHTParams CtrlID = C.V4L2_CID_STATELESS_FWHT_PARAMS
)

// Type-safe helper methods for ExtControls

// AddFWHTParams adds an FWHT Parameters control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddFWHTParams(params *ControlFWHTParams) error {
	if params == nil {
		return fmt.Errorf("FWHT params cannot be nil")
	}
	size := unsafe.Sizeof(*params)
	data := unsafe.Slice((*byte)(unsafe.Pointer(params)), size)
	ec.AddCompound(CtrlFWHTParams, data)
	return nil
}

// Type-safe helper methods for ExtControl (reading values back)

// GetFWHTParams retrieves the FWHT Parameters from a control value.
func (ec *ExtControl) GetFWHTParams() (*ControlFWHTParams, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlFWHTParams{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid FWHT params size: got %d, expected %d", len(data), expectedSize)
	}
	params := &ControlFWHTParams{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(params)), expectedSize), data)
	return params, nil
}
