package checker

func ContainUndeterminsticOps(wasmCode []byte) (bool, error) {
	// Check for floating point operations
	floatingPoint, err := containsFloatingPointOps(wasmCode)
	if err != nil {
		return false, err
	}
	if floatingPoint {
		return true, nil
	}

	// Check for SIMD operations
	simd, err := containsSIMDOps(wasmCode)
	if err != nil {
		return false, err
	}
	if simd {
		return true, nil
	}

	// Check for threading operations
	threading, err := containsThreadingOps(wasmCode)
	if err != nil {
		return false, err
	}
	if threading {
		return true, nil
	}

	return false, nil
}
