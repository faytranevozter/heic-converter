package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"

	"github.com/adrium/goheif"
)

type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

var verbose bool

func main() {
	inputPath := flag.String("input", "", "Path to input HEIC file")
	outputPath := flag.String("output", "", "Path to output JPEG file")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Usage: goheicconvert -input=image1.heic -output=sample.jpg [-verbose]")
		os.Exit(1)
	}

	logf("🔍 Input file:  %s", *inputPath)
	logf("💾 Output file: %s", *outputPath)

	err := convertHeicToJpg(*inputPath, *outputPath)
	if err != nil {
		log.Fatalf("❌ Conversion failed: %v", err)
	}

	info, err := os.Stat(*outputPath)
	if err == nil {
		fmt.Println("✅ Conversion successful!")
		fmt.Printf("📸 JPEG saved (%d bytes)\n", info.Size())
	} else {
		fmt.Println("⚠️ Conversion done, but failed to retrieve output file info.")
	}
}

func convertHeicToJpg(input, output string) error {
	fileInput, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer fileInput.Close()

	logf("📤 Extracting EXIF metadata...")
	exif, err := goheif.ExtractExif(fileInput)
	if err != nil {
		return fmt.Errorf("failed to extract EXIF: %w", err)
	}

	if _, err := fileInput.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset input file reader: %w", err)
	}

	logf("🖼️  Decoding HEIC image...")
	img, err := goheif.Decode(fileInput)
	if err != nil {
		return fmt.Errorf("failed to decode HEIC: %w", err)
	}

	fileOutput, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer fileOutput.Close()

	w, err := newWriterExif(fileOutput, exif)
	if err != nil {
		return fmt.Errorf("failed to write EXIF data: %w", err)
	}

	logf("📝 Encoding JPEG...")
	err = jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return nil
}

func (w *writerSkipper) Write(data []byte) (int, error) {
	if w.bytesToSkip <= 0 {
		return w.w.Write(data)
	}
	if len(data) < w.bytesToSkip {
		w.bytesToSkip -= len(data)
		return len(data), nil
	}

	n, err := w.w.Write(data[w.bytesToSkip:])
	if err != nil {
		return n, err
	}
	n += w.bytesToSkip
	w.bytesToSkip = 0
	return n, nil
}

func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	skipper := &writerSkipper{w: w, bytesToSkip: 2}

	soi := []byte{0xff, 0xd8}
	if _, err := w.Write(soi); err != nil {
		return nil, err
	}

	if exif != nil {
		app1Marker := 0xe1
		markerLen := 2 + len(exif)
		marker := []byte{0xff, uint8(app1Marker), uint8(markerLen >> 8), uint8(markerLen & 0xff)}

		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return skipper, nil
}

// logf prints formatted log output only if verbose is enabled
func logf(format string, a ...any) {
	if verbose {
		fmt.Printf(format+"\n", a...)
	}
}
