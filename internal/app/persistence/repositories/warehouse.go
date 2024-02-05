package repositories

import (
	"context"
	"lamoda_task/internal/app/models"
)

type WarehouseRepository interface {
	Get(ctx context.Context, warehouseId int) (*models.Warehouse, error)
	AddToStock(ctx context.Context, reservation []*models.ReservationItem) error
	RemoveFromStock(ctx context.Context, reservation []*models.ReservationItem) error
}
