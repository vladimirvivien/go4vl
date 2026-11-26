//go:build cgo_verify

package v4l2

/*
#include "internal/gen/compat.h"
#include <linux/media.h>
#include <linux/v4l2-controls.h>
#include <linux/v4l2-subdev.h>
#include <linux/media-bus-format.h>
#include <stddef.h>
#include <sys/time.h>
#include <time.h>

static size_t offsetof_v4l2_ext_control_value(void) {
	return offsetof(struct v4l2_ext_control, value);
}
static size_t offsetof_v4l2_format_fmt(void) {
	return offsetof(struct v4l2_format, fmt);
}
static size_t offsetof_v4l2_pix_format_ycbcr_enc(void) {
	return offsetof(struct v4l2_pix_format, ycbcr_enc);
}
static size_t offsetof_v4l2_buffer_m(void) {
	return offsetof(struct v4l2_buffer, m);
}
static size_t offsetof_v4l2_buffer_request_fd(void) {
	return offsetof(struct v4l2_buffer, request_fd);
}
static size_t offsetof_v4l2_event_u(void) {
	return offsetof(struct v4l2_event, u);
}
static size_t offsetof_v4l2_event_ctrl_value(void) {
	return offsetof(struct v4l2_event_ctrl, value);
}
static size_t offsetof_v4l2_frmsizeenum_discrete(void) {
	return offsetof(struct v4l2_frmsizeenum, discrete);
}
static size_t offsetof_v4l2_frmivalenum_discrete(void) {
	return offsetof(struct v4l2_frmivalenum, discrete);
}
static size_t offsetof_v4l2_streamparm_parm(void) {
	return offsetof(struct v4l2_streamparm, parm);
}
static size_t offsetof_v4l2_dv_timings_bt(void) {
	return offsetof(struct v4l2_dv_timings, bt);
}
static size_t offsetof_v4l2_dv_timings_cap_bt(void) {
	return offsetof(struct v4l2_dv_timings_cap, bt);
}
static size_t offsetof_v4l2_decoder_cmd_stop(void) {
	return offsetof(struct v4l2_decoder_cmd, stop);
}
static size_t offsetof_v4l2_querymenu_name(void) {
	return offsetof(struct v4l2_querymenu, name);
}
static size_t offsetof_v4l2_plane_m(void) {
	return offsetof(struct v4l2_plane, m);
}
static size_t offsetof_v4l2_plane_data_offset(void) {
	return offsetof(struct v4l2_plane, data_offset);
}
static size_t offsetof_v4l2_pix_format_mplane_num_planes(void) {
	return offsetof(struct v4l2_pix_format_mplane, num_planes);
}
static size_t offsetof_v4l2_create_buffers_format(void) {
	return offsetof(struct v4l2_create_buffers, format);
}
static size_t offsetof_v4l2_selection_r(void) {
	return offsetof(struct v4l2_selection, r);
}
static size_t offsetof_media_v2_topology_ptr_links(void) {
	return offsetof(struct media_v2_topology, ptr_links);
}
static size_t offsetof_media_v2_interface_devnode(void) {
	return offsetof(struct media_v2_interface, devnode);
}
static size_t offsetof_media_entity_desc_dev(void) {
	return offsetof(struct media_entity_desc, dev);
}
static size_t offsetof_v4l2_subdev_format_format(void) {
	return offsetof(struct v4l2_subdev_format, format);
}
static size_t offsetof_v4l2_subdev_selection_r(void) {
	return offsetof(struct v4l2_subdev_selection, r);
}
static size_t offsetof_v4l2_meta_format_buffersize(void) {
	return offsetof(struct v4l2_meta_format, buffersize);
}
static size_t offsetof_v4l2_edid_edid(void) {
	return offsetof(struct v4l2_edid, edid);
}
static size_t offsetof_v4l2_subdev_routing_routes(void) {
	return offsetof(struct v4l2_subdev_routing, routes);
}
static size_t offsetof_v4l2_ctrl_h264_sps_flags(void) {
	return offsetof(struct v4l2_ctrl_h264_sps, flags);
}
static size_t offsetof_v4l2_h264_dpb_entry_top_field_order_cnt(void) {
	return offsetof(struct v4l2_h264_dpb_entry, top_field_order_cnt);
}
static size_t offsetof_v4l2_h264_dpb_entry_flags(void) {
	return offsetof(struct v4l2_h264_dpb_entry, flags);
}
static size_t offsetof_v4l2_ctrl_h264_decode_params_flags(void) {
	return offsetof(struct v4l2_ctrl_h264_decode_params, flags);
}
static size_t offsetof_v4l2_ctrl_fwht_params_quantization(void) {
	return offsetof(struct v4l2_ctrl_fwht_params, quantization);
}
static size_t offsetof_v4l2_vp8_loop_filter_flags(void) {
	return offsetof(struct v4l2_vp8_loop_filter, flags);
}
static size_t offsetof_v4l2_vp8_entropy_mv_probs(void) {
	return offsetof(struct v4l2_vp8_entropy, mv_probs);
}
static size_t offsetof_v4l2_ctrl_vp8_frame_dct_part_sizes(void) {
	return offsetof(struct v4l2_ctrl_vp8_frame, dct_part_sizes);
}
static size_t offsetof_v4l2_ctrl_vp8_frame_flags(void) {
	return offsetof(struct v4l2_ctrl_vp8_frame, flags);
}
static size_t offsetof_v4l2_ctrl_mpeg2_picture_picture_coding_type(void) {
	return offsetof(struct v4l2_ctrl_mpeg2_picture, picture_coding_type);
}
*/
import "C"

