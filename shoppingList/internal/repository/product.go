package repository

import (
	"context"

	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"encore.app/shoppingList/internal"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) internal.ProductRepo {
	return ProductRepo{db: db}
}

func (r ProductRepo) Create(ctx context.Context, name string) (uint, error) {
	var prod = &entity.Product{Name: name}
	err := generic_queries.Create(ctx, r.db, prod)
	return prod.ID, err
}

func (r ProductRepo) List(ctx context.Context) ([]*entity.Product, error) {
	return generic_queries.List[*entity.Product](ctx, r.db)
}

func (r ProductRepo) Update(ctx context.Context, id uint, values map[string]any) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	return r.db.Where("pi_id = ?", piid).Updates(values).Error
}

func (r ProductRepo) Delete(ctx context.Context, id uint) error {
	return generic_queries.Delete[*entity.Product](ctx, r.db, id)
}
