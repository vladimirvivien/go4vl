package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func main() {
	devPath := flag.String("d", "/dev/video0", "V4L2 multi-planar device path")
	bufCount := flag.Uint("b", 4, "Number of buffers to request")
	exportPlane := flag.Bool("export", false, "Export plane 0 of each buffer as a DMA-BUF fd")
	flag.Parse()

	fd, err := v4l2.OpenDevice(*devPath, sys.O_RDWR|sys.O_NONBLOCK, 0)
	if err != nil {
		log.Fatalf("open %s: %v", *devPath, err)
	}
	defer v4l2.CloseDevice(fd)

	cap, err := v4l2.GetCapability(fd)
	if err != nil {
		log.Fatalf("querycap: %v", err)
	}
	fmt.Printf("Device: %s (%s) driver=%s\n", *devPath, cap.Card, cap.Driver)

	caps := cap.GetCapabilities()
	if caps&v4l2.CapVideoCaptureMPlane == 0 && caps&v4l2.CapVideoMem2MemMPlane == 0 {
		log.Fatalf("device %s does not advertise multi-planar capture capability (caps=0x%x)", *devPath, caps)
	}

	bufType := v4l2.BufTypeVideoCaptureMPlane

	descs, err := v4l2.GetAllFormatDescriptionsByType(fd, bufType)
	if err != nil && len(descs) == 0 {
		log.Fatalf("enum mplane formats: %v", err)
	}
	fmt.Printf("Supported MPLANE capture formats (%d):\n", len(descs))
	for _, d := range descs {
		fmt.Printf("  [%d] %s (0x%08x)\n", d.Index, d.Description, d.PixelFormat)
	}

	curFmt, err := v4l2.GetPixFormatMPlane(fd, bufType)
	if err != nil {
		log.Fatalf("get mplane format: %v", err)
	}
	fmt.Printf("Current MPLANE format: %dx%d fourcc=0x%08x planes=%d\n",
		curFmt.Width, curFmt.Height, curFmt.PixelFormat, curFmt.NumPlanes)
	for p := uint8(0); p < curFmt.NumPlanes; p++ {
		fmt.Printf("  plane %d: sizeimage=%d bytesperline=%d\n",
			p, curFmt.PlaneFmt[p].SizeImage, curFmt.PlaneFmt[p].BytesPerLine)
	}

	tried, err := v4l2.TryPixFormatMPlane(fd, bufType, curFmt)
	if err != nil {
		log.Fatalf("try mplane format: %v", err)
	}
	fmt.Printf("TryPixFormatMPlane negotiated: %dx%d planes=%d\n", tried.Width, tried.Height, tried.NumPlanes)

	req, err := v4l2.RequestBuffersForType(fd, bufType, v4l2.IOTypeMMAP, uint32(*bufCount))
	if err != nil {
		log.Fatalf("request buffers: %v", err)
	}
	defer v4l2.RequestBuffersForType(fd, bufType, v4l2.IOTypeMMAP, 0)
	fmt.Printf("Allocated %d multi-planar MMAP buffers (caps=0x%x)\n", req.Count, req.Capabilities)

	numPlanes := uint32(curFmt.NumPlanes)
	if numPlanes == 0 {
		numPlanes = 1
	}

	for i := uint32(0); i < req.Count; i++ {
		buf, err := v4l2.QueryBuffer(fd, bufType, v4l2.IOTypeMMAP, i, numPlanes)
		if err != nil {
			log.Fatalf("query buffer %d: %v", i, err)
		}
		fmt.Printf("Buffer %d: planes=%d flags=0x%x\n", buf.Index, len(buf.Planes), buf.Flags)
		for pIdx, pl := range buf.Planes {
			fmt.Printf("  plane %d: length=%d offset=0x%x data_offset=%d\n",
				pIdx, pl.Length, pl.Info.MemOffset, pl.DataOffset)
			mapped, err := v4l2.MapMemoryBuffer(fd, int64(pl.Info.MemOffset), int(pl.Length))
			if err != nil {
				log.Fatalf("mmap buffer %d plane %d: %v", i, pIdx, err)
			}
			_ = v4l2.UnmapMemoryBuffer(mapped)
		}

		if *exportPlane {
			dmaFD, err := v4l2.ExportBufferPlane(fd, bufType, i, 0, sys.O_CLOEXEC)
			if err != nil {
				log.Fatalf("export buffer %d plane 0: %v", i, err)
			}
			fmt.Printf("  exported plane 0 as DMA-BUF fd=%d\n", dmaFD)
			_ = sys.Close(int(dmaFD))
		}
	}
}
