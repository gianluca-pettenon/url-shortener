package urls

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("URL not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Insert(ctx context.Context, u URL) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO urls (id, original_url) VALUES ($1, $2)`, u.ID, u.OriginalURL)

	if err != nil {
		return fmt.Errorf("Insert URL: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (URL, error) {
	var u URL
	err := r.pool.QueryRow(ctx, `SELECT id, original_url, created_at FROM urls WHERE id = $1`, id).Scan(&u.ID, &u.OriginalURL, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return URL{}, ErrNotFound
	}

	if err != nil {
		return URL{}, fmt.Errorf("Get URL: %w", err)
	}

	return u, nil
}

func (r *Repository) List(ctx context.Context) ([]URL, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, original_url, created_at FROM urls ORDER BY created_at DESC`)

	if err != nil {
		return nil, fmt.Errorf("List URLs: %w", err)
	}

	defer rows.Close()

	var list []URL
	for rows.Next() {
		var u URL

		if err := rows.Scan(&u.ID, &u.OriginalURL, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("List URLs: %w", err)
		}

		list = append(list, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("List URLs: %w", err)
	}

	return list, nil
}
