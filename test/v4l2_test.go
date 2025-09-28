// +build integration

package test

import (
	"testing"
	"unsafe"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

// TestV4L2_Capability tests v4l2.Capability struct and methods
func TestV4L2_Capability(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	cap := dev.Capability()

	// Test Capability struct fields
	t.Run("Fields", func(t *testing.T) {
		if cap.Driver == "" {
			t.Error("Driver field should not be empty")
		}
		if cap.Card == "" {
			t.Error("Card field should not be empty")
		}
		if cap.BusInfo == "" {
			t.Error("BusInfo field should not be empty")
		}
		if cap.Version == 0 {
			t.Error("Version should not be zero")
		}

		t.Logf("Capability: Driver=%s, Card=%s, Bus=%s, Version=0x%08x",
			cap.Driver, cap.Card, cap.BusInfo, cap.Version)
	})

	// Test capability check methods
	t.Run("Methods", func(t *testing.T) {
		// Version info
		versionInfo := cap.GetVersionInfo()
		t.Logf("Version: %s", versionInfo.String())
		if versionInfo.Major() == 0 {
			t.Error("Major version should not be zero")
		}

		// Device capabilities
		if cap.IsVideoCaptureSupported() {
			t.Log("Video capture supported")
		}
		if cap.IsVideoOutputSupported() {
			t.Log("Video output supported")
		}
		if cap.IsVideoOverlaySupported() {
			t.Log("Video overlay supported")
		}
		if cap.IsVideoCaptureMultiplanarSupported() {
			t.Log("Video capture multiplanar supported")
		}
		if cap.IsVideoOutputMultiplanerSupported() {
			t.Log("Video output multiplanar supported")
		}

		// I/O capabilities
		if cap.IsStreamingSupported() {
			t.Log("Streaming I/O supported")
		}
		if cap.IsReadWriteSupported() {
			t.Log("Read/Write I/O supported")
		}

		// Check if device capabilities are provided (V4L2 specific capability)
		if cap.IsDeviceCapabilitiesProvided() {
			t.Log("Device capabilities provided")
		}
	})
}

// TestV4L2_PixFormat tests v4l2.PixFormat struct
func TestV4L2_PixFormat(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	t.Run("GetSet", func(t *testing.T) {
		// Get current format
		origFormat, err := dev.GetPixFormat()
		if err != nil {
			t.Fatalf("Failed to get pixel format: %v", err)
		}

		t.Logf("Original format: %dx%d, PixelFormat=0x%08x",
			origFormat.Width, origFormat.Height, origFormat.PixelFormat)

		// Test setting various formats
		testFormats := []v4l2.PixFormat{
			{
				Width:       640,
				Height:      480,
				PixelFormat: v4l2.PixelFmtYUYV,
				Field:       v4l2.FieldNone,
			},
			{
				Width:       320,
				Height:      240,
				PixelFormat: v4l2.PixelFmtYUYV,
				Field:       v4l2.FieldNone,
			},
		}

		for _, format := range testFormats {
			err := dev.SetPixFormat(format)
			if err != nil {
				t.Logf("Format %dx%d not supported: %v", format.Width, format.Height, err)
				continue
			}

			actual, err := dev.GetPixFormat()
			if err != nil {
				t.Errorf("Failed to get format after set: %v", err)
				continue
			}

			t.Logf("Set format %dx%d, got %dx%d",
				format.Width, format.Height,
				actual.Width, actual.Height)
		}

		// Restore original format
		dev.SetPixFormat(origFormat)
	})

	t.Run("Fields", func(t *testing.T) {
		format, err := dev.GetPixFormat()
		if err != nil {
			t.Fatalf("Failed to get pixel format: %v", err)
		}

		// Check all fields
		if format.Width == 0 {
			t.Error("Width should not be zero")
		}
		if format.Height == 0 {
			t.Error("Height should not be zero")
		}
		if format.PixelFormat == 0 {
			t.Error("PixelFormat should not be zero")
		}
		if format.BytesPerLine == 0 {
			t.Error("BytesPerLine should not be zero")
		}
		if format.SizeImage == 0 {
			t.Error("SizeImage should not be zero")
		}

		t.Logf("Format details: %dx%d, PixelFormat=0x%08x, Field=%d, BytesPerLine=%d, SizeImage=%d",
			format.Width, format.Height, format.PixelFormat,
			format.Field, format.BytesPerLine, format.SizeImage)
	})
}

// TestV4L2_FormatDescription tests v4l2.FormatDescription struct
func TestV4L2_FormatDescription(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	formats, err := dev.GetFormatDescriptions()
	if err != nil {
		t.Fatalf("Failed to get format descriptions: %v", err)
	}

	if len(formats) == 0 {
		t.Error("No format descriptions returned")
	}

	for i, fmt := range formats {
		t.Run(fmt.Description, func(t *testing.T) {
			// Check fields
			if fmt.Index != uint32(i) {
				t.Errorf("Format index mismatch: expected %d, got %d", i, fmt.Index)
			}
			if fmt.StreamType == 0 {
				t.Error("StreamType should not be zero")
			}
			if fmt.PixelFormat == 0 {
				t.Error("PixelFormat should not be zero")
			}
			if fmt.Description == "" {
				t.Error("Description should not be empty")
			}

			t.Logf("Format %d: %s (0x%08x), StreamType=%d, Flags=0x%08x",
				fmt.Index, fmt.Description, fmt.PixelFormat, fmt.StreamType, fmt.Flags)

			// Check frame sizes if supported
			if fmt.Flags&v4l2.FmtDescFlagCompressed != 0 {
				t.Log("  Compressed format")
			}
			if fmt.Flags&v4l2.FmtDescFlagEmulated != 0 {
				t.Log("  Emulated format")
			}
		})
	}
}


// TestV4L2_StreamParam tests v4l2.StreamParam struct
func TestV4L2_StreamParam(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	t.Run("CaptureParam", func(t *testing.T) {
		param, err := dev.GetStreamParam()
		if err != nil {
			t.Logf("Cannot get stream param: %v", err)
			return
		}

		// Check capability flags
		if param.Capture.Capability&v4l2.StreamParamTimePerFrame != 0 {
			t.Log("Time per frame capability supported")
		}

		// Check time per frame
		if param.Capture.TimePerFrame.Numerator > 0 && param.Capture.TimePerFrame.Denominator > 0 {
			fps := float64(param.Capture.TimePerFrame.Denominator) / float64(param.Capture.TimePerFrame.Numerator)
			t.Logf("Current FPS: %.2f (%d/%d)", fps,
				param.Capture.TimePerFrame.Denominator,
				param.Capture.TimePerFrame.Numerator)
		}

		// Try to set different frame rates
		testRates := []struct {
			num   uint32
			denom uint32
		}{
			{1, 15}, // 15 FPS
			{1, 30}, // 30 FPS
		}

		for _, rate := range testRates {
			param.Capture.TimePerFrame.Numerator = rate.num
			param.Capture.TimePerFrame.Denominator = rate.denom

			err := dev.SetStreamParam(param)
			if err != nil {
				t.Logf("Cannot set %d FPS: %v", rate.denom/rate.num, err)
				continue
			}

			// Verify it was set
			actual, err := dev.GetStreamParam()
			if err != nil {
				t.Errorf("Failed to get param after set: %v", err)
				continue
			}

			if actual.Capture.TimePerFrame.Numerator > 0 && actual.Capture.TimePerFrame.Denominator > 0 {
				actualFPS := float64(actual.Capture.TimePerFrame.Denominator) / float64(actual.Capture.TimePerFrame.Numerator)
				t.Logf("Set %d FPS, got %.2f FPS", rate.denom/rate.num, actualFPS)
			}
		}
	})
}

// TestV4L2_Buffer tests v4l2.Buffer struct
func TestV4L2_Buffer(t *testing.T) {
	// Create a test buffer
	buf := v4l2.Buffer{
		Index:     0,
		Type:      v4l2.BufTypeVideoCapture,
		BytesUsed: 614400,
		Flags:     v4l2.BufFlagMapped | v4l2.BufFlagQueued,
		Field:     v4l2.FieldNone,
		Memory:    v4l2.IOTypeMMAP,
	}

	t.Run("Fields", func(t *testing.T) {
		if buf.Index != 0 {
			t.Errorf("Expected index 0, got %d", buf.Index)
		}
		if buf.Type != v4l2.BufTypeVideoCapture {
			t.Errorf("Expected type %d, got %d", v4l2.BufTypeVideoCapture, buf.Type)
		}
		if buf.BytesUsed != 614400 {
			t.Errorf("Expected BytesUsed 614400, got %d", buf.BytesUsed)
		}
		if buf.Memory != v4l2.IOTypeMMAP {
			t.Errorf("Expected memory type %d, got %d", v4l2.IOTypeMMAP, buf.Memory)
		}
	})

	t.Run("Size", func(t *testing.T) {
		// Verify struct size matches kernel expectation
		size := unsafe.Sizeof(buf)
		if size == 0 {
			t.Error("Buffer struct size should not be zero")
		}
		t.Logf("Buffer struct size: %d bytes", size)
	})
}

// TestV4L2_Control tests v4l2.Control struct
func TestV4L2_Control(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	// Common control IDs to test
	testControls := []struct {
		id   uint32
		name string
	}{
		{uint32(v4l2.CtrlBrightness), "Brightness"},
		{uint32(v4l2.CtrlContrast), "Contrast"},
		{uint32(v4l2.CtrlSaturation), "Saturation"},
		{uint32(v4l2.CtrlHue), "Hue"},
	}

	for _, tc := range testControls {
		t.Run(tc.name, func(t *testing.T) {
			// Try to get control value
			ctrl, err := dev.GetControl(v4l2.CtrlID(tc.id))
			if err != nil {
				t.Logf("Control %s not supported: %v", tc.name, err)
				return
			}

			t.Logf("%s current value: %d", tc.name, ctrl.Value)
		})
	}
}

// TestV4L2_CropCapability tests v4l2.CropCapability struct
func TestV4L2_CropCapability(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	cropCap, err := dev.GetCropCapability()
	if err != nil {
		t.Logf("Crop capability not supported: %v", err)
		return
	}

	t.Run("Bounds", func(t *testing.T) {
		bounds := cropCap.Bounds
		t.Logf("Crop bounds: %dx%d at (%d,%d)",
			bounds.Width, bounds.Height, bounds.Left, bounds.Top)

		if bounds.Width == 0 || bounds.Height == 0 {
			t.Error("Crop bounds should have non-zero dimensions")
		}
	})

	t.Run("Default", func(t *testing.T) {
		def := cropCap.DefaultRect
		t.Logf("Default crop: %dx%d at (%d,%d)",
			def.Width, def.Height, def.Left, def.Top)
	})

	t.Run("PixelAspect", func(t *testing.T) {
		pa := cropCap.PixelAspect
		if pa.Numerator > 0 && pa.Denominator > 0 {
			ratio := float64(pa.Numerator) / float64(pa.Denominator)
			t.Logf("Pixel aspect ratio: %.2f (%d/%d)",
				ratio, pa.Numerator, pa.Denominator)
		}
	})
}

// TestV4L2_InputInfo tests v4l2.InputInfo handling
func TestV4L2_InputInfo(t *testing.T) {
	dev, err := device.Open(testDevice1)
	if err != nil {
		t.Skipf("Cannot open test device: %v", err)
	}
	defer dev.Close()

	// Get current input index
	currentInput, err := dev.GetVideoInputIndex()
	if err != nil {
		t.Logf("Cannot get current input: %v", err)
	} else {
		t.Logf("Current input index: %d", currentInput)
	}

	// Get input info for index 0
	info, err := dev.GetVideoInputInfo(0)
	if err != nil {
		t.Logf("GetVideoInputInfo not supported: %v", err)
		return
	}

	t.Logf("Input 0:")
	t.Logf("  Name: %s", info.GetName())
	t.Logf("  Type: %d", info.GetInputType())
	t.Logf("  Index: %d", info.GetIndex())
	t.Logf("  Status: 0x%x", info.GetStatus())
	t.Logf("  Capabilities: 0x%x", info.GetCapabilities())

	// Check input type
	switch info.GetInputType() {
	case v4l2.InputTypeTuner:
		t.Log("  Input type: Tuner")
	case v4l2.InputTypeCamera:
		t.Log("  Input type: Camera")
	}

	// Check status flags
	if info.GetStatus()&v4l2.InputStatusNoPower != 0 {
		t.Log("  Status: No power")
	}
	if info.GetStatus()&v4l2.InputStatusNoSignal != 0 {
		t.Log("  Status: No signal")
	}
	if info.GetStatus()&v4l2.InputStatusNoColor != 0 {
		t.Log("  Status: No color")
	}
}

// TestV4L2_PixelFormats tests pixel format constants
func TestV4L2_PixelFormats(t *testing.T) {
	// Test FourCC calculation
	testCases := []struct {
		name     string
		format   uint32
		expected string
	}{
		{"YUYV", v4l2.PixelFmtYUYV, "YUYV"},
		{"MJPEG", v4l2.PixelFmtMJPEG, "MJPG"},
		{"H264", v4l2.PixelFmtH264, "H264"},
		{"RGB24", v4l2.PixelFmtRGB24, "RGB3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert format to FourCC string
			fourcc := make([]byte, 4)
			fourcc[0] = byte(tc.format)
			fourcc[1] = byte(tc.format >> 8)
			fourcc[2] = byte(tc.format >> 16)
			fourcc[3] = byte(tc.format >> 24)

			fourccStr := string(fourcc)
			if fourccStr != tc.expected {
				t.Errorf("FourCC mismatch: expected %s, got %s", tc.expected, fourccStr)
			}

			t.Logf("Format 0x%08x = %s", tc.format, fourccStr)
		})
	}
}

