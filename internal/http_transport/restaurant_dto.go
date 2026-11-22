package httptransport

import (
	"time"

	"reservation/internal/domain"
)

type createRestaurantRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type restaurantResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

func toRestaurantResponse(r domain.Restaurant) restaurantResponse {
	return restaurantResponse{
		ID:        r.ID,
		Name:      r.Name,
		Address:   r.Address,
		CreatedAt: r.CreatedAt,
	}
}
