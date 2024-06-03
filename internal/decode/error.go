package decode

import (
	"errors"
	"fmt"

	"github.com/colin-davis/reedSolomon"
)

type ErrorCorrectionLevel uint8

const (
	ErrorCorrectionLevelMedium ErrorCorrectionLevel = iota
	ErrorCorrectionLevelLow
	ErrorCorrectionLevelHigh
	ErrorCorrectionLevelQuartile
)

func (ecl ErrorCorrectionLevel) String() string {
	switch ecl {
	case ErrorCorrectionLevelMedium:
		return "Medium"
	case ErrorCorrectionLevelLow:
		return "Low"
	case ErrorCorrectionLevelHigh:
		return "High"
	case ErrorCorrectionLevelQuartile:
		return "Quartile"
	}
	return fmt.Sprintf("Unknown(%d)", ecl)
}

func init() {
	// see https://github.com/colin-davis/reedSolomon#initializing-look-up-tables
	if err := reedSolomon.InitGaloisFields(285, 0); err != nil {
		panic(fmt.Errorf("failed to initialize Galois fields for QR-codes: %w", err))
	}
}

// Correct applies the Reed-Solomon error correction algorithm to the given bits.
// The corrected bits are returned upon success, or an error if the correction failed.
//
// This function uses http://github.com/colin-davis/reedSolomon for error correction.
// See https://en.m.wikiversity.org/wiki/Reed%E2%80%93Solomon_codes_for_coders for further details
// on how the algorithm works.
func Correct(bits []bool, version uint, errorCorrectionLevel ErrorCorrectionLevel) ([]bool, error) {
	blocksLayout := dataLayoutByVersionByErrorCorrectionLevel[version][errorCorrectionLevel]
	if len(blocksLayout) != 1 || blocksLayout[0].numberOfBlocks != 1 {
		// TODO: de-interleave error correction blocks and contents,
		// in the following cases:
		// - error correction level == low 		and version >= 6
		// - error correction level == medium 	and version >= 4
		// - error correction level == quartile and version >= 3
		// - error correction level == high 	and version >= 3
		// See https://www.thonky.com/qr-code-tutorial/error-correction-table

		return bits, errors.New("unable to correct message: data is interleaved, not handled yet")
	}
	totalLength, contentLength := blocksLayout[0].totalBlockBytes, blocksLayout[0].contentBlockBytes

	contentInt := make([]int, 0, totalLength)
	for i := 0; i < totalLength*8; i += 8 {
		n := BitsToUint16(bits[i : i+8])
		contentInt = append(contentInt, int(n))

	}

	numberECCSymbols := totalLength - contentLength
	errorLocations := []int{}
	correctedContent, _, err := reedSolomon.Decode(contentInt, numberECCSymbols, errorLocations)
	if err != nil {
		return bits, fmt.Errorf("failed to correct message: %w", err)
	}
	return intSliceToBits(correctedContent), nil
}

// intSliceToBits converts a list of integers (each one representing a byte, hence between 0 and 255)
// to the corresponding sequence of bits.
// Note that the result length is necessarily the input length times 8 (each integer represents a byte = 8 bits).
// e.g.
// [1]      => [false, false, false, false, false, false, false, true]
// [1, 255] => [false, false, false, false, false, false, false, true, /* 1 */ true, true, true, true, true, true, true /* 255 */]
func intSliceToBits(s []int) []bool {
	bits := make([]bool, 0, len(s)*8)
	for _, v := range s {
		for j := 7; j >= 0; j-- {
			b := v&(1<<j) != 0
			bits = append(bits, b)
		}
	}
	return bits
}

type dataLayout struct {
	// numberOfBlocks counts in how many pieces data is split.
	numberOfBlocks int
	// totalBlockBytes is the size of a block (in bytes), including ECC symbols.
	// Note that the the full bits list may exceed numberOfBlocks * totalBlockBytes by a few bits
	// (no more than 7), all set to 0 at the end of the QR-code
	totalBlockBytes int
	// contentBlockBytes is the size of a content block (in bytes), without ECC symbols.
	// The number of ECC symbols in each block is then totalBlockBytes - contentBlockBytes.
	contentBlockBytes int
}

