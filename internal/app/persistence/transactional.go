package persistence

import "context"

type Transactional interface {
	WithinTransaction(context.Context, func(context.Context) error) error
}
