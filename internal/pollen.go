package internal

import (
	"testing"

	"encore.app/entity"
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
		return err
	}
	updatedAt, err := uc.dwd.DwdStringToDate()
	if err != nil {
		return err
	}
	return uc.repo.Create(pollen.ToPollen(), updatedAt)
}

func (uc PollenUseCase) UseTestQuery(t *testing.T) {
	uc.dwd.UseTestQuery(t)
}
