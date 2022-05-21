package device

import (
	"testing"
)

func TestList(t *testing.T) {
	devices, err := GetAllDevicePaths()
	if err != nil {
		t.Error(err)
	}
	t.Logf("devices: %#v", devices)
}
