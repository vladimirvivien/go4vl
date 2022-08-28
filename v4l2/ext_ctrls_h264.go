package v4l2

import "C"

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
