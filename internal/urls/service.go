package urls

import (
	"context"
	"errors"
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
	code, err := hashid.Encode(encoded, s.salt)

	if err != nil {
		return "", err
	}

	return code, nil
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
