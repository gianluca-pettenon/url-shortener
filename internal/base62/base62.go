package base62

import "fmt"

const (
	alphabet      = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base          = len(alphabet)
	maxEncodedLen = 11
	invalidChar   = -1
)

var decodeTable [256]int8

func init() {
	for i := range decodeTable {
		decodeTable[i] = invalidChar
	}

	for i := range base {
		decodeTable[alphabet[i]] = int8(i)
	}
}

func Encode(num uint64) string {
	if num == 0 {
		return string(alphabet[0])
	}

	var buf [maxEncodedLen]byte
	pos := maxEncodedLen

	for num > 0 {
		pos--
		buf[pos] = alphabet[num%uint64(base)]
		num /= uint64(base)
	}

	return string(buf[pos:])
}

func Decode(s string) (uint64, error) {
	if s == "" {
		return 0, fmt.Errorf("Empty String")
	}

	var num uint64

	for i := range len(s) {

		digit := decodeTable[s[i]]

		if digit == invalidChar {
			return 0, fmt.Errorf("Invalid character: %q", s[i])
		}

		num = num*uint64(base) + uint64(digit)
	}

	return num, nil
}
