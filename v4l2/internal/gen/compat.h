/*
 * compat.h - Preprocessor fallbacks for V4L2 constants added after Linux 5.10.
 *
 * go4vl targets Linux 5.10+. Constants introduced in later kernel versions
 * are defined here so the library compiles against older kernel headers.
 * Each definition is guarded with #ifndef so it is only used when the
 * system headers do not already provide it.
 */

#ifndef GO4VL_COMPAT_H
#define GO4VL_COMPAT_H

#include <linux/videodev2.h>
#include <linux/v4l2-subdev.h>

/* ---- Control flags (control.go) ---- */

/* Linux 5.11 */
#ifndef V4L2_CTRL_FLAG_MODIFY_LAYOUT
#define V4L2_CTRL_FLAG_MODIFY_LAYOUT 0x0400
#endif

/* Linux 5.13 */
#ifndef V4L2_CTRL_FLAG_EXECUTE_ON_WRITE
#define V4L2_CTRL_FLAG_EXECUTE_ON_WRITE 0x0200
#endif

/* Linux 5.18 */
#ifndef V4L2_CTRL_FLAG_DYNAMIC_ARRAY
#define V4L2_CTRL_FLAG_DYNAMIC_ARRAY 0x0800
#endif

/* ---- Control classes (control_values.go) ---- */

/* V4L2_CTRL_CLASS_CODEC and V4L2_CID_CODEC_CLASS were renamed from MPEG variants in Linux 5.5 */
#ifndef V4L2_CTRL_CLASS_CODEC
#define V4L2_CTRL_CLASS_CODEC 0x00990000
#endif
#ifndef V4L2_CID_CODEC_CLASS
#define V4L2_CID_CODEC_CLASS (V4L2_CTRL_CLASS_CODEC | 1)
#endif

/* Linux 5.8 */
#ifndef V4L2_CTRL_CLASS_CODEC_STATELESS
#define V4L2_CTRL_CLASS_CODEC_STATELESS 0x00a30000
#endif

/* ---- Control types (control_values.go) ---- */

/* Linux 5.5 - stateless codec types */
#ifndef V4L2_CTRL_TYPE_H264_SPS
#define V4L2_CTRL_TYPE_H264_SPS 0x0200
#endif
#ifndef V4L2_CTRL_TYPE_H264_PPS
#define V4L2_CTRL_TYPE_H264_PPS 0x0201
#endif
#ifndef V4L2_CTRL_TYPE_H264_SCALING_MATRIX
#define V4L2_CTRL_TYPE_H264_SCALING_MATRIX 0x0202
#endif
#ifndef V4L2_CTRL_TYPE_H264_SLICE_PARAMS
#define V4L2_CTRL_TYPE_H264_SLICE_PARAMS 0x0203
#endif
#ifndef V4L2_CTRL_TYPE_H264_DECODE_PARAMS
#define V4L2_CTRL_TYPE_H264_DECODE_PARAMS 0x0204
#endif
#ifndef V4L2_CTRL_TYPE_FWHT_PARAMS
#define V4L2_CTRL_TYPE_FWHT_PARAMS 0x0220
#endif
#ifndef V4L2_CTRL_TYPE_VP8_FRAME
#define V4L2_CTRL_TYPE_VP8_FRAME 0x0240
#endif

/* Linux 5.7 - MPEG2 stateless codec types */
#ifndef V4L2_CTRL_TYPE_MPEG2_QUANTISATION
#define V4L2_CTRL_TYPE_MPEG2_QUANTISATION 0x0250
#endif
#ifndef V4L2_CTRL_TYPE_MPEG2_SEQUENCE
#define V4L2_CTRL_TYPE_MPEG2_SEQUENCE 0x0251
#endif
#ifndef V4L2_CTRL_TYPE_MPEG2_PICTURE
#define V4L2_CTRL_TYPE_MPEG2_PICTURE 0x0252
#endif

/* Linux 5.11 */
#ifndef V4L2_CTRL_TYPE_HDR10_CLL_INFO
#define V4L2_CTRL_TYPE_HDR10_CLL_INFO 0x0228
#endif
#ifndef V4L2_CTRL_TYPE_HDR10_MASTERING_DISPLAY
#define V4L2_CTRL_TYPE_HDR10_MASTERING_DISPLAY 0x0229
#endif

/* Linux 5.13 */
#ifndef V4L2_CTRL_TYPE_VP9_COMPRESSED_HDR
#define V4L2_CTRL_TYPE_VP9_COMPRESSED_HDR 0x0260
#endif
#ifndef V4L2_CTRL_TYPE_VP9_FRAME
#define V4L2_CTRL_TYPE_VP9_FRAME 0x0261
#endif

