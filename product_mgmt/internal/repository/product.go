package repository

import (
	"errors"
	"fmt"

	"encore.app/product_mgmt/entity"
)

type ProductRepo struct {
	apps     []entity.App
	products []entity.Product
}

func NewProductRepo(config Config) *ProductRepo {
	return &ProductRepo{apps: config.Apps, products: config.Products}
}

func (r ProductRepo) ListProducts() []entity.Product {
	return r.products
}

func (r ProductRepo) ListApps() []entity.App {
	return r.apps
}

func (r ProductRepo) Find(id string) (entity.Product, error) {
	for _, p := range r.products {
		fmt.Println(p)
		if p.ID == id {
			return p, nil
		}
	}
	return entity.Product{}, errors.New("not found")
}
