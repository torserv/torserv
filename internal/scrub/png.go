package scrub

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// ScrubPNG removes metadata chunks from a PNG file by rewriting it without tEXt, iTXt, zTXt, and pHYs chunks.
// The scrubbed image replaces the original file.
func ScrubPNG(path string) error {
	// Open the original PNG file
	input, err := os.Open(path)
	if err != nil {
		return err
	}
	defer input.Close()

	var output bytes.Buffer

	// Read and write the PNG signature (8 bytes)
	sig := make([]byte, 8)
	if _, err := io.ReadFull(input, sig); err != nil {
		return err
	}
	output.Write(sig)

	// Process each PNG chunk
	for {
		var length uint32
		// Read the length of the chunk data
		err := binary.Read(input, binary.BigEndian, &length)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Read chunk type (4 bytes)
		chunkType := make([]byte, 4)
		if _, err := io.ReadFull(input, chunkType); err != nil {
			return err
		}

		// Read chunk data
		chunkData := make([]byte, length)
		if _, err := io.ReadFull(input, chunkData); err != nil {
			return err
		}

		// Read CRC (4 bytes)
		crc := make([]byte, 4)
		if _, err := io.ReadFull(input, crc); err != nil {
			return err
		}

		typ := string(chunkType)
		// Skip metadata chunks
		if typ == "tEXt" || typ == "iTXt" || typ == "zTXt" || typ == "pHYs" {
			continue
		}

		// Write the chunk to output
		binary.Write(&output, binary.BigEndian, length)
		output.Write(chunkType)
		output.Write(chunkData)
		output.Write(crc)

		// Stop processing after IEND chunk
		if typ == "IEND" {
			break
		}
	}

	// Write scrubbed content to a temporary file
	tmpPath := fmt.Sprintf("%s.clean", path)
	if err := os.WriteFile(tmpPath, output.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write cleaned PNG to %s: %w", tmpPath, err)
	}

	// Replace the original file
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	return nil
}
