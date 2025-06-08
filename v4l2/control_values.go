package v4l2

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../include/
#include <linux/videodev2.h>
#include <linux/v4l2-controls.h>
*/
import "C"

// CtrlClass is a type alias for uint32, representing the class of a V4L2 control.
// Control classes are used to group related controls (e.g., user controls, camera controls).
// The class of a control can be queried using the CtrlTypeClass control type.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L56
type CtrlClass = uint32

// Control Class Constants
const (
	// CtrlClassUser defines the class for common user controls like brightness, contrast, etc.
	CtrlClassUser CtrlClass = C.V4L2_CTRL_CLASS_USER
	// CtrlClassCodec defines the class for codec-specific controls.
	CtrlClassCodec CtrlClass = C.V4L2_CTRL_CLASS_CODEC
	// CtrlClassCamera defines the class for camera-specific controls like exposure, focus, zoom.
	CtrlClassCamera CtrlClass = C.V4L2_CTRL_CLASS_CAMERA
	// CtrlClassFlash defines the class for camera flash controls.
	CtrlClassFlash CtrlClass = C.V4L2_CTRL_CLASS_FLASH
	// CtrlClassJPEG defines the class for JPEG compression controls.
	CtrlClassJPEG CtrlClass = C.V4L2_CTRL_CLASS_JPEG
	// CtrlClassImageSource defines the class for image source parameters (e.g., VBLANK).
	CtrlClassImageSource CtrlClass = C.V4L2_CTRL_CLASS_IMAGE_SOURCE
	// CtrlClassImageProcessing defines the class for image processing unit controls.
	CtrlClassImageProcessing CtrlClass = C.V4L2_CTRL_CLASS_IMAGE_PROC
	// CtrlClassDigitalVideo defines the class for Digital Video (DV) preset controls.
	CtrlClassDigitalVideo CtrlClass = C.V4L2_CTRL_CLASS_DV
	// CtrlClassDetection defines the class for motion or object detection controls.
	CtrlClassDetection CtrlClass = C.V4L2_CTRL_CLASS_DETECT
	// CtrlClassCodecStateless defines the class for stateless codec controls.
	CtrlClassCodecStateless CtrlClass = C.V4L2_CTRL_CLASS_CODEC_STATELESS
	// CtrlClassColorimitry defines the class for colorimetry controls (e.g. colorspace, transfer function).
	// Typo in source, should be CtrlClassColorimetry.
	CtrlClassColorimitry CtrlClass = C.V4L2_CTRL_CLASS_COLORIMETRY
)

var (
	// CtrlClasses is a slice containing all defined V4L2 control class constants.
	// This can be used for iterating or displaying available control classes.
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
	CtrlClassColorimitry, // Note: Typo in source, should be CtrlClassColorimetry
	}
)

// CtrlType is a type alias for uint32, representing the data type of a V4L2 control.
// The type determines how the control's value is interpreted (e.g., integer, boolean, menu).
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1799
// See also https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-queryctrl.html?highlight=v4l2_ctrl_type#c.V4L.v4l2_ctrl_type
type CtrlType uint32

