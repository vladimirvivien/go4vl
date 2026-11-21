package v4l2

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// Media Controller entity function constants (MEDIA_ENT_F_*).
const (
	MediaEntFuncUnknown                  uint32 = cMEDIA_ENT_F_UNKNOWN
	MediaEntFuncV4L2SubdevUnknown        uint32 = cMEDIA_ENT_F_V4L2_SUBDEV_UNKNOWN
	MediaEntFuncIOV4L                    uint32 = cMEDIA_ENT_F_IO_V4L
	MediaEntFuncCamSensor                uint32 = cMEDIA_ENT_F_CAM_SENSOR
	MediaEntFuncFlash                    uint32 = cMEDIA_ENT_F_FLASH
	MediaEntFuncLens                     uint32 = cMEDIA_ENT_F_LENS
	MediaEntFuncProcVideoComposer        uint32 = cMEDIA_ENT_F_PROC_VIDEO_COMPOSER
	MediaEntFuncProcVideoPixelFormatter  uint32 = cMEDIA_ENT_F_PROC_VIDEO_PIXEL_FORMATTER
	MediaEntFuncProcVideoPixelEncConv    uint32 = cMEDIA_ENT_F_PROC_VIDEO_PIXEL_ENC_CONV
	MediaEntFuncProcVideoLUT             uint32 = cMEDIA_ENT_F_PROC_VIDEO_LUT
	MediaEntFuncProcVideoScaler          uint32 = cMEDIA_ENT_F_PROC_VIDEO_SCALER
	MediaEntFuncProcVideoStatistics      uint32 = cMEDIA_ENT_F_PROC_VIDEO_STATISTICS
	MediaEntFuncProcVideoEncoder         uint32 = cMEDIA_ENT_F_PROC_VIDEO_ENCODER
	MediaEntFuncProcVideoDecoder         uint32 = cMEDIA_ENT_F_PROC_VIDEO_DECODER
	MediaEntFuncProcVideoISP             uint32 = cMEDIA_ENT_F_PROC_VIDEO_ISP
	MediaEntFuncVidMux                   uint32 = cMEDIA_ENT_F_VID_MUX
	MediaEntFuncVidIfBridge              uint32 = cMEDIA_ENT_F_VID_IF_BRIDGE
)

// MediaEntFunctionNames maps common MEDIA_ENT_F_* values to readable names.
var MediaEntFunctionNames = map[uint32]string{
	MediaEntFuncUnknown:                 "Unknown",
	MediaEntFuncV4L2SubdevUnknown:       "V4L2 Subdev",
	MediaEntFuncIOV4L:                   "V4L2 I/O",
	MediaEntFuncCamSensor:               "Camera Sensor",
	MediaEntFuncFlash:                   "Flash",
	MediaEntFuncLens:                    "Lens",
	MediaEntFuncProcVideoComposer:       "Video Composer",
	MediaEntFuncProcVideoPixelFormatter: "Video Pixel Formatter",
	MediaEntFuncProcVideoPixelEncConv:   "Video Pixel Encoding Converter",
	MediaEntFuncProcVideoLUT:            "Video LUT",
	MediaEntFuncProcVideoScaler:         "Video Scaler",
	MediaEntFuncProcVideoStatistics:     "Video Statistics",
	MediaEntFuncProcVideoEncoder:        "Video Encoder",
	MediaEntFuncProcVideoDecoder:        "Video Decoder",
	MediaEntFuncProcVideoISP:            "Video ISP",
	MediaEntFuncVidMux:                  "Video Mux",
	MediaEntFuncVidIfBridge:             "Video Interface Bridge",
}

// Media Controller entity ID and entity flags.
const (
	MediaEntIDFlagNext    uint32 = cMEDIA_ENT_ID_FLAG_NEXT
	MediaEntFlagDefault   uint32 = cMEDIA_ENT_FL_DEFAULT
	MediaEntFlagConnector uint32 = cMEDIA_ENT_FL_CONNECTOR
)

// Media Controller interface type constants (MEDIA_INTF_T_*).
const (
	MediaIntfTypeV4LVideo   uint32 = cMEDIA_INTF_T_V4L_VIDEO
	MediaIntfTypeV4LVBI     uint32 = cMEDIA_INTF_T_V4L_VBI
	MediaIntfTypeV4LRadio   uint32 = cMEDIA_INTF_T_V4L_RADIO
	MediaIntfTypeV4LSubdev  uint32 = cMEDIA_INTF_T_V4L_SUBDEV
	MediaIntfTypeV4LSWRadio uint32 = cMEDIA_INTF_T_V4L_SWRADIO
	MediaIntfTypeV4LTouch   uint32 = cMEDIA_INTF_T_V4L_TOUCH
)

