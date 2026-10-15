package domain

import "uuid"

type Table struct {
	ID            string
	RestaurantID  string
	NumberOfTable int
	Capacity      int
}

func NewTable(restaurantID string, numberOfTable, capacity int) (Table, error) {
	if capacity <= 0 {
		return Table{}, ErrCapacityInvalid
	}
	if numberOfTable <= 0 {
		return Table{}, ErrNumberOfTable
	}
	return Table{
		ID:            uuid.New().String(),
		RestaurantID:  restaurantID,
		NumberOfTable: numberOfTable,
		Capacity:      capacity,
	}, nil
}
