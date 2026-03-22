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
	return repo.db.Debug().Transaction(func(tx *gorm.DB) error {
		var event = entity.PollenEvent{}
		if err := gorm.G[entity.PollenEvent](tx).Create(ctx, &event); err != nil {
			return err
		}
		var dbPollens = []entity.Pollen{}
		for _, p := range pollens {
			dbPollens = append(dbPollens, entity.Pollen{
				PollenEventID: event.ID,
				Type:          p.Type,
				Intensity:     p.Intensity,
			})
		}
		return gorm.G[entity.Pollen](tx).CreateInBatches(ctx, &dbPollens, 100)
	})
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