// MediaIntfTypeNames maps MEDIA_INTF_T_* interface types to readable names.
var MediaIntfTypeNames = map[uint32]string{
	MediaIntfTypeV4LVideo:   "V4L2 Video",
	MediaIntfTypeV4LVBI:     "V4L2 VBI",
	MediaIntfTypeV4LRadio:   "V4L2 Radio",
	MediaIntfTypeV4LSubdev:  "V4L2 Subdev",
	MediaIntfTypeV4LSWRadio: "V4L2 SW Radio",
	MediaIntfTypeV4LTouch:   "V4L2 Touch",
}

// Media Controller pad flag constants (MEDIA_PAD_FL_*).
const (
	MediaPadFlagSink        uint32 = cMEDIA_PAD_FL_SINK
	MediaPadFlagSource      uint32 = cMEDIA_PAD_FL_SOURCE
	MediaPadFlagMustConnect uint32 = cMEDIA_PAD_FL_MUST_CONNECT
)

// Media Controller link flag constants (MEDIA_LNK_FL_*).
const (
	MediaLinkFlagEnabled       uint32 = cMEDIA_LNK_FL_ENABLED
	MediaLinkFlagImmutable     uint32 = cMEDIA_LNK_FL_IMMUTABLE
	MediaLinkFlagDynamic       uint32 = cMEDIA_LNK_FL_DYNAMIC
	MediaLinkFlagLinkTypeMask  uint32 = cMEDIA_LNK_FL_LINK_TYPE & 0xffffffff
	MediaLinkFlagDataLink      uint32 = cMEDIA_LNK_FL_DATA_LINK
	MediaLinkFlagInterfaceLink uint32 = cMEDIA_LNK_FL_INTERFACE_LINK
	MediaLinkFlagAncillaryLink uint32 = cMEDIA_LNK_FL_ANCILLARY_LINK
)

// MediaDeviceInfo (media_device_info)
// See https://www.kernel.org/doc/html/latest/userspace-api/media/mediactl/media-ioc-device-info.html
// See https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/media.h#L29
type MediaDeviceInfo struct {
	Driver           string
	Model            string
	Serial           string
	BusInfo          string
	MediaVersion     VersionInfo
	HardwareRevision uint32
	DriverVersion    VersionInfo
}

// MediaV2Entity represents an entity in the Media Controller v2 topology (media_v2_entity).
type MediaV2Entity struct {
	ID       uint32
	Name     string
	Function uint32
	Flags    uint32
}

// MediaV2Interface represents a userspace interface node in the v2 topology (media_v2_interface).
type MediaV2Interface struct {
	ID    uint32
	Type  uint32
	Flags uint32
	Major uint32
	Minor uint32
}

// MediaV2Pad represents an entity pad in the v2 topology (media_v2_pad).
type MediaV2Pad struct {
	ID       uint32
	EntityID uint32
	Flags    uint32
	Index    uint32
}

// IsSink reports whether the pad is a sink (input) pad.
func (p MediaV2Pad) IsSink() bool { return p.Flags&MediaPadFlagSink != 0 }

// IsSource reports whether the pad is a source (output) pad.
func (p MediaV2Pad) IsSource() bool { return p.Flags&MediaPadFlagSource != 0 }

// MediaV2Link represents a data, interface, or ancillary link in the v2 topology (media_v2_link).
type MediaV2Link struct {
	ID       uint32
	SourceID uint32
	SinkID   uint32
	Flags    uint32
}

// LinkType returns the masked link type (MediaLinkFlagDataLink, MediaLinkFlagInterfaceLink, or MediaLinkFlagAncillaryLink).
func (l MediaV2Link) LinkType() uint32 { return l.Flags & MediaLinkFlagLinkTypeMask }

// IsEnabled reports whether the link is enabled.
func (l MediaV2Link) IsEnabled() bool { return l.Flags&MediaLinkFlagEnabled != 0 }

// IsImmutable reports whether the link cannot be modified.
func (l MediaV2Link) IsImmutable() bool { return l.Flags&MediaLinkFlagImmutable != 0 }

// MediaTopology holds a complete snapshot of a Media Controller v2 graph (MEDIA_IOC_G_TOPOLOGY).
type MediaTopology struct {
	TopologyVersion uint64
	Entities        []MediaV2Entity
	Interfaces      []MediaV2Interface
	Pads            []MediaV2Pad
	Links           []MediaV2Link
}

