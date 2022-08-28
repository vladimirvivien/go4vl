# Device Format Example

The examples in this directory highlights the support for V4L2's device format API. It shows how to query format information and set video format for a selected device.

```go

func main() {

    device, err := dev.Open(
		devName,
		dev.WithPixFormat(v4l2.PixFormat{Width: uint32(width), Height: uint32(height), PixelFormat: fmtEnc, Field: v4l2.FieldNone}),
		dev.WithFPS(15),
	)

    ...

    currFmt, err := device.GetPixFormat()
	if err != nil {
		log.Fatalf("unable to get format: %s", err)
	}
	log.Printf("Current format: %s", currFmt)

    ...

    // FPS
	fps, err := device.GetFrameRate()
	if err != nil {
		log.Fatalf("failed to get fps: %s", err)
	}
	log.Printf("current frame rate: %d fps", fps)
	// update fps
	if fps < 30 {
		if err := device.SetFrameRate(30); err != nil {
			log.Fatalf("failed to set frame rate: %s", err)
		}
	}

}
```