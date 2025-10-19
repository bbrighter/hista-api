package repository

import (
	"context"
	"errors"

	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"gorm.io/gorm"
)

type ListRepo struct {
	db *gorm.DB
}

func NewListRepo(db *gorm.DB) ListRepo {
	return ListRepo{db: db}
}

func (r ListRepo) FirstOrCreate(ctx context.Context) (*entity.List, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var list = new(entity.List)
	err = r.db.Transaction(func(tx *gorm.DB) error {
		returnedList, err := gorm.G[*entity.List](tx).
			Where("pi_id = ?", piid).
			Where("deleted_at IS NULL").
			First(ctx)
		if err == nil {
			list = returnedList
			return nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return generic_queries.Create(ctx, tx, list)
		}
		return err
	})

	return list, err
}
