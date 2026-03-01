package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClickRepository struct {
	pool *pgxpool.Pool
}

func NewClickRepository(pool *pgxpool.Pool) *ClickRepository {
	return &ClickRepository{pool: pool}
}

func (r *ClickRepository) Insert(ctx context.Context, shortCode, userAgent, referer string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO click_events (short_code, user_agent, referrer) VALUES ($1, $2, $3)`,
		shortCode, userAgent, referer,
	)
	return err
}

func (r *ClickRepository) GetClickCount(ctx context.Context, shortCode string) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM click_events WHERE short_code = $1`,
		shortCode,
	).Scan(&count)
	return count, err
}
