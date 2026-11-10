package domain

import (
	"time"
)

type Restaurant struct {
	ID        int64
	Name      string
	Address   string
	CreatedAt time.Time
}

func NewRestaurant(name, address string) Restaurant {

	return Restaurant{
		Name:      name,
		Address:   address,
		CreatedAt: time.Now(),
	}
}
