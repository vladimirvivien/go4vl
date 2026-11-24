# stateless_hevc Example

Demonstrates binding stateless HEVC compound controls (`ControlHEVCSPS`, `ControlHEVCPPS`, `ControlHEVCDecodeParams`) to a `v4l2.MediaRequest` and validating them against a stateless hardware HEVC decoder (such as `rpi-hevc-dec` on Raspberry Pi 4) via `v4l2.TryExtControls`.

## Usage

```bash
go run ./examples/stateless_hevc -video /dev/video19 -media /dev/media0
```

## Example Output

```
Video device: /dev/video19 (rpi-hevc-dec, driver=rpi-hevc-dec)
  OUTPUT format: HEVC Parsed Slice Data (0x35363253)
Validated stateless HEVC SPS/PPS/DecodeParams with request FD 5
Reinitialized MediaRequest successfully
```

(Host-specific device paths and file descriptor numbers differ.)
