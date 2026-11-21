package v4l2

import (
	"errors"
	"fmt"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// SubdevFormatWhich selects whether a sub-device format or selection ioctl operates
// on the try state (per-file-handle staging) or the active hardware state.
type SubdevFormatWhich = uint32

const (
	// SubdevFormatTry validates and stores settings in per-FD staging state without touching hardware.
	SubdevFormatTry SubdevFormatWhich = cV4L2_SUBDEV_FORMAT_TRY

	// SubdevFormatActive reads or applies settings directly to active hardware state.
	SubdevFormatActive SubdevFormatWhich = cV4L2_SUBDEV_FORMAT_ACTIVE
)

// MBusFormatCode identifies a V4L2 media bus pixel format (MEDIA_BUS_FMT_*).
type MBusFormatCode = uint32

const (
	MBusFmtFixed          MBusFormatCode = cMEDIA_BUS_FMT_FIXED
	MBusFmtRGB565_2X8_LE  MBusFormatCode = cMEDIA_BUS_FMT_RGB565_2X8_LE
	MBusFmtRGB888_1X24    MBusFormatCode = cMEDIA_BUS_FMT_RGB888_1X24
	MBusFmtBGR888_1X24    MBusFormatCode = cMEDIA_BUS_FMT_BGR888_1X24
	MBusFmtY8_1X8         MBusFormatCode = cMEDIA_BUS_FMT_Y8_1X8
	MBusFmtY10_1X10       MBusFormatCode = cMEDIA_BUS_FMT_Y10_1X10
	MBusFmtY12_1X12       MBusFormatCode = cMEDIA_BUS_FMT_Y12_1X12
	MBusFmtUYVY8_1X16     MBusFormatCode = cMEDIA_BUS_FMT_UYVY8_1X16
	MBusFmtVYUY8_1X16     MBusFormatCode = cMEDIA_BUS_FMT_VYUY8_1X16
	MBusFmtYUYV8_1X16     MBusFormatCode = cMEDIA_BUS_FMT_YUYV8_1X16
	MBusFmtYVYU8_1X16     MBusFormatCode = cMEDIA_BUS_FMT_YVYU8_1X16
	MBusFmtUYVY8_2X8      MBusFormatCode = cMEDIA_BUS_FMT_UYVY8_2X8
	MBusFmtYUYV8_2X8      MBusFormatCode = cMEDIA_BUS_FMT_YUYV8_2X8
	MBusFmtSBGGR8_1X8     MBusFormatCode = cMEDIA_BUS_FMT_SBGGR8_1X8
	MBusFmtSGBRG8_1X8     MBusFormatCode = cMEDIA_BUS_FMT_SGBRG8_1X8
	MBusFmtSGRBG8_1X8     MBusFormatCode = cMEDIA_BUS_FMT_SGRBG8_1X8
	MBusFmtSRGGB8_1X8     MBusFormatCode = cMEDIA_BUS_FMT_SRGGB8_1X8
	MBusFmtSBGGR10_1X10   MBusFormatCode = cMEDIA_BUS_FMT_SBGGR10_1X10
	MBusFmtSGBRG10_1X10   MBusFormatCode = cMEDIA_BUS_FMT_SGBRG10_1X10
	MBusFmtSGRBG10_1X10   MBusFormatCode = cMEDIA_BUS_FMT_SGRBG10_1X10
	MBusFmtSRGGB10_1X10   MBusFormatCode = cMEDIA_BUS_FMT_SRGGB10_1X10
	MBusFmtSBGGR12_1X12   MBusFormatCode = cMEDIA_BUS_FMT_SBGGR12_1X12
	MBusFmtSGBRG12_1X12   MBusFormatCode = cMEDIA_BUS_FMT_SGBRG12_1X12
	MBusFmtSGRBG12_1X12   MBusFormatCode = cMEDIA_BUS_FMT_SGRBG12_1X12
	MBusFmtSRGGB12_1X12   MBusFormatCode = cMEDIA_BUS_FMT_SRGGB12_1X12
	MBusFmtSBGGR16_1X16   MBusFormatCode = cMEDIA_BUS_FMT_SBGGR16_1X16
	MBusFmtSGBRG16_1X16   MBusFormatCode = cMEDIA_BUS_FMT_SGBRG16_1X16
	MBusFmtSGRBG16_1X16   MBusFormatCode = cMEDIA_BUS_FMT_SGRBG16_1X16
	MBusFmtSRGGB16_1X16   MBusFormatCode = cMEDIA_BUS_FMT_SRGGB16_1X16
	MBusFmtJPEG_1X8       MBusFormatCode = cMEDIA_BUS_FMT_JPEG_1X8
)

// MBusFormatNames maps common MEDIA_BUS_FMT_* codes to readable names.
var MBusFormatNames = map[MBusFormatCode]string{
	MBusFmtFixed:         "FIXED",
	MBusFmtRGB565_2X8_LE: "RGB565_2X8_LE",
	MBusFmtRGB888_1X24:   "RGB888_1X24",
	MBusFmtBGR888_1X24:   "BGR888_1X24",
	MBusFmtY8_1X8:        "Y8_1X8",
	MBusFmtY10_1X10:      "Y10_1X10",
	MBusFmtY12_1X12:      "Y12_1X12",
	MBusFmtUYVY8_1X16:    "UYVY8_1X16",
	MBusFmtVYUY8_1X16:    "VYUY8_1X16",
	MBusFmtYUYV8_1X16:    "YUYV8_1X16",
	MBusFmtYVYU8_1X16:    "YVYU8_1X16",
	MBusFmtUYVY8_2X8:     "UYVY8_2X8",
	MBusFmtYUYV8_2X8:     "YUYV8_2X8",
	MBusFmtSBGGR8_1X8:    "SBGGR8_1X8",
	MBusFmtSGBRG8_1X8:    "SGBRG8_1X8",
	MBusFmtSGRBG8_1X8:    "SGRBG8_1X8",
	MBusFmtSRGGB8_1X8:    "SRGGB8_1X8",
	MBusFmtSBGGR10_1X10:  "SBGGR10_1X10",
	MBusFmtSGBRG10_1X10:  "SGBRG10_1X10",
	MBusFmtSGRBG10_1X10:  "SGRBG10_1X10",
	MBusFmtSRGGB10_1X10:  "SRGGB10_1X10",
	MBusFmtSBGGR12_1X12:  "SBGGR12_1X12",
	MBusFmtSGBRG12_1X12:  "SGBRG12_1X12",
	MBusFmtSGRBG12_1X12:  "SGRBG12_1X12",
	MBusFmtSRGGB12_1X12:  "SRGGB12_1X12",
	MBusFmtSBGGR16_1X16:  "SBGGR16_1X16",
	MBusFmtSGBRG16_1X16:  "SGBRG16_1X16",
	MBusFmtSGRBG16_1X16:  "SGRBG16_1X16",
	MBusFmtSRGGB16_1X16:  "SRGGB16_1X16",
	MBusFmtJPEG_1X8:      "JPEG_1X8",
}

// MBusFrameFormat represents a sub-device pad frame format (v4l2_mbus_framefmt).
type MBusFrameFormat struct {
	Width        uint32
	Height       uint32
	Code         MBusFormatCode
	Field        uint32
	Colorspace   uint32
	YcbcrEnc     uint16
	Quantization uint16
	XferFunc     uint16
	Flags        uint16
}

// SubdevFormat represents a pad format request/response (v4l2_subdev_format).
type SubdevFormat struct {
	Which  SubdevFormatWhich
	Pad    uint32
	Format MBusFrameFormat
	Stream uint32
}

// SubdevMBusCode represents an enumerated media bus code on a pad (v4l2_subdev_mbus_code_enum).
type SubdevMBusCode struct {
	Pad    uint32
	Index  uint32
	Code   MBusFormatCode
	Which  SubdevFormatWhich
	Flags  uint32
	Stream uint32
}

// SubdevFrameSize represents an enumerated frame size range on a pad (v4l2_subdev_frame_size_enum).
type SubdevFrameSize struct {
	Index     uint32
	Pad       uint32
	Code      MBusFormatCode
	MinWidth  uint32
	MaxWidth  uint32
	MinHeight uint32
	MaxHeight uint32
	Which     SubdevFormatWhich
	Stream    uint32
}

// SubdevFrameInterval represents a sub-device pad frame interval (v4l2_subdev_frame_interval).
type SubdevFrameInterval struct {
	Pad      uint32
	Interval Fract
	Stream   uint32
	Which    SubdevFormatWhich
}

// SubdevFrameIntervalEnum represents an enumerated frame interval on a pad (v4l2_subdev_frame_interval_enum).
type SubdevFrameIntervalEnum struct {
	Index    uint32
	Pad      uint32
	Code     MBusFormatCode
	Width    uint32
	Height   uint32
	Interval Fract
	Which    SubdevFormatWhich
	Stream   uint32
}

// SubdevSelection represents a sub-device pad selection rectangle (v4l2_subdev_selection).
type SubdevSelection struct {
	Which  SubdevFormatWhich
	Pad    uint32
	Target SelectionTarget
	Flags  SelectionFlag
	Rect   Rect
	Stream uint32
}

func fromCMBusFrameFmt(c cV4L2MBusFrameFmt) MBusFrameFormat {
	return MBusFrameFormat{
		Width:        c.Width,
		Height:       c.Height,
		Code:         c.Code,
		Field:        c.Field,
		Colorspace:   c.Colorspace,
		YcbcrEnc:     c.Ycbcr_enc,
		Quantization: c.Quantization,
		XferFunc:     c.Xfer_func,
		Flags:        c.Flags,
	}
}

func toCMBusFrameFmt(f MBusFrameFormat) cV4L2MBusFrameFmt {
	return cV4L2MBusFrameFmt{
		Width:        f.Width,
		Height:       f.Height,
		Code:         f.Code,
		Field:        f.Field,
		Colorspace:   f.Colorspace,
		Ycbcr_enc:    f.YcbcrEnc,
		Quantization: f.Quantization,
		Xfer_func:    f.XferFunc,
		Flags:        f.Flags,
	}
}

// GetSubdevFormat retrieves the media bus format for a sub-device pad via VIDIOC_SUBDEV_G_FMT.
func GetSubdevFormat(fd uintptr, pad uint32, which SubdevFormatWhich) (SubdevFormat, error) {
	cFmt := cV4L2SubdevFormat{
		Which: which,
		Pad:   pad,
	}
	if err := send(fd, cVIDIOC_SUBDEV_G_FMT, uintptr(unsafe.Pointer(&cFmt))); err != nil {
		return SubdevFormat{}, fmt.Errorf("subdev get format (pad %d): %w", pad, err)
	}
	return SubdevFormat{
		Which:  cFmt.Which,
		Pad:    cFmt.Pad,
		Format: fromCMBusFrameFmt(cFmt.Format),
		Stream: cFmt.Stream,
	}, nil
}

// SetSubdevFormat negotiates and sets the media bus format for a sub-device pad via
// VIDIOC_SUBDEV_S_FMT, returning the format adjusted by the driver.
func SetSubdevFormat(fd uintptr, format SubdevFormat) (SubdevFormat, error) {
	cFmt := cV4L2SubdevFormat{
		Which:  format.Which,
		Pad:    format.Pad,
		Format: toCMBusFrameFmt(format.Format),
		Stream: format.Stream,
	}
	if err := send(fd, cVIDIOC_SUBDEV_S_FMT, uintptr(unsafe.Pointer(&cFmt))); err != nil {
		return SubdevFormat{}, fmt.Errorf("subdev set format (pad %d): %w", format.Pad, err)
	}
	return SubdevFormat{
		Which:  cFmt.Which,
		Pad:    cFmt.Pad,
		Format: fromCMBusFrameFmt(cFmt.Format),
		Stream: cFmt.Stream,
	}, nil
}

// EnumSubdevMBusCodes enumerates all supported media bus format codes on a sub-device pad
// via VIDIOC_SUBDEV_ENUM_MBUS_CODE.
func EnumSubdevMBusCodes(fd uintptr, pad uint32, which SubdevFormatWhich) ([]SubdevMBusCode, error) {
	var codes []SubdevMBusCode
	for idx := uint32(0); ; idx++ {
		cEnum := cV4L2SubdevMBusCodeEnum{
			Pad:   pad,
			Index: idx,
			Which: which,
		}
		if err := send(fd, cVIDIOC_SUBDEV_ENUM_MBUS_CODE, uintptr(unsafe.Pointer(&cEnum))); err != nil {
			if errors.Is(err, ErrorBadArgument) || errors.Is(err, sys.EINVAL) {
				break
			}
			return codes, fmt.Errorf("subdev enum mbus code (pad %d, index %d): %w", pad, idx, err)
		}
		codes = append(codes, SubdevMBusCode{
			Pad:    cEnum.Pad,
			Index:  cEnum.Index,
			Code:   cEnum.Code,
			Which:  cEnum.Which,
			Flags:  cEnum.Flags,
			Stream: cEnum.Stream,
		})
	}
	return codes, nil
}

// EnumSubdevFrameSizes enumerates all supported frame size ranges for a pad and media bus code
// via VIDIOC_SUBDEV_ENUM_FRAME_SIZE.
func EnumSubdevFrameSizes(fd uintptr, pad uint32, code MBusFormatCode, which SubdevFormatWhich) ([]SubdevFrameSize, error) {
	var sizes []SubdevFrameSize
	for idx := uint32(0); ; idx++ {
		cEnum := cV4L2SubdevFrameSizeEnum{
			Index: idx,
			Pad:   pad,
			Code:  code,
			Which: which,
		}
		if err := send(fd, cVIDIOC_SUBDEV_ENUM_FRAME_SIZE, uintptr(unsafe.Pointer(&cEnum))); err != nil {
			if errors.Is(err, ErrorBadArgument) || errors.Is(err, sys.EINVAL) {
				break
			}
			return sizes, fmt.Errorf("subdev enum frame size (pad %d, index %d): %w", pad, idx, err)
		}
		sizes = append(sizes, SubdevFrameSize{
			Index:     cEnum.Index,
			Pad:       cEnum.Pad,
			Code:      cEnum.Code,
			MinWidth:  cEnum.Min_width,
			MaxWidth:  cEnum.Max_width,
			MinHeight: cEnum.Min_height,
			MaxHeight: cEnum.Max_height,
			Which:     cEnum.Which,
			Stream:    cEnum.Stream,
		})
	}
	return sizes, nil
}

// GetSubdevFrameInterval retrieves the active frame interval for a sub-device pad
// via VIDIOC_SUBDEV_G_FRAME_INTERVAL.
func GetSubdevFrameInterval(fd uintptr, pad uint32) (SubdevFrameInterval, error) {
	cIval := cV4L2SubdevFrameInterval{Pad: pad}
	if err := send(fd, cVIDIOC_SUBDEV_G_FRAME_INTERVAL, uintptr(unsafe.Pointer(&cIval))); err != nil {
		return SubdevFrameInterval{}, fmt.Errorf("subdev get frame interval (pad %d): %w", pad, err)
	}
	return SubdevFrameInterval{
		Pad:      cIval.Pad,
		Interval: Fract{Numerator: cIval.Interval.Numerator, Denominator: cIval.Interval.Denominator},
		Stream:   cIval.Stream,
		Which:    cIval.Which,
	}, nil
}

// SetSubdevFrameInterval configures the frame interval for a sub-device pad via
// VIDIOC_SUBDEV_S_FRAME_INTERVAL, returning the driver-adjusted interval.
func SetSubdevFrameInterval(fd uintptr, ival SubdevFrameInterval) (SubdevFrameInterval, error) {
	cIval := cV4L2SubdevFrameInterval{
		Pad:      ival.Pad,
		Interval: cV4L2Fract{Numerator: ival.Interval.Numerator, Denominator: ival.Interval.Denominator},
		Stream:   ival.Stream,
		Which:    ival.Which,
	}
	if err := send(fd, cVIDIOC_SUBDEV_S_FRAME_INTERVAL, uintptr(unsafe.Pointer(&cIval))); err != nil {
		return SubdevFrameInterval{}, fmt.Errorf("subdev set frame interval (pad %d): %w", ival.Pad, err)
	}
	return SubdevFrameInterval{
		Pad:      cIval.Pad,
		Interval: Fract{Numerator: cIval.Interval.Numerator, Denominator: cIval.Interval.Denominator},
		Stream:   cIval.Stream,
		Which:    cIval.Which,
	}, nil
}

// EnumSubdevFrameIntervals enumerates supported frame intervals for a pad, media bus code,
// and frame dimensions via VIDIOC_SUBDEV_ENUM_FRAME_INTERVAL.
func EnumSubdevFrameIntervals(fd uintptr, pad uint32, code MBusFormatCode, width, height uint32, which SubdevFormatWhich) ([]SubdevFrameIntervalEnum, error) {
	var intervals []SubdevFrameIntervalEnum
	for idx := uint32(0); ; idx++ {
		cEnum := cV4L2SubdevFrameIntervalEnum{
			Index:  idx,
			Pad:    pad,
			Code:   code,
			Width:  width,
			Height: height,
			Which:  which,
		}
		if err := send(fd, cVIDIOC_SUBDEV_ENUM_FRAME_INTERVAL, uintptr(unsafe.Pointer(&cEnum))); err != nil {
			if errors.Is(err, ErrorBadArgument) || errors.Is(err, sys.EINVAL) {
				break
			}
			return intervals, fmt.Errorf("subdev enum frame interval (pad %d, index %d): %w", pad, idx, err)
		}
		intervals = append(intervals, SubdevFrameIntervalEnum{
			Index:    cEnum.Index,
			Pad:      cEnum.Pad,
			Code:     cEnum.Code,
			Width:    cEnum.Width,
			Height:   cEnum.Height,
			Interval: Fract{Numerator: cEnum.Interval.Numerator, Denominator: cEnum.Interval.Denominator},
			Which:    cEnum.Which,
			Stream:   cEnum.Stream,
		})
	}
	return intervals, nil
}

// GetSubdevSelection queries a selection rectangle (e.g., SelTargetCrop, SelTargetCropBounds,
// SelTargetNativeSize) on a sub-device pad via VIDIOC_SUBDEV_G_SELECTION.
func GetSubdevSelection(fd uintptr, pad uint32, which SubdevFormatWhich, target SelectionTarget) (SubdevSelection, error) {
	cSel := cV4L2SubdevSelection{
		Which:  which,
		Pad:    pad,
		Target: target,
	}
	if err := send(fd, cVIDIOC_SUBDEV_G_SELECTION, uintptr(unsafe.Pointer(&cSel))); err != nil {
		return SubdevSelection{}, fmt.Errorf("subdev get selection (pad %d, target 0x%x): %w", pad, target, err)
	}
	return SubdevSelection{
		Which:  cSel.Which,
		Pad:    cSel.Pad,
		Target: cSel.Target,
		Flags:  cSel.Flags,
		Rect: Rect{
			Left:   cSel.R.Left,
			Top:    cSel.R.Top,
			Width:  cSel.R.Width,
			Height: cSel.R.Height,
		},
		Stream: cSel.Stream,
	}, nil
}

// SetSubdevSelection applies a selection rectangle on a sub-device pad via
// VIDIOC_SUBDEV_S_SELECTION, returning the driver-adjusted rectangle.
func SetSubdevSelection(fd uintptr, sel SubdevSelection) (SubdevSelection, error) {
	cSel := cV4L2SubdevSelection{
		Which:  sel.Which,
		Pad:    sel.Pad,
		Target: sel.Target,
		Flags:  sel.Flags,
		R: cV4L2Rect{
			Left:   sel.Rect.Left,
			Top:    sel.Rect.Top,
			Width:  sel.Rect.Width,
			Height: sel.Rect.Height,
		},
		Stream: sel.Stream,
	}
	if err := send(fd, cVIDIOC_SUBDEV_S_SELECTION, uintptr(unsafe.Pointer(&cSel))); err != nil {
		return SubdevSelection{}, fmt.Errorf("subdev set selection (pad %d, target 0x%x): %w", sel.Pad, sel.Target, err)
	}
	return SubdevSelection{
		Which:  cSel.Which,
		Pad:    cSel.Pad,
		Target: cSel.Target,
		Flags:  cSel.Flags,
		Rect: Rect{
			Left:   cSel.R.Left,
			Top:    cSel.R.Top,
			Width:  cSel.R.Width,
			Height: cSel.R.Height,
		},
		Stream: cSel.Stream,
	}, nil
}
