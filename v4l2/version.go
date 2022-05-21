package v4l2

import (
	"fmt"
)

type VersionInfo struct {
	value uint32
}

func (v VersionInfo) Major() uint32 {
	return v.value >> 16
}

func (v VersionInfo) Minor() uint32 {
	return (v.value >> 8) & 0xff
}

func (v VersionInfo) Patch() uint32 {
	return v.value & 0xff
}

// Value returns the raw numeric version value
func (v VersionInfo) Value() uint32 {
	return v.value
}

func (v VersionInfo) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major(), v.Minor(), v.Patch())
}
