package extract

import (
	"testing"
)

const (
	_0 = false
	_1 = true
)

var sampleDots = [][]bool{
	{_1, _1, _1, _1, _1, _1, _1, _0, _1, _1, _0, _0, _0, _1, _0, _0, _1, _0, _1, _1, _1, _1, _1, _1, _1},
	{_1, _0, _0, _0, _0, _0, _1, _0, _1, _1, _1, _0, _1, _1, _1, _0, _0, _0, _1, _0, _0, _0, _0, _0, _1},
	{_1, _0, _1, _1, _1, _0, _1, _0, _0, _0, _1, _0, _1, _1, _0, _0, _1, _0, _1, _0, _1, _1, _1, _0, _1},
	{_1, _0, _1, _1, _1, _0, _1, _0, _1, _0, _0, _1, _1, _0, _1, _0, _1, _0, _1, _0, _1, _1, _1, _0, _1},
	{_1, _0, _1, _1, _1, _0, _1, _0, _0, _1, _0, _0, _0, _1, _1, _0, _1, _0, _1, _0, _1, _1, _1, _0, _1},
	{_1, _0, _0, _0, _0, _0, _1, _0, _0, _1, _1, _0, _1, _0, _0, _0, _1, _0, _1, _0, _0, _0, _0, _0, _1},
	{_1, _1, _1, _1, _1, _1, _1, _0, _1, _0, _1, _0, _1, _0, _1, _0, _1, _0, _1, _1, _1, _1, _1, _1, _1},
	{_0, _0, _0, _0, _0, _0, _0, _0, _1, _0, _0, _1, _1, _0, _0, _0, _1, _0, _0, _0, _0, _0, _0, _0, _0},
	{_1, _0, _1, _1, _0, _1, _1, _0, _1, _0, _0, _1, _0, _1, _1, _0, _1, _0, _1, _0, _0, _1, _0, _1, _1},
	{_1, _0, _1, _0, _1, _1, _0, _0, _1, _1, _0, _1, _1, _0, _0, _1, _1, _1, _0, _1, _0, _0, _0, _1, _0},
	{_1, _0, _1, _1, _1, _0, _1, _1, _0, _0, _0, _1, _1, _1, _1, _0, _1, _0, _0, _1, _0, _0, _0, _0, _0},
	{_1, _1, _0, _0, _1, _1, _0, _1, _0, _1, _1, _0, _0, _1, _0, _1, _0, _0, _1, _0, _1, _1, _1, _0, _0},
	{_0, _0, _1, _1, _0, _1, _1, _0, _0, _0, _1, _1, _0, _1, _0, _0, _0, _1, _1, _0, _1, _0, _1, _1, _1},
	{_0, _0, _1, _1, _0, _0, _0, _0, _1, _1, _1, _1, _0, _1, _0, _1, _0, _1, _1, _1, _1, _0, _0, _0, _1},
	{_0, _1, _1, _0, _1, _1, _1, _0, _1, _0, _1, _0, _1, _1, _0, _0, _0, _0, _1, _0, _1, _0, _1, _1, _0},
	{_1, _0, _1, _1, _1, _0, _0, _1, _1, _1, _1, _0, _1, _0, _1, _0, _1, _0, _1, _1, _1, _0, _0, _0, _1},
	{_0, _0, _0, _1, _1, _1, _1, _1, _1, _1, _0, _1, _0, _0, _0, _1, _1, _1, _1, _1, _1, _1, _1, _1, _1},
	{_0, _0, _0, _0, _0, _0, _0, _0, _1, _1, _1, _0, _1, _1, _0, _0, _1, _0, _0, _0, _1, _0, _1, _0, _1},
	{_1, _1, _1, _1, _1, _1, _1, _0, _1, _0, _0, _1, _1, _0, _1, _0, _1, _0, _1, _0, _1, _0, _1, _1, _1},
	{_1, _0, _0, _0, _0, _0, _1, _0, _1, _1, _0, _0, _1, _0, _1, _1, _1, _0, _0, _0, _1, _0, _0, _1, _1},
	{_1, _0, _1, _1, _1, _0, _1, _0, _0, _1, _0, _1, _1, _1, _0, _0, _1, _1, _1, _1, _1, _1, _0, _1, _0},
	{_1, _0, _1, _1, _1, _0, _1, _0, _1, _0, _1, _0, _1, _1, _0, _1, _0, _0, _1, _0, _1, _1, _1, _1, _1},
	{_1, _0, _1, _1, _1, _0, _1, _0, _1, _0, _0, _1, _1, _1, _0, _1, _0, _1, _1, _0, _1, _0, _1, _1, _0},
	{_1, _0, _0, _0, _0, _0, _1, _0, _0, _0, _1, _0, _0, _0, _0, _1, _0, _1, _1, _0, _1, _0, _1, _0, _0},
	{_1, _1, _1, _1, _1, _1, _1, _0, _1, _1, _0, _0, _1, _0, _1, _0, _0, _0, _0, _1, _1, _1, _1, _1, _1},
}

