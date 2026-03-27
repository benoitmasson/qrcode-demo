package decode

import (
	"slices"
	"testing"
)

func TestIntSliceToBits(t *testing.T) {
	type test struct {
		name         string
		ints         []int
		expectedBits []bool
	}
	tests := []test{
		{
			name:         "nil",
			ints:         nil,
			expectedBits: []bool{},
		},
		{
			name:         "empty",
			ints:         []int{},
			expectedBits: []bool{},
		},
		{
			name:         "0 (all false)",
			ints:         []int{0},
			expectedBits: []bool{false, false, false, false, false, false, false, false},
		},
		{
			name:         "255 (all true)",
			ints:         []int{255},
			expectedBits: []bool{true, true, true, true, true, true, true, true},
		},
		{
			name:         "129 (start and end with true)",
			ints:         []int{129},
			expectedBits: []bool{true, false, false, false, false, false, false, true},
		},
		{
			name:         "126 (start and end with false)",
			ints:         []int{126},
			expectedBits: []bool{false, true, true, true, true, true, true, false},
		},
		{
			name:         "128 (start with true and end with false)",
			ints:         []int{128},
			expectedBits: []bool{true, false, false, false, false, false, false, false},
		},
		{
			name:         "1 (start with false and end with true)",
			ints:         []int{1},
			expectedBits: []bool{false, false, false, false, false, false, false, true},
		},
		{
			name:         "2 ints - 16 digits",
			ints:         []int{255, 255},
			expectedBits: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		},
		{
			name:         "letters 'https'",
			ints:         []int{105, 116, 116, 112, 115},
			expectedBits: []bool{false, true, true, false, true, false, false, true, false, true, true, true, false, true, false, false, false, true, true, true, false, true, false, false, false, true, true, true, false, false, false, false, false, true, true, true, false, false, true, true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualBits := intSliceToBits(test.ints)
			if !slices.Equal(actualBits, test.expectedBits) {
				t.Errorf("expected %v but got %v", test.expectedBits, actualBits)
			}
		})
	}
}
