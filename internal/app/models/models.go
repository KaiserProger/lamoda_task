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

func (item *ReservationItem) AsIntArgs() []int {
	return []int{
		item.ItemCode, item.WarehouseId, item.Quantity,
	}
}

func (*ReservationItem) MultipleIntArgs(reservation []*ReservationItem) [][]int {
	items := make([][]int, len(reservation))

	for _, item := range reservation {
		items = append(items, item.AsIntArgs())
	}

	return items
}
