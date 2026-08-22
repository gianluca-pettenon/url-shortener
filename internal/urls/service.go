package urls

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/gianluca-pettenon/url-shortener/internal/base62"
	"github.com/gianluca-pettenon/url-shortener/internal/idgen"
)

var (
	ErrInvalidURL  = errors.New("Invalid URL")
	ErrInvalidCode = errors.New("Invalid Short Code")
)

type Service struct {
	ids  *idgen.Generator
	repo *Repository
}

func NewService(ids *idgen.Generator, repo *Repository) *Service {
	return &Service{ids: ids, repo: repo}
}

func (s *Service) Create(ctx context.Context, rawURL string) (string, error) {
	original, err := normalizeURL(rawURL)

	if err != nil {
		return "", err
	}

	id, err := s.ids.Next(ctx)

	if err != nil {
		return "", err
	}

	if err := s.repo.Insert(ctx, URL{ID: id, OriginalURL: original}); err != nil {
		return "", err
	}

	return base62.Encode(id), nil
}

func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	id, err := base62.Decode(code)

	if err != nil {
		return "", ErrInvalidCode
	}

	u, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return "", err
	}

	return u.OriginalURL, nil
}

func normalizeURL(raw string) (string, error) {
	u, err := url.ParseRequestURI(strings.TrimSpace(raw))

	if err != nil {
		return "", ErrInvalidURL
	}

	if u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return "", ErrInvalidURL
	}

	return u.String(), nil
}
