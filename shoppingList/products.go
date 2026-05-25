package shoppinglist

import (
	"context"

	"encore.app/errors"
	shoppinglist "encore.app/shoppingList/internal/shoppingList"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
)

type ProductResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Archived bool   `json:"archived"`
}
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
}

func toProductListResponse(products []*shoppinglist.Product) ProductListResponse {
	var resps = []ProductResponse{}
	for _, product := range products {
		resps = append(resps, ProductResponse{
			ID:       product.ID,
			Name:     product.Name,
			Archived: product.Archived,
		})
	}
	return ProductListResponse{Products: resps}
}

// encore:api auth method=GET path=/piid/:piid/products
func (s *Service) GetProducts(ctx context.Context, piid uuid.UUID) (ProductListResponse, error) {
	products, err := s.sm.ListProducts(ctx)
	if err != nil {
		return ProductListResponse{}, errors.MapError(err)
	}
	return toProductListResponse(products), nil
}

type PatchProductParams struct {
	Name    option.Option[string] `json:"name" encore:"omitEmpty"`
	Archive option.Option[bool]   `json:"archive" encore:"omitEmpty"`
}

// encore:api auth method=PATCH path=/piid/:piid/products/:id
func (s *Service) PatchProduct(ctx context.Context, piid uuid.UUID, id uint, params PatchProductParams) error {
	err := s.sm.UpdateProduct(ctx, id, params.Name.PtrOrNil(), params.Archive.PtrOrNil())
	return errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/products/:id
func (s *Service) DeleteProduct(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(s.sm.DeleteProduct(ctx, id))
}
