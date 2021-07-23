package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	frames <-chan []byte
	fps    uint32 = 30
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
		if _, err := w.Write(frame); err != nil {
			log.Printf("failed to write image: %s", err)
			return
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
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.StringVar(&port, "p", port, "webcam service port")
	flag.Parse()

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
	if err := device.SetPixFormat(v4l2.PixFormat{
		Width:       640,
		Height:      480,
		PixelFormat: v4l2.PixelFmtMJPEG,
		Field:       v4l2.FieldNone,
	}); err != nil {
		log.Fatalf("failed to set format: %s", err)
	}

	// Setup and start stream capture
	if err := device.StartStream(15); err != nil {
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
