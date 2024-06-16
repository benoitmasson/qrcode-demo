package extract

// ReadBits extracts the contents bits from the QR-code, excluding markers and all special dots,
// applying the given mask, and putting everything in the right order.
// See https://www.thonky.com/qr-code-tutorial/module-placement-matrix#step-6-place-the-data-bits
// for a visual explanation.
func ReadBits(dots [][]bool, maskID MaskID) []bool {
	mask := masks[maskID]
	s := len(dots)
	bits := make([]bool, 0, s*len(dots)-8*8*3)

	direction := -1 // start moving up
	for j := s - 1; j >= 1; j -= 2 {
		if j == 6 {
			// vertical version column exception: ignore totally and move to previous column
			j--
		}
		start, end := 0, s
		if direction < 0 {
			start, end = s-1, -1
		}

		for i := start; i != end; i += direction {
			for offset := 0; offset <= 1; offset++ {
				if !isSignificantDot(i, j-offset, s) {
					continue
				}

				bits = append(bits, dots[i][j-offset] != mask(i, j-offset)) // != is the same as XOR for booleans
			}
		}
		direction *= -1 // invert direction
	}

	return bits
}

// isSignificantDot returns whether dot at position (i, j) represents a valid message bit,
// and not a specific pattern (marker, version, …)
func isSignificantDot(i, j, size int) bool {
	if (i <= 8 && j <= 8) ||
		(i <= 8 && j >= size-8) ||
		(i >= size-8 && j <= 8) {
		// ignore 3 finder markers in the corners + format info + dark spot
		return false
	}

	if i >= size-9 && i <= size-5 &&
		j >= size-9 && j <= size-5 {
		// ignore alignment pattern in the bottom-right
		// works only for a single pattern, i.e. version 6 and below
		return false
	}

	if j == 6 || i == 6 {
		// ignore version information
		return false
	}

	return true
}
