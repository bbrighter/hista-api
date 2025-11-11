package internal

import (
	"context"

	"encore.app/shoppingList/entity"
)

type IProductUseCase interface {
	List(ctx context.Context) (entity.Products, error)
}

type ProductUseCase struct {
	p ProductRepo
}

func NewProductUseCase(p ProductRepo) ProductUseCase {
	return ProductUseCase{p: p}
}

func (uc ProductUseCase) List(ctx context.Context) (entity.Products, error) {
	prods, err := uc.p.List(ctx)
	return prods, err
}
