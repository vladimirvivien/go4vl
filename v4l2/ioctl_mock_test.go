package v4l2

import (
	"errors"
	"strings"
	"testing"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

func withMockIoctl(t *testing.T, fn func(fd, req uintptr, arg unsafe.Pointer) sys.Errno) {
	t.Helper()
	orig := ioctlSyscall
	ioctlSyscall = fn
	t.Cleanup(func() {
		ioctlSyscall = orig
	})
}

func TestSend_EINTRRetryAndErrorMapping(t *testing.T) {
	t.Run("EINTR_retries_until_success", func(t *testing.T) {
		calls := 0
		withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
			calls++
			if calls < 3 {
				return sys.EINTR
			}
			return 0
		})
		if err := send(3, cVIDIOC_QUERYCAP, nil); err != nil {
			t.Fatalf("send() = %v, want nil after EINTR retries", err)
		}
		if calls != 3 {
			t.Fatalf("ioctl calls = %d, want 3", calls)
		}
	})

	cases := []struct {
		name     string
		errno    sys.Errno
		sentinel error
	}{
		{"EINVAL", sys.EINVAL, ErrorBadArgument},
		{"ENOTTY", sys.ENOTTY, ErrorUnsupported},
		{"EIO", sys.EIO, ErrorSystem},
		{"EAGAIN", sys.EAGAIN, sys.EAGAIN},
		{"ETIMEDOUT", sys.ETIMEDOUT, sys.ETIMEDOUT},
		{"ENODEV", sys.ENODEV, sys.ENODEV},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
				return tc.errno
			})
			err := send(3, cVIDIOC_QUERYCAP, nil)
			if err == nil {
				t.Fatalf("expected error for %s", tc.name)
			}
			if !errors.Is(err, tc.errno) {
				t.Errorf("errors.Is(%v, %v) = false", err, tc.errno)
			}
			if !errors.Is(err, tc.sentinel) {
				t.Errorf("errors.Is(%v, %v) = false", err, tc.sentinel)
			}
		})
	}
}

func TestSelection_MockIoctl(t *testing.T) {
	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		sel := (*cV4L2Selection)(arg)
		switch req {
		case cVIDIOC_G_SELECTION:
			sel.R = cV4L2Rect{Left: 10, Top: 20, Width: 640, Height: 480}
			return 0
		case cVIDIOC_S_SELECTION:
			// Driver aligns width to multiple of 16
			sel.R.Width = (sel.R.Width / 16) * 16
			return 0
		default:
			return sys.ENOTTY
		}
	})

	got, err := GetSelection(3, BufTypeVideoCapture, SelTargetCropBounds)
	if err != nil {
		t.Fatalf("GetSelection() error = %v", err)
	}
	if got.Rect.Left != 10 || got.Rect.Top != 20 || got.Rect.Width != 640 || got.Rect.Height != 480 {
		t.Fatalf("GetSelection() rect = %+v, want 10,20 640x480", got.Rect)
	}

	set, err := SetSelection(3, Selection{
		Type:   BufTypeVideoCapture,
		Target: SelTargetCrop,
		Flags:  SelFlagLE,
		Rect:   Rect{Left: 0, Top: 0, Width: 645, Height: 480},
	})
	if err != nil {
		t.Fatalf("SetSelection() error = %v", err)
	}
	if set.Rect.Width != 640 {
		t.Fatalf("SetSelection() adjusted width = %d, want 640", set.Rect.Width)
	}
}

