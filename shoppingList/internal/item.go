package internal

import (
	"context"

	"encore.app/shoppingList/entity"
)

type IItemRepo interface {
	List(ctx context.Context) ([]*entity.Item, error)
	Create(ctx context.Context, productId uint, listId uint) (uint, error)
	Delete(ctx context.Context, itemId uint) error
	Check(ctx context.Context, itemId uint) error
	Find(ctx context.Context, id uint) (*entity.Item, error)
}

type IProductRepo interface {
	Create(ctx context.Context, name string) (uint, error)
}

type IItemUseCase interface {
	AddItemByProductId(ctx context.Context, productId uint) (uint, error)
	AddItemByName(ctx context.Context, name string) (*entity.Item, error)
	CheckItem(ctx context.Context, itemId uint) error
	DeleteItem(ctx context.Context, itemId uint) error
}

type ItemUseCase struct {
	i IItemRepo
	p IProductRepo
}

func NewItemUseCase(i IItemRepo, p IProductRepo) ItemUseCase {
	return ItemUseCase{i: i, p: p}
}

func (uc ItemUseCase) AddItemByProductId(ctx context.Context, productId uint) (uint, error) {
	return uc.i.Create(ctx, productId, 1)
}
func (uc ItemUseCase) AddItemByName(ctx context.Context, name string) (*entity.Item, error) {
	// TODO: Or move this into one transaction?
	prodId, err := uc.p.Create(ctx, name)
	if err != nil {
		return &entity.Item{}, err
	}
	itemId, err := uc.i.Create(ctx, prodId, 1)
	if err != nil {
		return &entity.Item{}, err
	}
	return uc.i.Find(ctx, itemId)
}
func (uc ItemUseCase) CheckItem(ctx context.Context, itemId uint) error {
	return uc.i.Check(ctx, itemId)
}
func (uc ItemUseCase) DeleteItem(ctx context.Context, itemId uint) error {
	return uc.i.Delete(ctx, itemId)
}
