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

// TODO - Need to figure out how to import the proper header files for H264 support
const (
	H264NumDPBEntries uint32 = 16 // C.V4L2_H264_NUM_DPB_ENTRIES
	H264RefListLength uint32 = 32 // C.V4L2_H264_REF_LIST_LEN
)

// ControlH264SPS (v4l2_ctrl_h264_sps)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1308
type ControlH264SPS struct {
	ProfileIDC                     uint8
	ConstraintSetFlags             uint8
	LevelIDC                       uint8
	SequenceParameterSetID         uint8
	ChromaFormatIDC                uint8
	BitDepthLumaMinus8             uint8
	BitDepthChromaMinus8           uint8
	Log2MaxFrameNumMinus4          uint8
	PicOrderCntType                uint8
	Log2MaxPicOrderCntLsbMinus4    uint8
	MaxNumRefFrames                uint8
	NumRefFramesInPicOrderCntCycle uint8
	OffsetForRefFrame              [255]int32
	OffsetForNonRefPic             int32
	OffsetForTopToBottomField      int32
	PicWidthInMbsMinus1            uint16
	PicHeightInMapUnitsMinus1      uint16
	Falgs                          uint32
}

// ControlH264PPS (v4l2_ctrl_h264_pps)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1364
type ControlH264PPS struct {
	PicParameterSetID                uint8
	SeqParameterSetID                uint8
	NumSliceGroupsMinus1             uint8
	NumRefIndexL0DefaultActiveMinus1 uint8
	NumRefIndexL1DefaultActiveMinus1 uint8
	WeightedBipredIDC                uint8
	PicInitQPMinus26                 int8
	PicInitQSMinus26                 int8
	ChromaQPIndexOffset              int8
	SecondChromaQPIndexOffset        int8
	Flags                            uint16
}

// ControlH264ScalingMatrix (v4l2_ctrl_h264_scaling_matrix)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1396
type ControlH264ScalingMatrix struct {
	ScalingList4x4 [6][16]uint8
	ScalingList8x8 [6][64]uint8
}

type H264WeightFators struct {
	LumaWeight   [32]int16
	LumaOffset   [32]int16
	ChromaWeight [32][2]int16
	ChromaOffset [32][2]int16
}

// ControlH264PredictionWeights (4l2_ctrl_h264_pred_weights)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1426
type ControlH264PredictionWeights struct {
	LumaLog2WeightDenom   uint16
	ChromaLog2WeightDenom uint16
	WeightFactors         [2]H264WeightFators
}

// H264Reference (v4l2_h264_reference)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1452
type H264Reference struct {
	Fields uint8
	Index  uint8
}

// ControlH264SliceParams (v4l2_ctrl_h264_slice_params)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1499
type ControlH264SliceParams struct {
	HeaderBitSize              uint32
	FirstMBInSlice             uint32
	SliceType                  uint8
	ColorPlaneID               uint8
	RedundantPicCnt            uint8
	CabacInitIDC               uint8
	SliceQPDelta               int8
	SliceQSDelta               int8
	DisableDeblockingFilterIDC uint8
	SliceAlphaC0OffsetDiv2     int8
	SliceBetaOffsetDiv2        int8
	NumRefIdxL0ActiveMinus1    uint8
	NumRefIdxL1ActiveMinus1    uint8

	_ uint8 // reserved for padding

	RefPicList0 [H264RefListLength]H264Reference
	RefPicList1 [H264RefListLength]H264Reference

	Flags uint32
}

// H264DPBEntry (v4l2_h264_dpb_entry)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1544
type H264DPBEntry struct {
	ReferenceTS         uint64
	PicNum              uint32
	FrameNum            uint16
	Fields              uint8
	_                   [8]uint8 // reserved (padding field)
	TopFieldOrder       int32
	BottomFieldOrderCnt int32
	Flags               uint32
}

// ControlH264DecodeParams (v4l2_ctrl_h264_decode_params)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1581
type ControlH264DecodeParams struct {
	DPB                     [H264NumDPBEntries]H264DPBEntry
	NalRefIDC               uint16
	FrameNum                uint16
	TopFieldOrderCnt        int32
	BottomFieldOrderCnt     int32
	IDRPicID                uint16
	PicOrderCntLSB          uint16
	DeltaPicOrderCntBottom  int32
	DeltaPicOrderCnt0       int32
	DeltaPicOrderCnt1       int32
	DecRefPicMarkingBitSize uint32
	PicOrderCntBitSize      uint32
	SliceGroupChangeCycle   uint32
	_                       uint32 // reserved (padding)
	Flags                   uint32
}

// H.264 Stateless Codec Control IDs
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlH264SPS            CtrlID = C.V4L2_CID_STATELESS_H264_SPS
	CtrlH264PPS            CtrlID = C.V4L2_CID_STATELESS_H264_PPS
	CtrlH264ScalingMatrix  CtrlID = C.V4L2_CID_STATELESS_H264_SCALING_MATRIX
	CtrlH264SliceParams    CtrlID = C.V4L2_CID_STATELESS_H264_SLICE_PARAMS
	CtrlH264DecodeParams   CtrlID = C.V4L2_CID_STATELESS_H264_DECODE_PARAMS
	CtrlH264PredWeights    CtrlID = C.V4L2_CID_STATELESS_H264_PRED_WEIGHTS
	CtrlH264DecodeMode     CtrlID = C.V4L2_CID_STATELESS_H264_DECODE_MODE
	CtrlH264StartCode      CtrlID = C.V4L2_CID_STATELESS_H264_START_CODE
)

