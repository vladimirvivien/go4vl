package v4l2

import (
	"strings"
	"testing"
)

// TestCropCapability_StructFields tests that all CropCapability struct fields are accessible
func TestCropCapability_StructFields(t *testing.T) {
	cap := CropCapability{
		StreamType: BufTypeVideoCapture,
		Bounds: Rect{
			Left:   0,
			Top:    0,
			Width:  1920,
			Height: 1080,
		},
		DefaultRect: Rect{
			Left:   0,
			Top:    0,
			Width:  1920,
			Height: 1080,
		},
		PixelAspect: Fract{
			Numerator:   1,
			Denominator: 1,
		},
	}

	if cap.StreamType != BufTypeVideoCapture {
		t.Errorf("StreamType: expected %d, got %d", BufTypeVideoCapture, cap.StreamType)
	}
	if cap.Bounds.Width != 1920 {
		t.Errorf("Bounds.Width: expected 1920, got %d", cap.Bounds.Width)
	}
	if cap.DefaultRect.Height != 1080 {
		t.Errorf("DefaultRect.Height: expected 1080, got %d", cap.DefaultRect.Height)
	}
	if cap.PixelAspect.Numerator != 1 {
		t.Errorf("PixelAspect.Numerator: expected 1, got %d", cap.PixelAspect.Numerator)
	}
}

// TestCropCapability_String tests the String() method
func TestCropCapability_String(t *testing.T) {
	tests := []struct {
		name     string
		cap      CropCapability
		contains []string
	}{
		{
			name: "1080p full frame",
			cap: CropCapability{
				DefaultRect: Rect{
					Top:    0,
					Left:   0,
					Width:  1920,
					Height: 1080,
				},
				Bounds: Rect{
					Top:    0,
					Left:   0,
					Width:  1920,
					Height: 1080,
				},
				PixelAspect: Fract{
					Numerator:   1,
					Denominator: 1,
				},
			},
			contains: []string{"1920", "1080", "1:1"},
		},
		{
			name: "Cropped region",
			cap: CropCapability{
				DefaultRect: Rect{
					Top:    100,
					Left:   200,
					Width:  640,
					Height: 480,
				},
				Bounds: Rect{
					Top:    0,
					Left:   0,
					Width:  1920,
					Height: 1080,
				},
				PixelAspect: Fract{
					Numerator:   16,
					Denominator: 9,
				},
			},
			contains: []string{"640", "480", "100", "200", "1920", "1080", "16:9"},
		},
		{
			name: "4:3 aspect ratio",
			cap: CropCapability{
				DefaultRect: Rect{
					Width:  640,
					Height: 480,
				},
				Bounds: Rect{
					Width:  640,
					Height: 480,
				},
				PixelAspect: Fract{
					Numerator:   4,
					Denominator: 3,
				},
			},
			contains: []string{"640", "480", "4:3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cap.String()
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected string to contain %q, got: %s", substr, result)
				}
			}
		})
	}
}

// TestRect_StructFields tests Rect field accessibility
func TestRect_StructFields(t *testing.T) {
	r := Rect{
		Left:   100,
		Top:    50,
		Width:  1280,
		Height: 720,
	}

	if r.Left != 100 {
		t.Errorf("Left: expected 100, got %d", r.Left)
	}
	if r.Top != 50 {
		t.Errorf("Top: expected 50, got %d", r.Top)
	}
	if r.Width != 1280 {
		t.Errorf("Width: expected 1280, got %d", r.Width)
	}
	if r.Height != 720 {
		t.Errorf("Height: expected 720, got %d", r.Height)
	}
}

// TestRect_NegativeOffsets tests that Rect can handle negative offsets
func TestRect_NegativeOffsets(t *testing.T) {
	r := Rect{
		Left:   -100,
		Top:    -50,
		Width:  640,
		Height: 480,
	}

	if r.Left != -100 {
		t.Errorf("Left: expected -100, got %d", r.Left)
	}
	if r.Top != -50 {
		t.Errorf("Top: expected -50, got %d", r.Top)
	}
}

// TestRect_ZeroValues tests Rect with zero values
func TestRect_ZeroValues(t *testing.T) {
	r := Rect{}

	if r.Left != 0 || r.Top != 0 || r.Width != 0 || r.Height != 0 {
		t.Error("Zero-initialized Rect should have all fields at 0")
	}
}

// TestRect_CommonResolutions tests typical video crop rectangles
func TestRect_CommonResolutions(t *testing.T) {
	tests := []struct {
		name   string
		rect   Rect
		width  uint32
		height uint32
	}{
		{
			name:   "1080p full frame",
			rect:   Rect{Width: 1920, Height: 1080},
			width:  1920,
			height: 1080,
		},
		{
			name:   "720p full frame",
			rect:   Rect{Width: 1280, Height: 720},
			width:  1280,
			height: 720,
		},
		{
			name:   "480p full frame",
			rect:   Rect{Width: 640, Height: 480},
			width:  640,
			height: 480,
		},
		{
			name:   "4K full frame",
			rect:   Rect{Width: 3840, Height: 2160},
			width:  3840,
			height: 2160,
		},
		{
			name:   "Center crop 640x480 from 1920x1080",
			rect:   Rect{Left: 640, Top: 300, Width: 640, Height: 480},
			width:  640,
			height: 480,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.rect.Width != tt.width {
				t.Errorf("Width: expected %d, got %d", tt.width, tt.rect.Width)
			}
			if tt.rect.Height != tt.height {
				t.Errorf("Height: expected %d, got %d", tt.height, tt.rect.Height)
			}
		})
	}
}

