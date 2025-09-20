package product_mgmt

import (
	"context"
	"testing"

	"encore.app/product_mgmt/entity"
	"encore.app/product_mgmt/internal"
	"encore.app/product_mgmt/internal/repository"
	"encore.dev/et"
	"encore.dev/types/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initTestService(t *testing.T) *Service {
	sqlDb, err := et.NewTestDatabase(context.Background(), "product_mgmt_db")
	assert.NoError(t, err)
	db, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb.Stdlib(),
	}))

	app1 := entity.App{ID: "app1", Name: "App 1"}
	config := repository.Config{
		Products: []entity.Product{{ID: "product-id", Name: "Product Name", Apps: []entity.App{app1}}},
		Apps:     []entity.App{app1},
	}
	p := repository.NewProductRepo(config)
	i := repository.NewInstanceRepo(db)

	instance := internal.NewInstanceUseCase(i, p)

	var service = &Service{
		instance: instance,
	}
	return service
}

func TestCreateInstance(t *testing.T) {
	tests := map[string]struct {
		productId   string
		expectError bool
	}{
		"ok":        {productId: "product-id"},
		"not found": {productId: "does not exist", expectError: true},
	}

	ctx := t.Context()
	service := initTestService(t)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := service.CreateInstance(ctx, ProductInstanceParams{
				ProductId:    test.productId,
				InstanceName: "new instance",
			})

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindInstance(t *testing.T) {
	ctx := t.Context()
	service := initTestService(t)
	resp, err := service.CreateInstance(ctx, ProductInstanceParams{ProductId: "product-id", InstanceName: "name"})
	assert.NoError(t, err)
	guid := resp.ID
	otherGuid, _ := uuid.NewV4()

	tests := map[string]struct {
		instanceId  uuid.UUID
		expectError bool
	}{
		"ok":        {instanceId: guid},
		"not found": {instanceId: otherGuid, expectError: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			instance, err := service.FindInstance(ctx, test.instanceId)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, guid, instance.ID)
				assert.Equal(t, "product-id", instance.Product.ID)
				assert.Equal(t, "Product Name", instance.Product.Name)
				assert.Equal(t, "name", instance.Name)
				assert.Len(t, instance.Product.Apps, 1)
				assert.Equal(t, "app1", instance.Product.Apps[0].ID)
				assert.Equal(t, "App 1", instance.Product.Apps[0].Name)
			}
		})
	}

}
