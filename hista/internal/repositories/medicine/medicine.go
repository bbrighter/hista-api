package medicine

import (
	"context"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type MedicineRepo struct {
	db *gorm.DB
}

func NewMedicineRepo(db *gorm.DB) *MedicineRepo {
	return &MedicineRepo{db: db}
}

func (r MedicineRepo) List(ctx context.Context) ([]*entity.Medicine, error) {
	return generic_queries.List[*entity.Medicine](ctx, r.db)
}

func (r MedicineRepo) Create(ctx context.Context, name string) (uint, error) {
	medicine := entity.Medicine{Name: name}
	err := generic_queries.Create(ctx, r.db, &medicine)
	return medicine.ID, err
}

func (r MedicineRepo) Delete(ctx context.Context, id uint) error {
	return generic_queries.Delete[*entity.Medicine](ctx, r.db, id)
}

func (r MedicineRepo) Patch(ctx context.Context, id uint, column string, value any) error {
	return generic_queries.UpdateColumn[*entity.Medicine](ctx, r.db, id, column, value)
}
