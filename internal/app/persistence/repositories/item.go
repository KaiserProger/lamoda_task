package repositories

import (
	"context"
)

type StoredItem struct {
	ItemCode    int
	WarehouseId int
	Quantity    int
}

func (s *StoredItem) AsQueryArgs() []any {
	return []any{s.ItemCode, s.WarehouseId, s.Quantity}
}

type ItemRepository interface {
	// Get items with warehouses they are stored at.
	GetStoredAt(ctx context.Context, itemCodes []int) ([]*StoredItem, error)
}
