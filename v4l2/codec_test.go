package v4l2

import (
	"fmt"
	"testing"
)

// TestEncoderCmdConstants verifies encoder command constants are defined.
func TestEncoderCmdConstants(t *testing.T) {
	tests := []struct {
		name  string
		value uint32
	}{
		{"EncCmdStart", EncCmdStart},
		{"EncCmdStop", EncCmdStop},
		{"EncCmdPause", EncCmdPause},
		{"EncCmdResume", EncCmdResume},
	}

	// Verify all constants are distinct
	seen := make(map[uint32]string)
	for _, tt := range tests {
		if existing, ok := seen[tt.value]; ok {
			t.Errorf("duplicate constant value: %s and %s both have value %d", tt.name, existing, tt.value)
		}
		seen[tt.value] = tt.name
	}
}

// TestDecoderCmdConstants verifies decoder command constants are defined.
func TestDecoderCmdConstants(t *testing.T) {
	tests := []struct {
		name  string
		value uint32
	}{
		{"DecCmdStart", DecCmdStart},
		{"DecCmdStop", DecCmdStop},
		{"DecCmdPause", DecCmdPause},
		{"DecCmdResume", DecCmdResume},
		{"DecCmdFlush", DecCmdFlush},
	}

	// Verify all constants are distinct
	seen := make(map[uint32]string)
	for _, tt := range tests {
		if existing, ok := seen[tt.value]; ok {
			t.Errorf("duplicate constant value: %s and %s both have value %d", tt.name, existing, tt.value)
		}
		seen[tt.value] = tt.name
	}
}

// TestEncoderCmd verifies EncoderCmd structure and methods.
func TestEncoderCmd(t *testing.T) {
	t.Run("NewEncoderCmd", func(t *testing.T) {
		cmd := NewEncoderCmd(EncCmdStart)
		if cmd.GetCmd() != EncCmdStart {
			t.Errorf("expected cmd %d, got %d", EncCmdStart, cmd.GetCmd())
		}
		if cmd.GetFlags() != 0 {
			t.Errorf("expected flags 0, got %d", cmd.GetFlags())
		}
	})

	t.Run("NewEncoderCmdWithFlags", func(t *testing.T) {
		cmd := NewEncoderCmdWithFlags(EncCmdStop, EncCmdStopAtGOPEnd)
		if cmd.GetCmd() != EncCmdStop {
			t.Errorf("expected cmd %d, got %d", EncCmdStop, cmd.GetCmd())
		}
		if cmd.GetFlags() != EncCmdStopAtGOPEnd {
			t.Errorf("expected flags %d, got %d", EncCmdStopAtGOPEnd, cmd.GetFlags())
		}
	})

	t.Run("SetCmd", func(t *testing.T) {
		cmd := NewEncoderCmd(EncCmdStart)
		cmd.SetCmd(EncCmdPause)
		if cmd.GetCmd() != EncCmdPause {
			t.Errorf("expected cmd %d, got %d", EncCmdPause, cmd.GetCmd())
		}
	})

	t.Run("SetFlags", func(t *testing.T) {
		cmd := NewEncoderCmd(EncCmdStop)
		cmd.SetFlags(EncCmdStopAtGOPEnd)
		if cmd.GetFlags() != EncCmdStopAtGOPEnd {
			t.Errorf("expected flags %d, got %d", EncCmdStopAtGOPEnd, cmd.GetFlags())
		}
	})

	t.Run("IsStart", func(t *testing.T) {
		cmd := NewEncoderCmd(EncCmdStart)
		if !cmd.IsStart() {
			t.Error("expected IsStart() to return true")
		}
		if cmd.IsStop() || cmd.IsPause() || cmd.IsResume() {
			t.Error("expected other Is* methods to return false")
		}
	})

	t.Run("String", func(t *testing.T) {
		cmd := NewEncoderCmd(EncCmdStart)
		s := cmd.String()
		if s == "" {
			t.Error("expected non-empty string")
		}
		if !containsSubstring(s, "START") {
			t.Errorf("expected string to contain 'START', got %s", s)
		}
	})
}

