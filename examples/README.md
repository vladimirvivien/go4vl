# Examples

## Basic Capture Examples

* [snapshot](./snapshot/) - A simple example to capture a single frame and save it to a file.
* [capture0](./capture0) - Shows how to capture multiple images and saves them to files (legacy API).
* [capture_frames](./capture_frames) - **Recommended**: Capture frames with metadata and buffer pooling for high-performance streaming.
* [capture1](./capture1) - Shows how to capture multiple images using specified image format.

## Device Information and Control

* [device_info](./device_info) - Uses go4vl to query and print device and format information.
* [format](./format) - Shows how to query and apply device and format information.
* [video_inputs](./video_inputs/) - Enumerate and select video inputs (for devices with multiple inputs).
* [video_outputs](./video_outputs/) - Enumerate and select video outputs (for output devices and v4l2loopback).
* [ext_ctrls](./ext_ctrls/) Shows how to query and apply extended controls.
* [user_ctrl](./user_ctrl/) Shows how to query and apply user controls.

## Web Streaming Examples

* [simplecam](./simplecam/) A functional webcam program that streams video to web page.
* [webcam](./webcam) - Builds on simplecam and adds image control, format control, and face detection.

## Building the example code

For comprehensive build instructions including prerequisite installation, cross-compilation, and troubleshooting, see [docs/BUILD.md](../docs/BUILD.md).

There are three ways to build the code in the example directories.

### On-device build
One of the easiest ways to get started is to setup your Linux workstation (with camera attached), or device (such as Raspberry Pi), with Go to build your source code directly there.

Install required packages:
```shell
# Install C compiler and build tools
sudo apt install build-essential

# Install V4L2 kernel headers (required for CGO)
sudo apt install linux-libc-dev
```

Upgrade your system to pull down the latest OS packages:
```
sudo apt update
sudo apt full-upgrade
```

Build an example:
```shell
go build ./snapshot
```

### Cross-compile with Zig toolchain
If you would rather cross-compile the code from a different location (i.e. your MacOS laptop or x86 Linux machine), then
you will need tooling to do the CGo-enabled cross-compilation of the C code generated for the code.  One easy way to do this is with the [Zig language](https://ziglang.org/) toolchain.

Zig comes with a full C/C++ cross-compiler including flag compatibility with `gcc` and `clang`. It can be used as a drop-in replacement for those compilers allowing easy cross-compilation of C source code. Zig cross-compilers can be used for building CGo-enabled Go code with little fuss. Assuming you have the Zig build tools installed, you can easily cross-compile the code in this directory.

For instance, the following cross-compiles the [./simplecam](./simplecam/) example into a static binary:

```
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC="zig cc -target arm-linux-musleabihf" CXX="zig c++ -target arm-linux-musleabihf" go build -o simple-cam ./simplecam
```

The previous build command will create a static binary that can run on Linux/Arm/v7 architecture.

### Cross-compile with Docker
Another way you can achieve cross compilation is with Docker. If you already have Docker as part of your workflow, you will find some images that you can use to cross-compile the code in this directory. For instance, the simplecam example includes a [./simplecam/Dockerfile](./simplecam/Dockerfile) that uses image `crazymax/goxx` to cross-compile the go4vl code.

## Using Custom V4L2 Headers

By default, go4vl uses system V4L2 headers from `/usr/include`. To use custom or newer kernel headers:

### Specify Custom Headers

```bash
# Use headers from a specific kernel version
CGO_CFLAGS="-I/usr/src/linux-6.8/include" go build ./snapshot

# Use downloaded kernel headers
CGO_CFLAGS="-I$HOME/linux-headers/include" go build ./snapshot
```

### Cross-Compilation with Custom Headers

```bash
# Cross-compile for ARM with custom headers
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=arm \
GOARM=7 \
CC="zig cc -target arm-linux-musleabihf" \
CGO_CFLAGS="-I/path/to/arm/sysroot/usr/include" \
go build ./snapshot
```

For more details, see [docs/BUILD.md](../docs/BUILD.md).