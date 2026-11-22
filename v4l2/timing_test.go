package v4l2

import (
	"testing"
	"time"

	sys "golang.org/x/sys/unix"
)

func TestBufferTimestampModes(t *testing.T) {
	var mono sys.Timespec
	if err := sys.ClockGettime(sys.CLOCK_MONOTONIC, &mono); err != nil {
		t.Fatalf("ClockGettime(CLOCK_MONOTONIC): %v", err)
	}
	tv := sys.Timeval{
		Sec:  mono.Sec,
		Usec: mono.Nsec / 1000,
	}

	// Default mode preserves raw kernel timestamp.
	raw := BufferTimestamp(tv, BufFlagTimestampMonotonic, TimestampDefault)
	if raw.Unix() != int64(tv.Sec) {
		t.Fatalf("TimestampDefault sec = %d, want %d", raw.Unix(), tv.Sec)
	}

	// Auto mode converts monotonic timestamp to wall-clock time near time.Now().
	wallAuto := BufferTimestamp(tv, BufFlagTimestampMonotonic, TimestampAuto)
	if diff := time.Since(wallAuto).Abs(); diff > 250*time.Millisecond {
		t.Fatalf("TimestampAuto wall clock drift too large: %v", diff)
	}

	// Auto mode leaves non-monotonic timestamps untouched.
	epochTV := sys.Timeval{Sec: 1700000000, Usec: 500000}
	wallNonMono := BufferTimestamp(epochTV, BufFlagTimestampUnknown, TimestampAuto)
	if wallNonMono.Unix() != 1700000000 {
		t.Fatalf("TimestampAuto non-monotonic sec = %d, want 1700000000", wallNonMono.Unix())
	}

	// MonoToAbs converts regardless of flags.
	wallForced := BufferTimestamp(tv, 0, TimestampMonoToAbs)
	if diff := time.Since(wallForced).Abs(); diff > 250*time.Millisecond {
		t.Fatalf("TimestampMonoToAbs wall clock drift too large: %v", diff)
	}

	buf := Buffer{Timestamp: tv, Flags: BufFlagTimestampMonotonic}
	if buf.TimestampType() != BufFlagTimestampMonotonic {
		t.Fatalf("TimestampType = 0x%x, want 0x%x", buf.TimestampType(), BufFlagTimestampMonotonic)
	}
	if diff := time.Since(buf.Time(TimestampAuto)).Abs(); diff > 250*time.Millisecond {
		t.Fatalf("Buffer.Time(TimestampAuto) drift too large: %v", diff)
	}
}
