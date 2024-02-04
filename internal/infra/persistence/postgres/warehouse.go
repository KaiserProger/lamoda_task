package postgres

import (
	"lamoda_task/internal/app/models"
	app "lamoda_task/internal/app/persistence/repositories"
)

type _warehouseRepositoryImpl struct{}

func NewWarehouseRepository() app.WarehouseRepository {
	return &_warehouseRepositoryImpl{}
}

func (*_warehouseRepositoryImpl) Get(warehouseId int) (*models.Warehouse, error) {
	panic("unimplemented")
}

func (*_warehouseRepositoryImpl) Items(warehouseId int) ([]*models.Item, error) {
	panic("unimplemented")
}

func (*_warehouseRepositoryImpl) Update(*models.Warehouse) error {
	panic("unimplemented")
}