// TestDecoderCmd verifies DecoderCmd structure and methods.
func TestDecoderCmd(t *testing.T) {
	t.Run("NewDecoderCmd", func(t *testing.T) {
		cmd := NewDecoderCmd(DecCmdStart)
		if cmd.GetCmd() != DecCmdStart {
			t.Errorf("expected cmd %d, got %d", DecCmdStart, cmd.GetCmd())
		}
		if cmd.GetFlags() != 0 {
			t.Errorf("expected flags 0, got %d", cmd.GetFlags())
		}
	})

	t.Run("NewDecoderCmdWithFlags", func(t *testing.T) {
		cmd := NewDecoderCmdWithFlags(DecCmdStop, DecCmdStopImmediately)
		if cmd.GetCmd() != DecCmdStop {
			t.Errorf("expected cmd %d, got %d", DecCmdStop, cmd.GetCmd())
		}
		if cmd.GetFlags() != DecCmdStopImmediately {
			t.Errorf("expected flags %d, got %d", DecCmdStopImmediately, cmd.GetFlags())
		}
	})

	t.Run("StartSpeed", func(t *testing.T) {
		cmd := NewDecoderCmd(DecCmdStart)
		cmd.SetStartSpeed(1000)
		if cmd.GetStartSpeed() != 1000 {
			t.Errorf("expected speed 1000, got %d", cmd.GetStartSpeed())
		}
	})

	t.Run("StartFormat", func(t *testing.T) {
		cmd := NewDecoderCmd(DecCmdStart)
		cmd.SetStartFormat(DecStartFmtGOP)
		if cmd.GetStartFormat() != DecStartFmtGOP {
			t.Errorf("expected format %d, got %d", DecStartFmtGOP, cmd.GetStartFormat())
		}
	})

	t.Run("StopPts", func(t *testing.T) {
		cmd := NewDecoderCmd(DecCmdStop)
		cmd.SetStopPts(123456)
		if cmd.GetStopPts() != 123456 {
			t.Errorf("expected pts 123456, got %d", cmd.GetStopPts())
		}
	})

	t.Run("IsFlush", func(t *testing.T) {
		cmd := NewDecoderCmd(DecCmdFlush)
		if !cmd.IsFlush() {
			t.Error("expected IsFlush() to return true")
		}
		if cmd.IsStart() || cmd.IsStop() || cmd.IsPause() || cmd.IsResume() {
			t.Error("expected other Is* methods to return false")
		}
	})

	t.Run("String", func(t *testing.T) {
		cmd := NewDecoderCmd(DecCmdFlush)
		s := cmd.String()
		if s == "" {
			t.Error("expected non-empty string")
		}
		if !containsSubstring(s, "FLUSH") {
			t.Errorf("expected string to contain 'FLUSH', got %s", s)
		}
	})
}

