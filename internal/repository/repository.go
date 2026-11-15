package repository

import (
	"context"

	"reservation/internal/domain"
)

type RestaurantRepository interface {
	Create(ctx context.Context, r domain.Restaurant) (domain.Restaurant, error)
	GetByID(ctx context.Context, id int64) (domain.Restaurant, error)
}

type TableRepository interface {
	Create(ctx context.Context, t domain.Table) (domain.Table, error)
	GetByID(ctx context.Context, id int64) (domain.Table, error)
	ListByRestaurant(ctx context.Context, restaurantID int64) ([]domain.Table, error)
}

type GuestRepository interface {
	Create(ctx context.Context, g domain.Guest) (domain.Guest, error)
	GetByID(ctx context.Context, id int64) (domain.Guest, error)
	GetByPhone(ctx context.Context, phone string) (domain.Guest, error)
}

type BookingRepository interface {
	Create(ctx context.Context, b domain.Reservation) (domain.Reservation, error)
	GetByID(ctx context.Context, id int64) (domain.Reservation, error)
	Update(ctx context.Context, b domain.Reservation) error
}
