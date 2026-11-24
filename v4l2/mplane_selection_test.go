package v4l2

import (
	"testing"
)

func TestIsMultiPlanar(t *testing.T) {
	tests := []struct {
		name    string
		bufType BufType
		want    bool
	}{
		{"VideoCapture", BufTypeVideoCapture, false},
		{"VideoOutput", BufTypeVideoOutput, false},
		{"Overlay", BufTypeOverlay, false},
		{"MetaCapture", BufTypeMetaCapture, false},
		{"MetaOutput", BufTypeMetaOutput, false},
		{"VideoCaptureMPlane", BufTypeVideoCaptureMPlane, true},
		{"VideoOutputMPlane", BufTypeVideoOutputMPlane, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsMultiPlanar(tc.bufType); got != tc.want {
				t.Errorf("IsMultiPlanar(%d) = %v, want %v", tc.bufType, got, tc.want)
			}
		})
	}
}

func TestPixFormatMPlaneRoundTrip(t *testing.T) {
	in := PixFormatMPlane{
		Width:        1920,
		Height:       1080,
		PixelFormat:  PixelFmtNV12M,
		Field:        FieldNone,
		Colorspace:   ColorspaceREC709,
		NumPlanes:    2,
		Flags:        1,
		YcbcrEnc:     uint8(YCbCrEncoding709),
		Quantization: uint8(QuantizationLimitedRange),
		XferFunc:     uint8(XferFunc709),
	}
	in.PlaneFmt[0] = PlanePixFormat{SizeImage: 1920 * 1080, BytesPerLine: 1920}
	in.PlaneFmt[1] = PlanePixFormat{SizeImage: 1920 * 540, BytesPerLine: 1920}

	cFmt := pixFormatMPlaneToC(in)
	got := pixFormatMPlaneFromC(cFmt)

	if got.Width != in.Width || got.Height != in.Height || got.PixelFormat != in.PixelFormat {
		t.Fatalf("basic fields mismatch: got %+v, want %+v", got, in)
	}
	if got.NumPlanes != 2 {
		t.Fatalf("NumPlanes = %d, want 2", got.NumPlanes)
	}
	if got.PlaneFmt[0] != in.PlaneFmt[0] || got.PlaneFmt[1] != in.PlaneFmt[1] {
		t.Fatalf("PlaneFmt mismatch: got %+v / %+v", got.PlaneFmt[0], got.PlaneFmt[1])
	}
	if got.YcbcrEnc != in.YcbcrEnc || got.Quantization != in.Quantization || got.XferFunc != in.XferFunc {
		t.Fatalf("colorspace fields mismatch: got %+v, want %+v", got, in)
	}
}

func TestBufferInfoAndPlaneEncodeDecode(t *testing.T) {
	// Single-planar MMAP
	var cBuf cV4L2Buffer
	cBuf.Index = 2
	cBuf.Type = BufTypeVideoCapture
	cBuf.Memory = IOTypeMMAP
	cBuf.Length = 614400
	cBuf.Bytesused = 614400
	encodePlaneInfo(cBuf.M[:], IOTypeMMAP, PlaneInfo{MemOffset: 0x10000})
	b := makeBuffer(cBuf)
	if b.Info.Offset != 0x10000 || b.Info.UserPtr != 0 || b.Info.FD != 0 {
		t.Fatalf("MMAP BufferInfo mismatch: %+v", b.Info)
	}

	// Single-planar DMABUF
	cBuf.Memory = IOTypeDMABuf
	encodePlaneInfo(cBuf.M[:], IOTypeDMABuf, PlaneInfo{FD: 42})
	b = makeBuffer(cBuf)
	if b.Info.FD != 42 || b.Info.Offset != 0 {
		t.Fatalf("DMABUF BufferInfo mismatch: %+v", b.Info)
	}

	// Single-planar USERPTR
	cBuf.Memory = IOTypeUserPtr
	encodePlaneInfo(cBuf.M[:], IOTypeUserPtr, PlaneInfo{UserPtr: 0xdeadbeef})
	b = makeBuffer(cBuf)
	if b.Info.UserPtr != 0xdeadbeef {
		t.Fatalf("USERPTR BufferInfo mismatch: %+v", b.Info)
	}

	// Multi-planar MMAP
	cBuf.Type = BufTypeVideoCaptureMPlane
	cBuf.Memory = IOTypeMMAP
	cBuf.Length = 2
	cPlanes := []cV4L2Plane{
		{
			Bytesused: 1920 * 1080,
			Length:    1920 * 1080,
			Offset:    0,
		},
		{
			Bytesused: 1920 * 540,
			Length:    1920 * 540,
			Offset:    64,
		},
	}
	encodePlaneInfo(cPlanes[0].M[:], IOTypeMMAP, PlaneInfo{MemOffset: 0x0})
	encodePlaneInfo(cPlanes[1].M[:], IOTypeMMAP, PlaneInfo{MemOffset: 0x200000})
	mb := makeBufferWithPlanes(cBuf, cPlanes)
	if len(mb.Planes) != 2 {
		t.Fatalf("len(Planes) = %d, want 2", len(mb.Planes))
	}
	if mb.Planes[0].Length != 1920*1080 || mb.Planes[1].Info.MemOffset != 0x200000 || mb.Planes[1].DataOffset != 64 {
		t.Fatalf("unexpected decoded planes: %+v", mb.Planes)
	}
	if mb.Info.Planes != &mb.Planes[0] {
		t.Fatalf("expected Info.Planes to point to &Planes[0]")
	}
}

