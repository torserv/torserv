package scrub

import (
	"fmt"
	"os"

	"golang.org/x/image/bmp"
)

// ScrubBMP re-encodes a BMP image to remove any embedded metadata.
// The cleaned image replaces the original file on disk.
func ScrubBMP(path string) error {
	// Open the original BMP file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open error: %w", err)
	}
	defer file.Close()

	// Decode the BMP image into an in-memory image.Image
	img, err := bmp.Decode(file)
	if err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	// Create a temporary file to hold the cleaned image
	tmpPath := path + ".clean"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("temp file error: %w", err)
	}
	defer tmpFile.Close()

	// Re-encode the image without metadata
	if err := bmp.Encode(tmpFile, img); err != nil {
		return fmt.Errorf("encode error: %w", err)
	}

	// Replace the original file with the cleaned version
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replace error: %w", err)
	}

	return nil
}
