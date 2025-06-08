package v4l2

// Area defines a 2D area with a width and height.
// It corresponds to the `v4l2_area` struct in the Linux kernel.
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L424
type Area struct {
	// Width of the area.
	Width uint32
	// Height of the area.
	Height uint32
}

// Fract represents a fractional number, typically used for aspect ratios or frame rates.
// It corresponds to the `v4l2_fract` struct in the Linux kernel.
// See https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/vidioc-enumstd.html#c.v4l2_fract
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L419
type Fract struct {
	// Numerator is the numerator of the fraction.
	Numerator uint32
	// Denominator is the denominator of the fraction.
	Denominator uint32
}

// Rect defines a 2D rectangle with a position (Left, Top) and dimensions (Width, Height).
// It corresponds to the `v4l2_rect` struct in the Linux kernel.
// This is commonly used for defining cropping areas, selection rectangles, etc.
// See https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/dev-overlay.html?highlight=v4l2_rect#c.v4l2_rect
// See also https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L412
type Rect struct {
	// Left is the x-coordinate of the top-left corner of the rectangle.
	Left int32
	// Top is the y-coordinate of the top-left corner of the rectangle.
	Top int32
	// Width of the rectangle.
	Width uint32
	// Height of the rectangle.
	Height uint32
}
