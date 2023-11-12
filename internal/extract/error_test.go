package extract

import (
	"strconv"
	"testing"
)

func TestComputeFormatRemainder(t *testing.T) {
	type test struct {
		name              string
		format            uint16
		expectedRemainder uint16
	}
	tests := []test{
		{
			name:              "nominal test",
			format:            0b01100,
			expectedRemainder: 0b1000111101,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualRemainder := computeFormatRemainder(test.format)
			if actualRemainder != test.expectedRemainder {
				t.Errorf("expected %010s but got %010s", strconv.FormatInt(int64(test.expectedRemainder), 2), strconv.FormatInt(int64(actualRemainder), 2))
			}
		})
	}
}
