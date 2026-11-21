package v4l2

import "testing"

func TestMBusFrameFormatConversion(t *testing.T) {
	in := MBusFrameFormat{
		Width:        4056,
		Height:       3040,
		Code:         MBusFmtSRGGB12_1X12,
		Field:        FieldNone,
		Colorspace:   ColorspaceRaw,
		YcbcrEnc:     uint16(YCbCrEncodingDefault),
		Quantization: uint16(QuantizationFullRange),
		XferFunc:     uint16(XferFuncNone),
		Flags:        0,
	}
	roundTrip := fromCMBusFrameFmt(toCMBusFrameFmt(in))
	if roundTrip != in {
		t.Fatalf("round-trip mismatch:\n got: %+v\nwant: %+v", roundTrip, in)
	}

	if name := MBusFormatNames[MBusFmtSRGGB12_1X12]; name != "SRGGB12_1X12" {
		t.Fatalf("MBusFormatNames[SRGGB12_1X12] = %q, want SRGGB12_1X12", name)
	}
	if SubdevFormatTry != 0 || SubdevFormatActive != 1 {
		t.Fatalf("unexpected SubdevFormatWhich values: try=%d active=%d", SubdevFormatTry, SubdevFormatActive)
	}
}
