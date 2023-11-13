package decode

import (
	"errors"
	"fmt"
)

type Mode uint8

const (
	NumericMode      Mode = 0b0001
	AlphanumericMode Mode = 0b0010
	ByteMode         Mode = 0b0100
	KanjiMode        Mode = 0b1000
)

// GetMode extracts the contents type from the header.
// The first 4 bits are used.
// In case ECI character is found, ECI data is ignored, and the trimmed bits are returned.
// See https://www.thonky.com/qr-code-tutorial/data-encoding#step-3-add-the-mode-indicator
func GetMode(bits []bool) (Mode, []bool) {
	mode := BitsToUint16(bits[:4])
	if mode == 0b0111 {
		// ECI escape character found, it can be used for unicode encoding.
		// Skip character (no unicode support yet) and next byte, then start over reading mode.
		bits = bits[12:]
		mode = BitsToUint16(bits[:4])
	}
	return Mode(mode), bits
}

// GetContentLength extracts the contents length from the header.
// After removing the first 4 bits (mode), the number of bits used to encode the length
// is given by the QR-code version and mode.
// Also, the contents bits are returned after trimming the header.
// See https://www.thonky.com/qr-code-tutorial/data-encoding#step-4-add-the-character-count-indicator
func GetContentLength(bits []bool, version uint, mode Mode, errorCorrectionLevel ErrorCorrectionLevel) (uint, []bool, error) {
	nb := lengthBytes(version, mode)
	if nb == 0 {
		return 0, nil, fmt.Errorf("invalid version-mode (%d, %b) pair", version, mode)
	}
	if len(bits) < 4+nb {
		return 0, nil, errors.New("not enough bits in contents for metadata")
	}

	length := uint(BitsToUint16(bits[4 : 4+nb]))
	if length <= 0 || length > capacity(version, errorCorrectionLevel, mode) {
		return 0, nil, fmt.Errorf("invalid length %d", length)
	}

	return length, bits[4+nb:], nil
}

func lengthBytes(version uint, mode Mode) int {
	if version <= 9 {
		switch mode {
		case NumericMode:
			return 10
		case AlphanumericMode:
			return 9
		case ByteMode:
			return 8
		case KanjiMode:
			return 8
		default:
			return 0
		}
	}

	if version <= 26 {
		switch mode {
		case NumericMode:
			return 12
		case AlphanumericMode:
			return 11
		case ByteMode:
			return 16
		case KanjiMode:
			return 10
		default:
			return 0
		}
	}

	if version <= 40 {
		switch mode {
		case NumericMode:
			return 14
		case AlphanumericMode:
			return 13
		case ByteMode:
			return 16
		case KanjiMode:
			return 12
		default:
			return 0
		}
	}

	return 0
}

func capacity(version uint, errorCorrectionLevel ErrorCorrectionLevel, mode Mode) uint {
	if capacityByErrorCorrectionLevelByMode, ok := capacityByVersionByErrorCorrectionLevelByMode[version]; ok {
		if capacityByMode, ok := capacityByErrorCorrectionLevelByMode[errorCorrectionLevel]; ok {
			if capacity, ok := capacityByMode[mode]; ok {
				return capacity
			}
		}
	}
	return 0
}

