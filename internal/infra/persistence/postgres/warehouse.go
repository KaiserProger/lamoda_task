package postgres

import (
	"context"
	"database/sql"
	"errors"
	appErrors "lamoda_task/internal/app/errors"
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
		return nil, appErrors.ErrNotTransactional
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(getWarehouseBaseQuery, warehouseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Join(errors.New("query warehouse fail"), err)
	}
	defer rows.Close()

	var warehouse *models.Warehouse

	for rows.Next() {
		warehouse = &models.Warehouse{}
		if err := rows.Scan(
			&warehouse.Id,
			&warehouse.Name,
			&warehouse.Accessible); err != nil {
			return nil, errors.Join(errors.New("scan warehouse fail"), err)
		}
	}

	if warehouse == nil {
		return nil, appErrors.ErrNotFound
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

	return warehouse, nil
}

func (*_warehouseRepositoryImpl) AddToStock(ctx context.Context, reservation []*models.ReservationItem) error {
	txCtx, ok := ctx.(persistence.TransactionalContext)
	if !ok {
		return appErrors.ErrNotTransactional
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return err
	}

	argsArray := (&models.ReservationItem{}).MultipleIntArgs(reservation)

	args := make([]any, len(argsArray))
	for ix, arg := range argsArray {
		args[ix] = arg
	}

	query, err := GenBulkPlaceholders(addToStockQuery, args, 3)
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Join(errors.New("update stock fail"), err)
	}

	return nil
}

func (*_warehouseRepositoryImpl) RemoveFromStock(ctx context.Context, reservation []*models.ReservationItem) error {
	txCtx, ok := ctx.(persistence.TransactionalContext)
	if !ok {
		return appErrors.ErrNotTransactional
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return err
	}

	argsArray := (&models.ReservationItem{}).MultipleIntArgs(reservation)

	args := make([]any, len(argsArray))
	for ix, arg := range argsArray {
		args[ix] = arg
	}

	query, err := GenBulkPlaceholders(removeFromStockQuery, args, 3)
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Join(errors.New("update stock fail"), err)
	}

	return nil
}
