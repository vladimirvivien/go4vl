# Example: capture with V4L2 directly
:warning: 

This example is here to illustrate the complexity of v4l2.
Do not use it. 

If you want to play around with image capture, use the 
[examples/capture](../capture).

:warning:

The example in this directory shows most of the moving pieces that make
the V4L2 API works using Go.  It illustrates the steps, in detail, that
are required to communicate with a device driver to configure, initiate,
and capture images without using the Go v4l2 device type provided.