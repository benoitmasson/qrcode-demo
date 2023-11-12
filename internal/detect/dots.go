package detect

import (
	"fmt"
	"math"

	"gocv.io/x/gocv"
)

// QRCode is the representation of the code: 2-dimensional array of dots
// "true" means black dot, "false" means
type QRCode [][]bool

const (
	luminosityThreshold         = 500
	luminosityPersistenceOffset = 0 // used to favor sequence color persistence if in doubt
)

// GetDots scans the input image pixels to try to extract QR-code dots.
// Returns the grid of dots (true is black), and a boolean telling whether extraction was successful.
func GetDots(img gocv.Mat) (QRCode, bool) {
	firstColumn := getFirstColumnSequences(img)
	// fmt.Printf("First column: %v\n", firstColumn)

	firstColumnStartsAndEndsInBlack := len(firstColumn) >= 3 && len(firstColumn)%2 == 1
	firstAndLastBlocksHaveSimilarLength := len(firstColumn) >= 3 && nearlyEquals(firstColumn[0], firstColumn[len(firstColumn)-1], 5)
	if !firstColumnStartsAndEndsInBlack || !firstAndLastBlocksHaveSimilarLength {
		return nil, false
	}
	scale := float64(firstColumn[0]) / 7. // a marker is 7 dots high
	height := int(math.Round(float64(img.Rows()) / scale))
	if height%2 == 0 {
		return nil, false
	}

	fmt.Printf("Dots are %f pixels wide\n", scale)
	dots := scanDots(img, scale)
	return dots, true
}

// getFirstColumnSequences returns a sequence of luminosity values corresponding to pixels in the QR-code first column
func getFirstColumnSequences(img gocv.Mat) []int {
	sequences := make([]int, 0)
	const offset = 1

	isFirstSequence := true
	currentSequenceIsBlack := false
	currentSequenceLength := 0

	for row := offset; row < img.Rows()-offset; row++ {
		vec := img.GetVecbAt(row, offset)
		pixelLuminosity := int(vec[0]) + int(vec[1]) + int(vec[2])

		threshold := luminosityThreshold
		if currentSequenceIsBlack {
			threshold += luminosityPersistenceOffset
		} else {
			threshold -= luminosityPersistenceOffset
		}
		pixelIsBlack := true
		if pixelLuminosity >= threshold {
			pixelIsBlack = false
		}

		// fmt.Printf("[row %d] lum: %v / black: %v / sequence: %v\n", y, pixelLuminosity, pixelIsBlack, currentSequenceLength)
		if pixelIsBlack == currentSequenceIsBlack {
			currentSequenceLength++
		} else {
			if isFirstSequence && !currentSequenceIsBlack {
				isFirstSequence = false
			} else {
				sequences = append(sequences, currentSequenceLength)
			}
			currentSequenceIsBlack = !currentSequenceIsBlack
			currentSequenceLength = 1
		}
	}
	if !isFirstSequence {
		sequences = append(sequences, currentSequenceLength)
	}

	return sequences
}

// scanDots scans the input image step by step according to the given scale (dot size in pixel),
// and constructs the dots grid
func scanDots(img gocv.Mat, scale float64) QRCode {
	dots := make(QRCode, 0)

	i, j, row, col := 0, 0, 0, 0
	for {
		row = int((float64(i)+0.5)*scale + 0.5)
		if row >= img.Rows() {
			break
		}
		line := make([]bool, 0)
		for {
			col = int((float64(j)+0.5)*scale + 0.5)
			if col >= img.Cols() {
				break
			}
			vec := img.GetVecbAt(row, col)
			pixelLuminosity := int(vec[0]) + int(vec[1]) + int(vec[2])
			pixelIsBlack := true
			if pixelLuminosity >= luminosityThreshold {
				pixelIsBlack = false
			}
			line = append(line, pixelIsBlack)
			j++
		}
		dots = append(dots, line)
		i++
		j = 0
	}

	return dots
}
