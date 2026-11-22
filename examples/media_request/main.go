package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func main() {
	mediaPath := flag.String("media", "/dev/media0", "Media Controller device path supporting MEDIA_IOC_REQUEST_ALLOC")
	videoPath := flag.String("video", "", "optional paired stateless codec /dev/video* node (auto-discovered if omitted)")
	flag.Parse()

	mediaFD, err := v4l2.OpenDevice(*mediaPath, sys.O_RDWR|sys.O_CLOEXEC, 0)
	if err != nil {
		log.Fatalf("open %s: %v", *mediaPath, err)
	}
	defer v4l2.CloseDevice(mediaFD)

	info, err := v4l2.GetMediaDeviceInfo(mediaFD)
	if err != nil {
		log.Fatalf("get media info: %v", err)
	}
	fmt.Printf("Media device: %s (driver=%s model=%s)\n", *mediaPath, info.Driver, info.Model)

	req, err := v4l2.AllocMediaRequest(mediaFD)
	if err != nil {
		log.Fatalf("allocate media request: %v", err)
	}
	defer req.Close()
	fmt.Printf("Allocated Media Request FD: %d\n", req.Fd())

	if err := req.Reinit(); err != nil {
		log.Fatalf("reinit media request: %v", err)
	}
	fmt.Printf("Reinitialized Media Request FD %d for frame reuse\n", req.Fd())

	// Resolve paired V4L2 video interface from topology when -video is not specified.
	vPath := *videoPath
	if vPath == "" {
		if topo, err := v4l2.GetMediaTopology(mediaFD); err == nil {
			for _, intf := range topo.Interfaces {
				if intf.Type == v4l2.MediaIntfTypeV4LVideo {
					if resolved, err := v4l2.ResolveDevnodePath(intf.Major, intf.Minor); err == nil {
						vPath = resolved
						break
					}
				}
			}
		}
	}
	if vPath == "" {
		return
	}

	sess, err := v4l2.OpenM2M(vPath)
	if err != nil {
		log.Fatalf("open M2M %s: %v", vPath, err)
	}
	defer sess.Close()

	outType := v4l2.BufTypeVideoOutput
	if sess.Capability().GetCapabilities()&v4l2.CapVideoMem2MemMPlane != 0 {
		outType = v4l2.BufTypeVideoOutputMPlane
	}
	reqBufs, err := v4l2.RequestBuffersForType(sess.Fd(), outType, v4l2.IOTypeMMAP, 1)
	if err != nil {
		log.Fatalf("request output buffers on %s: %v", vPath, err)
	}
	defer v4l2.RequestBuffersForType(sess.Fd(), outType, v4l2.IOTypeMMAP, 0)

	supportsReq := reqBufs.Capabilities&v4l2.BufCapSupportsRequests != 0
	fmt.Printf("Paired M2M node %s: output buffers=%d supports_requests=%v\n",
		vPath, reqBufs.Count, supportsReq)

	// Demonstrate binding extended controls and a buffer descriptor to req.Fd().
	ctrls := v4l2.NewExtControls()
	ctrls.SetRequestFD(req.Fd())
	buf := v4l2.Buffer{
		Index:     0,
		Type:      outType,
		Memory:    v4l2.IOTypeMMAP,
		Flags:     v4l2.BufFlagRequestFD,
		RequestFD: req.Fd(),
	}
	fmt.Printf("Prepared request-bound ExtControls (which=0x%x request_fd=%d) and Buffer (flags=0x%x request_fd=%d)\n",
		ctrls.GetClass(), ctrls.GetRequestFD(), buf.Flags, buf.RequestFD)
}
