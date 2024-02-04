package postgres

import (
	"context"
	"errors"
	"lamoda_task/internal/app/models"
	"lamoda_task/internal/app/persistence"
	app "lamoda_task/internal/app/persistence/repositories"
)

type _warehouseRepositoryImpl struct{}

func NewWarehouseRepository() app.WarehouseRepository {
	return &_warehouseRepositoryImpl{}
}

func (*_warehouseRepositoryImpl) Get(ctx context.Context, warehouseId int) (*models.Warehouse, error) {
	txCtx, ok := ctx.(persistence.TransactionalContext)
	if !ok {
		return nil, errors.New("context is not transactional")
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return nil, errors.Join(errors.New("get transaction fail"), err)
	}

	rows, err := tx.Query(getWarehouseBaseQuery, warehouseId)
	if err != nil {
		return nil, errors.Join(errors.New("query warehouse fail"), err)
	}
	defer rows.Close()

	var warehouse models.Warehouse

	for rows.Next() {
		if err := rows.Scan(
			&warehouse.Id,
			&warehouse.Name,
			&warehouse.Accessible); err != nil {
			return nil, errors.Join(errors.New("scan warehouse fail"), err)
		}
	}

	itemsRows, err := tx.Query(getWarehouseItemsQuery, warehouseId)
	if err != nil {
		return nil, errors.Join(errors.New("query warehouse items fail"), err)
	}
	defer itemsRows.Close()

	for itemsRows.Next() {
		var item models.Item

		if err := itemsRows.Scan(
			&item.Code,
			&item.Name,
			&item.Size,
			&item.Quantity); err != nil {
			return nil, errors.Join(errors.New("scan warehouse item fail"), err)
		}

		warehouse.Items = append(warehouse.Items, &item)
	}

	return &warehouse, nil
}
