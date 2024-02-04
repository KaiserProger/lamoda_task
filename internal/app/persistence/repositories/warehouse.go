package repositories

import "lamoda_task/internal/app/models"

type WarehouseRepository interface {
	Get(warehouseId int) (*models.Warehouse, error)
}
