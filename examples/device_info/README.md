# Device info example

The example in this directory showcases `go4vl` support for device information. For instance, the following function prints driver information

```go
func main() {
    devName := '/dev/video0'
    device, err := device2.Open(devName)
    if err := printDeviceDriverInfo(device); err != nil {
        log.Fatal(err)
    }
}

func printDeviceDriverInfo(dev *device.Device) error {
	caps := dev.Capability()

	// print driver info
	fmt.Println("v4l2Device Info:")
	fmt.Printf(template, "Driver name", caps.Driver)
	fmt.Printf(template, "Card name", caps.Card)
	fmt.Printf(template, "Bus info", caps.BusInfo)

	fmt.Printf(template, "Driver version", caps.GetVersionInfo())

	fmt.Printf("\t%-16s : %0x\n", "Driver capabilities", caps.Capabilities)
	for _, desc := range caps.GetDriverCapDescriptions() {
		fmt.Printf("\t\t%s\n", desc.Desc)
	}

	fmt.Printf("\t%-16s : %0x\n", "v4l2Device capabilities", caps.Capabilities)
	for _, desc := range caps.GetDeviceCapDescriptions() {
		fmt.Printf("\t\t%s\n", desc.Desc)
	}

	return nil
}
```

> See the [complete example](./devinfo.go) and all available device information from go4vl.