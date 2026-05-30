package medicines

import (
	"context"
	"time"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MedicineRepo struct {
	db       *gorm.DB
	stepSize int
}

func NewMedicineRepo(db *gorm.DB) *MedicineRepo {
	return &MedicineRepo{db: db, stepSize: 100}
}

func (r *MedicineRepo) ListMedicines(ctx context.Context) ([]*Medicine, error) {
	return generic_queries.List[*Medicine](ctx, r.db)
}

func (r *MedicineRepo) CreateMedicine(ctx context.Context, medicine *Medicine) error {
	return generic_queries.Create(ctx, r.db, medicine)
}

func (r *MedicineRepo) DeleteMedicine(ctx context.Context, id uint) error {
	return generic_queries.Delete[*Medicine](ctx, r.db, id)
}

func (r *MedicineRepo) UpdateMedicine(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "medicines", id, values)
}

func (r *MedicineRepo) BulkShiftMedicine(ctx context.Context) error {
	_, err := gorm.G[Medicine](r.db).
		Scopes(wherePiid(ctx)).
		Update(ctx, "sort_order", gorm.Expr("sort_order + ?", r.stepSize))
	return err
}

func (r *MedicineRepo) ReorderMedicines(ctx context.Context) error {
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

func (r *MedicineRepo) ListIntakes(ctx context.Context) ([]Intake, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return []Intake{}, err
	}
	return gorm.G[Intake](r.db).
		Joins(clause.JoinTarget{Association: "Medicine"}, func(db gorm.JoinBuilder, joinTable, curTable clause.Table) error {
			db.Where("pi_id = ?", piid)
			return nil
		}).
		Find(ctx)
}

func (r *MedicineRepo) ListGroupedIntakes(ctx context.Context) (GroupedIntakeList, error) {
	var list GroupedIntakeList
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return list, err
	}

	err = r.db.
		Model(Intake{}).
		Scopes(isNotOld).
		Joins("LEFT JOIN medicines ON intakes.medicine_id = medicines.id").
		Select("count(*) as count", "medicines.id as medicine_id", "DATE(intakes.date) as date").
		Where("medicines.pi_id = ?", piid).
		Group("medicines.id").Group("DATE(intakes.date)").
		Scan(&list).Error
	return list, err
}

func (r *MedicineRepo) CreateIntake(ctx context.Context, intake *Intake) error {
	return gorm.G[Intake](r.db).Create(ctx, intake)
}

func (r *MedicineRepo) RemoveLastIntake(ctx context.Context, medicineId uint) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	sub := r.db.
		Model(&Intake{}).
		Where("medicine_id = ?", medicineId).
		Where("date > ?", today).
		Limit(1).
		Order("date DESC").
		Select("id")
	tx := r.db.
		Where("id = (?)", sub).
		Delete(&Intake{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
