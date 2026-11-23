package v4l2

import (
	"testing"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

func TestSubdevCapabilityAndRouting_MockIoctl(t *testing.T) {
	var clientCaps uint64
	storedRoutes := []cV4L2SubdevRoute{
		{Sink_pad: 0, Sink_stream: 0, Source_pad: 1, Source_stream: 0, Flags: SubdevRouteFlagActive},
		{Sink_pad: 0, Sink_stream: 1, Source_pad: 1, Source_stream: 1, Flags: 0},
	}

	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		switch req {
		case cVIDIOC_SUBDEV_QUERYCAP:
			c := (*cV4L2SubdevCapability)(arg)
			c.Version = (6 << 16) | (12 << 8)
			c.Capabilities = SubdevCapStreams
			return 0
		case cVIDIOC_SUBDEV_G_CLIENT_CAP:
			c := (*cV4L2SubdevClientCapability)(arg)
			c.Capabilities = clientCaps
			return 0
		case cVIDIOC_SUBDEV_S_CLIENT_CAP:
			c := (*cV4L2SubdevClientCapability)(arg)
			clientCaps = c.Capabilities & (SubdevClientCapStreams | SubdevClientCapIntervalUsesWhich)
			c.Capabilities = clientCaps
			return 0
		case cVIDIOC_SUBDEV_G_ROUTING:
			r := (*cV4L2SubdevRouting)(arg)
			r.Num_routes = uint32(len(storedRoutes))
			if r.Len_routes < r.Num_routes {
				return sys.ENOSPC
			}
			dst := unsafe.Slice((*cV4L2SubdevRoute)(*(*unsafe.Pointer)(unsafe.Pointer(&r.Routes))), r.Len_routes)
			copy(dst, storedRoutes)
			return 0
		case cVIDIOC_SUBDEV_S_ROUTING:
			r := (*cV4L2SubdevRouting)(arg)
			src := unsafe.Slice((*cV4L2SubdevRoute)(*(*unsafe.Pointer)(unsafe.Pointer(&r.Routes))), r.Num_routes)
			storedRoutes = append([]cV4L2SubdevRoute(nil), src...)
			return 0
		default:
			return sys.ENOTTY
		}
	})

	cap, err := GetSubdevCapability(3)
	if err != nil {
		t.Fatalf("GetSubdevCapability() = %v", err)
	}
	if !cap.SupportsStreams() || cap.IsReadOnly() {
		t.Fatalf("unexpected SubdevCapability: %+v", cap)
	}

	gotCaps, err := SetSubdevClientCapability(3, SubdevClientCapStreams|SubdevClientCapIntervalUsesWhich)
	if err != nil || gotCaps != (SubdevClientCapStreams|SubdevClientCapIntervalUsesWhich) {
		t.Fatalf("SetSubdevClientCapability() = 0x%x, %v", gotCaps, err)
	}
	readCaps, err := GetSubdevClientCapability(3)
	if err != nil || readCaps != gotCaps {
		t.Fatalf("GetSubdevClientCapability() = 0x%x, %v", readCaps, err)
	}

	routes, err := GetSubdevRouting(3, SubdevFormatActive)
	if err != nil {
		t.Fatalf("GetSubdevRouting() = %v", err)
	}
	if len(routes) != 2 || !routes[0].IsActive() || routes[1].IsActive() {
		t.Fatalf("GetSubdevRouting() = %+v", routes)
	}

	updated, err := SetSubdevRouting(3, SubdevFormatTry, []SubdevRoute{
		{SinkPad: 0, SinkStream: 0, SourcePad: 1, SourceStream: 0, Flags: SubdevRouteFlagActive},
	})
	if err != nil || len(updated) != 1 || updated[0].SourcePad != 1 {
		t.Fatalf("SetSubdevRouting() = %+v, %v", updated, err)
	}
}

func TestEDID_MockIoctl(t *testing.T) {
	stored := make([]byte, 2*EDIDBlockSize)
	// Standard 8-byte EDID header
	copy(stored[:8], []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00})
	stored[EDIDBlockSize] = 0x02 // CEA-861 extension tag

	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		e := (*cV4L2EDID)(arg)
		switch req {
		case cVIDIOC_G_EDID:
			totalBlocks := uint32(len(stored) / EDIDBlockSize)
			if e.Blocks == 0 {
				e.Blocks = totalBlocks
				return 0
			}
			if e.Blocks > totalBlocks {
				e.Blocks = totalBlocks
			}
			dst := unsafe.Slice(e.Edid, int(e.Blocks)*EDIDBlockSize)
			copy(dst, stored[int(e.Block)*EDIDBlockSize:])
			return 0
		case cVIDIOC_S_EDID:
			if e.Blocks == 0 {
				stored = nil
				return 0
			}
			src := unsafe.Slice(e.Edid, int(e.Blocks)*EDIDBlockSize)
			stored = append([]byte(nil), src...)
			return 0
		default:
			return sys.ENOTTY
		}
	})

	// Query block count (blocks=0)
	probe, err := GetEDID(3, 0, 0, 0)
	if err != nil || probe.Blocks != 2 || len(probe.Data) != 0 {
		t.Fatalf("GetEDID(probe) = %+v, %v", probe, err)
	}

	// Read 2 blocks
	got, err := GetEDID(3, 0, 0, 2)
	if err != nil || len(got.Data) != 2*EDIDBlockSize || got.Data[1] != 0xff || got.Data[EDIDBlockSize] != 0x02 {
		t.Fatalf("GetEDID(2) = blocks=%d len=%d err=%v", got.Blocks, len(got.Data), err)
	}

	// Short buffer validation
	if err := SetEDID(3, EDID{Pad: 0, Blocks: 1, Data: make([]byte, 10)}); err == nil {
		t.Fatal("expected error when EDID data is shorter than blocks*128")
	}

	// Subdev write + read round-trip
	newBlock := make([]byte, EDIDBlockSize)
	newBlock[0] = 0xAA
	if err := SetSubdevEDID(3, EDID{Pad: 0, Blocks: 1, Data: newBlock}); err != nil {
		t.Fatalf("SetSubdevEDID() = %v", err)
	}
	subGot, err := GetSubdevEDID(3, 0, 0, 1)
	if err != nil || len(subGot.Data) != EDIDBlockSize || subGot.Data[0] != 0xAA {
		t.Fatalf("GetSubdevEDID() = %+v, %v", subGot, err)
	}
}