// dataLayoutByVersionByErrorCorrectionLevel contains the data layout for each version and error correction level.
// The value is a non-empty sequence of data layouts. The first layout describes the first blocks of data, and so on.
// Content and ECC data is interleaved when the total number of blocks for a given (version, error correction level)
// is greater than 2.
// Source: https://github.com/Epik75/GLMF199/blob/97fc4ee02455e30dbb3da73f16358b7092fb38ae/Reperes/Preparer/qrcodestandard.py#L91
// (see also https://www.thonky.com/qr-code-tutorial/error-correction-table)
var dataLayoutByVersionByErrorCorrectionLevel = [41][4][]dataLayout{
	{{}, {}, {}, {}}, // index 0 is not used
	{ // Version 1
		{{1, 26, 16}} /* M */, {{1, 26, 19}} /* L */, {{1, 26, 9}} /* H */, {{1, 26, 13}}, /* Q */
	},
	{ // Version 2
		{{1, 44, 28}} /* M */, {{1, 44, 34}} /* L */, {{1, 44, 16}} /* H */, {{1, 44, 22}}, /* Q */
	},
	{ // Version 3
		{{1, 70, 44}} /* M */, {{1, 70, 55}} /* L */, {{2, 35, 13}} /* H */, {{2, 35, 17}}, /* Q */
	},
	{ // Version 4
		{{2, 50, 32}} /* M */, {{1, 100, 80}} /* L */, {{4, 25, 9}} /* H */, {{2, 50, 24}}, /* Q */
	},
	{ // Version 5
		{{2, 67, 43}} /* M */, {{1, 134, 108}} /* L */, {{2, 33, 11}, {2, 34, 12}} /* H */, {{2, 33, 15}, {2, 34, 16}}, /* Q */
	},
	{ // Version 6
		{{4, 43, 27}} /* M */, {{2, 86, 68}} /* L */, {{4, 43, 15}} /* H */, {{4, 43, 19}}, /* Q */
	},
	{ // Version 7
		{{4, 49, 31}} /* M */, {{2, 98, 78}} /* L */, {{4, 39, 13}, {1, 40, 14}} /* H */, {{2, 32, 14}, {4, 33, 15}}, /* Q */
	},
	{ // Version 8
		{{2, 60, 38}, {2, 61, 39}} /* M */, {{2, 121, 97}} /* L */, {{4, 40, 14}, {2, 41, 15}} /* H */, {{4, 40, 18}, {2, 41, 19}}, /* Q */
	},
	{ // Version 9
		{{3, 58, 36}, {2, 59, 37}} /* M */, {{2, 146, 116}} /* L */, {{4, 36, 12}, {4, 37, 13}} /* H */, {{4, 36, 16}, {4, 37, 17}}, /* Q */
	},
	{ // Version 10
		{{4, 69, 43}, {1, 70, 44}} /* M */, {{2, 86, 68}, {2, 87, 69}} /* L */, {{6, 43, 15}, {2, 44, 16}} /* H */, {{6, 43, 19}, {2, 44, 20}}, /* Q */
	},
	{ // Version 11
		{{1, 80, 50}, {4, 81, 51}} /* M */, {{4, 101, 81}} /* L */, {{3, 36, 12}, {8, 37, 13}} /* H */, {{4, 50, 22}, {4, 51, 23}}, /* Q */
	},
	{ // Version 12
		{{6, 58, 36}, {2, 59, 37}} /* M */, {{2, 116, 92}, {2, 117, 93}} /* L */, {{7, 42, 14}, {4, 43, 15}} /* H */, {{4, 46, 20}, {6, 47, 21}}, /* Q */
	},
	{ // Version 13
		{{8, 59, 37}, {1, 60, 38}} /* M */, {{4, 133, 107}} /* L */, {{12, 33, 11}, {4, 34, 12}} /* H */, {{8, 44, 20}, {4, 45, 21}}, /* Q */
	},
	{ // Version 14
		{{4, 64, 40}, {5, 65, 41}} /* M */, {{3, 145, 115}, {1, 146, 116}} /* L */, {{11, 36, 12}, {5, 37, 13}} /* H */, {{11, 36, 16}, {5, 37, 17}}, /* Q */
	},
	{ // Version 15
		{{5, 65, 41}, {5, 66, 42}} /* M */, {{5, 109, 87}, {1, 110, 88}} /* L */, {{11, 36, 12}, {7, 37, 13}} /* H */, {{5, 54, 24}, {7, 55, 25}}, /* Q */
	},
	{ // Version 16
		{{7, 73, 45}, {3, 74, 46}} /* M */, {{5, 122, 98}, {1, 123, 99}} /* L */, {{3, 45, 15}, {13, 46, 16}} /* H */, {{15, 43, 19}, {2, 44, 20}}, /* Q */
	},
	{ // Version 17
		{{10, 74, 46}, {1, 75, 47}} /* M */, {{1, 135, 107}, {5, 136, 108}} /* L */, {{2, 42, 14}, {17, 43, 15}} /* H */, {{1, 50, 22}, {15, 51, 23}}, /* Q */
	},
	{ // Version 18
		{{9, 69, 43}, {4, 70, 44}} /* M */, {{5, 150, 120}, {1, 151, 121}} /* L */, {{2, 42, 14}, {19, 43, 15}} /* H */, {{17, 50, 22}, {1, 51, 23}}, /* Q */
	},
	{ // Version 19
		{{3, 70, 44}, {11, 71, 45}} /* M */, {{3, 141, 113}, {4, 142, 114}} /* L */, {{9, 39, 13}, {16, 40, 14}} /* H */, {{17, 47, 21}, {4, 48, 22}}, /* Q */
	},
	{ // Version 20
		{{3, 67, 41}, {13, 68, 42}} /* M */, {{3, 135, 107}, {5, 136, 108}} /* L */, {{15, 43, 15}, {10, 44, 16}} /* H */, {{15, 54, 24}, {5, 55, 25}}, /* Q */
	},
	{ // Version 21
		{{17, 68, 42}} /* M */, {{4, 144, 116}, {4, 145, 117}} /* L */, {{19, 46, 16}, {6, 47, 17}} /* H */, {{17, 50, 22}, {6, 51, 23}}, /* Q */
	},
	{ // Version 22
		{{17, 74, 46}} /* M */, {{2, 139, 111}, {7, 140, 112}} /* L */, {{34, 37, 13}} /* H */, {{7, 54, 24}, {16, 55, 25}}, /* Q */
	},
	{ // Version 23
		{{4, 75, 47}, {14, 76, 48}} /* M */, {{4, 151, 121}, {5, 152, 122}} /* L */, {{16, 45, 15}, {14, 46, 16}} /* H */, {{11, 54, 24}, {14, 55, 25}}, /* Q */
	},
	{ // Version 24
		{{6, 73, 45}, {14, 74, 46}} /* M */, {{6, 147, 117}, {4, 148, 118}} /* L */, {{30, 46, 16}, {2, 47, 17}} /* H */, {{11, 54, 24}, {16, 55, 25}}, /* Q */
	},
	{ // Version 25
		{{8, 75, 47}, {13, 76, 48}} /* M */, {{8, 132, 106}, {4, 133, 107}} /* L */, {{22, 45, 15}, {13, 46, 16}} /* H */, {{7, 54, 24}, {22, 55, 25}}, /* Q */
	},
	{ // Version 26
		{{19, 74, 46}, {4, 75, 47}} /* M */, {{10, 142, 114}, {2, 143, 115}} /* L */, {{33, 46, 16}, {4, 47, 17}} /* H */, {{28, 50, 22}, {6, 51, 23}}, /* Q */
	},
	{ // Version 27
		{{22, 73, 45}, {3, 74, 46}} /* M */, {{8, 152, 122}, {4, 153, 123}} /* L */, {{12, 45, 15}, {28, 46, 16}} /* H */, {{8, 53, 23}, {26, 54, 24}}, /* Q */
	},
	{ // Version 28
		{{3, 73, 45}, {23, 74, 46}} /* M */, {{3, 147, 117}, {10, 148, 118}} /* L */, {{11, 45, 15}, {31, 46, 16}} /* H */, {{4, 54, 24}, {31, 55, 25}}, /* Q */
	},
	{ // Version 29
		{{21, 73, 45}, {7, 74, 46}} /* M */, {{7, 146, 116}, {7, 147, 117}} /* L */, {{19, 45, 15}, {26, 46, 16}} /* H */, {{1, 53, 23}, {37, 54, 24}}, /* Q */
	},
	{ // Version 30
		{{19, 75, 47}, {10, 76, 48}} /* M */, {{5, 145, 115}, {10, 146, 116}} /* L */, {{23, 45, 15}, {25, 46, 16}} /* H */, {{15, 54, 24}, {25, 55, 25}}, /* Q */
	},
	{ // Version 31
		{{2, 74, 46}, {29, 75, 47}} /* M */, {{13, 145, 115}, {3, 146, 116}} /* L */, {{23, 45, 15}, {28, 46, 16}} /* H */, {{42, 54, 24}, {1, 55, 25}}, /* Q */
	},
	{ // Version 32
		{{10, 74, 46}, {23, 75, 47}} /* M */, {{17, 145, 115}} /* L */, {{19, 45, 15}, {35, 46, 16}} /* H */, {{10, 54, 24}, {35, 55, 25}}, /* Q */
	},
	{ // Version 33
		{{14, 74, 46}, {21, 75, 47}} /* M */, {{17, 145, 115}, {1, 146, 116}} /* L */, {{11, 45, 15}, {46, 46, 16}} /* H */, {{29, 54, 24}, {19, 55, 25}}, /* Q */
	},
	{ // Version 34
		{{14, 74, 46}, {23, 75, 47}} /* M */, {{13, 145, 115}, {6, 146, 116}} /* L */, {{59, 46, 16}, {1, 47, 17}} /* H */, {{44, 54, 24}, {7, 55, 25}}, /* Q */
	},
	{ // Version 35
		{{12, 75, 47}, {26, 76, 48}} /* M */, {{12, 151, 121}, {7, 152, 122}} /* L */, {{22, 45, 15}, {41, 46, 16}} /* H */, {{39, 54, 24}, {14, 55, 25}}, /* Q */
	},
	{ // Version 36
		{{6, 75, 47}, {34, 76, 48}} /* M */, {{6, 151, 121}, {14, 152, 122}} /* L */, {{2, 45, 15}, {64, 46, 16}} /* H */, {{46, 54, 24}, {10, 55, 25}}, /* Q */
	},
	{ // Version 37
		{{29, 74, 46}, {14, 75, 47}} /* M */, {{17, 152, 122}, {4, 153, 123}} /* L */, {{24, 45, 15}, {46, 46, 16}} /* H */, {{49, 54, 24}, {10, 55, 25}}, /* Q */
	},
	{ // Version 38
		{{13, 74, 46}, {32, 75, 47}} /* M */, {{4, 152, 122}, {18, 153, 123}} /* L */, {{42, 45, 15}, {32, 46, 16}} /* H */, {{48, 54, 24}, {14, 55, 25}}, /* Q */
	},
	{ // Version 39
		{{40, 75, 47}, {7, 76, 48}} /* M */, {{20, 147, 117}, {4, 148, 118}} /* L */, {{10, 45, 15}, {67, 46, 16}} /* H */, {{43, 54, 24}, {22, 55, 25}}, /* Q */
	},
	{ // Version 40
		{{18, 75, 47}, {31, 76, 48}} /* M */, {{19, 148, 118}, {6, 149, 119}} /* L */, {{20, 45, 15}, {61, 46, 16}} /* H */, {{34, 54, 24}, {34, 55, 25}}, /* Q */
	},
}
