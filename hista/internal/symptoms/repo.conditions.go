package symptoms

import (
	"context"

	"encore.app/shared/generic_queries"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConditionRepo struct {
	db *gorm.DB
}

func NewConditionRepo(db *gorm.DB) *ConditionRepo {
	return &ConditionRepo{db: db}
}

func (r *ConditionRepo) ListConditionEvents(ctx context.Context) (ConditionEvents, error) {
	return generic_queries.List[*ConditionEvent](ctx, r.db)
}

func (r *ConditionRepo) ListConditionEventsAndDependencies(ctx context.Context) (ConditionEvents, error) {
	return gorm.G[*ConditionEvent](r.db).Scopes(wherePiid(ctx)).Preload("Conditions", nil).Preload("Conditions.Symptom", nil).Find(ctx)
}

func (r *ConditionRepo) FirstConditionEventAndConditions(ctx context.Context, id uint) (ConditionEvent, error) {
	return gorm.G[ConditionEvent](r.db).
		Scopes(wherePiid(ctx)).
		Where("id = ?", id).
		Preload("Conditions", nil).First(ctx)
}

func (r *ConditionRepo) UpdateConditionEvent(ctx context.Context, id uint, values map[string]any) error {
	return generic_queries.Updates(ctx, r.db, "condition_events", id, values)
}

func (r *ConditionRepo) DeleteConditionEvent(ctx context.Context, id uint) error {
	return generic_queries.Delete[*ConditionEvent](ctx, r.db, id)
}

func (r *ConditionRepo) CreateConditionEvent(ctx context.Context, event *ConditionEvent) error {
	return generic_queries.Create(ctx, r.db, event)
}

func (r *ConditionRepo) ListConditions(ctx context.Context, eventId uint) ([]Condition, error) {
	return gorm.G[Condition](r.db).
		Where("condition_event_id = ?", eventId).
		Find(ctx)
}

func (r *ConditionRepo) CreateCondition(ctx context.Context, condition *Condition) error {
	tx := r.db.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}, {Name: "symptom_id"}}}).Create(condition)
	return tx.Error
}

func (r *ConditionRepo) DeleteCondition(ctx context.Context, id uint) error {
	rows, err := gorm.G[Condition](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ConditionRepo) UpdateCondition(ctx context.Context, id uint, values map[string]any) error {
	rows, err := gorm.G[map[string]any](r.db).Table("conditions").Where("id = ?", id).Updates(ctx, values)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