// Control Type Constants
const (
	CtrlTypeInt                 CtrlType = C.V4L2_CTRL_TYPE_INTEGER            // Standard integer type.
	CtrlTypeBool                CtrlType = C.V4L2_CTRL_TYPE_BOOLEAN             // Boolean type (0 or 1).
	CtrlTypeMenu                CtrlType = C.V4L2_CTRL_TYPE_MENU                // Menu type, value is an index. Names are queried via VIDIOC_QUERYMENU.
	CtrlTypeButton              CtrlType = C.V4L2_CTRL_TYPE_BUTTON              // Button type, value is not used. Triggers an action.
	CtrlTypeInt64               CtrlType = C.V4L2_CTRL_TYPE_INTEGER64           // 64-bit integer type.
	CtrlTypeClass               CtrlType = C.V4L2_CTRL_TYPE_CTRL_CLASS          // Control class identifier.
	CtrlTypeString              CtrlType = C.V4L2_CTRL_TYPE_STRING              // String type.
	CtrlTypeBitMask             CtrlType = C.V4L2_CTRL_TYPE_BITMASK             // Bitmask type.
	CtrlTypeIntegerMenu         CtrlType = C.V4L2_CTRL_TYPE_INTEGER_MENU        // Integer menu type, value is an integer. Names are queried via VIDIOC_QUERYMENU.
	CtrlTypeCompoundTypes       CtrlType = C.V4L2_CTRL_COMPOUND_TYPES       // Marks start of compound types, not a type itself.
	CtrlTypeU8                  CtrlType = C.V4L2_CTRL_TYPE_U8                  // Unsigned 8-bit integer.
	CtrlTypeU16                 CtrlType = C.V4L2_CTRL_TYPE_U16                 // Unsigned 16-bit integer.
	CtrlTypeU32                 CtrlType = C.V4L2_CTRL_TYPE_U32                 // Unsigned 32-bit integer.
	CtrlTypeArear               CtrlType = C.V4L2_CTRL_TYPE_AREA                // Area type, for selections (e.g., focus area). Typo in source, should be CtrlTypeArea.
	CtrlTypeHDR10CLLInfo        CtrlType = C.V4L2_CTRL_TYPE_HDR10_CLL_INFO      // HDR10 Content Light Level information.
	CtrlTypeHDRMasteringDisplay CtrlType = C.V4L2_CTRL_TYPE_HDR10_MASTERING_DISPLAY // HDR10 Mastering Display information.
	CtrlTypeH264SPS             CtrlType = C.V4L2_CTRL_TYPE_H264_SPS             // H.264 Sequence Parameter Set.
	CtrlTypeH264PPS             CtrlType = C.V4L2_CTRL_TYPE_H264_PPS             // H.264 Picture Parameter Set.
	CtrlTypeH264ScalingMatrix   CtrlType = C.V4L2_CTRL_TYPE_H264_SCALING_MATRIX   // H.264 Scaling Matrix.
	CtrlTypeH264SliceParams     CtrlType = C.V4L2_CTRL_TYPE_H264_SLICE_PARAMS     // H.264 Slice Parameters.
	CtrlTypeH264DecodeParams    CtrlType = C.V4L2_CTRL_TYPE_H264_DECODE_PARAMS    // H.264 Decode Parameters.
	CtrlTypeFWHTParams          CtrlType = C.V4L2_CTRL_TYPE_FWHT_PARAMS          // FWHT (Fast Walsh-Hadamard Transform) codec parameters.
	CtrlTypeVP8Frame            CtrlType = C.V4L2_CTRL_TYPE_VP8_FRAME            // VP8 frame parameters.
	CtrlTypeMPEG2Quantization   CtrlType = C.V4L2_CTRL_TYPE_MPEG2_QUANTISATION   // MPEG-2 Quantization tables.
	CtrlTypeMPEG2Sequence       CtrlType = C.V4L2_CTRL_TYPE_MPEG2_SEQUENCE       // MPEG-2 Sequence parameters.
	CtrlTypeMPEG2Picture        CtrlType = C.V4L2_CTRL_TYPE_MPEG2_PICTURE        // MPEG-2 Picture parameters.
	CtrlTypeVP9CompressedHDR    CtrlType = C.V4L2_CTRL_TYPE_VP9_COMPRESSED_HDR    // VP9 Compressed HDR parameters.
	CtrlTypeVP9Frame            CtrlType = C.V4L2_CTRL_TYPE_VP9_FRAME            // VP9 Frame parameters.
)

// CtrlID is a type alias for uint32, representing the unique identifier of a V4L2 control.
// These IDs are used to query and modify specific controls.
type CtrlID = uint32

// PowerlineFrequency is a type alias for uint32, representing the power line frequency setting
// used to counteract flickering caused by artificial lighting.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L100
type PowerlineFrequency = uint32

// Powerline Frequency Enum Values for CtrlPowerlineFrequency control.
const (
	PowerlineFrequencyDisabled PowerlineFrequency = C.V4L2_CID_POWER_LINE_FREQUENCY_DISABLED // Powerline frequency filter disabled.
	PowerlineFrequency50Hz     PowerlineFrequency = C.V4L2_CID_POWER_LINE_FREQUENCY_50HZ     // 50 Hz powerline frequency.
	PowerlineFrequency60Hz     PowerlineFrequency = C.V4L2_CID_POWER_LINE_FREQUENCY_60HZ     // 60 Hz powerline frequency.
	PowerlineFrequencyAuto     PowerlineFrequency = C.V4L2_CID_POWER_LINE_FREQUENCY_AUTO     // Auto-detect powerline frequency.
)

