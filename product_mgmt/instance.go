package product_mgmt

import (
	"context"

	"encore.app/product_mgmt/entity"

	"encore.dev/types/uuid"
)

type ProductInstanceParams struct {
	ProductId    string `json:"productId"`
	InstanceName string `json:"instanceName"`
}

type UuidResponse struct {
	ID uuid.UUID `json:"id"`
}

// encore:api private method=POST path=/instance
func (s *Service) CreateInstance(ctx context.Context, params ProductInstanceParams) (UuidResponse, error) {
	id, err := s.instance.Create(ctx, params.InstanceName, params.ProductId)
	return UuidResponse{ID: id}, err
}

// encore:api private method=GET path=/instance/:id
func (s *Service) FindInstance(ctx context.Context, id uuid.UUID) (entity.ProductInstanceResponse, error) {
	instance, err := s.instance.Find(ctx, id)
	return instance.ToResponse(), err
}
