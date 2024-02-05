package postgres

import (
	"context"
	"errors"
	"fmt"
	appErrors "lamoda_task/internal/app/errors"
	"lamoda_task/internal/app/models"
	"lamoda_task/internal/app/persistence"
	app "lamoda_task/internal/app/persistence/repositories"
	"strings"
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

	argsArray := [][]int{}

	for _, order := range orders {
		argsArray = append(argsArray, order.AsIntArgs())
	}

	_, err = tx.Exec(makeReservationQuery, argsArray)
	if err != nil {
		return errors.Join(errors.New("insert reservation fail"), err)
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

	getReservationsArgs := []any{}
	rowsStrings := []string{}

	placeholderCounter := 1
	for itemCode, count := range itemsCount {
		getReservationsArgs = append(getReservationsArgs, itemCode, count)
		rowsStrings = append(rowsStrings, strings.Join([]string{fmt.Sprintf("($%d", placeholderCounter), fmt.Sprintf("$%d)", placeholderCounter+1)}, ", "))
		placeholderCounter += 2
	}

	reservationRows, err := tx.Query(strings.Replace(getReservations, "?", strings.Join(rowsStrings, ", "), 1), getReservationsArgs...)
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

	dereserveArgs := (&models.ReservationItem{}).MultipleIntArgs(reservation)

	_, err = tx.Exec(dereserveQuery, dereserveArgs)
	if err != nil {
		return errors.Join(errors.New("dereserve fail"), err)
	}

	return nil
}
