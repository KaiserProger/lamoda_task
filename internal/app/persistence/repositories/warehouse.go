package repositories

import (
	"context"
	"lamoda_task/internal/app/models"
)

type WarehouseRepository interface {
	Get(ctx context.Context, warehouseId int) (*models.Warehouse, error)
}
