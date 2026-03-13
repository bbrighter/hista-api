package shoppinglist

import (
	"context"

	"encore.app/errors"
	"encore.dev/types/uuid"
)

type PatchProductNameParams struct {
	Name string `json:"name"`
}

// encore:api auth method=PATCH path=/piid/:piid/product/:id
func (s *Service) PatchProductName(ctx context.Context, piid uuid.UUID, id uint, params PatchProductNameParams) error {
	return errors.MapError(s.prod.PatchName(ctx, id, params.Name))
}

// encore:api auth method=DELETE path=/piid/:piid/product/:id
func (s *Service) DeleteProduct(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(s.prod.Delete(ctx, id))
}

type PatchProductArchiveParams struct {
	Archive bool `json:"archive"`
}

// encore:api auth method=PATCH path=/piid/:piid/product/:id/archive
func (s *Service) PatchArchiveProduct(ctx context.Context, piid uuid.UUID, id uint, params PatchProductArchiveParams) error {
	return errors.MapError(s.prod.Archive(ctx, id, params.Archive))
}
