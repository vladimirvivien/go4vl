# Capture example

This example shows how to use the `go4vl` API to create a simple program that captures several video frames, from an attached input (camera) device, and save them to  files.

Firstly, the source code opens a device, `devName`, with a hard-coded pixel format (MPEG) and size. If the device does not support 
the specified format, the open operation will fail, returning an error.

```go
func main() {
    devName := "/dev/video0"
	// open device
	device, err := device.Open(
		devName,
		device.WithPixFormat(
			v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMPEG, Width: 640, Height: 480},
		),
	)
...
}
```

Next, the source code calls the `device.Start` method to start the input (capture) process.

```go
func main() {
...
	// start stream
	ctx, stop := context.WithCancel(context.TODO())
	if err := device.Start(ctx); err != nil {
		log.Fatalf("failed to start stream: %s", err)
	}
...

}
```

Once the device starts, the code sets up a loop to capture incoming video frame buffers from the input device and save each frame to a local file.

```go
func main() {
...
	for frame := range device.GetOutput() {
		fileName := fmt.Sprintf("capture_%d.jpg", count)
		file, err := os.Create(fileName)
		...
		if _, err := file.Write(frame); err != nil {
			log.Printf("failed to write file %s: %s", fileName, err)
			continue
		}
		
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s: %s", fileName, err)
		}
	}

}
```

> See the full source code [here](./capture0.go).