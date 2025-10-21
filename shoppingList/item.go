package shoppinglist

import (
	"context"

	"encore.app/errors"
	entity "encore.app/shoppingList/entity"
)

// encore:api auth method=POST path=/item/:productId
func (s *Service) PostItem(ctx context.Context, productId uint) (entity.IdResponse, error) {
	id, err := s.item.AddItemByProductId(ctx, productId)
	return entity.ToIdResponse(id), errors.MapError(err)
}

type ItemNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/item
func (s *Service) PostItemByName(ctx context.Context, params ItemNameParams) (entity.ItemResponse, error) {
	item, err := s.item.AddItemByName(ctx, params.Name)
	return item.ToResponse(), errors.MapError(err)
}

// encore:api auth method=PATCH path=/item/:itemId/check
func (s *Service) CheckItem(ctx context.Context, itemId uint) error {
	return errors.MapError(s.item.CheckItem(ctx, itemId))
}
