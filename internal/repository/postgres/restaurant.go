package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"reservation/internal/domain"
	"reservation/internal/repository"
)

type RestaurantRepo struct {
	pool *pgxpool.Pool
}

var _ repository.RestaurantRepository = (*RestaurantRepo)(nil)

func NewRestaurantRepo(pool *pgxpool.Pool) *RestaurantRepo {
	return &RestaurantRepo{pool: pool}
}

func (r *RestaurantRepo) Create(ctx context.Context, rest domain.Restaurant) (domain.Restaurant, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO restaurants (name, address, created_at) VALUES ($1, $2, $3) RETURNING id`,
		rest.Name, rest.Address, rest.CreatedAt,
	).Scan(&rest.ID)
	if err != nil {
		return domain.Restaurant{}, mapError(err)
	}
	return rest, nil
}

func (r *RestaurantRepo) GetByID(ctx context.Context, id int64) (domain.Restaurant, error) {
	var rest domain.Restaurant
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, address, created_at FROM restaurants WHERE id = $1`, id,
	).Scan(&rest.ID, &rest.Name, &rest.Address, &rest.CreatedAt)
	if err != nil {
		return domain.Restaurant{}, mapError(err)
	}
	return rest, nil
}
