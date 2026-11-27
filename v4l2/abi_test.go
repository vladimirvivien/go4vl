package v4l2

import (
	"testing"
	"unsafe"
)

type abiLayoutCase struct {
	name string
	got  uintptr
	want uintptr
}

func assertABILayouts(t *testing.T, cases []abiLayoutCase) {
	t.Helper()
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %d, want %d", tc.name, tc.got, tc.want)
		}
	}
}

func TestABIWrapperSizes(t *testing.T) {
	assertABILayouts(t, []abiLayoutCase{
		{"sizeof(Timecode)", unsafe.Sizeof(Timecode{}), unsafe.Sizeof(cV4L2Timecode{})},
		{"sizeof(Fract)", unsafe.Sizeof(Fract{}), unsafe.Sizeof(cV4L2Fract{})},
		{"sizeof(Rect)", unsafe.Sizeof(Rect{}), unsafe.Sizeof(cV4L2Rect{})},
		{"sizeof(CaptureParam)", unsafe.Sizeof(CaptureParam{}), unsafe.Sizeof(cV4L2CaptureParm{})},
		{"sizeof(OutputParam)", unsafe.Sizeof(OutputParam{}), unsafe.Sizeof(cV4L2OutputParm{})},
		{"sizeof(FrameSizeDiscrete)", unsafe.Sizeof(FrameSizeDiscrete{}), unsafe.Sizeof(cV4L2FrmSizeDiscrete{})},
		{"sizeof(FrameSize)", unsafe.Sizeof(FrameSize{}), unsafe.Sizeof(cV4L2FrmSizeStepwise{})},
		{"sizeof(FrameInterval)", unsafe.Sizeof(FrameInterval{}), unsafe.Sizeof(cV4L2FrmIvalStepwise{})},
		{"sizeof(RequestBuffers)", unsafe.Sizeof(RequestBuffers{}), unsafe.Sizeof(cV4L2RequestBuffers{})},
	})
}

func TestABIRegexPatchedAndTimeLayouts(t *testing.T) {
	const is64Bit = unsafe.Sizeof(uintptr(0)) == 8

	cases := []abiLayoutCase{
		// 64-bit time ABI (_TIME_BITS=64 across all supported Linux architectures)
		{"sizeof(cTimeval)", unsafe.Sizeof(cTimeval{}), 16},
		{"sizeof(cTimespec)", unsafe.Sizeof(cTimespec{}), 16},

		// Generator regex layout rewrites (anonymous C unions and packed timings)
		{"sizeof(cV4L2BTTimings)", unsafe.Sizeof(cV4L2BTTimings{}), 124},
		{"offsetof(cV4L2BTTimings.Pixelclock)", unsafe.Offsetof(cV4L2BTTimings{}.Pixelclock), 16},
		{"sizeof(cV4L2DecoderCmd)", unsafe.Sizeof(cV4L2DecoderCmd{}), 72},
		{"offsetof(cV4L2DecoderCmd.Stop)", unsafe.Offsetof(cV4L2DecoderCmd{}.Stop), 8},
		{"sizeof(cV4L2EncoderCmd)", unsafe.Sizeof(cV4L2EncoderCmd{}), 40},
		{"offsetof(cV4L2EncoderCmd.Raw)", unsafe.Offsetof(cV4L2EncoderCmd{}.Raw), 8},
		{"sizeof(cMediaEntityDesc)", unsafe.Sizeof(cMediaEntityDesc{}), 256},
		{"offsetof(cMediaEntityDesc.Dev)", unsafe.Offsetof(cMediaEntityDesc{}.Dev), 72},

		// Union offsets shared across 32-bit and 64-bit Linux targets
		{"offsetof(cV4L2ExtControl.Value)", unsafe.Offsetof(cV4L2ExtControl{}.Value), 12},
		{"sizeof(cV4L2ExtControl)", unsafe.Sizeof(cV4L2ExtControl{}), 20},
		{"offsetof(cV4L2Event.U)", unsafe.Offsetof(cV4L2Event{}.U), 8},
		{"sizeof(cV4L2Event)", unsafe.Sizeof(cV4L2Event{}), 136},
		{"offsetof(cV4L2Buffer.M)", unsafe.Offsetof(cV4L2Buffer{}.M), 64},
		{"offsetof(cV4L2Plane.M)", unsafe.Offsetof(cV4L2Plane{}.M), 8},
	}

	if is64Bit {
		cases = append(cases,
			abiLayoutCase{"sizeof(cV4L2Format)", unsafe.Sizeof(cV4L2Format{}), 208},
			abiLayoutCase{"offsetof(cV4L2Format.Fmt)", unsafe.Offsetof(cV4L2Format{}.Fmt), 8},
			abiLayoutCase{"sizeof(cV4L2Buffer)", unsafe.Sizeof(cV4L2Buffer{}), 88},
			abiLayoutCase{"sizeof(cV4L2Plane)", unsafe.Sizeof(cV4L2Plane{}), 64},
		)
	} else {
		cases = append(cases,
			abiLayoutCase{"sizeof(cV4L2Format)", unsafe.Sizeof(cV4L2Format{}), 204},
			abiLayoutCase{"offsetof(cV4L2Format.Fmt)", unsafe.Offsetof(cV4L2Format{}.Fmt), 4},
			abiLayoutCase{"sizeof(cV4L2Buffer)", unsafe.Sizeof(cV4L2Buffer{}), 80},
			abiLayoutCase{"sizeof(cV4L2Plane)", unsafe.Sizeof(cV4L2Plane{}), 60},
		)
	}

	assertABILayouts(t, cases)

	if MetaFmtBCM2835ISPStats != 0x41545342 {
		t.Errorf("MetaFmtBCM2835ISPStats = 0x%08x, want 0x41545342", MetaFmtBCM2835ISPStats)
	}
}
