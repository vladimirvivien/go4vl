# Capture example

In this capture example, the source code uses a more advanced approach (compared to [capture0/capture0.go](../capture0/capture0.go)) where it 
leverages the go4vl device description API to ensure that the device supports the selected preferred format and size.

First, the source code opens the device with `device.Open` function call. Unlike in the [previous example](../capture0/capture0.go), the call to
`Open` omits the pixel format option.

```go
func main() {
	devName := "/dev/video0"
	device, err := device.Open(devName)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()
}
```

Next, the source code defines a function that is used to search formats supported by the device.

```go
func main() {
...
	findPreferredFmt := func(fmts []v4l2.FormatDescription, pixEncoding v4l2.FourCCType) *v4l2.FormatDescription {
		for _, desc := range fmts {
			if desc.PixelFormat == pixEncoding{
				return &desc
			}
		}
		return nil
	}
}
```

Next, the code enumerates the formats supported by the device, `device.GetFormatDescriptions`, and used the search function
to test whether the device support one of several preferred formats.

```go
func main() {
...
	fmtDescs, err := device.GetFormatDescriptions()
	if err != nil{
		log.Fatal("failed to get format desc:", err)
	}

	// search for preferred formats
	preferredFmts := []v4l2.FourCCType{v4l2.PixelFmtMPEG, v4l2.PixelFmtMJPEG, v4l2.PixelFmtJPEG, v4l2.PixelFmtYUYV}
	var fmtDesc *v4l2.FormatDescription
	for _, preferredFmt := range preferredFmts{
		fmtDesc = findPreferredFmt(fmtDescs, preferredFmt)
		if fmtDesc != nil {
			break
		}
	}
}
```

Next, if one of the preferred formats is found, then it is assigned to `fmtDesc`. The next step is to search the device 
for an appropriate supported dimension (640x480) for the selected format which is stored in `frmSize`.

```go
func main() {
...
    frameSizes, err := v4l2.GetFormatFrameSizes(device.Fd(), fmtDesc.PixelFormat)

	// select size 640x480 for format
	var frmSize v4l2.FrameSizeEnum
	for _, size := range frameSizes {
		if size.Size.MinWidth == 640 && size.Size.MinHeight == 480 {
			frmSize = size
			break
		}
	}
}
```

At this point, the device can be assigned the selected pixel format and its associated size.

```go
func main() {
...
	if err := device.SetPixFormat(v4l2.PixFormat{
		Width:       frmSize.Size.MinWidth,
		Height:      frmSize.Size.MinHeight,
		PixelFormat: fmtDesc.PixelFormat,
		Field:       v4l2.FieldNone,
	}); err != nil {
		log.Fatalf("failed to set format: %s", err)
	}
}
```

Finally, the device can be started and the streaming buffers can be captured:

```go
fun main() {
...
	if err := device.Start(ctx); err != nil {
		log.Fatalf("failed to stream: %s", err)
	}

	for frame := range device.GetOutput() {
		fileName := fmt.Sprintf("capture_%d.jpg", count)
		file, err := os.Create(fileName)
		if err != nil {
			log.Printf("failed to create file %s: %s", fileName, err)
			continue
		}
		if _, err := file.Write(frame); err != nil {
			log.Printf("failed to write file %s: %s", fileName, err)
			continue
		}
        ...
	}
}
```

> See source code [here](./capture1.go).