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
		start, end := 0, s
		if direction < 0 {
			start, end = s-1, -1
		}
		for i := start; i != end; i += direction {
			if j == 6 || i == 6 {
				// ignore version information
				if j == 6 {
					// vertical version column exception: ignore totally and move to previous column
					j--
				}
				continue
			}
			if (i <= 8 && j <= 8) ||
				(i <= 8 && j >= s-8) ||
				(i >= s-8 && j <= 8) {
				// ignore 3 finder markers in the corners + format info
				continue
			}
			if i >= s-9 && i <= s-5 &&
				j >= s-8 && j <= s-5 {
				// ignore alignment pattern in the bottom-right
				continue
			}

			if i >= s-9 && i <= s-5 &&
				j == s-9 {
				// ignore alignment pattern in the bottom-right - special case where 2-bits column is half-cut
			} else {
				bits = append(bits, dots[i][j] != mask(i, j)) // != is the same as XOR for booleans
			}
			bits = append(bits, dots[i][j-1] != mask(i, j-1)) // != is the same as XOR for booleans
		}
		direction *= -1 // invert direction
	}

	return bits
}
