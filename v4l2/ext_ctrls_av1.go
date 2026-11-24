package v4l2

// Stateless AV1 control payload types mapped directly to the Linux kernel UAPI
// structs in <linux/v4l2-controls.h>.
type (
	AV1TileInfo             = cV4L2AV1TileInfo
	AV1Quantization         = cV4L2AV1Quantization
	AV1Segmentation         = cV4L2AV1Segmentation
	AV1LoopFilter           = cV4L2AV1LoopFilter
	AV1CDEF                 = cV4L2AV1CDEF
	AV1LoopRestoration      = cV4L2AV1LoopRestoration
	AV1GlobalMotion         = cV4L2AV1GlobalMotion
	ControlAV1Sequence      = cV4L2CtrlAV1Sequence
	ControlAV1TileGroupEntry = cV4L2CtrlAV1TileGroupEntry
	ControlAV1Frame         = cV4L2CtrlAV1Frame
	ControlAV1FilmGrain     = cV4L2CtrlAV1FilmGrain
)

// Stateless AV1 control IDs (V4L2_CID_STATELESS_AV1_*).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec-stateless.html
const (
	CtrlAV1Sequence       CtrlID = cV4L2_CID_STATELESS_AV1_SEQUENCE
	CtrlAV1TileGroupEntry CtrlID = cV4L2_CID_STATELESS_AV1_TILE_GROUP_ENTRY
	CtrlAV1Frame          CtrlID = cV4L2_CID_STATELESS_AV1_FRAME
	CtrlAV1FilmGrain      CtrlID = cV4L2_CID_STATELESS_AV1_FILM_GRAIN
)

// AddAV1Sequence adds a stateless AV1 sequence header control.
func (ec *ExtControls) AddAV1Sequence(seq *ControlAV1Sequence) error {
	return addCompoundStruct(ec, CtrlAV1Sequence, "AV1 sequence", seq)
}

// AddAV1TileGroupEntry adds a stateless AV1 tile group entry control.
func (ec *ExtControls) AddAV1TileGroupEntry(entry *ControlAV1TileGroupEntry) error {
	return addCompoundStruct(ec, CtrlAV1TileGroupEntry, "AV1 tile group entry", entry)
}

// AddAV1Frame adds a stateless AV1 frame parameters control.
func (ec *ExtControls) AddAV1Frame(frame *ControlAV1Frame) error {
	return addCompoundStruct(ec, CtrlAV1Frame, "AV1 frame", frame)
}

// AddAV1FilmGrain adds a stateless AV1 film grain parameters control.
func (ec *ExtControls) AddAV1FilmGrain(grain *ControlAV1FilmGrain) error {
	return addCompoundStruct(ec, CtrlAV1FilmGrain, "AV1 film grain", grain)
}

// GetAV1Sequence retrieves a stateless AV1 sequence header from a compound control value.
func (ec *ExtControl) GetAV1Sequence() (*ControlAV1Sequence, error) {
	return getCompoundStruct[ControlAV1Sequence](ec, "AV1 sequence")
}

// GetAV1TileGroupEntry retrieves a stateless AV1 tile group entry from a compound control value.
func (ec *ExtControl) GetAV1TileGroupEntry() (*ControlAV1TileGroupEntry, error) {
	return getCompoundStruct[ControlAV1TileGroupEntry](ec, "AV1 tile group entry")
}

// GetAV1Frame retrieves stateless AV1 frame parameters from a compound control value.
func (ec *ExtControl) GetAV1Frame() (*ControlAV1Frame, error) {
	return getCompoundStruct[ControlAV1Frame](ec, "AV1 frame")
}

// GetAV1FilmGrain retrieves stateless AV1 film grain parameters from a compound control value.
func (ec *ExtControl) GetAV1FilmGrain() (*ControlAV1FilmGrain, error) {
	return getCompoundStruct[ControlAV1FilmGrain](ec, "AV1 film grain")
}
