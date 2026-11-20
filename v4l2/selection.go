package v4l2

import (
	"fmt"
	"unsafe"
)

// SelectionTarget identifies the target rectangle for VIDIOC_G_SELECTION and VIDIOC_S_SELECTION.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/v4l2-selection-targets.html
type SelectionTarget = uint32

const (
	SelTargetCrop           SelectionTarget = cV4L2_SEL_TGT_CROP
	SelTargetCropDefault    SelectionTarget = cV4L2_SEL_TGT_CROP_DEFAULT
	SelTargetCropBounds     SelectionTarget = cV4L2_SEL_TGT_CROP_BOUNDS
	SelTargetNativeSize     SelectionTarget = cV4L2_SEL_TGT_NATIVE_SIZE
	SelTargetCompose        SelectionTarget = cV4L2_SEL_TGT_COMPOSE
	SelTargetComposeDefault SelectionTarget = cV4L2_SEL_TGT_COMPOSE_DEFAULT
	SelTargetComposeBounds  SelectionTarget = cV4L2_SEL_TGT_COMPOSE_BOUNDS
	SelTargetComposePadded  SelectionTarget = cV4L2_SEL_TGT_COMPOSE_PADDED
	SelTargetCropActive     SelectionTarget = cV4L2_SEL_TGT_CROP_ACTIVE
	SelTargetComposeActive  SelectionTarget = cV4L2_SEL_TGT_COMPOSE_ACTIVE
)

// SelectionFlag controls rounding and constraint behavior for VIDIOC_S_SELECTION.
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/v4l2-selection-flags.html
type SelectionFlag = uint32

const (
	SelFlagGE         SelectionFlag = cV4L2_SEL_FLAG_GE
	SelFlagLE         SelectionFlag = cV4L2_SEL_FLAG_LE
	SelFlagKeepConfig SelectionFlag = cV4L2_SEL_FLAG_KEEP_CONFIG
)

// Selection represents a V4L2 selection rectangle request or response (v4l2_selection).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-selection.html
type Selection struct {
	Type   BufType
	Target SelectionTarget
	Flags  SelectionFlag
	Rect   Rect
}

// GetSelection retrieves the selection rectangle for the specified buffer type and target (VIDIOC_G_SELECTION).
func GetSelection(fd uintptr, bufType BufType, target SelectionTarget) (Selection, error) {
	var sel cV4L2Selection
	sel.Type = bufType
	sel.Target = target

	if err := send(fd, cVIDIOC_G_SELECTION, uintptr(unsafe.Pointer(&sel))); err != nil {
		return Selection{}, fmt.Errorf("get selection: target 0x%x: %w", target, err)
	}

	return Selection{
		Type:   sel.Type,
		Target: sel.Target,
		Flags:  sel.Flags,
		Rect:   Rect(sel.R),
	}, nil
}

// SetSelection sets the selection rectangle (VIDIOC_S_SELECTION) and returns the driver-adjusted Selection.
func SetSelection(fd uintptr, sel Selection) (Selection, error) {
	cSel := cV4L2Selection{
		Type:   sel.Type,
		Target: sel.Target,
		Flags:  sel.Flags,
		R:      cV4L2Rect(sel.Rect),
	}

	if err := send(fd, cVIDIOC_S_SELECTION, uintptr(unsafe.Pointer(&cSel))); err != nil {
		return Selection{}, fmt.Errorf("set selection: target 0x%x: %w", sel.Target, err)
	}

	return Selection{
		Type:   cSel.Type,
		Target: cSel.Target,
		Flags:  cSel.Flags,
		Rect:   Rect(cSel.R),
	}, nil
}
