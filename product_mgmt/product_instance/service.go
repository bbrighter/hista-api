package product_instance

import (
	"context"

	"encore.app/errors"
	"encore.app/product_mgmt/instances"
	"encore.app/product_mgmt/product"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type ProductInstanceService struct {
	p *product.ProductRepo
	i *instances.InstanceRepo
}

func NewProductInstanceService(db *gorm.DB, products []product.Product, apps []product.App) *ProductInstanceService {
	i := instances.NewInstanceRepo(db)
	p := product.NewProductRepo(products, apps)
	return &ProductInstanceService{p: p, i: i}
}

func (s *ProductInstanceService) Create(ctx context.Context, name string, productId string) (uuid.UUID, error) {
	if _, err := s.p.Find(productId); err != nil {
		return uuid.UUID{}, errors.MapError(err)
	}

	instance, err := s.i.Create(ctx, name, productId)
	return instance.ID, errors.MapError(err)
}

func (s *ProductInstanceService) Find(ctx context.Context, id uuid.UUID) (ProductInstance, error) {
	instance, err := s.i.Find(ctx, id)
	if err != nil {
		return ProductInstance{}, errors.MapError(err)
	}

	product, err := s.p.Find(instance.ProductId)
	if err != nil {
		return ProductInstance{}, errors.MapError(err)
	}
	return toProductInstance(instance, product), nil
}

func (s *ProductInstanceService) List(ctx context.Context) ([]ProductInstance, error) {
	productInstances := []ProductInstance{}
	instances, err := s.i.List(ctx)
	for _, inst := range instances {
		prod, err := s.p.Find(inst.ProductId)
		if err != nil {
			return []ProductInstance{}, err
		}
		productInstances = append(productInstances, toProductInstance(inst, prod))

	}
	return productInstances, err
}

func toProductInstance(inst instances.Instance, product product.Product) ProductInstance {
	return ProductInstance{
		ID:        inst.ID,
		Name:      inst.Name,
		ProductId: inst.ProductId,
		Product:   product,
	}
}
