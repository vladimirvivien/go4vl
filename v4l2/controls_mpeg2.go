package v4l2

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