type abiPair struct {
	name string
	got  uintptr
	want uintptr
}

func cStructSizeChecks(goSizes map[string]uintptr) []abiPair {
	cSizes := map[string]uintptr{
		"cTimeval":             uintptr(C.sizeof_struct_timeval),
		"sys.Timeval":          uintptr(C.sizeof_struct_timeval),
		"cTimespec":            uintptr(C.sizeof_struct_timespec),
		"cV4L2Timecode":        uintptr(C.sizeof_struct_v4l2_timecode),
		"Timecode":             uintptr(C.sizeof_struct_v4l2_timecode),
		"cV4L2Fract":           uintptr(C.sizeof_struct_v4l2_fract),
		"Fract":                uintptr(C.sizeof_struct_v4l2_fract),
		"cV4L2Rect":            uintptr(C.sizeof_struct_v4l2_rect),
		"Rect":                 uintptr(C.sizeof_struct_v4l2_rect),
		"cMediaDeviceInfo":             uintptr(C.sizeof_struct_media_device_info),
		"cMediaV2Topology":             uintptr(C.sizeof_struct_media_v2_topology),
		"cMediaV2Entity":               uintptr(C.sizeof_struct_media_v2_entity),
		"cMediaV2Interface":            uintptr(C.sizeof_struct_media_v2_interface),
		"cMediaV2IntfDevnode":          uintptr(C.sizeof_struct_media_v2_intf_devnode),
		"cMediaV2Pad":                  uintptr(C.sizeof_struct_media_v2_pad),
		"cMediaV2Link":                 uintptr(C.sizeof_struct_media_v2_link),
		"cMediaEntityDesc":             uintptr(C.sizeof_struct_media_entity_desc),
		"cMediaPadDesc":                uintptr(C.sizeof_struct_media_pad_desc),
		"cMediaLinkDesc":               uintptr(C.sizeof_struct_media_link_desc),
		"cMediaLinksEnum":              uintptr(C.sizeof_struct_media_links_enum),
		"cV4L2MBusFrameFmt":            uintptr(C.sizeof_struct_v4l2_mbus_framefmt),
		"cV4L2SubdevFormat":            uintptr(C.sizeof_struct_v4l2_subdev_format),
		"cV4L2SubdevCrop":              uintptr(C.sizeof_struct_v4l2_subdev_crop),
		"cV4L2SubdevMBusCodeEnum":      uintptr(C.sizeof_struct_v4l2_subdev_mbus_code_enum),
		"cV4L2SubdevFrameSizeEnum":     uintptr(C.sizeof_struct_v4l2_subdev_frame_size_enum),
		"cV4L2SubdevFrameInterval":     uintptr(C.sizeof_struct_v4l2_subdev_frame_interval),
		"cV4L2SubdevFrameIntervalEnum": uintptr(C.sizeof_struct_v4l2_subdev_frame_interval_enum),
		"cV4L2SubdevSelection":         uintptr(C.sizeof_struct_v4l2_subdev_selection),
		"cV4L2SubdevCapability":        uintptr(C.sizeof_struct_v4l2_subdev_capability),
		"cV4L2SubdevClientCapability":  uintptr(C.sizeof_struct_v4l2_subdev_client_capability),
		"cV4L2SubdevRoute":             uintptr(C.sizeof_struct_v4l2_subdev_route),
		"cV4L2SubdevRouting":           uintptr(C.sizeof_struct_v4l2_subdev_routing),
		"cV4L2EDID":                    uintptr(C.sizeof_struct_v4l2_edid),
		"cV4L2Capability":      uintptr(C.sizeof_struct_v4l2_capability),
		"cV4L2Audio":           uintptr(C.sizeof_struct_v4l2_audio),
		"cV4L2AudioOut":        uintptr(C.sizeof_struct_v4l2_audioout),
		"cV4L2BTTimings":       uintptr(C.sizeof_struct_v4l2_bt_timings),
		"cV4L2BTTimingsCap":    uintptr(C.sizeof_struct_v4l2_bt_timings_cap),
		"cV4L2Buffer":          uintptr(C.sizeof_struct_v4l2_buffer),
		"cV4L2CaptureParm":     uintptr(C.sizeof_struct_v4l2_captureparm),
		"CaptureParam":         uintptr(C.sizeof_struct_v4l2_captureparm),
		"cV4L2Control":         uintptr(C.sizeof_struct_v4l2_control),
		"cV4L2CreateBuffers":   uintptr(C.sizeof_struct_v4l2_create_buffers),
		"cV4L2Crop":            uintptr(C.sizeof_struct_v4l2_crop),
		"cV4L2CropCap":         uintptr(C.sizeof_struct_v4l2_cropcap),
		"cV4L2DecoderCmd":      uintptr(C.sizeof_struct_v4l2_decoder_cmd),
		"cV4L2DVTimings":       uintptr(C.sizeof_struct_v4l2_dv_timings),
		"cV4L2DVTimingsCap":    uintptr(C.sizeof_struct_v4l2_dv_timings_cap),
		"cV4L2EncoderCmd":      uintptr(C.sizeof_struct_v4l2_encoder_cmd),
		"cV4L2EnumDVTimings":   uintptr(C.sizeof_struct_v4l2_enum_dv_timings),
		"cV4L2EventCtrl":       uintptr(C.sizeof_struct_v4l2_event_ctrl),
		"cV4L2EventFrameSync":  uintptr(C.sizeof_struct_v4l2_event_frame_sync),
		"cV4L2EventMotionDet":  uintptr(C.sizeof_struct_v4l2_event_motion_det),
		"cV4L2EventSrcChange":  uintptr(C.sizeof_struct_v4l2_event_src_change),
		"cV4L2EventVsync":      uintptr(C.sizeof_struct_v4l2_event_vsync),
		"cV4L2Event":           uintptr(C.sizeof_struct_v4l2_event),
		"cV4L2EventSubscription": uintptr(C.sizeof_struct_v4l2_event_subscription),
		"cV4L2ExportBuffer":    uintptr(C.sizeof_struct_v4l2_exportbuffer),
		"cV4L2ExtControl":      uintptr(C.sizeof_struct_v4l2_ext_control),
		"cV4L2ExtControls":     uintptr(C.sizeof_struct_v4l2_ext_controls),
		"cV4L2FmtDesc":         uintptr(C.sizeof_struct_v4l2_fmtdesc),
		"cV4L2Format":          uintptr(C.sizeof_struct_v4l2_format),
		"cV4L2Frequency":       uintptr(C.sizeof_struct_v4l2_frequency),
		"cV4L2FrequencyBand":   uintptr(C.sizeof_struct_v4l2_frequency_band),
		"cV4L2FrmSizeDiscrete": uintptr(C.sizeof_struct_v4l2_frmsize_discrete),
		"FrameSizeDiscrete":    uintptr(C.sizeof_struct_v4l2_frmsize_discrete),
		"cV4L2FrmSizeStepwise": uintptr(C.sizeof_struct_v4l2_frmsize_stepwise),
		"FrameSize":            uintptr(C.sizeof_struct_v4l2_frmsize_stepwise),
		"cV4L2FrmSizeEnum":     uintptr(C.sizeof_struct_v4l2_frmsizeenum),
		"cV4L2FrmIvalStepwise": uintptr(C.sizeof_struct_v4l2_frmival_stepwise),
		"FrameInterval":        uintptr(C.sizeof_struct_v4l2_frmival_stepwise),
		"cV4L2FrmIvalEnum":     uintptr(C.sizeof_struct_v4l2_frmivalenum),
		"cV4L2Input":           uintptr(C.sizeof_struct_v4l2_input),
		"cV4L2Modulator":       uintptr(C.sizeof_struct_v4l2_modulator),
		"cV4L2Output":          uintptr(C.sizeof_struct_v4l2_output),
		"cV4L2OutputParm":      uintptr(C.sizeof_struct_v4l2_outputparm),
		"OutputParam":          uintptr(C.sizeof_struct_v4l2_outputparm),
		"cV4L2PixFormat":       uintptr(C.sizeof_struct_v4l2_pix_format),
		"cV4L2PlanePixFormat":  uintptr(C.sizeof_struct_v4l2_plane_pix_format),
		"cV4L2PixFormatMPlane": uintptr(C.sizeof_struct_v4l2_pix_format_mplane),
		"cV4L2Plane":           uintptr(C.sizeof_struct_v4l2_plane),
		"cV4L2QueryCtrl":       uintptr(C.sizeof_struct_v4l2_queryctrl),
		"cV4L2QueryExtCtrl":    uintptr(C.sizeof_struct_v4l2_query_ext_ctrl),
		"cV4L2QueryMenu":       uintptr(C.sizeof_struct_v4l2_querymenu),
		"cV4L2RequestBuffers":  uintptr(C.sizeof_struct_v4l2_requestbuffers),
		"RequestBuffers":       uintptr(C.sizeof_struct_v4l2_requestbuffers),
		"cV4L2RemoveBuffers":   uintptr(C.sizeof_struct_v4l2_remove_buffers),
		"cV4L2Selection":       uintptr(C.sizeof_struct_v4l2_selection),
		"cV4L2Standard":        uintptr(C.sizeof_struct_v4l2_standard),
		"cV4L2StreamParm":            uintptr(C.sizeof_struct_v4l2_streamparm),
		"cV4L2Tuner":                 uintptr(C.sizeof_struct_v4l2_tuner),
		"cV4L2CtrlHEVCSPS":           uintptr(C.sizeof_struct_v4l2_ctrl_hevc_sps),
		"cV4L2CtrlHEVCPPS":           uintptr(C.sizeof_struct_v4l2_ctrl_hevc_pps),
		"cV4L2CtrlHEVCSliceParams":   uintptr(C.sizeof_struct_v4l2_ctrl_hevc_slice_params),
		"cV4L2CtrlHEVCScalingMatrix": uintptr(C.sizeof_struct_v4l2_ctrl_hevc_scaling_matrix),
		"cV4L2CtrlHEVCDecodeParams":  uintptr(C.sizeof_struct_v4l2_ctrl_hevc_decode_params),
		"cV4L2CtrlVP9Frame":          uintptr(C.sizeof_struct_v4l2_ctrl_vp9_frame),
		"cV4L2CtrlVP9CompressedHdr":  uintptr(C.sizeof_struct_v4l2_ctrl_vp9_compressed_hdr),
		"cV4L2CtrlAV1Sequence":       uintptr(C.sizeof_struct_v4l2_ctrl_av1_sequence),
		"cV4L2CtrlAV1TileGroupEntry": uintptr(C.sizeof_struct_v4l2_ctrl_av1_tile_group_entry),
		"cV4L2CtrlAV1Frame":          uintptr(C.sizeof_struct_v4l2_ctrl_av1_frame),
		"cV4L2CtrlAV1FilmGrain":      uintptr(C.sizeof_struct_v4l2_ctrl_av1_film_grain),
		"cV4L2CtrlH264SPS":           uintptr(C.sizeof_struct_v4l2_ctrl_h264_sps),
		"cV4L2CtrlH264PPS":           uintptr(C.sizeof_struct_v4l2_ctrl_h264_pps),
		"cV4L2CtrlH264ScalingMatrix": uintptr(C.sizeof_struct_v4l2_ctrl_h264_scaling_matrix),
		"cV4L2H264WeightFactors":     uintptr(C.sizeof_struct_v4l2_h264_weight_factors),
		"cV4L2CtrlH264PredWeights":   uintptr(C.sizeof_struct_v4l2_ctrl_h264_pred_weights),
		"cV4L2H264Reference":         uintptr(C.sizeof_struct_v4l2_h264_reference),
		"cV4L2CtrlH264SliceParams":   uintptr(C.sizeof_struct_v4l2_ctrl_h264_slice_params),
		"cV4L2H264DPBEntry":          uintptr(C.sizeof_struct_v4l2_h264_dpb_entry),
		"cV4L2CtrlH264DecodeParams":  uintptr(C.sizeof_struct_v4l2_ctrl_h264_decode_params),
		"cV4L2CtrlFWHTParams":        uintptr(C.sizeof_struct_v4l2_ctrl_fwht_params),
		"cV4L2VP8Segment":            uintptr(C.sizeof_struct_v4l2_vp8_segment),
		"cV4L2VP8LoopFilter":         uintptr(C.sizeof_struct_v4l2_vp8_loop_filter),
		"cV4L2VP8Quantization":       uintptr(C.sizeof_struct_v4l2_vp8_quantization),
		"cV4L2VP8Entropy":            uintptr(C.sizeof_struct_v4l2_vp8_entropy),
		"cV4L2VP8EntropyCoderState": uintptr(C.sizeof_struct_v4l2_vp8_entropy_coder_state),
		"cV4L2CtrlVP8Frame":          uintptr(C.sizeof_struct_v4l2_ctrl_vp8_frame),
		"cV4L2CtrlMPEG2Sequence":     uintptr(C.sizeof_struct_v4l2_ctrl_mpeg2_sequence),
		"cV4L2CtrlMPEG2Picture":      uintptr(C.sizeof_struct_v4l2_ctrl_mpeg2_picture),
		"cV4L2CtrlMPEG2Quantisation": uintptr(C.sizeof_struct_v4l2_ctrl_mpeg2_quantisation),
	}
	out := make([]abiPair, 0, len(goSizes))
	for name, got := range goSizes {
		out = append(out, abiPair{name: name, got: got, want: cSizes[name]})
	}
	return out
}

