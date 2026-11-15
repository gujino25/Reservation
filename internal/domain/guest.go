package domain

import "time"

type Guest struct {
	ID        int64
	Name      string
	Phone     string
	CreatedAt time.Time
}

func NewGuest(name, phone string) Guest {
	return Guest{
		Name:      name,
		Phone:     phone,
		CreatedAt: time.Now(),
	}
}
