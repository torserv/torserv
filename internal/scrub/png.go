package scrub

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

func ScrubPNG(path string) error {
	input, err := os.Open(path)
	if err != nil {
		return err
	}
	defer input.Close()

	var output bytes.Buffer

	// Copy PNG signature (8 bytes)
	sig := make([]byte, 8)
	if _, err := io.ReadFull(input, sig); err != nil {
		return err
	}
	output.Write(sig)

	for {
		var length uint32
		err := binary.Read(input, binary.BigEndian, &length)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunkType := make([]byte, 4)
		if _, err := io.ReadFull(input, chunkType); err != nil {
			return err
		}

		chunkData := make([]byte, length)
		if _, err := io.ReadFull(input, chunkData); err != nil {
			return err
		}

		crc := make([]byte, 4)
		if _, err := io.ReadFull(input, crc); err != nil {
			return err
		}

		typ := string(chunkType)
		if typ == "tEXt" || typ == "iTXt" || typ == "zTXt" || typ == "pHYs" {
			continue // skip metadata chunks
		}

		// Write cleaned chunk
		binary.Write(&output, binary.BigEndian, length)
		output.Write(chunkType)
		output.Write(chunkData)
		output.Write(crc)

		if typ == "IEND" {
			break
		}
	}

	tmpPath := path + ".clean"
	if err := os.WriteFile(tmpPath, output.Bytes(), 0644); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	return nil
}
