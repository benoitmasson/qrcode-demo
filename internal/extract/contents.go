package extract

// ReadBits extracts the contents bits from the QR-code, excluding markers and all special dots,
// applying the given mask, and putting everything in the right order.
// See https://www.thonky.com/qr-code-tutorial/module-placement-matrix#step-6-place-the-data-bits
// for a visual explanation.
func ReadBits(dots [][]bool, maskID MaskID) []bool {
	mask := masks[maskID]
	size := len(dots)
	output := make([]bool, 0, size*size)

	for col := size - 1; col >= 0; col -= 2 {
		// read from bottom to top
		for row := size - 1; row >= 0; row-- {
			if isSignificantDot(row, col, size) {
				output = append(output, dots[row][col] != mask(row, col))
			}

			if isSignificantDot(row, col-1, size) {
				output = append(output, dots[row][col-1] != mask(row, col-1))
			}
		}

		col -= 2
		if col == 6 {
			// vertical version column exception: ignore totally and move to previous column
			col--
		}

		// read from top to bottom
		for row := 0; row < size; row++ {
			if isSignificantDot(row, col, size) {
				output = append(output, dots[row][col] != mask(row, col))
			}

			if isSignificantDot(row, col-1, size) {
				output = append(output, dots[row][col-1] != mask(row, col-1))
			}
		}
	}

	return output
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
