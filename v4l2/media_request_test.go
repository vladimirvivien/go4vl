package v4l2

import (
	"context"
	"testing"
	"time"
)

func TestMediaRequestClosedSafety(t *testing.T) {
	req := &MediaRequest{fd: -1}
	if req.Fd() != -1 {
		t.Fatalf("Fd() = %d, want -1", req.Fd())
	}
	if err := req.Queue(); err == nil {
		t.Fatalf("Queue() on closed request should fail")
	}
	if err := req.Reinit(); err == nil {
		t.Fatalf("Reinit() on closed request should fail")
	}
	if _, err := req.Wait(context.Background(), 10*time.Millisecond); err == nil {
		t.Fatalf("Wait() on closed request should fail")
	}
	if err := req.Close(); err != nil {
		t.Fatalf("Close() on already-closed request should be no-op, got %v", err)
	}
}

func TestExtControlsRequestFDMarshaling(t *testing.T) {
	ctrls := NewExtControls()
	ctrls.SetRequestFD(17)
	ctrls.AddValue(CtrlBrightness, 100)

	if ctrls.GetRequestFD() != 17 {
		t.Fatalf("GetRequestFD() = %d, want 17", ctrls.GetRequestFD())
	}
	if ctrls.GetClass() != CtrlWhichRequestVal {
		t.Fatalf("GetClass() = 0x%x, want 0x%x", ctrls.GetClass(), CtrlWhichRequestVal)
	}

	cCtrls, cleanup, err := prepareForIoctl(ctrls)
	if err != nil {
		t.Fatalf("prepareForIoctl: %v", err)
	}
	defer cleanup()

	if cCtrls.Request_fd != 17 {
		t.Fatalf("cCtrls.Request_fd = %d, want 17", cCtrls.Request_fd)
	}
	if cCtrls.Ctrl_class != CtrlWhichRequestVal {
		t.Fatalf("cCtrls.Ctrl_class = 0x%x, want 0x%x", cCtrls.Ctrl_class, CtrlWhichRequestVal)
	}
}
