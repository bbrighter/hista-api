package internal

import (
	"testing"
	"time"

	"encore.app/hista/entity"
)

type (
	IPollenRepo interface {
		FindPollenWithSeverity(severity int) entity.PollenEvents
		Create(pollen entity.Pollens, lastUpdated time.Time) error
	}

	IDWDRepo interface {
		GetKarlsruheData() (entity.DWDPollen, error)
		DwdStringToDate() (time.Time, error)
		UseTestQuery(*testing.T)
	}

	IPollenUseCase interface {
		List() entity.PollenEvents
		Create() error
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

func (uc PollenUseCase) List() entity.PollenEvents {
	return uc.repo.FindPollenWithSeverity(0)
}

func (uc PollenUseCase) Create() error {
	pollen, err := uc.dwd.GetKarlsruheData()
	if err != nil {
		return errorMapper(err)
	}
	updatedAt, err := uc.dwd.DwdStringToDate()
	if err != nil {
		return errorMapper(err)
	}
	return errorMapper(uc.repo.Create(pollen.ToPollen(), updatedAt))
}

func (uc PollenUseCase) UseTestQuery(t *testing.T) {
	uc.dwd.UseTestQuery(t)
}
