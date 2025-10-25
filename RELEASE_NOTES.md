# Release Notes: v0.3.0

## Breaking Changes

### Removed Vendored Linux Headers (#88)

The vendored Linux kernel headers (`include/linux/` directory) have been removed from the repository. **go4vl now requires system V4L2 kernel headers to be installed.**

**Migration Required:**

Users must install V4L2 kernel headers on their system:

```bash
# Ubuntu/Debian
sudo apt install linux-libc-dev

# Fedora/RHEL/CentOS
sudo dnf install kernel-headers

# Arch Linux
sudo pacman -S linux-headers

# Alpine Linux
apk add linux-headers
```

**Why this change?**

- Follows Linux kernel documentation best practices for user-space applications
- Reduces repository size by ~5,500 lines of vendored code
- Simplifies maintenance and ensures headers stay up-to-date with the system
- Provides flexibility for users to choose kernel header versions

**Using Custom Headers:**

If you need specific kernel header versions or are cross-compiling, use the `CGO_CFLAGS` environment variable:

```bash
CGO_CFLAGS="-I/path/to/custom/headers" go build ./v4l2
```

See [docs/BUILD.md](./docs/BUILD.md) for comprehensive build instructions.

## Build System Enhancements (#88)

### New Build Documentation

- Added comprehensive [docs/BUILD.md](./docs/BUILD.md) with detailed build instructions
- Covers prerequisite installation for multiple Linux distributions
- Includes cross-compilation guides (GCC, Zig, Docker)
- Platform-specific instructions for Raspberry Pi, WSL2, and more
- Troubleshooting section for common build issues

### Updated Build Instructions

- Enhanced [README.md](./README.md) with new "Building" section
- Updated [examples/README.md](./examples/README.md) with header installation steps
- Documented `CGO_CFLAGS` override mechanism for custom headers
- Added examples for cross-compilation with custom headers

### CGo Configuration

- Centralized CGo directives in `v4l2/cgo.go`
- Default configuration uses system headers from `/usr/include`
- Simplified CGo setup with clear documentation on overrides
- Better support for cross-compilation scenarios

## Cleanup (#86)

### Removed Unsupported Code

- Removed `multipass` directory and related code
- Cleaned up legacy/deprecated packages
- Repository structure simplified for better maintainability

## Documentation Improvements

### Build System Documentation

- New comprehensive build guide covering all major Linux distributions
- Cross-compilation instructions for ARM (32-bit and 64-bit)
- Docker-based build examples
- Static linking and debug build options
- Example build scripts for automation

### Prerequisites Documentation

- Clear system requirements for building go4vl
- Package installation commands for all major distros
- User permission setup instructions
- V4L2 utilities installation guidance

## Files Changed

### Added
- `docs/BUILD.md` - Comprehensive build documentation (519 lines)

### Modified
- `README.md` - Added Building section with quick start commands
- `examples/README.md` - Updated with header installation instructions
- `v4l2/cgo.go` - Centralized CGo configuration with documentation
- `v4l2/doc.go` - Updated package documentation

### Removed
- `include/linux/videodev2.h` - 2,770 lines (now uses system headers)
- `include/linux/v4l2-controls.h` - 2,486 lines (now uses system headers)
- `include/linux/v4l2-common.h` - 108 lines (now uses system headers)
- `include/README.md` - Removed vendored headers documentation
- `multipass/` - Removed unsupported package

## Testing

All existing tests continue to pass with system headers:
- ✅ v4l2 package tests
- ✅ device package tests
- ✅ CGo compilation tests
- ✅ Example builds (snapshot, cgo_types, etc.)

Verified `CGO_CFLAGS` override mechanism works correctly for:
- Custom header paths
- Cross-compilation scenarios
- Target-specific sysroots

## Upgrade Guide

### For Users Building from Source

1. **Install V4L2 kernel headers** (if not already installed):
   ```bash
   # Ubuntu/Debian
   sudo apt install linux-libc-dev

   # Fedora/RHEL
   sudo dnf install kernel-headers

   # Arch Linux
   sudo pacman -S linux-headers
   ```

2. **Build normally** - no other changes required:
   ```bash
   go build ./v4l2
   ```

### For Cross-Compilation Users

Update your build scripts to include target headers:

```bash
# Example: Cross-compile for ARM64
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=arm64 \
CC="zig cc -target aarch64-linux-musl" \
CGO_CFLAGS="-I/path/to/arm64/sysroot/usr/include" \
go build ./v4l2
```

See [docs/BUILD.md](./docs/BUILD.md) for detailed cross-compilation instructions.

### For CI/CD Pipelines

Ensure your build containers/environments include the `linux-libc-dev` (or equivalent) package:

```dockerfile
FROM golang:1.21-bullseye

RUN apt-get update && apt-get install -y \
    linux-libc-dev \
    build-essential \
    && rm -rf /var/lib/apt/lists/*
```

## Known Issues

- Kernel headers from `/usr/src/linux-headers-*/include` may cause type conflicts. Use `/usr/include` (from `linux-libc-dev`) or kernel source tarballs instead.

## Contributors

- @vladimirvivien

## Full Changelog

**Merged Pull Requests:**
- #88: Issue #87 - Remove Linux headers, enhance build support
- #86: Issue #84 - Remove multipass directory

**Commits:** +620 additions, −5,586 deletions across 18 files

Compare: [v0.2.0...v0.3.0](https://github.com/vladimirvivien/go4vl/compare/v0.2.0...v0.3.0)
