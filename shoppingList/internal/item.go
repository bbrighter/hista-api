package internal

import (
	"context"
	"strings"

	"encore.app/shoppingList/entity"
	"gorm.io/gorm"
)

type IItemUseCase interface {
	AddItemByProductId(ctx context.Context, listId uint, productId uint) (uint, error)
	AddItemByName(ctx context.Context, listId uint, name string) (*entity.Item, error)
	CheckItem(ctx context.Context, itemId uint) error
	DeleteItem(ctx context.Context, itemIds []uint) error
	PatchItemQuantity(ctx context.Context, itemId uint, quantity *uint8) error
}

type ItemUseCase struct {
	i   ItemRepo
	p   ProductRepo
	uow UnitOfWork
	m   MomentRepo
}

func NewItemUseCase(i ItemRepo, p ProductRepo, uow UnitOfWork, m MomentRepo) ItemUseCase {
	return ItemUseCase{i: i, p: p, uow: uow, m: m}
}

func (uc ItemUseCase) AddItemByProductId(ctx context.Context, listId uint, productId uint) (uint, error) {
	if err := uc.m.Update(ctx); err != nil {
		return 0, err
	}
	return uc.i.Create(ctx, productId, listId)
}

func (uc ItemUseCase) AddItemByName(ctx context.Context, listId uint, name string) (*entity.Item, error) {
	var returnItem = new(entity.Item)
	err := uc.uow.WithTransaction(ctx, func(tx UnitOfWork) error {
		if err := uc.m.Update(ctx); err != nil {
			return err
		}
		trimmedName := strings.TrimSpace(name)
		prodId, err := uc.p.Create(ctx, trimmedName)
		if err != nil {
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
	if err := uc.m.Update(ctx); err != nil {
		return err
	}
	return uc.i.Patch(ctx, itemId, map[string]any{"checked": gorm.Expr("NOT checked")})
}
func (uc ItemUseCase) DeleteItem(ctx context.Context, itemIds []uint) error {
	if err := uc.m.Update(ctx); err != nil {
		return err
	}
	return uc.i.Delete(ctx, itemIds)
}

func (uc ItemUseCase) PatchItemQuantity(ctx context.Context, itemId uint, quantity *uint8) error {
	if err := uc.m.Update(ctx); err != nil {
		return err
	}
	return uc.i.Patch(ctx, itemId, map[string]any{"quantity": quantity})
}
