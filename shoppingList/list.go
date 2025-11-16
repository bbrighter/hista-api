package shoppinglist

import (
	"context"

	"encore.app/errors"
	"encore.app/shoppingList/entity"
	"encore.dev/types/uuid"
)

type DeleteListForceDeleteParam struct {
	Force bool `query:"force"`
}

// encore:api auth method=DELETE path=/piid/:piid/list/:listId
func (s *Service) DeleteList(ctx context.Context, piid uuid.UUID, listId uint, params DeleteListForceDeleteParam) error {
	return errors.MapError(s.list.Delete(ctx, listId, params.Force))
}

// encore:api auth method=POST path=/piid/:piid/list
func (s *Service) PostList(ctx context.Context, piid uuid.UUID) (entity.IdResponse, error) {
	list, err := s.list.Create(ctx)
	return entity.ToIdResponse(list.ID), errors.MapError(err)
}
