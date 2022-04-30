package v4l2

import (
	sys "golang.org/x/sys/unix"
)

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
