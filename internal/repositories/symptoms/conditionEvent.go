package symptoms

import (
	"context"
	"time"

	"encore.app/entity"
	"encore.app/generic_queries"
	"gorm.io/gorm"
)

func (repo *SymptomsRepo) CreateConditionEvent(ctx context.Context, event *entity.ConditionEvent) error {
	return generic_queries.Create(ctx, repo.db, event)
}

func (repo *SymptomsRepo) ListConditionEvents(ctx context.Context) ([]*entity.ConditionEvent, error) {
	return generic_queries.List[*entity.ConditionEvent](ctx, repo.db)
}

func (repo *SymptomsRepo) GetConditionEvent(ctx context.Context, id uint) (entity.ConditionEvent, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.ConditionEvent{}, err
	}
	return gorm.G[entity.ConditionEvent](repo.db).
		Where("pi_id = ?", piid).
		Where("id = ?", id).
		Preload("Conditions", nil).
		Preload("Conditions.Symptom", nil).
		First(ctx)
}

func (repo *SymptomsRepo) DeleteConditionEvent(ctx context.Context, id uint) error {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return err
	}
	event, err := gorm.G[entity.ConditionEvent](repo.db).Preload("Conditions", nil).Where("pi_id = ?", piid).Where("id = ?", id).First(ctx)
	if err != nil {
		return err
	}

	return repo.db.Transaction(func(tx *gorm.DB) error {
		for _, c := range event.Conditions {
			return deleteConditionAndSymptoms(ctx, tx, c.ID, piid)
		}
		return generic_queries.Delete[*entity.ConditionEvent](ctx, tx, id)
	})

}

func (repo *SymptomsRepo) PatchConditionEvent(ctx context.Context, eventId uint, date time.Time) error {
	return generic_queries.UpdateColumn[*entity.ConditionEvent](ctx, repo.db, eventId, "date", date)
}

func (repo *SymptomsRepo) ListConditionEventsAndDependencies(ctx context.Context) ([]*entity.ConditionEvent, error) {
	piid, err := generic_queries.PiidFromCtx(ctx)
	if err != nil {
		return entity.ConditionEvents{}, err
	}
	return gorm.G[*entity.ConditionEvent](repo.db).
		Where("pi_id = ?", piid).
		Preload("Conditions.Symptom", nil).
		Preload("Conditions", nil).
		Find(ctx)

}
