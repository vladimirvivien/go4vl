package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// CropCapability stores information about the cropping capabilities of a V4L2 device.
// It corresponds to the `v4l2_cropcap` struct in the Linux kernel.
// This structure defines the cropping boundaries, the default cropping rectangle,
// and the pixel aspect ratio for a given stream type (e.g., video capture, video output).
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-cropcap.html#c.v4l2_cropcap
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1221
type CropCapability struct {
	// StreamType is the type of data stream (e.g., v4l2_buf_type_video_capture).
	// This field is set by the application to specify which stream's capabilities to query.
	StreamType uint32
	// Bounds defines the outer limits of the cropping area.
	Bounds Rect
	// DefaultRect is the default cropping rectangle.
	DefaultRect Rect
	// PixelAspect is the pixel aspect ratio (width/height).
	PixelAspect Fract
	// reserved space in C struct
	_ [4]uint32
}

// GetCropCapability retrieves the cropping capabilities for a specified buffer type on the device.
// It takes the file descriptor of the V4L2 device and the buffer type (e.g., BufTypeVideoCapture).
// It returns a CropCapability struct populated with the device's cropping information and an error if the query fails.
//
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-cropcap.html#ioctl-vidioc-cropcap
func GetCropCapability(fd uintptr, bufType BufType) (CropCapability, error) {
	var cap C.struct_v4l2_cropcap
	cap._type = C.uint(bufType) // Application sets the type for which capabilities are requested.

	if err := send(fd, C.VIDIOC_CROPCAP, uintptr(unsafe.Pointer(&cap))); err != nil {
		return CropCapability{}, fmt.Errorf("crop capability: %w", err)
	}

	return *(*CropCapability)(unsafe.Pointer(&cap)), nil
}

// SetCropRect sets the current cropping rectangle for the device.
// It takes the file descriptor of the V4L2 device and a Rect defining the desired cropping area.
// The cropping rectangle is typically applied to video capture streams.
// The `_type` field in the underlying C struct is set to `V4L2_BUF_TYPE_VIDEO_CAPTURE` by default in this function.
//
// Returns an error if the VIDIOC_S_CROP ioctl call fails.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-crop.html#ioctl-vidioc-g-crop-vidioc-s-crop
func SetCropRect(fd uintptr, r Rect) error {
	var crop C.struct_v4l2_crop
	crop._type = C.uint(BufTypeVideoCapture) // Defaulting to video capture, adjust if other types are supported for cropping.
	crop.c = *(*C.struct_v4l2_rect)(unsafe.Pointer(&r))

	if err := send(fd, C.VIDIOC_S_CROP, uintptr(unsafe.Pointer(&crop))); err != nil {
		return fmt.Errorf("set crop: %w", err)
	}
	return nil
}

// String returns a human-readable string representation of the CropCapability struct.
// It includes the default cropping rectangle, bounds, and pixel aspect ratio.
func (c CropCapability) String() string {
	return fmt.Sprintf("default:{top=%d, left=%d, width=%d,height=%d};  bounds:{top=%d, left=%d, width=%d,height=%d}; pixel-aspect{%d:%d}",
		c.DefaultRect.Top,
		c.DefaultRect.Left,
		c.DefaultRect.Width,
		c.DefaultRect.Height,

		c.Bounds.Top,
		c.Bounds.Left,
		c.Bounds.Width,
		c.Bounds.Height,

		c.PixelAspect.Numerator,
		c.PixelAspect.Denominator,
	)
}
