package device

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	devices, err := GetAllDevicePaths()
	if err != nil {
		t.Error(err)
	}
	t.Logf("devices: %#v", devices)
}

func TestListDevicePathsSysfsPrimary(t *testing.T) {
	sysDir := t.TempDir()
	devDir := t.TempDir()

	for _, name := range []string{"video2", "video0", "v4l-subdev0"} {
		if err := os.Mkdir(filepath.Join(sysDir, name), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(devDir, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := listDevicePaths(sysDir, devDir, func(path string) (bool, error) {
		_, err := os.Stat(path)
		return err == nil, err
	})
	if err != nil {
		t.Fatalf("listDevicePaths: %v", err)
	}

	want := []string{
		filepath.Join(devDir, "v4l-subdev0"),
		filepath.Join(devDir, "video0"),
		filepath.Join(devDir, "video2"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestListDevicePathsDevFallback(t *testing.T) {
	devDir := t.TempDir()
	for _, name := range []string{"video0", "video10", "null", "random"} {
		if err := os.WriteFile(filepath.Join(devDir, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := listDevicePaths(filepath.Join(t.TempDir(), "nonexistent"), devDir, func(path string) (bool, error) {
		_, err := os.Stat(path)
		return err == nil, err
	})
	if err != nil {
		t.Fatalf("listDevicePaths fallback: %v", err)
	}

	want := []string{
		filepath.Join(devDir, "video0"),
		filepath.Join(devDir, "video10"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
