package unitofwork

import (
	"context"

	"encore.app/shoppingList/internal"
	"encore.app/shoppingList/internal/repository"
	"gorm.io/gorm"
)

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) internal.UnitOfWork {
	return UnitOfWork{db: db}
}

func (u UnitOfWork) List() internal.ListRepo {
	return repository.NewListRepo(u.db)
}

func (u UnitOfWork) Item() internal.ItemRepo {
	return repository.NewItemRepo(u.db)
}

func (u UnitOfWork) Product() internal.ProductRepo {
	return repository.NewProductRepo(u.db)
}

func (u UnitOfWork) WithTransaction(ctx context.Context, fn func(tx internal.UnitOfWork) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		r := NewUnitOfWork(tx)
		return fn(r)
	})
}
