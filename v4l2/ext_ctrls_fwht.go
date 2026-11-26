package v4l2

// ControlFWHTParams maps directly to Linux kernel struct v4l2_ctrl_fwht_params.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
type ControlFWHTParams = cV4L2CtrlFWHTParams

// FWHT Stateless Codec Control IDs
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlFWHTParams CtrlID = cV4L2_CID_STATELESS_FWHT_PARAMS
)

// AddFWHTParams adds an FWHT Parameters control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddFWHTParams(params *ControlFWHTParams) error {
	return addCompoundStruct(ec, CtrlFWHTParams, "FWHT params", params)
}

// GetFWHTParams retrieves the FWHT Parameters from a control value.
func (ec *ExtControl) GetFWHTParams() (*ControlFWHTParams, error) {
	return getCompoundStruct[ControlFWHTParams](ec, "FWHT params")
}
