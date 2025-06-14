package scrub

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var supported = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
var unsafe = []string{".doc", ".docx", ".pptx", ".xls", ".xlsx", ".odt", ".mp4", ".mp3", ".mkv", ".mov", ".pdf"}

func Init() error {
	return filepath.Walk("public", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
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
			return nil
		}
	})
}

func contains(list []string, target string) bool {
	for _, ext := range list {
		if ext == target {
			return true
		}
	}
	return false
}

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
	case ".webp":
		return ScrubWEBP(path)
	default:
		return nil
	}
}
