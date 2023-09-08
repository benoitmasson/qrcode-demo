package decode

// Inspired from https://www.thonky.com/qr-code-tutorial/mask-patterns

type MaskID uint8

// mask is a function which given 2 coordinates i and j, returns whether the dot
// at these coordinates should be switched (true) or kept as it is (false)
type mask func(i, j int) bool

// masks is the liste of all 8 masks
// See https://en.wikipedia.org/wiki/QR_code#/media/File:QR_Format_Information.svg
// for the repeated pattern
var masks = []mask{
	/* maskID == 0 */ func(i, j int) bool { return (i+j)%2 == 0 },
	/* maskID == 1 */ func(i, j int) bool { return i%2 == 0 },
	/* maskID == 2 */ func(i, j int) bool { return j%3 == 0 },
	/* maskID == 3 */ func(i, j int) bool { return (i+j)%3 == 0 },
	/* maskID == 4 */ func(i, j int) bool { return (i/2+j/3)%2 == 0 },
	/* maskID == 5 */ func(i, j int) bool { return (i*j)%2+(i*j)%3 == 0 },
	/* maskID == 6 */ func(i, j int) bool { return ((i*j)%3+i*j)%2 == 0 },
	/* maskID == 7 */ func(i, j int) bool { return ((i*j)%3+i+j)%2 == 0 },
}