// MediaEntityDesc represents a legacy v1 entity descriptor (media_entity_desc).
type MediaEntityDesc struct {
	ID       uint32
	Name     string
	Type     uint32
	Revision uint32
	Flags    uint32
	GroupID  uint32
	Pads     uint16
	Links    uint16
	Major    uint32
	Minor    uint32
}

// MediaPadDesc identifies a pad for MEDIA_IOC_ENUM_LINKS and MEDIA_IOC_SETUP_LINK (media_pad_desc).
type MediaPadDesc struct {
	Entity uint32
	Index  uint16
	Flags  uint32
}

// MediaLinkDesc describes a directed pad-to-pad link for MEDIA_IOC_SETUP_LINK and MEDIA_IOC_ENUM_LINKS.
type MediaLinkDesc struct {
	Source MediaPadDesc
	Sink   MediaPadDesc
	Flags  uint32
}

// GetMediaDeviceInfo retrieves media information for specified device, if supported.
func GetMediaDeviceInfo(fd uintptr) (MediaDeviceInfo, error) {
	var mdi cMediaDeviceInfo
	if err := send(fd, cMEDIA_IOC_DEVICE_INFO, uintptr(unsafe.Pointer(&mdi))); err != nil {
		return MediaDeviceInfo{}, fmt.Errorf("media device info: %w", err)
	}
	return MediaDeviceInfo{
		Driver:           cString(mdi.Driver[:]),
		Model:            cString(mdi.Model[:]),
		Serial:           cString(mdi.Serial[:]),
		BusInfo:          cString(mdi.Bus_info[:]),
		MediaVersion:     VersionInfo{value: mdi.Media_version},
		HardwareRevision: mdi.Hw_revision,
		DriverVersion:    VersionInfo{value: mdi.Driver_version},
	}, nil
}

// GetMediaTopology fetches an atomic snapshot of the Media Controller v2 graph
// (entities, interfaces, pads, and links) via MEDIA_IOC_G_TOPOLOGY.
func GetMediaTopology(fd uintptr) (MediaTopology, error) {
	const maxRetries = 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		var probe cMediaV2Topology
		if err := send(fd, cMEDIA_IOC_G_TOPOLOGY, uintptr(unsafe.Pointer(&probe))); err != nil {
			return MediaTopology{}, fmt.Errorf("media get topology (probe): %w", err)
		}

		cEntities := make([]cMediaV2Entity, probe.Num_entities)
		cInterfaces := make([]cMediaV2Interface, probe.Num_interfaces)
		cPads := make([]cMediaV2Pad, probe.Num_pads)
		cLinks := make([]cMediaV2Link, probe.Num_links)

		var pinner runtime.Pinner
		topo := cMediaV2Topology{
			Num_entities:   probe.Num_entities,
			Num_interfaces: probe.Num_interfaces,
			Num_pads:       probe.Num_pads,
			Num_links:      probe.Num_links,
		}
		if len(cEntities) > 0 {
			pinner.Pin(&cEntities[0])
			topo.Ptr_entities = uint64(uintptr(unsafe.Pointer(&cEntities[0])))
		}
		if len(cInterfaces) > 0 {
			pinner.Pin(&cInterfaces[0])
			topo.Ptr_interfaces = uint64(uintptr(unsafe.Pointer(&cInterfaces[0])))
		}
		if len(cPads) > 0 {
			pinner.Pin(&cPads[0])
			topo.Ptr_pads = uint64(uintptr(unsafe.Pointer(&cPads[0])))
		}
		if len(cLinks) > 0 {
			pinner.Pin(&cLinks[0])
			topo.Ptr_links = uint64(uintptr(unsafe.Pointer(&cLinks[0])))
		}

		err := send(fd, cMEDIA_IOC_G_TOPOLOGY, uintptr(unsafe.Pointer(&topo)))
		pinner.Unpin()
		runtime.KeepAlive(cEntities)
		runtime.KeepAlive(cInterfaces)
		runtime.KeepAlive(cPads)
		runtime.KeepAlive(cLinks)
		if err != nil {
			return MediaTopology{}, fmt.Errorf("media get topology: %w", err)
		}

		// Retry if the graph changed between probe and fill passes.
		if topo.Topology_version != probe.Topology_version {
			continue
		}

		res := MediaTopology{
			TopologyVersion: topo.Topology_version,
			Entities:        make([]MediaV2Entity, len(cEntities)),
			Interfaces:      make([]MediaV2Interface, len(cInterfaces)),
			Pads:            make([]MediaV2Pad, len(cPads)),
			Links:           make([]MediaV2Link, len(cLinks)),
		}
		for i, e := range cEntities {
			res.Entities[i] = MediaV2Entity{
				ID:       e.Id,
				Name:     cString(e.Name[:]),
				Function: e.Function,
				Flags:    e.Flags,
			}
		}
		for i, intf := range cInterfaces {
			res.Interfaces[i] = MediaV2Interface{
				ID:    intf.Id,
				Type:  intf.Type,
				Flags: intf.Flags,
				Major: intf.Devnode.Major,
				Minor: intf.Devnode.Minor,
			}
		}
		for i, p := range cPads {
			res.Pads[i] = MediaV2Pad{
				ID:       p.Id,
				EntityID: p.Entity_id,
				Flags:    p.Flags,
				Index:    p.Index,
			}
		}
		for i, l := range cLinks {
			res.Links[i] = MediaV2Link{
				ID:       l.Id,
				SourceID: l.Source_id,
				SinkID:   l.Sink_id,
				Flags:    l.Flags,
			}
		}
		return res, nil
	}

	return MediaTopology{}, fmt.Errorf("media get topology: topology version changed repeatedly")
}

