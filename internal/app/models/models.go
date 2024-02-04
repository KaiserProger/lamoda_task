package models

type Warehouse struct {
	Id         int
	Name       string
	Accessible bool
	Items      []*Item
}

type Item struct {
	Code     int
	Name     string
	Size     int
	Quantity int
	StoredAt []*Warehouse
}

type ReservationItem struct {
	ItemCode    int
	WarehouseId int
	Quantity    int
}

func (item *ReservationItem) AsIntArgs() [3]int {
	return [3]int{
		item.ItemCode, item.WarehouseId, item.Quantity,
	}
}

func (*ReservationItem) MultipleIntArgs(reservation []*ReservationItem) [][3]int {
	items := make([][3]int, len(reservation))

	for _, item := range reservation {
		items = append(items, item.AsIntArgs())
	}

	return items
}
