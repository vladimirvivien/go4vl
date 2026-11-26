package v4l2

const (
	VP8CoefficientProbabilityCount uint32 = cV4L2_VP8_COEFF_PROB_CNT
	VP8MVProbabilityCount          uint32 = cV4L2_VP8_MV_PROB_CNT
)

// Stateless VP8 control payload types mapped directly to the Linux kernel UAPI
// structs in <linux/v4l2-controls.h>.
type (
	VP8Segment           = cV4L2VP8Segment
	VP8LoopFilter        = cV4L2VP8LoopFilter
	VP8Quantization      = cV4L2VP8Quantization
	VP8Entropy           = cV4L2VP8Entropy
	VP8EntropyCoderState = cV4L2VP8EntropyCoderState
	ControlVP8Frame      = cV4L2CtrlVP8Frame
)

// VP8 Stateless Codec Control IDs
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h
const (
	CtrlVP8Frame CtrlID = cV4L2_CID_STATELESS_VP8_FRAME
)

// AddVP8Frame adds a VP8 Frame control.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
func (ec *ExtControls) AddVP8Frame(frame *ControlVP8Frame) error {
	return addCompoundStruct(ec, CtrlVP8Frame, "VP8 frame", frame)
}

// GetVP8Frame retrieves the VP8 Frame from a control value.
func (ec *ExtControl) GetVP8Frame() (*ControlVP8Frame, error) {
	return getCompoundStruct[ControlVP8Frame](ec, "VP8 frame")
}
