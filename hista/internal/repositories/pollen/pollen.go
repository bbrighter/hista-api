package pollen

import (
	"time"

	"encore.app/hista/entity"
	"encore.dev/rlog"
	"gorm.io/gorm/clause"
)

func (repo *PollenRepo) Create(pollen entity.Pollens, lastUpdated time.Time) error {
	var mustBeUpdated bool = repo.db.
		Where("created_at > ?", lastUpdated).
		First(&entity.PollenEvent{}).
		RowsAffected == 0
	var event = entity.PollenEvent{Pollens: pollen}
	rlog.Info("pollens to be inserted in repo.Create",
		"time", event.CreatedAt,
		"number of pollens", len(event.Pollens),
	)
	for _, p := range event.Pollens {
		rlog.Info("Detail", string(p.Type), p.Intensity)
	}
	if mustBeUpdated {
		err := repo.db.Debug().Create(&event).Error
		rlog.Info("id of event", "id", event.ID)
		return err
	}
	return nil
}

func (repo *PollenRepo) FindPollenWithSeverity(severity int) entity.PollenEvents {
	var pollenEvent entity.PollenEvents
	repo.db.Preload(clause.Associations, "intensity > ?", severity).Find(&pollenEvent)
	return pollenEvent
}
