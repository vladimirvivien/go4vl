# go4vl
A Go library for the Video for Linux user API (V4L2).

----

`go4vl` hides all the complexities of working with V4L2 and 
provides idiomatic Go types, like channels, to consume and process captured video frames.

## Features
* Capture and control video data from your Go programs
* Idiomatic Go API for device access and video capture
* Use familiar types such as channels to stream video data
* Exposes device enumeration and information
* Provides device capture control
* Access to video format information
* Streaming support using memory map (other methods coming soon)

## Getting started
To include `go4vl` in your own code, pull the package

```bash
go get github.com/vladimirvivien/go4vl/v4l2
```

## Example