// TestFract_StructFields tests Fract field accessibility
func TestFract_StructFields(t *testing.T) {
	f := Fract{
		Numerator:   16,
		Denominator: 9,
	}

	if f.Numerator != 16 {
		t.Errorf("Numerator: expected 16, got %d", f.Numerator)
	}
	if f.Denominator != 9 {
		t.Errorf("Denominator: expected 9, got %d", f.Denominator)
	}
}

// TestFract_CommonAspectRatios tests common pixel aspect ratios
func TestFract_CommonAspectRatios(t *testing.T) {
	tests := []struct {
		name  string
		fract Fract
		num   uint32
		denom uint32
	}{
		{
			name:  "1:1 (square pixels)",
			fract: Fract{Numerator: 1, Denominator: 1},
			num:   1,
			denom: 1,
		},
		{
			name:  "16:9 (HD)",
			fract: Fract{Numerator: 16, Denominator: 9},
			num:   16,
			denom: 9,
		},
		{
			name:  "4:3 (SD)",
			fract: Fract{Numerator: 4, Denominator: 3},
			num:   4,
			denom: 3,
		},
		{
			name:  "21:9 (ultrawide)",
			fract: Fract{Numerator: 21, Denominator: 9},
			num:   21,
			denom: 9,
		},
		{
			name:  "32:27 (NTSC 4:3 with non-square pixels)",
			fract: Fract{Numerator: 32, Denominator: 27},
			num:   32,
			denom: 27,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fract.Numerator != tt.num {
				t.Errorf("Numerator: expected %d, got %d", tt.num, tt.fract.Numerator)
			}
			if tt.fract.Denominator != tt.denom {
				t.Errorf("Denominator: expected %d, got %d", tt.denom, tt.fract.Denominator)
			}
		})
	}
}

// TestFract_ZeroValues tests Fract with zero values
func TestFract_ZeroValues(t *testing.T) {
	f := Fract{}

	if f.Numerator != 0 || f.Denominator != 0 {
		t.Error("Zero-initialized Fract should have all fields at 0")
	}
}

// TestCropCapability_BoundsVsDefaultRect tests the relationship between bounds and default
func TestCropCapability_BoundsVsDefaultRect(t *testing.T) {
	tests := []struct {
		name        string
		cap         CropCapability
		boundsLarger bool
		description string
	}{
		{
			name: "Bounds equal to default",
			cap: CropCapability{
				Bounds:      Rect{Width: 1920, Height: 1080},
				DefaultRect: Rect{Width: 1920, Height: 1080},
			},
			boundsLarger: false,
			description:  "Full frame capture, no crop",
		},
		{
			name: "Default smaller than bounds",
			cap: CropCapability{
				Bounds:      Rect{Width: 1920, Height: 1080},
				DefaultRect: Rect{Left: 640, Top: 300, Width: 640, Height: 480},
			},
			boundsLarger: true,
			description:  "Default crop region within larger sensor area",
		},
		{
			name: "Offset default rectangle",
			cap: CropCapability{
				Bounds:      Rect{Width: 1920, Height: 1080},
				DefaultRect: Rect{Left: 100, Top: 100, Width: 1720, Height: 880},
			},
			boundsLarger: false,
			description:  "Default crop with offset but similar size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boundsArea := tt.cap.Bounds.Width * tt.cap.Bounds.Height
			defaultArea := tt.cap.DefaultRect.Width * tt.cap.DefaultRect.Height

			if tt.boundsLarger && boundsArea <= defaultArea {
				t.Errorf("Expected bounds area (%d) > default area (%d)", boundsArea, defaultArea)
			}
		})
	}
}

