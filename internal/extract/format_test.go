package extract

import (
	"testing"

	"github.com/benoitmasson/qrcode-demo/internal/decode"
)

func TestFormat(t *testing.T) {
	maskID, errorLevel, err := Format(sampleDots)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if maskID != 3 {
		t.Errorf("expected mask to equal 3 but got %d", maskID)
	}
	if errorLevel != decode.ErrorCorrectionLevelMedium {
		t.Errorf("expected errorLevel to equal %s (%d) but got %s (%d)", decode.ErrorCorrectionLevelMedium, decode.ErrorCorrectionLevelMedium, errorLevel, errorLevel)
	}
}