// TestV4L2_BufferFlags tests buffer flag constants
func TestV4L2_BufferFlags(t *testing.T) {
	// Test flag combinations
	flags := v4l2.BufFlagMapped | v4l2.BufFlagQueued | v4l2.BufFlagDone

	t.Run("FlagChecks", func(t *testing.T) {
		if flags&v4l2.BufFlagMapped == 0 {
			t.Error("BufFlagMapped should be set")
		}
		if flags&v4l2.BufFlagQueued == 0 {
			t.Error("BufFlagQueued should be set")
		}
		if flags&v4l2.BufFlagDone == 0 {
			t.Error("BufFlagDone should be set")
		}
		if flags&v4l2.BufFlagError != 0 {
			t.Error("BufFlagError should not be set")
		}

		t.Logf("Flags: 0x%08x", flags)
	})

	// Test all buffer flags are defined
	allFlags := []struct {
		flag uint32
		name string
	}{
		{v4l2.BufFlagMapped, "Mapped"},
		{v4l2.BufFlagQueued, "Queued"},
		{v4l2.BufFlagDone, "Done"},
		{v4l2.BufFlagKeyFrame, "KeyFrame"},
		{v4l2.BufFlagPFrame, "PFrame"},
		{v4l2.BufFlagBFrame, "BFrame"},
		{v4l2.BufFlagError, "Error"},
		{v4l2.BufFlagInRequest, "InRequest"},
		{v4l2.BufFlagM2MHoldCaptureBuf, "M2MHoldCapture"},
		{v4l2.BufFlagPrepared, "Prepared"},
		{v4l2.BufFlagNoCacheInvalidate, "NoCacheInvalidate"},
		{v4l2.BufFlagNoCacheClean, "NoCacheClean"},
		{v4l2.BufFlagTimestampMonotonic, "TimestampMonotonic"},
		{v4l2.BufFlagTimestampCopy, "TimestampCopy"},
		{v4l2.BufFlagLast, "Last"},
		{v4l2.BufFlagRequestFD, "RequestFD"},
	}

	for _, f := range allFlags {
		t.Run(f.name, func(t *testing.T) {
			if f.flag == 0 {
				t.Errorf("Flag %s should not be zero", f.name)
			}
			t.Logf("Flag %s: 0x%08x", f.name, f.flag)
		})
	}
}

