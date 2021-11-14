package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/vladimirvivien/go4vl/imgsupport"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	frames <-chan []byte
	fps    uint32 = 30
	pixfmt v4l2.FourCCType
)

// servePage reads templated HTML
func servePage(w http.ResponseWriter, r *http.Request) {
	// Start HTTP response
	w.Header().Add("Content-Type", "text/html")
	pd := map[string]string{
		"fps":        fmt.Sprintf("%d fps", fps),
		"streamPath": fmt.Sprintf("/stream?%d", time.Now().UnixNano()),
	}
	t, err := template.ParseFiles("webcam.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute and return the template
	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, pd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// start http service
func serveVideoStream(w http.ResponseWriter, req *http.Request) {
	// Start HTTP Response
	const boundaryName = "Yt08gcU534c0p4Jqj0p0"

	// send multi-part header
	w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundaryName))
	w.WriteHeader(http.StatusOK)

	for frame := range frames {
		// start boundary
		io.WriteString(w, fmt.Sprintf("--%s\n", boundaryName))
		io.WriteString(w, "Content-Type: image/jpeg\n")
		io.WriteString(w, fmt.Sprintf("Content-Length: %d\n\n", len(frame)))

		// write frame
		switch pixfmt {
		case v4l2.PixelFmtMJPEG:
			if _, err := w.Write(frame); err != nil {
				log.Printf("failed to write image: %s", err)
				return
			}
		case v4l2.PixelFmtYUYV:
			data, err := imgsupport.Yuyv2Jpeg(640, 480, frame)
			if err != nil {
				log.Printf("failed to convert yuyv to jpeg: %s", err)
				continue
			}
			if _, err := w.Write(data); err != nil {
				log.Printf("failed to write image: %s", err)
				return
			}
		}
		// close boundary
		if _, err := io.WriteString(w, "\n"); err != nil {
			log.Printf("failed to write bounday: %s", err)
			return
		}
	}
}

func main() {
	port := ":9090"
	devName := "/dev/video0"
	defaultDev, err := v4l2.Open(devName)
	skipDefault := false
	if err != nil {
		skipDefault = true
	}

	width := 640
	height := 480
	format := "yuyv"
	if !skipDefault {
		pix, err := defaultDev.GetPixFormat()
		if err == nil {
			width = int(pix.Width)
			height = int(pix.Height)
			switch pix.PixelFormat {
			case v4l2.PixelFmtMJPEG:
				format = "mjpeg"
			case v4l2.PixelFmtH264:
				format = "h264"
			default:
				format = "yuyv"
			}
		}
	}

	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.IntVar(&width, "w", width, "capture width")
	flag.IntVar(&height, "h", height, "capture height")
	flag.StringVar(&format, "f", format, "pixel format")
	flag.StringVar(&port, "p", port, "webcam service port")
	flag.Parse()

	// close device used for default info
	if err := defaultDev.Close(); err != nil {
		// default device failed to close
	}

	// open device and setup device
	device, err := v4l2.Open(devName)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()
	caps, err := device.GetCapability()
	if err != nil {
		log.Println("failed to get device capabilities:", err)
	}
	log.Printf("device [%s] opened\n", devName)
	log.Printf("device info: %s", caps.String())

	// set device format
	currFmt, err := device.GetPixFormat()
	if err != nil {
		log.Fatalf("unable to get format: %s", err)
	}
	log.Printf("Current format: %s", currFmt)
	if err := device.SetPixFormat(updateFormat(currFmt, format, width, height)); err != nil {
		log.Fatalf("failed to set format: %s", err)
	}
	currFmt, err = device.GetPixFormat()
	if err != nil {
		log.Fatalf("unable to get format: %s", err)
	}
	pixfmt = currFmt.PixelFormat
	log.Printf("Updated format: %s", currFmt)

	// Setup and start stream capture
	if err := device.StartStream(2); err != nil {
		log.Fatalf("unable to start stream: %s", err)
	}

	// start capture
	ctx, cancel := context.WithCancel(context.TODO())
	f, err := device.Capture(ctx, fps)
	if err != nil {
		log.Fatalf("stream capture: %s", err)
	}
	defer func() {
		cancel()
		device.Close()
	}()
	frames = f // make frames available.
	log.Println("device capture started, frames available")

	log.Printf("starting server on port %s", port)
	log.Println("use url path /webcam")

	// setup http service
	http.HandleFunc("/webcam", servePage)        // returns an html page
	http.HandleFunc("/stream", serveVideoStream) // returns video feed
	http.Handle("/", http.FileServer(http.Dir(".")))
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func updateFormat(pix v4l2.PixFormat, fmtStr string, w, h int) v4l2.PixFormat {
	pix.Width = uint32(w)
	pix.Height = uint32(h)

	switch strings.ToLower(fmtStr) {
	case "mjpeg", "jpeg":
		pix.PixelFormat = v4l2.PixelFmtMJPEG
	case "h264", "h.264":
		pix.PixelFormat = v4l2.PixelFmtH264
	case "yuyv":
		pix.PixelFormat = v4l2.PixelFmtYUYV
	}

	return pix
}