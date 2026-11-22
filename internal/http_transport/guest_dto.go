package httptransport

import (
	"time"

	"reservation/internal/domain"
)

type createGuestRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type guestResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

func toGuestResponse(g domain.Guest) guestResponse {
	return guestResponse{
		ID:        g.ID,
		Name:      g.Name,
		Phone:     g.Phone,
		CreatedAt: g.CreatedAt,
	}
}
