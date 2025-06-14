package scrub

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"os"
)

func ScrubJPEG(path string) error {
	// Open original file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Decode JPEG image without preserving metadata
	img, err := jpeg.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode JPEG: %w", err)
	}

	// Re-encode to memory without metadata
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	// Overwrite original file
	tmpPath := path + ".scrubbed"
	tmp, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, &buf); err != nil {
		return fmt.Errorf("failed to write scrubbed JPEG: %w", err)
	}

	// Replace original
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to replace original: %w", err)
	}

	return nil
}
