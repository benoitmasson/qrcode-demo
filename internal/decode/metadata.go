package decode

import (
	"fmt"
)

type Mode uint8

const (
	NumericMode      Mode = 0b0001
	AlphanumericMode Mode = 0b0010
	ByteMode         Mode = 0b0100
	KanjiMode        Mode = 0b1000
)

func (m Mode) String() string {
	switch m {
	case NumericMode:
		return "Numeric"
	case AlphanumericMode:
		return "Alphanumeric"
	case ByteMode:
		return "Byte"
	case KanjiMode:
		return "Kanji"
	}
	return fmt.Sprintf("Unknown(%b)", m)
}

// GetMode extracts the contents type from the header.
// The first 4 bits are used.
func GetMode(bits []bool) Mode {
	// TODO (3.2): read mode from bits sequence
	return Mode(0)
}

// GetContentLength extracts the contents length from the header.
// After removing the first 4 bits (mode), the number of bits used to encode the length
// is given by the QR-code version and mode.
// Also, the contents bits are returned after trimming the header.
// See https://www.thonky.com/qr-code-tutorial/data-encoding#step-4-add-the-character-count-indicator
func GetContentLength(bits []bool, version uint, mode Mode, errorCorrectionLevel ErrorCorrectionLevel) (uint, []bool, error) {
	// TODO (3.2): read message length from bits sequence
	return 0, bits, nil
}

func lengthBytes(version uint, mode Mode) int {
	if version <= 9 {
		switch mode {
		case NumericMode:
			return 10
		case AlphanumericMode:
			return 9
		case ByteMode:
			return 8
		case KanjiMode:
			return 8
		default:
			return 0
		}
	}

	if version <= 26 {
		switch mode {
		case NumericMode:
			return 12
		case AlphanumericMode:
			return 11
		case ByteMode:
			return 16
		case KanjiMode:
			return 10
		default:
			return 0
		}
	}

	if version <= 40 {
		switch mode {
		case NumericMode:
			return 14
		case AlphanumericMode:
			return 13
		case ByteMode:
			return 16
		case KanjiMode:
			return 12
		default:
			return 0
		}
	}

	return 0
}
