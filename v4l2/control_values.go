package v4l2

// CtrlClass
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L56
type CtrlClass = uint32

const (
	CtrlClassUser            CtrlClass = cV4L2_CTRL_CLASS_USER
	CtrlClassCodec           CtrlClass = cV4L2_CTRL_CLASS_CODEC
	CtrlClassCamera          CtrlClass = cV4L2_CTRL_CLASS_CAMERA
	CtrlClassFlash           CtrlClass = cV4L2_CTRL_CLASS_FLASH
	CtrlClassJPEG            CtrlClass = cV4L2_CTRL_CLASS_JPEG
	CtrlClassImageSource     CtrlClass = cV4L2_CTRL_CLASS_IMAGE_SOURCE
	CtrlClassImageProcessing CtrlClass = cV4L2_CTRL_CLASS_IMAGE_PROC
	CtrlClassDigitalVideo    CtrlClass = cV4L2_CTRL_CLASS_DV
	CtrlClassDetection       CtrlClass = cV4L2_CTRL_CLASS_DETECT
	CtrlClassCodecStateless  CtrlClass = cV4L2_CTRL_CLASS_CODEC_STATELESS
	CtrlClassColorimitry     CtrlClass = cV4L2_CTRL_CLASS_COLORIMETRY
)

var (
	// CtrlClasses is a slice of all Control classes
	CtrlClasses = []CtrlClass{
		CtrlClassUser,
		CtrlClassCodec,
		CtrlClassCamera,
		CtrlClassFlash,
		CtrlClassJPEG,
		CtrlClassImageSource,
		CtrlClassDigitalVideo,
		CtrlClassDetection,
		CtrlClassCodecStateless,
		CtrlClassColorimitry,
	}
)

// CtrlType constants
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1799
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-queryctrl.html?highlight=v4l2_ctrl_type#c.V4L.v4l2_ctrl_type
type CtrlType uint32

const (
	CtrlTypeInt                 CtrlType = cV4L2_CTRL_TYPE_INTEGER
	CtrlTypeBool                CtrlType = cV4L2_CTRL_TYPE_BOOLEAN
	CtrlTypeMenu                CtrlType = cV4L2_CTRL_TYPE_MENU
	CtrlTypeButton              CtrlType = cV4L2_CTRL_TYPE_BUTTON
	CtrlTypeInt64               CtrlType = cV4L2_CTRL_TYPE_INTEGER64
	CtrlTypeClass               CtrlType = cV4L2_CTRL_TYPE_CTRL_CLASS
	CtrlTypeString              CtrlType = cV4L2_CTRL_TYPE_STRING
	CtrlTypeBitMask             CtrlType = cV4L2_CTRL_TYPE_BITMASK
	CtrlTypeIntegerMenu         CtrlType = cV4L2_CTRL_TYPE_INTEGER_MENU
	CtrlTypeCompoundTypes       CtrlType = cV4L2_CTRL_COMPOUND_TYPES
	CtrlTypeU8                  CtrlType = cV4L2_CTRL_TYPE_U8
	CtrlTypeU16                 CtrlType = cV4L2_CTRL_TYPE_U16
	CtrlTypeU32                 CtrlType = cV4L2_CTRL_TYPE_U32
	CtrlTypeArear               CtrlType = cV4L2_CTRL_TYPE_AREA
	CtrlTypeHDR10CLLInfo        CtrlType = cV4L2_CTRL_TYPE_HDR10_CLL_INFO
	CtrlTypeHDRMasteringDisplay CtrlType = cV4L2_CTRL_TYPE_HDR10_MASTERING_DISPLAY
	CtrlTypeH264SPS             CtrlType = cV4L2_CTRL_TYPE_H264_SPS
	CtrlTypeH264PPS             CtrlType = cV4L2_CTRL_TYPE_H264_PPS
	CtrlTypeH264ScalingMatrix   CtrlType = cV4L2_CTRL_TYPE_H264_SCALING_MATRIX
	CtrlTypeH264SliceParams     CtrlType = cV4L2_CTRL_TYPE_H264_SLICE_PARAMS
	CtrlTypeH264DecodeParams    CtrlType = cV4L2_CTRL_TYPE_H264_DECODE_PARAMS
	CtrlTypeFWHTParams          CtrlType = cV4L2_CTRL_TYPE_FWHT_PARAMS
	CtrlTypeVP8Frame            CtrlType = cV4L2_CTRL_TYPE_VP8_FRAME
	CtrlTypeMPEG2Quantization   CtrlType = cV4L2_CTRL_TYPE_MPEG2_QUANTISATION
	CtrlTypeMPEG2Sequence       CtrlType = cV4L2_CTRL_TYPE_MPEG2_SEQUENCE
	CtrlTypeMPEG2Picture        CtrlType = cV4L2_CTRL_TYPE_MPEG2_PICTURE
	CtrlTypeVP9CompressedHDR    CtrlType = cV4L2_CTRL_TYPE_VP9_COMPRESSED_HDR
	CtrlTypeVP9Frame            CtrlType = cV4L2_CTRL_TYPE_VP9_FRAME
)

