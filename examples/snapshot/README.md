# Snapshot

This program is a simple example that shows how to use the go4vl API to capture a single frame, from an attached camera device, and save it to a file. The example assumes the attached camera device supports JPEG (MJPEG) format inherently.

First, the `device` package is used to open a device with name `/dev/vide0`. If the device is not available or cant be opened (with the default configuration), the driver will return an error.

```go
func main() {
	dev, err := device.Open("/dev/video0")
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

    ...
}
```

Next, the device is started, with a context, and if no error is returned, it is ready to capture video data.

```go
func main() {
    ...
	if err := dev.Start(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
```

Next, the source code use variable `dev` to capture the frame and save the binary data to a file.

```go
func main() {
    ...
	frame := <-dev.GetOutput()

	file, err := os.Create("pic.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.Write(frame); err != nil {
		log.Fatal(err)
	}
}
```

> See the full [source code](./snap.go).