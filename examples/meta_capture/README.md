# meta_capture Example

Discovers and inspects a V4L2 metadata capture or output node (`V4L2_BUF_TYPE_META_CAPTURE` / `V4L2_BUF_TYPE_META_OUTPUT`) using `v4l2.GetMetaFormat` and `v4l2.TryMetaFormat`. Common targets include ISP 3A statistics/parameters nodes (such as `bcm2835-isp`) and UVC/CSI metadata nodes.

## Usage

Auto-discover the first metadata device on the host:

```bash
go run ./examples/meta_capture
```

Inspect a specific metadata node:

```bash
go run ./examples/meta_capture -d /dev/video16
```

## Example Output

```
Device: /dev/video16 (bcm2835-isp, driver=bcm2835-isp, bufType=13)
Supported metadata formats (1):
  [0] BCM2835 ISP Statistics (0x41545342)
Active format: BCM2835 ISP Statistics (BSTA) [buffersize=21744]
Try format:    BCM2835 ISP Statistics (BSTA) [buffersize=21744]
```

(Host-specific device paths and buffer sizes differ by driver and sensor.)