// CtrlID type for control values
type CtrlID = uint32

// PowerlineFrequency control enums
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L100
type PowerlineFrequency = uint32

const (
	PowerlineFrequencyDisabled PowerlineFrequency = cV4L2_CID_POWER_LINE_FREQUENCY_DISABLED
	PowerlineFrequency50Hz     PowerlineFrequency = cV4L2_CID_POWER_LINE_FREQUENCY_50HZ
	PowerlineFrequency60Hz     PowerlineFrequency = cV4L2_CID_POWER_LINE_FREQUENCY_60HZ
	PowerlineFrequencyAuto     PowerlineFrequency = cV4L2_CID_POWER_LINE_FREQUENCY_AUTO
)

// ColorFX control enums
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L114
type ColorFX = uint32

const (
	ColorFXNone         ColorFX = cV4L2_COLORFX_NONE
	ColorFXBlackWhite   ColorFX = cV4L2_COLORFX_BW
	ColorFXSepia        ColorFX = cV4L2_COLORFX_SEPIA
	ColorFXNegative     ColorFX = cV4L2_COLORFX_NEGATIVE
	ColorFXEmboss       ColorFX = cV4L2_COLORFX_EMBOSS
	ColorFXSketch       ColorFX = cV4L2_COLORFX_SKETCH
	ColorFXSkyBlue      ColorFX = cV4L2_COLORFX_SKY_BLUE
	ColorFXGrassGreen   ColorFX = cV4L2_COLORFX_GRASS_GREEN
	ColorFXSkinWhiten   ColorFX = cV4L2_COLORFX_SKIN_WHITEN
	ColorFXVivid        ColorFX = cV4L2_COLORFX_VIVID
	ColorFXAqua         ColorFX = cV4L2_COLORFX_AQUA
	ColorFXArtFreeze    ColorFX = cV4L2_COLORFX_ART_FREEZE
	ColorFXSilhouette   ColorFX = cV4L2_COLORFX_SILHOUETTE
	ColorFXSolarization ColorFX = cV4L2_COLORFX_SOLARIZATION
	ColorFXAntique      ColorFX = cV4L2_COLORFX_ANTIQUE
	ColorFXSetCBCR      ColorFX = cV4L2_COLORFX_SET_CBCR
	ColorFXSetRGB       ColorFX = cV4L2_COLORFX_SET_RGB
)

// User Controls IDs (CIDs)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L74
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html#control-id
const (
	CtrlBrightness              CtrlID = cV4L2_CID_BRIGHTNESS
	CtrlContrast                CtrlID = cV4L2_CID_CONTRAST
	CtrlSaturation              CtrlID = cV4L2_CID_SATURATION
	CtrlHue                     CtrlID = cV4L2_CID_HUE
	CtrlAutoWhiteBalance        CtrlID = cV4L2_CID_AUTO_WHITE_BALANCE
	CtrlDoWhiteBalance          CtrlID = cV4L2_CID_DO_WHITE_BALANCE
	CtrlRedBalance              CtrlID = cV4L2_CID_RED_BALANCE
	CtrlBlueBalance             CtrlID = cV4L2_CID_BLUE_BALANCE
	CtrlGamma                   CtrlID = cV4L2_CID_GAMMA
	CtrlExposure                CtrlID = cV4L2_CID_EXPOSURE
	CtrlAutogain                CtrlID = cV4L2_CID_AUTOGAIN
	CtrlGain                    CtrlID = cV4L2_CID_GAIN
	CtrlHFlip                   CtrlID = cV4L2_CID_HFLIP
	CtrlVFlip                   CtrlID = cV4L2_CID_VFLIP
	CtrlPowerlineFrequency      CtrlID = cV4L2_CID_POWER_LINE_FREQUENCY
	CtrlHueAuto                 CtrlID = cV4L2_CID_HUE_AUTO
	CtrlWhiteBalanceTemperature CtrlID = cV4L2_CID_WHITE_BALANCE_TEMPERATURE
	CtrlSharpness               CtrlID = cV4L2_CID_SHARPNESS
	CtrlBacklightCompensation   CtrlID = cV4L2_CID_BACKLIGHT_COMPENSATION
	CtrlChromaAutomaticGain     CtrlID = cV4L2_CID_CHROMA_AGC
	CtrlColorKiller             CtrlID = cV4L2_CID_COLOR_KILLER
	CtrlColorFX                 CtrlID = cV4L2_CID_COLORFX
	CtrlColorFXCBCR             CtrlID = cV4L2_CID_COLORFX_CBCR
	CtrlColorFXRGB              CtrlID = cV4L2_CID_COLORFX_RGB
	CtrlAutoBrightness          CtrlID = cV4L2_CID_AUTOBRIGHTNESS
	CtrlRotate                  CtrlID = cV4L2_CID_ROTATE
	CtrlBackgroundColor         CtrlID = cV4L2_CID_BG_COLOR
	CtrlMinimumCaptureBuffers   CtrlID = cV4L2_CID_MIN_BUFFERS_FOR_CAPTURE
	CtrlMinimumOutputBuffers    CtrlID = cV4L2_CID_MIN_BUFFERS_FOR_OUTPUT
	CtrlAlphaComponent          CtrlID = cV4L2_CID_ALPHA_COMPONENT
)

