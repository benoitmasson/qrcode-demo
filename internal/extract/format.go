package extract

import (
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
	// TODO (2.2): extract mask ID and error correction level from format
	return MaskID(0), decode.ErrorCorrectionLevel(0), nil
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
