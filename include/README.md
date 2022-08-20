# V4L2 Include Files

Here we provide a copy of the v4l2 header files as per the 
[kernel.org docs](https://docs.kernel.org/userspace-api/media/v4l/func-ioctl.html#description):

> Macros and defines specifying V4L2 ioctl requests are located in the `videodev2.h` header file.
> _**Applications should use their own copy, not include the version in the kernel sources on the system they compile on.**_

This helps to reduce compilation errors due to outdated headers installed in the user's system.

Headers have been obtained from `v4l-utils` repo as of commit [3b94a0ca4894d75de240b3ebb296071e551a261e](https://github.com/gjasny/v4l-utils/tree/3b94a0ca4894d75de240b3ebb296071e551a261e/include/linux)