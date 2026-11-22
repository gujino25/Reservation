package httptransport

import (
	"reservation/internal/domain"
)

type createTableRequest struct {
	Number   int `json:"number"`
	Capacity int `json:"capacity"`
}

type tableResponse struct {
	ID           int64 `json:"id"`
	RestaurantID int64 `json:"restaurant_id"`
	Number       int   `json:"number"`
	Capacity     int   `json:"capacity"`
}

func toTableResponse(t domain.Table) tableResponse {
	return tableResponse{
		ID:           t.ID,
		RestaurantID: t.RestaurantID,
		Number:       t.NumberOfTable,
		Capacity:     t.Capacity,
	}
}
