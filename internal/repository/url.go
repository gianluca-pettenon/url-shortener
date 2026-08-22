package repository

import (
	"context"
	"errors"

	"github.com/gianluca-pettenon/url-shortener/internal/domain"
)

var ErrNotFound = errors.New("URL not found")

type URLRepository interface {
	Insert(ctx context.Context, u domain.URL) error
	GetByID(ctx context.Context, id uint64) (domain.URL, error)
}
