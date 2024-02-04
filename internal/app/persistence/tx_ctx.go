package persistence

import (
	"context"
	"database/sql"
)

type TransactionalContext interface {
	context.Context
	GetTx() (*sql.Tx, error)
}
