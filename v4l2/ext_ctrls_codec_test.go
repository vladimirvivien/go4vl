package v4l2

import (
	"testing"
)

// Test H.264 codec helpers
func TestExtControls_H264Helpers(t *testing.T) {
	tests := []struct {
		name string
		add  func(*ExtControls) error
		get  func(*ExtControl) (interface{}, error)
	}{
		{
			name: "H264SPS",
			add: func(ec *ExtControls) error {
				sps := &ControlH264SPS{
					ProfileIDC: 100,
					LevelIDC:   51,
				}
				return ec.AddH264SPS(sps)
			},
			get: func(ec *ExtControl) (interface{}, error) {
				return ec.GetH264SPS()
			},
		},
		{
			name: "H264PPS",
			add: func(ec *ExtControls) error {
				pps := &ControlH264PPS{
					PicParameterSetID: 0,
					SeqParameterSetID: 0,
				}
				return ec.AddH264PPS(pps)
			},
			get: func(ec *ExtControl) (interface{}, error) {
				return ec.GetH264PPS()
			},
		},
		{
			name: "H264ScalingMatrix",
			add: func(ec *ExtControls) error {
				matrix := &ControlH264ScalingMatrix{}
				return ec.AddH264ScalingMatrix(matrix)
			},
			get: func(ec *ExtControl) (interface{}, error) {
				return ec.GetH264ScalingMatrix()
			},
		},
		{
			name: "H264SliceParams",
			add: func(ec *ExtControls) error {
				params := &ControlH264SliceParams{
					HeaderBitSize:  100,
					FirstMBInSlice: 0,
				}
				return ec.AddH264SliceParams(params)
			},
			get: func(ec *ExtControl) (interface{}, error) {
				return ec.GetH264SliceParams()
			},
		},
		{
			name: "H264DecodeParams",
			add: func(ec *ExtControls) error {
				params := &ControlH264DecodeParams{
					FrameNum: 1,
					NalRefIDC: 3,
				}
				return ec.AddH264DecodeParams(params)
			},
			get: func(ec *ExtControl) (interface{}, error) {
				return ec.GetH264DecodeParams()
			},
		},
		{
			name: "H264PredWeights",
			add: func(ec *ExtControls) error {
				weights := &ControlH264PredictionWeights{
					LumaLog2WeightDenom:   5,
					ChromaLog2WeightDenom: 5,
				}
				return ec.AddH264PredWeights(weights)
			},
			get: func(ec *ExtControl) (interface{}, error) {
				return ec.GetH264PredWeights()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrls := NewExtControls()

			// Test adding the control
			err := tt.add(ctrls)
			if err != nil {
				t.Fatalf("Failed to add control: %v", err)
			}

			// Verify it was added
			if len(ctrls.GetControls()) != 1 {
				t.Fatalf("Expected 1 control, got %d", len(ctrls.GetControls()))
			}

			// Test getting it back (simulated - would need actual device for real test)
			ctrl := ctrls.GetControls()[0]
			_, err = tt.get(ctrl)
			if err != nil {
				t.Fatalf("Failed to get control: %v", err)
			}
		})
	}
}

// Test MPEG2 codec helpers
func TestExtControls_MPEG2Helpers(t *testing.T) {
	tests := []struct {
		name string
		add  func(*ExtControls) error
		get  func(*ExtControl) (interface{}, error)
	}{
		{
			name: "MPEG2Sequence",
			add: func(ec *ExtControls) error {
				seq := &ControlMPEG2Sequence{
					HorizontalSize: 1920,
					VerticalSize:   1080,
				}
				return ec.AddMPEG2Sequence(seq)
			},
			get: func(ec *ExtControl) (interface{}, error) {
				return ec.GetMPEG2Sequence()
			},
		},
		{
			name: "MPEG2Picture",
			add: func(ec *ExtControls) error {
				pic := &ControlMPEG2Picture{
					BackwardRefTimestamp: 1000,
					ForwardRefTimestamp:  2000,
				}
				return ec.AddMPEG2Picture(pic)
			},
			get: func(ec *ExtControl) (interface{}, error) {
				return ec.GetMPEG2Picture()
			},
		},
		{
			name: "MPEG2Quantization",
			add: func(ec *ExtControls) error {
				quant := &ControlMPEG2Quantization{}
				return ec.AddMPEG2Quantization(quant)
			},
			get: func(ec *ExtControl) (interface{}, error) {
				return ec.GetMPEG2Quantization()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrls := NewExtControls()

			err := tt.add(ctrls)
			if err != nil {
				t.Fatalf("Failed to add control: %v", err)
			}

			if len(ctrls.GetControls()) != 1 {
				t.Fatalf("Expected 1 control, got %d", len(ctrls.GetControls()))
			}

			ctrl := ctrls.GetControls()[0]
			_, err = tt.get(ctrl)
			if err != nil {
				t.Fatalf("Failed to get control: %v", err)
			}
		})
	}
}

