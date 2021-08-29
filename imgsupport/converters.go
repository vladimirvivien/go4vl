package imgsupport

import (
	"bytes"
	"image"
	"image/jpeg"
)

func Yuyv2Jpeg(width, height int, frame []byte) ([]byte, error) {
	size := len(frame)
	ycbr := image.NewYCbCr(image.Rect(0, 0, width, height), image.YCbCrSubsampleRatio422)

	for i := 0; i < size; i += 4 {
		y1, u, y2, v := frame[i], frame[i+1], frame[i+2], frame[i+3]
		ycbr.Y[i]   = y1
		ycbr.Y[i+1] = y2
		ycbr.Cb[i]  = u
		ycbr.Cr[i]  = v
	}

	var jpgBuf bytes.Buffer
	if err := jpeg.Encode(&jpgBuf, ycbr, nil); err != nil {
		return nil, err
	}

	return jpgBuf.Bytes(), nil
}
