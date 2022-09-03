package v4l2

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../include/
#include <linux/videodev2.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// TODO - Implementation of extended controls (v4l2_ext_control) is paused for now,
// so that efforts can be focused on other parts of the API. This can resumed
// later when type v4l2_ext_control and v4l2_ext_controls are better understood.

// ExtControl (v4l2_ext_control)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1730
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ext-ctrls.html
type ExtControl struct {
	ID      uint32
	Value   int32
	Value64 int64

	// NOT USED YET
	Size               uint32
	String             string
	PU8                uint
	PU16               uint16
	PU32               uint32
	PArea              Area
	PH264SPS           ControlH264SPS
	PH264PPS           ControlH264PPS
	PH264ScalingMatrix ControlH264ScalingMatrix
	H264PredWeights    ControlH264PredictionWeights
	PH264SliceParams   ControlH264SliceParams
	PH264DecodeParams  ControlH264DecodeParams
	PFWHTParams        ControlFWHTParams
	PVP8Frame          ControlVP8Frame
	PMPEG2Sequence     ControlMPEG2Sequence
	PMPEG2Picture      ControlMPEG2Picture
	PMPEG2Quantization ControlMPEG2Quantization
}

func setExtCtrls(fd uintptr, whichCtrl uint32, ctrls []ExtControl) error {
	numCtrl := len(ctrls)

	var v4l2CtrlArray []C.struct_v4l2_ext_control

	for _, ctrl := range ctrls {
		var v4l2Ctrl C.struct_v4l2_ext_control
		v4l2Ctrl.id = C.uint(ctrl.ID)
		*(*C.uint)(unsafe.Pointer(&v4l2Ctrl.anon0[0])) = *(*C.uint)(unsafe.Pointer(&ctrl.Value))

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

// GetExtControls retrieve one or more controls
// func GetExtControls(fd uintptr, controls []ExtControl) (ExtControls, error) {
// 	if true {
// 		// TODO remove when supported
// 		return ExtControls{}, fmt.Errorf("unsupported")
// 	}

// 	var ctrls C.struct_v4l2_ext_controls
// 	ctrls.count = C.uint(len(controls))

// 	// prepare control requests
// 	var Cctrls []C.struct_v4l2_ext_control
// 	for _, control := range controls {
// 		var Cctrl C.struct_v4l2_ext_control
// 		Cctrl.id = C.uint(control.ID)
// 		Cctrl.size = C.uint(control.Size)
// 		*(*ExtControlUnion)(unsafe.Pointer(&Cctrl.anon0[0])) = control.Ctrl
// 		Cctrls = append(Cctrls, Cctrl)
// 	}
// 	ctrls.controls = (*C.struct_v4l2_ext_control)(unsafe.Pointer(&ctrls.controls))

// 	if err := send(fd, C.VIDIOC_G_EXT_CTRLS, uintptr(unsafe.Pointer(&ctrls))); err != nil {
// 		return ExtControls{}, fmt.Errorf("get ext controls: %w", err)
// 	}

// 	// gather returned controls
// 	retCtrls := ExtControls{
// 		Count:      uint32(ctrls.count),
// 		ErrorIndex: uint32(ctrls.error_idx),
// 	}
// 	// extract controls array
// 	Cctrls = *(*[]C.struct_v4l2_ext_control)(unsafe.Pointer(&ctrls.controls))
// 	for _, Cctrl := range Cctrls {
// 		extCtrl := ExtControl{
// 			ID:   uint32(Cctrl.id),
// 			Size: uint32(Cctrl.size),
// 			Ctrl: *(*ExtControlUnion)(unsafe.Pointer(&Cctrl.anon0[0])),
// 		}
// 		retCtrls.Controls = append(retCtrls.Controls, extCtrl)
// 	}

// 	return retCtrls, nil
// }
