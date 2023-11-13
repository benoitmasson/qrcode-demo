package decode

type ErrorCorrectionLevel uint8

const (
	ErrorCorrectionLevelMedium ErrorCorrectionLevel = iota
	ErrorCorrectionLevelLow
	ErrorCorrectionLevelHigh
	ErrorCorrectionLevelQuartile
)