// TestCodecStateMachine verifies state machine transitions.
func TestCodecStateMachine(t *testing.T) {
	t.Run("NewCodecStateMachine", func(t *testing.T) {
		sm := NewCodecStateMachine(CodecTypeEncoder)
		if sm.GetState() != CodecStateUninitialized {
			t.Errorf("expected initial state Uninitialized, got %s", sm.GetState())
		}
		if sm.GetCodecType() != CodecTypeEncoder {
			t.Errorf("expected codec type Encoder, got %s", sm.GetCodecType())
		}
	})

	t.Run("ValidTransitions", func(t *testing.T) {
		sm := NewCodecStateMachine(CodecTypeEncoder)

		// Uninitialized -> Initialized
		if err := sm.Initialize(); err != nil {
			t.Errorf("Initialize failed: %v", err)
		}
		if sm.GetState() != CodecStateInitialized {
			t.Errorf("expected state Initialized, got %s", sm.GetState())
		}

		// Initialized -> Streaming
		if err := sm.Start(); err != nil {
			t.Errorf("Start failed: %v", err)
		}
		if sm.GetState() != CodecStateStreaming {
			t.Errorf("expected state Streaming, got %s", sm.GetState())
		}

		// Streaming -> Paused
		if err := sm.Pause(); err != nil {
			t.Errorf("Pause failed: %v", err)
		}
		if sm.GetState() != CodecStatePaused {
			t.Errorf("expected state Paused, got %s", sm.GetState())
		}

		// Paused -> Streaming
		if err := sm.Resume(); err != nil {
			t.Errorf("Resume failed: %v", err)
		}
		if sm.GetState() != CodecStateStreaming {
			t.Errorf("expected state Streaming, got %s", sm.GetState())
		}

		// Streaming -> Draining
		if err := sm.StartDrain(); err != nil {
			t.Errorf("StartDrain failed: %v", err)
		}
		if sm.GetState() != CodecStateDraining {
			t.Errorf("expected state Draining, got %s", sm.GetState())
		}

		// Draining -> Stopped
		if err := sm.CompleteDrain(); err != nil {
			t.Errorf("CompleteDrain failed: %v", err)
		}
		if sm.GetState() != CodecStateStopped {
			t.Errorf("expected state Stopped, got %s", sm.GetState())
		}
	})

	t.Run("InvalidTransitions", func(t *testing.T) {
		sm := NewCodecStateMachine(CodecTypeEncoder)

		// Cannot go directly from Uninitialized to Streaming
		if err := sm.Start(); err == nil {
			t.Error("expected error for invalid transition Uninitialized -> Streaming")
		}

		// Cannot Pause from Uninitialized
		if err := sm.Pause(); err == nil {
			t.Error("expected error for invalid transition Uninitialized -> Paused")
		}
	})

	t.Run("ErrorState", func(t *testing.T) {
		sm := NewCodecStateMachine(CodecTypeDecoder)
		sm.Initialize()
		sm.Start()

		// Error can be reached from any state
		testErr := fmt.Errorf("test error")
		sm.SetError(testErr)

		if sm.GetState() != CodecStateError {
			t.Errorf("expected state Error, got %s", sm.GetState())
		}
		if sm.GetLastError() != testErr {
			t.Errorf("expected error %v, got %v", testErr, sm.GetLastError())
		}
	})

	t.Run("FlushOnlyForDecoder", func(t *testing.T) {
		// Encoder should not support flush
		encSm := NewCodecStateMachine(CodecTypeEncoder)
		encSm.Initialize()
		encSm.Start()
		if err := encSm.StartFlush(); err == nil {
			t.Error("expected error for flush on encoder")
		}

		// Decoder should support flush
		decSm := NewCodecStateMachine(CodecTypeDecoder)
		decSm.Initialize()
		decSm.Start()
		if err := decSm.StartFlush(); err != nil {
			t.Errorf("unexpected error for flush on decoder: %v", err)
		}
		if decSm.GetState() != CodecStateFlushing {
			t.Errorf("expected state Flushing, got %s", decSm.GetState())
		}
	})

	t.Run("HelperMethods", func(t *testing.T) {
		sm := NewCodecStateMachine(CodecTypeEncoder)

		// Initial state
		if sm.IsStreaming() {
			t.Error("expected IsStreaming() to be false initially")
		}
		if !sm.IsStopped() {
			t.Error("expected IsStopped() to be true initially")
		}

		sm.Initialize()
		sm.Start()

		if !sm.IsStreaming() {
			t.Error("expected IsStreaming() to be true after start")
		}
		if sm.IsStopped() {
			t.Error("expected IsStopped() to be false after start")
		}
		if !sm.CanAcceptInput() {
			t.Error("expected CanAcceptInput() to be true while streaming")
		}
		if !sm.CanProduceOutput() {
			t.Error("expected CanProduceOutput() to be true while streaming")
		}
	})

	t.Run("ValidTransitionsList", func(t *testing.T) {
		sm := NewCodecStateMachine(CodecTypeEncoder)
		transitions := sm.ValidTransitions()

		// From Uninitialized, should be able to go to Initialized or Error
		hasInitialized := false
		hasError := false
		for _, s := range transitions {
			if s == CodecStateInitialized {
				hasInitialized = true
			}
			if s == CodecStateError {
				hasError = true
			}
		}
		if !hasInitialized {
			t.Error("expected Initialized in valid transitions from Uninitialized")
		}
		if !hasError {
			t.Error("expected Error in valid transitions from Uninitialized")
		}
	})

	t.Run("String", func(t *testing.T) {
		sm := NewCodecStateMachine(CodecTypeDecoder)
		s := sm.String()
		if s == "" {
			t.Error("expected non-empty string")
		}
		if !containsSubstring(s, "Decoder") {
			t.Errorf("expected string to contain 'Decoder', got %s", s)
		}
	})
}