func TestFormatDescriptions_MockIoctl(t *testing.T) {
	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		if req != cVIDIOC_ENUM_FMT {
			return sys.ENOTTY
		}
		d := (*cV4L2FmtDesc)(arg)
		switch d.Index {
		case 0:
			copy(d.Description[:], "YUYV 4:2:2")
			d.Pixelformat = PixelFmtYUYV
			return 0
		case 1:
			copy(d.Description[:], "Motion-JPEG")
			d.Pixelformat = PixelFmtMJPEG
			d.Flags = FmtDescFlagCompressed
			return 0
		default:
			return sys.EINVAL
		}
	})

	single, err := GetFormatDescription(3, 0)
	if err != nil {
		t.Fatalf("GetFormatDescription(0) = %v", err)
	}
	if single.PixelFormat != PixelFmtYUYV || !strings.Contains(single.String(), "YUYV") {
		t.Fatalf("unexpected format desc: %+v (%q)", single, single.String())
	}

	all, err := GetAllFormatDescriptions(3)
	if err != nil {
		t.Fatalf("GetAllFormatDescriptions() = %v", err)
	}
	if len(all) != 2 || all[1].PixelFormat != PixelFmtMJPEG {
		t.Fatalf("GetAllFormatDescriptions() = %+v", all)
	}

	found, err := GetFormatDescriptionByEncoding(3, PixelFmtMJPEG)
	if err != nil || found.PixelFormat != PixelFmtMJPEG {
		t.Fatalf("GetFormatDescriptionByEncoding(MJPEG) = %+v, %v", found, err)
	}
	if _, err := GetFormatDescriptionByEncoding(3, PixelFmtH264); err == nil {
		t.Fatalf("expected error for unsupported encoding H264")
	}
}

func TestFrameSizesAndIntervals_MockIoctl(t *testing.T) {
	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		switch req {
		case cVIDIOC_ENUM_FRAMESIZES:
			fs := (*cV4L2FrmSizeEnum)(arg)
			switch fs.Index {
			case 0:
				fs.Type = FrameSizeTypeDiscrete
				*(*FrameSizeDiscrete)(unsafe.Pointer(&fs.Discrete)) = FrameSizeDiscrete{Width: 640, Height: 480}
				return 0
			case 1:
				fs.Type = FrameSizeTypeStepwise
				*(*FrameSize)(unsafe.Pointer(&fs.Discrete)) = FrameSize{
					MinWidth: 160, MaxWidth: 1920, StepWidth: 16,
					MinHeight: 120, MaxHeight: 1080, StepHeight: 16,
				}
				return 0
			default:
				return sys.EINVAL
			}
		case cVIDIOC_ENUM_FRAMEINTERVALS:
			fi := (*cV4L2FrmIvalEnum)(arg)
			switch fi.Index {
			case 0:
				fi.Type = FrameIntervalTypeDiscrete
				*(*Fract)(unsafe.Pointer(&fi.Discrete)) = Fract{Numerator: 1, Denominator: 30}
				return 0
			case 1:
				fi.Type = FrameIntervalTypeStepwise
				*(*FrameInterval)(unsafe.Pointer(&fi.Discrete)) = FrameInterval{
					Min:  Fract{Numerator: 1, Denominator: 60},
					Max:  Fract{Numerator: 1, Denominator: 15},
					Step: Fract{Numerator: 1, Denominator: 120},
				}
				return 0
			default:
				return sys.EINVAL
			}
		default:
			return sys.ENOTTY
		}
	})

	sizes, err := GetFormatFrameSizes(3, PixelFmtYUYV)
	if err != nil {
		t.Fatalf("GetFormatFrameSizes() = %v", err)
	}
	if len(sizes) != 2 || sizes[0].Size.MaxWidth != 640 || sizes[1].Size.MaxWidth != 1920 {
		t.Fatalf("GetFormatFrameSizes() = %+v", sizes)
	}

	iv0, err := GetFormatFrameInterval(3, 0, PixelFmtYUYV, 640, 480)
	if err != nil || iv0.Interval.Min.Denominator != 30 {
		t.Fatalf("GetFormatFrameInterval(0) = %+v, %v", iv0, err)
	}
	iv1, err := GetFormatFrameInterval(3, 1, PixelFmtYUYV, 640, 480)
	if err != nil || iv1.Interval.Max.Denominator != 15 {
		t.Fatalf("GetFormatFrameInterval(1) = %+v, %v", iv1, err)
	}
}

