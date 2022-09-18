# V4L2 video capture example in C

This an example in C showing a minimally required steps to capture video using V4L2. This is can be used to run tests on devices and compare results with the Go4VL code.

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