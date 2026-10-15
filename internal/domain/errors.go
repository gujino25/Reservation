package domain

import "errors"

var (
	ErrReservationTime    = errors.New("reservation must be in the future")
	ErrCapacityInvalid    = errors.New("capacity must be positive")
	ErrGuestAmountInvalid = errors.New("guest amount is invalid")
	ErrNumberOfTable      = errors.New("invalid number of table")
)
