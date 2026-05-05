package instances

import (
	"context"

	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type InstanceRepo struct {
	db *gorm.DB
}

func NewInstanceRepo(db *gorm.DB) *InstanceRepo {
	return &InstanceRepo{db: db}
}

func (r *InstanceRepo) Create(ctx context.Context, name string, productId string) (Instance, error) {
	guid, err := uuid.NewV4()
	if err != nil {
		return Instance{}, err
	}
	var instance = Instance{
		ID:        guid,
		Name:      name,
		ProductId: productId,
	}
	err = gorm.G[Instance](r.db).Create(ctx, &instance)
	return instance, err
}

func (r *InstanceRepo) List(ctx context.Context) ([]Instance, error) {
	return gorm.G[Instance](r.db).Find(ctx)
}

func (r *InstanceRepo) Find(ctx context.Context, id uuid.UUID) (Instance, error) {
	return gorm.G[Instance](r.db).Where("id = ?", id).First(ctx)
}
