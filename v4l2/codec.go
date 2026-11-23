package v4l2

// codec.go provides V4L2 stateful codec support for hardware-accelerated video encoding/decoding.
//
// V4L2 stateful codecs use a memory-to-memory (M2M) architecture where:
// - For encoding: raw video (e.g., NV12) goes in, compressed video (e.g., H.264) comes out
// - For decoding: compressed video goes in, raw video comes out
//
// The codec interface provides commands for controlling the codec state machine:
// - START: Begin encoding/decoding
// - STOP: Stop and drain remaining data
// - PAUSE: Pause processing
// - RESUME: Resume from pause
// - FLUSH: (decoder only) Flush decoder state for seeking
//
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/dev-encoder.html
// See: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/dev-decoder.html

import (
	"fmt"
	"unsafe"
)

// Encoder command constants for VIDIOC_ENCODER_CMD.
// These commands control the state of a V4L2 stateful video encoder.
//
// Reference: https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
const (
	// EncCmdStart starts the encoder. Encoding will begin with the next queued buffer.
	EncCmdStart uint32 = cV4L2_ENC_CMD_START

	// EncCmdStop stops the encoder. The encoder will finish encoding all queued buffers
	// and then signal end-of-stream on the capture queue.
	EncCmdStop uint32 = cV4L2_ENC_CMD_STOP

	// EncCmdPause pauses the encoder. Buffers in the queue will not be processed
	// until the encoder is resumed.
	EncCmdPause uint32 = cV4L2_ENC_CMD_PAUSE

	// EncCmdResume resumes a paused encoder.
	EncCmdResume uint32 = cV4L2_ENC_CMD_RESUME
)

// Encoder command flags for VIDIOC_ENCODER_CMD.
const (
	// EncCmdStopAtGOPEnd stops encoding at the next GOP boundary.
	// If not set, encoding stops immediately after the current frame.
	EncCmdStopAtGOPEnd uint32 = cV4L2_ENC_CMD_STOP_AT_GOP_END
)

// Decoder command constants for VIDIOC_DECODER_CMD.
// These commands control the state of a V4L2 stateful video decoder.
//
// Reference: https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
const (
	// DecCmdStart starts the decoder. Decoding will begin with the next queued buffer.
	DecCmdStart uint32 = cV4L2_DEC_CMD_START

	// DecCmdStop stops the decoder. The decoder will finish decoding all queued buffers
	// and then signal end-of-stream on the capture queue.
	DecCmdStop uint32 = cV4L2_DEC_CMD_STOP

	// DecCmdPause pauses the decoder. Buffers in the queue will not be processed
	// until the decoder is resumed.
	DecCmdPause uint32 = cV4L2_DEC_CMD_PAUSE

	// DecCmdResume resumes a paused decoder.
	DecCmdResume uint32 = cV4L2_DEC_CMD_RESUME

	// DecCmdFlush flushes the decoder. All pending decoded frames are returned,
	// and the decoder state is reset for seeking or stream changes.
	DecCmdFlush uint32 = cV4L2_DEC_CMD_FLUSH
)

// Decoder command flags for VIDIOC_DECODER_CMD.
const (
	// DecCmdStartMuteAudio mutes audio during playback start.
	DecCmdStartMuteAudio uint32 = cV4L2_DEC_CMD_START_MUTE_AUDIO

	// DecCmdPauseToBlack displays a black frame when pausing.
	DecCmdPauseToBlack uint32 = cV4L2_DEC_CMD_PAUSE_TO_BLACK

	// DecCmdStopToBlack displays a black frame when stopping.
	DecCmdStopToBlack uint32 = cV4L2_DEC_CMD_STOP_TO_BLACK

	// DecCmdStopImmediately stops immediately without draining.
	DecCmdStopImmediately uint32 = cV4L2_DEC_CMD_STOP_IMMEDIATELY
)

// Decoder start format requirements (returned by driver).
const (
	// DecStartFmtNone indicates the decoder has no special format requirements.
	DecStartFmtNone uint32 = cV4L2_DEC_START_FMT_NONE

	// DecStartFmtGOP indicates the decoder requires full GOPs.
	DecStartFmtGOP uint32 = cV4L2_DEC_START_FMT_GOP
)