// ColorFX is a type alias for uint32, representing color effect settings.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L114
type ColorFX = uint32

// Color Effects Enum Values for CtrlColorFX control.
const (
	ColorFXNone         ColorFX = C.V4L2_COLORFX_NONE         // No color effect.
	ColorFXBlackWhite   ColorFX = C.V4L2_COLORFX_BW           // Black and white (grayscale).
	ColorFXSepia        ColorFX = C.V4L2_COLORFX_SEPIA        // Sepia tone.
	ColorFXNegative     ColorFX = C.V4L2_COLORFX_NEGATIVE     // Negative colors.
	ColorFXEmboss       ColorFX = C.V4L2_COLORFX_EMBOSS       // Emboss effect.
	ColorFXSketch       ColorFX = C.V4L2_COLORFX_SKETCH       // Sketch effect.
	ColorFXSkyBlue      ColorFX = C.V4L2_COLORFX_SKY_BLUE     // Sky blue tint.
	ColorFXGrassGreen   ColorFX = C.V4L2_COLORFX_GRASS_GREEN  // Grass green tint.
	ColorFXSkinWhiten   ColorFX = C.V4L2_COLORFX_SKIN_WHITEN  // Skin whiten effect.
	ColorFXVivid        ColorFX = C.V4L2_COLORFX_VIVID        // Vivid colors.
	ColorFXAqua         ColorFX = C.V4L2_COLORFX_AQUA         // Aqua tint.
	ColorFXArtFreeze    ColorFX = C.V4L2_COLORFX_ART_FREEZE   // Art freeze effect.
	ColorFXSilhouette   ColorFX = C.V4L2_COLORFX_SILHOUETTE   // Silhouette effect.
	ColorFXSolarization ColorFX = C.V4L2_COLORFX_SOLARIZATION // Solarization effect.
	ColorFXAntique      ColorFX = C.V4L2_COLORFX_ANTIQUE      // Antique effect.
	ColorFXSetCBCR      ColorFX = C.V4L2_COLORFX_SET_CBCR     // Set Cb/Cr values directly.
	ColorFXSetRGB       ColorFX = C.V4L2_COLORFX_SET_RGB        // Set RGB values directly.
)

