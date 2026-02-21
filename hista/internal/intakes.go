package internal

import (
	"context"

	"encore.app/hista/entity"
)

type (
	IIntakeRepo interface {
		List(ctx context.Context) (entity.GroupedIntakeList, error)
		Create(ctx context.Context, medicineId uint) error
		Remove(ctx context.Context, medicineId uint) error
	}

	IntakeManager interface {
		List(ctx context.Context) (entity.GroupedIntakeList, error)
		Increment(ctx context.Context, medicineId uint) error
		Decrement(ctx context.Context, medicineId uint) error
	}
)

type IntakeMgmtUseCase struct {
	r IIntakeRepo
}

func NewIntakeMgmtUseCase(r IIntakeRepo) IntakeMgmtUseCase {
	return IntakeMgmtUseCase{r: r}
}

func (uc IntakeMgmtUseCase) List(ctx context.Context) (entity.GroupedIntakeList, error) {
	return uc.r.List(ctx)
}
func (uc IntakeMgmtUseCase) Increment(ctx context.Context, medicineId uint) error {
	return uc.r.Create(ctx, medicineId)
}
func (uc IntakeMgmtUseCase) Decrement(ctx context.Context, medicineId uint) error {
	return uc.r.Remove(ctx, medicineId)
}
