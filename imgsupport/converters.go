package imgsupport

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
)

// Yuyv2Jpeg converts a YUYV (YUV 4:2:2) formatted frame to JPEG format.
//
// YUYV is a packed pixel format where two pixels share chroma (color) information:
//   - Each 4 bytes represent 2 pixels: [Y0][U][Y1][V]
//   - Y0, Y1: Luminance values for pixel 0 and pixel 1
//   - U, V: Shared chroma values for both pixels
//
// Parameters:
//   - width: Frame width in pixels (must be even)
//   - height: Frame height in pixels
//   - frame: Raw YUYV data (expected size: width * height * 2 bytes)
//
// Returns:
//   - []byte: JPEG-encoded image data
//   - error: Conversion error if any
//
// Note: This function is currently experimental and disabled (returns unsupported error).
// The conversion implementation needs further testing and optimization.
//
// Example:
//
//	// Convert 640x480 YUYV frame to JPEG
//	jpegData, err := Yuyv2Jpeg(640, 480, yuyvFrame)
//	if err != nil {
//	    // Handle error - currently always returns unsupported
//	    log.Printf("Conversion not supported: %v", err)
//	}
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
