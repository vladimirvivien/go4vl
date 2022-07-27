# Example: using handcrafted types (Deprecated) 
:warning: 

This example is here to illustrate the complexity of v4l2 when using hand-crafted 
Go types to communicate with the driver.  This approach was abandoned in favor of 
cgo-generated types (see v4l2 package) for stability.

Do not use it. 

If you want to play around with image capture, use the 
[capture0](../capture0) or [capture1](../capture1) examples.

:warning:

The example in this directory shows most of the moving pieces that make
the V4L2 API works using Go.  It illustrates the steps, in detail, that
are required to communicate with a device driver to configure, initiate,
and capture images without using the Go v4l2 device type provided.