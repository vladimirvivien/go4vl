package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"

	pigo "github.com/esimov/pigo/core"
	"github.com/fogleman/gg"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	camera      *device.Device
	frames      <-chan []byte
	fps         uint32 = 30
	pixfmt      v4l2.FourCCType
	width       = 640
	height      = 480
	streamInfo  string
	faceEnabled bool
	faceFinder  *pigo.Pigo
)

type PageData struct {
	StreamInfo     string
	StreamPath     string
	ImgWidth       int
	ImgHeight      int
	ControlPath    string
	FaceDetectPath string
	FaceEnabled    bool
}

// servePage reads templated HTML
func servePage(w http.ResponseWriter, r *http.Request) {
	pd := PageData{
		StreamInfo:     streamInfo,
		StreamPath:     fmt.Sprintf("/stream?%d", time.Now().UnixNano()),
		ImgWidth:       width,
		ImgHeight:      height,
		ControlPath:    "/control",
		FaceDetectPath: "/face",
	}
	if faceFinder != nil {
		pd.FaceEnabled = true
	}

	// Start HTTP response
	w.Header().Add("Content-Type", "text/html")
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
	mimeWriter := multipart.NewWriter(w)
	w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", mimeWriter.Boundary()))
	partHeader := make(textproto.MIMEHeader)
	partHeader.Add("Content-Type", "image/jpeg")

	var frame []byte
	for frame = range frames {
		if len(frame) == 0 {
			log.Print("skipping empty frame")
			continue
		}

		partWriter, err := mimeWriter.CreatePart(partHeader)
		if err != nil {
			log.Printf("failed to create multi-part writer: %s", err)
			return
		}

		if faceEnabled {
			if err := runFaceDetect(partWriter, frame); err != nil {
				log.Printf("face detection failed: %s", err)
				continue
			}
		} else {
			if _, err := partWriter.Write(frame); err != nil {
				log.Printf("failed to write image: %s", err)
			}
		}

	}
}

type faceDetectRequest struct {
	Mode string
}

func faceDetectControl(w http.ResponseWriter, req *http.Request) {
	var face faceDetectRequest
	err := json.NewDecoder(req.Body).Decode(&face)
	if err != nil {
		log.Printf("failed to decode control: %s", err)
		return
	}
	log.Printf("Face mode = %s", face.Mode)
	switch face.Mode {
	case "true", "on", "enable":
		if faceFinder == nil {
			faceEnabled = false
			log.Println("face detection not enabled, re-run webcam with -face flag")
			return
		}
		faceEnabled = true
	case "off", "disabled":
		faceEnabled = false
	}
}

func initFaceDetect() error {
	model, err := os.ReadFile("./facefinder.model")
	if err != nil {
		return fmt.Errorf("failed to load face finder model: %s", err)
	}
	p := pigo.NewPigo()
	faceFinder, err = p.Unpack(model)
	if err != nil {
		faceFinder = nil
		return fmt.Errorf("failed to initialize face classifier: %s", err)
	}
	return nil
}

func runFaceDetect(w io.Writer, frame []byte) error {
	img, _, err := image.Decode(bytes.NewReader(frame))
	if err != nil {
		return err
	}

	src := img.(*image.YCbCr)
	bounds := img.Bounds()
	params := pigo.CascadeParams{
		MinSize:     100,
		MaxSize:     600,
		ShiftFactor: 0.15,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: src.Y,
			Rows:   bounds.Dy(),
			Cols:   bounds.Dx(),
			Dim:    bounds.Dx(),
		},
	}

	dets := faceFinder.RunCascade(params, 0.0)
	dets = faceFinder.ClusterDetections(dets, 0)

	drawer := gg.NewContext(bounds.Max.X, bounds.Max.Y)
	drawer.DrawImage(img, 0, 0)

	for _, det := range dets {
		if det.Q >= 5.0 {
			drawer.DrawRectangle(
				float64(det.Col-det.Scale/2),
				float64(det.Row-det.Scale/2),
				float64(det.Scale),
				float64(det.Scale),
			)

			drawer.SetLineWidth(3.0)
			drawer.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
			drawer.Stroke()
		}
	}

	return drawer.EncodeJPG(w, nil)
}

type controlRequest struct {
	Name  string
	Value string
}

