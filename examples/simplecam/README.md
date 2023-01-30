# simplecam

This is a simple example shows how easy it is to use go4vl to 
create a simple web application to stream camera images.

## Build
There are three ways to build the example.

### On device
You can always setup Go on your device and build your source code directly on the device.

### Cross-compile with Zig
The Zig language is (itself) full C/C++ cross-compiler. With flag compatibility with gcc and clang, 
Zig can be used as a drop-in replacement for those compilers allowing easy cross-compilation of source code.

It turns out that Zig can be used as the cross-compiler for building Go code with CGo enabled. Assuming you have
the Zig build tools installed, you can cross-compile the source code to target Linux/Arm/v7 as shown below:;

```
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC="zig cc -target arm-linux-musleabihf" CXX="zig c++ -target arm-linux-musleabihf" go build -o simple-cam .
```

The previous build command will create a staticly linked binary that can run on Linux for the Arm/v7 architecture:

```
simple-cam: ELF 32-bit LSB executable, ARM, EABI5 version 1 (SYSV), static-pie linked, Go BuildID=WYa4l3EGlIvd9EErrWkc/5Aa4CZdUXG8bERpToUcN/jjwKBQqSAfDbNfJGzSou/27sOKN7B1e0dtPc7PqmR, with debug_info, not stripped
```