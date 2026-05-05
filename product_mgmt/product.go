package product_mgmt

import (
	"context"

	"encore.app/errors"
	"encore.app/product_mgmt/product"
)

type ProductResponse struct {
	ID   string        `json:"id"`
	Name string        `json:"name"`
	Apps []AppResponse `json:"apps"`
}

type AppResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func toProductResponse(prod product.Product) ProductResponse {
	return ProductResponse{
		ID:   prod.ID,
		Name: prod.Name,
		Apps: toAppResponse(prod.Apps),
	}
}

func toAppResponse(apps []product.App) []AppResponse {
	resp := []AppResponse{}
	for _, app := range apps {
		resp = append(resp, AppResponse{ID: app.ID, Name: app.Name})
	}
	return resp
}

// encore:api private method=GET path=/internal/product/:productId
func (s Service) FindProduct(ctx context.Context, productId string) (ProductResponse, error) {
	prod, err := s.prod.Find(productId)
	if err != nil {
		return ProductResponse{}, errors.MapError(err)
	}
	return toProductResponse(prod), nil
}
