package v4l2

// audio_info.go provides audio input and output enumeration and selection.
//
// Audio inputs and outputs represent audio connections associated with video devices.
// For example, a TV tuner card might have multiple audio inputs (TV tuner audio, line-in),
// and a video device might have audio outputs (line-out, speaker).
//
// The V4L2 API provides the following operations:
//   - Enumerate available audio inputs/outputs
//   - Query audio capabilities and modes
//   - Select active audio input/output
//
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumaudio.html
// See: https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumaudout.html

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// AudioCapability represents audio capability flags
type AudioCapability = uint32

const (
	AudioCapStereo AudioCapability = C.V4L2_AUDCAP_STEREO // Stereo audio
	AudioCapAVL    AudioCapability = C.V4L2_AUDCAP_AVL    // Automatic Volume Level
)

// AudioMode represents audio mode flags
type AudioMode = uint32

const (
	AudioModeAVL AudioMode = C.V4L2_AUDMODE_AVL // Automatic Volume Level mode
)

// AudioInfo (v4l2_audio)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumaudio.html
type AudioInfo struct {
	v4l2Audio C.struct_v4l2_audio
}

func (a AudioInfo) GetIndex() uint32 {
	return uint32(a.v4l2Audio.index)
}

func (a AudioInfo) GetName() string {
	return C.GoString((*C.char)(unsafe.Pointer(&a.v4l2Audio.name[0])))
}

func (a AudioInfo) GetCapability() AudioCapability {
	return AudioCapability(a.v4l2Audio.capability)
}

func (a AudioInfo) GetMode() AudioMode {
	return AudioMode(a.v4l2Audio.mode)
}

// HasCapability checks if the audio input has a specific capability
func (a AudioInfo) HasCapability(cap AudioCapability) bool {
	return (a.v4l2Audio.capability & C.uint(cap)) != 0
}

// IsStereo returns true if the audio input supports stereo
func (a AudioInfo) IsStereo() bool {
	return a.HasCapability(AudioCapStereo)
}

// HasAVL returns true if the audio input supports Automatic Volume Level
func (a AudioInfo) HasAVL() bool {
	return a.HasCapability(AudioCapAVL)
}

// AudioOutInfo (v4l2_audioout)
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h
// https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumaudout.html
type AudioOutInfo struct {
	v4l2AudioOut C.struct_v4l2_audioout
}

func (a AudioOutInfo) GetIndex() uint32 {
	return uint32(a.v4l2AudioOut.index)
}

func (a AudioOutInfo) GetName() string {
	return C.GoString((*C.char)(unsafe.Pointer(&a.v4l2AudioOut.name[0])))
}

func (a AudioOutInfo) GetCapability() AudioCapability {
	return AudioCapability(a.v4l2AudioOut.capability)
}

func (a AudioOutInfo) GetMode() AudioMode {
	return AudioMode(a.v4l2AudioOut.mode)
}

// HasCapability checks if the audio output has a specific capability
func (a AudioOutInfo) HasCapability(cap AudioCapability) bool {
	return (a.v4l2AudioOut.capability & C.uint(cap)) != 0
}

// IsStereo returns true if the audio output supports stereo
func (a AudioOutInfo) IsStereo() bool {
	return a.HasCapability(AudioCapStereo)
}

// HasAVL returns true if the audio output supports Automatic Volume Level
func (a AudioOutInfo) HasAVL() bool {
	return a.HasCapability(AudioCapAVL)
}

// GetAudioInfo returns specified audio input information for the device
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumaudio.html
func GetAudioInfo(fd uintptr, index uint32) (AudioInfo, error) {
	var audio C.struct_v4l2_audio
	audio.index = C.uint(index)
	if err := send(fd, C.VIDIOC_ENUMAUDIO, uintptr(unsafe.Pointer(&audio))); err != nil {
		return AudioInfo{}, fmt.Errorf("audio input info: index %d: %w", index, err)
	}
	return AudioInfo{v4l2Audio: audio}, nil
}

// GetAllAudioInfo returns all audio input information for device by
// iterating from audio index = 0 until an error (EINVAL) is returned.
func GetAllAudioInfo(fd uintptr) (result []AudioInfo, err error) {
	index := uint32(0)
	for {
		var audio C.struct_v4l2_audio
		audio.index = C.uint(index)
		if err = send(fd, C.VIDIOC_ENUMAUDIO, uintptr(unsafe.Pointer(&audio))); err != nil {
			if errno, ok := err.(sys.Errno); ok && errno == sys.EINVAL && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("all audio input info: %w", err)
		}
		result = append(result, AudioInfo{v4l2Audio: audio})
		index++
	}
	return result, nil
}

