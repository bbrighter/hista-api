package shoppinglist

import (
	"context"

	"encore.app/errors"
	"encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/list
func (s *Service) GetOrCreateList(ctx context.Context, piid uuid.UUID) (entity.ListResponse, error) {
	list, err := s.list.FirstOrCreate(ctx)
	return list.ToResponse(), errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/list/:listId
func (s *Service) DeleteList(ctx context.Context, piid uuid.UUID, listId uint) error {
	return errors.MapError(s.list.Delete(ctx, listId))
}
