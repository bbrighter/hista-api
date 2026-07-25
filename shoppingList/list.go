package shoppinglist

import (
	"context"
	"sort"

	"encore.app/errors"
	shoppinglist "encore.app/shoppingList/internal/shoppingList"
	"encore.dev/types/uuid"
)

// encore:api auth method=DELETE path=/piid/:piid/list/:listId
func (s *Service) DeleteList(ctx context.Context, piid uuid.UUID, listId uint) error {
	err := s.sm.DeleteList(ctx, listId)
	return errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/list/:listId/force
func (s *Service) ForceDeleteList(ctx context.Context, piid uuid.UUID, listId uint) error {
	return errors.MapError(s.sm.ForceDeleteList(ctx, listId))
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
	sort.Slice(items, func(i, j int) bool {
		return items[j].createdAt.Before(items[i].createdAt)
	})
	return ListResponse{
		ID:    l.ID,
		Items: items,
	}
}

// encore:api auth method=DELETE path=/piid/:piid/list/:listId/move
func (s *Service) DeleteListAndMoveItems(ctx context.Context, piid uuid.UUID, listId uint) (ListResponse, error) {
	list, err := s.sm.DeleteListCreateNewAndMoveItems(ctx, listId)
	return toListResponse(list), errors.MapError(err)
}

// encore:api auth method=POST path=/piid/:piid/list
func (s *Service) PostOrGetList(ctx context.Context, piid uuid.UUID) (ListResponse, error) {
	list, err := s.sm.CreateOrFirstList(ctx)
	return toListResponse(list), errors.MapError(err)
}
