package v4l2

import (
	"fmt"

	sys "golang.org/x/sys/unix"
)

// StopStreamLoop unmaps allocated IO memory and signal device to stop streaming
func StopStreamLoop(dev StreamingDevice) error {
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

	fd := dev.Fd()

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
