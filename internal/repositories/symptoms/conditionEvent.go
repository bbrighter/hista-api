package symptoms

import (
	"time"

	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm/clause"
)

func (repo *SymptomsRepo) CreateConditionEvent(event *entity.ConditionEvent) error {
	return repo.db.Create(&event).Error
}

func (repo *SymptomsRepo) ListConditionEvents() entity.ConditionEvents {
	var events entity.ConditionEvents
	repo.db.Find(&events)
	return events
}

func (repo *SymptomsRepo) GetConditionEvent(id uint) (entity.ConditionEvent, error) {
	var event = entity.ConditionEvent{ID: id}
	if repo.db.Preload("Conditions.Symptom").Preload(clause.Associations).Find(&event).RowsAffected == 0 {
		return event, errors.ErrorNotFound
	}
	return event, nil
}

func (repo *SymptomsRepo) DeleteConditionEvent(id uint) error {
	var event = entity.ConditionEvent{ID: id}
	tx := repo.db.Delete(&event)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return nil
}

func (repo *SymptomsRepo) PatchConditionEvent(eventId uint, date time.Time) error {
	var event = entity.ConditionEvent{ID: eventId}
	if repo.db.Find(&event).RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return repo.db.Model(event).Update("date", date).Error
}

func (repo *SymptomsRepo) ListConditionEventsAndDependencies() entity.ConditionEvents {
	var events entity.ConditionEvents
	repo.db.Preload("Conditions.Symptom").
		Preload("Conditions").
		Find(&events)
	return events
}