// EncoderCmd represents a V4L2 encoder command (v4l2_encoder_cmd).
// Used with VIDIOC_ENCODER_CMD and VIDIOC_TRY_ENCODER_CMD ioctls.
//
// Reference: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-encoder-cmd.html
type EncoderCmd struct {
	// Cmd is the encoder command (EncCmdStart, EncCmdStop, etc.)
	Cmd uint32

	// Flags are command-specific flags
	Flags uint32

	// Raw is reserved space for future extensions
	Raw [8]uint32
}

// NewEncoderCmd creates a new encoder command.
func NewEncoderCmd(cmd uint32) *EncoderCmd {
	return &EncoderCmd{Cmd: cmd}
}

// NewEncoderCmdWithFlags creates a new encoder command with flags.
func NewEncoderCmdWithFlags(cmd, flags uint32) *EncoderCmd {
	return &EncoderCmd{Cmd: cmd, Flags: flags}
}

// GetCmd returns the command.
func (ec *EncoderCmd) GetCmd() uint32 {
	return ec.Cmd
}

// SetCmd sets the command.
func (ec *EncoderCmd) SetCmd(cmd uint32) {
	ec.Cmd = cmd
}

// GetFlags returns the flags.
func (ec *EncoderCmd) GetFlags() uint32 {
	return ec.Flags
}

// SetFlags sets the flags.
func (ec *EncoderCmd) SetFlags(flags uint32) {
	ec.Flags = flags
}

// IsStart returns true if this is a start command.
func (ec *EncoderCmd) IsStart() bool {
	return ec.Cmd == EncCmdStart
}

// IsStop returns true if this is a stop command.
func (ec *EncoderCmd) IsStop() bool {
	return ec.Cmd == EncCmdStop
}

// IsPause returns true if this is a pause command.
func (ec *EncoderCmd) IsPause() bool {
	return ec.Cmd == EncCmdPause
}

// IsResume returns true if this is a resume command.
func (ec *EncoderCmd) IsResume() bool {
	return ec.Cmd == EncCmdResume
}

// DecoderCmd represents a V4L2 decoder command (v4l2_decoder_cmd).
// Used with VIDIOC_DECODER_CMD and VIDIOC_TRY_DECODER_CMD ioctls.
//
// Reference: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-decoder-cmd.html
type DecoderCmd struct {
	// Cmd is the decoder command (DecCmdStart, DecCmdStop, etc.)
	Cmd uint32

	// Flags are command-specific flags
	Flags uint32

	// StopPts is the presentation timestamp for stop command
	StopPts int64

	// StartSpeed is the playback speed for start command (1000 = normal speed)
	StartSpeed int32

	// StartFormat is the required input format (DecStartFmtNone, DecStartFmtGOP)
	StartFormat uint32
}

// NewDecoderCmd creates a new decoder command.
func NewDecoderCmd(cmd uint32) *DecoderCmd {
	return &DecoderCmd{Cmd: cmd}
}

// NewDecoderCmdWithFlags creates a new decoder command with flags.
func NewDecoderCmdWithFlags(cmd, flags uint32) *DecoderCmd {
	return &DecoderCmd{Cmd: cmd, Flags: flags}
}

// GetCmd returns the command.
func (dc *DecoderCmd) GetCmd() uint32 {
	return dc.Cmd
}

// SetCmd sets the command.
func (dc *DecoderCmd) SetCmd(cmd uint32) {
	dc.Cmd = cmd
}

// GetFlags returns the flags.
func (dc *DecoderCmd) GetFlags() uint32 {
	return dc.Flags
}

// SetFlags sets the flags.
func (dc *DecoderCmd) SetFlags(flags uint32) {
	dc.Flags = flags
}

// IsStart returns true if this is a start command.
func (dc *DecoderCmd) IsStart() bool {
	return dc.Cmd == DecCmdStart
}

// IsStop returns true if this is a stop command.
func (dc *DecoderCmd) IsStop() bool {
	return dc.Cmd == DecCmdStop
}

// IsPause returns true if this is a pause command.
func (dc *DecoderCmd) IsPause() bool {
	return dc.Cmd == DecCmdPause
}

// IsResume returns true if this is a resume command.
func (dc *DecoderCmd) IsResume() bool {
	return dc.Cmd == DecCmdResume
}

// IsFlush returns true if this is a flush command.
func (dc *DecoderCmd) IsFlush() bool {
	return dc.Cmd == DecCmdFlush
}

// GetStopPts returns the PTS for stop command.
func (dc *DecoderCmd) GetStopPts() int64 {
	return dc.StopPts
}

