# Installation & Device Permissions

Setting up a Go project with `go4vl` requires no C compilers, `pkg-config` files, or Linux kernel headers on your development machine. Once the Go toolchain is installed, you add the module dependency and verify that your Linux user account has permission to open video character devices.

## Software requirements

Your build and target environments rely on a minimal stack:

* **Go 1.26 or later** (required for `iter.Seq2` support and current language features).
* **Target OS:** Linux with **kernel 5.10 or later** and V4L2 drivers enabled (standard on Ubuntu, Debian, Fedora, and Raspberry Pi OS).
* **Target Architectures:** `amd64`, `arm64`, `arm` (32-bit ARM), and `riscv64`.

## Adding `go4vl` to your module

Initialize a Go module if you do not have one yet, then fetch the library packages with `go get`:

```bash
go mod init example.com/mycam
go get github.com/vladimirvivien/go4vl/v4l2
```

Because `device` and `v4l2` live in the same Go module (`github.com/vladimirvivien/go4vl`), fetching the module makes both packages available for import in your Go code.

## Granting access to video devices

Linux creates character device nodes such as `/dev/video0`, `/dev/media0`, and `/dev/v4l-subdev0` owned by root and the `video` group. Opening these nodes without group membership fails with `permission denied`.

Add your user account to the `video` group cleanly on the target Linux host:

```bash
sudo usermod -a -G video "$USER"
```

Log out and start a new login session so the updated group membership applies, then confirm that your user can list and inspect the devices:

```bash
id -nG "$USER"
ls -l /dev/video*
```

On hardware with complex ISP pipelines (such as Raspberry Pi CSI-2 cameras), check both `/dev/media*` and `/dev/v4l-subdev*` ownership to confirm the same group permissions apply across all nodes your pipeline touches.

## Inspecting connected cameras with `v4l2-ctl`

Before writing code, inspecting your device's driver name, capabilities, and advertised pixel formats saves guesswork around supported resolutions and FourCC codes. Install `v4l-utils` if it is not present:

```bash
sudo apt-get update && sudo apt-get install -y v4l-utils
```

Query attached devices and inspect supported formats and frame rates on a target node:

```bash
v4l2-ctl --list-devices
v4l2-ctl -d /dev/video0 --all
v4l2-ctl -d /dev/video0 --list-formats-ext
```

Take note of the exact FourCC format (such as `MJPG` or `YUYV`) and supported resolutions printed by `--list-formats-ext`—you pass matching `v4l2.PixFormat` values when opening the device in the [Quick Start](quickstart.md).
