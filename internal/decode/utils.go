package decode

// BitsToUint16 converts a slice of booleans to the corresponding base-2 encoded integer.
// e.g.
// [true, false] => 10 (base 2) => 2 (uint16)
// [false, true] => 01 (base 2) => 1 (uint16)
func BitsToUint16(bits []bool) uint16 {
	i := uint16(0)
	for _, bit := range bits {
		i <<= 1
		if bit {
			i++
		}
	}
	return i
}
