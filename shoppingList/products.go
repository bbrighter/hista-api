package shoppinglist

import (
	"context"

	"encore.app/errors"
	"encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/products
func (s *Service) ListProducts(ctx context.Context, piid uuid.UUID) (entity.ProductListResponse, error) {
	prods, err := s.prod.List(ctx)
	return prods.ToResponse(), errors.MapError(err)
}
