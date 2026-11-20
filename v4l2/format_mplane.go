package v4l2

import (
	"fmt"
	"unsafe"
)

// MaxPlanes is the maximum number of memory planes in a multi-planar V4L2 buffer (VIDEO_MAX_PLANES).
const MaxPlanes = cVIDEO_MAX_PLANES

// PlanePixFormat describes size and stride layout for a single plane in a multi-planar pixel format (v4l2_plane_pix_format).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/pixfmt-v4l2-mplane.html
type PlanePixFormat struct {
	SizeImage    uint32
	BytesPerLine uint32
}

// PixFormatMPlane describes a multi-planar pixel format (v4l2_pix_format_mplane).
// Used with BufTypeVideoCaptureMPlane and BufTypeVideoOutputMPlane.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/pixfmt-v4l2-mplane.html
type PixFormatMPlane struct {
	Width        uint32
	Height       uint32
	PixelFormat  FourCCType
	Field        FieldType
	Colorspace   ColorspaceType
	PlaneFmt     [MaxPlanes]PlanePixFormat
	NumPlanes    uint8
	Flags        uint8
	YcbcrEnc     uint8
	HSVEnc       uint8
	Quantization uint8
	XferFunc     uint8
}

func pixFormatMPlaneFromC(src cV4L2PixFormatMPlane) PixFormatMPlane {
	out := PixFormatMPlane{
		Width:        src.Width,
		Height:       src.Height,
		PixelFormat:  src.Pixelformat,
		Field:        src.Field,
		Colorspace:   src.Colorspace,
		NumPlanes:    src.Num_planes,
		Flags:        src.Flags,
		YcbcrEnc:     src.Ycbcr_enc,
		HSVEnc:       src.Ycbcr_enc,
		Quantization: src.Quantization,
		XferFunc:     src.Xfer_func,
	}
	for i := 0; i < MaxPlanes; i++ {
		out.PlaneFmt[i] = PlanePixFormat{
			SizeImage:    src.Plane_fmt[i].Sizeimage,
			BytesPerLine: src.Plane_fmt[i].Bytesperline,
		}
	}
	return out
}

func pixFormatMPlaneToC(src PixFormatMPlane) cV4L2PixFormatMPlane {
	ycbcrEnc := src.YcbcrEnc
	if ycbcrEnc == 0 {
		ycbcrEnc = src.HSVEnc
	}
	out := cV4L2PixFormatMPlane{
		Width:        src.Width,
		Height:       src.Height,
		Pixelformat:  src.PixelFormat,
		Field:        src.Field,
		Colorspace:   src.Colorspace,
		Num_planes:   src.NumPlanes,
		Flags:        src.Flags,
		Ycbcr_enc:    ycbcrEnc,
		Quantization: src.Quantization,
		Xfer_func:    src.XferFunc,
	}
	for i := 0; i < MaxPlanes; i++ {
		out.Plane_fmt[i] = cV4L2PlanePixFormat{
			Sizeimage:    src.PlaneFmt[i].SizeImage,
			Bytesperline: src.PlaneFmt[i].BytesPerLine,
		}
	}
	return out
}

// GetPixFormatMPlane retrieves multi-planar pixel format information for the specified buffer type (VIDIOC_G_FMT).
// bufType should be BufTypeVideoCaptureMPlane or BufTypeVideoOutputMPlane.
func GetPixFormatMPlane(fd uintptr, bufType BufType) (PixFormatMPlane, error) {
	if !IsMultiPlanar(bufType) {
		return PixFormatMPlane{}, fmt.Errorf("get mplane format: buffer type %d is not multi-planar", bufType)
	}
	var v4l2Format cV4L2Format
	v4l2Format.Type = bufType

	if err := send(fd, cVIDIOC_G_FMT, uintptr(unsafe.Pointer(&v4l2Format))); err != nil {
		return PixFormatMPlane{}, fmt.Errorf("get mplane format: %w", err)
	}

	cFmt := *(*cV4L2PixFormatMPlane)(unsafe.Pointer(&v4l2Format.Fmt[0]))
	return pixFormatMPlaneFromC(cFmt), nil
}

// SetPixFormatMPlane configures a multi-planar pixel format for the specified buffer type (VIDIOC_S_FMT)
// and returns the driver-negotiated format.
func SetPixFormatMPlane(fd uintptr, bufType BufType, pixFmt PixFormatMPlane) (PixFormatMPlane, error) {
	if !IsMultiPlanar(bufType) {
		return PixFormatMPlane{}, fmt.Errorf("set mplane format: buffer type %d is not multi-planar", bufType)
	}
	var v4l2Format cV4L2Format
	v4l2Format.Type = bufType
	*(*cV4L2PixFormatMPlane)(unsafe.Pointer(&v4l2Format.Fmt[0])) = pixFormatMPlaneToC(pixFmt)

	if err := send(fd, cVIDIOC_S_FMT, uintptr(unsafe.Pointer(&v4l2Format))); err != nil {
		return PixFormatMPlane{}, fmt.Errorf("set mplane format: %w", err)
	}

	cFmt := *(*cV4L2PixFormatMPlane)(unsafe.Pointer(&v4l2Format.Fmt[0]))
	return pixFormatMPlaneFromC(cFmt), nil
}

// TryPixFormatMPlane negotiates a multi-planar pixel format without altering device state (VIDIOC_TRY_FMT)
// and returns the driver-adjusted format.
func TryPixFormatMPlane(fd uintptr, bufType BufType, pixFmt PixFormatMPlane) (PixFormatMPlane, error) {
	if !IsMultiPlanar(bufType) {
		return PixFormatMPlane{}, fmt.Errorf("try mplane format: buffer type %d is not multi-planar", bufType)
	}
	var v4l2Format cV4L2Format
	v4l2Format.Type = bufType
	*(*cV4L2PixFormatMPlane)(unsafe.Pointer(&v4l2Format.Fmt[0])) = pixFormatMPlaneToC(pixFmt)

	if err := send(fd, cVIDIOC_TRY_FMT, uintptr(unsafe.Pointer(&v4l2Format))); err != nil {
		return PixFormatMPlane{}, fmt.Errorf("try mplane format: %w", err)
	}

	cFmt := *(*cV4L2PixFormatMPlane)(unsafe.Pointer(&v4l2Format.Fmt[0]))
	return pixFormatMPlaneFromC(cFmt), nil
}
