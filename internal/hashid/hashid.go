package hashid

import (
	"fmt"

	"github.com/speps/go-hashids/v2"
)

func newHasher(salt string) (*hashids.HashID, error) {
	data := hashids.NewData()
	data.Salt = salt
	data.MinLength = 8

	return hashids.NewWithData(data)
}

func Encode(base62Value, salt string) (string, error) {
	hd, err := newHasher(salt)

	if err != nil {
		return "", fmt.Errorf("HashID: %w", err)
	}

	nums := make([]int, len(base62Value))

	for i := 0; i < len(base62Value); i++ {
		nums[i] = int(base62Value[i])
	}

	code, err := hd.Encode(nums)

	if err != nil {
		return "", fmt.Errorf("HashID encode: %w", err)
	}

	return code, nil
}

func Decode(code, salt string) (string, error) {
	hd, err := newHasher(salt)

	if err != nil {
		return "", fmt.Errorf("HashID: %w", err)
	}

	nums, err := hd.DecodeWithError(code)

	if err != nil {
		return "", fmt.Errorf("HashID decode: %w", err)
	}

	out := make([]byte, len(nums))

	for i, n := range nums {

		if n < 0 || n > 255 {
			return "", fmt.Errorf("HashID decode: invalid value")
		}

		out[i] = byte(n)
	}

	return string(out), nil
}
