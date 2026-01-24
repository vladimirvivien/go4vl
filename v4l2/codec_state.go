package v4l2

// codec_state.go provides a state machine for V4L2 stateful codec lifecycle management.
//
// V4L2 stateful codecs (encoders and decoders) follow a specific state machine:
//
//	Uninitialized → Initialized → Streaming → Draining → Stopped
//	                     ↓            ↓
//	                  Paused ←────────┘
//	                     ↓
//	                  Error (from any state)
//
// This state machine ensures correct operation sequence and provides callbacks
// for state change notifications, drain completion, resolution changes, and errors.
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/dev-encoder.html
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/dev-decoder.html

import (
	"fmt"
	"sync"
)

// CodecState represents the current state of a V4L2 stateful codec.
type CodecState int

// Codec states
const (
	// CodecStateUninitialized is the initial state before any setup.
	// The codec device is opened but not configured.
	CodecStateUninitialized CodecState = iota

	// CodecStateInitialized means the codec is configured but not streaming.
	// Format has been set, buffers may be allocated.
	CodecStateInitialized

	// CodecStateStreaming means the codec is actively processing data.
	// Buffers are being queued and dequeued.
	CodecStateStreaming

	// CodecStateDraining means the codec is draining remaining data.
	// No new input is accepted, waiting for output to complete.
	CodecStateDraining

	// CodecStatePaused means the codec is temporarily paused.
	// Can be resumed to continue processing.
	CodecStatePaused

	// CodecStateStopped means the codec has stopped.
	// Drain is complete or stop was immediate.
	CodecStateStopped

	// CodecStateFlushing means the codec is flushing (decoder only).
	// Used for seeking operations.
	CodecStateFlushing

	// CodecStateError means an unrecoverable error occurred.
	// The codec must be reset or closed.
	CodecStateError
)

// String returns a human-readable name for the codec state.
func (s CodecState) String() string {
	switch s {
	case CodecStateUninitialized:
		return "Uninitialized"
	case CodecStateInitialized:
		return "Initialized"
	case CodecStateStreaming:
		return "Streaming"
	case CodecStateDraining:
		return "Draining"
	case CodecStatePaused:
		return "Paused"
	case CodecStateStopped:
		return "Stopped"
	case CodecStateFlushing:
		return "Flushing"
	case CodecStateError:
		return "Error"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

// CodecCallbacks provides callback functions for codec state machine events.
// All callbacks are optional (can be nil).
type CodecCallbacks struct {
	// OnStateChange is called when the codec state changes.
	// Parameters are (oldState, newState).
	OnStateChange func(old, new CodecState)

	// OnDrainComplete is called when the drain operation completes.
	// The codec transitions to Stopped state after this callback.
	OnDrainComplete func()

	// OnFlushComplete is called when the flush operation completes (decoder only).
	// The codec transitions back to Streaming state after this callback.
	OnFlushComplete func()

	// OnResolutionChange is called when a resolution change is detected (decoder only).
	// This occurs when the decoder encounters a new resolution in the stream.
	// Parameters are (width, height, pixelFormat).
	OnResolutionChange func(width, height uint32, pixelFormat FourCCType)

	// OnError is called when an error occurs.
	// The codec transitions to Error state after this callback.
	OnError func(err error)

	// OnEOS is called when end-of-stream is detected.
	OnEOS func()
}

// CodecStateMachine manages the state transitions for a V4L2 codec.
// It enforces valid state transitions and provides thread-safe state management.
type CodecStateMachine struct {
	mu           sync.RWMutex
	currentState CodecState
	callbacks    CodecCallbacks
	lastError    error
	codecType    CodecType
}

// CodecType indicates whether the codec is an encoder or decoder.
type CodecType int

const (
	// CodecTypeEncoder for video encoders
	CodecTypeEncoder CodecType = iota
	// CodecTypeDecoder for video decoders
	CodecTypeDecoder
)

// String returns the codec type name.
func (t CodecType) String() string {
	switch t {
	case CodecTypeEncoder:
		return "Encoder"
	case CodecTypeDecoder:
		return "Decoder"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}

// NewCodecStateMachine creates a new codec state machine.
func NewCodecStateMachine(codecType CodecType) *CodecStateMachine {
	return &CodecStateMachine{
		currentState: CodecStateUninitialized,
		codecType:    codecType,
	}
}

// NewCodecStateMachineWithCallbacks creates a new codec state machine with callbacks.
func NewCodecStateMachineWithCallbacks(codecType CodecType, callbacks CodecCallbacks) *CodecStateMachine {
	return &CodecStateMachine{
		currentState: CodecStateUninitialized,
		callbacks:    callbacks,
		codecType:    codecType,
	}
}

// GetState returns the current state (thread-safe).
func (sm *CodecStateMachine) GetState() CodecState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.currentState
}

// GetCodecType returns the codec type (encoder or decoder).
func (sm *CodecStateMachine) GetCodecType() CodecType {
	return sm.codecType
}

// GetLastError returns the last error that occurred.
func (sm *CodecStateMachine) GetLastError() error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.lastError
}

// SetCallbacks sets the callback functions.
func (sm *CodecStateMachine) SetCallbacks(callbacks CodecCallbacks) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.callbacks = callbacks
}

