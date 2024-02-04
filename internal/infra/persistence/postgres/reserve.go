package postgres

import (
	"context"
	"errors"
	appErrors "lamoda_task/internal/app/errors"
	"lamoda_task/internal/app/persistence"
	app "lamoda_task/internal/app/persistence/repositories"
)

type _reserveRepositoryImpl struct{}

func NewReserveRepository() app.ReserveRepository {
	return &_reserveRepositoryImpl{}
}

func (*_reserveRepositoryImpl) MakeReservation(ctx context.Context, orders []*app.StoredItem) error {
	txCtx, ok := ctx.(persistence.TransactionalContext)
	if !ok {
		return appErrors.ErrNotTransactional
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return err
	}

	argsArray := [][]any{}

	for _, order := range orders {
		argsArray = append(argsArray, order.AsQueryArgs())
	}

	_, err = tx.Exec(makeReservationQuery, argsArray)
	if err != nil {
		return errors.Join(errors.New("insert reservation fail"), err)
	}

	return nil
}

func (*_reserveRepositoryImpl) FreeReservation(ctx context.Context, itemsCount map[int]int) error {
	txCtx, ok := ctx.(persistence.TransactionalContext)
	if !ok {
		return appErrors.ErrNotTransactional
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return err
	}

	getReservationsArgs := [][2]int{}

	for _, row := range getReservationsArgs {
		itemCode, count := row[0], row[1]
		getReservationsArgs = append(getReservationsArgs, [2]int{itemCode, count})
	}

	reservationRows, err := tx.Query(getReservations, getReservationsArgs)
	if err != nil {
		return errors.Join(errors.New("get reservation warehouse ids fail"), err)
	}
	defer reservationRows.Close()

	dereserveArgs := [][3]int{}

	for reservationRows.Next() {
		var itemCode, warehouseId, quantity int

		if err := reservationRows.Scan(
			&itemCode,
			&warehouseId,
			&quantity); err != nil {
			return errors.Join(errors.New("scan reservation rows fail"), err)
		}

		dereserveArgs = append(dereserveArgs, [3]int{itemCode, warehouseId, quantity})
	}

	_, err = tx.Exec(dereserveQuery, dereserveArgs)
	if err != nil {
		return errors.Join(errors.New("dereserve fail"), err)
	}

	_, err = tx.Exec(updateStockQuery, dereserveArgs)
	if err != nil {
		return errors.Join(errors.New("update stock fail"), err)
	}

	return nil
}
