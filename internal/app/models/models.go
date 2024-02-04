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
