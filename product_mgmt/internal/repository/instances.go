package repository

import (
	"context"

	"encore.app/product_mgmt/entity"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type InstanceRepo struct {
	db *gorm.DB
}

func NewInstanceRepo(db *gorm.DB) *InstanceRepo {
	return &InstanceRepo{db: db}
}

func (r *InstanceRepo) Create(ctx context.Context, name string, productId string) (entity.ProductInstance, error) {
	guid, err := uuid.NewV4()
	if err != nil {
		return entity.ProductInstance{}, err
	}
	var instance = entity.ProductInstance{
		ID:        guid,
		Name:      name,
		ProductId: productId,
	}
	err = gorm.G[entity.ProductInstance](r.db).Create(ctx, &instance)
	return instance, err
}

func (r *InstanceRepo) List(ctx context.Context) ([]entity.ProductInstance, error) {
	return gorm.G[entity.ProductInstance](r.db).Find(ctx)
}

func (r *InstanceRepo) Find(ctx context.Context, id uuid.UUID) (entity.ProductInstance, error) {
	return gorm.G[entity.ProductInstance](r.db).Where("id = ?", id).First(ctx)
}
