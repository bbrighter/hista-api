package repository

import (
	"context"

	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"encore.app/shoppingList/internal"
	"gorm.io/gorm"
)

type ItemRepo struct {
	db *gorm.DB
}

func NewItemRepo(db *gorm.DB) internal.ItemRepo {
	return ItemRepo{db: db}
}

func (r ItemRepo) CheckUniqueness(ctx context.Context, productId uint, listId uint) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	count, err := gorm.G[*entity.Item](r.db).
		Where("pi_id = ?", piid).
		Where("product_id = ?", productId).
		Where("list_id = ?", listId).
		Count(ctx, "*")
	if err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrCheckConstraintViolated
	}
	return nil
}

func (r ItemRepo) Create(ctx context.Context, productId uint, listId uint) (uint, error) {
	var item = &entity.Item{ProductId: productId, ListId: listId}
	err := generic_queries.Create(ctx, r.db, item)
	return item.ID, err
}

func (r ItemRepo) Delete(ctx context.Context, id uint) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	rows, err := gorm.G[*entity.Item](r.db.Unscoped()).Where("pi_id = ?", piid).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r ItemRepo) Check(ctx context.Context, id uint) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	rows, err := gorm.G[entity.Item](r.db).
		Where("pi_id = ?", piid).
		Where("id = ?", id).
		Update(ctx, "checked", gorm.Expr("NOT checked"))
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r ItemRepo) List(ctx context.Context, listId uint) ([]entity.Item, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return gorm.G[entity.Item](r.db).
		Where("pi_id = ?", piid).
		Where("list_id = ?", listId).
		Preload("Product", nil).
		Find(ctx)
}

func (r ItemRepo) Find(ctx context.Context, id uint) (*entity.Item, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return &entity.Item{}, err
	}
	return gorm.G[*entity.Item](r.db).
		Where("pi_id = ?", piid).
		Where("id = ?", id).
		Preload("Product", nil).
		First(ctx)
}
