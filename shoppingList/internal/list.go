package internal

import (
	"context"
	"fmt"

	"encore.app/errors"
	"encore.app/shoppingList/entity"
)

type IListUseCase interface {
	FirstOrCreate(ctx context.Context) (*entity.List, error)
	Delete(ctx context.Context, id uint) error
}

type ListUseCase struct {
	uow UnitOfWork
	li  ListRepo
	it  ItemRepo
}

func NewListUseCase(r ListRepo, it ItemRepo, uow UnitOfWork) ListUseCase {
	return ListUseCase{li: r, it: it, uow: uow}
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

func (l ListUseCase) Delete(ctx context.Context, id uint) error {
	err := l.uow.WithTransaction(ctx, func(tx UnitOfWork) error {
		items, err := tx.Item().List(ctx, id)
		if err != nil {
			return err
		}
		for _, item := range items {
			if !item.Checked {
				return errors.NewErrBadRequest(fmt.Sprintf("item unchecked: id = %d", item.ID))
			}
		}
		return tx.List().Delete(ctx, id)
	})

	return err
}
