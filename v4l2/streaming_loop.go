package v4l2

import (
	"context"
	"fmt"

	sys "golang.org/x/sys/unix"
)

// StartStreamLoop issue a streaming request for the device and sets up
// a loop to capture incoming buffers from the device.
func StartStreamLoop(ctx context.Context, dev Device) (chan []byte, error) {
	if err := StreamOn(dev); err != nil {
		return nil, fmt.Errorf("stream loop: driver stream on: %w", err)
	}

	dataChan := make(chan []byte, dev.BufferCount())

	go func() {
		defer close(dataChan)
		for {
			select {
			case <-WaitForRead(dev):
				//TODO add better error-handling, for now just panic
				frame, err  := CaptureFrame(dev)
				if err != nil {
					panic(fmt.Errorf("stream loop: frame capture: %s", err).Error())
				}
				select {
				case dataChan <-frame:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return dataChan, nil
}

// StopStreamLoop unmaps allocated IO memory and signal device to stop streaming
func StopStreamLoop(dev Device) error {
	if dev.Buffers() == nil {
		return fmt.Errorf("stop loop: failed to stop loop: buffers uninitialized")
	}

	if err := StreamOff(dev); err != nil {
		return fmt.Errorf("stop loop: stream off: %w", err)
	}
	return nil
}

func WaitForRead(dev Device) <-chan struct{} {
	sigChan := make(chan struct{})

	fd := dev.FileDescriptor()

	go func() {
		defer close(sigChan)
		var fdsRead sys.FdSet
		fdsRead.Set(int(fd))
		for {
			n, err := sys.Select(int(fd+1), &fdsRead, nil, nil, nil)
			if n == -1 {
				if err == sys.EINTR {
					continue
				}
			}
			sigChan <- struct{}{}
		}
	}()

	return sigChan
}
