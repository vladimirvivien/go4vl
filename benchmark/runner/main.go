package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"time"

	"github.com/vladimirvivien/go4vl/benchmark/loopback"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

type BenchmarkScenario struct {
	Name   string
	Width  int
	Height int
	FPS    int
	Format string
}

type BenchmarkConfig struct {
	DevicePath  string
	Width       uint32
	Height      uint32
	PixelFormat uint32
	FPS         uint32
	Duration    time.Duration
	BufferSize  uint32
	CPUProfile  string
	MemProfile  string
	TraceFile   string
	Verbose     bool
}

type BenchmarkResults struct {
	FramesCaptured   int
	FramesDropped    int
	Duration         time.Duration
	AvgFPS           float64
	MinFrameTime     time.Duration
	MaxFrameTime     time.Duration
	AvgFrameTime     time.Duration
	TotalBytes       uint64
	AvgBytesPerFrame uint64
	MemAllocBytes    uint64
	MemAllocObjects  uint64
	NumGC            uint32
	GCPauseTotal     time.Duration
}

var scenarios = []BenchmarkScenario{
	// Baseline benchmarks
	{Name: "baseline_480p_mjpeg", Width: 640, Height: 480, FPS: 30, Format: "MJPEG"},
	{Name: "baseline_720p_mjpeg", Width: 1280, Height: 720, FPS: 30, Format: "MJPEG"},
	{Name: "baseline_1080p_mjpeg", Width: 1920, Height: 1080, FPS: 30, Format: "MJPEG"},

	// Format comparison (at 480p)
	{Name: "format_480p_yuyv", Width: 640, Height: 480, FPS: 30, Format: "YUYV"},

	// Frame rate tests (at 720p)
	{Name: "fps_720p_15fps", Width: 1280, Height: 720, FPS: 15, Format: "MJPEG"},
	{Name: "fps_720p_30fps", Width: 1280, Height: 720, FPS: 30, Format: "MJPEG"},
	{Name: "fps_720p_60fps", Width: 1280, Height: 720, FPS: 60, Format: "MJPEG"},
}

func main() {
	// Mode selection (mutually exclusive)
	device := flag.String("device", "", "Use real device at path (e.g., /dev/video0)")

	// Loopback device parameters (default mode if -device not specified)
	loopbackNum := flag.Int("loopback-num", 50, "Loopback device number (default mode)")
	testPattern := flag.String("test-pattern", "testsrc", "FFmpeg test pattern: testsrc, smptebars, color=<color>")

	// Benchmark parameters
	duration := flag.String("duration", "10s", "Capture duration per benchmark")
	outputDir := flag.String("output", "", "Output directory (default: results_<timestamp>)")
	scenario := flag.String("scenario", "", "Run specific scenario by name (empty = run all)")
	list := flag.Bool("list", false, "List available scenarios and exit")
	single := flag.Bool("single", false, "Run single benchmark with custom parameters")

	// Single benchmark parameters
	width := flag.Uint("width", 640, "Frame width")
	height := flag.Uint("height", 480, "Frame height")
	fps := flag.Uint("fps", 30, "Frames per second")
	format := flag.String("format", "MJPEG", "Pixel format: MJPEG, YUYV, H264")
	buffers := flag.Uint("buffers", 4, "Number of buffers")
	verbose := flag.Bool("verbose", false, "Verbose output")
	cpuprofile := flag.String("cpuprofile", "", "Write CPU profile to file")
	memprofile := flag.String("memprofile", "", "Write memory profile to file")
	tracefile := flag.String("trace", "", "Write execution trace to file")

	flag.Parse()

	// List scenarios if requested
	if *list {
		fmt.Println("Available benchmark scenarios:")
		for _, s := range scenarios {
			fmt.Printf("  %s: %dx%d @ %d fps (%s)\n", s.Name, s.Width, s.Height, s.FPS, s.Format)
		}
		return
	}

	// Determine device path: use real device if specified, otherwise setup loopback
	var devicePath string
	var loopbackDev *loopback.Device

	if *device != "" {
		// User specified a real device
		devicePath = *device
		log.Printf("Using real device: %s", devicePath)
	} else {
		// Default: setup loopback device
		if !loopback.IsAvailable() {
			log.Fatal("Loopback mode requires ffmpeg and v4l2loopback-dkms to be installed.\nInstall with: sudo apt install ffmpeg v4l2loopback-dkms")
		}

		log.Printf("Setting up loopback device /dev/video%d...", *loopbackNum)
		var err error
		loopbackDev, err = loopback.Setup(*loopbackNum, int(*width), int(*height), int(*fps), *testPattern)
		if err != nil {
			log.Fatalf("Failed to setup loopback device: %v", err)
		}
		defer loopbackDev.Close()

		devicePath = loopbackDev.DevicePath
		log.Printf("Loopback device ready at %s", devicePath)
	}

	// Single benchmark mode
	if *single {
		config := BenchmarkConfig{
			DevicePath: devicePath,
			Width:      uint32(*width),
			Height:     uint32(*height),
			FPS:        uint32(*fps),
			Duration:   parseDuration(*duration),
			BufferSize: uint32(*buffers),
			CPUProfile: *cpuprofile,
			MemProfile: *memprofile,
			TraceFile:  *tracefile,
			Verbose:    *verbose,
		}

		// Map format string
		switch *format {
		case "MJPEG":
			config.PixelFormat = v4l2.PixelFmtMJPEG
		case "YUYV":
			config.PixelFormat = v4l2.PixelFmtYUYV
		case "H264":
			config.PixelFormat = v4l2.PixelFmtH264
		default:
			log.Fatalf("Unknown pixel format: %s", *format)
		}

		runSingleBenchmark(config)
		return
	}

	// Multi-scenario benchmark mode
	runMultiScenario(devicePath, *duration, *outputDir, *scenario)
}

