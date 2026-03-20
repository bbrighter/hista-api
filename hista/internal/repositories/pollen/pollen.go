package pollen

import (
	"context"
	"time"

	"encore.app/errors"
	"encore.app/hista/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *PollenRepo) Create(ctx context.Context, pollens entity.Pollens) error {
	var event = entity.PollenEvent{Pollens: pollens}
	return gorm.G[entity.PollenEvent](repo.db).Create(ctx, &event)
}

func (r *PollenRepo) DoesExistAfter(ctx context.Context, time time.Time) error {
	count, err := gorm.G[entity.PollenEvent](r.db).Where("created_at > ?", time).Count(ctx, "*")
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.ErrorAlreadyExists
	}
	return nil
}

func (repo *PollenRepo) FindPollenWithSeverity(severity int) entity.PollenEvents {
	var pollenEvent entity.PollenEvents
	repo.db.Preload(clause.Associations, "intensity > ?", severity).Find(&pollenEvent)
	return pollenEvent
}
