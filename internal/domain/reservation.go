package domain

import (
	"time"
)

type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "active"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusCompleted ReservationStatus = "completed"
)
const (
	defaultReserveDuration = 2 * time.Hour
)

type Reservation struct {
	ID           int64
	TableID      int64
	GuestID      int64
	GuestsAmount int
	CreatedAt    time.Time
	StartsAt     time.Time
	EndsAt       time.Time
	Status       ReservationStatus
}

func Create(tableID, guestID int64, startsAt time.Time, guests int) (Reservation, error) {
	if startsAt.Before(time.Now()) {
		return Reservation{}, ErrReservationTime
	}
	if guests <= 0 {
		return Reservation{}, ErrGuestAmountInvalid
	}
	endOfReservation := startsAt.Add(defaultReserveDuration)

	return Reservation{
		TableID:      tableID,
		GuestID:      guestID,
		GuestsAmount: guests,
		CreatedAt:    time.Now(),
		StartsAt:     startsAt,
		EndsAt:       endOfReservation,
		Status:       ReservationStatusActive,
	}, nil
}

func (r *Reservation) Cancel() error {
	if r.Status == ReservationStatusCancelled {
		return ErrReservationAlreadyCancelled
	}
	if !time.Now().Before(r.StartsAt) {
		return ErrReservationAlreadyStarted
	}
	r.Status = ReservationStatusCancelled
	return nil
}
