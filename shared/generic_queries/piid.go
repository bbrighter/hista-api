package generic_queries

import (
	"context"
	"errors"

	"encore.dev/types/uuid"
)

func PiidFromCtx(ctx context.Context) (uuid.UUID, error) {
	piid, ok := ctx.Value("piid").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("missing PIID in context")
	}
	return piid, nil
}
