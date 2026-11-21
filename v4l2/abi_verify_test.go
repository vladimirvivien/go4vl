//go:build cgo_verify

package v4l2

import (
	"testing"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

func TestABIStructSizes(t *testing.T) {
	goSizes := map[string]uintptr{
		"cTimeval":               unsafe.Sizeof(cTimeval{}),
		"sys.Timeval":            unsafe.Sizeof(sys.Timeval{}),
		"cTimespec":              unsafe.Sizeof(cTimespec{}),
		"cV4L2Timecode":          unsafe.Sizeof(cV4L2Timecode{}),
		"Timecode":               unsafe.Sizeof(Timecode{}),
		"cV4L2Fract":             unsafe.Sizeof(cV4L2Fract{}),
		"Fract":                  unsafe.Sizeof(Fract{}),
		"cV4L2Rect":              unsafe.Sizeof(cV4L2Rect{}),
		"Rect":                   unsafe.Sizeof(Rect{}),
		"cMediaDeviceInfo":             unsafe.Sizeof(cMediaDeviceInfo{}),
		"cMediaV2Topology":             unsafe.Sizeof(cMediaV2Topology{}),
		"cMediaV2Entity":               unsafe.Sizeof(cMediaV2Entity{}),
		"cMediaV2Interface":            unsafe.Sizeof(cMediaV2Interface{}),
		"cMediaV2IntfDevnode":          unsafe.Sizeof(cMediaV2IntfDevnode{}),
		"cMediaV2Pad":                  unsafe.Sizeof(cMediaV2Pad{}),
		"cMediaV2Link":                 unsafe.Sizeof(cMediaV2Link{}),
		"cMediaEntityDesc":             unsafe.Sizeof(cMediaEntityDesc{}),
		"cMediaPadDesc":                unsafe.Sizeof(cMediaPadDesc{}),
		"cMediaLinkDesc":               unsafe.Sizeof(cMediaLinkDesc{}),
		"cMediaLinksEnum":              unsafe.Sizeof(cMediaLinksEnum{}),
		"cV4L2MBusFrameFmt":            unsafe.Sizeof(cV4L2MBusFrameFmt{}),
		"cV4L2SubdevFormat":            unsafe.Sizeof(cV4L2SubdevFormat{}),
		"cV4L2SubdevCrop":              unsafe.Sizeof(cV4L2SubdevCrop{}),
		"cV4L2SubdevMBusCodeEnum":      unsafe.Sizeof(cV4L2SubdevMBusCodeEnum{}),
		"cV4L2SubdevFrameSizeEnum":     unsafe.Sizeof(cV4L2SubdevFrameSizeEnum{}),
		"cV4L2SubdevFrameInterval":     unsafe.Sizeof(cV4L2SubdevFrameInterval{}),
		"cV4L2SubdevFrameIntervalEnum": unsafe.Sizeof(cV4L2SubdevFrameIntervalEnum{}),
		"cV4L2SubdevSelection":         unsafe.Sizeof(cV4L2SubdevSelection{}),
		"cV4L2Capability":        unsafe.Sizeof(cV4L2Capability{}),
		"cV4L2Audio":             unsafe.Sizeof(cV4L2Audio{}),
		"cV4L2AudioOut":          unsafe.Sizeof(cV4L2AudioOut{}),
		"cV4L2BTTimings":         unsafe.Sizeof(cV4L2BTTimings{}),
		"cV4L2BTTimingsCap":      unsafe.Sizeof(cV4L2BTTimingsCap{}),
		"cV4L2Buffer":            unsafe.Sizeof(cV4L2Buffer{}),
		"cV4L2CaptureParm":       unsafe.Sizeof(cV4L2CaptureParm{}),
		"CaptureParam":           unsafe.Sizeof(CaptureParam{}),
		"cV4L2Control":           unsafe.Sizeof(cV4L2Control{}),
		"cV4L2CreateBuffers":     unsafe.Sizeof(cV4L2CreateBuffers{}),
		"cV4L2Crop":              unsafe.Sizeof(cV4L2Crop{}),
		"cV4L2CropCap":           unsafe.Sizeof(cV4L2CropCap{}),
		"cV4L2DecoderCmd":        unsafe.Sizeof(cV4L2DecoderCmd{}),
		"cV4L2DVTimings":         unsafe.Sizeof(cV4L2DVTimings{}),
		"cV4L2DVTimingsCap":      unsafe.Sizeof(cV4L2DVTimingsCap{}),
		"cV4L2EncoderCmd":        unsafe.Sizeof(cV4L2EncoderCmd{}),
		"cV4L2EnumDVTimings":     unsafe.Sizeof(cV4L2EnumDVTimings{}),
		"cV4L2EventCtrl":         unsafe.Sizeof(cV4L2EventCtrl{}),
		"cV4L2EventFrameSync":    unsafe.Sizeof(cV4L2EventFrameSync{}),
		"cV4L2EventMotionDet":    unsafe.Sizeof(cV4L2EventMotionDet{}),
		"cV4L2EventSrcChange":    unsafe.Sizeof(cV4L2EventSrcChange{}),
		"cV4L2EventVsync":        unsafe.Sizeof(cV4L2EventVsync{}),
		"cV4L2Event":             unsafe.Sizeof(cV4L2Event{}),
		"cV4L2EventSubscription": unsafe.Sizeof(cV4L2EventSubscription{}),
		"cV4L2ExportBuffer":      unsafe.Sizeof(cV4L2ExportBuffer{}),
		"cV4L2ExtControl":        unsafe.Sizeof(cV4L2ExtControl{}),
		"cV4L2ExtControls":       unsafe.Sizeof(cV4L2ExtControls{}),
		"cV4L2FmtDesc":           unsafe.Sizeof(cV4L2FmtDesc{}),
		"cV4L2Format":            unsafe.Sizeof(cV4L2Format{}),
		"cV4L2Frequency":         unsafe.Sizeof(cV4L2Frequency{}),
		"cV4L2FrequencyBand":     unsafe.Sizeof(cV4L2FrequencyBand{}),
		"cV4L2FrmSizeDiscrete":   unsafe.Sizeof(cV4L2FrmSizeDiscrete{}),
		"FrameSizeDiscrete":      unsafe.Sizeof(FrameSizeDiscrete{}),
		"cV4L2FrmSizeStepwise":   unsafe.Sizeof(cV4L2FrmSizeStepwise{}),
		"FrameSize":              unsafe.Sizeof(FrameSize{}),
		"cV4L2FrmSizeEnum":       unsafe.Sizeof(cV4L2FrmSizeEnum{}),
		"cV4L2FrmIvalStepwise":   unsafe.Sizeof(cV4L2FrmIvalStepwise{}),
		"FrameInterval":          unsafe.Sizeof(FrameInterval{}),
		"cV4L2FrmIvalEnum":       unsafe.Sizeof(cV4L2FrmIvalEnum{}),
		"cV4L2Input":             unsafe.Sizeof(cV4L2Input{}),
		"cV4L2Modulator":         unsafe.Sizeof(cV4L2Modulator{}),
		"cV4L2Output":            unsafe.Sizeof(cV4L2Output{}),
		"cV4L2OutputParm":        unsafe.Sizeof(cV4L2OutputParm{}),
		"OutputParam":            unsafe.Sizeof(OutputParam{}),
		"cV4L2PixFormat":         unsafe.Sizeof(cV4L2PixFormat{}),
		"cV4L2PlanePixFormat":    unsafe.Sizeof(cV4L2PlanePixFormat{}),
		"cV4L2PixFormatMPlane":   unsafe.Sizeof(cV4L2PixFormatMPlane{}),
		"cV4L2Plane":             unsafe.Sizeof(cV4L2Plane{}),
		"cV4L2QueryCtrl":         unsafe.Sizeof(cV4L2QueryCtrl{}),
		"cV4L2QueryExtCtrl":      unsafe.Sizeof(cV4L2QueryExtCtrl{}),
		"cV4L2QueryMenu":         unsafe.Sizeof(cV4L2QueryMenu{}),
		"cV4L2RequestBuffers":    unsafe.Sizeof(cV4L2RequestBuffers{}),
		"RequestBuffers":         unsafe.Sizeof(RequestBuffers{}),
		"cV4L2Selection":         unsafe.Sizeof(cV4L2Selection{}),
		"cV4L2Standard":          unsafe.Sizeof(cV4L2Standard{}),
		"cV4L2StreamParm":        unsafe.Sizeof(cV4L2StreamParm{}),
		"cV4L2Tuner":             unsafe.Sizeof(cV4L2Tuner{}),
	}

	for _, tc := range cStructSizeChecks(goSizes) {
		if tc.got != tc.want {
			t.Errorf("sizeof(%s): got %d, want %d", tc.name, tc.got, tc.want)
		}
	}
}

func TestABIStructOffsets(t *testing.T) {
	goOffsets := map[string]uintptr{
		"cV4L2ExtControl.Value":     unsafe.Offsetof(cV4L2ExtControl{}.Value),
		"cV4L2Format.Fmt":           unsafe.Offsetof(cV4L2Format{}.Fmt),
		"cV4L2PixFormat.Ycbcr_enc":  unsafe.Offsetof(cV4L2PixFormat{}.Ycbcr_enc),
		"cV4L2Buffer.M":             unsafe.Offsetof(cV4L2Buffer{}.M),
		"cV4L2Buffer.Fd":            unsafe.Offsetof(cV4L2Buffer{}.Fd),
		"cV4L2Event.U":              unsafe.Offsetof(cV4L2Event{}.U),
		"cV4L2EventCtrl.Value":      unsafe.Offsetof(cV4L2EventCtrl{}.Value),
		"cV4L2FrmSizeEnum.Discrete": unsafe.Offsetof(cV4L2FrmSizeEnum{}.Discrete),
		"cV4L2FrmIvalEnum.Discrete": unsafe.Offsetof(cV4L2FrmIvalEnum{}.Discrete),
		"cV4L2StreamParm.Parm":      unsafe.Offsetof(cV4L2StreamParm{}.Parm),
		"cV4L2DVTimings.Pad_cgo_0":  unsafe.Offsetof(cV4L2DVTimings{}.Pad_cgo_0),
		"cV4L2DVTimingsCap.Bt":      unsafe.Offsetof(cV4L2DVTimingsCap{}.Bt),
		"cV4L2DecoderCmd.Stop":      unsafe.Offsetof(cV4L2DecoderCmd{}.Stop),
		"cV4L2QueryMenu.Name":             unsafe.Offsetof(cV4L2QueryMenu{}.Name),
		"cV4L2Plane.M":                    unsafe.Offsetof(cV4L2Plane{}.M),
		"cV4L2Plane.Offset":               unsafe.Offsetof(cV4L2Plane{}.Offset),
		"cV4L2PixFormatMPlane.Num_planes": unsafe.Offsetof(cV4L2PixFormatMPlane{}.Num_planes),
		"cV4L2CreateBuffers.Format":       unsafe.Offsetof(cV4L2CreateBuffers{}.Format),
		"cV4L2Selection.R":                unsafe.Offsetof(cV4L2Selection{}.R),
		"cMediaV2Topology.Ptr_links":      unsafe.Offsetof(cMediaV2Topology{}.Ptr_links),
		"cMediaV2Interface.Devnode":       unsafe.Offsetof(cMediaV2Interface{}.Devnode),
		"cMediaEntityDesc.Dev":            unsafe.Offsetof(cMediaEntityDesc{}.Dev),
		"cV4L2SubdevFormat.Format":        unsafe.Offsetof(cV4L2SubdevFormat{}.Format),
		"cV4L2SubdevSelection.R":          unsafe.Offsetof(cV4L2SubdevSelection{}.R),
	}

	for _, tc := range cStructOffsetChecks(goOffsets) {
		if tc.got != tc.want {
			t.Errorf("offsetof(%s): got %d, want %d", tc.name, tc.got, tc.want)
		}
	}
}

func TestABIIoctlConstants(t *testing.T) {
	for _, tc := range cIoctlChecks() {
		if tc.got != tc.want {
			t.Errorf("%s: got 0x%x, want 0x%x", tc.name, tc.got, tc.want)
		}
	}
}

func TestABIValueConstants(t *testing.T) {
	checks, stdAll := cValueConstantChecks()
	for _, tc := range checks {
		if tc.got != tc.want {
			t.Errorf("%s: got 0x%x, want 0x%x", tc.name, tc.got, tc.want)
		}
	}
	if StdAll != stdAll {
		t.Errorf("V4L2_STD_ALL: got 0x%x, want 0x%x", StdAll, stdAll)
	}
}
