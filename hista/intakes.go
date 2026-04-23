package hista

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/internal/medicines"
	"encore.dev/types/uuid"
)

type IntakeResponse struct {
	Date       time.Time `json:"date"`
	MedicineId uint      `json:"medicineId"`
	Count      int64     `json:"count"`
}

type IntakeResponseList struct {
	Intakes []IntakeResponse `json:"intakes"`
}

func toIntakeResponseList(gil medicines.GroupedIntakeList) IntakeResponseList {
	intakes := []IntakeResponse{}
	for _, gi := range gil {
		intakes = append(intakes,
			IntakeResponse{
				Date:       gi.Date,
				MedicineId: gi.MedicineId,
				Count:      gi.Count,
			})
	}
	return IntakeResponseList{Intakes: intakes}
}

// encore:api auth method=GET path=/piid/:piid/intakes
func (s *Service) ListIntakes(ctx context.Context, piid uuid.UUID) (IntakeResponseList, error) {
	list, err := s.meds.ListGroupedIntakes(ctx)
	if err != nil {
		return IntakeResponseList{}, errors.MapError(err)
	}
	return toIntakeResponseList(list), nil
}

// encore:api auth method=POST path=/piid/:piid/intakes/medicines/:medicineId/increment
func (s *Service) IncrementIntake(ctx context.Context, piid uuid.UUID, medicineId uint) error {
	return errors.MapError(s.meds.IncrementIntake(ctx, medicineId))
}

// encore:api auth method=POST path=/piid/:piid/intakes/medicines/:medicineId/decrement
func (s *Service) DecrementIntake(ctx context.Context, piid uuid.UUID, medicineId uint) error {
	return errors.MapError(s.meds.DecrementIntake(ctx, medicineId))
}
