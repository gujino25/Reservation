package service

import (
	"context"
	"time"

	"reservation/internal/domain"
)

type tableGetter interface {
	GetByID(ctx context.Context, id int64) (domain.Table, error)
}

type guestGetter interface {
	GetByID(ctx context.Context, id int64) (domain.Guest, error)
}

type bookingStore interface {
	Create(ctx context.Context, b domain.Reservation) (domain.Reservation, error)
	GetByID(ctx context.Context, id int64) (domain.Reservation, error)
	Cancel(ctx context.Context, id int64) error
}

type BookingService struct {
	tables   tableGetter
	guests   guestGetter
	bookings bookingStore
}

func NewBookingService(tables tableGetter, guests guestGetter, bookings bookingStore) *BookingService {
	return &BookingService{
		tables:   tables,
		guests:   guests,
		bookings: bookings,
	}
}

func (s *BookingService) Create(ctx context.Context, tableID, guestID int64, startsAt time.Time, guestsCount int) (domain.Reservation, error) {
	res, err := domain.Create(tableID, guestID, startsAt, guestsCount)
	if err != nil {
		return domain.Reservation{}, err
	}

	table, err := s.tables.GetByID(ctx, tableID)
	if err != nil {
		return domain.Reservation{}, err
	}

	if _, err := s.guests.GetByID(ctx, guestID); err != nil {
		return domain.Reservation{}, err
	}

	if guestsCount > table.Capacity {
		return domain.Reservation{}, domain.ErrCapacityExceeded
	}

	return s.bookings.Create(ctx, res)
}

func (s *BookingService) Cancel(ctx context.Context, id int64) (domain.Reservation, error) {
	b, err := s.bookings.GetByID(ctx, id)
	if err != nil {
		return domain.Reservation{}, err
	}

	if err := b.Cancel(); err != nil {
		return domain.Reservation{}, err
	}

	if err := s.bookings.Cancel(ctx, id); err != nil {
		return domain.Reservation{}, err
	}
	return b, nil
}

func (s *BookingService) GetByID(ctx context.Context, id int64) (domain.Reservation, error) {
	return s.bookings.GetByID(ctx, id)
}