// Camera control values
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-camera.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L897
const (
	CtrlCameraClass                   CtrlID = cV4L2_CID_CAMERA_CLASS
	CtrlCameraExposureAuto            CtrlID = cV4L2_CID_EXPOSURE_AUTO
	CtrlCameraExposureAbsolute        CtrlID = cV4L2_CID_EXPOSURE_ABSOLUTE
	CtrlCameraExposureAutoPriority    CtrlID = cV4L2_CID_EXPOSURE_AUTO_PRIORITY
	CtrlCameraPanRelative             CtrlID = cV4L2_CID_PAN_RELATIVE
	CtrlCameraTiltRelative            CtrlID = cV4L2_CID_TILT_RELATIVE
	CtrlCameraPanReset                CtrlID = cV4L2_CID_PAN_RESET
	CtrlCameraTiltReset               CtrlID = cV4L2_CID_TILT_RESET
	CtrlCameraPanAbsolute             CtrlID = cV4L2_CID_PAN_ABSOLUTE
	CtrlCameraTiltAbsolute            CtrlID = cV4L2_CID_TILT_ABSOLUTE
	CtrlCameraFocusAbsolute           CtrlID = cV4L2_CID_FOCUS_ABSOLUTE
	CtrlCameraFocusRelative           CtrlID = cV4L2_CID_FOCUS_RELATIVE
	CtrlCameraFocusAuto               CtrlID = cV4L2_CID_FOCUS_AUTO
	CtrlCameraZoomAbsolute            CtrlID = cV4L2_CID_ZOOM_ABSOLUTE
	CtrlCameraZoomRelative            CtrlID = cV4L2_CID_ZOOM_RELATIVE
	CtrlCameraZoomContinuous          CtrlID = cV4L2_CID_ZOOM_CONTINUOUS
	CtrlCameraPrivacy                 CtrlID = cV4L2_CID_PRIVACY
	CtrlCameraIrisAbsolute            CtrlID = cV4L2_CID_IRIS_ABSOLUTE
	CtrlCameraIrisRelative            CtrlID = cV4L2_CID_IRIS_RELATIVE
	CtrlCameraAutoExposureBias        CtrlID = cV4L2_CID_AUTO_EXPOSURE_BIAS
	CtrlCameraAutoNPresetWhiteBalance CtrlID = cV4L2_CID_AUTO_N_PRESET_WHITE_BALANCE
	CtrlCameraWideDynamicRange        CtrlID = cV4L2_CID_WIDE_DYNAMIC_RANGE
	CtrlCameraImageStabilization      CtrlID = cV4L2_CID_IMAGE_STABILIZATION
	CtrlCameraIsoSensitivity          CtrlID = cV4L2_CID_ISO_SENSITIVITY
	CtrlCameraIsoSensitivityAuto      CtrlID = cV4L2_CID_ISO_SENSITIVITY_AUTO
	CtrlCameraExposureMetering        CtrlID = cV4L2_CID_EXPOSURE_METERING
	CtrlCameraSceneMode               CtrlID = cV4L2_CID_SCENE_MODE
	CtrlCamera3ALock                  CtrlID = cV4L2_CID_3A_LOCK
	CtrlCameraAutoFocusStart          CtrlID = cV4L2_CID_AUTO_FOCUS_START
	CtrlCameraAutoFocusStop           CtrlID = cV4L2_CID_AUTO_FOCUS_STOP
	CtrlCameraAutoFocusStatus         CtrlID = cV4L2_CID_AUTO_FOCUS_STATUS
	CtrlCameraAutoFocusRange          CtrlID = cV4L2_CID_AUTO_FOCUS_RANGE
	CtrlCameraPanSpeed                CtrlID = cV4L2_CID_PAN_SPEED
	CtrlCameraTiltSpeed               CtrlID = cV4L2_CID_TILT_SPEED
	CtrlCameraCameraOrientation       CtrlID = cV4L2_CID_CAMERA_ORIENTATION
	CtrlCameraCameraSensorRotation    CtrlID = cV4L2_CID_CAMERA_SENSOR_ROTATION
)

