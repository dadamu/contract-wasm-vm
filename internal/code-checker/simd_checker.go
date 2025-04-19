package checker

import (
	"bytes"
	"fmt"
	"io"
)

// containsSIMDOps checks if a Wasm binary contains any SIMD operations
func containsSIMDOps(wasmCode []byte) (bool, error) {
	// Check Wasm magic number
	if len(wasmCode) < 8 {
		return false, fmt.Errorf("invalid Wasm binary: too short")
	}

	magic := wasmCode[0:4]
	if !bytes.Equal(magic, []byte{0x00, 0x61, 0x73, 0x6D}) {
		return false, fmt.Errorf("invalid Wasm binary: wrong magic number")
	}

	// Check for SIMD in the type section
	// First, look for the type section (section ID = 1)
	reader := bytes.NewReader(wasmCode)

	// Skip magic number and version
	_, err := reader.Seek(8, io.SeekStart)
	if err != nil {
		return false, err
	}

	// Scan for SIMD opcodes
	// SIMD instructions use the 0xFD prefix
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}

		// Look for the SIMD prefix (0xFD)
		for i := 0; i < n; i++ {
			if buf[i] == 0xFD {
				return true, nil
			}
		}
	}

	// Check if v128 type is used in the types or function signatures
	// Vector type v128 is represented by the value 0x7B in the binary format
	reader.Seek(8, io.SeekStart) // Reset to beginning after magic+version

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}

		for i := 0; i < n; i++ {
			if buf[i] == 0x7B { // v128 type identifier
				return true, nil
			}
		}
	}

	return false, nil
}
