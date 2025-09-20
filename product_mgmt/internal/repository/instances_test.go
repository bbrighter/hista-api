package repository

import (
	"context"
	"testing"

	"encore.app/product_mgmt/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const productTestId string = "prod"

func initTest(t *testing.T) (*InstanceRepo, context.Context) {
	ctx := context.Background()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(
		&entity.ProductInstance{},
	)
	assert.NoError(t, err)
	db.Exec("PRAGMA foreign_keys = ON;")

	return NewInstanceRepo(db), ctx
}

func TestCreateInstance(t *testing.T) {
	r, ctx := initTest(t)

	tests := map[string]struct {
		productId     string
		expectedError bool
	}{
		"ok": {productId: productTestId},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			instance, err := r.Create(ctx, "instance name", test.productId)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "instance name", instance.Name)
				assert.NotNil(t, instance.ID)
				instances, err := gorm.G[entity.ProductInstance](r.db).Find(ctx)
				assert.NoError(t, err)
				assert.Len(t, instances, 1)
			}
		})
	}
}
