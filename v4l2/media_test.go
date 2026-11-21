package v4l2

import "testing"

func TestMediaTopologyHelpers(t *testing.T) {
	topo := MediaTopology{
		TopologyVersion: 1,
		Entities: []MediaV2Entity{
			{ID: 1, Name: "imx477 10-001a", Function: MediaEntFuncCamSensor},
			{ID: 2, Name: "unicam-image", Function: MediaEntFuncIOV4L},
		},
		Interfaces: []MediaV2Interface{
			{ID: 100, Type: MediaIntfTypeV4LSubdev, Major: 81, Minor: 16},
			{ID: 101, Type: MediaIntfTypeV4LVideo, Major: 81, Minor: 0},
		},
		Pads: []MediaV2Pad{
			{ID: 10, EntityID: 1, Index: 0, Flags: MediaPadFlagSource},
			{ID: 11, EntityID: 2, Index: 0, Flags: MediaPadFlagSink},
		},
		Links: []MediaV2Link{
			{ID: 20, SourceID: 10, SinkID: 11, Flags: MediaLinkFlagDataLink | MediaLinkFlagEnabled | MediaLinkFlagImmutable},
			{ID: 21, SourceID: 100, SinkID: 1, Flags: MediaLinkFlagInterfaceLink},
			{ID: 22, SourceID: 101, SinkID: 2, Flags: MediaLinkFlagInterfaceLink},
		},
	}

	sensor, ok := topo.FindEntityByName("imx477 10-001a")
	if !ok || sensor.ID != 1 {
		t.Fatalf("FindEntityByName failed: got (%+v, %v)", sensor, ok)
	}

	sensors := topo.FindEntitiesByFunction(MediaEntFuncCamSensor)
	if len(sensors) != 1 || sensors[0].ID != 1 {
		t.Fatalf("FindEntitiesByFunction: got %+v", sensors)
	}

	pads := topo.EntityPads(1)
	if len(pads) != 1 || !pads[0].IsSource() || pads[0].IsSink() {
		t.Fatalf("EntityPads(1): got %+v", pads)
	}

	intfs := topo.EntityInterfaces(1)
	if len(intfs) != 1 || intfs[0].Type != MediaIntfTypeV4LSubdev || intfs[0].Minor != 16 {
		t.Fatalf("EntityInterfaces(1): got %+v", intfs)
	}

	padLinks := topo.PadLinks(10)
	if len(padLinks) != 1 || !padLinks[0].IsEnabled() || !padLinks[0].IsImmutable() {
		t.Fatalf("PadLinks(10): got %+v", padLinks)
	}

	if len(topo.DataLinks()) != 1 {
		t.Fatalf("DataLinks count = %d, want 1", len(topo.DataLinks()))
	}
	if len(topo.InterfaceLinks()) != 2 {
		t.Fatalf("InterfaceLinks count = %d, want 2", len(topo.InterfaceLinks()))
	}
}
