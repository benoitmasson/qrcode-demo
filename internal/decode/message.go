package decode

import (
	"bytes"
	"fmt"
)

// Message decodes the given binary contents, of the given length and for the given mode, to text.
// Kanji, Numeric and Alphanumeric modes are not implemented.
func Message(mode Mode, length uint, contents []bool) (string, error) {
	if mode != ByteMode {
		return "", fmt.Errorf("%s mode is not supported", mode.String())
	}

	var buffer bytes.Buffer
	for nb := range length {
		i := nb * 8
		char := byte(BitsToUint16(contents[i : i+8]))
		buffer.WriteByte(char)
	}

	return buffer.String(), nil
}
