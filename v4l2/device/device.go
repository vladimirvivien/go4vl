package device

import (
	"context"
	"fmt"
	"os"
	sys "syscall"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
)

type Device struct {
	path             string
	file             *os.File
	fd               uintptr
	bufType          v4l2.BufType
	cap              v4l2.Capability
	cropCap          v4l2.CropCapability
	selectedFormat   v4l2.PixFormat
	supportedFormats []v4l2.PixFormat
	buffers          [][]byte
	requestedBuf     v4l2.RequestBuffers
	streaming        bool
}

// Open creates opens the underlying device at specified path
// and returns a *Device or an error if unable to open device.
func Open(path string) (*Device, error) {
	file, err := os.OpenFile(path, sys.O_RDWR|sys.O_NONBLOCK, 0644)
	if err != nil {
		return nil, fmt.Errorf("device open: %w", err)
	}
	dev := &Device{path: path, file: file, fd: file.Fd()}

	cap, err := v4l2.GetCapability(file.Fd())
	if err != nil {
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("device %s: closing after failure: %s", path, err)
		}
		return nil, fmt.Errorf("device %s: capability: %w", path, err)
	}
	dev.cap = cap

	switch {
	case cap.IsVideoCaptureSupported():
		dev.bufType = v4l2.BufTypeVideoCapture
	case cap.IsVideoOutputSupported():
		dev.bufType = v4l2.BufTypeVideoOutput
	default:
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("device %s: closing after failure: %s", path, err)
		}
		return nil, fmt.Errorf("unsupported device type")
	}

	cropCap, err := v4l2.GetCropCapability(file.Fd())
	if err != nil {
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("device %s: closing after failure: %s", path, err)
		}
		return nil, fmt.Errorf("device %s: crop capability: %w", path, err)
	}
	dev.cropCap = cropCap

	return dev, nil
}

// Close closes the underlying device associated with `d` .
func (d *Device) Close() error {
	if d.streaming {
		if err := d.StopStream(); err != nil {
			return err
		}
	}

	return d.file.Close()
}

// Name returns the device name (or path)
func (d *Device) Name() uintptr {
	return d.fd
}

// FileDescriptor returns the file descriptor value for the device
func (d *Device) FileDescriptor() uintptr {
	return d.fd
}

// Capability returns device capability info.
func (d *Device) Capability() v4l2.Capability {
	return d.cap
}

// GetCropCapability returns cropping info for device
func (d *Device) GetCropCapability() (v4l2.CropCapability, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.CropCapability{}, v4l2.ErrorUnsupportedFeature
	}
	return d.cropCap, nil
}

// SetCropRect crops the video dimension for the device
func (d *Device) SetCropRect(r v4l2.Rect) error {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.ErrorUnsupportedFeature
	}
	if err := v4l2.SetCropRect(d.fd, r); err != nil {
		return fmt.Errorf("device: %w", err)
	}
	return nil
}

// GetPixFormat retrieves pixel format info for device
func (d *Device) GetPixFormat() (v4l2.PixFormat, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.PixFormat{}, v4l2.ErrorUnsupportedFeature
	}
	pixFmt, err := v4l2.GetPixFormat(d.fd)
	if err != nil {
		return v4l2.PixFormat{}, fmt.Errorf("device: %w", err)
	}
	d.selectedFormat = pixFmt
	return pixFmt, nil
}

// SetPixFormat sets the pixel format for the associated device.
func (d *Device) SetPixFormat(pixFmt v4l2.PixFormat) error {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.ErrorUnsupportedFeature
	}

	if err := v4l2.SetPixFormat(d.fd, pixFmt); err != nil {
		return fmt.Errorf("device: %w", err)
	}
	d.selectedFormat = pixFmt
	return nil
}

// GetFormatDescription returns a format description for the device at specified format index
func (d *Device) GetFormatDescription(idx uint32) (v4l2.FormatDescription, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.FormatDescription{}, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetFormatDescription(d.fd, idx)
}

// GetFormatDescriptions returns all possible format descriptions for device
func (d *Device) GetFormatDescriptions() ([]v4l2.FormatDescription, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return nil, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetAllFormatDescriptions(d.fd)
}

// GetVideoInputIndex returns current video input index for device
func (d *Device) GetVideoInputIndex() (int32, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return 0, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetCurrentVideoInputIndex(d.fd)
}

// GetVideoInputInfo returns video input info for device
func (d *Device) GetVideoInputInfo(index uint32) (v4l2.InputInfo, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.InputInfo{}, v4l2.ErrorUnsupportedFeature
	}

	return v4l2.GetVideoInputInfo(d.fd, index)
}

