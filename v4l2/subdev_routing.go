package v4l2

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// Sub-device capability flags (V4L2_SUBDEV_CAP_*).
const (
	// SubdevCapReadOnly indicates the sub-device node is registered in read-only mode.
	SubdevCapReadOnly uint32 = cV4L2_SUBDEV_CAP_RO_SUBDEV

	// SubdevCapStreams indicates the sub-device supports routing and multiplexed streams.
	SubdevCapStreams uint32 = cV4L2_SUBDEV_CAP_STREAMS
)

// Sub-device route flags (V4L2_SUBDEV_ROUTE_FL_*).
const (
	// SubdevRouteFlagActive indicates the route is active when streaming is enabled.
	SubdevRouteFlagActive uint32 = cV4L2_SUBDEV_ROUTE_FL_ACTIVE
)

// Sub-device client capability flags (V4L2_SUBDEV_CLIENT_CAP_*).
const (
	// SubdevClientCapStreams signals to the kernel that the client is stream-aware.
	SubdevClientCapStreams uint64 = cV4L2_SUBDEV_CLIENT_CAP_STREAMS

	// SubdevClientCapIntervalUsesWhich signals that the client uses the which field in frame interval ioctls.
	SubdevClientCapIntervalUsesWhich uint64 = cV4L2_SUBDEV_CLIENT_CAP_INTERVAL_USES_WHICH
)

// SubdevCapability describes a sub-device's version and capabilities (v4l2_subdev_capability).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-subdev-querycap.html
type SubdevCapability struct {
	Version      uint32
	Capabilities uint32
}

// IsReadOnly reports whether the sub-device node is read-only.
func (c SubdevCapability) IsReadOnly() bool {
	return c.Capabilities&SubdevCapReadOnly != 0
}

// SupportsStreams reports whether the sub-device supports routing and multiplexed streams.
func (c SubdevCapability) SupportsStreams() bool {
	return c.Capabilities&SubdevCapStreams != 0
}

// SubdevRoute describes an internal stream route inside a sub-device (v4l2_subdev_route).
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-subdev-g-routing.html
type SubdevRoute struct {
	SinkPad      uint32
	SinkStream   uint32
	SourcePad    uint32
	SourceStream uint32
	Flags        uint32
}

// IsActive reports whether the route has SubdevRouteFlagActive set.
func (r SubdevRoute) IsActive() bool {
	return r.Flags&SubdevRouteFlagActive != 0
}

// GetSubdevCapability queries a sub-device's version and capabilities (VIDIOC_SUBDEV_QUERYCAP).
func GetSubdevCapability(fd uintptr) (SubdevCapability, error) {
	var cCap cV4L2SubdevCapability
	if err := send(fd, cVIDIOC_SUBDEV_QUERYCAP, unsafe.Pointer(&cCap)); err != nil {
		return SubdevCapability{}, fmt.Errorf("subdev querycap: %w", err)
	}
	return SubdevCapability{
		Version:      cCap.Version,
		Capabilities: cCap.Capabilities,
	}, nil
}

// GetSubdevClientCapability retrieves the active client capability bitmask (VIDIOC_SUBDEV_G_CLIENT_CAP).
func GetSubdevClientCapability(fd uintptr) (uint64, error) {
	var cCap cV4L2SubdevClientCapability
	if err := send(fd, cVIDIOC_SUBDEV_G_CLIENT_CAP, unsafe.Pointer(&cCap)); err != nil {
		return 0, fmt.Errorf("subdev get client cap: %w", err)
	}
	return cCap.Capabilities, nil
}

// SetSubdevClientCapability sets the client capability bitmask (VIDIOC_SUBDEV_S_CLIENT_CAP)
// and returns the effective capability bitmask acknowledged by the kernel.
func SetSubdevClientCapability(fd uintptr, caps uint64) (uint64, error) {
	cCap := cV4L2SubdevClientCapability{Capabilities: caps}
	if err := send(fd, cVIDIOC_SUBDEV_S_CLIENT_CAP, unsafe.Pointer(&cCap)); err != nil {
		return 0, fmt.Errorf("subdev set client cap: %w", err)
	}
	return cCap.Capabilities, nil
}

