package v4l2

import (
	"fmt"
	"unsafe"
)

// Metadata pixel format constants (V4L2_META_FMT_*).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/meta-formats.html
const (
	MetaFmtVSP1HGO         FourCCType = cV4L2_META_FMT_VSP1_HGO
	MetaFmtVSP1HGT         FourCCType = cV4L2_META_FMT_VSP1_HGT
	MetaFmtUVC             FourCCType = cV4L2_META_FMT_UVC
	MetaFmtD4XX            FourCCType = cV4L2_META_FMT_D4XX
	MetaFmtVivid           FourCCType = cV4L2_META_FMT_VIVID
	MetaFmtRKISP1Params    FourCCType = cV4L2_META_FMT_RK_ISP1_PARAMS
	MetaFmtRKISP1Stat3A    FourCCType = cV4L2_META_FMT_RK_ISP1_STAT_3A
	MetaFmtBCM2835ISPStats FourCCType = cV4L2_META_FMT_BCM2835_ISP_STATS
	MetaFmtRPiBECfg        FourCCType = cV4L2_META_FMT_RPI_BE_CFG
	MetaFmtRPiFECfg        FourCCType = cV4L2_META_FMT_RPI_FE_CFG
	MetaFmtRPiFEStats      FourCCType = cV4L2_META_FMT_RPI_FE_STATS
)

// MetaFormats maps known V4L2_META_FMT_* FourCC codes to readable descriptions.
var MetaFormats = map[FourCCType]string{
	MetaFmtVSP1HGO:         "R-Car VSP1 1-D Histogram",
	MetaFmtVSP1HGT:         "R-Car VSP1 2-D Histogram",
	MetaFmtUVC:             "UVC Payload Header Metadata",
	MetaFmtD4XX:            "Intel D4xx Payload Header Metadata",
	MetaFmtVivid:           "Vivid Metadata",
	MetaFmtRKISP1Params:    "Rockchip ISP1 3A Parameters",
	MetaFmtRKISP1Stat3A:    "Rockchip ISP1 3A Statistics",
	MetaFmtBCM2835ISPStats: "BCM2835 ISP Statistics",
	MetaFmtRPiBECfg:        "Raspberry Pi PiSP BE Config",
	MetaFmtRPiFECfg:        "Raspberry Pi PiSP FE Config",
	MetaFmtRPiFEStats:      "Raspberry Pi PiSP FE Stats",
}

// MetaFormat describes the data format of a V4L2 metadata stream (v4l2_meta_format).
// Used with BufTypeMetaCapture and BufTypeMetaOutput.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/dev-meta.html
type MetaFormat struct {
	DataFormat   FourCCType
	BufferSize   uint32
	Width        uint32
	Height       uint32
	BytesPerLine uint32
}

// String returns a human-readable summary of the metadata format.
func (m MetaFormat) String() string {
	fourcc := string([]byte{
		byte(m.DataFormat),
		byte(m.DataFormat >> 8),
		byte(m.DataFormat >> 16),
		byte(m.DataFormat >> 24),
	})
	if desc, ok := MetaFormats[m.DataFormat]; ok {
		return fmt.Sprintf("%s (%s) [buffersize=%d]", desc, fourcc, m.BufferSize)
	}
	return fmt.Sprintf("0x%08x (%s) [buffersize=%d]", m.DataFormat, fourcc, m.BufferSize)
}

// IsMetaBufType reports whether bufType is a metadata capture or output buffer type.
func IsMetaBufType(bufType BufType) bool {
	return bufType == BufTypeMetaCapture || bufType == BufTypeMetaOutput
}

func metaFormatFromC(src cV4L2MetaFormat) MetaFormat {
	return MetaFormat{
		DataFormat:   src.Dataformat,
		BufferSize:   src.Buffersize,
		Width:        src.Width,
		Height:       src.Height,
		BytesPerLine: src.Bytesperline,
	}
}

func metaFormatToC(src MetaFormat) cV4L2MetaFormat {
	return cV4L2MetaFormat{
		Dataformat:   src.DataFormat,
		Buffersize:   src.BufferSize,
		Width:        src.Width,
		Height:       src.Height,
		Bytesperline: src.BytesPerLine,
	}
}

// GetMetaFormat retrieves the active metadata format for a metadata capture or
// output device (VIDIOC_G_FMT with BufTypeMetaCapture or BufTypeMetaOutput).
func GetMetaFormat(fd uintptr, bufType BufType) (MetaFormat, error) {
	if !IsMetaBufType(bufType) {
		return MetaFormat{}, fmt.Errorf("get meta format: buffer type %d is not metadata", bufType)
	}
	var cFmt cV4L2Format
	cFmt.Type = bufType
	if err := send(fd, cVIDIOC_G_FMT, unsafe.Pointer(&cFmt)); err != nil {
		return MetaFormat{}, fmt.Errorf("get meta format: %w", err)
	}
	cMeta := *(*cV4L2MetaFormat)(unsafe.Pointer(&cFmt.Fmt[0]))
	return metaFormatFromC(cMeta), nil
}

// SetMetaFormat configures the metadata format on a metadata capture or output
// device (VIDIOC_S_FMT) and returns the driver-negotiated MetaFormat.
func SetMetaFormat(fd uintptr, bufType BufType, metaFmt MetaFormat) (MetaFormat, error) {
	if !IsMetaBufType(bufType) {
		return MetaFormat{}, fmt.Errorf("set meta format: buffer type %d is not metadata", bufType)
	}
	var cFmt cV4L2Format
	cFmt.Type = bufType
	*(*cV4L2MetaFormat)(unsafe.Pointer(&cFmt.Fmt[0])) = metaFormatToC(metaFmt)
	if err := send(fd, cVIDIOC_S_FMT, unsafe.Pointer(&cFmt)); err != nil {
		return MetaFormat{}, fmt.Errorf("set meta format: %w", err)
	}
	cMeta := *(*cV4L2MetaFormat)(unsafe.Pointer(&cFmt.Fmt[0]))
	return metaFormatFromC(cMeta), nil
}

// TryMetaFormat negotiates a metadata format without modifying device state
// (VIDIOC_TRY_FMT) and returns the driver-adjusted MetaFormat.
func TryMetaFormat(fd uintptr, bufType BufType, metaFmt MetaFormat) (MetaFormat, error) {
	if !IsMetaBufType(bufType) {
		return MetaFormat{}, fmt.Errorf("try meta format: buffer type %d is not metadata", bufType)
	}
	var cFmt cV4L2Format
	cFmt.Type = bufType
	*(*cV4L2MetaFormat)(unsafe.Pointer(&cFmt.Fmt[0])) = metaFormatToC(metaFmt)
	if err := send(fd, cVIDIOC_TRY_FMT, unsafe.Pointer(&cFmt)); err != nil {
		return MetaFormat{}, fmt.Errorf("try meta format: %w", err)
	}
	cMeta := *(*cV4L2MetaFormat)(unsafe.Pointer(&cFmt.Fmt[0]))
	return metaFormatFromC(cMeta), nil
}