// User Control IDs (CtrlID constants for common user-adjustable controls).
// These typically belong to the CtrlClassUser class.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L74
// See also https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/control.html#control-id
const (
	// CtrlBrightness adjusts image brightness.
	CtrlBrightness CtrlID = C.V4L2_CID_BRIGHTNESS
	// CtrlContrast adjusts image contrast.
	CtrlContrast CtrlID = C.V4L2_CID_CONTRAST
	// CtrlSaturation adjusts image color saturation.
	CtrlSaturation CtrlID = C.V4L2_CID_SATURATION
	// CtrlHue adjusts image hue.
	CtrlHue CtrlID = C.V4L2_CID_HUE
	// CtrlAutoWhiteBalance enables or disables automatic white balance.
	CtrlAutoWhiteBalance CtrlID = C.V4L2_CID_AUTO_WHITE_BALANCE
	// CtrlDoWhiteBalance triggers a one-shot white balance adjustment.
	CtrlDoWhiteBalance CtrlID = C.V4L2_CID_DO_WHITE_BALANCE
	// CtrlRedBalance adjusts the red color component.
	CtrlRedBalance CtrlID = C.V4L2_CID_RED_BALANCE
	// CtrlBlueBalance adjusts the blue color component.
	CtrlBlueBalance CtrlID = C.V4L2_CID_BLUE_BALANCE
	// CtrlGamma adjusts image gamma.
	CtrlGamma CtrlID = C.V4L2_CID_GAMMA
	// CtrlExposure is an alias for CtrlCameraExposureAbsolute for older UVC drivers.
	CtrlExposure CtrlID = C.V4L2_CID_EXPOSURE
	// CtrlAutogain enables or disables automatic gain control.
	CtrlAutogain CtrlID = C.V4L2_CID_AUTOGAIN
	// CtrlGain adjusts image gain.
	CtrlGain CtrlID = C.V4L2_CID_GAIN
	// CtrlHFlip flips the image horizontally.
	CtrlHFlip CtrlID = C.V4L2_CID_HFLIP
	// CtrlVFlip flips the image vertically.
	CtrlVFlip CtrlID = C.V4L2_CID_VFLIP
	// CtrlPowerlineFrequency sets the power line frequency (e.g., 50Hz, 60Hz) to reduce flicker.
	CtrlPowerlineFrequency CtrlID = C.V4L2_CID_POWER_LINE_FREQUENCY
	// CtrlHueAuto enables or disables automatic hue adjustment.
	CtrlHueAuto CtrlID = C.V4L2_CID_HUE_AUTO
	// CtrlWhiteBalanceTemperature sets the white balance temperature.
	CtrlWhiteBalanceTemperature CtrlID = C.V4L2_CID_WHITE_BALANCE_TEMPERATURE
	// CtrlSharpness adjusts image sharpness.
	CtrlSharpness CtrlID = C.V4L2_CID_SHARPNESS
	// CtrlBacklightCompensation enables or disables backlight compensation.
	CtrlBacklightCompensation CtrlID = C.V4L2_CID_BACKLIGHT_COMPENSATION
	// CtrlChromaAutomaticGain enables or disables automatic chroma gain. (DEPRECATED)
	CtrlChromaAutomaticGain CtrlID = C.V4L2_CID_CHROMA_AGC
	// CtrlColorKiller enables or disables color killer (forces B&W).
	CtrlColorKiller CtrlID = C.V4L2_CID_COLOR_KILLER
	// CtrlColorFX selects a color effect (e.g., sepia, negative). See ColorFX constants.
	CtrlColorFX CtrlID = C.V4L2_CID_COLORFX
	// CtrlColorFXCBCR sets Cb and Cr components for V4L2_COLORFX_SET_CBCR.
	CtrlColorFXCBCR CtrlID = C.V4L2_CID_COLORFX_CBCR
	// CtrlColorFXRGB sets R, G and B components for V4L2_COLORFX_SET_RGB.
	CtrlColorFXRGB CtrlID = C.V4L2_CID_COLORFX_RGB
	// CtrlAutoBrightness enables or disables automatic brightness adjustment.
	CtrlAutoBrightness CtrlID = C.V4L2_CID_AUTOBRIGHTNESS
	// CtrlRotate rotates the image by a specified angle (e.g., 90, 180, 270 degrees).
	CtrlRotate CtrlID = C.V4L2_CID_ROTATE
	// CtrlBackgroundColor sets the background color for overlay.
	CtrlBackgroundColor CtrlID = C.V4L2_CID_BG_COLOR
	// CtrlMinimumCaptureBuffers reports the minimum number of buffers required for capture. (Read-only)
	CtrlMinimumCaptureBuffers CtrlID = C.V4L2_CID_MIN_BUFFERS_FOR_CAPTURE
	// CtrlMinimumOutputBuffers reports the minimum number of buffers required for output. (Read-only)
	CtrlMinimumOutputBuffers CtrlID = C.V4L2_CID_MIN_BUFFERS_FOR_OUTPUT
	// CtrlAlphaComponent sets the global alpha component value.
	CtrlAlphaComponent CtrlID = C.V4L2_CID_ALPHA_COMPONENT
)

