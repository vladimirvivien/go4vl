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
	EventAll          EventType = cV4L2_EVENT_ALL           // All events
	EventVSync        EventType = cV4L2_EVENT_VSYNC         // Vertical sync
	EventEOS          EventType = cV4L2_EVENT_EOS           // End of stream
	EventCtrl         EventType = cV4L2_EVENT_CTRL          // Control changed
	EventFrameSync    EventType = cV4L2_EVENT_FRAME_SYNC    // Frame sync
	EventSourceChange EventType = cV4L2_EVENT_SOURCE_CHANGE // Source resolution/format changed
	EventMotionDet    EventType = cV4L2_EVENT_MOTION_DET    // Motion detected
	EventPrivateStart EventType = cV4L2_EVENT_PRIVATE_START // Start of driver-specific events
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
	EventSubFlagSendInitial    EventSubscriptionFlags = cV4L2_EVENT_SUB_FL_SEND_INITIAL     // Send initial event on subscribe
	EventSubFlagAllowFeedback  EventSubscriptionFlags = cV4L2_EVENT_SUB_FL_ALLOW_FEEDBACK   // Allow feedback events
)

// EventCtrlChanges represents changes to a control
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type EventCtrlChanges = uint32

const (
	EventCtrlChValue      EventCtrlChanges = cV4L2_EVENT_CTRL_CH_VALUE      // Control value changed
	EventCtrlChFlags      EventCtrlChanges = cV4L2_EVENT_CTRL_CH_FLAGS      // Control flags changed
	EventCtrlChRange      EventCtrlChanges = cV4L2_EVENT_CTRL_CH_RANGE      // Control range changed
	EventCtrlChDimensions EventCtrlChanges = cV4L2_EVENT_CTRL_CH_DIMENSIONS // Control dimensions changed
)

// EventSrcChanges represents source change types
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type EventSrcChanges = uint32

const (
	EventSrcChResolution EventSrcChanges = cV4L2_EVENT_SRC_CH_RESOLUTION // Resolution changed
)

// EventMotionDetFlags represents motion detection event flags
type EventMotionDetFlags = uint32

const (
	EventMotionDetFlagHaveFrameSeq EventMotionDetFlags = cV4L2_EVENT_MD_FL_HAVE_FRAME_SEQ // Have frame sequence number
)

// EventSubscription represents an event subscription (v4l2_event_subscription).
//
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
type EventSubscription struct {
	v4l2EventSubscription cV4L2EventSubscription
}

// NewEventSubscription creates a new event subscription.
func NewEventSubscription(eventType EventType) *EventSubscription {
	es := &EventSubscription{}
	es.v4l2EventSubscription.Type = eventType
	return es
}

// NewControlEventSubscription creates an event subscription for a specific control.
func NewControlEventSubscription(ctrlID CtrlID) *EventSubscription {
	es := NewEventSubscription(EventCtrl)
	es.v4l2EventSubscription.Id = ctrlID
	return es
}

// GetType returns the event type.
func (es *EventSubscription) GetType() EventType {
	return EventType(es.v4l2EventSubscription.Type)
}

// SetType sets the event type.
func (es *EventSubscription) SetType(eventType EventType) {
	es.v4l2EventSubscription.Type = eventType
}

// GetID returns the event ID (e.g., control ID for control events).
func (es *EventSubscription) GetID() uint32 {
	return es.v4l2EventSubscription.Id
}

// SetID sets the event ID.
func (es *EventSubscription) SetID(id uint32) {
	es.v4l2EventSubscription.Id = id
}

// GetFlags returns the subscription flags.
func (es *EventSubscription) GetFlags() EventSubscriptionFlags {
	return EventSubscriptionFlags(es.v4l2EventSubscription.Flags)
}

// SetFlags sets the subscription flags.
func (es *EventSubscription) SetFlags(flags EventSubscriptionFlags) {
	es.v4l2EventSubscription.Flags = flags
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
	v4l2Event cV4L2Event
}

// GetType returns the event type.
func (e *Event) GetType() EventType {
	return EventType(e.v4l2Event.Type)
}

