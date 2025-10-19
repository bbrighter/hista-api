package shoppinglist

import "context"

type IdResponse struct {
	ID uint `json:"id"`
}

// encore:api method=POST path=/item/:productId
func (s *Service) PostItem(ctx context.Context, productId uint) (IdResponse, error) {
	id, err := s.uc.AddItemByProductId(ctx, productId)
	return IdResponse{ID: id}, err
}

type ItemNameParams struct {
	Name string `json:"name"`
}

type ItemResponse struct {
	ID        uint `json:"id"`
	ProductId uint `json:"productId"`
}

// encore:api method=POST path=/item
func (s *Service) PostItemByName(ctx context.Context, params ItemNameParams) (ItemResponse, error) {
	item, err := s.uc.AddItemByName(ctx, params.Name)
	return ItemResponse{ID: item.ID, ProductId: item.ProductId}, err
}
