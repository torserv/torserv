package scrub

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// List of supported image extensions for scrubbing
var supported = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

// List of known unsafe file types to reject
var unsafe = []string{
	".doc", ".docx", ".pptx", ".xls", ".xlsx", ".odt",
	".mp4", ".mp3", ".mkv", ".mov", ".pdf", ".webp",
}

// Init walks the "public" directory, scrubs supported image files,
// and errors out if any unsafe files are detected.
func Init() error {
	return filepath.Walk("public", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil // Skip unreadable paths or directories
		}

		ext := strings.ToLower(filepath.Ext(path))

		switch {
		case contains(unsafe, ext):
			fmt.Printf("[!] UNSAFE FILE: %s — unsupported type\n", path)
			return fmt.Errorf("unsafe file found: %s", path)

		case contains(supported, ext):
			fmt.Printf("[*] Scrubbing: %s\n", path)
			return ScrubFile(path, ext)

		default:
			// Skip files with unknown extensions
			return nil
		}
	})
}

// contains checks if the given target extension exists in the provided list, returns if true.
func contains(list []string, target string) bool {
	for _, ext := range list {
		if ext == target {
			return true
		}
	}
	return false
}

// ScrubFile dispatches the appropriate scrubber function based on file extension.
func ScrubFile(path, ext string) error {
	switch ext {
	case ".jpg", ".jpeg":
		return ScrubJPEG(path)
	case ".png":
		return ScrubPNG(path)
	case ".gif":
		return ScrubGIF(path)
	case ".bmp":
		return ScrubBMP(path)
	default:
		// Unsupported extensions are ignored
		return nil
	}
}
