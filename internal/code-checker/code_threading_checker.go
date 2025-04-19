package checker

import (
	"bytes"
	"fmt"
	"io"
)

func containsThreadingOps(wasmCode []byte) (bool, error) {
	// Check Wasm magic number
	if len(wasmCode) < 8 {
		return false, fmt.Errorf("invalid Wasm binary: too short")
	}

	magic := wasmCode[0:4]
	if !bytes.Equal(magic, []byte{0x00, 0x61, 0x73, 0x6D}) {
		return false, fmt.Errorf("invalid Wasm binary: wrong magic number")
	}

	reader := bytes.NewReader(wasmCode)

	// Skip magic number and version
	_, err := reader.Seek(8, io.SeekStart)
	if err != nil {
		return false, err
	}

	// 1. Scan for atomic operations prefixed with 0xFE
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}

		for i := 0; i < n; i++ {
			if buf[i] == 0xFE {
				return true, nil // Found atomic operation prefix
			}
		}
	}

	// 2. Check for shared memory flag in memory section
	// Reset to after magic number and version
	_, err = reader.Seek(8, io.SeekStart)
	if err != nil {
		return false, err
	}

	// Iterate through all sections
	for {
		// Read section ID
		sectionID, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}

		// Read section size (LEB128 encoded)
		sectionSize, err := readLEB128(reader)
		if err != nil {
			return false, err
		}

		// If this is the memory section (section ID = 5), check for shared flag
		if sectionID == 5 {
			// Read number of memory definitions
			memCount, err := readLEB128(reader)
			if err != nil {
				return false, err
			}

			for i := uint64(0); i < memCount; i++ {
				// Read memory type
				memType, err := reader.ReadByte()
				if err != nil {
					return false, err
				}

				// Check for shared bit (0x02)
				if memType&0x02 != 0 {
					return true, nil // Found shared memory
				}

				// Skip over memory limits
				// Handle minimum value
				_, err = readLEB128(reader)
				if err != nil {
					return false, err
				}

				// If maximum exists, handle maximum value
				if memType&0x01 != 0 {
					_, err = readLEB128(reader)
					if err != nil {
						return false, err
					}
				}
			}

			break // Finished processing memory section
		} else {
			// Skip other sections
			_, err = reader.Seek(int64(sectionSize), io.SeekCurrent)
			if err != nil {
				return false, err
			}
		}
	}

	return false, nil
}

// Read an unsigned LEB128 encoded integer
func readLEB128(r *bytes.Reader) (uint64, error) {
	var result uint64
	var shift uint

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		// Add the 7 bits from this byte to the result
		result |= uint64(b&0x7F) << shift
		shift += 7

		// If the continuation bit is not set, we're done
		if b&0x80 == 0 {
			break
		}
	}

	return result, nil
}
