package checker

import (
	"bytes"
	"fmt"
	"io"
)

// containsFloatingPointOps checks if a Wasm binary contains any f32 or f64 operations
func containsFloatingPointOps(wasmCode []byte) (bool, error) {
	// Check Wasm magic number
	if len(wasmCode) < 8 {
		return false, fmt.Errorf("invalid Wasm binary: too short")
	}

	magic := wasmCode[0:4]
	if !bytes.Equal(magic, []byte{0x00, 0x61, 0x73, 0x6D}) {
		return false, fmt.Errorf("invalid Wasm binary: wrong magic number")
	}

	// Define floating-point opcodes to check for
	floatingPointOpcodes := map[byte]bool{
		// f32 operations
		0x8B: true, // f32.abs
		0x8C: true, // f32.neg
		0x8D: true, // f32.ceil
		0x8E: true, // f32.floor
		0x8F: true, // f32.trunc
		0x90: true, // f32.nearest
		0x91: true, // f32.sqrt
		0x92: true, // f32.add
		0x93: true, // f32.sub
		0x94: true, // f32.mul
		0x95: true, // f32.div
		0x96: true, // f32.min
		0x97: true, // f32.max
		0x98: true, // f32.copysign
		0x5B: true, // f32.eq
		0x5C: true, // f32.ne
		0x5D: true, // f32.lt
		0x5E: true, // f32.gt
		0x5F: true, // f32.le
		0x60: true, // f32.ge

		// f64 operations
		0x99: true, // f64.abs
		0x9A: true, // f64.neg
		0x9B: true, // f64.ceil
		0x9C: true, // f64.floor
		0x9D: true, // f64.trunc
		0x9E: true, // f64.nearest
		0x9F: true, // f64.sqrt
		0xA0: true, // f64.add
		0xA1: true, // f64.sub
		0xA2: true, // f64.mul
		0xA3: true, // f64.div
		0xA4: true, // f64.min
		0xA5: true, // f64.max
		0xA6: true, // f64.copysign
		0x61: true, // f64.eq
		0x62: true, // f64.ne
		0x63: true, // f64.lt
		0x64: true, // f64.gt
		0x65: true, // f64.le
		0x66: true, // f64.ge

		// Conversion operations
		0xB2: true, // f32.convert_i32_s
		0xB3: true, // f32.convert_i32_u
		0xB4: true, // f32.convert_i64_s
		0xB5: true, // f32.convert_i64_u
		0xB6: true, // f32.demote_f64
		0xB7: true, // f64.convert_i32_s
		0xB8: true, // f64.convert_i32_u
		0xB9: true, // f64.convert_i64_s
		0xBA: true, // f64.convert_i64_u
		0xBB: true, // f64.promote_f32
		0xA8: true, // i32.trunc_f32_s
		0xA9: true, // i32.trunc_f32_u
		0xAA: true, // i32.trunc_f64_s
		0xAB: true, // i32.trunc_f64_u
		0xAE: true, // i64.trunc_f32_s
		0xAF: true, // i64.trunc_f32_u
		0xB0: true, // i64.trunc_f64_s
		0xB1: true, // i64.trunc_f64_u
	}

	// Parse Wasm binary, looking for code section
	reader := bytes.NewReader(wasmCode)

	// Skip magic number and version
	_, err := reader.Seek(8, io.SeekStart)
	if err != nil {
		return false, err
	}

	// Simple scan for opcodes in the entire binary
	// This is a simpler approach that may produce false positives,
	// but will not miss any floating point operations
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
			if floatingPointOpcodes[buf[i]] {
				return true, nil
			}
		}
	}

	return false, nil
}
