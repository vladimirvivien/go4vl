# Tuner Example

This example demonstrates how to enumerate and control tuners on V4L2 devices.

## Overview

Tuners are used to receive radio frequency signals. Common devices with tuners include:
- **TV tuner cards** (analog and digital TV)
- **FM/AM radio receivers**
- **Software Defined Radio (SDR) devices**
- **RF receivers**

## Features Demonstrated

- Enumerate all available tuners
- Display tuner capabilities (stereo, RDS, hardware seek, etc.)
- Query current frequency
- Set tuner frequency
- Enumerate frequency bands (if supported)
- Display signal strength and AFC (Automatic Frequency Control)

## Usage

```bash
# Build the example
go build

# List all tuners on the default device
./tuner

# Use a specific device
./tuner -d /dev/video1

# Show frequency bands
./tuner -b

# Set frequency (example for FM radio at 100.5 MHz)
# For tuners with TunerCapLow: 100.5 MHz = 100500 kHz = 100500 * 16 = 1,608,000 units
./tuner -f 1608000

# Combine options
./tuner -d /dev/radio0 -b -f 1608000
```

## Command-line Options

- `-d <device>` - Device path (default: `/dev/video0`)
- `-f <frequency>` - Set frequency in device-specific units
- `-b` - Show frequency bands

## Understanding Frequency Units

V4L2 uses two different frequency unit systems depending on the tuner capability:

### Low Frequency Units (TunerCapLow)
- **Unit**: 1/16000 kHz = 62.5 Hz
- **Common for**: FM/AM radio tuners
- **Example**: To tune to 100.5 MHz FM radio:
  - 100.5 MHz = 100,500 kHz
  - Frequency units = 100,500 × 16 = 1,608,000

### High Frequency Units (Normal)
- **Unit**: 1/16 MHz = 62.5 kHz
- **Common for**: TV tuners
- **Example**: To tune to 474 MHz (TV channel):
  - Frequency units = 474 × 16 = 7,584

## Example Output

```
Device: /dev/radio0
Driver: radio-si470x
Card: Silicon Labs Si470x FM Radio Receiver

Available tuners (1):
================================================================================
[0] FM Radio
    Type:       Radio
    Capability: 0x00001091
    Range:      1408000 - 1728000 (units: 62.5 Hz)
                88.000 - 108.000 MHz
    Signal:     32768 / 65535 (50.0%)
    Audio Mode: Stereo
    RxSubchans: 0x2 (Stereo)
    Features:   Stereo, RDS, HW Seek

Current Frequency (Tuner 0):
================================================================================
Tuner:      0
Type:       Radio
Frequency:  1608000 (100.500 MHz)

Tip: Use -f <frequency> to tune to a specific frequency
     Units are 1/16000 kHz (62.5 Hz)
     Example: -f 1608000 for 100.5 MHz (FM radio)
Tip: Use -b to show available frequency bands
```

## Tuner Capabilities

The example displays the following tuner capabilities:

- **Stereo**: Tuner supports stereo audio reception
- **RDS**: Radio Data System support (FM radio)
- **HW Seek**: Hardware-assisted frequency scanning
- **Freq Bands**: Multiple frequency bands support

## Frequency Bands

Some tuners support multiple frequency bands with different characteristics. Use the `-b` flag to enumerate them:

```bash
./tuner -b
```

Example output:
```
Frequency Bands (Tuner 0):
================================================================================
Band 0:
  Range:      1408000 - 1728000 (88.000 - 108.000 MHz)
  Capability: 0x00001091
  Modulation: 0x4 (FM)
```

## Common Use Cases

### FM Radio Tuning

```bash
# Tune to your favorite FM station (e.g., 95.5 MHz)
# Calculation: 95.5 MHz = 95500 kHz × 16 = 1,528,000 units
./tuner -d /dev/radio0 -f 1528000
```

### Checking Signal Strength

```bash
# View current signal strength
./tuner -d /dev/radio0
```

The signal strength is displayed as a value from 0 to 65535, where:
- 0 = No signal
- 65535 = Maximum signal

## Notes

- Not all V4L2 devices have tuners
- Frequency units vary by device (check `IsLowFreq()`)
- Some devices may restrict frequency setting for regulatory reasons
- RDS data reception requires additional V4L2 event handling (not shown in this example)

## See Also

- [V4L2 Tuner API Documentation](https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-tuner.html)
- [V4L2 Frequency API Documentation](https://linuxtv.org/downloads/v4l-dvb-apis/userspace-api/v4l/vidioc-g-frequency.html)
- Related example: `examples/modulator` - For RF transmission
