package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

func main() {
	devName := "/dev/video0"
	totalFrames := 3
	width := 640
	height := 480
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.IntVar(&totalFrames, "c", totalFrames, "number of frames to caputure")
	flag.IntVar(&width, "w", width, "picture width")
	flag.IntVar(&height, "h", height, "picture height")
	flag.Parse()

	// open device
	device, err := device.Open(devName)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()

	fps, err := device.GetFrameRate()
	if err != nil {
		log.Fatalf("failed to get framerate: %s", err)
	}

	// helper function to search for format descriptions
	findPreferredFmt := func(fmts []v4l2.FormatDescription, pixEncoding v4l2.FourCCType) *v4l2.FormatDescription {
		for _, desc := range fmts {
			if desc.PixelFormat == pixEncoding {
				return &desc
			}
		}
		return nil
	}

	// get supported format descriptions
	fmtDescs, err := device.GetFormatDescriptions()
	if err != nil {
		log.Fatal("failed to get format desc:", err)
	}

	// search for preferred formats
	preferredFmts := []v4l2.FourCCType{v4l2.PixelFmtMPEG, v4l2.PixelFmtMJPEG, v4l2.PixelFmtJPEG, v4l2.PixelFmtYUYV}
	var fmtDesc *v4l2.FormatDescription
	for _, preferredFmt := range preferredFmts {
		fmtDesc = findPreferredFmt(fmtDescs, preferredFmt)
		if fmtDesc != nil {
			break
		}
	}

	// no preferred pix fmt supported
	if fmtDesc == nil {
		log.Fatalf("device does not support any of %#v", preferredFmts)
	}
	log.Printf("Found preferred fmt: %s", fmtDesc)

	frameSizes, err := v4l2.GetFormatFrameSizes(device.Fd(), fmtDesc.PixelFormat)
	if err != nil {
		log.Fatalf("failed to get framesize info: %s", err)
	}

	// select size 640x480 for format
	var frmSize v4l2.FrameSizeEnum
	for _, size := range frameSizes {
		if size.Size.MinWidth == uint32(width) && size.Size.MinHeight == uint32(height) {
			frmSize = size
			break
		}
	}

	if frmSize.Size.MinWidth == 0 {
		log.Fatalf("Size %dx%d not supported for fmt: %s", width, height, fmtDesc)
	}

	log.Printf("Found preferred size: %#v", frmSize)

	// configure device with preferred fmt

	if err := device.SetPixFormat(v4l2.PixFormat{
		Width:       frmSize.Size.MinWidth,
		Height:      frmSize.Size.MinHeight,
		PixelFormat: fmtDesc.PixelFormat,
		Field:       v4l2.FieldNone,
	}); err != nil {
		log.Fatalf("failed to set format: %s", err)
	}

	pixFmt, err := device.GetPixFormat()
	if err != nil {
		log.Fatalf("failed to get format: %s", err)
	}
	log.Printf("Pixel format set to [%s]", pixFmt)

	// start stream
	ctx, cancel := context.WithCancel(context.TODO())
	if err := device.Start(ctx); err != nil {
		log.Fatalf("failed to stream: %s", err)
	}

	// process frames from capture channel
	count := 0
	log.Printf("Capturing %d frames (buffers: %d, %d fps)...", totalFrames, device.BufferCount(), fps)
	for frame := range device.GetOutput() {
		if count >= totalFrames {
			break
		}
		count++

		if len(frame) == 0 {
			log.Println("received frame size 0")
			continue
		}

		log.Printf("captured %d bytes", len(frame))
		img, fmtName, err := image.Decode(bytes.NewReader(frame))
		if err != nil {
			log.Printf("failed to decode jpeg: %s", err)
			continue
		}
		log.Printf("decoded image format: %s", fmtName)

		var imgBuf bytes.Buffer
		if err := jpeg.Encode(&imgBuf, img, nil); err != nil {
			log.Printf("failed to encode jpeg: %s", err)
			continue
		}

		fileName := fmt.Sprintf("capture_%d.jpg", count)
		file, err := os.Create(fileName)
		if err != nil {
			log.Printf("failed to create file %s: %s", fileName, err)
			continue
		}

		if _, err := file.Write(frame); err != nil {
			log.Printf("failed to write file %s: %s", fileName, err)
			file.Close()
			continue
		}
		log.Printf("Saved file: %s", fileName)
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s: %s", fileName, err)
		}
	}

	cancel() // stop capture
	if err := device.Stop(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done.")

}