// SetStopPts sets the PTS for stop command.
func (dc *DecoderCmd) SetStopPts(pts int64) {
	dc.StopPts = pts
}

// GetStartSpeed returns the playback speed.
// 0 or 1000 = normal speed, 1 = forward single step, -1 = backward single step.
func (dc *DecoderCmd) GetStartSpeed() int32 {
	return dc.StartSpeed
}

// SetStartSpeed sets the playback speed.
func (dc *DecoderCmd) SetStartSpeed(speed int32) {
	dc.StartSpeed = speed
}

// GetStartFormat returns the required input format.
func (dc *DecoderCmd) GetStartFormat() uint32 {
	return dc.StartFormat
}

// SetStartFormat sets the required input format.
func (dc *DecoderCmd) SetStartFormat(format uint32) {
	dc.StartFormat = format
}

// SendEncoderCmd sends an encoder command to the device.
// This ioctl controls the state of a V4L2 video encoder.
//
// Parameters:
//   - fd: File descriptor of an opened encoder device
//   - cmd: Encoder command to send
//
// Returns:
//   - error: An error if the ioctl fails
//
// This function issues the VIDIOC_ENCODER_CMD ioctl.
// Reference: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-encoder-cmd.html
func SendEncoderCmd(fd uintptr, cmd *EncoderCmd) error {
	var v4l2Cmd cV4L2EncoderCmd
	v4l2Cmd.Cmd = cmd.Cmd
	v4l2Cmd.Flags = cmd.Flags

	if err := send(fd, cVIDIOC_ENCODER_CMD, unsafe.Pointer(&v4l2Cmd)); err != nil {
		return fmt.Errorf("encoder cmd: %w", err)
	}

	// Update command with any values modified by driver
	cmd.Cmd = v4l2Cmd.Cmd
	cmd.Flags = v4l2Cmd.Flags

	return nil
}

// TryEncoderCmd tests an encoder command without executing it.
// This can be used to check if a command is supported by the device.
//
// Parameters:
//   - fd: File descriptor of an opened encoder device
//   - cmd: Encoder command to test
//
// Returns:
//   - error: An error if the command is not supported
//
// This function issues the VIDIOC_TRY_ENCODER_CMD ioctl.
func TryEncoderCmd(fd uintptr, cmd *EncoderCmd) error {
	var v4l2Cmd cV4L2EncoderCmd
	v4l2Cmd.Cmd = cmd.Cmd
	v4l2Cmd.Flags = cmd.Flags

	if err := send(fd, cVIDIOC_TRY_ENCODER_CMD, unsafe.Pointer(&v4l2Cmd)); err != nil {
		return fmt.Errorf("try encoder cmd: %w", err)
	}

	// Update command with any values modified by driver
	cmd.Cmd = v4l2Cmd.Cmd
	cmd.Flags = v4l2Cmd.Flags

	return nil
}

// SendDecoderCmd sends a decoder command to the device.
// This ioctl controls the state of a V4L2 video decoder.
//
// Parameters:
//   - fd: File descriptor of an opened decoder device
//   - cmd: Decoder command to send
//
// Returns:
//   - error: An error if the ioctl fails
//
// This function issues the VIDIOC_DECODER_CMD ioctl.
// Reference: https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-decoder-cmd.html
func SendDecoderCmd(fd uintptr, cmd *DecoderCmd) error {
	var v4l2Cmd cV4L2DecoderCmd
	v4l2Cmd.Cmd = cmd.Cmd
	v4l2Cmd.Flags = cmd.Flags

	// Set union fields based on command type
	switch cmd.Cmd {
	case DecCmdStop:
		// Set stop.pts field
		stopPtr := (*uint64)(unsafe.Pointer(&v4l2Cmd.Stop[0]))
		*stopPtr = uint64(cmd.StopPts)
	case DecCmdStart:
		// Set start.speed and start.format fields
		speedPtr := (*int32)(unsafe.Pointer(&v4l2Cmd.Stop[0]))
		formatPtr := (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&v4l2Cmd.Stop[0])) + 4))
		*speedPtr = cmd.StartSpeed
		*formatPtr = cmd.StartFormat
	}

	if err := send(fd, cVIDIOC_DECODER_CMD, unsafe.Pointer(&v4l2Cmd)); err != nil {
		return fmt.Errorf("decoder cmd: %w", err)
	}

	// Update command with any values modified by driver
	cmd.Cmd = v4l2Cmd.Cmd
	cmd.Flags = v4l2Cmd.Flags

	// Read back union fields
	switch cmd.Cmd {
	case DecCmdStop:
		stopPtr := (*uint64)(unsafe.Pointer(&v4l2Cmd.Stop[0]))
		cmd.StopPts = int64(*stopPtr)
	case DecCmdStart:
		speedPtr := (*int32)(unsafe.Pointer(&v4l2Cmd.Stop[0]))
		formatPtr := (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&v4l2Cmd.Stop[0])) + 4))
		cmd.StartSpeed = *speedPtr
		cmd.StartFormat = *formatPtr
	}

	return nil
}

