# Read/Write Capture Example

Demonstrates frame capture using the **read/write I/O method**, an alternative to the default streaming mode.

## Read/Write vs Streaming

| | Streaming (default) | Read/Write |
|---|---|---|
| **API** | `Start → GetFrames → Stop` | `Read` or `ReadFrame` |
| **Flow** | Continuous, channel-based | Synchronous, one frame at a time |
| **Performance** | Zero-copy (MMAP) | Extra copy per frame |
| **Complexity** | Higher (buffer lifecycle) | Lower (one syscall) |
| **Capability** | `V4L2_CAP_STREAMING` | `V4L2_CAP_READWRITE` |

Use read/write mode when:
- You need simple, synchronous frame access
- The device doesn't support streaming
- Performance is not critical

## Usage

```bash
# Capture 5 frames using ReadFrame() (high-level API with metadata)
go run . -d /dev/video0 -n 5

# Capture 5 frames using Read() (low-level API)
go run . -d /dev/video0 -n 5 -raw
```

## Flags

- `-d` — Device path (default: `/dev/video0`)
- `-n` — Number of frames to capture (default: 5)
- `-raw` — Use low-level `Read()` instead of `ReadFrame()`

## Output

Frames are saved as `frame_0.jpg`, `frame_1.jpg`, etc.

```
Format: 640x480, SizeImage: 614400
Frame 0: seq=0, 12543 bytes, ts=2026-04-05 10:30:00 -> frame_0.jpg
Frame 1: seq=1, 12601 bytes, ts=2026-04-05 10:30:00 -> frame_1.jpg
```
