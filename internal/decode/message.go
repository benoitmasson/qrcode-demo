package decode

import (
	"bytes"
	"errors"
	"fmt"
)

// Message decodes the given binary contents, of the given length and for the given mode, to text.
// Kanji mode is currently not implemented.
// Error correction is currently not performed, only the raw contents is decoded.
func Message(mode Mode, length uint, contents []bool, errorCorrectionLevel ErrorCorrectionLevel) (string, error) {
	if mode == KanjiMode {
		return "", errors.New("Kanji mode is not supported")
	}

	var blockSize, charactersPerBlock uint
	switch mode {
	case NumericMode:
		blockSize, charactersPerBlock = 10, 3 // 3 numeric characters (0-9) on 10 bits
	case AlphanumericMode:
		blockSize, charactersPerBlock = 11, 2 // 2 alpha-numeric characters (0-9 + uppercase letters + 9 symbols) on 10 bits
	case ByteMode:
		blockSize, charactersPerBlock = 8, 1 // 1 ASCII character on 8 bits
	}
	if len(contents)*int(charactersPerBlock) < int(length*blockSize) {
		return "", fmt.Errorf("missing data in contents, not enough bits to encode the expected %d characters", length)
	}

	var buffer bytes.Buffer
	i := 0
	for nbCharacters := uint(0); nbCharacters < length; nbCharacters += charactersPerBlock {
		if mode == NumericMode {
			if length-nbCharacters == 1 {
				blockSize, charactersPerBlock = 4, 1 // last numeric character encoded on 4 bits
			} else if length-nbCharacters == 2 {
				blockSize, charactersPerBlock = 7, 2 // last 2 numeric characters encoded on 7 bits
			}
		} else if mode == AlphanumericMode && length-nbCharacters == 1 {
			blockSize, charactersPerBlock = 6, 1 // last alphanumeric character encoded on 6 bits
		}

		characters := decodeCharacters(mode, contents[i:i+int(blockSize)], charactersPerBlock)
		buffer.Write(characters)

		i += int(blockSize)
	}

	return buffer.String(), nil
}

const alphanumericCharacters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"

func decodeCharacters(mode Mode, charactersBits []bool, nbCharacters uint) []byte {
	val := BitsToUint16(charactersBits)

	if mode == NumericMode {
		format := fmt.Sprintf("%%0%dd", nbCharacters) // add leading 0's if needed
		return []byte(fmt.Sprintf(format, val))
	} else if mode == AlphanumericMode {
		if nbCharacters == 1 {
			return []byte{alphanumericCharacters[val%45]} // last character, ignore the high-value bits
		}
		return []byte{alphanumericCharacters[val/45], alphanumericCharacters[val%45]}
	} else if mode == ByteMode {
		return []byte{byte(val)}
	}

	return nil
}
