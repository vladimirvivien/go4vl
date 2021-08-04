package v4l2

import "bytes"

// toGoString encodes C null-terminated string as a Go string
func toGoString(s []byte) string {
	null := bytes.Index(s, []byte{0})
	if null < 0 {
		return ""
	}
	return string(s[:null])
}
