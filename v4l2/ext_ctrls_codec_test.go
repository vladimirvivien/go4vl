package v4l2

import (
	"testing"
)

// Test H.264 codec helpers
func TestExtControls_H264Helpers(t *testing.T) {
	tests := []struct {
		name string
		add  func(*ExtControls) error
		get  func(*ExtControl) (any, error)
	}{
		{
			name: "H264SPS",
			add: func(ec *ExtControls) error {
				sps := &ControlH264SPS{
					Profile_idc: 100,
					Level_idc:   51,
				}
				return ec.AddH264SPS(sps)
			},
			get: func(ec *ExtControl) (any, error) {
				return ec.GetH264SPS()
			},
		},
		{
			name: "H264PPS",
			add: func(ec *ExtControls) error {
				pps := &ControlH264PPS{
					Pic_parameter_set_id: 0,
					Seq_parameter_set_id: 0,
				}
				return ec.AddH264PPS(pps)
			},
			get: func(ec *ExtControl) (any, error) {
				return ec.GetH264PPS()
			},
		},
		{
			name: "H264ScalingMatrix",
			add: func(ec *ExtControls) error {
				matrix := &ControlH264ScalingMatrix{}
				return ec.AddH264ScalingMatrix(matrix)
			},
			get: func(ec *ExtControl) (any, error) {
				return ec.GetH264ScalingMatrix()
			},
		},
		{
			name: "H264SliceParams",
			add: func(ec *ExtControls) error {
				params := &ControlH264SliceParams{
					Header_bit_size:   100,
					First_mb_in_slice: 0,
				}
				return ec.AddH264SliceParams(params)
			},
			get: func(ec *ExtControl) (any, error) {
				return ec.GetH264SliceParams()
			},
		},
		{
			name: "H264DecodeParams",
			add: func(ec *ExtControls) error {
				params := &ControlH264DecodeParams{
					Frame_num:   1,
					Nal_ref_idc: 3,
				}
				return ec.AddH264DecodeParams(params)
			},
			get: func(ec *ExtControl) (any, error) {
				return ec.GetH264DecodeParams()
			},
		},
		{
			name: "H264PredWeights",
			add: func(ec *ExtControls) error {
				weights := &ControlH264PredictionWeights{
					Luma_log2_weight_denom:   5,
					Chroma_log2_weight_denom: 5,
				}
				return ec.AddH264PredWeights(weights)
			},
			get: func(ec *ExtControl) (any, error) {
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
		get  func(*ExtControl) (any, error)
	}{
		{
			name: "MPEG2Sequence",
			add: func(ec *ExtControls) error {
				seq := &ControlMPEG2Sequence{
					Horizontal_size: 1920,
					Vertical_size:   1080,
				}
				return ec.AddMPEG2Sequence(seq)
			},
			get: func(ec *ExtControl) (any, error) {
				return ec.GetMPEG2Sequence()
			},
		},
		{
			name: "MPEG2Picture",
			add: func(ec *ExtControls) error {
				pic := &ControlMPEG2Picture{
					Backward_ref_ts: 1000,
					Forward_ref_ts:  2000,
				}
				return ec.AddMPEG2Picture(pic)
			},
			get: func(ec *ExtControl) (any, error) {
				return ec.GetMPEG2Picture()
			},
		},
		{
			name: "MPEG2Quantization",
			add: func(ec *ExtControls) error {
				quant := &ControlMPEG2Quantization{}
				return ec.AddMPEG2Quantization(quant)
			},
			get: func(ec *ExtControl) (any, error) {
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
		Width:          1920,
		Height:         1080,
		Version:        0,
		Dct_part_sizes: [8]uint32{128, 256},
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

	if retrieved.Width != 1920 || retrieved.Height != 1080 || retrieved.Dct_part_sizes[1] != 256 {
		t.Errorf("VP8 frame data mismatch: got %dx%d (dct[1]=%d)",
			retrieved.Width, retrieved.Height, retrieved.Dct_part_sizes[1])
	}
}

// Test FWHT codec helpers
func TestExtControls_FWHTHelpers(t *testing.T) {
	ctrls := NewExtControls()

	params := &ControlFWHTParams{
		Width:        1920,
		Height:       1080,
		Version:      1,
		Quantization: 2,
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

	if retrieved.Width != 1920 || retrieved.Height != 1080 || retrieved.Version != 1 || retrieved.Quantization != 2 {
		t.Errorf("FWHT params data mismatch: %+v", retrieved)
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
	sps := &ControlH264SPS{Profile_idc: 100}
	if err := ctrls.AddH264SPS(sps); err != nil {
		t.Fatalf("Failed to add H264 SPS: %v", err)
	}

	// Add H.264 PPS
	pps := &ControlH264PPS{Pic_parameter_set_id: 0}
	if err := ctrls.AddH264PPS(pps); err != nil {
		t.Fatalf("Failed to add H264 PPS: %v", err)
	}

	// Add MPEG2 sequence
	seq := &ControlMPEG2Sequence{Horizontal_size: 1920}
	if err := ctrls.AddMPEG2Sequence(seq); err != nil {
		t.Fatalf("Failed to add MPEG2 sequence: %v", err)
	}

	// Verify all were added
	if len(ctrls.GetControls()) != 3 {
		t.Errorf("Expected 3 controls, got %d", len(ctrls.GetControls()))
	}
}

func TestExtControls_GenericStructHelpers(t *testing.T) {
	ctrls := NewExtControls()
	in := &ControlH264SPS{
		Profile_idc: 100,
		Level_idc:   41,
		Flags:       0x03,
	}

	if err := AddExtControlStruct(ctrls, CtrlH264SPS, in); err != nil {
		t.Fatalf("AddExtControlStruct() = %v", err)
	}

	got, err := GetExtControlStruct[ControlH264SPS](ctrls, CtrlH264SPS)
	if err != nil {
		t.Fatalf("GetExtControlStruct() = %v", err)
	}
	if got.Profile_idc != in.Profile_idc || got.Level_idc != in.Level_idc || got.Flags != in.Flags {
		t.Fatalf("round-trip mismatch: got %+v, want %+v", got, in)
	}

	if err := AddExtControlStruct[ControlH264SPS](nil, CtrlH264SPS, in); err == nil {
		t.Fatalf("expected error for nil ExtControls in AddExtControlStruct")
	}
	if err := AddExtControlStruct[ControlH264SPS](ctrls, CtrlH264SPS, nil); err == nil {
		t.Fatalf("expected error for nil value in AddExtControlStruct")
	}
	if _, err := GetExtControlStruct[ControlH264SPS](nil, CtrlH264SPS); err == nil {
		t.Fatalf("expected error for nil ExtControls in GetExtControlStruct")
	}
	if _, err := GetExtControlStruct[ControlH264PPS](ctrls, CtrlH264PPS); err == nil {
		t.Fatalf("expected error for missing control ID in GetExtControlStruct")
	}

	bad := NewExtControls()
	bad.AddCompound(CtrlH264SPS, []byte{1, 2, 3})
	if _, err := GetExtControlStruct[ControlH264SPS](bad, CtrlH264SPS); err == nil {
		t.Fatalf("expected error for truncated compound payload in GetExtControlStruct")
	}
}
