package medicine

import (
	"context"
	"time"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

type IntakeRepo struct {
	db *gorm.DB
}

func NewIntakeRepo(db *gorm.DB) *IntakeRepo {
	return &IntakeRepo{db: db}
}

func isIntakesPiid(piid uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("intakes.pi_id = ?", piid)
	}
}

func isNotOld(tx *gorm.DB) *gorm.DB {
	return tx.Where("intakes.date > current_date - 7")
}

func (r IntakeRepo) List(ctx context.Context) ([]*entity.Intake, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return []*entity.Intake{}, err
	}
	return gorm.G[*entity.Intake](r.db).
		Where("pi_id = ?", piid).
		Preload("Medicine", nil).
		Find(ctx)
}

func (r IntakeRepo) ListGrouped(ctx context.Context) (entity.GroupedIntakeList, error) {
	var list entity.GroupedIntakeList
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return list, err
	}

	err = r.db.
		Model(entity.Intake{}).
		Scopes(isIntakesPiid(piid), isNotOld).
		Joins("LEFT JOIN medicines ON intakes.medicine_id = medicines.id and intakes.medicine_pi_id = medicines.pi_id").
		Select("count(*) as count", "medicines.id as medicine_id", "DATE(intakes.date) as date").
		Group("medicines.id").Group("DATE(intakes.date)").
		Scan(&list).Error
	return list, err
}

func (r IntakeRepo) Create(ctx context.Context, medicineId uint) error {
	return generic_queries.Create(ctx, r.db, &entity.Intake{Date: time.Now(), MedicineID: medicineId})
}

func (r IntakeRepo) Remove(ctx context.Context, medicineId uint) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	sub := r.db.
		Model(&entity.Intake{}).
		Scopes(isIntakesPiid(piid)).
		Where("medicine_id = ?", medicineId).
		Where("date > ?", today).
		Limit(1).
		Order("date DESC").
		Select("id")
	tx := r.db.
		Where("id = (?)", sub).
		Delete(&entity.Intake{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