func TestStreamParamAndPixFormat_MockIoctl(t *testing.T) {
	var storedParm cV4L2StreamParm
	var storedFmt cV4L2Format

	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		switch req {
		case cVIDIOC_S_PARM:
			storedParm = *(*cV4L2StreamParm)(arg)
			return 0
		case cVIDIOC_G_PARM:
			p := (*cV4L2StreamParm)(arg)
			*p = storedParm
			return 0
		case cVIDIOC_S_FMT, cVIDIOC_TRY_FMT:
			f := (*cV4L2Format)(arg)
			if f.Type == BufTypeVideoCapture || f.Type == BufTypeVideoOutput {
				pix := (*cV4L2PixFormat)(unsafe.Pointer(&f.Fmt[0]))
				pix.Bytesperline = pix.Width * 2
				pix.Sizeimage = pix.Bytesperline * pix.Height
				if req == cVIDIOC_S_FMT {
					storedFmt = *f
				}
				return 0
			}
			if f.Type == BufTypeVideoCaptureMPlane || f.Type == BufTypeVideoOutputMPlane {
				mp := (*cV4L2PixFormatMPlane)(unsafe.Pointer(&f.Fmt[0]))
				mp.Num_planes = 2
				mp.Plane_fmt[0] = cV4L2PlanePixFormat{Bytesperline: mp.Width, Sizeimage: mp.Width * mp.Height}
				mp.Plane_fmt[1] = cV4L2PlanePixFormat{Bytesperline: mp.Width, Sizeimage: mp.Width * mp.Height / 2}
				if req == cVIDIOC_S_FMT {
					storedFmt = *f
				}
				return 0
			}
			return sys.EINVAL
		case cVIDIOC_G_FMT:
			f := (*cV4L2Format)(arg)
			*f = storedFmt
			return 0
		default:
			return sys.ENOTTY
		}
	})

	wantParm := StreamParam{
		Type: BufTypeVideoCapture,
		Capture: CaptureParam{
			Capability:   StreamParamTimePerFrame,
			TimePerFrame: Fract{Numerator: 1, Denominator: 30},
		},
	}
	if err := SetStreamParam(3, BufTypeVideoCapture, wantParm); err != nil {
		t.Fatalf("SetStreamParam() = %v", err)
	}
	gotParm, err := GetStreamParam(3, BufTypeVideoCapture)
	if err != nil || gotParm.Capture.TimePerFrame.Denominator != 30 {
		t.Fatalf("GetStreamParam() = %+v, %v", gotParm, err)
	}

	if err := SetPixFormat(3, PixFormat{Width: 640, Height: 480, PixelFormat: PixelFmtYUYV}); err != nil {
		t.Fatalf("SetPixFormat() = %v", err)
	}
	gotFmt, err := GetPixFormat(3)
	if err != nil || gotFmt.SizeImage != 640*480*2 {
		t.Fatalf("GetPixFormat() = %+v, %v", gotFmt, err)
	}
	tried, err := TryPixFormat(3, BufTypeVideoCapture, PixFormat{Width: 320, Height: 240, PixelFormat: PixelFmtYUYV})
	if err != nil || tried.SizeImage != 320*240*2 {
		t.Fatalf("TryPixFormat() = %+v, %v", tried, err)
	}

	mpSet, err := SetPixFormatMPlane(3, BufTypeVideoCaptureMPlane, PixFormatMPlane{
		Width: 1280, Height: 720, PixelFormat: PixelFmtNV12, NumPlanes: 2,
	})
	if err != nil || mpSet.NumPlanes != 2 || mpSet.PlaneFmt[0].SizeImage != 1280*720 {
		t.Fatalf("SetPixFormatMPlane() = %+v, %v", mpSet, err)
	}
	mpGot, err := GetPixFormatMPlane(3, BufTypeVideoCaptureMPlane)
	if err != nil || mpGot.Width != 1280 || mpGot.PlaneFmt[1].SizeImage != 1280*720/2 {
		t.Fatalf("GetPixFormatMPlane() = %+v, %v", mpGot, err)
	}
	mpTry, err := TryPixFormatMPlane(3, BufTypeVideoCaptureMPlane, mpGot)
	if err != nil || mpTry.NumPlanes != 2 {
		t.Fatalf("TryPixFormatMPlane() = %+v, %v", mpTry, err)
	}
}

