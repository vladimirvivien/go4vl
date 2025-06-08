package v4l2

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../include/
#include <linux/videodev2.h>
#include <linux/v4l2-controls.h>
*/
import "C"

// MPEGStreamType is a type alias for uint32, representing the type of an MPEG stream.
// Used with the CtrlMPEGStreamType control ID.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L237
type MPEGStreamType = uint32

// MPEG Stream Type Enum Values
const (
	MPEGStreamTypeMPEG2ProgramStream   MPEGStreamType = C.V4L2_MPEG_STREAM_TYPE_MPEG2_PS  // MPEG-2 Program Stream.
	MPEGStreamTypeMPEG2TransportStream MPEGStreamType = C.V4L2_MPEG_STREAM_TYPE_MPEG2_TS  // MPEG-2 Transport Stream.
	MPEGStreamTypeMPEG1SystemStream    MPEGStreamType = C.V4L2_MPEG_STREAM_TYPE_MPEG1_SS  // MPEG-1 System Stream.
	MPEGStreamTypeMPEG2DVD             MPEGStreamType = C.V4L2_MPEG_STREAM_TYPE_MPEG2_DVD // MPEG-2 DVD-compatible stream.
	MPEGStreamTypeMPEG1VCD             MPEGStreamType = C.V4L2_MPEG_STREAM_TYPE_MPEG1_VCD // MPEG-1 VCD-compatible stream.
	MPEGStreamTypeMPEG2SVCD            MPEGStreamType = C.V4L2_MPEG_STREAM_TYPE_MPEG2_SVCD // MPEG-2 SVCD-compatible stream.
)

// MPEGVideoEncoding is a type alias for uint32, representing the video encoding format for MPEG streams.
// Used with the CtrlMPEGVideoEncoding control ID.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L251
type MPEGVideoEncoding = uint32

// MPEG Video Encoding Enum Values
const (
	MPEGVideoEncodingMPEG1    MPEGVideoEncoding = C.V4L2_MPEG_VIDEO_ENCODING_MPEG_1      // MPEG-1 video encoding.
	MPEGVideoEncodingMPEG2    MPEGVideoEncoding = C.V4L2_MPEG_VIDEO_ENCODING_MPEG_2      // MPEG-2 video encoding.
	MPEGVideoEncodingMPEG4AVC MPEGVideoEncoding = C.V4L2_MPEG_VIDEO_ENCODING_MPEG_4_AVC // MPEG-4 AVC (H.264) video encoding.
)

// MPEGVideoAspect is a type alias for uint32, representing the aspect ratio of MPEG video.
// Used with the CtrlMPEGVideoAspect control ID.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L260
type MPEGVideoAspect = uint32

// MPEG Video Aspect Ratio Enum Values
const (
	MPEGVideoAspect1x1     MPEGVideoAspect = C.V4L2_MPEG_VIDEO_ASPECT_1x1      // 1:1 aspect ratio.
	MPEGVideoAspect4x3     MPEGVideoAspect = C.V4L2_MPEG_VIDEO_ASPECT_4x3      // 4:3 aspect ratio.
	MPEGVideoAspect16x9    MPEGVideoAspect = C.V4L2_MPEG_VIDEO_ASPECT_16x9     // 16:9 aspect ratio.
	MPEGVideoAspect221x100 MPEGVideoAspect = C.V4L2_MPEG_VIDEO_ASPECT_221x100 // 2.21:1 aspect ratio.
)

// MPEGVideoBitrateMode is a type alias for uint32, representing the bitrate control mode for MPEG video.
// Used with the CtrlMPEGVideoBitrateMode control ID.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L307
type MPEGVideoBitrateMode = uint32

// MPEG Video Bitrate Mode Enum Values
const (
	MPEGVideoBitrateModeVBR MPEGVideoBitrateMode = C.V4L2_MPEG_VIDEO_BITRATE_MODE_VBR // Variable Bitrate Mode.
	MPEGVideoBitrateModeCBR MPEGVideoBitrateMode = C.V4L2_MPEG_VIDEO_BITRATE_MODE_CBR // Constant Bitrate Mode.
	MPEGVideoBitrateModeCQ  MPEGVideoBitrateMode = C.V4L2_MPEG_VIDEO_BITRATE_MODE_CQ  // Constant Quality Mode.
)

