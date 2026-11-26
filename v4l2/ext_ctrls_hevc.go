package v4l2

// Stateless HEVC control payload types mapped directly to the Linux kernel UAPI
// structs in <linux/v4l2-controls.h>.
type (
	HEVCPredWeightTable        = cV4L2HEVCPredWeightTable
	HEVCDPBEntry               = cV4L2HEVCDPBEntry
	ControlHEVCSPS             = cV4L2CtrlHEVCSPS
	ControlHEVCPPS             = cV4L2CtrlHEVCPPS
	ControlHEVCSliceParams     = cV4L2CtrlHEVCSliceParams
	ControlHEVCScalingMatrix   = cV4L2CtrlHEVCScalingMatrix
	ControlHEVCDecodeParams    = cV4L2CtrlHEVCDecodeParams
)

// Stateless HEVC control IDs (V4L2_CID_STATELESS_HEVC_*).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
const (
	CtrlHEVCSPS               CtrlID = cV4L2_CID_STATELESS_HEVC_SPS
	CtrlHEVCPPS               CtrlID = cV4L2_CID_STATELESS_HEVC_PPS
	CtrlHEVCSliceParams       CtrlID = cV4L2_CID_STATELESS_HEVC_SLICE_PARAMS
	CtrlHEVCScalingMatrix     CtrlID = cV4L2_CID_STATELESS_HEVC_SCALING_MATRIX
	CtrlHEVCDecodeParams      CtrlID = cV4L2_CID_STATELESS_HEVC_DECODE_PARAMS
	CtrlHEVCDecodeMode        CtrlID = cV4L2_CID_STATELESS_HEVC_DECODE_MODE
	CtrlHEVCStartCode         CtrlID = cV4L2_CID_STATELESS_HEVC_START_CODE
	CtrlHEVCEntryPointOffsets CtrlID = cV4L2_CID_STATELESS_HEVC_ENTRY_POINT_OFFSETS
)

// AddHEVCSPS adds a stateless HEVC SPS control.
func (ec *ExtControls) AddHEVCSPS(sps *ControlHEVCSPS) error {
	return addCompoundStruct(ec, CtrlHEVCSPS, "HEVC SPS", sps)
}

// AddHEVCPPS adds a stateless HEVC PPS control.
func (ec *ExtControls) AddHEVCPPS(pps *ControlHEVCPPS) error {
	return addCompoundStruct(ec, CtrlHEVCPPS, "HEVC PPS", pps)
}

// AddHEVCSliceParams adds a stateless HEVC slice parameters control.
func (ec *ExtControls) AddHEVCSliceParams(params *ControlHEVCSliceParams) error {
	return addCompoundStruct(ec, CtrlHEVCSliceParams, "HEVC slice params", params)
}

// AddHEVCScalingMatrix adds a stateless HEVC scaling matrix control.
func (ec *ExtControls) AddHEVCScalingMatrix(matrix *ControlHEVCScalingMatrix) error {
	return addCompoundStruct(ec, CtrlHEVCScalingMatrix, "HEVC scaling matrix", matrix)
}

// AddHEVCDecodeParams adds a stateless HEVC decode parameters control.
func (ec *ExtControls) AddHEVCDecodeParams(params *ControlHEVCDecodeParams) error {
	return addCompoundStruct(ec, CtrlHEVCDecodeParams, "HEVC decode params", params)
}

// GetHEVCSPS retrieves a stateless HEVC SPS from a compound control value.
func (ec *ExtControl) GetHEVCSPS() (*ControlHEVCSPS, error) {
	return getCompoundStruct[ControlHEVCSPS](ec, "HEVC SPS")
}

// GetHEVCPPS retrieves a stateless HEVC PPS from a compound control value.
func (ec *ExtControl) GetHEVCPPS() (*ControlHEVCPPS, error) {
	return getCompoundStruct[ControlHEVCPPS](ec, "HEVC PPS")
}

// GetHEVCSliceParams retrieves stateless HEVC slice parameters from a compound control value.
func (ec *ExtControl) GetHEVCSliceParams() (*ControlHEVCSliceParams, error) {
	return getCompoundStruct[ControlHEVCSliceParams](ec, "HEVC slice params")
}

// GetHEVCScalingMatrix retrieves a stateless HEVC scaling matrix from a compound control value.
func (ec *ExtControl) GetHEVCScalingMatrix() (*ControlHEVCScalingMatrix, error) {
	return getCompoundStruct[ControlHEVCScalingMatrix](ec, "HEVC scaling matrix")
}

// GetHEVCDecodeParams retrieves stateless HEVC decode parameters from a compound control value.
func (ec *ExtControl) GetHEVCDecodeParams() (*ControlHEVCDecodeParams, error) {
	return getCompoundStruct[ControlHEVCDecodeParams](ec, "HEVC decode params")
}
