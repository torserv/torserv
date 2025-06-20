package scrub

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"os"
)

// ScrubJPEG decodes and re-encodes a JPEG file to strip all metadata.
// The scrubbed image overwrites the original file.
func ScrubJPEG(path string) error {
	// Open the original JPEG file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Decode the image — this implicitly drops metadata
	img, err := jpeg.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode JPEG: %w", err)
	}

	// Re-encode the image into memory without metadata
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	// Write scrubbed image to a temporary file
	tmpPath := path + ".scrubbed"
	tmp, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, &buf); err != nil {
		return fmt.Errorf("failed to write scrubbed JPEG: %w", err)
	}

	// Replace the original file with the scrubbed version
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to replace original: %w", err)
	}

	return nil
}