// Codec Control IDs (CtrlID constants specific to codecs, primarily MPEG).
// These typically belong to the CtrlClassCodec class.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec.html
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L228
const (
	// CtrlCodecClass identifies the codec control class. This is a class identifier, not a control itself.
	CtrlCodecClass CtrlID = C.V4L2_CID_CODEC_CLASS

	// CtrlMPEGStreamType sets the type of the MPEG stream (e.g., Program Stream, Transport Stream).
	// Note: In the C header, V4L2_CID_MPEG_STREAM_TYPE is a CtrlID. Here it's typed as the enum itself.
	CtrlMPEGStreamType MPEGStreamType = C.V4L2_CID_MPEG_STREAM_TYPE
	// CtrlMPEGStreamPIDPMT sets the Program Map Table PID for the MPEG stream.
	CtrlMPEGStreamPIDPMT CtrlID = C.V4L2_CID_MPEG_STREAM_PID_PMT
	// CtrlMPEGStreamPIDAudio sets the Audio PID for the MPEG stream.
	CtrlMPEGStreamPIDAudio CtrlID = C.V4L2_CID_MPEG_STREAM_PID_AUDIO
	// CtrlMPEGStreamPIDVideo sets the Video PID for the MPEG stream.
	CtrlMPEGStreamPIDVideo CtrlID = C.V4L2_CID_MPEG_STREAM_PID_VIDEO
	// CtrlMPEGStreamPIDPCR sets the Program Clock Reference PID for the MPEG stream.
	CtrlMPEGStreamPIDPCR CtrlID = C.V4L2_CID_MPEG_STREAM_PID_PCR
	// CtrlMPEGStreamPIDPESAudio sets the Audio PES (Packetized Elementary Stream) ID.
	CtrlMPEGStreamPIDPESAudio CtrlID = C.V4L2_CID_MPEG_STREAM_PES_ID_AUDIO
	// CtrlMPEGStreamPESVideo sets the Video PES ID.
	CtrlMPEGStreamPESVideo CtrlID = C.V4L2_CID_MPEG_STREAM_PES_ID_VIDEO
	// CtrlMEPGStreamVBIFormat sets the VBI (Vertical Blanking Interval) data format in the MPEG stream.
	CtrlMEPGStreamVBIFormat CtrlID = C.V4L2_CID_MPEG_STREAM_VBI_FMT

	// CtrlMPEGVideoEncoding sets the video encoding format (e.g., MPEG-1, MPEG-2, H.264).
	// Note: In the C header, V4L2_CID_MPEG_VIDEO_ENCODING is a CtrlID. Here it's typed as the enum itself.
	CtrlMPEGVideoEncoding MPEGVideoEncoding = C.V4L2_CID_MPEG_VIDEO_ENCODING
	// CtrlMPEGVideoAspect sets the video aspect ratio (e.g., 4:3, 16:9).
	// Note: In the C header, V4L2_CID_MPEG_VIDEO_ASPECT is a CtrlID. Here it's typed as the enum itself.
	CtrlMPEGVideoAspect MPEGVideoAspect = C.V4L2_CID_MPEG_VIDEO_ASPECT
	// CtrlMPEGVideoBFrames sets the number of B-frames between I/P frames.
	CtrlMPEGVideoBFrames CtrlID = C.V4L2_CID_MPEG_VIDEO_B_FRAMES
	// CtrlMPEGVideoGOPSize sets the Group of Pictures (GOP) size.
	CtrlMPEGVideoGOPSize CtrlID = C.V4L2_CID_MPEG_VIDEO_GOP_SIZE
	// CtrlMPEGVideoGOPClosure sets whether the GOP is closed or open.
	CtrlMPEGVideoGOPClosure CtrlID = C.V4L2_CID_MPEG_VIDEO_GOP_CLOSURE
	// CtrlMPEGVideoPulldown enables or disables 3:2 pulldown.
	CtrlMPEGVideoPulldown CtrlID = C.V4L2_CID_MPEG_VIDEO_PULLDOWN
	// CtrlMPEGVideoBitrateMode sets the video bitrate mode (e.g., VBR, CBR).
	// Note: In the C header, V4L2_CID_MPEG_VIDEO_BITRATE_MODE is a CtrlID. Here it's typed as the enum itself.
	CtrlMPEGVideoBitrateMode MPEGVideoBitrateMode = C.V4L2_CID_MPEG_VIDEO_BITRATE_MODE
	// CtrlMPEGVideoBitrate sets the video bitrate in bits per second.
	CtrlMPEGVideoBitrate CtrlID = C.V4L2_CID_MPEG_VIDEO_BITRATE
	// CtrlMPEGVideoBitratePeak sets the peak video bitrate for VBR mode.
	CtrlMPEGVideoBitratePeak CtrlID = C.V4L2_CID_MPEG_VIDEO_BITRATE_PEAK
	// CtrlMPEGVideoTemporalDecimation sets the temporal decimation factor.
	CtrlMPEGVideoTemporalDecimation CtrlID = C.V4L2_CID_MPEG_VIDEO_TEMPORAL_DECIMATION
	// CtrlMPEGVideoMute mutes or unmutes the video.
	CtrlMPEGVideoMute CtrlID = C.V4L2_CID_MPEG_VIDEO_MUTE
	// CtrlMPEGVideoMutYUV sets the YUV value to use when video is muted.
	CtrlMPEGVideoMutYUV CtrlID = C.V4L2_CID_MPEG_VIDEO_MUTE_YUV // Typo in original V4L2 define (MutYUV -> MuteYUV)

	// CtrlMPEGVideoDecoderSliceInterface enables/disables the slice interface for decoders.
	CtrlMPEGVideoDecoderSliceInterface CtrlID = C.V4L2_CID_MPEG_VIDEO_DECODER_SLICE_INTERFACE
	// CtrlMPEGVideoDecoderMPEG4DeblockFilter enables/disables the MPEG-4 deblocking filter for decoders.
	CtrlMPEGVideoDecoderMPEG4DeblockFilter CtrlID = C.V4L2_CID_MPEG_VIDEO_DECODER_MPEG4_DEBLOCK_FILTER
	// CtrlMPEGVideoCyclicIntraRefreshMB sets the number of macroblocks per cyclic intra refresh.
	CtrlMPEGVideoCyclicIntraRefreshMB CtrlID = C.V4L2_CID_MPEG_VIDEO_CYCLIC_INTRA_REFRESH_MB
	// CtrlMPEGVideoFrameRCEnable enables/disables frame-level rate control.
	CtrlMPEGVideoFrameRCEnable CtrlID = C.V4L2_CID_MPEG_VIDEO_FRAME_RC_ENABLE
	// CtrlMPEGVideoHeaderMode sets the video header mode (e.g., separate headers, joined).
	CtrlMPEGVideoHeaderMode CtrlID = C.V4L2_CID_MPEG_VIDEO_HEADER_MODE
	// CtrlMPEGVideoMaxRefPic sets the maximum number of reference pictures.
	CtrlMPEGVideoMaxRefPic CtrlID = C.V4L2_CID_MPEG_VIDEO_MAX_REF_PIC
	// CtrlMPEGVideoMBRCEnable enables/disables macroblock-level rate control.
	CtrlMPEGVideoMBRCEnable CtrlID = C.V4L2_CID_MPEG_VIDEO_MB_RC_ENABLE
	// CtrlMPEGVideoMultiSliceMaxBytes sets the maximum bytes per slice for multi-slice mode.
	CtrlMPEGVideoMultiSliceMaxBytes CtrlID = C.V4L2_CID_MPEG_VIDEO_MULTI_SLICE_MAX_BYTES
	// CtrlMPEGVideoMultiSliceMaxMB sets the maximum macroblocks per slice for multi-slice mode.
	CtrlMPEGVideoMultiSliceMaxMB CtrlID = C.V4L2_CID_MPEG_VIDEO_MULTI_SLICE_MAX_MB
	// CtrlMPEGVideoMultiSliceMode sets the multi-slice mode.
	CtrlMPEGVideoMultiSliceMode CtrlID = C.V4L2_CID_MPEG_VIDEO_MULTI_SLICE_MODE

	// TODO (vladimir) add remainder codec controls from <linux/v4l2-controls.h>, there are a lot more!
)
