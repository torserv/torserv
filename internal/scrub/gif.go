package scrub

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// ScrubGIF removes comment extensions and non-animation application extensions from a GIF file.
// It retains only essential blocks and specific app extensions like "NETSCAPE2.0" and "ANIMEXTS1.0".
func ScrubGIF(path string) error {
	// Read the full file into memory
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	var out bytes.Buffer
	r := bytes.NewReader(data)

	// Copy GIF Header (6 bytes) + Logical Screen Descriptor (7 bytes)
	header := make([]byte, 13)
	if _, err := io.ReadFull(r, header); err != nil {
		return fmt.Errorf("invalid gif header: %w", err)
	}
	out.Write(header)

	// Process each block in the GIF
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		switch b {
		case 0x21: // Extension Introducer
			label, err := r.ReadByte()
			if err != nil {
				return err
			}

			switch label {
			case 0xFE: // Comment Extension — skip entirely
				skipSubBlocks(r)
				continue

			case 0xFF: // Application Extension
				app := make([]byte, 11)
				if _, err := io.ReadFull(r, app); err != nil {
					return err
				}
				appName := string(app)

				// Skip any non-animation-related application extensions
				if appName != "NETSCAPE2.0" && appName != "ANIMEXTS1.0" {
					skipSubBlocks(r)
					continue
				}

				// Write the application extension block
				out.WriteByte(0x21)
				out.WriteByte(label)
				out.WriteByte(0x0B) // Block size
				out.Write(app)
				copySubBlocks(r, &out)
				continue

			default:
				// Preserve other extensions (e.g., graphic control extension)
				out.WriteByte(0x21)
				out.WriteByte(label)
				copySubBlocks(r, &out)
			}

		case 0x2C: // Image Descriptor — start of an image frame
			out.WriteByte(b)
			io.Copy(&out, r) // Copy rest of file (assumes well-formed GIF)
			break

		case 0x3B: // Trailer — end of GIF
			out.WriteByte(b)
			break

		default:
			// Unknown or unsupported blocks — preserve
			out.WriteByte(b)
		}
	}

	// Write cleaned GIF back to file
	if err := os.WriteFile(path, out.Bytes(), 0644); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}

// skipSubBlocks advances the reader past a sequence of GIF sub-blocks (used for skipping comments/app data).
func skipSubBlocks(r *bytes.Reader) {
	for {
		size, err := r.ReadByte()
		if err != nil || size == 0 {
			break
		}
		r.Seek(int64(size), io.SeekCurrent)
	}
}

// copySubBlocks reads sub-blocks from the reader and writes them to the buffer.
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
