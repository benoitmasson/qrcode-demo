package extract

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/benoitmasson/qrcode-demo/internal/decode"
	"github.com/benoitmasson/qrcode-demo/internal/detect"
)

// Inspired from https://www.thonky.com/qr-code-tutorial/format-version-information

const formatMask = 0b101010000010010 // 21522

// Format returns the QR-code "format", i.e. the mask ID used for the data dots
// and the error correction level.
// It uses both occurrences of the format and its error correction code, and returns
// the more likely value among all the encoded values. It fails when the format cannot
// be clearly recovered from the error correction codes.
func Format(dots detect.QRCode) (MaskID, decode.ErrorCorrectionLevel, error) {
	topLeftFormat := topLeftFormat(dots) ^ formatMask
	bottomRightFormat := bottomRightFormat(dots) ^ formatMask
	slog.Debug(fmt.Sprintf("Scanned formats: %015b | %015b", topLeftFormat, bottomRightFormat))

	format1 := uint16(topLeftFormat >> 10)     // first 5 bits
	format2 := uint16(bottomRightFormat >> 10) // first 5 bits
	if format1 != format2 {
		return 0, 0, errors.New("format 1 and format 2 do not match")
	}

	return maskIDFromFormat(format1), errorCorrectionLevelFromFormat(format1), nil
}

func topLeftFormat(dots detect.QRCode) uint16 {
	bits := dots[8][0:6]
	bits = append(bits, dots[8][7:9]...)
	bits = append(bits, dots[7][8], dots[5][8], dots[4][8], dots[3][8], dots[2][8], dots[1][8], dots[0][8])
	return decode.BitsToUint16(bits)
}

func bottomRightFormat(dots detect.QRCode) uint16 {
	l := len(dots)
	bits := []bool{dots[l-1][8], dots[l-2][8], dots[l-3][8], dots[l-4][8], dots[l-5][8], dots[l-6][8], dots[l-7][8]}
	bits = append(bits, dots[8][l-8:]...)
	return decode.BitsToUint16(bits)
}

func errorCorrectionLevelFromFormat(format uint16) decode.ErrorCorrectionLevel {
	return decode.ErrorCorrectionLevel(format >> 3) // use the first 2 bits
}

func maskIDFromFormat(format uint16) MaskID {
	return MaskID(format % (1 << 3)) // use the last 3 bits
}
