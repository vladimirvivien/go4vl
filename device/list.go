package device

import (
	"fmt"
	"os"
	"regexp"
)

var (
	root = "/dev"
)

// devPattern is device directory name pattern on Linux (i.e. video0, video10, vbi0, etc)
var devPattern = regexp.MustCompile(fmt.Sprintf(`%s/(video|radio|vbi|swradio|v4l-subdev|v4l-touch|media)[0-9]+`, root))

// IsDevice tests whether the path matches a V4L device name and is a device file
func IsDevice(devpath string) (bool, error) {
	stat, err := os.Stat(devpath)
	if err != nil {
		return false, err
	}
	if stat.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(devpath)
		if err != nil {
			return false, err
		}
		return IsDevice(target)
	}
	if stat.Mode()&os.ModeDevice != 0 {
		return true, nil
	}
	return false, nil
}

// GetAllDevicePaths return a slice of all mounted v4l2 devices
func GetAllDevicePaths() ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, entry := range entries {
		dev := fmt.Sprintf("%s/%s", root, entry.Name())
		if !devPattern.MatchString(dev) {
			continue
		}
		ok, err := IsDevice(dev)
		if err != nil {
			return result, err
		}
		if ok {
			result = append(result, dev)
		}
	}
	return result, nil
}