func runSingleBenchmark(config BenchmarkConfig) {
	// Start profiling if requested
	if config.CPUProfile != "" {
		f, err := os.Create(config.CPUProfile)
		if err != nil {
			log.Fatalf("Could not create CPU profile: %v", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("Could not start CPU profile: %v", err)
		}
		defer pprof.StopCPUProfile()
		log.Printf("CPU profiling enabled: %s", config.CPUProfile)
	}

	if config.TraceFile != "" {
		f, err := os.Create(config.TraceFile)
		if err != nil {
			log.Fatalf("Could not create trace file: %v", err)
		}
		defer f.Close()
		if err := trace.Start(f); err != nil {
			log.Fatalf("Could not start trace: %v", err)
		}
		defer trace.Stop()
		log.Printf("Execution trace enabled: %s", config.TraceFile)
	}

	// Run benchmark
	results := runDeviceBenchmark(config)

	// Write memory profile if requested
	if config.MemProfile != "" {
		f, err := os.Create(config.MemProfile)
		if err != nil {
			log.Fatalf("Could not create memory profile: %v", err)
		}
		defer f.Close()
		runtime.GC() // Force GC before taking heap snapshot
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("Could not write memory profile: %v", err)
		}
		log.Printf("Memory profile written: %s", config.MemProfile)
	}

	// Print results
	printResults(config, results)
}

