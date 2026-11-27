//go:build ignore

package v4l2

/*
#ifndef _FILE_OFFSET_BITS
#define _FILE_OFFSET_BITS 64
#endif
#ifndef _TIME_BITS
#define _TIME_BITS 64
#endif
#include <linux/videodev2.h>
#include <linux/v4l2-controls.h>
#include <linux/v4l2-common.h>
#include <linux/media.h>
#include <linux/v4l2-subdev.h>
#include <linux/media-bus-format.h>
*/
import "C"

const (
	cMEDIA_IOC_DEVICE_INFO uintptr = C.MEDIA_IOC_DEVICE_INFO
	cMEDIA_IOC_ENUM_ENTITIES uintptr = C.MEDIA_IOC_ENUM_ENTITIES
	cMEDIA_IOC_ENUM_LINKS uintptr = C.MEDIA_IOC_ENUM_LINKS
	cMEDIA_IOC_G_TOPOLOGY uintptr = C.MEDIA_IOC_G_TOPOLOGY
	cMEDIA_IOC_REQUEST_ALLOC uintptr = C.MEDIA_IOC_REQUEST_ALLOC
	cMEDIA_IOC_SETUP_LINK uintptr = C.MEDIA_IOC_SETUP_LINK
	cMEDIA_REQUEST_IOC_QUEUE uintptr = C.MEDIA_REQUEST_IOC_QUEUE
	cMEDIA_REQUEST_IOC_REINIT uintptr = C.MEDIA_REQUEST_IOC_REINIT
	cVIDIOC_CREATE_BUFS uintptr = C.VIDIOC_CREATE_BUFS
	cVIDIOC_CROPCAP uintptr = C.VIDIOC_CROPCAP
	cVIDIOC_DECODER_CMD uintptr = C.VIDIOC_DECODER_CMD
	cVIDIOC_DQBUF uintptr = C.VIDIOC_DQBUF
	cVIDIOC_DQEVENT uintptr = C.VIDIOC_DQEVENT
	cVIDIOC_DV_TIMINGS_CAP uintptr = C.VIDIOC_DV_TIMINGS_CAP
	cVIDIOC_ENCODER_CMD uintptr = C.VIDIOC_ENCODER_CMD
	cVIDIOC_ENUMAUDIO uintptr = C.VIDIOC_ENUMAUDIO
	cVIDIOC_ENUMAUDOUT uintptr = C.VIDIOC_ENUMAUDOUT
	cVIDIOC_ENUMINPUT uintptr = C.VIDIOC_ENUMINPUT
	cVIDIOC_ENUMOUTPUT uintptr = C.VIDIOC_ENUMOUTPUT
	cVIDIOC_ENUMSTD uintptr = C.VIDIOC_ENUMSTD
	cVIDIOC_ENUM_DV_TIMINGS uintptr = C.VIDIOC_ENUM_DV_TIMINGS
	cVIDIOC_ENUM_FMT uintptr = C.VIDIOC_ENUM_FMT
	cVIDIOC_ENUM_FRAMEINTERVALS uintptr = C.VIDIOC_ENUM_FRAMEINTERVALS
	cVIDIOC_ENUM_FRAMESIZES uintptr = C.VIDIOC_ENUM_FRAMESIZES
	cVIDIOC_ENUM_FREQ_BANDS uintptr = C.VIDIOC_ENUM_FREQ_BANDS
	cVIDIOC_EXPBUF uintptr = C.VIDIOC_EXPBUF
	cVIDIOC_G_AUDIO uintptr = C.VIDIOC_G_AUDIO
	cVIDIOC_G_AUDOUT uintptr = C.VIDIOC_G_AUDOUT
	cVIDIOC_G_CROP uintptr = C.VIDIOC_G_CROP
	cVIDIOC_G_CTRL uintptr = C.VIDIOC_G_CTRL
	cVIDIOC_G_DV_TIMINGS uintptr = C.VIDIOC_G_DV_TIMINGS
	cVIDIOC_G_EDID uintptr = C.VIDIOC_G_EDID
	cVIDIOC_G_EXT_CTRLS uintptr = C.VIDIOC_G_EXT_CTRLS
	cVIDIOC_G_FMT uintptr = C.VIDIOC_G_FMT
	cVIDIOC_G_FREQUENCY uintptr = C.VIDIOC_G_FREQUENCY
	cVIDIOC_G_INPUT uintptr = C.VIDIOC_G_INPUT
	cVIDIOC_G_MODULATOR uintptr = C.VIDIOC_G_MODULATOR
	cVIDIOC_G_OUTPUT uintptr = C.VIDIOC_G_OUTPUT
	cVIDIOC_G_PARM uintptr = C.VIDIOC_G_PARM
	cVIDIOC_G_SELECTION uintptr = C.VIDIOC_G_SELECTION
	cVIDIOC_G_STD uintptr = C.VIDIOC_G_STD
	cVIDIOC_G_TUNER uintptr = C.VIDIOC_G_TUNER
	cVIDIOC_PREPARE_BUF uintptr = C.VIDIOC_PREPARE_BUF
	cVIDIOC_QBUF uintptr = C.VIDIOC_QBUF
	cVIDIOC_QUERYBUF uintptr = C.VIDIOC_QUERYBUF
	cVIDIOC_QUERYCAP uintptr = C.VIDIOC_QUERYCAP
	cVIDIOC_QUERYCTRL uintptr = C.VIDIOC_QUERYCTRL
	cVIDIOC_QUERYMENU uintptr = C.VIDIOC_QUERYMENU
	cVIDIOC_QUERYSTD uintptr = C.VIDIOC_QUERYSTD
	cVIDIOC_QUERY_DV_TIMINGS uintptr = C.VIDIOC_QUERY_DV_TIMINGS
	cVIDIOC_QUERY_EXT_CTRL uintptr = C.VIDIOC_QUERY_EXT_CTRL
	cVIDIOC_REMOVE_BUFS uintptr = C.VIDIOC_REMOVE_BUFS
	cVIDIOC_REQBUFS uintptr = C.VIDIOC_REQBUFS
	cVIDIOC_STREAMOFF uintptr = C.VIDIOC_STREAMOFF
	cVIDIOC_STREAMON uintptr = C.VIDIOC_STREAMON
	cVIDIOC_SUBDEV_ENUM_FRAME_INTERVAL uintptr = C.VIDIOC_SUBDEV_ENUM_FRAME_INTERVAL
	cVIDIOC_SUBDEV_ENUM_FRAME_SIZE uintptr = C.VIDIOC_SUBDEV_ENUM_FRAME_SIZE
	cVIDIOC_SUBDEV_ENUM_MBUS_CODE uintptr = C.VIDIOC_SUBDEV_ENUM_MBUS_CODE
	cVIDIOC_SUBDEV_G_CLIENT_CAP uintptr = C.VIDIOC_SUBDEV_G_CLIENT_CAP
	cVIDIOC_SUBDEV_G_CROP uintptr = C.VIDIOC_SUBDEV_G_CROP
	cVIDIOC_SUBDEV_G_EDID uintptr = C.VIDIOC_SUBDEV_G_EDID
	cVIDIOC_SUBDEV_G_FMT uintptr = C.VIDIOC_SUBDEV_G_FMT
	cVIDIOC_SUBDEV_G_FRAME_INTERVAL uintptr = C.VIDIOC_SUBDEV_G_FRAME_INTERVAL
	cVIDIOC_SUBDEV_G_ROUTING uintptr = C.VIDIOC_SUBDEV_G_ROUTING
	cVIDIOC_SUBDEV_G_SELECTION uintptr = C.VIDIOC_SUBDEV_G_SELECTION
	cVIDIOC_SUBDEV_QUERYCAP uintptr = C.VIDIOC_SUBDEV_QUERYCAP
	cVIDIOC_SUBDEV_S_CLIENT_CAP uintptr = C.VIDIOC_SUBDEV_S_CLIENT_CAP
	cVIDIOC_SUBDEV_S_CROP uintptr = C.VIDIOC_SUBDEV_S_CROP
	cVIDIOC_SUBDEV_S_EDID uintptr = C.VIDIOC_SUBDEV_S_EDID
	cVIDIOC_SUBDEV_S_FMT uintptr = C.VIDIOC_SUBDEV_S_FMT
	cVIDIOC_SUBDEV_S_FRAME_INTERVAL uintptr = C.VIDIOC_SUBDEV_S_FRAME_INTERVAL
	cVIDIOC_SUBDEV_S_ROUTING uintptr = C.VIDIOC_SUBDEV_S_ROUTING
	cVIDIOC_SUBDEV_S_SELECTION uintptr = C.VIDIOC_SUBDEV_S_SELECTION
	cVIDIOC_SUBSCRIBE_EVENT uintptr = C.VIDIOC_SUBSCRIBE_EVENT
	cVIDIOC_S_AUDIO uintptr = C.VIDIOC_S_AUDIO
	cVIDIOC_S_AUDOUT uintptr = C.VIDIOC_S_AUDOUT
	cVIDIOC_S_CROP uintptr = C.VIDIOC_S_CROP
	cVIDIOC_S_CTRL uintptr = C.VIDIOC_S_CTRL
	cVIDIOC_S_DV_TIMINGS uintptr = C.VIDIOC_S_DV_TIMINGS
	cVIDIOC_S_EDID uintptr = C.VIDIOC_S_EDID
	cVIDIOC_S_EXT_CTRLS uintptr = C.VIDIOC_S_EXT_CTRLS
	cVIDIOC_S_FMT uintptr = C.VIDIOC_S_FMT
	cVIDIOC_S_FREQUENCY uintptr = C.VIDIOC_S_FREQUENCY
	cVIDIOC_S_INPUT uintptr = C.VIDIOC_S_INPUT
	cVIDIOC_S_MODULATOR uintptr = C.VIDIOC_S_MODULATOR
	cVIDIOC_S_OUTPUT uintptr = C.VIDIOC_S_OUTPUT
	cVIDIOC_S_PARM uintptr = C.VIDIOC_S_PARM
	cVIDIOC_S_SELECTION uintptr = C.VIDIOC_S_SELECTION
	cVIDIOC_S_STD uintptr = C.VIDIOC_S_STD
	cVIDIOC_S_TUNER uintptr = C.VIDIOC_S_TUNER
	cVIDIOC_TRY_DECODER_CMD uintptr = C.VIDIOC_TRY_DECODER_CMD
	cVIDIOC_TRY_ENCODER_CMD uintptr = C.VIDIOC_TRY_ENCODER_CMD
	cVIDIOC_TRY_EXT_CTRLS uintptr = C.VIDIOC_TRY_EXT_CTRLS
	cVIDIOC_TRY_FMT uintptr = C.VIDIOC_TRY_FMT
	cVIDIOC_UNSUBSCRIBE_EVENT uintptr = C.VIDIOC_UNSUBSCRIBE_EVENT
)

