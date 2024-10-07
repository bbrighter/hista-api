package status

import (
	"encore.app/entity"
)

func (repo *StatusRepo) Create(event *entity.Status) error {
	err := repo.db.Create(&event).Error
	return err
}

func (repo *StatusRepo) Find() (events entity.Statuses) {
	repo.db.Find(&events)
	return events
}
