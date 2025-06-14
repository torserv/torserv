package scrub

import (
	"fmt"
	"os"

	"golang.org/x/image/bmp"
)

// ScrubBMP re-encodes the BMP image to strip metadata.
func ScrubBMP(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open error: %w", err)
	}
	defer file.Close()

	img, err := bmp.Decode(file)
	if err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	tmpPath := path + ".clean"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("temp file error: %w", err)
	}
	defer tmpFile.Close()

	if err := bmp.Encode(tmpFile, img); err != nil {
		return fmt.Errorf("encode error: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replace error: %w", err)
	}

	return nil
}
