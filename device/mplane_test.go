package device

import (
	"bytes"
	"context"
	"testing"

	"github.com/vladimirvivien/go4vl/v4l2"
)

func TestWithPixFormatMPlane(t *testing.T) {
	var cfg config
	want := v4l2.PixFormatMPlane{
		Width:       1920,
		Height:      1080,
		PixelFormat: v4l2.PixelFmtNV12M,
		NumPlanes:   2,
	}
	WithPixFormatMPlane(want)(&cfg)
	if cfg.pixFormatMPlane == nil || *cfg.pixFormatMPlane != want {
		t.Fatalf("pixFormatMPlane = %+v, want %+v", cfg.pixFormatMPlane, want)
	}
}

func TestFrame_CloneMPlane(t *testing.T) {
	y := []byte{1, 2, 3, 4}
	uv := []byte{10, 20}
	orig := &Frame{
		Data:     y,
		Planes:   [][]byte{y, uv},
		Sequence: 7,
	}

	cloned := orig.Clone()
	if len(cloned.Planes) != 2 {
		t.Fatalf("cloned Planes len = %d, want 2", len(cloned.Planes))
	}
	if !bytes.Equal(cloned.Planes[0], y) || !bytes.Equal(cloned.Planes[1], uv) {
		t.Fatalf("cloned plane contents mismatch")
	}
	if !bytes.Equal(cloned.Data, y) {
		t.Fatalf("cloned Data mismatch")
	}

	// Mutating originals must not affect clone.
	y[0] = 99
	uv[0] = 88
	if cloned.Planes[0][0] != 1 || cloned.Planes[1][0] != 10 {
		t.Fatalf("cloned Planes share backing memory with original")
	}
}

func TestDevice_WriteFramePlanesValidation(t *testing.T) {
	d := &Device{
		bufType: v4l2.BufTypeVideoCaptureMPlane,
		config:  config{ioMethod: IOMethodStreaming},
	}
	if err := d.WriteFramePlanes(context.Background(), [][]byte{{1}}); err == nil {
		t.Fatalf("expected error for capture mplane device")
	}

	dOut := &Device{
		bufType:   v4l2.BufTypeVideoOutputMPlane,
		numPlanes: 2,
		config:    config{ioMethod: IOMethodStreaming},
	}
	if err := dOut.WriteFramePlanes(context.Background(), [][]byte{{1}}); err == nil {
		t.Fatalf("expected error before Start()")
	}
}
