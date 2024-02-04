package postgres

import (
	"context"
	"database/sql"
	appErrors "lamoda_task/internal/app/errors"
	application "lamoda_task/internal/app/persistence"
)

type _transactionalContextImpl struct {
	context.Context
}

type txKey struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return NewTransactionalContext(context.WithValue(ctx, txKey{}, tx), tx)
}

func NewTransactionalContext(ctx context.Context, tx *sql.Tx) application.TransactionalContext {
	return &_transactionalContextImpl{
		ctx,
	}
}

func (txCtx *_transactionalContextImpl) GetTx() (*sql.Tx, error) {
	tx, ok := txCtx.Value(txKey{}).(*sql.Tx)
	if !ok {
		return nil, appErrors.ErrNoTxInCtx
	}

	return tx, nil
}