// Test VP8 codec helpers
func TestExtControls_VP8Helpers(t *testing.T) {
	ctrls := NewExtControls()

	frame := &ControlVP8Frame{
		Width:  1920,
		Height: 1080,
		Version: 0,
	}

	err := ctrls.AddVP8Frame(frame)
	if err != nil {
		t.Fatalf("Failed to add VP8 frame: %v", err)
	}

	if len(ctrls.GetControls()) != 1 {
		t.Fatalf("Expected 1 control, got %d", len(ctrls.GetControls()))
	}

	ctrl := ctrls.GetControls()[0]
	retrieved, err := ctrl.GetVP8Frame()
	if err != nil {
		t.Fatalf("Failed to get VP8 frame: %v", err)
	}

	if retrieved.Width != 1920 || retrieved.Height != 1080 {
		t.Errorf("VP8 frame data mismatch: got %dx%d, expected 1920x1080",
			retrieved.Width, retrieved.Height)
	}
}

// Test FWHT codec helpers
func TestExtControls_FWHTHelpers(t *testing.T) {
	ctrls := NewExtControls()

	params := &ControlFWHTParams{
		Width:   1920,
		Height:  1080,
		Version: 1,
	}

	err := ctrls.AddFWHTParams(params)
	if err != nil {
		t.Fatalf("Failed to add FWHT params: %v", err)
	}

	if len(ctrls.GetControls()) != 1 {
		t.Fatalf("Expected 1 control, got %d", len(ctrls.GetControls()))
	}

	ctrl := ctrls.GetControls()[0]
	retrieved, err := ctrl.GetFWHTParams()
	if err != nil {
		t.Fatalf("Failed to get FWHT params: %v", err)
	}

	if retrieved.Width != 1920 || retrieved.Height != 1080 || retrieved.Version != 1 {
		t.Errorf("FWHT params data mismatch")
	}
}

// Test nil parameter handling
func TestExtControls_CodecHelpers_NilCheck(t *testing.T) {
	ctrls := NewExtControls()

	tests := []struct {
		name string
		fn   func() error
	}{
		{"H264SPS", func() error { return ctrls.AddH264SPS(nil) }},
		{"H264PPS", func() error { return ctrls.AddH264PPS(nil) }},
		{"H264ScalingMatrix", func() error { return ctrls.AddH264ScalingMatrix(nil) }},
		{"H264SliceParams", func() error { return ctrls.AddH264SliceParams(nil) }},
		{"H264DecodeParams", func() error { return ctrls.AddH264DecodeParams(nil) }},
		{"H264PredWeights", func() error { return ctrls.AddH264PredWeights(nil) }},
		{"MPEG2Sequence", func() error { return ctrls.AddMPEG2Sequence(nil) }},
		{"MPEG2Picture", func() error { return ctrls.AddMPEG2Picture(nil) }},
		{"MPEG2Quantization", func() error { return ctrls.AddMPEG2Quantization(nil) }},
		{"VP8Frame", func() error { return ctrls.AddVP8Frame(nil) }},
		{"FWHTParams", func() error { return ctrls.AddFWHTParams(nil) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Errorf("Expected error for nil parameter, got nil")
			}
		})
	}
}

// Test multiple codec controls in one collection
func TestExtControls_MultipleCodecControls(t *testing.T) {
	ctrls := NewExtControls()

	// Add H.264 SPS
	sps := &ControlH264SPS{ProfileIDC: 100}
	if err := ctrls.AddH264SPS(sps); err != nil {
		t.Fatalf("Failed to add H264 SPS: %v", err)
	}

	// Add H.264 PPS
	pps := &ControlH264PPS{PicParameterSetID: 0}
	if err := ctrls.AddH264PPS(pps); err != nil {
		t.Fatalf("Failed to add H264 PPS: %v", err)
	}

	// Add MPEG2 sequence
	seq := &ControlMPEG2Sequence{HorizontalSize: 1920}
	if err := ctrls.AddMPEG2Sequence(seq); err != nil {
		t.Fatalf("Failed to add MPEG2 sequence: %v", err)
	}

	// Verify all were added
	if len(ctrls.GetControls()) != 3 {
		t.Errorf("Expected 3 controls, got %d", len(ctrls.GetControls()))
	}
}
