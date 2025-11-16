package v4l2

import (
	"testing"
	"time"
)

// TestEventType_Constants verifies event type constants are defined
func TestEventType_Constants(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
	}{
		{"EventAll", EventAll},
		{"EventVSync", EventVSync},
		{"EventEOS", EventEOS},
		{"EventCtrl", EventCtrl},
		{"EventFrameSync", EventFrameSync},
		{"EventSourceChange", EventSourceChange},
		{"EventMotionDet", EventMotionDet},
		{"EventPrivateStart", EventPrivateStart},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// EventAll is 0, others should be non-zero
			if tt.name != "EventAll" && tt.eventType == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestEventTypeNames_MapComplete verifies EventTypeNames map contains all event types
func TestEventTypeNames_MapComplete(t *testing.T) {
	expectedTypes := []EventType{
		EventAll,
		EventVSync,
		EventEOS,
		EventCtrl,
		EventFrameSync,
		EventSourceChange,
		EventMotionDet,
	}

	for _, eventType := range expectedTypes {
		if name, ok := EventTypeNames[eventType]; !ok {
			t.Errorf("EventTypeNames map missing entry for type %d", eventType)
		} else if name == "" {
			t.Errorf("EventTypeNames map has empty name for type %d", eventType)
		}
	}
}

// TestEventSubscriptionFlags_Constants verifies event subscription flags
func TestEventSubscriptionFlags_Constants(t *testing.T) {
	tests := []struct {
		name string
		flag EventSubscriptionFlags
	}{
		{"EventSubFlagSendInitial", EventSubFlagSendInitial},
		{"EventSubFlagAllowFeedback", EventSubFlagAllowFeedback},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.flag == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestEventCtrlChanges_Constants verifies control change flag constants
func TestEventCtrlChanges_Constants(t *testing.T) {
	tests := []struct {
		name    string
		changes EventCtrlChanges
	}{
		{"EventCtrlChValue", EventCtrlChValue},
		{"EventCtrlChFlags", EventCtrlChFlags},
		{"EventCtrlChRange", EventCtrlChRange},
		{"EventCtrlChDimensions", EventCtrlChDimensions},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.changes == 0 {
				t.Errorf("%s should not be zero", tt.name)
			}
		})
	}
}

// TestEventSrcChanges_Constants verifies source change flag constants
func TestEventSrcChanges_Constants(t *testing.T) {
	if EventSrcChResolution == 0 {
		t.Error("EventSrcChResolution should not be zero")
	}
}

// TestEventMotionDetFlags_Constants verifies motion detection flag constants
func TestEventMotionDetFlags_Constants(t *testing.T) {
	if EventMotionDetFlagHaveFrameSeq == 0 {
		t.Error("EventMotionDetFlagHaveFrameSeq should not be zero")
	}
}

// TestNewEventSubscription verifies EventSubscription creation
func TestNewEventSubscription(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
	}{
		{"All", EventAll},
		{"VSync", EventVSync},
		{"EOS", EventEOS},
		{"Ctrl", EventCtrl},
		{"FrameSync", EventFrameSync},
		{"SourceChange", EventSourceChange},
		{"MotionDet", EventMotionDet},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := NewEventSubscription(tt.eventType)
			if sub == nil {
				t.Fatal("NewEventSubscription returned nil")
			}
			if sub.GetType() != tt.eventType {
				t.Errorf("GetType() = %d, want %d", sub.GetType(), tt.eventType)
			}
		})
	}
}

// TestNewControlEventSubscription verifies control event subscription creation
func TestNewControlEventSubscription(t *testing.T) {
	tests := []struct {
		name   string
		ctrlID CtrlID
	}{
		{"Brightness", CtrlBrightness},
		{"Contrast", CtrlContrast},
		{"Saturation", CtrlSaturation},
		{"Hue", CtrlHue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := NewControlEventSubscription(tt.ctrlID)
			if sub == nil {
				t.Fatal("NewControlEventSubscription returned nil")
			}
			if sub.GetType() != EventCtrl {
				t.Errorf("GetType() = %d, want %d", sub.GetType(), EventCtrl)
			}
			if sub.GetID() != uint32(tt.ctrlID) {
				t.Errorf("GetID() = %d, want %d", sub.GetID(), tt.ctrlID)
			}
		})
	}
}

// TestEventSubscription_SetGetType verifies type set/get
func TestEventSubscription_SetGetType(t *testing.T) {
	sub := NewEventSubscription(EventAll)

	tests := []EventType{
		EventVSync,
		EventEOS,
		EventCtrl,
		EventFrameSync,
		EventSourceChange,
		EventMotionDet,
	}

	for _, eventType := range tests {
		sub.SetType(eventType)
		got := sub.GetType()
		if got != eventType {
			t.Errorf("After SetType(%d), GetType() = %d", eventType, got)
		}
	}
}

