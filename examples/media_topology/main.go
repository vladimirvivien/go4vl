package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func main() {
	mediaPath := flag.String("d", "/dev/media0", "Media Controller device path")
	flag.Parse()

	fd, err := v4l2.OpenDevice(*mediaPath, sys.O_RDWR|sys.O_CLOEXEC, 0)
	if err != nil {
		log.Fatalf("open %s: %v", *mediaPath, err)
	}
	defer v4l2.CloseDevice(fd)

	info, err := v4l2.GetMediaDeviceInfo(fd)
	if err != nil {
		log.Fatalf("get media info: %v", err)
	}
	fmt.Printf("Media device: %s\n", *mediaPath)
	fmt.Printf("  Driver:   %s\n", info.Driver)
	fmt.Printf("  Model:    %s\n", info.Model)
	fmt.Printf("  BusInfo:  %s\n", info.BusInfo)

	topo, err := v4l2.GetMediaTopology(fd)
	if err != nil {
		log.Fatalf("get topology: %v", err)
	}
	fmt.Printf("Topology v%d: %d entities, %d interfaces, %d pads, %d links\n\n",
		topo.TopologyVersion, len(topo.Entities), len(topo.Interfaces), len(topo.Pads), len(topo.Links))

	padByID := make(map[uint32]v4l2.MediaV2Pad, len(topo.Pads))
	for _, p := range topo.Pads {
		padByID[p.ID] = p
	}
	entityByID := make(map[uint32]v4l2.MediaV2Entity, len(topo.Entities))
	for _, e := range topo.Entities {
		entityByID[e.ID] = e
	}

	for _, ent := range topo.Entities {
		fnName := v4l2.MediaEntFunctionNames[ent.Function]
		if fnName == "" {
			fnName = fmt.Sprintf("0x%x", ent.Function)
		}
		fmt.Printf("- Entity %d: %s (function=%s)\n", ent.ID, ent.Name, fnName)

		for _, intf := range topo.EntityInterfaces(ent.ID) {
			devPath, err := v4l2.ResolveDevnodePath(intf.Major, intf.Minor)
			if err != nil {
				devPath = fmt.Sprintf("char %d:%d", intf.Major, intf.Minor)
			}
			intfName := v4l2.MediaIntfTypeNames[intf.Type]
			if intfName == "" {
				intfName = fmt.Sprintf("0x%x", intf.Type)
			}
			fmt.Printf("    interface: %s (%s)\n", devPath, intfName)
		}

		for _, pad := range topo.EntityPads(ent.ID) {
			dir := "Pad"
			if pad.IsSource() {
				dir = "Source"
			} else if pad.IsSink() {
				dir = "Sink"
			}
			fmt.Printf("    pad%d [%s] (id=%d)\n", pad.Index, dir, pad.ID)
			for _, link := range topo.PadLinks(pad.ID) {
				flags := formatLinkFlags(link)
				if link.SourceID == pad.ID {
					sinkPad := padByID[link.SinkID]
					sinkEnt := entityByID[sinkPad.EntityID]
					fmt.Printf("      -> %q:%d [%s]\n", sinkEnt.Name, sinkPad.Index, flags)
				} else {
					srcPad := padByID[link.SourceID]
					srcEnt := entityByID[srcPad.EntityID]
					fmt.Printf("      <- %q:%d [%s]\n", srcEnt.Name, srcPad.Index, flags)
				}
			}
		}
	}
}

func formatLinkFlags(l v4l2.MediaV2Link) string {
	var parts []string
	if l.IsEnabled() {
		parts = append(parts, "ENABLED")
	}
	if l.IsImmutable() {
		parts = append(parts, "IMMUTABLE")
	}
	if l.Flags&v4l2.MediaLinkFlagDynamic != 0 {
		parts = append(parts, "DYNAMIC")
	}
	if len(parts) == 0 {
		return "DISABLED"
	}
	return strings.Join(parts, ",")
}