// GetID returns the event ID (e.g., control ID for control events).
func (e *Event) GetID() uint32 {
	return e.v4l2Event.Id
}

// GetPending returns the number of pending events of this type.
func (e *Event) GetPending() uint32 {
	return e.v4l2Event.Pending
}

// GetSequence returns the event sequence number.
func (e *Event) GetSequence() uint32 {
	return e.v4l2Event.Sequence
}

// GetTimestamp returns the event timestamp.
func (e *Event) GetTimestamp() time.Time {
	ts := e.v4l2Event.Timestamp
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

// GetCtrlData returns the control event data (valid if Type is EventCtrl).
func (e *Event) GetCtrlData() EventCtrlData {
	ctrlPtr := (*cV4L2EventCtrl)(unsafe.Pointer(&e.v4l2Event.U[0]))

	// The value field is a union of value/value64
	value := ctrlPtr.Value
	value64 := *(*int64)(unsafe.Pointer(&ctrlPtr.Value))

	return EventCtrlData{
		Changes: EventCtrlChanges(ctrlPtr.Changes),
		Type:    CtrlType(ctrlPtr.Type),
		Value:   value,
		Value64: value64,
		Flags:   ctrlPtr.Flags,
		Minimum: ctrlPtr.Minimum,
		Maximum: ctrlPtr.Maximum,
		Step:    ctrlPtr.Step,
		Default: ctrlPtr.Default_value,
	}
}

// GetSrcChangeData returns the source change event data (valid if Type is EventSourceChange).
func (e *Event) GetSrcChangeData() EventSrcChangeData {
	srcPtr := (*cV4L2EventSrcChange)(unsafe.Pointer(&e.v4l2Event.U[0]))
	return EventSrcChangeData{
		Changes: EventSrcChanges(srcPtr.Changes),
	}
}

// GetFrameSyncData returns the frame sync event data (valid if Type is EventFrameSync).
func (e *Event) GetFrameSyncData() EventFrameSyncData {
	fsPtr := (*cV4L2EventFrameSync)(unsafe.Pointer(&e.v4l2Event.U[0]))
	return EventFrameSyncData{
		FrameSequence: fsPtr.Sequence,
	}
}

// GetMotionDetData returns the motion detection event data (valid if Type is EventMotionDet).
func (e *Event) GetMotionDetData() EventMotionDetData {
	mdPtr := (*cV4L2EventMotionDet)(unsafe.Pointer(&e.v4l2Event.U[0]))
	return EventMotionDetData{
		Flags:         EventMotionDetFlags(mdPtr.Flags),
		FrameSequence: mdPtr.Frame_sequence,
		RegionMask:    mdPtr.Region_mask,
	}
}

// GetRawData returns the raw event data as a byte slice (for unknown event types).
func (e *Event) GetRawData() []byte {
	data := make([]byte, 64)
	copy(data, e.v4l2Event.U[:])
	return data
}

// SubscribeEvent subscribes to an event type.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-subscribe-event.html
func SubscribeEvent(fd uintptr, sub *EventSubscription) error {
	if err := send(fd, cVIDIOC_SUBSCRIBE_EVENT, unsafe.Pointer(&sub.v4l2EventSubscription)); err != nil {
		return fmt.Errorf("subscribe event: type %d: %w", sub.GetType(), err)
	}
	return nil
}

// UnsubscribeEvent unsubscribes from an event type.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-subscribe-event.html
func UnsubscribeEvent(fd uintptr, sub *EventSubscription) error {
	if err := send(fd, cVIDIOC_UNSUBSCRIBE_EVENT, unsafe.Pointer(&sub.v4l2EventSubscription)); err != nil {
		return fmt.Errorf("unsubscribe event: type %d: %w", sub.GetType(), err)
	}
	return nil
}

// DequeueEvent dequeues a pending event.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-dqevent.html
func DequeueEvent(fd uintptr) (*Event, error) {
	event := &Event{}
	if err := send(fd, cVIDIOC_DQEVENT, unsafe.Pointer(&event.v4l2Event)); err != nil {
		return nil, fmt.Errorf("dequeue event: %w", err)
	}
	return event, nil
}