func runMultiScenario(devicePath, duration, outputDir, scenarioName string) {
	// Create output directory
	if outputDir == "" {
		outputDir = fmt.Sprintf("results_%s", time.Now().Format("20060102_150405"))
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("go4vl Frame Capture Benchmark Suite")
	fmt.Println("========================================")
	fmt.Printf("Device:   %s\n", devicePath)
	fmt.Printf("Duration: %s\n", duration)
	fmt.Printf("Output:   %s\n", outputDir)
	fmt.Println()

	// Filter scenarios if specific one requested
	var toRun []BenchmarkScenario
	if scenarioName != "" {
		found := false
		for _, s := range scenarios {
			if s.Name == scenarioName {
				toRun = []BenchmarkScenario{s}
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("Unknown scenario: %s (use -list to see available scenarios)", scenarioName)
		}
	} else {
		toRun = scenarios
	}

	// Run benchmarks
	results := make(map[string]string)
	for i, s := range toRun {
		fmt.Printf("[%d/%d] Running: %s (%dx%d @ %d fps, %s)\n",
			i+1, len(toRun), s.Name, s.Width, s.Height, s.FPS, s.Format)

		resultFile := filepath.Join(outputDir, fmt.Sprintf("%s_results.txt", s.Name))
		cpuProfile := filepath.Join(outputDir, fmt.Sprintf("%s_cpu.prof", s.Name))
		memProfile := filepath.Join(outputDir, fmt.Sprintf("%s_mem.prof", s.Name))
		traceFile := filepath.Join(outputDir, fmt.Sprintf("%s_trace.out", s.Name))

		// Parse format
		var pixelFormat uint32
		switch s.Format {
		case "MJPEG":
			pixelFormat = v4l2.PixelFmtMJPEG
		case "YUYV":
			pixelFormat = v4l2.PixelFmtYUYV
		case "H264":
			pixelFormat = v4l2.PixelFmtH264
		default:
			log.Printf("Warning: Unknown format %s, skipping", s.Format)
			results[s.Name] = "SKIPPED"
			continue
		}

		config := BenchmarkConfig{
			DevicePath:  devicePath,
			Width:       uint32(s.Width),
			Height:      uint32(s.Height),
			FPS:         uint32(s.FPS),
			PixelFormat: pixelFormat,
			Duration:    parseDuration(duration),
			BufferSize:  4,
			CPUProfile:  cpuProfile,
			MemProfile:  memProfile,
			TraceFile:   traceFile,
			Verbose:     false,
		}

		// Capture output to buffer
		var buf bytes.Buffer
		oldStdout := os.Stdout
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stdout = w
		os.Stderr = w

		// Run in goroutine to capture output
		done := make(chan bool)
		go func() {
			buf.ReadFrom(r)
			done <- true
		}()

		// Start profiling
		var cpuFile, traceF *os.File
		if config.CPUProfile != "" {
			cpuFile, _ = os.Create(config.CPUProfile)
			if cpuFile != nil {
				pprof.StartCPUProfile(cpuFile)
			}
		}
		if config.TraceFile != "" {
			traceF, _ = os.Create(config.TraceFile)
			if traceF != nil {
				trace.Start(traceF)
			}
		}

		// Run benchmark
		benchResults := runDeviceBenchmark(config)

		// Stop profiling
		if cpuFile != nil {
			pprof.StopCPUProfile()
			cpuFile.Close()
		}
		if traceF != nil {
			trace.Stop()
			traceF.Close()
		}

		// Write memory profile
		if config.MemProfile != "" {
			memFile, _ := os.Create(config.MemProfile)
			if memFile != nil {
				runtime.GC()
				pprof.WriteHeapProfile(memFile)
				memFile.Close()
			}
		}

		// Print results to buffer
		printResults(config, benchResults)

		// Restore stdout/stderr
		w.Close()
		<-done
		os.Stdout = oldStdout
		os.Stderr = oldStderr

		// Save output
		output := buf.String()
		if err := os.WriteFile(resultFile, []byte(output), 0644); err != nil {
			log.Printf("Warning: Failed to write result file: %v", err)
			results[s.Name] = "FAILED"
		} else {
			fmt.Println("  ✓ Complete")
			results[s.Name] = "SUCCESS"
		}
		fmt.Println()
	}

	// Print summary
	fmt.Println("========================================")
	fmt.Println("Benchmarks Complete!")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Results:")
	for _, s := range toRun {
		status := results[s.Name]
		fmt.Printf("  [%s] %s\n", status, s.Name)
	}
	fmt.Println()
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Println()
	fmt.Println("View results:")
	fmt.Printf("  cat %s/*_results.txt\n", outputDir)
	fmt.Println()
	fmt.Println("Analyze CPU profiles:")
	fmt.Printf("  go tool pprof %s/*_cpu.prof\n", outputDir)
	fmt.Println()
	fmt.Println("Analyze memory profiles:")
	fmt.Printf("  go tool pprof %s/*_mem.prof\n", outputDir)
	fmt.Println()
	fmt.Println("View execution traces:")
	fmt.Printf("  go tool trace %s/*_trace.out\n", outputDir)
	fmt.Println()

	// Generate summary file
	generateSummary(outputDir, toRun)
}

func runDeviceBenchmark(config BenchmarkConfig) BenchmarkResults {
	log.Printf("Opening device: %s", config.DevicePath)
	log.Printf("Format: %dx%d @ %d fps", config.Width, config.Height, config.FPS)
	log.Printf("Duration: %v, Buffers: %d", config.Duration, config.BufferSize)

	// Open device
	dev, err := device.Open(
		config.DevicePath,
		device.WithBufferSize(config.BufferSize),
		device.WithPixFormat(v4l2.PixFormat{
			Width:       config.Width,
			Height:      config.Height,
			PixelFormat: config.PixelFormat,
			Field:       v4l2.FieldNone,
		}),
		device.WithFPS(config.FPS),
	)
	if err != nil {
		log.Fatalf("Failed to open device: %v", err)
	}
	defer dev.Close()

	// Verify actual format
	actualFormat, err := dev.GetPixFormat()
	if err != nil {
		log.Fatalf("Failed to get format: %v", err)
	}
	log.Printf("Actual format: %dx%d, bytesperline=%d, sizeimage=%d",
		actualFormat.Width, actualFormat.Height,
		actualFormat.BytesPerLine, actualFormat.SizeImage)

	// Collect memory stats before
	var memStatsBefore runtime.MemStats
	runtime.ReadMemStats(&memStatsBefore)

	// Start capture
	ctx, cancel := context.WithTimeout(context.Background(), config.Duration)
	defer cancel()

	if err := dev.Start(ctx); err != nil {
		log.Fatalf("Failed to start capture: %v", err)
	}
	defer dev.Stop()

	log.Println("Capturing frames...")

	// Capture loop with timing
	results := BenchmarkResults{}
	frameTimes := make([]time.Duration, 0, 1000)
	lastFrameTime := time.Now()
	startTime := time.Now()

	for frame := range dev.GetOutput() {
		now := time.Now()
		frameTime := now.Sub(lastFrameTime)

		if len(frame) == 0 {
			results.FramesDropped++
			if config.Verbose {
				log.Printf("Frame %d: DROPPED", results.FramesCaptured+results.FramesDropped)
			}
		} else {
			results.FramesCaptured++
			results.TotalBytes += uint64(len(frame))
			frameTimes = append(frameTimes, frameTime)

			if config.Verbose && results.FramesCaptured%100 == 0 {
				log.Printf("Captured %d frames (%.1f fps)", results.FramesCaptured,
					float64(results.FramesCaptured)/time.Since(startTime).Seconds())
			}
		}

		lastFrameTime = now
	}

	results.Duration = time.Since(startTime)

	// Collect memory stats after
	var memStatsAfter runtime.MemStats
	runtime.ReadMemStats(&memStatsAfter)

	// Calculate statistics
	if results.FramesCaptured > 0 {
		results.AvgFPS = float64(results.FramesCaptured) / results.Duration.Seconds()
		results.AvgBytesPerFrame = results.TotalBytes / uint64(results.FramesCaptured)

		// Calculate frame time statistics
		if len(frameTimes) > 1 {
			results.MinFrameTime = frameTimes[1] // Skip first
			results.MaxFrameTime = frameTimes[1]
			var totalFrameTime time.Duration

			for i := 1; i < len(frameTimes); i++ { // Skip first frame
				ft := frameTimes[i]
				totalFrameTime += ft
				if ft < results.MinFrameTime {
					results.MinFrameTime = ft
				}
				if ft > results.MaxFrameTime {
					results.MaxFrameTime = ft
				}
			}
			results.AvgFrameTime = totalFrameTime / time.Duration(len(frameTimes)-1)
		}
	}

	// Memory statistics
	results.MemAllocBytes = memStatsAfter.TotalAlloc - memStatsBefore.TotalAlloc
	results.MemAllocObjects = memStatsAfter.Mallocs - memStatsBefore.Mallocs
	results.NumGC = memStatsAfter.NumGC - memStatsBefore.NumGC
	results.GCPauseTotal = time.Duration(memStatsAfter.PauseTotalNs - memStatsBefore.PauseTotalNs)

	return results
}

func printResults(config BenchmarkConfig, r BenchmarkResults) {
	separator := strings.Repeat("=", 70)
	fmt.Println("\n" + separator)
	fmt.Println("BENCHMARK RESULTS")
	fmt.Println(separator)

	fmt.Println("\nConfiguration:")
	fmt.Printf("  Device:        %s\n", config.DevicePath)
	fmt.Printf("  Resolution:    %dx%d\n", config.Width, config.Height)
	fmt.Printf("  Target FPS:    %d\n", config.FPS)
	fmt.Printf("  Duration:      %v\n", config.Duration)
	fmt.Printf("  Buffers:       %d\n", config.BufferSize)

	fmt.Println("\nCapture Statistics:")
	fmt.Printf("  Frames Captured:   %d\n", r.FramesCaptured)
	fmt.Printf("  Frames Dropped:    %d\n", r.FramesDropped)
	fmt.Printf("  Actual Duration:   %v\n", r.Duration)
	fmt.Printf("  Average FPS:       %.2f\n", r.AvgFPS)
	fmt.Printf("  Total Data:        %.2f MB\n", float64(r.TotalBytes)/(1024*1024))
	fmt.Printf("  Avg Bytes/Frame:   %d\n", r.AvgBytesPerFrame)

	fmt.Println("\nTiming Statistics:")
	fmt.Printf("  Min Frame Time:    %v\n", r.MinFrameTime)
	fmt.Printf("  Avg Frame Time:    %v\n", r.AvgFrameTime)
	fmt.Printf("  Max Frame Time:    %v\n", r.MaxFrameTime)
	if config.FPS > 0 {
		targetFrameTime := time.Second / time.Duration(config.FPS)
		fmt.Printf("  Target Frame Time: %v\n", targetFrameTime)
	}

	fmt.Println("\nMemory Statistics:")
	fmt.Printf("  Total Allocated:   %.2f MB\n", float64(r.MemAllocBytes)/(1024*1024))
	if r.FramesCaptured > 0 {
		fmt.Printf("  Allocs per Frame:  %.2f MB\n", float64(r.MemAllocBytes)/float64(r.FramesCaptured)/(1024*1024))
	}
	fmt.Printf("  Total Allocations: %d\n", r.MemAllocObjects)
	if r.FramesCaptured > 0 {
		fmt.Printf("  Allocs per Frame:  %.0f\n", float64(r.MemAllocObjects)/float64(r.FramesCaptured))
	}
	fmt.Printf("  GC Runs:           %d\n", r.NumGC)
	fmt.Printf("  GC Pause Total:    %v\n", r.GCPauseTotal)
	if r.NumGC > 0 {
		fmt.Printf("  Avg GC Pause:      %v\n", r.GCPauseTotal/time.Duration(r.NumGC))
	}

	fmt.Println("\nPerformance Metrics:")
	if r.AvgFPS > 0 && r.FramesCaptured > 0 {
		cpuTimePerFrame := r.Duration / time.Duration(r.FramesCaptured)
		fmt.Printf("  CPU Time/Frame:    %v\n", cpuTimePerFrame)
		fmt.Printf("  Throughput:        %.2f MB/s\n",
			float64(r.TotalBytes)/(1024*1024)/r.Duration.Seconds())
	}

	fmt.Println("\n" + separator)

	if config.CPUProfile != "" {
		fmt.Printf("\nCPU Profile: %s\n", config.CPUProfile)
		fmt.Println("Analyze with: go tool pprof " + config.CPUProfile)
	}
	if config.MemProfile != "" {
		fmt.Printf("Memory Profile: %s\n", config.MemProfile)
		fmt.Println("Analyze with: go tool pprof " + config.MemProfile)
	}
	if config.TraceFile != "" {
		fmt.Printf("Trace File: %s\n", config.TraceFile)
		fmt.Println("Analyze with: go tool trace " + config.TraceFile)
	}
}

func generateSummary(outputDir string, scenarios []BenchmarkScenario) {
	summaryFile := filepath.Join(outputDir, "summary.txt")

	summary := "# Benchmark Summary\n\n"
	summary += fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339))

	for _, s := range scenarios {
		resultFile := filepath.Join(outputDir, fmt.Sprintf("%s_results.txt", s.Name))
		if _, err := os.Stat(resultFile); os.IsNotExist(err) {
			continue
		}

		content, err := os.ReadFile(resultFile)
		if err != nil {
			continue
		}

		summary += fmt.Sprintf("## %s\n\n", s.Name)
		summary += string(content)
		summary += "\n\n"
	}

	if err := os.WriteFile(summaryFile, []byte(summary), 0644); err != nil {
		log.Printf("Warning: Failed to write summary: %v", err)
	} else {
		fmt.Printf("Summary written to: %s\n", summaryFile)
	}
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("Invalid duration: %s", s)
	}
	return d
}
