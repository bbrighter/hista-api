package internal

import (
	"context"

	"encore.app/product_mgmt/entity"
	"encore.dev/types/uuid"
)

type (
	ProductInstanceStoreRepo interface {
		Find(ctx context.Context, id uuid.UUID) (entity.ProductInstance, error)
		Create(ctx context.Context, name string, productId string) (entity.ProductInstance, error)
	}
	InstanceStore interface {
		Create(ctx context.Context, name string, productId string) (uuid.UUID, error)
		Find(ctx context.Context, id uuid.UUID) (entity.ProductInstance, error)
	}
)

type InstanceUseCase struct {
	pf ProductFinderRepo
	is ProductInstanceStoreRepo
}

func NewInstanceUseCase(is ProductInstanceStoreRepo, pf ProductFinderRepo) InstanceUseCase {
	return InstanceUseCase{
		is: is,
		pf: pf,
	}
}

func (uc InstanceUseCase) Create(ctx context.Context, name string, productId string) (uuid.UUID, error) {
	if _, err := uc.pf.Find(productId); err != nil {
		return uuid.UUID{}, err
	}

	instance, err := uc.is.Create(ctx, name, productId)
	return instance.ID, err
}

func (uc InstanceUseCase) Find(ctx context.Context, id uuid.UUID) (entity.ProductInstance, error) {
	instance, err := uc.is.Find(ctx, id)
	if err != nil {
		return entity.ProductInstance{}, err
	}

	product, err := uc.pf.Find(instance.ProductId)
	if err != nil {
		return entity.ProductInstance{}, err
	}
	instance.Product = product
	return instance, nil
}
