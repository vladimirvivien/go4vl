package v4l2

// Area (v4l2_area)
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L424
type Area struct {
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

// Rect (v4l2_rect)
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/dev-overlay.html?highlight=v4l2_rect#c.v4l2_rect
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L412
type Rect struct {
	Left   int32
	Top    int32
	Width  uint32
	Height uint32
}
