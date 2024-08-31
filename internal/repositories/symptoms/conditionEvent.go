package symptoms

import (
	"time"

	"encore.app/entity"
	"encore.app/errors"
	"gorm.io/gorm/clause"
)

func (repo *SymptomsRepo) CreateConditionEvent(date time.Time) (uint, error) {
	var event = entity.ConditionEvent{Date: date}
	err := repo.db.Create(&event).Error
	return event.ID, err
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
	if tx.RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return tx.Error
}

func (repo *SymptomsRepo) PatchConditionEvent(eventId uint, date time.Time) error {
	var event = entity.ConditionEvent{ID: eventId}
	if repo.db.Find(&event).RowsAffected == 0 {
		return errors.ErrorNotFound
	}
	return repo.db.Model(event).Update("date", date).Error
}

func (repo *SymptomsRepo) GetConditionEventsAndDependencies() (entity.ConditionEvents, entity.SymptomCategories, error) {
	var events entity.ConditionEvents
	var cats entity.SymptomCategories
	if err := repo.db.Preload("Conditions.Symptom").
		Preload("Conditions").
		Find(&events).
		Error; err != nil {
		return events, cats, err
	}
	if err := repo.db.Find(&cats).Error; err != nil {
		return events, cats, err
	}
	return events, cats, nil
}