func cStructOffsetChecks(goOffsets map[string]uintptr) []abiPair {
	cOffsets := map[string]uintptr{
		"cV4L2ExtControl.Value":       uintptr(C.offsetof_v4l2_ext_control_value()),
		"cV4L2Format.Fmt":             uintptr(C.offsetof_v4l2_format_fmt()),
		"cV4L2PixFormat.Ycbcr_enc":    uintptr(C.offsetof_v4l2_pix_format_ycbcr_enc()),
		"cV4L2Buffer.M":               uintptr(C.offsetof_v4l2_buffer_m()),
		"cV4L2Buffer.Fd":              uintptr(C.offsetof_v4l2_buffer_request_fd()),
		"cV4L2Event.U":                uintptr(C.offsetof_v4l2_event_u()),
		"cV4L2EventCtrl.Value":        uintptr(C.offsetof_v4l2_event_ctrl_value()),
		"cV4L2FrmSizeEnum.Discrete":   uintptr(C.offsetof_v4l2_frmsizeenum_discrete()),
		"cV4L2FrmIvalEnum.Discrete":   uintptr(C.offsetof_v4l2_frmivalenum_discrete()),
		"cV4L2StreamParm.Parm":        uintptr(C.offsetof_v4l2_streamparm_parm()),
		"cV4L2DVTimings.Pad_cgo_0":    uintptr(C.offsetof_v4l2_dv_timings_bt()),
		"cV4L2DVTimingsCap.Bt":        uintptr(C.offsetof_v4l2_dv_timings_cap_bt()),
		"cV4L2DecoderCmd.Stop":        uintptr(C.offsetof_v4l2_decoder_cmd_stop()),
		"cV4L2QueryMenu.Name":              uintptr(C.offsetof_v4l2_querymenu_name()),
		"cV4L2Plane.M":                     uintptr(C.offsetof_v4l2_plane_m()),
		"cV4L2Plane.Offset":                uintptr(C.offsetof_v4l2_plane_data_offset()),
		"cV4L2PixFormatMPlane.Num_planes":  uintptr(C.offsetof_v4l2_pix_format_mplane_num_planes()),
		"cV4L2CreateBuffers.Format":        uintptr(C.offsetof_v4l2_create_buffers_format()),
		"cV4L2Selection.R":                 uintptr(C.offsetof_v4l2_selection_r()),
		"cMediaV2Topology.Ptr_links":       uintptr(C.offsetof_media_v2_topology_ptr_links()),
		"cMediaV2Interface.Devnode":        uintptr(C.offsetof_media_v2_interface_devnode()),
		"cMediaEntityDesc.Dev":             uintptr(C.offsetof_media_entity_desc_dev()),
		"cV4L2SubdevFormat.Format":         uintptr(C.offsetof_v4l2_subdev_format_format()),
		"cV4L2SubdevSelection.R":           uintptr(C.offsetof_v4l2_subdev_selection_r()),
		"cV4L2MetaFormat.Buffersize":       uintptr(C.offsetof_v4l2_meta_format_buffersize()),
		"cV4L2EDID.Edid":                   uintptr(C.offsetof_v4l2_edid_edid()),
		"cV4L2SubdevRouting.Routes":        uintptr(C.offsetof_v4l2_subdev_routing_routes()),
		"cV4L2CtrlH264SPS.Flags":                       uintptr(C.offsetof_v4l2_ctrl_h264_sps_flags()),
		"cV4L2H264DPBEntry.Top_field_order_cnt":        uintptr(C.offsetof_v4l2_h264_dpb_entry_top_field_order_cnt()),
		"cV4L2H264DPBEntry.Flags":                      uintptr(C.offsetof_v4l2_h264_dpb_entry_flags()),
		"cV4L2CtrlH264DecodeParams.Flags":              uintptr(C.offsetof_v4l2_ctrl_h264_decode_params_flags()),
		"cV4L2CtrlFWHTParams.Quantization":             uintptr(C.offsetof_v4l2_ctrl_fwht_params_quantization()),
		"cV4L2VP8LoopFilter.Flags":                     uintptr(C.offsetof_v4l2_vp8_loop_filter_flags()),
		"cV4L2VP8Entropy.Mv_probs":                     uintptr(C.offsetof_v4l2_vp8_entropy_mv_probs()),
		"cV4L2CtrlVP8Frame.Dct_part_sizes":             uintptr(C.offsetof_v4l2_ctrl_vp8_frame_dct_part_sizes()),
		"cV4L2CtrlVP8Frame.Flags":                      uintptr(C.offsetof_v4l2_ctrl_vp8_frame_flags()),
		"cV4L2CtrlMPEG2Picture.Picture_coding_type":    uintptr(C.offsetof_v4l2_ctrl_mpeg2_picture_picture_coding_type()),
	}
	out := make([]abiPair, 0, len(goOffsets))
	for name, got := range goOffsets {
		out = append(out, abiPair{name: name, got: got, want: cOffsets[name]})
	}
	return out
}

