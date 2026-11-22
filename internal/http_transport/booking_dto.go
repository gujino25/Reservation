package httptransport

import (
	"time"

	"reservation/internal/domain"
)

type createBookingRequest struct {
	TableID     int64     `json:"table_id"`
	GuestID     int64     `json:"guest_id"`
	StartsAt    time.Time `json:"starts_at"`
	GuestsCount int       `json:"guests_count"`
}

type bookingResponse struct {
	ID          int64     `json:"id"`
	TableID     int64     `json:"table_id"`
	GuestID     int64     `json:"guest_id"`
	GuestsCount int       `json:"guests_count"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func toBookingResponse(b domain.Reservation) bookingResponse {
	return bookingResponse{
		ID:          b.ID,
		TableID:     b.TableID,
		GuestID:     b.GuestID,
		GuestsCount: b.GuestsAmount,
		StartsAt:    b.StartsAt,
		EndsAt:      b.EndsAt,
		Status:      string(b.Status),
		CreatedAt:   b.CreatedAt,
	}
}
