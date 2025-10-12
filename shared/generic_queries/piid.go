package generic_queries

import (
	"context"

	"encore.app/errors"
	"encore.dev/types/uuid"
)

func PiidFromCtx(ctx context.Context) (uuid.UUID, error) {
	piid, ok := ctx.Value("piid").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.PiidMissing
	}
	return piid, nil
}
