package v4l2

// events.go provides V4L2 event subscription and handling support.
//
// Events allow applications to be notified of changes in device state, such as:
// - Control value changes (V4L2_EVENT_CTRL)
// - Vertical sync events (V4L2_EVENT_VSYNC)
// - End of stream (V4L2_EVENT_EOS)
// - Frame sync events (V4L2_EVENT_FRAME_SYNC)
// - Source change events (V4L2_EVENT_SOURCE_CHANGE)
// - Motion detection events (V4L2_EVENT_MOTION_DET)
//
// Applications can subscribe to events using VIDIOC_SUBSCRIBE_EVENT,
// dequeue events using VIDIOC_DQEVENT, and unsubscribe using VIDIOC_UNSUBSCRIBE_EVENT.
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-subscribe-event.html
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-dqevent.html

/*
#include <linux/videodev2.h>
#include <string.h>
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

// EventType represents the type of V4L2 event
type EventType = uint32

// Event type constants
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
const (
	EventAll          EventType = C.V4L2_EVENT_ALL           // All events
	EventVSync        EventType = C.V4L2_EVENT_VSYNC         // Vertical sync
	EventEOS          EventType = C.V4L2_EVENT_EOS           // End of stream
	EventCtrl         EventType = C.V4L2_EVENT_CTRL          // Control changed
	EventFrameSync    EventType = C.V4L2_EVENT_FRAME_SYNC    // Frame sync
	EventSourceChange EventType = C.V4L2_EVENT_SOURCE_CHANGE // Source resolution/format changed
	EventMotionDet    EventType = C.V4L2_EVENT_MOTION_DET    // Motion detected
	EventPrivateStart EventType = C.V4L2_EVENT_PRIVATE_START // Start of driver-specific events
)

// Event type name mapping
var EventTypeNames = map[EventType]string{
	EventAll:          "All",
	EventVSync:        "VSync",
	EventEOS:          "End of Stream",
	EventCtrl:         "Control",
	EventFrameSync:    "Frame Sync",
	EventSourceChange: "Source Change",
	EventMotionDet:    "Motion Detection",
}

// EventSubscriptionFlags
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type EventSubscriptionFlags = uint32

const (
	EventSubFlagSendInitial    EventSubscriptionFlags = C.V4L2_EVENT_SUB_FL_SEND_INITIAL     // Send initial event on subscribe
	EventSubFlagAllowFeedback  EventSubscriptionFlags = C.V4L2_EVENT_SUB_FL_ALLOW_FEEDBACK   // Allow feedback events
)

// EventCtrlChanges represents changes to a control
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type EventCtrlChanges = uint32

const (
	EventCtrlChValue      EventCtrlChanges = C.V4L2_EVENT_CTRL_CH_VALUE      // Control value changed
	EventCtrlChFlags      EventCtrlChanges = C.V4L2_EVENT_CTRL_CH_FLAGS      // Control flags changed
	EventCtrlChRange      EventCtrlChanges = C.V4L2_EVENT_CTRL_CH_RANGE      // Control range changed
	EventCtrlChDimensions EventCtrlChanges = C.V4L2_EVENT_CTRL_CH_DIMENSIONS // Control dimensions changed
)

// EventSrcChanges represents source change types
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type EventSrcChanges = uint32

const (
	EventSrcChResolution EventSrcChanges = C.V4L2_EVENT_SRC_CH_RESOLUTION // Resolution changed
)

// EventMotionDetFlags represents motion detection event flags
type EventMotionDetFlags = uint32

const (
	EventMotionDetFlagHaveFrameSeq EventMotionDetFlags = C.V4L2_EVENT_MD_FL_HAVE_FRAME_SEQ // Have frame sequence number
)

// EventSubscription represents an event subscription (v4l2_event_subscription).
//
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type EventSubscription struct {
	v4l2EventSubscription C.struct_v4l2_event_subscription
}

// NewEventSubscription creates a new event subscription.
func NewEventSubscription(eventType EventType) *EventSubscription {
	es := &EventSubscription{}
	es.v4l2EventSubscription._type = C.__u32(eventType)
	return es
}

// NewControlEventSubscription creates an event subscription for a specific control.
func NewControlEventSubscription(ctrlID CtrlID) *EventSubscription {
	es := NewEventSubscription(EventCtrl)
	es.v4l2EventSubscription.id = C.__u32(ctrlID)
	return es
}

// GetType returns the event type.
func (es *EventSubscription) GetType() EventType {
	return EventType(es.v4l2EventSubscription._type)
}

// SetType sets the event type.
func (es *EventSubscription) SetType(eventType EventType) {
	es.v4l2EventSubscription._type = C.__u32(eventType)
}

// GetID returns the event ID (e.g., control ID for control events).
func (es *EventSubscription) GetID() uint32 {
	return uint32(es.v4l2EventSubscription.id)
}

// SetID sets the event ID.
func (es *EventSubscription) SetID(id uint32) {
	es.v4l2EventSubscription.id = C.__u32(id)
}

// GetFlags returns the subscription flags.
func (es *EventSubscription) GetFlags() EventSubscriptionFlags {
	return EventSubscriptionFlags(es.v4l2EventSubscription.flags)
}

// SetFlags sets the subscription flags.
func (es *EventSubscription) SetFlags(flags EventSubscriptionFlags) {
	es.v4l2EventSubscription.flags = C.__u32(flags)
}

// EventCtrlData represents control event data.
type EventCtrlData struct {
	Changes EventCtrlChanges // What changed
	Type    CtrlType         // Control type
	Value   int32            // Current value (for 32-bit controls)
	Value64 int64            // Current value (for 64-bit controls)
	Flags   uint32           // Control flags
	Minimum int32            // Minimum value
	Maximum int32            // Maximum value
	Step    int32            // Step value
	Default int32            // Default value
}

// EventSrcChangeData represents source change event data.
type EventSrcChangeData struct {
	Changes EventSrcChanges // What changed
}

// EventFrameSyncData represents frame sync event data.
type EventFrameSyncData struct {
	FrameSequence uint32 // Frame sequence number
}

// EventMotionDetData represents motion detection event data.
type EventMotionDetData struct {
	Flags         EventMotionDetFlags // Flags
	FrameSequence uint32              // Frame sequence number (if flag set)
	RegionMask    uint32              // Bitmask of regions with detected motion
}

// Event represents a V4L2 event (v4l2_event).
//
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type Event struct {
	v4l2Event C.struct_v4l2_event
}

// GetType returns the event type.
func (e *Event) GetType() EventType {
	return EventType(e.v4l2Event._type)
}

// GetID returns the event ID (e.g., control ID for control events).
func (e *Event) GetID() uint32 {
	return uint32(e.v4l2Event.id)
}

// GetPending returns the number of pending events of this type.
func (e *Event) GetPending() uint32 {
	return uint32(e.v4l2Event.pending)
}

// GetSequence returns the event sequence number.
func (e *Event) GetSequence() uint32 {
	return uint32(e.v4l2Event.sequence)
}

// GetTimestamp returns the event timestamp.
func (e *Event) GetTimestamp() time.Time {
	ts := e.v4l2Event.timestamp
	return time.Unix(int64(ts.tv_sec), int64(ts.tv_nsec))
}

// GetCtrlData returns the control event data (valid if Type is EventCtrl).
func (e *Event) GetCtrlData() EventCtrlData {
	// Access the ctrl union member through unsafe pointer
	ctrlPtr := (*C.struct_v4l2_event_ctrl)(unsafe.Pointer(&e.v4l2Event.u[0]))

	// The value field is a union of value/value64
	value := int32(*(*C.__s32)(unsafe.Pointer(&ctrlPtr.anon0[0])))
	value64 := int64(*(*C.__s64)(unsafe.Pointer(&ctrlPtr.anon0[0])))

	return EventCtrlData{
		Changes: EventCtrlChanges(ctrlPtr.changes),
		Type:    CtrlType(ctrlPtr._type),
		Value:   value,
		Value64: value64,
		Flags:   uint32(ctrlPtr.flags),
		Minimum: int32(ctrlPtr.minimum),
		Maximum: int32(ctrlPtr.maximum),
		Step:    int32(ctrlPtr.step),
		Default: int32(ctrlPtr.default_value),
	}
}

// GetSrcChangeData returns the source change event data (valid if Type is EventSourceChange).
func (e *Event) GetSrcChangeData() EventSrcChangeData {
	srcPtr := (*C.struct_v4l2_event_src_change)(unsafe.Pointer(&e.v4l2Event.u[0]))
	return EventSrcChangeData{
		Changes: EventSrcChanges(srcPtr.changes),
	}
}

// GetFrameSyncData returns the frame sync event data (valid if Type is EventFrameSync).
func (e *Event) GetFrameSyncData() EventFrameSyncData {
	fsPtr := (*C.struct_v4l2_event_frame_sync)(unsafe.Pointer(&e.v4l2Event.u[0]))
	return EventFrameSyncData{
		FrameSequence: uint32(fsPtr.frame_sequence),
	}
}

// GetMotionDetData returns the motion detection event data (valid if Type is EventMotionDet).
func (e *Event) GetMotionDetData() EventMotionDetData {
	mdPtr := (*C.struct_v4l2_event_motion_det)(unsafe.Pointer(&e.v4l2Event.u[0]))
	return EventMotionDetData{
		Flags:         EventMotionDetFlags(mdPtr.flags),
		FrameSequence: uint32(mdPtr.frame_sequence),
		RegionMask:    uint32(mdPtr.region_mask),
	}
}

// GetRawData returns the raw event data as a byte slice (for unknown event types).
func (e *Event) GetRawData() []byte {
	data := make([]byte, 64)
	for i := 0; i < 64; i++ {
		data[i] = byte(e.v4l2Event.u[i])
	}
	return data
}

// SubscribeEvent subscribes to an event type.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-subscribe-event.html
func SubscribeEvent(fd uintptr, sub *EventSubscription) error {
	if err := send(fd, C.VIDIOC_SUBSCRIBE_EVENT, uintptr(unsafe.Pointer(&sub.v4l2EventSubscription))); err != nil {
		return fmt.Errorf("subscribe event: type %d: %w", sub.GetType(), err)
	}
	return nil
}

// UnsubscribeEvent unsubscribes from an event type.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-subscribe-event.html
func UnsubscribeEvent(fd uintptr, sub *EventSubscription) error {
	if err := send(fd, C.VIDIOC_UNSUBSCRIBE_EVENT, uintptr(unsafe.Pointer(&sub.v4l2EventSubscription))); err != nil {
		return fmt.Errorf("unsubscribe event: type %d: %w", sub.GetType(), err)
	}
	return nil
}

// DequeueEvent dequeues a pending event.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-dqevent.html
func DequeueEvent(fd uintptr) (*Event, error) {
	event := &Event{}
	if err := send(fd, C.VIDIOC_DQEVENT, uintptr(unsafe.Pointer(&event.v4l2Event))); err != nil {
		return nil, fmt.Errorf("dequeue event: %w", err)
	}
	return event, nil
}
