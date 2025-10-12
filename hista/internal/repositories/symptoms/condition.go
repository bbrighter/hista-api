package symptoms

import (
	"context"

	"encore.app/hista/entity"
	"encore.app/shared/generic_queries"
	"encore.dev/types/uuid"
	"gorm.io/gorm"
)

// Get all conditions including their conditionTypes
func (repo *SymptomsRepo) ListConditions(ctx context.Context, eventId uint) ([]*entity.Condition, error) {
	var conditions []*entity.Condition
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return conditions, err
	}
	return gorm.G[*entity.Condition](repo.db).
		Where("pi_id = ?", piid).
		Where(entity.Condition{ConditionEventID: eventId}).
		// Preload(clause.Associations, nil).
		Preload("Symptom", nil).
		Find(ctx)
}

// Create a new condition based on the name. New symptoms are only created if the name in the corresponding category doesn't exist.
// Requires a conditionEventID
func (repo *SymptomsRepo) CreateConditionBySymptomName(ctx context.Context, eventId uint, symptomName string, categoryId uint) (uint, error) {
	count, err := generic_queries.Count[*entity.ConditionEvent](ctx, repo.db, eventId)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	symptomId, err := repo.CreateOrReplace(ctx, symptomName, categoryId)
	if err != nil {
		return 0, err
	}
	condition := &entity.Condition{SymptomID: symptomId, ConditionEventID: eventId, Severity: entity.MediumSeverity}
	err = generic_queries.Create(ctx, repo.db, condition)
	return condition.ID, err
}

// Create a new condition by SymptomID.
// Requires a ConditionEventID and SymptomID
func (repo *SymptomsRepo) CreateConditionBySymptomID(ctx context.Context, eventId uint, symptomId uint) (uint, error) {
	count, err := generic_queries.Count[*entity.ConditionEvent](ctx, repo.db, eventId)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	symptom, err := generic_queries.First[*entity.Symptom](ctx, repo.db, symptomId)
	if err != nil {
		return 0, err
	}
	condition := &entity.Condition{Symptom: *symptom, ConditionEventID: eventId}
	err = generic_queries.Create(ctx, repo.db, condition)
	return condition.ID, err
}

func (repo *SymptomsRepo) DeleteCondition(ctx context.Context, conditionId uint) error {
	// Delete condition and drop all symptoms if it was the last condition using this symptom.
	// If it was the last symptom in a category, drop it as well.

	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	return deleteConditionAndSymptoms(ctx, repo.db, conditionId, piid)

	// sym, err := gorm.G[entity.Condition](repo.db).
	// 	Preload("Symptom", nil).
	// 	Where("pi_id = ?", piid).
	// 	Where("id = ?", conditionId).
	// 	First(ctx)
	// if err != nil {
	// 	return err
	// }
	// countUsageOfCondition, err := gorm.G[entity.Symptom](repo.db).Where("pi_id = ?", piid).Where("id = ?", sym.SymptomID).Count(ctx, "*")
	// if err != nil {
	// 	return err
	// }
	// return repo.db.Transaction(func(tx *gorm.DB) error {
	// 	if countUsageOfCondition == 1 {
	// 		if err := generic_queries.Delete[*entity.Symptom](ctx, tx, sym.SymptomID); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return generic_queries.Delete[*entity.Condition](ctx, tx, conditionId)
	// })

}

func (repo *SymptomsRepo) ChangeSeverity(ctx context.Context, conditionId uint, newSeverity entity.Severity) error {
	return generic_queries.UpdateColumn[*entity.Condition](ctx, repo.db, conditionId, "Severity", newSeverity)
}

func (repo *SymptomsRepo) GetCondition(ctx context.Context, id uint) (entity.Condition, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.Condition{}, err
	}
	cond, err := gorm.G[entity.Condition](repo.db).
		Preload("Symptom", nil).
		Where("pi_id = ?", piid).
		Where("id = ?", id).
		First(ctx)
	return cond, err
}

func deleteConditionAndSymptoms(ctx context.Context, db *gorm.DB, conditionId uint, piid uuid.UUID) error {
	con, err := gorm.G[entity.Condition](db).
		Preload("Symptom", nil).
		Where("pi_id = ?", piid).
		Where("id = ?", conditionId).
		First(ctx)
	if err != nil {
		return err
	}
	countUsageOfCondition, err := gorm.G[entity.Symptom](db).Where("pi_id = ?", piid).Where("id = ?", con.SymptomID).Count(ctx, "*")
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := generic_queries.Delete[*entity.Condition](ctx, tx, conditionId); err != nil {
			return err
		}
		if countUsageOfCondition == 1 {
			return generic_queries.Delete[*entity.Symptom](ctx, tx, con.SymptomID)
		}
		return nil
	})
}