func TestStreamingAndControlsAndSubdev_MockIoctl(t *testing.T) {
	withMockIoctl(t, func(fd, req uintptr, arg unsafe.Pointer) sys.Errno {
		switch req {
		case cVIDIOC_REQBUFS:
			r := (*cV4L2RequestBuffers)(arg)
			r.Capabilities = cV4L2_BUF_CAP_SUPPORTS_MMAP
			return 0
		case cVIDIOC_CREATE_BUFS:
			cb := (*cV4L2CreateBuffers)(arg)
			cb.Index = 2
			cb.Capabilities = cV4L2_BUF_CAP_SUPPORTS_MMAP
			return 0
		case cVIDIOC_PREPARE_BUF:
			b := (*cV4L2Buffer)(arg)
			b.Flags |= BufFlagPrepared
			return 0
		case cVIDIOC_REMOVE_BUFS:
			return 0
		case cVIDIOC_QUERYBUF, cVIDIOC_QBUF, cVIDIOC_DQBUF:
			b := (*cV4L2Buffer)(arg)
			b.Sequence = 7
			b.Bytesused = 1024
			if IsMultiPlanar(b.Type) && b.Length > 0 {
				planesPtr := *(*unsafe.Pointer)(unsafe.Pointer(&b.M[0]))
				planes := unsafe.Slice((*cV4L2Plane)(planesPtr), b.Length)
				for i := range planes {
					planes[i].Bytesused = 512
					planes[i].Length = 2048
				}
			} else {
				b.Length = 4096
			}
			return 0
		case cVIDIOC_STREAMON, cVIDIOC_STREAMOFF:
			return 0
		case cVIDIOC_EXPBUF:
			eb := (*cV4L2ExportBuffer)(arg)
			eb.Fd = 42
			return 0
		case cVIDIOC_QUERYCTRL:
			qc := (*cV4L2QueryCtrl)(arg)
			if qc.Id == CtrlBrightness || qc.Id == CtrlFlagNextCtrl {
				qc.Id = CtrlBrightness
				qc.Type = uint32(CtrlTypeInt)
				qc.Minimum = 0
				qc.Maximum = 255
				qc.Step = 1
				qc.Value = 128
				copy(qc.Name[:], "Brightness")
				return 0
			}
			return sys.EINVAL
		case cVIDIOC_G_CTRL:
			c := (*cV4L2Control)(arg)
			c.Value = 128
			return 0
		case cVIDIOC_S_CTRL:
			return 0
		case cVIDIOC_SUBDEV_G_FMT, cVIDIOC_SUBDEV_S_FMT:
			f := (*cV4L2SubdevFormat)(arg)
			f.Format.Width = 1920
			f.Format.Height = 1080
			f.Format.Code = 0x300f
			return 0
		case cVIDIOC_SUBDEV_G_SELECTION, cVIDIOC_SUBDEV_S_SELECTION:
			s := (*cV4L2SubdevSelection)(arg)
			s.R = cV4L2Rect{Width: 1920, Height: 1080}
			return 0
		case cVIDIOC_SUBDEV_ENUM_MBUS_CODE:
			e := (*cV4L2SubdevMBusCodeEnum)(arg)
			if e.Index == 0 {
				e.Code = 0x300f
				return 0
			}
			return sys.EINVAL
		case cVIDIOC_SUBDEV_ENUM_FRAME_SIZE:
			e := (*cV4L2SubdevFrameSizeEnum)(arg)
			if e.Index == 0 {
				e.Min_width, e.Max_width = 1920, 1920
				e.Min_height, e.Max_height = 1080, 1080
				return 0
			}
			return sys.EINVAL
		default:
			return sys.ENOTTY
		}
	})

	req, err := RequestBuffersForType(3, BufTypeVideoCapture, IOTypeMMAP, 4)
	if err != nil || req.Count != 4 {
		t.Fatalf("RequestBuffersForType() = %+v, %v", req, err)
	}
	cb, err := CreateBuffers(3, 2, IOTypeMMAP, BufTypeVideoCapture, &PixFormat{Width: 640, Height: 480}, nil)
	if err != nil || cb.Index != 2 {
		t.Fatalf("CreateBuffers() = %+v, %v", cb, err)
	}
	prep, err := PrepareBuffer(3, Buffer{Type: BufTypeVideoCapture, Memory: IOTypeMMAP, Index: 0})
	if err != nil || prep.Flags&BufFlagPrepared == 0 {
		t.Fatalf("PrepareBuffer() = %+v, %v", prep, err)
	}
	if err := RemoveBuffers(3, BufTypeVideoCapture, 2, 2); err != nil {
		t.Fatalf("RemoveBuffers() = %v", err)
	}
	qbuf, err := QueryBuffer(3, BufTypeVideoCapture, IOTypeMMAP, 0, 0)
	if err != nil || qbuf.Length != 4096 {
		t.Fatalf("QueryBuffer() = %+v, %v", qbuf, err)
	}
	if _, err := QueueBufferWith(3, Buffer{Type: BufTypeVideoCapture, Memory: IOTypeMMAP, Index: 0}); err != nil {
		t.Fatalf("QueueBufferWith() = %v", err)
	}
	dq, err := DequeueBufferMPlane(3, IOTypeMMAP, BufTypeVideoCaptureMPlane, 2)
	if err != nil || dq.Sequence != 7 || len(dq.Planes) != 2 || dq.Planes[0].BytesUsed != 512 {
		t.Fatalf("DequeueBufferMPlane() = %+v, %v", dq, err)
	}
	if _, err := DequeueBufferMPlane(3, IOTypeMMAP, BufTypeVideoCapture, 1); err == nil {
		t.Fatalf("expected error for single-planar type in DequeueBufferMPlane")
	}
	if err := StreamOnType(3, BufTypeVideoCapture); err != nil {
		t.Fatalf("StreamOnType() = %v", err)
	}
	if err := StreamOffType(3, BufTypeVideoCapture); err != nil {
		t.Fatalf("StreamOffType() = %v", err)
	}
	expFD, err := ExportBufferPlane(3, BufTypeVideoCapture, 0, 0, sys.O_CLOEXEC)
	if err != nil || expFD != 42 {
		t.Fatalf("ExportBufferPlane() = %d, %v", expFD, err)
	}

	// Controls
	val, err := GetControlValue(3, CtrlBrightness)
	if err != nil || val != 128 {
		t.Fatalf("GetControlValue() = %d, %v", val, err)
	}
	if err := SetControlValue(3, CtrlBrightness, 200); err != nil {
		t.Fatalf("SetControlValue() = %v", err)
	}
	if err := SetControlValue(3, CtrlBrightness, 999); err == nil {
		t.Fatalf("SetControlValue(999) expected out-of-range error")
	}

	// Subdev
	sfmt, err := GetSubdevFormat(3, 0, SubdevFormatActive)
	if err != nil || sfmt.Format.Width != 1920 {
		t.Fatalf("GetSubdevFormat() = %+v, %v", sfmt, err)
	}
	if _, err := SetSubdevFormat(3, sfmt); err != nil {
		t.Fatalf("SetSubdevFormat() = %v", err)
	}
	codes, err := EnumSubdevMBusCodes(3, 0, SubdevFormatActive)
	if err != nil || len(codes) != 1 || codes[0].Code != 0x300f {
		t.Fatalf("EnumSubdevMBusCodes() = %+v, %v", codes, err)
	}
	fsizes, err := EnumSubdevFrameSizes(3, 0, 0x300f, SubdevFormatActive)
	if err != nil || len(fsizes) != 1 || fsizes[0].MaxWidth != 1920 {
		t.Fatalf("EnumSubdevFrameSizes() = %+v, %v", fsizes, err)
	}
	ssel, err := GetSubdevSelection(3, 0, SubdevFormatActive, SelTargetCropBounds)
	if err != nil || ssel.Rect.Width != 1920 {
		t.Fatalf("GetSubdevSelection() = %+v, %v", ssel, err)
	}
	if _, err := SetSubdevSelection(3, ssel); err != nil {
		t.Fatalf("SetSubdevSelection() = %v", err)
	}
}
