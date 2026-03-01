package internal

import (
	"context"
	"fmt"
	"slices"

	"encore.app/hista/entity"
)

type (
	IMedicineRepo interface {
		List(ctx context.Context) ([]*entity.Medicine, error)
		Create(ctx context.Context, name string) (uint, error)
		Delete(ctx context.Context, id uint) error
		Patch(ctx context.Context, id uint, column string, value any) error
		BulkShift(ctx context.Context) error
		Reorder(ctx context.Context) error
	}

	MedicineLister interface {
		List(ctx context.Context) (entity.Medicines, error)
	}

	MedicineManager interface {
		Create(ctx context.Context, name string) (uint, error)
		Delete(ctx context.Context, id uint) error
		Rename(ctx context.Context, id uint, name string) error
		Archive(ctx context.Context, id uint, archive bool) error
		Reorder(ctx context.Context, id uint, previousId, nextId *uint) error
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
	r        IMedicineRepo
	stepSize int
	uow      UnitOfWork
}

func NewMedicineMgtmUseCase(r IMedicineRepo, uow UnitOfWork) MedicineMgmtUseCase {
	return MedicineMgmtUseCase{r: r, uow: uow}
}
func (uc MedicineMgmtUseCase) Create(ctx context.Context, name string) (uint, error) {
	var returnId = new(uint)
	err := uc.uow.WithTransaction(ctx, func(tx UnitOfWork) error {
		if err := tx.Medicine().BulkShift(ctx); err != nil {
			fmt.Printf("Error: %s", err)
			return err
		}
		id, err := tx.Medicine().Create(ctx, name)
		returnId = &id
		return err
	})
	return *returnId, err

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

func (uc MedicineMgmtUseCase) Reorder(ctx context.Context, id uint, previousId, nextId *uint) error {
	return uc.uow.WithTransaction(ctx, func(tx UnitOfWork) error {
		medicines, err := tx.Medicine().List(ctx)
		if err != nil {
			return err
		}
		prevOrder, prevOrderFound := findOrder(medicines, previousId)
		nextOrder, nextOrderFound := findOrder(medicines, nextId)
		minOrder, maxOrder := findMinMaxOrder(medicines)

		newOrder, recalculate := calculateNewOrder(prevOrder, prevOrderFound, nextOrder, nextOrderFound, minOrder, maxOrder)
		if recalculate {
			if err := tx.Medicine().Reorder(ctx); err != nil {
				return err
			}
			return uc.Reorder(ctx, id, previousId, nextId)
		}
		return tx.Medicine().Patch(ctx, id, "sort_order", newOrder)
	})
}

func findOrder(meds entity.Medicines, id *uint) (int, bool) {
	if id == nil {
		return 0, false
	}
	for _, m := range meds {
		if m.ID == *id {
			return m.SortOrder, true
		}
	}
	return 0, false
}

func findMinMaxOrder(meds entity.Medicines) (int, int) {
	minOrder := slices.MinFunc(meds, func(a, b *entity.Medicine) int {
		return a.SortOrder - b.SortOrder
	}).SortOrder
	maxOrder := slices.MaxFunc(meds, func(a, b *entity.Medicine) int {
		return a.SortOrder - b.SortOrder
	}).SortOrder
	return minOrder, maxOrder
}

func calculateNewOrder(prevOrder int, prevFound bool, nextOrder int, nextFound bool, minOrder int, maxOrder int) (int, bool) {
	const STEP_SIZE int = 100
	if !prevFound && nextFound {
		if minOrder <= 1 {
			return 0, true
		}
		return minOrder / 2, false
	}
	if !nextFound && prevFound {
		return maxOrder + STEP_SIZE, false
	}
	if prevFound && nextFound {
		if nextOrder-prevOrder <= 1 {
			return 0, true
		}
		return prevOrder + (nextOrder-prevOrder)/2, false
	}
	return STEP_SIZE, false
}
