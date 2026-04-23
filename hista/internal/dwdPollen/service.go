package dwdPollen

import (
	"context"
	"time"

	"encore.app/hista/internal/dwd"
	"encore.app/hista/internal/pollen"
	"gorm.io/gorm"
)

type (
	KarlsruheDataGetter interface {
		GetKarlsruheData() (dwd.DWDPollen, time.Time, error)
	}
	Pollens interface {
		ListPollen(context.Context) ([]pollen.PollenEvent, error)
		DoesExistAfter(context.Context, time.Time) error
		CreatePollen(context.Context, *pollen.PollenEvent) error
	}
)

type DWDPollenService struct {
	dwd KarlsruheDataGetter
	p   Pollens
}

func NewDWDPollenService(db *gorm.DB, client KarlsruheDataGetter) *DWDPollenService {
	p := pollen.NewPollenRepo(db)
	return &DWDPollenService{p: p, dwd: client}
}

func (s *DWDPollenService) ListPollen(ctx context.Context) ([]pollen.PollenEvent, error) {
	return s.p.ListPollen(ctx)
}

func (s *DWDPollenService) CreatePollen(ctx context.Context) error {
	dwdPollen, updatedAt, err := s.dwd.GetKarlsruheData()
	if err != nil {
		return err
	}
	if err := s.p.DoesExistAfter(ctx, updatedAt); err != nil {
		return err
	}

	pollens := dwdToPollens(dwdPollen)
	event := pollen.PollenEvent{Pollens: pollens}
	return s.p.CreatePollen(ctx, &event)
}
