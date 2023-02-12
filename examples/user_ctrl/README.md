# Device user control

The go4vl API has support for querying and setting values for device user control as demonstrated in this example.
For instance, the two functions below uses the go4vl API to set a user control and retrieve all user controls respectively.

```go
func setUserControlValue(device *dev.Device, ctrlID v4l2.CtrlID, val int) error {
	if ctrlID == 0 {
		return fmt.Errorf("invalid control specified")
	}
	return device.SetControlValue(ctrlID, v4l2.CtrlValue(val))
}

func listUserControls(device *dev.Device) {
	ctrls, err := device.QueryAllControls()
	if err != nil {
		log.Fatalf("query controls: %s", err)
	}

	for _, ctrl := range ctrls {
		printUserControl(ctrl)
	}
}
```

> See complete [source code](./ctrl.go).

