// +build integration

package test

import (
	"os"
	"testing"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// Common codec device paths to check
var codecDevicePaths = []string{
	"/dev/video10", // Raspberry Pi decoder
	"/dev/video11", // Raspberry Pi encoder
	"/dev/video0",  // Rockchip decoder
	"/dev/video1",  // Rockchip encoder
}

// findCodecDevice searches for a codec device
func findCodecDevice(t *testing.T, wantEncoder bool) string {
	t.Helper()

	for _, path := range codecDevicePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		fd, err := v4l2.OpenDevice(path, 0, 0)
		if err != nil {
			continue
		}

		cap, err := v4l2.GetCapability(fd)
		v4l2.CloseDevice(fd)
		if err != nil {
			continue
		}

		if wantEncoder && cap.IsEncoderSupported() {
			return path
		}
		if !wantEncoder && cap.IsDecoderSupported() {
			return path
		}
	}

	return ""
}

// TestCodec_EncoderCapabilityCheck tests encoder capability detection
func TestCodec_EncoderCapabilityCheck(t *testing.T) {
	encoderPath := findCodecDevice(t, true)
	if encoderPath == "" {
		t.Skip("No encoder device found - skipping test")
	}

	fd, err := v4l2.OpenDevice(encoderPath, 0, 0)
	if err != nil {
		t.Fatalf("Failed to open encoder device: %v", err)
	}
	defer v4l2.CloseDevice(fd)

	cap, err := v4l2.GetCapability(fd)
	if err != nil {
		t.Fatalf("Failed to get capabilities: %v", err)
	}

	t.Logf("Encoder device: %s", encoderPath)
	t.Logf("  Driver: %s", cap.Driver)
	t.Logf("  Card: %s", cap.Card)

	if !cap.IsM2MSupported() {
		t.Error("Encoder should support M2M")
	}
	if !cap.IsEncoderSupported() {
		t.Error("Device should be detected as encoder")
	}

	t.Log("Encoder capabilities verified")
}

// TestCodec_DecoderCapabilityCheck tests decoder capability detection
func TestCodec_DecoderCapabilityCheck(t *testing.T) {
	decoderPath := findCodecDevice(t, false)
	if decoderPath == "" {
		t.Skip("No decoder device found - skipping test")
	}

	fd, err := v4l2.OpenDevice(decoderPath, 0, 0)
	if err != nil {
		t.Fatalf("Failed to open decoder device: %v", err)
	}
	defer v4l2.CloseDevice(fd)

	cap, err := v4l2.GetCapability(fd)
	if err != nil {
		t.Fatalf("Failed to get capabilities: %v", err)
	}

	t.Logf("Decoder device: %s", decoderPath)
	t.Logf("  Driver: %s", cap.Driver)
	t.Logf("  Card: %s", cap.Card)

	if !cap.IsM2MSupported() {
		t.Error("Decoder should support M2M")
	}
	if !cap.IsDecoderSupported() {
		t.Error("Device should be detected as decoder")
	}

	t.Log("Decoder capabilities verified")
}

// TestCodec_EncoderCommands tests encoder command support
func TestCodec_EncoderCommands(t *testing.T) {
	encoderPath := findCodecDevice(t, true)
	if encoderPath == "" {
		t.Skip("No encoder device found - skipping test")
	}

	fd, err := v4l2.OpenDevice(encoderPath, 2, 0) // O_RDWR
	if err != nil {
		t.Fatalf("Failed to open encoder device: %v", err)
	}
	defer v4l2.CloseDevice(fd)

	// Test TRY_ENCODER_CMD for each command
	commands := []struct {
		name string
		cmd  uint32
	}{
		{"START", v4l2.EncCmdStart},
		{"STOP", v4l2.EncCmdStop},
		{"PAUSE", v4l2.EncCmdPause},
		{"RESUME", v4l2.EncCmdResume},
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			cmd := v4l2.NewEncoderCmd(tc.cmd)
			err := v4l2.TryEncoderCmd(fd, cmd)
			if err != nil {
				// Some drivers don't support TRY_ENCODER_CMD - not a failure
				t.Logf("TRY_ENCODER_CMD %s: %v (may not be supported)", tc.name, err)
			} else {
				t.Logf("TRY_ENCODER_CMD %s: supported", tc.name)
			}
		})
	}
}

// TestCodec_DecoderCommands tests decoder command support
func TestCodec_DecoderCommands(t *testing.T) {
	decoderPath := findCodecDevice(t, false)
	if decoderPath == "" {
		t.Skip("No decoder device found - skipping test")
	}

	fd, err := v4l2.OpenDevice(decoderPath, 2, 0) // O_RDWR
	if err != nil {
		t.Fatalf("Failed to open decoder device: %v", err)
	}
	defer v4l2.CloseDevice(fd)

	// Test TRY_DECODER_CMD for each command
	commands := []struct {
		name string
		cmd  uint32
	}{
		{"START", v4l2.DecCmdStart},
		{"STOP", v4l2.DecCmdStop},
		{"PAUSE", v4l2.DecCmdPause},
		{"RESUME", v4l2.DecCmdResume},
		{"FLUSH", v4l2.DecCmdFlush},
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			cmd := v4l2.NewDecoderCmd(tc.cmd)
			err := v4l2.TryDecoderCmd(fd, cmd)
			if err != nil {
				// Some drivers don't support TRY_DECODER_CMD - not a failure
				t.Logf("TRY_DECODER_CMD %s: %v (may not be supported)", tc.name, err)
			} else {
				t.Logf("TRY_DECODER_CMD %s: supported", tc.name)
			}
		})
	}
}

