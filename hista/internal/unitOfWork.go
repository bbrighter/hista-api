package internal

import "context"

type UnitOfWork interface {
	WithTransaction(ctx context.Context, fn func(tx UnitOfWork) error) error
	Medicine() IMedicineRepo
}
