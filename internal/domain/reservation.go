package domain

import (
	"time"
	"uuid"
)

type ReservationStatus string

const (
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)
const (
	defaultReserveDuration = 2 * time.Hour
)

type Reservation struct {
	ID           string
	TableID      string
	GuestID      string
	GuestsAmount int
	CreatedAt    time.Time
	StartsAt     time.Time
	EndsAt       time.Time
	Status       ReservationStatus
}

func Create(tableID, guestID string, startsAt time.Time, guests int) (Reservation, error) {
	if startsAt.Before(time.Now()) {
		return Reservation{}, ErrReservationTime
	}
	if guests <= 0 {
		return Reservation{}, ErrGuestAmountInvalid
	}
	endOfReservation := startsAt.Add(defaultReserveDuration)

	return Reservation{
		ID:           uuid.New().String(),
		TableID:      tableID,
		GuestID:      guestID,
		GuestsAmount: guests,
		CreatedAt:    time.Now(),
		StartsAt:     startsAt,
		EndsAt:       endOfReservation,
		Status:       ReservationStatusConfirmed,
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