// source: https://www.thonky.com/qr-code-tutorial/character-capacities
var capacityByVersionByErrorCorrectionLevelByMode = map[uint]map[ErrorCorrectionLevel]map[Mode]uint{
	1: {
		ErrorCorrectionLevelLow:      {NumericMode: 41, AlphanumericMode: 25, ByteMode: 17, KanjiMode: 10},
		ErrorCorrectionLevelMedium:   {NumericMode: 34, AlphanumericMode: 20, ByteMode: 14, KanjiMode: 8},
		ErrorCorrectionLevelQuartile: {NumericMode: 27, AlphanumericMode: 16, ByteMode: 11, KanjiMode: 7},
		ErrorCorrectionLevelHigh:     {NumericMode: 17, AlphanumericMode: 10, ByteMode: 7, KanjiMode: 4},
	},
	2: {
		ErrorCorrectionLevelLow:      {NumericMode: 77, AlphanumericMode: 47, ByteMode: 32, KanjiMode: 20},
		ErrorCorrectionLevelMedium:   {NumericMode: 63, AlphanumericMode: 38, ByteMode: 26, KanjiMode: 16},
		ErrorCorrectionLevelQuartile: {NumericMode: 48, AlphanumericMode: 29, ByteMode: 20, KanjiMode: 12},
		ErrorCorrectionLevelHigh:     {NumericMode: 34, AlphanumericMode: 20, ByteMode: 14, KanjiMode: 8},
	},
	3: {
		ErrorCorrectionLevelLow:      {NumericMode: 127, AlphanumericMode: 77, ByteMode: 53, KanjiMode: 32},
		ErrorCorrectionLevelMedium:   {NumericMode: 101, AlphanumericMode: 61, ByteMode: 42, KanjiMode: 26},
		ErrorCorrectionLevelQuartile: {NumericMode: 77, AlphanumericMode: 47, ByteMode: 32, KanjiMode: 20},
		ErrorCorrectionLevelHigh:     {NumericMode: 58, AlphanumericMode: 35, ByteMode: 24, KanjiMode: 15},
	},
	4: {
		ErrorCorrectionLevelLow:      {NumericMode: 187, AlphanumericMode: 114, ByteMode: 78, KanjiMode: 48},
		ErrorCorrectionLevelMedium:   {NumericMode: 149, AlphanumericMode: 90, ByteMode: 62, KanjiMode: 38},
		ErrorCorrectionLevelQuartile: {NumericMode: 111, AlphanumericMode: 67, ByteMode: 46, KanjiMode: 28},
		ErrorCorrectionLevelHigh:     {NumericMode: 82, AlphanumericMode: 50, ByteMode: 34, KanjiMode: 21},
	},
	5: {
		ErrorCorrectionLevelLow:      {NumericMode: 255, AlphanumericMode: 154, ByteMode: 106, KanjiMode: 65},
		ErrorCorrectionLevelMedium:   {NumericMode: 202, AlphanumericMode: 122, ByteMode: 84, KanjiMode: 52},
		ErrorCorrectionLevelQuartile: {NumericMode: 144, AlphanumericMode: 87, ByteMode: 60, KanjiMode: 37},
		ErrorCorrectionLevelHigh:     {NumericMode: 106, AlphanumericMode: 64, ByteMode: 44, KanjiMode: 27},
	},
	6: {
		ErrorCorrectionLevelLow:      {NumericMode: 322, AlphanumericMode: 195, ByteMode: 134, KanjiMode: 82},
		ErrorCorrectionLevelMedium:   {NumericMode: 255, AlphanumericMode: 154, ByteMode: 106, KanjiMode: 65},
		ErrorCorrectionLevelQuartile: {NumericMode: 178, AlphanumericMode: 108, ByteMode: 74, KanjiMode: 45},
		ErrorCorrectionLevelHigh:     {NumericMode: 139, AlphanumericMode: 84, ByteMode: 58, KanjiMode: 36},
	},
	7: {
		ErrorCorrectionLevelLow:      {NumericMode: 370, AlphanumericMode: 224, ByteMode: 154, KanjiMode: 95},
		ErrorCorrectionLevelMedium:   {NumericMode: 293, AlphanumericMode: 178, ByteMode: 122, KanjiMode: 75},
		ErrorCorrectionLevelQuartile: {NumericMode: 207, AlphanumericMode: 125, ByteMode: 86, KanjiMode: 53},
		ErrorCorrectionLevelHigh:     {NumericMode: 154, AlphanumericMode: 93, ByteMode: 64, KanjiMode: 39},
	},
	8: {
		ErrorCorrectionLevelLow:      {NumericMode: 461, AlphanumericMode: 279, ByteMode: 192, KanjiMode: 118},
		ErrorCorrectionLevelMedium:   {NumericMode: 365, AlphanumericMode: 221, ByteMode: 152, KanjiMode: 93},
		ErrorCorrectionLevelQuartile: {NumericMode: 259, AlphanumericMode: 157, ByteMode: 108, KanjiMode: 66},
		ErrorCorrectionLevelHigh:     {NumericMode: 202, AlphanumericMode: 122, ByteMode: 84, KanjiMode: 52},
	},
	9: {
		ErrorCorrectionLevelLow:      {NumericMode: 552, AlphanumericMode: 335, ByteMode: 230, KanjiMode: 141},
		ErrorCorrectionLevelMedium:   {NumericMode: 432, AlphanumericMode: 262, ByteMode: 180, KanjiMode: 111},
		ErrorCorrectionLevelQuartile: {NumericMode: 312, AlphanumericMode: 189, ByteMode: 130, KanjiMode: 80},
		ErrorCorrectionLevelHigh:     {NumericMode: 235, AlphanumericMode: 143, ByteMode: 98, KanjiMode: 60},
	},
	10: {
		ErrorCorrectionLevelLow:      {NumericMode: 652, AlphanumericMode: 395, ByteMode: 271, KanjiMode: 167},
		ErrorCorrectionLevelMedium:   {NumericMode: 513, AlphanumericMode: 311, ByteMode: 213, KanjiMode: 131},
		ErrorCorrectionLevelQuartile: {NumericMode: 364, AlphanumericMode: 221, ByteMode: 151, KanjiMode: 93},
		ErrorCorrectionLevelHigh:     {NumericMode: 288, AlphanumericMode: 174, ByteMode: 119, KanjiMode: 74},
	},
	11: {
		ErrorCorrectionLevelLow:      {NumericMode: 772, AlphanumericMode: 468, ByteMode: 321, KanjiMode: 198},
		ErrorCorrectionLevelMedium:   {NumericMode: 604, AlphanumericMode: 366, ByteMode: 251, KanjiMode: 155},
		ErrorCorrectionLevelQuartile: {NumericMode: 427, AlphanumericMode: 259, ByteMode: 177, KanjiMode: 109},
		ErrorCorrectionLevelHigh:     {NumericMode: 331, AlphanumericMode: 200, ByteMode: 137, KanjiMode: 85},
	},
	12: {
		ErrorCorrectionLevelLow:      {NumericMode: 883, AlphanumericMode: 535, ByteMode: 367, KanjiMode: 226},
		ErrorCorrectionLevelMedium:   {NumericMode: 691, AlphanumericMode: 419, ByteMode: 287, KanjiMode: 177},
		ErrorCorrectionLevelQuartile: {NumericMode: 489, AlphanumericMode: 296, ByteMode: 203, KanjiMode: 125},
		ErrorCorrectionLevelHigh:     {NumericMode: 374, AlphanumericMode: 227, ByteMode: 155, KanjiMode: 96},
	},
	13: {
		ErrorCorrectionLevelLow:      {NumericMode: 1022, AlphanumericMode: 619, ByteMode: 425, KanjiMode: 262},
		ErrorCorrectionLevelMedium:   {NumericMode: 796, AlphanumericMode: 483, ByteMode: 331, KanjiMode: 204},
		ErrorCorrectionLevelQuartile: {NumericMode: 580, AlphanumericMode: 352, ByteMode: 241, KanjiMode: 149},
		ErrorCorrectionLevelHigh:     {NumericMode: 427, AlphanumericMode: 259, ByteMode: 177, KanjiMode: 109},
	},
	14: {
		ErrorCorrectionLevelLow:      {NumericMode: 1101, AlphanumericMode: 667, ByteMode: 458, KanjiMode: 282},
		ErrorCorrectionLevelMedium:   {NumericMode: 871, AlphanumericMode: 528, ByteMode: 362, KanjiMode: 223},
		ErrorCorrectionLevelQuartile: {NumericMode: 621, AlphanumericMode: 376, ByteMode: 258, KanjiMode: 159},
		ErrorCorrectionLevelHigh:     {NumericMode: 468, AlphanumericMode: 283, ByteMode: 194, KanjiMode: 120},
	},
	15: {
		ErrorCorrectionLevelLow:      {NumericMode: 1250, AlphanumericMode: 758, ByteMode: 520, KanjiMode: 320},
		ErrorCorrectionLevelMedium:   {NumericMode: 991, AlphanumericMode: 600, ByteMode: 412, KanjiMode: 254},
		ErrorCorrectionLevelQuartile: {NumericMode: 703, AlphanumericMode: 426, ByteMode: 292, KanjiMode: 180},
		ErrorCorrectionLevelHigh:     {NumericMode: 530, AlphanumericMode: 321, ByteMode: 220, KanjiMode: 136},
	},
	16: {
		ErrorCorrectionLevelLow:      {NumericMode: 1408, AlphanumericMode: 854, ByteMode: 586, KanjiMode: 361},
		ErrorCorrectionLevelMedium:   {NumericMode: 1082, AlphanumericMode: 656, ByteMode: 450, KanjiMode: 277},
		ErrorCorrectionLevelQuartile: {NumericMode: 775, AlphanumericMode: 470, ByteMode: 322, KanjiMode: 198},
		ErrorCorrectionLevelHigh:     {NumericMode: 602, AlphanumericMode: 365, ByteMode: 250, KanjiMode: 154},
	},
	17: {
		ErrorCorrectionLevelLow:      {NumericMode: 1548, AlphanumericMode: 938, ByteMode: 644, KanjiMode: 397},
		ErrorCorrectionLevelMedium:   {NumericMode: 1212, AlphanumericMode: 734, ByteMode: 504, KanjiMode: 310},
		ErrorCorrectionLevelQuartile: {NumericMode: 876, AlphanumericMode: 531, ByteMode: 364, KanjiMode: 224},
		ErrorCorrectionLevelHigh:     {NumericMode: 674, AlphanumericMode: 408, ByteMode: 280, KanjiMode: 173},
	},
	18: {
		ErrorCorrectionLevelLow:      {NumericMode: 1725, AlphanumericMode: 1046, ByteMode: 718, KanjiMode: 442},
		ErrorCorrectionLevelMedium:   {NumericMode: 1346, AlphanumericMode: 816, ByteMode: 560, KanjiMode: 345},
		ErrorCorrectionLevelQuartile: {NumericMode: 948, AlphanumericMode: 574, ByteMode: 394, KanjiMode: 243},
		ErrorCorrectionLevelHigh:     {NumericMode: 746, AlphanumericMode: 452, ByteMode: 310, KanjiMode: 191},
	},
	19: {
		ErrorCorrectionLevelLow:      {NumericMode: 1903, AlphanumericMode: 1153, ByteMode: 792, KanjiMode: 488},
		ErrorCorrectionLevelMedium:   {NumericMode: 1500, AlphanumericMode: 909, ByteMode: 624, KanjiMode: 384},
		ErrorCorrectionLevelQuartile: {NumericMode: 1063, AlphanumericMode: 644, ByteMode: 442, KanjiMode: 272},
		ErrorCorrectionLevelHigh:     {NumericMode: 813, AlphanumericMode: 493, ByteMode: 338, KanjiMode: 208},
	},
	20: {
		ErrorCorrectionLevelLow:      {NumericMode: 2061, AlphanumericMode: 1249, ByteMode: 858, KanjiMode: 528},
		ErrorCorrectionLevelMedium:   {NumericMode: 1600, AlphanumericMode: 970, ByteMode: 666, KanjiMode: 410},
		ErrorCorrectionLevelQuartile: {NumericMode: 1159, AlphanumericMode: 702, ByteMode: 482, KanjiMode: 297},
		ErrorCorrectionLevelHigh:     {NumericMode: 919, AlphanumericMode: 557, ByteMode: 382, KanjiMode: 235},
	},
	21: {
		ErrorCorrectionLevelLow:      {NumericMode: 2232, AlphanumericMode: 1352, ByteMode: 929, KanjiMode: 572},
		ErrorCorrectionLevelMedium:   {NumericMode: 1708, AlphanumericMode: 1035, ByteMode: 711, KanjiMode: 438},
		ErrorCorrectionLevelQuartile: {NumericMode: 1224, AlphanumericMode: 742, ByteMode: 509, KanjiMode: 314},
		ErrorCorrectionLevelHigh:     {NumericMode: 969, AlphanumericMode: 587, ByteMode: 403, KanjiMode: 248},
	},
	22: {
		ErrorCorrectionLevelLow:      {NumericMode: 2409, AlphanumericMode: 1460, ByteMode: 1003, KanjiMode: 618},
		ErrorCorrectionLevelMedium:   {NumericMode: 1872, AlphanumericMode: 1134, ByteMode: 779, KanjiMode: 480},
		ErrorCorrectionLevelQuartile: {NumericMode: 1358, AlphanumericMode: 823, ByteMode: 565, KanjiMode: 348},
		ErrorCorrectionLevelHigh:     {NumericMode: 1056, AlphanumericMode: 640, ByteMode: 439, KanjiMode: 270},
	},
	23: {
		ErrorCorrectionLevelLow:      {NumericMode: 2620, AlphanumericMode: 1588, ByteMode: 1091, KanjiMode: 672},
		ErrorCorrectionLevelMedium:   {NumericMode: 2059, AlphanumericMode: 1248, ByteMode: 857, KanjiMode: 528},
		ErrorCorrectionLevelQuartile: {NumericMode: 1468, AlphanumericMode: 890, ByteMode: 611, KanjiMode: 376},
		ErrorCorrectionLevelHigh:     {NumericMode: 1108, AlphanumericMode: 672, ByteMode: 461, KanjiMode: 284},
	},
	24: {
		ErrorCorrectionLevelLow:      {NumericMode: 2812, AlphanumericMode: 1704, ByteMode: 1171, KanjiMode: 721},
		ErrorCorrectionLevelMedium:   {NumericMode: 2188, AlphanumericMode: 1326, ByteMode: 911, KanjiMode: 561},
		ErrorCorrectionLevelQuartile: {NumericMode: 1588, AlphanumericMode: 963, ByteMode: 661, KanjiMode: 407},
		ErrorCorrectionLevelHigh:     {NumericMode: 1228, AlphanumericMode: 744, ByteMode: 511, KanjiMode: 315},
	},
	25: {
		ErrorCorrectionLevelLow:      {NumericMode: 3057, AlphanumericMode: 1853, ByteMode: 1273, KanjiMode: 784},
		ErrorCorrectionLevelMedium:   {NumericMode: 2395, AlphanumericMode: 1451, ByteMode: 997, KanjiMode: 614},
		ErrorCorrectionLevelQuartile: {NumericMode: 1718, AlphanumericMode: 1041, ByteMode: 715, KanjiMode: 440},
		ErrorCorrectionLevelHigh:     {NumericMode: 1286, AlphanumericMode: 779, ByteMode: 535, KanjiMode: 330},
	},
	26: {
		ErrorCorrectionLevelLow:      {NumericMode: 3283, AlphanumericMode: 1990, ByteMode: 1367, KanjiMode: 842},
		ErrorCorrectionLevelMedium:   {NumericMode: 2544, AlphanumericMode: 1542, ByteMode: 1059, KanjiMode: 652},
		ErrorCorrectionLevelQuartile: {NumericMode: 1804, AlphanumericMode: 1094, ByteMode: 751, KanjiMode: 462},
		ErrorCorrectionLevelHigh:     {NumericMode: 1425, AlphanumericMode: 864, ByteMode: 593, KanjiMode: 365},
	},
	27: {
		ErrorCorrectionLevelLow:      {NumericMode: 3517, AlphanumericMode: 2132, ByteMode: 1465, KanjiMode: 902},
		ErrorCorrectionLevelMedium:   {NumericMode: 2701, AlphanumericMode: 1637, ByteMode: 1125, KanjiMode: 692},
		ErrorCorrectionLevelQuartile: {NumericMode: 1933, AlphanumericMode: 1172, ByteMode: 805, KanjiMode: 496},
		ErrorCorrectionLevelHigh:     {NumericMode: 1501, AlphanumericMode: 910, ByteMode: 625, KanjiMode: 385},
	},
	28: {
		ErrorCorrectionLevelLow:      {NumericMode: 3669, AlphanumericMode: 2223, ByteMode: 1528, KanjiMode: 940},
		ErrorCorrectionLevelMedium:   {NumericMode: 2857, AlphanumericMode: 1732, ByteMode: 1190, KanjiMode: 732},
		ErrorCorrectionLevelQuartile: {NumericMode: 2085, AlphanumericMode: 1263, ByteMode: 868, KanjiMode: 534},
		ErrorCorrectionLevelHigh:     {NumericMode: 1581, AlphanumericMode: 958, ByteMode: 658, KanjiMode: 405},
	},
	29: {
		ErrorCorrectionLevelLow:      {NumericMode: 3909, AlphanumericMode: 2369, ByteMode: 1628, KanjiMode: 1002},
		ErrorCorrectionLevelMedium:   {NumericMode: 3035, AlphanumericMode: 1839, ByteMode: 1264, KanjiMode: 778},
		ErrorCorrectionLevelQuartile: {NumericMode: 2181, AlphanumericMode: 1322, ByteMode: 908, KanjiMode: 559},
		ErrorCorrectionLevelHigh:     {NumericMode: 1677, AlphanumericMode: 1016, ByteMode: 698, KanjiMode: 430},
	},
	30: {
		ErrorCorrectionLevelLow:      {NumericMode: 4158, AlphanumericMode: 2520, ByteMode: 1732, KanjiMode: 1066},
		ErrorCorrectionLevelMedium:   {NumericMode: 3289, AlphanumericMode: 1994, ByteMode: 1370, KanjiMode: 843},
		ErrorCorrectionLevelQuartile: {NumericMode: 2358, AlphanumericMode: 1429, ByteMode: 982, KanjiMode: 604},
		ErrorCorrectionLevelHigh:     {NumericMode: 1782, AlphanumericMode: 1080, ByteMode: 742, KanjiMode: 457},
	},
	31: {
		ErrorCorrectionLevelLow:      {NumericMode: 4417, AlphanumericMode: 2677, ByteMode: 1840, KanjiMode: 1132},
		ErrorCorrectionLevelMedium:   {NumericMode: 3486, AlphanumericMode: 2113, ByteMode: 1452, KanjiMode: 894},
		ErrorCorrectionLevelQuartile: {NumericMode: 2473, AlphanumericMode: 1499, ByteMode: 1030, KanjiMode: 634},
		ErrorCorrectionLevelHigh:     {NumericMode: 1897, AlphanumericMode: 1150, ByteMode: 790, KanjiMode: 486},
	},
	32: {
		ErrorCorrectionLevelLow:      {NumericMode: 4686, AlphanumericMode: 2840, ByteMode: 1952, KanjiMode: 1201},
		ErrorCorrectionLevelMedium:   {NumericMode: 3693, AlphanumericMode: 2238, ByteMode: 1538, KanjiMode: 947},
		ErrorCorrectionLevelQuartile: {NumericMode: 2670, AlphanumericMode: 1618, ByteMode: 1112, KanjiMode: 684},
		ErrorCorrectionLevelHigh:     {NumericMode: 2022, AlphanumericMode: 1226, ByteMode: 842, KanjiMode: 518},
	},
	33: {
		ErrorCorrectionLevelLow:      {NumericMode: 4965, AlphanumericMode: 3009, ByteMode: 2068, KanjiMode: 1273},
		ErrorCorrectionLevelMedium:   {NumericMode: 3909, AlphanumericMode: 2369, ByteMode: 1628, KanjiMode: 1002},
		ErrorCorrectionLevelQuartile: {NumericMode: 2805, AlphanumericMode: 1700, ByteMode: 1168, KanjiMode: 719},
		ErrorCorrectionLevelHigh:     {NumericMode: 2157, AlphanumericMode: 1307, ByteMode: 898, KanjiMode: 553},
	},
	34: {
		ErrorCorrectionLevelLow:      {NumericMode: 5253, AlphanumericMode: 3183, ByteMode: 2188, KanjiMode: 1347},
		ErrorCorrectionLevelMedium:   {NumericMode: 4134, AlphanumericMode: 2506, ByteMode: 1722, KanjiMode: 1060},
		ErrorCorrectionLevelQuartile: {NumericMode: 2949, AlphanumericMode: 1787, ByteMode: 1228, KanjiMode: 756},
		ErrorCorrectionLevelHigh:     {NumericMode: 2301, AlphanumericMode: 1394, ByteMode: 958, KanjiMode: 590},
	},
	35: {
		ErrorCorrectionLevelLow:      {NumericMode: 5529, AlphanumericMode: 3351, ByteMode: 2303, KanjiMode: 1417},
		ErrorCorrectionLevelMedium:   {NumericMode: 4343, AlphanumericMode: 2632, ByteMode: 1809, KanjiMode: 1113},
		ErrorCorrectionLevelQuartile: {NumericMode: 3081, AlphanumericMode: 1867, ByteMode: 1283, KanjiMode: 790},
		ErrorCorrectionLevelHigh:     {NumericMode: 2361, AlphanumericMode: 1431, ByteMode: 983, KanjiMode: 605},
	},
	36: {
		ErrorCorrectionLevelLow:      {NumericMode: 5836, AlphanumericMode: 3537, ByteMode: 2431, KanjiMode: 1496},
		ErrorCorrectionLevelMedium:   {NumericMode: 4588, AlphanumericMode: 2780, ByteMode: 1911, KanjiMode: 1176},
		ErrorCorrectionLevelQuartile: {NumericMode: 3244, AlphanumericMode: 1966, ByteMode: 1351, KanjiMode: 832},
		ErrorCorrectionLevelHigh:     {NumericMode: 2524, AlphanumericMode: 1530, ByteMode: 1051, KanjiMode: 647},
	},
	37: {
		ErrorCorrectionLevelLow:      {NumericMode: 6153, AlphanumericMode: 3729, ByteMode: 2563, KanjiMode: 1577},
		ErrorCorrectionLevelMedium:   {NumericMode: 4775, AlphanumericMode: 2894, ByteMode: 1989, KanjiMode: 1224},
		ErrorCorrectionLevelQuartile: {NumericMode: 3417, AlphanumericMode: 2071, ByteMode: 1423, KanjiMode: 876},
		ErrorCorrectionLevelHigh:     {NumericMode: 2625, AlphanumericMode: 1591, ByteMode: 1093, KanjiMode: 673},
	},
	38: {
		ErrorCorrectionLevelLow:      {NumericMode: 6479, AlphanumericMode: 3927, ByteMode: 2699, KanjiMode: 1661},
		ErrorCorrectionLevelMedium:   {NumericMode: 5039, AlphanumericMode: 3054, ByteMode: 2099, KanjiMode: 1292},
		ErrorCorrectionLevelQuartile: {NumericMode: 3599, AlphanumericMode: 2181, ByteMode: 1499, KanjiMode: 923},
		ErrorCorrectionLevelHigh:     {NumericMode: 2735, AlphanumericMode: 1658, ByteMode: 1139, KanjiMode: 701},
	},
	39: {
		ErrorCorrectionLevelLow:      {NumericMode: 6743, AlphanumericMode: 4087, ByteMode: 2809, KanjiMode: 1729},
		ErrorCorrectionLevelMedium:   {NumericMode: 5313, AlphanumericMode: 3220, ByteMode: 2213, KanjiMode: 1362},
		ErrorCorrectionLevelQuartile: {NumericMode: 3791, AlphanumericMode: 2298, ByteMode: 1579, KanjiMode: 972},
		ErrorCorrectionLevelHigh:     {NumericMode: 2927, AlphanumericMode: 1774, ByteMode: 1219, KanjiMode: 750},
	},
	40: {
		ErrorCorrectionLevelLow:      {NumericMode: 7089, AlphanumericMode: 4296, ByteMode: 2953, KanjiMode: 1817},
		ErrorCorrectionLevelMedium:   {NumericMode: 5596, AlphanumericMode: 3391, ByteMode: 2331, KanjiMode: 1435},
		ErrorCorrectionLevelQuartile: {NumericMode: 3993, AlphanumericMode: 2420, ByteMode: 1663, KanjiMode: 1024},
		ErrorCorrectionLevelHigh:     {NumericMode: 3057, AlphanumericMode: 1852, ByteMode: 1273, KanjiMode: 784},
	},
}
