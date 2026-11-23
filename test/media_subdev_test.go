//go:build integration

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func TestIntegration_MediaTopologyAndEnum(t *testing.T) {
	matches, _ := filepath.Glob("/dev/media*")
	if len(matches) == 0 {
		t.Skip("No /dev/media* devices available")
	}

	for _, path := range matches {
		fd, err := v4l2.OpenDevice(path, sys.O_RDWR|sys.O_CLOEXEC, 0)
		if err != nil {
			t.Logf("skip %s: %v", path, err)
			continue
		}

		info, err := v4l2.GetMediaDeviceInfo(fd)
		if err != nil {
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("GetMediaDeviceInfo(%s): %v", path, err)
		}
		topo, err := v4l2.GetMediaTopology(fd)
		if err != nil {
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("GetMediaTopology(%s): %v", path, err)
		}
		if len(topo.Entities) == 0 {
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("GetMediaTopology(%s) returned 0 entities", path)
		}

		entities, err := v4l2.EnumAllMediaEntities(fd)
		if err != nil {
			_ = v4l2.CloseDevice(fd)
			t.Fatalf("EnumAllMediaEntities(%s): %v", path, err)
		}
		if len(entities) > 0 && (entities[0].Pads > 0 || entities[0].Links > 0) {
			if _, _, err := v4l2.EnumMediaLinks(fd, entities[0].ID, entities[0].Pads, entities[0].Links); err != nil {
				_ = v4l2.CloseDevice(fd)
				t.Fatalf("EnumMediaLinks(%s, entity %d): %v", path, entities[0].ID, err)
			}
		}

		for _, intf := range topo.Interfaces {
			if intf.Major > 0 {
				if devPath, err := v4l2.ResolveDevnodePath(intf.Major, intf.Minor); err == nil {
					t.Logf("%s (%s): interface %d:%d -> %s", path, info.Driver, intf.Major, intf.Minor, devPath)
				}
			}
		}
		_ = v4l2.CloseDevice(fd)
	}
}

func TestIntegration_SubdevPadInspection(t *testing.T) {
	const subdevPath = "/dev/v4l-subdev0"
	if _, err := os.Stat(subdevPath); err != nil {
		t.Skipf("Sub-device %s not available: %v", subdevPath, err)
	}

	fd, err := v4l2.OpenDevice(subdevPath, sys.O_RDWR|sys.O_CLOEXEC, 0)
	if err != nil {
		t.Fatalf("OpenDevice(%s): %v", subdevPath, err)
	}
	defer v4l2.CloseDevice(fd)

	activeFmt, err := v4l2.GetSubdevFormat(fd, 0, v4l2.SubdevFormatActive)
	if err != nil {
		t.Fatalf("GetSubdevFormat(ACTIVE): %v", err)
	}
	if activeFmt.Format.Width == 0 || activeFmt.Format.Height == 0 {
		t.Fatalf("unexpected zero active format: %+v", activeFmt.Format)
	}

	// Negotiate a TRY format without mutating active state.
	tryFmt := activeFmt
	tryFmt.Which = v4l2.SubdevFormatTry
	if _, err := v4l2.SetSubdevFormat(fd, tryFmt); err != nil {
		t.Fatalf("SetSubdevFormat(TRY): %v", err)
	}

	codes, err := v4l2.EnumSubdevMBusCodes(fd, 0, v4l2.SubdevFormatActive)
	if err != nil {
		t.Fatalf("EnumSubdevMBusCodes: %v", err)
	}
	if len(codes) == 0 {
		t.Fatalf("expected at least 1 media bus code on %s", subdevPath)
	}
	sizes, err := v4l2.EnumSubdevFrameSizes(fd, 0, codes[0].Code, v4l2.SubdevFormatActive)
	if err != nil {
		t.Fatalf("EnumSubdevFrameSizes: %v", err)
	}
	t.Logf("%s active=%dx%d code=0x%04x, %d codes, %d sizes",
		subdevPath, activeFmt.Format.Width, activeFmt.Format.Height, activeFmt.Format.Code, len(codes), len(sizes))

	if sel, err := v4l2.GetSubdevSelection(fd, 0, v4l2.SubdevFormatActive, v4l2.SelTargetCropBounds); err == nil {
		t.Logf("%s CROP_BOUNDS: %+v", subdevPath, sel.Rect)
	}
}
