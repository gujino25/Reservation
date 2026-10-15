package domain

import "uuid"

type Guest struct {
	ID   string
	Name string
}

func NewGuest(name string) Guest {
	return Guest{
		ID:   uuid.New().String(),
		Name: name,
	}
}
