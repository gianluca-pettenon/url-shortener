package hashid

import (
	"fmt"
	"math"

	"github.com/gianluca-pettenon/url-shortener/internal/base62"
	"github.com/speps/go-hashids/v2"
)

type Coder struct {
	hd *hashids.HashID
}

func New(salt string) (*Coder, error) {
	data := hashids.NewData()
	data.Salt = salt
	data.MinLength = 4

	hd, err := hashids.NewWithData(data)

	if err != nil {
		return nil, fmt.Errorf("HashID: %w", err)
	}

	return &Coder{hd: hd}, nil
}

func (c *Coder) Encode(base62Value string) (string, error) {
	id, err := base62.Decode(base62Value)

	if err != nil {
		return "", fmt.Errorf("HashID encode: %w", err)
	}

	if id > uint64(math.MaxInt) {
		return "", fmt.Errorf("HashID encode: value out of range")
	}

	code, err := c.hd.Encode([]int{int(id)})

	if err != nil {
		return "", fmt.Errorf("HashID encode: %w", err)
	}

	return code, nil
}

func (c *Coder) Decode(code string) (string, error) {
	nums, err := c.hd.DecodeWithError(code)

	if err != nil {
		return "", fmt.Errorf("HashID decode: %w", err)
	}

	if len(nums) != 1 || nums[0] < 0 {
		return "", fmt.Errorf("HashID decode: invalid value")
	}

	return base62.Encode(uint64(nums[0])), nil
}