type cTimeval C.struct_timeval
type cTimespec C.struct_timespec
type cV4L2Timecode C.struct_v4l2_timecode
type cV4L2Fract C.struct_v4l2_fract
type cV4L2Rect C.struct_v4l2_rect
type cMediaDeviceInfo C.struct_media_device_info
type cMediaV2Topology C.struct_media_v2_topology
type cMediaV2Entity C.struct_media_v2_entity
type cMediaV2Interface C.struct_media_v2_interface
type cMediaV2IntfDevnode C.struct_media_v2_intf_devnode
type cMediaV2Pad C.struct_media_v2_pad
type cMediaV2Link C.struct_media_v2_link
type cMediaEntityDesc C.struct_media_entity_desc
type cMediaPadDesc C.struct_media_pad_desc
type cMediaLinkDesc C.struct_media_link_desc
type cMediaLinksEnum C.struct_media_links_enum
type cV4L2MBusFrameFmt C.struct_v4l2_mbus_framefmt
type cV4L2SubdevFormat C.struct_v4l2_subdev_format
type cV4L2SubdevCrop C.struct_v4l2_subdev_crop
type cV4L2SubdevMBusCodeEnum C.struct_v4l2_subdev_mbus_code_enum
type cV4L2SubdevFrameSizeEnum C.struct_v4l2_subdev_frame_size_enum
type cV4L2SubdevFrameInterval C.struct_v4l2_subdev_frame_interval
type cV4L2SubdevFrameIntervalEnum C.struct_v4l2_subdev_frame_interval_enum
type cV4L2SubdevSelection C.struct_v4l2_subdev_selection
type cV4L2SubdevCapability C.struct_v4l2_subdev_capability
type cV4L2SubdevClientCapability C.struct_v4l2_subdev_client_capability
type cV4L2SubdevRoute C.struct_v4l2_subdev_route
type cV4L2SubdevRouting C.struct_v4l2_subdev_routing
type cV4L2EDID C.struct_v4l2_edid
type cV4L2MetaFormat C.struct_v4l2_meta_format
type cV4L2Capability C.struct_v4l2_capability
type cV4L2Audio C.struct_v4l2_audio
type cV4L2AudioOut C.struct_v4l2_audioout
type cV4L2BTTimings C.struct_v4l2_bt_timings
type cV4L2BTTimingsCap C.struct_v4l2_bt_timings_cap
type cV4L2Buffer C.struct_v4l2_buffer
type cV4L2CaptureParm C.struct_v4l2_captureparm
type cV4L2Control C.struct_v4l2_control
type cV4L2CreateBuffers C.struct_v4l2_create_buffers
type cV4L2Crop C.struct_v4l2_crop
type cV4L2CropCap C.struct_v4l2_cropcap
type cV4L2DecoderCmd C.struct_v4l2_decoder_cmd
type cV4L2DVTimings C.struct_v4l2_dv_timings
type cV4L2DVTimingsCap C.struct_v4l2_dv_timings_cap
type cV4L2EncoderCmd C.struct_v4l2_encoder_cmd
type cV4L2EnumDVTimings C.struct_v4l2_enum_dv_timings
type cV4L2EventCtrl C.struct_v4l2_event_ctrl
type cV4L2EventFrameSync C.struct_v4l2_event_frame_sync
type cV4L2EventMotionDet C.struct_v4l2_event_motion_det
type cV4L2EventSrcChange C.struct_v4l2_event_src_change
type cV4L2EventVsync C.struct_v4l2_event_vsync
type cV4L2Event C.struct_v4l2_event
type cV4L2EventSubscription C.struct_v4l2_event_subscription
type cV4L2ExportBuffer C.struct_v4l2_exportbuffer
type cV4L2ExtControl C.struct_v4l2_ext_control
type cV4L2ExtControls C.struct_v4l2_ext_controls
type cV4L2FmtDesc C.struct_v4l2_fmtdesc
type cV4L2Format C.struct_v4l2_format
type cV4L2Frequency C.struct_v4l2_frequency
type cV4L2FrequencyBand C.struct_v4l2_frequency_band
type cV4L2FrmSizeDiscrete C.struct_v4l2_frmsize_discrete
type cV4L2FrmSizeStepwise C.struct_v4l2_frmsize_stepwise
type cV4L2FrmSizeEnum C.struct_v4l2_frmsizeenum
type cV4L2FrmIvalStepwise C.struct_v4l2_frmival_stepwise
type cV4L2FrmIvalEnum C.struct_v4l2_frmivalenum
type cV4L2Input C.struct_v4l2_input
type cV4L2Modulator C.struct_v4l2_modulator
type cV4L2Output C.struct_v4l2_output
type cV4L2OutputParm C.struct_v4l2_outputparm
type cV4L2PixFormat C.struct_v4l2_pix_format
type cV4L2PlanePixFormat C.struct_v4l2_plane_pix_format
type cV4L2PixFormatMPlane C.struct_v4l2_pix_format_mplane
type cV4L2Plane C.struct_v4l2_plane
type cV4L2QueryCtrl C.struct_v4l2_queryctrl
type cV4L2QueryExtCtrl C.struct_v4l2_query_ext_ctrl
type cV4L2QueryMenu C.struct_v4l2_querymenu
type cV4L2RequestBuffers C.struct_v4l2_requestbuffers
type cV4L2RemoveBuffers C.struct_v4l2_remove_buffers
type cV4L2Selection C.struct_v4l2_selection
type cV4L2Standard C.struct_v4l2_standard
type cV4L2StreamParm C.struct_v4l2_streamparm
type cV4L2Tuner C.struct_v4l2_tuner
type cV4L2HEVCPredWeightTable C.struct_v4l2_hevc_pred_weight_table
type cV4L2HEVCDPBEntry C.struct_v4l2_hevc_dpb_entry
type cV4L2CtrlHEVCSPS C.struct_v4l2_ctrl_hevc_sps
type cV4L2CtrlHEVCPPS C.struct_v4l2_ctrl_hevc_pps
type cV4L2CtrlHEVCSliceParams C.struct_v4l2_ctrl_hevc_slice_params
type cV4L2CtrlHEVCScalingMatrix C.struct_v4l2_ctrl_hevc_scaling_matrix
type cV4L2CtrlHEVCDecodeParams C.struct_v4l2_ctrl_hevc_decode_params
type cV4L2VP9LoopFilter C.struct_v4l2_vp9_loop_filter
type cV4L2VP9Quantization C.struct_v4l2_vp9_quantization
type cV4L2VP9Segmentation C.struct_v4l2_vp9_segmentation
type cV4L2VP9MVProbs C.struct_v4l2_vp9_mv_probs
type cV4L2CtrlVP9Frame C.struct_v4l2_ctrl_vp9_frame
type cV4L2CtrlVP9CompressedHdr C.struct_v4l2_ctrl_vp9_compressed_hdr
type cV4L2AV1TileInfo C.struct_v4l2_av1_tile_info
type cV4L2AV1Quantization C.struct_v4l2_av1_quantization
type cV4L2AV1Segmentation C.struct_v4l2_av1_segmentation
type cV4L2AV1LoopFilter C.struct_v4l2_av1_loop_filter
type cV4L2AV1CDEF C.struct_v4l2_av1_cdef
type cV4L2AV1LoopRestoration C.struct_v4l2_av1_loop_restoration
type cV4L2AV1GlobalMotion C.struct_v4l2_av1_global_motion
type cV4L2CtrlAV1Sequence C.struct_v4l2_ctrl_av1_sequence
type cV4L2CtrlAV1TileGroupEntry C.struct_v4l2_ctrl_av1_tile_group_entry
type cV4L2CtrlAV1Frame C.struct_v4l2_ctrl_av1_frame
type cV4L2CtrlAV1FilmGrain C.struct_v4l2_ctrl_av1_film_grain
type cV4L2CtrlH264SPS C.struct_v4l2_ctrl_h264_sps
type cV4L2CtrlH264PPS C.struct_v4l2_ctrl_h264_pps
type cV4L2CtrlH264ScalingMatrix C.struct_v4l2_ctrl_h264_scaling_matrix
type cV4L2H264WeightFactors C.struct_v4l2_h264_weight_factors
type cV4L2CtrlH264PredWeights C.struct_v4l2_ctrl_h264_pred_weights
type cV4L2H264Reference C.struct_v4l2_h264_reference
type cV4L2CtrlH264SliceParams C.struct_v4l2_ctrl_h264_slice_params
type cV4L2H264DPBEntry C.struct_v4l2_h264_dpb_entry
type cV4L2CtrlH264DecodeParams C.struct_v4l2_ctrl_h264_decode_params
type cV4L2CtrlFWHTParams C.struct_v4l2_ctrl_fwht_params
type cV4L2VP8Segment C.struct_v4l2_vp8_segment
type cV4L2VP8LoopFilter C.struct_v4l2_vp8_loop_filter
type cV4L2VP8Quantization C.struct_v4l2_vp8_quantization
type cV4L2VP8Entropy C.struct_v4l2_vp8_entropy
type cV4L2VP8EntropyCoderState C.struct_v4l2_vp8_entropy_coder_state
type cV4L2CtrlVP8Frame C.struct_v4l2_ctrl_vp8_frame
type cV4L2CtrlMPEG2Sequence C.struct_v4l2_ctrl_mpeg2_sequence
type cV4L2CtrlMPEG2Picture C.struct_v4l2_ctrl_mpeg2_picture
type cV4L2CtrlMPEG2Quantisation C.struct_v4l2_ctrl_mpeg2_quantisation
