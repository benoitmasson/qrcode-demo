package extract

import (
	"strconv"
	"testing"
)

func TestBitsToUint16(t *testing.T) {
	type test struct {
		name           string
		bits           []bool
		expectedUint16 uint16
	}
	tests := []test{
		{
			name:           "only 0s",
			bits:           []bool{false, false, false, false, false},
			expectedUint16: 0b00000,
		},
		{
			name:           "only 1s",
			bits:           []bool{true, true, true, true, true},
			expectedUint16: 0b11111,
		},
		{
			name:           "start and end with 1",
			bits:           []bool{true, false, false, false, true},
			expectedUint16: 0b10001,
		},
		{
			name:           "start and end with 0",
			bits:           []bool{false, true, true, true, false},
			expectedUint16: 0b01110,
		},
		{
			name:           "start with 1 and end with 0",
			bits:           []bool{true, false, true, true, false},
			expectedUint16: 0b10110,
		},
		{
			name:           "start with 0 and end with 1",
			bits:           []bool{false, true, true, false, true},
			expectedUint16: 0b01101,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualUint16 := bitsToUint16(test.bits)
			if actualUint16 != test.expectedUint16 {
				t.Errorf("expected %016s but got %016s", strconv.FormatInt(int64(test.expectedUint16), 2), strconv.FormatInt(int64(actualUint16), 2))
			}
		})
	}
}
