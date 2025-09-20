package internal

import (
	"encore.app/product_mgmt/entity"
)

type (
	ProductFinderRepo interface {
		Find(id string) (entity.Product, error)
	}
)

type ProductFinder struct {
	pf ProductFinderRepo
}

func NewProductFinder(pf ProductFinderRepo) ProductFinder {
	return ProductFinder{pf: pf}
}

func (p ProductFinder) Find(id string) (entity.Product, error) {
	return p.pf.Find(id)
}