// GetCaptureParam returns streaming capture parameter information
func (d *Device) GetCaptureParam() (v4l2.CaptureParam, error) {
	if !d.cap.IsVideoCaptureSupported() {
		return v4l2.CaptureParam{}, v4l2.ErrorUnsupportedFeature
	}
	return v4l2.GetStreamCaptureParam(d.fd)
}

// GetMediaInfo returns info for a device that supports the Media API
func (d *Device) GetMediaInfo() (v4l2.MediaDeviceInfo, error) {
	return v4l2.GetMediaDeviceInfo(d.fd)
}

func (d *Device) StartStream(buffSize uint32) error {
	if d.streaming {
		return nil
	}

	// allocate device buffers
	bufReq, err := v4l2.InitBuffers(d.fd, buffSize)
	if err != nil {
		return fmt.Errorf("device: start stream: %w", err)
	}
	d.requestedBuf = bufReq

	// for each device buff allocated, prepare local mapped buffer
	bufCount := int(d.requestedBuf.Count)
	d.buffers = make([][]byte, d.requestedBuf.Count)
	for i := 0; i < bufCount; i++ {
		buffer, err := v4l2.GetBuffer(d.fd, uint32(i))
		if err != nil {
			return fmt.Errorf("device start stream: %w", err)
		}

		offset := buffer.Info.Offset
		length := buffer.Length
		mappedBuf, err := v4l2.MapMemoryBuffer(d.fd, int64(offset), int(length))
		if err != nil {
			return fmt.Errorf("device start stream: %w", err)
		}
		d.buffers[i] = mappedBuf
	}

	// Initial enqueue of buffers for capture
	for i := 0; i < bufCount; i++ {
		_, err := v4l2.QueueBuffer(d.fd, uint32(i))
		if err != nil {
			return fmt.Errorf("device start stream: %w", err)
		}
	}

	// turn on device stream
	if err := v4l2.StreamOn(d.fd); err != nil {
		return fmt.Errorf("device start stream: %w", err)
	}

	d.streaming = true

	return nil
}

// Capture captures video buffer from device and emit
// each buffer on channel.
func (d *Device) Capture(ctx context.Context, fps uint32) (<-chan []byte, error) {
	if !d.streaming {
		return nil, fmt.Errorf("device: capture: streaming not started")
	}
	if ctx == nil {
		return nil, fmt.Errorf("device: context nil")
	}

	bufCount := int(d.requestedBuf.Count)
	dataChan := make(chan []byte, bufCount)

	if fps == 0 {
		fps = 10
	}

	// delay duration based on frame per second
	fpsDelay := time.Duration((float64(1) / float64(fps)) * float64(time.Second))

	go func() {
		defer close(dataChan)

		// capture forever or until signaled to stop
		for {
			// capture bufCount frames
			for i := 0; i < bufCount; i++ {
				//TODO add better error-handling during capture, for now just panic
				if err := v4l2.WaitForDeviceRead(d.fd, 2*time.Second); err != nil {
					panic(fmt.Errorf("device: capture: %w", err).Error())
				}

				// dequeue the device buf
				bufInfo, err := v4l2.DequeueBuffer(d.fd)
				if err != nil {
					panic(fmt.Errorf("device: capture: %w", err).Error())
				}

				// assert dequeued buffer is in proper range
				if !(int(bufInfo.Index) < bufCount) {
					panic(fmt.Errorf("device: capture: unexpected device buffer index: %d", bufInfo.Index).Error())
				}

				select {
				case dataChan <- d.buffers[bufInfo.Index][:bufInfo.BytesUsed]:
				case <-ctx.Done():
					return
				}
				// enqueu used buffer, prepare for next read
				if _, err := v4l2.QueueBuffer(d.fd, bufInfo.Index); err != nil {
					panic(fmt.Errorf("device capture: %w", err).Error())
				}

				time.Sleep(fpsDelay)
			}
		}
	}()

	return dataChan, nil
}

func (d *Device) StopStream() error {
	d.streaming = false
	for i := 0; i < len(d.buffers); i++ {
		if err := v4l2.UnmapMemoryBuffer(d.buffers[i]); err != nil {
			return fmt.Errorf("device: stop stream: %w", err)
		}
	}
	if err := v4l2.StreamOff(d.fd); err != nil {
		return fmt.Errorf("device: stop stream: %w", err)
	}
	return nil
}