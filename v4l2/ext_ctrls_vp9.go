package v4l2

// Stateless VP9 control payload types mapped directly to the Linux kernel UAPI
// structs in <linux/v4l2-controls.h>.
type (
	VP9LoopFilter           = cV4L2VP9LoopFilter
	VP9Quantization         = cV4L2VP9Quantization
	VP9Segmentation         = cV4L2VP9Segmentation
	VP9MVProbs              = cV4L2VP9MVProbs
	ControlVP9Frame         = cV4L2CtrlVP9Frame
	ControlVP9CompressedHdr = cV4L2CtrlVP9CompressedHdr
)

// Stateless VP9 control IDs (V4L2_CID_STATELESS_VP9_*).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
const (
	CtrlVP9Frame         CtrlID = cV4L2_CID_STATELESS_VP9_FRAME
	CtrlVP9CompressedHdr CtrlID = cV4L2_CID_STATELESS_VP9_COMPRESSED_HDR
)

// AddVP9Frame adds a stateless VP9 frame parameters control.
func (ec *ExtControls) AddVP9Frame(frame *ControlVP9Frame) error {
	return addCompoundStruct(ec, CtrlVP9Frame, "VP9 frame", frame)
}

// AddVP9CompressedHdr adds a stateless VP9 compressed header probabilities control.
func (ec *ExtControls) AddVP9CompressedHdr(hdr *ControlVP9CompressedHdr) error {
	return addCompoundStruct(ec, CtrlVP9CompressedHdr, "VP9 compressed header", hdr)
}

// GetVP9Frame retrieves stateless VP9 frame parameters from a compound control value.
func (ec *ExtControl) GetVP9Frame() (*ControlVP9Frame, error) {
	return getCompoundStruct[ControlVP9Frame](ec, "VP9 frame")
}

// GetVP9CompressedHdr retrieves stateless VP9 compressed header probabilities from a compound control value.
func (ec *ExtControl) GetVP9CompressedHdr() (*ControlVP9CompressedHdr, error) {
	return getCompoundStruct[ControlVP9CompressedHdr](ec, "VP9 compressed header")
}
