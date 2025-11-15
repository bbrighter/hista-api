package internal

import (
	"context"
	"fmt"

	"encore.app/errors"
	"encore.app/shoppingList/entity"
)

type IListUseCase interface {
	FirstOrCreate(ctx context.Context) (*entity.List, error)
	Delete(ctx context.Context, id uint, force bool) error
}

type ListUseCase struct {
	uow UnitOfWork
	li  ListRepo
	it  ItemRepo
	m   MomentRepo
}

func NewListUseCase(r ListRepo, it ItemRepo, uow UnitOfWork, m MomentRepo) ListUseCase {
	return ListUseCase{li: r, it: it, uow: uow, m: m}
}

func (l ListUseCase) FirstOrCreate(ctx context.Context) (*entity.List, error) {
	var returnList = new(entity.List)
	err := l.uow.WithTransaction(ctx, func(tx UnitOfWork) error {
		list, err := tx.List().First(ctx)
		if err == nil {
			returnList = list
			return nil
		}
		listId, err := tx.List().Create(ctx)
		returnList = &entity.List{ID: listId}
		return err
	})
	return returnList, err
}

func (l ListUseCase) Delete(ctx context.Context, id uint, force bool) error {
	err := l.uow.WithTransaction(ctx, func(tx UnitOfWork) error {
		if err := l.m.Update(ctx); err != nil {
			return err
		}
		items, err := tx.Item().List(ctx, id)
		if err != nil {
			return err
		}
		var uncheckedIds []uint
		for _, item := range items {
			if !item.Checked {
				uncheckedIds = append(uncheckedIds, item.ID)
			}
		}
		if len(uncheckedIds) > 0 && !force {
			return fmt.Errorf("%w with item id %v", errors.ErrUncheckedItems, uncheckedIds)
		}
		if len(uncheckedIds) > 0 && force {
			if err := tx.Item().Delete(ctx, uncheckedIds); err != nil {
				return err
			}
		}

		return tx.List().Delete(ctx, id)
	})

	return err
}
