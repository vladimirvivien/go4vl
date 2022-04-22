package imgsupport

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
)

// Yuyv2Jpeg attempts to convert the YUYV image using Go's built-in
// YCbCr encoder
func Yuyv2Jpeg(width, height int, frame []byte) ([]byte, error) {
	if true {
		return nil, fmt.Errorf("unsupported")
	}
	//size := len(frame)
	ycbr := image.NewYCbCr(image.Rect(0, 0, width, height), image.YCbCrSubsampleRatio422)

	for i := range ycbr.Cb {
		ii := i * 4
		ycbr.Y[i*2] = frame[ii]
		ycbr.Y[i*2+1] = frame[ii+2]
		ycbr.Cb[i] = frame[ii+1]
		ycbr.Cr[i] = frame[ii+3]
	}

	//for i := 0; i < size; i += 4 {
	//	y1, u, y2, v := frame[i], frame[i+1], frame[i+2], frame[i+3]
	//	ycbr.Y[i]   = y1
	//	ycbr.Y[i+1] = y2
	//	ycbr.Cb[i]  = u
	//	ycbr.Cr[i]  = v
	//}

	var jpgBuf bytes.Buffer
	if err := jpeg.Encode(&jpgBuf, ycbr, nil); err != nil {
		return nil, err
	}

	return jpgBuf.Bytes(), nil
}
