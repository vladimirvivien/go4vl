package device

import (
	"fmt"

	"github.com/vladimirvivien/go4vl/v4l2"
)

// GetControl queries the device for information about the specified control id.
func (d *Device) GetControl(ctrlID v4l2.CtrlID) (v4l2.Control, error) {
	ctlr, err := v4l2.GetControl(d.fd, ctrlID)
	if err != nil {
		return v4l2.Control{}, fmt.Errorf("device: %s: %w", d.path, err)
	}
	return ctlr, nil
}

// SetControlValue updates the value of the specified control id.
func (d *Device) SetControlValue(ctrlID v4l2.CtrlID, val v4l2.CtrlValue) error {
	err := v4l2.SetControlValue(d.fd, ctrlID, val)
	if err != nil {
		return fmt.Errorf("device: %s: %w", d.path, err)
	}
	return nil
}

// QueryAllControls fetches all supported device controls and their current values.
func (d *Device) QueryAllControls() ([]v4l2.Control, error) {
	ctrls, err := v4l2.QueryAllControls(d.fd)
	if err != nil {
		return nil, fmt.Errorf("device: %s: %w", d.path, err)
	}
	return ctrls, nil
}

// SetControlBrightness is a convenience method for setting value for control v4l2.CtrlBrightness
func (d *Device) SetControlBrightness(val v4l2.CtrlValue) error {
	return d.SetControlValue(v4l2.CtrlBrightness, val)
}

// SetControlContrast is a convenience method for setting value for control v4l2.CtrlContrast
func (d *Device) SetControlContrast(val v4l2.CtrlValue) error {
	return d.SetControlValue(v4l2.CtrlContrast, val)
}

// SetControlSaturation is a convenience method for setting value for control v4l2.CtrlSaturation
func (d *Device) SetControlSaturation(val v4l2.CtrlValue) error {
	return d.SetControlValue(v4l2.CtrlSaturation, val)
}

// SetControlHue is a convenience method for setting value for control v4l2.CtrlHue
func (d *Device) SetControlHue(val v4l2.CtrlValue) error {
	return d.SetControlValue(v4l2.CtrlHue, val)
}
