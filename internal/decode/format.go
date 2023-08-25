package decode

import (
	"errors"

	"github.com/benoitmasson/qrcode-demo/internal/detect"
)

// Inspired from https://www.thonky.com/qr-code-tutorial/format-version-information

const formatMask = 0b101010000010010 // 21522

// Format returns the QR-code "format", i.e. the mask ID used for the data dots
// and the error correction level.
// It uses both occurrences of the format and its error correction code, and returns
// the more likely value among all the encoded values. It fails when the format cannot
// be clearly recovered from the error correction codes.
func Format(dots detect.QRCode) (MaskID, ErrorCorrectionLevel, error) {
	format1 := topLeftFormat(dots)
	format2 := bottomRightFormat(dots)
	// fmt.Printf("scanned formats:\n- %015s\n- %015s\n", strconv.FormatInt(int64(format1), 2), strconv.FormatInt(int64(format2), 2))

	possibleFormats := make(map[uint16]int, 10)

	format1 ^= formatMask
	rawFormat1 := uint16(format1 >> 10) // first 5 bits
	possibleFormats[rawFormat1]++
	eccFormat1 := uint16(format1 % (1 << 10)) // last 10 bits
	for _, format := range decodeFormat(eccFormat1) {
		possibleFormats[format]++
	}

	format2 ^= formatMask
	rawFormat2 := uint16(format2 >> 10) // first 5 bits
	possibleFormats[rawFormat2]++
	eccFormat2 := uint16(format2 % (1 << 10)) // last 10 bits
	for _, format := range decodeFormat(eccFormat2) {
		possibleFormats[format]++
	}

	format, err := findMoreFrequent(possibleFormats)
	if err != nil {
		return 0, 0, err
	}
	// fmt.Printf("selected format: %05s\n", strconv.FormatInt(int64(format), 2))

	return maskIDFromFormat(format), errorCorrectionLevelFromFormat(format), nil
}

func topLeftFormat(dots detect.QRCode) uint16 {
	bits := dots[8][0:6]
	bits = append(bits, dots[8][7:9]...)
	bits = append(bits, dots[7][8], dots[5][8], dots[4][8], dots[3][8], dots[2][8], dots[1][8], dots[0][8])
	return bitsToUint16(bits)
}

func bottomRightFormat(dots detect.QRCode) uint16 {
	l := len(dots)
	bits := []bool{dots[l-1][8], dots[l-2][8], dots[l-3][8], dots[l-4][8], dots[l-5][8], dots[l-6][8], dots[l-7][8]}
	bits = append(bits, dots[8][l-8:]...)
	return bitsToUint16(bits)
}

func findMoreFrequent(m map[uint16]int) (uint16, error) {
	bestValues := make([]uint16, 0, len(m))
	bestCount := 0
	for val, count := range m {
		if count == bestCount {
			bestValues = append(bestValues, val)
			continue
		}
		if count > bestCount {
			bestCount = count
			bestValues = make([]uint16, 0, len(m))
			bestValues = append(bestValues, val)
			continue
		}
	}

	if len(bestValues) > 1 {
		return 0, errors.New("ambiguous value for format")
	}
	return bestValues[0], nil
}

func errorCorrectionLevelFromFormat(format uint16) ErrorCorrectionLevel {
	return ErrorCorrectionLevel(format >> 3) // use the first 2 bits
}

func maskIDFromFormat(format uint16) MaskID {
	return MaskID(format % (1 << 3)) // use the last 3 bits
}

func bitsToUint16(bits []bool) uint16 {
	i := uint16(0)
	for _, bit := range bits {
		i <<= 1
		if bit {
			i++
		}
	}
	return i
}