// TestCropCapability_TypicalUseCases tests typical cropping scenarios
func TestCropCapability_TypicalUseCases(t *testing.T) {
	tests := []struct {
		name        string
		cap         CropCapability
		description string
	}{
		{
			name: "HD webcam",
			cap: CropCapability{
				StreamType: BufTypeVideoCapture,
				Bounds: Rect{
					Width:  1920,
					Height: 1080,
				},
				DefaultRect: Rect{
					Width:  1920,
					Height: 1080,
				},
				PixelAspect: Fract{
					Numerator:   1,
					Denominator: 1,
				},
			},
			description: "Standard 1080p webcam with square pixels",
		},
		{
			name: "Center crop for portrait mode",
			cap: CropCapability{
				StreamType: BufTypeVideoCapture,
				Bounds: Rect{
					Width:  1920,
					Height: 1080,
				},
				DefaultRect: Rect{
					Left:   720,
					Top:    0,
					Width:  720,
					Height: 1080,
				},
				PixelAspect: Fract{
					Numerator:   1,
					Denominator: 1,
				},
			},
			description: "Portrait mode: vertical center crop",
		},
		{
			name: "Digital zoom (center crop 2x)",
			cap: CropCapability{
				StreamType: BufTypeVideoCapture,
				Bounds: Rect{
					Width:  1920,
					Height: 1080,
				},
				DefaultRect: Rect{
					Left:   480,
					Top:    270,
					Width:  960,
					Height: 540,
				},
				PixelAspect: Fract{
					Numerator:   1,
					Denominator: 1,
				},
			},
			description: "2x digital zoom via center crop",
		},
		{
			name: "SD capture with non-square pixels",
			cap: CropCapability{
				StreamType: BufTypeVideoCapture,
				Bounds: Rect{
					Width:  720,
					Height: 480,
				},
				DefaultRect: Rect{
					Width:  720,
					Height: 480,
				},
				PixelAspect: Fract{
					Numerator:   32,
					Denominator: 27,
				},
			},
			description: "NTSC DV with non-square pixels",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify StreamType is set
			if tt.cap.StreamType == 0 {
				t.Error("StreamType should be set")
			}

			// Verify bounds has non-zero dimensions
			if tt.cap.Bounds.Width == 0 || tt.cap.Bounds.Height == 0 {
				t.Error("Bounds should have non-zero dimensions")
			}

			// Verify default rect has non-zero dimensions
			if tt.cap.DefaultRect.Width == 0 || tt.cap.DefaultRect.Height == 0 {
				t.Error("DefaultRect should have non-zero dimensions")
			}

			// Verify pixel aspect is set
			if tt.cap.PixelAspect.Numerator == 0 || tt.cap.PixelAspect.Denominator == 0 {
				t.Error("PixelAspect should have non-zero values")
			}
		})
	}
}

// TestCropCapability_StringFormatting tests the String() output format
func TestCropCapability_StringFormatting(t *testing.T) {
	cap := CropCapability{
		DefaultRect: Rect{
			Top:    10,
			Left:   20,
			Width:  640,
			Height: 480,
		},
		Bounds: Rect{
			Top:    0,
			Left:   0,
			Width:  1920,
			Height: 1080,
		},
		PixelAspect: Fract{
			Numerator:   16,
			Denominator: 9,
		},
	}

	result := cap.String()

	// Verify it contains key sections
	requiredParts := []string{
		"default:",
		"bounds:",
		"pixel-aspect",
		"top=10",
		"left=20",
		"width=640",
		"height=480",
		"top=0",
		"left=0",
		"width=1920",
		"height=1080",
		"16:9",
	}

	for _, part := range requiredParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected String() output to contain %q, got: %s", part, result)
		}
	}
}

// TestArea_StructFields tests Area field accessibility
func TestArea_StructFields(t *testing.T) {
	a := Area{
		Width:  1920,
		Height: 1080,
	}

	if a.Width != 1920 {
		t.Errorf("Width: expected 1920, got %d", a.Width)
	}
	if a.Height != 1080 {
		t.Errorf("Height: expected 1080, got %d", a.Height)
	}
}

// TestArea_CommonSizes tests typical video areas
func TestArea_CommonSizes(t *testing.T) {
	tests := []struct {
		name   string
		area   Area
		width  uint32
		height uint32
	}{
		{"1080p", Area{Width: 1920, Height: 1080}, 1920, 1080},
		{"720p", Area{Width: 1280, Height: 720}, 1280, 720},
		{"480p", Area{Width: 640, Height: 480}, 640, 480},
		{"4K", Area{Width: 3840, Height: 2160}, 3840, 2160},
		{"8K", Area{Width: 7680, Height: 4320}, 7680, 4320},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.area.Width != tt.width {
				t.Errorf("Width: expected %d, got %d", tt.width, tt.area.Width)
			}
			if tt.area.Height != tt.height {
				t.Errorf("Height: expected %d, got %d", tt.height, tt.area.Height)
			}
		})
	}
}

// TestArea_ZeroValues tests Area with zero values
func TestArea_ZeroValues(t *testing.T) {
	a := Area{}

	if a.Width != 0 || a.Height != 0 {
		t.Error("Zero-initialized Area should have all fields at 0")
	}
}

// TestCropCapability_StreamTypeValues tests different stream types
func TestCropCapability_StreamTypeValues(t *testing.T) {
	tests := []struct {
		name       string
		streamType uint32
	}{
		{"VideoCapture", BufTypeVideoCapture},
		{"VideoOutput", BufTypeVideoOutput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cap := CropCapability{
				StreamType: tt.streamType,
			}

			if cap.StreamType != tt.streamType {
				t.Errorf("StreamType: expected %d, got %d", tt.streamType, cap.StreamType)
			}
		})
	}
}
