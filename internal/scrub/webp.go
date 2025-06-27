package scrub

import (
	"fmt"
	"os"

	"github.com/chai2010/webp"
)

// ScrubWEBP re-encodes a WEBP image in lossless mode to strip metadata.
// The cleaned image replaces the original file.
func ScrubWEBP(path string) error {
	// Open the original WEBP file
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open error: %w", err)
	}
	defer f.Close()

	// Decode the image to remove metadata
	img, err := webp.Decode(f)
	if err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	// Create a temporary file to hold the scrubbed image
	tmpPath := fmt.Sprintf("%s.clean", path)
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file %s: %w", tmpPath, err)
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close temp file %s: %v\n", tmpPath, cerr)
		}
	}()

	// Encode image in lossless mode
	if err := webp.Encode(out, img, &webp.Options{Lossless: true}); err != nil {
		return fmt.Errorf("encode error: %w", err)
	}

	// Replace the original file with the scrubbed version
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename error: %w", err)
	}

	return nil
}
