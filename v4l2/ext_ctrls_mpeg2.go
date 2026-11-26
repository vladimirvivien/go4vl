package v4l2

// Stateless MPEG2 control payload types mapped directly to the Linux kernel UAPI
// structs in <linux/v4l2-controls.h>.
type (
	ControlMPEG2Sequence     = cV4L2CtrlMPEG2Sequence
	ControlMPEG2Picture      = cV4L2CtrlMPEG2Picture
	ControlMPEG2Quantization = cV4L2CtrlMPEG2Quantisation
)

// MPEG2 Stateless Codec Control IDs
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlMPEG2Sequence     CtrlID = cV4L2_CID_STATELESS_MPEG2_SEQUENCE
	CtrlMPEG2Picture      CtrlID = cV4L2_CID_STATELESS_MPEG2_PICTURE
	CtrlMPEG2Quantization CtrlID = cV4L2_CID_STATELESS_MPEG2_QUANTISATION
)

// AddMPEG2Sequence adds an MPEG2 Sequence control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddMPEG2Sequence(seq *ControlMPEG2Sequence) error {
	return addCompoundStruct(ec, CtrlMPEG2Sequence, "MPEG2 sequence", seq)
}

// AddMPEG2Picture adds an MPEG2 Picture control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddMPEG2Picture(pic *ControlMPEG2Picture) error {
	return addCompoundStruct(ec, CtrlMPEG2Picture, "MPEG2 picture", pic)
}

// AddMPEG2Quantization adds an MPEG2 Quantization control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddMPEG2Quantization(quant *ControlMPEG2Quantization) error {
	return addCompoundStruct(ec, CtrlMPEG2Quantization, "MPEG2 quantization", quant)
}

// GetMPEG2Sequence retrieves the MPEG2 Sequence from a control value.
func (ec *ExtControl) GetMPEG2Sequence() (*ControlMPEG2Sequence, error) {
	return getCompoundStruct[ControlMPEG2Sequence](ec, "MPEG2 sequence")
}

// GetMPEG2Picture retrieves the MPEG2 Picture from a control value.
func (ec *ExtControl) GetMPEG2Picture() (*ControlMPEG2Picture, error) {
	return getCompoundStruct[ControlMPEG2Picture](ec, "MPEG2 picture")
}

// GetMPEG2Quantization retrieves the MPEG2 Quantization from a control value.
func (ec *ExtControl) GetMPEG2Quantization() (*ControlMPEG2Quantization, error) {
	return getCompoundStruct[ControlMPEG2Quantization](ec, "MPEG2 quantization")
}
