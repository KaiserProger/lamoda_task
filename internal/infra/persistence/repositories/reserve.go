package repositories

import (
	"context"
	"errors"
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
		return errors.New("context is not transactional")
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return errors.Join(errors.New("get transaction from context fail"), err)
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

func (*_reserveRepositoryImpl) FreeReservation(ctx context.Context, itemCodes []int) error {
	txCtx, ok := ctx.(persistence.TransactionalContext)
	if !ok {
		return errors.New("context is not transactional")
	}

	tx, err := txCtx.GetTx()
	if err != nil {
		return errors.Join(errors.New("get transaction from context fail"), err)
	}

	for _, itemCode := range itemCodes {
		tx.Exec(dereserveQuery, itemCode)
	}
}