// TryDecoderCmd tests a decoder command without executing it.
// This can be used to check if a command is supported by the device.
//
// Parameters:
//   - fd: File descriptor of an opened decoder device
//   - cmd: Decoder command to test
//
// Returns:
//   - error: An error if the command is not supported
//
// This function issues the VIDIOC_TRY_DECODER_CMD ioctl.
func TryDecoderCmd(fd uintptr, cmd *DecoderCmd) error {
	var v4l2Cmd cV4L2DecoderCmd
	v4l2Cmd.Cmd = cmd.Cmd
	v4l2Cmd.Flags = cmd.Flags

	// Set union fields based on command type
	switch cmd.Cmd {
	case DecCmdStop:
		stopPtr := (*uint64)(unsafe.Pointer(&v4l2Cmd.Stop[0]))
		*stopPtr = uint64(cmd.StopPts)
	case DecCmdStart:
		speedPtr := (*int32)(unsafe.Pointer(&v4l2Cmd.Stop[0]))
		formatPtr := (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&v4l2Cmd.Stop[0])) + 4))
		*speedPtr = cmd.StartSpeed
		*formatPtr = cmd.StartFormat
	}

	if err := send(fd, cVIDIOC_TRY_DECODER_CMD, unsafe.Pointer(&v4l2Cmd)); err != nil {
		return fmt.Errorf("try decoder cmd: %w", err)
	}

	// Update command with any values modified by driver
	cmd.Cmd = v4l2Cmd.Cmd
	cmd.Flags = v4l2Cmd.Flags

	// Read back union fields
	switch cmd.Cmd {
	case DecCmdStop:
		stopPtr := (*uint64)(unsafe.Pointer(&v4l2Cmd.Stop[0]))
		cmd.StopPts = int64(*stopPtr)
	case DecCmdStart:
		speedPtr := (*int32)(unsafe.Pointer(&v4l2Cmd.Stop[0]))
		formatPtr := (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&v4l2Cmd.Stop[0])) + 4))
		cmd.StartSpeed = *speedPtr
		cmd.StartFormat = *formatPtr
	}

	return nil
}

// StartEncoder sends the START command to an encoder.
// This is a convenience wrapper around SendEncoderCmd.
func StartEncoder(fd uintptr) error {
	cmd := NewEncoderCmd(EncCmdStart)
	return SendEncoderCmd(fd, cmd)
}

// StopEncoder sends the STOP command to an encoder.
// The encoder will drain any remaining buffers and signal EOS.
// This is a convenience wrapper around SendEncoderCmd.
func StopEncoder(fd uintptr) error {
	cmd := NewEncoderCmd(EncCmdStop)
	return SendEncoderCmd(fd, cmd)
}

// StopEncoderAtGOPEnd sends the STOP command with the GOP end flag.
// The encoder will stop at the next GOP boundary.
func StopEncoderAtGOPEnd(fd uintptr) error {
	cmd := NewEncoderCmdWithFlags(EncCmdStop, EncCmdStopAtGOPEnd)
	return SendEncoderCmd(fd, cmd)
}

// PauseEncoder sends the PAUSE command to an encoder.
// This is a convenience wrapper around SendEncoderCmd.
func PauseEncoder(fd uintptr) error {
	cmd := NewEncoderCmd(EncCmdPause)
	return SendEncoderCmd(fd, cmd)
}

// ResumeEncoder sends the RESUME command to an encoder.
// This is a convenience wrapper around SendEncoderCmd.
func ResumeEncoder(fd uintptr) error {
	cmd := NewEncoderCmd(EncCmdResume)
	return SendEncoderCmd(fd, cmd)
}

