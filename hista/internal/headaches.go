package internal

import (
	"context"
	"time"

	"encore.app/hista/entity"
)

type (
	IHeadacheRepo interface {
		ListHeadaches(ctx context.Context) ([]*entity.Headache, error)
		CreateHeadache(ctx context.Context, ha *entity.Headache) error
		DeleteHeadache(ctx context.Context, haId uint) error
		GetHeadache(ctx context.Context, haId uint) (*entity.Headache, error)
		PatchHeadache(ctx context.Context, haId uint, date *time.Time, severity *entity.HeadacheSeverity, types *entity.HeadacheTypes, positions *entity.HeadachePositions, symptoms *entity.HeadacheSymptoms, description *string) error
	}

	IHeadacheUseCase interface {
		List(ctx context.Context) (entity.Headaches, error)
		Create(ctx context.Context, date time.Time, severity entity.HeadacheSeverity) (uint, error)
		Delete(ctx context.Context, haId uint) error
		Get(ctx context.Context, haId uint) (*entity.Headache, error)
		PatchDate(ctx context.Context, haId uint, date time.Time) error
		PatchSeverity(ctx context.Context, haId uint, severity entity.HeadacheSeverity) error
		PatchTypes(ctx context.Context, haId uint, types entity.HeadacheTypes) error
		PatchPositions(ctx context.Context, haId uint, positions entity.HeadachePositions) error
		PatchSymptoms(ctx context.Context, haId uint, symptoms entity.HeadacheSymptoms) error
		PatchDescription(ctx context.Context, haId uint, description string) error
	}
)

type HeadacheUseCase struct {
	hRepo IHeadacheRepo
}

func NewHeadacheUseCase(hRepo IHeadacheRepo) HeadacheUseCase {
	return HeadacheUseCase{hRepo: hRepo}
}

func (uc HeadacheUseCase) List(ctx context.Context) (entity.Headaches, error) {
	return uc.hRepo.ListHeadaches(ctx)
}
func (uc HeadacheUseCase) Create(ctx context.Context, date time.Time, severity entity.HeadacheSeverity) (uint, error) {
	var headache = entity.Headache{
		Date:     date,
		Severity: severity,
	}
	err := uc.hRepo.CreateHeadache(ctx, &headache)
	return headache.ID, errorMapper(err)
}
func (uc HeadacheUseCase) Delete(ctx context.Context, haId uint) error {
	return errorMapper(uc.hRepo.DeleteHeadache(ctx, haId))
}
func (uc HeadacheUseCase) Get(ctx context.Context, haId uint) (*entity.Headache, error) {
	h, err := uc.hRepo.GetHeadache(ctx, haId)
	return h, errorMapper(err)
}
func (uc HeadacheUseCase) PatchDate(ctx context.Context, haId uint, date time.Time) error {
	return errorMapper(uc.hRepo.PatchHeadache(ctx, haId, &date, nil, nil, nil, nil, nil))
}
func (uc HeadacheUseCase) PatchSeverity(ctx context.Context, haId uint, severity entity.HeadacheSeverity) error {
	return errorMapper(uc.hRepo.PatchHeadache(ctx, haId, nil, &severity, nil, nil, nil, nil))
}
func (uc HeadacheUseCase) PatchTypes(ctx context.Context, haId uint, types entity.HeadacheTypes) error {
	return errorMapper(uc.hRepo.PatchHeadache(ctx, haId, nil, nil, &types, nil, nil, nil))
}
func (uc HeadacheUseCase) PatchPositions(ctx context.Context, haId uint, positions entity.HeadachePositions) error {
	return errorMapper(uc.hRepo.PatchHeadache(ctx, haId, nil, nil, nil, &positions, nil, nil))
}
func (uc HeadacheUseCase) PatchSymptoms(ctx context.Context, haId uint, symptoms entity.HeadacheSymptoms) error {
	return errorMapper(uc.hRepo.PatchHeadache(ctx, haId, nil, nil, nil, nil, &symptoms, nil))
}
func (uc HeadacheUseCase) PatchDescription(ctx context.Context, haId uint, description string) error {
	return errorMapper(uc.hRepo.PatchHeadache(ctx, haId, nil, nil, nil, nil, nil, &description))
}