func cIoctlChecks() []abiPair {
	return []abiPair{
		{"cMEDIA_IOC_DEVICE_INFO", cMEDIA_IOC_DEVICE_INFO, uintptr(C.MEDIA_IOC_DEVICE_INFO)},
		{"cMEDIA_IOC_ENUM_ENTITIES", cMEDIA_IOC_ENUM_ENTITIES, uintptr(C.MEDIA_IOC_ENUM_ENTITIES)},
		{"cMEDIA_IOC_ENUM_LINKS", cMEDIA_IOC_ENUM_LINKS, uintptr(C.MEDIA_IOC_ENUM_LINKS)},
		{"cMEDIA_IOC_G_TOPOLOGY", cMEDIA_IOC_G_TOPOLOGY, uintptr(C.MEDIA_IOC_G_TOPOLOGY)},
		{"cMEDIA_IOC_REQUEST_ALLOC", cMEDIA_IOC_REQUEST_ALLOC, uintptr(C.MEDIA_IOC_REQUEST_ALLOC)},
		{"cMEDIA_IOC_SETUP_LINK", cMEDIA_IOC_SETUP_LINK, uintptr(C.MEDIA_IOC_SETUP_LINK)},
		{"cMEDIA_REQUEST_IOC_QUEUE", cMEDIA_REQUEST_IOC_QUEUE, uintptr(C.MEDIA_REQUEST_IOC_QUEUE)},
		{"cMEDIA_REQUEST_IOC_REINIT", cMEDIA_REQUEST_IOC_REINIT, uintptr(C.MEDIA_REQUEST_IOC_REINIT)},
		{"cVIDIOC_SUBDEV_ENUM_FRAME_INTERVAL", cVIDIOC_SUBDEV_ENUM_FRAME_INTERVAL, uintptr(C.VIDIOC_SUBDEV_ENUM_FRAME_INTERVAL)},
		{"cVIDIOC_SUBDEV_ENUM_FRAME_SIZE", cVIDIOC_SUBDEV_ENUM_FRAME_SIZE, uintptr(C.VIDIOC_SUBDEV_ENUM_FRAME_SIZE)},
		{"cVIDIOC_SUBDEV_ENUM_MBUS_CODE", cVIDIOC_SUBDEV_ENUM_MBUS_CODE, uintptr(C.VIDIOC_SUBDEV_ENUM_MBUS_CODE)},
		{"cVIDIOC_SUBDEV_G_CROP", cVIDIOC_SUBDEV_G_CROP, uintptr(C.VIDIOC_SUBDEV_G_CROP)},
		{"cVIDIOC_SUBDEV_G_FMT", cVIDIOC_SUBDEV_G_FMT, uintptr(C.VIDIOC_SUBDEV_G_FMT)},
		{"cVIDIOC_SUBDEV_G_FRAME_INTERVAL", cVIDIOC_SUBDEV_G_FRAME_INTERVAL, uintptr(C.VIDIOC_SUBDEV_G_FRAME_INTERVAL)},
		{"cVIDIOC_SUBDEV_G_SELECTION", cVIDIOC_SUBDEV_G_SELECTION, uintptr(C.VIDIOC_SUBDEV_G_SELECTION)},
		{"cVIDIOC_SUBDEV_S_CROP", cVIDIOC_SUBDEV_S_CROP, uintptr(C.VIDIOC_SUBDEV_S_CROP)},
		{"cVIDIOC_SUBDEV_S_FMT", cVIDIOC_SUBDEV_S_FMT, uintptr(C.VIDIOC_SUBDEV_S_FMT)},
		{"cVIDIOC_SUBDEV_S_FRAME_INTERVAL", cVIDIOC_SUBDEV_S_FRAME_INTERVAL, uintptr(C.VIDIOC_SUBDEV_S_FRAME_INTERVAL)},
		{"cVIDIOC_SUBDEV_S_SELECTION", cVIDIOC_SUBDEV_S_SELECTION, uintptr(C.VIDIOC_SUBDEV_S_SELECTION)},
		{"cVIDIOC_SUBDEV_QUERYCAP", cVIDIOC_SUBDEV_QUERYCAP, uintptr(C.VIDIOC_SUBDEV_QUERYCAP)},
		{"cVIDIOC_SUBDEV_G_CLIENT_CAP", cVIDIOC_SUBDEV_G_CLIENT_CAP, uintptr(C.VIDIOC_SUBDEV_G_CLIENT_CAP)},
		{"cVIDIOC_SUBDEV_S_CLIENT_CAP", cVIDIOC_SUBDEV_S_CLIENT_CAP, uintptr(C.VIDIOC_SUBDEV_S_CLIENT_CAP)},
		{"cVIDIOC_SUBDEV_G_ROUTING", cVIDIOC_SUBDEV_G_ROUTING, uintptr(C.VIDIOC_SUBDEV_G_ROUTING)},
		{"cVIDIOC_SUBDEV_S_ROUTING", cVIDIOC_SUBDEV_S_ROUTING, uintptr(C.VIDIOC_SUBDEV_S_ROUTING)},
		{"cVIDIOC_SUBDEV_G_EDID", cVIDIOC_SUBDEV_G_EDID, uintptr(C.VIDIOC_SUBDEV_G_EDID)},
		{"cVIDIOC_SUBDEV_S_EDID", cVIDIOC_SUBDEV_S_EDID, uintptr(C.VIDIOC_SUBDEV_S_EDID)},
		{"cVIDIOC_G_EDID", cVIDIOC_G_EDID, uintptr(C.VIDIOC_G_EDID)},
		{"cVIDIOC_S_EDID", cVIDIOC_S_EDID, uintptr(C.VIDIOC_S_EDID)},
		{"cVIDIOC_CREATE_BUFS", cVIDIOC_CREATE_BUFS, uintptr(C.VIDIOC_CREATE_BUFS)},
		{"cVIDIOC_CROPCAP", cVIDIOC_CROPCAP, uintptr(C.VIDIOC_CROPCAP)},
		{"cVIDIOC_DECODER_CMD", cVIDIOC_DECODER_CMD, uintptr(C.VIDIOC_DECODER_CMD)},
		{"cVIDIOC_DQBUF", cVIDIOC_DQBUF, uintptr(C.VIDIOC_DQBUF)},
		{"cVIDIOC_DQEVENT", cVIDIOC_DQEVENT, uintptr(C.VIDIOC_DQEVENT)},
		{"cVIDIOC_DV_TIMINGS_CAP", cVIDIOC_DV_TIMINGS_CAP, uintptr(C.VIDIOC_DV_TIMINGS_CAP)},
		{"cVIDIOC_ENCODER_CMD", cVIDIOC_ENCODER_CMD, uintptr(C.VIDIOC_ENCODER_CMD)},
		{"cVIDIOC_ENUMAUDIO", cVIDIOC_ENUMAUDIO, uintptr(C.VIDIOC_ENUMAUDIO)},
		{"cVIDIOC_ENUMAUDOUT", cVIDIOC_ENUMAUDOUT, uintptr(C.VIDIOC_ENUMAUDOUT)},
		{"cVIDIOC_ENUMINPUT", cVIDIOC_ENUMINPUT, uintptr(C.VIDIOC_ENUMINPUT)},
		{"cVIDIOC_ENUMOUTPUT", cVIDIOC_ENUMOUTPUT, uintptr(C.VIDIOC_ENUMOUTPUT)},
		{"cVIDIOC_ENUMSTD", cVIDIOC_ENUMSTD, uintptr(C.VIDIOC_ENUMSTD)},
		{"cVIDIOC_ENUM_DV_TIMINGS", cVIDIOC_ENUM_DV_TIMINGS, uintptr(C.VIDIOC_ENUM_DV_TIMINGS)},
		{"cVIDIOC_ENUM_FMT", cVIDIOC_ENUM_FMT, uintptr(C.VIDIOC_ENUM_FMT)},
		{"cVIDIOC_ENUM_FRAMEINTERVALS", cVIDIOC_ENUM_FRAMEINTERVALS, uintptr(C.VIDIOC_ENUM_FRAMEINTERVALS)},
		{"cVIDIOC_ENUM_FRAMESIZES", cVIDIOC_ENUM_FRAMESIZES, uintptr(C.VIDIOC_ENUM_FRAMESIZES)},
		{"cVIDIOC_ENUM_FREQ_BANDS", cVIDIOC_ENUM_FREQ_BANDS, uintptr(C.VIDIOC_ENUM_FREQ_BANDS)},
		{"cVIDIOC_EXPBUF", cVIDIOC_EXPBUF, uintptr(C.VIDIOC_EXPBUF)},
		{"cVIDIOC_G_AUDIO", cVIDIOC_G_AUDIO, uintptr(C.VIDIOC_G_AUDIO)},
		{"cVIDIOC_G_AUDOUT", cVIDIOC_G_AUDOUT, uintptr(C.VIDIOC_G_AUDOUT)},
		{"cVIDIOC_G_CROP", cVIDIOC_G_CROP, uintptr(C.VIDIOC_G_CROP)},
		{"cVIDIOC_G_CTRL", cVIDIOC_G_CTRL, uintptr(C.VIDIOC_G_CTRL)},
		{"cVIDIOC_G_DV_TIMINGS", cVIDIOC_G_DV_TIMINGS, uintptr(C.VIDIOC_G_DV_TIMINGS)},
		{"cVIDIOC_G_EXT_CTRLS", cVIDIOC_G_EXT_CTRLS, uintptr(C.VIDIOC_G_EXT_CTRLS)},
		{"cVIDIOC_G_FMT", cVIDIOC_G_FMT, uintptr(C.VIDIOC_G_FMT)},
		{"cVIDIOC_G_FREQUENCY", cVIDIOC_G_FREQUENCY, uintptr(C.VIDIOC_G_FREQUENCY)},
		{"cVIDIOC_G_INPUT", cVIDIOC_G_INPUT, uintptr(C.VIDIOC_G_INPUT)},
		{"cVIDIOC_G_MODULATOR", cVIDIOC_G_MODULATOR, uintptr(C.VIDIOC_G_MODULATOR)},
		{"cVIDIOC_G_OUTPUT", cVIDIOC_G_OUTPUT, uintptr(C.VIDIOC_G_OUTPUT)},
		{"cVIDIOC_G_PARM", cVIDIOC_G_PARM, uintptr(C.VIDIOC_G_PARM)},
		{"cVIDIOC_G_SELECTION", cVIDIOC_G_SELECTION, uintptr(C.VIDIOC_G_SELECTION)},
		{"cVIDIOC_G_STD", cVIDIOC_G_STD, uintptr(C.VIDIOC_G_STD)},
		{"cVIDIOC_G_TUNER", cVIDIOC_G_TUNER, uintptr(C.VIDIOC_G_TUNER)},
		{"cVIDIOC_PREPARE_BUF", cVIDIOC_PREPARE_BUF, uintptr(C.VIDIOC_PREPARE_BUF)},
		{"cVIDIOC_QBUF", cVIDIOC_QBUF, uintptr(C.VIDIOC_QBUF)},
		{"cVIDIOC_QUERYBUF", cVIDIOC_QUERYBUF, uintptr(C.VIDIOC_QUERYBUF)},
		{"cVIDIOC_QUERYCAP", cVIDIOC_QUERYCAP, uintptr(C.VIDIOC_QUERYCAP)},
		{"cVIDIOC_QUERYCTRL", cVIDIOC_QUERYCTRL, uintptr(C.VIDIOC_QUERYCTRL)},
		{"cVIDIOC_QUERYMENU", cVIDIOC_QUERYMENU, uintptr(C.VIDIOC_QUERYMENU)},
		{"cVIDIOC_QUERYSTD", cVIDIOC_QUERYSTD, uintptr(C.VIDIOC_QUERYSTD)},
		{"cVIDIOC_QUERY_DV_TIMINGS", cVIDIOC_QUERY_DV_TIMINGS, uintptr(C.VIDIOC_QUERY_DV_TIMINGS)},
		{"cVIDIOC_QUERY_EXT_CTRL", cVIDIOC_QUERY_EXT_CTRL, uintptr(C.VIDIOC_QUERY_EXT_CTRL)},
		{"cVIDIOC_REMOVE_BUFS", cVIDIOC_REMOVE_BUFS, uintptr(C.VIDIOC_REMOVE_BUFS)},
		{"cVIDIOC_REQBUFS", cVIDIOC_REQBUFS, uintptr(C.VIDIOC_REQBUFS)},
		{"cVIDIOC_STREAMOFF", cVIDIOC_STREAMOFF, uintptr(C.VIDIOC_STREAMOFF)},
		{"cVIDIOC_STREAMON", cVIDIOC_STREAMON, uintptr(C.VIDIOC_STREAMON)},
		{"cVIDIOC_SUBSCRIBE_EVENT", cVIDIOC_SUBSCRIBE_EVENT, uintptr(C.VIDIOC_SUBSCRIBE_EVENT)},
		{"cVIDIOC_S_AUDIO", cVIDIOC_S_AUDIO, uintptr(C.VIDIOC_S_AUDIO)},
		{"cVIDIOC_S_AUDOUT", cVIDIOC_S_AUDOUT, uintptr(C.VIDIOC_S_AUDOUT)},
		{"cVIDIOC_S_CROP", cVIDIOC_S_CROP, uintptr(C.VIDIOC_S_CROP)},
		{"cVIDIOC_S_CTRL", cVIDIOC_S_CTRL, uintptr(C.VIDIOC_S_CTRL)},
		{"cVIDIOC_S_DV_TIMINGS", cVIDIOC_S_DV_TIMINGS, uintptr(C.VIDIOC_S_DV_TIMINGS)},
		{"cVIDIOC_S_EXT_CTRLS", cVIDIOC_S_EXT_CTRLS, uintptr(C.VIDIOC_S_EXT_CTRLS)},
		{"cVIDIOC_S_FMT", cVIDIOC_S_FMT, uintptr(C.VIDIOC_S_FMT)},
		{"cVIDIOC_S_FREQUENCY", cVIDIOC_S_FREQUENCY, uintptr(C.VIDIOC_S_FREQUENCY)},
		{"cVIDIOC_S_INPUT", cVIDIOC_S_INPUT, uintptr(C.VIDIOC_S_INPUT)},
		{"cVIDIOC_S_MODULATOR", cVIDIOC_S_MODULATOR, uintptr(C.VIDIOC_S_MODULATOR)},
		{"cVIDIOC_S_OUTPUT", cVIDIOC_S_OUTPUT, uintptr(C.VIDIOC_S_OUTPUT)},
		{"cVIDIOC_S_PARM", cVIDIOC_S_PARM, uintptr(C.VIDIOC_S_PARM)},
		{"cVIDIOC_S_SELECTION", cVIDIOC_S_SELECTION, uintptr(C.VIDIOC_S_SELECTION)},
		{"cVIDIOC_S_STD", cVIDIOC_S_STD, uintptr(C.VIDIOC_S_STD)},
		{"cVIDIOC_S_TUNER", cVIDIOC_S_TUNER, uintptr(C.VIDIOC_S_TUNER)},
		{"cVIDIOC_TRY_DECODER_CMD", cVIDIOC_TRY_DECODER_CMD, uintptr(C.VIDIOC_TRY_DECODER_CMD)},
		{"cVIDIOC_TRY_ENCODER_CMD", cVIDIOC_TRY_ENCODER_CMD, uintptr(C.VIDIOC_TRY_ENCODER_CMD)},
		{"cVIDIOC_TRY_EXT_CTRLS", cVIDIOC_TRY_EXT_CTRLS, uintptr(C.VIDIOC_TRY_EXT_CTRLS)},
		{"cVIDIOC_TRY_FMT", cVIDIOC_TRY_FMT, uintptr(C.VIDIOC_TRY_FMT)},
		{"cVIDIOC_UNSUBSCRIBE_EVENT", cVIDIOC_UNSUBSCRIBE_EVENT, uintptr(C.VIDIOC_UNSUBSCRIBE_EVENT)},
	}
}

