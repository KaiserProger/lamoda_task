package postgres

import (
	"context"
	"errors"
	appErrors "lamoda_task/internal/app/errors"
	"lamoda_task/internal/app/models"
	"lamoda_task/internal/app/persistence"
	app "lamoda_task/internal/app/persistence/repositories"
)

type _reserveRepositoryImpl struct{}

func NewReserveRepository() app.ReserveRepository {
	return &_reserveRepositoryImpl{}
}

func (*_reserveRepositoryImpl) MakeReservation(ctx context.Context, orders []*models.ReservationItem) error {
	txCtx, ok := ctx.(persistence.TransactionalContext)
	if !ok {
		return appErrors.ErrNotTransactional
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return err
	}

	argsArray := (&models.ReservationItem{}).MultipleIntArgs(orders)

	args := make([]any, len(argsArray))
	for ix, arg := range argsArray {
		args[ix] = arg
	}

	query, err := GenBulkPlaceholders(makeReservationQuery, args, 3)
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Join(errors.New("insert reservation fail"), err)
	}

	query, err = GenBulkPlaceholders(removeFromStockQuery, args, 3)
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Join(errors.New("remove from stock fail"), err)
	}

	return nil
}

func (*_reserveRepositoryImpl) GetReservation(ctx context.Context, itemsCount map[int]int) ([]*models.ReservationItem, error) {
	txCtx, ok := ctx.(persistence.TransactionalContext)
	if !ok {
		return nil, appErrors.ErrNotTransactional
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return nil, err
	}

	reservationItems := make([]*models.ReservationItem, 0)

	// Prepare get reservations query as it contains IN.
	// This is a HUGE crutch, but I NEEDED IN to work with arrays.

	itemCodes := make([]int, 0, len(itemsCount))
	for k := range itemsCount {
		itemCodes = append(itemCodes, k)
	}

	reservationRows, err := tx.Query(getReservationsQuery, itemCodes)
	if err != nil {
		return nil, errors.Join(errors.New("get reservation warehouse ids fail"), err)
	}
	defer reservationRows.Close()

	// Crutch end.

	for reservationRows.Next() {
		var reservation models.ReservationItem

		if err := reservationRows.Scan(
			&reservation.ItemCode,
			&reservation.WarehouseId,
			&reservation.Quantity); err != nil {
			return nil, errors.Join(errors.New("get reservations fail"), err)
		}

		reservationItems = append(reservationItems, &reservation)
	}

	return reservationItems, nil
}

func (*_reserveRepositoryImpl) FreeReservation(ctx context.Context, reservation []*models.ReservationItem) error {
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

	query, err := GenBulkPlaceholders(dereserveQuery, args, 3)
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Join(errors.New("dereserve fail"), err)
	}

	return nil
}
