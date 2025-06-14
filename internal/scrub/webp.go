package scrub

import (
	"fmt"
	"os"

	"github.com/chai2010/webp"
)

func ScrubWEBP(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open error: %w", err)
	}
	defer f.Close()

	img, err := webp.Decode(f)
	if err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	tmpPath := path + ".clean"
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	defer out.Close()

	if err := webp.Encode(out, img, &webp.Options{Lossless: true}); err != nil {
		return fmt.Errorf("encode error: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename error: %w", err)
	}

	return nil
}
