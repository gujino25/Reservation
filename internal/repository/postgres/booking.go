package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"reservation/internal/domain"
	"reservation/internal/repository"
)

type BookingRepo struct {
	pool *pgxpool.Pool
}

var _ repository.BookingRepository = (*BookingRepo)(nil)

func NewBookingRepo(pool *pgxpool.Pool) *BookingRepo {
	return &BookingRepo{pool: pool}
}

func (r *BookingRepo) Create(ctx context.Context, b domain.Reservation) (domain.Reservation, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO bookings (table_id, guest_id, starts_at, ends_at, guests_count, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		b.TableID, b.GuestID, b.StartsAt, b.EndsAt, b.GuestsAmount, b.Status, b.CreatedAt,
	).Scan(&b.ID)
	if err != nil {
		return domain.Reservation{}, mapError(err)
	}
	return b, nil
}

func (r *BookingRepo) GetByID(ctx context.Context, id int64) (domain.Reservation, error) {
	var b domain.Reservation
	err := r.pool.QueryRow(ctx,
		`SELECT id, table_id, guest_id, guests_count, created_at, starts_at, ends_at, status
		 FROM bookings WHERE id = $1`, id,
	).Scan(&b.ID, &b.TableID, &b.GuestID, &b.GuestsAmount, &b.CreatedAt, &b.StartsAt, &b.EndsAt, &b.Status)
	if err != nil {
		return domain.Reservation{}, mapError(err)
	}
	return b, nil
}

func (r *BookingRepo) Update(ctx context.Context, b domain.Reservation) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE bookings
		 SET table_id = $2, guest_id = $3, starts_at = $4, ends_at = $5,
		     guests_count = $6, status = $7, created_at = $8
		 WHERE id = $1`,
		b.ID, b.TableID, b.GuestID, b.StartsAt, b.EndsAt, b.GuestsAmount, b.Status, b.CreatedAt,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Cancel changes the status only while the booking is still active. The
// conditional update makes concurrent cancellation requests deterministic.
func (r *BookingRepo) Cancel(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE bookings SET status = $2 WHERE id = $1 AND status = $3 AND starts_at > now()`,
		id, domain.ReservationStatusCancelled, domain.ReservationStatusActive,
	)
	if err != nil {
		return mapError(err)
	}
	if tag.RowsAffected() == 0 {
		var status domain.ReservationStatus
		var started bool
		if err := r.pool.QueryRow(ctx, `SELECT status, starts_at <= now() FROM bookings WHERE id = $1`, id).Scan(&status, &started); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrNotFound
			}
			return err
		}
		if started {
			return domain.ErrReservationAlreadyStarted
		}
		return domain.ErrReservationAlreadyCancelled
	}
	return nil
}
