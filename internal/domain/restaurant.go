package domain

import (
	"uuid"
)

type Restaurant struct {
	ID      string
	Name    string
	Address string
}

func NewRestaurant(name, address string) Restaurant {

	return Restaurant{
		ID:      uuid.New().String(),
		Name:    name,
		Address: address,
	}
}
