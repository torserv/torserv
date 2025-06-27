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

	// Create a temporary file in the same directory to hold the cleaned image
	tmpPath := path + ".clean"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file %s: %w", tmpPath, err)
	}
	defer func() {
		if cerr := tmpFile.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close temp file %s: %v\n", tmpPath, cerr)
		}
	}()

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