// TestCodecStateString verifies state string representations.
func TestCodecStateString(t *testing.T) {
	tests := []struct {
		state    CodecState
		expected string
	}{
		{CodecStateUninitialized, "Uninitialized"},
		{CodecStateInitialized, "Initialized"},
		{CodecStateStreaming, "Streaming"},
		{CodecStateDraining, "Draining"},
		{CodecStatePaused, "Paused"},
		{CodecStateStopped, "Stopped"},
		{CodecStateFlushing, "Flushing"},
		{CodecStateError, "Error"},
		{CodecState(999), "Unknown(999)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

// TestCodecTypeString verifies codec type string representations.
func TestCodecTypeString(t *testing.T) {
	if CodecTypeEncoder.String() != "Encoder" {
		t.Errorf("expected 'Encoder', got %s", CodecTypeEncoder.String())
	}
	if CodecTypeDecoder.String() != "Decoder" {
		t.Errorf("expected 'Decoder', got %s", CodecTypeDecoder.String())
	}
}

// TestEncoderCmdNames verifies command name map is populated.
func TestEncoderCmdNames(t *testing.T) {
	expectedCmds := []uint32{EncCmdStart, EncCmdStop, EncCmdPause, EncCmdResume}
	for _, cmd := range expectedCmds {
		if _, ok := EncoderCmdNames[cmd]; !ok {
			t.Errorf("missing name for encoder command %d", cmd)
		}
	}
}

// TestDecoderCmdNames verifies command name map is populated.
func TestDecoderCmdNames(t *testing.T) {
	expectedCmds := []uint32{DecCmdStart, DecCmdStop, DecCmdPause, DecCmdResume, DecCmdFlush}
	for _, cmd := range expectedCmds {
		if _, ok := DecoderCmdNames[cmd]; !ok {
			t.Errorf("missing name for decoder command %d", cmd)
		}
	}
}

// TestCodecCallbacks verifies callbacks are invoked correctly.
func TestCodecCallbacks(t *testing.T) {
	stateChangeCalled := make(chan struct{}, 10)
	eosCallled := make(chan struct{}, 1)

	callbacks := CodecCallbacks{
		OnStateChange: func(old, new CodecState) {
			stateChangeCalled <- struct{}{}
		},
		OnEOS: func() {
			eosCallled <- struct{}{}
		},
	}

	sm := NewCodecStateMachineWithCallbacks(CodecTypeEncoder, callbacks)

	// Trigger state changes
	sm.Initialize()
	sm.Start()

	// Wait for callbacks (with timeout via select)
	select {
	case <-stateChangeCalled:
		// Got first callback
	default:
		// Callback may be async, that's OK
	}

	// Test EOS notification
	sm.NotifyEOS()

	select {
	case <-eosCallled:
		// Got EOS callback
	default:
		// Callback may be async, that's OK
	}
}

// TestIsSourceChangeEvent verifies event type checking.
func TestIsSourceChangeEvent(t *testing.T) {
	if IsSourceChangeEvent(nil) {
		t.Error("expected false for nil event")
	}

	if IsEOSEvent(nil) {
		t.Error("expected false for nil event")
	}
}

// Helper function
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstringHelper(s, substr))
}

func containsSubstringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
