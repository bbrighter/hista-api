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
