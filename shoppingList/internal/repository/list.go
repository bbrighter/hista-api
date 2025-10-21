package repository

import (
	"context"

	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"encore.app/shoppingList/internal"
	"gorm.io/gorm"
)

type ListRepo struct {
	db *gorm.DB
}

func NewListRepo(db *gorm.DB) internal.ListRepo {
	return ListRepo{db: db}
}

func (r ListRepo) First(ctx context.Context) (*entity.List, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	returnedList, err := gorm.G[*entity.List](r.db).
		Where("pi_id = ?", piid).
		Where("deleted_at IS NULL").
		Preload("Items", nil).
		First(ctx)
	return returnedList, err
}

func (r ListRepo) Create(ctx context.Context) (uint, error) {
	var list = entity.List{}
	err := generic_queries.Create(ctx, r.db, &list)
	return list.ID, err
}

func (r ListRepo) Delete(ctx context.Context, id uint) error {
	return generic_queries.Delete[*entity.List](ctx, r.db, id)
}