// GetSubdevRouting retrieves the sub-device's routing table (VIDIOC_SUBDEV_G_ROUTING).
// It first queries the route count and then fetches all routes in a second ioctl.
func GetSubdevRouting(fd uintptr, which SubdevFormatWhich) ([]SubdevRoute, error) {
	var cRouting cV4L2SubdevRouting
	cRouting.Which = which

	if err := send(fd, cVIDIOC_SUBDEV_G_ROUTING, unsafe.Pointer(&cRouting)); err != nil {
		// Linux returns ENOSPC when len_routes < num_routes and fills num_routes.
		if !errors.Is(err, sys.ENOSPC) {
			return nil, fmt.Errorf("subdev get routing (size query): %w", err)
		}
	}
	if cRouting.Num_routes == 0 {
		return nil, nil
	}

	raw := make([]cV4L2SubdevRoute, cRouting.Num_routes)
	var pinner runtime.Pinner
	pinner.Pin(&raw[0])
	defer pinner.Unpin()

	cRouting.Len_routes = cRouting.Num_routes
	cRouting.Routes = uint64(uintptr(unsafe.Pointer(&raw[0])))

	if err := send(fd, cVIDIOC_SUBDEV_G_ROUTING, unsafe.Pointer(&cRouting)); err != nil {
		return nil, fmt.Errorf("subdev get routing: %w", err)
	}

	count := int(cRouting.Num_routes)
	if count > len(raw) {
		count = len(raw)
	}
	out := make([]SubdevRoute, count)
	for i := 0; i < count; i++ {
		out[i] = SubdevRoute{
			SinkPad:      raw[i].Sink_pad,
			SinkStream:   raw[i].Sink_stream,
			SourcePad:    raw[i].Source_pad,
			SourceStream: raw[i].Source_stream,
			Flags:        raw[i].Flags,
		}
	}
	return out, nil
}

// SetSubdevRouting configures the sub-device's routing table (VIDIOC_SUBDEV_S_ROUTING)
// and returns the driver-negotiated routes.
func SetSubdevRouting(fd uintptr, which SubdevFormatWhich, routes []SubdevRoute) ([]SubdevRoute, error) {
	var cRouting cV4L2SubdevRouting
	cRouting.Which = which

	var raw []cV4L2SubdevRoute
	var pinner runtime.Pinner
	if len(routes) > 0 {
		raw = make([]cV4L2SubdevRoute, len(routes))
		for i, r := range routes {
			raw[i] = cV4L2SubdevRoute{
				Sink_pad:      r.SinkPad,
				Sink_stream:   r.SinkStream,
				Source_pad:    r.SourcePad,
				Source_stream: r.SourceStream,
				Flags:         r.Flags,
			}
		}
		pinner.Pin(&raw[0])
		defer pinner.Unpin()
		cRouting.Len_routes = uint32(len(raw))
		cRouting.Num_routes = uint32(len(raw))
		cRouting.Routes = uint64(uintptr(unsafe.Pointer(&raw[0])))
	}

	if err := send(fd, cVIDIOC_SUBDEV_S_ROUTING, unsafe.Pointer(&cRouting)); err != nil {
		return nil, fmt.Errorf("subdev set routing: %w", err)
	}

	count := int(cRouting.Num_routes)
	if count > len(raw) {
		count = len(raw)
	}
	out := make([]SubdevRoute, count)
	for i := 0; i < count; i++ {
		out[i] = SubdevRoute{
			SinkPad:      raw[i].Sink_pad,
			SinkStream:   raw[i].Sink_stream,
			SourcePad:    raw[i].Source_pad,
			SourceStream: raw[i].Source_stream,
			Flags:        raw[i].Flags,
		}
	}
	return out, nil
}
