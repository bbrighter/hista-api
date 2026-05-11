package shoppingmoments

import (
	"context"

	sl "encore.app/shoppingList/internal/shoppingList"
	"gorm.io/gorm"
)

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db.Debug()}
}

func (uow *UnitOfWork) ShoppingList() *sl.ShoppingListRepo {
	return sl.NewShoppingListRepo(uow.db)
}

func (uow *UnitOfWork) WithTransaction(ctx context.Context, fn func(*UnitOfWork) error) error {
	return uow.db.Transaction(func(tx *gorm.DB) error {
		r := NewUnitOfWork(tx)
		return fn(r)
	})
}
