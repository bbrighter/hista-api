package internal

import (
	"context"
	"errors"
	"testing"

	"encore.app/product_mgmt/entity"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
)

type testInstanceStoreOk struct{}

func (r testInstanceStoreOk) Create(ctx context.Context, name string, productId string) (entity.ProductInstance, error) {
	guid, _ := uuid.FromString("a821afdc-5e80-4ec5-8f7b-ee42f7635304")
	return entity.ProductInstance{ID: guid}, nil
}
func (r testInstanceStoreOk) Find(ctx context.Context, id uuid.UUID) (entity.ProductInstance, error) {
	return entity.ProductInstance{}, nil
}
func (r testInstanceStoreOk) List(ctx context.Context) ([]entity.ProductInstance, error) {
	return []entity.ProductInstance{}, nil
}

type testProductFinderOk struct{}

func (r testProductFinderOk) Find(id string) (entity.Product, error) {
	return entity.Product{ID: id}, nil
}

type testProductFinderNotOk struct{}

func (r testProductFinderNotOk) Find(id string) (entity.Product, error) {
	return entity.Product{}, errors.New("not found")
}

func TestCreateInstances(t *testing.T) {
	uc := NewInstanceUseCase(testInstanceStoreOk{}, testProductFinderOk{})
	ctx := context.Background()

	guid, err := uc.Create(ctx, "new name", "hista")

	assert.NoError(t, err)
	assert.Equal(t, "a821afdc-5e80-4ec5-8f7b-ee42f7635304", guid.String())
}

func TestCreateInstancesNotOk(t *testing.T) {
	uc := NewInstanceUseCase(testInstanceStoreOk{}, testProductFinderNotOk{})
	ctx := context.Background()

	_, err := uc.Create(ctx, "new name", "hista")

	assert.Error(t, err)
}
