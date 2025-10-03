package repository

import (
	"testing"

	"encore.app/product_mgmt/entity"
	"github.com/stretchr/testify/assert"
)

func newTestProductRepo() *ProductRepo {
	app1 := entity.App{ID: "app1", Name: "App 1"}
	app2 := entity.App{ID: "app2", Name: "App 2"}
	return &ProductRepo{
		apps: []entity.App{app1, app2},
		products: []entity.Product{
			{ID: "test-id", Name: "test-name", Apps: []entity.App{app1}},
		},
	}
}

func TestProductList(t *testing.T) {
	r := newTestProductRepo()

	products := r.ListProducts()

	assert.Len(t, products, 1)

	var histaProduct entity.Product
	for _, p := range products {
		if p.ID == "test-id" {
			histaProduct = p
			break
		}
	}
	assert.Equal(t, "test-name", histaProduct.Name)
	assert.Len(t, histaProduct.Apps, 1)

	var userMgmtApp entity.App
	for _, a := range histaProduct.Apps {
		if a.ID == "app1" {
			userMgmtApp = a
			break
		}
	}
	assert.Equal(t, "App 1", userMgmtApp.Name)
}
