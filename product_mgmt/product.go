package product_mgmt

import (
	"context"

	entity "encore.app/product_mgmt/entity"
)

// encore:api private method=GET path=/internal/product/:productId
func (s Service) FindProduct(ctx context.Context, productId string) (entity.Product, error) {
	prod, err := s.product.Find(productId)
	return prod, err
}
