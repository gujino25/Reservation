package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"reservation/internal/domain"
	"reservation/internal/repository"
)

type GuestRepo struct {
	pool *pgxpool.Pool
}

var _ repository.GuestRepository = (*GuestRepo)(nil)

func NewGuestRepo(pool *pgxpool.Pool) *GuestRepo {
	return &GuestRepo{pool: pool}
}

func (r *GuestRepo) Create(ctx context.Context, g domain.Guest) (domain.Guest, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO guests (name, phone, created_at) VALUES ($1, NULLIF($2, ''), $3) RETURNING id`,
		g.Name, g.Phone, g.CreatedAt,
	).Scan(&g.ID)
	if err != nil {
		return domain.Guest{}, mapError(err)
	}
	return g, nil
}

func (r *GuestRepo) GetByID(ctx context.Context, id int64) (domain.Guest, error) {
	var g domain.Guest
	var phone *string
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, phone, created_at FROM guests WHERE id = $1`, id,
	).Scan(&g.ID, &g.Name, &phone, &g.CreatedAt)
	if err != nil {
		return domain.Guest{}, mapError(err)
	}
	if phone != nil {
		g.Phone = *phone
	}
	return g, nil
}

func (r *GuestRepo) GetByPhone(ctx context.Context, phone string) (domain.Guest, error) {
	var g domain.Guest
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, phone, created_at FROM guests WHERE phone = $1`, phone,
	).Scan(&g.ID, &g.Name, &g.Phone, &g.CreatedAt)
	if err != nil {
		return domain.Guest{}, mapError(err)
	}
	return g, nil
}