/* ---- Control classes and IDs (control_values.go) ---- */

/* Linux 6.1 */
#ifndef V4L2_CTRL_CLASS_COLORIMETRY
#define V4L2_CTRL_CLASS_COLORIMETRY 0x00a40000
#endif
#ifndef V4L2_CID_CAMERA_ORIENTATION
#define V4L2_CID_CAMERA_ORIENTATION (V4L2_CID_CAMERA_CLASS_BASE+34)
#endif
#ifndef V4L2_CID_CAMERA_SENSOR_ROTATION
#define V4L2_CID_CAMERA_SENSOR_ROTATION (V4L2_CID_CAMERA_CLASS_BASE+35)
#endif
#ifndef V4L2_CID_IMAGE_SOURCE_CLASS
#define V4L2_CID_IMAGE_SOURCE_CLASS 0x009e0000
#endif
#ifndef V4L2_CID_UNIT_CELL_SIZE
#define V4L2_CID_UNIT_CELL_SIZE (V4L2_CID_IMAGE_SOURCE_CLASS+1)
#endif

/* V4L2_CID_COLORFX_RGB - added in Linux 5.2 */
#ifndef V4L2_CID_COLORFX_RGB
#define V4L2_CID_COLORFX_RGB (V4L2_CID_USER_BASE + 0x1032)
#endif

/* V4L2_COLORFX_SET_RGB - added in Linux 5.2 */
#ifndef V4L2_COLORFX_SET_RGB
#define V4L2_COLORFX_SET_RGB 16
#endif

/* V4L2_CID_NOTIFY_GAINS - added in Linux 5.9 */
#ifndef V4L2_CID_NOTIFY_GAINS
#define V4L2_CID_NOTIFY_GAINS (V4L2_CID_IMAGE_SOURCE_CLASS + 2)
#endif

/* Linux 6.0 */
#ifndef V4L2_CID_DEINTERLACING_MODE
#define V4L2_CID_DEINTERLACING_MODE (V4L2_CID_IMAGE_PROC_CLASS_BASE+4)
#endif
#ifndef V4L2_CID_DIGITAL_GAIN
#define V4L2_CID_DIGITAL_GAIN (V4L2_CID_IMAGE_PROC_CLASS_BASE+5)
#endif

/* ---- Stateless codec CIDs (ext_ctrls_*.go) ---- */

/* Codec stateless base (needed by CIDs below) */
#ifndef V4L2_CID_CODEC_STATELESS_BASE
#define V4L2_CID_CODEC_STATELESS_BASE (V4L2_CTRL_CLASS_CODEC_STATELESS | 0x900)
#endif

/* Linux 5.5 - H264 stateless */
#ifndef V4L2_CID_STATELESS_H264_SPS
#define V4L2_CID_STATELESS_H264_SPS (V4L2_CID_CODEC_STATELESS_BASE + 0)
#endif
#ifndef V4L2_CID_STATELESS_H264_PPS
#define V4L2_CID_STATELESS_H264_PPS (V4L2_CID_CODEC_STATELESS_BASE + 1)
#endif
#ifndef V4L2_CID_STATELESS_H264_SCALING_MATRIX
#define V4L2_CID_STATELESS_H264_SCALING_MATRIX (V4L2_CID_CODEC_STATELESS_BASE + 2)
#endif
#ifndef V4L2_CID_STATELESS_H264_SLICE_PARAMS
#define V4L2_CID_STATELESS_H264_SLICE_PARAMS (V4L2_CID_CODEC_STATELESS_BASE + 3)
#endif
#ifndef V4L2_CID_STATELESS_H264_DECODE_PARAMS
#define V4L2_CID_STATELESS_H264_DECODE_PARAMS (V4L2_CID_CODEC_STATELESS_BASE + 4)
#endif
#ifndef V4L2_CID_STATELESS_H264_PRED_WEIGHTS
#define V4L2_CID_STATELESS_H264_PRED_WEIGHTS (V4L2_CID_CODEC_STATELESS_BASE + 5)
#endif
#ifndef V4L2_CID_STATELESS_H264_DECODE_MODE
#define V4L2_CID_STATELESS_H264_DECODE_MODE (V4L2_CID_CODEC_STATELESS_BASE + 6)
#endif
#ifndef V4L2_CID_STATELESS_H264_START_CODE
#define V4L2_CID_STATELESS_H264_START_CODE (V4L2_CID_CODEC_STATELESS_BASE + 7)
#endif

/* Linux 5.5 - FWHT stateless */
#ifndef V4L2_CID_STATELESS_FWHT_PARAMS
#define V4L2_CID_STATELESS_FWHT_PARAMS (V4L2_CID_CODEC_STATELESS_BASE + 100)
#endif

