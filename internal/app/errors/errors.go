package errors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNotTransactional = errors.New("context is not transactional")
var ErrNoTxInCtx = errors.New("cannot get transaction from context")
var ErrImpossibleReserve = errors.New("impossible to reserve: not enough items")

const ImpossibleStatusCode = 432
