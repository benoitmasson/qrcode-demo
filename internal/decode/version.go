package decode

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

	version, err := verticalVersion(dots)
	if err != nil {
		// fallback to horizontal if not found vertically
		version, err = horizontalVersion(dots)
		return version, err
	}
	return version, nil
}

// verticalVersion detects alternating dots in the 7-th column, between markers
func verticalVersion(dots detect.QRCode) (uint, error) {
	previous := true // black
	for row := 7; row < len(dots)-7; row++ {
		if dots[row][6] == previous {
			return 0, errors.New("version not found in dots")
		}
		previous = !previous
	}
	patternLength := uint((len(dots) - 15) / 2)
	version := (patternLength - 1) / 2
	return version, nil
}

// horizontalVersion detects alternating dots in the 7-th row, between markers
func horizontalVersion(dots detect.QRCode) (uint, error) {
	previous := true // black
	for col := 7; col < len(dots[0])-7; col++ {
		if dots[6][col] == previous {
			return 0, errors.New("version not found in dots")
		}
		previous = !previous
	}
	patternLength := uint((len(dots) - 15) / 2)
	version := (patternLength - 1) / 2
	return version, nil
}
