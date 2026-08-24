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

const (
	MaxCreateMany    = 10_000_000
	DefaultListLimit = 50
	MaxListLimit     = 1_000
)

var (
	ErrInvalidURL  = errors.New("Invalid URL")
	ErrInvalidCode = errors.New("Invalid Short Code")
)

type Service struct {
	ids   *idgen.Generator
	repo  *Repository
	coder *hashid.Coder
}

func NewService(ids *idgen.Generator, repo *Repository, salt string) (*Service, error) {
	coder, err := hashid.New(salt)

	if err != nil {
		return nil, err
	}

	return &Service{ids: ids, repo: repo, coder: coder}, nil
}

func (s *Service) Create(ctx context.Context, rawURL string) (string, error) {
	original, err := normalizeURL(rawURL)

	if err != nil {
		return "", err
	}

	n, err := s.ids.Next(ctx)

	if err != nil {
		return "", err
	}

	code, err := s.encode(n)

	if err != nil {
		return "", err
	}

	if err := s.repo.Insert(ctx, URL{ID: code, OriginalURL: original}); err != nil {
		return "", err
	}

	return code, nil
}

func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	if _, err := s.coder.Decode(code); err != nil {
		return "", ErrInvalidCode
	}

	u, err := s.repo.GetByID(ctx, code)

	if err != nil {
		return "", err
	}

	return u.OriginalURL, nil
}

func (s *Service) CreateMany(ctx context.Context, rawURL string, n int) (string, string, error) {
	if n < 1 || n > MaxCreateMany {
		return "", "", fmt.Errorf("Times must be an integer from 1 to %d", MaxCreateMany)
	}

	original, err := normalizeURL(rawURL)

	if err != nil {
		return "", "", err
	}

	first, last, err := s.ids.NextN(ctx, n)

	if err != nil {
		return "", "", err
	}

	var firstCode, lastCode string
	id := first

	err = s.repo.InsertRange(ctx, original, func() (string, bool, error) {
		if id > last {
			return "", false, nil
		}

		code, err := s.encode(id)

		if err != nil {
			return "", false, err
		}

		if id == first {
			firstCode = code
		}

		if id == last {
			lastCode = code
		}

		id++

		return code, true, nil
	})

	if err != nil {
		return "", "", err
	}

	return firstCode, lastCode, nil
}

func (s *Service) List(ctx context.Context, limit int) ([]URL, int64, error) {
	if limit < 1 || limit > MaxListLimit {
		return nil, 0, fmt.Errorf("Limit must be an integer from 1 to %d", MaxListLimit)
	}

	total, err := s.repo.Count(ctx)

	if err != nil {
		return nil, 0, err
	}

	items, err := s.repo.List(ctx, limit)

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *Service) encode(n uint64) (string, error) {
	return s.coder.Encode(base62.Encode(n))
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
