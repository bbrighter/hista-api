package internal

import (
	"time"

	"encore.app/entity"
)

type HeadacheUseCase struct {
	hRepo IHeadacheRepo
}

func NewHeadacheUseCase(hRepo IHeadacheRepo) HeadacheUseCase {
	return HeadacheUseCase{hRepo: hRepo}
}

func (uc HeadacheUseCase) List() entity.Headaches {
	return uc.hRepo.ListHeadaches()
}
func (uc HeadacheUseCase) Create(date time.Time, severity entity.HeadacheSeverity) (uint, error) {
	var headache = entity.Headache{
		Date:     date,
		Severity: severity,
	}
	err := uc.hRepo.CreateHeadache(&headache)
	return headache.ID, err
}
func (uc HeadacheUseCase) Delete(haId uint) error {
	return uc.hRepo.DeleteHeadache(haId)
}
func (uc HeadacheUseCase) Get(haId uint) (entity.Headache, error) {
	return uc.hRepo.GetHeadache(haId)
}
func (uc HeadacheUseCase) PatchDate(haId uint, date time.Time) error {
	return uc.hRepo.PatchHeadache(haId, &date, nil, nil, nil, nil)
}
func (uc HeadacheUseCase) PatchSeverity(haId uint, severity entity.HeadacheSeverity) error {
	return uc.hRepo.PatchHeadache(haId, nil, &severity, nil, nil, nil)
}
func (uc HeadacheUseCase) PatchTypes(haId uint, types entity.HeadacheTypes) error {
	return uc.hRepo.PatchHeadache(haId, nil, nil, &types, nil, nil)
}
func (uc HeadacheUseCase) PatchPositions(haId uint, positions entity.HeadachePositions) error {
	return uc.hRepo.PatchHeadache(haId, nil, nil, nil, &positions, nil)
}
func (uc HeadacheUseCase) PatchSymptoms(haId uint, symptoms entity.HeadacheSymptoms) error {
	return uc.hRepo.PatchHeadache(haId, nil, nil, nil, nil, &symptoms)
}