/* Linux 5.5 - VP8 stateless */
#ifndef V4L2_CID_STATELESS_VP8_FRAME
#define V4L2_CID_STATELESS_VP8_FRAME (V4L2_CID_CODEC_STATELESS_BASE + 200)
#endif

/* Linux 5.7 - MPEG2 stateless */
#ifndef V4L2_CID_STATELESS_MPEG2_SEQUENCE
#define V4L2_CID_STATELESS_MPEG2_SEQUENCE (V4L2_CID_CODEC_STATELESS_BASE + 220)
#endif
#ifndef V4L2_CID_STATELESS_MPEG2_PICTURE
#define V4L2_CID_STATELESS_MPEG2_PICTURE (V4L2_CID_CODEC_STATELESS_BASE + 221)
#endif
#ifndef V4L2_CID_STATELESS_MPEG2_QUANTISATION
#define V4L2_CID_STATELESS_MPEG2_QUANTISATION (V4L2_CID_CODEC_STATELESS_BASE + 222)
#endif

/* Linux 5.13 - VP9 stateless */
#ifndef V4L2_CID_STATELESS_VP9_FRAME
#define V4L2_CID_STATELESS_VP9_FRAME (V4L2_CID_CODEC_STATELESS_BASE + 300)
#endif

/* ---- Pixel formats (format.go) ---- */

/* Linux 5.0 */
#ifndef V4L2_PIX_FMT_P010
#define V4L2_PIX_FMT_P010 v4l2_fourcc('P', '0', '1', '0')
#endif
#ifndef V4L2_PIX_FMT_P012
#define V4L2_PIX_FMT_P012 v4l2_fourcc('P', '0', '1', '2')
#endif

/* Linux 5.5 */
#ifndef V4L2_PIX_FMT_H264_SLICE
#define V4L2_PIX_FMT_H264_SLICE v4l2_fourcc('S', '2', '6', '4')
#endif
#ifndef V4L2_PIX_FMT_HEVC_SLICE
#define V4L2_PIX_FMT_HEVC_SLICE v4l2_fourcc('S', '2', '6', '5')
#endif
#ifndef V4L2_PIX_FMT_VP8_FRAME
#define V4L2_PIX_FMT_VP8_FRAME v4l2_fourcc('V', 'P', '8', 'F')
#endif

/* Linux 5.13 */
#ifndef V4L2_PIX_FMT_VP9_FRAME
#define V4L2_PIX_FMT_VP9_FRAME v4l2_fourcc('V', 'P', '9', 'F')
#endif

/* Linux 6.0 */
#ifndef V4L2_PIX_FMT_AV1_FRAME
#define V4L2_PIX_FMT_AV1_FRAME v4l2_fourcc('A', 'V', '1', 'F')
#endif

/* ---- Format flags (format.go, format_desc.go) ---- */

/* Linux 5.11 */
#ifndef V4L2_FMT_FLAG_CSC_COLORSPACE
#define V4L2_FMT_FLAG_CSC_COLORSPACE 0x0020
#endif
#ifndef V4L2_FMT_FLAG_CSC_XFER_FUNC
#define V4L2_FMT_FLAG_CSC_XFER_FUNC 0x0040
#endif
#ifndef V4L2_FMT_FLAG_CSC_YCBCR_ENC
#define V4L2_FMT_FLAG_CSC_YCBCR_ENC 0x0080
#endif
#ifndef V4L2_FMT_FLAG_CSC_HSV_ENC
#define V4L2_FMT_FLAG_CSC_HSV_ENC V4L2_FMT_FLAG_CSC_YCBCR_ENC
#endif
#ifndef V4L2_FMT_FLAG_CSC_QUANTIZATION
#define V4L2_FMT_FLAG_CSC_QUANTIZATION 0x0100
#endif

/* ---- Buffer flags (streaming.go) ---- */

/* Linux 5.14 */
#ifndef V4L2_BUF_FLAG_NO_CACHE_INVALIDATE
#define V4L2_BUF_FLAG_NO_CACHE_INVALIDATE 0x00000800
#endif
#ifndef V4L2_BUF_FLAG_NO_CACHE_CLEAN
#define V4L2_BUF_FLAG_NO_CACHE_CLEAN 0x00001000
#endif

/* Linux 5.15 */
#ifndef V4L2_BUF_FLAG_LAST
#define V4L2_BUF_FLAG_LAST 0x00100000
#endif

/* ---- Event flags (events.go) ---- */

/* Linux 5.11 */
#ifndef V4L2_EVENT_CTRL_CH_DIMENSIONS
#define V4L2_EVENT_CTRL_CH_DIMENSIONS 0x0004
#endif

