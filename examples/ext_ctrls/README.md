# Extended Controls Example

This example demonstrates go4vl's comprehensive support for V4L2 extended controls with automatic memory management and high-level convenience methods.

## Features Demonstrated

1. **High-level convenience methods** - Simple one-liners for common controls
2. **Automatic memory management** - No manual `Free()` calls needed
3. **Atomic multi-control operations** - Set multiple controls in one atomic operation
4. **Control validation** - Test values before applying them
5. **Control enumeration** - List all available extended controls

## Usage

### List all available extended controls

```bash
go run main.go -l
```

This will display all extended controls supported by your device, including their ranges, types, and current values.

### Get current control values

```bash
go run main.go -d /dev/video0
```

Shows current brightness, contrast, saturation, and hue values.

### Set individual controls

```bash
# Set brightness to 150
go run main.go -brightness 150

# Set multiple controls at once
go run main.go -brightness 128 -contrast 100 -saturation 64

# Use a different device
go run main.go -d /dev/video1 -brightness 200
```

### Interactive mode

```bash
go run main.go -i
```

Demonstrates:
- Atomic multi-control operations
- Control value validation before applying
- Proper error handling
- Restoring original values

## Code Examples

### High-level API (Simple)

The easiest way to control brightness, contrast, saturation, and hue:

```go
dev, _ := device.Open("/dev/video0")
defer dev.Close()

// Get current brightness
brightness, err := dev.GetBrightness()
if err == nil {
    fmt.Printf("Current brightness: %d\n", brightness)
}

// Set brightness
err = dev.SetBrightness(128)
if err != nil {
    log.Printf("Failed to set brightness: %v", err)
}
```

### Mid-level API (Atomic operations)

For setting multiple controls atomically:

```go
dev, _ := device.Open("/dev/video0")
defer dev.Close()

// Set multiple controls in one atomic operation
ctrls := v4l2.NewExtControls()
ctrls.AddValue(v4l2.CtrlBrightness, 128)
ctrls.AddValue(v4l2.CtrlContrast, 100)
ctrls.AddValue(v4l2.CtrlSaturation, 64)

err := dev.SetExtControls(ctrls)  // All or nothing - atomic
if err != nil {
    log.Printf("Failed to set controls: %v", err)
}
```

### Control Validation

Test values before applying them:

```go
// Create controls to test
testCtrls := v4l2.NewExtControls()
testCtrls.AddValue(v4l2.CtrlBrightness, 200)
testCtrls.AddValue(v4l2.CtrlContrast, 150)

// Validate without applying
err := dev.TryExtControls(testCtrls)
if err != nil {
    fmt.Printf("Values would be rejected: %v\n", err)
} else {
    // Values are valid, apply them
    applyCtrls := v4l2.NewExtControls()
    applyCtrls.AddValue(v4l2.CtrlBrightness, 200)
    applyCtrls.AddValue(v4l2.CtrlContrast, 150)
    dev.SetExtControls(applyCtrls)
}
```

### Low-level API (Full control)

For advanced use cases with complete control:

```go
dev, _ := device.Open("/dev/video0")
defer dev.Close()

// Query all available extended controls
ctrls, err := v4l2.QueryAllExtControls(dev.Fd())
if err != nil {
    log.Fatalf("Failed to query controls: %v", err)
}

for _, ctrl := range ctrls {
    fmt.Printf("Control: %s (ID: 0x%08x)\n", ctrl.Name, ctrl.ID)
    fmt.Printf("  Range: [%d - %d], Default: %d\n",
        ctrl.Minimum, ctrl.Maximum, ctrl.Default)

    if ctrl.IsReadOnly() {
        fmt.Println("  [READ-ONLY]")
    }
}
```

## Key Improvements

### Before (Old API - manual memory management)

```go
ctrls := v4l2.NewExtControls()
defer ctrls.Free()  // Easy to forget!

ctrl := v4l2.NewExtControlWithValue(v4l2.CtrlBrightness, 128)
defer ctrl.Free()   // More cleanup to remember

ctrls.Add(ctrl)
v4l2.SetExtControls(fd, ctrls)
```

### After (New API - automatic memory management)

```go
ctrls := v4l2.NewExtControls()
ctrls.AddValue(v4l2.CtrlBrightness, 128)
v4l2.SetExtControls(fd, ctrls)  // Memory automatically managed!
```

Or even simpler:

```go
dev.SetBrightness(128)  // One line!
```

## Control Types

Extended controls support various data types:

- **int32** - Standard integer controls (brightness, contrast, etc.)
- **int64** - 64-bit values for large ranges
- **string** - Text controls (device names, etc.)
- **[]byte** - Compound controls (codec parameters, structures)

Example with different types:

```go
ctrls := v4l2.NewExtControls()

// 32-bit integer
ctrls.AddValue(v4l2.CtrlBrightness, 128)

// 64-bit integer
ctrls.AddValue64(someCtrlID, 1000000000)

// String
ctrls.AddString(nameCtrlID, "MyCamera")

// Compound data (raw bytes)
ctrls.AddCompound(h264SpsID, spsBytes)

dev.SetExtControls(ctrls)
```

## Type-Safe Codec Controls

For stateless codec controls (H.264, MPEG2, VP8, FWHT), go4vl provides type-safe helper methods that work with Go structs instead of raw byte slices:

### H.264 Codec Controls

