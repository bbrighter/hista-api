package api

import (
	"context"

	"encore.dev/types/uuid"
)

// encore:api private method=PATCH path=/internal/product-instance/move/:fromPiid/:toPiid
func (s *Service) MovePiid(ctx context.Context, fromPiid, toPiid uuid.UUID) error {
	return s.move.Move(ctx, fromPiid, toPiid)
}

// encore:api private method=PATCH path=/internal/product-instance/move/:toPiid
func (s *Service) MoveNullPiid(ctx context.Context, toPiid uuid.UUID) error {
	return s.move.MoveNull(ctx, toPiid)
}