// StartDecoder sends the START command to a decoder.
// This is a convenience wrapper around SendDecoderCmd.
func StartDecoder(fd uintptr) error {
	cmd := NewDecoderCmd(DecCmdStart)
	return SendDecoderCmd(fd, cmd)
}

// StopDecoder sends the STOP command to a decoder.
// The decoder will drain any remaining buffers and signal EOS.
// This is a convenience wrapper around SendDecoderCmd.
func StopDecoder(fd uintptr) error {
	cmd := NewDecoderCmd(DecCmdStop)
	return SendDecoderCmd(fd, cmd)
}

// StopDecoderImmediately sends the STOP command with immediate flag.
// The decoder will stop without draining remaining buffers.
func StopDecoderImmediately(fd uintptr) error {
	cmd := NewDecoderCmdWithFlags(DecCmdStop, DecCmdStopImmediately)
	return SendDecoderCmd(fd, cmd)
}

// PauseDecoder sends the PAUSE command to a decoder.
// This is a convenience wrapper around SendDecoderCmd.
func PauseDecoder(fd uintptr) error {
	cmd := NewDecoderCmd(DecCmdPause)
	return SendDecoderCmd(fd, cmd)
}

// ResumeDecoder sends the RESUME command to a decoder.
// This is a convenience wrapper around SendDecoderCmd.
func ResumeDecoder(fd uintptr) error {
	cmd := NewDecoderCmd(DecCmdResume)
	return SendDecoderCmd(fd, cmd)
}

// FlushDecoder sends the FLUSH command to a decoder.
// This flushes all pending decoded frames and resets the decoder state.
// Useful for seeking operations.
// This is a convenience wrapper around SendDecoderCmd.
func FlushDecoder(fd uintptr) error {
	cmd := NewDecoderCmd(DecCmdFlush)
	return SendDecoderCmd(fd, cmd)
}

// M2M (Memory-to-Memory) Helper Functions
// These functions support M2M devices like encoders/decoders that have two queues:
// - OUTPUT queue: Where raw/compressed data is sent TO the device
// - CAPTURE queue: Where processed data is received FROM the device

// GetPixFormatOutput retrieves the pixel format for the OUTPUT queue.
// For encoders: This is the raw video input format (e.g., NV12).
// For decoders: This is the compressed input format (e.g., H.264).
func GetPixFormatOutput(fd uintptr) (PixFormat, error) {
	return getPixFormatForType(fd, BufTypeVideoOutput)
}

// SetPixFormatOutput sets the pixel format for the OUTPUT queue.
// For encoders: Set the raw video input format.
// For decoders: Set the compressed input format.
func SetPixFormatOutput(fd uintptr, pixFmt PixFormat) error {
	return setPixFormatForType(fd, BufTypeVideoOutput, pixFmt)
}

// GetPixFormatCapture retrieves the pixel format for the CAPTURE queue.
// For encoders: This is the compressed output format (e.g., H.264).
// For decoders: This is the raw video output format (e.g., NV12).
func GetPixFormatCapture(fd uintptr) (PixFormat, error) {
	return getPixFormatForType(fd, BufTypeVideoCapture)
}

// SetPixFormatCapture sets the pixel format for the CAPTURE queue.
// For encoders: Set the compressed output format.
// For decoders: Set the raw video output format.
func SetPixFormatCapture(fd uintptr, pixFmt PixFormat) error {
	return setPixFormatForType(fd, BufTypeVideoCapture, pixFmt)
}

// getPixFormatForType retrieves the pixel format for a specific buffer type.
func getPixFormatForType(fd uintptr, bufType BufType) (PixFormat, error) {
	var v4l2Format cV4L2Format
	v4l2Format.Type = bufType

	if err := send(fd, cVIDIOC_G_FMT, unsafe.Pointer(&v4l2Format)); err != nil {
		return PixFormat{}, fmt.Errorf("get pix format: %w", err)
	}

	v4l2PixFmt := *(*cV4L2PixFormat)(unsafe.Pointer(&v4l2Format.Fmt[0]))
	return PixFormat{
		Width:        v4l2PixFmt.Width,
		Height:       v4l2PixFmt.Height,
		PixelFormat:  v4l2PixFmt.Pixelformat,
		Field:        v4l2PixFmt.Field,
		BytesPerLine: v4l2PixFmt.Bytesperline,
		SizeImage:    v4l2PixFmt.Sizeimage,
		Colorspace:   v4l2PixFmt.Colorspace,
		Priv:         v4l2PixFmt.Priv,
		Flags:        v4l2PixFmt.Flags,
		YcbcrEnc:     v4l2PixFmt.Ycbcr_enc,
		Quantization: v4l2PixFmt.Quantization,
		XferFunc:     v4l2PixFmt.Xfer_func,
	}, nil
}