// IsStreaming returns true if the codec is in a streaming state.
func (sm *CodecStateMachine) IsStreaming() bool {
	state := sm.GetState()
	return state == CodecStateStreaming || state == CodecStateDraining || state == CodecStateFlushing
}

// IsStopped returns true if the codec is stopped or in error state.
func (sm *CodecStateMachine) IsStopped() bool {
	state := sm.GetState()
	return state == CodecStateStopped || state == CodecStateError || state == CodecStateUninitialized
}

// CanAcceptInput returns true if the codec can accept new input data.
func (sm *CodecStateMachine) CanAcceptInput() bool {
	state := sm.GetState()
	return state == CodecStateStreaming || state == CodecStatePaused
}

// CanProduceOutput returns true if the codec may produce output data.
func (sm *CodecStateMachine) CanProduceOutput() bool {
	state := sm.GetState()
	return state == CodecStateStreaming || state == CodecStateDraining || state == CodecStateFlushing
}

// transitionTo attempts to transition to a new state.
// Returns an error if the transition is not valid.
func (sm *CodecStateMachine) transitionTo(newState CodecState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	oldState := sm.currentState

	if !sm.isValidTransition(oldState, newState) {
		return fmt.Errorf("invalid state transition: %s -> %s", oldState, newState)
	}

	sm.currentState = newState

	// Call state change callback (outside lock to avoid deadlock)
	if sm.callbacks.OnStateChange != nil {
		go sm.callbacks.OnStateChange(oldState, newState)
	}

	return nil
}

// isValidTransition checks if a state transition is valid.
func (sm *CodecStateMachine) isValidTransition(from, to CodecState) bool {
	// Error state can be reached from any state
	if to == CodecStateError {
		return true
	}

	// Define valid transitions
	switch from {
	case CodecStateUninitialized:
		return to == CodecStateInitialized

	case CodecStateInitialized:
		return to == CodecStateStreaming || to == CodecStateUninitialized

	case CodecStateStreaming:
		return to == CodecStateDraining || to == CodecStatePaused ||
			to == CodecStateStopped || to == CodecStateFlushing

	case CodecStateDraining:
		return to == CodecStateStopped

	case CodecStatePaused:
		return to == CodecStateStreaming || to == CodecStateStopped ||
			to == CodecStateDraining

	case CodecStateStopped:
		return to == CodecStateInitialized || to == CodecStateUninitialized

	case CodecStateFlushing:
		return to == CodecStateStreaming || to == CodecStateStopped

	case CodecStateError:
		return to == CodecStateUninitialized
	}

	return false
}

// Initialize transitions from Uninitialized to Initialized state.
// Call this after setting up the codec format and allocating buffers.
func (sm *CodecStateMachine) Initialize() error {
	return sm.transitionTo(CodecStateInitialized)
}

// Start transitions from Initialized to Streaming state.
// Call this when starting the codec with VIDIOC_STREAMON.
func (sm *CodecStateMachine) Start() error {
	return sm.transitionTo(CodecStateStreaming)
}

// Pause transitions from Streaming to Paused state.
// Call this when pausing the codec.
func (sm *CodecStateMachine) Pause() error {
	return sm.transitionTo(CodecStatePaused)
}