// SetupMediaLink modifies the flags (e.g., MediaLinkFlagEnabled) of a mutable media link
// via MEDIA_IOC_SETUP_LINK.
func SetupMediaLink(fd uintptr, link MediaLinkDesc) error {
	cLink := cMediaLinkDesc{
		Source: cMediaPadDesc{
			Entity: link.Source.Entity,
			Index:  link.Source.Index,
			Flags:  link.Source.Flags,
		},
		Sink: cMediaPadDesc{
			Entity: link.Sink.Entity,
			Index:  link.Sink.Index,
			Flags:  link.Sink.Flags,
		},
		Flags: link.Flags,
	}
	if err := send(fd, cMEDIA_IOC_SETUP_LINK, uintptr(unsafe.Pointer(&cLink))); err != nil {
		return fmt.Errorf("media setup link: %w", err)
	}
	return nil
}

// EnumMediaEntity queries a single entity by ID via MEDIA_IOC_ENUM_ENTITIES.
// Set MediaEntIDFlagNext on id to query the next entity with ID > (id &^ MediaEntIDFlagNext).
func EnumMediaEntity(fd uintptr, id uint32) (MediaEntityDesc, error) {
	var desc cMediaEntityDesc
	desc.Id = id
	if err := send(fd, cMEDIA_IOC_ENUM_ENTITIES, uintptr(unsafe.Pointer(&desc))); err != nil {
		return MediaEntityDesc{}, fmt.Errorf("media enum entity 0x%x: %w", id, err)
	}
	return MediaEntityDesc{
		ID:       desc.Id,
		Name:     cString(desc.Name[:]),
		Type:     desc.Type,
		Revision: desc.Revision,
		Flags:    desc.Flags,
		GroupID:  desc.Group_id,
		Pads:     desc.Pads,
		Links:    desc.Links,
		Major:    desc.Dev[0],
		Minor:    desc.Dev[1],
	}, nil
}

// EnumAllMediaEntities enumerates all entities on a Media Controller device via MEDIA_IOC_ENUM_ENTITIES.
func EnumAllMediaEntities(fd uintptr) ([]MediaEntityDesc, error) {
	var entities []MediaEntityDesc
	var id uint32
	for {
		ent, err := EnumMediaEntity(fd, id|MediaEntIDFlagNext)
		if err != nil {
			if errors.Is(err, ErrorBadArgument) || errors.Is(err, sys.EINVAL) {
				break
			}
			return entities, err
		}
		entities = append(entities, ent)
		id = ent.ID
	}
	return entities, nil
}

