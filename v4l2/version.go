package v4l2

import (
	"fmt"
)

// VersionInfo represents a V4L2 driver version number.
// The version is encoded as a 32-bit value with major, minor, and patch components.
// Format: (major << 16) | (minor << 8) | patch
//
// Example: Version 5.10.42 would be encoded as 0x050A2A
type VersionInfo struct {
	value uint32
}

// Major returns the major version number (bits 31-16).
// For kernel 5.10.42, this returns 5.
func (v VersionInfo) Major() uint32 {
	return v.value >> 16
}

// Minor returns the minor version number (bits 15-8).
// For kernel 5.10.42, this returns 10.
func (v VersionInfo) Minor() uint32 {
	return (v.value >> 8) & 0xff
}

// Patch returns the patch version number (bits 7-0).
// For kernel 5.10.42, this returns 42.
func (v VersionInfo) Patch() uint32 {
	return v.value & 0xff
}

// Value returns the raw 32-bit encoded version value.
// This can be used for version comparisons or when the raw value is needed.
func (v VersionInfo) Value() uint32 {
	return v.value
}

// String returns a human-readable version string in the format "vMAJOR.MINOR.PATCH".
// For example: "v5.10.42"
func (v VersionInfo) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major(), v.Minor(), v.Patch())
}
