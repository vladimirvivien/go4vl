package device

import (
	"testing"
)

func TestList(t *testing.T){
	devices, err := List()
	if err != nil {
		t.Error(err)
	}
	t.Logf("devices: %#v", devices)
}