// setPixFormatForType sets the pixel format for a specific buffer type.
func setPixFormatForType(fd uintptr, bufType BufType, pixFmt PixFormat) error {
	var v4l2Format cV4L2Format
	v4l2Format.Type = bufType

	v4l2PixFmt := (*cV4L2PixFormat)(unsafe.Pointer(&v4l2Format.Fmt[0]))
	v4l2PixFmt.Width = pixFmt.Width
	v4l2PixFmt.Height = pixFmt.Height
	v4l2PixFmt.Pixelformat = pixFmt.PixelFormat
	v4l2PixFmt.Field = pixFmt.Field
	v4l2PixFmt.Bytesperline = pixFmt.BytesPerLine
	v4l2PixFmt.Sizeimage = pixFmt.SizeImage
	v4l2PixFmt.Colorspace = pixFmt.Colorspace

	if err := send(fd, cVIDIOC_S_FMT, unsafe.Pointer(&v4l2Format)); err != nil {
		return fmt.Errorf("set pix format: %w", err)
	}
	return nil
}

// RequestBuffersOutput allocates buffers for the OUTPUT queue.
func RequestBuffersOutput(fd uintptr, count uint32, ioType IOType) (RequestBuffers, error) {
	return requestBuffersForType(fd, BufTypeVideoOutput, count, ioType)
}

// RequestBuffersCapture allocates buffers for the CAPTURE queue.
func RequestBuffersCapture(fd uintptr, count uint32, ioType IOType) (RequestBuffers, error) {
	return requestBuffersForType(fd, BufTypeVideoCapture, count, ioType)
}

// requestBuffersForType allocates buffers for a specific buffer type.
func requestBuffersForType(fd uintptr, bufType BufType, count uint32, ioType IOType) (RequestBuffers, error) {
	var req cV4L2RequestBuffers
	req.Count = count
	req.Type = bufType
	req.Memory = ioType

	if err := send(fd, cVIDIOC_REQBUFS, unsafe.Pointer(&req)); err != nil {
		return RequestBuffers{}, fmt.Errorf("request buffers: %w", err)
	}

	return *(*RequestBuffers)(unsafe.Pointer(&req)), nil
}

// MapMemoryBuffersOutput maps buffers for the OUTPUT queue.
func MapMemoryBuffersOutput(fd uintptr, count uint32) ([][]byte, error) {
	return mapMemoryBuffersForType(fd, BufTypeVideoOutput, count)
}

// MapMemoryBuffersCapture maps buffers for the CAPTURE queue.
func MapMemoryBuffersCapture(fd uintptr, count uint32) ([][]byte, error) {
	return mapMemoryBuffersForType(fd, BufTypeVideoCapture, count)
}

// mapMemoryBuffersForType maps buffers for a specific buffer type.
func mapMemoryBuffersForType(fd uintptr, bufType BufType, count uint32) ([][]byte, error) {
	buffers := make([][]byte, count)
	for i := uint32(0); i < count; i++ {
		var v4l2Buf cV4L2Buffer
		v4l2Buf.Type = bufType
		v4l2Buf.Memory = IOTypeMMAP
		v4l2Buf.Index = i

		if err := send(fd, cVIDIOC_QUERYBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
			return nil, fmt.Errorf("query buffer %d: %w", i, err)
		}

		offset := *(*uint32)(unsafe.Pointer(&v4l2Buf.M[0]))
		length := v4l2Buf.Length

		data, err := mapMemoryBuffer(fd, int64(offset), int(length))
		if err != nil {
			return nil, fmt.Errorf("map buffer %d: %w", i, err)
		}
		buffers[i] = data
	}
	return buffers, nil
}

