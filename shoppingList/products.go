package shoppinglist

import (
	"context"

	"encore.app/errors"
	"encore.dev/types/option"
	"encore.dev/types/uuid"
)

type PatchProductParams struct {
	Name    option.Option[string] `json:"name" encore:"omitEmpty"`
	Archive option.Option[bool]   `json:"archive" encore:"omitEmpty"`
}

// encore:api auth method=PATCH path=/piid/:piid/product/:id
func (s *Service) PatchProduct(ctx context.Context, piid uuid.UUID, id uint, params PatchProductParams) error {
	err := s.sm.UpdateProduct(ctx, id, params.Name.PtrOrNil(), params.Archive.PtrOrNil())
	return errors.MapError(err)
}

// encore:api auth method=DELETE path=/piid/:piid/product/:id
func (s *Service) DeleteProduct(ctx context.Context, piid uuid.UUID, id uint) error {
	return errors.MapError(s.sm.DeleteProduct(ctx, id))
}
