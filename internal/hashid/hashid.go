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
	data.Alphabet = base62.Alphabet
	data.Salt = salt
	data.MinLength = 0

	hd, err := hashids.NewWithData(data)

	if err != nil {
		return nil, fmt.Errorf("HashID: %w", err)
	}

	return &Coder{hd: hd}, nil
}

func (c *Coder) Encode(id uint64) (string, error) {
	if id > uint64(math.MaxInt) {
		return "", fmt.Errorf("HashID encode: value out of range")
	}

	code, err := c.hd.Encode([]int{int(id)})

	if err != nil {
		return "", fmt.Errorf("HashID encode: %w", err)
	}

	return code, nil
}
