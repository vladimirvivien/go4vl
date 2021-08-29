package v4l2

import (
	"context"
	"fmt"
	"os"
	sys "syscall"
	"time"
)

type Device struct {
	path         string
	file         *os.File
	fd           uintptr
	cap          *Capability
	cropCap      *CropCapability
	pixFormat    PixFormat
	buffers      [][]byte
	requestedBuf RequestBuffers
	streaming    bool
}

// Open creates opens the underlying device at specified path
// and returns a *Device or an error if unable to open device.
func Open(path string) (*Device, error) {
	file, err := os.OpenFile(path, sys.O_RDWR|sys.O_NONBLOCK, 0666)
	if err != nil {
		return nil, fmt.Errorf("device open: %w", err)
	}
	return &Device{path: path, file: file, fd: file.Fd()}, nil
}

// Close closes the underlying device associated with `d` .
func (d *Device) Close() error {
	if d.streaming{
		if err := d.StopStream(); err != nil{
			return err
		}
	}

	return d.file.Close()
}

// GetCapability retrieves device capability info and
// caches it for future capability check.
func (d *Device) GetCapability() (*Capability, error) {
	if d.cap != nil {
		return d.cap, nil
	}
	cap, err := GetCapability(d.fd)
	if err != nil {
		return nil, fmt.Errorf("device: %w", err)
	}
	d.cap = &cap
	return d.cap, nil
}

// GetCropCapability returns cropping info for device `d`
// and caches it for future capability check.
func (d *Device) GetCropCapability() (CropCapability, error) {
	if d.cropCap != nil {
		return *d.cropCap, nil
	}
	if err := d.assertVideoCaptureSupport(); err != nil {
		return CropCapability{}, fmt.Errorf("device: %w", err)
	}

	cropCap, err := GetCropCapability(d.fd)
	if err != nil {
		return CropCapability{}, fmt.Errorf("device: %w", err)
	}
	d.cropCap = &cropCap
	return cropCap, nil
}

// SetCropRect crops the video dimension for the device
func (d *Device) SetCropRect(r Rect) error {
	if err := d.assertVideoCaptureSupport(); err != nil {
		return fmt.Errorf("device: %w", err)
	}
	if err := SetCropRect(d.fd, r); err != nil {
		return fmt.Errorf("device: %w", err)
	}
	return nil
}

// GetPixFormat retrieves pixel format info for device
func (d *Device) GetPixFormat() (PixFormat, error) {
	if err := d.assertVideoCaptureSupport(); err != nil {
		return PixFormat{}, fmt.Errorf("device: %w", err)
	}
	pixFmt, err := GetPixFormat(d.fd)
	if err != nil {
		return PixFormat{}, fmt.Errorf("device: %w", err)
	}
	return pixFmt, nil
}

// SetPixFormat sets the pixel format for the associated device.
func (d *Device) SetPixFormat(pixFmt PixFormat) error {
	if err := d.assertVideoCaptureSupport(); err != nil {
		return fmt.Errorf("device: %w", err)
	}

	if err := SetPixFormat(d.fd, pixFmt); err != nil {
		return fmt.Errorf("device: %w", err)
	}
	return nil
}

func (d *Device) GetFormatDescriptions() ([]FormatDescription, error) {
	if err := d.assertVideoCaptureSupport(); err != nil {
		return nil, fmt.Errorf("device: %w", err)
	}

	return GetAllFormatDescriptions(d.fd)
}

func (d *Device) StartStream(buffSize uint32) error {
	if d.streaming {
		return nil
	}
	if err := d.assertVideoStreamSupport(); err != nil {
		return fmt.Errorf("device: %w", err)
	}

	// allocate device buffers
	bufReq, err := AllocateBuffers(d.fd, buffSize)
	if err != nil {
		return fmt.Errorf("device: start stream: %w", err)
	}
	d.requestedBuf = bufReq

	// for each device buff allocated, prepare local mapped buffer
	bufCount := int(d.requestedBuf.Count)
	d.buffers = make([][]byte, d.requestedBuf.Count)
	for i := 0; i < bufCount; i++ {
		bufInfo, err := GetBufferInfo(d.fd, uint32(i))
		if err != nil {
			return fmt.Errorf("device: start stream: %w", err)
		}

		offset := bufInfo.GetService().Offset
		length := bufInfo.Length
		mappedBuf, err := MapMemoryBuffer(d.fd, int64(offset), int(length))
		if err != nil {
			return fmt.Errorf("device: start stream: %w", err)
		}
		d.buffers[i] = mappedBuf
	}

	// Initial enqueue of buffers for capture
	for i := 0; i < bufCount; i++ {
		_, err := QueueBuffer(d.fd, uint32(i))
		if err != nil {
			return fmt.Errorf("device: start stream: %w", err)
		}
	}

	// turn on device stream
	if err := StreamOn(d.fd); err != nil {
		return fmt.Errorf("device: start stream: %w", err)
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
				if err := WaitForDeviceRead(d.fd, 2*time.Second); err != nil {
					panic(fmt.Errorf("device: capture: %w", err).Error())
				}

				// dequeue the device buf
				bufInfo, err := DequeueBuffer(d.fd)
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
				if _, err := QueueBuffer(d.fd, bufInfo.Index); err != nil {
					panic(fmt.Errorf("device capture: %w", err).Error())
				}

				time.Sleep(fpsDelay)
			}
		}
	}()

	return dataChan, nil
}

func (d *Device) StopStream() error{
	d.streaming = false
	for i := 0; i < len(d.buffers); i++ {
		if err := UnmapMemoryBuffer(d.buffers[i]); err != nil {
			return fmt.Errorf("device: stop stream: %w", err)
		}
	}
	if err := StreamOff(d.fd); err != nil {
		return fmt.Errorf("device: stop stream: %w", err)
	}
	return nil
}

func (d *Device) assertVideoCaptureSupport() error {
	cap, err := d.GetCapability()
	if err != nil {
		return fmt.Errorf("device capability: %w", err)
	}
	if !cap.IsVideoCaptureSupported() {
		return fmt.Errorf("device capability: video capture not supported")
	}
	return nil
}

func (d *Device) assertVideoStreamSupport() error {
	cap, err := d.GetCapability()
	if err != nil {
		return fmt.Errorf("device capability: %w", err)
	}
	if !cap.IsVideoCaptureSupported() {
		return fmt.Errorf("device capability: video capture not supported")
	}
	if !cap.IsStreamingSupported() {
		return fmt.Errorf("device capability: streaming not supported")
	}
	return nil
}