// Resume transitions from Paused to Streaming state.
// Call this when resuming the codec.
func (sm *CodecStateMachine) Resume() error {
	return sm.transitionTo(CodecStateStreaming)
}

// StartDrain transitions to Draining state.
// Call this when initiating a drain operation (stop command).
func (sm *CodecStateMachine) StartDrain() error {
	return sm.transitionTo(CodecStateDraining)
}

// StartFlush transitions to Flushing state (decoder only).
// Call this when initiating a flush operation for seeking.
func (sm *CodecStateMachine) StartFlush() error {
	if sm.codecType != CodecTypeDecoder {
		return fmt.Errorf("flush is only supported for decoders")
	}
	return sm.transitionTo(CodecStateFlushing)
}

// CompleteDrain transitions from Draining to Stopped state.
// Call this when the drain operation completes (EOS received).
func (sm *CodecStateMachine) CompleteDrain() error {
	err := sm.transitionTo(CodecStateStopped)
	if err != nil {
		return err
	}

	if sm.callbacks.OnDrainComplete != nil {
		go sm.callbacks.OnDrainComplete()
	}

	return nil
}

// CompleteFlush transitions from Flushing to Streaming state.
// Call this when the flush operation completes.
func (sm *CodecStateMachine) CompleteFlush() error {
	err := sm.transitionTo(CodecStateStreaming)
	if err != nil {
		return err
	}

	if sm.callbacks.OnFlushComplete != nil {
		go sm.callbacks.OnFlushComplete()
	}

	return nil
}

// Stop transitions to Stopped state.
// Can be called from Streaming, Paused, or Flushing states.
func (sm *CodecStateMachine) Stop() error {
	return sm.transitionTo(CodecStateStopped)
}

// Reset transitions from Stopped or Error to Initialized state.
// Call this to prepare the codec for reuse.
func (sm *CodecStateMachine) Reset() error {
	return sm.transitionTo(CodecStateInitialized)
}

// Uninitialize transitions to Uninitialized state.
// Call this when closing the codec.
func (sm *CodecStateMachine) Uninitialize() error {
	return sm.transitionTo(CodecStateUninitialized)
}

// SetError transitions to Error state with the given error.
func (sm *CodecStateMachine) SetError(err error) {
	sm.mu.Lock()
	sm.lastError = err
	oldState := sm.currentState
	sm.currentState = CodecStateError
	sm.mu.Unlock()

	if sm.callbacks.OnStateChange != nil {
		go sm.callbacks.OnStateChange(oldState, CodecStateError)
	}

	if sm.callbacks.OnError != nil {
		go sm.callbacks.OnError(err)
	}
}

// NotifyResolutionChange notifies of a resolution change (decoder only).
// This does not change the state, but calls the callback.
func (sm *CodecStateMachine) NotifyResolutionChange(width, height uint32, pixelFormat FourCCType) {
	if sm.callbacks.OnResolutionChange != nil {
		go sm.callbacks.OnResolutionChange(width, height, pixelFormat)
	}
}

// NotifyEOS notifies of end-of-stream.
// This does not change the state, but calls the callback.
func (sm *CodecStateMachine) NotifyEOS() {
	if sm.callbacks.OnEOS != nil {
		go sm.callbacks.OnEOS()
	}
}

// ValidTransitions returns a list of valid next states from the current state.
func (sm *CodecStateMachine) ValidTransitions() []CodecState {
	sm.mu.RLock()
	current := sm.currentState
	sm.mu.RUnlock()

	var valid []CodecState
	allStates := []CodecState{
		CodecStateUninitialized,
		CodecStateInitialized,
		CodecStateStreaming,
		CodecStateDraining,
		CodecStatePaused,
		CodecStateStopped,
		CodecStateFlushing,
		CodecStateError,
	}

	for _, state := range allStates {
		if sm.isValidTransition(current, state) {
			valid = append(valid, state)
		}
	}

	return valid
}

// String returns a string representation of the state machine.
func (sm *CodecStateMachine) String() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return fmt.Sprintf("CodecStateMachine{type: %s, state: %s}", sm.codecType, sm.currentState)
}