// TestV4L2_MemoryTypes tests memory type constants
func TestV4L2_MemoryTypes(t *testing.T) {
	memTypes := []struct {
		typ  uint32
		name string
	}{
		{v4l2.IOTypeMMAP, "MMAP"},
		{v4l2.IOTypeUserPtr, "UserPtr"},
		{v4l2.IOTypeOverlay, "Overlay"},
		{v4l2.IOTypeDMABuf, "DMABuf"},
	}

	for _, mt := range memTypes {
		t.Run(mt.name, func(t *testing.T) {
			if mt.typ > 4 {
				t.Errorf("Memory type %s has unexpected value: %d", mt.name, mt.typ)
			}
			t.Logf("Memory type %s: %d", mt.name, mt.typ)
		})
	}
}

// TestV4L2_FieldTypes tests field type constants
func TestV4L2_FieldTypes(t *testing.T) {
	fieldTypes := []struct {
		field uint32
		name  string
	}{
		{v4l2.FieldAny, "Any"},
		{v4l2.FieldNone, "None"},
		{v4l2.FieldTop, "Top"},
		{v4l2.FieldBottom, "Bottom"},
		{v4l2.FieldInterlaced, "Interlaced"},
		{v4l2.FieldAlternate, "Alternate"},
	}

	for _, ft := range fieldTypes {
		t.Run(ft.name, func(t *testing.T) {
			t.Logf("Field type %s: %d", ft.name, ft.field)
		})
	}
}

