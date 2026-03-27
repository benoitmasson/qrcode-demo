package decode

import "testing"

const (
	_0 = false
	_1 = true
)

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
			actualMessage, err := Message(ByteMode, test.length, test.contents)
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
