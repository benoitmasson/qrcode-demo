package extract

import (
	"fmt"
	"log/slog"
	"math/bits"
)

func init() {
	initFormatRemainders()
}

// Inspired from https://connect.ed-diamond.com/GNU-Linux-Magazine/glmf-198/reparer-un-code-qr

const formatGenerator uint16 = 0b10100110111 // 1335

var formatRemainders = make([]uint16, 32)

func initFormatRemainders() {
	for format := uint16(0); format <= 0b11111; format++ {
		formatRemainders[format] = computeFormatRemainder(format)
	}
}

// computeFormatRemainder is an implementation of the Reed-Solomon algorithm, in the particular
// case of length-5 strings of 0's and 1's (represented by n), using the formatGenerator polynom.
func computeFormatRemainder(n uint16) uint16 {
	mod := uint(formatGenerator) << 4 // pad generator with trailing 0's to have 15 bits => 0b101001101110000
	val := uint(n) << 10              // pad format with 10 0's to have 15 bits
	for i := 0; i < 5; i++ {
		slog.Debug(fmt.Sprintf("Value: %015b | Modulo: %015b", val, mod))
		if bits.Len(uint(val)) >= bits.Len(mod) {
			val ^= mod
		}
		mod >>= 1
	}
	return uint16(val)
}

// decodeFormat returns the 5-bits formats which fit the better with the given 10-bits ECC code.
func decodeFormat(eccFormat uint16) []uint16 {
	formatsByHammingDistance := make([][]uint16, 11)
	for format := uint16(0); format <= 0b11111; format++ {
		r := formatRemainders[format]
		hammingDistance := uint16(bits.OnesCount16(eccFormat ^ r)) // number of differences between the given eccFormat and the candidate r
		formatsByHammingDistance[hammingDistance] = append(formatsByHammingDistance[hammingDistance], format)
	}

	for _, formats := range formatsByHammingDistance {
		if len(formats) == 0 {
			continue
		}
		return formats
	}
	return nil
}
