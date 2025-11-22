package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	devicePath   = flag.String("d", "/dev/video0", "device path")
	listAll      = flag.Bool("all", false, "list all supported formats")
	listCategory = flag.String("category", "", "list formats by category (rgb, yuv, greyscale, bayer, compressed)")
	showCurrent  = flag.Bool("current", false, "show current format")
	detailed     = flag.Bool("detailed", false, "show detailed format information")
	testFormat   = flag.String("test", "", "test if format is supported (e.g., YUYV, MJPEG, NV12)")
)

type formatInfo struct {
	fourcc      v4l2.FourCCType
	description string
	flags       v4l2.PixFormatFlag
}

func main() {
	flag.Parse()

	if !*listAll && *listCategory == "" && !*showCurrent && *testFormat == "" {
		fmt.Println("V4L2 Pixel Formats Example")
		fmt.Println("==========================")
		fmt.Println()
		fmt.Println("This example demonstrates pixel format enumeration and querying.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  -d <device>         Device path (default: /dev/video0)")
		fmt.Println("  -all                List all supported formats")
		fmt.Println("  -category <name>    List formats by category:")
		fmt.Println("                        rgb, yuv, greyscale, bayer, compressed, jpeg, h264, hevc, mpeg, vp")
		fmt.Println("  -current            Show current format")
		fmt.Println("  -detailed           Show detailed format information")
		fmt.Println("  -test <fourcc>      Test if format is supported (e.g., YUYV, MJPEG, NV12)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  ./pixel_formats -current")
		fmt.Println("  ./pixel_formats -all")
		fmt.Println("  ./pixel_formats -all -detailed")
		fmt.Println("  ./pixel_formats -category yuv")
		fmt.Println("  ./pixel_formats -test YUYV")
		return
	}

	dev, err := device.Open(*devicePath)
	if err != nil {
		log.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	if *showCurrent {
		showCurrentFormat(dev)
		return
	}

	if *testFormat != "" {
		testFormatSupport(dev, *testFormat)
		return
	}

	if *listAll {
		listAllFormats(dev)
		return
	}

	if *listCategory != "" {
		listFormatsByCategory(dev, *listCategory)
		return
	}
}

func showCurrentFormat(dev *device.Device) {
	fmt.Printf("Current Format - Device: %s\n", *devicePath)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	pixFmt, err := v4l2.GetPixFormat(dev.Fd())
	if err != nil {
		log.Fatalf("Failed to get current format: %v", err)
	}

	fmt.Printf("Pixel Format:     %s (FourCC: %s)\n",
		v4l2.PixelFormats[pixFmt.PixelFormat],
		fourccToString(pixFmt.PixelFormat))
	fmt.Printf("Resolution:       %d x %d\n", pixFmt.Width, pixFmt.Height)
	fmt.Printf("Field:            %s\n", v4l2.Fields[pixFmt.Field])
	fmt.Printf("Bytes Per Line:   %d\n", pixFmt.BytesPerLine)
	fmt.Printf("Size Image:       %d bytes\n", pixFmt.SizeImage)
	fmt.Printf("Colorspace:       %s\n", v4l2.Colorspaces[pixFmt.Colorspace])
	fmt.Printf("YCbCr Encoding:   %s\n", v4l2.YCbCrEncodings[pixFmt.YcbcrEnc])
	fmt.Printf("Quantization:     %s\n", v4l2.Quantizations[pixFmt.Quantization])
	fmt.Printf("Transfer Func:    %s\n", v4l2.XferFunctions[pixFmt.XferFunc])
	fmt.Println()

	// Show category
	fmt.Printf("Category:         %s\n", pixFmt.GetCategory())

	// Show BPP if available
	if bpp := pixFmt.GetBitsPerPixel(); bpp > 0 {
		fmt.Printf("Bits Per Pixel:   %d\n", bpp)
	}

	// Show flags
	flags := pixFmt.GetFlags()
	if len(flags) > 0 {
		fmt.Printf("Flags:            %s\n", strings.Join(flags, ", "))
	}
}

func listAllFormats(dev *device.Device) {
	fmt.Printf("Supported Formats - Device: %s\n", *devicePath)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	formats, err := enumerateFormats(dev)
	if err != nil {
		log.Fatalf("Failed to enumerate formats: %v", err)
	}

	if len(formats) == 0 {
		fmt.Println("No formats found")
		return
	}

	if *detailed {
		printDetailedFormats(formats)
	} else {
		printCompactFormats(formats)
	}
}

func listFormatsByCategory(dev *device.Device, category string) {
	fmt.Printf("Formats by Category: %s - Device: %s\n", strings.ToUpper(category), *devicePath)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	formats, err := enumerateFormats(dev)
	if err != nil {
		log.Fatalf("Failed to enumerate formats: %v", err)
	}

	// Filter by category
	var filtered []formatInfo
	categoryLower := strings.ToLower(category)

	for _, fmtInfo := range formats {
		// Create a temp PixFormat to use the category helpers
		pixFmt := v4l2.PixFormat{PixelFormat: fmtInfo.fourcc, Flags: fmtInfo.flags}
		fmtCategory := strings.ToLower(pixFmt.GetCategory())

		match := false
		switch categoryLower {
		case "rgb":
			match = pixFmt.IsRGB()
		case "yuv":
			match = pixFmt.IsYUV()
		case "greyscale", "grayscale", "grey", "gray":
			match = pixFmt.IsGreyscale()
		case "bayer":
			match = pixFmt.IsBayer()
		case "compressed":
			match = pixFmt.IsCompressed()
		case "jpeg":
			match = pixFmt.IsJPEG()
		case "h264", "h.264":
			match = pixFmt.IsH264()
		case "hevc", "h265", "h.265":
			match = pixFmt.IsHEVC()
		case "mpeg":
			match = pixFmt.IsMPEG()
		case "vp":
			match = pixFmt.IsVP()
		default:
			// Fallback to string matching
			match = strings.Contains(fmtCategory, categoryLower)
		}

		if match {
			filtered = append(filtered, fmtInfo)
		}
	}

	if len(filtered) == 0 {
		fmt.Printf("No formats found in category: %s\n", category)
		return
	}

	fmt.Printf("Found %d format(s) in category '%s'\n\n", len(filtered), category)

	if *detailed {
		printDetailedFormats(filtered)
	} else {
		printCompactFormats(filtered)
	}
}

func testFormatSupport(dev *device.Device, fourccStr string) {
	fmt.Printf("Testing Format Support: %s - Device: %s\n", fourccStr, *devicePath)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	formats, err := enumerateFormats(dev)
	if err != nil {
		log.Fatalf("Failed to enumerate formats: %v", err)
	}

	fourccUpper := strings.ToUpper(fourccStr)

	for _, fmtInfo := range formats {
		if fourccToString(fmtInfo.fourcc) == fourccUpper {
			fmt.Printf("✓ Format %s IS SUPPORTED\n", fourccUpper)
			fmt.Println()
			fmt.Printf("Description:  %s\n", fmtInfo.description)
			fmt.Printf("FourCC:       %s (0x%08x)\n", fourccToString(fmtInfo.fourcc), fmtInfo.fourcc)

			// Show category
			pixFmt := v4l2.PixFormat{PixelFormat: fmtInfo.fourcc, Flags: fmtInfo.flags}
			fmt.Printf("Category:     %s\n", pixFmt.GetCategory())

			// Show BPP if available
			if bpp := pixFmt.GetBitsPerPixel(); bpp > 0 {
				fmt.Printf("Bits/Pixel:   %d\n", bpp)
			}

			// Show flags
			if fmtInfo.flags != 0 {
				flags := pixFmt.GetFlags()
				if len(flags) > 0 {
					fmt.Printf("Flags:        %s\n", strings.Join(flags, ", "))
				}
			}

			return
		}
	}

	fmt.Printf("✗ Format %s IS NOT SUPPORTED\n", fourccUpper)
}

func enumerateFormats(dev *device.Device) ([]formatInfo, error) {
	var formats []formatInfo

	for i := uint32(0); ; i++ {
		fmtDesc, err := v4l2.GetFormatDescription(dev.Fd(), i)
		if err != nil {
			break // No more formats
		}

		formats = append(formats, formatInfo{
			fourcc:      fmtDesc.PixelFormat,
			description: fmtDesc.Description,
			flags:       fmtDesc.Flags,
		})
	}

	// Sort by FourCC for consistent display
	sort.Slice(formats, func(i, j int) bool {
		return formats[i].fourcc < formats[j].fourcc
	})

	return formats, nil
}

func printCompactFormats(formats []formatInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FourCC\tCategory\tDescription\tFlags")
	fmt.Fprintln(w, "------\t--------\t-----------\t-----")

	for _, fmtInfo := range formats {
		pixFmt := v4l2.PixFormat{PixelFormat: fmtInfo.fourcc, Flags: fmtInfo.flags}
		fourcc := fourccToString(fmtInfo.fourcc)
		category := pixFmt.GetCategory()
		desc := fmtInfo.description

		flagsStr := ""
		if fmtInfo.flags != 0 {
			flags := pixFmt.GetFlags()
			if len(flags) > 0 {
				flagsStr = strings.Join(flags, ", ")
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", fourcc, category, desc, flagsStr)
	}

	w.Flush()
	fmt.Printf("\nTotal: %d format(s)\n", len(formats))
}

func printDetailedFormats(formats []formatInfo) {
	for i, fmtInfo := range formats {
		if i > 0 {
			fmt.Println()
		}

		pixFmt := v4l2.PixFormat{PixelFormat: fmtInfo.fourcc, Flags: fmtInfo.flags}

		fmt.Printf("Format %d:\n", i+1)
		fmt.Println(strings.Repeat("-", 60))
		fmt.Printf("  FourCC:       %s (0x%08x)\n", fourccToString(fmtInfo.fourcc), fmtInfo.fourcc)
		fmt.Printf("  Description:  %s\n", fmtInfo.description)
		fmt.Printf("  Category:     %s\n", pixFmt.GetCategory())

		if bpp := pixFmt.GetBitsPerPixel(); bpp > 0 {
			fmt.Printf("  Bits/Pixel:   %d\n", bpp)
		}

		if fmtInfo.flags != 0 {
			flags := pixFmt.GetFlags()
			if len(flags) > 0 {
				fmt.Printf("  Flags:        %s\n", strings.Join(flags, ", "))
			}
		}

		// Show format characteristics
		var characteristics []string
		if pixFmt.IsRGB() {
			characteristics = append(characteristics, "RGB")
		}
		if pixFmt.IsYUV() {
			if pixFmt.IsYUVPacked() {
				characteristics = append(characteristics, "YUV Packed")
			}
			if pixFmt.IsYUVPlanar() {
				characteristics = append(characteristics, "YUV Planar")
			}
			if pixFmt.IsYUVSemiPlanar() {
				characteristics = append(characteristics, "YUV Semi-Planar")
			}
		}
		if pixFmt.IsGreyscale() {
			characteristics = append(characteristics, "Greyscale")
		}
		if pixFmt.IsBayer() {
			characteristics = append(characteristics, "Bayer Pattern")
		}
		if pixFmt.IsCompressed() {
			characteristics = append(characteristics, "Compressed")
		}
		if pixFmt.IsEmulated() {
			characteristics = append(characteristics, "Emulated (Software)")
		}

		if len(characteristics) > 0 {
			fmt.Printf("  Type:         %s\n", strings.Join(characteristics, ", "))
		}
	}

	fmt.Printf("\nTotal: %d format(s)\n", len(formats))
}

func fourccToString(fourcc v4l2.FourCCType) string {
	return string([]byte{
		byte(fourcc & 0xff),
		byte((fourcc >> 8) & 0xff),
		byte((fourcc >> 16) & 0xff),
		byte((fourcc >> 24) & 0xff),
	})
}