func TestMPlaneValidationErrors(t *testing.T) {
	if _, err := GetPixFormatMPlane(0, BufTypeVideoCapture); err == nil {
		t.Error("expected error for non-mplane buffer type in GetPixFormatMPlane")
	}
	if _, err := SetPixFormatMPlane(0, BufTypeVideoOutput, PixFormatMPlane{}); err == nil {
		t.Error("expected error for non-mplane buffer type in SetPixFormatMPlane")
	}
	if _, err := TryPixFormatMPlane(0, BufTypeVideoCapture, PixFormatMPlane{}); err == nil {
		t.Error("expected error for non-mplane buffer type in TryPixFormatMPlane")
	}
	if _, err := QueryBuffer(0, BufTypeVideoCaptureMPlane, IOTypeMMAP, 0, 0); err == nil {
		t.Error("expected error for 0 planes in QueryBuffer")
	}
	if _, err := QueryBuffer(0, BufTypeVideoCaptureMPlane, IOTypeMMAP, 0, MaxPlanes+1); err == nil {
		t.Error("expected error for >MaxPlanes in QueryBuffer")
	}
	if _, err := DequeueBufferMPlane(0, IOTypeMMAP, BufTypeVideoCapture, 1); err == nil {
		t.Error("expected error for non-mplane buffer type in DequeueBufferMPlane")
	}
	if _, err := DequeueBufferMPlane(0, IOTypeMMAP, BufTypeVideoCaptureMPlane, 0); err == nil {
		t.Error("expected error for 0 planes in DequeueBufferMPlane")
	}
	if _, err := QueueBufferWith(0, Buffer{Type: BufTypeVideoCaptureMPlane, Memory: IOTypeMMAP}); err == nil {
		t.Error("expected error for empty Planes in QueueBufferWith")
	}
	if _, err := CreateBuffers(0, 2, IOTypeMMAP, BufTypeVideoCapture, nil, nil); err == nil {
		t.Error("expected error when pixFmt is nil for single-planar CreateBuffers")
	}
	if _, err := CreateBuffers(0, 2, IOTypeMMAP, BufTypeVideoCaptureMPlane, nil, nil); err == nil {
		t.Error("expected error when pixFmtMPlane is nil for multi-planar CreateBuffers")
	}
}

func TestSelectionConstants(t *testing.T) {
	if SelTargetCrop == SelTargetCompose {
		t.Error("SelTargetCrop and SelTargetCompose must differ")
	}
	if SelTargetCropBounds == 0 || SelTargetNativeSize == 0 {
		t.Error("selection target constants must be non-zero")
	}
	if SelFlagGE == 0 || SelFlagLE == 0 || SelFlagKeepConfig == 0 {
		t.Error("selection flag constants must be non-zero")
	}
}