/* Linux 6.10 */
#ifndef VIDIOC_REMOVE_BUFS
struct v4l2_remove_buffers {
	__u32 index;
	__u32 count;
	__u32 type;
	__u32 reserved[13];
};
#define VIDIOC_REMOVE_BUFS _IOWR('V', 104, struct v4l2_remove_buffers)
#endif

/* ---- Metadata formats (format_meta.go) ---- */

#ifndef V4L2_META_FMT_VSP1_HGO
#define V4L2_META_FMT_VSP1_HGO v4l2_fourcc('V', 'S', 'P', 'H')
#endif
#ifndef V4L2_META_FMT_VSP1_HGT
#define V4L2_META_FMT_VSP1_HGT v4l2_fourcc('V', 'S', 'P', 'T')
#endif
#ifndef V4L2_META_FMT_UVC
#define V4L2_META_FMT_UVC v4l2_fourcc('U', 'V', 'C', 'H')
#endif
#ifndef V4L2_META_FMT_D4XX
#define V4L2_META_FMT_D4XX v4l2_fourcc('D', '4', 'X', 'X')
#endif
#ifndef V4L2_META_FMT_VIVID
#define V4L2_META_FMT_VIVID v4l2_fourcc('V', 'I', 'V', 'D')
#endif
#ifndef V4L2_META_FMT_RK_ISP1_PARAMS
#define V4L2_META_FMT_RK_ISP1_PARAMS v4l2_fourcc('R', 'K', '1', 'P')
#endif
#ifndef V4L2_META_FMT_RK_ISP1_STAT_3A
#define V4L2_META_FMT_RK_ISP1_STAT_3A v4l2_fourcc('R', 'K', '1', 'S')
#endif
#ifndef V4L2_META_FMT_BCM2835_ISP_STATS
#define V4L2_META_FMT_BCM2835_ISP_STATS v4l2_fourcc('B', 'S', 'T', 'A')
#endif
#ifndef V4L2_META_FMT_RPI_BE_CFG
#define V4L2_META_FMT_RPI_BE_CFG v4l2_fourcc('R', 'P', 'B', 'C')
#endif
#ifndef V4L2_META_FMT_RPI_FE_CFG
#define V4L2_META_FMT_RPI_FE_CFG v4l2_fourcc('R', 'P', 'F', 'C')
#endif
#ifndef V4L2_META_FMT_RPI_FE_STATS
#define V4L2_META_FMT_RPI_FE_STATS v4l2_fourcc('R', 'P', 'F', 'S')
#endif

/* ---- Sub-device capabilities & routing (subdev.go) ---- */

#ifndef V4L2_SUBDEV_CAP_RO_SUBDEV
#define V4L2_SUBDEV_CAP_RO_SUBDEV 0x00000001
#endif
#ifndef V4L2_SUBDEV_CAP_STREAMS
#define V4L2_SUBDEV_CAP_STREAMS 0x00000002
#endif
#ifndef V4L2_SUBDEV_ROUTE_FL_ACTIVE
#define V4L2_SUBDEV_ROUTE_FL_ACTIVE (1U << 0)
#endif
#ifndef VIDIOC_SUBDEV_G_ROUTING
struct v4l2_subdev_route {
	__u32 sink_pad;
	__u32 sink_stream;
	__u32 source_pad;
	__u32 source_stream;
	__u32 flags;
	__u32 reserved[5];
};
struct v4l2_subdev_routing {
	__u32 which;
	__u32 len_routes;
	__u64 routes;
	__u32 num_routes;
	__u32 reserved[11];
};
#define VIDIOC_SUBDEV_G_ROUTING _IOWR('V', 38, struct v4l2_subdev_routing)
#define VIDIOC_SUBDEV_S_ROUTING _IOWR('V', 39, struct v4l2_subdev_routing)
#endif
#ifndef V4L2_SUBDEV_CLIENT_CAP_STREAMS
#define V4L2_SUBDEV_CLIENT_CAP_STREAMS (1ULL << 0)
#endif
#ifndef V4L2_SUBDEV_CLIENT_CAP_INTERVAL_USES_WHICH
#define V4L2_SUBDEV_CLIENT_CAP_INTERVAL_USES_WHICH (1ULL << 1)
#endif
#ifndef VIDIOC_SUBDEV_G_CLIENT_CAP
struct v4l2_subdev_client_capability {
	__u64 capabilities;
};
#define VIDIOC_SUBDEV_G_CLIENT_CAP _IOR('V', 101, struct v4l2_subdev_client_capability)
#define VIDIOC_SUBDEV_S_CLIENT_CAP _IOWR('V', 102, struct v4l2_subdev_client_capability)
#endif

#endif /* GO4VL_COMPAT_H */
