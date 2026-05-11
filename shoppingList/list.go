package shoppinglist

import (
	"context"

	"encore.app/errors"
	shoppinglist "encore.app/shoppingList/internal/shoppingList"
	"encore.dev/types/uuid"
)

type DeleteListForceDeleteParam struct {
	Force bool `query:"force"`
}

// encore:api auth method=DELETE path=/piid/:piid/list/:listId
func (s *Service) DeleteList(ctx context.Context, piid uuid.UUID, listId uint, params DeleteListForceDeleteParam) error {
	var err error
	if params.Force {
		err = s.sm.ForceDeleteList(ctx, listId)
	} else {
		err = s.sm.DeleteList(ctx, listId)
	}
	return errors.MapError(err)
}

type ListResponse struct {
	ID    uint           `json:"id"`
	Items []ItemResponse `json:"items"`
}

func toListResponse(l shoppinglist.List) ListResponse {
	var items = []ItemResponse{}
	for _, i := range l.Items {
		items = append(items, toItemResponse(i))
	}
	return ListResponse{
		ID:    l.ID,
		Items: items,
	}
}

// encore:api auth method=POST path=/piid/:piid/list
func (s *Service) PostOrGetList(ctx context.Context, piid uuid.UUID) (ListResponse, error) {
	list, err := s.sm.CreateOrFirstList(ctx)
	return toListResponse(list), errors.MapError(err)
}
