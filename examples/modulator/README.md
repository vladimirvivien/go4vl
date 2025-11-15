# Modulator Example

This example demonstrates how to enumerate and control modulators on V4L2 devices.

## Overview

Modulators are used to transmit radio frequency signals. Common devices with modulators include:
- **RF modulators** for video output (TV channel modulation)
- **TV transmitters**
- **FM/AM radio transmitters**
- **Software Defined Radio (SDR) devices** with transmission capability

**⚠️ Warning**: Operating radio transmitters may require licensing in your jurisdiction. Always comply with local regulations regarding radio frequency transmission.

## Features Demonstrated

- Enumerate all available modulators
- Display modulator capabilities (stereo, RDS, etc.)
- Query current transmission frequency
- Set modulator frequency
- Enumerate frequency bands (if supported)
- Display transmission subchannel configuration

## Usage

```bash
# Build the example
go build

# List all modulators on the default device
./modulator

# Use a specific device
./modulator -d /dev/video1

# Show frequency bands
./modulator -b

# Set transmission frequency (example for UHF TV channel)
./modulator -f 7584  # 474 MHz (example)

# Combine options
./modulator -d /dev/video0 -b -f 7584
```

## Command-line Options

- `-d <device>` - Device path (default: `/dev/video0`)
- `-f <frequency>` - Set transmission frequency in device-specific units
- `-b` - Show frequency bands

## Understanding Frequency Units

V4L2 uses two different frequency unit systems depending on the modulator capability:

### Low Frequency Units (TunerCapLow)
- **Unit**: 1/16000 kHz = 62.5 Hz
- **Common for**: FM/AM radio modulators
- **Example**: To transmit on 100.5 MHz FM:
  - 100.5 MHz = 100,500 kHz
  - Frequency units = 100,500 × 16 = 1,608,000

### High Frequency Units (Normal)
- **Unit**: 1/16 MHz = 62.5 kHz
- **Common for**: TV modulators
- **Example**: To transmit on 474 MHz (UHF TV):
  - Frequency units = 474 × 16 = 7,584

## Example Output

```
Device: /dev/video0
Driver: example-modulator
Card: Example RF Modulator

Available modulators (1):
================================================================================
[0] RF Modulator
    Type:       Analog TV
    Capability: 0x00001012
    Range:      880 - 13760 (units: 62.5 kHz)
                55.0 - 860.0 MHz
    TxSubchans: 0x1 (Mono)
    Features:   None

Current Frequency (Modulator 0):
================================================================================
Modulator:  0
Type:       Analog TV
Frequency:  7584 (474.0 MHz)

Tip: Use -f <frequency> to set modulator frequency
     Units are 1/16 MHz (62.5 kHz)
Tip: Use -b to show available frequency bands
```

## Modulator Capabilities

The example displays the following modulator capabilities:

- **Stereo**: Modulator supports stereo audio transmission
- **RDS**: Radio Data System support (FM transmission)
- **Freq Bands**: Multiple frequency bands support

## Frequency Bands

Some modulators support multiple frequency bands with different characteristics. Use the `-b` flag to enumerate them:

```bash
./modulator -b
```

Example output:
```
Frequency Bands (Modulator 0):
================================================================================
Band 0:
  Range:      880 - 1520 (55.0 - 95.0 MHz)
  Capability: 0x00001002
  Modulation: 0x4 (FM)

Band 1:
  Range:      2880 - 4320 (180.0 - 270.0 MHz)
  Capability: 0x00001002
  Modulation: 0x2 (VSB)
```

## Common Use Cases

### RF Video Modulator

RF modulators are commonly used to transmit composite video on a TV channel:

```bash
# Set modulator to channel 3 (typically 61.25 MHz for NTSC)
# 61.25 MHz = 61250 kHz × 16 = 980,000 units (if TunerCapLow)
# OR 61.25 MHz ÷ 62.5 kHz = 980 units (if normal)
./modulator -d /dev/video0 -f 980
```

### FM Radio Transmission

```bash
# Transmit on 88.5 MHz FM (where legal)
# 88.5 MHz = 88500 kHz × 16 = 1,416,000 units (if TunerCapLow)
./modulator -d /dev/radio0 -f 1416000
```

## Regulatory Compliance

**IMPORTANT**:
- Radio frequency transmission is regulated in most countries
- Operating a transmitter without proper authorization may be illegal
- Always check your local regulations before transmitting
- This example is for educational purposes and authorized use only
- Unlicensed low-power devices may have specific frequency and power restrictions

## Notes

- Not all V4L2 devices have modulators (modulators are much less common than tuners)
- Frequency units vary by device (check `IsLowFreq()`)
- Some devices may restrict frequency setting for regulatory compliance
- Transmission subchannels (`TxSubchans`) control audio mode (mono/stereo/RDS)

## See Also

- [V4L2 Modulator API Documentation](https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-modulator.html)
- [V4L2 Frequency API Documentation](https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-frequency.html)
- Related example: `examples/tuner` - For RF reception