func cValueConstantChecks() ([]abiPair, uint64) {
	return []abiPair{
		{"V4L2_PIX_FMT_YUYV", uintptr(PixelFmtYUYV), uintptr(uint32(C.V4L2_PIX_FMT_YUYV))},
		{"V4L2_PIX_FMT_MJPEG", uintptr(PixelFmtMJPEG), uintptr(uint32(C.V4L2_PIX_FMT_MJPEG))},
		{"V4L2_PIX_FMT_H264", uintptr(PixelFmtH264), uintptr(uint32(C.V4L2_PIX_FMT_H264))},
		{"V4L2_PIX_FMT_NV12", uintptr(PixelFmtNV12), uintptr(uint32(C.V4L2_PIX_FMT_NV12))},
		{"V4L2_BUF_TYPE_VIDEO_CAPTURE", uintptr(BufTypeVideoCapture), uintptr(uint32(C.V4L2_BUF_TYPE_VIDEO_CAPTURE))},
		{"V4L2_BUF_TYPE_VIDEO_OUTPUT", uintptr(BufTypeVideoOutput), uintptr(uint32(C.V4L2_BUF_TYPE_VIDEO_OUTPUT))},
		{"V4L2_BUF_TYPE_VIDEO_CAPTURE_MPLANE", uintptr(BufTypeVideoCaptureMPlane), uintptr(uint32(C.V4L2_BUF_TYPE_VIDEO_CAPTURE_MPLANE))},
		{"V4L2_BUF_TYPE_VIDEO_OUTPUT_MPLANE", uintptr(BufTypeVideoOutputMPlane), uintptr(uint32(C.V4L2_BUF_TYPE_VIDEO_OUTPUT_MPLANE))},
		{"V4L2_SEL_TGT_CROP", uintptr(SelTargetCrop), uintptr(uint32(C.V4L2_SEL_TGT_CROP))},
		{"V4L2_SEL_TGT_COMPOSE", uintptr(SelTargetCompose), uintptr(uint32(C.V4L2_SEL_TGT_COMPOSE))},
		{"VIDEO_MAX_PLANES", uintptr(MaxPlanes), uintptr(uint32(C.VIDEO_MAX_PLANES))},
		{"V4L2_MEMORY_MMAP", uintptr(IOTypeMMAP), uintptr(uint32(C.V4L2_MEMORY_MMAP))},
		{"V4L2_MEMORY_USERPTR", uintptr(IOTypeUserPtr), uintptr(uint32(C.V4L2_MEMORY_USERPTR))},
		{"V4L2_MEMORY_DMABUF", uintptr(IOTypeDMABuf), uintptr(uint32(C.V4L2_MEMORY_DMABUF))},
		{"V4L2_CTRL_FLAG_NEXT_CTRL", uintptr(CtrlFlagNextCtrl), uintptr(uint32(C.V4L2_CTRL_FLAG_NEXT_CTRL))},
		{"V4L2_CID_BRIGHTNESS", uintptr(CtrlBrightness), uintptr(uint32(C.V4L2_CID_BRIGHTNESS))},
		{"V4L2_CID_MPEG_VIDEO_ENCODING", uintptr(CtrlMPEGVideoEncoding), uintptr(uint32(C.V4L2_CID_MPEG_VIDEO_ENCODING))},
		{"V4L2_CAP_VIDEO_CAPTURE", uintptr(CapVideoCapture), uintptr(uint32(C.V4L2_CAP_VIDEO_CAPTURE))},
		{"V4L2_CAP_STREAMING", uintptr(CapStreaming), uintptr(uint32(C.V4L2_CAP_STREAMING))},
		{"V4L2_EVENT_CTRL", uintptr(EventCtrl), uintptr(uint32(C.V4L2_EVENT_CTRL))},
		{"V4L2_ENC_CMD_START", uintptr(EncCmdStart), uintptr(uint32(C.V4L2_ENC_CMD_START))},
		{"V4L2_DEC_CMD_STOP", uintptr(DecCmdStop), uintptr(uint32(C.V4L2_DEC_CMD_STOP))},
	}, uint64(C.V4L2_STD_ALL)
}
