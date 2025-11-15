# Audio Input Enumeration Example

This example demonstrates how to enumerate and select audio inputs on V4L2 devices using go4vl.

## Features

- List all available audio inputs
- Display audio input capabilities (stereo, AVL)
- Show current audio input
- Select a specific audio input

## Usage

```bash
# List audio inputs for default device
go run main.go

# List audio inputs for specific device
go run main.go -d /dev/video1

# Select a specific audio input
go run main.go -s 1
```

## Audio Input Capabilities

Audio inputs can have the following capabilities:

- **Stereo**: The audio input supports stereo (left/right channels)
- **AVL**: Automatic Volume Level control is available

## Common Devices with Audio Inputs

- TV tuner cards (composite audio, tuner audio, line-in)
- Webcams with built-in microphones
- Video capture cards with audio inputs
- Professional video equipment

## Example Output

```
Device: /dev/video0
Driver: uvcvideo
Card: HD Webcam

Current audio input: [0] Microphone

Available audio inputs (1):
================================================================================
[0] Microphone ** ACTIVE **
    Capability: 0x00000001
    Mode:       0x00000000
    ✓ Stereo

Current Audio Input Details:
================================================================================
Index:      0
Name:       Microphone
Capability: 0x00000001
Mode:       0x00000000
Stereo:     true
AVL:        false

Tip: Use -s <index> to select a different audio input
```

## Notes

- Not all V4L2 devices support audio inputs
- Most modern webcams have a single built-in microphone
- TV tuner cards typically have multiple audio inputs
- Some drivers may not fully implement audio input selection
