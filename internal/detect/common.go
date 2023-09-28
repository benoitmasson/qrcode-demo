package detect

// nearlyEquals returns whether i == j, with a given tolerance
func nearlyEquals(i, j, tolerance int) bool {
	return j-i >= -tolerance && j-i <= tolerance
}
