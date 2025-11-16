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

const (
	VP8CoefficientProbabilityCount uint32 = 11 // C.V4L2_VP8_COEFF_PROB_CNT
	VP8MVProbabilityCount          uint32 = 19 // C.V4L2_VP8_MV_PROB_CNT
)

// VP8Segment (v4l2_vp8_segment)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1692
type VP8Segment struct {
	QuantUpdate          [4]int8
	LoopFilterUpdate     [4]int8
	SegmentProbabilities [3]uint8
	_                    uint8 // padding
	Flags                uint32
}

// VP8LoopFilter (v4l2_vp8_loop_filter)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1719
type VP8LoopFilter struct {
	ReferenceFrameDelta int8
	MBModeDelta         int8
	SharpnessLevel      uint8
	Level               uint8
	_                   uint16 // padding
	Flags               uint32
}

// VP8Quantization (v4l2_vp8_quantization)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1744
type VP8Quantization struct {
	YACQIndex uint8
	YDCDelta  int8
	Y2DCDelta int8
	Y2ACDelta int8
	UVDCDelta int8
	UVACDelta int8
	_         uint16
}

// VP8Entropy (v4l2_vp8_entropy)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1771
type VP8Entropy struct {
	CoefficientProbabilities [4][8][3][VP8CoefficientProbabilityCount]uint8
	YModeProbabilities       uint8
	UVModeProbabilities      uint8
	MVProbabilities          [2][VP8MVProbabilityCount]uint8
	_                        [3]uint8 // padding
}

// VP8EntropyCoderState (v4l2_vp8_entropy_coder_state)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1790
type VP8EntropyCoderState struct {
	Range    uint8
	Value    uint8
	BitCount uint8
	_        uint8 // padding
}

// ControlVP8Frame (v4l2_ctrl_vp8_frame)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1836
type ControlVP8Frame struct {
	Segment           VP8Segment
	LoopFilter        VP8LoopFilter
	Quantization      VP8Quantization
	Entropy           VP8Entropy
	EntropyCoderState VP8EntropyCoderState

	Width  uint16
	Height uint16

	HorizontalScale uint8
	VerticalScale   uint8

	Version       uint8
	ProbSkipFalse uint8
	PropIntra     uint8
	PropLast      uint8
	ProbGF        uint8
	NumDCTParts   uint8

	FirstPartSize   uint32
	FirstPartHeader uint32
	DCTPartSize     uint32

	LastFrameTimestamp   uint64
	GoldenFrameTimestamp uint64
	AltFrameTimestamp    uint64

	Flags uint64
}

// VP8 Stateless Codec Control IDs
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlVP8Frame CtrlID = C.V4L2_CID_STATELESS_VP8_FRAME
)

// Type-safe helper methods for ExtControls

// AddVP8Frame adds a VP8 Frame control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddVP8Frame(frame *ControlVP8Frame) error {
	if frame == nil {
		return fmt.Errorf("VP8 frame cannot be nil")
	}
	size := unsafe.Sizeof(*frame)
	data := unsafe.Slice((*byte)(unsafe.Pointer(frame)), size)
	ec.AddCompound(CtrlVP8Frame, data)
	return nil
}

// Type-safe helper methods for ExtControl (reading values back)

// GetVP8Frame retrieves the VP8 Frame from a control value.
func (ec *ExtControl) GetVP8Frame() (*ControlVP8Frame, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlVP8Frame{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid VP8 frame size: got %d, expected %d", len(data), expectedSize)
	}
	frame := &ControlVP8Frame{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(frame)), expectedSize), data)
	return frame, nil
}
