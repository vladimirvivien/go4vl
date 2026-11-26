package v4l2

const (
	H264NumDPBEntries uint32 = cV4L2_H264_NUM_DPB_ENTRIES
	H264RefListLength uint32 = cV4L2_H264_REF_LIST_LEN
)

// Stateless H.264 control payload types mapped directly to the Linux kernel UAPI
// structs in <linux/v4l2-controls.h>.
type (
	ControlH264SPS               = cV4L2CtrlH264SPS
	ControlH264PPS               = cV4L2CtrlH264PPS
	ControlH264ScalingMatrix     = cV4L2CtrlH264ScalingMatrix
	H264WeightFactors            = cV4L2H264WeightFactors
	// Deprecated: Use H264WeightFactors.
	H264WeightFators             = cV4L2H264WeightFactors
	ControlH264PredictionWeights = cV4L2CtrlH264PredWeights
	H264Reference                = cV4L2H264Reference
	ControlH264SliceParams       = cV4L2CtrlH264SliceParams
	H264DPBEntry                 = cV4L2H264DPBEntry
	ControlH264DecodeParams      = cV4L2CtrlH264DecodeParams
)

// H.264 Stateless Codec Control IDs
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlH264SPS           CtrlID = cV4L2_CID_STATELESS_H264_SPS
	CtrlH264PPS           CtrlID = cV4L2_CID_STATELESS_H264_PPS
	CtrlH264ScalingMatrix CtrlID = cV4L2_CID_STATELESS_H264_SCALING_MATRIX
	CtrlH264SliceParams   CtrlID = cV4L2_CID_STATELESS_H264_SLICE_PARAMS
	CtrlH264DecodeParams  CtrlID = cV4L2_CID_STATELESS_H264_DECODE_PARAMS
	CtrlH264PredWeights   CtrlID = cV4L2_CID_STATELESS_H264_PRED_WEIGHTS
	CtrlH264DecodeMode    CtrlID = cV4L2_CID_STATELESS_H264_DECODE_MODE
	CtrlH264StartCode     CtrlID = cV4L2_CID_STATELESS_H264_START_CODE
)

// AddH264SPS adds an H.264 Sequence Parameter Set control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264SPS(sps *ControlH264SPS) error {
	return addCompoundStruct(ec, CtrlH264SPS, "H.264 SPS", sps)
}

// AddH264PPS adds an H.264 Picture Parameter Set control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264PPS(pps *ControlH264PPS) error {
	return addCompoundStruct(ec, CtrlH264PPS, "H.264 PPS", pps)
}

// AddH264ScalingMatrix adds an H.264 Scaling Matrix control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264ScalingMatrix(matrix *ControlH264ScalingMatrix) error {
	return addCompoundStruct(ec, CtrlH264ScalingMatrix, "H.264 scaling matrix", matrix)
}

// AddH264SliceParams adds H.264 Slice Parameters control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264SliceParams(params *ControlH264SliceParams) error {
	return addCompoundStruct(ec, CtrlH264SliceParams, "H.264 slice params", params)
}

// AddH264DecodeParams adds H.264 Decode Parameters control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264DecodeParams(params *ControlH264DecodeParams) error {
	return addCompoundStruct(ec, CtrlH264DecodeParams, "H.264 decode params", params)
}

// AddH264PredWeights adds H.264 Prediction Weights control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddH264PredWeights(weights *ControlH264PredictionWeights) error {
	return addCompoundStruct(ec, CtrlH264PredWeights, "H.264 prediction weights", weights)
}

// GetH264SPS retrieves the H.264 SPS from a control value.
func (ec *ExtControl) GetH264SPS() (*ControlH264SPS, error) {
	return getCompoundStruct[ControlH264SPS](ec, "H.264 SPS")
}

// GetH264PPS retrieves the H.264 PPS from a control value.
func (ec *ExtControl) GetH264PPS() (*ControlH264PPS, error) {
	return getCompoundStruct[ControlH264PPS](ec, "H.264 PPS")
}

// GetH264ScalingMatrix retrieves the H.264 Scaling Matrix from a control value.
func (ec *ExtControl) GetH264ScalingMatrix() (*ControlH264ScalingMatrix, error) {
	return getCompoundStruct[ControlH264ScalingMatrix](ec, "H.264 scaling matrix")
}

// GetH264SliceParams retrieves the H.264 Slice Parameters from a control value.
func (ec *ExtControl) GetH264SliceParams() (*ControlH264SliceParams, error) {
	return getCompoundStruct[ControlH264SliceParams](ec, "H.264 slice params")
}

// GetH264DecodeParams retrieves the H.264 Decode Parameters from a control value.
func (ec *ExtControl) GetH264DecodeParams() (*ControlH264DecodeParams, error) {
	return getCompoundStruct[ControlH264DecodeParams](ec, "H.264 decode params")
}

// GetH264PredWeights retrieves the H.264 Prediction Weights from a control value.
func (ec *ExtControl) GetH264PredWeights() (*ControlH264PredictionWeights, error) {
	return getCompoundStruct[ControlH264PredictionWeights](ec, "H.264 prediction weights")
}