func controlVideo(w http.ResponseWriter, req *http.Request) {
	var ctrl controlRequest
	err := json.NewDecoder(req.Body).Decode(&ctrl)
	if err != nil {
		log.Printf("failed to decode control: %s", err)
		return
	}

	val, err := strconv.Atoi(ctrl.Value)
	if err != nil {
		log.Printf("failed to set brightness: %s", err)
		return
	}

	switch ctrl.Name {
	case "brightness":
		if err := camera.SetControlBrightness(int32(val)); err != nil {
			log.Printf("failed to set brightness: %s", err)
			return
		}
	case "contrast":
		if err := camera.SetControlContrast(int32(val)); err != nil {
			log.Printf("failed to set contrast: %s", err)
			return
		}
	case "saturation":
		if err := camera.SetControlSaturation(int32(val)); err != nil {
			log.Printf("failed to set saturation: %s", err)
			return
		}
	}

	log.Printf("applied control %#v", ctrl)

}

func main() {
	port := ":9090"
	devName := "/dev/video0"
	frameRate := int(fps)
	buffSize := 4
	defaultDev, err := device.Open(devName)
	skipDefault := false
	face := false
	if err != nil {
		skipDefault = true
	}

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
	// close device used for default info
	if err := defaultDev.Close(); err != nil {
		log.Fatalf("failed to close default device: %s", err)
	}

	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.IntVar(&width, "w", width, "capture width")
	flag.IntVar(&height, "h", height, "capture height")
	flag.StringVar(&format, "f", format, "pixel format")
	flag.StringVar(&port, "p", port, "webcam service port")
	flag.IntVar(&frameRate, "r", frameRate, "frames per second (fps)")
	flag.IntVar(&buffSize, "b", buffSize, "device buffer size")
	flag.BoolVar(&face, "face", face, "turns on face detection mode")
	flag.Parse()

	// if face enabled, force fmt, buff size, and frame rate to low.
	if face {
		if err := initFaceDetect(); err != nil {
			log.Printf("failed to initialize face detection: %s", err)
		}
		format = "mjpeg"
		buffSize = 1
		frameRate = 5
	}

	// open camera and setup camera
	camera, err = device.Open(devName,
		device.WithIOType(v4l2.IOTypeMMAP),
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: getFormatType(format), Width: uint32(width), Height: uint32(height), Field: v4l2.FieldAny}),
		device.WithFPS(uint32(frameRate)),
		device.WithBufferSize(uint32(buffSize)),
	)

	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer camera.Close()

	caps := camera.Capability()
	log.Printf("device [%s] opened\n", devName)
	log.Printf("device info: %s", caps.String())

	// set device format
	currFmt, err := camera.GetPixFormat()
	if err != nil {
		log.Fatalf("unable to get format: %s", err)
	}
	log.Printf("Current format: %s", currFmt)
	pixfmt = currFmt.PixelFormat
	streamInfo = fmt.Sprintf("%s - %s [%dx%d] %d fps",
		caps.Card,
		v4l2.PixelFormats[currFmt.PixelFormat],
		currFmt.Width, currFmt.Height, frameRate,
	)

	// start capture
	ctx, cancel := context.WithCancel(context.TODO())
	if err := camera.Start(ctx); err != nil {
		log.Fatalf("stream capture: %s", err)
	}
	defer func() {
		cancel()
		camera.Close()
	}()

	// video stream
	frames = camera.GetOutput()

	log.Printf("device capture started (buffer size set %d)", camera.BufferCount())
	log.Printf("starting server on port %s", port)
	log.Println("use url path /webcam")

	// setup http service
	http.HandleFunc("/webcam", servePage)        // returns an html page
	http.HandleFunc("/stream", serveVideoStream) // returns video feed
	http.HandleFunc("/control", controlVideo)    // applies video controls
	http.HandleFunc("/face", faceDetectControl)  // controls face detection
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func getFormatType(fmtStr string) v4l2.FourCCType {
	switch strings.ToLower(fmtStr) {
	case "jpeg":
		return v4l2.PixelFmtJPEG
	case "mpeg":
		return v4l2.PixelFmtMPEG
	case "mjpeg":
		return v4l2.PixelFmtMJPEG
	case "h264", "h.264":
		return v4l2.PixelFmtH264
	case "yuyv":
		return v4l2.PixelFmtYUYV
	case "rgb":
		return v4l2.PixelFmtRGB24
	}
	return v4l2.PixelFmtMPEG
}
