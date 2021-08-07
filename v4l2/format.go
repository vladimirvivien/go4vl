package v4l2

import (
	"errors"
	"fmt"
	"unsafe"
)

// FourCCEncoding represents the four character encoding value
type FourCCEncoding = uint32

// Some Predefined pixel format definitions
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/pixfmt.html
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L518
var (
	PixFmtRGB24   = fourcc('R', 'G', 'B', '3') // V4L2_PIX_FMT_RGB24
	PixFmtGrey    = fourcc('G', 'R', 'E', 'Y') // V4L2_PIX_FMT_GREY
	PixelFmtYUYV  = fourcc('Y', 'U', 'Y', 'V') // V4L2_PIX_FMT_YUYV
	PixelFmtYYUV  = fourcc('Y', 'Y', 'U', 'V') // V4L2_PIX_FMT_YYUV
	PixelFmtYVYU  = fourcc('Y', 'V', 'Y', 'U') // V4L2_PIX_FMT_YVYU
	PixelFmtUYVY  = fourcc('U', 'Y', 'V', 'Y') // V4L2_PIX_FMT_UYVY
	PixelFmtVYUY  = fourcc('V', 'Y', 'U', 'Y') // V4L2_PIX_FMT_VYUY
	PixelFmtMJPEG = fourcc('M', 'J', 'P', 'G') // V4L2_PIX_FMT_MJPEG
	PixelFmtJPEG  = fourcc('J', 'P', 'E', 'G') // V4L2_PIX_FMT_JPEG
	PixelFmtMPEG  = fourcc('M', 'P', 'E', 'G') // V4L2_PIX_FMT_MPEG
	PixelFmtH264  = fourcc('H', '2', '6', '4') // V4L2_PIX_FMT_H264
	PixelFmtMPEG4 = fourcc('M', 'P', 'G', '4') // V4L2_PIX_FMT_MPEG4
)

// fourcc implements the four character code encoding found
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L81
// #define v4l2_fourcc(a, b, c, d)\
// 	 ((__u32)(a) | ((__u32)(b) << 8) | ((__u32)(c) << 16) | ((__u32)(d) << 24))
func fourcc(a, b, c, d uint32) FourCCEncoding {
	return (a | b<<8) | c<<16 | d<<24
}

// YcbcrEncoding (v4l2_ycbcr_encoding)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_ycbcr_encoding
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L300
type YcbcrEncoding = uint32

const (
	YcbcrEncDefault YcbcrEncoding = iota // V4L2_YCBCR_ENC_DEFAULT
	YcbcrEnc601                          // V4L2_YCBCR_ENC_601
	YcbcrEnc709                          // V4L2_YCBCR_ENC_709
	YcbcrEncXV601                        // V4L2_YCBCR_ENC_XV601
	YcbcrEncXV709                        // V4L2_YCBCR_ENC_XV709
)

// Quantization (v4l2_quantization)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_quantization#c.V4L.v4l2_quantization
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L372
type Quantization = uint32

const (
	QuantizationDefault   Quantization = iota // V4L2_QUANTIZATION_DEFAULT
	QuantizationFullRange                     // V4L2_QUANTIZATION_FULL_RANGE
	QuantizationLimRange                      // V4L2_QUANTIZATION_LIM_RANGE
)

// XferFunction (v4l2_xfer_func)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/colorspaces-defs.html?highlight=v4l2_xfer_func#c.V4L.v4l2_xfer_func
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L259
type XferFunction = uint32

const (
	XferFuncDefault   XferFunction = iota // V4L2_XFER_FUNC_DEFAULT     = 0
	XferFunc709                           // V4L2_XFER_FUNC_709         = 1,
	ferFuncSRGB                           // V4L2_XFER_FUNC_SRGB        = 2,
	XferFuncOpRGB                         // V4L2_XFER_FUNC_OPRGB       = 3,
	XferFuncSmpte240M                     // V4L2_XFER_FUNC_SMPTE240M   = 4,
	XferFuncNone                          // V4L2_XFER_FUNC_NONE        = 5,
	XferFuncDciP3                         // V4L2_XFER_FUNC_DCI_P3      = 6,
	XferFuncSmpte2084                     // V4L2_XFER_FUNC_SMPTE2084   = 7,
)