```go
dev, _ := device.Open("/dev/video0")
defer dev.Close()

ctrls := v4l2.NewExtControls()

// Add H.264 SPS (Sequence Parameter Set) using type-safe API
sps := &v4l2.ControlH264SPS{
    ProfileIDC:              100, // High Profile
    LevelIDC:                51,  // Level 5.1
    ChromaFormatIDC:         1,   // 4:2:0
    Log2MaxFrameNumMinus4:   0,
    PicOrderCntType:         0,
    Log2MaxPicOrderCntLsbMinus4: 2,
    MaxNumRefFrames:         4,
    PicWidthInMbsMinus1:     119, // 1920 pixels
    PicHeightInMapUnitsMinus1: 67, // 1080 pixels
}
ctrls.AddH264SPS(sps)

// Add H.264 PPS (Picture Parameter Set)
pps := &v4l2.ControlH264PPS{
    PicParameterSetID: 0,
    SeqParameterSetID: 0,
    NumRefIndexL0DefaultActiveMinus1: 3,
    NumRefIndexL1DefaultActiveMinus1: 0,
}
ctrls.AddH264PPS(pps)

// Add H.264 Scaling Matrix (optional)
matrix := &v4l2.ControlH264ScalingMatrix{}
// Fill in scaling lists...
ctrls.AddH264ScalingMatrix(matrix)

// Set all controls atomically
err := dev.SetExtControls(ctrls)
if err != nil {
    log.Fatalf("Failed to set H.264 controls: %v", err)
}

// Read controls back
readCtrls := v4l2.NewExtControls()
readCtrls.Add(v4l2.NewExtControl(v4l2.CtrlH264SPS))
dev.GetExtControls(readCtrls)

retrievedSPS, err := readCtrls.GetControls()[0].GetH264SPS()
if err == nil {
    fmt.Printf("Profile: %d, Level: %d\n", retrievedSPS.ProfileIDC, retrievedSPS.LevelIDC)
}
```

### MPEG2 Codec Controls

```go
ctrls := v4l2.NewExtControls()

// MPEG2 Sequence Header
seq := &v4l2.ControlMPEG2Sequence{
    HorizontalSize: 1920,
    VerticalSize:   1080,
    VBVBufferSize:  224,
    ProfileAndLevelIndication: 0x85, // Main Profile @ Main Level
    ChromaFormat:   1,                // 4:2:0
}
ctrls.AddMPEG2Sequence(seq)

// MPEG2 Picture Header
pic := &v4l2.ControlMPEG2Picture{
    BackwardRefTimestamp: 0,
    ForwardRefTimestamp:  0,
    PictureCodingType:    1, // I-frame
    FCode:                [[2][2]uint8{{15, 15}, {15, 15}}],
}
ctrls.AddMPEG2Picture(pic)

// MPEG2 Quantization Matrices
quant := &v4l2.ControlMPEG2Quantization{}
// Fill quantization matrices...
ctrls.AddMPEG2Quantization(quant)

dev.SetExtControls(ctrls)
```

### VP8 Codec Controls

```go
ctrls := v4l2.NewExtControls()

// VP8 Frame Parameters
frame := &v4l2.ControlVP8Frame{
    Width:         1920,
    Height:        1080,
    Version:       0,
    HorizontalScale: 0,
    VerticalScale:   0,
    Flags:          0,
}
ctrls.AddVP8Frame(frame)

dev.SetExtControls(ctrls)
```

### FWHT Codec Controls

```go
ctrls := v4l2.NewExtControls()

// FWHT (Fast Walsh Hadamard Transform) Parameters
params := &v4l2.ControlFWHTParams{
    Width:   1920,
    Height:  1080,
    Version: 3,
    Flags:   0,
}
ctrls.AddFWHTParams(params)

dev.SetExtControls(ctrls)
```

### Available Type-Safe Codec Helpers

**H.264 Stateless Codec:**
- `AddH264SPS()` / `GetH264SPS()` - Sequence Parameter Set
- `AddH264PPS()` / `GetH264PPS()` - Picture Parameter Set
- `AddH264ScalingMatrix()` / `GetH264ScalingMatrix()` - Scaling matrices
- `AddH264SliceParams()` / `GetH264SliceParams()` - Slice parameters
- `AddH264DecodeParams()` / `GetH264DecodeParams()` - Decode parameters
- `AddH264PredWeights()` / `GetH264PredWeights()` - Prediction weights

**MPEG2 Stateless Codec:**
- `AddMPEG2Sequence()` / `GetMPEG2Sequence()` - Sequence header
- `AddMPEG2Picture()` / `GetMPEG2Picture()` - Picture header
- `AddMPEG2Quantization()` / `GetMPEG2Quantization()` - Quantization matrices

**VP8 Stateless Codec:**
- `AddVP8Frame()` / `GetVP8Frame()` - Frame parameters

**FWHT Stateless Codec:**
- `AddFWHTParams()` / `GetFWHTParams()` - FWHT parameters

### Benefits of Type-Safe Codec API

1. **Type Safety**: Compile-time checking of field names and types
2. **Self-Documenting**: Go struct fields describe the codec parameters
3. **IDE Support**: Auto-completion and inline documentation
4. **Less Error-Prone**: No manual byte manipulation required
5. **Automatic Memory Management**: No need to manually marshal/unmarshal data

## Notes

- The example gracefully handles devices that don't support certain controls
- All memory management is automatic - no `Free()` calls needed
- Controls are set atomically - either all succeed or all fail
- The device must support extended controls (most modern webcams do)
- Some controls may be read-only depending on the device

## Typical Devices with Extended Controls

- USB webcams (Logitech, Microsoft, etc.)
- Video capture cards
- Virtual video devices (v4l2loopback)
- Hardware video encoders/decoders
- Professional broadcast equipment

## See Also

- [V4L2 Extended Controls Documentation](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/extended-controls.html)
- [V4L2 Control API](https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-ext-ctrls.html)
