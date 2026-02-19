package internal

import (
	"context"

	"encore.app/hista/entity"
)

type (
	IMedicineRepo interface {
		List(ctx context.Context) ([]*entity.Medicine, error)
		Create(ctx context.Context, name string) (uint, error)
		Delete(ctx context.Context, id uint) error
		Patch(ctx context.Context, id uint, column string, value any) error
	}

	MedicineLister interface {
		List(ctx context.Context) (entity.Medicines, error)
	}

	MedicineManager interface {
		Create(ctx context.Context, name string) (uint, error)
		Delete(ctx context.Context, id uint) error
		Rename(ctx context.Context, id uint, name string) error
		Archive(ctx context.Context, id uint, archive bool) error
	}
)

type MedicineListUseCase struct {
	r IMedicineRepo
}

func NewMedicineListUseCase(r IMedicineRepo) MedicineListUseCase {
	return MedicineListUseCase{r: r}
}

func (uc MedicineListUseCase) List(ctx context.Context) (entity.Medicines, error) {
	return uc.r.List(ctx)
}

type MedicineMgmtUseCase struct {
	r IMedicineRepo
}

func NewMedicineMgtmUseCase(r IMedicineRepo) MedicineMgmtUseCase {
	return MedicineMgmtUseCase{r: r}
}
func (uc MedicineMgmtUseCase) Create(ctx context.Context, name string) (uint, error) {
	return uc.r.Create(ctx, name)
}

func (uc MedicineMgmtUseCase) Delete(ctx context.Context, id uint) error {
	return uc.r.Delete(ctx, id)
}
func (uc MedicineMgmtUseCase) Rename(ctx context.Context, id uint, name string) error {
	return uc.r.Patch(ctx, id, "name", name)
}
func (uc MedicineMgmtUseCase) Archive(ctx context.Context, id uint, archive bool) error {
	return uc.r.Patch(ctx, id, "is_archived", archive)
}
