package device

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/vladimirvivien/go4vl/v4l2"
	sys "golang.org/x/sys/unix"
)

func TestFrame_Clone(t *testing.T) {
	origData := []byte{1, 2, 3, 4, 5}
	ts := time.Unix(1700000000, 123456000)
	orig := &Frame{
		Data:      origData,
		Timestamp: ts,
		Sequence:  42,
		Flags:     v4l2.BufFlagKeyFrame,
		Index:     3,
	}

	cloned := orig.Clone()
	if !bytes.Equal(cloned.Data, orig.Data) {
		t.Fatalf("Clone Data mismatch: got %v want %v", cloned.Data, orig.Data)
	}
	if cloned.Sequence != orig.Sequence || cloned.Flags != orig.Flags || cloned.Index != orig.Index || !cloned.Timestamp.Equal(orig.Timestamp) {
		t.Fatalf("Clone metadata mismatch: got %+v want %+v", cloned, orig)
	}

	// Mutating original slice must not affect clone.
	origData[0] = 99
	if cloned.Data[0] != 1 {
		t.Fatalf("Clone shares backing array with original")
	}
}

func TestWithBufferType(t *testing.T) {
	var cfg config
	WithBufferType(v4l2.BufTypeVideoOutput)(&cfg)
	if cfg.bufType != v4l2.BufTypeVideoOutput {
		t.Fatalf("bufType = %d, want %d", cfg.bufType, v4l2.BufTypeVideoOutput)
	}
}

func TestDevice_FramesValidation(t *testing.T) {
	// Read/write mode should yield an error.
	d := &Device{config: config{ioMethod: IOMethodReadWrite}}
	var gotErr error
	for _, err := range d.Frames(context.Background()) {
		gotErr = err
		break
	}
	if gotErr == nil {
		t.Fatalf("expected error in read/write mode")
	}

	// Streaming mode before Start() should yield an error.
	d2 := &Device{config: config{ioMethod: IOMethodStreaming}}
	gotErr = nil
	for _, err := range d2.Frames(context.Background()) {
		gotErr = err
		break
	}
	if gotErr == nil {
		t.Fatalf("expected error before Start()")
	}

	// Mutual exclusivity with GetFrames (mode 2).
	d3 := &Device{config: config{ioMethod: IOMethodStreaming}}
	d3.streaming.Store(true)
	d3.streamActive.Store(true)
	d3.streamingMode.Store(2)
	gotErr = nil
	for _, err := range d3.Frames(context.Background()) {
		gotErr = err
		break
	}
	if gotErr == nil {
		t.Fatalf("expected error when GetFrames mode is active")
	}
}

func TestDevice_ReadWriteContextPipe(t *testing.T) {
	var fds [2]int
	if err := sys.Pipe2(fds[:], sys.O_NONBLOCK|sys.O_CLOEXEC); err != nil {
		t.Fatalf("Pipe2: %v", err)
	}
	defer sys.Close(fds[0])
	defer sys.Close(fds[1])

	reader := &Device{
		fd:      uintptr(fds[0]),
		bufType: v4l2.BufTypeVideoCapture,
		config:  config{ioMethod: IOMethodReadWrite, pixFormat: v4l2.PixFormat{SizeImage: 16}},
	}
	writer := &Device{
		fd:      uintptr(fds[1]),
		bufType: v4l2.BufTypeVideoOutput,
		config:  config{ioMethod: IOMethodReadWrite},
	}

	payload := []byte("hello-v4l2-out")
	go func() {
		time.Sleep(15 * time.Millisecond)
		_, _ = writer.WriteContext(context.Background(), payload)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	frame, err := reader.ReadFrameContext(ctx)
	if err != nil {
		t.Fatalf("ReadFrameContext: %v", err)
	}
	if !bytes.Equal(frame.Data, payload) {
		t.Fatalf("ReadFrameContext data = %q, want %q", frame.Data, payload)
	}
}

func TestDevice_WriteFrameValidation(t *testing.T) {
	d := &Device{
		bufType: v4l2.BufTypeVideoCapture,
		config:  config{ioMethod: IOMethodStreaming},
	}
	if err := d.WriteFrame(context.Background(), []byte{1, 2, 3}); err == nil {
		t.Fatalf("expected error for capture device")
	}

	dOut := &Device{
		bufType: v4l2.BufTypeVideoOutput,
		config:  config{ioMethod: IOMethodStreaming},
	}
	if err := dOut.WriteFrame(context.Background(), []byte{1, 2, 3}); err == nil {
		t.Fatalf("expected error before Start()")
	}
}
