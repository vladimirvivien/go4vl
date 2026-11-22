package device

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var (
	devRoot  = "/dev"
	sysV4L2  = "/sys/class/video4linux"
	devRegex = regexp.MustCompile(`^(video|radio|vbi|swradio|v4l-subdev|v4l-touch|media)[0-9]+$`)
)

// IsDevice tests whether devpath points to a character or block device file,
// following symlinks if necessary.
func IsDevice(devpath string) (bool, error) {
	stat, err := os.Stat(devpath)
	if err != nil {
		return false, err
	}
	return stat.Mode()&os.ModeDevice != 0, nil
}

// GetAllDevicePaths returns a sorted slice of mounted V4L2 device paths (e.g.,
// "/dev/video0", "/dev/v4l-subdev0").
//
// It enumerates entries from /sys/class/video4linux and resolves each to /dev/<name>.
// If /sys/class/video4linux is unavailable (for example, in a minimal container without
// sysfs), it falls back to scanning /dev.
func GetAllDevicePaths() ([]string, error) {
	return listDevicePaths(sysV4L2, devRoot, IsDevice)
}

func listDevicePaths(sysDir, devDir string, checkDevice func(string) (bool, error)) ([]string, error) {
	if entries, err := os.ReadDir(sysDir); err == nil {
		result := make([]string, 0, len(entries))
		for _, entry := range entries {
			dev := filepath.Join(devDir, entry.Name())
			ok, err := checkDevice(dev)
			if err != nil {
				continue
			}
			if ok {
				result = append(result, dev)
			}
		}
		sort.Strings(result)
		return result, nil
	}

	entries, err := os.ReadDir(devDir)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	var result []string
	for _, entry := range entries {
		if !devRegex.MatchString(entry.Name()) {
			continue
		}
		dev := filepath.Join(devDir, entry.Name())
		ok, err := checkDevice(dev)
		if err != nil {
			continue
		}
		if ok {
			result = append(result, dev)
		}
	}
	sort.Strings(result)
	return result, nil
}
