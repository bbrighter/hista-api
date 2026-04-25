package headaches

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type HeadacheService struct {
	h *HeadacheRepo
}

func NewHeadacheService(db *gorm.DB) *HeadacheService {
	h := NewHeadacheRepo(db)
	return &HeadacheService{h: h}
}

func (s *HeadacheService) ListHeadaches(ctx context.Context) ([]*Headache, error) {
	return s.h.ListHeadaches(ctx)
}

func (s *HeadacheService) FirstHeadache(ctx context.Context, id uint) (*Headache, error) {
	return s.h.FirstHeadache(ctx, id)
}

func (s *HeadacheService) DeleteHeadache(ctx context.Context, id uint) error {
	return s.h.DeleteHeadache(ctx, id)
}

func (s *HeadacheService) CreateHeadache(ctx context.Context, date time.Time, severity uint8) (uint, error) {
	headache := &Headache{Date: date, Severity: severity}
	err := s.h.CreateHeadache(ctx, headache)
	return headache.ID, err
}

type HeadacheUpdateParams struct {
	Date        *time.Time
	Positions   *HeadachePositions
	Severity    *uint8
	Symptoms    *HeadacheSymptoms
	Types       *HeadacheTypes
	Description *string
}

func (s *HeadacheService) UpdateHeadache(ctx context.Context, id uint, params HeadacheUpdateParams) error {
	values := make(map[string]any)
	if params.Date != nil {
		values["date"] = *params.Date
	}
	if params.Positions != nil {
		values["positions"] = *params.Positions
	}
	if params.Severity != nil {
		values["severity"] = *params.Severity
	}
	if params.Symptoms != nil {
		values["symptoms"] = *params.Symptoms
	}
	if params.Types != nil {
		values["types"] = *params.Types
	}
	if params.Description != nil {
		values["description"] = *params.Description
	}
	return s.h.UpdateHeadache(ctx, id, values)
}
