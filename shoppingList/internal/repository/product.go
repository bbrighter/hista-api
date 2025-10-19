package repository

import (
	"context"

	"encore.app/shared/generic_queries"
	"encore.app/shoppingList/entity"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return ProductRepo{db: db}
}

func (r ProductRepo) Create(ctx context.Context, name string) (uint, error) {
	var prod = &entity.Product{Name: name}
	err := generic_queries.Create(ctx, r.db, prod)
	return prod.ID, err
}
