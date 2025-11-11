package internal

import (
	"context"
	"strings"

	"encore.app/errors"
	"encore.app/shoppingList/entity"
)

type IItemUseCase interface {
	AddItemByProductId(ctx context.Context, listId uint, productId uint) (uint, error)
	AddItemByName(ctx context.Context, listId uint, name string) (*entity.Item, error)
	CheckItem(ctx context.Context, itemId uint) error
	DeleteItem(ctx context.Context, itemIds []uint) error
}

type ItemUseCase struct {
	i   ItemRepo
	p   ProductRepo
	uow UnitOfWork
}

func NewItemUseCase(i ItemRepo, p ProductRepo, uow UnitOfWork) ItemUseCase {
	return ItemUseCase{i: i, p: p, uow: uow}
}

func (uc ItemUseCase) AddItemByProductId(ctx context.Context, listId uint, productId uint) (uint, error) {
	if err := uc.i.CheckUniqueness(ctx, productId, listId); err != nil {
		return 0, err
	}
	id, err := uc.i.Create(ctx, productId, listId)
	return id, errors.MapError(err)
}
func (uc ItemUseCase) AddItemByName(ctx context.Context, listId uint, name string) (*entity.Item, error) {

	var returnItem = new(entity.Item)
	err := uc.uow.WithTransaction(ctx, func(tx UnitOfWork) error {
		trimmedName := strings.TrimSpace(name)
		prodId, err := uc.p.Create(ctx, trimmedName)
		if err != nil {
			return err
		}
		if err := uc.i.CheckUniqueness(ctx, prodId, listId); err != nil {
			return err
		}
		itemId, err := uc.i.Create(ctx, prodId, listId)
		if err != nil {
			return err
		}
		item, err := uc.i.Find(ctx, itemId)
		if err != nil {
			return err
		}
		returnItem = item
		return nil
	})

	return returnItem, err
}
func (uc ItemUseCase) CheckItem(ctx context.Context, itemId uint) error {
	return uc.i.Check(ctx, itemId)
}
func (uc ItemUseCase) DeleteItem(ctx context.Context, itemIds []uint) error {
	return uc.i.Delete(ctx, itemIds)
}
