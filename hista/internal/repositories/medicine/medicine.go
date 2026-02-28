package medicine

import (
	"context"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
)

type MedicineRepo struct {
	db       *gorm.DB
	stepSize int
}

func NewMedicineRepo(db *gorm.DB) *MedicineRepo {
	return &MedicineRepo{db: db, stepSize: 100}
}

func (r MedicineRepo) List(ctx context.Context) ([]*entity.Medicine, error) {
	return generic_queries.List[*entity.Medicine](ctx, r.db)
}

func (r MedicineRepo) Create(ctx context.Context, name string) (uint, error) {
	medicine := entity.Medicine{Name: name, SortOrder: r.stepSize}
	err := generic_queries.Create(ctx, r.db, &medicine)
	return medicine.ID, err
}

func (r MedicineRepo) Delete(ctx context.Context, id uint) error {
	return generic_queries.Delete[*entity.Medicine](ctx, r.db, id)
}

func (r MedicineRepo) Patch(ctx context.Context, id uint, column string, value any) error {
	return generic_queries.UpdateColumn[*entity.Medicine](ctx, r.db, id, column, value)
}

func (r MedicineRepo) BulkShift(ctx context.Context) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	rows, err := gorm.G[entity.Medicine](r.db).
		Where("pi_id = ?", piid).
		Update(ctx, "sort_order", gorm.Expr("sort_order + ?", r.stepSize))
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r MedicineRepo) Reorder(ctx context.Context) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	return r.db.Exec(`
	UPDATE medicines m
	SET sort_order = o.new_order * ?
	FROM (
		SELECT 
			id,
			pi_id,
			ROW_NUMBER() OVER (ORDER BY sort_order) as new_order
		FROM medicines
		WHERE pi_id = ?
	) as o
	 WHERE o.id = m.id AND o.pi_id = m.pi_id`, r.stepSize, piid).Error
}