// TestCodec_StateMachineTransitions tests the codec state machine
func TestCodec_StateMachineTransitions(t *testing.T) {
	// This test doesn't require hardware

	t.Run("Encoder", func(t *testing.T) {
		sm := v4l2.NewCodecStateMachine(v4l2.CodecTypeEncoder)

		// Valid transition sequence
		steps := []struct {
			action string
			fn     func() error
			expect v4l2.CodecState
		}{
			{"Initialize", sm.Initialize, v4l2.CodecStateInitialized},
			{"Start", sm.Start, v4l2.CodecStateStreaming},
			{"Pause", sm.Pause, v4l2.CodecStatePaused},
			{"Resume", sm.Resume, v4l2.CodecStateStreaming},
			{"StartDrain", sm.StartDrain, v4l2.CodecStateDraining},
			{"CompleteDrain", sm.CompleteDrain, v4l2.CodecStateStopped},
		}

		for _, step := range steps {
			if err := step.fn(); err != nil {
				t.Errorf("%s failed: %v", step.action, err)
			}
			if sm.GetState() != step.expect {
				t.Errorf("After %s: expected %s, got %s", step.action, step.expect, sm.GetState())
			}
		}
	})

	t.Run("Decoder", func(t *testing.T) {
		sm := v4l2.NewCodecStateMachine(v4l2.CodecTypeDecoder)

		// Test decoder-specific flush
		sm.Initialize()
		sm.Start()

		if err := sm.StartFlush(); err != nil {
			t.Errorf("StartFlush failed: %v", err)
		}
		if sm.GetState() != v4l2.CodecStateFlushing {
			t.Errorf("Expected Flushing state, got %s", sm.GetState())
		}

		if err := sm.CompleteFlush(); err != nil {
			t.Errorf("CompleteFlush failed: %v", err)
		}
		if sm.GetState() != v4l2.CodecStateStreaming {
			t.Errorf("Expected Streaming state after flush, got %s", sm.GetState())
		}
	})
}

// TestCodec_FormatHelpers tests M2M format helper functions
func TestCodec_FormatHelpers(t *testing.T) {
	// Test with any codec device (encoder or decoder)
	codecPath := findCodecDevice(t, true)
	if codecPath == "" {
		codecPath = findCodecDevice(t, false)
	}
	if codecPath == "" {
		t.Skip("No codec device found - skipping test")
	}

	fd, err := v4l2.OpenDevice(codecPath, 2, 0) // O_RDWR
	if err != nil {
		t.Fatalf("Failed to open codec device: %v", err)
	}
	defer v4l2.CloseDevice(fd)

	t.Run("GetOutputFormat", func(t *testing.T) {
		fmt, err := v4l2.GetPixFormatOutput(fd)
		if err != nil {
			t.Logf("GetPixFormatOutput: %v (may not be configured yet)", err)
		} else {
			t.Logf("Output format: %dx%d, %s",
				fmt.Width, fmt.Height, v4l2.PixelFormats[fmt.PixelFormat])
		}
	})

	t.Run("GetCaptureFormat", func(t *testing.T) {
		fmt, err := v4l2.GetPixFormatCapture(fd)
		if err != nil {
			t.Logf("GetPixFormatCapture: %v (may not be configured yet)", err)
		} else {
			t.Logf("Capture format: %dx%d, %s",
				fmt.Width, fmt.Height, v4l2.PixelFormats[fmt.PixelFormat])
		}
	})
}

// TestCodec_EventHelpers tests codec event helper functions
func TestCodec_EventHelpers(t *testing.T) {
	codecPath := findCodecDevice(t, false) // Prefer decoder for source change events
	if codecPath == "" {
		codecPath = findCodecDevice(t, true)
	}
	if codecPath == "" {
		t.Skip("No codec device found - skipping test")
	}

	fd, err := v4l2.OpenDevice(codecPath, 2, 0) // O_RDWR
	if err != nil {
		t.Fatalf("Failed to open codec device: %v", err)
	}
	defer v4l2.CloseDevice(fd)

	t.Run("SubscribeSourceChange", func(t *testing.T) {
		err := v4l2.SubscribeSourceChangeEvent(fd)
		if err != nil {
			t.Logf("Subscribe source change: %v (may not be supported)", err)
		} else {
			t.Log("Subscribed to source change events")
			v4l2.UnsubscribeSourceChangeEvent(fd)
		}
	})

	t.Run("SubscribeEOS", func(t *testing.T) {
		err := v4l2.SubscribeEOSEvent(fd)
		if err != nil {
			t.Logf("Subscribe EOS: %v (may not be supported)", err)
		} else {
			t.Log("Subscribed to EOS events")
			v4l2.UnsubscribeEOSEvent(fd)
		}
	})
}
