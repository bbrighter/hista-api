package product

import (
	"fmt"

	"gorm.io/gorm"
)

type ProductRepo struct {
	apps     []App
	products []Product
}

func NewProductRepo(products []Product, apps []App) *ProductRepo {
	return &ProductRepo{apps: apps, products: products}
}

func (r ProductRepo) ListProducts() []Product {
	return r.products
}

func (r ProductRepo) ListApps() []App {
	return r.apps
}

func (r ProductRepo) Find(id string) (Product, error) {
	for _, p := range r.products {
		fmt.Println(p)
		if p.ID == id {
			return p, nil
		}
	}
	return Product{}, gorm.ErrRecordNotFound
}
