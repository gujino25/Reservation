package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"reservation/internal/domain"
	"reservation/internal/repository"
)

type TableRepo struct {
	pool *pgxpool.Pool
}

var _ repository.TableRepository = (*TableRepo)(nil)

func NewTableRepo(pool *pgxpool.Pool) *TableRepo {
	return &TableRepo{pool: pool}
}

func (r *TableRepo) Create(ctx context.Context, t domain.Table) (domain.Table, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO tables (restaurant_id, number, capacity) VALUES ($1, $2, $3) RETURNING id`,
		t.RestaurantID, t.NumberOfTable, t.Capacity,
	).Scan(&t.ID)
	if err != nil {
		return domain.Table{}, mapError(err)
	}
	return t, nil
}

func (r *TableRepo) GetByID(ctx context.Context, id int64) (domain.Table, error) {
	var t domain.Table
	err := r.pool.QueryRow(ctx,
		`SELECT id, restaurant_id, number, capacity FROM tables WHERE id = $1`, id,
	).Scan(&t.ID, &t.RestaurantID, &t.NumberOfTable, &t.Capacity)
	if err != nil {
		return domain.Table{}, mapError(err)
	}
	return t, nil
}

func (r *TableRepo) ListByRestaurant(ctx context.Context, restaurantID int64) ([]domain.Table, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, restaurant_id, number, capacity FROM tables WHERE restaurant_id = $1 ORDER BY number`,
		restaurantID,
	)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	tables := make([]domain.Table, 0)
	for rows.Next() {
		var t domain.Table
		if err := rows.Scan(&t.ID, &t.RestaurantID, &t.NumberOfTable, &t.Capacity); err != nil {
			return nil, mapError(err)
		}
		tables = append(tables, t)
	}
	return tables, mapError(rows.Err())
}
