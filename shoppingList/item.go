package shoppinglist

import (
	"context"
	"time"

	"encore.app/errors"
	shoppinglist "encore.app/shoppingList/internal/shoppingList"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
)

type ItemResponse struct {
	ID        uint `json:"id"`
	ProductId uint `json:"productId"`

	Checked   bool                 `json:"checked"`
	Quantity  option.Option[uint8] `json:"quantity"`
	createdAt time.Time
}

func toItemResponse(i shoppinglist.Item) ItemResponse {
	quantity := option.None[uint8]()
	if i.Quantity != nil {
		quantity = option.Some(*i.Quantity)
	}
	return ItemResponse{
		ID:        i.ID,
		ProductId: i.ProductId,
		Checked:   i.Checked,
		Quantity:  quantity,
		createdAt: i.CreatedAt,
	}
}

type IdResponse struct {
	ID uint `json:"id"`
}

// encore:api auth method=POST path=/piid/:piid/list/:listId/item/:productId
func (s *Service) PostItem(ctx context.Context, piid uuid.UUID, listId uint, productId uint) (IdResponse, error) {
	id, err := s.sm.AddItemByProductId(ctx, listId, productId)
	return IdResponse{ID: id}, errors.MapError(err)
}

type ItemNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=POST path=/piid/:piid/list/:listId/item
func (s *Service) PostItemByName(ctx context.Context, piid uuid.UUID, listId uint, params ItemNameParams) (ItemResponse, error) {
	item, _, err := s.sm.AddItemByName(ctx, listId, params.Name)
	return toItemResponse(item), errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/item/:itemId
func (s *Service) DeleteItem(ctx context.Context, piid uuid.UUID, itemId uint) error {
	return errors.MapError(s.sm.DeleteItems(ctx, []uint{itemId}))
}

type ItemPatchParams struct {
	Quantity option.Option[uint8] `json:"quantity"`
	Checked  option.Option[bool]  `json:"checked"`
}

// encore:api auth method=PATCH path=/piid/:piid/item/:itemId
func (s *Service) PatchItem(ctx context.Context, piid uuid.UUID, itemId uint, params ItemPatchParams) error {
	err := s.sm.UpdateItem(ctx, itemId, params.Quantity.PtrOrNil(), params.Checked.PtrOrNil())
	return errors.MapError(err)
}
