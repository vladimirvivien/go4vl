package v4l2

import "testing"

func TestExtControls_StatelessHEVC(t *testing.T) {
	ctrls := NewExtControls()
	ctrls.SetRequestFD(7)

	if err := ctrls.AddHEVCSPS(nil); err == nil {
		t.Fatal("expected error for nil HEVC SPS")
	}

	sps := &ControlHEVCSPS{Pic_width_in_luma_samples: 1920, Pic_height_in_luma_samples: 1080}
	pps := &ControlHEVCPPS{Init_qp_minus26: -4}
	slice := &ControlHEVCSliceParams{Bit_size: 4096, Nal_unit_type: 19}
	matrix := &ControlHEVCScalingMatrix{}
	matrix.List_4x4[0][0] = 16
	decode := &ControlHEVCDecodeParams{Pic_order_cnt_val: 1, Num_active_dpb_entries: 1}

	if err := ctrls.AddHEVCSPS(sps); err != nil {
		t.Fatalf("AddHEVCSPS: %v", err)
	}
	if err := ctrls.AddHEVCPPS(pps); err != nil {
		t.Fatalf("AddHEVCPPS: %v", err)
	}
	if err := ctrls.AddHEVCSliceParams(slice); err != nil {
		t.Fatalf("AddHEVCSliceParams: %v", err)
	}
	if err := ctrls.AddHEVCScalingMatrix(matrix); err != nil {
		t.Fatalf("AddHEVCScalingMatrix: %v", err)
	}
	if err := ctrls.AddHEVCDecodeParams(decode); err != nil {
		t.Fatalf("AddHEVCDecodeParams: %v", err)
	}

	list := ctrls.GetControls()
	if len(list) != 5 {
		t.Fatalf("expected 5 HEVC controls, got %d", len(list))
	}

	gotSPS, err := list[0].GetHEVCSPS()
	if err != nil || gotSPS.Pic_width_in_luma_samples != 1920 {
		t.Fatalf("GetHEVCSPS = %+v, %v", gotSPS, err)
	}
	gotPPS, err := list[1].GetHEVCPPS()
	if err != nil || gotPPS.Init_qp_minus26 != -4 {
		t.Fatalf("GetHEVCPPS = %+v, %v", gotPPS, err)
	}
	gotSlice, err := list[2].GetHEVCSliceParams()
	if err != nil || gotSlice.Bit_size != 4096 {
		t.Fatalf("GetHEVCSliceParams = %+v, %v", gotSlice, err)
	}
	gotMatrix, err := list[3].GetHEVCScalingMatrix()
	if err != nil || gotMatrix.List_4x4[0][0] != 16 {
		t.Fatalf("GetHEVCScalingMatrix = %+v, %v", gotMatrix, err)
	}
	gotDecode, err := list[4].GetHEVCDecodeParams()
	if err != nil || gotDecode.Pic_order_cnt_val != 1 {
		t.Fatalf("GetHEVCDecodeParams = %+v, %v", gotDecode, err)
	}
}

func TestExtControls_StatelessVP9AndAV1(t *testing.T) {
	ctrls := NewExtControls()

	vp9Frame := &ControlVP9Frame{Frame_width_minus_1: 1279, Frame_height_minus_1: 719, Profile: 0}
	vp9Hdr := &ControlVP9CompressedHdr{Tx_mode: 2}
	av1Seq := &ControlAV1Sequence{Seq_profile: 0, Bit_depth: 8, Max_frame_width_minus_1: 1919}
	av1Tile := &ControlAV1TileGroupEntry{Offset: 12, Size: 512, Row: 0, Col: 1}
	av1Frame := &ControlAV1Frame{Frame_width_minus_1: 1919, Frame_height_minus_1: 1079}
	av1Grain := &ControlAV1FilmGrain{Grain_seed: 1234}

	if err := ctrls.AddVP9Frame(vp9Frame); err != nil {
		t.Fatalf("AddVP9Frame: %v", err)
	}
	if err := ctrls.AddVP9CompressedHdr(vp9Hdr); err != nil {
		t.Fatalf("AddVP9CompressedHdr: %v", err)
	}
	if err := ctrls.AddAV1Sequence(av1Seq); err != nil {
		t.Fatalf("AddAV1Sequence: %v", err)
	}
	if err := ctrls.AddAV1TileGroupEntry(av1Tile); err != nil {
		t.Fatalf("AddAV1TileGroupEntry: %v", err)
	}
	if err := ctrls.AddAV1Frame(av1Frame); err != nil {
		t.Fatalf("AddAV1Frame: %v", err)
	}
	if err := ctrls.AddAV1FilmGrain(av1Grain); err != nil {
		t.Fatalf("AddAV1FilmGrain: %v", err)
	}

	list := ctrls.GetControls()
	if len(list) != 6 {
		t.Fatalf("expected 6 controls, got %d", len(list))
	}

	if f, err := list[0].GetVP9Frame(); err != nil || f.Frame_width_minus_1 != 1279 {
		t.Fatalf("GetVP9Frame = %+v, %v", f, err)
	}
	if h, err := list[1].GetVP9CompressedHdr(); err != nil || h.Tx_mode != 2 {
		t.Fatalf("GetVP9CompressedHdr = %+v, %v", h, err)
	}
	if s, err := list[2].GetAV1Sequence(); err != nil || s.Max_frame_width_minus_1 != 1919 {
		t.Fatalf("GetAV1Sequence = %+v, %v", s, err)
	}
	if te, err := list[3].GetAV1TileGroupEntry(); err != nil || te.Size != 512 {
		t.Fatalf("GetAV1TileGroupEntry = %+v, %v", te, err)
	}
	if af, err := list[4].GetAV1Frame(); err != nil || af.Frame_height_minus_1 != 1079 {
		t.Fatalf("GetAV1Frame = %+v, %v", af, err)
	}
	if fg, err := list[5].GetAV1FilmGrain(); err != nil || fg.Grain_seed != 1234 {
		t.Fatalf("GetAV1FilmGrain = %+v, %v", fg, err)
	}

	// Verify short compound data returns an error.
	short := NewExtControls()
	short.AddCompound(CtrlAV1Sequence, []byte{1, 2})
	if _, err := short.GetControls()[0].GetAV1Sequence(); err == nil {
		t.Fatal("expected error for truncated compound control data")
	}
}
