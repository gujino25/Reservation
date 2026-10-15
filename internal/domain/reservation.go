package domain

import (
	"time"
	"uuid"
)

type Reservation struct {
	ID           string
	TableID      string
	GuestID      string
	GuestsAmount int
	CreatedAt    time.Time
	StartsAt     time.Time
	EndsAt       time.Time
	Reserved     bool
}

const (
	defaultReserveDuration = 2 * time.Hour
)

func NewReservation(tableID, guestID string, startsAt time.Time, guests int) (Reservation, error) {
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
		Reserved:     false,
	}, nil
}
