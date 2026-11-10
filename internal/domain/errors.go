package domain

import "errors"

var (
	ErrReservationTime             = errors.New("reservation must be in the future")
	ErrCapacityInvalid             = errors.New("capacity must be positive")
	ErrGuestAmountInvalid          = errors.New("guest amount is invalid")
	ErrNumberOfTableInvalid        = errors.New("invalid number of table")
	ErrReservationAlreadyCancelled = errors.New("reservation is already cancelled")
	ErrReservationAlreadyStarted   = errors.New("reservation has already started")
	ErrNotFound                    = errors.New("not found")
)
