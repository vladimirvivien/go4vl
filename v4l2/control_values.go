package v4l2

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../include/
#include <linux/videodev2.h>
#include <linux/v4l2-controls.h>
*/
import "C"

// ControlClass
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L56
type ControlClass = uint32

const (
	UserControlClass            ControlClass = C.V4L2_CTRL_CLASS_USER
	CodecControlClass           ControlClass = C.V4L2_CTRL_CLASS_CODEC
	CameraControlClass          ControlClass = C.V4L2_CTRL_CLASS_CAMERA
	CameraFlashControlClass     ControlClass = C.V4L2_CTRL_CLASS_FLASH
	JPEGControlClass            ControlClass = C.V4L2_CTRL_CLASS_JPEG
	ImageSourceControlClass     ControlClass = C.V4L2_CTRL_CLASS_IMAGE_SOURCE
	ImageProcessingControlClass ControlClass = C.V4L2_CTRL_CLASS_IMAGE_PROC
	DigitalVideoControlClass    ControlClass = C.V4L2_CTRL_CLASS_DV
	DetectionControlClass       ControlClass = C.V4L2_CTRL_CLASS_DETECT
	CodecStatelessControlClass  ControlClass = C.V4L2_CTRL_CLASS_CODEC_STATELESS
	ColorimitryControlClass     ControlClass = C.V4L2_CTRL_CLASS_COLORIMETRY
)

// CtrlID type for control values
type CtrlID = uint32

// Powerline frequency control enums
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L100
const (
	PowerlineFrequencyDisabled uint32 = iota
	PowerlineFrequency50Hz
	PowerlineFrequency60Hz
	PowerlineFrequencyAuto
)

// Color FX control enums
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L114
const (
	ColorFXNone uint32 = iota
	ColorFXBlackWhite
	ColorFXSepia
	ColorFXNegative
	ColorFXEmboss
	ColorFXSketch
	ColorFXSkyBlue
	ColorFXGrassGreen
	ColorFXSkinWhiten
	ColorFXVivid
	ColorFXAqua
	ColorFXArtFreeze
	ColorFXSilhouette
	ColorFXSolarization
	ColorFXAntique
	ColorFXSetCBCR
	ColorFXSetRGB
)

// User Controls
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L74
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html#control-id
const (
	CtrlBrightness              CtrlID = C.V4L2_CID_BRIGHTNESS
	CtrlContrast                CtrlID = C.V4L2_CID_CONTRAST
	CtrlSaturation              CtrlID = C.V4L2_CID_SATURATION
	CtrlHue                     CtrlID = C.V4L2_CID_HUE
	CtrlAutoWhiteBalance        CtrlID = C.V4L2_CID_AUTO_WHITE_BALANCE
	CtrlDoWhiteBalance          CtrlID = C.V4L2_CID_DO_WHITE_BALANCE
	CtrlRedBalance              CtrlID = C.V4L2_CID_RED_BALANCE
	CtrlBlueBalance             CtrlID = C.V4L2_CID_BLUE_BALANCE
	CtrlGamma                   CtrlID = C.V4L2_CID_GAMMA
	CtrlExposure                CtrlID = C.V4L2_CID_EXPOSURE
	CtrlAutogain                CtrlID = C.V4L2_CID_AUTOGAIN
	CtrlGain                    CtrlID = C.V4L2_CID_GAIN
	CtrlHFlip                   CtrlID = C.V4L2_CID_HFLIP
	CtrlVFlip                   CtrlID = C.V4L2_CID_VFLIP
	CtrlPowerlineFrequency      CtrlID = C.V4L2_CID_POWER_LINE_FREQUENCY
	CtrlHueAuto                 CtrlID = C.V4L2_CID_HUE_AUTO
	CtrlWhiteBalanceTemperature CtrlID = C.V4L2_CID_WHITE_BALANCE_TEMPERATURE
	CtrlSharpness               CtrlID = C.V4L2_CID_SHARPNESS
	CtrlBacklightCompensation   CtrlID = C.V4L2_CID_BACKLIGHT_COMPENSATION
	CtrlChromaAutomaticGain     CtrlID = C.V4L2_CID_CHROMA_AGC
	CtrlColorKiller             CtrlID = C.V4L2_CID_COLOR_KILLER
	CtrlColorFX                 CtrlID = C.V4L2_CID_COLORFX
	CtrlColorFXCBCR             CtrlID = C.V4L2_CID_COLORFX_CBCR
	CtrlColorFXRGB              CtrlID = C.V4L2_CID_COLORFX_RGB
	CtrlAutoBrightness          CtrlID = C.V4L2_CID_AUTOBRIGHTNESS
	CtrlRotate                  CtrlID = C.V4L2_CID_ROTATE
	CtrlBackgroundColor         CtrlID = C.V4L2_CID_BG_COLOR
	CtrlMinimumCaptureBuffers   CtrlID = C.V4L2_CID_MIN_BUFFERS_FOR_CAPTURE
	CtrlMinimumOutputBuffers    CtrlID = C.V4L2_CID_MIN_BUFFERS_FOR_OUTPUT
	CtrlAlphaComponent          CtrlID = C.V4L2_CID_ALPHA_COMPONENT
)

