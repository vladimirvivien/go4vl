package v4l2

import (
	"strings"
	"testing"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

func TestMetaFormat_MockIoctl(t *testing.T) {
	var stored cV4L2MetaFormat
	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		f := (*cV4L2Format)(arg)
		if !IsMetaBufType(f.Type) {
			return sys.EINVAL
		}
		meta := (*cV4L2MetaFormat)(unsafe.Pointer(&f.Fmt[0]))
		switch req {
		case cVIDIOC_S_FMT:
			meta.Buffersize = 16384
			stored = *meta
			return 0
		case cVIDIOC_TRY_FMT:
			meta.Buffersize = 16384
			return 0
		case cVIDIOC_G_FMT:
			*meta = stored
			return 0
		default:
			return sys.ENOTTY
		}
	})

	// Reject non-metadata buffer types before syscall.
	if _, err := GetMetaFormat(3, BufTypeVideoCapture); err == nil {
		t.Fatal("expected error for BufTypeVideoCapture in GetMetaFormat")
	}
	if _, err := SetMetaFormat(3, BufTypeVideoOutput, MetaFormat{}); err == nil {
		t.Fatal("expected error for BufTypeVideoOutput in SetMetaFormat")
	}
	if _, err := TryMetaFormat(3, BufTypeVideoCaptureMPlane, MetaFormat{}); err == nil {
		t.Fatal("expected error for BufTypeVideoCaptureMPlane in TryMetaFormat")
	}

	reqFmt := MetaFormat{
		DataFormat:   MetaFmtBCM2835ISPStats,
		Width:        640,
		Height:       1,
		BytesPerLine: 640,
	}
	set, err := SetMetaFormat(3, BufTypeMetaCapture, reqFmt)
	if err != nil {
		t.Fatalf("SetMetaFormat() = %v", err)
	}
	if set.DataFormat != MetaFmtBCM2835ISPStats || set.BufferSize != 16384 {
		t.Fatalf("SetMetaFormat() = %+v", set)
	}

	got, err := GetMetaFormat(3, BufTypeMetaCapture)
	if err != nil {
		t.Fatalf("GetMetaFormat() = %v", err)
	}
	if got != set {
		t.Fatalf("GetMetaFormat() = %+v, want %+v", got, set)
	}

	tried, err := TryMetaFormat(3, BufTypeMetaOutput, MetaFormat{DataFormat: MetaFmtUVC})
	if err != nil {
		t.Fatalf("TryMetaFormat() = %v", err)
	}
	if tried.DataFormat != MetaFmtUVC || tried.BufferSize != 16384 {
		t.Fatalf("TryMetaFormat() = %+v", tried)
	}

	if !strings.Contains(set.String(), "BCM2835") || !strings.Contains(tried.String(), "UVCH") {
		t.Errorf("unexpected MetaFormat.String(): %q / %q", set.String(), tried.String())
	}
}
