package v4l2

import (
	"fmt"
	"unsafe"
)

// Rect (v4l2_rect)
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/dev-overlay.html?highlight=v4l2_rect#c.v4l2_rect
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L412
type Rect struct {
	Left   int32
	Top    int32
	Width  uint32
	Height uint32
}

// Fract (v4l2_fract)
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/vidioc-enumstd.html#c.v4l2_fract
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L419
type Fract struct {
	Numerator   uint32
	Denominator uint32
}

// CropCapability (v4l2_cropcap)
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/vidioc-cropcap.html#c.v4l2_cropcap
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1221
type CropCapability struct {
	StreamType  uint32
	Bounds      Rect
	DefaultRect Rect
	PixelAspect Fract
}

// Crop (v4l2_crop)
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/vidioc-g-crop.html#c.v4l2_crop
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1228
type Crop struct {
	StreamType uint32
	Rect       Rect
}

// GetCropCapability  retrieves cropping info for specified device
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-cropcap.html#ioctl-vidioc-cropcap
func GetCropCapability(fd uintptr) (CropCapability, error) {
	cropCap := CropCapability{}
	cropCap.StreamType = BufTypeVideoCapture
	if err := Send(fd, VidiocCropCap, uintptr(unsafe.Pointer(&cropCap))); err != nil {
		return CropCapability{}, fmt.Errorf("crop capability: %w", err)
	}
	return cropCap, nil
}

// SetCropRect sets the cropping dimension for specified device
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-crop.html#ioctl-vidioc-g-crop-vidioc-s-crop
func SetCropRect(fd uintptr, r Rect) error {
	crop := Crop{Rect: r, StreamType: BufTypeVideoCapture}
	if err := Send(fd, VidiocSetCrop, uintptr(unsafe.Pointer(&crop))); err != nil {
		return fmt.Errorf("set crop: %w", err)
	}
	return nil
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
