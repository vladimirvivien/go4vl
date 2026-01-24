package v4l2

// codec_events.go provides codec-specific event helpers for V4L2 stateful codecs.
//
// Codecs use events to signal important state changes:
// - EventSourceChange: Resolution or format change detected (decoders)
// - EventEOS: End of stream reached (both encoders and decoders)
//
// These events are essential for proper codec operation, especially for decoders
// that need to handle dynamic resolution changes in the stream.
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-subscribe-event.html

import "fmt"

// SubscribeSourceChangeEvent subscribes to source change events.
// Source change events are critical for decoders to detect resolution changes.
// When the decoder encounters a new resolution in the stream, it signals this event.
//
// Parameters:
//   - fd: File descriptor of an opened codec device
//
// Returns:
//   - error: An error if subscription fails
//
// Usage: After receiving this event, the application should:
//  1. Stop streaming on the capture queue
//  2. Query the new format with VIDIOC_G_FMT
//  3. Reallocate capture buffers if needed
//  4. Resume streaming
func SubscribeSourceChangeEvent(fd uintptr) error {
	sub := NewEventSubscription(EventSourceChange)
	if err := SubscribeEvent(fd, sub); err != nil {
		return fmt.Errorf("subscribe source change: %w", err)
	}
	return nil
}

// SubscribeSourceChangeEventWithInitial subscribes to source change events
// and requests an initial event to be sent immediately.
// This is useful to get the current source state when subscribing.
func SubscribeSourceChangeEventWithInitial(fd uintptr) error {
	sub := NewEventSubscription(EventSourceChange)
	sub.SetFlags(EventSubFlagSendInitial)
	if err := SubscribeEvent(fd, sub); err != nil {
		return fmt.Errorf("subscribe source change: %w", err)
	}
	return nil
}

// UnsubscribeSourceChangeEvent unsubscribes from source change events.
func UnsubscribeSourceChangeEvent(fd uintptr) error {
	sub := NewEventSubscription(EventSourceChange)
	if err := UnsubscribeEvent(fd, sub); err != nil {
		return fmt.Errorf("unsubscribe source change: %w", err)
	}
	return nil
}

// SubscribeEOSEvent subscribes to end-of-stream events.
// EOS events indicate that the codec has finished processing all data
// and no more output will be produced.
//
// Parameters:
//   - fd: File descriptor of an opened codec device
//
// Returns:
//   - error: An error if subscription fails
//
// Usage: After receiving this event, the application should:
//  1. Dequeue any remaining buffers
//  2. Stop streaming
//  3. Clean up resources
func SubscribeEOSEvent(fd uintptr) error {
	sub := NewEventSubscription(EventEOS)
	if err := SubscribeEvent(fd, sub); err != nil {
		return fmt.Errorf("subscribe eos: %w", err)
	}
	return nil
}

// UnsubscribeEOSEvent unsubscribes from end-of-stream events.
func UnsubscribeEOSEvent(fd uintptr) error {
	sub := NewEventSubscription(EventEOS)
	if err := UnsubscribeEvent(fd, sub); err != nil {
		return fmt.Errorf("unsubscribe eos: %w", err)
	}
	return nil
}

// SubscribeCodecEvents subscribes to all codec-relevant events.
// This is a convenience function that subscribes to:
// - Source change events (for resolution changes)
// - EOS events (for end of stream)
//
// Parameters:
//   - fd: File descriptor of an opened codec device
//
// Returns:
//   - error: An error if any subscription fails
func SubscribeCodecEvents(fd uintptr) error {
	if err := SubscribeSourceChangeEvent(fd); err != nil {
		return err
	}
	if err := SubscribeEOSEvent(fd); err != nil {
		// Try to clean up the source change subscription
		_ = UnsubscribeSourceChangeEvent(fd)
		return err
	}
	return nil
}

// UnsubscribeCodecEvents unsubscribes from all codec-relevant events.
func UnsubscribeCodecEvents(fd uintptr) error {
	var firstErr error

	if err := UnsubscribeSourceChangeEvent(fd); err != nil && firstErr == nil {
		firstErr = err
	}
	if err := UnsubscribeEOSEvent(fd); err != nil && firstErr == nil {
		firstErr = err
	}

	return firstErr
}

// IsSourceChangeEvent checks if an event is a source change event.
func IsSourceChangeEvent(event *Event) bool {
	return event != nil && event.GetType() == EventSourceChange
}

// IsEOSEvent checks if an event is an end-of-stream event.
func IsEOSEvent(event *Event) bool {
	return event != nil && event.GetType() == EventEOS
}

// IsResolutionChangeEvent checks if a source change event indicates a resolution change.
func IsResolutionChangeEvent(event *Event) bool {
	if !IsSourceChangeEvent(event) {
		return false
	}
	srcData := event.GetSrcChangeData()
	return srcData.Changes&EventSrcChResolution != 0
}

// WaitForEvent waits for and returns the next event.
// This is a convenience wrapper around DequeueEvent.
// For production use, consider using select/poll on the file descriptor.
func WaitForEvent(fd uintptr) (*Event, error) {
	return DequeueEvent(fd)
}

// PollForEvent attempts to dequeue an event without blocking.
// Returns nil if no event is available (not an error condition).
// This function should be called after poll/select indicates an event is ready.
func PollForEvent(fd uintptr) (*Event, error) {
	event, err := DequeueEvent(fd)
	if err != nil {
		// EAGAIN means no event available - not an error
		if err == ErrorTemporary || err == ErrorInterrupted {
			return nil, nil
		}
		return nil, err
	}
	return event, nil
}
