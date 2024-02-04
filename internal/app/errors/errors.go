package errors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNotTransactional = errors.New("context is not transactional")
var ErrNoTxInCtx = errors.New("cannot get transaction from context")