// QueueBufferOutput queues a buffer on the OUTPUT queue with data.
func QueueBufferOutput(fd uintptr, ioType IOType, index uint32, bytesUsed uint32) (Buffer, error) {
	var v4l2Buf cV4L2Buffer
	v4l2Buf.Type = BufTypeVideoOutput
	v4l2Buf.Memory = ioType
	v4l2Buf.Index = index
	v4l2Buf.Bytesused = bytesUsed

	if err := send(fd, cVIDIOC_QBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
		return Buffer{}, fmt.Errorf("queue buffer: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// QueueBufferCapture queues a buffer on the CAPTURE queue.
func QueueBufferCapture(fd uintptr, ioType IOType, index uint32) (Buffer, error) {
	var v4l2Buf cV4L2Buffer
	v4l2Buf.Type = BufTypeVideoCapture
	v4l2Buf.Memory = ioType
	v4l2Buf.Index = index

	if err := send(fd, cVIDIOC_QBUF, unsafe.Pointer(&v4l2Buf)); err != nil {
		return Buffer{}, fmt.Errorf("queue buffer: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// DequeueBufferOutput dequeues a buffer from the OUTPUT queue.
func DequeueBufferOutput(fd uintptr, ioType IOType) (Buffer, error) {
	return DequeueBuffer(fd, ioType, BufTypeVideoOutput)
}

// DequeueBufferCapture dequeues a buffer from the CAPTURE queue.
func DequeueBufferCapture(fd uintptr, ioType IOType) (Buffer, error) {
	return DequeueBuffer(fd, ioType, BufTypeVideoCapture)
}

// StreamOnOutput starts streaming on the OUTPUT queue.
func StreamOnOutput(fd uintptr) error {
	bufType := BufTypeVideoOutput
	if err := send(fd, cVIDIOC_STREAMON, unsafe.Pointer(&bufType)); err != nil {
		return fmt.Errorf("stream on output: %w", err)
	}
	return nil
}

// StreamOnCapture starts streaming on the CAPTURE queue.
func StreamOnCapture(fd uintptr) error {
	bufType := BufTypeVideoCapture
	if err := send(fd, cVIDIOC_STREAMON, unsafe.Pointer(&bufType)); err != nil {
		return fmt.Errorf("stream on capture: %w", err)
	}
	return nil
}

// StreamOffOutput stops streaming on the OUTPUT queue.
func StreamOffOutput(fd uintptr) error {
	bufType := BufTypeVideoOutput
	if err := send(fd, cVIDIOC_STREAMOFF, unsafe.Pointer(&bufType)); err != nil {
		return fmt.Errorf("stream off output: %w", err)
	}
	return nil
}

// StreamOffCapture stops streaming on the CAPTURE queue.
func StreamOffCapture(fd uintptr) error {
	bufType := BufTypeVideoCapture
	if err := send(fd, cVIDIOC_STREAMOFF, unsafe.Pointer(&bufType)); err != nil {
		return fmt.Errorf("stream off capture: %w", err)
	}
	return nil
}

// IsTemporaryError checks if an error is temporary (EAGAIN, EINTR).
func IsTemporaryError(err error) bool {
	return err == ErrorTemporary || err == ErrorInterrupted
}

// Encoder and decoder command name maps for debugging and display
var (
	// EncoderCmdNames maps encoder command values to human-readable names.
	EncoderCmdNames = map[uint32]string{
		EncCmdStart:  "START",
		EncCmdStop:   "STOP",
		EncCmdPause:  "PAUSE",
		EncCmdResume: "RESUME",
	}

	// DecoderCmdNames maps decoder command values to human-readable names.
	DecoderCmdNames = map[uint32]string{
		DecCmdStart:  "START",
		DecCmdStop:   "STOP",
		DecCmdPause:  "PAUSE",
		DecCmdResume: "RESUME",
		DecCmdFlush:  "FLUSH",
	}
)

// String returns a human-readable representation of the encoder command.
func (ec *EncoderCmd) String() string {
	name, ok := EncoderCmdNames[ec.Cmd]
	if !ok {
		name = fmt.Sprintf("UNKNOWN(%d)", ec.Cmd)
	}
	return fmt.Sprintf("EncoderCmd{Cmd: %s, Flags: 0x%x}", name, ec.Flags)
}

// String returns a human-readable representation of the decoder command.
func (dc *DecoderCmd) String() string {
	name, ok := DecoderCmdNames[dc.Cmd]
	if !ok {
		name = fmt.Sprintf("UNKNOWN(%d)", dc.Cmd)
	}
	return fmt.Sprintf("DecoderCmd{Cmd: %s, Flags: 0x%x}", name, dc.Flags)
}
