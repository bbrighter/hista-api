package status

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type StatusService struct {
	s *statusRepo
}

func NewStatusService(db *gorm.DB) *StatusService {
	s := newStatusRepo(db)
	return &StatusService{s: s}
}

func (s *StatusService) ListStatuses(ctx context.Context) ([]*Status, error) {
	return s.s.ListStatus(ctx)
}

func (s *StatusService) CreateStatus(ctx context.Context, date time.Time) (*Status, error) {
	status := &Status{Date: date}
	err := s.s.CreateStatus(ctx, status)
	return status, err
}

func (s *StatusService) DeleteStatus(ctx context.Context, id uint) error {
	return s.s.DeleteStatus(ctx, id)
}

func (s *StatusService) UpdateStatus(
	ctx context.Context,
	id uint,
	params struct {
		Date           *time.Time
		MorningFitness *int
		MorningSleep   *int
		EveningFitness *int
	},
) error {
	values := make(map[string]any)
	if params.Date != nil {
		values["date"] = *params.Date
	}
	if params.MorningFitness != nil {
		values["morning_fitness"] = *params.MorningFitness
	}
	if params.MorningSleep != nil {
		values["morning_sleep"] = *params.MorningSleep
	}
	if params.EveningFitness != nil {
		values["evening_fitness"] = *params.EveningFitness
	}
	return s.s.UpdateStatus(ctx, id, values)
}
