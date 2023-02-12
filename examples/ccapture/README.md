# V4L2 video capture example in C

This an example in C showing the minimally required steps to capture video using the V4L2 framework. This C code is used as a test tool to compare results between C and the Go4VL Go code.

## Build and run
On a Linux machine, run the following:

```
gcc -o capture capture.c
```

Run the program using:

```
./capture
```

Or, run `--help` to see available flags:

```
./capture --help
```

## Debugging with `strace`

To view the ioctl calls made when running the capture program:

```
strace -o trace.log -e trace=ioctl  ./capture
```