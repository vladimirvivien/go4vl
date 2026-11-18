package v4l2

// MPEGStreamType represents v4l2_mpeg_stream_type
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L237
type MPEGStreamType = uint32

const (
	MPEGStreamTypeMPEG2ProgramStream   MPEGStreamType = cV4L2_MPEG_STREAM_TYPE_MPEG2_PS
	MPEGStreamTypeMPEG2TransportStream MPEGStreamType = cV4L2_MPEG_STREAM_TYPE_MPEG2_TS
	MPEGStreamTypeMPEG1SystemStream    MPEGStreamType = cV4L2_MPEG_STREAM_TYPE_MPEG1_SS
	MPEGStreamTypeMPEG2DVD             MPEGStreamType = cV4L2_MPEG_STREAM_TYPE_MPEG2_DVD
	MPEGStreamTypeMPEG1VCD             MPEGStreamType = cV4L2_MPEG_STREAM_TYPE_MPEG1_VCD
	MPEGStreamTypeMPEG2SVCD            MPEGStreamType = cV4L2_MPEG_STREAM_TYPE_MPEG2_SVCD
)

type MPEGVideoEncoding = uint32

const (
	MPEGVideoEncodingMPEG1    MPEGVideoEncoding = cV4L2_MPEG_VIDEO_ENCODING_MPEG_1
	MPEGVideoEncodingMPEG2    MPEGVideoEncoding = cV4L2_MPEG_VIDEO_ENCODING_MPEG_2
	MPEGVideoEncodingMPEG4AVC MPEGVideoEncoding = cV4L2_MPEG_VIDEO_ENCODING_MPEG_4_AVC
)

type MPEGVideoAspect = uint32

const (
	MPEGVideoAspect1x1     MPEGVideoAspect = cV4L2_MPEG_VIDEO_ASPECT_1x1
	MPEGVideoAspect4x3     MPEGVideoAspect = cV4L2_MPEG_VIDEO_ASPECT_4x3
	MPEGVideoAspect16x9    MPEGVideoAspect = cV4L2_MPEG_VIDEO_ASPECT_16x9
	MPEGVideoAspect221x100 MPEGVideoAspect = cV4L2_MPEG_VIDEO_ASPECT_221x100
)

type MPEGVideoBitrateMode = uint32

const (
	MPEGVideoBitrateModeVBR = cV4L2_MPEG_VIDEO_BITRATE_MODE_VBR
	MPEGVideoBitrateModeCBR = cV4L2_MPEG_VIDEO_BITRATE_MODE_CBR
	MPEGVideoBitrateModeCQ  = cV4L2_MPEG_VIDEO_BITRATE_MODE_CQ
)

// Codec control values
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/ext-ctrls-codec.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/v4l2-controls.h#L228
const (
	CtrlCodecClass                         CtrlID               = cV4L2_CID_CODEC_CLASS
	CtrlMPEGStreamType                     MPEGStreamType       = cV4L2_CID_MPEG_STREAM_TYPE
	CtrlMPEGStreamPIDPMT                   CtrlID               = cV4L2_CID_MPEG_STREAM_PID_PMT
	CtrlMPEGStreamPIDAudio                 CtrlID               = cV4L2_CID_MPEG_STREAM_PID_AUDIO
	CtrlMPEGStreamPIDVideo                 CtrlID               = cV4L2_CID_MPEG_STREAM_PID_VIDEO
	CtrlMPEGStreamPIDPCR                   CtrlID               = cV4L2_CID_MPEG_STREAM_PID_PCR
	CtrlMPEGStreamPIDPESAudio              CtrlID               = cV4L2_CID_MPEG_STREAM_PES_ID_AUDIO
	CtrlMPEGStreamPESVideo                 CtrlID               = cV4L2_CID_MPEG_STREAM_PES_ID_VIDEO
	CtrlMEPGStreamVBIFormat                CtrlID               = cV4L2_CID_MPEG_STREAM_VBI_FMT
	CtrlMPEGVideoEncoding                  MPEGVideoEncoding    = cV4L2_CID_MPEG_VIDEO_ENCODING
	CtrlMPEGVideoAspect                    MPEGVideoAspect      = cV4L2_CID_MPEG_VIDEO_ASPECT
	CtrlMPEGVideoBFrames                   CtrlID               = cV4L2_CID_MPEG_VIDEO_B_FRAMES
	CtrlMPEGVideoGOPSize                   CtrlID               = cV4L2_CID_MPEG_VIDEO_GOP_SIZE
	CtrlMPEGVideoGOPClosure                CtrlID               = cV4L2_CID_MPEG_VIDEO_GOP_CLOSURE
	CtrlMPEGVideoPulldown                  CtrlID               = cV4L2_CID_MPEG_VIDEO_PULLDOWN
	CtrlMPEGVideoBitrateMode               MPEGVideoBitrateMode = cV4L2_CID_MPEG_VIDEO_BITRATE_MODE
	CtrlMPEGVideoBitrate                   CtrlID               = cV4L2_CID_MPEG_VIDEO_BITRATE
	CtrlMPEGVideoBitratePeak               CtrlID               = cV4L2_CID_MPEG_VIDEO_BITRATE_PEAK
	CtrlMPEGVideoTemporalDecimation        CtrlID               = cV4L2_CID_MPEG_VIDEO_TEMPORAL_DECIMATION
	CtrlMPEGVideoMute                      CtrlID               = cV4L2_CID_MPEG_VIDEO_MUTE
	CtrlMPEGVideoMutYUV                    CtrlID               = cV4L2_CID_MPEG_VIDEO_MUTE_YUV
	CtrlMPEGVideoDecoderSliceInterface     CtrlID               = cV4L2_CID_MPEG_VIDEO_DECODER_SLICE_INTERFACE
	CtrlMPEGVideoDecoderMPEG4DeblockFilter CtrlID               = cV4L2_CID_MPEG_VIDEO_DECODER_MPEG4_DEBLOCK_FILTER
	CtrlMPEGVideoCyclicIntraRefreshMB      CtrlID               = cV4L2_CID_MPEG_VIDEO_CYCLIC_INTRA_REFRESH_MB
	CtrlMPEGVideoFrameRCEnable             CtrlID               = cV4L2_CID_MPEG_VIDEO_FRAME_RC_ENABLE
	CtrlMPEGVideoHeaderMode                CtrlID               = cV4L2_CID_MPEG_VIDEO_HEADER_MODE
	CtrlMPEGVideoMaxRefPic                 CtrlID               = cV4L2_CID_MPEG_VIDEO_MAX_REF_PIC
	CtrlMPEGVideoMBRCEnable                CtrlID               = cV4L2_CID_MPEG_VIDEO_MB_RC_ENABLE
	CtrlMPEGVideoMultiSliceMaxBytes        CtrlID               = cV4L2_CID_MPEG_VIDEO_MULTI_SLICE_MAX_BYTES
	CtrlMPEGVideoMultiSliceMaxMB           CtrlID               = cV4L2_CID_MPEG_VIDEO_MULTI_SLICE_MAX_MB
	CtrlMPEGVideoMultiSliceMode            CtrlID               = cV4L2_CID_MPEG_VIDEO_MULTI_SLICE_MODE

	// TODO (vladimir) add remainder codec, there are a lot more!
)
