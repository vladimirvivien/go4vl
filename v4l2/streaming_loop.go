package v4l2

import (
	sys "golang.org/x/sys/unix"
)

// WaitForRead returns a channel that can be used to be notified when
// a device's is ready to be read.
func WaitForRead(dev Device) <-chan struct{} {
	sigChan := make(chan struct{})

	go func(fd uintptr) {
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
	}(dev.Fd())

	return sigChan
}