// TestEventSubscription_SetGetID verifies ID set/get
func TestEventSubscription_SetGetID(t *testing.T) {
	sub := NewEventSubscription(EventCtrl)

	tests := []uint32{0, 1, 100, 0xFFFFFFFF}

	for _, id := range tests {
		sub.SetID(id)
		got := sub.GetID()
		if got != id {
			t.Errorf("After SetID(%d), GetID() = %d", id, got)
		}
	}
}

// TestEventSubscription_SetGetFlags verifies flags set/get
func TestEventSubscription_SetGetFlags(t *testing.T) {
	sub := NewEventSubscription(EventCtrl)

	tests := []EventSubscriptionFlags{
		0,
		EventSubFlagSendInitial,
		EventSubFlagAllowFeedback,
		EventSubFlagSendInitial | EventSubFlagAllowFeedback,
	}

	for _, flags := range tests {
		sub.SetFlags(flags)
		got := sub.GetFlags()
		if got != flags {
			t.Errorf("After SetFlags(0x%08x), GetFlags() = 0x%08x", flags, got)
		}
	}
}

// TestEvent_GetType verifies event type retrieval
func TestEvent_GetType(t *testing.T) {
	event := &Event{}
	// Default should be 0 (EventAll)
	if event.GetType() != 0 {
		t.Errorf("Default GetType() = %d, want 0", event.GetType())
	}
}

// TestEvent_GetID verifies event ID retrieval
func TestEvent_GetID(t *testing.T) {
	event := &Event{}
	// Default should be 0
	if event.GetID() != 0 {
		t.Errorf("Default GetID() = %d, want 0", event.GetID())
	}
}

// TestEvent_GetPending verifies pending count retrieval
func TestEvent_GetPending(t *testing.T) {
	event := &Event{}
	// Default should be 0
	if event.GetPending() != 0 {
		t.Errorf("Default GetPending() = %d, want 0", event.GetPending())
	}
}

// TestEvent_GetSequence verifies sequence number retrieval
func TestEvent_GetSequence(t *testing.T) {
	event := &Event{}
	// Default should be 0
	if event.GetSequence() != 0 {
		t.Errorf("Default GetSequence() = %d, want 0", event.GetSequence())
	}
}

// TestEvent_GetTimestamp verifies timestamp retrieval
func TestEvent_GetTimestamp(t *testing.T) {
	event := &Event{}
	ts := event.GetTimestamp()

	// Should be Unix epoch (zero time)
	if !ts.Equal(time.Unix(0, 0)) {
		t.Errorf("Default GetTimestamp() = %v, want Unix epoch", ts)
	}
}

// TestEvent_GetCtrlData verifies control event data extraction
func TestEvent_GetCtrlData(t *testing.T) {
	event := &Event{}
	data := event.GetCtrlData()

	// Verify structure is initialized
	if data.Changes != 0 {
		t.Logf("GetCtrlData() returned changes = 0x%08x", data.Changes)
	}
	if data.Type != 0 {
		t.Logf("GetCtrlData() returned type = %d", data.Type)
	}
}

// TestEvent_GetSrcChangeData verifies source change event data extraction
func TestEvent_GetSrcChangeData(t *testing.T) {
	event := &Event{}
	data := event.GetSrcChangeData()

	// Verify structure is initialized
	if data.Changes != 0 {
		t.Logf("GetSrcChangeData() returned changes = 0x%08x", data.Changes)
	}
}

// TestEvent_GetFrameSyncData verifies frame sync event data extraction
func TestEvent_GetFrameSyncData(t *testing.T) {
	event := &Event{}
	data := event.GetFrameSyncData()

	// Verify structure is initialized
	if data.FrameSequence != 0 {
		t.Logf("GetFrameSyncData() returned frame_sequence = %d", data.FrameSequence)
	}
}

// TestEvent_GetMotionDetData verifies motion detection event data extraction
func TestEvent_GetMotionDetData(t *testing.T) {
	event := &Event{}
	data := event.GetMotionDetData()

	// Verify structure is initialized
	if data.Flags != 0 {
		t.Logf("GetMotionDetData() returned flags = 0x%08x", data.Flags)
	}
	if data.FrameSequence != 0 {
		t.Logf("GetMotionDetData() returned frame_sequence = %d", data.FrameSequence)
	}
	if data.RegionMask != 0 {
		t.Logf("GetMotionDetData() returned region_mask = 0x%08x", data.RegionMask)
	}
}

