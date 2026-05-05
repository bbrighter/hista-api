package product_mgmt

import (
	"context"

	"encore.app/errors"
	"encore.app/product_mgmt/product_instance"

	"encore.dev/types/uuid"
)

type ProductInstanceParams struct {
	ProductId    string `json:"productId"`
	InstanceName string `json:"instanceName"`
}

type UuidResponse struct {
	ID uuid.UUID `json:"id"`
}

type ProductInstanceResponse struct {
	ID      uuid.UUID       `json:"id"`
	Name    string          `json:"name"`
	Product ProductResponse `json:"product"`
}

type ProductInstanceListResponse struct {
	Instances []ProductInstanceResponse `json:"instances"`
}

func toProductInstanceResponse(p product_instance.ProductInstance) ProductInstanceResponse {
	return ProductInstanceResponse{
		ID:      p.ID,
		Name:    p.Name,
		Product: toProductResponse(p.Product),
	}
}

func toProductInstanceListResponse(ps []product_instance.ProductInstance) ProductInstanceListResponse {
	var instances = []ProductInstanceResponse{}
	for _, p := range ps {
		instances = append(instances, toProductInstanceResponse(p))
	}
	return ProductInstanceListResponse{Instances: instances}
}

// encore:api private method=POST path=/internal/instance
func (s *Service) CreateInstance(ctx context.Context, params ProductInstanceParams) (UuidResponse, error) {
	id, err := s.pi.Create(ctx, params.InstanceName, params.ProductId)
	return UuidResponse{ID: id}, err
}

// encore:api private method=GET path=/internal/instance
func (s *Service) ListInstances(ctx context.Context) (ProductInstanceListResponse, error) {
	instances, err := s.pi.List(ctx)
	if err != nil {
		return ProductInstanceListResponse{}, errors.MapError(err)
	}
	return toProductInstanceListResponse(instances), nil
}

// encore:api private method=GET path=/internal/instance/:id
func (s *Service) FindInstance(ctx context.Context, id uuid.UUID) (ProductInstanceResponse, error) {
	instance, err := s.pi.Find(ctx, id)
	if err != nil {
		return ProductInstanceResponse{}, errors.MapError(err)
	}
	return toProductInstanceResponse(instance), nil
}
