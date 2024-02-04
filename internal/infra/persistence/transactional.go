package persistence

import (
	"context"
	"database/sql"
	"errors"
	application "lamoda_task/internal/app/persistence"
)

type _transactionalImpl struct {
	pool *sql.DB
}

func NewTransactional(pool *sql.DB) application.Transactional {
	return &_transactionalImpl{
		pool: pool,
	}
}

func (impl *_transactionalImpl) WithinTransaction(ctx context.Context, action func(context.Context) error) error {
	tx, err := impl.pool.Begin()
	if err != nil {
		return errors.Join(errors.New("start transaction fail"), err)
	}
	defer tx.Rollback()

	err = action(injectTx(ctx, tx))
	if err != nil {
		return errors.Join(errors.New("action fail"), err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(errors.New("commit transaction fail"), err)
	}
	return nil
}
