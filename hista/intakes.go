package hista

import (
	"context"

	"encore.app/errors"
	"encore.app/hista/entity"
	"encore.dev/types/uuid"
)

// encore:api auth method=GET path=/piid/:piid/intakes
func (s *Service) ListIntakes(ctx context.Context, piid uuid.UUID) (entity.IntakeResponseList, error) {
	list, err := s.intake.List(ctx)
	return list.ToResponse(), errors.MapError(err)
}

// encore:api auth method=POST path=/piid/:piid/intakes/medicines/:medicineId/increment
func (s *Service) IncrementIntake(ctx context.Context, piid uuid.UUID, medicineId uint) error {
	return errors.MapError(s.intake.Increment(ctx, medicineId))
}

// encore:api auth method=POST path=/piid/:piid/intakes/medicines/:medicineId/decrement
func (s *Service) DecrementIntake(ctx context.Context, piid uuid.UUID, medicineId uint) error {
	return errors.MapError(s.intake.Decrement(ctx, medicineId))
}
