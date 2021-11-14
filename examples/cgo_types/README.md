# Example: capture with cgo-generated types
:warning:

This example illustrates how to use ioctl directly
to communicate to device using cgo-generated types.

## Do not use it ## 

Use package `v4l2` to do realtime image capture, as shown in example
[examples/capture](../capture).

:warning:

The example in this directory shows most of the moving pieces that make
the V4L2 API works using Go.  It illustrates the steps, in detail, that
are required to communicate with a device driver to configure, initiate,
and capture images without using the Go v4l2 device type provided.