// EnumMediaLinks enumerates the pads and outbound links of an entity via MEDIA_IOC_ENUM_LINKS.
func EnumMediaLinks(fd uintptr, entityID uint32, numPads, numLinks uint16) ([]MediaPadDesc, []MediaLinkDesc, error) {
	cPads := make([]cMediaPadDesc, numPads)
	cLinks := make([]cMediaLinkDesc, numLinks)

	var pinner runtime.Pinner
	enum := cMediaLinksEnum{Entity: entityID}
	if len(cPads) > 0 {
		pinner.Pin(&cPads[0])
		enum.Pads = &cPads[0]
	}
	if len(cLinks) > 0 {
		pinner.Pin(&cLinks[0])
		enum.Links = &cLinks[0]
	}

	err := send(fd, cMEDIA_IOC_ENUM_LINKS, uintptr(unsafe.Pointer(&enum)))
	pinner.Unpin()
	runtime.KeepAlive(cPads)
	runtime.KeepAlive(cLinks)
	if err != nil {
		return nil, nil, fmt.Errorf("media enum links (entity %d): %w", entityID, err)
	}

	pads := make([]MediaPadDesc, len(cPads))
	for i, p := range cPads {
		pads[i] = MediaPadDesc{Entity: p.Entity, Index: p.Index, Flags: p.Flags}
	}
	links := make([]MediaLinkDesc, len(cLinks))
	for i, l := range cLinks {
		links[i] = MediaLinkDesc{
			Source: MediaPadDesc{Entity: l.Source.Entity, Index: l.Source.Index, Flags: l.Source.Flags},
			Sink:   MediaPadDesc{Entity: l.Sink.Entity, Index: l.Sink.Index, Flags: l.Sink.Flags},
			Flags:  l.Flags,
		}
	}
	return pads, links, nil
}

// FindEntityByName returns the first entity in the topology matching name.
func (t MediaTopology) FindEntityByName(name string) (MediaV2Entity, bool) {
	for _, e := range t.Entities {
		if e.Name == name {
			return e, true
		}
	}
	return MediaV2Entity{}, false
}

// FindEntitiesByFunction returns all entities matching the given MEDIA_ENT_F_* function.
func (t MediaTopology) FindEntitiesByFunction(fn uint32) []MediaV2Entity {
	var out []MediaV2Entity
	for _, e := range t.Entities {
		if e.Function == fn {
			out = append(out, e)
		}
	}
	return out
}

// EntityPads returns all pads belonging to entityID.
func (t MediaTopology) EntityPads(entityID uint32) []MediaV2Pad {
	var out []MediaV2Pad
	for _, p := range t.Pads {
		if p.EntityID == entityID {
			out = append(out, p)
		}
	}
	return out
}

// EntityInterfaces returns all interfaces linked to entityID via interface links.
func (t MediaTopology) EntityInterfaces(entityID uint32) []MediaV2Interface {
	intfByID := make(map[uint32]MediaV2Interface, len(t.Interfaces))
	for _, intf := range t.Interfaces {
		intfByID[intf.ID] = intf
	}
	var out []MediaV2Interface
	for _, l := range t.Links {
		if l.LinkType() == MediaLinkFlagInterfaceLink && l.SinkID == entityID {
			if intf, ok := intfByID[l.SourceID]; ok {
				out = append(out, intf)
			}
		}
	}
	return out
}

// PadLinks returns all data links where padID is either the source or sink pad.
func (t MediaTopology) PadLinks(padID uint32) []MediaV2Link {
	var out []MediaV2Link
	for _, l := range t.Links {
		if l.LinkType() == MediaLinkFlagDataLink && (l.SourceID == padID || l.SinkID == padID) {
			out = append(out, l)
		}
	}
	return out
}

// DataLinks returns all pad-to-pad data links in the topology.
func (t MediaTopology) DataLinks() []MediaV2Link {
	var out []MediaV2Link
	for _, l := range t.Links {
		if l.LinkType() == MediaLinkFlagDataLink {
			out = append(out, l)
		}
	}
	return out
}

// InterfaceLinks returns all interface-to-entity links in the topology.
func (t MediaTopology) InterfaceLinks() []MediaV2Link {
	var out []MediaV2Link
	for _, l := range t.Links {
		if l.LinkType() == MediaLinkFlagInterfaceLink {
			out = append(out, l)
		}
	}
	return out
}

// ResolveDevnodePath resolves a character device's (major, minor) pair to its
// /dev/<name> path via /sys/dev/char/<major>:<minor>/uevent.
func ResolveDevnodePath(major, minor uint32) (string, error) {
	ueventPath := fmt.Sprintf("/sys/dev/char/%d:%d/uevent", major, minor)
	f, err := os.Open(ueventPath)
	if err != nil {
		return "", fmt.Errorf("resolve devnode %d:%d: %w", major, minor, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if devName, ok := strings.CutPrefix(line, "DEVNAME="); ok && devName != "" {
			return "/dev/" + devName, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("resolve devnode %d:%d: %w", major, minor, err)
	}
	return "", fmt.Errorf("resolve devnode %d:%d: DEVNAME not found in %s", major, minor, ueventPath)
}
