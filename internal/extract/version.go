package extract

import (
	"errors"

	"github.com/benoitmasson/qrcode-demo/internal/detect"
)

// Version gets the QR-code version.
// In our demo, it is not so useful, but helps detecting false positives.
func Version(dots detect.QRCode) (uint, error) {
	if !dots[len(dots)-8][8] {
		return 0, errors.New("dark spot not found")
	}

	verticalVersion, err := verticalVersion(dots)
	if err == nil {
		return verticalVersion, nil
	}
	horizontalVersion, err := horizontalVersion(dots)
	if err == nil {
		return horizontalVersion, nil
	}
	return 0, err
}

// verticalVersion detects alternating dots in the 7-th column, between markers
func verticalVersion(dots detect.QRCode) (uint, error) {
	for index := 7; index < len(dots)-7; index++ {
		if dots[6][index] == dots[6][index-1] {
			return 0, errors.New("invalid vertical version")
		}
	}
	return ((uint(len(dots)) - 14) - 1) / 4, nil
}

// horizontalVersion detects alternating dots in the 7-th row, between markers
func horizontalVersion(dots detect.QRCode) (uint, error) {
	for index := 7; index < len(dots)-7; index++ {
		if dots[index][6] == dots[index-1][6] {
			return 0, errors.New("invalid horizontal version")
		}
	}
	return ((uint(len(dots)) - 14) - 1) / 4, nil
}
