package extract

import (
	"testing"
)

func TestVersion(t *testing.T) {
	version, err := Version(sampleDots)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if version != 2 {
		t.Errorf("expected version to equal 2 but got %d", version)
	}
}