// Flash control values
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-flash.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1060
const (
	CtrlFlashClass              CtrlID = cV4L2_CID_FLASH_CLASS
	CtrlFlashLEDMode            CtrlID = cV4L2_CID_FLASH_LED_MODE
	CtrlFlashStrobeSource       CtrlID = cV4L2_CID_FLASH_STROBE_SOURCE
	CtrlFlashStrobe             CtrlID = cV4L2_CID_FLASH_STROBE
	CtrlFlashStrobeStop         CtrlID = cV4L2_CID_FLASH_STROBE_STOP
	CtrlFlashStrobeStatus       CtrlID = cV4L2_CID_FLASH_STROBE_STATUS
	CtrlFlashTimeout            CtrlID = cV4L2_CID_FLASH_TIMEOUT
	CtrlFlashIntensity          CtrlID = cV4L2_CID_FLASH_INTENSITY
	CtrlFlashTorchIntensity     CtrlID = cV4L2_CID_FLASH_TORCH_INTENSITY
	CtrlFlashIndicatorIntensity CtrlID = cV4L2_CID_FLASH_INDICATOR_INTENSITY
	CtrlFlashFault              CtrlID = cV4L2_CID_FLASH_FAULT
	CtrlFlashChargeCurrent      CtrlID = cV4L2_CID_FLASH_CHARGE
	CtrlFlashReady              CtrlID = cV4L2_CID_FLASH_READY
)

// JPEG control values
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-jpeg.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1104
const (
	CtrlJPEGClass          CtrlID = cV4L2_CID_JPEG_CLASS
	CtrlJPEGChromaSampling CtrlID = cV4L2_CID_JPEG_CHROMA_SUBSAMPLING
	// TODO add all JPEG flash controls
)

// Image source controls
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1127
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-image-source.html
const (
	CtrlImgSrcClass            CtrlID = cV4L2_CID_IMAGE_SOURCE_CLASS
	CtrlImgSrcVerticalBlank    CtrlID = cV4L2_CID_VBLANK
	CtrlImgSrcHorizontalBlank  CtrlID = cV4L2_CID_HBLANK
	CtrlImgSrcAnalogueGain     CtrlID = cV4L2_CID_ANALOGUE_GAIN
	CtrlImgSrcTestPatternRed   CtrlID = cV4L2_CID_TEST_PATTERN_RED
	CtrlImgSrcTestPatternGreenR CtrlID = cV4L2_CID_TEST_PATTERN_GREENR
	CtrlImgSrcTestPatternBlue  CtrlID = cV4L2_CID_TEST_PATTERN_BLUE
	CtrlImgSrcTestPatternGreenB CtrlID = cV4L2_CID_TEST_PATTERN_GREENB
	CtrlImgSrcUnitCellSize     CtrlID = cV4L2_CID_UNIT_CELL_SIZE
	CtrlImgSrcNotifyGains      CtrlID = cV4L2_CID_NOTIFY_GAINS
)

// Image process controls
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-image-process.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1144
const (
	CtrlImgProcClass          CtrlID = cV4L2_CID_IMAGE_PROC_CLASS
	CtrlImgProcLinkFreq       CtrlID = cV4L2_CID_LINK_FREQ
	CtrlImgProcPixelRate      CtrlID = cV4L2_CID_PIXEL_RATE
	CtrlImgProcTestPattern    CtrlID = cV4L2_CID_TEST_PATTERN
	CtrlImgProcDeinterlacing  CtrlID = cV4L2_CID_DEINTERLACING_MODE
	CtrlImgProcDigitalGain    CtrlID = cV4L2_CID_DIGITAL_GAIN
)

// Codec controls are already defined in ext_ctrls_*.go files
// See: ext_ctrls_h264.go, ext_ctrls_mpeg2.go, ext_ctrls_vp8.go, ext_ctrls_fwht.go
// Stateless codec controls (h264, vp8, fwht, mpeg2, etc)
