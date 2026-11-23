package v4l2

import (
	"testing"
	"time"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

var (
	benchTimeSink   time.Time
	benchBufferSink Buffer
)

func withMockIoctlBench(b *testing.B, fn func(fd, req uintptr, arg unsafe.Pointer) sys.Errno) {
	b.Helper()
	orig := ioctlSyscall
	ioctlSyscall = fn
	b.Cleanup(func() {
		ioctlSyscall = orig
	})
}

func BenchmarkBufferTimestamp_Monotonic(b *testing.B) {
	tv := sys.Timeval{Sec: 12345, Usec: 678900}
	b.ReportAllocs()
	b.ResetTimer()

	var sink time.Time
	for i := 0; i < b.N; i++ {
		sink = BufferTimestamp(tv, BufFlagTimestampMonotonic, TimestampDefault)
	}
	benchTimeSink = sink
}

func BenchmarkBufferTimestamp_WallClock(b *testing.B) {
	tv := sys.Timeval{Sec: 12345, Usec: 678900}
	b.ReportAllocs()
	b.ResetTimer()

	var sink time.Time
	for i := 0; i < b.N; i++ {
		sink = BufferTimestamp(tv, BufFlagTimestampMonotonic, TimestampAuto)
	}
	benchTimeSink = sink
}

func BenchmarkQueueBufferWith_SinglePlanar(b *testing.B) {
	withMockIoctlBench(b, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		return 0
	})

	buf := Buffer{
		Index:     0,
		Type:      BufTypeVideoCapture,
		Memory:    IOTypeMMAP,
		BytesUsed: 614400,
		Length:    614400,
	}

	b.ReportAllocs()
	b.ResetTimer()

	var sink Buffer
	for i := 0; i < b.N; i++ {
		res, err := QueueBufferWith(3, buf)
		if err != nil {
			b.Fatal(err)
		}
		sink = res
	}
	benchBufferSink = sink
}

func BenchmarkQueueBufferWith_MPlane(b *testing.B) {
	withMockIoctlBench(b, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		return 0
	})

	buf := Buffer{
		Index:  0,
		Type:   BufTypeVideoCaptureMPlane,
		Memory: IOTypeMMAP,
		Planes: []Plane{
			{BytesUsed: 640 * 480, Length: 640 * 480},
			{BytesUsed: 640 * 240, Length: 640 * 240},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	var sink Buffer
	for i := 0; i < b.N; i++ {
		res, err := QueueBufferWith(3, buf)
		if err != nil {
			b.Fatal(err)
		}
		sink = res
	}
	benchBufferSink = sink
}

func BenchmarkDequeueBufferMPlane_TwoPlanes(b *testing.B) {
	withMockIoctlBench(b, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		cb := (*cV4L2Buffer)(arg)
		cb.Sequence = 42
		if cb.Length > 0 {
			planesPtr := *(*unsafe.Pointer)(unsafe.Pointer(&cb.M[0]))
			planes := unsafe.Slice((*cV4L2Plane)(planesPtr), cb.Length)
			for i := range planes {
				planes[i].Bytesused = 307200
				planes[i].Length = 307200
			}
		}
		return 0
	})

	b.ReportAllocs()
	b.ResetTimer()

	var sink Buffer
	for i := 0; i < b.N; i++ {
		res, err := DequeueBufferMPlane(3, IOTypeMMAP, BufTypeVideoCaptureMPlane, 2)
		if err != nil {
			b.Fatal(err)
		}
		sink = res
	}
	benchBufferSink = sink
}
