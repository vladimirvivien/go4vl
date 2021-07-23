package v4l2

import "unsafe"

// v4l2 ioctl commands
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L2510
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/user-func.html
var (
	vidiocQueryCap   = encodeRead('V', 0, uintptr(unsafe.Sizeof(v4l2Capability{})))       // VIDIOC_QUERYCAP
	vidiocGetFormat  = encodeReadWrite('V', 4, uintptr(unsafe.Sizeof(Format{})))          // VIDIOC_G_FMT
	vidiocSetFormat  = encodeReadWrite('V', 5, uintptr(unsafe.Sizeof(Format{})))          // VIDIOC_S_FMT
	vidiocReqBufs    = encodeReadWrite('V', 8, uintptr(unsafe.Sizeof(RequestBuffers{})))  // VIDIOC_REQBUFS
	vidiocQueryBuf   = encodeReadWrite('V', 9, uintptr(unsafe.Sizeof(BufferInfo{})))      // VIDIOC_QUERYBUF
	vidiocQueueBuf   = encodeReadWrite('V', 15, uintptr(unsafe.Sizeof(BufferInfo{})))     // VIDIOC_QBUF
	vidiocDequeueBuf = encodeReadWrite('V', 17, uintptr(unsafe.Sizeof(BufferInfo{})))     // VIDIOC_DQBUF
	vidiocStreamOn   = encodeWrite('V', 18, uintptr(unsafe.Sizeof(int32(0))))             // VIDIOC_STREAMON
	vidiocStreamOff  = encodeWrite('V', 19, uintptr(unsafe.Sizeof(int32(0))))             // VIDIOC_STREAMOFF
	vidiocCropCap    = encodeReadWrite('V', 58, uintptr(unsafe.Sizeof(CropCapability{}))) // VIDIOC_CROPCAP
	vidiocSetCrop    = encodeWrite('V', 60, uintptr(unsafe.Sizeof(Crop{})))               // VIDIOC_S_CROP
)
