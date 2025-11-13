package internal

import (
	"context"

	"encore.app/shoppingList/entity"
)

type (
	UnitOfWork interface {
		WithTransaction(ctx context.Context, fn func(tx UnitOfWork) error) error
		List() ListRepo
		Item() ItemRepo
		Product() ProductRepo
	}

	ListRepo interface {
		First(ctx context.Context) (*entity.List, error)
		Create(ctx context.Context) (uint, error)
		Delete(ctx context.Context, id uint) error
	}

	ItemRepo interface {
		Create(ctx context.Context, productId uint, listId uint) (uint, error)
		Delete(ctx context.Context, itemIds []uint) error
		Find(ctx context.Context, id uint) (*entity.Item, error)
		List(ctx context.Context, listId uint) ([]entity.Item, error)
		Patch(ctx context.Context, id uint, values map[string]any) error
	}

	ProductRepo interface {
		Create(ctx context.Context, name string) (uint, error)
		List(ctx context.Context) ([]*entity.Product, error)
	}
)