// GetCurrentAudio returns the currently selected audio input
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-audio.html
func GetCurrentAudio(fd uintptr) (AudioInfo, error) {
	var audio C.struct_v4l2_audio
	if err := send(fd, C.VIDIOC_G_AUDIO, uintptr(unsafe.Pointer(&audio))); err != nil {
		return AudioInfo{}, fmt.Errorf("audio input get: %w", err)
	}
	return AudioInfo{v4l2Audio: audio}, nil
}

// SetAudio sets the current audio input by index
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-s-audio.html
func SetAudio(fd uintptr, index uint32) error {
	var audio C.struct_v4l2_audio
	audio.index = C.uint(index)
	if err := send(fd, C.VIDIOC_S_AUDIO, uintptr(unsafe.Pointer(&audio))); err != nil {
		return fmt.Errorf("audio input set: index %d: %w", index, err)
	}
	return nil
}

// SetAudioMode sets the audio input mode for the current audio input
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-s-audio.html
func SetAudioMode(fd uintptr, mode AudioMode) error {
	// Get current audio first
	current, err := GetCurrentAudio(fd)
	if err != nil {
		return err
	}

	// Set the mode
	var audio C.struct_v4l2_audio
	audio.index = C.uint(current.GetIndex())
	audio.mode = C.uint(mode)

	if err := send(fd, C.VIDIOC_S_AUDIO, uintptr(unsafe.Pointer(&audio))); err != nil {
		return fmt.Errorf("audio input set mode: %w", err)
	}
	return nil
}

// GetAudioOutInfo returns specified audio output information for the device
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-enumaudout.html
func GetAudioOutInfo(fd uintptr, index uint32) (AudioOutInfo, error) {
	var audioOut C.struct_v4l2_audioout
	audioOut.index = C.uint(index)
	if err := send(fd, C.VIDIOC_ENUMAUDOUT, uintptr(unsafe.Pointer(&audioOut))); err != nil {
		return AudioOutInfo{}, fmt.Errorf("audio output info: index %d: %w", index, err)
	}
	return AudioOutInfo{v4l2AudioOut: audioOut}, nil
}

// GetAllAudioOutInfo returns all audio output information for device by
// iterating from audio output index = 0 until an error (EINVAL) is returned.
func GetAllAudioOutInfo(fd uintptr) (result []AudioOutInfo, err error) {
	index := uint32(0)
	for {
		var audioOut C.struct_v4l2_audioout
		audioOut.index = C.uint(index)
		if err = send(fd, C.VIDIOC_ENUMAUDOUT, uintptr(unsafe.Pointer(&audioOut))); err != nil {
			if errno, ok := err.(sys.Errno); ok && errno == sys.EINVAL && len(result) > 0 {
				break
			}
			return result, fmt.Errorf("all audio output info: %w", err)
		}
		result = append(result, AudioOutInfo{v4l2AudioOut: audioOut})
		index++
	}
	return result, nil
}

// GetCurrentAudioOut returns the currently selected audio output
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-audout.html
func GetCurrentAudioOut(fd uintptr) (AudioOutInfo, error) {
	var audioOut C.struct_v4l2_audioout
	if err := send(fd, C.VIDIOC_G_AUDOUT, uintptr(unsafe.Pointer(&audioOut))); err != nil {
		return AudioOutInfo{}, fmt.Errorf("audio output get: %w", err)
	}
	return AudioOutInfo{v4l2AudioOut: audioOut}, nil
}

// SetAudioOut sets the current audio output by index
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-s-audout.html
func SetAudioOut(fd uintptr, index uint32) error {
	var audioOut C.struct_v4l2_audioout
	audioOut.index = C.uint(index)
	if err := send(fd, C.VIDIOC_S_AUDOUT, uintptr(unsafe.Pointer(&audioOut))); err != nil {
		return fmt.Errorf("audio output set: index %d: %w", index, err)
	}
	return nil
}

// SetAudioOutMode sets the audio output mode for the current audio output
// See https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-s-audout.html
func SetAudioOutMode(fd uintptr, mode AudioMode) error {
	// Get current audio output first
	current, err := GetCurrentAudioOut(fd)
	if err != nil {
		return err
	}

	// Set the mode
	var audioOut C.struct_v4l2_audioout
	audioOut.index = C.uint(current.GetIndex())
	audioOut.mode = C.uint(mode)

	if err := send(fd, C.VIDIOC_S_AUDOUT, uintptr(unsafe.Pointer(&audioOut))); err != nil {
		return fmt.Errorf("audio output set mode: %w", err)
	}
	return nil
}
