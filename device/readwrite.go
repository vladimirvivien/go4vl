package device

import (
	"errors"
	"fmt"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

// Read reads a single frame from the device into buf.
// Returns the number of bytes read.
//
// Only available when the device is opened with IOMethodReadWrite.
// The buffer should be at least SizeImage bytes (from GetPixFormat().SizeImage).
//
// Read blocks until a frame is available. The device is opened with O_NONBLOCK,
// so EAGAIN is returned if no frame is ready yet.
//
// Example:
//
//	dev, _ := device.Open("/dev/video0", device.WithIOMethod(device.IOMethodReadWrite))
//	defer dev.Close()
//
//	pf, _ := dev.GetPixFormat()
//	buf := make([]byte, pf.SizeImage)
//	n, err := dev.Read(buf)
//	// buf[:n] contains the frame data
func (d *Device) Read(buf []byte) (int, error) {
	if d.config.ioMethod != IOMethodReadWrite {
		return 0, fmt.Errorf("device: Read() not supported in streaming IO mode; use Start/GetFrames instead")
	}

	n, err := v4l2.ReadDevice(d.fd, buf)
	if err != nil {
		if errors.Is(err, sys.EAGAIN) || errors.Is(err, sys.EINTR) {
			return 0, err
		}
		return 0, fmt.Errorf("device: read: %w", err)
	}
	return n, nil
}

// ReadFrame reads a single frame from the device and returns it as a *Frame
// with metadata, providing parity with the streaming GetFrames() API.
//
// Each call allocates a fresh buffer. The returned Frame.Data is valid until
// the caller discards it (no Release() needed, unlike streaming mode frames).
//
// Only available when the device is opened with IOMethodReadWrite.
//
// Example:
//
//	dev, _ := device.Open("/dev/video0", device.WithIOMethod(device.IOMethodReadWrite))
//	defer dev.Close()
//
//	for i := 0; i < 10; i++ {
//	    frame, err := dev.ReadFrame()
//	    if err != nil { break }
//	    fmt.Printf("Frame %d: %d bytes at %v\n", frame.Sequence, len(frame.Data), frame.Timestamp)
//	}
func (d *Device) ReadFrame() (*Frame, error) {
	if d.config.ioMethod != IOMethodReadWrite {
		return nil, fmt.Errorf("device: ReadFrame() not supported in streaming IO mode; use Start/GetFrames instead")
	}

	buf := make([]byte, d.config.pixFormat.SizeImage)
	n, err := v4l2.ReadDevice(d.fd, buf)
	if err != nil {
		if errors.Is(err, sys.EAGAIN) || errors.Is(err, sys.EINTR) {
			return nil, err
		}
		return nil, fmt.Errorf("device: read frame: %w", err)
	}

	frame := &Frame{
		Data:      buf[:n],
		Timestamp: time.Now(),
		Sequence:  d.readSeq,
	}
	d.readSeq++

	return frame, nil
}
