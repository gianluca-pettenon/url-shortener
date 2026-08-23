package hashid

import (
	"fmt"

	"github.com/gianluca-pettenon/url-shortener/internal/base62"
	"github.com/speps/go-hashids/v2"
)

func newHasher(salt string) (*hashids.HashID, error) {
	data := hashids.NewData()
	data.Salt = salt
	data.MinLength = 4

	return hashids.NewWithData(data)
}

func Encode(base62Value, salt string) (string, error) {
	hd, err := newHasher(salt)

	if err != nil {
		return "", fmt.Errorf("HashID: %w", err)
	}

	id, err := base62.Decode(base62Value)

	if err != nil {
		return "", fmt.Errorf("HashID encode: %w", err)
	}

	code, err := hd.Encode([]int{int(id)})

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

	if len(nums) != 1 || nums[0] < 0 {
		return "", fmt.Errorf("HashID decode: invalid value")
	}

	return base62.Encode(uint64(nums[0])), nil
}
