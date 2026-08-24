package base62

const Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(n uint64) string {
	if n == 0 {
		return string(Alphabet[0])
	}

	var buf [11]byte
	i := len(buf)

	for n > 0 {
		i--
		buf[i] = Alphabet[n%62]
		n /= 62
	}

	return string(buf[i:])
}
