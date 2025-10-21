package shoppinglist

import (
	"context"

	"encore.app/errors"
	"encore.app/shoppingList/entity"
)

// encore:api auth method=GET path=/products
func (s *Service) ListProducts(ctx context.Context) (entity.ProductListResponse, error) {
	prods, err := s.prod.List(ctx)
	return prods.ToResponse(), errors.MapError(err)
}
