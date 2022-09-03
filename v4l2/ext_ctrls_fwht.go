package v4l2

// ControlFWHTParams (v4l2_ctrl_fwht_params)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1659
type ControlFWHTParams struct {
	BackwardRefTimestamp uint64
	Version              uint32
	Width                uint32
	Height               uint32
	Flags                uint32
	// Colorspace           ColorspaceType
	// XFerFunc             XferFunctionType
	// YCbCrEncoding        YCbCrEncodingType
	// Quantization         QuantizationType
}
