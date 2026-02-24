package repository

import (
	"context"
	"errors"

	"github.com/gabrieldrouin/shortly/shortener-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateShortCode = errors.New("short code already exists")

type URLRepository struct {
	pool *pgxpool.Pool
}

func NewURLRepository(pool *pgxpool.Pool) *URLRepository {
	return &URLRepository{pool: pool}
}

func (r *URLRepository) Insert(ctx context.Context, shortCode, originalURL string) (*model.URL, error) {
	u := &model.URL{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO urls (short_code, original_url) VALUES ($1, $2)
		 RETURNING id, short_code, original_url, created_at, expires_at`,
		shortCode, originalURL,
	).Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt, &u.ExpiresAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateShortCode
		}
		return nil, err
	}
	return u, nil
}

// GetByShortCode retrieves a URL by its short code.
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
