package urls

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gianluca-pettenon/url-shortener/internal/base62"
	"github.com/gianluca-pettenon/url-shortener/internal/hashid"
	"github.com/gianluca-pettenon/url-shortener/internal/idgen"
)

var (
	ErrInvalidURL  = errors.New("Invalid URL")
	ErrInvalidCode = errors.New("Invalid Short Code")
)

type Service struct {
	ids  *idgen.Generator
	repo *Repository
	salt string
}

func NewService(ids *idgen.Generator, repo *Repository, salt string) *Service {
	return &Service{ids: ids, repo: repo, salt: salt}
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

	encoded := base62.Encode(id)

	return hashid.Encode(encoded, s.salt)
}

func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	encoded, err := hashid.Decode(code, s.salt)

	if err != nil {
		return "", ErrInvalidCode
	}

	id, err := base62.Decode(encoded)

	if err != nil {
		return "", ErrInvalidCode
	}

	u, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return "", err
	}

	return u.OriginalURL, nil
}

func (s *Service) Code(id uint64) (string, error) {
	return hashid.Encode(base62.Encode(id), s.salt)
}

func (s *Service) CreateMany(ctx context.Context, rawURL string, n int) (uint64, uint64, error) {
	if n < 1 {
		return 0, 0, fmt.Errorf("times must be at least 1")
	}

	original, err := normalizeURL(rawURL)

	if err != nil {
		return 0, 0, err
	}

	first, last, err := s.ids.NextN(ctx, n)

	if err != nil {
		return 0, 0, err
	}

	if err := s.repo.InsertRange(ctx, first, last, original); err != nil {
		return 0, 0, err
	}

	return first, last, nil
}

func (s *Service) List(ctx context.Context) ([]URL, error) {
	return s.repo.List(ctx)
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
