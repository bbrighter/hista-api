package pollen

import (
	"time"

	"encore.app/entity"
	"gorm.io/gorm/clause"
)

func (repo *PollenRepo) Create(pollen entity.Pollens, lastUpdated time.Time) error {
	var mustBeUpdated bool = repo.db.
		Where("created_at > ?", lastUpdated).
		First(&entity.PollenEvent{}).
		RowsAffected == 0
	if mustBeUpdated {
		var event = entity.PollenEvent{Pollens: pollen}
		return repo.db.Create(&event).Error
	}
	return nil
}

func (repo *PollenRepo) FindPollenWithSeverity(severity int) entity.PollenEvents {
	var pollenEvent entity.PollenEvents
	repo.db.Preload(clause.Associations, "intensity > ?", severity).Find(&pollenEvent)
	return pollenEvent
}
