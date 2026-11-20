package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	devName := flag.String("d", "/dev/video12", "V4L2 M2M device path (e.g., hardware scaler/ISP/codec)")
	bufCount := flag.Uint("b", 2, "buffers per queue")
	flag.Parse()

	sess, err := v4l2.OpenM2M(*devName)
	if err != nil {
		log.Fatalf("open M2M %s: %v", *devName, err)
	}
	defer sess.Close()

	cap := sess.Capability()
	fmt.Printf("Opened M2M device %s: driver=%s card=%s\n", *devName, cap.Driver, cap.Card)

	outType := v4l2.BufTypeVideoOutput
	capType := v4l2.BufTypeVideoCapture
	if cap.GetCapabilities()&v4l2.CapVideoMem2MemMPlane != 0 {
		outType = v4l2.BufTypeVideoOutputMPlane
		capType = v4l2.BufTypeVideoCaptureMPlane
	}

	outFmts, err := v4l2.GetAllFormatDescriptionsByType(sess.Fd(), outType)
	if err == nil {
		fmt.Printf("OUTPUT formats (%d):\n", len(outFmts))
		for _, f := range outFmts {
			fmt.Printf("  - %s\n", f.Description)
		}
	}
	capFmts, err := v4l2.GetAllFormatDescriptionsByType(sess.Fd(), capType)
	if err == nil {
		fmt.Printf("CAPTURE formats (%d):\n", len(capFmts))
		for _, f := range capFmts {
			fmt.Printf("  - %s\n", f.Description)
		}
	}

	if err := sess.SetupOutput(v4l2.M2MQueueConfig{
		BufType: outType,
		Memory:  v4l2.IOTypeMMAP,
		Count:   uint32(*bufCount),
	}); err != nil {
		log.Fatalf("setup output queue: %v", err)
	}
	if err := sess.SetupCapture(v4l2.M2MQueueConfig{
		BufType: capType,
		Memory:  v4l2.IOTypeMMAP,
		Count:   uint32(*bufCount),
	}); err != nil {
		log.Fatalf("setup capture queue: %v", err)
	}

	poller, err := v4l2.NewPoller()
	if err != nil {
		log.Fatalf("new poller: %v", err)
	}
	defer poller.Close()

	if err := poller.Add(sess.Fd(), v4l2.PollIn|v4l2.PollOut|v4l2.PollPri); err != nil {
		log.Fatalf("poller add: %v", err)
	}

	fmt.Printf("Initialized M2M session: output=%d buffers (%d planes), capture=%d buffers (%d planes)\n",
		sess.Output.Count, sess.Output.NumPlanes, sess.Capture.Count, sess.Capture.NumPlanes)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	results := make([]v4l2.PollResult, 4)
	n, err := poller.Wait(ctx, results, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("Poller wait completed: %v\n", err)
	} else {
		fmt.Printf("Poller ready FDs: %d\n", n)
	}
}
