package internal

import (
	"context"

	"encore.dev/types/uuid"
)

type (
	PiidMoveRepo interface {
		Move(ctx context.Context, fromPiid uuid.UUID, toPiid uuid.UUID) error
	}
	PiidMover interface {
		Move(ctx context.Context, fromPiid uuid.UUID, toPiid uuid.UUID) error
		MoveNull(ctx context.Context, toPiid uuid.UUID) error
	}
)

type PiidMoveUseCase struct {
	r PiidMoveRepo
}

func NewPiidMoveUseCase(r PiidMoveRepo) PiidMoveUseCase {
	return PiidMoveUseCase{r: r}
}

func (uc PiidMoveUseCase) Move(ctx context.Context, fromPiid, toPiid uuid.UUID) error {
	return errorMapper(uc.r.Move(ctx, fromPiid, toPiid))
}

func (uc PiidMoveUseCase) MoveNull(ctx context.Context, toPiid uuid.UUID) error {
	return errorMapper(uc.r.Move(ctx, uuid.Nil, toPiid))
}
