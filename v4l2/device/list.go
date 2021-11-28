package device

import (
	"fmt"
	"os"
	"regexp"
)

var (
	root = "/dev"
)

// devPattern is device directory name pattern on Linux
var devPattern = regexp.MustCompile(fmt.Sprintf(`%s/(video|radio|vbi|swradio|v4l-subdev|v4l-touch|media)[0-9]+`, root))

// IsDevice tests whether the path matches a V4L device name
func IsDevice(devpath string) bool {
	return devPattern.MatchString(devpath)
}

func List() ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, entry := range entries {
		dev := fmt.Sprintf("%s/%s", root, entry.Name())
		if entry.Type() & os.ModeDevice != 0 && IsDevice(dev) {
			result = append(result, dev)
		}
	}

	return result,  nil
}