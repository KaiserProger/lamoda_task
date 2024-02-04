package models

type Warehouse struct {
	Id            int
	Name          string
	Accessibility int
	Items         []*StoredItem
}

type Item struct {
	Code     int
	Name     string
	Size     int
	Quantity int
	StoredAt []*Warehouse
}

type StoredItem struct {
	WarehouseItemId [2]int
	Warehouse       *Warehouse
	Item            *Item
}
