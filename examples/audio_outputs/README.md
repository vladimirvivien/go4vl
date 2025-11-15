# Audio Output Enumeration Example

This example demonstrates how to enumerate and select audio outputs on V4L2 devices using go4vl.

## Features

- List all available audio outputs
- Display audio output capabilities (stereo, AVL)
- Show current audio output
- Select a specific audio output

## Usage

```bash
# List audio outputs for default device
go run main.go

# List audio outputs for specific device
go run main.go -d /dev/video1

# Select a specific audio output
go run main.go -s 1
```

## Audio Output Capabilities

Audio outputs can have the following capabilities:

- **Stereo**: The audio output supports stereo (left/right channels)
- **AVL**: Automatic Volume Level control is available

## Common Devices with Audio Outputs

- Video output devices
- TV tuner cards with audio output
- Professional video equipment with audio routing

## Example Output

```
Device: /dev/video0
Driver: vivid
Card: VIVID Virtual Device

Current audio output: [0] Line-Out

Available audio outputs (2):
================================================================================
[0] Line-Out ** ACTIVE **
    Capability: 0x00000001
    Mode:       0x00000000
    ✓ Stereo

[1] Speaker
    Capability: 0x00000001
    Mode:       0x00000000
    ✓ Stereo

Current Audio Output Details:
================================================================================
Index:      0
Name:       Line-Out
Capability: 0x00000001
Mode:       0x00000000
Stereo:     true
AVL:        false

Tip: Use -s <index> to select a different audio output
```

## Notes

- Audio outputs are less common than audio inputs in V4L2 devices
- Most capture devices do not have audio outputs
- Video output devices and TV tuners are more likely to have audio outputs
- Some drivers may not fully implement audio output selection
