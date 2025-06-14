package scrub

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// ScrubGIF removes comment and non-animation app extensions from GIFs.
func ScrubGIF(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	var out bytes.Buffer
	r := bytes.NewReader(data)

	// Copy header and logical screen descriptor
	header := make([]byte, 13)
	if _, err := io.ReadFull(r, header); err != nil {
		return fmt.Errorf("invalid gif header: %w", err)
	}
	out.Write(header)

	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		switch b {
		case 0x21: // Extension block
			label, err := r.ReadByte()
			if err != nil {
				return err
			}

			if label == 0xFE { // Comment Extension
				skipSubBlocks(r)
				continue
			}

			if label == 0xFF { // Application Extension
				app := make([]byte, 11)
				if _, err := io.ReadFull(r, app); err != nil {
					return err
				}
				appName := string(app)
				if appName != "NETSCAPE2.0" && appName != "ANIMEXTS1.0" {
					skipSubBlocks(r)
					continue
				}
				// Write extension introducer, label, block size, and app identifier
				out.WriteByte(0x21)
				out.WriteByte(label)
				out.WriteByte(0x0B)
				out.Write(app)
				copySubBlocks(r, &out)
				continue
			}

			// Preserve other extensions
			out.WriteByte(0x21)
			out.WriteByte(label)
			copySubBlocks(r, &out)

		case 0x2C: // Image Descriptor
			out.WriteByte(b)
			io.Copy(&out, r)
			break

		case 0x3B: // Trailer
			out.WriteByte(b)
			break

		default:
			// Just in case, preserve unknown blocks
			out.WriteByte(b)
		}
	}

	if err := os.WriteFile(path, out.Bytes(), 0644); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}

func skipSubBlocks(r *bytes.Reader) {
	for {
		size, err := r.ReadByte()
		if err != nil || size == 0 {
			break
		}
		r.Seek(int64(size), io.SeekCurrent)
	}
}

func copySubBlocks(r *bytes.Reader, w *bytes.Buffer) {
	for {
		size, err := r.ReadByte()
		if err != nil {
			break
		}
		w.WriteByte(size)
		if size == 0 {
			break
		}
		buf := make([]byte, size)
		if _, err := io.ReadFull(r, buf); err != nil {
			break
		}
		w.Write(buf)
	}
}
