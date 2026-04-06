# User Pointer Capture Example

Demonstrates frame capture using **User Pointer (USERPTR) streaming I/O**, an alternative to the default MMAP streaming mode.

## MMAP vs USERPTR

| | MMAP (default) | USERPTR |
|---|---|---|
| **Buffer allocation** | Kernel allocates | Application allocates |
| **Buffer mapping** | `mmap()` into userspace | Direct pointer access |
| **Performance** | Zero-copy from kernel | One copy per frame |
| **Use case** | General purpose | Custom allocators, shared memory |

The consumer API (`Start/GetFrames/Stop`) is identical for both modes.

## Usage

```bash
# Capture 5 frames using USERPTR streaming
go run . -d /dev/video0 -n 5
```

## Flags

- `-d` — Device path (default: `/dev/video0`)
- `-n` — Number of frames to capture (default: 5)

## Output

Frames are saved as `frame_0.jpg`, `frame_1.jpg`, etc.

```
Frame 0: seq=0, 12543 bytes, ts=2026-04-05 10:30:00 -> frame_0.jpg
Frame 1: seq=1, 12601 bytes, ts=2026-04-05 10:30:00 -> frame_1.jpg
Captured 5 frames
```
