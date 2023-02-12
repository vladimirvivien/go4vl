# Extended controls

This example shows go4vl support for V4L2 extended controls. The API allows users to query and set control values.
For instance, the following snippet shows how to retrieve all extended controls for a given device.

```go
func main() {
	devName := "/dev/video0"
	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.Parse()

	device, err := dev.Open(devName)
	if err != nil {
		log.Fatalf("failed to open device: %s", err)
	}
	defer device.Close()

	ctrls, err := v4l2.QueryAllExtControls(device.Fd())
	if err != nil {
		log.Fatalf("failed to get ext controls: %s", err)
	}
	if len(ctrls) == 0 {
		log.Println("Device does not have extended controls")
		os.Exit(0)
	}
	for _, ctrl := range ctrls {
		printControl(ctrl)
	}
}
```

> See full example [source code](./extctrls.go).