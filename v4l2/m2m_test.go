package v4l2

import (
	"testing"
)

func TestM2MQueueValidation(t *testing.T) {
	sess := &M2MSession{
		cap: Capability{Capabilities: CapVideoMem2MemMPlane | CapStreaming},
	}

	if err := sess.SetupOutput(M2MQueueConfig{BufType: BufTypeVideoCapture}); err == nil {
		t.Fatalf("SetupOutput with capture BufType should fail")
	}
	if err := sess.SetupCapture(M2MQueueConfig{BufType: BufTypeVideoOutput}); err == nil {
		t.Fatalf("SetupCapture with output BufType should fail")
	}
	if err := sess.SetupOutput(M2MQueueConfig{
		BufType:   BufTypeVideoOutputMPlane,
		NumPlanes: MaxPlanes + 1,
	}); err == nil {
		t.Fatalf("SetupOutput with NumPlanes > MaxPlanes should fail")
	}
	if err := sess.StartOutput(); err == nil {
		t.Fatalf("StartOutput before SetupOutput should fail")
	}
	if err := sess.StartCapture(); err == nil {
		t.Fatalf("StartCapture before SetupCapture should fail")
	}
	// Stopping inactive queues should be no-op.
	if err := sess.StopOutput(); err != nil {
		t.Fatalf("StopOutput on idle queue: %v", err)
	}
	if err := sess.StopCapture(); err != nil {
		t.Fatalf("StopCapture on idle queue: %v", err)
	}
}
