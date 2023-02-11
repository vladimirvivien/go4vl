# Examples

* [snapshot](./snapshot/) - A simple example to capture a single frame and save it to a file.
* [capture0](./capture0) - Shows how to capture multiple images and saves them to files.
* [capture1](./capture1) - Shows how to capture multiple images using specified image format.
* [device_info](./device_info) - Uses go4vl to query and print device and format information.
* [format](./format) - Shows how to query and apply device and format information.
* [ext_ctrls](./ext_ctrls/) Shows how to query and apply extended controls.
* [user_ctrl](./user_ctrl/) Shows how to query and apply user controls.
* [simplecam](./simplecam/) A functional webcam program that streams video to web page.
* [webcam](./webcam) - Builds on simplecam and adds image control, format control, and face detection.

## Building your code

There are three ways to build the example.

### On device builds
One of the easiest ways to get started is to setup your Linux workstation, or device, with Go to build your source code directly there.

Install the `build-essential` package to install required C compilers:
```shell
sudo apt install build-essential
```
Also, upgrade your system to pull down the latest OS packages (follow directions for your system-specific steps):

```
sudo apt update
sudo apt full-upgrade
```

### Cross-compile with Zig toolchain
The Zig language comes with a full C/C++ cross-compiler including flag compatibility with gcc and clang. It can be used as a drop-in replacement for those compilers allowing easy cross-compilation of source code.

Zig cross-compilers can be used for building CGo-enabled Go code with little fuss. Assuming you have the Zig build tools installed, you can cross-compile the code in this directory as shown below:

```
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC="zig cc -target arm-linux-musleabihf" CXX="zig c++ -target arm-linux-musleabihf" go build -o simple-cam ./simplecam
```

The previous build command will create a static binary that can run on Linux for the Arm/v7 architecture.

### Cross-compile with Docker
There are some Docker images that you can use to properly cross-compile your code. Example [./simplecam](./simplecam/Dockerfile) includes a Dockerfile that uses image `crazymax/goxx` to cross-compile go4vl code.