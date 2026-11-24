package v4l2

import (
	"encoding/binary"
	"testing"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// mockGetExtCtrls answers VIDIOC_G_EXT_CTRLS with the 64-bit value of each
// requested control, as a driver fills value64.
func mockGetExtCtrls(t *testing.T, values map[CtrlID]int64) {
	t.Helper()
	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		if req != cVIDIOC_G_EXT_CTRLS {
			return sys.ENOTTY
		}
		c := (*cV4L2ExtControls)(arg)
		ctrls := unsafe.Slice(c.Controls, c.Count)
		for i := range ctrls {
			ctrl := &ctrls[i]
			v, ok := values[ctrl.Id]
			if !ok {
				c.Error_idx = uint32(i)
				return sys.EINVAL
			}
			b := (*[8]byte)(unsafe.Pointer(&ctrl.Value))
			binary.NativeEndian.PutUint64(b[:], uint64(v))
		}
		return 0
	})
}

func TestGetExtControlValue64(t *testing.T) {
	const pixelRate = int64(1_188_000_000) // above 2^31
	mockGetExtCtrls(t, map[CtrlID]int64{CtrlImgProcPixelRate: pixelRate})

	got, err := GetExtControlValue64(3, CtrlImgProcPixelRate)
	if err != nil {
		t.Fatalf("GetExtControlValue64() error = %v", err)
	}
	if got != pixelRate {
		t.Fatalf("GetExtControlValue64() = %d, want %d", got, pixelRate)
	}
}

func TestGetExtControls_Value64(t *testing.T) {
	const big = int64(-1) << 40
	mockGetExtCtrls(t, map[CtrlID]int64{CtrlImgProcPixelRate: big, CtrlBrightness: 42})

	ctrls := NewExtControls()
	ctrls.AddValue64(CtrlImgProcPixelRate, 0)
	ctrls.Add(NewExtControl(CtrlBrightness))
	if err := GetExtControls(3, ctrls); err != nil {
		t.Fatalf("GetExtControls() error = %v", err)
	}

	got := ctrls.GetControls()
	if v := got[0].GetValue64(); v != big {
		t.Errorf("64-bit control = %d, want %d", v, big)
	}
	// A control added without a value is still read as an int32.
	if v, ok := got[1].Value.(int32); !ok || v != 42 {
		t.Errorf("32-bit control = %#v, want int32(42)", got[1].Value)
	}
}

func TestGetExtControlValue_TruncatesValue64(t *testing.T) {
	// GetExtControlValue reads an int32: the low 32 bits of a 64-bit control.
	mockGetExtCtrls(t, map[CtrlID]int64{CtrlImgProcPixelRate: 1<<32 | 7})

	got, err := GetExtControlValue(3, CtrlImgProcPixelRate)
	if err != nil {
		t.Fatalf("GetExtControlValue() error = %v", err)
	}
	if got != 7 {
		t.Fatalf("GetExtControlValue() = %d, want 7", got)
	}
}

func TestQueryExtControlInfo_Range64(t *testing.T) {
	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		if req != cVIDIOC_QUERY_EXT_CTRL {
			return sys.ENOTTY
		}
		q := (*cV4L2QueryExtCtrl)(arg)
		q.Type = uint32(CtrlTypeInt64)
		q.Minimum = 1
		q.Maximum = 2_376_000_000 // above 2^31
		q.Step = 1
		q.Default_value = 1_188_000_000
		q.Flags = uint32(CtrlFlagReadOnly)
		return 0
	})

	c, err := QueryExtControlInfo(3, CtrlImgProcPixelRate)
	if err != nil {
		t.Fatalf("QueryExtControlInfo() error = %v", err)
	}
	if c.Minimum64 != 1 || c.Maximum64 != 2_376_000_000 || c.Step64 != 1 || c.Default64 != 1_188_000_000 {
		t.Errorf("64-bit range = %d..%d step %d default %d, want 1..2376000000 step 1 default 1188000000",
			c.Minimum64, c.Maximum64, c.Step64, c.Default64)
	}
	if c.Default != 1_188_000_000 {
		t.Errorf("Default = %d, want 1188000000 (fits in 32 bits)", c.Default)
	}
}