// ControlType stores values fround in enum v4l2_ctrl_type
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1799
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-queryctrl.html?highlight=v4l2_ctrl_type#c.V4L.v4l2_ctrl_type
type CtrlType uint32

// CtrlType constants
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1799
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-queryctrl.html?highlight=v4l2_ctrl_type#c.V4L.v4l2_ctrl_type
const (
	CtrlTypeInt                 CtrlType = C.V4L2_CTRL_TYPE_INTEGER
	CtrlTypeBool                CtrlType = C.V4L2_CTRL_TYPE_BOOLEAN
	CtrlTypeMenu                CtrlType = C.V4L2_CTRL_TYPE_MENU
	CtrlTypeButton              CtrlType = C.V4L2_CTRL_TYPE_BUTTON
	CtrlTypeInt64               CtrlType = C.V4L2_CTRL_TYPE_INTEGER64
	CtrlTypeClass               CtrlType = C.V4L2_CTRL_TYPE_CTRL_CLASS
	CtrlTypeString              CtrlType = C.V4L2_CTRL_TYPE_STRING
	CtrlTypeBitMask             CtrlType = C.V4L2_CTRL_TYPE_BITMASK
	CtrlTypeIntegerMenu         CtrlType = C.V4L2_CTRL_TYPE_INTEGER_MENU
	CtrlTypeCompoundTypes       CtrlType = C.V4L2_CTRL_COMPOUND_TYPES
	CtrlTypeU8                  CtrlType = C.V4L2_CTRL_TYPE_U8
	CtrlTypeU16                 CtrlType = C.V4L2_CTRL_TYPE_U16
	CtrlTypeU32                 CtrlType = C.V4L2_CTRL_TYPE_U32
	CtrlTypeArear               CtrlType = C.V4L2_CTRL_TYPE_AREA
	CtrlTypeHDR10CLLInfo        CtrlType = C.V4L2_CTRL_TYPE_HDR10_CLL_INFO
	CtrlTypeHDRMasteringDisplay CtrlType = C.V4L2_CTRL_TYPE_HDR10_MASTERING_DISPLAY
	CtrlTypeH264SPS             CtrlType = C.V4L2_CTRL_TYPE_H264_SPS
	CtrlTypeH264PPS             CtrlType = C.V4L2_CTRL_TYPE_H264_PPS
	CtrlTypeH264ScalingMatrix   CtrlType = C.V4L2_CTRL_TYPE_H264_SCALING_MATRIX
	CtrlTypeH264SliceParams     CtrlType = C.V4L2_CTRL_TYPE_H264_SLICE_PARAMS
	CtrlTypeH264DecodeParams    CtrlType = C.V4L2_CTRL_TYPE_H264_DECODE_PARAMS
	CtrlTypeFWHTParams          CtrlType = C.V4L2_CTRL_TYPE_FWHT_PARAMS
	CtrlTypeVP8Frame            CtrlType = C.V4L2_CTRL_TYPE_VP8_FRAME
	CtrlTypeMPEG2Quantization   CtrlType = C.V4L2_CTRL_TYPE_MPEG2_QUANTISATION
	CtrlTypeMPEG2Sequence       CtrlType = C.V4L2_CTRL_TYPE_MPEG2_SEQUENCE
	CtrlTypeMPEG2Picture        CtrlType = C.V4L2_CTRL_TYPE_MPEG2_PICTURE
	CtrlTypeVP9CompressedHDR    CtrlType = C.V4L2_CTRL_TYPE_VP9_COMPRESSED_HDR
	CtrlTypeVP9Frame            CtrlType = C.V4L2_CTRL_TYPE_VP9_FRAME
)