// Field (v4l2_field)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/field-order.html?highlight=v4l2_field#c.v4l2_field
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L88
type Field = uint32

const (
	FieldAny          Field = iota // V4L2_FIELD_ANY
	FieldNone                      // V4L2_FIELD_NONE
	FieldTop                       // V4L2_FIELD_TOP
	FieldBottom                    // V4L2_FIELD_BOTTOM
	FieldInterlaced                // V4L2_FIELD_INTERLACED
	FieldSeqTb                     // V4L2_FIELD_SEQ_TB
	FieldSeqBt                     // V4L2_FIELD_SEQ_BT
	FieldAlternate                 // V4L2_FIELD_ALTERNATE
	FieldInterlacedTb              // V4L2_FIELD_INTERLACED_TB
	FieldInterlacedBt              // V4L2_FIELD_INTERLACED_BT
)

// PixFormat (v4l2_pix_format)
// https://www.kernel.org/doc/html/v4.9/media/uapi/v4l/pixfmt-002.html?highlight=v4l2_pix_format
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L496
type PixFormat struct {
	Width        uint32
	Height       uint32
	PixelFormat  FourCCEncoding
	Field        Field
	BytesPerLine uint32
	SizeImage    uint32
	Colorspace   uint32
	Priv         uint32
	Flags        uint32
	YcbcrEnc     YcbcrEncoding
	Quantization Quantization
	XferFunc     XferFunction
}

// v4l2Format (v4l2_format)
// https://www.kernel.org/doc/html/v4.9/media/uapi/v4l/vidioc-g-fmt.html?highlight=v4l2_format
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2303
//
// field fmt is a union, thus it's constructed as an appropriately sized array:
//
// struct v4l2_format {
// 	__u32	 type;
// 	union {
// 		struct v4l2_pix_format		    pix;
// 		struct v4l2_pix_format_mplane	pix_mp;
// 		struct v4l2_window		        win;
// 		struct v4l2_vbi_format		    vbi;
// 		struct v4l2_sliced_vbi_format	sliced;
// 		struct v4l2_sdr_format	 	    sdr;
// 		struct v4l2_meta_format		    meta;
// 		__u8	raw_data[200];   /* user-defined */
// 	} fmt;
// };
type v4l2Format struct {
	StreamType uint32
	fmt        [200]byte
}

// getPixFormat returns the PixFormat by casting the pointer to the union type
func (f v4l2Format) getPixFormat() PixFormat {
	pixfmt := (*PixFormat)(unsafe.Pointer(&f.fmt[0]))
	return *pixfmt
}

// setPixFormat sets the PixFormat by casting the pointer to the fmt union and set its value
func (f v4l2Format) setPixFormat(newPix PixFormat) {
	*(*PixFormat)(unsafe.Pointer(&f.fmt[0])) = newPix
}

// GetPixFormat retrieves pixel information for the specified driver
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-fmt.html#ioctl-vidioc-g-fmt-vidioc-s-fmt-vidioc-try-fmt
func GetPixFormat(fd uintptr) (PixFormat, error) {
	format := v4l2Format{StreamType: BufTypeVideoCapture}
	if err := Send(fd, VidiocGetFormat, uintptr(unsafe.Pointer(&format))); err != nil {
		switch {
		case errors.Is(err, ErrorUnsupported):
			return PixFormat{}, fmt.Errorf("pix format: unsupported: %w", err)
		default:
			return PixFormat{}, fmt.Errorf("pix format failed: %w", err)
		}
	}

	return format.getPixFormat(), nil
}

// SetPixFormat sets the pixel format information for the specified driver
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-fmt.html#ioctl-vidioc-g-fmt-vidioc-s-fmt-vidioc-try-fmt
func SetPixFormat(fd uintptr, pixFmt PixFormat) error {
	format := v4l2Format{StreamType: BufTypeVideoCapture}
	format.setPixFormat(pixFmt)

	if err := Send(fd, VidiocSetFormat, uintptr(unsafe.Pointer(&format))); err != nil {
		switch {
		case errors.Is(err, ErrorUnsupported):
			return fmt.Errorf("pix format: unsupported operation: %w", err)
		default:
			return fmt.Errorf("pix format failed: %w", err)
		}
	}
	return nil
}
