package internal

import (
	"context"
	"testing"
	"time"

	"encore.app/hista/entity"
	"encore.dev/rlog"
)

type (
	IPollenRepo interface {
		FindPollenWithSeverity(severity int) entity.PollenEvents
		DoesExistAfter(ctx context.Context, time time.Time) error
		Create(ctx context.Context, pollens entity.Pollens) error
	}

	IDWDRepo interface {
		GetKarlsruheData() (entity.DWDPollen, error)
		DwdStringToDate() (time.Time, error)
		UseTestQuery(*testing.T)
	}

	IPollenUseCase interface {
		List(ctx context.Context) (entity.PollenEvents, error)
		Create(ctx context.Context) error
		UseTestQuery(*testing.T)
	}
)

type PollenUseCase struct {
	repo IPollenRepo
	dwd  IDWDRepo
}

func NewPollenUseCase(repo IPollenRepo, dwd IDWDRepo) PollenUseCase {
	return PollenUseCase{repo: repo, dwd: dwd}
}

func (uc PollenUseCase) List(ctx context.Context) (entity.PollenEvents, error) {
	return uc.repo.FindPollenWithSeverity(0), nil
}

func (uc PollenUseCase) Create(ctx context.Context) error {
	pollen, err := uc.dwd.GetKarlsruheData()
	if err != nil {
		return err
	}
	updatedAt, err := uc.dwd.DwdStringToDate()
	if err != nil {
		return err
	}
	if err := uc.repo.DoesExistAfter(ctx, updatedAt); err != nil {
		return err
	}

	dwdPollen := pollen.ToPollen()

	args := []any{
		"updatedAt", updatedAt,
	}

	for _, p := range dwdPollen {
		args = append(args, string(p.Type), p.Intensity.String())
	}

	rlog.Info("pollen to be saved", args...)

	return uc.repo.Create(ctx, dwdPollen)
}

func (uc PollenUseCase) UseTestQuery(t *testing.T) {
	uc.dwd.UseTestQuery(t)
}
