package product_mgmt

import "context"

// encore:api private method=GET path=/internal/product/:productId
func (s Service) FindProduct(ctx context.Context, productId string) error {
	_, err := s.product.Find(productId)
	return err
}
