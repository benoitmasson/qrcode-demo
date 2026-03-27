package extract

import (
	"github.com/benoitmasson/qrcode-demo/internal/detect"
)

// Version gets the QR-code version.
// In our demo, it is not so useful, but helps detecting false positives.
func Version(dots detect.QRCode) (uint, error) {
	// TODO (2.1): check dark spot and read version
	return 0, nil
}

// verticalVersion detects alternating dots in the 7-th column, between markers
func verticalVersion(dots detect.QRCode) (uint, error) {
	// TODO (2.1): read vertical (7th column) version
	return 0, nil
}

// horizontalVersion detects alternating dots in the 7-th row, between markers
func horizontalVersion(dots detect.QRCode) (uint, error) {
	// TODO (2.1): read horizontal (7th row) version
	return 0, nil
}