// Type-safe helper methods for ExtControls

// AddH264SPS adds an H.264 Sequence Parameter Set control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264SPS(sps *ControlH264SPS) error {
	if sps == nil {
		return fmt.Errorf("H.264 SPS cannot be nil")
	}
	// Convert Go struct to byte slice using unsafe
	size := unsafe.Sizeof(*sps)
	data := unsafe.Slice((*byte)(unsafe.Pointer(sps)), size)
	ec.AddCompound(CtrlH264SPS, data)
	return nil
}

// AddH264PPS adds an H.264 Picture Parameter Set control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264PPS(pps *ControlH264PPS) error {
	if pps == nil {
		return fmt.Errorf("H.264 PPS cannot be nil")
	}
	size := unsafe.Sizeof(*pps)
	data := unsafe.Slice((*byte)(unsafe.Pointer(pps)), size)
	ec.AddCompound(CtrlH264PPS, data)
	return nil
}

// AddH264ScalingMatrix adds an H.264 Scaling Matrix control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264ScalingMatrix(matrix *ControlH264ScalingMatrix) error {
	if matrix == nil {
		return fmt.Errorf("H.264 scaling matrix cannot be nil")
	}
	size := unsafe.Sizeof(*matrix)
	data := unsafe.Slice((*byte)(unsafe.Pointer(matrix)), size)
	ec.AddCompound(CtrlH264ScalingMatrix, data)
	return nil
}

// AddH264SliceParams adds H.264 Slice Parameters control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264SliceParams(params *ControlH264SliceParams) error {
	if params == nil {
		return fmt.Errorf("H.264 slice params cannot be nil")
	}
	size := unsafe.Sizeof(*params)
	data := unsafe.Slice((*byte)(unsafe.Pointer(params)), size)
	ec.AddCompound(CtrlH264SliceParams, data)
	return nil
}

// AddH264DecodeParams adds H.264 Decode Parameters control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264DecodeParams(params *ControlH264DecodeParams) error {
	if params == nil {
		return fmt.Errorf("H.264 decode params cannot be nil")
	}
	size := unsafe.Sizeof(*params)
	data := unsafe.Slice((*byte)(unsafe.Pointer(params)), size)
	ec.AddCompound(CtrlH264DecodeParams, data)
	return nil
}

// AddH264PredWeights adds H.264 Prediction Weights control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264PredWeights(weights *ControlH264PredictionWeights) error {
	if weights == nil {
		return fmt.Errorf("H.264 prediction weights cannot be nil")
	}
	size := unsafe.Sizeof(*weights)
	data := unsafe.Slice((*byte)(unsafe.Pointer(weights)), size)
	ec.AddCompound(CtrlH264PredWeights, data)
	return nil
}

// Type-safe helper methods for ExtControl (reading values back)

// GetH264SPS retrieves the H.264 SPS from a control value.
func (ec *ExtControl) GetH264SPS() (*ControlH264SPS, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlH264SPS{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid H.264 SPS size: got %d, expected %d", len(data), expectedSize)
	}
	// Create a copy to avoid issues with data lifetime
	sps := &ControlH264SPS{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(sps)), expectedSize), data)
	return sps, nil
}

// GetH264PPS retrieves the H.264 PPS from a control value.
func (ec *ExtControl) GetH264PPS() (*ControlH264PPS, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlH264PPS{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid H.264 PPS size: got %d, expected %d", len(data), expectedSize)
	}
	pps := &ControlH264PPS{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(pps)), expectedSize), data)
	return pps, nil
}

// GetH264ScalingMatrix retrieves the H.264 Scaling Matrix from a control value.
func (ec *ExtControl) GetH264ScalingMatrix() (*ControlH264ScalingMatrix, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlH264ScalingMatrix{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid H.264 scaling matrix size: got %d, expected %d", len(data), expectedSize)
	}
	matrix := &ControlH264ScalingMatrix{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(matrix)), expectedSize), data)
	return matrix, nil
}

// GetH264SliceParams retrieves the H.264 Slice Parameters from a control value.
func (ec *ExtControl) GetH264SliceParams() (*ControlH264SliceParams, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlH264SliceParams{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid H.264 slice params size: got %d, expected %d", len(data), expectedSize)
	}
	params := &ControlH264SliceParams{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(params)), expectedSize), data)
	return params, nil
}

// GetH264DecodeParams retrieves the H.264 Decode Parameters from a control value.
func (ec *ExtControl) GetH264DecodeParams() (*ControlH264DecodeParams, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlH264DecodeParams{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid H.264 decode params size: got %d, expected %d", len(data), expectedSize)
	}
	params := &ControlH264DecodeParams{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(params)), expectedSize), data)
	return params, nil
}

// GetH264PredWeights retrieves the H.264 Prediction Weights from a control value.
func (ec *ExtControl) GetH264PredWeights() (*ControlH264PredictionWeights, error) {
	data := ec.GetCompoundData()
	if data == nil {
		return nil, fmt.Errorf("no compound data in control")
	}
	expectedSize := int(unsafe.Sizeof(ControlH264PredictionWeights{}))
	if len(data) < expectedSize {
		return nil, fmt.Errorf("invalid H.264 prediction weights size: got %d, expected %d", len(data), expectedSize)
	}
	weights := &ControlH264PredictionWeights{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(weights)), expectedSize), data)
	return weights, nil
}
