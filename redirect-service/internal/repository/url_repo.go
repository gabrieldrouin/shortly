package repository

import (
	"context"
	"errors"

	"github.com/gabrieldrouin/shortly/redirect-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	pool *pgxpool.Pool
}

func NewURLRepository(pool *pgxpool.Pool) *URLRepository {
	return &URLRepository{pool: pool}
}

func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	u := &model.URL{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, short_code, original_url, created_at, expires_at FROM urls WHERE short_code = $1`,
		shortCode,
	).Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt, &u.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}
