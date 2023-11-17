package decode

import "testing"

const (
	_0 = false
	_1 = true
)

func TestMessage_Numeric(t *testing.T) {
	type test struct {
		name            string
		length          uint
		contents        []bool
		expectedMessage string
	}
	tests := []test{
		{
			name:   "multiple of 3 digits",
			length: 6,
			contents: []bool{
				_0, _0, _0, _1, _1, _1, _1, _0, _1, _1,
				_0, _1, _1, _1, _0, _0, _1, _0, _0, _0,
			},
			expectedMessage: "123456",
		},
		{
			name:   "multiple of 3 digits with leading 0's",
			length: 6,
			contents: []bool{
				_0, _0, _0, _1, _1, _1, _1, _0, _1, _1,
				_0, _0, _0, _0, _0, _0, _0, _1, _1, _1,
			},
			expectedMessage: "123007",
		},
		{
			name:   "multiple of 3 digits plus 1",
			length: 7,
			contents: []bool{
				_0, _0, _0, _1, _1, _1, _1, _0, _1, _1,
				_0, _1, _1, _1, _0, _0, _1, _0, _0, _0,
				_0, _1, _1, _1,
			},
			expectedMessage: "1234567",
		},
		{
			name:   "multiple of 3 digits plus 1 with 0",
			length: 7,
			contents: []bool{
				_0, _0, _0, _1, _1, _1, _1, _0, _1, _1,
				_0, _1, _1, _1, _0, _0, _1, _0, _0, _0,
				_0, _0, _0, _0,
			},
			expectedMessage: "1234560",
		},
		{
			name:   "multiple of 3 digits plus 2",
			length: 8,
			contents: []bool{
				_0, _0, _0, _1, _1, _1, _1, _0, _1, _1,
				_0, _1, _1, _1, _0, _0, _1, _0, _0, _0,
				_1, _0, _0, _1, _1, _1, _0,
			},
			expectedMessage: "12345678",
		},
		{
			name:   "multiple of 3 digits plus 2 with leading 0",
			length: 8,
			contents: []bool{
				_0, _0, _0, _1, _1, _1, _1, _0, _1, _1,
				_0, _1, _1, _1, _0, _0, _1, _0, _0, _0,
				_0, _0, _0, _0, _1, _1, _1,
			},
			expectedMessage: "12345607",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualMessage, err := Message(NumericMode, test.length, test.contents, ErrorCorrectionLevelLow)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if actualMessage != test.expectedMessage {
				t.Errorf("expected %q but got %q", test.expectedMessage, actualMessage)
			}
		})
	}
}

func TestMessage_Alphanumeric(t *testing.T) {
	type test struct {
		name            string
		length          uint
		contents        []bool
		expectedMessage string
	}
	tests := []test{
		{
			name:   "multiple of 2 characters",
			length: 6,
			contents: []bool{
				_0, _0, _1, _1, _1, _0, _0, _0, _0, _1, _1,
				_1, _1, _1, _0, _0, _0, _0, _1, _0, _1, _0,
				_1, _1, _1, _1, _0, _1, _1, _1, _1, _1, _1,
			},
			expectedMessage: "A1+2:3",
		},
		{
			name:   "multiple of 2 characters plus 1",
			length: 7,
			contents: []bool{
				_0, _0, _1, _1, _1, _1, _1, _0, _0, _0, _0,
				_1, _1, _1, _0, _0, _0, _0, _1, _0, _1, _0,
				_1, _1, _1, _1, _0, _1, _1, _1, _1, _0, _0,
				_0, _0, _0, _0, _1, _1,
			},
			expectedMessage: "B1+2:03",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualMessage, err := Message(AlphanumericMode, test.length, test.contents, ErrorCorrectionLevelLow)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if actualMessage != test.expectedMessage {
				t.Errorf("expected %q but got %q", test.expectedMessage, actualMessage)
			}
		})
	}
}

func TestMessage_Byte(t *testing.T) {
	type test struct {
		name            string
		length          uint
		contents        []bool
		expectedMessage string
	}
	tests := []test{
		{
			name:   "url",
			length: 11,
			contents: []bool{
				_0, _1, _1, _0, _1, _0, _0, _0,
				_0, _1, _1, _1, _0, _1, _0, _0,
				_0, _1, _1, _1, _0, _1, _0, _0,
				_0, _1, _1, _1, _0, _0, _0, _0,
				_0, _1, _1, _1, _0, _0, _1, _1,
				_0, _0, _1, _1, _1, _0, _1, _0,
				_0, _0, _1, _0, _1, _1, _1, _1,
				_0, _0, _1, _0, _1, _1, _1, _1,
				_0, _1, _1, _1, _0, _1, _1, _1,
				_0, _1, _1, _1, _0, _1, _1, _1,
				_0, _1, _1, _1, _0, _1, _1, _1,
			},
			expectedMessage: "https://www",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualMessage, err := Message(ByteMode, test.length, test.contents, ErrorCorrectionLevelLow)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if actualMessage != test.expectedMessage {
				t.Errorf("expected %q but got %q", test.expectedMessage, actualMessage)
			}
		})
	}
}