// Camera Control IDs (CtrlID constants specific to camera devices).
// These typically belong to the CtrlClassCamera class.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-camera.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L897
const (
	// CtrlCameraClass identifies the camera control class.
	CtrlCameraClass CtrlID = C.V4L2_CID_CAMERA_CLASS
	// CtrlCameraExposureAuto sets the auto exposure mode.
	CtrlCameraExposureAuto CtrlID = C.V4L2_CID_EXPOSURE_AUTO
	// CtrlCameraExposureAbsolute sets the absolute exposure time.
	CtrlCameraExposureAbsolute CtrlID = C.V4L2_CID_EXPOSURE_ABSOLUTE
	// CtrlCameraExposureAutoPriority enables or disables auto exposure priority mode (for UVC cameras).
	CtrlCameraExposureAutoPriority CtrlID = C.V4L2_CID_EXPOSURE_AUTO_PRIORITY
	// CtrlCameraPanRelative sets the relative pan (horizontal movement).
	CtrlCameraPanRelative CtrlID = C.V4L2_CID_PAN_RELATIVE
	// CtrlCameraTiltRelative sets the relative tilt (vertical movement).
	CtrlCameraTiltRelative CtrlID = C.V4L2_CID_TILT_RELATIVE
	// CtrlCameraPanReset resets the pan position.
	CtrlCameraPanReset CtrlID = C.V4L2_CID_PAN_RESET
	// CtrlCameraTiltReset resets the tilt position.
	CtrlCameraTiltReset CtrlID = C.V4L2_CID_TILT_RESET
	// CtrlCameraPanAbsolute sets the absolute pan position.
	CtrlCameraPanAbsolute CtrlID = C.V4L2_CID_PAN_ABSOLUTE
	// CtrlCameraTiltAbsolute sets the absolute tilt position.
	CtrlCameraTiltAbsolute CtrlID = C.V4L2_CID_TILT_ABSOLUTE
	// CtrlCameraFocusAbsolute sets the absolute focus position.
	CtrlCameraFocusAbsolute CtrlID = C.V4L2_CID_FOCUS_ABSOLUTE
	// CtrlCameraFocusRelative sets the relative focus adjustment.
	CtrlCameraFocusRelative CtrlID = C.V4L2_CID_FOCUS_RELATIVE
	// CtrlCameraFocusAuto enables or disables automatic focus.
	CtrlCameraFocusAuto CtrlID = C.V4L2_CID_FOCUS_AUTO
	// CtrlCameraZoomAbsolute sets the absolute zoom level.
	CtrlCameraZoomAbsolute CtrlID = C.V4L2_CID_ZOOM_ABSOLUTE
	// CtrlCameraZoomRelative sets the relative zoom adjustment.
	CtrlCameraZoomRelative CtrlID = C.V4L2_CID_ZOOM_RELATIVE
	// CtrlCameraZoomContinuous sets the continuous zoom speed.
	CtrlCameraZoomContinuous CtrlID = C.V4L2_CID_ZOOM_CONTINUOUS
	// CtrlCameraPrivacy enables or disables the privacy shutter.
	CtrlCameraPrivacy CtrlID = C.V4L2_CID_PRIVACY
	// CtrlCameraIrisAbsolute sets the absolute iris aperture.
	CtrlCameraIrisAbsolute CtrlID = C.V4L2_CID_IRIS_ABSOLUTE
	// CtrlCameraIrisRelative sets the relative iris adjustment.
	CtrlCameraIrisRelative CtrlID = C.V4L2_CID_IRIS_RELATIVE
	// CtrlCameraAutoExposureBias sets the auto exposure bias.
	CtrlCameraAutoExposureBias CtrlID = C.V4L2_CID_AUTO_EXPOSURE_BIAS
	// CtrlCameraAutoNPresetWhiteBalance sets the auto white balance preset.
	CtrlCameraAutoNPresetWhiteBalance CtrlID = C.V4L2_CID_AUTO_N_PRESET_WHITE_BALANCE
	// CtrlCameraWideDynamicRange enables or disables wide dynamic range.
	CtrlCameraWideDynamicRange CtrlID = C.V4L2_CID_WIDE_DYNAMIC_RANGE
	// CtrlCameraImageStabilization enables or disables image stabilization.
	CtrlCameraImageStabilization CtrlID = C.V4L2_CID_IMAGE_STABILIZATION
	// CtrlCameraIsoSensitivity sets the ISO sensitivity.
	CtrlCameraIsoSensitivity CtrlID = C.V4L2_CID_ISO_SENSITIVITY
	// CtrlCameraIsoSensitivityAuto sets the auto ISO sensitivity mode.
	CtrlCameraIsoSensitivityAuto CtrlID = C.V4L2_CID_ISO_SENSITIVITY_AUTO
	// CtrlCameraExposureMetering sets the exposure metering mode.
	CtrlCameraExposureMetering CtrlID = C.V4L2_CID_EXPOSURE_METERING
	// CtrlCameraSceneMode sets the scene mode (e.g., sports, night).
	CtrlCameraSceneMode CtrlID = C.V4L2_CID_SCENE_MODE
	// CtrlCamera3ALock locks or unlocks auto exposure, auto white balance, and auto focus.
	CtrlCamera3ALock CtrlID = C.V4L2_CID_3A_LOCK
	// CtrlCameraAutoFocusStart starts a one-shot auto focus operation.
	CtrlCameraAutoFocusStart CtrlID = C.V4L2_CID_AUTO_FOCUS_START
	// CtrlCameraAutoFocusStop stops a one-shot auto focus operation.
	CtrlCameraAutoFocusStop CtrlID = C.V4L2_CID_AUTO_FOCUS_STOP
	// CtrlCameraAutoFocusStatus reports the status of auto focus. (Read-only)
	CtrlCameraAutoFocusStatus CtrlID = C.V4L2_CID_AUTO_FOCUS_STATUS
	// CtrlCameraAutoFocusRange sets the auto focus range.
	CtrlCameraAutoFocusRange CtrlID = C.V4L2_CID_AUTO_FOCUS_RANGE
	// CtrlCameraPanSpeed sets the pan speed for continuous pan.
	CtrlCameraPanSpeed CtrlID = C.V4L2_CID_PAN_SPEED
	// CtrlCameraTiltSpeed sets the tilt speed for continuous tilt.
	CtrlCameraTiltSpeed CtrlID = C.V4L2_CID_TILT_SPEED
	// CtrlCameraCameraOrientation reports the physical orientation of the camera sensor on the host system. (Read-only)
	CtrlCameraCameraOrientation CtrlID = C.V4L2_CID_CAMERA_ORIENTATION
	// CtrlCameraCameraSensorRotation reports the rotation of the camera sensor relative to the camera housing. (Read-only)
	CtrlCameraCameraSensorRotation CtrlID = C.V4L2_CID_CAMERA_SENSOR_ROTATION
)

