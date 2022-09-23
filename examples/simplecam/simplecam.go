package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	frames <-chan []byte
)

func imageServ(w http.ResponseWriter, req *http.Request) {
	mimeWriter := multipart.NewWriter(w)
	w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", mimeWriter.Boundary()))
	partHeader := make(textproto.MIMEHeader)
	partHeader.Add("Content-Type", "image/jpeg")

	var frame []byte
	for frame = range frames {
		partWriter, err := mimeWriter.CreatePart(partHeader)
		if err != nil {
			log.Printf("failed to create multi-part writer: %s", err)
			return
		}

		if _, err := partWriter.Write(frame); err != nil {
			log.Printf("failed to write image: %s", err)
		}
	}
}

func main() {
	port := ":9090"
	devName := "/dev/video0"
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.StringVar(&port, "p", port, "webcam service port")

	camera, err := device.Open(
		devName,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 640, Height: 480}),
	)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer camera.Close()

	if err := camera.Start(context.TODO()); err != nil {
		log.Fatalf("camera start: %s", err)
	}

	frames = camera.GetOutput()

	log.Printf("Serving images: [%s/stream]", port)
	http.HandleFunc("/stream", imageServ)
	log.Fatal(http.ListenAndServe(port, nil))
}