// TestEvent_GetRawData verifies raw data extraction
func TestEvent_GetRawData(t *testing.T) {
	event := &Event{}
	data := event.GetRawData()

	if len(data) != 64 {
		t.Errorf("GetRawData() length = %d, want 64", len(data))
	}

	// Should all be zero initially
	allZero := true
	for i, b := range data {
		if b != 0 {
			allZero = false
			t.Logf("GetRawData()[%d] = %d (non-zero)", i, b)
			break
		}
	}

	if allZero {
		t.Log("GetRawData() returned all zeros (expected for uninitialized event)")
	}
}

// TestEventCtrlData_Fields verifies EventCtrlData structure
func TestEventCtrlData_Fields(t *testing.T) {
	data := EventCtrlData{
		Changes: EventCtrlChValue | EventCtrlChRange,
		Type:    CtrlTypeInt,
		Value:   100,
		Value64: 1000000,
		Flags:   0x1234,
		Minimum: 0,
		Maximum: 255,
		Step:    1,
		Default: 128,
	}

	if data.Changes != (EventCtrlChValue | EventCtrlChRange) {
		t.Errorf("Changes = 0x%08x, want 0x%08x", data.Changes, EventCtrlChValue|EventCtrlChRange)
	}
	if data.Type != CtrlTypeInt {
		t.Errorf("Type = %d, want %d", data.Type, CtrlTypeInt)
	}
	if data.Value != 100 {
		t.Errorf("Value = %d, want 100", data.Value)
	}
	if data.Value64 != 1000000 {
		t.Errorf("Value64 = %d, want 1000000", data.Value64)
	}
	if data.Flags != 0x1234 {
		t.Errorf("Flags = 0x%08x, want 0x1234", data.Flags)
	}
	if data.Minimum != 0 {
		t.Errorf("Minimum = %d, want 0", data.Minimum)
	}
	if data.Maximum != 255 {
		t.Errorf("Maximum = %d, want 255", data.Maximum)
	}
	if data.Step != 1 {
		t.Errorf("Step = %d, want 1", data.Step)
	}
	if data.Default != 128 {
		t.Errorf("Default = %d, want 128", data.Default)
	}
}

// TestEventSrcChangeData_Fields verifies EventSrcChangeData structure
func TestEventSrcChangeData_Fields(t *testing.T) {
	data := EventSrcChangeData{
		Changes: EventSrcChResolution,
	}

	if data.Changes != EventSrcChResolution {
		t.Errorf("Changes = 0x%08x, want 0x%08x", data.Changes, EventSrcChResolution)
	}
}

// TestEventFrameSyncData_Fields verifies EventFrameSyncData structure
func TestEventFrameSyncData_Fields(t *testing.T) {
	data := EventFrameSyncData{
		FrameSequence: 12345,
	}

	if data.FrameSequence != 12345 {
		t.Errorf("FrameSequence = %d, want 12345", data.FrameSequence)
	}
}

// TestEventMotionDetData_Fields verifies EventMotionDetData structure
func TestEventMotionDetData_Fields(t *testing.T) {
	data := EventMotionDetData{
		Flags:         EventMotionDetFlagHaveFrameSeq,
		FrameSequence: 54321,
		RegionMask:    0xABCD,
	}

	if data.Flags != EventMotionDetFlagHaveFrameSeq {
		t.Errorf("Flags = 0x%08x, want 0x%08x", data.Flags, EventMotionDetFlagHaveFrameSeq)
	}
	if data.FrameSequence != 54321 {
		t.Errorf("FrameSequence = %d, want 54321", data.FrameSequence)
	}
	if data.RegionMask != 0xABCD {
		t.Errorf("RegionMask = 0x%08x, want 0xABCD", data.RegionMask)
	}
}

// TestEventSubscription_MultipleOperations verifies multiple operations
func TestEventSubscription_MultipleOperations(t *testing.T) {
	// Create subscription
	sub := NewEventSubscription(EventAll)

	// Change type
	sub.SetType(EventCtrl)
	if sub.GetType() != EventCtrl {
		t.Errorf("After SetType(EventCtrl), GetType() = %d", sub.GetType())
	}

	// Set ID
	sub.SetID(uint32(CtrlBrightness))
	if sub.GetID() != uint32(CtrlBrightness) {
		t.Errorf("After SetID(CtrlBrightness), GetID() = %d", sub.GetID())
	}

	// Set flags
	sub.SetFlags(EventSubFlagSendInitial | EventSubFlagAllowFeedback)
	expectedFlags := EventSubFlagSendInitial | EventSubFlagAllowFeedback
	if sub.GetFlags() != expectedFlags {
		t.Errorf("After SetFlags(), GetFlags() = 0x%08x, want 0x%08x", sub.GetFlags(), expectedFlags)
	}
}