// Flash Control IDs (CtrlID constants specific to camera flash units).
// These typically belong to the CtrlClassFlash class.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-flash.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1060
const (
	// CtrlFlashClass identifies the flash control class.
	CtrlFlashClass CtrlID = C.V4L2_CID_FLASH_CLASS
	// CtrlFlashLEDMode sets the LED flash mode.
	CtrlFlashLEDMode CtrlID = C.V4L2_CID_FLASH_LED_MODE
	// TODO add all flash control const values from <linux/v4l2-controls.h>
)

// JPEG Control IDs (CtrlID constants specific to JPEG compression).
// These typically belong to the CtrlClassJPEG class.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-jpeg.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1104
const (
	// CtrlJPEGClass identifies the JPEG control class.
	CtrlJPEGClass CtrlID = C.V4L2_CID_JPEG_CLASS
	// CtrlJPEGChromaSampling sets the chroma subsampling format (e.g., 4:4:4, 4:2:2, 4:2:0).
	CtrlJPEGChromaSampling CtrlID = C.V4L2_CID_JPEG_CHROMA_SUBSAMPLING
	// TODO add all JPEG control const values from <linux/v4l2-controls.h>
)

// Image Source Control IDs (CtrlID constants for image source parameters).
// These typically belong to the CtrlClassImageSource class.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1127
// See also https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-image-source.html
const (
	// CtrlImgSrcClass identifies the image source control class.
	CtrlImgSrcClass CtrlID = C.V4L2_CID_IMAGE_SOURCE_CLASS
	// CtrlImgSrcVerticalBlank sets the vertical blanking interval.
	CtrlImgSrcVerticalBlank CtrlID = C.V4L2_CID_VBLANK
)

// Image Process Control IDs (CtrlID constants for image processing unit parameters).
// These typically belong to the CtrlClassImageProcessing class.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-image-process.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L1144
const (
	// CtrlImgProcClass identifies the image processing control class.
	CtrlImgProcClass = C.V4L2_CID_IMAGE_PROC_CLASS
	// TODO implement all image process control const values from <linux/v4l2-controls.h>
)

// TODO add code for the following controls:
// Stateless codec controls (h264, vp8, fwht, mpeg2, etc)