func TestReadbits(t *testing.T) {
	output := []bool{
		_0, _1, _0, _0, _0, _0, _0, _1, _0, _1,
		_1, _0, _0, _1, _1, _0, _1, _0, _0, _0,
		_0, _1, _1, _1, _0, _1, _0, _0, _0, _1,
		_1, _1, _0, _1, _0, _0, _0, _1, _1, _1,
		_0, _0, _0, _0, _0, _1, _1, _1, _0, _0,
		_1, _1, _0, _0, _1, _1, _1, _0, _1, _0,
		_0, _0, _1, _0, _1, _1, _1, _1, _0, _0,
		_1, _0, _1, _1, _1, _1, _0, _1, _1, _0,
		_1, _1, _1, _1, _0, _1, _1, _1, _0, _1,
		_1, _0, _0, _1, _1, _0, _1, _0, _0, _0,
		_0, _0, _1, _0, _1, _1, _1, _0, _0, _1,
		_1, _1, _0, _1, _0, _0, _0, _1, _1, _0,
		_1, _1, _1, _1, _0, _0, _1, _0, _1, _1,
		_1, _1, _0, _1, _0, _0, _0, _1, _0, _0,
		_0, _1, _0, _1, _1, _0, _0, _0, _0, _1,
		_1, _0, _1, _0, _1, _0, _0, _1, _0, _0,
		_1, _0, _1, _1, _0, _1, _0, _0, _0, _0,
		_1, _0, _0, _0, _1, _1, _1, _0, _0, _1,
		_0, _1, _0, _0, _0, _1, _0, _0, _0, _0,
		_0, _0, _1, _1, _1, _0, _1, _1, _0, _0,
		_0, _0, _0, _1, _0, _0, _0, _1, _1, _1,
		_1, _0, _1, _1, _0, _0, _0, _0, _0, _1,
		_0, _0, _0, _1, _1, _0, _0, _1, _0, _0,
		_1, _0, _1, _1, _0, _1, _0, _0, _1, _0,
		_1, _1, _0, _1, _1, _1, _0, _0, _0, _1,
		_1, _0, _0, _0, _1, _0, _1, _0, _1, _0,
		_0, _1, _0, _0, _0, _0, _0, _0, _1, _1,
		_1, _1, _1, _0, _0, _1, _0, _0, _1, _0,
		_1, _0, _1, _0, _1, _1, _1, _1, _1, _0,
		_0, _0, _0, _0, _0, _0, _1, _1, _1, _1,
		_1, _0, _0, _0, _0, _0, _0, _1, _1, _1,
		_0, _1, _1, _1, _1, _1, _0, _1, _0, _1,
		_0, _1, _0, _0, _1, _0, _1, _1, _1, _0,
		_1, _0, _1, _1, _0, _0, _1, _0, _0, _1,
		_0, _1, _1, _0, _0, _0, _1, _0, _1, _0,
		_1, _0, _0, _0, _0, _0, _0, _0, _0,
	}

	bits := ReadBits(sampleDots, MaskID(3))

	compare := compareSlices(bits, output)
	if compare >= 0 {
		t.Errorf("bits mismatch at position %d", compare)
	}
}

func compareSlices[T comparable](a, b []T) int {
	if len(a) != len(b) {
		if len(a) < len(b) {
			return len(a)
		} else {
			return len(b)
		}
	}
	for i := range a {
		if a[i] != b[i] {
			return i
		}
	}
	return -1
}