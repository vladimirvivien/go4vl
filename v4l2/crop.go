package v4l2

import (
	"fmt"
	"unsafe"
)

// CropCapability (v4l2_cropcap)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-cropcap.html#c.v4l2_cropcap
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1221
type CropCapability struct {
	StreamType  uint32
	Bounds      Rect
	DefaultRect Rect
	PixelAspect Fract
	_           [4]uint32
}

// GetCropCapability  retrieves cropping info for specified device
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-cropcap.html#ioctl-vidioc-cropcap
func GetCropCapability(fd uintptr, bufType BufType) (CropCapability, error) {
	var cap cV4L2CropCap
	cap.Type = bufType

	if err := send(fd, cVIDIOC_CROPCAP, uintptr(unsafe.Pointer(&cap))); err != nil {
		return CropCapability{}, fmt.Errorf("crop capability: %w", err)
	}

	return CropCapability{
		StreamType:  cap.Type,
		Bounds:      Rect(cap.Bounds),
		DefaultRect: Rect(cap.Defrect),
		PixelAspect: Fract(cap.Pixelaspect),
	}, nil
}

// GetCropRect retrieves the current cropping rectangle for the specified buffer type.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-crop.html
func GetCropRect(fd uintptr, bufType BufType) (Rect, error) {
	var crop cV4L2Crop
	crop.Type = bufType

	if err := send(fd, cVIDIOC_G_CROP, uintptr(unsafe.Pointer(&crop))); err != nil {
		return Rect{}, fmt.Errorf("get crop: %w", err)
	}
	return Rect(crop.C), nil
}

// SetCropRectByType sets the cropping rectangle for the specified buffer type.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-crop.html#ioctl-vidioc-g-crop-vidioc-s-crop
func SetCropRectByType(fd uintptr, bufType BufType, r Rect) error {
	var crop cV4L2Crop
	crop.Type = bufType
	crop.C = cV4L2Rect(r)

	if err := send(fd, cVIDIOC_S_CROP, uintptr(unsafe.Pointer(&crop))); err != nil {
		return fmt.Errorf("set crop: %w", err)
	}
	return nil
}

// SetCropRect sets the cropping dimension for specified device
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-crop.html#ioctl-vidioc-g-crop-vidioc-s-crop
func SetCropRect(fd uintptr, r Rect) error {
	return SetCropRectByType(fd, BufTypeVideoCapture, r)
}

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
