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

// ControlMPEG2Sequence (v4l2_ctrl_mpeg2_sequence)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1892
type ControlMPEG2Sequence struct {
	HorizontalSize            uint16
	VerticalSize              uint16
	VBVBufferSize             uint32
	ProfileAndLevelIndication uint16
	ChromaFormat              uint8
	Flags                     uint8
}

// ControlMPEG2Picture (v4l2_ctrl_mpeg2_picture)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1939
type ControlMPEG2Picture struct {
	BackwardRefTimestamp uint64
	ForwardRefTimestamp  uint64
	Flags                uint32
	FCode                [2][2]uint8
	PictureCodingType    uint8
	PictureStructure     uint8
	IntraDCPrecision     uint8
	_                    [5]uint8 // padding
}

// ControlMPEG2Quantization (v4l2_ctrl_mpeg2_quantisation)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1972
type ControlMPEG2Quantization struct {
	IntraQuantizerMatrix          [64]uint8
	NonIntraQuantizerMatrix       [64]uint8
	ChromaIntraQuantizerMatrix    [64]uint8
	ChromaNonIntraQuantizerMatrix [64]uint8
}

// MPEG2 Stateless Codec Control IDs
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlMPEG2Sequence      CtrlID = C.V4L2_CID_STATELESS_MPEG2_SEQUENCE
	CtrlMPEG2Picture       CtrlID = C.V4L2_CID_STATELESS_MPEG2_PICTURE
	CtrlMPEG2Quantization  CtrlID = C.V4L2_CID_STATELESS_MPEG2_QUANTISATION
)

// Type-safe helper methods for ExtControls

// AddMPEG2Sequence adds an MPEG2 Sequence control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddMPEG2Sequence(seq *ControlMPEG2Sequence) error {
	if seq == nil {
		return fmt.Errorf("MPEG2 sequence cannot be nil")
	}
	size := unsafe.Sizeof(*seq)
	data := unsafe.Slice((*byte)(unsafe.Pointer(seq)), size)
	ec.AddCompound(CtrlMPEG2Sequence, data)
	return nil
}

// AddMPEG2Picture adds an MPEG2 Picture control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddMPEG2Picture(pic *ControlMPEG2Picture) error {
	if pic == nil {
		return fmt.Errorf("MPEG2 picture cannot be nil")
	}
	size := unsafe.Sizeof(*pic)
	data := unsafe.Slice((*byte)(unsafe.Pointer(pic)), size)
	ec.AddCompound(CtrlMPEG2Picture, data)
	return nil
}

// AddMPEG2Quantization adds an MPEG2 Quantization control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddMPEG2Quantization(quant *ControlMPEG2Quantization) error {
	if quant == nil {
		return fmt.Errorf("MPEG2 quantization cannot be nil")
	}
	size := unsafe.Sizeof(*quant)
	data := unsafe.Slice((*byte)(unsafe.Pointer(quant)), size)
	ec.AddCompound(CtrlMPEG2Quantization, data)
	return nil
}

// Type-safe helper methods for ExtControl (reading values back)

// GetMPEG2Sequence retrieves the MPEG2 Sequence from a control value.
func (ec *ExtControl) GetMPEG2Sequence() (*ControlMPEG2Sequence, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlMPEG2Sequence{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid MPEG2 sequence size: got %d, expected %d", len(data), expectedSize)
	}
	seq := &ControlMPEG2Sequence{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(seq)), expectedSize), data)
	return seq, nil
}

// GetMPEG2Picture retrieves the MPEG2 Picture from a control value.
func (ec *ExtControl) GetMPEG2Picture() (*ControlMPEG2Picture, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlMPEG2Picture{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid MPEG2 picture size: got %d, expected %d", len(data), expectedSize)
	}
	pic := &ControlMPEG2Picture{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(pic)), expectedSize), data)
	return pic, nil
}

// GetMPEG2Quantization retrieves the MPEG2 Quantization from a control value.
func (ec *ExtControl) GetMPEG2Quantization() (*ControlMPEG2Quantization, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlMPEG2Quantization{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid MPEG2 quantization size: got %d, expected %d", len(data), expectedSize)
	}
	quant := &ControlMPEG2Quantization{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(quant)), expectedSize), data)
	return quant, nil
}
