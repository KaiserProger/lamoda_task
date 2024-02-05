package postgres

import (
	"context"
	"errors"
	appErrors "lamoda_task/internal/app/errors"
	"lamoda_task/internal/app/persistence"
	app "lamoda_task/internal/app/persistence/repositories"
)

type _itemRepositoryImpl struct {
}

func NewItemRepository() app.ItemRepository {
	return &_itemRepositoryImpl{}
}

func (*_itemRepositoryImpl) GetStoredAt(ctx context.Context, itemCodes []int) ([]*app.StoredItem, error) {
	storedItems := make([]*app.StoredItem, 0)
	txCtx, ok := ctx.(persistence.TransactionalContext)
	if !ok {
		return nil, appErrors.ErrNotTransactional
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(getStoredItemsQuery, itemCodes)
	if err != nil {
		return nil, errors.Join(errors.New("get stored items fail"), err)
	}
	defer rows.Close()

	for rows.Next() {
		storedItem := &app.StoredItem{}
		if err := rows.Scan(
			&storedItem.ItemCode,
			&storedItem.WarehouseId,
			&storedItem.Quantity); err != nil {
			return nil, errors.Join(errors.New("scan stored items fail"), err)
		}

		storedItems = append(storedItems, storedItem)
	}

	return storedItems, nil
}
