package shoppinglist

import (
	"context"

	"encore.app/errors"
	entity "encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=POST path=/piid/:piid/list/:listId/item/:productId
func (s *Service) PostItem(ctx context.Context, piid uuid.UUID, listId uint, productId uint) (entity.IdResponse, error) {
	id, err := s.item.AddItemByProductId(ctx, listId, productId)
	return entity.ToIdResponse(id), errors.MapError(err)
}

type ItemNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/piid/:piid/list/:listId/item
func (s *Service) PostItemByName(ctx context.Context, piid uuid.UUID, listId uint, params ItemNameParams) (entity.ItemResponse, error) {
	item, err := s.item.AddItemByName(ctx, listId, params.Name)
	return item.ToResponse(), errors.MapError(err)
}

// encore:api auth method=PATCH path=/piid/:piid/item/:itemId/check
func (s *Service) CheckItem(ctx context.Context, piid uuid.UUID, itemId uint) error {
	return errors.MapError(s.item.CheckItem(ctx, itemId))
}

// encore:api auth method=DELETE path=/piid/:piid/item/:itemId
func (s *Service) DeleteItem(ctx context.Context, piid uuid.UUID, itemId uint) error {
	return errors.MapError(s.item.DeleteItem(ctx, []uint{itemId}))
}
