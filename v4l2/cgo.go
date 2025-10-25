package v4l2

/*
#cgo linux CFLAGS: -I/usr/include

#include <linux/videodev2.h>
#include <linux/v4l2-controls.h>
#include <linux/v4l2-common.h>
*/
import "C"

// This file centralizes all CGO compiler directives for the v4l2 package.
//
// The default configuration uses system-provided V4L2 kernel headers from /usr/include.
// These headers are typically provided by the linux-libc-dev package (Debian/Ubuntu),
// kernel-headers package (RHEL/Fedora), or linux-headers package (Arch Linux).
//
// To use custom or newer kernel headers, override the include path using the CGO_CFLAGS
// environment variable:
//
//	CGO_CFLAGS="-I/path/to/custom/headers" go build
//
// For cross-compilation, point CGO_CFLAGS to your target's sysroot headers:
//
//	CGO_CFLAGS="-I/path/to/sysroot/usr/include" \
//	CC=aarch64-linux-gnu-gcc \
//	GOOS=linux GOARCH=arm64 \
//	go build
//
// The V4L2 headers are kernel UAPI (user-space API) headers and do not provide
// pkg-config files. Direct inclusion via CFLAGS is the standard approach.
