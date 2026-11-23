package v4l2

import (
	"fmt"
	"runtime"
	"unsafe"
)

// EDIDBlockSize is the size in bytes of a single VESA EDID block.
const EDIDBlockSize = 128

// EDID represents an Extended Display Identification Data block payload (v4l2_edid).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-edid.html
type EDID struct {
	Pad        uint32
	StartBlock uint32
	Blocks     uint32
	Data       []byte
}

func getEDID(fd uintptr, req uintptr, pad, startBlock, blocks uint32) (EDID, error) {
	cEDID := cV4L2EDID{
		Pad:    pad,
		Block:  startBlock,
		Blocks: blocks,
	}

	var data []byte
	var pinner runtime.Pinner
	if blocks > 0 {
		data = make([]byte, int(blocks)*EDIDBlockSize)
		pinner.Pin(&data[0])
		defer pinner.Unpin()
		cEDID.Edid = &data[0]
	}

	if err := send(fd, req, unsafe.Pointer(&cEDID)); err != nil {
		return EDID{}, fmt.Errorf("get edid (pad=%d start=%d blocks=%d): %w", pad, startBlock, blocks, err)
	}

	outLen := int(cEDID.Blocks) * EDIDBlockSize
	if outLen > len(data) {
		outLen = len(data)
	}
	return EDID{
		Pad:        cEDID.Pad,
		StartBlock: cEDID.Block,
		Blocks:     cEDID.Blocks,
		Data:       data[:outLen],
	}, nil
}

func setEDID(fd uintptr, req uintptr, edid EDID) error {
	if edid.Blocks > 0 && len(edid.Data) < int(edid.Blocks)*EDIDBlockSize {
		return fmt.Errorf("set edid: data length %d smaller than %d blocks (%d bytes)",
			len(edid.Data), edid.Blocks, int(edid.Blocks)*EDIDBlockSize)
	}

	cEDID := cV4L2EDID{
		Pad:    edid.Pad,
		Block:  edid.StartBlock,
		Blocks: edid.Blocks,
	}

	var pinner runtime.Pinner
	if edid.Blocks > 0 {
		pinner.Pin(&edid.Data[0])
		defer pinner.Unpin()
		cEDID.Edid = &edid.Data[0]
	}

	if err := send(fd, req, unsafe.Pointer(&cEDID)); err != nil {
		return fmt.Errorf("set edid (pad=%d start=%d blocks=%d): %w", edid.Pad, edid.StartBlock, edid.Blocks, err)
	}
	return nil
}

// GetEDID reads EDID blocks from a V4L2 video device (VIDIOC_G_EDID).
// Pass blocks=0 to query the total number of available EDID blocks without reading data.
func GetEDID(fd uintptr, pad, startBlock, blocks uint32) (EDID, error) {
	return getEDID(fd, cVIDIOC_G_EDID, pad, startBlock, blocks)
}

// SetEDID writes EDID blocks to a V4L2 video device (VIDIOC_S_EDID).
// Pass Blocks=0 to clear the EDID on the specified pad.
func SetEDID(fd uintptr, edid EDID) error {
	return setEDID(fd, cVIDIOC_S_EDID, edid)
}

// GetSubdevEDID reads EDID blocks from a V4L2 sub-device (VIDIOC_SUBDEV_G_EDID).
func GetSubdevEDID(fd uintptr, pad, startBlock, blocks uint32) (EDID, error) {
	return getEDID(fd, cVIDIOC_SUBDEV_G_EDID, pad, startBlock, blocks)
}

// SetSubdevEDID writes EDID blocks to a V4L2 sub-device (VIDIOC_SUBDEV_S_EDID).
func SetSubdevEDID(fd uintptr, edid EDID) error {
	return setEDID(fd, cVIDIOC_SUBDEV_S_EDID, edid)
}
