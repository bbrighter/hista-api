package medicines

import (
	"context"
	"slices"
	"time"

	"gorm.io/gorm"
)

const STEP_SIZE = 100

type MedicinesService struct {
	m        *MedicineRepo
	stepSize int
	db       *gorm.DB
}

func NewMedicineService(db *gorm.DB) *MedicinesService {
	m := NewMedicineRepo(db)
	return &MedicinesService{m: m, stepSize: STEP_SIZE, db: db}
}

func (s *MedicinesService) ListMedicines(ctx context.Context) (Medicines, error) {
	return s.m.ListMedicines(ctx)
}

func (s *MedicinesService) CreateMedicine(ctx context.Context, name string) (uint, error) {
	medicine := &Medicine{Name: name, SortOrder: s.stepSize}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		r := NewMedicineRepo(tx)

		if err := r.BulkShiftMedicine(ctx); err != nil {
			return err
		}
		err := r.CreateMedicine(ctx, medicine)
		return err
	})
	return medicine.ID, err
}

func (s *MedicinesService) DeleteMedicine(ctx context.Context, id uint) error {
	return s.m.DeleteMedicine(ctx, id)
}

func (s *MedicinesService) UpdateMedicine(ctx context.Context, id uint, name *string, archive *bool) error {
	var values = make(map[string]any)
	if name != nil {
		values["name"] = *name
	}
	if archive != nil {
		values["is_archived"] = *archive
	}
	return s.m.UpdateMedicine(ctx, id, values)
}

func (s *MedicinesService) ReorderMedicine(ctx context.Context, id uint, prevId *uint, nextId *uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		r := NewMedicineRepo(tx)
		medicines, err := r.ListMedicines(ctx)
		if err != nil {
			return err
		}
		prevOrder, prevOrderFound := findOrder(medicines, prevId)
		nextOrder, nextOrderFound := findOrder(medicines, nextId)
		minOrder, maxOrder := findMinMaxOrder(medicines)

		newOrder, recalculate := calculateNewOrder(prevOrder, prevOrderFound, nextOrder, nextOrderFound, minOrder, maxOrder)
		if recalculate {
			if err := r.ReorderMedicines(ctx); err != nil {
				return err
			}
			return s.ReorderMedicine(ctx, id, prevId, nextId)
		}
		return r.UpdateMedicine(ctx, id, map[string]any{"sort_order": newOrder})
	})
}

func (s *MedicinesService) RenameMedicine(ctx context.Context, id uint, name string) error {
	return s.m.UpdateMedicine(ctx, id, map[string]any{"name": name})
}

func (s *MedicinesService) ArchiveMedicine(ctx context.Context, id uint, archive bool) error {
	return s.m.UpdateMedicine(ctx, id, map[string]any{"is_archived": archive})
}

func findOrder(meds Medicines, id *uint) (int, bool) {
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

func findMinMaxOrder(meds Medicines) (int, int) {
	minOrder := slices.MinFunc(meds, func(a, b *Medicine) int {
		return a.SortOrder - b.SortOrder
	}).SortOrder
	maxOrder := slices.MaxFunc(meds, func(a, b *Medicine) int {
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

func (s *MedicinesService) ListGroupedIntakes(ctx context.Context) (GroupedIntakeList, error) {
	return s.m.ListGroupedIntakes(ctx)
}

func (s *MedicinesService) IncrementIntake(ctx context.Context, medicineId uint) error {
	intake := Intake{MedicineID: medicineId, Date: time.Now()}
	return s.m.CreateIntake(ctx, &intake)
}

func (s *MedicinesService) DecrementIntake(ctx context.Context, medicineId uint) error {
	return s.m.RemoveLastIntake(ctx, medicineId)
}
