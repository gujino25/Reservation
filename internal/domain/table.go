package domain

type Table struct {
	ID            int64
	RestaurantID  int64
	NumberOfTable int
	Capacity      int
}

func NewTable(restaurantID int64, numberOfTable, capacity int) (Table, error) {
	if capacity <= 0 {
		return Table{}, ErrCapacityInvalid
	}
	if numberOfTable <= 0 {
		return Table{}, ErrNumberOfTableInvalid
	}
	return Table{
		RestaurantID:  restaurantID,
		NumberOfTable: numberOfTable,
		Capacity:      capacity,
	}, nil
}
